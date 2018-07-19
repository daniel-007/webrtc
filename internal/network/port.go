package network

import (
	"net"
	"strconv"

	"github.com/pions/pkg/stun"
	"github.com/pions/webrtc/internal/dtls"
	"github.com/pions/webrtc/pkg/ice"
	"github.com/pkg/errors"
	"golang.org/x/net/ipv4"
)

// Port represents a UDP listener that handles incoming/outgoing traffic
type Port struct {
	ICEState ice.ConnectionState

	conn          *ipv4.PacketConn
	listeningAddr *stun.TransportAddr
	seenPeers     map[string]bool

	sharedState *State
}

// NewPort creates a new Port
func NewPort(address string, state *State) (*Port, error) {
	if state == nil {
		return nil, errors.Errorf("network.State must not be nil")
	}

	listener, err := net.ListenPacket("udp4", address)
	if err != nil {
		return nil, err
	}

	addr, err := stun.NewTransportAddr(listener.LocalAddr())
	if err != nil {
		return nil, err
	}

	srcString := addr.IP.String() + ":" + strconv.Itoa(addr.Port)
	conn := ipv4.NewPacketConn(listener)
	dtls.AddListener(srcString, conn)

	p := &Port{
		listeningAddr: addr,
		conn:          conn,
		sharedState:   state,
		seenPeers:     make(map[string]bool),
	}

	go p.networkLoop()
	return p, nil
}

// Close closes the listening port and cleans up any state
func (p *Port) Close() error {
	return p.conn.Close()
}
