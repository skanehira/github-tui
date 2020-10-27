package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
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
