package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	UI *ui
)

type Primitive interface {
	focus()
	blur()
	tview.Primitive
}

type ui struct {
	app          *tview.Application
	pages        *tview.Pages
	current      int
	primitives   []Primitive
	primitiveLen int
	updater      chan func()
}

func New() *ui {
	ui := &ui{
		app: tview.NewApplication(),
	}

	ui.updater = make(chan func(), 100)

	UI = ui

	return ui
}

func (ui *ui) toNextUI() {
	ui.primitives[ui.current].blur()
	if ui.primitiveLen-1 > ui.current {
		ui.current++
	} else {
		ui.current = 0
	}
	p := ui.primitives[ui.current]
	p.focus()
	ui.app.SetFocus(p)
}

func (ui *ui) toPrevUI() {
	ui.primitives[ui.current].blur()
	if ui.current == 0 {
		ui.current = ui.primitiveLen - 1
	} else {
		ui.current--
	}
	p := ui.primitives[ui.current]
	p.focus()
	ui.app.SetFocus(p)
}

func (ui *ui) Capture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlN:
		UI.toNextUI()
	case tcell.KeyCtrlP:
		UI.toPrevUI()
	}

	return event
}

func (ui *ui) Start() error {
	view := NewViewUI()
	issueUI := NewIssueUI()
	labelUI := NewLabelsUI()
	milestoneUI := NewMilestoneUI()
	projectUI := NewProjectUI()
	assigneesUI := NewAssignableUI()
	filterUI := NewFilterUI()

	ui.primitives = []Primitive{filterUI, issueUI, view, assigneesUI, projectUI, labelUI, milestoneUI}
	ui.primitiveLen = len(ui.primitives)

	grid := tview.NewGrid().SetRows(3).
		AddItem(filterUI, 0, 0, 1, 4, 0, 0, true).
		AddItem(issueUI, 1, 0, 1, 4, 0, 0, true).
		AddItem(view, 2, 0, 3, 2, 0, 0, true).
		AddItem(assigneesUI, 2, 2, 2, 1, 0, 0, true).
		AddItem(labelUI, 2, 3, 2, 1, 0, 0, true).
		AddItem(milestoneUI, 4, 3, 1, 1, 0, 0, true).
		AddItem(projectUI, 4, 2, 1, 1, 0, 0, true)

	ui.pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	ui.app.SetRoot(ui.pages, true)

	ui.current = 1
	ui.app.SetFocus(issueUI)
	issueUI.focus()

	go func() {
		for f := range UI.updater {
			go ui.app.QueueUpdateDraw(f)
		}
	}()

	if err := ui.app.Run(); err != nil {
		ui.app.Stop()
		return err
	}

	return nil
}
