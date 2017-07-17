package app

import (
	"os"

	"github.com/paalth/kless/pkg/kless/cmd"
	cmdutil "github.com/paalth/kless/pkg/kless/cmd/util"
)

func Run() error {
	cmd := cmd.NewKlessCommand(cmdutil.NewFactory(), os.Stdin, os.Stdout, os.Stderr)
	return cmd.Execute()
}
