package tasks

import (
	"flag"
	"os"

	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-cf-experimental/veritas/tasks/submit_task"
)

func SubmitTaskCommand() common.Command {
	var (
		receptorEndpointFlag string
		receptorUsernameFlag string
		receptorPasswordFlag string
	)

	flagSet := flag.NewFlagSet("submit-task", flag.ExitOnError)
	flagSet.StringVar(&receptorEndpointFlag, "receptorEndpoint", "", "receptor url (e.g. 127.0.0.1:8888)")
	flagSet.StringVar(&receptorUsernameFlag, "receptorUsername", "", "receptor username")
	flagSet.StringVar(&receptorPasswordFlag, "receptorPassword", "", "receptor password")

	return common.Command{
		Name:        "submit-task",
		Description: "[file] - submits a task to the Diego API",
		FlagSet:     flagSet,
		Run: func(args []string) {
			receptorClient, err := config_finder.FindReceptor(receptorEndpointFlag, receptorUsernameFlag, receptorPasswordFlag)
			common.ExitIfError("Could not find API", err)

			if len(args) == 0 {
				err := submit_task.SubmitTask(receptorClient, nil)
				common.ExitIfError("Failed to submit task", err)
			} else {
				f, err := os.Open(args[0])
				common.ExitIfError("Could not open file", err)

				err = submit_task.SubmitTask(receptorClient, f)
				common.ExitIfError("Failed to submit task", err)
			}
		},
	}
}
