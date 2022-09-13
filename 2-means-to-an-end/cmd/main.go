package main

import (
	"log"

	meanstoanend "github.com/nint8835/protohackers/2-means-to-an-end"
	"github.com/nint8835/protohackers/pkg/server"
)

func main() {
	err := server.New(meanstoanend.HandleConn).Start()
	if err != nil {
		log.Fatalf("Error running server: %s", err)
	}
}
