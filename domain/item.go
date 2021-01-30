package domain

import "github.com/gdamore/tcell/v2"

type Item interface {
	Key() string
	Fields() []Field
}

type Field struct {
	Text  string
	Color tcell.Color
}
