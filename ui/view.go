package ui

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/skanehira/ght/utils"
)

var (
	IssueViewUI   *ViewUI
	CommentViewUI *ViewUI
	CommonViewUI  *ViewUI
)

type ViewUI struct {
	*tview.TextView
	regionIndex  int
	regionLength int
	regionIDs    []string
	uiKind       UIKind
	setFocus     func()
}

func NewViewUI(uiKind UIKind) {
	ui := &ViewUI{
		TextView: tview.NewTextView(),
		uiKind:   uiKind,
	}

	ui.SetBorder(true).SetTitle(string(uiKind)).SetTitleAlign(tview.AlignLeft)
	ui.SetDynamicColors(true).SetWordWrap(false).SetRegions(true)
	ui.SetBorderPadding(1, 1, 1, 1)

	var setFocus func()

	switch uiKind {
	case UIKindIssueView:
		IssueViewUI = ui
		setFocus = func() {
			UI.app.SetFocus(IssueViewUI)
		}
	case UIKindCommentView:
		CommentViewUI = ui
		setFocus = func() {
			UI.app.SetFocus(CommentViewUI)
		}
	case UIKindCommonView:
		CommonViewUI = ui
	}

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
		case 'o':
			if ui.uiKind == UIKindCommonView {
				UI.pages.SwitchToPage("main")
				ui.setFocus()
				return event
			}
			UI.FullScreenPreview(ui.GetText(true), setFocus)
		}

		//switch event.Key() {
		//}
		return event
	})

}

func (ui *ViewUI) updateView(text string) {
	UI.updater <- func() {
		ui.SetText(text).ScrollToBeginning()
		//out, err := glamour.Render(text, "dark")
		//if err != nil {
		//	out = err.Error()
		//}
		//ui.SetText(tview.TranslateANSI(out)).ScrollToBeginning()
	}
}

func (v *ViewUI) focus() {}

func (v *ViewUI) blur() {}
