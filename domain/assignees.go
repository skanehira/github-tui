package domain

import "github.com/gdamore/tcell/v2"

type AssignableUser struct {
	Login string
}

func (a *AssignableUser) Key() string {
	return a.Login
}

func (a *AssignableUser) Fields() []Field {
	return []Field{
		{Text: a.Login, Color: tcell.ColorFuchsia},
	}
}
