#!/bin/bash

IMAGE_PULL_SECRETS=""

if [[ ! -z "$KLESS_SRC_REGISTRY_USERNAME" ]]; then
  kubectl create secret docker-registry src-registry-key --docker-server=$KLESS_SRC_REGISTRY --docker-username=$KLESS_SRC_REGISTRY_USERNAME --docker-password=$KLESS_SRC_REGISTRY_PASSWORD --docker-email=$KLESS_SRC_REGISTRY_EMAIL -n $KLESS_NAMESPACE
  IMAGE_PULL_SECRETS='{"name": "src-registry-key"}'
fi

if [[ ! -z "$KLESS_SRC_REGISTRY_QUAY_USERNAME" ]]; then
  kubectl create secret docker-registry src-registry-key-quay --docker-server=$KLESS_SRC_REGISTRY_QUAY --docker-username=$KLESS_SRC_REGISTRY_QUAY_USERNAME --docker-password=$KLESS_SRC_REGISTRY_QUAY_PASSWORD --docker-email=$KLESS_SRC_REGISTRY_QUAY_EMAIL -n $KLESS_NAMESPACE

  if [[ ! -z "$IMAGE_PULL_SECRETS" ]]; then
    IMAGE_PULL_SECRETS=$IMAGE_PULL_SECRETS','
  fi
  IMAGE_PULL_SECRETS=$IMAGE_PULL_SECRETS'{"name": "src-registry-key-quay"}'
fi

if [[ ! -z "$KLESS_SRC_REGISTRY_GCR_USERNAME" ]]; then
  kubectl create secret docker-registry src-registry-key-gcr --docker-server=$KLESS_SRC_REGISTRY_GCR --docker-username=$KLESS_SRC_REGISTRY_GCR_USERNAME --docker-password=$KLESS_SRC_REGISTRY_GCR_PASSWORD --docker-email=$KLESS_SRC_REGISTRY_GCR_EMAIL -n $KLESS_NAMESPACE

  if [[ ! -z "$IMAGE_PULL_SECRETS" ]]; then
    IMAGE_PULL_SECRETS=$IMAGE_PULL_SECRETS','
  fi
  IMAGE_PULL_SECRETS=$IMAGE_PULL_SECRETS'{"name": "src-registry-key-gcr"}'
fi

if [[ ! -z "$KLESS_DEST_REGISTRY_USERNAME" ]]; then
  kubectl create secret docker-registry dest-registry-key --docker-server=$KLESS_DEST_REGISTRY --docker-username=$KLESS_DEST_REGISTRY_USERNAME --docker-password=$KLESS_DEST_REGISTRY_PASSWORD --docker-email=$KLESS_DEST_REGISTRY_EMAIL -n $KLESS_NAMESPACE
  
  if [[ ! -z "$IMAGE_PULL_SECRETS" ]]; then
    IMAGE_PULL_SECRETS=$IMAGE_PULL_SECRETS','
  fi
  IMAGE_PULL_SECRETS=$IMAGE_PULL_SECRETS'{"name": "dest-registry-key"}'
fi

PATCH_COMMAND='{"imagePullSecrets": ['$IMAGE_PULL_SECRETS']}' 

echo "Adding image pull secrets: $PATCH_COMMAND"

kubectl patch serviceaccount default -p "$PATCH_COMMAND" -n $KLESS_NAMESPACE