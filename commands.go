package main

import (
	"fmt"
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
	err := s.Config.SetUser(cmd.Arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("User %s was logged in.\n", cmd.Arguments[0])
	return nil
}
