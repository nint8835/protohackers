package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Database struct {
	conn *net.UDPConn

	values map[string]string
	valuesMux *sync.Mutex
}

func (d *Database) Insert(key string, value string) {
	d.valuesMux.Lock()
	defer d.valuesMux.Unlock()

	if key == "version" {
		return
	}

	d.values[key] = value
}

func (d *Database) Retrieve(key string) string {
	return d.values[key]
}

func (d *Database) HandleMessage(msg string, addr *net.UDPAddr) {
	log.Printf("(%v): %s", addr, msg)

	key, val, isInsert := strings.Cut(msg, "=")

	if isInsert {
		d.Insert(key, val)
	} else {
		d.conn.WriteToUDP([]byte(fmt.Sprintf("%s=%s", key, d.Retrieve(key))), addr)
	}
}

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 9191,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		panic(err)
	}

	db := &Database{
		conn: conn,
		values: map[string]string{
			"version": "nint8835/protohackers (4-unusual-database-program)",
		},
		valuesMux: new(sync.Mutex),
	}

	for {
		buffer := make([]byte, 1024)

		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("error reading: %s", err)
			continue
		}

		go db.HandleMessage(string(buffer[:n]), addr)
	}
}
