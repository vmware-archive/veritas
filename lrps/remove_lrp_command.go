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
		etcdClusterFlag   string
		consulClusterFlag string
	)

	flagSet := flag.NewFlagSet("remove-lrp", flag.ExitOnError)
	flagSet.StringVar(&etcdClusterFlag, "etcdCluster", "", "comma-separated etcd cluster urls")
	flagSet.StringVar(&consulClusterFlag, "consulCluster", "", "comma-separated consul cluster urls")

	return common.Command{
		Name:        "remove-lrp",
		Description: "process-guid - undesired an lrp",
		FlagSet:     flagSet,
		Run: func(args []string) {
			veritasBBS, _, err := config_finder.ConstructBBS(etcdClusterFlag, consulClusterFlag)
			common.ExitIfError("Could not construct BBS", err)

			if len(args) == 0 {
				say.Fprintln(os.Stderr, 0, say.Red("You must specify a process-guid"))
				os.Exit(1)
			} else {
				err = remove_lrp.RemoveLRP(veritasBBS, args[0])
				common.ExitIfError("Failed to remove lrp", err)
			}
		},
	}
}
