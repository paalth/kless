#!/bin/bash

set -o errexit
set -o pipefail

KLESS_ROOT=$(dirname "${BASH_SOURCE}")/..

IMGNAME=klessv1/klessserver

TAG=$IMGNAME:$BUILD_ID

if [[ ! -z "$KLESS_DEST_REGISTRY" ]]; then
  TAG=$KLESS_DEST_REGISTRY/$TAG
fi

echo "Building image with tag $TAG"

echo $KLESS_ROOT

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s" -a -installsuffix cgo -o klessserver cmd/klessserver/klessserver.go

mv klessserver cmd/klessserver

cd cmd/klessserver

echo "Logging into docker registry $KLESS_DEST_REGISTRY"
sudo docker login -u $KLESS_DEST_REGISTRY_USERNAME -p "$KLESS_DEST_REGISTRY_PASSWORD" $KLESS_DEST_REGISTRY
sudo docker build -f Dockerfile --build-arg KLESS_VERSION=0.0.1 --build-arg KLESS_MAINTAINER=paal@thorstensen.org -t $TAG .
sudo docker push $TAG

rm klessserver

cd ../..
