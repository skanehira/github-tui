package domain

import "errors"

var (
	ErrCommentBodyIsEmpty = errors.New("comment body is empty")
	ErrNotFoundComment    = errors.New("not found comment")
	ErrNotFoundIssue      = errors.New("not found issue")
)
