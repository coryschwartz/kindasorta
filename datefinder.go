package main

import (
	"encoding/json"
	"time"

	"github.com/araddon/dateparse"
)

type datefinder func(string) (time.Time, bool)

var (
	DATE_FINDER_ORDER = []datefinder{
		dateFindZapJson,
		dateFindRecursive,
	}
)

// expensive search for a date somewhere in the string
func dateFindRecursive(datestr string) (time.Time, bool) {
	if len(datestr) == 0 {
		return time.Time{}, false
	}
	t, err := dateparse.ParseAny(datestr)
	if err != nil {
		return dateFindRecursive(datestr[1:])
	}
	return t, true
}

// go.uber.com/zap
// by default, NewProductionEncoderConfig will have the timestamp keyed by "ts"
// and NewDevelopmentEncoderConfig will use "T" as the TimeKey.
// Try to find either of those keys and decode the date from the value
func dateFindZapJson(datestr string) (time.Time, bool) {
	var anyjson map[string]interface{}
	if err := json.Unmarshal([]byte(datestr), &anyjson); err != nil {
		return time.Time{}, false
	}
	for _, timeKey := range []string{"ts", "T"} {
		if ts, ok := anyjson[timeKey]; ok {
			return dateFindRecursive(ts.(string))
		}

	}
	return time.Time{}, false
}

// Try to detect the date using several methods.
func findDate(datestr string) time.Time {
	for _, fdr := range DATE_FINDER_ORDER {
		if t, found := fdr(datestr); found {
			return t
		}
	}
	// no date found. This will sort low and print the line as quickly as possible.
	return time.Time{}
}
