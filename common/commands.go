package common

import "flag"

type Command struct {
	Name        string
	Description string
	FlagSet     *flag.FlagSet
	Run         func(args []string)
}

type CommandGroup struct {
	Name        string
	Description string
	Commands    []Command
}
