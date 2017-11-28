#!/bin/bash

sed -e "s/KLESS_NAMESPACE/$KLESS_NAMESPACE/g" -e "s/KLESS_SRC_REGISTRY_QUAY/$KLESS_SRC_REGISTRY_QUAY/g" deploy/prereqs/etcd/kless-etcd.yaml > /tmp/kless-etcd.yaml
kubectl delete -f /tmp/kless-etcd.yaml
rm /tmp/kless-etcd.yaml