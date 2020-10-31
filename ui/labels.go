package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/config"
	"github.com/skanehira/ght/github"
)

type Label struct {
	Name        string
	Description string
}

func (l *Label) Key() string {
	return l.Name
}

func (l *Label) Fields() []Field {
	return []Field{
		{Text: l.Name, Color: tcell.ColorYellow},
	}
}

func NewLabelsUI(updater func(f func())) *SelectListUI {
	getList := func(cursor *string) ([]List, github.PageInfo) {
		v := map[string]interface{}{
			"owner":  githubv4.String(config.GitHub.Owner),
			"name":   githubv4.String(config.GitHub.Repo),
			"first":  githubv4.Int(100),
			"cursor": (*githubv4.String)(cursor),
		}
		resp, err := github.GetRepoLabels(v)
		if err != nil {
			log.Println(err)
			return nil, github.PageInfo{}
		}

		labels := make([]List, len(resp.Nodes))
		for i, l := range resp.Nodes {
			name := string(l.Name)
			description := string(l.Description)
			labels[i] = &Label{
				Name:        name,
				Description: description,
			}
		}
		return labels, resp.PageInfo
	}

	capture := func(event *tcell.EventKey) *tcell.EventKey {
		return event
	}

	return NewSelectListUI("label list", nil, updater, tcell.ColorYellow, getList, capture)
}
