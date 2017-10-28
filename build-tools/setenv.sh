# Kubernetes namespace where kless will be installed
export KLESS_NAMESPACE=kless

# Registry where prerequisite images will be pulled from normally hosted on Docker Hub
export KLESS_SRC_REGISTRY=therepo
export KLESS_SRC_REGISTRY_USERNAME=username
export KLESS_SRC_REGISTRY_PASSWORD=password
export KLESS_SRC_REGISTRY_EMAIL=email

# Registry where prerequisite images will be pulled from normally hosted on Quay.io
export KLESS_SRC_REGISTRY_QUAY=therepo
export KLESS_SRC_REGISTRY_QUAY_USERNAME=username
export KLESS_SRC_REGISTRY_QUAY_PASSWORD=password
export KLESS_SRC_REGISTRY_QUAY_EMAIL=email

# Registry where prerequisite images will be pulled from normally hosted on gcr.io
export KLESS_SRC_REGISTRY_GCR=therepo
export KLESS_SRC_REGISTRY_GCR_USERNAME=username
export KLESS_SRC_REGISTRY_GCR_PASSWORD=password
export KLESS_SRC_REGISTRY_GCR_EMAIL=email

# Registry where built images will be pushed
# When using Jenkins the username and password needs to be stored as credentials named DEST_REGISTRY_CREDENTIALS
export DEST_REGISTRY=therepo
export DEST_REGISTRY_USERNAME=username
export DEST_REGISTRY_PASSWORD=password
export DEST_REGISTRY_EMAIL=email

export BUILD_ID=123
