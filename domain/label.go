package domain

import "github.com/gdamore/tcell/v2"

type Label struct {
	Name        string
	Description string
}

func (l *Label) Key() string {
	return l.Name
}

func (l *Label) Fields() []Field {
	return []Field{
		{Text: l.Name, Color: tcell.ColorLightYellow},
	}
}
