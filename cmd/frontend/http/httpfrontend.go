package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"os"

	"github.com/paalth/kless/pkg/influxdbinterface"
)

//HTTPFrontendProxy contains information needed by the proxy
type HTTPFrontendProxy struct {
	target       *url.URL
	proxy        *httputil.ReverseProxy
	influxdb     *influxdbinterface.InfluxdbInterface
	namespace    string
	eventHandler string
	version      string
	podName      string
}

//NewProxy creates a new proxy
func NewProxy(target string, ns string, handler string, version string, pod string) *HTTPFrontendProxy {
	url, err := url.Parse(target)
	if nil != err {
		fmt.Printf("Unable to parse URL, %v", err)
	}

	return &HTTPFrontendProxy{target: url,
		proxy:        httputil.NewSingleHostReverseProxy(url),
		influxdb:     &influxdbinterface.InfluxdbInterface{},
		namespace:    ns,
		eventHandler: handler,
		version:      version,
		podName:      pod,
	}
}

func (p *HTTPFrontendProxy) handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Before proxy")

	reqSize := "0"

	reqSizeHeader := r.Header.Get("Content-Length")
	if "" != reqSizeHeader {
		reqSize = reqSizeHeader
	}

	startTime := time.Now()

	p.proxy.ServeHTTP(w, r)

	respTime := time.Since(startTime).Nanoseconds() / 1e6

	respSize := w.Header().Get("Content-Length")

	fmt.Printf("After proxy, request size = %s, response size = %s, response time = %v\n", reqSize, respSize, respTime)

	requestSize, _ := strconv.ParseInt(reqSize, 10, 64)
	responseSize, _ := strconv.ParseInt(respSize, 10, 64)

	klessEvent := &influxdbinterface.KlessEvent{
		Namespace:    p.namespace,
		EventHandler: p.eventHandler,
		Version:      p.version,
		PodName:      p.podName,
		RequestSize:  requestSize,
		ResponseSize: responseSize,
		ResponseTime: respTime,
	}

	fmt.Println("Adding event to influxdb")

	err := p.influxdb.AddEvent(klessEvent)
	if nil != err {
		fmt.Printf("Error: %v", err)
	}
}

func main() {
	eventHandler := os.Getenv("KLESS_EVENT_HANDLER_NAME")
	namespace := os.Getenv("KLESS_EVENT_HANDLER_NAMESPACE")
	version := os.Getenv("KLESS_EVENT_HANDLER_VERSION")

	httpPort := os.Getenv("KLESS_FRONTEND_HTTP_PORT")

	listenAddress := "0.0.0.0:" + httpPort
	targetAddress := "http://localhost:8080"

	podName, err := os.Hostname()
	if nil != err {
		fmt.Printf("Unable to retrieve pod name, error = %v. Setting pod name to 'default'", err)
		podName = "default"
	}

	fmt.Printf("Starting HTTP frontend, namespace = %s, event handler name = %s, version = %s, pod name = %s, listen address = %s, target address = %s\n", namespace, eventHandler, version, podName, listenAddress, targetAddress)

	proxy := NewProxy(targetAddress, namespace, eventHandler, version, podName)

	http.HandleFunc("/", proxy.handle)
	http.ListenAndServe(listenAddress, nil)
}
