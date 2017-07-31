package apiserver

import (
	"fmt"

	klessapi "github.com/paalth/kless/pkg/klessserver/grpc"

	"golang.org/x/net/context"
)

//GetServerVersion retrieves the server version
func (s *APIServer) GetServerVersion(ctx context.Context, in *klessapi.GetServerVersionRequest) (*klessapi.GetServerVersionReply, error) {
	fmt.Printf("Entering GetServerVersion\n")

	fmt.Printf("Leaving GetServerVersion\n")

	return &klessapi.GetServerVersionReply{Serverversion: "0.0.1"}, nil
}
