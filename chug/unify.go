package chug

import (
	"io"
	"time"

	"github.com/pivotal-golang/lager/chug"
)

func Unify(files []io.Reader, out io.Writer, minTime time.Time, maxTime time.Time) error {
	chans := make([]chan chug.Entry, len(files))
	for i, file := range files {
		out := ChugWithFilter(file, minTime, maxTime)
		chans[i] = out
	}

	entries := make([]*chug.Entry, len(files))
	for {
		for i, c := range chans {
			if entries[i] != nil {
				continue
			}
			entry, ok := <-c
			if ok {
				entries[i] = &entry
			}
		}

		winningIndex := -1
		winningTime := time.Unix(1e12, 0) //very distant future
		for i, entry := range entries {
			if entry == nil {
				continue
			}
			if !entry.IsLager {
				winningIndex = i
				break
			}
			if entry.Log.Timestamp.Before(winningTime) {
				winningTime = entry.Log.Timestamp
				winningIndex = i
			}
		}

		if winningIndex == -1 {
			return nil
		}

		out.Write(entries[winningIndex].Raw)
		out.Write([]byte("\n"))
		entries[winningIndex] = nil
	}

	return nil
}
