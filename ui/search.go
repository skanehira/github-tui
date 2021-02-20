package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var SearchUI *searchUI

type SearchFunc func(text string)

type searchUI struct {
	*tview.InputField
	SearchFunc SearchFunc
	FocusFunc  func()
}

func NewSearchUI() {
	ui := &searchUI{
		InputField: tview.NewInputField(),
		SearchFunc: func(text string) {},
		FocusFunc:  func() {},
	}

	ui.SetDoneFunc(func(key tcell.Key) {
		ui.FocusFunc()
	})

	SearchUI = ui
}

func (s *searchUI) SetSerachFunc(f SearchFunc) {
	s.SetChangedFunc(f)
}

func (s *searchUI) SetFocusFunc(f func()) {
	s.FocusFunc = f
}

func (s *searchUI) focus() {

}

func (s *searchUI) blur() {

}
