FROM KLESS_SRC_REGISTRYubuntu:16.04

ARG KLESS_MAINTAINER
ARG KLESS_VERSION

LABEL maintainer=$KLESS_MAINTAINER
LABEL version=$KLESS_VERSION

# Install Docker client and curl
RUN apt-get -y update && \
    apt-get install -y apt-transport-https ca-certificates gnupg2 && \
    apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D && \
    echo "deb https://apt.dockerproject.org/repo debian-jessie main" >> /etc/apt/sources.list.d/docker.list && \
    apt-get -y update && \
    apt-get install -y docker-engine=DOCKER_ENGINE_VER-0~jessie && \
    apt-get install -y software-properties-common && \
    add-apt-repository ppa:webupd8team/java -y && \
    apt-get update && \
    apt-get install -y curl && \
    apt-get install -y wget && \
    apt-get clean

COPY buildEventHandlerPython.sh /tmp
COPY buildEventHandlerPythonDockerfile /tmp/Dockerfile

# Run command to build the Python event handler
CMD ["/tmp/buildEventHandlerPython.sh"]

