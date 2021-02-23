package github

import "github.com/shurcooL/githubv4"

type MutateDeleteComment struct {
	DeleteIssueComment struct {
		ClientMutationId githubv4.String
	} `graphql:"deleteIssueComment(input: $input)"`
}

type MutateUpdateIssueComment struct {
	UpdateIssueComment struct {
		ClientMutationId githubv4.String
	} `graphql:"updateIssueComment(input: $input)"`
}
