package components

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/common"
	"github.com/cloudfoundry-incubator/veritas/components/vitals"
	"github.com/cloudfoundry-incubator/veritas/config_finder"
)

func VitalsCommand() common.Command {
	var (
		addrs string
	)

	flagSet := flag.NewFlagSet("vitals", flag.ExitOnError)
	flagSet.StringVar(&addrs, "vitalsAddrs", "", "debug addresses: name:addr:port,...")

	return common.Command{
		Name:        "vitals",
		Description: "[file] - Fetch vitals for passed in golang processes",
		FlagSet:     flagSet,
		Run: func(args []string) {
			vitalsAddrs, err := config_finder.FindVitalsAddrs(addrs)
			common.ExitIfError("Could not find vitals addrs", err)

			if len(args) == 0 {
				err := vitals.Vitals(vitalsAddrs, os.Stdout)
				common.ExitIfError("Failed to fetch vitals", err)
			} else {
				f, err := os.Create(args[0])
				common.ExitIfError("Could not create file", err)

				err = vitals.Vitals(vitalsAddrs, f)
				common.ExitIfError("Failed to fetch vitals", err)

				f.Close()
			}
		},
	}
}
