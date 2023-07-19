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
	fmt.Printf("Now listening %s\n\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Failed to accept connection")
		}

		buf := make([]byte, 1500)

		go func() {
			defer conn.Close()
			fmt.Printf("[Remote Address]\n%s\n\n", conn.RemoteAddr())

			for {
				reqMessage := ""
				for {
					n, _ := conn.Read(buf)
					reqMessage += string(buf[:n])
					if strings.HasSuffix(reqMessage, "\r\n\r\n") {
						break
					}
				}
				fmt.Printf("[Request Message]\n%s", reqMessage)

				req, isValid := http.CheckRequest(reqMessage)

				res, contentType, err := http.ReadFile(req.Path)
				if !isValid {
					resMessage := http.GenerateResponse(400, "", "", req.Connection)
					fmt.Printf("[Response Status Code] %d\n\n", 400)
					//fmt.Printf("[Response Message]\n%s\n\n", resMessage)
					conn.Write([]byte(resMessage))
				} else if err != nil {
					notFoundRes, contentType, err := http.ReadFile("/404.html")
					if err != nil {
						log.Fatal("Failed to load 404.html")
					}
					resMessage := http.GenerateResponse(404, contentType, notFoundRes, req.Connection)
					fmt.Printf("[Response Status Code] %d\n\n", 404)
					//fmt.Printf("[Response Message]\n%s\n\n", resMessage)
					conn.Write([]byte(resMessage))
				} else {
					resMessage := http.GenerateResponse(200, contentType, res, req.Connection)
					fmt.Printf("[Response Status Code] %d\n\n", 200)
					//fmt.Printf("[Response Message]\n%s\n\n", resMessage)
					conn.Write([]byte(resMessage))
				}

				if req.Connection == "Close" {
					break
				}
			}
		}()
	}
}
