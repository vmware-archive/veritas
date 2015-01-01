package print_store

import (
	"github.com/pivotal-cf-experimental/veritas/say"
	"github.com/pivotal-cf-experimental/veritas/veritas_models"
)

func printServices(verbose bool, services veritas_models.VeritasServices) {
	say.Println(0, say.Green("Cells"))
	for _, cell := range services.Cells {
		say.Println(1, "[%s] %s %s %s", say.Cyan(cell.Stack), say.Yellow(cell.Zone), say.Green(cell.CellID), cell.RepAddress)
	}
	say.Println(0, "%s: %s", say.Green("Auctioneer"), services.AuctioneerAddress)
}
