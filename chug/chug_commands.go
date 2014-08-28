package chug

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/common"
)

func ChugCommand() common.Command {
	var (
		rel          string
		data         string
		hideNonLager bool
	)

	flagSet := flag.NewFlagSet("chug", flag.ExitOnError)
	flagSet.StringVar(&rel, "rel", "", "render timestamps as durations relative to: 'first', 'now', or a number interpreted as a unix timestamp")
	flagSet.StringVar(&data, "data", "short", "render data: 'none', 'short', 'long'")
	flagSet.BoolVar(&hideNonLager, "hide", false, "hide non-lager logs")

	return common.Command{
		Name:        "chug",
		Description: "[file] - Prettify lager logs",
		FlagSet:     flagSet,
		Run: func(args []string) {
			if len(args) == 0 {
				err := Prettify(rel, data, hideNonLager, os.Stdin)
				common.ExitIfError("Failed to chug", err)
			} else {
				f, err := os.Open(args[0])
				common.ExitIfError("Could not open file", err)

				err = Prettify(rel, data, hideNonLager, f)
				common.ExitIfError("Failed to chug", err)

				f.Close()
			}
		},
	}
}

func ServeChugCommand() common.Command {
	var (
		addr string
		dev  bool
	)

	flagSet := flag.NewFlagSet("chug-serve", flag.ExitOnError)
	flagSet.StringVar(&addr, "addr", "127.0.0.1:0", "address to serve chug")
	flagSet.BoolVar(&dev, "dev", false, "dev mode")

	return common.Command{
		Name:        "chug-serve",
		Description: "[file] - Serve up pretty lager logs",
		FlagSet:     flagSet,
		Run: func(args []string) {
			if len(args) == 0 {
				err := ServeLogs(addr, dev, os.Stdin)
				common.ExitIfError("Failed to serve chug", err)
			} else {
				f, err := os.Open(args[0])
				common.ExitIfError("Could not open file", err)

				err = ServeLogs(addr, dev, f)
				common.ExitIfError("Failed to serve chug", err)

				f.Close()
			}
		},
	}
}
