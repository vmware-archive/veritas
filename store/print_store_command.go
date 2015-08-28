package store

import (
	"flag"
	"os"

	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/store/print_store"
)

func PrintStoreCommand() common.Command {
	var (
		tasks   bool
		lrps    bool
		verbose bool
	)

	flagSet := flag.NewFlagSet("print-store", flag.ExitOnError)
	flagSet.BoolVar(&verbose, "v", false, "be verbose")
	flagSet.BoolVar(&tasks, "tasks", true, "print tasks")
	flagSet.BoolVar(&lrps, "lrps", true, "print lrps")

	return common.Command{
		Name:        "print-store",
		Description: "[file] - Print previously fetched contents of the BBS.  If file is blank, reads from stdin.",
		FlagSet:     flagSet,
		Run: func(args []string) {
			if len(args) == 0 {
				err := print_store.PrintStore(verbose, tasks, lrps, false, os.Stdin)
				common.ExitIfError("Failed to print store", err)
			} else {
				f, err := os.Open(args[0])
				common.ExitIfError("Could not open file", err)

				err = print_store.PrintStore(verbose, tasks, lrps, false, f)
				common.ExitIfError("Failed to print store", err)
			}
		},
	}
}
