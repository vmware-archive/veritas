package main

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/config_finder"
)

func AutodetectCommand() Command {
	flagSet := flag.NewFlagSet("autodetect", flag.ExitOnError)

	return Command{
		Name:        "autodetect",
		Description: " - autodetect configuration **must be run on a bosh vm**",
		FlagSet:     flagSet,
		Run: func(args []string) {
			err := config_finder.Autodetect(os.Stdout)
			ExitIfError("Autodetect failed", err)
		},
	}
}
