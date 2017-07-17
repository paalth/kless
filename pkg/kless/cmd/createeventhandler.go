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

//NewCmdCreateEventHandler creates an event handler
func NewCmdCreateEventHandler(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handler",
		Short: "Create event handler",
		Run: func(cmd *cobra.Command, args []string) {
			RunCreateEventHandler(f, out, cmd)
		},
	}
	cmd.Flags().StringP("source", "s", "", "File containing event handler source code")
	cmd.Flags().StringP("frontend", "f", "", "Event handler frontend")
	cmd.Flags().StringP("url", "u", "", "URL of event handler source code")
	cmd.Flags().StringP("version", "v", "1.0", "Event handler version")
	cmd.Flags().StringP("dependencies-file", "d", "", "File containing event handler dependencies")
	cmd.Flags().StringP("dependencies-url", "l", "", "URL of event handler source dependencies")

	return cmd
}

//RunCreateEventHandler creates an event handler on the server
func RunCreateEventHandler(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	address, err := cmdutil.GetServerAddress(cmd)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
		return nil
	}

	eventHandlerName := cmdutil.GetFlagString(cmd, "eventhandler")
	eventHandlerNamespace := cmdutil.GetFlagString(cmd, "namespace")
	eventHandlerBuilder := cmdutil.GetFlagString(cmd, "eventhandlerbuilder")
	eventHandlerSourceCodeFile := cmdutil.GetFlagString(cmd, "source")
	eventHandlerSourceCodeURL := cmdutil.GetFlagString(cmd, "url")
	eventHandlerVersion := cmdutil.GetFlagString(cmd, "version")
	eventHandlerDependenciesFile := cmdutil.GetFlagString(cmd, "dependencies-file")
	eventHandlerDependenciesURL := cmdutil.GetFlagString(cmd, "dependencies-url")
	eventHandlerFrontend := cmdutil.GetFlagString(cmd, "frontend")

	var eventHandlerSourceCode []byte

	if "" != eventHandlerSourceCodeFile {
		eventHandlerSourceCode, err = ioutil.ReadFile(eventHandlerSourceCodeFile)
		if nil != err {
			log.Fatalf("could not read source code file %s: %v", eventHandlerSourceCodeFile, err)
			return nil
		}
	}

	var eventHandlerDependencies []byte

	if "" != eventHandlerDependenciesFile {
		eventHandlerDependencies, err = ioutil.ReadFile(eventHandlerDependenciesFile)
		if nil != err {
			log.Fatalf("could not read dependecies file %s: %v", eventHandlerDependenciesFile, err)
			return nil
		}
	}

	fmt.Fprintf(out, "Create event handler = %s version %s in namespace = %s using event handler builder = %s with frontend = %s\n", eventHandlerName, eventHandlerVersion, eventHandlerNamespace, eventHandlerBuilder, eventHandlerFrontend)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := klessapi.NewKlessAPIClient(conn)

	// Contact the server and print out its response.
	r, err := c.CreateEventHandler(context.Background(), &klessapi.CreateEventHandlerRequest{EventHandlerName: eventHandlerName,
		EventHandlerNamespace:       eventHandlerNamespace,
		EventHandlerBuilder:         eventHandlerBuilder,
		EventHandlerSourceCode:      eventHandlerSourceCode,
		EventHandlerSourceCodeURL:   eventHandlerSourceCodeURL,
		EventHandlerVersion:         eventHandlerVersion,
		EventHandlerFrontend:        eventHandlerFrontend,
		EventHandlerDependencies:    eventHandlerDependencies,
		EventHandlerDependenciesURL: eventHandlerDependenciesURL,
	})
	if err != nil {
		log.Fatalf("could not create event handler: %v", err)
	}

	fmt.Fprintf(out, "Server response: %s\n", r.Response)

	return nil
}
