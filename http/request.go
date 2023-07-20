package http

import (
	"bufio"
	"log"
	"strings"
)

type Request struct {
	Method          string
	Path            string
	HTTPVersion     string
	Host            string
	Connection      string
	IfModifiedSince string
}

func CheckRequest(reqMessage string) (*Request, bool) {
	req := &Request{
		Method:          "",
		Path:            "",
		HTTPVersion:     "",
		Host:            "",
		Connection:      "",
		IfModifiedSince: "",
	}
	scanner := bufio.NewScanner(strings.NewReader(reqMessage))

	scanner.Scan()
	topLine := scanner.Text()
	topLineFields := strings.Split(topLine, " ")
	if len(topLineFields) != 3 {
		return req, false
	}

	if topLineFields[0] != "GET" || !strings.Contains(topLineFields[2], "HTTP/1.1") {
		return req, false
	}
	req.Method = topLineFields[0]
	req.Path = topLineFields[1]
	req.HTTPVersion = topLineFields[2]

	for scanner.Scan() {
		reqLine := scanner.Text()
		if strings.HasPrefix(reqLine, "Host: ") {
			req.Host = strings.Replace(reqLine, "Host: ", "", 1)
		} else if strings.HasPrefix(reqLine, "Connection: ") {
			req.Connection = strings.Replace(reqLine, "Connection: ", "", 1)
		} else if strings.HasPrefix(reqLine, "If-Modified-Since: ") {
			req.IfModifiedSince = strings.Replace(reqLine, "If-Modified-Since: ", "", 1)
		} else if reqLine == "" {
			break
		}
	}

	if len(req.Host) == 0 {
		return req, false
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to scan: %s\n", err)
	}

	return req, true
}
