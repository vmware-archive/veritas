package main

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/config_finder"
	"github.com/cloudfoundry-incubator/veritas/warden_commands"
)

func WardenContainersCommand() Command {
	var (
		raw           bool
		wardenAddr    string
		wardenNetwork string
	)

	flagSet := flag.NewFlagSet("warden-containers", flag.ExitOnError)
	flagSet.BoolVar(&raw, "raw", false, "display raw response")
	flagSet.StringVar(&wardenAddr, "wardenAddr", "", "warden API address")
	flagSet.StringVar(&wardenNetwork, "wardenNetwork", "", "warden API network (unix/tcp)")

	return Command{
		Name:        "warden-containers",
		Description: "[file] - Fetch warden containers",
		FlagSet:     flagSet,
		Run: func(args []string) {
			wardenAddr, wardenNetwork, err := config_finder.FindWardenAddr(wardenAddr, wardenNetwork)
			ExitIfError("Could not find warden", err)

			if len(args) == 0 {
				err := warden_commands.WardenContainers(wardenAddr, wardenNetwork, raw, os.Stdout)
				ExitIfError("Failed to fetch warden containers", err)
			} else {
				f, err := os.Create(args[0])
				ExitIfError("Could not create file", err)

				err = warden_commands.WardenContainers(wardenAddr, wardenNetwork, raw, f)
				ExitIfError("Failed to fetch warden containers", err)

				f.Close()
			}
		},
	}
}
