package github

type MutateOpenIsseue struct {
	ReopenIssue struct {
		Issue struct {
			ID string
		}
	} `graphql:"reopenIssue(input: $input)"`
}
