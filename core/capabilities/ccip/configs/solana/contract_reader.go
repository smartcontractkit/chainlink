package solana

import (
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
)

var HomeChainReaderConfigRaw = config.ChainReader{} //  TODO update the home chain reader configuration
var DestReaderConfig = config.ChainReader{}         //  TODO update the Dest chain reader configuration
var SourceReaderConfig = config.ChainReader{}       // TODO update the Source chain reader configuration

func MergeReaderConfigs(configs ...config.ChainReader) config.ChainReader {
	allNamespaces := make(map[string]config.ChainReaderMethods)
	for _, c := range configs {
		for namespace, method := range c.Namespaces {
			allNamespaces[namespace] = method
		}
	}

	return config.ChainReader{Namespaces: allNamespaces}
}
