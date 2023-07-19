package http

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Request struct {
	Method      string
	Path        string
	HTTPVersion string
	Host        string
	Connection  string
}

func CheckRequest(reqMessage string) (*Request, bool) {
	req := &Request{
		Method:      "",
		Path:        "",
		HTTPVersion: "",
		Host:        "",
		Connection:  "",
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
		} else if reqLine == "" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		return req, false
	}

	return req, true
}
