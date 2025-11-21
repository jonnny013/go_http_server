package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")

	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		line, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				log.Println("stdin closed (EOF), exiting")
				return
			}
			conn.Write([]byte(fmt.Sprintf("err: %s", err.Error())))
			continue
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Printf("write error: %v", err)
		} else {
			log.Printf("sent %d bytes", len(line))
		}
	}
}
