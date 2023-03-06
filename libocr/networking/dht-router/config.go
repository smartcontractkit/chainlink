package dhtrouter

import (
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	ocr1types "github.com/smartcontractkit/libocr/offchainreporting/types"
)

type DHTNodeConfig struct {
	prefix         protocol.ID
	extension      protocol.ID
	bootstrapNodes []peer.AddrInfo
	logger         loghelper.LoggerWithContext

	// node will check connections to all bootstrap nodes at this interval
	bootstrapCheckInterval time.Duration
	failureThreshold       int
	extendedDHTLogging     bool
	announcementUserPrefix uint32
}

func genExtension(cd ocr1types.ConfigDigest) protocol.ID {
	return protocol.ID(fmt.Sprintf("/%x", cd))
}

func BuildConfig(
	bootstrapNodes []peer.AddrInfo,
	prefix protocol.ID,
	configDigest ocr1types.ConfigDigest,
	logger loghelper.LoggerWithContext,
	bootstrapConnectionCheckInterval time.Duration,
	failureThreshold int,
	extendedDHTLogging bool,
	announcementUserPrefix uint32,
) DHTNodeConfig {
	c := DHTNodeConfig{
		bootstrapNodes:         bootstrapNodes,
		prefix:                 prefix,
		extension:              genExtension(configDigest),
		bootstrapCheckInterval: bootstrapConnectionCheckInterval,
		failureThreshold:       failureThreshold,
		extendedDHTLogging:     extendedDHTLogging,
		announcementUserPrefix: announcementUserPrefix,
	}

	c.logger = logger.MakeChild(commontypes.LogFields{
		"id":              "DHT",
		"protocolID":      c.ProtocolID(),
		"F":               failureThreshold,
		"extendedLogging": extendedDHTLogging,
	})

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
