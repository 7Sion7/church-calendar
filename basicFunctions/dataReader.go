package basicfunctions

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/dslipak/pdf"
)

func ManageCLI() {
	args := os.Args[1:]
	if len(args) == 1 {
		log.Println("Running Query from Command Line!")
	}
	log.Printf("No/Too many Queries on the Command Line, app running normally.Commands --> %s\n", args)
	return
}

func ExtractFileID(url string) string {
	re := regexp.MustCompile(`/d/([^/]+)/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func DownloadFile(url string) ([]byte, error) {
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

func ParseFile(pdfPath string) string {
	f, r, err := pdf.Open(pdfPath)
    
	defer f.Close()
	if err != nil {
		log.Println("Error at ParseFile: unable to Open pdf file->",err.Error())
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		log.Println("Error at ParseFile: unable to GetPlainText->",err.Error())
	}

	_, err = buf.ReadFrom(b)
	if err != nil {
		log.Println("Error at ParseFile: unable to ReadFrom io.Reader->",err.Error())
	}

	return buf.String()
}



