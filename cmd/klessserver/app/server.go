package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/paalth/kless/pkg/etcdinterface"
	k "github.com/paalth/kless/pkg/k8sinterface"
	apiserver "github.com/paalth/kless/pkg/klessserver/apiserver"
	klessapi "github.com/paalth/kless/pkg/klessserver/grpc"

	"google.golang.org/grpc"
)

// Run starts the GRPC server functionality....
func Run() error {
	go StartServer()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	klessapi.RegisterKlessAPIServer(s, &apiserver.APIServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}

// StartServer starts the REST etcd content and configuration serving functionality...
func StartServer() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/etcd", etcdHandler)
	serverMux.HandleFunc("/k8s", k8sHandler)
	serverMux.HandleFunc("/cfg", cfgHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8010", serverMux))
}

func cfgHandler(w http.ResponseWriter, r *http.Request) {
	op := r.URL.Query().Get("op")
	key := r.URL.Query().Get("key")

	fmt.Printf("cfg handler, op = %s key = %s\n", op, key)

	switch op {
	case "get":
		switch key {
		case "REGISTRY_USERNAME":
			value := os.Getenv("REGISTRY_USERNAME")
			fmt.Fprintf(w, "%s", value)
		case "REGISTRY_PASSWORD":
			value := os.Getenv("REGISTRY_PASSWORD")
			fmt.Fprintf(w, "%s", value)
		case "REGISTRY_HOSTPORT":
			value := os.Getenv("REGISTRY_HOSTPORT")
			fmt.Fprintf(w, "%s", value)
		default:
			fmt.Fprintf(w, "unsupported key parameter = %s", key)
		}
	default:
		fmt.Fprintf(w, "op query parameter must currently be 'get'")
	}

}

func k8sHandler(w http.ResponseWriter, r *http.Request) {
	op := r.URL.Query().Get("op")
	name := r.URL.Query().Get("name")
	namespace := r.URL.Query().Get("namespace")

	fmt.Printf("k8s handler, op = %s name = %s namespace = %s\n", op, name, namespace)

	k8s := &k.K8sInterface{}

	switch op {
	case "getsecret":
		fmt.Printf("get secret from k8s\n")
		secret, err := k8s.GetSecret(name, namespace)
		if nil != err {
			log.Fatal("Unable to get secret")
			fmt.Fprint(w, "{ status : failure }")
		} else {
			fmt.Printf("got secret from k8")
			for k, v := range secret {
				fmt.Fprintf(w, "%s: %s\n", k, v)
			}
		}
	default:
		fmt.Fprintf(w, "op query parameter must be one of: 'getsecret'")
	}

}

func etcdHandler(w http.ResponseWriter, r *http.Request) {
	op := r.URL.Query().Get("op")
	key := r.URL.Query().Get("key")
	builder := r.URL.Query().Get("builder")
	handler := r.URL.Query().Get("handler")
	status := r.URL.Query().Get("status")

	fmt.Printf("etcd handler, op = %s builder = %s key = %s handler = %s status = %s\n", op, builder, key, handler, status)

	e := &etcdinterface.EtcdInterface{}

	var etcdkey string

	switch op {
	case "get":
		etcdkey = "/kless/buildinfo/" + builder + "/" + key
		fmt.Printf("get value from etcd, key = %s\n", etcdkey)
		value, err := e.GetValue(etcdkey)
		if nil != err {
			log.Fatal("Unable to get value")
			fmt.Fprint(w, "{ status : failure }")
		} else {
			fmt.Printf("got value from etcd, key = %s value = %s\n", etcdkey, value)
			fmt.Fprint(w, value)
		}
	case "put":
		etcdkey = "/kless/buildinfo/" + builder + "/" + key
		fmt.Printf("put value to etcd, key = %s\n", etcdkey)
		value, err := ioutil.ReadAll(r.Body)
		if nil != err {
			log.Fatal("Unable to read request body")
		}

		err = e.SetValue(etcdkey, string(value))
		if nil != err {
			log.Fatal("Unable to set value")
		}

		fmt.Fprintf(w, "{ status: OK }")
	case "sethandlerstatus":
		etcdkey = "/kless/handlerstatus/" + handler
		fmt.Printf("put value to etcd, key = %s\n", etcdkey)

		err := e.SetValue(etcdkey, status)
		if nil != err {
			log.Fatal("Unable to set value")
		}

		fmt.Fprintf(w, "{ status: OK }")
	case "gethandlerstatus":
		etcdkey = "/kless/handlerstatus/" + handler
		fmt.Printf("get value from etcd, key = %s\n", etcdkey)
		value, err := e.GetValue(etcdkey)
		if nil != err {
			log.Fatal("Unable to get value")
			fmt.Fprint(w, "{ status : failure }")
		} else {
			fmt.Printf("got value from etcd, key = %s value = %s\n", etcdkey, value)
			fmt.Fprint(w, value)
		}
	case "setbuildoutput":
		etcdkey = "/kless/handlerbuildoutput/" + handler
		fmt.Printf("put value to etcd, key = %s\n", etcdkey)
		value, err := ioutil.ReadAll(r.Body)
		if nil != err {
			log.Fatal("Unable to read request body")
		}

		err = e.SetValue(etcdkey, string(value))
		if nil != err {
			log.Fatal("Unable to set value")
		}

		fmt.Fprintf(w, "{ status: OK }")
	case "getbuildoutput":
		etcdkey = "/kless/handlerbuildoutput/" + handler
		fmt.Printf("get value from etcd, key = %s\n", etcdkey)
		value, err := e.GetValue(etcdkey)
		if nil != err {
			log.Fatal("Unable to get value")
			fmt.Fprint(w, "{ status : failure }")
		} else {
			fmt.Printf("got value from etcd, key = %s value = %s\n", etcdkey, value)
			fmt.Fprint(w, value)
		}
	case "getsource":
		etcdkey = "/kless/source/" + key
		fmt.Printf("get source from etcd, key = %s\n", etcdkey)
		source, err := e.GetValue(etcdkey)
		if nil != err {
			log.Fatal("Unable to get source")
			fmt.Fprint(w, "{ status : failure }")
		} else {
			fmt.Printf("got source from etcd, key = %s source = %s\n", etcdkey, source)
			fmt.Fprint(w, source)
		}
	default:
		fmt.Fprintf(w, "op query parameter must be one of: 'get', 'put', 'sethandlerstatus', 'gethandlerstatus', 'setbuildoutput', 'getbuildoutput', 'getsource'")
	}

}
