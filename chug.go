package main

import (
	"flag"
	"io"
	"os"

	"github.com/cloudfoundry-incubator/veritas/chug_commands"
)

func ChugCommand() Command {
	var (
		rel          string
		data         string
		hideNonLager bool
	)

	flagSet := flag.NewFlagSet("chug", flag.ExitOnError)
	flagSet.StringVar(&rel, "rel", "", "render timestamps as durations relative to: 'first', 'now', or a number interpreted as a unix timestamp")
	flagSet.StringVar(&data, "data", "short", "render data: 'none', 'short', 'long'")
	flagSet.BoolVar(&hideNonLager, "hide", false, "hide non-lager logs")

	return Command{
		Name:        "chug",
		Description: "[file1, file2,...] - Prettify lager logs",
		FlagSet:     flagSet,
		Run: func(args []string) {
			if len(args) == 0 {
				err := chug_commands.Prettify(rel, data, hideNonLager, os.Stdin)
				ExitIfError("Failed to chug", err)
			} else {
				readers := []io.Reader{}
				for _, arg := range args {
					f, err := os.Open(args[0])
					ExitIfError("Could not open file", err)
					readers = append(readers, f)
				}

				err = chug_commands.Prettify(rel, data, hideNonLager, readers...)
				ExitIfError("Failed to chug", err)
			}
		},
	}
}

func ServeChugCommand() Command {
	var (
		addr string
		dev  bool
	)

	flagSet := flag.NewFlagSet("chug-serve", flag.ExitOnError)
	flagSet.StringVar(&addr, "addr", ":8080", "address to serve chug")
	flagSet.BoolVar(&dev, "dev", false, "dev mode")

	return Command{
		Name:        "chug-serve",
		Description: "[file] - Serve up pretty lager logs",
		FlagSet:     flagSet,
		Run: func(args []string) {
			if len(args) == 0 {
				err := chug_commands.ServeLogs(addr, dev, os.Stdin)
				ExitIfError("Failed to serve chug", err)
			} else {
				f, err := os.Open(args[0])
				ExitIfError("Could not open file", err)

				err = chug_commands.ServeLogs(addr, dev, f)
				ExitIfError("Failed to serve chug", err)

				f.Close()
			}
		},
	}
}
