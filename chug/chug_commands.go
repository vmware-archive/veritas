package chug

import (
	"flag"
	"io"
	"os"

	"github.com/cloudfoundry-incubator/veritas/common"
	"github.com/cloudfoundry-incubator/veritas/say"
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

func UnifyChugCommand() common.Command {
	var (
		minTimeFlag string
		maxTimeFlag string
	)

	flagSet := flag.NewFlagSet("chug-unify", flag.ExitOnError)
	flagSet.StringVar(&minTimeFlag, "min", "", "only include entries logged after this time (either a unix timestamp or a chug-formatted time)")
	flagSet.StringVar(&maxTimeFlag, "max", "", "only include entries logged before this time (either a unix timestamp or a chug-formatted time)")

	return common.Command{
		Name:        "chug-unify",
		Description: "file1, file2,... - Combine lager files in temporal order",
		FlagSet:     flagSet,
		Run: func(args []string) {
			minTime, err := ParseTimeFlag(minTimeFlag)
			common.ExitIfError("Failed to parse min", err)
			maxTime, err := ParseTimeFlag(maxTimeFlag)
			common.ExitIfError("Failed to parse max", err)

			if len(args) == 0 {
				say.Println(0, say.Red("You must pass chug-unify files to combine"))
				os.Exit(1)
			} else {
				files := []io.Reader{}
				for _, arg := range args {
					f, err := os.Open(arg)
					common.ExitIfError("Could not open file", err)
					files = append(files, f)
				}

				err := Unify(files, os.Stdout, minTime, maxTime)
				common.ExitIfError("Failed to chug-unify", err)
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
