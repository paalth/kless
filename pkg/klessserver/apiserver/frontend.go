package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/paalth/kless/pkg/etcdinterface"
	klessapi "github.com/paalth/kless/pkg/klessserver/grpc"
	klesshandlers "github.com/paalth/kless/pkg/klessserver/servicehandler"

	"golang.org/x/net/context"
)

//CreateEventHandlerFrontend adds a new event handler frontend
func (s *APIServer) CreateEventHandlerFrontend(ctx context.Context, in *klessapi.CreateEventHandlerFrontendRequest) (*klessapi.CreateEventHandlerFrontendReply, error) {
	fmt.Printf("Entering CreateEventHandlerFrontend\n")

	fmt.Printf("Event handler frontend name = %s using event handler frontend type %s, comment = %s\n", in.EventHandlerFrontendName, in.EventHandlerFrontendType, in.Comment)

	etcdTypeKey := "/kless/frontendtypes/" + in.EventHandlerFrontendType

	fmt.Printf("Get event handler frontend type repository URL from etcd with key = %s\n", etcdTypeKey)

	e := &etcdinterface.EtcdInterface{}

	eventHandlerFrontendTypeURL, err := e.GetValue(etcdTypeKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Event handler frontent type URL = %s for frontend with key = %s\n", eventHandlerFrontendTypeURL, etcdTypeKey)

	if eventHandlerFrontendTypeURL == "" {
		return &klessapi.CreateEventHandlerFrontendReply{Response: "Event handler frontend type not found"}, nil
	}

	var frontendInformation map[string]string

	if nil != in.EventHandlerFrontendInformation {
		frontendInformation = make(map[string]string)

		for k, v := range in.EventHandlerFrontendInformation {
			fmt.Printf("Frontend information: key = %s, value = %s\n", k, v)
			frontendInformation[k] = v
		}
	}

	eventHandlerFrontendInfo := &klesshandlers.EventHandlerFrontendInfo{
		Name:        in.EventHandlerFrontendName,
		Type:        in.EventHandlerFrontendType,
		Information: frontendInformation,
		Comment:     in.Comment,
	}

	etcdFrontendKey := "/kless/frontend/" + in.EventHandlerFrontendName

	fmt.Printf("Adding event handler frontend to etcd with key = %s\n", etcdFrontendKey)

	eventHandlerFrontendInfoJSON, err := json.Marshal(eventHandlerFrontendInfo)
	if err != nil {
		log.Fatal(err)
	}

	e.SetValue(etcdFrontendKey, string(eventHandlerFrontendInfoJSON))

	fmt.Printf("Leaving CreateEventHandlerFrontend\n")

	return &klessapi.CreateEventHandlerFrontendReply{Response: "OK"}, nil
}

//GetEventHandlerFrontends retrieves list of available frontends
func (s *APIServer) GetEventHandlerFrontends(in *klessapi.GetEventHandlerFrontendsRequest, stream klessapi.KlessAPI_GetEventHandlerFrontendsServer) error {
	fmt.Printf("Entering GetEventHandlerFrontends\n")

	e := &etcdinterface.EtcdInterface{}

	frontendJSON, _ := e.GetValuesFromPrefix("/kless/frontend/")

	for i := 0; i < len(frontendJSON); i++ {
		eventHandlerFrontendInfo := klesshandlers.EventHandlerFrontendInfo{}
		err := json.Unmarshal([]byte(frontendJSON[i]), &eventHandlerFrontendInfo)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Frontend: name = %s, type = %s, comment = %s\n", eventHandlerFrontendInfo.Name, eventHandlerFrontendInfo.Type, eventHandlerFrontendInfo.Comment)

		stream.Send(&klessapi.FrontendInformation{
			EventHandlerFrontendName: eventHandlerFrontendInfo.Name,
			EventHandlerFrontendType: eventHandlerFrontendInfo.Type,
			Comment:                  eventHandlerFrontendInfo.Comment,
		})
	}

	fmt.Printf("Leaving GetEventHandlerFrontends\n")

	return nil
}

//DeleteEventHandlerFrontend removes an existing frontend
func (s *APIServer) DeleteEventHandlerFrontend(ctx context.Context, in *klessapi.DeleteEventHandlerFrontendRequest) (*klessapi.DeleteEventHandlerFrontendReply, error) {
	fmt.Printf("Entering DeleteEventHandlerFrontend\n")

	fmt.Printf("Event handler frontend %s\n", in.EventHandlerFrontendName)

	etcdKey := "/kless/frontend/" + in.EventHandlerFrontendName

	fmt.Printf("Delete event handler frontend from etcd with key = %s\n", etcdKey)

	e := &etcdinterface.EtcdInterface{}
	err := e.Delete(etcdKey)
	if err != nil {
		log.Fatal(err)
		return &klessapi.DeleteEventHandlerFrontendReply{Response: "Unable to delete event handler frontend"}, nil
	}

	fmt.Printf("Leaving DeleteEventHandlerFrontend\n")

	return &klessapi.DeleteEventHandlerFrontendReply{Response: "OK"}, nil
}

//SubmitEventHandlerFrontendRequest submits a request to a frontend for test purposes
func (s *APIServer) SubmitEventHandlerFrontendRequest(ctx context.Context, in *klessapi.SubmitEventHandlerFrontendRequestRequest) (*klessapi.SubmitEventHandlerFrontendRequestReply, error) {

	e := &etcdinterface.EtcdInterface{}

	frontendJSON, _ := e.GetValue("/kless/frontend/" + in.EventHandlerFrontendName)

	eventHandlerFrontendInfo := klesshandlers.EventHandlerFrontendInfo{}
	err := json.Unmarshal([]byte(frontendJSON), &eventHandlerFrontendInfo)
	if err != nil {
		log.Fatal(err)
	}

	status := "N/A"
	var response []byte

	buildersJSON, _ := e.GetValuesFromPrefix("/kless/handlers/")

	for i := 0; i < len(buildersJSON); i++ {
		eventHandlerInfo := klesshandlers.EventHandlerInfo{}
		err := json.Unmarshal([]byte(buildersJSON[i]), &eventHandlerInfo)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Checking for matching frontend with handler: %s\n", eventHandlerInfo.Name)

		if eventHandlerInfo.Frontend == in.EventHandlerFrontendName {
			fmt.Printf("Found handler matching frontend name\n")

			handler := &klesshandlers.ServiceHandler{}

			url := "http://" + handler.GetServiceName(&eventHandlerInfo) + "." + eventHandlerInfo.Namespace + ":" + eventHandlerFrontendInfo.Information["KLESS_FRONTEND_HTTP_PORT"]

			fmt.Printf("URL = %s\n", url)

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(in.Request))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))

			status = "OK"
			response = body
		}
	}

	return &klessapi.SubmitEventHandlerFrontendRequestReply{Status: status, Response: response}, nil
}
