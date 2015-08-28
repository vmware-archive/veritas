package fetch_store

import (
	"encoding/json"
	"io"
	"strconv"

	"github.com/cloudfoundry-incubator/bbs"
	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/pivotal-cf-experimental/veritas/veritas_models"
)

func Fetch(bbsClient bbs.Client, w io.Writer) error {
	desiredLRPs, err := bbsClient.DesiredLRPs(models.DesiredLRPFilter{})
	if err != nil {
		return err
	}

	actualLRPGroups, err := bbsClient.ActualLRPGroups(models.ActualLRPFilter{})
	if err != nil {
		return err
	}

	tasks, err := bbsClient.Tasks()
	if err != nil {
		return err
	}

	domains, err := bbsClient.Domains()
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
		actual, _ := actualLRPGroup.Resolve()
		lrp := dump.LRPS.Get(actual.ProcessGuid)
		index := strconv.Itoa(int(actual.Index))
		lrp.ActualLRPGroupsByIndex[index] = actualLRPGroup
	}

	for _, task := range tasks {
		dump.Tasks[task.Domain] = append(dump.Tasks[task.Domain], task)
	}

	encoder := json.NewEncoder(w)
	return encoder.Encode(dump)
}
