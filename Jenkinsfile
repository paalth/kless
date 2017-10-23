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

                   echo 'Build kless CLI'
                   sh 'cd $GOPATH/src/$REPO_NAME; make client'

                   echo 'Build kless server'
                   sh 'cd $GOPATH/src/$REPO_NAME; make'
                   echo 'Build complete'
                }
            }
        }

        stage('Test') {
            steps {
                echo 'Testing..'
                echo 'Tests TBD...'
                echo 'Test complete'
            }
        }

        stage('Deploy') {
            steps {
                echo 'Deploying....'

                echo 'Deploy kless CLI'

                echo 'Deploy kless server'

                echo 'Deploy complete'
            }
        }
    }
}
