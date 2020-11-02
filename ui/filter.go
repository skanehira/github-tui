package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	filterQuery string
)

type FilterUI struct {
	*tview.Form
}

func NewFilterUI() *FilterUI {
	ui := &FilterUI{
		Form: tview.NewForm().AddInputField("Filters", "", 100, nil, nil),
	}

	ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			filterQuery = ui.GetFormItem(0).(*tview.InputField).GetText()
			IssueUI.GetList()
		}
		return UI.Capture(event)
	})

	return ui
}

func (ui *FilterUI) focus() {
}

func (ui *FilterUI) blur() {
}
