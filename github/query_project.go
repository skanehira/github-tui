package github

import (
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/domain"
)

type Project struct {
	Name githubv4.String
	URL  githubv4.URI
}

func (p *Project) ToDomain() *domain.Project {
	project := &domain.Project{
		Name: string(p.Name),
		URL:  p.URL.String(),
	}
	return project
}

type Projects struct {
	Nodes    []Project
	PageInfo PageInfo
}
