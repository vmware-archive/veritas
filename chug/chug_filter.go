package chug

import (
	"io"
	"regexp"
	"time"

	"github.com/pivotal-golang/lager/chug"
)

func ChugWithFilter(reader io.Reader, minTime time.Time, maxTime time.Time, match *regexp.Regexp, exclude *regexp.Regexp) chan chug.Entry {
	chugOut := make(chan chug.Entry)
	go chug.Chug(reader, chugOut)
	filteredOut := make(chan chug.Entry)
	go filter(filteredOut, chugOut, minTime, maxTime, match, exclude)
	return filteredOut
}

func filter(out chan<- chug.Entry, in <-chan chug.Entry, minTime time.Time, maxTime time.Time, match *regexp.Regexp, exclude *regexp.Regexp) {
	defer close(out)
	lastLagerTime := time.Time{}
	for entry := range in {
		if entry.IsLager {
			lastLagerTime = entry.Log.Timestamp
		}

		if isAfterMin(minTime, lastLagerTime) && isBeforeMax(maxTime, lastLagerTime) && matchesFilters(entry.Raw, match, exclude) {
			out <- entry
		}

		if !isBeforeMax(maxTime, lastLagerTime) {
			return
		}
	}
}

func isAfterMin(minTime time.Time, lastLagerTime time.Time) bool {
	if minTime.IsZero() {
		return true
	}
	if lastLagerTime.IsZero() {
		return false
	}
	return lastLagerTime.After(minTime) || lastLagerTime.Equal(minTime)
}

func isBeforeMax(maxTime time.Time, lastLagerTime time.Time) bool {
	if maxTime.IsZero() {
		return true
	}
	if lastLagerTime.IsZero() {
		return true
	}
	return lastLagerTime.Before(maxTime) || lastLagerTime.Equal(maxTime)
}

func matchesFilters(raw []byte, match *regexp.Regexp, exclude *regexp.Regexp) bool {
	matches := true

	if match != nil {
		matches = match.Match(raw)
	}

	if exclude != nil {
		matches = matches && !exclude.Match(raw)
	}

	return matches
}
