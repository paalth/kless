#!/bin/bash

if [[ ! -z "$KLESS_SRC_REGISTRY" ]]; then
  sed -e "s/KLESS_NAMESPACE/$KLESS_NAMESPACE/g" -e "s/KLESS_SRC_REGISTRY/$KLESS_SRC_REGISTRY\//g" deploy/prereqs/influxdb/kless-influxdb.yaml > /tmp/kless-influxdb.yaml
else
  sed -e "s/KLESS_NAMESPACE/$KLESS_NAMESPACE/g" -e "s/KLESS_SRC_REGISTRY//g" deploy/prereqs/influxdb/kless-influxdb.yaml > /tmp/kless-influxdb.yaml
fi

kubectl create -f /tmp/kless-influxdb.yaml
rm /tmp/kless-influxdb.yaml
#curl -XPOST -G 'http://10.245.1.3:31734/query' --data-urlencode "q=CREATE DATABASE klessdb" -- k8s