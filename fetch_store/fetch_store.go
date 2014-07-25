package fetch_store

import (
	"io"

	"github.com/cloudfoundry-incubator/veritas/models"
	"github.com/cloudfoundry/gunk/timeprovider"
	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/cloudfoundry/storeadapter/workerpool"
)

func Fetch(cluster []string, raw bool, w io.Writer) error {
	adapter := etcdstoreadapter.NewETCDStoreAdapter(cluster, workerpool.NewWorkerPool(10))
	err := adapter.Connect()
	if err != nil {
		return err
	}

	store := bbs.NewVeritasBBS(adapter, timeprovider.NewTimeProvider(), steno.NewLogger("veritas"))

	desiredLRPs, err := store.GetAllDesiredLRPs()
	if err != nil {
		return err
	}

	actualLRPS, err := store.GetAllActualLRPs()
	if err != nil {
		return err
	}

	lrpStartAuctions, err := store.GetAllLRPStartAuctions()
	if err != nil {
		return err
	}

	lrpStopAuctions, err := store.GetAllLRPStopAuctions()
	if err != nil {
		return err
	}

	stopLRPInstance, err := store.GetAllStopLRPInstances()
	if err != nil {
		return err
	}

	lrps := models.VeritasLRPS{}

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

}
