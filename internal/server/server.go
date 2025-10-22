package server

import (
	"fmt"
	"log"
	"net"
	"to-tcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	server := &Server{
		listener: listener,
		closed:   false,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true

	err := s.listener.Close()

	return err
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	// defer conn.Close()

	headers := response.GetDefaultHeaders(0)
	response.WriteStatusLine(conn, response.StatusOK)
	response.WriteHeaders(conn, headers)

}
