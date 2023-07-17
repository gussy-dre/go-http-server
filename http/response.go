package http

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func GenerateResponse(statusCode int, contentType string, body string) string {
	message := ""
	if statusCode == 200 {
		message = "OK"
	} else if statusCode == 400 {
		message = "Bad Request"
	} else if statusCode == 404 {
		message = "Not Found"
	}

	response := fmt.Sprintf("HTTP/1.1 %d %s\n", statusCode, message)
	if len(body) > 0 {
		response += fmt.Sprintf("Content-Type: %s\n\n%s", contentType, body)
	}

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
