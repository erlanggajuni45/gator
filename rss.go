package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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
	// if len(cmd.Args) < 1 {
	// 	return fmt.Errorf("usage: agg <feed_url>")
	// }

	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	fmt.Printf("Feed Title: %s\n", feed.Channel.Title)
	fmt.Printf("Feed Link: %s\n", feed.Channel.Link)
	fmt.Printf("Feed Description: %s\n", feed.Channel.Description)
	fmt.Println("Items:")
	for _, item := range feed.Channel.Items {
		fmt.Printf("- Title: %s\n", item.Title)
		fmt.Printf("  Link: %s\n", item.Link)
		fmt.Printf("  Description: %s\n", item.Description)
		fmt.Printf("  PubDate: %s\n", item.PubDate)
	}

	return nil
}
