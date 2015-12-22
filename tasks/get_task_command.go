package tasks

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
)

func GetTaskCommand() common.Command {
	var (
		bbsConfig config_finder.BBSConfig
	)

	flagSet := flag.NewFlagSet("get-task", flag.ExitOnError)
	bbsConfig.PopulateFlags(flagSet)

	return common.Command{
		Name:        "get-task",
		Description: "<task-guid> - get a Task",
		FlagSet:     flagSet,
		Run: func(args []string) {
			bbsClient, err := config_finder.NewBBS(bbsConfig)
			common.ExitIfError("Could not construct BBS", err)

			if len(args) == 0 {
				say.Fprintln(os.Stderr, 0, say.Red("missing task"))
				os.Exit(1)
			}

			task, err := bbsClient.TaskByGuid(args[0])
			common.ExitIfError("Failed to fetch Task", err)

			preview, _ := json.MarshalIndent(task, "", "  ")
			say.Println(0, string(preview))
		},
	}
}
