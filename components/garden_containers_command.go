package components

import (
	"flag"
	"os"

	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/components/garden"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
)

func GardenContainersCommand() common.Command {
	var (
		raw           bool
		gardenAddr    string
		gardenNetwork string
	)

	flagSet := flag.NewFlagSet("garden-containers", flag.ExitOnError)
	flagSet.BoolVar(&raw, "raw", false, "display raw response")
	flagSet.StringVar(&gardenAddr, "gardenAddr", "", "garden API address")
	flagSet.StringVar(&gardenNetwork, "gardenNetwork", "", "garden API network (unix/tcp)")

	return common.Command{
		Name:        "garden-containers",
		Description: "[file] - Fetch garden containers",
		FlagSet:     flagSet,
		Run: func(args []string) {
			gardenAddr, gardenNetwork, err := config_finder.FindGardenAddr(gardenAddr, gardenNetwork)
			common.ExitIfError("Could not find garden", err)

			if len(args) == 0 {
				err := garden.GardenContainers(gardenAddr, gardenNetwork, raw, os.Stdout)
				common.ExitIfError("Failed to fetch garden containers", err)
			} else {
				f, err := os.Create(args[0])
				common.ExitIfError("Could not create file", err)

				err = garden.GardenContainers(gardenAddr, gardenNetwork, raw, f)
				common.ExitIfError("Failed to fetch garden containers", err)

				f.Close()
			}
		},
	}
}
