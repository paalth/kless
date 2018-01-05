#!/bin/bash

set -o errexit
set -o pipefail

KLESS_ROOT=$(dirname "${BASH_SOURCE}")/..

IMGNAME=klessv1/klessfrontendrabbitmq

TAG=$IMGNAME:$BUILD_ID

if [[ ! -z "$KLESS_DEST_REGISTRY" ]]; then
  TAG=$KLESS_DEST_REGISTRY/$TAG
fi

echo "Building image with tag $TAG"

echo $KUBELESS_ROOT

cd cmd/frontend/rabbitmq

mkdir common-libs
cp ../libs/* common-libs

./build.sh

if [[ ! -z "$KLESS_SRC_REGISTRY" ]]; then
  sed -e "s/KLESS_SRC_REGISTRY/${KLESS_SRC_REGISTRY}\//g" Dockerfile > Dockerfile.tmp
else 
  sed -e "s/KLESS_SRC_REGISTRY//g" Dockerfile > Dockerfile.tmp
fi

if [[ ! -z "$KLESS_DEST_REGISTRY_USERNAME" ]]; then
  echo "Logging into docker registry $KLESS_DEST_REGISTRY"
  sudo docker login -u $KLESS_DEST_REGISTRY_USERNAME -p $KLESS_DEST_REGISTRY_PASSWORD $KLESS_DEST_REGISTRY
fi
sudo docker build -f Dockerfile.tmp --build-arg KLESS_VERSION=0.0.1 --build-arg KLESS_MAINTAINER=paal@thorstensen.org -t $TAG .
sudo docker push $TAG

rm Dockerfile.tmp
rm io/kless/frontend/*.class
rm kless-frontend-rabbitmq.jar 
rm common-libs/*
rmdir common-libs
rm ../libs/kless-frontend-utils.jar

cd ../../..
