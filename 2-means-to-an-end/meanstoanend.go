package meanstoanend

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/nint8835/protohackers/pkg/server"
)

type MessageType byte

const (
	MessageTypeInsert MessageType = 'I'
	MessageTypeQuery  MessageType = 'Q'
)

type Message struct {
	Type MessageType
	Arg1 int32
	Arg2 int32
}

type Price struct {
	Timestamp int32
	Price     int32
}

type Connection struct {
	Connection io.ReadWriteCloser
	Prices     []Price
}

func (conn *Connection) ReadMessage() (Message, error) {
	var message Message
	err := binary.Read(conn.Connection, binary.BigEndian, &message)
	if err != nil {
		return Message{}, fmt.Errorf("error reading data: %w", err)
	}

	return message, nil
}

func (conn *Connection) handleInsert(message Message) {
	conn.Prices = append(conn.Prices, Price{message.Arg1, message.Arg2})
}

func (conn *Connection) handleQuery(message Message) []byte {
	resp := bytes.NewBuffer([]byte{})

	if message.Arg2 < message.Arg1 {
		binary.Write(resp, binary.BigEndian, int32(0))
		return resp.Bytes()
	}

	sum := 0
	priceCount := 0

	for _, price := range conn.Prices {
		if price.Timestamp >= message.Arg1 && price.Timestamp <= message.Arg2 {
			sum += int(price.Price)
			priceCount++
		}
	}

	if priceCount == 0 {
		binary.Write(resp, binary.BigEndian, int32(0))
		return resp.Bytes()
	} else {
		binary.Write(resp, binary.BigEndian, int32(sum/priceCount))
		return resp.Bytes()
	}
}

func (conn *Connection) HandleMessage(message Message) []byte {
	switch message.Type {
	case MessageTypeInsert:
		conn.handleInsert(message)
		return []byte{}
	case MessageTypeQuery:
		return conn.handleQuery(message)
	default:
		log.Printf("Invalid message type: %v", message.Type)
		return []byte{}
	}
}

func HandleConn(conn server.Connection) {
	connection := Connection{
		Connection: conn,
		Prices:     []Price{},
	}

	for {
		message, err := connection.ReadMessage()
		if err != nil {
			if errors.Is(err, io.EOF) {
				conn.Close()
			} else {
				log.Printf("Unexpected error reading message: %s", err)
			}
			return
		}
		connection.Connection.Write(connection.HandleMessage(message))
	}
}
