package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"server/http"
)

var (
	ip   = flag.String("ip", "127.0.0.1", "IP Address")
	port = flag.String("port", "8888", "Port")
)

func main() {
	flag.Parse()

	address := fmt.Sprintf("%s:%s", *ip, *port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Faild to listen %s\n", address)
	}
	defer listener.Close()
	log.Printf("Now listening %s\n\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept connection: %s\n", err)
		}

		buf := make([]byte, 1500)

		go func() {
			defer conn.Close()
			log.Printf("[Remote Address] %s\n\n", conn.RemoteAddr())

			reqMessage := ""
			for {
				n, _ := conn.Read(buf)
				reqMessage += string(buf[:n])
				if strings.HasSuffix(reqMessage, "\r\n\r\n") {
					break
				}
			}
			log.Printf("[Request Message]\n%s", reqMessage)

			statusCode := 200
			req, isValid := http.CheckRequest(reqMessage)
			if !isValid {
				statusCode = 400
				resMessage := http.GenerateResponse(statusCode, "", "")
				log.Printf("[Response Status Code] %d\n\n", statusCode)
				//log.Printf("[Response Message]\n%s\n\n", resMessage)
				conn.Write([]byte(resMessage))
				return
			}

			content, contentType, isFound, err := http.ReadFile(req.Path)
			if err != nil {
				log.Fatalf("Failed to load file: %s\n", err)
			} else if isFound {
				statusCode = 404
			}
			resMessage := http.GenerateResponse(statusCode, contentType, content)
			log.Printf("[Response Status Code] %d\n\n", statusCode)
			//log.Printf("[Response Message]\n%s\n\n", resMessage)
			conn.Write([]byte(resMessage))
		}()
	}
}
