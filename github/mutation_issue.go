package github

import "github.com/shurcooL/githubv4"

type MutateOpenIsseue struct {
	ReopenIssue struct {
		Issue struct {
			ID githubv4.String
		}
	} `graphql:"reopenIssue(input: $input)"`
}

type MutateCoseIssue struct {
	CloseIssue struct {
		Issue struct {
			ID githubv4.String
		}
	} `graphql:"closeIssue(input: $input)"`
}

type MutateCreateIssue struct {
	CreateIssue struct {
		Issue struct {
			ID githubv4.String
		}
	} `graphql:"createIssue(input: $input)"`
}

type MutateUpdateIssue struct {
	UpdateIssue struct {
		Issue struct {
			ID githubv4.ID
		}
	} `graphql:"updateIssue(input: $input)"`
}
