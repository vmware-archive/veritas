package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry-incubator/veritas/say"
)

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

func main() {
	commandGroups := []CommandGroup{
		CommandGroup{
			Name:        "Setup",
			Description: "Commands to set veritas up on a BOSH Job",
			Commands: []Command{
				AutodetectCommand(),
			},
		},

		CommandGroup{
			Name:        "BBS",
			Description: "Commands to fetch from the BBS",
			Commands: []Command{
				DumpStoreCommand(),
				FetchStoreCommand(),
				PrintStoreCommand(),
			},
		},
		CommandGroup{
			Name:        "Chug",
			Description: "Commands to prettify lager logs",
			Commands: []Command{
				ChugCommand(),
			},
		},
		CommandGroup{
			Name:        "Executor & Warden",
			Description: "Commands to fetch information from executor and warden",
			Commands: []Command{
				ExecutorResourcesCommand(),
				ExecutorContainersCommand(),
				WardenContainersCommand(),
			},
		},
		CommandGroup{
			Name:        "Vitals",
			Description: "Commands to fetch vitals for components",
			Commands: []Command{
				VitalsCommand(),
			},
		},
		CommandGroup{
			Name:        "Loggregator",
			Description: "Commands to stream loggregator logs",
			Commands: []Command{
				StreamLogsCommand(),
			},
		},
		CommandGroup{
			Name:        "DesiredLRPS " + say.Red("[DANGER]"),
			Description: "Commands to add and remove DesiredLRPs",
			Commands: []Command{
				SubmitLRPCommand(),
				RemoveLRPCommand(),
			},
		},
	}

	if len(os.Args) == 1 || os.Args[1] == "help" {
		usage(commandGroups)
		os.Exit(1)
	}

	if os.Args[1] == "completions" {
		completions(commandGroups)
		os.Exit(0)
	}

	for _, commandGroup := range commandGroups {
		for _, command := range commandGroup.Commands {
			if command.Name == os.Args[1] {
				command.FlagSet.Parse(os.Args[2:])
				command.Run(command.FlagSet.Args())
				os.Exit(0)
			}
		}
	}

	say.Println(0, say.Red("Unkown command: %s", os.Args[1]))
	usage(commandGroups)
}

func completions(commandGroups []CommandGroup) {
	availableCommands := []string{}
	for _, commands := range commandGroups {
		for _, command := range commands.Commands {
			availableCommands = append(availableCommands, command.Name)
		}
	}

	out := fmt.Sprintf(`
function _veritas() {
	local cur prev commands
	COMPREPLY=()
	cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
	commands="%s"

	if [[ "${COMP_CWORD}" == "1" ]] ; then
		COMPREPLY=( $(compgen -W "${commands} help completions" -- ${cur}) );
	elif [[ "${prev}" == "help" ]] ; then
		COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) );
	else
		COMPREPLY=( $(compgen -f ${cur}) );
	fi

	return 0
}

complete -F _veritas veritas
`, strings.Join(availableCommands, " "))

	say.Println(0, out)
}

func usage(commandGroups []CommandGroup) {
	if len(os.Args) > 2 {
		matcher := strings.ToLower(os.Args[2])
		for _, commandGroup := range commandGroups {
			if strings.HasPrefix(strings.ToLower(commandGroup.Name), matcher) {
				usageForCommandGroup(commandGroup, true)
				return
			}

			for _, command := range commandGroup.Commands {
				if strings.HasPrefix(strings.ToLower(command.Name), matcher) {
					usageForCommand(0, command, true)
					return
				}
			}
		}
		say.Fprintln(os.Stderr, 0, say.Red("Unkown command: %s", os.Args[2]))
	}

	say.Fprintln(os.Stderr, 0, "%s", say.Cyan("Help and Autocompletion"))
	say.Fprintln(os.Stderr, 0, strings.Repeat("-", len("Help and Autocompletion")))
	say.Fprintln(os.Stderr, 1, "%s %s", say.Green("help"), say.LightGray("[command] - Show this help, or detailed help for the passed in command"))
	say.Fprintln(os.Stderr, 1, "%s %s", say.Green("completions"), say.LightGray("Generate BASH Completions for veritas"))
	say.Fprintln(os.Stderr, 0, "")

	for _, commandGroup := range commandGroups {
		usageForCommandGroup(commandGroup, false)
		say.Println(0, "")
	}
}

func usageForCommandGroup(commandGroup CommandGroup, includeFlags bool) {
	say.Fprintln(os.Stderr, 0, "%s - %s", say.Cyan(commandGroup.Name), say.LightGray(commandGroup.Description))
	say.Fprintln(os.Stderr, 0, strings.Repeat("-", len(commandGroup.Name)+3+len(commandGroup.Description)))
	for _, command := range commandGroup.Commands {
		usageForCommand(1, command, includeFlags)
	}
}

func usageForCommand(indentation int, command Command, includeFlags bool) {
	say.Fprintln(os.Stderr, indentation, "%s %s", say.Green(command.Name), say.LightGray(command.Description))
	if includeFlags {
		buffer := &bytes.Buffer{}
		command.FlagSet.SetOutput(buffer)
		command.FlagSet.PrintDefaults()
		say.Fprintln(os.Stderr, indentation, buffer.String())
	}
}
