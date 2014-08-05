package config_finder

import (
	"fmt"
	"os"
)

func FindLoggregatorAddr(loggregatorAddr string) (string, error) {
	if loggregatorAddr != "" {
		return loggregatorAddr, nil
	}

	loggregatorAddr = os.Getenv("LOGGREGATOR_ADDR")
	if loggregatorAddr != "" {
		return loggregatorAddr, nil
	}

	return "", fmt.Errorf("For now, you must either specify --loggregatorAddr or set LOGGREGATOR_ADDR")
}
