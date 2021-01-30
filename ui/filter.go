package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var IssueFilterUI *FilterUI

type (
	SetFilterOpt func(ui *FilterUI)
	FilterUI     struct {
		*tview.Form
	}
)

func NewFilterUI() {
	ui := &FilterUI{
		Form: tview.NewForm().AddInputField("Filters", "", 100, nil, nil),
	}

	ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			go IssueUI.GetList()
		}
		return event
	})
	IssueFilterUI = ui
}

func (ui *FilterUI) GetInputField() *tview.InputField {
	return ui.GetFormItem(0).(*tview.InputField)
}

func (ui *FilterUI) SetQuery(query string) {
	ui.GetInputField().SetText(query)
}

func (ui *FilterUI) GetQuery() string {
	return ui.GetInputField().GetText()
}

func (ui *FilterUI) focus() {
}

func (ui *FilterUI) blur() {
}
