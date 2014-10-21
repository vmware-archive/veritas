package config_finder

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/receptor"
)

func FindReceptor(receptorEndpoint, receptorUsername, receptorPassword string) (receptor.Client, error) {
	var endpoint, username, password string

	if receptorEndpoint != "" {
		endpoint = receptorEndpoint
	} else if os.Getenv("RECEPTOR_ENDPOINT") != "" {
		endpoint = os.Getenv("RECEPTOR_ENDPOINT")
	} else {
		return nil, fmt.Errorf("for now, you must either specify --receptorEndpoint or set RECEPTOR_ENDPOINT")
	}

	if receptorUsername != "" {
		username = receptorUsername
	} else if os.Getenv("RECEPTOR_USERNAME") != "" {
		username = os.Getenv("RECEPTOR_USERNAME")
	}

	if receptorPassword != "" {
		password = receptorPassword
	} else if os.Getenv("RECEPTOR_PASSWORD") != "" {
		password = os.Getenv("RECEPTOR_PASSWORD")
	}

	return receptor.NewClient(endpoint, username, password), nil

}
