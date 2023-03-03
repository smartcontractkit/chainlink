package protocol

import (
	"log"

	"github.com/smartcontractkit/libocr/commontypes"
)

// NetworkSender sends messages to other oracles
type NetworkSender interface {
	// SendTo(msg, to) sends msg to "to"
	SendTo(msg Message, to commontypes.OracleID)
	// Broadcast(msg) sends msg to all oracles
	Broadcast(msg Message)
}

// NetworkEndpoint sends & receives messages to/from other oracles
type NetworkEndpoint interface {
	NetworkSender
	// Receive returns channel which carries all messages sent to this oracle
	Receive() <-chan MessageWithSender

	// Start must be called before Receive. Calling Start more than once causes
	// panic.
	Start() error

	// Close must be called before receive. Close can be called multiple times.
	// Close can be called even on an unstarted NetworkEndpoint.
	Close() error
}

// SimpleNetwork is a strawman (in-memory) implementation of the Network
// interface. Network channels are buffered and can queue up to 100 messages
// before blocking.
type SimpleNetwork struct {
	chs []chan MessageWithSender // i'th channel models oracle i's network
}

// NewSimpleNetwork returns a SimpleNetwork for n oracles
func NewSimpleNetwork(n int) *SimpleNetwork {
	s := SimpleNetwork{}
	for i := 0; i < n; i++ {
		s.chs = append(s.chs, make(chan MessageWithSender, 100))
	}
	return &s
}

// Endpoint returns the interface for oracle id's networking facilities
func (net *SimpleNetwork) Endpoint(id commontypes.OracleID) (NetworkEndpoint, error) {
	return SimpleNetworkEndpoint{
		net,
		id,
	}, nil
}

// SimpleNetworkEndpoint is a strawman (in-memory) implementation of
// NetworkEndpoint, used by SimpleNetwork
type SimpleNetworkEndpoint struct {
	net *SimpleNetwork       // Reference back to network for all participants
	id  commontypes.OracleID // Index of oracle this endpoint pertains to
}

var _ NetworkEndpoint = (*SimpleNetworkEndpoint)(nil)

// SendTo sends msg to oracle "to"
func (end SimpleNetworkEndpoint) SendTo(msg Message, to commontypes.OracleID) {
	log.Printf("[%v] sending to %v: %T\n", end.id, to, msg)
	end.net.chs[to] <- MessageWithSender{msg, end.id}
}

// Broadcast sends msg to all participating oracles
func (end SimpleNetworkEndpoint) Broadcast(msg Message) {
	log.Printf("[%v] broadcasting: %T\n", end.id, msg)
	for _, ch := range end.net.chs {
		ch <- MessageWithSender{msg, end.id}
	}
}

// Receive returns a channel which carries all messages sent to the oracle
func (end SimpleNetworkEndpoint) Receive() <-chan MessageWithSender {
	return end.net.chs[end.id]
}

// Start satisfies the interface
func (SimpleNetworkEndpoint) Start() error { return nil }

// Close satisfies the interface
func (SimpleNetworkEndpoint) Close() error { return nil }
