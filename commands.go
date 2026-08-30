package main

import (
	"context"
	"fmt"
	"time"

	"github.com/SlothEfficiency/Gator/internal/database"
	"github.com/SlothEfficiency/Gator/internal/rss"
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
	rssFeed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(*rssFeed)
	return nil
}
