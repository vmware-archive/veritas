package store

import (
	"flag"
	"os"

	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-cf-experimental/veritas/store/fetch_store"
)

func FetchStoreCommand() common.Command {
	var (
		bbsEndpointFlag string
	)

	flagSet := flag.NewFlagSet("fetch-store", flag.ExitOnError)
	flagSet.StringVar(&bbsEndpointFlag, "bbsEndpoint", "", "bbs endpoint")

	return common.Command{
		Name:        "fetch-store",
		Description: "[file] - Fetch contents of the BBS",
		FlagSet:     flagSet,
		Run: func(args []string) {
			bbsClient, err := config_finder.ConstructBBS(bbsEndpointFlag)
			common.ExitIfError("Could not construct BBS", err)

			if len(args) == 0 {
				err := fetch_store.Fetch(bbsClient, os.Stdout)
				common.ExitIfError("Failed to fetch store", err)
			} else {
				f, err := os.Create(args[0])
				common.ExitIfError("Could not create file", err)

				err = fetch_store.Fetch(bbsClient, f)
				common.ExitIfError("Failed to fetch store", err)

				f.Close()
			}
		},
	}
}
