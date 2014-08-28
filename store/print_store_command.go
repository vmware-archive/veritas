package store

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/common"
	"github.com/cloudfoundry-incubator/veritas/store/print_store"
)

func PrintStoreCommand() common.Command {
	var (
		tasks    bool
		lrps     bool
		services bool
		verbose  bool
	)

	flagSet := flag.NewFlagSet("print-store", flag.ExitOnError)
	flagSet.BoolVar(&verbose, "v", false, "be verbose")
	flagSet.BoolVar(&tasks, "tasks", true, "print tasks")
	flagSet.BoolVar(&lrps, "lrps", true, "print lrps")
	flagSet.BoolVar(&services, "services", true, "print services")

	return common.Command{
		Name:        "print-store",
		Description: "[file] - Print previously fetched contents of the BBS.  If file is blank, reads from stdin.",
		FlagSet:     flagSet,
		Run: func(args []string) {
			if len(args) == 0 {
				err := print_store.PrintStore(verbose, tasks, lrps, services, false, os.Stdin)
				common.ExitIfError("Failed to print store", err)
			} else {
				f, err := os.Open(args[0])
				common.ExitIfError("Could not open file", err)

				err = print_store.PrintStore(verbose, tasks, lrps, services, false, f)
				common.ExitIfError("Failed to print store", err)
			}
		},
	}
}
