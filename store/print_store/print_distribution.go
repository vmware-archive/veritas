package print_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudfoundry-incubator/runtime-schema/models"

	"github.com/pivotal-cf-experimental/veritas/say"
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

	for _, tasks := range dump.Tasks {
		for _, task := range tasks {
			nTasks[task.CellID]++
		}
	}

	for _, lrp := range dump.LRPS {
		for _, actualLRPGroup := range lrp.ActualLRPGroupsByIndex {
			if actualLRPGroup.Instance != nil {
				if actualLRPGroup.Instance.State == models.ActualLRPStateClaimed {
					nLRPsClaimed[actualLRPGroup.Instance.CellID]++
				} else {
					nLRPsRunning[actualLRPGroup.Instance.CellID]++
				}
			}
			if actualLRPGroup.Evacuating != nil {
				nLRPsEvacuating[actualLRPGroup.Evacuating.CellID]++
			}
		}
	}

	buffer := &bytes.Buffer{}
	if clear {
		say.Fclear(buffer)
	}
	say.Fprintln(buffer, 0, "Distribution")
	for _, cell := range dump.Services.Cells {
		numTasks := nTasks[cell.CellID]
		numLRPs := nLRPsClaimed[cell.CellID] + nLRPsRunning[cell.CellID] + nLRPsEvacuating[cell.CellID]
		var content string
		if numTasks == 0 && numLRPs == 0 {
			content = say.Red("Empty")
		} else {
			content = fmt.Sprintf("%s%s%s",
				say.Yellow(strings.Repeat("•", nTasks[cell.CellID])),
				say.Green(strings.Repeat("•", nLRPsRunning[cell.CellID])),
				say.Gray(strings.Repeat("•", nLRPsClaimed[cell.CellID])),
				say.Red(strings.Repeat("•", nLRPsEvacuating[cell.CellID])),
			)
		}
		say.Fprintln(buffer, 0, "%s %s: %s", say.Yellow(cell.Zone), say.Green("%12s", cell.CellID), content)
	}

	buffer.WriteTo(os.Stdout)
}
