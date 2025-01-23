package solana

import (
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
)

var DestReaderConfig = config.ContractReader{}   //  TODO update the Dest chain reader configuration
var SourceReaderConfig = config.ContractReader{} // TODO update the Source chain reader configuration

func MergeReaderConfigs(configs ...config.ContractReader) config.ContractReader {
	allNamespaces := make(map[string]config.ChainContractReader)
	for _, c := range configs {
		for namespace, method := range c.Namespaces {
			allNamespaces[namespace] = method
		}
	}

	return config.ContractReader{Namespaces: allNamespaces}
}
