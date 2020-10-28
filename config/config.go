package config

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var conf struct {
		GitHub github `yaml:"github"`
	}

	if err := yaml.Unmarshal(b, &conf); err != nil {
		log.Fatal(err)
	}

	if conf.GitHub.Token == "" {
		log.Fatal("github token is empty")
	}

	flag.Parse()
	if len(flag.Args()) > 0 {
		args := strings.Split(flag.Arg(0), "/")
		if len(args) < 2 {
			log.Fatal("invalid args")
		}
		conf.GitHub.Owner = args[0]
		conf.GitHub.Repo = args[1]
	}

	App.File = configFile

	GitHub = conf.GitHub
}
