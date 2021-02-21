package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var IssueFilterUI *FilterUI

type (
	SetFilterOpt func(ui *FilterUI)
	FilterUI     struct {
		*tview.InputField
	}
)

func NewFilterUI() {
	ui := &FilterUI{
		InputField: tview.NewInputField().SetLabel("Filters").SetLabelWidth(8),
	}
	ui.SetBorderPadding(0, 0, 1, 0)

	ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			go IssueUI.GetList()
		}
		return event
	})
	IssueFilterUI = ui
}

func (ui *FilterUI) SetQuery(query string) {
	ui.SetText(query)
}

func (ui *FilterUI) GetQuery() string {
	return ui.GetText()
}

func (ui *FilterUI) focus() {
}

func (ui *FilterUI) blur() {
}
