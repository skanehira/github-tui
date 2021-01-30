package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var FilterUI *filterUI

var (
	filterQuery string
)

type filterUI struct {
	*tview.Form
}

func NewFilterUI() {
	ui := &filterUI{
		Form: tview.NewForm().AddInputField("Filters", filterQuery, 100, nil, nil),
	}

	ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			filterQuery = ui.GetFormItem(0).(*tview.InputField).GetText()
			go IssueUI.GetList()
		}
		return event
	})
	FilterUI = ui
}

func (ui *filterUI) focus() {
}

func (ui *filterUI) blur() {
}
