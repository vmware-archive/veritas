package store

import (
	"flag"
	"io"
	"time"

	"github.com/cloudfoundry/gunk/workpool"
	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-cf-experimental/veritas/say"
	"github.com/pivotal-cf-experimental/veritas/store/fetch_store"
	"github.com/pivotal-cf-experimental/veritas/store/print_store"
)

func DistributionCommand() common.Command {
	var (
		etcdClusterFlag string
		tasks           bool
		lrps            bool
		rate            time.Duration
	)

	flagSet := flag.NewFlagSet("distribution", flag.ExitOnError)
	flagSet.StringVar(&etcdClusterFlag, "etcdCluster", "", "comma-separated etcd cluster urls")
	flagSet.BoolVar(&tasks, "tasks", true, "print tasks")
	flagSet.BoolVar(&lrps, "lrps", true, "print lrps")
	flagSet.DurationVar(&rate, "rate", time.Duration(0), "rate at which to poll the store")

	return common.Command{
		Name:        "distribution",
		Description: "- Fetch and print distribution of Tasks and LRPs",
		FlagSet:     flagSet,
		Run: func(args []string) {
			etcdCluster, err := config_finder.FindETCDCluster(etcdClusterFlag)
			common.ExitIfError("Could not find etcd cluster", err)

			adapter := etcdstoreadapter.NewETCDStoreAdapter(etcdCluster, workpool.NewWorkPool(10))
			err = adapter.Connect()
			common.ExitIfError("Could not connect to etcd cluster", err)

			if rate == 0 {
				err = distribution(adapter, tasks, lrps, false)
				common.ExitIfError("Failed to print distribution", err)
				return
			}

			ticker := time.NewTicker(rate)
			for {
				<-ticker.C
				err = distribution(adapter, tasks, lrps, true)
				if err != nil {
					say.Println(0, say.Red("Failed to print distribution: %s", err.Error()))
				}
			}
		},
	}
}

func distribution(adapter *etcdstoreadapter.ETCDStoreAdapter, tasks bool, lrps bool, clear bool) error {
	reader, writer := io.Pipe()

	errs := make(chan error)
	go func() {
		err := fetch_store.Fetch(adapter, false, writer)
		errs <- err
	}()
	go func() {
		err := print_store.PrintDistribution(tasks, lrps, clear, reader)
		errs <- err
	}()

	err1 := <-errs
	if err1 != nil {
		return err1
	}
	err2 := <-errs
	if err2 != nil {
		return err2
	}
	return nil
}
