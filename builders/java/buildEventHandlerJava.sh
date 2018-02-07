#!/bin/bash

cd /tmp

echo "Starting build of Java event handler"

CURRENT_STATUS="etcd?op=sethandlerstatus&handler=$KLESS_EVENT_HANDLER_NAME:$KLESS_EVENT_HANDLER_VERSION&status=BuildInit"
curl -s $KLESS_SERVER/$CURRENT_STATUS 

mkdir -p io/kless

echo "Starting download of source files"

curl -o io/kless/EventHandler1.java $KLESS_SERVER/$KLESS_EVENT_HANDLER_SOURCE
curl -o io/kless/EventHandlerInterface.java $KLESS_SERVER/$KLESS_INTERFACE_SOURCE
curl -o io/kless/InvokeEventHandler.java $KLESS_SERVER/$KLESS_ENTRYPOINT_SOURCE
curl -o io/kless/Context.java $KLESS_SERVER/$KLESS_CONTEXT_SOURCE
curl -o io/kless/Request.java $KLESS_SERVER/$KLESS_REQUEST_SOURCE
curl -o io/kless/Response.java $KLESS_SERVER/$KLESS_RESPONSE_SOURCE

echo "Source downloaded"

CP="InvokeEventHandler.jar"

mkdir klesslib

if [[ ! -z "$KLESS_DEPENDENCIES_URL" ]]; then
  cd klesslib

  echo "Retrieving dependencies from $KLESS_DEPENDENCIES_URL"

  # Download list of dependencies
  curl -o deps.txt $KLESS_DEPENDENCIES_URL

  # Download all dependencies listed in file
  awk '{ cmd="curl -O "$1; system(cmd) }' deps.txt

  rm deps.txt
  cd ..

  # Determine classpath
  for file in $( find . -type f -name "*.jar" )
  do
    echo $file
    CP=$CP:$file
  done

  echo "Dependencies retrieved"
else 
  echo "No dependencies"
fi

echo "Classpath:"
echo $CP

echo "Compiling"
javac -cp $CP io/kless/*.java > stdout.txt &> stderr.txt
if [ "$?" -ne 0 ]; then 
  echo "Compilation failed, error:";
  cat stderr.txt;
  CURRENT_STATUS="etcd?op=sethandlerstatus&handler=$KLESS_EVENT_HANDLER_NAME:$KLESS_EVENT_HANDLER_VERSION&status=BuildError";
  curl -s $KLESS_SERVER/$CURRENT_STATUS;
  BUILD_OUTPUT="etcd?op=setbuildoutput&handler=$KLESS_EVENT_HANDLER_NAME:$KLESS_EVENT_HANDLER_VERSION";
  curl -s --request POST --data-binary @stderr.txt --header "Content-Type:application/octet-stream" $KLESS_SERVER/$BUILD_OUTPUT;
  exit 0; 
fi
jar cvf InvokeEventHandler.jar io/kless/*.class

echo "Build complete"

TAG=$KLESS_REPO/$KLESS_EVENT_HANDLER_NAME:$KLESS_EVENT_HANDLER_VERSION

echo "Container tag = $TAG"

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
docker build -f Dockerfile.tmp --build-arg KLESS_CP=$CP -t $TAG .

echo "Pushing image to registry"
docker push $TAG

rm Dockerfile.tmp

echo "Reporting complete status"
CURRENT_STATUS="etcd?op=sethandlerstatus&handler=$KLESS_EVENT_HANDLER_NAME:$KLESS_EVENT_HANDLER_VERSION&status=BuildComplete"
curl -s $KLESS_SERVER/$CURRENT_STATUS 

echo "Image creation complete"
