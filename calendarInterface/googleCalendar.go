package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	//"time"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	//"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	//"google.golang.org/api/option"
)


type CalendarEvents struct{
	Events *calendar.Events
}

type NewService struct {
	srv *calendar.Service
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
        tokFile := "token.json"
        tok, err := tokenFromFile(tokFile)
        if err != nil {
                tok = getTokenFromWeb(config)
                saveToken(tokFile, tok)
        }
        return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
                log.Fatalf("Unable to read authorization code: %v", err)
        }

        tok, err := config.Exchange(context.TODO(), authCode)
        if err != nil {
                log.Fatalf("Unable to retrieve token from web: %v", err)
        }
        return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
        f, err := os.Open(file)
        if err != nil {
                return nil, err
        }
        defer f.Close()
        tok := &oauth2.Token{}
        err = json.NewDecoder(f).Decode(tok)
        return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
        fmt.Printf("Saving credential file to: %s\n", path)
        f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
        if err != nil {
                log.Fatalf("Unable to cache oauth token: %v", err)
        }
        defer f.Close()
        json.NewEncoder(f).Encode(token)
}

func loadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	loadEnvVariables()
    ctx := context.Background()

	var events CalendarEvents

	var newService NewService

	ManageCLI()

	newService.srv = NewCalendarService(ctx)

	now := time.Now()
	startTime := now.Format(time.RFC3339)
	fmt.Println(startTime)
	endTime := now.Add(24 * time.Hour).Format(time.RFC3339)

	events.Events = newService.ListEvents(startTime, endTime);


	
	
	events.DeleteEventByName(newService.srv, "Description of Event")

	
}

func NewEvent(summary, location, description string, startTime, endTime string) *calendar.Event {
	
	event := &calendar.Event{
		Summary:     summary,
		Location:    location,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: startTime,
			TimeZone: "UTC", // Adjust the timezone as needed
		},
		End: &calendar.EventDateTime{
			DateTime: endTime,
			TimeZone: "UTC", // Adjust the timezone as needed
		},
		// Add other fields as needed
	}

	return event
}

func CreateConfigFile(creds []byte) *oauth2.Config{
	config, err := google.ConfigFromJSON(creds, calendar.CalendarEventsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	} 
	return config
}


func ReadCredentials() []byte {
	b, err := os.ReadFile(os.Getenv("CREDS_FILEPATH"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	return b
}

func NewCalendarService(ctx context.Context) *calendar.Service {
	creds := ReadCredentials()

	config := CreateConfigFile(creds)

	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	return srv
}


func ManageCLI() {
	args := os.Args[1:]
	if len(args) == 1 {
		log.Println("Running Query from Command Line!")
	}
	log.Printf("No/Too many Queries on the Command Line, app running normally.Commands --> %s\n", args[0])
	return
}

func GetEventByDescription(events *calendar.Events, description string) string{
	for _, event := range events.Items {
		if event.Description == description {
			return event.Id
		}
	}
	log.Printf("No Events with that description found, description:%s\n", description)
	return ""
}

func (service NewService) ListEvents(startTime, endTime string) (*calendar.Events) {
	events, err := service.srv.Events.List("primary").
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


func (events CalendarEvents)DeleteEventByName(srv *calendar.Service, eventDescription string)  {
	eventId := GetEventByDescription(events.Events, eventDescription)

	err := srv.Events.Delete("primary", eventId).Do()
	if err != nil {
		log.Fatalf("Unable to delete event: %v", err)
	}
	log.Println("Event deleted successfully!")
}