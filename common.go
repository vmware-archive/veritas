package veritas

import (
	"os"

	"github.com/cloudfoundry-incubator/veritas/say"
)

func ExitIfError(context string, err error) {
	if err != nil {
		say.Fprintln(say.Red(context))
		say.Fprintln(say.Red(err.Error()))
		os.Exit(1)
	}
}
