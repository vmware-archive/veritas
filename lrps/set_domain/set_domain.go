package set_domain

import (
	"time"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry/gunk/timeprovider"
	"github.com/cloudfoundry/gunk/workpool"
	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/pivotal-cf-experimental/veritas/say"
	"github.com/pivotal-golang/lager"
)

func SetDomain(cluster []string, domain string, ttl time.Duration) error {
	adapter := etcdstoreadapter.NewETCDStoreAdapter(cluster, workpool.NewWorkPool(10))
	err := adapter.Connect()
	if err != nil {
		return err
	}

	logger := lager.NewLogger("veritas")
	store := bbs.NewVeritasBBS(adapter, timeprovider.NewTimeProvider(), logger)

	say.Println(0, say.Green("Setting Domain %s with TTL %ds", domain, int(ttl.Seconds())))

	return store.UpsertDomain(domain, int(ttl.Seconds()))
}
