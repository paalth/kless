FROM KLESS_SRC_REGISTRYgolang:1.7.4

# This Dockerfile builds an image that is used to build event handler builders implemented in Go
#
# For now the procees to build builders is pretty simple and just uses a shell script
# The shell script uses the Docker Engine client which obviously is a less than optimal solution
#
# All builders are invoked by a simple entrypoint called 'InvokeEventHandler' which is itself a small piece of Go code

ARG KLESS_MAINTAINER
ARG KLESS_VERSION

LABEL maintainer=$KLESS_MAINTAINER
LABEL version=$KLESS_VERSION

# Install Docker client
RUN apt-get -y update && \
    apt-get install -y apt-transport-https ca-certificates gnupg2 && \
    apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D && \
    echo "deb https://apt.dockerproject.org/repo debian-jessie main" >> /etc/apt/sources.list.d/docker.list && \
    apt-get -y update && \
    apt-get install -y docker-engine=DOCKER_ENGINE_VER-0~jessie

COPY buildEventHandlerGo.sh /tmp
COPY buildEventHandlerGoDockerfile.tmp /tmp/Dockerfile

# Run command to build the Go event handler
CMD ["/tmp/buildEventHandlerGo.sh"]

