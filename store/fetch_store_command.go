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
		verbose           bool
		etcdClusterFlag   string
		consulClusterFlag string
	)

	flagSet := flag.NewFlagSet("fetch-store", flag.ExitOnError)
	flagSet.BoolVar(&verbose, "v", false, "fetch raw store dump")
	flagSet.StringVar(&etcdClusterFlag, "etcdCluster", "", "comma-separated etcd cluster urls")
	flagSet.StringVar(&consulClusterFlag, "consulCluster", "", "comma-separated consul cluster urls")

	return common.Command{
		Name:        "fetch-store",
		Description: "[file] - Fetch contents of the BBS",
		FlagSet:     flagSet,
		Run: func(args []string) {
			veritasBBS, etcdStore, err := config_finder.ConstructBBS(etcdClusterFlag, consulClusterFlag)
			common.ExitIfError("Could not construct BBS", err)

			if len(args) == 0 {
				err := fetch_store.Fetch(veritasBBS, etcdStore, verbose, os.Stdout)
				common.ExitIfError("Failed to fetch store", err)
			} else {
				f, err := os.Create(args[0])
				common.ExitIfError("Could not create file", err)

				err = fetch_store.Fetch(veritasBBS, etcdStore, verbose, f)
				common.ExitIfError("Failed to fetch store", err)

				f.Close()
			}
		},
	}
}
