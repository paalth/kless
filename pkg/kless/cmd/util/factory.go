package util

import (
	"github.com/spf13/pflag"
)

type Factory interface {
	FlagSet() *pflag.FlagSet
}

type factory struct {
	flags *pflag.FlagSet
}

func NewFactory() Factory {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)

	f := &factory{
	  flags: flags,
	}

	return f
}

func (f *factory) FlagSet() *pflag.FlagSet {
	return f.flags
}

