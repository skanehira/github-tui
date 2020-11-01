package ui

import (
	"fmt"
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
		{Text: i.Author, Color: tcell.ColorYellow},
		{Text: i.Title, Color: tcell.ColorWhite},
	}

	return f
}

var (
	IssueUI *SelectListUI
)

func NewIssueUI(viewUpdater func(text string)) *SelectListUI {
	queries := []string{
		fmt.Sprintf("repo:%s/%s", config.GitHub.Owner, config.GitHub.Repo),
		"is:issue",
	}

	getList := func(cursor *string) ([]List, github.PageInfo) {
		queries := queries
		for _, q := range strings.Split(filterQuery, " ") {
			if strings.Contains(q, "type:pr") || strings.Contains(q, "is:pr") {
				continue
			}
			queries = append(queries, q)
		}
		query := strings.Join(queries, " ")

		v := map[string]interface{}{
			"query":  githubv4.String(query),
			"first":  githubv4.Int(30),
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
				Number: strconv.Itoa(int(node.Issue.Number)),
				State:  string(node.Issue.State),
				Author: string(node.Issue.Author.Login),
				Title:  string(node.Issue.Title),
				Body:   string(node.Issue.Body),
			}

			labels := make([]string, len(node.Issue.Labels.Nodes))
			for i, l := range node.Issue.Labels.Nodes {
				labels[i] = string(l.Name)
			}
			issue.Labels = labels

			assignees := make([]string, len(node.Issue.Assignees.Nodes))
			for i, a := range node.Issue.Assignees.Nodes {
				assignees[i] = string(a.Login)
			}
			issue.Assignees = assignees
			issues[i] = issue
		}
		return issues, resp.PageInfo
	}

	capture := func(event *tcell.EventKey) *tcell.EventKey {
		return UI.Capture(event)
	}

	header := []string{
		"",
		"Number",
		"State",
		"Author",
		"Title",
	}

	init := func(ui *SelectListUI) {
		UI.updater(func() {
			if len(ui.list) > 0 {
				viewUpdater(ui.list[0].(*Issue).Body)
			}
		})
	}

	ui := NewSelectListUI("issue list", header, tcell.ColorBlue, getList, capture, init)

	ui.SetSelectionChangedFunc(func(row, col int) {
		if row > 0 {
			UI.updater(func() {
				viewUpdater(ui.list[row-1].(*Issue).Body)
			})
		}
	})

	IssueUI = ui
	return ui
}
