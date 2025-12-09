package server

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

type serverState string

const (
	StateDone  serverState = "done"
	StateError serverState = "error"
	StateOn    serverState = "listening"
)

type Server struct {
	listener net.Listener
	state    serverState
	closed   atomic.Bool
	wg       sync.WaitGroup
}

func Serve(port int) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: lis,
		state:    StateOn,
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	err := s.listener.Close()
	s.wg.Wait()
	return err
}

func (s *Server) listen() {

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			s.state = StateError
			return
		}
		s.wg.Add(1)

		go func(c net.Conn) {
			defer s.wg.Done()
			s.handle(c)
		}(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	statusLine := "HTTP/1.1 200 OK"
	headers := "Content-Type: text/plain\r\n"
	contentLen := "Content-Length: 13\r\n\r\n"
	body := "Hello World!"
	conn.Write([]byte(statusLine + headers + contentLen + body))
}
