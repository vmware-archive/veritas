package veritas_models

import "github.com/cloudfoundry-incubator/runtime-schema/models"

type StoreDump struct {
	LRPS      VeritasLRPS
	Tasks     VeritasTasks
	Services  VeritasServices
	Freshness []models.Freshness
}
