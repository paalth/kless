#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

KLESS_ROOT=$(dirname "${BASH_SOURCE}")/..

IMGNAME=klessv1/klessserver
VER=0.0.1

TAG=$DEST_REPO/$IMGNAME:$BUILD_ID

echo $KLESS_ROOT

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s" -a -installsuffix cgo -o klessserver cmd/klessserver/klessserver.go

mv klessserver cmd/klessserver

cd cmd/klessserver

sudo docker login -u $DEST_REPO_USERNAME -p $DEST_REPO_PASSWORD $DEST_REPO
sudo docker build -f Dockerfile --build-arg KLESS_VERSION=0.0.1 --build-arg KLESS_MAINTAINER=paal@thorstensen.org -t $TAG .
sudo docker push $TAG

rm klessserver

cd ../..
