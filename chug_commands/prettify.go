package chug_commands

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/veritas/say"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/chug"
)

var colorLookup = map[string]string{
	"executor":       "\x1b91m",
	"rep":            "\x1b92m",
	"converger":      "\x1b93m",
	"auctioneer":     "\x1b94m",
	"route-emitter":  "\x1b95m",
	"tps":            "\x1b96m",
	"nsync-listener": "\x1b97m",
	"file-server":    "\x1b90m",
	"router":         "\x1b32m",
	"loggregator":    "\x1b33m",
	"stager":         "\x1b94m",
	"warden-linux":   "\x1b31m",
}

func Prettify(relativeTime string, data string, src io.Reader) error {
	out := make(chan chug.Entry)
	go chug.Chug(src, out)

	if data != "none" && data != "short" && data != "long" {
		return fmt.Errorf("invalid data specification: %s", data)
	}

	s := &stenographer{
		Data: data,
	}

	switch relativeTime {
	case "":
		s.Absolute = true
	case "now":
		s.RelativeTime = time.Now()
	case "first":
		break
	default:
		seconds, err := strconv.ParseFloat(relativeTime, 64)
		if err != nil {
			return fmt.Errorf("invalid relative time specification: %s", relativeTime)
		}
		s.RelativeTime = time.Unix(0, int64(seconds*1e9))
	}

	for entry := range out {
		s.PrettyPrint(entry)
	}

	return nil
}

type stenographer struct {
	Absolute     bool
	RelativeTime time.Time
	Data         string
}

func (s *stenographer) PrettyPrint(entry chug.Entry) {
	if !s.Absolute && s.RelativeTime.IsZero() && entry.IsLager {
		s.RelativeTime = entry.Log.Timestamp
	}

	if entry.IsLager {
		s.PrettyPrintLog(entry.Log)
	} else {
		s.PrettyPrintRaw(entry.Raw)
	}
}

func (s *stenographer) PrettyPrintRaw(raw []byte) {
	say.Println(0, say.Gray(string(raw)))
}

func (s *stenographer) PrettyPrintLog(log chug.LogEntry) {
	components := []string{}

	color, ok := colorLookup[log.Source]
	if !ok {
		color = say.DefaultStyle
	}

	level := ""
	switch log.LogLevel {
	case lager.INFO:
		level = say.Green("%-7s", "[INFO]")
	case lager.DEBUG:
		level = say.Gray("%-7s", "[DEBUG]")
	case lager.ERROR:
		level = say.Red("%-7s", "[ERROR]")
	case lager.FATAL:
		level = say.Red("%-7s", "[FATAL]")
	}

	var timestamp string
	if s.Absolute {
		timestamp = log.Timestamp.Format("01/_2 15:04:05.00")
	} else {
		timestamp = log.Timestamp.Sub(s.RelativeTime).String()
		timestamp = fmt.Sprintf("%17s", timestamp)
	}

	components = append(components, say.Colorize(color, "%-16s", log.Source))
	components = append(components, level)
	components = append(components, timestamp)
	components = append(components, say.Gray("%-10s", log.Session))
	components = append(components, say.Colorize(color, log.Message))

	say.Println(0, strings.Join(components, " "))

	if log.Error != nil {
		say.Println(31, say.Red("Error: %s", log.Error.Error()))
	}

	if log.Trace != "" {
		say.Println(31, say.Red(log.Trace))
	}

	if len(log.Data) > 0 && s.Data == "short" {
		dataJSON, _ := json.Marshal(log.Data)
		say.Println(31, say.Colorize(color, string(dataJSON)))
	}

	if len(log.Data) > 0 && s.Data == "long" {
		dataJSON, _ := json.MarshalIndent(log.Data, "", "  ")
		say.Println(31, say.Colorize(color, string(dataJSON)))
	}
}
