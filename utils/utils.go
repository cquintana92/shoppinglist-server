package utils

import (
	"time"
	"unicode"
)

const (
	format = "2006-01-02 15:04:05.999999999 -0700 MST"
)

func Now() string {
	now := time.Now()
	return now.Format(format)
}

func CapitalizeName(name string) string {
	for i, v := range name {
		return string(unicode.ToUpper(v)) + name[i+1:]
	}
	return ""
}

func DateFrom(date string) (time.Time, error) {
	return time.Parse(format, date)
}
