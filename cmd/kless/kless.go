package main

import (
	"os"

	"github.com/paalth/kless/cmd/kless/app"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize()
}

func main() {
	if err := app.Run(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
