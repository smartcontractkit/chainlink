package blockchain

import (
	"github.com/rs/zerolog/log"
)

type ClientImplementation string

const (
	// Ethereum uses the standard EVM implementation, and is considered default
	EthereumClientImplementation ClientImplementation = "Ethereum"
	// MetisClientImplementation uses a client tailored for Metis EVM networks
	MetisClientImplementation ClientImplementation = "Metis"
	// KlaytnClientImplementation uses a client tailored for Klaytn EVM networks
	KlaytnClientImplementation ClientImplementation = "Klaytn"
	// OptimismClientImplementation uses a client tailored for Optimism EVM networks
	OptimismClientImplementation ClientImplementation = "Optimism"
	// ArbitrumClientImplementation uses a client tailored for Arbitrum EVM networks
	ArbitrumClientImplementation ClientImplementation = "Arbitrum"
	// PolygonClientImplementation uses a client tailored for Polygon EVM networks
	PolygonClientImplementation ClientImplementation = "Polygon"
	// RSKClientImplementation uses a client tailored for RSK EVM networks
	RSKClientImplementation ClientImplementation = "RSK"
)

// wrapSingleClient Wraps a single EVM client in its appropriate implementation, based on the chain ID
func wrapSingleClient(networkSettings EVMNetwork, client *EthereumClient) EVMClient {
	var wrappedEc EVMClient
	switch networkSettings.ClientImplementation {
	case EthereumClientImplementation:
		wrappedEc = client
	case MetisClientImplementation:
		wrappedEc = &MetisClient{client}
	case PolygonClientImplementation:
		wrappedEc = &PolygonClient{client}
	case KlaytnClientImplementation:
		wrappedEc = &KlaytnClient{client}
	case ArbitrumClientImplementation:
		wrappedEc = &ArbitrumClient{client}
	case OptimismClientImplementation:
		wrappedEc = &OptimismClient{client}
	case RSKClientImplementation:
		wrappedEc = &RSKClient{client}
	default:
		wrappedEc = client
	}
	return wrappedEc
}

// wrapMultiClient Wraps a multi-node EVM client in its appropriate implementation, based on the chain ID
func wrapMultiClient(networkSettings EVMNetwork, client *EthereumMultinodeClient) EVMClient {
	var wrappedEc EVMClient
	logMsg := log.Info().Str("Network", networkSettings.Name)
	switch networkSettings.ClientImplementation {
	case EthereumClientImplementation:
		logMsg.Msg("Using Standard Ethereum Client")
		wrappedEc = client
	case PolygonClientImplementation:
		logMsg.Msg("Using Polygon Client")
		wrappedEc = &PolygonMultinodeClient{client}
	case MetisClientImplementation:
		logMsg.Msg("Using Metis Client")
		wrappedEc = &MetisMultinodeClient{client}
	case KlaytnClientImplementation:
		logMsg.Msg("Using Klaytn Client")
		wrappedEc = &KlaytnMultinodeClient{client}
	case ArbitrumClientImplementation:
		logMsg.Msg("Using Arbitrum Client")
		wrappedEc = &ArbitrumMultinodeClient{client}
	case OptimismClientImplementation:
		logMsg.Msg("Using Optimism Client")
		wrappedEc = &OptimismMultinodeClient{client}
	case RSKClientImplementation:
		logMsg.Msg("Using RSK Client")
		wrappedEc = &RSKMultinodeClient{client}
	default:
		log.Warn().Str("Network", networkSettings.Name).Msg("Unknown client implementation, defaulting to standard Ethereum client")
		wrappedEc = client
	}
	return wrappedEc
}
