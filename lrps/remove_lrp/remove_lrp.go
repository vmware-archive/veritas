package remove_lrp

import (
	"github.com/cloudfoundry-incubator/runtime-schema/bbs"

	"github.com/cloudfoundry/gunk/workpool"
	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
)

func RemoveLRP(cluster []string, guid string) error {
	adapter := etcdstoreadapter.NewETCDStoreAdapter(cluster, workpool.NewWorkPool(10))
	err := adapter.Connect()
	if err != nil {
		return err
	}

	logger := lager.NewLogger("veritas")
	store := bbs.NewVeritasBBS(adapter, clock.NewClock(), logger)

	return store.RemoveDesiredLRPByProcessGuid(logger, guid)
}
