package config_finder

import (
	"github.com/cloudfoundry-incubator/consuladapter"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry/gunk/workpool"
	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
)

func ConstructBBS(etcdClusterFlag string, consulClusterFlag string) (bbs.VeritasBBS, *etcdstoreadapter.ETCDStoreAdapter, error) {
	etcdCluster, err := FindETCDCluster(etcdClusterFlag)
	if err != nil {
		return nil, nil, err
	}

	etcdAdapter := etcdstoreadapter.NewETCDStoreAdapter(etcdCluster, workpool.NewWorkPool(10))
	err = etcdAdapter.Connect()
	if err != nil {
		return nil, nil, err
	}

	consulAdapter, err := constructConsulAdapter(consulClusterFlag)
	if err != nil {
		return nil, nil, err
	}

	logger := lager.NewLogger("veritas")
	store := bbs.NewVeritasBBS(etcdAdapter, consulAdapter, clock.NewClock(), logger)
	return store, etcdAdapter, nil
}

func constructConsulAdapter(consulClusterFlag string) (*consuladapter.Adapter, error) {
	consulCluster, err := FindConsulCluster(consulClusterFlag)
	if err != nil {
		return nil, nil
	}

	consulScheme, consulAddresses, err := consuladapter.Parse(consulCluster)
	if err != nil {
		return nil, err
	}

	consulAdapter, err := consuladapter.NewAdapter(consulAddresses, consulScheme)
	if err != nil {
		return nil, err
	}

	return consulAdapter, nil
}
