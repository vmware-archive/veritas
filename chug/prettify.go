package chug

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/veritas/say"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/chug"
)

var colorLookup = map[string]string{
	"executor":       "\x1b[92m",
	"rep":            "\x1b[93m",
	"converger":      "\x1b[94m",
	"auctioneer":     "\x1b[95m",
	"route-emitter":  "\x1b[96m",
	"tps":            "\x1b[97m",
	"nsync-listener": "\x1b[98m",
	"file-server":    "\x1b[34m",
	"router":         "\x1b[32m",
	"loggregator":    "\x1b[33m",
	"stager":         "\x1b[36m",
	"warden-linux":   "\x1b[35m",
}

func Prettify(relativeTime string, data string, hideNonLager bool, minTime time.Time, maxTime time.Time, match *regexp.Regexp, exclude *regexp.Regexp, src io.Reader) error {
	out := ChugWithFilter(src, minTime, maxTime, match, exclude)

	if data != "none" && data != "short" && data != "long" {
		return fmt.Errorf("invalid data specification: %s", data)
	}

	s := &stenographer{
		Data:         data,
		HideNonLager: hideNonLager,
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
	HideNonLager bool
}

func (s *stenographer) PrettyPrint(entry chug.Entry) {
	if isEmptyInigoLog(entry) {
		return
	}

	if !s.Absolute && s.RelativeTime.IsZero() && entry.IsLager {
		s.RelativeTime = entry.Log.Timestamp
	}

	if entry.IsLager {
		s.PrettyPrintLog(entry.Log)
	} else {
		if !s.HideNonLager {
			s.PrettyPrintRaw(entry.Raw)
		}
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
		level = say.Colorize(color, "%-7s", "[INFO]")
	case lager.DEBUG:
		level = say.Gray("%-7s", "[DEBUG]")
	case lager.ERROR:
		level = say.Red("%-7s", "[ERROR]")
	case lager.FATAL:
		level = say.Red("%-7s", "[FATAL]")
	}

	var timestamp string
	if s.Absolute {
		timestamp = log.Timestamp.Format("01/02 15:04:05.00")
	} else {
		timestamp = log.Timestamp.Sub(s.RelativeTime).String()
		timestamp = fmt.Sprintf("%17s", timestamp)
	}

	components = append(components, say.Colorize(color, "%-16s", log.Source))
	components = append(components, level)
	components = append(components, say.Colorize(color, timestamp))
	components = append(components, say.Gray("%-10s", log.Session))
	components = append(components, say.Colorize(color, log.Message))

	say.Println(0, strings.Join(components, " "))

	if log.Error != nil {
		say.Println(27, say.Red("Error: %s", log.Error.Error()))
	}

	if log.Trace != "" {
		say.Println(27, say.Red(log.Trace))
	}

	if len(log.Data) > 0 && s.Data == "short" {
		dataJSON, _ := json.Marshal(log.Data)
		say.Println(27, string(dataJSON))
	}

	if len(log.Data) > 0 && s.Data == "long" {
		dataJSON, _ := json.MarshalIndent(log.Data, "", "  ")
		say.Println(27, string(dataJSON))
	}
}
