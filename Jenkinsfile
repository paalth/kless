pipeline {
    agent any

    environment {
      REPO_NAME = 'github.com/paalth/kless'
      GOPATH = '/opt/go'
      GOBIN = '/opt/go/bin'
    }

    stages {
        stage('Build') {
            steps {
                withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'DEST_REPO_CREDENTIALS', usernameVariable: 'DEST_REPO_USERNAME', passwordVariable: 'DEST_REPO_PASSWORD']]) {
                   echo 'Building..'
                   sh 'printenv | sort'
                   sh 'rm -rf $GOPATH/*'
                   sh 'mkdir -p $GOPATH/src/$REPO_NAME'
                   sh 'mv * $GOPATH/src/$REPO_NAME'
                   sh 'cd $GOPATH/src/$REPO_NAME; make client'
                   sh 'cd $GOPATH/src/$REPO_NAME; make'
                   echo 'Build complete'
                }
            }
        }

        stage('Test') {
            steps {
                echo 'Testing..'
            }
        }

        stage('Deploy') {
            steps {
                echo 'Deploying....'
            }
        }
    }
}
