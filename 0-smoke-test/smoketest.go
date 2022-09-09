package main

import (
	"io"
	"log"
	"net"
)

func handleConn(conn net.Conn) {
	log.Printf("New connection: %s", conn.RemoteAddr())

	_, err := io.Copy(conn, conn)
	if err != nil {
		log.Printf("Error copying data: %s", err)
	}

	conn.Close()
}

func main() {
	listener, _ := net.Listen("tcp", ":3000")

	for {
		conn, _ := listener.Accept()
		go handleConn(conn)
	}
}
