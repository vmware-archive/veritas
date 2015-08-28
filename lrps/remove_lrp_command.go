package lrps

import (
	"flag"
	"os"

	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-cf-experimental/veritas/lrps/remove_lrp"
)

func RemoveLRPCommand() common.Command {
	var (
		bbsEndpointFlag string
	)

	flagSet := flag.NewFlagSet("remove-lrp", flag.ExitOnError)
	flagSet.StringVar(&bbsEndpointFlag, "bbsEndpoint", "", "bbs endpoint")

	return common.Command{
		Name:        "remove-lrp",
		Description: "process-guid - undesired an lrp",
		FlagSet:     flagSet,
		Run: func(args []string) {
			bbsClient, err := config_finder.ConstructBBS(bbsEndpointFlag)
			common.ExitIfError("Could not construct BBS", err)

			if len(args) == 0 {
				say.Fprintln(os.Stderr, 0, say.Red("You must specify a process-guid"))
				os.Exit(1)
			} else {
				err = remove_lrp.RemoveLRP(bbsClient, args[0])
				common.ExitIfError("Failed to remove lrp", err)
			}
		},
	}
}
