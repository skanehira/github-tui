package github

import (
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/domain"
)

type Label struct {
	ID          githubv4.ID
	Name        githubv4.String
	Description githubv4.String
	Color       githubv4.String
}

func (l *Label) ToDomain() *domain.Label {
	label := &domain.Label{
		Name:        string(l.Name),
		Description: string(l.Description),
	}
	return label
}

type Labels struct {
	Nodes    []Label
	PageInfo PageInfo
}
