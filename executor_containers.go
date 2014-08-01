package main

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/config_finder"
	"github.com/cloudfoundry-incubator/veritas/executor_commands"
)

func ExecutorContainersCommand() Command {
	var (
		raw          bool
		executorAddr string
	)

	flagSet := flag.NewFlagSet("executor-containers", flag.ExitOnError)
	flagSet.BoolVar(&raw, "raw", false, "display raw response")
	flagSet.StringVar(&executorAddr, "executorAddr", "", "executor API url")

	return Command{
		Name:        "executor-containers",
		Description: "[file] - Fetch containers as the executor sees them",
		FlagSet:     flagSet,
		Run: func(args []string) {
			executorAddr, err := config_finder.FindExecutorAddr(executorAddr)
			ExitIfError("Could not find executor", err)

			if len(args) == 0 {
				err := executor_commands.ExecutorContainers(executorAddr, raw, os.Stdout)
				ExitIfError("Failed to fetch executor containers", err)
			} else {
				f, err := os.Create(args[0])
				ExitIfError("Could not create file", err)

				err = executor_commands.ExecutorContainers(executorAddr, raw, f)
				ExitIfError("Failed to fetch executor containers", err)

				f.Close()
			}
		},
	}
}
