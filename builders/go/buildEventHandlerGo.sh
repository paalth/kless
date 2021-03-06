#!/bin/bash

echo "Starting build of Go event handler, builder version = 0.0.1"

CURRENT_STATUS="etcd?op=sethandlerstatus&handler=$KLESS_EVENT_HANDLER_NAME:$KLESS_EVENT_HANDLER_VERSION&status=BuildInit"
curl -s $KLESS_SERVER/$CURRENT_STATUS 

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

go install /go/src/InvokeEventHandler.go > stdout.txt &> stderr.txt
if [ "$?" -ne 0 ]; then 
  echo "Compilation failed, error:";
  cat stderr.txt;
  CURRENT_STATUS="etcd?op=sethandlerstatus&handler=$KLESS_EVENT_HANDLER_NAME:$KLESS_EVENT_HANDLER_VERSION&status=BuildError";
  curl -s $KLESS_SERVER/$CURRENT_STATUS;
  BUILD_OUTPUT="etcd?op=setbuildoutput&handler=$KLESS_EVENT_HANDLER_NAME:$KLESS_EVENT_HANDLER_VERSION";
  curl -s --request POST --data-binary @stderr.txt --header "Content-Type:application/octet-stream" $KLESS_SERVER/$BUILD_OUTPUT;
  exit 0; 
fi

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

if [[ ! -z "$REGISTRY_HOSTPORT" ]]; then
  echo $REGISTRY_HOSTPORT | grep '/$'
  if [[ $? -eq 0 ]]; then
    REGISTRY_HOSTPORT=${REGISTRY_HOSTPORT%?};
  fi
  sed -e "s/REGISTRY_HOSTPORT/${REGISTRY_HOSTPORT}\//g" Dockerfile > Dockerfile.tmp
else 
  sed -e "s/REGISTRY_HOSTPORT//g" Dockerfile > Dockerfile.tmp
fi

echo "Dockerfile updated"

if [[ ! -z "$REGISTRY_USERNAME"  ]]; then
  echo "Logging into Docker registry $REGISTRY_HOSTPORT"
  docker login -u $REGISTRY_USERNAME -p $REGISTRY_PASSWORD $REGISTRY_HOSTPORT
fi

echo "Building image"
docker build -f Dockerfile.tmp -t $TAG .

echo "Pushing image to registry"
docker push $TAG

echo "Reporting complete status"
CURRENT_STATUS="etcd?op=sethandlerstatus&handler=$KLESS_EVENT_HANDLER_NAME:$KLESS_EVENT_HANDLER_VERSION&status=BuildComplete"
curl -s $KLESS_SERVER/$CURRENT_STATUS 

echo "Requesting deployment of built handler"
DEPLOY_REQ="api?op=deploy&handlerName=$KLESS_EVENT_HANDLER_NAME&handlerNamespace=$KLESS_EVENT_HANDLER_NAMESPACE&handlerVersion=$KLESS_EVENT_HANDLER_VERSION&handlerId=$KLESS_EVENT_HANDLER_ID"
curl -s $KLESS_SERVER/$DEPLOY_REQ

echo "Image creation complete"
