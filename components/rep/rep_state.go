package rep

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/onsi/say"

	"github.com/cloudfoundry-incubator/auction/auctiontypes"
	"github.com/cloudfoundry-incubator/auction/communication/http/auction_http_client"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/cloudfoundry/gunk/workpool"
	"github.com/pivotal-golang/lager"
)

type stateResponse struct {
	ID       string
	Duration time.Duration
	Error    error
	State    auctiontypes.CellState
}

type stateResponses []stateResponse

func (r stateResponses) Len() int           { return len(r) }
func (r stateResponses) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r stateResponses) Less(i, j int) bool { return r[i].ID < r[j].ID }

func RepState(bbs bbs.VeritasBBS, out io.Writer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panicked -- you most likely aren't pointing at consul correctly")
			return
		}
	}()

	cells, err := bbs.Cells()
	if err != nil {
		return err
	}

	states := fetchCellStates(cells)
	sort.Sort(states)

	for _, state := range states {
		if state.Error != nil {
			say.Println(0, "%s [%s] - Error:%s", say.Red(state.ID), state.Duration, say.Red(state.Error.Error()))
			continue
		}

		name := say.Green(state.ID)
		if state.State.Evacuating {
			name = say.Red("%s - EVAC -", state.ID)
		}

		say.Println(0, "%s [%s] - Zone:%s C:%d/%d M:%d/%d D:%d/%d",
			name,
			state.Duration,
			say.Cyan(state.State.Zone),
			state.State.AvailableResources.Containers,
			state.State.TotalResources.Containers,
			state.State.AvailableResources.MemoryMB,
			state.State.TotalResources.MemoryMB,
			state.State.AvailableResources.DiskMB,
			state.State.TotalResources.DiskMB,
		)

		rootFSes := []string{}
		for key := range state.State.RootFSProviders {
			if key != "preloaded" {
				rootFSes = append(rootFSes, say.Yellow(key))
			}
		}

		for key := range state.State.RootFSProviders["preloaded"].(auctiontypes.FixedSetRootFSProvider).FixedSet {
			rootFSes = append(rootFSes, say.Green("preloaded:%s", key))
		}

		say.Println(1, "%s Tasks, %s LRPs, %s",
			say.Cyan("%d", len(state.State.Tasks)),
			say.Cyan("%d", len(state.State.LRPs)),
			strings.Join(rootFSes, ","),
		)
	}

	return nil
}

func fetchCellStates(cells []models.CellPresence) stateResponses {
	lock := &sync.Mutex{}
	responses := []stateResponse{}

	wp := workpool.NewWorkPool(20)
	wg := &sync.WaitGroup{}
	wg.Add(len(cells))
	for _, cell := range cells {
		cell := cell
		wp.Submit(func() {
			defer wg.Done()
			client := auction_http_client.New(&http.Client{
				Timeout: 5 * time.Second,
			}, cell.CellID, cell.RepAddress, lager.NewLogger("veritas"))

			t := time.Now()
			state, err := client.State()
			dt := time.Since(t)
			lock.Lock()
			responses = append(responses, stateResponse{
				ID:       cell.CellID,
				Duration: dt,
				Error:    err,
				State:    state,
			})
			lock.Unlock()
		})
	}
	wg.Wait()

	return responses
}
