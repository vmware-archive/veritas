package dump

import "github.com/cloudfoundry-incubator/runtime-schema/bbs"

func DumpServices(bbs *bbs.BBS, c Config) {
	executors, err := bbs.GetAllExecutors()
	panicIfErr(err)

	fileservers, err := bbs.GetAllFileServers()
	panicIfErr(err)

	c.S.printBanner(c.S.colorize("Services", greenColor), "~")

	c.S.println(0, c.S.colorize(greenColor, "Executors"))
	for _, executor := range executors {
		c.S.println(1, "%s (%s)", executor.ExecutorID, executor.Stack)
	}

	c.S.println(0, c.S.colorize(greenColor, "File Servers"))
	for _, fileServer := range fileservers {
		c.S.println(1, fileServer)
	}
}
