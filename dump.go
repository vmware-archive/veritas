package main

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
)

func Dump(bbs *bbs.BBS, c Config) {
	c.S.printBanner(fmt.Sprintf("Dump %s", time.Now()), "=")
	if c.Tasks {
		DumpTasks(bbs, c)
	}

	if c.LRPs {
		DumpLRPs(bbs, c)
	}

	if c.Services {
		DumpServices(bbs, c)
	}
}
