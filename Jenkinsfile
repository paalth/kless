pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                make client
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
