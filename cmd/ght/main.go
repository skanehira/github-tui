package main

import (
	"log"

	"github.com/skanehira/ght/config"
	"github.com/skanehira/ght/github"
	"github.com/skanehira/ght/ui"
)

func main() {
	config.Init()
	github.NewClient(config.GitHub.Token)
	if err := ui.New().Start(); err != nil {
		log.Fatal(err)
	}
}
