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

type viewUI struct {
	*tview.TextView
}

func NewViewUI(uiKind UIKind) {
	ui := &viewUI{
		TextView: tview.NewTextView(),
	}

	ui.SetBorder(true).SetTitle(string(uiKind)).SetTitleAlign(tview.AlignLeft)
	ui.SetDynamicColors(true).SetWordWrap(false)

	ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	switch uiKind {
	case UIKindIssueView:
		IssueViewUI = ui
	case UIKindCommentView:
		CommentViewUI = ui
	}
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
