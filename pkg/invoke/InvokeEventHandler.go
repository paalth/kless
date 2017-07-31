package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	evt "eventhandler"

	"github.com/paalth/kless/pkg/interface/klessgo"
)

func main() {
	fmt.Printf("Starting event handler...\n")

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Processing incoming request\n")

	context := getContext()

	req := klessgo.Request{
		Headers: getRequestHeaders(r),
		Body:    getRequestBody(r),
	}

	resp := klessgo.Response{}

	handler := evt.EventHandler{}

	invokeEventHandler(handler, context, &resp, &req)

	writeResponse(w, resp)

	fmt.Printf("Incoming request processed\n")
}

func getContext() (c *klessgo.Context) {

	c = &klessgo.Context{}

	eventHandlerNamespace := os.Getenv("KLESS_EVENT_HANDLER_NAMESPACE")
	eventHandlerName := os.Getenv("KLESS_EVENT_HANDLER_NAME")
	eventHandlerVersion := os.Getenv("KLESS_EVENT_HANDLER_VERSION")

	fmt.Printf("Event handler name = %s, namespace = %s, version = %s\n", eventHandlerName, eventHandlerNamespace, eventHandlerVersion)

	c.Info = make(map[string]string)

	c.Info["namespace"] = eventHandlerNamespace
	c.Info["name"] = eventHandlerName
	c.Info["version"] = eventHandlerVersion

	return c
}

func getRequestHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)

	headers["kless-method"] = r.Method
	headers["kless-url"] = r.URL.String()
	headers["kless-proto"] = r.Proto
	for name, httpHeaders := range r.Header {
		for _, h := range httpHeaders {
			headers[name] = h
		}
	}

	return headers
}

func getRequestBody(r *http.Request) []byte {

	bodyBuffer, err := ioutil.ReadAll(r.Body)
	if nil != err {
		log.Fatal(err)
	}

	return bodyBuffer
}

func invokeEventHandler(eventHandler klessgo.KlessHandler, c *klessgo.Context, resp *klessgo.Response, req *klessgo.Request) {
	eventHandler.Handler(c, resp, req)
}

func writeResponse(w http.ResponseWriter, resp klessgo.Response) error {
	if nil != resp.Headers {
		for headerKey, headerValue := range resp.Headers {
			w.Header().Add(headerKey, headerValue)
		}
	}

	if nil != resp.Body {
		fmt.Fprint(w, string(resp.Body))
	}

	return nil
}
