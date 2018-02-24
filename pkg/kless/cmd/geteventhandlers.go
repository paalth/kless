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

//NewCmdGetEventHandlers gets the event handlers
func NewCmdGetEventHandlers(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handlers",
		Short: "Get event handlers",
		Run: func(cmd *cobra.Command, args []string) {
			RunGetEventHandlers(f, out, cmd)
		},
	}

	return cmd
}

//RunGetEventHandlers gets the event handlers from the server
func RunGetEventHandlers(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	address, err := cmdutil.GetServerAddress(cmd)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
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
	stream, err := c.GetEventHandlers(context.Background(), &klessapi.GetEventHandlersRequest{Clientversion: "0.0.1-alpha1"})
	if err != nil {
		log.Fatalf("Could not get event handlers: %v", err)
		return nil
	}

	fmt.Fprintf(out, "%-25s %-15s %-10s %-15s %-15s %-15s %-17s %-50s\n", "NAME", "NAMESPACE", "VERSION", "BUILDER", "FRONTEND", "BUILD STATUS", "HANDLER AVAILABLE", "COMMENT")
	for {
		eventHandler, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not get event handlers: %v", err)
			break
		}
		fmt.Fprintf(out, "%-25s %-15s %-10s %-15s %-15s %-15s %-17s %-50s\n", eventHandler.EventHandlerName, eventHandler.EventHandlerNamespace, eventHandler.EventHandlerVersion, eventHandler.EventHandlerBuilder, eventHandler.Frontend, eventHandler.BuildStatus, eventHandler.EventHandlerAvailable, eventHandler.Comment)
	}

	return nil
}
