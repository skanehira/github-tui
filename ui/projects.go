package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/config"
	"github.com/skanehira/ght/github"
)

type Project struct {
	Name        string
	Description string
	Color       string
}

type ProjectUI struct {
	updater func(f func())
	*tview.Table
}

func NewProjectUI(updater func(f func())) *ProjectUI {
	ui := &ProjectUI{
		Table:   tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(0, 0),
		updater: updater,
	}
	ui.SetBorder(true).SetTitle("project list").SetTitleAlign(tview.AlignLeft)
	ui.updateProjectList()
	return ui
}

func (ui *ProjectUI) updateProjectList() {
	table := ui.Clear()

	ui.updater(func() {
		v := map[string]interface{}{
			"owner":  githubv4.String(config.GitHub.Owner),
			"name":   githubv4.String(config.GitHub.Repo),
			"first":  githubv4.Int(100),
			"cursor": (*githubv4.String)(nil),
		}
		resp, err := github.GetRepoProjects(v)
		if err != nil {
			log.Println(err)
			return
		}

		labels := make([]Project, len(resp.Nodes))

		for i, p := range resp.Nodes {
			name := string(p.Name)
			labels[i] = Project{
				Name: name,
			}

			table.SetCell(i, 1, tview.NewTableCell(name).
				SetTextColor(tcell.ColorSnow).SetExpansion(1))
		}

		ui.ScrollToBeginning()
	})
}
