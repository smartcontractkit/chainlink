package types

import (
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

const PeerIDLength = 32

type PeerID = ragetypes.PeerID

type Peer interface {
	services.Service
	ID() PeerID
	UpdateConnections(peers map[PeerID]StreamConfig) error
	Send(peerID PeerID, msg []byte) error
	Receive() <-chan Message
}

type PeerWrapper interface {
	services.Service
	GetPeer() Peer
}

type Signer interface {
	Sign(data []byte) ([]byte, error)
}

type Message struct {
	Sender  PeerID
	Payload []byte
}

type StreamConfig struct {
	IncomingMessageBufferSize int
	OutgoingMessageBufferSize int
	MaxMessageLenBytes        int
	MessageRateLimiter        ragep2p.TokenBucketParams
	BytesRateLimiter          ragep2p.TokenBucketParams
}
