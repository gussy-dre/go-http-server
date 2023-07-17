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

func generateResponse(req string) (int, string) {
	reqTopLine := strings.Split(req, "\n")[0]
	httpArray := strings.Split(reqTopLine, " ")
	if len(httpArray) != 3 {
		return 400, "Bad Request"
	}

	if httpArray[0] != "GET" || !strings.Contains(httpArray[2], "HTTP/1.1") {
		return 400, "Bad Request"
	}

	return 200, "OK"
}

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
			fmt.Printf("[Remote Address]\n%s\n", conn.RemoteAddr())

			n, _ := conn.Read(buf)
			req := string(buf[:n])
			fmt.Printf("[Message]\n%s", req)
			statusCode, mes := generateResponse(req)
			res := fmt.Sprintf("HTTP/1.1 %d %s\nContent-Type: text/html\nHello!", statusCode, mes)
			fmt.Println(res)

			conn.Write([]byte(res))

			conn.Close()
		}()
	}
}
