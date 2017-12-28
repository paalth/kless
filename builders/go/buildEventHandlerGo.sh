#!/bin/bash

echo "Starting build of Go event handler"

mkdir /go/src/eventhandler
mkdir -p /go/src/github.com/paalth/kless/pkg/interface/klessgo

echo "Downloading source files from kless server"

wget $KLESS_SERVER/$KLESS_EVENT_HANDLER_SOURCE -O /go/src/eventhandler/EventHandler.go
wget $KLESS_SERVER/$KLESS_INTERFACE_SOURCE -O /go/src/github.com/paalth/kless/pkg/interface/klessgo/Interface.go
wget $KLESS_SERVER/$KLESS_ENTRYPOINT_SOURCE -O /go/src/InvokeEventHandler.go

echo "Source files downloaded"

export GOBIN=$GOPATH/bin

cd /go/src

echo "Retrieving dependencies"

go get ./...

echo "Building event handler"

go install /go/src/InvokeEventHandler.go

cp /go/bin/InvokeEventHandler /tmp

echo "Build complete"

TAG=$KLESS_REPO/$KLESS_EVENT_HANDLER_NAME:$KLESS_EVENT_HANDLER_VERSION

echo "Container tag = $TAG"

cd /tmp

echo "Retrieving registry information from kless server"

# Retrieve information so we can log into the registry we're pulling the base image from if needed
REGISTRY_USERNAME=`curl "$KLESS_SERVER/cfg?op=get&key=REGISTRY_USERNAME"`
REGISTRY_PASSWORD=`curl "$KLESS_SERVER/cfg?op=get&key=REGISTRY_PASSWORD"`
REGISTRY_HOSTPORT=`curl "$KLESS_SERVER/cfg?op=get&key=REGISTRY_HOSTPORT"`

echo "Registry info retrieved, registry = $REGISTRY_HOSTPORT, username = $REGISTRY_USERNAME"

sed -e "s/REGISTRY_HOSTPORT/$REGISTRY_HOSTPORT/g" Dockerfile > Dockerfile.tmp

echo "Dockerfile updated"

if [[ ! -z "$REGISTRY_USERNAME"  ]]; then
  echo "Logging into Docker registry $REGISTRY_HOSTPORT"
  docker login -u $REGISTRY_USERNAME -p $REGISTRY_PASSWORD $REGISTRY_HOSTPORT
fi

echo "Building image"
docker build -f Dockerfile.tmp -t $TAG .

echo "Pushing image to registry"
docker push $TAG

echo "Image creation complete"
