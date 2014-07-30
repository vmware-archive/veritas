package main

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/config_finder"
	"github.com/cloudfoundry-incubator/veritas/submit_lrp"
)

func SubmitLRPCommand() Command {
	var (
		etcdClusterFlag string
	)

	flagSet := flag.NewFlagSet("submit-lrp", flag.ExitOnError)
	flagSet.StringVar(&etcdClusterFlag, "etcdCluster", "", "comma-separated etcd cluster urls")

	return Command{
		Name:        "submit-lrp",
		Description: "[file] - a json representation of an lrp",
		FlagSet:     flagSet,
		Run: func(args []string) {
			etcdCluster, err := config_finder.FindETCDCluster(etcdClusterFlag)
			ExitIfError("Could not find etcd cluster", err)

			if len(args) == 0 {
				err := submit_lrp.SubmitLRP(etcdCluster, nil)
				ExitIfError("Failed to fetch store", err)
			} else {
				f, err := os.Open(args[0])
				ExitIfError("Could not open file", err)

				err = submit_lrp.SubmitLRP(etcdCluster, f)
				ExitIfError("Failed to fetch store", err)
			}
		},
	}
}
