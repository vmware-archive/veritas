package print_store

import (
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/onsi/gomega/format"
	"github.com/pivotal-cf-experimental/veritas/say"
	"github.com/pivotal-cf-experimental/veritas/veritas_models"
)

func printLRPS(verbose bool, lrps veritas_models.VeritasLRPS) {
	say.PrintBanner(say.Green("LRPs"), "~")

	sortedLRPS := lrps.SortedByProcessGuid()
	for _, lrp := range sortedLRPS {
		if verbose {
			printVerboseLRP(lrp)
		} else {
			printLRP(lrp)
		}
	}
}

func printFreshness(freshnesses []models.Freshness) {
	say.PrintBanner(say.Green("Freshness"), "~")
	if len(freshnesses) == 0 {
		say.Println(1, say.Red("None"))
		return
	}
	for _, freshness := range freshnesses {
		say.Println(1, say.Green("%s - %d", freshness.Domain, freshness.TTLInSeconds))
	}
}

func printVerboseLRP(lrp *veritas_models.VeritasLRP) {
	say.Println(0, format.Object(lrp, 1))
}

func printLRP(lrp *veritas_models.VeritasLRP) {
	say.Println(1, say.Green(lrp.ProcessGuid))
	if lrp.DesiredLRP.ProcessGuid != "" {
		say.Println(
			2,
			"Desired: %s on %s (%d MB, %d MB, %d CPU) %s",
			say.Green("%d", lrp.DesiredLRP.Instances),
			say.Green(lrp.DesiredLRP.Stack),
			lrp.DesiredLRP.MemoryMB,
			lrp.DesiredLRP.DiskMB,
			lrp.DesiredLRP.CPUWeight,
			say.Yellow(strings.Join(lrp.DesiredLRP.Routes, ",")),
		)
	} else {
		say.Println(2, say.Red("UNDESIRED"))
	}

	orderedActualIndices := lrp.OrderedActualLRPIndices()
	for _, index := range orderedActualIndices {
		actual := lrp.ActualLRPsByIndex[index]
		say.Println(
			2,
			"%7s: %s on %s [%s for %s]",
			index,
			actual.InstanceGuid,
			actual.CellID,
			time.Since(time.Unix(0, actual.Since)),
			actualState(actual),
		)
	}
}

func actualState(actual models.ActualLRP) string {
	switch actual.State {
	case models.ActualLRPStateUnclaimed:
		return say.LightGray("UNCLAIMED")
	case models.ActualLRPStateClaimed:
		return say.Yellow("CLAIMED")
	case models.ActualLRPStateRunning:
		return say.Green("RUNNING")
	default:
		return say.Red("INVALID")
	}
}
