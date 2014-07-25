package models

import (
	"sort"

	"github.com/cloudfoundry-incubator/runtime-schema/models"
)

type VeritasLRP struct {
	ProcessGuid       string
	DesiredLRP        models.DesiredLRP
	ActualLRPsByIndex map[int][]models.ActualLRP
	StartAuctions     map[int]models.LRPStartAuction
	StopAuctions      map[int]models.LRPStopAuction
	StopInstances     map[int][]models.StopLRPInstance
}

func (l *VeritasLRP) OrderedActualLRPIndices() []int {
	indices := []int{}
	for index := range l.ActualLRPsByIndex {
		indices = append(indices, index)
	}
	sort.Ints(indices)
	return indices
}

func (l *VeritasLRP) OrderedStartAuctionIndices() []int {
	indices := []int{}
	for index := range l.StartAuctions {
		indices = append(indices, index)
	}
	sort.Ints(indices)
	return indices
}

func (l *VeritasLRP) OrderedStopAuctionIndices() []int {
	indices := []int{}
	for index := range l.StopAuctions {
		indices = append(indices, index)
	}
	sort.Ints(indices)
	return indices
}

func (l *VeritasLRP) OrderedStopIndices() []int {
	indices := []int{}
	for index := range l.StopInstances {
		indices = append(indices, index)
	}
	sort.Ints(indices)
	return indices
}
