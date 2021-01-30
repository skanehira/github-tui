package ui

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/config"
	"github.com/skanehira/ght/github"
	"github.com/skanehira/ght/utils"
)

type Issue struct {
	Number    string
	State     string
	Title     string
	Body      string
	Author    string
	URL       string
	Labels    []Item
	Assignees []Item
	Comments  []Item
	MileStone []Item
	Projects  []Item
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
	IssueUI *SelectUI
)

func NewIssueUI() {
	queries := []string{
		fmt.Sprintf("repo:%s/%s", config.GitHub.Owner, config.GitHub.Repo),
		"is:issue",
		"state:open",
	}

	getList := func(cursor *string) ([]Item, *github.PageInfo) {
		queries := queries
		for _, q := range strings.Split(filterQuery, " ") {
			// execlude PR
			if strings.Contains(q, "type:pr") || strings.Contains(q, "is:pr") {
				continue
			}

			// give priority to filterQuery's state
			if strings.Contains(q, "state:") {
				queries = append(queries[:2], q)
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
			return nil, nil
		}

		issues := make([]Item, len(resp.Nodes))
		for i, node := range resp.Nodes {
			issue := &Issue{
				Number: strconv.Itoa(int(node.Issue.Number)),
				State:  string(node.Issue.State),
				Author: string(node.Issue.Author.Login),
				URL:    node.Issue.URL.String(),
				Title:  string(node.Issue.Title),
				Body:   string(node.Issue.Body),
			}

			labels := make([]Item, len(node.Issue.Labels.Nodes))
			for i, l := range node.Issue.Labels.Nodes {
				labels[i] = &Label{
					Name: string(l.Name),
				}
			}
			issue.Labels = labels

			assignees := make([]Item, len(node.Issue.Assignees.Nodes))
			for i, a := range node.Issue.Assignees.Nodes {
				assignees[i] = &AssignableUser{
					Login: string(a.Login),
				}
			}
			issue.Assignees = assignees

			comments := make([]Item, len(node.Issue.Comments.Nodes))
			for i, c := range node.Issue.Comments.Nodes {
				comments[i] = &Comment{
					ID:        string(c.ID),
					Author:    string(c.Author.Login),
					UpdatedAt: c.UpdatedAt.Local().Format("2006/01/02 15:04:05"),
					URL:       c.URL.String(),
					Body:      string(c.Body),
				}
			}
			issue.Comments = comments

			if !reflect.ValueOf(node.Issue.Milestone).IsZero() {
				issue.MileStone = append(issue.MileStone, &Milestone{
					ID:    string(node.Issue.Milestone.ID),
					Title: string(node.Issue.Milestone.Title),
					URL:   node.Issue.Milestone.URL.String(),
				})
			}

			projects := make([]Item, len(node.Issue.ProjectCards.Nodes))
			for i, card := range node.Issue.ProjectCards.Nodes {
				projects[i] = &Project{
					Name: string(card.Project.Name),
					URL:  card.Project.URL.String(),
				}
			}
			issue.Projects = projects

			issues[i] = issue

		}
		return issues, &resp.PageInfo
	}

	capture := func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlO:
			var urls []string
			if len(IssueUI.selected) == 0 {
				data := IssueUI.GetSelect()
				if data != nil {
					urls = append(urls, data.(*Issue).URL)
				}
			} else {
				for _, s := range IssueUI.selected {
					urls = append(urls, s.(*Issue).URL)
				}
			}

			for _, url := range urls {
				if err := utils.OpenBrowser(url); err != nil {
					log.Println(err)
				}
			}
		}
		return UI.Capture(event)
	}

	header := []string{
		"",
		"Number",
		"State",
		"Author",
		"Title",
	}

	init := func(ui *SelectUI) {
		if len(ui.items) > 0 {
			issue := ui.items[0].(*Issue)
			IssueViewUI.updateView(issue.Body)
			if len(issue.Comments) > 0 {
				CommentUI.SetList(issue.Comments)
				CommentViewUI.updateView(issue.Comments[0].(*Comment).Body)
			}
			if len(issue.Assignees) > 0 {
				AssigneesUI.SetList(issue.Assignees)
			}

			if len(issue.Labels) > 0 {
				LabelUI.SetList(issue.Labels)
			}

			if len(issue.MileStone) > 0 {
				MilestoneUI.SetList(issue.MileStone)
			}

			if len(issue.Projects) > 0 {
				ProjectUI.SetList(issue.Projects)
			}
		}
	}

	ui := NewSelectListUI("issue list", header, tcell.ColorBlue, getList, capture, init)

	ui.SetSelectionChangedFunc(func(row, col int) {
		if row > 0 {
			issue := ui.items[row-1].(*Issue)
			IssueViewUI.updateView(issue.Body)

			if len(issue.Comments) > 0 {
				CommentUI.SetList(issue.Comments)
				CommentViewUI.updateView(issue.Comments[0].(*Comment).Body)
			} else {
				CommentUI.ClearView()
				CommentViewUI.Clear()
			}

			if len(issue.Assignees) > 0 {
				AssigneesUI.SetList(issue.Assignees)
			} else {
				AssigneesUI.ClearView()
			}

			if len(issue.Labels) > 0 {
				LabelUI.SetList(issue.Labels)
			} else {
				LabelUI.ClearView()
			}

			if len(issue.MileStone) > 0 {
				MilestoneUI.SetList(issue.MileStone)
			} else {
				MilestoneUI.ClearView()
			}

			if len(issue.Projects) > 0 {
				ProjectUI.SetList(issue.Projects)
			} else {
				ProjectUI.ClearView()
			}
		}
	})

	IssueUI = ui
}
