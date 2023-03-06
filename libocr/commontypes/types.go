package commontypes

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

// OracleID is an index over the oracles, used as a succinct attribution to an
// oracle in communication with the on-chain contract. It is not a cryptographic
// commitment to the oracle's private key, like a public key is.
type OracleID uint8

// BootstrapperLocator contains information for locating a bootstrapper on the network.
// It is encoded like PeerID@Addr[0]/Addr[1]/.../Addr[len(Addr)-1].
// Sample encoding:
// 12D3KooWQzePGqHw66cV1Qsm71eGZKiPEgALYYM3inPtFYibZ67e@192.168.1.1:1234/192.168.1.2:2345
type BootstrapperLocator struct {
	// PeerID is the libp2p-style peer ID of the bootstrapper
	PeerID string

	// Addrs contains the addresses of the bootstrapper. An address must be of the form "<host>:<port>",
	// such as "52.49.198.28:80" or "chain.link:443".
	Addrs []string
}

func NewBootstrapperLocator(peerID string, addrs []string) (*BootstrapperLocator, error) {
	if err := (&ragetypes.PeerID{}).UnmarshalText([]byte(peerID)); err != nil {
		return nil, fmt.Errorf("invalid peer id (%q): %w", peerID, err)
	}
	for _, address := range addrs {
		_, _, err := net.SplitHostPort(address)
		if err != nil {
			return nil, fmt.Errorf("invalid address (%q) for bootstrapper (%q): %w", address, peerID, err)
		}
	}
	return &BootstrapperLocator{peerID, addrs}, nil
}

func (b *BootstrapperLocator) MarshalText() ([]byte, error) {
	var bs bytes.Buffer
	bs.WriteString(b.PeerID)
	bs.WriteRune('@')
	for i, addr := range b.Addrs {
		if i != 0 {
			bs.WriteRune('/')
		}
		bs.WriteString(addr)
	}
	return bs.Bytes(), nil
}

func (b *BootstrapperLocator) UnmarshalText(text []byte) error {
	parts := strings.Split(string(text), "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid BootstrapperLocator, expected format is PeerID@Host1:Port1/Host2:Port2")
	}
	peerID, addrsJoined := parts[0], parts[1]
	addrs := []string{}
	if len(addrsJoined) > 0 {
		addrs = strings.Split(addrsJoined, "/")
	}
	bNew, err := NewBootstrapperLocator(peerID, addrs)
	if err != nil {
		return err
	}
	b.PeerID = bNew.PeerID
	b.Addrs = bNew.Addrs
	return nil
}

// BinaryMessageWithSender contains the information from a Receive() channel
// message: The binary representation of the message, and the ID of its sender.
type BinaryMessageWithSender struct {
	Msg    []byte
	Sender OracleID
}

// BinaryNetworkEndpoint contains the network methods a consumer must implement
// SendTo and Broadcast must not block. They should buffer messages and
// (optionally) drop the oldest buffered messages if the buffer reaches capacity.
//
// The protocol trusts the sender in BinaryMessageWithSender. Implementors of
// this interface are responsible for securely authenticating that messages come
// from their indicated senders.
//
// All its functions should be thread-safe.
type BinaryNetworkEndpoint interface {
	// SendTo(msg, to) sends msg to "to"
	SendTo(payload []byte, to OracleID)
	// Broadcast(msg) sends msg to all oracles
	Broadcast(payload []byte)
	// Receive returns channel which carries all messages sent to this oracle.
	Receive() <-chan BinaryMessageWithSender
	// Start starts the endpoint
	Start() error
	// Close stops the endpoint. Calling this multiple times may return an
	// error, but must not panic.
	Close() error
}

// Bootstrapper helps nodes find each other on the network level by providing
// peer-discovery services.
//
// All its functions should be thread-safe.
type Bootstrapper interface {
	Start() error
	// Close closes the bootstrapper. Calling this multiple times may return an
	// error, but must not panic.
	Close() error
}

// MonitoringEndpoint is where the OCR protocol sends monitoring output
//
// All its functions should be thread-safe.
type MonitoringEndpoint interface {
	SendLog(log []byte)
}
