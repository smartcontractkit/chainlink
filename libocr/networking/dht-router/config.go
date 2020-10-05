package dhtrouter

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

type DHTNodeConfig struct {
	prefix             protocol.ID
	extension          protocol.ID
	bootstrapNodes     []peer.AddrInfo
	logger             types.Logger
	failureThreshold   int
	extendedDHTLogging bool
}

func BuildConfig(bootstrapNodes []peer.AddrInfo, prefix protocol.ID, configDigest types.ConfigDigest,
	logger types.Logger, failureThreshold int, extendedDHTLogging bool) DHTNodeConfig {
	extension := protocol.ID(fmt.Sprintf("/%x", configDigest))

	c := DHTNodeConfig{
		bootstrapNodes:     bootstrapNodes,
		prefix:             prefix,
		extension:          extension,
		failureThreshold:   failureThreshold,
		extendedDHTLogging: extendedDHTLogging,
	}

	logger = loghelper.MakeLoggerWithContext(logger, types.LogFields{
		"id":              "DHT",
		"protocolID":      c.ProtocolID(),
		"F":               failureThreshold,
		"extendedLogging": extendedDHTLogging,
	})

	c.logger = logger

	return c
}

func (config DHTNodeConfig) ProtocolID() protocol.ID {
	return protocol.ID(fmt.Sprintf("%s%s/kad/1.0.0", config.prefix, config.extension))
}

func (config DHTNodeConfig) String() string {
	s := ""
	if len(config.bootstrapNodes) > 0 {
		s += "bootnodes: "
		for _, b := range config.bootstrapNodes {
			s += b.String()
			s += ","
		}
		s += "; "
	}

	s += fmt.Sprintf("ns=%s", config.prefix)

	return s
}

func (config *DHTNodeConfig) AddBootstrapNodes(addrs []peer.AddrInfo) {
	config.bootstrapNodes = append(config.bootstrapNodes, addrs...)
}
