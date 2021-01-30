package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/skanehira/ght/utils"
)

var CommentUI *SelectUI

type Comment struct {
	ID        string
	Author    string
	UpdatedAt string
	URL       string
	Body      string
}

func (c *Comment) Key() string {
	return c.ID
}

func (c *Comment) Fields() []Field {
	f := []Field{
		{Text: c.Author, Color: tcell.ColorYellow},
		{Text: c.UpdatedAt, Color: tcell.ColorWhite},
	}

	return f
}

func NewCommentUI() {
	setOpt := func(ui *SelectUI) {
		ui.capture = func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyCtrlO:
				for _, s := range CommentUI.selected {
					comment := s.(*Comment)
					if err := utils.OpenBrowser(comment.URL); err != nil {
						log.Println(err)
					}
				}
				if len(CommentUI.selected) == 0 {
					data := CommentUI.GetSelect()
					if data != nil {
						if err := utils.OpenBrowser(data.(*Comment).URL); err != nil {
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

	ui := NewSelectListUI("comment list", tcell.ColorYellow, setOpt)

	ui.SetSelectionChangedFunc(func(row, col int) {
		if row > 0 {
			CommentViewUI.updateView(ui.items[row-1].(*Comment).Body)
		}
	})

	CommentUI = ui
}
