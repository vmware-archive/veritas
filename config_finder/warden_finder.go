package config_finder

import (
	"fmt"
	"os"
)

func FindWardenAddr(wardenAddr string, wardenNetwork string) (string, string, error) {
	if wardenAddr != "" && wardenNetwork != "" {
		return wardenAddr, wardenNetwork, nil
	}

	if wardenAddr == "" {
		wardenAddr = os.Getenv("WARDEN_ADDR")
	}

	if wardenNetwork == "" {
		wardenNetwork = os.Getenv("WARDEN_NETWORK")
	}

	if wardenAddr != "" && wardenNetwork != "" {
		return wardenAddr, wardenNetwork, nil
	}

	return "", "", fmt.Errorf("For now, you must either specify --wardenAddr and --wardenNetwork or set WARDEN_ADDR and WARDEN_NETWORK")
}
