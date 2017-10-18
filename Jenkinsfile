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
                echo 'Building..'
                sh 'printenv | sort'
                sh 'rm -rf $GOPATH/*'
                sh 'mkdir -p $GOPATH/src/$REPO_NAME'
                sh 'ln -svf * $GOPATH/src/$REPO_NAME'
                sh 'cd $GOPATH/src/$REPO_NAME; make client'
                echo 'Build complete'
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
