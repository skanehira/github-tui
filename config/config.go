package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type github struct {
	Owner string
	Repo  string
	Token string `yaml:"token"`
}

type app struct {
	File string `yaml:"file"`
}

var (
	GitHub github
	App    app
)

func Init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	configFile := filepath.Join(configDir, "ght", "config.yaml")

	b, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logFile := filepath.Join(configDir, "ght", "debug.log")
	output, err := os.Create(logFile)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(output)

	var conf struct {
		GitHub github `yaml:"github"`
	}

	if err := yaml.Unmarshal(b, &conf); err != nil {
		log.Fatal(err)
	}

	if conf.GitHub.Token == "" {
		log.Fatal("github token is empty")
	}

	App.File = configFile

	GitHub = conf.GitHub
}
