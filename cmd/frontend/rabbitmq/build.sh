javac -cp libs/amqp-client-4.0.1.jar:../libs/kless-frontend-utils.jar:../libs/influxdb-java-2.5-SNAPSHOT.jar io/kless/frontend/*.java

jar cvf kless-frontend-rabbitmq.jar io/kless/frontend/*.class
