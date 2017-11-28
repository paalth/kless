#!/bin/bash

sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/KLESS_SRC_REGISTRY_GCR/${KLESS_SRC_REGISTRY_GCR}/g" deploy/prereqs/grafana/kless-grafana.yaml > /tmp/kless-grafana.yaml
kubectl create -f /tmp/kless-grafana.yaml
rm /tmp/kless-grafana.yaml