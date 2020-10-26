package main

import "github.com/shurcooL/githubv4"

type RepositoryOwner struct {
	Repositories struct {
		Nodes []struct {
			NameWithOwner githubv4.String
		}
		PageInfo struct {
			EndCursor   githubv4.String
			HasNextPage githubv4.Boolean
		}
	} `graphql:"repositories(first: $first, after: $cursor, orderBy: {field: CREATED_AT, direction: DESC})"`
}
