// Package transports implements TCP and UDP transport layers
package transports

import "net"

// Transport defines all the methods a transport protocol should have
type Transport interface {
	OnPacket(handler func(connection net.Conn, data []byte))
	Listen(port uint16)
}
