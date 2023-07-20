package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

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

	const timeout = 30 * time.Second
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept connection: %s\n", err)
		}

		err = conn.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			log.Fatalf("Failed to set read deadline: %s\n", err)
		}

		reqBuffer := make([]byte, 1500)

		go func() {
			defer conn.Close()
			log.Printf("[Remote Address] %s\n\n", conn.RemoteAddr())

			for {
				reqMessage := ""
				for {
					n, err := conn.Read(reqBuffer)
					if err != nil {
						if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
							log.Printf("read timeout: %s\n", err)
						} else {
							log.Printf("Failed to read connection: %s\n", err)
						}
						return
					}
					reqMessage += string(reqBuffer[:n])
					if strings.HasSuffix(reqMessage, "\r\n\r\n") {
						break
					}
				}
				log.Printf("[Request Message]\n%s", reqMessage)

				statusCode := 200
				req, isValid := http.CheckRequest(reqMessage)
				if !isValid {
					statusCode = 400
					resHeaderString := http.GenerateResponseHeader(statusCode, "", req.Connection)
					log.Printf("[Response Status Code] %d\n\n", statusCode)
					//log.Printf("[Response Message]\n%s\n\n", resMessage)
					conn.Write([]byte(resHeaderString))
					break
				}

				content, contentType, isFound, err := http.ReadFile(req.Path)
				if err != nil {
					log.Fatalf("Failed to load file: %s\n", err)
				} else if !isFound {
					statusCode = 404
				}

				resHeaderString := http.GenerateResponseHeader(statusCode, contentType, req.Connection)
				log.Printf("[Response Status Code] %d\n\n", statusCode)
				//log.Printf("[Response Message]\n%s\n\n", resMessage)
				conn.Write([]byte(resHeaderString))

				chunkBuffer := bytes.NewBufferString(content)
				for chunkBuffer.Len() != 0 {
					chunkBytes := chunkBuffer.Next(20)
					chunk := fmt.Sprintf("%x\r\n%s\r\n", len(chunkBytes), string(chunkBytes))
					conn.Write([]byte(chunk))
				}
				body := "0\r\n\r\n"
				conn.Write([]byte(body))

				if req.Connection == "Close" {
					break
				}
			}
		}()
	}
}
