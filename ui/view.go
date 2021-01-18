package ui

import (
	"github.com/charmbracelet/glamour"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	IssueViewUI   *viewUI
	CommentViewUI *viewUI
)

const (
	IssuePreview   = "issue preview"
	CommentPreview = "comment preview"
)

type viewUI struct {
	*tview.TextView
}

func NewViewUI(previewType string) *viewUI {
	ui := &viewUI{
		TextView: tview.NewTextView(),
	}

	ui.SetBorder(true).SetTitle(previewType).SetTitleAlign(tview.AlignLeft)
	ui.SetDynamicColors(true).SetWordWrap(false)

	ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return UI.Capture(event)
	})

	switch previewType {
	case IssuePreview:
		IssueViewUI = ui
	case CommentPreview:
		CommentViewUI = ui
	}
	return ui
}

func (ui *viewUI) updateView(text string) {
	UI.updater <- func() {
		out, err := glamour.Render(text, "dark")
		if err != nil {
			out = err.Error()
		}
		ui.SetText(tview.TranslateANSI(out)).ScrollToBeginning()
	}
}

func (v *viewUI) focus() {}

func (v *viewUI) blur() {}
