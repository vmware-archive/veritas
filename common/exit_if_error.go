package common

import (
	"os"

	"github.com/cloudfoundry-incubator/veritas/say"
)

func ExitIfError(context string, err error) {
	if err != nil {
		say.Fprintln(os.Stderr, 0, say.Red(context))
		say.Fprintln(os.Stderr, 0, say.Red(err.Error()))
		os.Exit(1)
	}
}
