package veritas

import (
	"flag"
	"os"
	"strings"

	"github.com/cloudfoundry-incubator/veritas/config_finder"
	"github.com/cloudfoundry-incubator/veritas/fetch_store"
)

func FetchStoreCommand() Command {
	var (
		raw             bool
		etcdClusterFlag string
	)

	flagSet := flag.NewFlagSet("fetch-store", flag.ExitOnError)
	flagSet.BoolVar(&raw, "raw", false, "fetch raw store dump")
	flagSet.StringVar(&etcdClusterFlag, "etcdCluster", "", "comma-separated etcd cluster urls")

	return Command{
		Name:        "fetch-store",
		Description: "[file] - Fetch contents of the BBS",
		FlagSet:     flagSet,
		Run: func(args []string) {
			etcdCluster, err := config_finder.FindETCDCluster(strings.Split(etcdClusterFlag, ","))
			ExitIfError("Could not find etcd cluster", err)

			if len(args) == 0 {
				err := fetch_store.Fetch(etcdCluster, raw, os.Stdout)
				ExitIfError("Failed to fetch store", err)
			} else {
				f, err := os.Create(args[0])
				ExitIfError("Could not create file", err)

				err = fetch_store.Fetch(etcdCluster, raw, f)
				ExitIfError("Failed to fetch store", err)

				f.Close()
			}
		},
	}
}
