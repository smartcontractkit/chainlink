package ocrcommon

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/common/config"
)

type Config interface {
	ChainType() config.ChainType
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

// GetValidatedBootstrapPeers will error unless at least one valid bootstrap peer is found
func GetValidatedBootstrapPeers(specPeers []string, configPeers []commontypes.BootstrapperLocator) ([]commontypes.BootstrapperLocator, error) {
	bootstrapPeers, err := ParseBootstrapPeers(specPeers)
	if err != nil {
		return nil, err
	}
	if len(bootstrapPeers) == 0 {
		if len(configPeers) == 0 {
			return nil, errors.New("no bootstrappers found")
		}
		return configPeers, nil
	}
	return bootstrapPeers, nil
}
