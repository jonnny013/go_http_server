package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")

	if err != nil {
		log.Fatal(err)
	}

	curLine := ""
	for {
		data := make([]byte, 8)
		n, err := f.Read(data)

		if err != nil {
			break
		}
		data = data[:n]
		if i := bytes.IndexByte(data, '\n'); i != -1 {
			curLine += string(data[:i])
			data = data[i+1:]
			fmt.Printf("read: %s\n", curLine)
			curLine = ""
		}
		curLine += string(data)
	}

	if len(curLine) != 0 {
		fmt.Printf("read: %s\n", curLine)
	}
}
