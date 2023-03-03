package dht

import (
	"github.com/libp2p/go-libp2p-core/protocol"
)

var (
	// ProtocolDHT is the default DHT protocol.
	ProtocolDHT protocol.ID = "/ipfs/kad/1.0.0"
	// DefaultProtocols spoken by the DHT.
	DefaultProtocols = []protocol.ID{ProtocolDHT}
)
