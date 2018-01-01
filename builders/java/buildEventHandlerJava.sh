#!/bin/bash

echo "Starting build of Java event handler"

mkdir -p io/kless

wget $KLESS_SERVER/$KLESS_EVENT_HANDLER_SOURCE -O io/kless/EventHandler1.java
wget $KLESS_SERVER/$KLESS_INTERFACE_SOURCE -O io/kless/EventHandlerInterface.java
wget $KLESS_SERVER/$KLESS_ENTRYPOINT_SOURCE -O io/kless/InvokeEventHandler.java
wget $KLESS_SERVER/$KLESS_CONTEXT_SOURCE -O io/kless/Context.java
wget $KLESS_SERVER/$KLESS_REQUEST_SOURCE -O io/kless/Request.java
wget $KLESS_SERVER/$KLESS_RESPONSE_SOURCE -O io/kless/Response.java

echo "Source downloaded"

javac io/kless/*.java
jar cvf InvokeEventHandler.jar io/kless/*.class
cp InvokeEventHandler.jar /tmp

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

rm Dockerfile.tmp

echo "Image creation complete"
