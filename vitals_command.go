package main

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/config_finder"
	"github.com/cloudfoundry-incubator/veritas/vitals_commands"
)

func VitalsCommand() Command {
	var (
		addrs string
	)

	flagSet := flag.NewFlagSet("vitals", flag.ExitOnError)
	flagSet.StringVar(&addrs, "vitalsAddrs", "", "debug addresses: name:addr:port,...")

	return Command{
		Name:        "vitals",
		Description: "[file] - Fetch vitals for passed in golang processes",
		FlagSet:     flagSet,
		Run: func(args []string) {
			vitalsAddrs, err := config_finder.FindVitalsAddrs(addrs)
			ExitIfError("Could not find vitals addrs", err)

			if len(args) == 0 {
				err := vitals_commands.Vitals(vitalsAddrs, os.Stdout)
				ExitIfError("Failed to fetch vitals", err)
			} else {
				f, err := os.Create(args[0])
				ExitIfError("Could not create file", err)

				err = vitals_commands.Vitals(vitalsAddrs, f)
				ExitIfError("Failed to fetch vitals", err)

				f.Close()
			}
		},
	}
}
