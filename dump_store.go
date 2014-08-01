package main

import (
	"flag"
	"io"
	"time"

	"github.com/cloudfoundry-incubator/veritas/config_finder"
	"github.com/cloudfoundry-incubator/veritas/fetch_store"
	"github.com/cloudfoundry-incubator/veritas/print_store"
	"github.com/cloudfoundry-incubator/veritas/say"
)

func DumpStoreCommand() Command {
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

	return Command{
		Name:        "dump-store",
		Description: "- Fetch and print contents of the BBS",
		FlagSet:     flagSet,
		Run: func(args []string) {
			etcdCluster, err := config_finder.FindETCDCluster(etcdClusterFlag)
			ExitIfError("Could not find etcd cluster", err)

			if rate == 0 {
				err = dump(etcdCluster, verbose, tasks, lrps, services)
				ExitIfError("Failed to dump", err)
				return
			}

			ticker := time.NewTicker(rate)
			for {
				<-ticker.C
				say.Clear()
				err = dump(etcdCluster, verbose, tasks, lrps, services)
				if err != nil {
					say.Println(0, say.Red("Failed to dump: %s", err.Error()))
				}
			}
		},
	}
}

func dump(etcdCluster []string, verbose bool, tasks bool, lrps bool, services bool) error {
	reader, writer := io.Pipe()

	errs := make(chan error)
	go func() {
		err := fetch_store.Fetch(etcdCluster, false, writer)
		errs <- err
	}()
	go func() {
		err := print_store.PrintStore(verbose, tasks, lrps, services, reader)
		errs <- err
	}()

	err1 := <-errs
	err2 := <-errs
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}
