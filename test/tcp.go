package main

import (
	"fmt"
	"log"
	"net"

	"github.com/PandoraStream/irtsp-go"
	"github.com/PandoraStream/irtsp-go/messages"
	"github.com/PandoraStream/irtsp-go/transports"
)

func main() {
	transport := transports.NewTCPTransport()
	server := irtsp.NewServer(transport)

	server.OnStartRequest = onStartRequest

	go server.Listen(8080)

	clientTest()
}

func onStartRequest(conn net.Conn, request *messages.StartRequest) {
	response := messages.NewStartResponseMessage(200, nil)

	response.SequenceNumber = request.SequenceNumber
	response.Scheme = request.Scheme
	response.Timings = []uint64{45850, 1806955874, 1806955874}

	conn.Write(response.Bytes())
}

func clientTest() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	request := messages.NewStartRequestMessage(nil)

	request.SequenceNumber = 0
	request.Scheme = ""
	request.Timings = []uint64{1454326}

	_, err = conn.Write(request.Bytes())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	buffer := make([]byte, 1024)

	for {
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Received: %s\n", buffer[:bytesRead])
	}
}
