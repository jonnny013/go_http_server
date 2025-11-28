package main

import (
	"fmt"
	"log"

	"github.com/jonnny013/go_html_server/internal/headers"
)

func main() {
	intro := []byte("Hello: my name is Jon.\r\noops: hi\r\n\r\n")

	header := headers.NewHeaders()

	_, _, err := header.Parse(intro)
	if err != nil {
		log.Fatalf("error %s\n", err.Error())
	}

	for key, value := range header.GetAll() {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}
}
