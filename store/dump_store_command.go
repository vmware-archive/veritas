package store

import (
	"flag"
	"io"
	"time"

	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/cloudfoundry/storeadapter/workerpool"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-cf-experimental/veritas/say"
	"github.com/pivotal-cf-experimental/veritas/store/fetch_store"
	"github.com/pivotal-cf-experimental/veritas/store/print_store"
)

func DumpStoreCommand() common.Command {
	var (
		etcdClusterFlag string
		tasks           bool
		lrps            bool
		services        bool
		verbose         bool
		rate            time.Duration
	)

	flagSet := flag.NewFlagSet("dump-store", flag.ExitOnError)
	flagSet.StringVar(&etcdClusterFlag, "etcdCluster", "", "comma-separated etcd cluster urls")
	flagSet.BoolVar(&verbose, "v", false, "be verbose")
	flagSet.BoolVar(&tasks, "tasks", true, "print tasks")
	flagSet.BoolVar(&lrps, "lrps", true, "print lrps")
	flagSet.BoolVar(&services, "services", true, "print services")
	flagSet.DurationVar(&rate, "rate", time.Duration(0), "rate at which to poll the store")

	return common.Command{
		Name:        "dump-store",
		Description: "- Fetch and print contents of the BBS",
		FlagSet:     flagSet,
		Run: func(args []string) {
			etcdCluster, err := config_finder.FindETCDCluster(etcdClusterFlag)
			common.ExitIfError("Could not find etcd cluster", err)

			adapter := etcdstoreadapter.NewETCDStoreAdapter(etcdCluster, workerpool.NewWorkerPool(10))
			err = adapter.Connect()
			common.ExitIfError("Could not connect to etcd cluster", err)

			if rate == 0 {
				err = dump(adapter, verbose, tasks, lrps, services, false)
				common.ExitIfError("Failed to dump", err)
				return
			}

			ticker := time.NewTicker(rate)
			for {
				<-ticker.C
				err = dump(adapter, verbose, tasks, lrps, services, true)
				if err != nil {
					say.Println(0, say.Red("Failed to dump: %s", err.Error()))
				}
			}
		},
	}
}

func dump(adapter *etcdstoreadapter.ETCDStoreAdapter, verbose bool, tasks bool, lrps bool, services bool, clear bool) error {
	reader, writer := io.Pipe()

	errs := make(chan error)
	go func() {
		err := fetch_store.Fetch(adapter, false, writer)
		errs <- err
	}()
	go func() {
		err := print_store.PrintStore(verbose, tasks, lrps, services, clear, reader)
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
