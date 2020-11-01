package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/config"
	"github.com/skanehira/ght/github"
)

type AssignableUser struct {
	Login string
}

func (a *AssignableUser) Key() string {
	return a.Login
}

func (a *AssignableUser) Fields() []Field {
	return []Field{
		{Text: a.Login, Color: tcell.ColorBlue},
	}
}

func NewAssignableUI(updater func(f func())) *SelectListUI {
	getList := func(cursor *string) ([]List, github.PageInfo) {
		v := map[string]interface{}{
			"owner":  githubv4.String(config.GitHub.Owner),
			"name":   githubv4.String(config.GitHub.Repo),
			"first":  githubv4.Int(100),
			"cursor": (*githubv4.String)(cursor),
		}
		resp, err := github.GetRepoAssignableUsers(v)
		if err != nil {
			return nil, github.PageInfo{}
		}

		assignees := make([]List, len(resp.Nodes))
		for i, p := range resp.Nodes {
			assignees[i] = &AssignableUser{
				Login: string(p.Login),
			}
		}
		return assignees, resp.PageInfo
	}

	capture := func(event *tcell.EventKey) *tcell.EventKey {
		return event
	}

	return NewSelectListUI("assibnable user list", nil, updater, tcell.ColorBlue, getList, capture, nil)
}
