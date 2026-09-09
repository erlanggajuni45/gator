package main

import (
	"context"
	"fmt"
	"gator/internal/database"
)

func feedFollowsHandler(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("missing feed URL")
	}

	feedURL := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)

	newFeedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}

	fmt.Printf("Current User: %s\n", newFeedFollow.UserName)
	fmt.Printf("Feed Name: %s\n", newFeedFollow.FeedName)
	return nil
}

func followingHandler(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feed follows: %w", err)
	}

	fmt.Printf("Feed that you are following:\n")
	for _, follow := range follows {
		fmt.Printf("- %s\n", follow.FeedName)
	}
	return nil
}

func deleteFeedFollowHandler(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("missing feed URL")
	}

	feedURL := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("failed to get feed: %w", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete feed follow: %w", err)
	}

	return nil
}
