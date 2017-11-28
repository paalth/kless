#!/bin/bash

sed -e "s/KLESS_NAMESPACE/$KLESS_NAMESPACE/g" -e "s/KLESS_SRC_REGISTRY/$KLESS_SRC_REGISTRY/g" deploy/prereqs/influxdb/kless-influxdb.yaml > /tmp/kless-influxdb.yaml
kubectl delete -f /tmp/kless-influxdb.yaml
rm /tmp/kless-influxdb.yaml