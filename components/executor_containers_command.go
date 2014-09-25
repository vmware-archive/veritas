package components

import (
	"flag"
	"os"

	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/components/executor"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
)

func ExecutorContainersCommand() common.Command {
	var (
		raw          bool
		executorAddr string
	)

	flagSet := flag.NewFlagSet("executor-containers", flag.ExitOnError)
	flagSet.BoolVar(&raw, "raw", false, "display raw response")
	flagSet.StringVar(&executorAddr, "executorAddr", "", "executor API url")

	return common.Command{
		Name:        "executor-containers",
		Description: "[file] - Fetch containers as the executor sees them",
		FlagSet:     flagSet,
		Run: func(args []string) {
			executorAddr, err := config_finder.FindExecutorAddr(executorAddr)
			common.ExitIfError("Could not find executor", err)

			if len(args) == 0 {
				err := executor.ExecutorContainers(executorAddr, raw, os.Stdout)
				common.ExitIfError("Failed to fetch executor containers", err)
			} else {
				f, err := os.Create(args[0])
				common.ExitIfError("Could not create file", err)

				err = executor.ExecutorContainers(executorAddr, raw, f)
				common.ExitIfError("Failed to fetch executor containers", err)

				f.Close()
			}
		},
	}
}
