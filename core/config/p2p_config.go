package config

import (
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
)

type P2P interface {
	V2() V2
	PeerID() p2pkey.PeerID
	IncomingMessageBufferSize() int
	OutgoingMessageBufferSize() int
	TraceLogging() bool
	EnableExperimentalRageP2P() bool
	Enabled() bool
}
