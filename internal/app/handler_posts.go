package app

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/uLuKaiDev/Gator/internal/database"
)

func HandlerBrowse(s *State, cmd Command, user database.User) error {
	limit := 2
	if len(cmd.Args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = parsedLimit
	}

	ctx := context.Background()

	posts, err := s.DB.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error fetching posts: %w", err)
	}
	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Println("==========================")
		fmt.Printf("Feed: %s\n", post.FeedName)
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Printf("Published At: %s\n", humanizeTime(post.PublishedAt))
		fmt.Println()
	}

	return nil
}

func humanizeTime(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration < 48*time.Hour:
		return "yesterday"
	default:
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	}
}
