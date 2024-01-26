package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	pattern "houses-data/patternRecogniser"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	//"os"
	"regexp"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"

	"github.com/gocolly/colly/v2"
)


var extractedData []*pattern.Month

var dateRegex = regexp.MustCompile(`((?:Mon|Tue|Wed|Thu|Fri|Sat|Sun))`)

var hourRegex = regexp.MustCompile(`(\d{2}:\d{2})`)

func main() {
	var IDs []string

	dirPath := "unparsedDocuments"

	defer deleteFile(dirPath)

	CreateProvisionalDir(dirPath)

    const baseURL = "http://www.sourozh.org/cathedral-timetable-old/"

	thisMonth := time.Now().Local().Month().String()
	nextMonth := time.Now().Local().AddDate(0, 1, 0).Month().String()

	c:= colly.NewCollector()

	c.OnHTML("#content a", func(h *colly.HTMLElement) {
		if strings.Contains(h.Text, strings.ToUpper(thisMonth)) || strings.Contains(h.Text, strings.ToUpper(nextMonth)) {
			IDs = append(IDs, extractFileID(h.Attr("href")))
		}
	})

	c.Visit(baseURL)

	for number, Id := range IDs {
		content, err := downloadFile(fmt.Sprintf("https://docs.google.com/uc?export=download&id=%s", Id))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		filePath := fmt.Sprintf("%s/downloaded_file%s.pdf", dirPath, strconv.Itoa(number+1))
		err = os.WriteFile(filePath, content, 0644)
		if err != nil {
			fmt.Println("Error saving file:", err)
			return
		}

	}

	filePaths, err := getFilePathsInFolder(dirPath)
	if err != nil {
		fmt.Printf("Error getting file paths: %s\n", err)
		return
	}

	for _, path := range filePaths {

		parsedFile, err := parseFile(strings.Replace(path, "\\", "/", 1))
		if err != nil {
			fmt.Println(err.Error())
		}

		sortedData := pattern.FormatData(parsedFile)

		jsonCal, err := json.MarshalIndent(sortedData, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile("calendar.json", jsonCal, 0644)
		if err != nil {
			log.Fatal(err)
		}

		deleteFile("parsedDocument")
	}
}


func CreateProvisionalDir(dirPath string) {
	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %s\n", err)
		return
	}

	fmt.Printf("Directory created: %s\n", dirPath)
}

func deleteFile(dirPath string) {
	err := os.RemoveAll(dirPath)
	if err != nil {
		fmt.Printf("Error deleting file %s: %s\n", dirPath, err)
		return
	}
	fmt.Printf("Directory deleted: %s\n", dirPath)
}

func getFilePathsInFolder(folderPath string) ([]string, error) {
	var filePaths []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {

		if err == nil && !info.IsDir() {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	return filePaths, err
}

func parseFile(pdfPath string) (text string, err error) {
	f, r, err := pdf.Open(pdfPath)
	defer f.Close()
	if err != nil {
		return
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return
	}

	buf.ReadFrom(b)
	text = buf.String()
	return text, nil
}

func extractFileID(url string) string {
	re := regexp.MustCompile(`/d/([^/]+)/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func downloadFile(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
