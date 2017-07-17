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

// NewCmdDeleteEventHandlerBuilder deletes an event handler builder
func NewCmdDeleteEventHandlerBuilder(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "builder",
		Short: "Delete event handler builder",
		Run: func(cmd *cobra.Command, args []string) {
			RunDeleteEventHandlerBuilder(f, out, cmd)
		},
	}

	return cmd
}

// RunDeleteEventHandlerBuilder executes the event handler builder deletion
func RunDeleteEventHandlerBuilder(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	address, err := cmdutil.GetServerAddress(cmd)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
		return nil
	}

	eventHandlerBuilder := cmdutil.GetFlagString(cmd, "eventhandlerbuilder")

	fmt.Fprintf(out, "Delete event handler builder = %s\n", eventHandlerBuilder)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := klessapi.NewKlessAPIClient(conn)

	// Contact the server and print out its response.
	r, err := c.DeleteEventHandlerBuilder(context.Background(), &klessapi.DeleteEventHandlerBuilderRequest{EventHandlerBuilderName: eventHandlerBuilder})
	if err != nil {
		log.Fatalf("could not delete event handler builder: %v", err)
	}

	fmt.Fprintf(out, "Server response: %s\n", r.Response)

	return nil
}
