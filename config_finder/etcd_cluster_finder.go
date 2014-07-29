package config_finder

import (
	"fmt"
	"os"
	"strings"
)

func FindETCDCluster(cluster string) ([]string, error) {
	if cluster != "" {
		return strings.Split(cluster, ","), nil
	}

	cluster = os.Getenv("ETCD_CLUSTER")
	if cluster != "" {
		return strings.Split(cluster, ","), nil
	}

	return nil, fmt.Errorf("For now, you must either specify an etcd cluster or set ETCD_CLUSTER")
}
