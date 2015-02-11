package veritas_models

import (
	"sort"
	"strconv"

	"github.com/cloudfoundry-incubator/runtime-schema/models"
)

type VeritasLRP struct {
	ProcessGuid            string
	DesiredLRP             models.DesiredLRP
	ActualLRPGroupsByIndex map[string]models.ActualLRPGroup
}

func (l *VeritasLRP) OrderedActualLRPIndices() []string {
	indicesAsStrings := []string{}
	for index := range l.ActualLRPGroupsByIndex {
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
