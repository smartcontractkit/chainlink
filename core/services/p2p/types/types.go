package types

import (
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

const PeerIDLength = 32

type PeerID = ragetypes.PeerID

//go:generate mockery --quiet --name Peer --output ./mocks/ --case=underscore
type Peer interface {
	services.Service
	ID() PeerID
	UpdateConnections(peers map[PeerID]StreamConfig) error
	Send(peerID PeerID, msg []byte) error
	Receive() <-chan Message
}

//go:generate mockery --quiet --name PeerWrapper --output ./mocks/ --case=underscore
type PeerWrapper interface {
	services.Service
	GetPeer() Peer
}

//go:generate mockery --quiet --name Signer --output ./mocks/ --case=underscore
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
