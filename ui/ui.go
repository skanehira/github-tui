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

func (ui *ui) canFocus() bool {
	fs := ui.app.GetFocus()
	if fs == nil {
		return false
	}
	switch fs.(type) {
	case *FilterUI, *SelectUI, *viewUI:
		return true
	}
	return false
}

func (ui *ui) toNextUI() {
	if !ui.canFocus() {
		return
	}
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
	if !ui.canFocus() {
		return
	}
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

func (ui *ui) Modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func (ui *ui) Message(msg string, focusFunc func()) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			ui.pages.RemovePage("message").ShowPage("main")
			focusFunc()
		})
	ui.pages.AddAndSwitchToPage("message", ui.Modal(modal, 80, 29), true).ShowPage("main")
}

func (ui *ui) Confirm(msg, doLabel string, doFunc func() error, focusFunc func()) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{doLabel, "Cancel"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			ui.pages.RemovePage("modal").ShowPage("main")
			focusFunc()
			if buttonLabel == doLabel {
				if err := doFunc(); err != nil {
					ui.Message(err.Error(), func() {
						focusFunc()
					})
				}
			}
		})
	ui.pages.AddAndSwitchToPage("modal", ui.Modal(modal, 80, 29), true).ShowPage("main")
}

func (ui *ui) Start() error {
	NewFilterUI()
	NewViewUI(UIKindIssueView)
	NewViewUI(UIKindCommentView)
	NewIssueUI()
	NewLabelsUI()
	NewMilestoneUI()
	NewProjectUI()
	NewAssignableUI()
	NewCommentUI()
	NewSearchUI()

	ui.primitives = []Primitive{IssueFilterUI, AssigneesUI, LabelUI, MilestoneUI,
		ProjectUI, IssueUI, IssueViewUI, CommentUI, CommentViewUI}
	ui.primitiveLen = len(ui.primitives)

	// for readability
	row, col, rowSpan, colSpan := 0, 0, 0, 0

	grid := tview.NewGrid().SetRows(1, 0, 0, 0, 0, 0, 0, 0, 0, 1).
		AddItem(IssueFilterUI, row, col, rowSpan+1, colSpan+3, 0, 0, true).
		AddItem(IssueUI, row+1, col+1, rowSpan+4, colSpan+3, 0, 0, true).
		AddItem(AssigneesUI, row+1, col, rowSpan+1, colSpan+1, 0, 0, true).
		AddItem(LabelUI, row+2, col, rowSpan+1, colSpan+1, 0, 0, true).
		AddItem(MilestoneUI, row+3, col, rowSpan+1, colSpan+1, 0, 0, true).
		AddItem(ProjectUI, row+4, col, rowSpan+1, colSpan+1, 0, 0, true).
		AddItem(CommentUI, row+5, col, rowSpan+4, colSpan+4, 0, 0, true).
		AddItem(IssueViewUI, row+1, col+4, rowSpan+4, colSpan+3, 0, 0, true).
		AddItem(CommentViewUI, row+5, col+4, rowSpan+4, colSpan+3, 0, 0, true).
		AddItem(SearchUI, row+9, col, rowSpan+1, colSpan+7, 0, 0, true)

	ui.pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	ui.app.SetRoot(ui.pages, true)

	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlN:
			UI.toNextUI()
		case tcell.KeyCtrlP:
			UI.toPrevUI()
		}
		return event
	})

	ui.current = 5
	ui.app.SetFocus(IssueUI)
	IssueUI.focus()

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
