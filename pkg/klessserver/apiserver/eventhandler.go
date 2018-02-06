package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/paalth/kless/pkg/etcdinterface"
	"github.com/paalth/kless/pkg/influxdbinterface"
	klessapi "github.com/paalth/kless/pkg/klessserver/grpc"
	klesshandlers "github.com/paalth/kless/pkg/klessserver/servicehandler"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

//CreateEventHandler adds a new event handler
func (s *APIServer) CreateEventHandler(ctx context.Context, in *klessapi.CreateEventHandlerRequest) (*klessapi.CreateEventHandlerReply, error) {
	fmt.Printf("Entering CreateEventHandler\n")

	fmt.Printf("Event handler name = %s in namespace %s using event handler builder %s, comment = %s\n", in.EventHandlerName, in.EventHandlerNamespace, in.EventHandlerBuilder, in.Comment)

	etcdBuilderKey := "/kless/builders/" + in.EventHandlerBuilder

	fmt.Printf("Get event handler builder URL for builder = %s\n", etcdBuilderKey)

	e := &etcdinterface.EtcdInterface{}

	eventHandlerBuilderURL, err := e.GetValue(etcdBuilderKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Event handler builder URL = %s for builder = %s\n", eventHandlerBuilderURL, etcdBuilderKey)

	if eventHandlerBuilderURL == "" {
		return &klessapi.CreateEventHandlerReply{Response: "Event handler builder not found"}, nil
	}

	eventHandlerID := uuid.NewV4().String()

	etcdSourceKey := "/kless/source/" + eventHandlerID

	eventHandlerSourceCode := in.EventHandlerSourceCode

	if len(eventHandlerSourceCode) == 0 {
		resp, err := http.Get(in.EventHandlerSourceCodeURL)
		if nil != err {
			log.Fatal(err)
			return &klessapi.CreateEventHandlerReply{Response: "HTTP get failed"}, nil
		}
		defer resp.Body.Close()

		eventHandlerSourceCode, err = ioutil.ReadAll(resp.Body)
		if nil != err {
			log.Fatal(err)
			return &klessapi.CreateEventHandlerReply{Response: "Could not read HTTP response"}, nil
		}
	}

	sourceCode := string(eventHandlerSourceCode)

	fmt.Printf("Adding event handler source code to etcd with key = %s, source code:\n%s\n", etcdSourceKey, sourceCode)

	e.SetValue(etcdSourceKey, sourceCode)

	handler := &klesshandlers.ServiceHandler{}

	eventHandlerInfo := &klesshandlers.EventHandlerInfo{
		ID:                     eventHandlerID,
		Name:                   in.EventHandlerName,
		Namespace:              in.EventHandlerNamespace,
		Version:                in.EventHandlerVersion,
		EventHandlerBuilder:    in.EventHandlerBuilder,
		EventHandlerBuilderURL: eventHandlerBuilderURL,
		Frontend:               in.EventHandlerFrontend,
		DependenciesURL:        in.EventHandlerDependenciesURL,
		Comment:                in.Comment,
	}

	etcdFrontendKey := "/kless/frontend/" + in.EventHandlerFrontend

	fmt.Printf("Getting event handler frontend from etcd with key = %s\n", etcdFrontendKey)

	eventHandlerFrontendInfoJSON, err := e.GetValue(etcdFrontendKey)

	if nil != err {
		return &klessapi.CreateEventHandlerReply{Response: "Unable to retrieve event handler frontend"}, nil
	}

	if eventHandlerFrontendInfoJSON == "" {
		return &klessapi.CreateEventHandlerReply{Response: "Event handler frontend not found"}, nil
	}

	eventHandlerFrontendInfo := &klesshandlers.EventHandlerFrontendInfo{}

	err = json.Unmarshal([]byte(eventHandlerFrontendInfoJSON), &eventHandlerFrontendInfo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Event handler frontend name = %s has event handler frontend type %s\n", eventHandlerFrontendInfo.Name, eventHandlerFrontendInfo.Type)

	etcdTypeKey := "/kless/frontendtypes/" + eventHandlerFrontendInfo.Type

	fmt.Printf("Get event handler frontend type repository URL from etcd with key = %s\n", etcdTypeKey)

	eventHandlerFrontendTypeURL, err := e.GetValue(etcdTypeKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Event handler frontent type URL = %s for frontend with key = %s\n", eventHandlerFrontendTypeURL, etcdTypeKey)

	if eventHandlerFrontendTypeURL == "" {
		return &klessapi.CreateEventHandlerReply{Response: "Event handler frontend type not found"}, nil
	}

	handler.CreateEventHandler(eventHandlerInfo, eventHandlerFrontendInfo, eventHandlerFrontendTypeURL)

	etcdHandlerKey := "/kless/handlers/" + in.EventHandlerName

	fmt.Printf("Adding event handler to etcd with key = %s\n", etcdHandlerKey)

	eventHandlerInfoJSON, err := json.Marshal(eventHandlerInfo)
	if err != nil {
		log.Fatal(err)
	}

	e.SetValue(etcdHandlerKey, string(eventHandlerInfoJSON))

	fmt.Printf("Leaving CreateEventHandler\n")

	return &klessapi.CreateEventHandlerReply{Response: "OK"}, nil
}

//GetEventHandlers retrieves a list of all defined event handlers
func (s *APIServer) GetEventHandlers(in *klessapi.GetEventHandlersRequest, stream klessapi.KlessAPI_GetEventHandlersServer) error {
	fmt.Printf("Entering GetEventHandlers\n")

	e := &etcdinterface.EtcdInterface{}

	buildersJSON, _ := e.GetValuesFromPrefix("/kless/handlers/")

	for i := 0; i < len(buildersJSON); i++ {
		eventHandlerInfo := klesshandlers.EventHandlerInfo{}
		err := json.Unmarshal([]byte(buildersJSON[i]), &eventHandlerInfo)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Builder: %s\n", eventHandlerInfo.Name)

		stream.Send(&klessapi.EventHandlerInformation{
			EventHandlerId:         eventHandlerInfo.ID,
			EventHandlerName:       eventHandlerInfo.Name,
			EventHandlerNamespace:  eventHandlerInfo.Namespace,
			EventHandlerVersion:    eventHandlerInfo.Version,
			EventHandlerBuilder:    eventHandlerInfo.EventHandlerBuilder,
			EventHandlerBuilderURL: eventHandlerInfo.EventHandlerBuilderURL,
			Frontend:               eventHandlerInfo.Frontend,
			Comment:                eventHandlerInfo.Comment,
		})
	}

	fmt.Printf("Leaving GetEventHandlers\n")

	return nil
}

//GetEventHandlerStatistics retrieves event handler statistics, for now all stats are retrieved which needs to be fixed
func (s *APIServer) GetEventHandlerStatistics(in *klessapi.GetEventHandlerStatisticsRequest, stream klessapi.KlessAPI_GetEventHandlerStatisticsServer) error {
	fmt.Printf("Entering GetEventHandlerStatistics\n")

	i := &influxdbinterface.InfluxdbInterface{}

	events, _ := i.GetEvents()
	for _, e := range events {
		stream.Send(&klessapi.EventHandlerStatisticsInformation{
			Timestamp:             e.Timestamp,
			EventHandlerName:      e.EventHandler,
			EventHandlerNamespace: e.Namespace,
			EventHandlerVersion:   e.Version,
			PodName:               e.PodName,
			RequestSize:           e.RequestSize,
			ResponseSize:          e.ResponseSize,
			ResponseTime:          e.ResponseTime,
		})
	}

	fmt.Printf("Leaving GetEventHandlerStatistics\n")

	return nil
}

//DeleteEventHandler removes an event handler
func (s *APIServer) DeleteEventHandler(ctx context.Context, in *klessapi.DeleteEventHandlerRequest) (*klessapi.DeleteEventHandlerReply, error) {
	fmt.Printf("Entering DeleteEventHandler\n")

	fmt.Printf("Event handler name = %s in namespace %s\n", in.EventHandlerName, in.EventHandlerNamespace)

	handler := &klesshandlers.ServiceHandler{}

	eventHandlerInfo := &klesshandlers.EventHandlerInfo{
		Name:      in.EventHandlerName,
		Namespace: in.EventHandlerNamespace,
		Comment:   in.Comment,
	}

	handler.DeleteEventHandler(eventHandlerInfo)

	etcdHandlerKey := "/kless/handlers/" + in.EventHandlerName

	fmt.Printf("Delete event handler from etcd with key = %s\n", etcdHandlerKey)

	e := &etcdinterface.EtcdInterface{}
	err := e.Delete(etcdHandlerKey)
	if err != nil {
		log.Fatal(err)
		return &klessapi.DeleteEventHandlerReply{Response: "Unable to delete event handler"}, nil
	}

	fmt.Printf("Leaving DeleteEventHandler\n")

	return &klessapi.DeleteEventHandlerReply{Response: "OK"}, nil
}
