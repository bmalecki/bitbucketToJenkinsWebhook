package bitbucket

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func assert(expect, actual interface{}) {
	if expect != actual {
		panic(fmt.Errorf("Expect: %v, but was: %v", expect, actual))
	}
}

type mockHTTPClient struct {
	t *testing.T
}

func (c *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	assert("http://testurl.exaple.com:7990/rest/api/1.0/projects/testProjectKey/repos/testRepoName/pull-requests?state=OPEN",
		req.URL.String())

	data, err := os.Open("../test/list_pr.json")
	if err != nil {
		c.t.Errorf("Error %v. ", err)
	}

	return &http.Response{
		Body: data,
	}, nil
}

func TestGetAllOpenPullRequests(t *testing.T) {
	httpClient = &mockHTTPClient{t}

	var bitbucketClient = NewBitbucketClient(BitbucketClientConfig{
		User:     "testUser",
		Password: "testUserPass",
		URL:      "http://testurl.exaple.com:7990",
	})

	prList := bitbucketClient.GetAllOpenPullRequests("testProjectKey", "testRepoName")
	assert(2, prList.Size)
	assert(3, prList.Values[0].ID)
	assert("refs/heads/testPR2", prList.Values[0].FromRef.ID)
	assert("example", prList.Values[0].FromRef.Repository.Name)
	assert("TEST", prList.Values[0].FromRef.Repository.Project.Key)
	assert(2, prList.Values[1].ID)
	assert("refs/heads/PR", prList.Values[1].FromRef.ID)
	assert("example", prList.Values[1].FromRef.Repository.Name)
	assert("TEST", prList.Values[1].FromRef.Repository.Project.Key)
}
