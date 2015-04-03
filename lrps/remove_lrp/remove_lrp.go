package remove_lrp

import (
	"github.com/cloudfoundry-incubator/runtime-schema/bbs"

	"github.com/pivotal-golang/lager"
)

func RemoveLRP(store bbs.VeritasBBS, guid string) error {
	logger := lager.NewLogger("veritas")
	return store.RemoveDesiredLRPByProcessGuid(logger, guid)
}
