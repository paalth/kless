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
                sh 'kubectl config set-credentials k8s-user --client-certificate=$K8S_CLIENT_CERT_PATH --client-key=$K8S_CLIENT_KEY_PATH'
                sh 'kubectl config set-cluster k8s-cluster --insecure-skip-tls-verify=true --server=$K8S_SERVER_URL'
                sh 'kubectl config set-context k8s --cluster=k8s-cluster --user=k8s-user --namespace=$KLESS_NAMESPACE'
                sh 'kubectl config use-context k8s'
                sh 'kubectl set image kless/kless-server kless-server=$DEST_REPO/klessv1/klessserver:$BUILD_ID'

                echo 'Deploy complete'
            }
        }
    }
}
