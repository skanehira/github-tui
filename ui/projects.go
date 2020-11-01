package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/config"
	"github.com/skanehira/ght/github"
)

type Project struct {
	Name string
}

func (p *Project) Key() string {
	return p.Name
}

func (p *Project) Fields() []Field {
	return []Field{
		{Text: p.Name, Color: tcell.ColorLightSalmon},
	}
}

func NewProjectUI() *SelectListUI {
	getList := func(cursor *string) ([]List, github.PageInfo) {
		v := map[string]interface{}{
			"owner":  githubv4.String(config.GitHub.Owner),
			"name":   githubv4.String(config.GitHub.Repo),
			"first":  githubv4.Int(100),
			"cursor": (*githubv4.String)(cursor),
		}
		resp, err := github.GetRepoProjects(v)
		if err != nil {
			return nil, github.PageInfo{}
		}

		projects := make([]List, len(resp.Nodes))
		for i, m := range resp.Nodes {
			projects[i] = &Project{
				Name: string(m.Name),
			}
		}

		return projects, resp.PageInfo
	}

	capture := func(event *tcell.EventKey) *tcell.EventKey {
		return UI.Capture(event)
	}

	return NewSelectListUI("project list", nil, tcell.ColorLightSalmon, getList, capture, nil)
}
