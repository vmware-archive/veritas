package components

import (
	"flag"
	"os"

	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/components/rep"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
)

func RepStateCommand() common.Command {
	var (
		etcdClusterFlag   string
		consulClusterFlag string
	)

	flagSet := flag.NewFlagSet("executor-containers", flag.ExitOnError)
	flagSet.StringVar(&etcdClusterFlag, "etcdCluster", "", "comma-separated etcd cluster urls")
	flagSet.StringVar(&consulClusterFlag, "consulCluster", "", "comma-separated consul cluster urls")

	return common.Command{
		Name:        "rep-states",
		Description: "- Fetch all rep states",
		FlagSet:     flagSet,
		Run: func(args []string) {
			veritasBBS, _, err := config_finder.ConstructBBS(etcdClusterFlag, consulClusterFlag)
			common.ExitIfError("Could not construct BBS", err)

			err = rep.RepState(veritasBBS, os.Stdout)
			common.ExitIfError("Failed to fetch rep states", err)
		},
	}
}
