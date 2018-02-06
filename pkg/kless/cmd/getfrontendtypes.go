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

//NewCmdGetFrontendTypes gets the frontend types
func NewCmdGetFrontendTypes(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frontendtypes",
		Short: "Get event handler frontend types",
		Run: func(cmd *cobra.Command, args []string) {
			RunGetFrontendTypes(f, out, cmd)
		},
	}

	return cmd
}

//RunGetFrontendTypes gets the frontend types from the server
func RunGetFrontendTypes(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

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
	stream, err := c.GetEventHandlerFrontendTypes(context.Background(), &klessapi.GetEventHandlerFrontendTypesRequest{Clientversion: "0.0.1-alpha1"})
	if err != nil {
		log.Fatalf("Could not get event handler frontend types: %v", err)
		return nil
	}

	fmt.Fprintf(out, "%-20s %-75s %-50s\n", "TYPE", "URL", "COMMENT")
	for {
		eventHandlerFrontendType, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not get event handler frontend types: %v", err)
			break
		}
		fmt.Fprintf(out, "%-20s %-75s %-50s\n", eventHandlerFrontendType.EventHandlerFrontendType, eventHandlerFrontendType.EventHandlerFrontendTypeURL, eventHandlerFrontendType.Comment)
	}

	return nil
}
