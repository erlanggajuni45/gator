package command

import (
	"fmt"
	"gator/internal/config"
)

type state struct {
	Config *config.Config
}

type command struct {
	Name string
	Args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("username is required")
	}

	username := cmd.Args[0]
	err := s.Config.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}

	fmt.Printf("User '%s' has been set.\n", username)
	return nil
}
