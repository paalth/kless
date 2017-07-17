package cmd

import (
	"fmt"
	"io"
	"log"
	"strings"

	cmdutil "github.com/paalth/kless/pkg/kless/cmd/util"
	klessapi "github.com/paalth/kless/pkg/klessserver/grpc"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//NewCmdCreateFrontend creates a frontend
func NewCmdCreateFrontend(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frontend",
		Short: "Create event handler frontend",
		Run: func(cmd *cobra.Command, args []string) {
			RunCreateFrontend(f, out, cmd)
		},
	}
	cmd.Flags().StringP("frontend-information", "i", "", "Event handler frontend information in the form key1=value1,key2=value2...")
	cmd.Flags().StringP("secret", "s", "", "Kubernetes secret to provide to the event handler frontend for credentials")

	return cmd
}

//RunCreateFrontend creates the frontend on the server
func RunCreateFrontend(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	address, err := cmdutil.GetServerAddress(cmd)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
		return nil
	}

	eventHandlerFrontend := cmdutil.GetFlagString(cmd, "frontend")
	eventHandlerFrontendType := cmdutil.GetFlagString(cmd, "frontendtype")
	eventHandlerFrontendInformation := cmdutil.GetFlagString(cmd, "frontend-information")
	eventHandlerFrontendSecret := cmdutil.GetFlagString(cmd, "secret")

	var frontendInformation map[string]string

	if "" != eventHandlerFrontendInformation {
		frontendInformation = make(map[string]string)

		frontendInfoStrings := strings.Split(eventHandlerFrontendInformation, ",")

		for _, frontendInfoString := range frontendInfoStrings {
			frontendInfo := strings.Split(frontendInfoString, "=")

			frontendInformation[frontendInfo[0]] = frontendInfo[1]
		}
	}

	fmt.Fprintf(out, "Create event handler frontend = %s with frontend type = %s\n", eventHandlerFrontend, eventHandlerFrontendType)
	if nil != frontendInformation {
		fmt.Fprint(out, "Additional information provided: ")
		for k := range frontendInformation {
			fmt.Fprintf(out, "%s ", k)
		}
		fmt.Fprint(out, "\n")
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect to server: %v", err)
	}
	defer conn.Close()
	c := klessapi.NewKlessAPIClient(conn)

	// Contact the server and print out its response.
	r, err := c.CreateEventHandlerFrontend(context.Background(), &klessapi.CreateEventHandlerFrontendRequest{EventHandlerFrontendName: eventHandlerFrontend,
		EventHandlerFrontendType:        eventHandlerFrontendType,
		EventHandlerFrontendInformation: frontendInformation,
		EventHandlerFrontendSecret:      eventHandlerFrontendSecret,
	})
	if err != nil {
		log.Fatalf("Could not create event handler frontend: %v", err)
	}

	fmt.Fprintf(out, "Server response: %s\n", r.Response)

	return nil
}
