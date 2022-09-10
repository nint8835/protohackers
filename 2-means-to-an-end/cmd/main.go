package main

import (
	"net"

	meanstoanend "github.com/nint8835/protohackers/2-means-to-an-end"
)

func main() {
	listener, _ := net.Listen("tcp", ":3000")

	for {
		conn, _ := listener.Accept()
		go meanstoanend.HandleConn(conn)
	}
}
