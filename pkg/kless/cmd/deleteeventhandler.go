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

// NewCmdDeleteEventHandler deletes an event handler
func NewCmdDeleteEventHandler(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handler",
		Short: "Delete event handler",
		Run: func(cmd *cobra.Command, args []string) {
			RunDeleteEventHandler(f, out, cmd)
		},
	}

	return cmd
}

// RunDeleteEventHandler executes the event handler deletion
func RunDeleteEventHandler(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	address, err := cmdutil.GetServerAddress(cmd)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
		return nil
	}

	eventHandlerName := cmdutil.GetFlagString(cmd, "eventhandler")
	eventHandlerNamespace := cmdutil.GetFlagString(cmd, "namespace")

	fmt.Fprintf(out, "Delete event handler name = %s in namespace = %s\n", eventHandlerName, eventHandlerNamespace)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := klessapi.NewKlessAPIClient(conn)

	// Contact the server and print out its response.
	r, err := c.DeleteEventHandler(context.Background(), &klessapi.DeleteEventHandlerRequest{EventHandlerName: eventHandlerName,
		EventHandlerNamespace: eventHandlerNamespace})
	if err != nil {
		log.Fatalf("could not delete event handler: %v", err)
	}

	fmt.Fprintf(out, "Server response: %s\n", r.Response)

	return nil
}
