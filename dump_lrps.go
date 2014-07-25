package main

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/onsi/gomega/format"
)

func DumpLRPs(bbs *bbs.BBS, c Config) {
	desiredLRPS, err := bbs.GetAllDesiredLRPs()
	panicIfErr(err)

	actualLRPS, err := bbs.GetAllActualLRPs()
	panicIfErr(err)

	lrpStartAuctions, err := bbs.GetAllLRPStartAuctions()
	panicIfErr(err)

	lrpStopAuctions, err := bbs.GetAllLRPStopAuctions()
	panicIfErr(err)

	stopLRPInstance, err := bbs.GetAllStopLRPInstances()
	panicIfErr(err)

	lrps := LRPS{}

	for _, desired := range desiredLRPS {
		lrps.Get(desired.ProcessGuid).DesiredLRP = desired
	}

	for _, actual := range actualLRPS {
		lrp := lrps.Get(actual.ProcessGuid)
		lrp.ActualLRPsByIndex[actual.Index] = append(lrp.ActualLRPsByIndex[actual.Index], actual)
	}

	for _, startAuction := range lrpStartAuctions {
		lrps.Get(startAuction.ProcessGuid).StopAuctions[startAuction.Index] = startAuction
	}

	for _, stopAuction := range lrpStopAuctions {
		lrps.Get(stopAuction.ProcessGuid).StopAuctions[stopAuction.Index] = stopAuction
	}

	for _, stopInstance := range stopLRPInstance {
		lrp := lrps.Get(stopInstance.ProcessGuid)
		lrp.StopInstances[stopInstance.Index] = append(lrp.StopInstances[stopInstance.Index], stopInstance)
	}

	c.S.printBanner(c.S.colorize("LRPs", greenColor), "~")
	sortedLRPS := lrps.SortedByProcessGuid()
	for _, lrp := range sortedLRPS {
		if c.Verbose {
			dumpVerboseLRP(lrp, c)
		} else {
			dumpLRP(lrp, c)
		}
	}
}

func dumpVerboseLRP(lrp *LRP, c Config) {
	c.S.println(format.Object(lrp, 1))
}

func dumpLRP(lrp *LRP, c Config) {
	c.S.println(1, c.colorize(greenColor, lrp.ProcessGuid))
	c.S.println(
		2,
		"Desired: %d on %s (%d MB, %d MB)",
		lrp.DesiredLRP.Instances,
		lrp.DesiredLRP.Stack,
		lrp.DesiredLRP.MemoryMB,
		lrp.DesiredLRP.DiskMB,
	)

	orderedActualIndices := lrp.OrderedActualLRPIndices()
	for _, index := range orderedActualIndices {
		for i, actual := range lrp.ActualLRPsByIndex[index] {
			prefix := "    "
			if i == 0 {
				prefix = fmt.Sprintf("%3d:", index)
			}
			c.S.println(
				2,
				"%s %s on %s [%s for %s]",
				prefix,
				actual.InstanceGuid,
				actual.ExecutorID,
				time.Since(time.Unix(0, actual.Since)),
				actualState(actual, c),
			)
		}
	}

	orderedStartAuctionIndices := lrp.OrderedStartAuctionIndices()
	if len(orderedStartAuctionIndices) > 0 {
		c.S.println(2, "Start Auctions:")
		for _, index := range orderedStartAuctionIndices {
			auction := lrp.StartAuctions[index]
			c.S.println(
				3,
				"%3d: %s [%s for %s]",
				index,
				auction.InstanceGuid,
				startAuctionState(auction, c),
				time.Since(time.Unix(0, auction.UpdatedAt)),
			)
		}
	}

	orderedStopAuctionIndices := lrp.OrderedStopAuctionIndices()
	if len(orderedStopAuctionIndices) > 0 {
		c.S.println(2, "Stop Auctions:")
		for _, index := range orderedStopAuctionIndices {
			auction := lrp.StopAuctions[index]
			c.S.println(
				3,
				"%3d: [%s for %s]",
				index,
				stopAuctionState(auction, c),
				time.Since(time.Unix(0, auction.UpdatedAt)),
			)
		}
	}

	orderedStopIndices := lrp.OrderdStopIndices()
	if len(orderedStopIndices) > 0 {
		c.S.println(2, "Stop Instances:")
		for _, index := range orderedStopIndices {
			for i, stop := range lrp.StopInstances[index] {
				prefix := "    "
				if i == 0 {
					prefix = fmt.Sprintf("%3d:", index)
				}
				c.S.println(
					3,
					"%s %s",
					prefix,
					stop.InstanceGuid,
				)
			}
		}
	}
}

func actualState(actual models.ActualLRP, c Config) string {
	switch actual.State {
	case models.ActualLRPStateStarting:
		return c.S.colorize(lightGrayColor, "STARTING")
	case models.ActualLRPStateRunning:
		return c.S.colorize(greenColor, "RUNNING")
	default:
		return c.S.colorize(redColor, "RUNNING")
	}
}

func startAuctionState(startAuction models.LRPStartAuction, c Config) string {
	switch startAuction.State {
	case models.LRPStartAuctionStatePending:
		return c.S.colorize(lightGrayColor, "PENDING")
	case models.LRPStartAuctionStateClaimed:
		return c.S.colorize(greenColor, "CLAIMED")
	default:
		return c.S.colorize(redColor, "RUNNING")
	}
}

func stopAuctionState(stopAuction models.LRPStopAuction, c Config) string {
	switch stopAuction.State {
	case models.LRPStopAuctionStatePending:
		return c.S.colorize(lightGrayColor, "PENDING")
	case models.LRPStopAuctionStateClaimed:
		return c.S.colorize(greenColor, "CLAIMED")
	default:
		return c.S.colorize(redColor, "RUNNING")
	}
}
