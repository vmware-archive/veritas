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

func DistributionCommand() Command {
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

	return Command{
		Name:        "distribution",
		Description: "- Fetch and print distribution of Tasks and LRPs",
		FlagSet:     flagSet,
		Run: func(args []string) {
			etcdCluster, err := config_finder.FindETCDCluster(etcdClusterFlag)
			ExitIfError("Could not find etcd cluster", err)

			if rate == 0 {
				err = distribution(etcdCluster, tasks, lrps, false)
				ExitIfError("Failed to print distribution", err)
				return
			}

			ticker := time.NewTicker(rate)
			for {
				<-ticker.C
				err = distribution(etcdCluster, tasks, lrps, true)
				if err != nil {
					say.Println(0, say.Red("Failed to print distribution: %s", err.Error()))
				}
			}
		},
	}
}

func distribution(etcdCluster []string, tasks bool, lrps bool, clear bool) error {
	reader, writer := io.Pipe()

	errs := make(chan error)
	go func() {
		err := fetch_store.Fetch(etcdCluster, false, writer)
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
