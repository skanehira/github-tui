package ui

import "github.com/rivo/tview"

type UI struct {
	app     *tview.Application
	pages   *tview.Pages
	updater func(f func())
}

func New() *UI {
	ui := &UI{
		app: tview.NewApplication(),
	}

	ui.updater = func(f func()) {
		go ui.app.QueueUpdateDraw(f)
	}

	return ui
}

func (ui *UI) Start() error {
	view, viewUpdater := NewViewUI()
	labelUI := NewLabelsUI(ui.updater)
	issueUI := NewIssueUI(ui.updater, viewUpdater)
	grid := tview.NewGrid().SetRows(0, 0, 0, 0).SetColumns(0, 0, 0, 0).
		AddItem(issueUI, 0, 0, 1, 4, 0, 0, true).
		AddItem(view, 1, 0, 3, 2, 0, 0, true).
		AddItem(labelUI, 1, 2, 3, 1, 0, 0, true)

	ui.pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	ui.app.SetRoot(ui.pages, true)

	if err := ui.app.Run(); err != nil {
		ui.app.Stop()
		return err
	}

	return nil
}
