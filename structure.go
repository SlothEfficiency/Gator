package main

import (
	"github.com/SlothEfficiency/Gator/internal/config"
)

type state struct {
	Config *config.Config
}

type command struct {
	Name      string
	Arguments []string
}
