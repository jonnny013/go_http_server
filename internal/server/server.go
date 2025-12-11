package server

import (
	"fmt"
	"io"
	"net"

	"github.com/jonnny013/go_html_server/internal/request"
	"github.com/jonnny013/go_html_server/internal/response"
)

type Server struct {
	closed  bool
	handler Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()

	resWriter := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)

	if err != nil {
		s.closed = true
		return
	}

	s.handler(resWriter, req)

}

func runServer(s *Server, listener net.Listener) {

	go func() {
		for {
			conn, err := listener.Accept()
			if s.closed {
				return
			}
			if err != nil {
				return
			}
			go runConnection(s, conn)
		}
	}()

}

func Serve(port int, handler Handler) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		handler: handler,
		closed:  false,
	}

	go runServer(server, lis)

	return server, err

}

func (s *Server) Close() {
	s.closed = true

}

func Response200() []byte {
	return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}

func Response400() []byte {
	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}

func Response500() []byte {
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}
