package cmd

import (
	"fmt"
	"io"
	"log"

	cmdutil "github.com/paalth/kless/pkg/kless/cmd/util"
	klessapi "github.com/paalth/kless/pkg/klessserver/grpc"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//NewCmdGetFrontends gets the frontends
func NewCmdGetFrontends(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frontends",
		Short: "Get event handler frontends",
		Run: func(cmd *cobra.Command, args []string) {
			RunGetFrontends(f, out, cmd)
		},
	}

	return cmd
}

//RunGetFrontends gets the frontends from the server
func RunGetFrontends(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	address, err := cmdutil.GetServerAddress(cmd)
	if err != nil {
		log.Fatalf("Unable to get server address: %v", err)
		return nil
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := klessapi.NewKlessAPIClient(conn)

	// Contact the server and print out its response.
	stream, err := c.GetEventHandlerFrontends(context.Background(), &klessapi.GetEventHandlerFrontendsRequest{Clientversion: "0.0.1-alpha1"})
	if err != nil {
		log.Fatalf("Could not get event handler frontends: %v", err)
		return nil
	}

	fmt.Fprintf(out, "%-15s %-15s\n", "NAME", "TYPE")
	for {
		eventHandlerFrontend, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not get event handler frontends: %v", err)
			break
		}
		fmt.Fprintf(out, "%-15s %-15s\n", eventHandlerFrontend.EventHandlerFrontendName, eventHandlerFrontend.EventHandlerFrontendType)
	}

	return nil
}
