package lrps

import (
	"flag"
	"os"
	"time"

	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-cf-experimental/veritas/lrps/set_domain"
)

func SetDomainCommand() common.Command {
	var (
		bbsConfig config_finder.BBSConfig
	)

	flagSet := flag.NewFlagSet("set-domain", flag.ExitOnError)
	bbsConfig.PopulateFlags(flagSet)

	return common.Command{
		Name:        "set-domain",
		Description: "domain ttl - updates the domain in the BBS (ttl is a duration)",
		FlagSet:     flagSet,
		Run: func(args []string) {
			bbsClient, err := config_finder.NewBBS(bbsConfig)
			common.ExitIfError("Could not construct BBS", err)

			if len(args) != 2 {
				say.Fprintln(os.Stderr, 0, say.Red("Expected domain and ttl"))
				os.Exit(1)
			}
			ttl, err := time.ParseDuration(args[1])
			common.ExitIfError("Failed to parse TTL", err)

			err = set_domain.SetDomain(bbsClient, args[0], ttl)
			common.ExitIfError("Failed to submit lrp", err)
		},
	}
}
