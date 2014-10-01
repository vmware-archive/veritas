package store

import (
	"flag"
	"os"

	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/cloudfoundry/storeadapter/workerpool"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-cf-experimental/veritas/store/fetch_store"
)

func FetchStoreCommand() common.Command {
	var (
		raw             bool
		etcdClusterFlag string
	)

	flagSet := flag.NewFlagSet("fetch-store", flag.ExitOnError)
	flagSet.BoolVar(&raw, "raw", false, "fetch raw store dump")
	flagSet.StringVar(&etcdClusterFlag, "etcdCluster", "", "comma-separated etcd cluster urls")

	return common.Command{
		Name:        "fetch-store",
		Description: "[file] - Fetch contents of the BBS",
		FlagSet:     flagSet,
		Run: func(args []string) {
			etcdCluster, err := config_finder.FindETCDCluster(etcdClusterFlag)
			common.ExitIfError("Could not find etcd cluster", err)

			adapter := etcdstoreadapter.NewETCDStoreAdapter(etcdCluster, workerpool.NewWorkerPool(10))
			err = adapter.Connect()
			common.ExitIfError("Could not connect to etcd cluster", err)

			if len(args) == 0 {
				err := fetch_store.Fetch(adapter, raw, os.Stdout)
				common.ExitIfError("Failed to fetch store", err)
			} else {
				f, err := os.Create(args[0])
				common.ExitIfError("Could not create file", err)

				err = fetch_store.Fetch(adapter, raw, f)
				common.ExitIfError("Failed to fetch store", err)

				f.Close()
			}
		},
	}
}
