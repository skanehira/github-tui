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

func GetIssue(variables map[string]interface{}) (*Issues, error) {
	var q struct {
		Repository struct {
			Issues `graphql:"issues(first: $first, after: $cursor, orderBy: {field: CREATED_AT, direction: DESC})"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}
	if err := client.Query(context.Background(), &q, variables); err != nil {
		return nil, err
	}
	return &q.Repository.Issues, nil
}

func GetRepoLabels(variables map[string]interface{}) (*Labels, error) {
	var q struct {
		Repository struct {
			Labels `graphql:"labels(first: $first, after: $cursor, orderBy: {field: CREATED_AT, direction: DESC})"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}
	if err := client.Query(context.Background(), &q, variables); err != nil {
		return nil, err
	}
	return &q.Repository.Labels, nil
}

func GetRepoMillestones(variables map[string]interface{}) (*Milestones, error) {
	var q struct {
		Repository struct {
			Milestones `graphql:"milestones(first: $first, after: $cursor, orderBy: {field: CREATED_AT, direction: DESC})"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}
	if err := client.Query(context.Background(), &q, variables); err != nil {
		return nil, err
	}
	return &q.Repository.Milestones, nil
}
