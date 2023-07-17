package main

import (
	"flag"
	"fmt"
	"log"
	"net"

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
	fmt.Printf("Now listening %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Failed to accept connection")
		}

		buf := make([]byte, 1500)

		go func() {
			defer conn.Close()
			fmt.Printf("[Remote Address]\n%s\n", conn.RemoteAddr())

			n, _ := conn.Read(buf)
			reqMessage := string(buf[:n])
			fmt.Printf("[Message]\n%s", reqMessage)

			req, isValid := http.CheckRequest(reqMessage)

			res, contentType, err := http.ReadFile(req.Path)
			if !isValid {
				conn.Write([]byte(http.GenerateResponse(400, "", "")))
			} else if err != nil {
				notFoundRes, contentType, err := http.ReadFile("404.html")
				if err != nil {
					log.Fatal("Failed to load 404.html")
				}
				conn.Write([]byte(http.GenerateResponse(404, contentType, notFoundRes)))
			} else {
				conn.Write([]byte(http.GenerateResponse(200, contentType, res)))
			}
		}()
	}
}
