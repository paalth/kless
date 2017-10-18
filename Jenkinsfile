pipeline {
    agent any

    environment {
      REPO_NAME = 'github.com/paalth/kless'
      GOPATH = '/opt/go'
      GOBIN = '$GOPATH/bin'
    }

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                sh 'printenv'
                sh 'rm -rf ${env.GOPATH}
                sh 'mkdir -p ${env.GOPATH}/src/${env.REPO_NAME}'
                sh 'ln -svf * ${env.GOPATH}/src/${env.REPO_NAME}'
                sh 'cd ${env.GOPATH}/src/${env.REPO_NAME}'
                sh 'make client'
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
