package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonnny013/go_html_server/internal/request"
	"github.com/jonnny013/go_html_server/internal/server"
)

const port = 42069

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	var errorHandler = &server.HandlerError{}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		errorHandler.StatusCode = 400
		errorHandler.Message = "Your problem is not my problem\n"
	case "/myproblem":
		errorHandler.StatusCode = 500
		errorHandler.Message = "Woopsie, my bad\n"
	default:
		errorHandler.StatusCode = 200
		errorHandler.Message = "All good\n"
	}

	w.Write([]byte(errorHandler.Message))
	return errorHandler

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
