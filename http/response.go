package http

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Response struct {
	StatusCode        int
	Message           string
	HTTPVersion       string
	Connection        string
	TransferEncording string
	ContentType       string
	Body              string
}

func GenerateResponseHeader(statusCode int, contentType string, body string, connection string) (string, *Response) {
	res := &Response{
		StatusCode:        statusCode,
		Message:           generateResponseMessage(statusCode),
		HTTPVersion:       "HTTP/1.1",
		Connection:        "Keep-Alive",
		TransferEncording: "chunked",
		ContentType:       contentType,
		Body:              body,
	}
	resHeader := fmt.Sprintf("%s %d %s\r\n", res.HTTPVersion, res.StatusCode, res.Message)

	if res.StatusCode == 400 {
		resHeader += "\r\n"
		return resHeader, res
	}

	if len(res.ContentType) > 0 {
		resHeader += fmt.Sprintf("Content-Type: %s\r\n", contentType)
	}

	if connection == "Close" {
		res.Connection = connection
	}
	resHeader += fmt.Sprintf("Connection: %s\r\n", res.Connection)

	resHeader += "Transfer-Encoding: chunked\r\n\r\n"
	return resHeader, res
}

func ReadFile(path string) (string, string, bool, error) {
	isFound := true

	if path == "/" {
		path = "/index.html"
	}
	f, err := os.Open(fmt.Sprintf("public%s", path))
	if err != nil {
		path = "/404.html"
		isFound = false
		f, err = os.Open(fmt.Sprintf("public%s", path))
		if err != nil {
			log.Printf("Failed to open %s: %s\n", path, err)
			return "", "", false, err
		}
		log.Printf("Failed to open %s: %s\n", path, err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		log.Printf("Failed to read %s: %s\n", path, err)
		return "", "", false, err
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
