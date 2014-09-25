package loggregator_logs

import (
	"flag"
	"os"

	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-cf-experimental/veritas/say"
)

func StreamLogsCommand() common.Command {
	var (
		loggregatorAddr string
	)

	flagSet := flag.NewFlagSet("stream-logs", flag.ExitOnError)
	flagSet.StringVar(&loggregatorAddr, "loggregatorAddr", "", "loggregator OUT addr")

	return common.Command{
		Name:        "stream-logs",
		Description: "app-id - Fetch loggregator-logs for the given app-id",
		FlagSet:     flagSet,
		Run: func(args []string) {
			if len(args) == 0 {
				say.Fprintln(os.Stderr, 0, say.Red("You must specify an app-id"))
				os.Exit(1)
			}
			loggregatorAddr, err := config_finder.FindLoggregatorAddr(loggregatorAddr)
			common.ExitIfError("Could not find loggregator", err)

			err = StreamLogs(loggregatorAddr, args[0], os.Stdout)
			common.ExitIfError("Failed to stream logs", err)
		},
	}
}
