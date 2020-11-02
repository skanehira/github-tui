package ui

import (
	"github.com/charmbracelet/glamour"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var ViewUI *viewUI

type viewUI struct {
	*tview.TextView
}

func NewViewUI() *viewUI {
	ui := &viewUI{
		TextView: tview.NewTextView(),
	}

	ui.SetBorder(true).SetTitle("view").SetTitleAlign(tview.AlignLeft)
	ui.SetDynamicColors(true).SetWordWrap(false)

	ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return UI.Capture(event)
	})

	ViewUI = ui
	return ui
}

func viewUpdater(text string) {
	UI.updater <- func() {
		out, err := glamour.Render(text, "dark")
		if err != nil {
			out = err.Error()
		}
		ViewUI.SetText(tview.TranslateANSI(out)).ScrollToBeginning()
	}
}

func (v *viewUI) focus() {}

func (v *viewUI) blur() {}
