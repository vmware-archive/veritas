package remove_lrp

import "github.com/cloudfoundry-incubator/bbs"

func RemoveLRP(bbsClient bbs.Client, guid string) error {
	return bbsClient.RemoveDesiredLRP(guid)
}
