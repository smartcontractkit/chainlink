package config

import (
	"time"

	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

// P2PNetworking is a subset of global config relevant to p2p networking.
type P2PNetworking interface {
	P2PNetworkingStack() (n ocrnetworking.NetworkingStack)
	P2PNetworkingStackRaw() string
	P2PPeerID() p2pkey.PeerID
	P2PPeerIDRaw() string
	P2PIncomingMessageBufferSize() int
	P2POutgoingMessageBufferSize() int
}

type P2PDeprecated interface {
	// DEPRECATED - HERE FOR BACKWARDS COMPATIBILITY
	ocrNewStreamTimeout() time.Duration
	ocrBootstrapCheckInterval() time.Duration
	ocrDHTLookupInterval() int
	ocrIncomingMessageBufferSize() int
	ocrOutgoingMessageBufferSize() int
}
