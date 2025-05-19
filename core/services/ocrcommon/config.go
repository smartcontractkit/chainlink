package ocrcommon

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-evm/pkg/config/chaintype"
)

type Config interface {
	ChainType() chaintype.ChainType
}

func ParseBootstrapPeers(peers []string) (bootstrapPeers []commontypes.BootstrapperLocator, err error) {
	for _, bs := range peers {
		var bsl commontypes.BootstrapperLocator
		err = bsl.UnmarshalText([]byte(bs))
		if err != nil {
			return nil, err
		}
		bootstrapPeers = append(bootstrapPeers, bsl)
	}
	return
}

// GetValidatedBootstrapPeers will error unless at least one valid bootstrap is found or
// no bootstrappers are found and allowNoBootstrappers is true.
func GetValidatedBootstrapPeers(specPeers []string, configPeers []commontypes.BootstrapperLocator, allowNoBootstrappers bool) ([]commontypes.BootstrapperLocator, error) {
	bootstrapPeers, err := ParseBootstrapPeers(specPeers)
	if err != nil {
		return nil, err
	}
	if len(bootstrapPeers) == 0 {
		if len(configPeers) == 0 {
			if !allowNoBootstrappers {
				return nil, errors.New("no bootstrappers found")
			}
			// Bootstrappers may be empty if the node is not configured to conduct consensus (i.e. f = 0 and n = 1).
			// This is useful for testing and development.
			return nil, nil
		}
		return configPeers, nil
	}
	return bootstrapPeers, nil
}
