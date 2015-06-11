package fetch_store

import (
	"bytes"
	"encoding/json"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/shared"
	"github.com/cloudfoundry/storeadapter"
	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/onsi/gomega/format"
	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/veritas_models"
	"github.com/pivotal-golang/lager"
)

func Fetch(store bbs.VeritasBBS, adapter *etcdstoreadapter.ETCDStoreAdapter, raw bool, w io.Writer) error {
	if raw {
		node, err := adapter.ListRecursively(shared.DataSchemaRoot)
		if err != nil {
			return err
		}
		printNode(0, node, w)
		return nil
	}

	logger := lager.NewLogger("veritas")

	desiredLRPs, err := store.DesiredLRPs(logger)
	if err != nil {
		return err
	}

	actualLRPGroups, err := store.ActualLRPGroups(logger)
	if err != nil {
		return err
	}

	tasks, err := store.Tasks(logger)
	if err != nil {
		return err
	}

	domains, err := store.Domains()
	if err != nil {
		return err
	}

	dump := veritas_models.StoreDump{
		Domains: domains,
		LRPS:    veritas_models.VeritasLRPS{},
		Tasks:   veritas_models.VeritasTasks{},
	}

	for _, desired := range desiredLRPs {
		dump.LRPS.Get(desired.ProcessGuid).DesiredLRP = desired
	}

	for _, actualLRPGroup := range actualLRPGroups {
		actual, _, err := actualLRPGroup.Resolve()
		if err != nil {
			continue
		}
		lrp := dump.LRPS.Get(actual.ProcessGuid)
		index := strconv.Itoa(actual.Index)
		lrp.ActualLRPGroupsByIndex[index] = actualLRPGroup
	}

	for _, task := range tasks {
		dump.Tasks[task.Domain] = append(dump.Tasks[task.Domain], task)
	}

	dump.Services, err = tryToFetchServices(store)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(w)
	return encoder.Encode(dump)
}

func tryToFetchServices(store bbs.VeritasBBS) (services veritas_models.VeritasServices, err error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	services.Cells, err = store.Cells()
	if err != nil {
		return services, err
	}

	sort.Sort(veritas_models.CellsByZoneAndID(services.Cells))

	services.AuctioneerAddress, err = store.AuctioneerAddress()
	if err != nil {
		return services, err
	}

	return services, nil
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
