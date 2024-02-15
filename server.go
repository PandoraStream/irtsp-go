// Package irtsp implements the iRTSP protocol created by Ubitus
package irtsp

import (
	"net"
	"strings"

	"github.com/PandoraStream/irtsp-go/messages"
	"github.com/PandoraStream/irtsp-go/transports"
)

// Server implements the lowest level parts of the iRTSP protocol
type Server struct {
	transport          transports.Transport
	OnMessageRequest   func(conn net.Conn, message *messages.Message)
	OnStartRequest     func(conn net.Conn, message *messages.StartRequest)
	OnSetupRequest     func(conn net.Conn, message *messages.SetupRequest)
	OnKnockRequest     func(conn net.Conn, message *messages.KnockRequest)
	OnPushInfoRequest  func(conn net.Conn, message *messages.PushInfoRequest)
	OnUserInfoRequest  func(conn net.Conn, message *messages.UserInfoRequest)
	OnRemoteCTLRequest func(conn net.Conn, message *messages.RemoteCTLRequest)
	OnSetRequest       func(conn net.Conn, message *messages.SetRequest)
	OnPlayRequest      func(conn net.Conn, message *messages.PlayRequest)
	OnEchoRequest      func(conn net.Conn, message *messages.EchoRequest)
	OnTeardownRequest  func(conn net.Conn, message *messages.TeardownRequest)
}

// Listen starts the server on the given port using the configured transport protocol
func (s *Server) Listen(port uint16) {
	s.transport.Listen(port)
}

func (s *Server) onPacket(conn net.Conn, data []byte) {
	// TODO - Move the nex-go packet reordering system to here, since RTSP can be over UDP
	// TODO - Don't assume only one message per packet. Move the nex-go multi-packet scanning system here
	// TODO - State management. Use a Client struct?

	contents := string(data)

	// TODO - These are dirty checks, replace with something better?
	if strings.Contains(contents, "SET/START") && s.OnStartRequest != nil {
		s.OnStartRequest(conn, messages.NewStartRequestMessage(data))
	}
}

// NewServer returns a new iRTSP server
func NewServer(transport transports.Transport) *Server {
	server := &Server{
		transport: transport,
	}

	server.transport.OnPacket(server.onPacket)

	return server
}
