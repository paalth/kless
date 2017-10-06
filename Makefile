
SHELL := /bin/bash
MAKEFLAGS += --no-builtin-rules

export DEST_REPO=192.168.1.12:32768

all:
	build-tools/build.sh

deploy-registry:
	kubectl create -f prereqs/repository

undeploy-registry:
	kubectl delete -f prereqs/repository

deploy-server:
	kubectl create -f deploy/klessserver/klessserver-rc.yaml
	kubectl create -f deploy/klessserver/klessserver-svc.yaml
	kubectl describe svc kless-server -n kless

redeploy-server:
	kubectl delete -f deploy/klessserver/klessserver-rc.yaml
	kubectl delete -f deploy/klessserver/klessserver-svc.yaml
	kubectl create -f deploy/klessserver/klessserver-rc.yaml
	kubectl create -f deploy/klessserver/klessserver-svc.yaml
	kubectl describe svc kless-server -n kless

genrpc:
	build-tools/genrpc.sh

client:
	build-tools/build-client.sh
	
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

