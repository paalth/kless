package cmd

import (
	"fmt"
	"io"

	cmdutil "github.com/paalth/kless/pkg/kless/cmd/util"
	//	kubelessapi "github.com/paalth/kless/pkg/klessserver/grpc"

	"github.com/spf13/cobra"
)

//NewCmdUpdateEventHandlerBuilder updates event handler builder
func NewCmdUpdateEventHandlerBuilder(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "builder",
		Short: "Update event handler builder",
		Run: func(cmd *cobra.Command, args []string) {
			RunUpdateEventHandlerBuilder(f, out, cmd)
		},
	}
	cmd.Flags().StringP("url", "u", "", "Event handler builder URL")
	//cmd.Flags().MarkShorthandDeprecated("URL", "please use --URL instead")

	return cmd
}

//RunUpdateEventHandlerBuilder will update the event handler builder on the server
func RunUpdateEventHandlerBuilder(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	/*
		address, err := cmdutil.GetServerAddress(cmd)
		if err != nil {
			log.Fatalf("Unable to connect to server: %v", err)
			return nil
		}
	*/

	eventHandlerBuilder := cmdutil.GetFlagString(cmd, "eventhandlerbuilder")
	eventHandlerBuilderURL := cmdutil.GetFlagString(cmd, "url")

	fmt.Fprintf(out, "Update event handler builder = %s with URL = %s\n", eventHandlerBuilder, eventHandlerBuilderURL)
	/*
		// Set up a connection to the server.
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := kubelessapi.NewKubelessAPIClient(conn)

		// Contact the server and print out its response.
		r, err := c.CreateEventHandler(context.Background(), &kubelessapi.CreateEventHandlerRequest{EventHandlerName: eventHandlerName,
			EventHandlerNamespace: eventHandlerNamespace,
			EventHandlerBuilder:   "go",
			EventHandlerSource:    "sourceUrl",
			EventHandlerVersion:   "1.0",
			EventHanlderEventType: "http"})
		if err != nil {
			log.Fatalf("could not create event handler: %v", err)
		}

		fmt.Fprintf(out, "Server response: %s\n", r.Response)
	*/

	return nil
}
