package print_store

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/onsi/gomega/format"
	"github.com/onsi/say"
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

		routesString := routes(lrp.DesiredLRP.Routes)
		if routesString != "" {
			routesString = "\n" + say.Indent(1, routesString)
		}

		say.Println(
			2,
			"%s %s%s (%d MB, %d MB, %d CPU)%s",
			say.Green("%d", lrp.DesiredLRP.Instances),
			say.Cyan(lrp.DesiredLRP.RootFs),
			privileged,
			lrp.DesiredLRP.MemoryMb,
			lrp.DesiredLRP.DiskMb,
			lrp.DesiredLRP.CpuWeight,
			routesString,
		)
	} else {
		say.Println(2, say.Red("UNDESIRED"))
	}

	orderedActualIndices := lrp.OrderedActualLRPIndices()
	for _, index := range orderedActualIndices {
		actualLRPGroup := lrp.ActualLRPGroupsByIndex[index]
		if instance := actualLRPGroup.Instance; instance != nil {
			if instance.State == models.ActualLRPStateUnclaimed || instance.State == models.ActualLRPStateCrashed {
				say.Println(
					3,
					"%2s: [%s for %s]",
					index,
					actualState(instance),
					time.Since(time.Unix(0, instance.Since)),
				)
			} else {
				say.Println(
					3,
					"%2s: %s %s [%s for %s]",
					index,
					instance.InstanceGuid,
					say.Yellow(instance.CellId),
					actualState(instance),
					time.Since(time.Unix(0, instance.Since)),
				)
			}
		}
		if evacuating := actualLRPGroup.Evacuating; evacuating != nil {
			say.Println(
				3,
				"%s: %s %s [%s for %s] - %s",
				say.Red("%2s", index),
				say.Red(evacuating.InstanceGuid),
				say.Yellow(evacuating.CellId),
				actualState(evacuating),
				time.Since(time.Unix(0, evacuating.Since)),
				say.Red("EVACUATING"),
			)
		}
	}
}

func actualState(actual *models.ActualLRP) string {
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
		return say.Red("CRASHED (%d - %s)", actual.CrashCount, strings.Replace(actual.CrashReason, "\n", " ", -1))
	default:
		return say.Red("INVALID")
	}
}

const CF_ROUTER = "cf-router"

type CFRoutes []CFRoute

type CFRoute struct {
	Hostnames []string `json:"hostnames"`
	Port      uint16   `json:"port"`
}

func routes(info *models.Routes) string {
	if info == nil {
		return ""
	}

	data, found := (*info)[CF_ROUTER]
	if !found || data == nil {
		return ""
	}

	routes := CFRoutes{}
	err := json.Unmarshal(*data, &routes)

	if err != nil {
		return ""
	}

	out := ""

	for _, route := range routes {
		out += fmt.Sprintf("%s => %s ", say.Yellow("%d", route.Port), say.Green(strings.Join(route.Hostnames, " ")))
	}

	return out
}
