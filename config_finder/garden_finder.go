package config_finder

import (
	"fmt"
	"os"
)

func FindGardenAddr(gardenAddr string, gardenNetwork string) (string, string, error) {
	if gardenAddr != "" && gardenNetwork != "" {
		return gardenAddr, gardenNetwork, nil
	}

	if gardenAddr == "" {
		gardenAddr = os.Getenv("GARDEN_ADDR")
	}

	if gardenNetwork == "" {
		gardenNetwork = os.Getenv("GARDEN_NETWORK")
	}

	if gardenAddr != "" && gardenNetwork != "" {
		return gardenAddr, gardenNetwork, nil
	}

	return "", "", fmt.Errorf("For now, you must either specify --gardenAddr and --gardenNetwork or set GARDEN_ADDR and GARDEN_NETWORK")
}
