package print_store

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/cloudfoundry-incubator/veritas/say"
	"github.com/cloudfoundry-incubator/veritas/veritas_models"
	"github.com/onsi/gomega/format"
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

func printVerboseLRP(lrp *veritas_models.VeritasLRP) {
	say.Println(0, format.Object(lrp, 1))
}

func printLRP(lrp *veritas_models.VeritasLRP) {
	say.Println(1, say.Green(lrp.ProcessGuid))
	if lrp.DesiredLRP.ProcessGuid != "" {
		say.Println(
			2,
			"Desired: %d on %s (%d MB, %d MB)",
			lrp.DesiredLRP.Instances,
			lrp.DesiredLRP.Stack,
			lrp.DesiredLRP.MemoryMB,
			lrp.DesiredLRP.DiskMB,
		)
	} else {
		say.Println(2, say.Red("UNDESIRED"))
	}

	orderedActualIndices := lrp.OrderedActualLRPIndices()
	for _, index := range orderedActualIndices {
		for i, actual := range lrp.ActualLRPsByIndex[index] {
			prefix := "    "
			if i == 0 {
				prefix = fmt.Sprintf("%3s:", index)
			}
			say.Println(
				2,
				"%s %s on %s [%s for %s]",
				prefix,
				actual.InstanceGuid,
				actual.ExecutorID,
				time.Since(time.Unix(0, actual.Since)),
				actualState(actual),
			)
		}
	}

	orderedStartAuctionIndices := lrp.OrderedStartAuctionIndices()
	if len(orderedStartAuctionIndices) > 0 {
		say.Println(2, "Start Auctions:")
		for _, index := range orderedStartAuctionIndices {
			auction := lrp.StartAuctions[index]
			say.Println(
				3,
				"%3s: %s [%s for %s]",
				index,
				auction.InstanceGuid,
				startAuctionState(auction),
				time.Since(time.Unix(0, auction.UpdatedAt)),
			)
		}
	}

	orderedStopAuctionIndices := lrp.OrderedStopAuctionIndices()
	if len(orderedStopAuctionIndices) > 0 {
		say.Println(2, "Stop Auctions:")
		for _, index := range orderedStopAuctionIndices {
			auction := lrp.StopAuctions[index]
			say.Println(
				3,
				"%3s: [%s for %s]",
				index,
				stopAuctionState(auction),
				time.Since(time.Unix(0, auction.UpdatedAt)),
			)
		}
	}

	orderedStopIndices := lrp.OrderedStopIndices()
	if len(orderedStopIndices) > 0 {
		say.Println(2, "Stop Instances:")
		for _, index := range orderedStopIndices {
			for i, stop := range lrp.StopInstances[index] {
				prefix := "    "
				if i == 0 {
					prefix = fmt.Sprintf("%3s:", index)
				}
				say.Println(
					3,
					"%s %s",
					prefix,
					stop.InstanceGuid,
				)
			}
		}
	}
}

func actualState(actual models.ActualLRP) string {
	switch actual.State {
	case models.ActualLRPStateStarting:
		return say.LightGray("STARTING")
	case models.ActualLRPStateRunning:
		return say.Green("RUNNING")
	default:
		return say.Red("INVALID")
	}
}

func startAuctionState(startAuction models.LRPStartAuction) string {
	switch startAuction.State {
	case models.LRPStartAuctionStatePending:
		return say.LightGray("PENDING")
	case models.LRPStartAuctionStateClaimed:
		return say.Green("CLAIMED")
	default:
		return say.Red("INVALID")
	}
}

func stopAuctionState(stopAuction models.LRPStopAuction) string {
	switch stopAuction.State {
	case models.LRPStopAuctionStatePending:
		return say.LightGray("PENDING")
	case models.LRPStopAuctionStateClaimed:
		return say.Green("CLAIMED")
	default:
		return say.Red("INVALID")
	}
}
