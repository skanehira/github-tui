package github

import "github.com/shurcooL/githubv4"

type Repository struct {
	ID               githubv4.ID
	Name             githubv4.String
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

type Repositories struct {
	Nodes    []Repository
	PageInfo PageInfo
}
