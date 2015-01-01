package veritas_models

import "github.com/cloudfoundry-incubator/runtime-schema/models"

type VeritasServices struct {
	Cells             []models.CellPresence
	AuctioneerAddress string
}

type CellsByZoneAndID []models.CellPresence

func (a CellsByZoneAndID) Len() int      { return len(a) }
func (a CellsByZoneAndID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a CellsByZoneAndID) Less(i, j int) bool {
	if a[i].Zone == a[j].Zone {
		return a[i].CellID < a[j].CellID
	} else {
		return a[i].Zone < a[j].Zone
	}
}
