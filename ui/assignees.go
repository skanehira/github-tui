package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/config"
	"github.com/skanehira/ght/github"
)

type AssignableUsers struct {
	Login string
}

type AssigneesUI struct {
	updater func(f func())
	*tview.Table
}

func NewAssignableUI(updater func(f func())) *AssigneesUI {
	ui := &AssigneesUI{
		Table:   tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(0, 0),
		updater: updater,
	}
	ui.SetBorder(true).SetTitle("assignabel user list").SetTitleAlign(tview.AlignLeft)
	go ui.updateAssignees()
	return ui
}

func (ui *AssigneesUI) updateAssignees() {
	table := ui.Clear()
	v := map[string]interface{}{
		"owner":  githubv4.String(config.GitHub.Owner),
		"name":   githubv4.String(config.GitHub.Repo),
		"first":  githubv4.Int(100),
		"cursor": (*githubv4.String)(nil),
	}
	resp, err := github.GetRepoAssignableUsers(v)
	if err != nil {
		log.Println(err)
		return
	}
	ui.updater(func() {
		assignees := make([]AssignableUsers, len(resp.Nodes))

		for i, p := range resp.Nodes {
			login := string(p.Login)
			assignees[i] = AssignableUsers{
				Login: login,
			}

			table.SetCell(i, 0, tview.NewTableCell(login).
				SetTextColor(tcell.ColorYellowGreen).SetExpansion(1))
		}

		ui.ScrollToBeginning()
	})
}
