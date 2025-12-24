package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jirbthagoras/tcp-to-http/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", "error", err)
		}

		fmt.Printf("Request Line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- HTTP Version: %s\n", r.RequestLine.HttpVersion)
		fmt.Printf("- Request Target: %s\n", r.RequestLine.RequestTarget)
	}
}
