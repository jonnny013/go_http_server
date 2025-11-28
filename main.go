package main

import (
	"fmt"
	"log"

	"github.com/jonnny013/go_html_server/internal/headers"
)

func main() {
	intro := []byte("Hello: my name is Jon.\r\noops: hi\r\n\r\n")

	header := headers.NewHeaders()
	curByte := 0
	for {
		bytes, done, err := header.Parse(intro[curByte:])
		if err != nil {
			log.Fatalf("error %s\n", err.Error())
		}

		if done {
			break
		}
		curByte += bytes
	}

	for key, value := range header {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}
}
