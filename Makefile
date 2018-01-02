
SHELL := /bin/bash
MAKEFLAGS += --no-builtin-rules

# TODO: do something smarter than sed...

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
	deploy/scripts/create-ns.sh

delete-ns:
	deploy/scripts/delete-ns.sh

create-imagepullsecrets:
	deploy/scripts/create-imagepullsecrets.sh
	
deploy-registry:
	deploy/scripts/deploy-registry.sh

undeploy-registry:
	deploy/scripts/undeploy-registry.sh

deploy-etcd:
	deploy/scripts/deploy-etcd.sh

undeploy-etcd:
	deploy/scripts/undeploy-etcd.sh

deploy-influxdb:
	deploy/scripts/deploy-influxdb.sh

undeploy-influxdb:
	deploy/scripts/undeploy-influxdb.sh

deploy-grafana:
	deploy/scripts/deploy-grafana.sh

undeploy-grafana:
	deploy/scripts/undeploy-grafana.sh

deploy-server:
	deploy/scripts/deploy-server.sh

undeploy-server:
	deploy/scripts/undeploy-server.sh

ehb-go:
	build-tools/ehbgo.sh

ehb-java:
	build-tools/ehbjava.sh

ehb-python:
	build-tools/ehbpython.sh

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

