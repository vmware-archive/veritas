package config_finder

import (
	"fmt"
	"os"
)

func FindWardenAddr(wardenAddr string) (string, error) {
	if wardenAddr != "" {
		return wardenAddr, nil
	}

	wardenAddr = os.Getenv("WARDEN_ADDR")
	if wardenAddr != "" {
		return wardenAddr, nil
	}

	return "", fmt.Errorf("For now, you must either specify --wardenAddr or set WARDEN_ADDR")
}
