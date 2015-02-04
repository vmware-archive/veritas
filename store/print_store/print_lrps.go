package print_store

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/receptor"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/onsi/gomega/format"
	"github.com/pivotal-cf-experimental/veritas/say"
	"github.com/pivotal-cf-experimental/veritas/veritas_models"
)

func printLRPS(verbose bool, lrps veritas_models.VeritasLRPS) {
	say.Println(0, say.Green("LRPs"))

	sortedLRPS := lrps.SortedByProcessGuid()
	for _, lrp := range sortedLRPS {
		if verbose {
			printVerboseLRP(lrp)
		} else {
			printLRP(lrp)
		}
	}
}

func printDomains(domains []string) {
	say.Println(0, say.Green("Domains"))
	if len(domains) == 0 {
		say.Println(1, say.Red("None"))
		return
	}
	for _, domain := range domains {
		say.Println(1, say.Green("%s", domain))
	}
}

func printVerboseLRP(lrp *veritas_models.VeritasLRP) {
	say.Println(0, format.Object(lrp, 1))
}

func printLRP(lrp *veritas_models.VeritasLRP) {
	say.Println(1, say.Green(lrp.ProcessGuid))
	if lrp.DesiredLRP.ProcessGuid != "" {
		privileged := ""
		if lrp.DesiredLRP.Privileged {
			privileged = say.Red(" PRIVILEGED")
		}
		say.Println(
			2,
			"%s on %s%s (%d MB, %d MB, %d CPU)\n%s",
			say.Green("%d", lrp.DesiredLRP.Instances),
			say.Green(lrp.DesiredLRP.Stack),
			privileged,
			lrp.DesiredLRP.MemoryMB,
			lrp.DesiredLRP.DiskMB,
			lrp.DesiredLRP.CPUWeight,
			say.Indent(1, routes(lrp.DesiredLRP.Routes)),
		)
	} else {
		say.Println(2, say.Red("UNDESIRED"))
	}

	orderedActualIndices := lrp.OrderedActualLRPIndices()
	for _, index := range orderedActualIndices {
		actual := lrp.ActualLRPsByIndex[index]
		if actual.State == models.ActualLRPStateUnclaimed || actual.State == models.ActualLRPStateCrashed {
			say.Println(
				3,
				"%2s: [%s for %s]",
				index,
				actualState(actual),
				time.Since(time.Unix(0, actual.Since)),
			)
		} else {
			say.Println(
				3,
				"%2s: %s %s [%s for %s]",
				index,
				actual.InstanceGuid,
				say.Yellow(actual.CellID),
				actualState(actual),
				time.Since(time.Unix(0, actual.Since)),
			)
		}
	}
}

func actualState(actual models.ActualLRP) string {
	switch actual.State {
	case models.ActualLRPStateUnclaimed:
		if actual.PlacementError == "" {
			return say.LightGray("UNCLAIMED")
		} else {
			return say.Red("UNCLAIMED (%s)", actual.PlacementError)
		}
	case models.ActualLRPStateClaimed:
		return say.Yellow("CLAIMED")
	case models.ActualLRPStateRunning:
		return say.Green("RUNNING")
	case models.ActualLRPStateCrashed:
		return say.Red("CRASHED (%d)", actual.CrashCount)
	default:
		return say.Red("INVALID")
	}
}

func routes(info map[string]*json.RawMessage) string {
	if info == nil {
		return ""
	}

	if info[receptor.CFRouter] == nil {
		return ""
	}

	var routerRoutes receptor.CFRoutes
	json.Unmarshal(*info[receptor.CFRouter], &routerRoutes)

	out := ""

	for _, cfRoute := range routerRoutes {
		out += fmt.Sprintf("%s => %s ", say.Yellow("%d", cfRoute.Port), say.Green(strings.Join(cfRoute.Hostnames, " ")))
	}

	return out
}
