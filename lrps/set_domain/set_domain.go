package set_domain

import (
	"time"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/onsi/say"
)

func SetDomain(store bbs.VeritasBBS, domain string, ttl time.Duration) error {
	say.Println(0, say.Green("Setting Domain %s with TTL %ds", domain, int(ttl.Seconds())))

	return store.UpsertDomain(domain, int(ttl.Seconds()))
}
