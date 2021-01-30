package github

import (
	"reflect"
	"strconv"

	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/domain"
)

type Issue struct {
	ID         githubv4.String
	Repository struct {
		Name githubv4.String
	}
	Number githubv4.Int
	Body   githubv4.String
	State  githubv4.String
	Author struct {
		Login githubv4.String
	}
	Title     githubv4.String
	URL       githubv4.URI
	Labels    Labels `graphql:"labels(first: 10)"`
	Assignees struct {
		Nodes []AssignableUser
	} `graphql:"assignees(first: 10)"`
	ProjectCards struct {
		Nodes []struct {
			Project Project
		}
	} `graphql:"projectCards(first: 10)"`
	Milestone Milestone
	Comments  struct {
		Nodes []Comment
	} `graphql:"comments(first: 100)"`
}

func (i *Issue) ToDomain() *domain.Issue {
	issue := &domain.Issue{
		ID:     string(i.ID),
		Repo:   string(i.Repository.Name),
		Number: strconv.Itoa(int(i.Number)),
		State:  string(i.State),
		Author: string(i.Author.Login),
		URL:    i.URL.String(),
		Title:  string(i.Title),
		Body:   string(i.Body),
	}

	labels := make([]domain.Item, len(i.Labels.Nodes))
	for i, label := range i.Labels.Nodes {
		labels[i] = label.ToDomain()
	}
	issue.Labels = labels

	assignees := make([]domain.Item, len(i.Assignees.Nodes))
	for i, a := range i.Assignees.Nodes {
		assignees[i] = a.ToDomain()
	}
	issue.Assignees = assignees

	comments := make([]domain.Item, len(i.Comments.Nodes))
	for i, comment := range i.Comments.Nodes {
		comments[i] = comment.ToDomain()
	}
	issue.Comments = comments

	if !reflect.ValueOf(i.Milestone).IsZero() {
		issue.MileStone = append(issue.MileStone, i.Milestone.ToDomain())
	}

	projects := make([]domain.Item, len(i.ProjectCards.Nodes))
	for i, card := range i.ProjectCards.Nodes {
		projects[i] = card.Project.ToDomain()
	}
	issue.Projects = projects
	return issue
}

type Issues struct {
	Nodes []struct {
		Issue Issue `graphql:"... on Issue"`
	}
	PageInfo PageInfo
}
