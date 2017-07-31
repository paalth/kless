# Add etcd-operator TPR to create etcd cluster, eventually move this to helm when it can pass the verification step
kubectl create -f kless/etcd-op.txt

# Add InfluxDB database, eventually move this to helm with a post-install hook or something like that 

# Docker registry to retrieve images from
NAMESPACE=kless
REGISTRY=192.168.1.12:32768
USER=test
PASS=test
EMAIL=test@test.com

kubectl create secret docker-registry regsecret -n $NAMESPACE --docker-server=$REGISTRY --docker-username=$USER --docker-password=$PASS --docker-email=$EMAIL
kubectl patch serviceaccount default -n $NAMESPACE -p '{"imagePullSecrets": [{"name": "regsecret"}]}'

# Add event handler builders 
kless create builder -b go-1.7.4 -u $REGISTRY/eventhandlerbuildergo:0.0.1 -i "KlessInterface=../pkg/interface/klessgo/Interface.go,InvokeEventHandler=../pkg/invoke/InvokeEventHandler.go"

# Add frontend types
kless create frontendtype -t http -u $REGISTRY/klessfrontendhttp:0.0.1
