package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/SlothEfficiency/Gator/internal/config"
	"github.com/SlothEfficiency/Gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("postgres", c.DatabaseURL)
	dbQueries := database.New(db)

	s := state{
		db:     dbQueries,
		Config: &c,
	}

	commands := commands{
		commands: map[string]func(*state, command) error{},
	}

	// register all possible commands
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handleAgg)

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
