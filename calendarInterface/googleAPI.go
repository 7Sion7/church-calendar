package googleapi

import (
	pattern "church-calendar/patternRecogniser"
	"context"
	"log"

	// "log"

	"time"
	//"golang.org/x/oauth2/google"
	//"google.golang.org/api/option"
)




var zone, _ = time.Now().Local().Zone()

const tokFile = "json/token.json"

func CallAPI(month pattern.Month) {

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




