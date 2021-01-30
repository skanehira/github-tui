package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/skanehira/ght/domain"
	"github.com/skanehira/ght/utils"
)

var CommentUI *SelectUI

func NewCommentUI() {
	setOpt := func(ui *SelectUI) {
		ui.capture = func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyCtrlO:
				for _, s := range CommentUI.selected {
					comment := s.(*domain.Comment)
					if err := utils.OpenBrowser(comment.URL); err != nil {
						log.Println(err)
					}
				}
				if len(CommentUI.selected) == 0 {
					data := CommentUI.GetSelect()
					if data != nil {
						if err := utils.OpenBrowser(data.(*domain.Comment).URL); err != nil {
							log.Println(err)
						}
					}
				}
			}
			return event
		}

		ui.header = []string{
			"",
			"Author",
			"UpdatedAt",
		}
		ui.hasHeader = len(ui.header) > 0
	}

	ui := NewSelectListUI(UIKindComment, tcell.ColorYellow, setOpt)

	ui.SetSelectionChangedFunc(func(row, col int) {
		if row > 0 {
			CommentViewUI.updateView(ui.items[row-1].(*domain.Comment).Body)
		}
	})

	CommentUI = ui
}
