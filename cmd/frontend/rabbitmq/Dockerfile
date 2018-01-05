FROM KLESS_SRC_REGISTRYopenjdk:8u151-jdk-alpine

ARG KLESS_MAINTAINER
ARG KLESS_VERSION

MAINTAINER $KLESS_MAINTAINER

LABEL version=$KLESS_VERSION

COPY libs/amqp-client-4.0.1.jar amqp-client-4.0.1.jar

COPY common-libs/kless-frontend-utils.jar kless-frontend-utils.jar
COPY common-libs/slf4j-api-1.7.22.jar slf4j-api-1.7.22.jar
COPY common-libs/slf4j-simple-1.7.22.jar slf4j-simple-1.7.22.jar
COPY common-libs/influxdb-java-2.5-SNAPSHOT.jar influxdb-java-2.5-SNAPSHOT.jar
COPY common-libs/guava-20.0.jar guava-20.0.jar
COPY common-libs/moshi-1.2.0.jar moshi-1.2.0.jar
COPY common-libs/okhttp-3.5.0.jar okhttp-3.5.0.jar
COPY common-libs/okio-1.11.0.jar okio-1.11.0.jar
COPY common-libs/retrofit-2.1.0.jar retrofit-2.1.0.jar 
COPY common-libs/converter-moshi-2.1.0.jar converter-moshi-2.1.0.jar
COPY common-libs/logging-interceptor-3.5.0.jar logging-interceptor-3.5.0.jar

COPY kless-frontend-rabbitmq.jar kless-frontend-rabbitmq.jar

CMD ["java", "-cp", "amqp-client-4.0.1.jar:kless-frontend-utils.jar:slf4j-api-1.7.22.jar:slf4j-simple-1.7.22.jar:influxdb-java-2.5-SNAPSHOT.jar:guava-20.0.jar:moshi-1.2.0.jar:okhttp-3.5.0.jar:okio-1.11.0.jar:retrofit-2.1.0.jar:converter-moshi-2.1.0.jar:logging-interceptor-3.5.0.jar:kless-frontend-rabbitmq.jar", "io.kless.frontend.KlessRabbitMQConsumer"]

