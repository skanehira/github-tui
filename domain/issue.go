package domain

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type Issue struct {
	ID        string
	Repo      string
	RepoOwner string
	Number    string
	State     string
	Title     string
	Body      string
	Author    string
	URL       string
	Labels    []Item
	Assignees []Item
	Comments  []Item
	MileStone []Item
	Projects  []Item
}

func (i *Issue) Key() string {
	return i.ID
}

func (i *Issue) Fields() []Field {
	stateColor := tcell.ColorGreen
	if i.State == "CLOSED" {
		stateColor = tcell.ColorRed
	}

	f := []Field{
		{Text: fmt.Sprintf("%s/%s", i.RepoOwner, i.Repo), Color: tcell.ColorLightSalmon},
		{Text: i.Number, Color: tcell.ColorBlue},
		{Text: i.State, Color: stateColor},
		{Text: i.Author, Color: tcell.ColorYellow},
		{Text: i.Title, Color: tcell.ColorWhite},
	}

	return f
}
