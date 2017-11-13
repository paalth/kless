#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

KLESS_ROOT=$(dirname "${BASH_SOURCE}")/..

echo $KLESS_ROOT

TAG=$KLESS_DEST_REGISTRY/eventhandlerbuildergo_$DOCKER_ENGINE_VER:$BUILD_ID

echo "Building image with tag $TAG"

go install ./...

cd builders/go

sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/KLESS_SRC_REGISTRY/${KLESS_SRC_REGISTRY}/g" -e "s/DOCKER_ENGINE_VER/${DOCKER_ENGINE_VER}/g" Dockerfile > Dockerfile.tmp
sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/KLESS_SRC_REGISTRY/${KLESS_SRC_REGISTRY}/g" buildEventHandlerGoDockerfile > buildEventHandlerGoDockerfile.tmp

if [[ ! -z "$KLESS_DEST_REGISTRY_USERNAME" ]]; then
  echo "Logging into docker registry $KLESS_DEST_REGISTRY"
  docker login -u $KLESS_DEST_REGISTRY_USERNAME -p $KLESS_DEST_REGISTRY_PASSWORD $KLESS_DEST_REGISTRY
fi
docker build -f Dockerfile.tmp --build-arg KLESS_VERSION=$BUILD_ID --build-arg KLESS_MAINTAINER=paal@thorstensen.org -t $TAG .
docker push $TAG

rm Dockerfile.tmp
rm buildEventHandlerGoDockerfile.tmp

cd ../..
