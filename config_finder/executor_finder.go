package config_finder

import (
	"fmt"
	"os"
)

func FindExecutorAddr(executorAddr string) (string, error) {
	if executorAddr != "" {
		return executorAddr, nil
	}

	executorAddr = os.Getenv("EXECUTOR_ADDR")
	if executorAddr != "" {
		return executorAddr, nil
	}

	return "", fmt.Errorf("For now, you must either specify --executorAddr or set EXECUTOR_ADDR")
}
