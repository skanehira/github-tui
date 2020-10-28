package github

import "github.com/shurcooL/githubv4"

type Repositories struct {
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
}

type Issues struct {
	Nodes []struct {
		Number githubv4.Int
		Body   githubv4.String
		State  githubv4.String
		Author struct {
			Login githubv4.String
		}
		Title     githubv4.String
		URL       githubv4.URI
		Labels    Labels `graphql:"labels(first: $first)"`
		Assignees struct {
			Nodes []struct {
				Login githubv4.String
			}
		} `graphql:"assignees(first: $first)"`
	}
	PageInfo struct {
		EndCursor   githubv4.String
		HasNextPage githubv4.Boolean
	}
}

type Labels struct {
	Nodes []struct {
		Name        githubv4.String
		Description githubv4.String
		Color       githubv4.String
	}
}

type Milestones struct {
	Nodes []struct {
		Title       githubv4.String
		State       githubv4.String
		Description githubv4.String
	}
}

type Projects struct {
	Nodes []struct {
		Name githubv4.String
	}
}
