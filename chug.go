package main

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/chug_commands"
)

func ChugCommand() Command {
	var (
		rel  string
		data string
	)

	flagSet := flag.NewFlagSet("chug", flag.ExitOnError)
	flagSet.StringVar(&rel, "rel", "", "render timestamps as durations relative to: 'first', 'now', or a number interpreted as a unix timestamp")
	flagSet.StringVar(&data, "data", "short", "render data: 'none', 'short', 'long'")

	return Command{
		Name:        "chug",
		Description: "[file] - Prettify lager logs",
		FlagSet:     flagSet,
		Run: func(args []string) {
			if len(args) == 0 {
				err := chug_commands.Prettify(rel, data, os.Stdin)
				ExitIfError("Failed to chug", err)
			} else {
				f, err := os.Open(args[0])
				ExitIfError("Could not open file", err)

				err = chug_commands.Prettify(rel, data, f)
				ExitIfError("Failed to chug", err)

				f.Close()
			}
		},
	}
}