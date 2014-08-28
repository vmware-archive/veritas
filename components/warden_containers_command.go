package components

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/common"
	"github.com/cloudfoundry-incubator/veritas/components/warden"
	"github.com/cloudfoundry-incubator/veritas/config_finder"
)

func WardenContainersCommand() common.Command {
	var (
		raw           bool
		wardenAddr    string
		wardenNetwork string
	)

	flagSet := flag.NewFlagSet("warden-containers", flag.ExitOnError)
	flagSet.BoolVar(&raw, "raw", false, "display raw response")
	flagSet.StringVar(&wardenAddr, "wardenAddr", "", "warden API address")
	flagSet.StringVar(&wardenNetwork, "wardenNetwork", "", "warden API network (unix/tcp)")

	return common.Command{
		Name:        "warden-containers",
		Description: "[file] - Fetch warden containers",
		FlagSet:     flagSet,
		Run: func(args []string) {
			wardenAddr, wardenNetwork, err := config_finder.FindWardenAddr(wardenAddr, wardenNetwork)
			common.ExitIfError("Could not find warden", err)

			if len(args) == 0 {
				err := warden.WardenContainers(wardenAddr, wardenNetwork, raw, os.Stdout)
				common.ExitIfError("Failed to fetch warden containers", err)
			} else {
				f, err := os.Create(args[0])
				common.ExitIfError("Could not create file", err)

				err = warden.WardenContainers(wardenAddr, wardenNetwork, raw, f)
				common.ExitIfError("Failed to fetch warden containers", err)

				f.Close()
			}
		},
	}
}
