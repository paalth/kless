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

//NewCmdDeleteEventHandlerFrontendType deletes a frontend type
func NewCmdDeleteEventHandlerFrontendType(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frontendtype",
		Short: "Delete event handler frontend type",
		Run: func(cmd *cobra.Command, args []string) {
			RunDeleteEventHandlerFrontendType(f, out, cmd)
		},
	}

	return cmd
}

//RunDeleteEventHandlerFrontendType deletes a frontend type on the server
func RunDeleteEventHandlerFrontendType(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	address, err := cmdutil.GetServerAddress(cmd)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
		return nil
	}

	eventHandlerNameFrontendType := cmdutil.GetFlagString(cmd, "frontendtype")

	fmt.Fprintf(out, "Delete event handler frontend type = %s\n", eventHandlerNameFrontendType)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := klessapi.NewKlessAPIClient(conn)

	// Contact the server and print out its response.
	r, err := c.DeleteEventHandlerFrontendType(context.Background(), &klessapi.DeleteEventHandlerFrontendTypeRequest{EventHandlerFrontendType: eventHandlerNameFrontendType})
	if err != nil {
		log.Fatalf("could not delete event handler frontend type: %v", err)
	}

	fmt.Fprintf(out, "Server response: %s\n", r.Response)

	return nil
}
