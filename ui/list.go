package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/skanehira/ght/github"
)

type ListData interface {
	Key() string
	Fields() []string
}

type (
	GetListFunc func(cursor *string) ([]ListData, github.PageInfo)
	CaptureFunc func(event *tcell.EventKey) *tcell.EventKey
)

type SelectListUI struct {
	cursor    *string
	hasNext   bool
	getList   GetListFunc
	capture   CaptureFunc
	list      []ListData
	selected  map[string]struct{}
	textColor tcell.Color
	updater   func(func())
	*tview.Table
}

func NewSelectListUI(title string, updater func(func()), textColor tcell.Color, getList GetListFunc, capture CaptureFunc) *SelectListUI {
	ui := &SelectListUI{
		hasNext:   true,
		getList:   getList,
		capture:   capture,
		selected:  make(map[string]struct{}),
		textColor: textColor,
		updater:   updater,
		Table:     tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(0, 0),
	}
	ui.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)

	go ui.Init()
	return ui
}

func (ui *SelectListUI) GetList() {
	if ui.hasNext {
		list, pageInfo := ui.getList(ui.cursor)
		ui.hasNext = bool(pageInfo.HasNextPage)
		cursor := string(pageInfo.EndCursor)
		ui.list = append(ui.list, list...)
		ui.cursor = &cursor
		ui.UpdateList()
	}
}

func (ui *SelectListUI) UpdateList() {
	ui.updater(func() {
		ui.Clear()
		for i, data := range ui.list {
			if _, ok := ui.selected[data.Key()]; ok {
				ui.SetCell(i, 0, tview.NewTableCell("◉").SetTextColor(ui.textColor))
			} else {
				ui.SetCell(i, 0, tview.NewTableCell("◯").SetTextColor(ui.textColor))
			}
			for j, f := range data.Fields() {
				ui.SetCell(i, j+1, tview.NewTableCell(f).SetTextColor(ui.textColor))
			}
		}
	})
}

func (ui *SelectListUI) Init() {
	ui.GetList()
	ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlJ:
			row, col := ui.GetSelection()
			leng := len(ui.list)
			if leng >= row+1 {
				ui.toggleSelected(row)
			}
			if leng > row+1 {
				ui.Select(row+1, col)
			}
		case tcell.KeyCtrlK:
			row, col := ui.GetSelection()
			if row >= 0 {
				ui.toggleSelected(row)
			}
			if row > 0 {
				ui.Select(row-1, col)
			}
		}

		switch event.Rune() {
		case 'G':
			ui.GetList()
		}

		return ui.capture(event)
	})
}

func (ui *SelectListUI) toggleSelected(row int) {
	data := ui.list[row]
	if _, ok := ui.selected[data.Key()]; ok {
		delete(ui.selected, data.Key())
		ui.SetCell(row, 0, tview.NewTableCell("◯").SetTextColor(ui.textColor))
	} else {
		ui.selected[data.Key()] = struct{}{}
		ui.SetCell(row, 0, tview.NewTableCell("◉").SetTextColor(ui.textColor))
	}
}
