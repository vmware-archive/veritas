package lrps

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-golang/lager"
)

func GetDesiredLRPCommand() common.Command {
	var (
		bbsConfig config_finder.BBSConfig
	)

	logger := lager.NewLogger("veritas")
	flagSet := flag.NewFlagSet("get-desired-lrp", flag.ExitOnError)
	bbsConfig.PopulateFlags(flagSet)

	return common.Command{
		Name:        "get-desired-lrp",
		Description: "<process-guid> - get a DesiredLRP",
		FlagSet:     flagSet,
		Run: func(args []string) {
			bbsClient, err := config_finder.NewBBS(bbsConfig)
			common.ExitIfError("Could not construct BBS", err)

			if len(args) == 0 {
				say.Fprintln(os.Stderr, 0, say.Red("missing process-guid"))
				os.Exit(1)
			}

			desiredLRP, err := bbsClient.DesiredLRPByProcessGuid(logger, args[0])
			common.ExitIfError("Failed to fetch DesiredLRP", err)

			preview, _ := json.MarshalIndent(desiredLRP, "", "  ")
			say.Println(0, string(preview))
		},
	}
}
