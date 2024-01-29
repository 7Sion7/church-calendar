package googleapi

import (
	"context"
	funcs "houses-data/basicFunctions"
	pattern "houses-data/patternRecogniser"
	"log"
	"time"
	//"golang.org/x/oauth2/google"
	//"google.golang.org/api/option"
)


type Date struct {
	date time.Time
}

var zone, _ = time.Now().Local().Zone()

const (
	tokFile = "token.json"
)


func CallAPI(month pattern.Month) {

	funcs.LoadEnvVariables()
	
    ctx := context.Background()

	creds := ReadCredentials()

	config := CreateConfigFile(creds)

	client := getClient(config)

	calendarSrv := newCalendarService(ctx, client)

	searchSrv := newSearchService(ctx, client)

	// Calculate date
	now := time.Now()
	for _, day := range month.Days {
        ImgURL := getImageURL(searchSrv, day.CommemoratedSaint + " orthodox icon")

        d := getEventDate(now, month, day.DayOfTheMonth)

	    d.createDayEvent(calendarSrv, day, ImgURL) // send day struct with imgURL of commemorated saint
	
	    d.addDayEvents(calendarSrv, day)
	}
	
	log.Println("Finished Calendar Shenanigans!")
}



