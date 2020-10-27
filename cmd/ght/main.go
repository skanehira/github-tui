package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/github"
)

var Config config

type config struct {
	GitHub struct {
		Token string `yaml:"token"`
	} `yaml:"github"`
	ConfigFile string
}

func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	configFile := filepath.Join(configDir, "ght", "config.yaml")

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var conf config
	if err := yaml.Unmarshal(b, &conf); err != nil {
		log.Fatal(err)
	}

	if conf.GitHub.Token == "" {
		log.Fatal("github token is empty")
	}

	conf.ConfigFile = configFile

	Config = conf
}

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
