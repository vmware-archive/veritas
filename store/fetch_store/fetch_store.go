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

	desiredLRPs, err := store.DesiredLRPs()
	if err != nil {
		return err
	}

	actualLRPs, err := store.ActualLRPs()
	if err != nil {
		return err
	}

	tasks, err := store.Tasks()
	if err != nil {
		return err
	}

	cells, err := store.Cells()
	if err != nil {
		return err
	}

	auctioneerAddress, err := store.AuctioneerAddress()
	if err != nil {
		return err
	}

	domains, err := store.Domains()
	if err != nil {
		return err
	}

	dump := veritas_models.StoreDump{
		Domains:  domains,
		LRPS:     veritas_models.VeritasLRPS{},
		Tasks:    veritas_models.VeritasTasks{},
		Services: veritas_models.VeritasServices{},
	}

	for _, desired := range desiredLRPs {
		dump.LRPS.Get(desired.ProcessGuid).DesiredLRP = desired
	}

	for _, actual := range actualLRPs {
		lrp := dump.LRPS.Get(actual.ProcessGuid)
		index := strconv.Itoa(actual.Index)
		lrp.ActualLRPsByIndex[index] = actual
	}

	for _, task := range tasks {
		dump.Tasks[task.Domain] = append(dump.Tasks[task.Domain], task)
	}

	dump.Services.Cells = cells
	dump.Services.AuctioneerAddress = auctioneerAddress

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
