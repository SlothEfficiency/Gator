package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/SlothEfficiency/Gator/internal/database"
	"github.com/SlothEfficiency/Gator/internal/rss"
)

func loginCheck(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, c, user)
	}
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	markParameters := database.MarkFeedFetchedParams{
		ID:        nextFeed.ID,
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	err = s.db.MarkFeedFetched(context.Background(), markParameters)
	if err != nil {
		return err
	}

	rssFeed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	for _, rssItem := range rssFeed.Channel.Item {
		fmt.Printf("Title: %s\n", rssItem.Title)
	}
	return err

}
