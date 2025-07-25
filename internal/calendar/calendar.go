package calendar

import (
	"fmt"
	"strings"
)

type FirstWeekday int

const (
	Sunday FirstWeekday = iota // default
	Monday
)

func (d *FirstWeekday) String() string {
	switch *d {
	case Sunday:
		return "sunday"
	case Monday:
		return "monday"
	}
	return "Unknown"
}

func (d *FirstWeekday) Set(value string) error {
	switch strings.ToLower(value) {
	case "sunday":
		*d = Sunday
	case "monday":
		*d = Monday
	default:
		return fmt.Errorf("Invalid first weekday")
	}
	return nil
}

var weekdays = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

func GenerateWeekdayMaps(firstDay FirstWeekday) (map[string]int, map[int]string) {
	WEEKDAY_INDEX := make(map[string]int)
	INDEX_WEEKDAY := make(map[int]string)

	start := int(firstDay)
	for i := range 7 {
		day := weekdays[(start+i)%7]
		WEEKDAY_INDEX[day] = i
		INDEX_WEEKDAY[i] = day
	}

	return WEEKDAY_INDEX, INDEX_WEEKDAY
}
