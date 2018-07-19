package network

import (
	"fmt"
	"sync"

	"github.com/pions/webrtc/internal/datachannel"
	"github.com/pions/webrtc/internal/dtls"
	"github.com/pions/webrtc/internal/sctp"
	"github.com/pions/webrtc/internal/srtp"
	"github.com/pions/webrtc/pkg/rtp"
	"github.com/pkg/errors"
)

// State is all the network state (DTLS, SRTP) that is shared between connections
type State struct {
	icePwd      []byte
	iceNotifier ICENotifier

	dtlsState *dtls.State
	certPair  *dtls.CertPair

	BufferTransportGenerator BufferTransportGenerator
	bufferTransports         map[uint32]chan<- *rtp.Packet

	// https://tools.ietf.org/html/rfc3711#section-3.2.3
	// A cryptographic context SHALL be uniquely identified by the triplet
	//  <SSRC, destination network address, destination transport port number>
	// contexts are keyed by IP:PORT:SSRC
	srtpContextsLock *sync.Mutex
	srtpContexts     map[string]*srtp.Context

	sctpAssociation *sctp.Association
}

// StateArgs are the arguments required to initalize a new network State
type StateArgs struct {
	ICEPwd                   []byte
	BufferTransportGenerator BufferTransportGenerator
	DataChannelEventHandler  DataChannelEventHandler
}

// NewState creates a new network state, that is shared across all ports
func NewState(args *StateArgs) *State {
	s := &State{}

	s.sctpAssociation = sctp.NewAssocation(func(pkt *sctp.Packet) {
		raw, err := pkt.Marshal()
		if err != nil {
			fmt.Println(errors.Wrap(err, "Failed to Marshal SCTP packet"))
			return
		}

		fmt.Println(raw)
		// TODO loop over all ports, use first port with good ICE state
		// d.Send(raw)
	}, func(data []byte, streamIdentifier uint16) {
		msg, err := datachannel.Parse(data)
		if err != nil {
			fmt.Println(errors.Wrap(err, "Failed to parse DataChannel packet"))
			return
		}
		switch m := msg.(type) {
		case *datachannel.ChannelOpen:
			args.DataChannelEventHandler(&DataChannelCreated{streamIdentifier: streamIdentifier, Label: string(m.Label)})
		case *datachannel.Data:
			args.DataChannelEventHandler(&DataChannelMessage{streamIdentifier: streamIdentifier, Body: m.Data})
		default:
			fmt.Println("Unhandled DataChannel message", m)
		}
	})

	return nil
}

// AddPort links the state to a network port
func (s *State) AddPort(port *Port) {
}

// Close cleans up all the allocated state
func (s *State) Close() {
	s.sctpAssociation.Close()
	s.dtlsState.Close()
}
