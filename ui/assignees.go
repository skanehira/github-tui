package ui

import (
	"github.com/gdamore/tcell/v2"
)

var (
	AssigneesUI *SelectUI
)

type AssignableUser struct {
	Login string
}

func (a *AssignableUser) Key() string {
	return a.Login
}

func (a *AssignableUser) Fields() []Field {
	return []Field{
		{Text: a.Login, Color: tcell.ColorFuchsia},
	}
}

func NewAssignableUI() *SelectUI {
	//getList := func(cursor *string) ([]List, github.PageInfo) {
	//	v := map[string]interface{}{
	//		"owner":  githubv4.String(config.GitHub.Owner),
	//		"name":   githubv4.String(config.GitHub.Repo),
	//		"first":  githubv4.Int(100),
	//		"cursor": (*githubv4.String)(cursor),
	//	}
	//	resp, err := github.GetRepoAssignableUsers(v)
	//	if err != nil {
	//		return nil, github.PageInfo{}
	//	}

	//	assignees := make([]List, len(resp.Nodes))
	//	for i, p := range resp.Nodes {
	//		assignees[i] = &AssignableUser{
	//			Login: string(p.Login),
	//		}
	//	}
	//	return assignees, resp.PageInfo
	//}

	capture := func(event *tcell.EventKey) *tcell.EventKey {
		return UI.Capture(event)
	}

	ui := NewSelectListUI("assignable user list", nil, tcell.ColorFuchsia, nil, capture, nil)
	AssigneesUI = ui
	return ui
}
