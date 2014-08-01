package main

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/config_finder"
	"github.com/cloudfoundry-incubator/veritas/executor_commands"
)

func ExecutorResourcesCommand() Command {
	var (
		raw          bool
		executorAddr string
	)

	flagSet := flag.NewFlagSet("executor-resources", flag.ExitOnError)
	flagSet.BoolVar(&raw, "raw", false, "display raw response")
	flagSet.StringVar(&executorAddr, "executorAddr", "", "executor API url")

	return Command{
		Name:        "executor-resources",
		Description: "[file] - Fetch initial and available resources for an executor",
		FlagSet:     flagSet,
		Run: func(args []string) {
			executorAddr, err := config_finder.FindExecutorAddr(executorAddr)
			ExitIfError("Could not find executor", err)

			if len(args) == 0 {
				err := executor_commands.ExecutorResources(executorAddr, raw, os.Stdout)
				ExitIfError("Failed to fetch executor resources", err)
			} else {
				f, err := os.Create(args[0])
				ExitIfError("Could not create file", err)

				err = executor_commands.ExecutorResources(executorAddr, raw, f)
				ExitIfError("Failed to fetch executor resources", err)

				f.Close()
			}
		},
	}
}
