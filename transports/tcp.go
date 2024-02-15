package transports

import (
	"fmt"
	"log"
	"net"
)

// TCPTransport handles iRTSP connections over TCP
type TCPTransport struct {
	onPacketHandler func(connection net.Conn, data []byte)
}

// OnPacket sets the handler to be called when a new message packet is received
func (tp *TCPTransport) OnPacket(handler func(connection net.Conn, data []byte)) {
	tp.onPacketHandler = handler
}

// Listen starts the TCP transport on the given port
func (tp *TCPTransport) Listen(port uint16) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go tp.handleConnetion(conn)
	}
}

func (tp *TCPTransport) handleConnetion(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}

		if bytesRead > 0 {
			message := buffer[:bytesRead]

			tp.onPacketHandler(conn, message)
		}
	}
}

// NewTCPTransport returns a new TCPTransport
func NewTCPTransport() *TCPTransport {
	return &TCPTransport{}
}
