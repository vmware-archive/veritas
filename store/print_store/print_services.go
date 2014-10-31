package print_store

import (
	"github.com/pivotal-cf-experimental/veritas/say"
	"github.com/pivotal-cf-experimental/veritas/veritas_models"
)

func printServices(verbose bool, services veritas_models.VeritasServices) {
	say.PrintBanner(say.Green("Services"), "~")
	say.Println(0, say.Green("Executors"))
	for _, executor := range services.Executors {
		say.Println(1, "%s (%s)", executor.ExecutorID, executor.Stack)
	}
}
