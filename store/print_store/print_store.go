package print_store

import (
	"encoding/json"
	"io"

	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/veritas_models"
)

func PrintStore(verbose bool, tasks bool, lrps bool, clear bool, f io.Reader) error {
	decoder := json.NewDecoder(f)
	var dump veritas_models.StoreDump
	err := decoder.Decode(&dump)
	if err != nil {
		return err
	}

	if clear {
		say.Clear()
	}

	if tasks {
		printTasks(verbose, dump.Tasks)
	}

	if lrps {
		printLRPS(verbose, dump.LRPS)
		printDomains(dump.Domains)
	}

	return nil
}
