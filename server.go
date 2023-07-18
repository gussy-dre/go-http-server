package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
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

			// 個々を実装する
			resMessage := reqMessage

			log.Printf("[Response Status Code] %s\n\n", "200")
			//log.Printf("[Response Message]\n%s\n\n", resMessage)
			conn.Write([]byte(resMessage))
		}()
	}
}
