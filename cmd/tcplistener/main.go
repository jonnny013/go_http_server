package main

import (
	"fmt"
	"log"
	"net"

	request "github.com/jonnny013/go_html_server/internal"
)

const port = ":42069"

func main() {
	l, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("error %s\n", err.Error())
	}

	defer l.Close()

	fmt.Println("Listening for TCP traffic on", port)

	for {
		c, err := l.Accept()

		if err != nil {
			log.Fatalf("error %s\n", err.Error())
		}

		fmt.Println("Connection accepted from", c.RemoteAddr())
		req, err := request.RequestFromReader(c)
		if err != nil {
			log.Fatalf("error %s\n", err.Error())
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		fmt.Println("Connection to ", c.RemoteAddr(), "closed")
	}

}
