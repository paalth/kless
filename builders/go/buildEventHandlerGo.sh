#!/bin/bash

echo "Starting build of Go event handler"

mkdir /go/src/eventhandler
mkdir -p /go/src/github.com/paalth/kless/pkg/interface/klessgo

wget $KLESS_EVENT_HANDLER_SOURCE -O /go/src/eventhandler/EventHandler.go
wget $KLESS_INTERFACE_SOURCE -O /go/src/github.com/paalth/kless/pkg/interface/klessgo/Interface.go
wget $KLESS_ENTRYPOINT_SOURCE -O /go/src/InvokeEventHandler.go

echo "Source downloaded"

export GOBIN=$GOPATH/bin

cd /go/src
go get ./...

go install /go/src/InvokeEventHandler.go

cp /go/bin/InvokeEventHandler /tmp

echo "Build complete"

TAG=$KLESS_REPO/$KLESS_EVENT_HANDLER_NAME:$KLESS_EVENT_HANDLER_VERSION

echo "Container tag = $TAG"

cd /tmp
docker build -t $TAG .
docker push $TAG

echo "Image creation complete"
