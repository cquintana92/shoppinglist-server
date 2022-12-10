package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
)

const (
	format = "2006-01-02 15:04:05.999999999 -0700 MST"
)

var (
	replacements []replacement
)

type replacement struct {
	from string
	into string
}

func SetReplacements(input string) error {
	parsed, err := parseReplacements(input)
	if err != nil {
		return err
	}
	replacements = parsed

	return nil
}

func parseReplacements(input string) ([]replacement, error) {
	parts := strings.Split(input, ",")
	replacements := make([]replacement, 0)
	for _, part := range parts {
		replacementParts := strings.Split(part, "=")
		if len(replacementParts) != 2 {
			return nil, errors.New(fmt.Sprintf("invalid replacement, not following the form a=b: %s", part))
		}
		replacements = append(replacements, replacement{
			from: replacementParts[0],
			into: replacementParts[1],
		})
	}
	return replacements, nil
}

func SanitizeName(name string) string {
	return performSanitize(name, replacements)
}

func performSanitize(input string, replacements []replacement) string {
	splitted := strings.Split(input, " ")
	res := ""
	for idx, word := range splitted {
		if idx > 0 {
			res += " "
		}

		wordAsLower := strings.ToLower(word)
		wordRes := word
		for _, replacement := range replacements {
			if wordAsLower == replacement.from {
				wordRes = replacement.into
			}
		}
		res += wordRes
	}

	return CapitalizeName(res)
}

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
