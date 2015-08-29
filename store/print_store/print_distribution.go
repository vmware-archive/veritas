package print_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/veritas_models"
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
	nTasks := map[string]int{}
	nLRPsClaimed := map[string]int{}
	nLRPsRunning := map[string]int{}
	nLRPsEvacuating := map[string]int{}
	cells := []string{}
	knownCells := map[string]bool{}

	for _, tasks := range dump.Tasks {
		for _, task := range tasks {
			cellId := task.GetCellId()
			if cellId == "" {
				continue
			}
			nTasks[cellId]++
			if !knownCells[cellId] {
				knownCells[cellId] = true
				cells = append(cells, cellId)
			}
		}
	}

	for _, lrp := range dump.LRPS {
		for _, actualLRPGroup := range lrp.ActualLRPGroupsByIndex {
			if actualLRPGroup.Instance != nil {
				cellId := actualLRPGroup.Instance.GetCellId()
				if cellId == "" {
					continue
				}
				if actualLRPGroup.Instance.State == models.ActualLRPStateClaimed {
					nLRPsClaimed[cellId]++
				} else {
					nLRPsRunning[cellId]++
				}

				if !knownCells[cellId] {
					knownCells[cellId] = true
					cells = append(cells, cellId)
				}
			}
			if actualLRPGroup.Evacuating != nil {
				cellId := actualLRPGroup.Evacuating.GetCellId()
				if cellId == "" {
					continue
				}
				nLRPsEvacuating[cellId]++

				if !knownCells[cellId] {
					knownCells[cellId] = true
					cells = append(cells, cellId)
				}
			}
		}
	}

	sort.Strings(cells)

	buffer := &bytes.Buffer{}
	if clear {
		say.Fclear(buffer)
	}
	say.Fprintln(buffer, 0, "Distribution")
	for _, cell := range cells {
		numTasks := nTasks[cell]
		numLRPs := nLRPsClaimed[cell] + nLRPsRunning[cell] + nLRPsEvacuating[cell]
		var content string
		if numTasks == 0 && numLRPs == 0 {
			content = say.Red("Empty")
		} else {
			content = fmt.Sprintf("%s%s%s%s",
				say.Yellow(strings.Repeat("•", nTasks[cell])),
				say.Green(strings.Repeat("•", nLRPsRunning[cell])),
				say.Gray(strings.Repeat("•", nLRPsClaimed[cell])),
				say.Red(strings.Repeat("•", nLRPsEvacuating[cell])),
			)
		}
		say.Fprintln(buffer, 0, "%s: %s", say.Green("%12s", cell), content)
	}

	buffer.WriteTo(os.Stdout)
}
