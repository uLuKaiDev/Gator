package app

import (
	"context"
	"fmt"

	"github.com/uLuKaiDev/Gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(s *State, cmd Command) error {
	return func(s *State, cmd Command) error {
		if s.Config.CurrentUserName == "" {
			return fmt.Errorf("user not set, please set a user first")
		}

		user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		return handler(s, cmd, user)
	}
}
