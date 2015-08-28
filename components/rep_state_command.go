package components

import (
	"flag"
	"os"

	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/components/rep"
)

func RepStateCommand() common.Command {
	flagSet := flag.NewFlagSet("rep-state", flag.ExitOnError)

	return common.Command{
		Name:        "rep-state",
		Description: "- Fetch state for rep on localhost",
		FlagSet:     flagSet,
		Run: func(args []string) {
			err := rep.RepState(os.Stdout)
			common.ExitIfError("Failed to fetch rep state", err)
		},
	}
}
