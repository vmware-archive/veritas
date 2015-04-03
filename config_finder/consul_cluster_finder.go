package config_finder

import (
	"fmt"
	"os"
)

func FindConsulCluster(cluster string) (string, error) {
	if cluster != "" {
		return cluster, nil
	}

	cluster = os.Getenv("CONSUL_CLUSTER")
	if cluster != "" {
		return cluster, nil
	}

	return "", fmt.Errorf("You must either specify --consulCluster or set CONSUL_CLUSTER")
}
