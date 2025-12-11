package server

import (
	"fmt"
	"io"
	"net"

	"github.com/jonnny013/go_html_server/internal/response"
)

type Server struct {
	closed bool
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	err := response.WriteStatusLine(conn, 200)
	if err != nil {
		s.closed = true
		return
	}
	h := response.GetDefaultHeaders(0)

	err = response.WriteHeaders(conn, h)
	if err != nil {
		s.closed = true
		return
	}
	conn.Close()
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

func Serve(port int) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{}
	go runServer(server, lis)

	return server, err

}

func (s *Server) Close() {
	s.closed = true

}
