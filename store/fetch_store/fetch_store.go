package fetch_store

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/shared"
	"github.com/cloudfoundry/gunk/timeprovider"
	"github.com/cloudfoundry/storeadapter"
	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/onsi/gomega/format"
	"github.com/pivotal-cf-experimental/veritas/say"
	"github.com/pivotal-cf-experimental/veritas/veritas_models"
	"github.com/pivotal-golang/lager"
)

func Fetch(adapter *etcdstoreadapter.ETCDStoreAdapter, raw bool, w io.Writer) error {
	if raw {
		node, err := adapter.ListRecursively(shared.SchemaRoot)
		if err != nil {
			return err
		}
		printNode(0, node, w)
		return nil
	}

	store := bbs.NewVeritasBBS(adapter, timeprovider.NewTimeProvider(), lager.NewLogger("veritas"))

	desiredLRPs, err := store.GetAllDesiredLRPs()
	if err != nil {
		return err
	}

	actualLRPs, err := store.GetAllActualLRPs()
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

	tasks, err := store.GetAllTasks()
	if err != nil {
		return err
	}

	cells, err := store.GetAllCells()
	if err != nil {
		return err
	}

	freshness, err := store.Freshnesses()
	if err != nil {
		return err
	}

	dump := veritas_models.StoreDump{
		Freshness: freshness,
		LRPS:      veritas_models.VeritasLRPS{},
		Tasks:     veritas_models.VeritasTasks{},
		Services:  veritas_models.VeritasServices{},
	}

	for _, desired := range desiredLRPs {
		dump.LRPS.Get(desired.ProcessGuid).DesiredLRP = desired
	}

	for _, actual := range actualLRPs {
		lrp := dump.LRPS.Get(actual.ProcessGuid)
		index := strconv.Itoa(actual.Index)
		lrp.ActualLRPsByIndex[index] = append(lrp.ActualLRPsByIndex[index], actual)
	}

	for _, startAuction := range lrpStartAuctions {
		index := strconv.Itoa(startAuction.Index)
		dump.LRPS.Get(startAuction.DesiredLRP.ProcessGuid).StartAuctions[index] = startAuction
	}

	for _, stopAuction := range lrpStopAuctions {
		index := strconv.Itoa(stopAuction.Index)
		dump.LRPS.Get(stopAuction.ProcessGuid).StopAuctions[index] = stopAuction
	}

	for _, stopInstance := range stopLRPInstance {
		lrp := dump.LRPS.Get(stopInstance.ProcessGuid)
		index := strconv.Itoa(stopInstance.Index)
		lrp.StopInstances[index] = append(lrp.StopInstances[index], stopInstance)
	}

	for _, task := range tasks {
		dump.Tasks[task.Domain] = append(dump.Tasks[task.Domain], task)
	}

	dump.Services.Cells = cells

	encoder := json.NewEncoder(w)
	return encoder.Encode(dump)
}

func printNode(indentation int, node storeadapter.StoreNode, w io.Writer) {
	if node.TTL != 0 {
		say.Fprintln(w, indentation, "%s [%d]", node.Key, node.TTL)
	} else {
		say.Fprintln(w, indentation, node.Key)
	}
	if len(node.ChildNodes) > 0 {
		for _, node := range node.ChildNodes {
			printNode(indentation+1, node, w)
		}
	} else {
		b := &bytes.Buffer{}
		err := json.Indent(b, node.Value, "", strings.Repeat(format.Indent, indentation))
		if err == nil {
			b.WriteTo(w)
			say.Fprintln(w, 0, "")
		} else {
			say.Fprintln(w, indentation, string(node.Value))
		}
	}
}
