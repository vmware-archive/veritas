package main

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/veritas/say"
)

type Command struct {
	Name        string
	Description string
	FlagSet     *flag.FlagSet
	Run         func(args []string)
}

func main() {
	commands := []Command{
		FetchStoreCommand(),
		PrintStoreCommand(),
	}

	if len(os.Args) == 1 || os.Args[1] == "help" {
		usage(commands)
		os.Exit(1)
	}

	for _, command := range commands {
		if command.Name == os.Args[1] {
			command.FlagSet.Parse(os.Args[1:])
			command.Run(command.FlagSet.Args())
		}
	}
}

func usage(commands []Command) {
	say.FprintBanner(os.Stderr, "Veritas", "=")
	for _, command := range commands {
		say.Fprintln(os.Stderr, 0, "%s %s", say.Green(command.Name), say.LightGray(command.Description))
		command.FlagSet.PrintDefaults()
		say.Fprintln(os.Stderr, 0, "")
	}
}
