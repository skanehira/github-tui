package ui

import (
	"github.com/charmbracelet/glamour"
	"github.com/rivo/tview"
)

type ViewUI struct {
	*tview.TextView
}

func NewViewUI() (*ViewUI, func(text string)) {
	ui := &ViewUI{
		TextView: tview.NewTextView(),
	}

	ui.SetBorder(true).SetTitle("view").SetTitleAlign(tview.AlignLeft)
	ui.SetDynamicColors(true)

	viewUpdater := func(text string) {
		out, err := glamour.Render(text, "dark")
		if err != nil {
			out = err.Error()
		}
		ui.SetText(tview.TranslateANSI(out)).ScrollToBeginning()
	}

	return ui, viewUpdater
}
