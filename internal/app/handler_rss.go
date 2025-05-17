package app

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/uLuKaiDev/Gator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

const RSSFeedUrl = "https://www.wagslane.dev/index.xml"

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading body: %w", err)
	}

	var rssFeed RSSFeed
	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
	}

	return &rssFeed, nil
}

func HandlerAgg(s *State, cmd Command) error {
	feed, err := fetchFeed(context.Background(), RSSFeedUrl)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	fmt.Printf("Fetched feed: %+v\n", feed)
	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("missing arguments, usage: gator addfeed <name> <feed_url>")
	}

	if s.Config.CurrentUserName == "" {
		return fmt.Errorf("user not set, please set a user first")
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	_, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	id := uuid.New()
	now := time.Now()

	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}
	fmt.Printf("feed created: %v\n", feed)

	_, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}

	return nil
}

func HandlerListFeeds(s *State, cmd Command) error {
	if s.Config.CurrentUserName == "" {
		return fmt.Errorf("user not set, please register a user first")
	}

	feeds, err := s.DB.GetFeedsWithUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds with users: %w", err)
	}

	fmt.Printf("Feeds stored in the database:\n")
	for _, feed := range feeds {
		fmt.Printf("Name: %s, URL: %s, User: %s\n", feed.FeedName, feed.FeedUrl, feed.UserName)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: gator follow <feed_url>")
	}

	feedURL := cmd.Args[0]

	feed, err := s.DB.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("feed not found: %w", err)
	}

	id := uuid.New()
	now := time.Now()

	feedFollow, err := s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}

	fmt.Printf("User %s is now following feed %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	follows, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feed follows: %w", err)
	}

	fmt.Println("Feeds you are following:")
	for _, follow := range follows {
		fmt.Println("-", follow.FeedName)
	}
	return nil
}
