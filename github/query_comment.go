package github

import (
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/domain"
)

type Comment struct {
	ID     githubv4.String
	Author struct {
		Login githubv4.String
	}
	UpdatedAt githubv4.DateTime
	Body      githubv4.String
	URL       githubv4.URI
}

func (c *Comment) ToDomain() *domain.Comment {
	comment := &domain.Comment{
		ID:        string(c.ID),
		Author:    string(c.Author.Login),
		UpdatedAt: c.UpdatedAt.Local().Format("2006/01/02 15:04:05"),
		URL:       c.URL.String(),
		Body:      string(c.Body),
	}
	return comment
}
