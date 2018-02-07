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

//NewCmdDescribeEventHandler gets event handler details
func NewCmdDescribeEventHandler(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handler",
		Short: "Get detailed event handler information",
		Run: func(cmd *cobra.Command, args []string) {
			RunDescribeEventHandler(f, out, cmd)
		},
	}
	cmd.Flags().StringP("version", "v", "1.0", "Event handler version")

	return cmd
}

//RunDescribeEventHandler gets the event handler details from the server
func RunDescribeEventHandler(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	eventHandlerName := cmdutil.GetFlagString(cmd, "eventhandler")
	eventHandlerNamespace := cmdutil.GetFlagString(cmd, "namespace")
	eventHandlerVersion := cmdutil.GetFlagString(cmd, "version")

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
	r, err := c.DescribeEventHandler(context.Background(), &klessapi.DescribeEventHandlerRequest{Clientversion: "0.0.1-alpha1", EventHandlerName: eventHandlerName, EventHandlerNamespace: eventHandlerNamespace, EventHandlerVersion: eventHandlerVersion})
	if err != nil {
		log.Fatalf("Could not describe event handler: %v", err)
		return nil
	}

	if "OK" == r.Response {
		fmt.Fprintf(out, "ID:          %s\n", r.EventHandlerInformation.EventHandlerId)
		fmt.Fprintf(out, "Name:        %s\n", r.EventHandlerInformation.EventHandlerName)
		fmt.Fprintf(out, "Namespace:   %s\n", r.EventHandlerInformation.EventHandlerNamespace)
		fmt.Fprintf(out, "Version:     %s\n", r.EventHandlerInformation.EventHandlerVersion)
		fmt.Fprintf(out, "Builder:     %s\n", r.EventHandlerInformation.EventHandlerBuilder)
		fmt.Fprintf(out, "Builder URL: %s\n", r.EventHandlerInformation.EventHandlerBuilderURL)
		fmt.Fprintf(out, "Frontend:    %s\n", r.EventHandlerInformation.Frontend)
		fmt.Fprintf(out, "Comment:     %s\n", r.EventHandlerInformation.Comment)
		fmt.Fprintf(out, "Status:      %s\n", r.EventHandlerInformation.Status)

		sourceCode := string(r.SourceCode)
		fmt.Fprintf(out, "Source code:\n\n")
		fmt.Fprintf(out, "%s\n", sourceCode)

		buildOutput := string(r.BuildOutput)
		fmt.Fprintf(out, "Build output:\n\n")
		fmt.Fprintf(out, "%s\n", buildOutput)
	} else {
		fmt.Fprintf(out, "%s\n", r.Response)
	}

	return nil
}
