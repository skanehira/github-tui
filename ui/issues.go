package ui

import (
	"log"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/config"
	"github.com/skanehira/ght/github"
)

type Issue struct {
	Number    string
	State     string
	Title     string
	Body      string
	Author    string
	Labels    []string
	Assignees []string
}

func (i *Issue) Key() string {
	return i.Number
}

func (i *Issue) Fields() []Field {
	stateColor := tcell.ColorGreen
	if i.State == "CLOSED" {
		stateColor = tcell.ColorRed
	}

	f := []Field{
		{Text: i.Number, Color: tcell.ColorBlue},
		{Text: i.State, Color: stateColor},
		{Text: i.Title, Color: tcell.ColorWhite},
		{Text: i.Author, Color: tcell.ColorYellow},
		{Text: strings.Join(i.Labels, ","), Color: tcell.ColorAqua},
		{Text: strings.Join(i.Assignees, ","), Color: tcell.ColorOlive},
	}

	return f
}

func NewIssueUI(updater func(f func()), viewUpdater func(text string)) *SelectListUI {
	getList := func(cursor *string) ([]List, github.PageInfo) {
		v := map[string]interface{}{
			"owner":  githubv4.String(config.GitHub.Owner),
			"name":   githubv4.String(config.GitHub.Repo),
			"first":  githubv4.Int(100),
			"cursor": (*githubv4.String)(cursor),
		}
		resp, err := github.GetIssue(v)
		if err != nil {
			log.Println(err)
			return nil, github.PageInfo{}
		}

		issues := make([]List, len(resp.Nodes))
		for i, node := range resp.Nodes {
			issue := &Issue{
				Number: strconv.Itoa(int(node.Number)),
				State:  string(node.State),
				Author: string(node.Author.Login),
				Title:  string(node.Title),
				Body:   string(node.Body),
			}

			labels := make([]string, len(node.Labels.Nodes))
			for i, l := range node.Labels.Nodes {
				labels[i] = string(l.Name)
			}
			issue.Labels = labels

			assignees := make([]string, len(node.Assignees.Nodes))
			for i, a := range node.Assignees.Nodes {
				assignees[i] = string(a.Login)
			}
			issue.Assignees = assignees
			issues[i] = issue
		}
		return issues, resp.PageInfo
	}

	capture := func(event *tcell.EventKey) *tcell.EventKey {
		return event
	}

	header := []string{
		"",
		"Number",
		"State",
		"Title",
		"Author",
		"Labels",
		"Assignees",
	}

	init := func(ui *SelectListUI) {
		ui.updater(func() {
			viewUpdater(ui.list[0].(*Issue).Body)
		})
	}

	ui := NewSelectListUI("issue list", header, updater, tcell.ColorBlue, getList, capture, init)

	ui.SetSelectionChangedFunc(func(row, col int) {
		if row > 0 {
			ui.updater(func() {
				viewUpdater(ui.list[row-1].(*Issue).Body)
			})
		}
	})
	return ui
}
