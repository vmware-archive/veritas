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

	color := say.DefaultStyle
	ornamentColor := say.GrayColor
	level := "I"

	switch log.LogLevel {
	case lager.INFO:
	case lager.DEBUG:
		color = say.GrayColor
		level = "D"
	case lager.ERROR:
		color = say.RedColor
		ornamentColor = say.RedColor
		level = "E"
	case lager.FATAL:
		color = say.RedColor
		ornamentColor = say.RedColor
		level = "F"
	}

	var timestamp string
	components = append(components, say.Colorize(color, level))
	if s.Absolute {
		timestamp = log.Timestamp.Format(time.StampMilli)
		components = append(components, say.Colorize(color, timestamp))
	} else {
		timestamp = log.Timestamp.Sub(s.RelativeTime).String()
		components = append(components, say.Colorize(color, "%12s", timestamp))
	}
	components = append(components, fmt.Sprintf("%-16s", "["+log.Source+"]"))

	session := log.Session()
	components = append(components, say.Colorize(ornamentColor, "%-10s", session))

	components = append(components, say.Colorize(color, log.Message()))

	if log.Error != nil {
		components = append(components, say.Red(" - Error: "+log.Error.Error()))
	}

	if len(log.Data) > 0 && s.Data == "short" {
		dataJSON, _ := json.Marshal(log.Data)
		components = append(components, say.Colorize(color, string(dataJSON)))
	}

	say.Println(0, strings.Join(components, " "))

	if len(log.Data) > 0 && s.Data == "long" {
		dataJSON, _ := json.MarshalIndent(log.Data, "", "  ")
		say.Println(2, say.Colorize(color, string(dataJSON)))
	}
}
