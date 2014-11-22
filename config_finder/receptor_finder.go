package config_finder

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/receptor"
)

func FindReceptor(receptorEndpoint string) (receptor.Client, error) {
	var endpoint string

	if receptorEndpoint != "" {
		endpoint = receptorEndpoint
	} else if os.Getenv("RECEPTOR_ENDPOINT") != "" {
		endpoint = os.Getenv("RECEPTOR_ENDPOINT")
	} else {
		return nil, fmt.Errorf("for now, you must either specify --receptorEndpoint or set RECEPTOR_ENDPOINT")
	}

	return receptor.NewClient(endpoint), nil

}
