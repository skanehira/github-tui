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

const readThisMessage = "read this https://github.com/skanehira/github-tui?tab=readme-ov-file#settings to know more"

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
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}

		log.Fatalf("Could not find configuration file, %s", readThisMessage)
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
		log.Fatalf("cannot deserialize config file: %s", err.Error())
	}

	if conf.GitHub.Token == "" {
		log.Fatalf("github token is empty, %s", readThisMessage)
	}

	App.File = configFile
	GitHub = conf.GitHub
}
