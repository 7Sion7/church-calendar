package googleapi

import (
	"context"
	"fmt"
	funcs "houses-data/basicFunctions"
	pattern "houses-data/patternRecogniser"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func newCalendarService(ctx context.Context, client *http.Client) *calendar.Service {
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	return srv
}


func newDayEvent(event, date, endTime string) *calendar.Event {
	return &calendar.Event{
		Summary: event,
		Reminders: &calendar.EventReminders{
			UseDefault:      false,
			ForceSendFields: []string{"UseDefault"},
		},
		Start: &calendar.EventDateTime{
			DateTime: date,
			TimeZone: zone,
		},
		End: &calendar.EventDateTime{
			DateTime: endTime,
			TimeZone: zone,
		},
	}
}



func (d Date) createDayEvent(srv *calendar.Service, day pattern.Day, image string) { 
	// function for the whole day, with name of commemorated saint as summary
	weekday := d.date.Weekday().String()
	if !strings.Contains(weekday, day.WeekDay) {
		log.Fatalf("Provided weekday (%s) does not match the calculated weekday (%s) for the given date.", day.WeekDay, weekday)
	}
	endTime := d.date.Add(24 * time.Hour)
	
	newEvent := newDayEvent(day.CommemoratedSaint, d.date.Format(time.RFC3339), endTime.Format(time.RFC3339))
	
	newEvent.Description = fmt.Sprintf(`%s: %s`, day.CommemoratedSaint, image)
	newEvent.Reminders.Overrides = []*calendar.EventReminder{
		&calendar.EventReminder{
			Method:  "popup",
			Minutes: -1,
		},
	}
	
	_, err := srv.Events.Insert("primary", newEvent).Do()
	if err != nil {
		log.Fatalf("Unable to create event: %v", err)
	}
}

func (d Date) addDayEvents(srv *calendar.Service, dayEvents pattern.Day) {
	// function for adding events at church and their timings
	for _, e := range dayEvents.Events {

		hourOfTheEvent := e.HourOfTheEvent

		event := e.Event
		
		date := d.date.Add(time.Duration(funcs.ParseTime(hourOfTheEvent)) * time.Hour)
		
		endTime := date.Add(1 * time.Hour)
		
		newEvent := newDayEvent(event, date.Format(time.RFC3339), endTime.Format(time.RFC3339))
		
		newEvent.Description = fmt.Sprintf("Day's Saint: %s,\nEvent: %s at %s", dayEvents.CommemoratedSaint, event, date.Format(time.RFC3339)[11:16])
		
		newEvent.Reminders.Overrides = []*calendar.EventReminder{ //Setting the reminders for the events in the day
			&calendar.EventReminder{
				Method:  "popup",
				Minutes: 5 * 60,
			},
			&calendar.EventReminder{
				Method:  "popup",
				Minutes: 90,
			},
		}
		
		_, err := srv.Events.Insert("primary", newEvent).Do()
		if err != nil {
			log.Fatalf("Unable to create event: %v", err)
		}
	}
}

func getEventDate(now time.Time, m pattern.Month, dayOfMonth string) Date {
	var d Date
	d.date = time.Date(now.Year(), funcs.ParseMonth(m.Month), funcs.Atoi(dayOfMonth), 0, 0, 0, 0, now.Location())
	return d
}

func GetEventByTitle(events *calendar.Events, description string) string {
	for _, event := range events.Items {
		if event.Description == description {
			return event.Id
		}
	}
	log.Printf("No Events with that description found, description:%s\n", description)
	return ""
}

func ListEvents(srv *calendar.Service, startTime, endTime string) *calendar.Events {
	events, err := srv.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(startTime).
		TimeMax(endTime).
		MaxResults(10).
		OrderBy("startTime").
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve events: %v", err)
	}
	return events
}

func DeleteEventByName(srv *calendar.Service, events *calendar.Events, eventDescription string) {
	eventId := GetEventByTitle(events, eventDescription)

	err := srv.Events.Delete("primary", eventId).Do()
	if err != nil {
		log.Fatalf("Unable to delete event: %v", err)
	}
	log.Println("Event deleted successfully!")
}