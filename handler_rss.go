package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/erlanggajuni45/gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	// set User-Agent header to gator
	req.Header.Set("User-Agent", "gator")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var feed RSSFeed
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Items {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Items[i] = item
	}

	return &feed, nil
}

func handleFetchFeed(state *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid time_between_reqs: %v", err)
	}

	log.Printf("Collecting feed every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(state)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("error fetching feed: %v", err)
		return
	}

	log.Println("found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("error marking feed fetched: %v", err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("error fetching feed: %v", err)
		return
	}

	for _, item := range feedData.Channel.Items {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
}

func handleAddFeed(state *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: add <feed_name> <feed_url>")
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	newFeed, err := state.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   feedName,
		Url:    feedURL,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed: %v", err)
	}

	fmt.Printf("Feed added: %s\n", newFeed.Name)
	fmt.Printf("Feed URL: %s\n", newFeed.Url)
	fmt.Printf("Feed ID: %s\n", newFeed.ID)
	fmt.Printf("Feed CreatedAt: %s\n", newFeed.CreatedAt)
	fmt.Printf("Feed UpdatedAt: %s\n", newFeed.UpdatedAt)
	fmt.Printf("Feed UserID: %s\n", newFeed.UserID)
	fmt.Printf("* LastFetchedAt: %v\n", newFeed.LastFetchedAt.Time)

	_, err = state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		FeedID: newFeed.ID,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error following feed: %v", err)
	}
	return nil
}

func handleListFeeds(state *state, cmd command) error {
	feeds, err := state.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %v", err)
	}

	for _, feed := range feeds {
		fmt.Printf("- Name: %s\n", feed.Name)
		fmt.Printf("  URL: %s\n", feed.Url)
		fmt.Printf("  User: %s\n", feed.User)
	}

	return nil
}
