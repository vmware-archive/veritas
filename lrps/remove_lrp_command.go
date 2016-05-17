package lrps

import (
	"flag"
	"os"

	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-golang/lager"
)

func RemoveLRPCommand() common.Command {
	var (
		bbsConfig config_finder.BBSConfig
	)

	flagSet := flag.NewFlagSet("remove-lrp", flag.ExitOnError)
	bbsConfig.PopulateFlags(flagSet)
	logger := lager.NewLogger("veritas")

	return common.Command{
		Name:        "remove-lrp",
		Description: "<process-guid> - remove an lrp",
		FlagSet:     flagSet,
		Run: func(args []string) {
			bbsClient, err := config_finder.NewBBS(bbsConfig)
			common.ExitIfError("Could not construct BBS", err)

			if len(args) == 0 {
				say.Fprintln(os.Stderr, 0, say.Red("You must specify a process-guid"))
				os.Exit(1)
			} else {
				err := bbsClient.RemoveDesiredLRP(logger, args[0])
				common.ExitIfError("Failed to remove lrp", err)
			}
		},
	}
}
