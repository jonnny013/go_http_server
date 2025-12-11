package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonnny013/go_html_server/internal/request"
	"github.com/jonnny013/go_html_server/internal/response"
	"github.com/jonnny013/go_html_server/internal/server"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {
	var body []byte

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.StatusBadRequest)
		body = server.Response400()
	case "/myproblem":
		w.WriteStatusLine(response.StatusSystemError)
		body = server.Response500()
	default:
		w.WriteStatusLine(response.StatusOk)
		body = server.Response200()
	}

	headers := response.GetDefaultHeaders(len(body))
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
