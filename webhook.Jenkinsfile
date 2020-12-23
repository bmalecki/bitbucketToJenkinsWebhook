pipeline {
    agent any
    parameters {
        string(name: 'ProjectKey')
        string(name: 'RepoName')
        string(name: 'RefID')
        string(name: 'PullRequestID')
    }
    environment {
        BITBUCKET_URL = 'http://172.29.80.1:7990'
        BITBUCKET_USER = 'user'
        BITBUCKET_PASSWORD = 'user'
    }
    options {
        disableConcurrentBuilds()
    }
    stages {
        stage('Build') {
            steps {
                // mock build
                sh "echo git clone ${BITBUCKET_URL}/scm/${ProjectKey}/${RepoName}.git"
                sh "echo git checkout ${RefID}"
                sh "exit 0" // successful build
                // sh "exit 1" // failed build
            }
        }
        stage('Approve PR') {
            steps {
                sh "echo Approve PR"
                sh """
                curl -v -u ${BITBUCKET_USER}:${BITBUCKET_PASSWORD} -XPUT -H "Content-Type: application/json" \
                    --data '{"status": "APPROVED"}' \
                    ${BITBUCKET_URL}/rest/api/1.0/projects/${ProjectKey}/repos/${RepoName}/pull-requests/${PullRequestID}/participants/user
                """
            }
        }
    }
    post { 
        failure { 
            sh "echo Decline PR"
            sh """
            curl -v -u ${BITBUCKET_USER}:${BITBUCKET_PASSWORD}  -H "Content-Type: application/json" \
                --data '{"text": "Build failed. Please test your branch locally."}' \
                ${BITBUCKET_URL}/rest/api/1.0/projects/${ProjectKey}/repos/${RepoName}/pull-requests/${PullRequestID}/comments
             """
        }
    }
}
