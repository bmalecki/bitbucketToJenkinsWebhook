package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	b "github.com/bmalecki/bitbucketToJenkinsWebhook/bitbucket"
	j "github.com/bmalecki/bitbucketToJenkinsWebhook/jenkins"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var bitbucketClient = b.NewBitbucketClient(b.BitbucketClientConfig{
	User:     getEnv("BITBUCKET_USER", "user"),
	Password: getEnv("BITBUCKET_PASSWORD", "user"),
	URL:      getEnv("BITBUCKET_URL", "http://localhost:7990"),
})

var jenkinsClient = j.NewJenkinsClient(j.JenkinsClientConfig{
	User:  getEnv("JENKINS_USER", "admin"),
	Token: getEnv("JENKINS_TOKEN", "11cef4962f699ef4e4d4f9093e63445a2e"),
	URL:   getEnv("JENKINS_URL", "http://localhost:8080/job/test/buildWithParameters"),
})

func notFound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	logger.Printf("%s %s\n", req.Method, req.URL.String())
	fmt.Fprintf(w, "")
}

func findPullRequestForRefsChanged(prList *b.BitbucketPullRequestListResponse,
	body *b.BitbucketRequestBody) []*b.PullRequestInfo {

	var s []*b.PullRequestInfo

	for _, pr := range prList.Values {
		if pr.FromRef.ID == body.Changes[0].RefID {
			s = append(s, &b.PullRequestInfo{
				ProjectKey:    body.Repository.Project.Key,
				RepoName:      body.Repository.Name,
				RefID:         body.Changes[0].RefID,
				PullRequestID: pr.ID,
			})
		}
	}

	return s
}

func parseBitbucketWebhookRequest(r io.Reader) (*b.BitbucketRequestBody, error) {
	bodyBuffer, _ := ioutil.ReadAll(r)
	body := &b.BitbucketRequestBody{}
	if err := json.Unmarshal(bodyBuffer, body); err != nil {
		return nil, err
	}
	return body, nil
}

func pullRequestOpened(body *b.BitbucketRequestBody) {
	prInfo := &b.PullRequestInfo{
		ProjectKey:    body.PullRequest.FromRef.Repository.Project.Key,
		RepoName:      body.PullRequest.FromRef.Repository.Name,
		RefID:         body.PullRequest.FromRef.ID,
		PullRequestID: body.PullRequest.ID,
	}
	logger.Printf("Opened PR. %v", prInfo)
	bitbucketClient.SetNeedsWorkStatus(prInfo)
	bitbucketClient.AddComment(prInfo)
	jenkinsClient.RunJenkinsJob(prInfo)
}

func refsChanged(body *b.BitbucketRequestBody) {
	logger.Printf("Refs in repo changed.")
	prList := bitbucketClient.GetAllOpenPullRequests(body.Repository.Project.Key, body.Repository.Name)
	prInfos := findPullRequestForRefsChanged(prList, body)
	for _, prInfo := range prInfos {
		logger.Printf("Updated open PR: %v", prInfo)
		bitbucketClient.SetNeedsWorkStatus(prInfo)
		bitbucketClient.AddComment(prInfo)
		jenkinsClient.RunJenkinsJob(prInfo)
	}
}

func webhook(w http.ResponseWriter, req *http.Request) {
	logger.Printf("%s %s\n", req.Method, req.URL.String())

	/*
		Securing your webhook
		https://confluence.atlassian.com/bitbucketserver0610/managing-webhooks-in-bitbucket-server-989761415.html?utm_campaign=in-app-help&utm_medium=in-app-help&utm_source=stash#ManagingwebhooksinBitbucketServer-webhooksecrets
	*/
	// logger.Println(req.Header["X-Hub-Signature"][0])

	body, err := parseBitbucketWebhookRequest(req.Body)
	if err != nil {
		panic(err)
	}

	switch body.EventKey {
	case "pr:opened":
		pullRequestOpened(body)
	case "repo:refs_changed":
		refsChanged(body)
	}

	fmt.Fprintf(w, "")
}

func main() {
	logger.Println("Start serving...")

	http.HandleFunc("/", notFound)
	http.HandleFunc("/bitbucketWebhook", webhook)
	http.ListenAndServe(":8090", nil)
}
