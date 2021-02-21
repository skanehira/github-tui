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
		getSelectedComments := func() []*domain.Comment {
			var comments []*domain.Comment
			if len(CommentUI.selected) == 0 {
				data := CommentUI.GetSelect()
				comments = append(comments, data.(*domain.Comment))
			} else {
				for _, item := range CommentUI.selected {
					comments = append(comments, item.(*domain.Comment))
				}
			}
			return comments
		}

		ui.capture = func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyCtrlO:
				for _, comment := range getSelectedComments() {
					if err := utils.Open(comment.URL); err != nil {
						log.Println(err)
					}
				}
				CommentUI.ClearSelected()
				CommentUI.UpdateView()
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
