package server

import (
	"fmt"
	"io"
	"net"
)

type Connection interface {
	io.ReadWriteCloser
}

type Handler func(Connection)

type Server struct {
	Handler Handler
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		return fmt.Errorf("error listening for connections: %w", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}

		go s.Handler(conn)
	}
}

func New(handler Handler) *Server {
	return &Server{handler}
}
