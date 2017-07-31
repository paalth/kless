package main

import (
	"fmt"
	"os"

	"github.com/paalth/kless/cmd/klessserver/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
