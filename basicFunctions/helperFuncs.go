package basicfunctions

import (
	"log"
	"strconv"
	"time"
)

func ParseMonth(month string) time.Month {
	months := map[string]time.Month{
		"January":   time.January,
		"February":  time.February,
		"March":     time.March,
		"April":     time.April,
		"May":       time.May,
		"June":      time.June,
		"July":      time.July,
		"August":    time.August,
		"September": time.September,
		"October":   time.October,
		"November":  time.November,
		"December":  time.December,
	}

	// Return the time.Month value
	return months[month]
}

// Helper function to parse time from string
func ParseTime(s string) int {
	t, err := time.Parse("15:04", s)
	if err != nil {
		log.Fatalf("Error parsing time: %v", err)
	}
	return t.Hour()
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Error converting string to int: %v", err)
	}
	return i
}
