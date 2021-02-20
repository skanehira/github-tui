package github

import (
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/domain"
)

type AssignableUser struct {
	Login githubv4.String
}

func (a *AssignableUser) ToDomain() *domain.AssignableUser {
	assignableUser := &domain.AssignableUser{
		Login: string(a.Login),
	}
	return assignableUser
}

type AssignableUsers struct {
	Nodes []struct {
		ID    githubv4.ID
		Login githubv4.String
	}
	PageInfo PageInfo
}
