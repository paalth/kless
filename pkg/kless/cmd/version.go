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

// NewCmdVersion displays the client and optionally the server version
func NewCmdVersion(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the client and server version",
		Run: func(cmd *cobra.Command, args []string) {
			RunVersion(f, out, cmd)
		},
	}
	cmd.Flags().BoolP("client", "", false, "Client version only")
	cmd.Flags().MarkShorthandDeprecated("client", "please use --client instead")

	return cmd
}

// RunVersion executes the version display
func RunVersion(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {
	clientVersion := "0.0.1-alpha1"

	fmt.Fprintf(out, "Client Version: %s\n", clientVersion)
	if cmdutil.GetFlagBool(cmd, "client") {
		return nil
	}

	address, err := cmdutil.GetServerAddress(cmd)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
		return nil
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not not connect to server: %v", err)
	}
	defer conn.Close()

	c := klessapi.NewKlessAPIClient(conn)

	r, err := c.GetServerVersion(context.Background(), &klessapi.GetServerVersionRequest{Clientversion: clientVersion})
	if err != nil {
		log.Fatalf("Could not get server version: %v", err)
	}

	fmt.Fprintf(out, "Server Version: %s\n", r.Serverversion)

	return nil
}
