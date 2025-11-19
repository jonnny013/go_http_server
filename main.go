package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
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
				break
			}
			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				curLine += string(data[:i])
				ch <- curLine
				data = data[i+1:]
				curLine = ""
			}
			curLine += string(data)
		}
		if len(curLine)!= 0 {
			ch <- curLine
		}
	}()

	return ch
}

func main() {
	f, err := os.Open("messages.txt")

	if err != nil {
		log.Fatal(err)
	}

	ch := getLinesChannel(f)

	for l := range ch {
		fmt.Printf("read: %s\n", l)
	}

}
