package lrps

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"

	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-golang/lager"
)

func GetActualLRPCommand() common.Command {
	var (
		bbsConfig config_finder.BBSConfig
	)

	flagSet := flag.NewFlagSet("get-actual-lrp", flag.ExitOnError)
	bbsConfig.PopulateFlags(flagSet)
	logger := lager.NewLogger("veritas")

	return common.Command{
		Name:        "get-actual-lrp",
		Description: "<process-guid> <optional: index> - get an ActualLRP",
		FlagSet:     flagSet,
		Run: func(args []string) {
			bbsClient, err := config_finder.NewBBS(bbsConfig)
			common.ExitIfError("Could not construct BBS", err)

			var index = -1

			if len(args) == 0 {
				say.Fprintln(os.Stderr, 0, say.Red("missing process-guid"))
				os.Exit(1)
			}

			processGuid := args[0]

			if len(args) == 2 {
				index, err = strconv.Atoi(args[1])
				common.ExitIfError("Could not parse index", err)
			}

			if index == -1 {
				actualLRPGroups, err := bbsClient.ActualLRPGroupsByProcessGuid(logger, processGuid)
				common.ExitIfError("Could not fetch ActualLRPs", err)

				for _, actualLRPGroup := range actualLRPGroups {
					actualLRP, _ := actualLRPGroup.Resolve()
					preview, _ := json.MarshalIndent(actualLRP, "", "  ")
					say.Println(0, string(preview))
				}
			} else {
				actualLRPGroup, err := bbsClient.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
				common.ExitIfError("Could not fetch ActualLRP", err)

				actualLRP, _ := actualLRPGroup.Resolve()
				preview, _ := json.MarshalIndent(actualLRP, "", "  ")
				say.Println(0, string(preview))
			}
		},
	}
}
