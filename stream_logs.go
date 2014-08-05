package main

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/config_finder"
	"github.com/cloudfoundry-incubator/veritas/log_commands"
	"github.com/cloudfoundry-incubator/veritas/say"
)

func StreamLogsCommand() Command {
	var (
		loggregatorAddr string
	)

	flagSet := flag.NewFlagSet("stream-logs", flag.ExitOnError)
	flagSet.StringVar(&loggregatorAddr, "loggregatorAddr", "", "loggregator OUT addr")

	return Command{
		Name:        "stream-logs",
		Description: "app-id - Fetch loggregator-logs for the given app-id",
		FlagSet:     flagSet,
		Run: func(args []string) {
			if len(args) == 0 {
				say.Fprintln(os.Stderr, 0, say.Red("You must specify an app-id"))
				os.Exit(1)
			}
			loggregatorAddr, err := config_finder.FindLoggregatorAddr(loggregatorAddr)
			ExitIfError("Could not find loggregator", err)

			err = log_commands.StreamLogs(loggregatorAddr, args[0], os.Stdout)
			ExitIfError("Failed to stream logs", err)
		},
	}
}
