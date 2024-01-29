package basicfunctions

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)


func CreateProvisionalDir(dirPath string) {
	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		log.Printf("Directory already exists, proceeding to read files inside it...")
		return
	}
	fmt.Printf("Directory created: %s\n", dirPath)
}

func CreateProvisionalFiles(dirPath string, Id string) {	
	content, err := DownloadFile(fmt.Sprintf("https://docs.google.com/uc?export=download&id=%s", Id))
	if err != nil {
		log.Println("Error at CreateProvisionalFiles:", err)
		return
	}
	filePath := fmt.Sprintf("%s/downloaded_file.pdf", dirPath)
	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Println("Error saving file:", err)
		return
	}
}

func WriteJson(data interface{}) {
	fileName := "calendar.json"
	jsonCal, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal("Error at WriteJson: unable to marshal data:", err)
	}
	err = os.WriteFile(fileName, jsonCal, 0644)
	if err != nil {
		log.Fatal("Error at WriteJson: unable to write json file:",err)
	}

	log.Println(fileName, "file written successfully!")

	DeleteFile("parsedDocument")
}