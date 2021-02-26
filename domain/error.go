package domain

import "errors"

var (
	ErrCommentBodyIsEmpty      = errors.New("comment body is empty")
	ErrCommentBodyIsNotChanged = errors.New("issu body is not changed")
	ErrNotFoundComment         = errors.New("not found comment")
	ErrNotFoundIssue           = errors.New("not found issue")
)
