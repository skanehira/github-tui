package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/skanehira/ght/domain"
	"github.com/skanehira/ght/github"
)

const (
	unselected = "\u25ef"
	selected   = "\u25c9"
)

type UIKind string

const (
	UIKindIssue       UIKind = "issues"
	UIKindAssignee           = "assignees"
	UIKindComment            = "comments"
	UIKindLabel              = "labels"
	UIKindMilestones         = "milestones"
	UIKindProject            = "projects"
	UIKindIssueView          = "issue preview"
	UIKindCommentView        = "comment preview"
	UIKindCommonView         = "preview"
)

type (
	SetSelectUIOpt func(ui *SelectUI)
	GetListFunc    func(cursor *string) ([]domain.Item, *github.PageInfo)
	CaptureFunc    func(event *tcell.EventKey) *tcell.EventKey
)

type SelectUI struct {
	uiKind      UIKind
	cursor      *string
	hasNext     bool
	getList     GetListFunc
	capture     CaptureFunc
	header      []string
	hasHeader   bool
	originItems []domain.Item
	items       []domain.Item
	selected    map[string]domain.Item
	boxColor    tcell.Color
	searchWord  string
	*tview.Table
}

func NewSelectListUI(uiKind UIKind, boxColor tcell.Color, setOpt SetSelectUIOpt) *SelectUI {
	ui := &SelectUI{
		uiKind:   uiKind,
		hasNext:  true,
		selected: make(map[string]domain.Item),
		boxColor: boxColor,
		Table:    tview.NewTable().SetSelectable(false, false),
	}

	ui.SetBorder(true).SetTitle(string(uiKind)).SetTitleAlign(tview.AlignLeft)
	ui.SetBorderColor(boxColor)

	setOpt(ui)

	go ui.Init()
	return ui
}

func (ui *SelectUI) GetList() {
	if ui.getList != nil {
		list, pageInfo := ui.getList(nil)
		if pageInfo != nil {
			ui.hasNext = bool(pageInfo.HasNextPage)
			cursor := string(pageInfo.EndCursor)
			ui.originItems = list
			ui.cursor = &cursor
			ui.Select(0, 0)
			ui.UpdateView()
		}
	}
}

func (ui *SelectUI) SetList(list []domain.Item) {
	ui.originItems = list
	ui.selected = make(map[string]domain.Item)
	ui.Select(0, 0)
	ui.UpdateView()
}

func (ui *SelectUI) FetchList() {
	if ui.hasNext && ui.getList != nil {
		list, pageInfo := ui.getList(ui.cursor)
		ui.hasNext = bool(pageInfo.HasNextPage)
		cursor := string(pageInfo.EndCursor)
		ui.originItems = append(ui.originItems, list...)
		ui.cursor = &cursor
		ui.UpdateView()
	}
}

func (ui *SelectUI) UpdateView() {
	UI.updater <- func() {
		ui.Clear()
		for i, h := range ui.header {
			ui.SetCell(0, i, &tview.TableCell{
				Text:            h,
				NotSelectable:   true,
				Align:           tview.AlignLeft,
				Color:           tcell.ColorWhite,
				BackgroundColor: tcell.ColorDefault,
				Attributes:      tcell.AttrBold | tcell.AttrUnderline,
			})
		}

		if len(ui.originItems) < 1 {
			return
		}

		h := 0
		if ui.hasHeader {
			h++
			ui.SetFixed(1, 0)
		}

		selectColor := ui.originItems[0].Fields()[0].Color

		ui.items = []domain.Item{}
		if ui.searchWord != "" {
			for _, data := range ui.originItems {
				for _, f := range data.Fields() {
					if strings.Contains(strings.ToLower(f.Text), strings.ToLower(ui.searchWord)) {
						ui.items = append(ui.items, data)
						break
					}
				}
			}
		} else {
			ui.items = ui.originItems
		}

		for i, data := range ui.items {
			if _, ok := ui.selected[data.Key()]; ok {
				ui.SetCell(i+h, 0, tview.NewTableCell(selected).SetTextColor(selectColor))
			} else {
				ui.SetCell(i+h, 0, tview.NewTableCell(unselected).SetTextColor(selectColor))
			}
			for j, f := range data.Fields() {
				ui.SetCell(i+h, j+1, tview.NewTableCell(f.Text).SetTextColor(f.Color))
			}
		}
		ui.ScrollToBeginning()

		// when update filter, then update ui related issue primitives
		if ui.uiKind == UIKindIssue {
			row, _ := ui.GetSelection()
			if row == 0 {
				row = 1
			}
			updateUIRelatedIssue(ui, row)
		}
	}
}

func (ui *SelectUI) Init() {
	ui.GetList()

	searchFunc := func(text string) {
		ui.searchWord = text
		ui.UpdateView()
	}

	ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlJ:
			row, col := ui.GetSelection()
			max := len(ui.items)
			if ui.hasHeader {
				max++
			}
			if row < max {
				ui.toggleSelected(row)
			}

			if row+1 < max {
				ui.Select(row+1, col)
			}
		case tcell.KeyCtrlK:
			row, col := ui.GetSelection()
			min := 0
			if ui.hasHeader {
				min++
			}
			if row > min {
				ui.toggleSelected(row - 1)
			}
			if row > min {
				ui.Select(row-1, col)
			}
		}

		switch event.Rune() {
		case 'f':
			go ui.FetchList()
		case '/':
			SearchUI.SetSerachFunc(searchFunc)
			SearchUI.SetFocusFunc(func() {
				UI.app.SetFocus(ui)
			})
			UI.app.SetFocus(SearchUI)
		}

		return ui.capture(event)
	})
}

func (ui *SelectUI) toggleSelected(row int) {
	var data domain.Item
	if ui.hasHeader {
		data = ui.items[row-1]
	} else {
		data = ui.items[row]
	}
	selectColor := ui.items[0].Fields()[0].Color
	if _, ok := ui.selected[data.Key()]; ok {
		delete(ui.selected, data.Key())
		ui.SetCell(row, 0, tview.NewTableCell(unselected).SetTextColor(selectColor))
	} else {
		ui.selected[data.Key()] = data
		ui.SetCell(row, 0, tview.NewTableCell(selected).SetTextColor(selectColor))
	}
}

func (ui *SelectUI) UpdateItem(item domain.Item) {
	for i, t := range ui.originItems {
		if t.Key() == item.Key() {
			ui.originItems[i] = item
		}
	}
}

func (ui *SelectUI) GetSelect() domain.Item {
	row, _ := ui.GetSelection()
	if ui.hasHeader {
		row = row - 1
	}
	if len(ui.items) > row {
		id := ui.items[row].Key()
		for _, item := range ui.originItems {
			if item.Key() == id {
				return item
			}
		}
	}
	return nil
}

func (ui *SelectUI) focus() {
	ui.SetSelectable(true, false)
}

func (ui *SelectUI) blur() {
	ui.SetSelectable(false, false)
}

func (ui *SelectUI) ClearView() {
	ui.Clear()
	ui.ClearSelected()
}

func (ui *SelectUI) ClearSelected() {
	ui.selected = make(map[string]domain.Item)
}
