package chug_commands

import (
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
	go Chug(src, out)

	if data != "none" && data != "short" && data != "full" {
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
		continue
	default:
		seconds, err := strconv.ParseFloat(relativeTime, 64)
		if err != nil {
			return fmt.Errorf("invalid relative time specification: %s", relativeTime)
		}
		s.RelativeTime = time.Unix(0, seconds*1e9)
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
	say.Println(0, say.Gray(raw))
}

func (s *stenographer) PrettyPringLog(log chug.LogEntry) {
	color := say.DefaultStyle

	components := []string{}

	var timestamp string
	if s.Absolute {
		timestamp = log.Timestamp.Format(time.StampMilli)
	} else {
		timestamp = log.Timestamp.Sub(s.RelativeTime).String()
	}

	components = append(components, fmt.Sprintf("%12s", timestamp))
	components = append(components, fmt.Sprintf("[%s]", log.Source))

	messageComponents := []string{}
	for _, task := range log.Tasks {
		messageComponents = append(messageComponents, task)
	}
	messageComponents = append(messageComponents, log.Action)
	message := strings.Join(messageComponents, ".")

	switch log.LogLevel {
	case lager.INFO:

	case lager.DEBUG:
		color = say.LightGrayColor
	case lager.ERROR:
		color = say.RedColor
	case lager.FATAL:
		color = say.RedColor
	}

	say.Println(0, "%12s [%s] %s", timestamp, log.Source, message)

}
