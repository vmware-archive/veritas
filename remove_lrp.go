package main

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/config_finder"
	"github.com/cloudfoundry-incubator/veritas/remove_lrp"
	"github.com/cloudfoundry-incubator/veritas/say"
)

func RemoveLRPCommand() Command {
	var (
		etcdClusterFlag string
	)

	flagSet := flag.NewFlagSet("remove-lrp", flag.ExitOnError)
	flagSet.StringVar(&etcdClusterFlag, "etcdCluster", "", "comma-separated etcd cluster urls")

	return Command{
		Name:        "remove-lrp",
		Description: "process-guid - undesired an lrp",
		FlagSet:     flagSet,
		Run: func(args []string) {
			etcdCluster, err := config_finder.FindETCDCluster(etcdClusterFlag)
			ExitIfError("Could not find etcd cluster", err)

			if len(args) == 0 {
				say.Red("You must specify a process-guid")
				os.Exit(1)
			} else {
				err = remove_lrp.RemoveLRP(etcdCluster, args[0])
				ExitIfError("Failed to remove lrp", err)
			}
		},
	}
}
