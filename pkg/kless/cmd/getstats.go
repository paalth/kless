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

// NewCmdGetStats retrieves event handler statistics from the server
func NewCmdGetStats(f cmdutil.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Display event handler statistics",
		Run: func(cmd *cobra.Command, args []string) {
			RunGetStats(f, out, cmd)
		},
	}
	cmd.Flags().StringP("version", "v", "1.0", "Event handler version")
	cmd.Flags().StringP("podname", "p", "", "Event handler pod name")
	cmd.Flags().StringP("counts", "c", "", "Return summary event count information")
	cmd.Flags().StringP("last", "l", "", "Return information for some time period up untill the present, values can be 1h, 1d,...")
	cmd.Flags().String("starttime", "", "Start time (format YYYY-MM-DDTHH:MM:SS.MMMZ)")
	cmd.Flags().String("endtime", "", "End time (format YYYY-MM-DDTHH:MM:SS.MMMZ)")

	return cmd
}

// RunGetStats retrieves stats from the klessserver
func RunGetStats(f cmdutil.Factory, out io.Writer, cmd *cobra.Command) error {

	address, err := cmdutil.GetServerAddress(cmd)
	if err != nil {
		log.Fatalf("Unable to get server connection information: %v", err)
		return nil
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not not connect to server: %v", err)
	}
	defer conn.Close()

	c := klessapi.NewKlessAPIClient(conn)

	stream, err := c.GetEventHandlerStatistics(context.Background(), &klessapi.GetEventHandlerStatisticsRequest{Clientversion: "0.0.1-alpha1"})
	if err != nil {
		log.Fatalf("Could not get event handler statistics: %v", err)
	}

	fmt.Fprintf(out, "%-25s %-15s %-25s %-15s %-40s %-12s %-12s %-12s\n", "TIMESTAMP", "NAMESPACE", "EVENTHANDLER", "VERSION", "PODNAME", "REQUESTSIZE", "RESPONSESIZE", "RESPONSETIME")
	for {
		s, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not get event handler statistics: %v", err)
			break
		}
		fmt.Fprintf(out, "%-25s %-15s %-25s %-15s %-40s %-12v %-12v %-12v\n", s.Timestamp, s.EventHandlerNamespace, s.EventHandlerName, s.EventHandlerVersion, s.PodName, s.RequestSize, s.ResponseSize, s.ResponseTime)
	}

	return nil
}
