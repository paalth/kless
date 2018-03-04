package cmd

import (
	"io"

	cmdutil "github.com/paalth/kless/pkg/kless/cmd/util"

	"github.com/spf13/cobra"
)

// NewKlessCommand returns the available commands...
func NewKlessCommand(f cmdutil.Factory, in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "kless",
		Short: "kless interacts with the Kless Server",
		Long:  "",
		Run:   runHelp,
	}

	cmds.PersistentFlags().StringP("eventhandler", "e", "eventhandler", "Event handler name")
	cmds.PersistentFlags().StringP("eventhandlerbuilder", "b", "eventhandlerbuilder", "Event handler builder name")
	cmds.PersistentFlags().StringP("namespace", "n", "kless", "Event handler namespace")
	cmds.PersistentFlags().StringP("frontendtype", "t", "kless", "Event handler frontend type")
	cmds.PersistentFlags().StringP("frontend", "f", "kless", "Event handler frontend")
	cmds.PersistentFlags().StringP("comment", "c", "", "Comment")
	cmds.PersistentFlags().String("serveraddress", "", "Kless server address (name:port)")

	createCmds := &cobra.Command{
		Use:   "create",
		Short: "create commands",
		Long:  "Create Kless objects",
	}

	getCmds := &cobra.Command{
		Use:   "get",
		Short: "get commands",
		Long:  "Get Kless information",
	}

	updateCmds := &cobra.Command{
		Use:   "update",
		Short: "update commands",
		Long:  "Update Kless objects",
	}

	deleteCmds := &cobra.Command{
		Use:   "delete",
		Short: "delete commands",
		Long:  "Delete Kless objects",
	}

	describeCmds := &cobra.Command{
		Use:   "describe",
		Short: "describe commands",
		Long:  "Get detailed Kless object information",
	}

	submitCmds := &cobra.Command{
		Use:   "submit",
		Short: "submit commands",
		Long:  "Submit request to Kless object",
	}

	cmds.AddCommand(createCmds)
	cmds.AddCommand(getCmds)
	cmds.AddCommand(updateCmds)
	cmds.AddCommand(deleteCmds)
	cmds.AddCommand(describeCmds)
	cmds.AddCommand(submitCmds)

	createCmds.AddCommand(NewCmdCreateEventHandler(f, out))
	createCmds.AddCommand(NewCmdCreateEventHandlerBuilder(f, out))
	createCmds.AddCommand(NewCmdCreateFrontendType(f, out))
	createCmds.AddCommand(NewCmdCreateFrontend(f, out))

	getCmds.AddCommand(NewCmdGetEventHandlers(f, out))
	getCmds.AddCommand(NewCmdGetEventHandlerBuilders(f, out))
	getCmds.AddCommand(NewCmdGetFrontendTypes(f, out))
	getCmds.AddCommand(NewCmdGetFrontends(f, out))
	getCmds.AddCommand(NewCmdGetStats(f, out))

	updateCmds.AddCommand(NewCmdUpdateEventHandlerBuilder(f, out))

	deleteCmds.AddCommand(NewCmdDeleteEventHandler(f, out))
	deleteCmds.AddCommand(NewCmdDeleteEventHandlerBuilder(f, out))
	deleteCmds.AddCommand(NewCmdDeleteEventHandlerFrontendType(f, out))
	deleteCmds.AddCommand(NewCmdDeleteEventHandlerFrontend(f, out))

	describeCmds.AddCommand(NewCmdDescribeEventHandler(f, out))

	submitCmds.AddCommand(NewCmdSubmitFrontendRequest(f, out))

	cmds.AddCommand(NewCmdVersion(f, out))

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
