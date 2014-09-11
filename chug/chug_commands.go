package chug

import (
	"flag"
	"io"
	"os"
	"time"

	"github.com/cloudfoundry-incubator/veritas/common"
	"github.com/cloudfoundry-incubator/veritas/say"
)

func baseFlagSet(command string, minTimeFlag *string, maxTimeFlag *string) *flag.FlagSet {
	flagSet := flag.NewFlagSet(command, flag.ExitOnError)
	flagSet.StringVar(minTimeFlag, "min", "", "only include entries logged after this time: either a unix timestamp, a chug-formatted time, or a duration (relative to now)")
	flagSet.StringVar(maxTimeFlag, "max", "", "only include entries logged before this time: either a unix timestamp, a chug-formatted time, or a duration (relative to now)")

	return flagSet
}

func parseBaseFlags(minTimeFlag, maxTimeFlag string) (time.Time, time.Time) {
	minTime, err := ParseTimeFlag(minTimeFlag)
	common.ExitIfError("Failed to parse min", err)
	maxTime, err := ParseTimeFlag(maxTimeFlag)
	common.ExitIfError("Failed to parse max", err)
	return minTime, maxTime
}

func ChugCommand() common.Command {
	var (
		minTimeFlag  string
		maxTimeFlag  string
		rel          string
		data         string
		hideNonLager bool
	)

	flagSet := baseFlagSet("chug", &minTimeFlag, &maxTimeFlag)
	flagSet.StringVar(&rel, "rel", "", "render timestamps as durations relative to: 'first', 'now', or a number interpreted as a unix timestamp")
	flagSet.StringVar(&data, "data", "short", "render data: 'none', 'short', 'long'")
	flagSet.BoolVar(&hideNonLager, "hide", false, "hide non-lager logs")

	return common.Command{
		Name:        "chug",
		Description: "[file] - Prettify lager logs",
		FlagSet:     flagSet,
		Run: func(args []string) {
			minTime, maxTime := parseBaseFlags(minTimeFlag, maxTimeFlag)

			if len(args) == 0 {
				err := Prettify(rel, data, hideNonLager, minTime, maxTime, os.Stdin)
				common.ExitIfError("Failed to chug", err)
			} else {
				f, err := os.Open(args[0])
				common.ExitIfError("Could not open file", err)

				err = Prettify(rel, data, hideNonLager, minTime, maxTime, f)
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

	flagSet := baseFlagSet("chug-unify", &minTimeFlag, &maxTimeFlag)

	return common.Command{
		Name:        "chug-unify",
		Description: "file1, file2,... - Combine lager files in temporal order",
		FlagSet:     flagSet,
		Run: func(args []string) {
			minTime, maxTime := parseBaseFlags(minTimeFlag, maxTimeFlag)

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
		minTimeFlag string
		maxTimeFlag string
		addr        string
		dev         bool
	)

	flagSet := baseFlagSet("chug-serve", &minTimeFlag, &maxTimeFlag)
	flagSet.StringVar(&addr, "addr", "127.0.0.1:0", "address to serve chug")
	flagSet.BoolVar(&dev, "dev", false, "dev mode")

	return common.Command{
		Name:        "chug-serve",
		Description: "[file] - Serve up pretty lager logs",
		FlagSet:     flagSet,
		Run: func(args []string) {
			minTime, maxTime := parseBaseFlags(minTimeFlag, maxTimeFlag)

			if len(args) == 0 {
				err := ServeLogs(addr, dev, minTime, maxTime, os.Stdin)
				common.ExitIfError("Failed to serve chug", err)
			} else {
				f, err := os.Open(args[0])
				common.ExitIfError("Could not open file", err)

				err = ServeLogs(addr, dev, minTime, maxTime, f)
				common.ExitIfError("Failed to serve chug", err)

				f.Close()
			}
		},
	}
}
