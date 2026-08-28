package main

import (
	"github.com/SlothEfficiency/Gator/internal/config"
	"github.com/SlothEfficiency/Gator/internal/database"
)

type state struct {
	db     *database.Queries
	Config *config.Config
}

type command struct {
	Name      string
	Arguments []string
}
