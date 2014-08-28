package chug

import (
	"regexp"

	"github.com/pivotal-golang/lager/chug"
)

var colorRe = regexp.MustCompile("\x1b" + `\[\d+m`)
var emptyRe = regexp.MustCompile(`\[[oe]\]\[[\w-]+\]\s*$`)

func isEmptyInigoLog(entry chug.Entry) bool {
	return !entry.IsLager && emptyRe.Match(colorless(entry.Raw))
}

func colorless(raw []byte) []byte {
	return colorRe.ReplaceAll(raw, []byte(""))
}
