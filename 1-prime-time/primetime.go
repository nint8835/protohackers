package main

import (
	"bufio"
	"encoding/json"
	"log"
	"math/big"

	"github.com/nint8835/protohackers/pkg/server"
)

var malformedResp = []byte("{}\n")

func handleConn(conn server.Connection) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		var data map[string]any

		err := json.Unmarshal(scanner.Bytes(), &data)
		if err != nil {
			log.Printf("Error decoding data: %s", err)
			conn.Write(malformedResp)
			continue
		}

		method, hasMethod := data["method"]
		number, hasNumber := data["number"]
		numberFloat, numberIsFloat := number.(float64)

		if !hasMethod || !hasNumber || !numberIsFloat || method != "isPrime" {
			log.Println(hasMethod, hasNumber, numberIsFloat, method)
			conn.Write(malformedResp)
			continue
		}

		resp := map[string]any{
			"method": "isPrime",
			"prime":  big.NewInt(int64(numberFloat)).ProbablyPrime(0),
		}

		respBytes, _ := json.Marshal(resp)
		conn.Write(respBytes)
		conn.Write([]byte("\n"))
	}

	conn.Close()
}

func main() {
	err := server.New(handleConn).Start()
	if err != nil {
		log.Fatalf("Error running server: %s", err)
	}
}
