package server

import (
	"fmt"
	"log"
	"net"
	"to-tcp/internal/request"
	"to-tcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   bool
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}
type Handler func(w *response.Writer, req *request.Request)

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	server := &Server{
		listener: listener,
		closed:   false,
		handler:  handler,
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
	defer conn.Close()

	rw := response.NewWriter(conn)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		rw.WriteStatusLine(response.StatusBadRequest)
		rw.WriteHeaders(response.GetDefaultHeaders(0))
		return
	}

	s.handler(rw, r)
}
