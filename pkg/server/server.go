package server

import (
	"fmt"
	"io"
	"net"
)

type Connection interface {
	io.ReadWriteCloser

	RemoteAddr() net.Addr
}

type Handler func(Connection)

type Server struct {
	Handler Handler
	Addr    string
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.Addr)
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
	return &Server{Handler: handler, Addr: ":3000"}
}
