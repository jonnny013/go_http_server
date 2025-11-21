package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		curLine := ""

		for {
			data := make([]byte, 8)
			n, err := f.Read(data)

			if err != nil {
				if curLine != "" {
					ch <- curLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}

			str := string(data[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				ch <- fmt.Sprintf("%s%s", curLine, parts[i])
				curLine = ""
			}
			curLine += parts[len(parts)-1]
		}

	}()

	return ch
}

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

		for l := range getLinesChannel(c) {
			fmt.Println(l)
		}
		fmt.Println("Connection to ", c.RemoteAddr(), "closed")
	}

}
