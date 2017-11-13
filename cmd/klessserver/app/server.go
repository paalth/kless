package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/paalth/kless/pkg/etcdinterface"
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

func etcdHandler(w http.ResponseWriter, r *http.Request) {
	op := r.URL.Query().Get("op")
	key := r.URL.Query().Get("key")
	builder := r.URL.Query().Get("builder")

	fmt.Printf("etcd handler, op = %s builder = %s key = %s\n", op, builder, key)

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
		fmt.Fprintf(w, "op query parameter must be one of: 'get', 'put', 'getsource'")
	}

}
