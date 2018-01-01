FROM KLESS_SRC_REGISTRYubuntu:16.04

ARG KLESS_MAINTAINER
ARG KLESS_VERSION

LABEL maintainer=$KLESS_MAINTAINER
LABEL version=$KLESS_VERSION

# Install JDK8 and Docker client
RUN apt-get -y update && \
    apt-get install -y apt-transport-https ca-certificates gnupg2 && \
    apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D && \
    echo "deb https://apt.dockerproject.org/repo debian-jessie main" >> /etc/apt/sources.list.d/docker.list && \
    apt-get -y update && \
    apt-get install -y docker-engine=DOCKER_ENGINE_VER-0~jessie && \
    apt-get install -y software-properties-common && \
    add-apt-repository ppa:webupd8team/java -y && \
    apt-get update && \
    echo oracle-java7-installer shared/accepted-oracle-license-v1-1 select true | /usr/bin/debconf-set-selections && \
    apt-get install -y oracle-java8-installer && \
    apt-get install -y curl && \
    apt-get clean

COPY buildEventHandlerJava.sh /tmp
COPY buildEventHandlerJavaDockerfile /tmp/Dockerfile

# Run command to build the Java event handler
CMD ["/tmp/buildEventHandlerJava.sh"]

