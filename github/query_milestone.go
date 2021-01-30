package github

import (
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/domain"
)

type Milestone struct {
	ID          githubv4.String
	Title       githubv4.String
	State       githubv4.String
	Description githubv4.String
	URL         githubv4.URI
}

func (m *Milestone) ToDomain() *domain.Milestone {
	milestone := &domain.Milestone{
		ID:          string(m.ID),
		Title:       string(m.Title),
		State:       string(m.State),
		Description: string(m.Description),
		URL:         m.URL.String(),
	}
	return milestone
}

type Milestones struct {
	Nodes    []Milestone
	PageInfo PageInfo
}
