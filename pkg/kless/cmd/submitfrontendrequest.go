package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	cmdutil "github.com/paalth/kless/pkg/kless/cmd/util"
	klessapi "github.com/paalth/kless/pkg/klessserver/grpc"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//NewCmdSubmitFrontendRequest creates a frontend
func NewCmdSubmitFrontendRequest(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frontend",
		Short: "Submit request to event handler frontend",
		Run: func(cmd *cobra.Command, args []string) {
			RunSubmitFrontendRequest(f, out, cmd)
		},
	}
	cmd.Flags().StringP("request", "r", "", "Request file to submit")
	cmd.Flags().StringP("url", "u", "", "URL of request to submit")

	return cmd
}

//RunSubmitFrontendRequest creates the frontend on the server
func RunSubmitFrontendRequest(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	address, err := cmdutil.GetServerAddress(cmd)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
		return nil
	}

	eventHandlerFrontend := cmdutil.GetFlagString(cmd, "frontend")
	request := cmdutil.GetFlagString(cmd, "request")
	//url := cmdutil.GetFlagString(cmd, "url")

	var requestBody []byte

	if "" != request {
		requestBody, err = ioutil.ReadFile(request)
		if nil != err {
			log.Fatalf("could not read request from file %s: %v", request, err)
			return nil
		}
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect to server: %v", err)
	}
	defer conn.Close()
	c := klessapi.NewKlessAPIClient(conn)

	// Contact the server and print out its response.
	r, err := c.SubmitEventHandlerFrontendRequest(context.Background(), &klessapi.SubmitEventHandlerFrontendRequestRequest{
		EventHandlerFrontendName: eventHandlerFrontend,
		Request:                  requestBody,
	})
	if err != nil {
		log.Fatalf("Could not submit request to handler frontend: %v", err)
	}

	fmt.Fprintf(out, "Server status: %s\n", r.Status)
	fmt.Fprintf(out, "Server response: %s\n", string(r.Response))

	return nil
}
