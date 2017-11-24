#!/bin/bash

set -o errexit
set -o pipefail

KLESS_ROOT=$(dirname "${BASH_SOURCE}")/..

IMGNAME=klessv1/klessfrontendhttp

TAG=$IMGNAME:$BUILD_ID

if [[ ! -z "$KLESS_DEST_REGISTRY" ]]; then
  TAG=$KLESS_DEST_REGISTRY/$TAG
fi

echo "Building image with tag $TAG"

echo $KLESS_ROOT

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s" -a -installsuffix cgo -o httpfrontend cmd/frontend/http/httpfrontend.go

mv httpfrontend cmd/frontend/http

cd cmd/frontend/http

if [[ ! -z "$KLESS_SRC_REGISTRY" ]]; then
  KLESS_SRC_REGISTRY=KLESS_SRC_REGISTRY/
fi

sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/KLESS_SRC_REGISTRY/${KLESS_SRC_REGISTRY}/g" Dockerfile > Dockerfile.tmp

if [[ ! -z "$KLESS_DEST_REGISTRY_USERNAME" ]]; then
  echo "Logging into docker registry $KLESS_DEST_REGISTRY"
  sudo docker login -u $KLESS_DEST_REGISTRY_USERNAME -p $KLESS_DEST_REGISTRY_PASSWORD $KLESS_DEST_REGISTRY
fi
sudo docker build -f Dockerfile.tmp --build-arg KLESS_VERSION=0.0.1 --build-arg KLESS_MAINTAINER=paal@thorstensen.org -t $TAG .
sudo docker push $TAG

rm Dockerfile.tmp

rm httpfrontend

cd ../../..
