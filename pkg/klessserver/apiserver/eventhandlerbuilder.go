package apiserver

import (
	"fmt"
	"log"

	"github.com/paalth/kless/pkg/etcdinterface"
	klessapi "github.com/paalth/kless/pkg/klessserver/grpc"

	"golang.org/x/net/context"
)

//CreateEventHandlerBuilder adds a new event handler builder
func (s *APIServer) CreateEventHandlerBuilder(ctx context.Context, in *klessapi.CreateEventHandlerBuilderRequest) (*klessapi.CreateEventHandlerBuilderReply, error) {
	fmt.Printf("Entering CreateEventHandlerBuilder\n")

	fmt.Printf("Event handler builder name = %s with URL = %s\n", in.EventHandlerBuilderName, in.EventHandlerBuilderURL)

	//TODO: return error is event handler builder already exists

	etcdBuilderKey := "/kless/builders/" + in.EventHandlerBuilderName

	fmt.Printf("Event handler builder etcd key = %s\n", etcdBuilderKey)

	e := &etcdinterface.EtcdInterface{}
	e.SetValue(etcdBuilderKey, in.EventHandlerBuilderURL)

	if nil != in.EventHandlerBuilderInformation {
		for key, val := range in.EventHandlerBuilderInformation {
			etcdBuildInfoKey := "/kless/buildinfo/" + in.EventHandlerBuilderName + "/" + key

			fmt.Printf("Adding build info with etcd key = %s for builder = %s\n", etcdBuildInfoKey, in.EventHandlerBuilderName)

			e.SetValue(etcdBuildInfoKey, string(val))
		}
	}

	fmt.Printf("Leaving CreateEventHandlerBuilder\n")

	return &klessapi.CreateEventHandlerBuilderReply{Response: "OK"}, nil
}

//GetEventHandlerBuilders retrieves the list of available event handler builders
func (s *APIServer) GetEventHandlerBuilders(in *klessapi.GetEventHandlerBuildersRequest, stream klessapi.KlessAPI_GetEventHandlerBuildersServer) error {
	fmt.Printf("Entering GetEventHandlerBuilders\n")

	e := &etcdinterface.EtcdInterface{}

	builders, _ := e.GetKeysValuesFromPrefix("/kless/builders/")

	for key, value := range builders {
		fmt.Printf("Builder name = %s, URL: %s\n", key, value)
		stream.Send(&klessapi.EventHandlerBuilderInformation{
			EventHandlerBuilderName: key,
			EventHandlerBuilderURL:  value,
		})
	}

	fmt.Printf("Leaving GetEventHandlerBuilders\n")

	return nil
}

//DeleteEventHandlerBuilder removes an existing event handler builder
func (s *APIServer) DeleteEventHandlerBuilder(ctx context.Context, in *klessapi.DeleteEventHandlerBuilderRequest) (*klessapi.DeleteEventHandlerBuilderReply, error) {
	fmt.Printf("Entering DeleteEventHandlerBuilder\n")

	fmt.Printf("Event handler builder name = %s\n", in.EventHandlerBuilderName)

	etcdBuilderKey := "/kless/builders/" + in.EventHandlerBuilderName

	fmt.Printf("Event handler builder etcd key = %s\n", etcdBuilderKey)

	e := &etcdinterface.EtcdInterface{}
	err := e.Delete(etcdBuilderKey)
	if err != nil {
		log.Fatal(err)
		return &klessapi.DeleteEventHandlerBuilderReply{Response: "Unable to delete event handler"}, nil
	}

	buildInfo, err := e.GetKeysFromPrefix("/kless/buildinfo/" + in.EventHandlerBuilderName)

	if err != nil {
		log.Fatal(err)
		return &klessapi.DeleteEventHandlerBuilderReply{Response: "Unable to retrieve event handler builder info"}, nil
	}

	if nil != buildInfo {
		for _, key := range buildInfo {
			etcdBuildInfoKey := "/kless/buildinfo/" + in.EventHandlerBuilderName + key
			fmt.Printf("Removing build info from etcd with key = %s\n", etcdBuildInfoKey)

			err := e.Delete(etcdBuildInfoKey)
			if err != nil {
				log.Fatal(err)
				return &klessapi.DeleteEventHandlerBuilderReply{Response: "Unable to delete event handler build info"}, nil
			}

		}
	}

	fmt.Printf("Leaving DeleteEventHandlerBuilder\n")

	return &klessapi.DeleteEventHandlerBuilderReply{Response: "OK"}, nil
}
