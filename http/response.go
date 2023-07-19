package http

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Response struct {
	StatusCode    int
	Message       string
	HTTPVersion   string
	Connection    string
	ContentLength int
	ContentType   string
	Body          string
}

func GenerateResponse(statusCode int, contentType string, body string, connection string) string {
	res := &Response{
		StatusCode:    statusCode,
		Message:       generateResponseMessage(statusCode),
		HTTPVersion:   "HTTP/1.1",
		Connection:    "Keep-Alive",
		ContentLength: len(body),
		ContentType:   contentType,
		Body:          body,
	}
	resString := fmt.Sprintf("%s %d %s\r\n", res.HTTPVersion, res.StatusCode, res.Message)

	if res.StatusCode == 400 {
		resString += "\r\n"
		return resString
	}

	if len(res.ContentType) > 0 {
		resString += fmt.Sprintf("Content-Type: %s\r\n", res.ContentType)
	}

	if connection == "Close" {
		res.Connection = connection
	}
	resString += fmt.Sprintf("Connection: %s\r\n", res.Connection)

	if res.ContentLength != 0 {
		resString += fmt.Sprintf("Content-Length: %d\r\n", res.ContentLength)
	}

	if len(res.Body) > 0 {
		resString += fmt.Sprintf("\r\n%s", body)
	}

	resString += "\r\n"
	return resString
}

func ReadFile(path string) (string, string, bool, error) {
	isFound := true

	if path == "/" {
		path = "/index.html"
	}
	f, err := os.Open(fmt.Sprintf("public%s", path))
	if err != nil {
		log.Printf("Failed to open %s: %s\n", path, err)
		path = "/404.html"
		isFound = false
		f, err = os.Open(fmt.Sprintf("public%s", path))
		if err != nil {
			log.Fatalf("Failed to open %s: %s\n", path, err)
		}
		log.Printf("Failed to open %s: %s\n", path, err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("Failed to read %s: %s\n", path, err)
	}

	extension := filepath.Ext(path)
	contentType := "text/plain"
	if extension == ".html" {
		contentType = "text/html"
	} else if extension == ".png" {
		contentType = "image/png"
	} else if extension == ".jpeg" {
		contentType = "image/jpeg"
	}

	return string(b), contentType, isFound, err
}

func generateResponseMessage(statusCode int) string {
	switch statusCode {
	case 200:
		return "OK"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	default:
		return "Internal Server Error"
	}
}
