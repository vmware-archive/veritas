package veritas_models

import (
	"sort"
	"strconv"

	"github.com/cloudfoundry-incubator/runtime-schema/models"
)

type VeritasLRP struct {
	ProcessGuid       string
	DesiredLRP        models.DesiredLRP
	ActualLRPsByIndex map[string][]models.ActualLRP
	StartAuctions     map[string]models.LRPStartAuction
	StopAuctions      map[string]models.LRPStopAuction
	StopInstances     map[string][]models.StopLRPInstance
}

func (l *VeritasLRP) OrderedActualLRPIndices() []string {
	indicesAsStrings := []string{}
	for index := range l.ActualLRPsByIndex {
		indicesAsStrings = append(indicesAsStrings, index)
	}

	sort.Sort(ByNumericalValue(indicesAsStrings))
	return indicesAsStrings
}

func (l *VeritasLRP) OrderedStartAuctionIndices() []string {
	indicesAsStrings := []string{}
	for index := range l.StartAuctions {
		indicesAsStrings = append(indicesAsStrings, index)
	}

	sort.Sort(ByNumericalValue(indicesAsStrings))
	return indicesAsStrings
}

func (l *VeritasLRP) OrderedStopAuctionIndices() []string {
	indicesAsStrings := []string{}
	for index := range l.StopAuctions {
		indicesAsStrings = append(indicesAsStrings, index)
	}

	sort.Sort(ByNumericalValue(indicesAsStrings))
	return indicesAsStrings
}

func (l *VeritasLRP) OrderedStopIndices() []string {
	indicesAsStrings := []string{}
	for index := range l.StopInstances {
		indicesAsStrings = append(indicesAsStrings, index)
	}

	sort.Sort(ByNumericalValue(indicesAsStrings))
	return indicesAsStrings
}

type ByNumericalValue []string

func (a ByNumericalValue) Len() int      { return len(a) }
func (a ByNumericalValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByNumericalValue) Less(i, j int) bool {
	ai, _ := strconv.Atoi(a[i])
	aj, _ := strconv.Atoi(a[j])

	return ai < aj
}
