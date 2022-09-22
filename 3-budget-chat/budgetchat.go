package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"

	"github.com/nint8835/protohackers/pkg/server"
)

var ValidUsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

type ClientState int

const (
	ClientStateAwaitingUsername ClientState = iota
	ClientStateJoined
	ClientStateDisconnected
)

type Client struct {
	conn  server.Connection
	name  string
	state ClientState

	server *Server
}

func (c *Client) WriteLine(line string) {
	_, err := c.conn.Write([]byte(line + "\n"))
	if err != nil {
		log.Printf("Error writing line to client: %s", err)
	}
}

func (c *Client) HandleConn() {
	scanner := bufio.NewScanner(c.conn)

	c.WriteLine("Welcome to budget-chat, please provide a username.")

	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("Line: %s", line)

		switch c.state {
		case ClientStateAwaitingUsername:
			if !ValidUsernameRegex.Match([]byte(line)) {
				c.WriteLine("Invalid username.")
				return
			}
			c.name = line
			c.state = ClientStateJoined
			log.Printf("Client set username to %s", line)

			var clientNameList []string

			for _, client := range c.server.clients {
				if client.state == ClientStateJoined && client.name != c.name {
					clientNameList = append(clientNameList, client.name)
					client.WriteLine(fmt.Sprintf("* %s has entered the room", c.name))
				}
			}

			c.WriteLine(fmt.Sprintf("* This room contains: %s", strings.Join(clientNameList, ", ")))
		case ClientStateJoined:
			for _, client := range c.server.clients {
				if client.state == ClientStateJoined && client.name != c.name {
					client.WriteLine(fmt.Sprintf("[%s] %s", c.name, line))
				}
			}
		}
	}
}

type Server struct {
	clients     map[net.Addr]*Client
	clientsLock *sync.Mutex
}

func (s *Server) HandleConn(conn server.Connection) {
	client := &Client{
		conn:   conn,
		name:   "",
		state:  ClientStateAwaitingUsername,
		server: s,
	}
	s.clientsLock.Lock()
	s.clients[conn.RemoteAddr()] = client
	s.clientsLock.Unlock()

	log.Printf("Client %s connected", client.conn.RemoteAddr())
	client.HandleConn()
	log.Printf("Client %s disconnected", client.conn.RemoteAddr())
	client.conn.Close()

	s.clientsLock.Lock()
	delete(s.clients, client.conn.RemoteAddr())
	s.clientsLock.Unlock()

	if client.state == ClientStateJoined {
		for _, existingClient := range s.clients {
			existingClient.WriteLine(fmt.Sprintf("* %s has left the room", client.name))
		}
	}

	client.state = ClientStateDisconnected
}

func main() {
	serverInst := Server{
		clients:     map[net.Addr]*Client{},
		clientsLock: &sync.Mutex{},
	}

	err := server.New(serverInst.HandleConn).Start()
	if err != nil {
		log.Fatalf("Error running server: %s", err)
	}
}
