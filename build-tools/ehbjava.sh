#!/bin/bash

set -o errexit
set -o pipefail

KLESS_ROOT=$(dirname "${BASH_SOURCE}")/..

echo $KLESS_ROOT

IMGNAME=klessv1/eventhandlerbuilderjava_$DOCKER_ENGINE_VER

TAG=$IMGNAME:$BUILD_ID

if [[ ! -z "$KLESS_DEST_REGISTRY" ]]; then
  TAG=$KLESS_DEST_REGISTRY/$TAG
fi

echo "Building image with tag $TAG"

cd builders/java

if [[ ! -z "$KLESS_SRC_REGISTRY" ]]; then
  sed -e "s/KLESS_SRC_REGISTRY/${KLESS_SRC_REGISTRY}\//g" -e "s/DOCKER_ENGINE_VER/${DOCKER_ENGINE_VER}/g" Dockerfile > Dockerfile.tmp
else 
  sed -e "s/KLESS_SRC_REGISTRY//g" -e "s/DOCKER_ENGINE_VER/${DOCKER_ENGINE_VER}/g" Dockerfile > Dockerfile.tmp
fi

if [[ ! -z "$KLESS_DEST_REGISTRY_USERNAME" ]]; then
  echo "Logging into docker registry $KLESS_DEST_REGISTRY"
  sudo docker login -u $KLESS_DEST_REGISTRY_USERNAME -p $KLESS_DEST_REGISTRY_PASSWORD $KLESS_DEST_REGISTRY
fi
sudo docker build -f Dockerfile.tmp --build-arg KLESS_VERSION=$BUILD_ID --build-arg KLESS_MAINTAINER=paal@thorstensen.org -t $TAG .
sudo docker push $TAG

rm Dockerfile.tmp

cd ../..
