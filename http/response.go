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

	response := fmt.Sprintf("%s %d %s\r\n", res.HTTPVersion, res.StatusCode, res.Message)

	if len(res.ContentType) > 0 {
		response += fmt.Sprintf("Content-Type: %s\r\n", contentType)
	}

	if connection == "Close" {
		res.Connection = connection
	}
	response += fmt.Sprintf("Connection: %s\r\n", res.Connection)

	if res.ContentLength != 0 {
		response += fmt.Sprintf("Content-Length: %d\r\n", res.ContentLength)
	}

	if len(res.Body) > 0 {
		response += fmt.Sprintf("\r\n%s", body)
	}

	response += "\r\n"
	return response
}

func ReadFile(path string) (string, string, error) {
	if path == "/" {
		path = "/index.html"
	}
	f, err := os.Open(fmt.Sprintf("public%s", path))
	if err != nil {
		log.Print(err, path)
		return "", "", err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	extension := filepath.Ext(path)
	contentType := "text/plain"
	if extension == ".html" {
		contentType = "text/html"
	} else if extension == ".png" {
		contentType = "image/png"
	} else if extension == ".jpeg" {
		contentType = "image/jpeg"
	}

	return string(b), contentType, err
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
