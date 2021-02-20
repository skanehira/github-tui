package ui

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/skanehira/ght/utils"
)

var (
	IssueViewUI   *viewUI
	CommentViewUI *viewUI
)

type viewUI struct {
	*tview.TextView
	regionIndex  int
	regionLength int
	regionIDs    []string
}

func NewViewUI(uiKind UIKind) {
	ui := &viewUI{
		TextView: tview.NewTextView(),
	}

	ui.SetBorder(true).SetTitle(string(uiKind)).SetTitleAlign(tview.AlignLeft)
	ui.SetDynamicColors(true).SetWordWrap(false).SetRegions(true)
	ui.SetBorderPadding(1, 1, 1, 1)

	searchFunc := func(input string) {
		text := ui.GetText(true)
		if input != "" {
			ui.regionIDs, text = utils.Replace(text, input, `[#ff0000]["%d"]`+input+`[""][white]`, -1)
			ui.regionLength = len(ui.regionIDs)
			if ui.regionLength > 0 {
				ui.regionIndex = 0
				ui.Highlight(ui.regionIDs[0]).ScrollToHighlight()
			}
		}

		go ui.updateView(text)
	}

	ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case '/':
			SearchUI.SetText("")
			SearchUI.SetSerachFunc(searchFunc)
			SearchUI.SetFocusFunc(func() {
				UI.app.SetFocus(ui)
			})
			UI.app.SetFocus(SearchUI)
		case 'n':
			if ui.regionLength > 0 {
				ui.regionIndex = (ui.regionIndex + 1) % ui.regionLength
				ui.Highlight(strconv.Itoa(ui.regionIndex)).ScrollToHighlight()
			}
		case 'N':
			if ui.regionLength > 0 {
				ui.regionIndex = (ui.regionIndex - 1 + ui.regionLength) % ui.regionLength
				ui.Highlight(strconv.Itoa(ui.regionIndex)).ScrollToHighlight()
			}
		}
		switch event.Key() {
		case tcell.KeyCtrlO:
			// TODO toggle full screen
		}
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
		ui.SetText(text).ScrollToBeginning()
		//out, err := glamour.Render(text, "dark")
		//if err != nil {
		//	out = err.Error()
		//}
		//ui.SetText(tview.TranslateANSI(out)).ScrollToBeginning()
	}
}

func (v *viewUI) focus() {}

func (v *viewUI) blur() {}
