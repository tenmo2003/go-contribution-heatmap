package time

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseRelativeTime(relativeTime string) (from time.Time, to time.Time, err error) {
	relativeTimeParts := regexp.MustCompile(`(?i)(\d+)\s*(day|week|month|year)s?`).FindStringSubmatch(relativeTime)
	if len(relativeTimeParts) != 3 {
		return time.Time{}, time.Time{}, fmt.Errorf("Invalid relative time")
	}

	value, err := strconv.Atoi(relativeTimeParts[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	now := time.Now()

	switch strings.ToLower(relativeTimeParts[2]) {
	case "day":
		return now.AddDate(0, 0, -1*value).AddDate(0, 0, 1), now, nil
	case "week":
		return now.AddDate(0, 0, -7*value).AddDate(0, 0, 1), now, nil
	case "month":
		return now.AddDate(0, -1*value, 0).AddDate(0, 0, 1), now, nil
	case "year":
		return now.AddDate(-1*value, 0, 0).AddDate(0, 0, 1), now, nil
	}

	return time.Time{}, time.Time{}, fmt.Errorf("Invalid relative time unit")
}

func ParseYear(year int) (from time.Time, to time.Time, err error) {
	if year < 0 {
		return time.Time{}, time.Time{}, fmt.Errorf("Year must be positive")
	}
	from = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	to = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	return from, to, nil
}
