package print_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/cloudfoundry-incubator/runtime-schema/models"

	"github.com/cloudfoundry-incubator/veritas/say"
	"github.com/cloudfoundry-incubator/veritas/veritas_models"
)

func PrintDistribution(tasks bool, lrps bool, clear bool, f io.Reader) error {
	decoder := json.NewDecoder(f)
	var dump veritas_models.StoreDump
	err := decoder.Decode(&dump)
	if err != nil {
		return err
	}

	printDistribution(dump, tasks, lrps, clear)

	return nil
}

func printDistribution(dump veritas_models.StoreDump, includeTasks bool, includeLRPS bool, clear bool) {
	executorIDs := []string{}
	for _, executorPresence := range dump.Services.Executors {
		executorIDs = append(executorIDs, executorPresence.ExecutorID)
	}

	sort.Strings(executorIDs)

	nTasks := map[string]int{}
	nLRPsStarting := map[string]int{}
	nLRPsRunning := map[string]int{}

	for _, tasks := range dump.Tasks {
		for _, task := range tasks {
			nTasks[task.ExecutorID]++
		}
	}

	for _, lrp := range dump.LRPS {
		for _, actuals := range lrp.ActualLRPsByIndex {
			for _, actual := range actuals {
				if actual.State == models.ActualLRPStateStarting {
					nLRPsStarting[actual.ExecutorID]++
				} else {
					nLRPsRunning[actual.ExecutorID]++
				}
			}
		}
	}

	buffer := &bytes.Buffer{}
	say.Println(0, "Distribution")
	for _, executorID := range executorIDs {
		numTasks := nTasks[executorID]
		numLRPs := nLRPsStarting[executorID] + nLRPsRunning[executorID]
		var content string
		if numTasks == 0 && numLRPs == 0 {
			content = say.Red("Empty")
		} else {
			content = fmt.Sprintf("%s%s", say.Yellow(strings.Repeat("•", nTasks[executorID])), say.Green(strings.Repeat("•", nLRPsRunning[executorID])), say.Gray(strings.Repeat("•", nLRPsStarting[executorID])))
		}
		say.Fprintln(buffer, 0, "%12s: %s", executorID, content)
	}

	if clear {
		say.Clear()
	}
	buffer.WriteTo(os.Stdout)
}
