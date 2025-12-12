package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jonnny013/go_html_server/internal/request"
	"github.com/jonnny013/go_html_server/internal/response"
	"github.com/jonnny013/go_html_server/internal/server"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {
	path := req.RequestLine.RequestTarget
	var body []byte

	headers := response.GetDefaultHeaders()

	if path == "/yourproblem" {
		w.WriteStatusLine(response.StatusBadRequest)
		body = server.Response400()
	} else if path == "/myproblem" {
		w.WriteStatusLine(response.StatusSystemError)
		body = server.Response500()
	} else if strings.HasPrefix(path, "/httpbin") {

		res, err := http.Get("https://httpbin.org/" + strings.TrimPrefix(path, "/httpbin/"))

		if err != nil {
			w.WriteStatusLine(response.StatusSystemError)
			body = server.Response500()
		} else {

			w.WriteStatusLine(response.StatusOk)
			headers.Set("Transfer-Encoding", "chunked")
			w.WriteHeaders(headers)
			for {
				data := make([]byte, 32)
				n, err := res.Body.Read(data)
				if err != nil {
					break
				}
				w.WriteBody(fmt.Appendf(nil, "%x\r\n", n))
				w.WriteBody(data[:n])
				w.WriteBody([]byte("\r\n"))
			}
			w.WriteChunkedBodyDone()
			return
		}

		w.WriteStatusLine(response.StatusStream)
	} else {
		w.WriteStatusLine(response.StatusOk)
		body = server.Response200()
	}

	response.GetContentLengthHeader(headers, len(body))

	w.WriteHeaders(headers)

	w.WriteBody(body)

}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
