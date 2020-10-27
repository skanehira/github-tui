package github

import "github.com/shurcooL/githubv4"

type RepositoryOwner struct {
	Repositories struct {
		Nodes []struct {
			NameWithOwner    githubv4.String
			CreatedAt        githubv4.DateTime
			DefaultBranchRef struct {
				Name githubv4.String
			}
			Description githubv4.String
			LicenseInfo struct {
				Name githubv4.String
			}
			StargazerCount githubv4.Int
			URL            githubv4.URI
			SSHURL         githubv4.String
		}
		PageInfo struct {
			EndCursor   githubv4.String
			HasNextPage githubv4.Boolean
		}
	} `graphql:"repositories(first: $first, after: $cursor, orderBy: {field: CREATED_AT, direction: DESC})"`
}
