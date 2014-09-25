package chug

import (
	"flag"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/say"
)

func baseFlagSet(command string, minTimeFlag, maxTimeFlag, matchFlag, excludeFlag *string) *flag.FlagSet {
	flagSet := flag.NewFlagSet(command, flag.ExitOnError)
	flagSet.StringVar(minTimeFlag, "min", "", "only include entries logged after this time: either a unix timestamp, a chug-formatted time, or a duration (relative to now)")
	flagSet.StringVar(maxTimeFlag, "max", "", "only include entries logged before this time: either a unix timestamp, a chug-formatted time, or a duration (relative to now)")
	flagSet.StringVar(matchFlag, "match", "", "only include entries that match this regular expression")
	flagSet.StringVar(excludeFlag, "exclude", "", "exclude entries that match this regular expression")

	return flagSet
}

func parseBaseFlags(minTimeFlag, maxTimeFlag, matchFlag, excludeFlag string) (time.Time, time.Time, *regexp.Regexp, *regexp.Regexp) {
	minTime, err := ParseTimeFlag(minTimeFlag)
	common.ExitIfError("Failed to parse -min", err)
	maxTime, err := ParseTimeFlag(maxTimeFlag)
	common.ExitIfError("Failed to parse -max", err)
	match, err := regexp.Compile(matchFlag)
	common.ExitIfError("Failed to parse -match", err)
	if matchFlag == "" {
		match = nil
	}
	exclude, err := regexp.Compile(excludeFlag)
	common.ExitIfError("Failed to parse -match", err)
	if excludeFlag == "" {
		exclude = nil
	}
	return minTime, maxTime, match, exclude
}

func ChugCommand() common.Command {
	var (
		minTimeFlag string
		maxTimeFlag string
		matchFlag   string
		excludeFlag string

		rel          string
		data         string
		hideNonLager bool
	)

	flagSet := baseFlagSet("chug", &minTimeFlag, &maxTimeFlag, &matchFlag, &excludeFlag)
	flagSet.StringVar(&rel, "rel", "", "render timestamps as durations relative to: 'first', 'now', or a number interpreted as a unix timestamp")
	flagSet.StringVar(&data, "data", "short", "render data: 'none', 'short', 'long'")
	flagSet.BoolVar(&hideNonLager, "hide", false, "hide non-lager logs")

	return common.Command{
		Name:        "chug",
		Description: "[file] - Prettify lager logs",
		FlagSet:     flagSet,
		Run: func(args []string) {
			minTime, maxTime, match, exclude := parseBaseFlags(minTimeFlag, maxTimeFlag, matchFlag, excludeFlag)

			if len(args) == 0 {
				err := Prettify(rel, data, hideNonLager, minTime, maxTime, match, exclude, os.Stdin)
				common.ExitIfError("Failed to chug", err)
			} else {
				f, err := os.Open(args[0])
				common.ExitIfError("Could not open file", err)

				err = Prettify(rel, data, hideNonLager, minTime, maxTime, match, exclude, f)
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
		matchFlag   string
		excludeFlag string
	)

	flagSet := baseFlagSet("chug-unify", &minTimeFlag, &maxTimeFlag, &matchFlag, &excludeFlag)

	return common.Command{
		Name:        "chug-unify",
		Description: "file1, file2,... - Combine lager files in temporal order",
		FlagSet:     flagSet,
		Run: func(args []string) {
			minTime, maxTime, match, exclude := parseBaseFlags(minTimeFlag, maxTimeFlag, matchFlag, excludeFlag)

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

				err := Unify(files, os.Stdout, minTime, maxTime, match, exclude)
				common.ExitIfError("Failed to chug-unify", err)
			}
		},
	}
}

func ServeChugCommand() common.Command {
	var (
		minTimeFlag string
		maxTimeFlag string
		matchFlag   string
		excludeFlag string
		addr        string
		dev         bool
	)

	flagSet := baseFlagSet("chug-serve", &minTimeFlag, &maxTimeFlag, &matchFlag, &excludeFlag)
	flagSet.StringVar(&addr, "addr", "127.0.0.1:0", "address to serve chug")
	flagSet.BoolVar(&dev, "dev", false, "dev mode")

	return common.Command{
		Name:        "chug-serve",
		Description: "[file] - Serve up pretty lager logs",
		FlagSet:     flagSet,
		Run: func(args []string) {
			minTime, maxTime, match, exclude := parseBaseFlags(minTimeFlag, maxTimeFlag, matchFlag, excludeFlag)

			if len(args) == 0 {
				err := ServeLogs(addr, dev, minTime, maxTime, match, exclude, os.Stdin)
				common.ExitIfError("Failed to serve chug", err)
			} else {
				f, err := os.Open(args[0])
				common.ExitIfError("Could not open file", err)

				err = ServeLogs(addr, dev, minTime, maxTime, match, exclude, f)
				common.ExitIfError("Failed to serve chug", err)

				f.Close()
			}
		},
	}
}
