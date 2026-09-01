package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/SlothEfficiency/gator/internal/database"
	"github.com/SlothEfficiency/gator/internal/rss"
	"github.com/google/uuid"
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
	if rssFeed == nil {
		return fmt.Errorf("RSS Feed was not reached.")
	}
	for _, rssItem := range rssFeed.Channel.Item {
		//Check whether this post is already saved
		posts, err := s.db.CheckForPostUrl(context.Background(), rssItem.Link)
		if err != nil {
			return err
		}
		if len(posts) == 1 {
			fmt.Printf("RSS Item \"%s\" was already saved.\n", rssItem.Title)
			continue
		}

		// Prepare post parameters
		valid := true
		if rssItem.Description == "" {
			valid = false
		}
		pubAt, err := time.Parse("Mon, 02 Jan 2006 15:04:05 +0000", rssItem.PubDate)
		if err != nil {
			return err
		}
		parameters := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     rssItem.Title,
			Url:       rssItem.Link,
			Description: sql.NullString{
				String: rssItem.Description,
				Valid:  valid,
			},
			PublishedAt: pubAt,
			FeedID:      nextFeed.ID,
		}
		post, err := s.db.CreatePost(context.Background(), parameters)
		fmt.Printf("Post with Title %s was saved.\n", post.Title)
	}
	return err

}
