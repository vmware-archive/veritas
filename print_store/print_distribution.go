package print_store

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

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

	if clear {
		say.Clear()
	}

	printDistribution(dump, tasks, lrps)

	return nil
}

func printDistribution(dump veritas_models.StoreDump, includeTasks bool, includeLRPS bool) {
	executorIDs := []string{}
	for _, executorPresence := range dump.Services.Executors {
		executorIDs = append(executorIDs, executorPresence.ExecutorID)
	}

	sort.Strings(executorIDs)

	nTasks := map[string]int{}
	nLRPs := map[string]int{}

	for _, tasks := range dump.Tasks {
		for _, task := range tasks {
			nTasks[task.ExecutorID]++
		}
	}

	for _, lrp := range dump.LRPS {
		for _, actuals := range lrp.ActualLRPsByIndex {
			for _, actual := range actuals {
				nLRPs[actual.ExecutorID]++
			}
		}
	}

	say.Println(0, "Distribution")
	for _, executorID := range executorIDs {
		numTasks := nTasks[executorID]
		numLRPs := nLRPs[executorID]
		var content string
		if numTasks == 0 && numLRPs == 0 {
			content = say.Red("Empty")
		} else {
			content = fmt.Sprintf("%s%s", say.Yellow(strings.Repeat("•", nTasks[executorID])), say.Green(strings.Repeat("•", nLRPs[executorID])))
		}
		say.Println(0, "%12s: %s", executorID, content)
	}
}
