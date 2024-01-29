package basicfunctions

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	envFilePath := filepath.Join(currentDir, ".env") 

	err = godotenv.Load(envFilePath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}


func DeleteFile(dirPath string) {
	err := os.RemoveAll(dirPath)
	if err != nil {
		log.Printf("Error deleting file %s: %s\n", dirPath, err)
		return
	}
	log.Printf("Directory deleted: %s\n", dirPath)
}

func GetFilePathInFolder(folderPath string) (string) {
	var filePath string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			filePath = path
		}
		return nil
	})
	if err != nil {
		log.Printf("Error getting file paths: %s\n", err)
		return ""
	}

	return filePath
}