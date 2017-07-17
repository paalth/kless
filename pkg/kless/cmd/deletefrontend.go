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

//NewCmdDeleteEventHandlerFrontend deletes a frontend
func NewCmdDeleteEventHandlerFrontend(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frontend",
		Short: "Delete event handler frontend",
		Run: func(cmd *cobra.Command, args []string) {
			RunDeleteEventHandlerFrontend(f, out, cmd)
		},
	}

	return cmd
}

//RunDeleteEventHandlerFrontend deletes a frontend on the server
func RunDeleteEventHandlerFrontend(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	address, err := cmdutil.GetServerAddress(cmd)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
		return nil
	}

	eventHandlerNameFrontend := cmdutil.GetFlagString(cmd, "frontend")

	fmt.Fprintf(out, "Delete event handler frontend = %s\n", eventHandlerNameFrontend)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := klessapi.NewKlessAPIClient(conn)

	// Contact the server and print out its response.
	r, err := c.DeleteEventHandlerFrontend(context.Background(), &klessapi.DeleteEventHandlerFrontendRequest{EventHandlerFrontendName: eventHandlerNameFrontend})
	if err != nil {
		log.Fatalf("could not delete event handler frontend: %v", err)
	}

	fmt.Fprintf(out, "Server response: %s\n", r.Response)

	return nil
}
