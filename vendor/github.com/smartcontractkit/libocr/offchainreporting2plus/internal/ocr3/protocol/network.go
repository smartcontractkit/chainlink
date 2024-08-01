package protocol

import (
	"log"

	"github.com/smartcontractkit/libocr/commontypes"
)

// NetworkSender sends messages to other oracles
type NetworkSender[RI any] interface {
	// SendTo(msg, to) sends msg to "to"
	SendTo(msg Message[RI], to commontypes.OracleID)
	// Broadcast(msg) sends msg to all oracles
	Broadcast(msg Message[RI])
}

// NetworkEndpoint sends & receives messages to/from other oracles
type NetworkEndpoint[RI any] interface {
	NetworkSender[RI]
	// Receive returns channel which carries all messages sent to this oracle
	Receive() <-chan MessageWithSender[RI]

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
type SimpleNetwork[RI any] struct {
	chs []chan MessageWithSender[RI] // i'th channel models oracle i's network
}

// NewSimpleNetwork returns a SimpleNetwork for n oracles
func NewSimpleNetwork[RI any](n int) *SimpleNetwork[RI] {
	s := SimpleNetwork[RI]{}
	for i := 0; i < n; i++ {
		s.chs = append(s.chs, make(chan MessageWithSender[RI], 100))
	}
	return &s
}

// Endpoint returns the interface for oracle id's networking facilities
func (net *SimpleNetwork[RI]) Endpoint(id commontypes.OracleID) (NetworkEndpoint[RI], error) {
	return SimpleNetworkEndpoint[RI]{
		net,
		id,
	}, nil
}

// SimpleNetworkEndpoint is a strawman (in-memory) implementation of
// NetworkEndpoint, used by SimpleNetwork
type SimpleNetworkEndpoint[RI any] struct {
	net *SimpleNetwork[RI]   // Reference back to network for all participants
	id  commontypes.OracleID // Index of oracle this endpoint pertains to
}

var _ NetworkEndpoint[struct{}] = (*SimpleNetworkEndpoint[struct{}])(nil)

// SendTo sends msg to oracle "to"
func (end SimpleNetworkEndpoint[RI]) SendTo(msg Message[RI], to commontypes.OracleID) {
	log.Printf("[%v] sending to %v: %T\n", end.id, to, msg)
	end.net.chs[to] <- MessageWithSender[RI]{msg, end.id}
}

// Broadcast sends msg to all participating oracles
func (end SimpleNetworkEndpoint[RI]) Broadcast(msg Message[RI]) {
	log.Printf("[%v] broadcasting: %T\n", end.id, msg)
	for _, ch := range end.net.chs {
		ch <- MessageWithSender[RI]{msg, end.id}
	}
}

// Receive returns a channel which carries all messages sent to the oracle
func (end SimpleNetworkEndpoint[RI]) Receive() <-chan MessageWithSender[RI] {
	return end.net.chs[end.id]
}

// Start satisfies the interface
func (SimpleNetworkEndpoint[RI]) Start() error { return nil }

// Close satisfies the interface
func (SimpleNetworkEndpoint[RI]) Close() error { return nil }
