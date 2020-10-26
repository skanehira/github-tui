package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	GitHub struct {
		Token string `yaml:"token"`
	} `yaml:"github"`
}

func NewConfig() Config {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	configFile := filepath.Join(configDir, "ght", "config.yaml")

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	if err := yaml.Unmarshal(b, &config); err != nil {
		log.Fatal(err)
	}

	if config.GitHub.Token == "" {
		log.Fatal("github token is empty")
	}

	return config
}
