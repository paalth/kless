package util

import (
	"errors"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func GetServerAddress(cmd *cobra.Command) (string, error) {
	address := GetFlagString(cmd, "serveraddress")

	if address == "" {
		address = os.Getenv("KLESS_SERVERADDRESS")
	}

	if address == "" {
		return "", errors.New("Kless server address not available")
	}

	return address, nil
}

func GetFlagBool(cmd *cobra.Command, flag string) bool {
	b, err := cmd.Flags().GetBool(flag)
	if err != nil {
		glog.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}

	return b
}

func GetFlagString(cmd *cobra.Command, flag string) string {
	s, err := cmd.Flags().GetString(flag)
	if err != nil {
		glog.Fatalf("error accessing flag %s for command %s: %v", flag, cmd.Name(), err)
	}

	return s
}
