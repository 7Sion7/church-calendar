package pattern

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Month struct{
	Month string `json:"month"`
	Days []Day   `json:"days"`
}

type Day struct {
	DayOfTheMonth     string  `json:"day_of_the_month"`
	WeekDay           string  `json:"week_day"`
	CommemoratedSaint string  `json:"commemorated_saint"`
	Events            []Event  `json:"events"`
}

type Event struct {
	HourOfTheEvent string `json:"hour_of_the_event"`
	Event          string `json:"event"`
}

func FormatData(text string) Month {
	calendar, newMonth := GetCalendar(strings.Fields(text))
	separatedByDay := SeparateByDay(0, calendar)
	for _, dayInfo := range separatedByDay {
		SortDayInfo(dayInfo, &newMonth)
	}
	return newMonth
}

func SortDayInfo(dayInfo []string, month *Month) {
	var day Day
	var commemorated_saint []string

	for index, info := range dayInfo {
		_, err := strconv.Atoi(info)
		if err == nil {
			day.DayOfTheMonth = info
			continue
		}
		if isDayOfWeek(info) {
			day.WeekDay = info
			continue
		}
		if !isHour(info) {
			commemorated_saint = append(commemorated_saint, info)
			continue
		} else {
			day.CommemoratedSaint = strings.Join(commemorated_saint, " ")
			day.Events = GetEventsOfDay(dayInfo[index:])
			break
		}
	}

	month.Days = append(month.Days, day)
}

func GetEventsOfDay(slice []string) []Event{
	var events []Event

	var indexOfTime int
	for i, word := range slice {
		if i == len(slice)-1 {
			events = append(events, Event{
				HourOfTheEvent: slice[indexOfTime],
				Event: strings.Join(slice[indexOfTime+1:], " "),
			})
			break
		}

		if isHour(word) && i != 0 {
			events = append(events, Event{
				HourOfTheEvent: slice[0],
				Event: strings.Join(slice[indexOfTime+1:i], " "),
			})
			indexOfTime = i
		}

	}
	return events
}

func SeparateByDay(index int, calendar []string)  [][]string{
	var separated [][]string
    _, err := strconv.Atoi(calendar[index])
    if err == nil && index != 0 {
        separated = append(separated, calendar[:index])
        separated = append(separated, SeparateByDay(0, calendar[index:])...)
    } else if index == len(calendar)-1 {
        separated = append(separated, calendar)
    } else {
        separated = append(separated, SeparateByDay(index+1, calendar)...)
    }

    return separated
}


func GetCalendar(words []string) ([]string, Month){
	indexOfStartOfCalendar := slices.Index(words, strconv.Itoa(time.Now().Year())) + 1
	month := words[indexOfStartOfCalendar-2]
	month = string(month[0]) + strings.ToLower(month[1:]) //Basically to title func
	newMonth := Month{
		Month: month,
	}
	
	calendar := words[indexOfStartOfCalendar:]
	return calendar, newMonth
}

func isHour(wordAfter string) bool{
	clockForm := regexp.MustCompile(`\b\d{2}:\d{2}\b`)
	if clockForm.MatchString(wordAfter) {
		return true
	}

	return false
}

func isDayOfWeek(word string) bool {
	now := time.Now()
	// Iterate over the next 7 days to get the weekdays
	for i := 0; i < 7; i++ {
		// Calculate the date for the current iteration
		currentDate := now.AddDate(0, 0, i)

		// Get the weekday for the current date
		weekday := currentDate.Weekday().String()

		if word == weekday[:3] {
			return true
		}
	}
	return false
}