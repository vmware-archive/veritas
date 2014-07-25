package config_finder

import "fmt"

func FindETCDCluster(cluster []string) ([]string, err) {
	if len(cluster) == 0 {
		return fmt.Errorf("For now, you must specify an etcd cluster")
	}

	return cluster
}
