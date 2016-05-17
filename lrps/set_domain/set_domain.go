package set_domain

import (
	"time"

	"github.com/cloudfoundry-incubator/bbs"
	"github.com/onsi/say"
	"github.com/pivotal-golang/lager"
)

func SetDomain(logger lager.Logger, bbsClient bbs.Client, domain string, ttl time.Duration) error {
	say.Println(0, say.Green("Setting Domain %s with TTL %ds", domain, int(ttl.Seconds())))

	return bbsClient.UpsertDomain(logger, domain, ttl)
}
