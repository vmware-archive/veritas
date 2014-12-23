package veritas_models

type StoreDump struct {
	LRPS     VeritasLRPS
	Tasks    VeritasTasks
	Services VeritasServices
	Domains  []string
}
