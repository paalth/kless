javac -cp ../libs/influxdb-java-2.5-SNAPSHOT.jar io/kless/frontend/utils/*.java

jar cvf ../libs/kless-frontend-utils.jar io/kless/frontend/utils/*.class
