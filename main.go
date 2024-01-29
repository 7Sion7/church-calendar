package main

import (
	"fmt"
	funcs "houses-data/basicFunctions"
	googleapi "houses-data/calendarInterface"
	pattern "houses-data/patternRecogniser"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {

	now := time.Now().Local()

	var Id string

	dirPath := "unparsedDocuments"

	defer funcs.DeleteFile(dirPath)

	funcs.CreateProvisionalDir(dirPath)

    const baseURL = "http://www.sourozh.org/cathedral-timetable-old/"

	nextMonth := now.AddDate(0, 1, 0).Month().String()

	c:= colly.NewCollector()

	c.OnHTML("#content a", func(h *colly.HTMLElement) {
		if  strings.Contains(h.Text, strings.ToUpper(nextMonth)) {
			Id = funcs.ExtractFileID(h.Attr("href"))
		}
	})

	c.Visit(baseURL)

	funcs.CreateProvisionalFiles(dirPath, Id) //Write provisional pdf files

	path := funcs.GetFilePathInFolder(dirPath)
	
	parsedFile := funcs.ParseFile(strings.Replace(path, "\\", "/", 1))

	Month := pattern.FormatData(parsedFile)

    googleapi.CallAPI(Month)

	elapsedTime := time.Since(now)
	fmt.Printf("Time taken: %s\n", elapsedTime)
}

