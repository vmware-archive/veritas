package lrps

import (
	"flag"
	"os"

	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-cf-experimental/veritas/lrps/remove_lrp"
)

func RemoveLRPCommand() common.Command {
	var (
		etcdClusterFlag string
	)

	flagSet := flag.NewFlagSet("remove-lrp", flag.ExitOnError)
	flagSet.StringVar(&etcdClusterFlag, "etcdCluster", "", "comma-separated etcd cluster urls")

	return common.Command{
		Name:        "remove-lrp",
		Description: "process-guid - undesired an lrp",
		FlagSet:     flagSet,
		Run: func(args []string) {
			etcdCluster, err := config_finder.FindETCDCluster(etcdClusterFlag)
			common.ExitIfError("Could not find etcd cluster", err)

			if len(args) == 0 {
				say.Fprintln(os.Stderr, 0, say.Red("You must specify a process-guid"))
				os.Exit(1)
			} else {
				err = remove_lrp.RemoveLRP(etcdCluster, args[0])
				common.ExitIfError("Failed to remove lrp", err)
			}
		},
	}
}
