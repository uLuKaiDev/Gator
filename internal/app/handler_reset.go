package app

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pressly/goose/v3"
)

func HandlerDBReset(s *State, cmd Command) error {
	if len(cmd.Args) == 0 || cmd.Args[0] != "--force" {
		fmt.Println("This will DROP and RECREATE the database schema.")
		fmt.Println("If you're sure, run: gator db-reset --force")
		return nil
	}

	migrationsDir, err := filepath.Abs("sql/schema")
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	goose.SetDialect("postgres")

	if err := goose.Down(s.DBConn, migrationsDir); err != nil {
		return fmt.Errorf("goose down failed: %w", err)
	}

	if err := goose.Up(s.DBConn, migrationsDir); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}

	fmt.Printf("Database reset successfully\n")
	return nil
}

func HandlerDeleteUsers(s *State, cmd Command) error {
	err := s.DB.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete users: %w", err)
	}
	fmt.Printf("All users deleted successfully\n")
	return nil
}
