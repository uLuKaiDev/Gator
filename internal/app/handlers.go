package app

import (
	"errors"
	"fmt"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("missing username")
	}
	username := cmd.Args[0]

	err := s.Config.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	fmt.Printf("User set to %s\n", username)
	return nil
}
