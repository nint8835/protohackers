package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/dlclark/regexp2"

	"github.com/nint8835/protohackers/pkg/server"
)

var messageRegex = regexp2.MustCompile(`^\[(.*)] (.*)$`, 0)
var boguscoinRegex = regexp2.MustCompile(`(?:^|(?<= ))(7[a-zA-Z0-9]{25,34})(?:$|(?= ))`, 0)

const targetAddr = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"

const remoteAddr = "chat.protohackers.com:16963"

//const remoteAddr = "localhost:3000"

type ProxyClient struct {
	inConn    server.Connection
	outConn   server.Connection
	connected bool
}

func (p *ProxyClient) HandleIn() {
	reader := bufio.NewReader(p.inConn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from client: %s", err)
			break
		}
		line = strings.TrimSuffix(line, "\n")
		log.Printf("Client -> Server: %s", line)

		if !p.connected {
			p.connected = true
		} else {
			updatedContent, err := boguscoinRegex.Replace(line, targetAddr, -1, -1)
			if err != nil {
				log.Printf("error updating content: %s", err)
			}
			line = updatedContent
		}
		p.outConn.Write([]byte(fmt.Sprintf("%s\n", line)))
	}

	p.outConn.Close()
}

func (p *ProxyClient) HandleOut() {
	scanner := bufio.NewScanner(p.outConn)

	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("Server -> Client: %s", line)
		messageContentMatches, err := messageRegex.FindStringMatch(line)
		if err != nil {
			log.Printf("error running message content regex: %s", err)
		}

		if messageContentMatches != nil {
			username := messageContentMatches.Groups()[1].String()
			messageContent := messageContentMatches.Groups()[2].String()

			updatedContent, err := boguscoinRegex.Replace(messageContent, targetAddr, -1, -1)
			if err != nil {
				log.Printf("error updating content: %s", err)
			}

			log.Printf("Old content: %s", messageContent)
			log.Printf("New content: %s", updatedContent)

			line = fmt.Sprintf("[%s] %s", username, updatedContent)
		}
		p.inConn.Write([]byte(fmt.Sprintf("%s\n", line)))
	}

	p.inConn.Close()
}

func (p *ProxyClient) Run() {
	go p.HandleIn()
	p.HandleOut()
}

func handleConn(conn server.Connection) {
	outConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("Error dialing outbound connection: %s", err)
	}

	client := ProxyClient{
		inConn:  conn,
		outConn: outConn,
	}
	client.Run()
}

func main() {
	s := server.New(handleConn)
	s.Addr = ":3001"
	err := s.Start()
	if err != nil {
		log.Fatalf("Error running server: %s", err)
	}
}
