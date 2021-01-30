package domain

import "github.com/gdamore/tcell/v2"

type Milestone struct {
	ID          string
	Title       string
	State       string
	Description string
	URL         string
}

func (m *Milestone) Key() string {
	return m.Title
}

func (m *Milestone) Fields() []Field {
	return []Field{
		{Text: m.Title, Color: tcell.ColorGreen},
	}
}
