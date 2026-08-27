package main

import (
	"fmt"

	"github.com/SlothEfficiency/Gator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c)
	config.SetUser("Pascal")

	d, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)
}
