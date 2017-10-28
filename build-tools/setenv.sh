# Kubernetes namespace where kless will be installed
export KLESS_NAMESPACE=kless

# Registry where prerequisite images will be pulled from
export SRC_REGISTRY=therepo
export SRC_REGISTRY_USERNAME=username
export SRC_REGISTRY_PASSWORD=password
export SRC_REGISTRY_EMAIL=email

# Registry where built images will be pushed
# When using Jenkins the username and password needs to be stored as credentials named  DEST_REGISTRY_CREDENTIALS
export DEST_REGISTRY=therepo
export DEST_REGISTRY_USERNAME=username
export DEST_REGISTRY_PASSWORD=password
export DEST_REGISTRY_EMAIL=email

export BUILD_ID=123
