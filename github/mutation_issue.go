package github

type MutateOpenIsseue struct {
	ReopenIssue struct {
		Issue struct {
			ID string
		}
	} `graphql:"reopenIssue(input: $input)"`
}

type MutateCoseIssue struct {
	CloseIssue struct {
		Issue struct {
			ID string
		}
	} `graphql:"closeIssue(input: $input)"`
}
