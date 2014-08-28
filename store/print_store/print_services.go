package print_store

import (
	"github.com/cloudfoundry-incubator/veritas/say"
	"github.com/cloudfoundry-incubator/veritas/veritas_models"
)

func printServices(verbose bool, services veritas_models.VeritasServices) {
	say.PrintBanner(say.Green("Services"), "~")
	say.Println(0, say.Green("Executors"))
	for _, executor := range services.Executors {
		say.Println(1, "%s (%s)", executor.ExecutorID, executor.Stack)
	}

	say.Println(0, say.Green("File Servers"))
	for _, fileServer := range services.FileServers {
		say.Println(1, fileServer)
	}
}
