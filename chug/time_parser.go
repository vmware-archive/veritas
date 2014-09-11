package chug

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func ParseTimeFlag(t string) (time.Time, error) {
	if t == "" {
		return time.Time{}, nil
	}
	tAsFloat, err := strconv.ParseFloat(t, 64)
	if err == nil {
		seconds := int64(math.Ceil(tAsFloat))
		nanoseconds := int64((tAsFloat - float64(seconds)) * 1e9)
		return time.Unix(seconds, nanoseconds), nil
	}

	duration, err := time.ParseDuration(t)
	if err == nil {
		return time.Now().Add(duration), err
	}

	out, err := time.Parse("01/02 15:04:05.00", t)
	if err == nil {
		return out, err
	}

	return time.Time{}, fmt.Errorf("Invalid time: %s", t)
}
