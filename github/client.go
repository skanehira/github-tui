package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var client *githubv4.Client

func NewClient(token string) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client = githubv4.NewClient(httpClient)
}

func GetRepos(variables map[string]interface{}) (*Repositories, error) {
	var q struct {
		RepositoryOwner struct {
			Repositories `graphql:"repositories(first: $first, after: $cursor, orderBy: {field: CREATED_AT, direction: DESC})"`
		} `graphql:"repositoryOwner(login: $login)"`
	}

	if err := client.Query(context.Background(), &q, variables); err != nil {
		return nil, err
	}
	return &q.RepositoryOwner.Repositories, nil
}
