pipeline {
    agent any
    parameters {
        string(name: 'ProjectKey')
        string(name: 'RepoName')
        string(name: 'RefID')
        string(name: 'PullRequestID')
    }
    options {
        disableConcurrentBuilds()
    }
    stages {
        stage('Build') {
            steps {
                sh "echo git clone http://localhost:7990/scm/${ProjectKey}/${RepoName}.git"
                sh "echo git checkout ${RefID}"
                sh "exit 0"
            }
        }
        stage('Approve PR') {
            steps {
                sh "echo Approve PR"
                sh """
                curl -v -u user:user -XPUT -H "Content-Type: application/json" \
                    --data '{"status": "APPROVED"}' \
                    http://172.29.80.1:7990/rest/api/1.0/projects/${ProjectKey}/repos/${RepoName}/pull-requests/${PullRequestID}/participants/user
                """
            }
        }
    }
    post { 
        failure { 
            sh "echo Decline PR"
            sh """
            curl -v -u user:user  -H "Content-Type: application/json" \
                --data '{"text": "Build failed. Please test your branch locally."}' \
                http://172.29.80.1:7990/rest/api/1.0/projects/${ProjectKey}/repos/${RepoName}/pull-requests/${PullRequestID}/comments
             """
        }
    }
}
