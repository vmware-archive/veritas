package veritas_models

import "github.com/cloudfoundry-incubator/runtime-schema/models"

type VeritasServices struct {
	Cells             []models.CellPresence
	AuctioneerAddress string
}
