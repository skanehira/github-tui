package ui

import "github.com/rivo/tview"

type UI struct {
	app     *tview.Application
	pages   *tview.Pages
	updater chan func()
}

func New() *UI {
	return &UI{
		app:     tview.NewApplication(),
		updater: make(chan func(), 0),
	}
}

func (ui *UI) Start() error {
	issue := newIssueUI(ui.updater)
	grid := tview.NewGrid().SetRows(0).
		AddItem(issue, 0, 0, 1, 1, 0, 0, true)

	ui.pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	ui.app.SetRoot(ui.pages, true)

	go func() {
		for f := range ui.updater {
			ui.app.QueueUpdateDraw(f)
		}
	}()

	if err := ui.app.Run(); err != nil {
		ui.app.Stop()
		return err
	}

	return nil
}
