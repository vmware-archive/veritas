package config_finder

import (
	"fmt"
	"os"
)

func FindBBSEndpoint(endpoint string) (string, error) {
	if endpoint != "" {
		return endpoint, nil
	}

	endpoint = os.Getenv("BBS_ENDPOINT")
	if endpoint != "" {
		return endpoint, nil
	}

	return "", fmt.Errorf("You must either specify --bbsEndpoint or set BBS_ENDPOINT")
}
