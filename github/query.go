package github

import "github.com/shurcooL/githubv4"

type PageInfo struct {
	EndCursor   githubv4.String
	HasNextPage githubv4.Boolean
}

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
	PageInfo PageInfo
}

type Issue struct {
	Number githubv4.Int
	Body   githubv4.String
	State  githubv4.String
	Author struct {
		Login githubv4.String
	}
	Title     githubv4.String
	URL       githubv4.URI
	Labels    Labels `graphql:"labels(first: 10)"`
	Assignees struct {
		Nodes []struct {
			Login githubv4.String
		}
	} `graphql:"assignees(first: 10)"`
	ProjectCards struct {
		Nodes []struct {
			Project struct {
				Name githubv4.String
			}
		}
	} `graphql:"projectCards(first: 10)"`
	Milestone struct {
		ID    githubv4.String
		Title githubv4.String
	}
	Comments struct {
		Nodes []struct {
			ID     githubv4.String
			Author struct {
				Login githubv4.String
			}
			UpdatedAt githubv4.DateTime
			BodyText  githubv4.String
			URL       githubv4.URI
		}
	} `graphql:"comments(first: 100)"`
}

type Issues struct {
	Nodes []struct {
		Issue Issue `graphql:"... on Issue"`
	}
	PageInfo PageInfo
}

type AssignableUsers struct {
	Nodes []struct {
		Login githubv4.String
	}
	PageInfo PageInfo
}

type Labels struct {
	Nodes []struct {
		Name        githubv4.String
		Description githubv4.String
		Color       githubv4.String
	}
	PageInfo PageInfo
}

type Milestones struct {
	Nodes []struct {
		Title       githubv4.String
		State       githubv4.String
		Description githubv4.String
	}
	PageInfo PageInfo
}

type Projects struct {
	Nodes []struct {
		Name githubv4.String
	}
	PageInfo PageInfo
}
