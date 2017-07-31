package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	cmdutil "github.com/paalth/kless/pkg/kless/cmd/util"
	klessapi "github.com/paalth/kless/pkg/klessserver/grpc"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// NewCmdCreateEventHandlerBuilder enables registration of a new event handler builder on the server
func NewCmdCreateEventHandlerBuilder(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "builder",
		Short: "Create event handler builder",
		Run: func(cmd *cobra.Command, args []string) {
			RunCreateEventHandlerBuilder(f, out, cmd)
		},
	}
	cmd.Flags().StringP("url", "u", "", "Event handler builder URL")
	cmd.Flags().StringP("build-information", "i", "", "Event handler builder information (source code etc.) in the form key1=filename1,key2=filename2...")

	return cmd
}

// RunCreateEventHandlerBuilder executes the registration of a new event handler builder on the kless server
func RunCreateEventHandlerBuilder(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	address, err := cmdutil.GetServerAddress(cmd)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
		return nil
	}

	eventHandlerBuilder := cmdutil.GetFlagString(cmd, "eventhandlerbuilder")
	eventHandlerBuilderURL := cmdutil.GetFlagString(cmd, "url")
	buildInformation := cmdutil.GetFlagString(cmd, "build-information")

	if eventHandlerBuilder == "" {
		fmt.Printf("builder flag is required\n")
		return nil
	}
	if eventHandlerBuilderURL == "" {
		fmt.Printf("url flag is required\n")
		return nil
	}

	var eventHandlerBuilderInformation map[string][]byte

	if "" != buildInformation {
		eventHandlerBuilderInformation = make(map[string][]byte)

		buildInfoStrings := strings.Split(buildInformation, ",")

		for _, buildInfoString := range buildInfoStrings {
			buildInfo := strings.Split(buildInfoString, "=")

			fileContent, err := ioutil.ReadFile(buildInfo[1])
			if nil != err {
				log.Fatalf("could not read file %s: %v", buildInfo[1], err)
				return nil
			}

			eventHandlerBuilderInformation[buildInfo[0]] = fileContent
		}
	}

	fmt.Fprintf(out, "Create event handler builder = %s with URL = %s\n", eventHandlerBuilder, eventHandlerBuilderURL)
	if nil != eventHandlerBuilderInformation {
		fmt.Fprint(out, "Additional information provided: ")
		for k := range eventHandlerBuilderInformation {
			fmt.Fprintf(out, "%s ", k)
		}
		fmt.Fprint(out, "\n")
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := klessapi.NewKlessAPIClient(conn)

	// Contact the server and print out its response.
	r, err := c.CreateEventHandlerBuilder(context.Background(), &klessapi.CreateEventHandlerBuilderRequest{EventHandlerBuilderName: eventHandlerBuilder,
		EventHandlerBuilderURL:         eventHandlerBuilderURL,
		EventHandlerBuilderInformation: eventHandlerBuilderInformation,
	})
	if err != nil {
		log.Fatalf("could not create event handler builder: %v", err)
	}

	fmt.Fprintf(out, "Server response: %s\n", r.Response)

	return nil
}
