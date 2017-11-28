#!/bin/bash

sed -e "s/KLESS_NAMESPACE/$KLESS_NAMESPACE/g" deploy/prereqs/namespace/kless-ns.yaml > /tmp/kless-ns.yaml
kubectl create -f /tmp/kless-ns.yaml
rm /tmp/kless-ns.yaml