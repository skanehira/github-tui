package domain

import "github.com/gdamore/tcell/v2"

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
