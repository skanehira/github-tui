package ui

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/config"
	"github.com/skanehira/ght/domain"
	"github.com/skanehira/ght/github"
	"github.com/skanehira/ght/utils"
)

var IssueUI *SelectUI

func NewIssueUI() {
	opt := func(ui *SelectUI) {
		// initial query
		queries := []string{
			fmt.Sprintf("repo:%s/%s", config.GitHub.Owner, config.GitHub.Repo),
			"state:open",
		}

		IssueFilterUI.SetQuery(strings.Join(queries, " "))

		ui.getList = func(cursor *string) ([]domain.Item, *github.PageInfo) {
			var queries []string
			query := IssueFilterUI.GetQuery()

			if !strings.Contains(query, "is:issue") {
				queries = append(queries, "is:issue")
			}

			for _, q := range strings.Split(query, " ") {
				// execlude Pull request
				if strings.Contains(q, "type:pr") || strings.Contains(q, "is:pr") {
					continue
				}

				queries = append(queries, q)
			}
			query = strings.Join(queries, " ")
			IssueFilterUI.SetQuery(query)

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

			issues := make([]domain.Item, len(resp.Nodes))
			for i, node := range resp.Nodes {
				issues[i] = node.Issue.ToDomain()
			}
			return issues, &resp.PageInfo
		}

		getSelectedIssues := func() []*domain.Issue {
			var issues []*domain.Issue
			if len(IssueUI.selected) == 0 {
				data := IssueUI.GetSelect()
				if data != nil {
					issues = append(issues, data.(*domain.Issue))
				}
			} else {
				for _, item := range IssueUI.selected {
					issues = append(issues, item.(*domain.Issue))
				}
			}
			return issues
		}

		ui.capture = func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'y':
				var urls []string
				for _, issue := range getSelectedIssues() {
					urls = append(urls, issue.URL)
				}

				url := strings.Join(urls, "\n")
				if err := clipboard.WriteAll(url); err != nil {
					log.Println(err)
				}
				IssueUI.ClearSelected()
			case 'o':
				go func() {
					var wg sync.WaitGroup
					for _, issue := range getSelectedIssues() {
						wg.Add(1)
						go func(issue *domain.Issue) {
							defer wg.Done()
							if err := github.ReopenIssue(issue.ID); err != nil {
								log.Println(err)
								return
							}
							issue.State = "OPEN"
						}(issue)
					}
					wg.Wait()
					IssueUI.ClearSelected()
					IssueUI.UpdateView()
				}()
			case 'c':
				go func() {
					var wg sync.WaitGroup
					for _, issue := range getSelectedIssues() {
						wg.Add(1)
						go func(issue *domain.Issue) {
							defer wg.Done()
							if err := github.CloseIssue(issue.ID); err != nil {
								log.Println(err)
								return
							}
							issue.State = "CLOSED"
						}(issue)
					}
					wg.Wait()
					IssueUI.ClearSelected()
					IssueUI.UpdateView()
				}()
			}
			switch event.Key() {
			case tcell.KeyCtrlO:
				for _, issue := range getSelectedIssues() {
					if err := utils.Open(issue.URL); err != nil {
						log.Println(err)
					}
				}
			}
			return event
		}

		ui.header = []string{
			"",
			"Repo",
			"Number",
			"State",
			"Author",
			"Title",
		}

		ui.hasHeader = len(ui.header) > 0
	}

	ui := NewSelectListUI(UIKindIssue, tcell.ColorBlue, opt)

	ui.SetSelectionChangedFunc(func(row, col int) {
		updateUIRelatedIssue(ui, row)
	})

	IssueUI = ui
}

func updateUIRelatedIssue(ui *SelectUI, row int) {
	if row > 0 && row < len(ui.items) {
		issue := ui.items[row-1].(*domain.Issue)
		IssueViewUI.updateView(issue.Body)

		if len(issue.Comments) > 0 {
			CommentUI.SetList(issue.Comments)
			CommentViewUI.updateView(issue.Comments[0].(*domain.Comment).Body)
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
}
