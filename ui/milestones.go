package ui

import (
	"github.com/gdamore/tcell/v2"
)

var (
	MilestoneUI *SelectListUI
)

type Milestone struct {
	ID    string
	Title string
}

func (m *Milestone) Key() string {
	return m.Title
}

func (m *Milestone) Fields() []Field {
	return []Field{
		{Text: m.Title, Color: tcell.ColorGreen},
	}
}

func NewMilestoneUI() *SelectListUI {
	//getList := func(cursor *string) ([]List, github.PageInfo) {
	//	v := map[string]interface{}{
	//		"owner":  githubv4.String(config.GitHub.Owner),
	//		"name":   githubv4.String(config.GitHub.Repo),
	//		"first":  githubv4.Int(100),
	//		"cursor": (*githubv4.String)(cursor),
	//	}
	//	resp, err := github.GetRepoMillestones(v)
	//	if err != nil {
	//		return nil, github.PageInfo{}
	//	}

	//	milestones := make([]List, len(resp.Nodes))
	//	for i, m := range resp.Nodes {
	//		milestones[i] = &Milestone{
	//			Title: string(m.Title),
	//		}
	//	}

	//	return milestones, resp.PageInfo
	//}

	capture := func(event *tcell.EventKey) *tcell.EventKey {
		return UI.Capture(event)
	}

	ui := NewSelectListUI("milestone list", nil, tcell.ColorGreen, nil, capture, nil)
	MilestoneUI = ui
	return ui
}
