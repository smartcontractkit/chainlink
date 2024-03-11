package types

import (
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

//go:generate mockery --quiet --name Peer --output ./mocks/ --case=underscore
type Peer interface {
	services.Service
	UpdateConnections(peers map[ragetypes.PeerID]StreamConfig) error
	Send(peerID ragetypes.PeerID, msg []byte) error
	Receive() <-chan Message
}

type Message struct {
	Sender  ragetypes.PeerID
	Payload []byte
}

type StreamConfig struct {
	IncomingMessageBufferSize int
	OutgoingMessageBufferSize int
	MaxMessageLenBytes        int
	MessageRateLimiter        ragep2p.TokenBucketParams
	BytesRateLimiter          ragep2p.TokenBucketParams
}
