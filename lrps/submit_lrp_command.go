package lrps

import (
	"flag"
	"os"

	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-cf-experimental/veritas/lrps/submit_lrp"
)

func SubmitLRPCommand() common.Command {
	var (
		etcdClusterFlag string
	)

	flagSet := flag.NewFlagSet("submit-lrp", flag.ExitOnError)
	flagSet.StringVar(&etcdClusterFlag, "etcdCluster", "", "comma-separated etcd cluster urls")

	return common.Command{
		Name:        "submit-lrp",
		Description: "[file] - submits a desired lrp to the bbs",
		FlagSet:     flagSet,
		Run: func(args []string) {
			etcdCluster, err := config_finder.FindETCDCluster(etcdClusterFlag)
			common.ExitIfError("Could not find etcd cluster", err)

			if len(args) == 0 {
				err := submit_lrp.SubmitLRP(etcdCluster, nil)
				common.ExitIfError("Failed to submit lrp", err)
			} else {
				f, err := os.Open(args[0])
				common.ExitIfError("Could not open file", err)

				err = submit_lrp.SubmitLRP(etcdCluster, f)
				common.ExitIfError("Failed to submit lrp", err)
			}
		},
	}
}
