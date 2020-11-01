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
	updater      func(f func())
}

func New() *ui {
	ui := &ui{
		app: tview.NewApplication(),
	}

	ui.updater = func(f func()) {
		go ui.app.QueueUpdateDraw(f)
	}

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
	switch event.Rune() {
	case 'l':
		UI.toNextUI()
	case 'h':
		UI.toPrevUI()
	}

	return event
}

func (ui *ui) Start() error {
	view, viewUpdater := NewViewUI()
	issueUI := NewIssueUI(viewUpdater)
	labelUI := NewLabelsUI()
	milestoneUI := NewMilestoneUI()
	projectUI := NewProjectUI()
	assigneesUI := NewAssignableUI()

	ui.primitives = []Primitive{issueUI, view, assigneesUI, projectUI, labelUI, milestoneUI}
	ui.primitiveLen = len(ui.primitives)
	issueUI.focus()

	grid := tview.NewGrid().
		AddItem(issueUI, 0, 0, 1, 4, 0, 0, true).
		AddItem(view, 1, 0, 3, 2, 0, 0, true).
		AddItem(assigneesUI, 1, 2, 2, 1, 0, 0, true).
		AddItem(labelUI, 1, 3, 2, 1, 0, 0, true).
		AddItem(milestoneUI, 3, 3, 1, 1, 0, 0, true).
		AddItem(projectUI, 3, 2, 1, 1, 0, 0, true)

	ui.pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	ui.app.SetRoot(ui.pages, true)

	if err := ui.app.Run(); err != nil {
		ui.app.Stop()
		return err
	}

	return nil
}
