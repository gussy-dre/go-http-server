package http

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type ResponseHeader struct {
	StatusCode        int
	Message           string
	HTTPVersion       string
	Connection        string
	LastModified      string
	CacheControl      string
	TransferEncording string
	ContentType       string
}

func GenerateResponseHeader(statusCode int, contentType string, connection string, modString string) string {
	resHeader := &ResponseHeader{
		StatusCode:        statusCode,
		Message:           generateResponseMessage(statusCode),
		HTTPVersion:       "HTTP/1.1",
		Connection:        "Keep-Alive",
		LastModified:      modString,
		CacheControl:      "max-age=3600",
		TransferEncording: "chunked",
		ContentType:       contentType,
	}
	resHeaderString := fmt.Sprintf("%s %d %s\r\n", resHeader.HTTPVersion, resHeader.StatusCode, resHeader.Message)

	if resHeader.StatusCode == 400 {
		resHeaderString += "\r\n"
		return resHeaderString
	}

	if connection == "Close" {
		resHeader.Connection = connection
	}
	resHeaderString += fmt.Sprintf("Connection: %s\r\n", resHeader.Connection)

	if resHeader.StatusCode != 404 {
		resHeaderString += fmt.Sprintf("Last-Modified: %s\r\n", resHeader.LastModified)
		resHeaderString += fmt.Sprintf("Cache-Control: %s\r\n", resHeader.CacheControl)
	}

	if resHeader.StatusCode != 304 {
		resHeaderString += fmt.Sprintf("Content-Type: %s\r\n", resHeader.ContentType)
		resHeaderString += "Transfer-Encoding: chunked\r\n"
	}

	resHeaderString += "\r\n"
	return resHeaderString
}

func ReadFile(path string) (string, string, string, bool, error) {
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

	if path == "/404.html" {
		return string(b), contentType, "", isFound, err
	}
	finfo, err := f.Stat()
	if err != nil {
		log.Fatalf("Failed to get %s info: %s\n", path, err)
	}
	fmod := finfo.ModTime()

	loc, err := time.LoadLocation("GMT")
	if err != nil {
		log.Fatalf("Failed to load time location: %s\n", err)
	}
	modString := fmod.In(loc).Format(time.RFC1123)

	return string(b), contentType, modString, isFound, err
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
