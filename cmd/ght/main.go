package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/skanehira/ght/config"
	"github.com/skanehira/ght/github"
	"github.com/skanehira/ght/ui"
)

type Repo struct {
	Owner string
	Name  string
}

func main() {
	config.Init()
	getRepoInfo()
	github.NewClient(config.GitHub.Token)
	if err := ui.New().Start(); err != nil {
		log.Fatal(err)
	}
}

func getRepoInfo() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		args := strings.Split(flag.Arg(0), "/")
		if len(args) < 2 {
			log.Fatal("invalid args")
		}
		config.GitHub.Owner = args[0]
		config.GitHub.Repo = args[1]
	} else {
		repo, err := getOwnerRepo()
		if err != nil {
			log.Fatalf("invalid repo: %s", err)
		}
		config.GitHub.Owner = repo.Owner
		config.GitHub.Repo = repo.Name
	}
}

func getOwnerRepo() (*Repo, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return nil, err
	}
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return nil, errors.New("current directory is not git repository")
	}
	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.CombinedOutput()

	result := strings.TrimRight(string(out), "\r\n")
	if err != nil {
		return nil, err
	}

	return parseRemote(result)
}

func parseRemote(remote string) (*Repo, error) {
	if strings.HasSuffix(remote, ".git") {
		remote = strings.TrimRight(remote, ".git")
	}
	var ownerRepo []string
	if strings.HasPrefix(remote, "ssh") {
		p := strings.Split(remote, "/")
		if len(p) < 1 {
			return nil, fmt.Errorf("cannot get owner/repo from remote: %s", remote)
		}
		ownerRepo = p[len(p)-2:]
	} else if strings.HasPrefix(remote, "git") {
		p := strings.Split(remote, ":")
		if len(p) < 1 {
			return nil, fmt.Errorf("cannot get owner/repo from remote: %s", remote)
		}
		ownerRepo = strings.Split(p[1], "/")
	} else if strings.HasPrefix(remote, "http") || strings.HasPrefix(remote, "https") {
		p := strings.Split(remote, "/")
		if len(p) < 1 {
			return nil, fmt.Errorf("cannot get owner/repo from remote: %s", remote)
		}
		ownerRepo = p[len(p)-2:]
	}

	repo := Repo{
		Owner: ownerRepo[0],
		Name:  ownerRepo[1],
	}

	return &repo, nil
}
