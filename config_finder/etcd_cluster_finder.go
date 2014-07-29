package config_finder

import "fmt"

func FindETCDCluster(cluster []string) ([]string, error) {
	if len(cluster) == 0 {
		return nil, fmt.Errorf("For now, you must specify an etcd cluster")
	}

	return cluster, nil
}
