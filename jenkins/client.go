package jenkins

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/bmalecki/bitbucketToJenkinsWebhook/bitbucket"
)

type JenkinsClient interface {
	RunJenkinsJob(prInfo *bitbucket.PullRequestInfo)
}

type JenkinsClientConfig struct {
	URL   string
	User  string
	Token string
}

type jenkinsClientImpl struct {
	config JenkinsClientConfig
}

func NewJenkinsClient(config JenkinsClientConfig) JenkinsClient {
	return &jenkinsClientImpl{config}
}

var httpClient = &http.Client{}
var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

func (c *jenkinsClientImpl) RunJenkinsJob(prInfo *bitbucket.PullRequestInfo) {
	request, _ := http.NewRequest(http.MethodPost, c.config.URL, nil)
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(c.config.User, c.config.Token)

	v := reflect.ValueOf(*prInfo)
	typeOfprInfo := v.Type()
	q := request.URL.Query()

	for i := 0; i < v.NumField(); i++ {
		q.Add(typeOfprInfo.Field(i).Name, fmt.Sprintf("%v", v.Field(i).Interface()))
	}

	request.URL.RawQuery = q.Encode()
	resp, err := httpClient.Do(request)

	if err != nil {
		panic(err)
	}

	logger.Printf("Jenkins status code: %v", resp.StatusCode)

	if resp.StatusCode != 201 {
		panic(fmt.Errorf("Status Code from Jenkins: %d", resp.StatusCode))
	}
}
