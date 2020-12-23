package bitbucket

type Project struct {
	Key string
}

type Repository struct {
	Name    string
	Project Project
}

type FromRef struct {
	ID         string
	Repository Repository
}

type PullRequest struct {
	ID      int
	Version int
	FromRef FromRef
}

type Change struct {
	RefID string
}

type BitbucketRequestBody struct {
	EventKey    string
	Repository  Repository
	PullRequest PullRequest
	Changes     []Change
}

type BitbucketPullRequestListResponse struct {
	Size       int
	IsLastPage bool
	Values     []PullRequest
}

type PullRequestInfo struct {
	ProjectKey    string
	RepoName      string
	RefID         string
	PullRequestID int
}
