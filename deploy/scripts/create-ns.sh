#!/bin/bash

if [[ ! -z "$KLESS_NAMESPACE_EXISTS"  ]]; then
  echo "Creating namespace: $KLESS_NAMESPACE"
  sed -e "s/KLESS_NAMESPACE/$KLESS_NAMESPACE/g" deploy/prereqs/namespace/kless-ns.yaml > /tmp/kless-ns.yaml
  kubectl create -f /tmp/kless-ns.yaml
  rm /tmp/kless-ns.yaml
else
  echo "Not creating namespace"
fi