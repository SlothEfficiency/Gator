package main

import (
	"github.com/SlothEfficiency/gator/internal/config"
	"github.com/SlothEfficiency/gator/internal/database"
)

type state struct {
	db     *database.Queries
	Config *config.Config
}

type command struct {
	Name      string
	Arguments []string
}
