package github

import "github.com/shurcooL/githubv4"

type PageInfo struct {
	EndCursor   githubv4.String
	HasNextPage githubv4.Boolean
}
