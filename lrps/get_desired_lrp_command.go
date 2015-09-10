package lrps

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
)

func GetDesiredLRPCommand() common.Command {
	var (
		bbsEndpointFlag string
	)

	flagSet := flag.NewFlagSet("get-desired-lrp", flag.ExitOnError)
	flagSet.StringVar(&bbsEndpointFlag, "bbsEndpoint", "", "bbs endpoint")

	return common.Command{
		Name:        "get-desired-lrp",
		Description: "<process-guid> - get a DesiredLRP",
		FlagSet:     flagSet,
		Run: func(args []string) {
			bbsClient, err := config_finder.ConstructBBS(bbsEndpointFlag)
			common.ExitIfError("Could not construct BBS", err)

			if len(args) == 0 {
				say.Fprintln(os.Stderr, 0, say.Red("missing process-guid"))
				os.Exit(1)
			}

			desiredLRP, err := bbsClient.DesiredLRPByProcessGuid(args[0])
			common.ExitIfError("Failed to fetch DesiredLRP", err)

			preview, _ := json.MarshalIndent(desiredLRP, "", "  ")
			say.Println(0, string(preview))
		},
	}
}
