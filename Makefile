
SHELL := /bin/bash
MAKEFLAGS += --no-builtin-rules

all:
	build-tools/build.sh

client:
	build-tools/build-client.sh
	
genrpc:
	build-tools/genrpc.sh

install-prereqs: create-ns create-imagepullsecrets deploy-registry deploy-etcd deploy-influxdb deploy-grafana
	echo "Installation of prerequisites complete"

uninstall-prereqs: delete-ns
	echo "Uninstallation of prerequisites complete"

install: deploy-server
	echo "Installation complete"

uninstall: undeploy-server
	echo "Uninstallation complete"

create-ns:
	sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" deploy/prereqs/namespace/kless-ns.yaml > /tmp/kless-ns.yaml
	kubectl create -f /tmp/kless-ns.yaml
	rm /tmp/kless-ns.yaml

delete-ns:
	sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" deploy/prereqs/namespace/kless-ns.yaml > /tmp/kless-ns.yaml
	kubectl delete -f /tmp/kless-ns.yaml
	rm /tmp/kless-ns.yaml

create-imagepullsecrets:
	kubectl create secret docker-registry src-registry-key --docker-server=${KLESS_SRC_REGISTRY} --docker-username=${KLESS_SRC_REGISTRY_USERNAME} --docker-password=${KLESS_SRC_REGISTRY_PASSWORD} --docker-email=${KLESS_SRC_REGISTRY_EMAIL} -n ${KLESS_NAMESPACE}
	kubectl create secret docker-registry dest-registry-key --docker-server=${DEST_REPO} --docker-username=${DEST_REPO_USERNAME} --docker-password=${DEST_REPO_PASSWORD} --docker-email=${DEST_REPO_EMAIL} -n ${KLESS_NAMESPACE}
	kubectl patch serviceaccount default -p '{"imagePullSecrets": [{"name": "src-registry-key"}, {"name": "dest-registry-key"}]}' -n ${KLESS_NAMESPACE}
	
deploy-registry:
	sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/KLESS_SRC_REGISTRY/${KLESS_SRC_REGISTRY}/g" deploy/prereqs/registry/kless-registry.yaml > /tmp/kless-registry.yaml
	kubectl create -f /tmp/kless-registry.yaml
	rm /tmp/kless-registry.yaml

undeploy-registry:
	sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/KLESS_SRC_REGISTRY/${KLESS_SRC_REGISTRY}/g" deploy/prereqs/registry/kless-registry.yaml > /tmp/kless-registry.yaml
	kubectl delete -f /tmp/kless-registry.yaml
	rm /tmp/kless-registry.yaml

deploy-etcd:
	sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/KLESS_SRC_REGISTRY_QUAY/${KLESS_SRC_REGISTRY_QUAY}/g" deploy/prereqs/etcd/kless-etcd.yaml > /tmp/kless-etcd.yaml
	kubectl create -f /tmp/kless-etcd.yaml
	rm /tmp/kless-etcd.yaml

undeploy-etcd:
	sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/KLESS_SRC_REGISTRY_QUAY/${KLESS_SRC_REGISTRY_QUAY}/g" deploy/prereqs/etcd/kless-etcd.yaml > /tmp/kless-etcd.yaml
	kubectl delete -f /tmp/kless-etcd.yaml
	rm /tmp/kless-etcd.yaml

deploy-influxdb:
	sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/KLESS_SRC_REGISTRY/${KLESS_SRC_REGISTRY}/g" deploy/prereqs/influxdb/kless-influxdb.yaml > /tmp/kless-influxdb.yaml
	kubectl create -f /tmp/kless-influxdb.yaml
	rm /tmp/kless-influxdb.yaml

undeploy-influxdb:
	sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/KLESS_SRC_REGISTRY/${KLESS_SRC_REGISTRY}/g" deploy/prereqs/influxdb/kless-influxdb.yaml > /tmp/kless-influxdb.yaml
	kubectl delete -f /tmp/kless-influxdb.yaml
	rm /tmp/kless-influxdb.yaml

deploy-grafana:
	sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/KLESS_SRC_REGISTRY_GCR/${KLESS_SRC_REGISTRY_GCR}/g" deploy/prereqs/grafana/kless-grafana.yaml > /tmp/kless-grafana.yaml
	kubectl create -f /tmp/kless-grafana.yaml
	rm /tmp/kless-grafana.yaml

undeploy-grafana:
	sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/KLESS_SRC_REGISTRY_GCR/${KLESS_SRC_REGISTRY_GCR}/g" deploy/prereqs/influxdb/kless-grafana.yaml > /tmp/kless-grafana.yaml
	kubectl delete -f /tmp/kless-grafana.yaml
	rm /tmp/kless-grafana.yaml

deploy-server:
	sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/DEST_REPO/${DEST_REPO}/g" -e "s/BUILD_ID/${BUILD_ID}/g" deploy/kless-server/kless-server.yaml > /tmp/kless-server.yaml
	kubectl create -f /tmp/kless-server.yaml
	rm /tmp/kless-server.yaml
	kubectl describe svc kless-server -n ${KLESS_NAMESPACE}

undeploy-server:
	sed -e "s/KLESS_NAMESPACE/${KLESS_NAMESPACE}/g" -e "s/DEST_REPO/${DEST_REPO}/g" -e "s/BUILD_ID/${BUILD_ID}/g" deploy/kless-server/kless-server.yaml > /tmp/kless-server.yaml
	kubectl delete -f /tmp/kless-server.yaml
	rm /tmp/kless-server.yaml

ehb-go:
	build-tools/ehbgo.sh

ehb-go-fedora:
	build-tools/ehbgo-fedora.sh

ehb-java:
	build-tools/ehbjava.sh

frontend-utils:
	build-tools/frontend-utils.sh

frontend-http:
	build-tools/frontend-http.sh

frontend-nats:
	build-tools/frontend-nats.sh

frontend-rabbitmq: frontend-utils
	build-tools/frontend-rabbitmq.sh

frontend-kafka: frontend-utils
	build-tools/frontend-kafka.sh

test:
	build-tools/test.sh

test-integration:
	build-tools/test-integration.sh

test-e2e:
	build-tools/test-e2e.sh

release:
	build-tools/release.sh

clean:
	build-tools/make-clean.sh

