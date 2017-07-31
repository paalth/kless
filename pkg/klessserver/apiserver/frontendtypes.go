package apiserver

import (
	"fmt"
	"log"

	"github.com/paalth/kless/pkg/etcdinterface"
	klessapi "github.com/paalth/kless/pkg/klessserver/grpc"

	"golang.org/x/net/context"
)

//CreateEventHandlerFrontendType adds a new frontend type
func (s *APIServer) CreateEventHandlerFrontendType(ctx context.Context, in *klessapi.CreateEventHandlerFrontendTypeRequest) (*klessapi.CreateEventHandlerFrontendTypeReply, error) {
	fmt.Printf("Entering CreateEventHandlerFrontendType\n")

	fmt.Printf("Event handler frontend type = %s with URL = %s\n", in.EventHandlerFrontendType, in.EventHandlerFrontendTypeURL)

	etcdKey := "/kless/frontendtypes/" + in.EventHandlerFrontendType

	fmt.Printf("Adding event handler frontend type with key = %s\n", etcdKey)

	e := &etcdinterface.EtcdInterface{}
	e.SetValue(etcdKey, in.EventHandlerFrontendTypeURL)

	fmt.Printf("Leaving CreateEventHandlerFrontendType\n")

	return &klessapi.CreateEventHandlerFrontendTypeReply{Response: "OK"}, nil
}

//GetEventHandlerFrontendTypes retrieves the list of available frontend types
func (s *APIServer) GetEventHandlerFrontendTypes(in *klessapi.GetEventHandlerFrontendTypesRequest, stream klessapi.KlessAPI_GetEventHandlerFrontendTypesServer) error {
	fmt.Printf("Entering GetEventHandlerFrontendTypes\n")

	e := &etcdinterface.EtcdInterface{}

	frontendTypes, _ := e.GetKeysValuesFromPrefix("/kless/frontendtypes/")

	for key, value := range frontendTypes {
		fmt.Printf("Event handler frontend type = %s, URL = %s\n", key, value)

		stream.Send(&klessapi.FrontendTypeInformation{
			EventHandlerFrontendType:    key,
			EventHandlerFrontendTypeURL: value,
		})
	}

	fmt.Printf("Leaving GetEventHandlerFrontendTypes\n")

	return nil
}

//DeleteEventHandlerFrontendType removes an existing frontend type
func (s *APIServer) DeleteEventHandlerFrontendType(ctx context.Context, in *klessapi.DeleteEventHandlerFrontendTypeRequest) (*klessapi.DeleteEventHandlerFrontendTypeReply, error) {
	fmt.Printf("Entering DeleteEventHandlerFrontendType\n")

	fmt.Printf("Event handler frontent type = %s\n", in.EventHandlerFrontendType)

	etcdKey := "/kless/frontendtypes/" + in.EventHandlerFrontendType

	fmt.Printf("Delete event handler frontend type from etcd with key = %s\n", etcdKey)

	e := &etcdinterface.EtcdInterface{}
	err := e.Delete(etcdKey)
	if err != nil {
		log.Fatal(err)
		return &klessapi.DeleteEventHandlerFrontendTypeReply{Response: "Unable to delete event handler frontend type"}, nil
	}

	fmt.Printf("Leaving DeleteEventHandlerFrontendType\n")

	return &klessapi.DeleteEventHandlerFrontendTypeReply{Response: "OK"}, nil
}
