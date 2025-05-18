package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/uLuKaiDev/Gator/internal/database"
)

func scrapeFeeds(s *State) {
	ctx := context.Background()

	feed, err := s.DB.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Printf("failed to get next feed: %v", err)
		return
	}

	err = s.DB.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		log.Printf("failed to mark feed as fetched: %v", err)
		return
	}

	fetchedFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		log.Printf("failed to fetch feed: %v", err)
		return
	}

	// Save the fetched feeds to the database
	for _, item := range fetchedFeed.Channel.Item {
		pubDate, err := parsePubDate(item.PubDate)
		if err != nil {
			log.Printf("failed to parse pubDate: %v", err)
			continue
		}

		_, err = s.DB.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		})
		if err != nil {
			// Check if it's a unique violation error on posts_feed_id_url_key
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" && pgErr.Constraint == "posts_feed_id_url_key" {
				// duplicate post, just ignore
				continue
			} else {
				log.Printf("failed to create post: %v", err)
				continue
			}
		}
	}
}

func parsePubDate(dateStr string) (time.Time, error) {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unrecognized time format: %s", dateStr)
}
