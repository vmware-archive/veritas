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

func CreateDesiredLRPCommand() common.Command {
	var (
		bbsConfig config_finder.BBSConfig
	)

	flagSet := flag.NewFlagSet("desire-lrp", flag.ExitOnError)
	bbsConfig.PopulateFlags(flagSet)

	return common.Command{
		Name:        "desire-lrp",
		Description: "<path to json file> - create a DesiredLRP",
		FlagSet:     flagSet,
		Run: func(args []string) {
			bbsClient, err := config_finder.NewBBS(bbsConfig)
			common.ExitIfError("Could not construct BBS", err)

			var raw = []byte{}

			if len(args) == 0 {
				say.Fprintln(os.Stderr, 0, "Reading from stdin...")
				raw, err = ioutil.ReadAll(os.Stdin)
				common.ExitIfError("Failed to read from stdin", err)
			} else {
				raw, err = ioutil.ReadFile(args[0])
				common.ExitIfError("Failed to read from file", err)
			}

			desiredLRP := &models.DesiredLRP{}

			err = json.Unmarshal(raw, desiredLRP)
			common.ExitIfError("Failed to decode DesiredLRP", err)

			say.Println(0, "Desiring:")
			preview, _ := json.MarshalIndent(desiredLRP, "", "  ")
			say.Println(0, string(preview))

			err = bbsClient.DesireLRP(desiredLRP)
			common.ExitIfError("Failed to desire DesiredLRP", err)
		},
	}
}
