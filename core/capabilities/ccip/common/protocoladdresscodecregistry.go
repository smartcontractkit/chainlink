package common

import (
	"fmt"
	"maps"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

var registeredProtocolAddressCodecs = make(map[string]ProtocolAddressCodec)

// RegisterProtocolAddressCodec registers a ProtocolAddressCodec for a chain family.
// Called from init() in each ccip<chain> package.
func RegisterProtocolAddressCodec(family string, codec ProtocolAddressCodec) {
	registeredProtocolAddressCodecs[family] = codec
}

// ProtocolAddressCodecs returns a registry populated from all registered codecs.
func ProtocolAddressCodecs() ProtocolAddressCodecRegistry {
	return ProtocolAddressCodecRegistry{registry: maps.Clone(registeredProtocolAddressCodecs)}
}

// ProtocolAddressCodecRegistry is the 2-method registry used by the bootstrap path.
type ProtocolAddressCodecRegistry struct {
	registry map[string]ProtocolAddressCodec
}

// OracleIDAsAddressBytes returns valid address bytes for a chain selector and oracle ID.
func (r ProtocolAddressCodecRegistry) OracleIDAsAddressBytes(oracleID uint8, chainSelector cciptypes.ChainSelector) ([]byte, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}
	codec, ok := r.registry[family]
	if !ok {
		return nil, fmt.Errorf("no ProtocolAddressCodec registered for family %s", family)
	}
	return codec.OracleIDAsAddressBytes(oracleID)
}

// TransmitterBytesToString converts a transmitter account from bytes to string.
func (r ProtocolAddressCodecRegistry) TransmitterBytesToString(addr cciptypes.UnknownAddress, chainSelector cciptypes.ChainSelector) (string, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return "", fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}
	codec, ok := r.registry[family]
	if !ok {
		return "", fmt.Errorf("no ProtocolAddressCodec registered for family %s", family)
	}
	return codec.TransmitterBytesToString(addr)
}
