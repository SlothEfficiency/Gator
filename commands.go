package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/SlothEfficiency/gator/internal/database"
	"github.com/google/uuid"
)

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	err := c.commands[cmd.Name](s, cmd)
	return err
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("Missing argument: Username is required")
	}
	if len(cmd.Arguments) >= 2 {
		return fmt.Errorf("Too many argument: Only 1 username is required")
	}

	name := cmd.Arguments[0]

	context := context.Background()
	user, err := s.db.GetUser(context, name)

	if err != nil {
		return err
	}

	err = s.Config.SetUser(cmd.Arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("User %s was logged in.\n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("Missing argument: Username is required")
	}
	if len(cmd.Arguments) >= 2 {
		return fmt.Errorf("Too many argument: Only 1 username is required")
	}
	context := context.Background()
	parameters := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.Arguments[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user, err := s.db.CreateUser(context, parameters)

	if err != nil {
		return err
	}

	fmt.Printf("%s was created as User.\n", user.Name)
	fmt.Println(user)
	return s.Config.SetUser(user.Name)
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("All users deleted.")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func handleAgg(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("agg takes exactly 1 argument (time between rss fetches).")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return err
	}

	if timeBetweenReqs < time.Second {
		return fmt.Errorf("Time has to be at least 1s")
	}
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func handlerAddFeed(s *state, cmd command, u database.User) error {
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("addfeed takes exactly 2 arguments (name, url).")
	}

	parameters := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
		UserID:    u.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), parameters)
	if err != nil {
		return err
	}

	feed_follow_param := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    feed.UserID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feed_follow_param)

	return err
}

func handleListFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("List of all feeds:")
	for _, feed := range feeds {
		fmt.Printf("Name: %s, url: %s, user who created the feed: %s\n", feed.Name, feed.Url, feed.Username)
	}

	return nil
}

func handleFollow(s *state, cmd command, u database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("follow takes exactly 1 argument (url).")
	}

	feedId, err := s.db.GetFeedID(context.Background(), cmd.Arguments[0])
	if err != nil {
		return err
	}

	parameters := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feedId,
		UserID:    u.ID,
	}
	feed, err := s.db.CreateFeedFollow(context.Background(), parameters)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Id: %s, CreatedAt: %s, UpdatedAt: %s, FeedID: %s, UserID: %s, FeedName: %s, UserName: %s\n",
		feed.ID,
		feed.CreatedAt,
		feed.UpdatedAt,
		feed.FeedID,
		feed.UserID,
		feed.FeedName,
		feed.UserName,
	)
	return nil
}

func handleFollowing(s *state, cmd command) error {
	feeds, err := s.db.GetFeedFollowsPerUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("Feed: %s, Followed by: %s\n", feed.FeedName, feed.UserName)
	}
	return nil
}

func handleUnfollow(s *state, cmd command, u database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("unfollow takes exactly 1 argument (url).")
	}

	feedId, err := s.db.GetFeedID(context.Background(), cmd.Arguments[0])
	if err != nil {
		return err
	}

	parameters := database.DeleteFeedFollowParams{
		UserID: u.ID,
		FeedID: feedId,
	}
	err = s.db.DeleteFeedFollow(context.Background(), parameters)
	return err
}

func handleBrowse(s *state, cmd command, u database.User) error {
	limit := 2
	var err error

	if len(cmd.Arguments) == 1 {
		limit, err = strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return err
		}
	} else if len(cmd.Arguments) >= 2 {
		return fmt.Errorf("browse takes maximum 1 argument (limit)")
	}

	parameters := database.GetPostPerUserParams{
		Limit:  int32(limit),
		UserID: u.ID,
	}
	posts, err := s.db.GetPostPerUser(context.Background(), parameters)
	if err != nil {
		return err
	}

	for _, post := range posts {
		description := ""
		if post.Description.Valid == true {
			description = post.Description.String
		}
		fmt.Printf("Title: \"%s\", URL: \"%s\", Description: \"%s\", Published At: \"%v\"\n", post.Title, post.Url, description, post.PublishedAt)
	}
	return nil
}
