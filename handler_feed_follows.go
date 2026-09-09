package main

import (
	"context"
	"fmt"
	"gator/internal/database"
)

func feedFollowsHandler(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("missing feed URL")
	}

	feedURL := cmd.Args[0]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

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

func followingHandler(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

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
