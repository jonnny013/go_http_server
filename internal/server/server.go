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
	closed bool
}

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func runConnection(s *Server, conn io.ReadWriteCloser, handler Handler) {

	req, err := request.RequestFromReader(conn)

	if err != nil {
		s.closed = true
		return
	}

	buf := new(bytes.Buffer)

	handlerError := handler(buf, req)

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

	conn.Close()
}

func runServer(s *Server, listener net.Listener, handler Handler) {

	go func() {
		for {
			conn, err := listener.Accept()
			if s.closed {
				return
			}
			if err != nil {
				return
			}
			go runConnection(s, conn, handler)
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
	server := &Server{}
	go runServer(server, lis, handler)

	return server, err

}

func (s *Server) Close() {
	s.closed = true

}
