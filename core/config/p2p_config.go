package config

import (
	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type P2P interface {
	V2() V2
	V1() V1
	NetworkStack() (n ocrnetworking.NetworkingStack)
	PeerID() p2pkey.PeerID
	IncomingMessageBufferSize() int
	OutgoingMessageBufferSize() int
	TraceLogging() bool
	Enabled() bool
}
