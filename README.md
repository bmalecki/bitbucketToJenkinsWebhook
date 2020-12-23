
Jenkins

1 Trigger builds remotely (e.g., from scripts)
use token from:
Dashboard -> <User> -> Configure -> API Token

curl -u admin:11cef4962f699ef4e4d4f9093e63445a2e -XPOST http://localhost:8080/job/test/build




/home/bartek/projects/poc

Middleware:
curl http://172.29.81.89:8090/


docker run -v bitbucketVolume:/var/atlassian/application-data/bitbucket --name="bitbucket" --rm -p 7990:7990 -p 7999:7999 atlassian/bitbucket-server:6.10

docker run -p 8080:8080 -p 50000:50000 -v jenkins_home:/var/jenkins_home jenkins/jenkins:lts

2. Add Webhook on bitbucket

- Go to repository settings
- Select 'Webhooks'
- Click 'Create Webhook'


List PR:
https://docs.atlassian.com/bitbucket-server/rest/6.10.0/bitbucket-rest.html



curl -u user:user  http://localhost:7990/rest/api/1.0/projects/TEST/repos/example/pull-requests?state=OPEN

curl -v -u user:user  -H "Content-Type: application/json" \
--data '{"text": "An insightful general comment on a pull request."}' \
http://localhost:7990/rest/api/1.0/projects/TEST/repos/example/pull-requests/2/comments

curl -v -u user:user -XPUT -H "Content-Type: application/json" \
--data '{"status": "NEEDS_WORK"}' \
http://localhost:7990/rest/api/1.0/projects/TEST/repos/Example/pull-requests/2/participants/user

curl -v -u user:user -XPUT -H "Content-Type: application/json" \
--data '{"status": "APPROVED"}' \
http://localhost:7990/rest/api/1.0/projects/TEST/repos/Example/pull-requests/2/participants/user

curl -u user:user -XPOST -H "Content-Type: application/json" \
http://localhost:7990/rest/api/1.0/projects/TEST/repos/example/pull-requests/2/approve

curl -u user:user -XPOST -H "Content-Type: application/json" \
http://localhost:7990/rest/api/1.0/projects/TEST/repos/example/pull-requests/2/decline


Test go webhook:

curl -i -X POST host:port/post-file \
  -H "Content-Type: text/xml" \
  --data-binary "@path/to/file"
