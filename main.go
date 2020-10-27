package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/github"
)

func main() {
	github.NewClient(Config.GitHub.Token)

	variables := map[string]interface{}{
		"login":  githubv4.String("skanehira"),
		"first":  githubv4.Int(10),
		"cursor": (*githubv4.String)(nil),
	}

	repos, err := github.GetRepos(variables)
	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(repos); err != nil {
		log.Fatal(err)
	}
}
