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

//NewCmdCreateFrontendType creates a new frontend type
func NewCmdCreateFrontendType(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frontendtype",
		Short: "Create event handler frontend type",
		Run: func(cmd *cobra.Command, args []string) {
			RunCreateFrontendType(f, out, cmd)
		},
	}
	cmd.Flags().StringP("url", "u", "", "URL of event handler source code")

	return cmd
}

//RunCreateFrontendType creates a new frontend type on the server
func RunCreateFrontendType(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	address, err := cmdutil.GetServerAddress(cmd)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
		return nil
	}

	eventHandlerFrontendType := cmdutil.GetFlagString(cmd, "frontendtype")
	eventHandlerFrontendTypeURL := cmdutil.GetFlagString(cmd, "url")

	fmt.Fprintf(out, "Create event handler frontend type = %s with repository URL = %s\n", eventHandlerFrontendType, eventHandlerFrontendTypeURL)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect to server: %v", err)
	}
	defer conn.Close()
	c := klessapi.NewKlessAPIClient(conn)

	// Contact the server and print out its response.
	r, err := c.CreateEventHandlerFrontendType(context.Background(), &klessapi.CreateEventHandlerFrontendTypeRequest{EventHandlerFrontendType: eventHandlerFrontendType,
		EventHandlerFrontendTypeURL: eventHandlerFrontendTypeURL,
	})
	if err != nil {
		log.Fatalf("Could not create event handler frontend type: %v", err)
	}

	fmt.Fprintf(out, "Server response: %s\n", r.Response)

	return nil
}
