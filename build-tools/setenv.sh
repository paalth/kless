# Kubernetes namespace where kless will be installed
export KLESS_NAMESPACE=kless

# Registry where prerequisite images will be pulled from normally hosted on Docker Hub
export KLESS_SRC_REGISTRY=index.docker.io
export KLESS_SRC_REGISTRY_USERNAME=
export KLESS_SRC_REGISTRY_PASSWORD=
export KLESS_SRC_REGISTRY_EMAIL=

# Registry where prerequisite images will be pulled from normally hosted on quay.io
export KLESS_SRC_REGISTRY_QUAY=quay.io
export KLESS_SRC_REGISTRY_QUAY_USERNAME=
export KLESS_SRC_REGISTRY_QUAY_PASSWORD=
export KLESS_SRC_REGISTRY_QUAY_EMAIL=

# Registry where prerequisite images will be pulled from normally hosted on gcr.io
export KLESS_SRC_REGISTRY_GCR=gcr.io
export KLESS_SRC_REGISTRY_GCR_USERNAME=
export KLESS_SRC_REGISTRY_GCR_PASSWORD=
export KLESS_SRC_REGISTRY_GCR_EMAIL=

# Registry where built images will be pushed
# When using Jenkins the username and password needs to be stored as credentials named DEST_REGISTRY_CREDENTIALS
export KLESS_DEST_REGISTRY=therepo
export KLESS_DEST_REGISTRY_USERNAME=username
export KLESS_DEST_REGISTRY_PASSWORD=password
export KLESS_DEST_REGISTRY_EMAIL=email

# Credentials for cluster where kless will be deployed
# When using Jenkins these needs to be set as env vars on the build node
export K8S_CLIENT_CERT_PATH=pathtocacert
export K8S_CLIENT_KEY_PATH=pathtoclientkey
export K8S_SERVER_URL=theurl

# Version of Docker Engine used to build the builders
export DOCKER_ENGINE_VER=1.11.1

# Set manually if build not run by Jenkins
export BUILD_ID=123

# Set if build should send CLI to Google bucket
#export KLESS_UPLOAD_CLI=true
