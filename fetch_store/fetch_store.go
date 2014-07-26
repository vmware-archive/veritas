package fetch_store

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs/shared"
	"github.com/cloudfoundry-incubator/veritas/models"
	"github.com/cloudfoundry-incubator/veritas/say"
	"github.com/cloudfoundry/gunk/timeprovider"
	"github.com/cloudfoundry/storeadapter"
	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/cloudfoundry/storeadapter/workerpool"
	"github.com/onsi/gomega/format"
)

func Fetch(cluster []string, raw bool, w io.Writer) error {
	adapter := etcdstoreadapter.NewETCDStoreAdapter(cluster, workerpool.NewWorkerPool(10))
	err := adapter.Connect()
	if err != nil {
		return err
	}

	if raw {
		node, err := adapter.ListRecursively(shared.SchemaRoot)
		if err != nil {
			return err
		}
		printNodes(0, node, w)
		return
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

	tasks, err := bbs.GetAllTasks()
	if err != nil {
		return err
	}

	executors, err := bbs.GetAllExecutors()
	if err != nil {
		return err
	}

	fileservers, err := bbs.GetAllFileServers()
	if err != nil {
		return err
	}

	dump := models.StoreDump{
		LRPS:     models.VeritasLRPS{},
		Tasks:    models.VeritasTasks{},
		Services: models.VeritasServices{},
	}

	for _, desired := range desiredLRPS {
		dump.LRPS.Get(desired.ProcessGuid).DesiredLRP = desired
	}

	for _, actual := range actualLRPS {
		lrp := dump.LRPS.Get(actual.ProcessGuid)
		lrp.ActualLRPsByIndex[actual.Index] = append(lrp.ActualLRPsByIndex[actual.Index], actual)
	}

	for _, startAuction := range lrpStartAuctions {
		dump.LRPS.Get(startAuction.ProcessGuid).StopAuctions[startAuction.Index] = startAuction
	}

	for _, stopAuction := range lrpStopAuctions {
		dump.LRPS.Get(stopAuction.ProcessGuid).StopAuctions[stopAuction.Index] = stopAuction
	}

	for _, stopInstance := range stopLRPInstance {
		lrp := dump.LRPS.Get(stopInstance.ProcessGuid)
		lrp.StopInstances[stopInstance.Index] = append(lrp.StopInstances[stopInstance.Index], stopInstance)
	}

	for _, task := range tasks {
		dump.Tasks[task.TaskType] = append(dump.Tasks[task.TaskType], task)
	}

	dump.Services.Executors = executors
	dump.Services.FileServers = fileservers

	encoder := json.NewEncoder(w)
	return encoder.Encode(dump)
}

func printNode(indentation int, node storeadapter.StoreNode, w io.Writer) {
	if node.TTL != 0 {
		say.Fprintln(w, indentation, "%s [%d]", node.Key, node.TTL)
	} else {
		say.Fprintln(w, indentation, node.Key)
	}
	if node.ChildNodes {
		for _, node := range node.ChildNodes {
			say.Fprintln(w, indentation+1, node)
		}
	} else {
		b := bytes.Buffer{}
		err := json.Indent(b, node.Value, "", strings.Repeat(format.Indent, indentation))
		if err == nil {
			b.WriteTo(w)
			say.Fprintln(w, 0, "")
		} else {
			say.Fprintln(w, indentation, string(node.Value))
		}
	}
}
