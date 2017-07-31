#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

KLESS_ROOT=$(dirname "${BASH_SOURCE}")/..

EHB=eventhandlerbuildergo
VERSION=0.0.1

TAG=$REPO/$EHB:$VERSION

echo $KLESS_ROOT

go install ./...

cd builders/go

docker build -f Dockerfile --build-arg KLESS_VERSION=0.0.1 --build-arg KLESS_MAINTAINER=paal@thorstensen.org -t $TAG .
docker push $TAG

cd ../..
