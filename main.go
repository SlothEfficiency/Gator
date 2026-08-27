package main

import (
	"fmt"
	"os"

	"github.com/SlothEfficiency/Gator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	s := state{
		Config: &c,
	}
	commands := commands{
		commands: map[string]func(*state, command) error{},
	}

	// register all possible commands
	commands.register("login", handlerLogin)

	input := os.Args
	if len(input) < 2 {
		println("Require a command.")
		os.Exit(1)
	}

	cmd := command{
		Name:      input[1],
		Arguments: input[2:],
	}

	err = commands.run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
