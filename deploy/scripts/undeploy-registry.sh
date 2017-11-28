#!/bin/bash

if [[ ! -z "$KLESS_SRC_REGISTRY" ]]; then
  sed -e "s/KLESS_NAMESPACE/$KLESS_NAMESPACE/g" -e "s/KLESS_SRC_REGISTRY/$KLESS_SRC_REGISTRY\//g" deploy/prereqs/registry/kless-registry.yaml > /tmp/kless-registry.yaml
else
  sed -e "s/KLESS_NAMESPACE/$KLESS_NAMESPACE/g" -e "s/KLESS_SRC_REGISTRY//g" deploy/prereqs/registry/kless-registry.yaml > /tmp/kless-registry.yaml
fi

kubectl delete -f /tmp/kless-registry.yaml
rm /tmp/kless-registry.yaml