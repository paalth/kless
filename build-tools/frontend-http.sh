#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

KLESS_ROOT=$(dirname "${BASH_SOURCE}")/..

IMGNAME=klessfrontendhttp
VER=0.0.1

TAG=$REPO/$IMGNAME:$VER

echo $KLESS_ROOT

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s" -a -installsuffix cgo -o httpfrontend cmd/frontend/http/httpfrontend.go

mv httpfrontend cmd/frontend/http

cd cmd/frontend/http

docker build -f Dockerfile --build-arg KLESS_VERSION=0.0.1 --build-arg KLESS_MAINTAINER=paal@thorstensen.org -t $TAG .
docker push $TAG

rm httpfrontend

cd ../../..
