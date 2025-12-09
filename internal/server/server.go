package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	closed bool
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!")
	conn.Write(out)
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
