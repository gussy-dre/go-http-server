package http

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type ResponseHeader struct {
	StatusCode        int
	Message           string
	HTTPVersion       string
	Connection        string
	TransferEncording string
	ContentType       string
}

func GenerateResponseHeader(statusCode int, contentType string, connection string) string {
	resHeader := &ResponseHeader{
		StatusCode:        statusCode,
		Message:           generateResponseMessage(statusCode),
		HTTPVersion:       "HTTP/1.1",
		Connection:        "Keep-Alive",
		TransferEncording: "chunked",
		ContentType:       contentType,
	}
	resHeaderString := fmt.Sprintf("%s %d %s\r\n", resHeader.HTTPVersion, resHeader.StatusCode, resHeader.Message)

	if resHeader.StatusCode == 400 {
		resHeaderString += "\r\n"
		return resHeaderString
	}

	if len(resHeader.ContentType) > 0 {
		resHeaderString += fmt.Sprintf("Content-Type: %s\r\n", contentType)
	}

	if connection == "Close" {
		resHeader.Connection = connection
	}
	resHeaderString += fmt.Sprintf("Connection: %s\r\n", resHeader.Connection)

	resHeaderString += "Transfer-Encoding: chunked\r\n\r\n"
	return resHeaderString
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
