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

//NewCmdGetEventHandlerBuilders gets the available event handler builders
func NewCmdGetEventHandlerBuilders(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "builders",
		Short: "Get event handler builders",
		Run: func(cmd *cobra.Command, args []string) {
			RunGetEventHandlerBuilders(f, out, cmd)
		},
	}

	return cmd
}

//RunGetEventHandlerBuilders gets the event handler builders from the server
func RunGetEventHandlerBuilders(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

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
	stream, err := c.GetEventHandlerBuilders(context.Background(), &klessapi.GetEventHandlerBuildersRequest{Clientversion: "0.0.1-alpha1"})
	if err != nil {
		log.Fatalf("Could not get event handler builders: %v", err)
		return nil
	}

	fmt.Fprintf(out, "%-25s %-75s %-50s\n", "NAME", "URL", "COMMENT")
	for {
		eventHandlerBuilder, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not get event handler builders: %v", err)
			break
		}
		fmt.Fprintf(out, "%-25s %-75s %-50s\n", eventHandlerBuilder.EventHandlerBuilderName, eventHandlerBuilder.EventHandlerBuilderURL, eventHandlerBuilder.Comment)
	}

	return nil
}
