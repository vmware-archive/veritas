package config_finder

import (
	"fmt"
	"os"
	"strings"
)

func FindVitalsAddrs(addrs string) (map[string]string, error) {
	if addrs == "" {
		addrs = os.Getenv("VITALS_ADDRS")
	}

	if addrs == "" {
		return nil, fmt.Errorf("For now, you must either specify --vitalsAddrs or set VITALS_ADDRS")
	}

	components := strings.Split(addrs, ",")
	out := map[string]string{}
	for _, component := range components {
		subcomponents := strings.SplitN(component, ":", 2)
		out[subcomponents[0]] = subcomponents[1]
	}

	return out, nil
}
