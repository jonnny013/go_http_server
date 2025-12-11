package server

import (
	"bytes"
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

type Handler func(w io.Writer, req *request.Request) *HandlerError

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)

	if err != nil {
		s.closed = true
		return
	}

	buf := new(bytes.Buffer)

	handlerError := s.handler(buf, req)

	err = response.WriteStatusLine(conn, response.StatusCode(handlerError.StatusCode))

	if err != nil {
		s.closed = true
		return
	}

	h := response.GetDefaultHeaders(len(handlerError.Message))

	err = response.WriteHeaders(conn, h)

	if err != nil {
		s.closed = true
		return
	}

	conn.Write(buf.Bytes())

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

func (e *HandlerError) NewHandlerError(w io.Writer) error {
	return response.WriteStatusLine(w, response.StatusCode(e.StatusCode))
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
