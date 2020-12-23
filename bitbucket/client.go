package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/bmalecki/bitbucketToJenkinsWebhook/common"
)

type BitbucketClient interface {
	SetNeedsWorkStatus(prInfo *PullRequestInfo)
	AddComment(prInfo *PullRequestInfo)
	GetAllOpenPullRequests(projectKey, repoName string) *BitbucketPullRequestListResponse
}

type BitbucketClientConfig struct {
	URL      string
	User     string
	Password string
}

type bitbucketClientImpl struct {
	config BitbucketClientConfig
}

func NewBitbucketClient(config BitbucketClientConfig) BitbucketClient {
	return &bitbucketClientImpl{config}
}

var httpClient common.HTTPClient = &http.Client{}
var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

func (c *bitbucketClientImpl) SetNeedsWorkStatus(prInfo *PullRequestInfo) {
	prURL := fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/pull-requests/%d/participants/%s",
		c.config.URL, prInfo.ProjectKey, prInfo.RepoName, prInfo.PullRequestID, c.config.User)
	needsWorkPayload, _ := json.Marshal(struct {
		Status string `json:"status"`
	}{"NEEDS_WORK"})

	request, _ := http.NewRequest(http.MethodPut, prURL, bytes.NewBuffer(needsWorkPayload))
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(c.config.User, c.config.Password)
	resp, err := httpClient.Do(request)

	if err != nil {
		panic(err)
	}

	logger.Printf("Need works status code: %v", resp.StatusCode)

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("Status Code from Bitbucket: %d", resp.StatusCode))
	}
}

func (c *bitbucketClientImpl) AddComment(prInfo *PullRequestInfo) {
	commentURL := fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/pull-requests/%d/comments",
		c.config.URL, prInfo.ProjectKey, prInfo.RepoName, prInfo.PullRequestID)
	commentPayload, _ := json.Marshal(struct {
		Text string `json:"text"`
	}{"Waiting for completion of build and functional tests"})

	request, _ := http.NewRequest(http.MethodPost, commentURL, bytes.NewBuffer(commentPayload))
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(c.config.User, c.config.Password)
	resp, err := httpClient.Do(request)

	if err != nil {
		panic(err)
	}

	logger.Printf("Add comment status code: %v", resp.StatusCode)

	if resp.StatusCode != 201 {
		panic(fmt.Errorf("Status Code from Bitbucket: %d", resp.StatusCode))
	}
}

func (c *bitbucketClientImpl) GetAllOpenPullRequests(projectKey, repoName string) *BitbucketPullRequestListResponse {
	prListURL := fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/pull-requests?state=OPEN",
		c.config.URL, projectKey, repoName)
	request, _ := http.NewRequest(http.MethodGet, prListURL, nil)
	request.SetBasicAuth(c.config.User, c.config.Password)
	resp, err := httpClient.Do(request)

	if err != nil {
		panic(err)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	body := BitbucketPullRequestListResponse{}
	if err := json.Unmarshal(bodyBuffer, &body); err != nil {
		panic(err)
	}

	return &body
}
