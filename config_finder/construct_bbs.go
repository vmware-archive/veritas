package config_finder

import "github.com/cloudfoundry-incubator/bbs"

func ConstructBBS(bbsEndpointFlag string) (bbs.Client, error) {
	bbsEndpoint, err := FindBBSEndpoint(bbsEndpointFlag)
	if err != nil {
		return nil, err
	}

	return bbs.NewClient(bbsEndpoint), nil
}
