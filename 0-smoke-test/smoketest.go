package main

import (
	"io"
	"log"

	"github.com/nint8835/protohackers/pkg/server"
)

func handleConn(conn server.Connection) {
	_, err := io.Copy(conn, conn)
	if err != nil {
		log.Printf("Error copying data: %s", err)
	}

	conn.Close()
}

func main() {
	err := server.New(handleConn).Start()
	if err != nil {
		log.Fatalf("Error running server: %s", err)
	}
}
