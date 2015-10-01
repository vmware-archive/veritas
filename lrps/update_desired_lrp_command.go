package lrps

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
)

func UpdateDesiredLRPCommand() common.Command {
	var (
		bbsConfig config_finder.BBSConfig
	)

	flagSet := flag.NewFlagSet("update-lrp", flag.ExitOnError)
	bbsConfig.PopulateFlags(flagSet)

	return common.Command{
		Name:        "update-lrp",
		Description: "<process-guid> <path to json file> - update a DesiredLRP",
		FlagSet:     flagSet,
		Run: func(args []string) {
			bbsClient, err := config_finder.NewBBS(bbsConfig)
			common.ExitIfError("Could not construct BBS", err)

			var raw = []byte{}

			if len(args) == 0 {
				say.Fprintln(os.Stderr, 0, say.Red("missing process-guid"))
				os.Exit(1)
			} else if len(args) == 1 {
				say.Fprintln(os.Stderr, 0, "Reading from stdin...")
				raw, err = ioutil.ReadAll(os.Stdin)
				common.ExitIfError("Failed to read from stdin", err)
			} else {
				raw, err = ioutil.ReadFile(args[1])
				common.ExitIfError("Failed to read from file", err)
			}

			desiredLRPUpdate := &models.DesiredLRPUpdate{}

			err = json.Unmarshal(raw, desiredLRPUpdate)
			common.ExitIfError("Failed to decode DesiredLRPUpdate", err)

			say.Println(0, "Updating %s:", args[0])
			preview, _ := json.MarshalIndent(desiredLRPUpdate, "", "  ")
			say.Println(0, string(preview))

			err = bbsClient.UpdateDesiredLRP(args[0], desiredLRPUpdate)
			common.ExitIfError("Failed to update DesiredLRP", err)
		},
	}
}
