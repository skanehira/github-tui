package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/config"
	"github.com/skanehira/ght/github"
)

type Milestone struct {
	Title       string
	Description string
}

type MilestoneUI struct {
	updater func(f func())
	*tview.Table
}

func NewMilestoneUI(updater func(f func())) *MilestoneUI {
	ui := &MilestoneUI{
		Table:   tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(0, 0),
		updater: updater,
	}
	ui.SetBorder(true).SetTitle("millestone list").SetTitleAlign(tview.AlignLeft)
	go ui.updateMilestoneList()
	return ui
}

func (ui *MilestoneUI) updateMilestoneList() {
	table := ui.Clear()
	v := map[string]interface{}{
		"owner":  githubv4.String(config.GitHub.Owner),
		"name":   githubv4.String(config.GitHub.Repo),
		"first":  githubv4.Int(100),
		"cursor": (*githubv4.String)(nil),
	}
	resp, err := github.GetRepoMillestones(v)
	if err != nil {
		log.Println(err)
		return
	}

	ui.updater(func() {
		labels := make([]Milestone, len(resp.Nodes))

		for i, m := range resp.Nodes {
			title := string(m.Title)
			description := string(m.Description)
			labels[i] = Milestone{
				Title:       title,
				Description: description,
			}

			table.SetCell(i, 0, tview.NewTableCell(title).
				SetTextColor(tcell.ColorPowderBlue).SetExpansion(1))
		}

		ui.ScrollToBeginning()
	})
}
