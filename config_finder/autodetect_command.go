package config_finder

import (
	"flag"
	"os"

	"github.com/pivotal-cf-experimental/veritas/common"
)

func AutodetectCommand() common.Command {
	flagSet := flag.NewFlagSet("autodetect", flag.ExitOnError)

	return common.Command{
		Name:        "autodetect",
		Description: "- autodetect configuration **must be run on a bosh vm**",
		FlagSet:     flagSet,
		Run: func(args []string) {
			err := Autodetect(os.Stdout)
			common.ExitIfError("Autodetect failed", err)
		},
	}
}
