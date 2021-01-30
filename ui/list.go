package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/skanehira/ght/github"
)

const (
	unselected = "\u25ef"
	selected   = "\u25c9"
)

type Item interface {
	Key() string
	Fields() []Field
}

type Field struct {
	Text  string
	Color tcell.Color
}

type (
	SetSelectUIOpt func(ui *SelectUI)
	GetListFunc    func(cursor *string) ([]Item, *github.PageInfo)
	CaptureFunc    func(event *tcell.EventKey) *tcell.EventKey
	InitFunc       func(ui *SelectUI)
)

type SelectUI struct {
	cursor    *string
	hasNext   bool
	getList   GetListFunc
	capture   CaptureFunc
	init      InitFunc
	header    []string
	hasHeader bool
	items     []Item
	selected  map[string]Item
	boxColor  tcell.Color
	*tview.Table
}

func NewSelectListUI(title string, boxColor tcell.Color, setOpt SetSelectUIOpt) *SelectUI {
	ui := &SelectUI{
		hasNext:  true,
		selected: make(map[string]Item),
		boxColor: boxColor,
		Table:    tview.NewTable().SetSelectable(false, false),
	}

	ui.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
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
			ui.items = list
			ui.cursor = &cursor
			ui.Select(0, 0)
			ui.UpdateView()
		}
	}
}

func (ui *SelectUI) SetList(list []Item) {
	ui.items = list
	ui.selected = make(map[string]Item)
	ui.Select(0, 0)
	ui.UpdateView()
}

func (ui *SelectUI) FetchList() {
	if ui.hasNext {
		list, pageInfo := ui.getList(ui.cursor)
		ui.hasNext = bool(pageInfo.HasNextPage)
		cursor := string(pageInfo.EndCursor)
		ui.items = append(ui.items, list...)
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

		h := 0
		if ui.hasHeader {
			h++
		}
		for i, data := range ui.items {
			if _, ok := ui.selected[data.Key()]; ok {
				ui.SetCell(i+h, 0, tview.NewTableCell(selected).SetTextColor(ui.boxColor))
			} else {
				ui.SetCell(i+h, 0, tview.NewTableCell(unselected).SetTextColor(ui.boxColor))
			}
			for j, f := range data.Fields() {
				ui.SetCell(i+h, j+1, tview.NewTableCell(f.Text).SetTextColor(f.Color))
			}
		}
		ui.ScrollToBeginning()
	}
}

func (ui *SelectUI) Init() {
	ui.GetList()
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
		case 'G':
			go ui.FetchList()
		}

		return ui.capture(event)
	})

	if ui.init != nil {
		ui.init(ui)
	}
}

func (ui *SelectUI) toggleSelected(row int) {
	var data Item
	if ui.hasHeader {
		data = ui.items[row-1]
	} else {
		data = ui.items[row]
	}
	if _, ok := ui.selected[data.Key()]; ok {
		delete(ui.selected, data.Key())
		ui.SetCell(row, 0, tview.NewTableCell(unselected).SetTextColor(ui.boxColor))
	} else {
		ui.selected[data.Key()] = data
		ui.SetCell(row, 0, tview.NewTableCell(selected).SetTextColor(ui.boxColor))
	}
}

func (ui *SelectUI) GetSelect() Item {
	row, _ := ui.GetSelection()
	if ui.hasHeader {
		row = row - 1
	}
	if len(ui.items) > row {
		return ui.items[row]
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
	ui.selected = make(map[string]Item)
}
