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

func CreateIssue(input githubv4.CreateIssueInput) error {
	var m MutateCreateIssue
	return client.Mutate(context.Background(), &m, input, nil)
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

func GetRepo(variables map[string]interface{}) (*Repository, error) {
	var q struct {
		Repository `graphql:"repository(owner: $owner, name: $name)"`
	}
	if err := client.Query(context.Background(), &q, variables); err != nil {
		return nil, err
	}
	return &q.Repository, nil
}

func GetIssue(variables map[string]interface{}) (*Issues, error) {
	var q struct {
		Search Issues `graphql:"search(query: $query, type: ISSUE, first: $first, after: $cursor)"`
	}
	if err := client.Query(context.Background(), &q, variables); err != nil {
		return nil, err
	}

	issues := &Issues{
		Nodes:    q.Search.Nodes,
		PageInfo: q.Search.PageInfo,
	}
	return issues, nil
}

func GetIssueTemplates(variables map[string]interface{}) ([]IssueTemplate, error) {
	var q struct {
		Repository struct {
			IssueTemplates []IssueTemplate
		} `graphql:"repository(name: $name, owner: $owner)"`
	}
	if err := client.Query(context.Background(), &q, variables); err != nil {
		return nil, err
	}
	return q.Repository.IssueTemplates, nil
}

func ReopenIssue(id string) error {
	input := githubv4.ReopenIssueInput{
		IssueID: githubv4.String(id),
	}

	var m MutateOpenIsseue

	return client.Mutate(context.Background(), &m, input, nil)
}

func CloseIssue(id string) error {
	input := githubv4.CloseIssueInput{
		IssueID: githubv4.String(id),
	}

	var m MutateCoseIssue
	return client.Mutate(context.Background(), &m, input, nil)
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

func GetRepoProjects(variables map[string]interface{}) (*Projects, error) {
	var q struct {
		Repository struct {
			Projects `graphql:"projects(first: $first, after: $cursor, orderBy: {field: CREATED_AT, direction: DESC})"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}
	if err := client.Query(context.Background(), &q, variables); err != nil {
		return nil, err
	}
	return &q.Repository.Projects, nil
}

func GetRepoAssignableUsers(variables map[string]interface{}) (*AssignableUsers, error) {
	var q struct {
		Repository struct {
			AssignableUsers `graphql:"assignableUsers(first: $first, after: $cursor)"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}
	if err := client.Query(context.Background(), &q, variables); err != nil {
		return nil, err
	}
	return &q.Repository.AssignableUsers, nil
}
