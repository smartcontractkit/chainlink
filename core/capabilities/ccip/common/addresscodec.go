package common

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// AddressCodec is a struct that holds the chain specific address codecs
type AddressCodec struct {
	registeredAddressCodecMap map[string]ChainSpecificAddressCodec
}

// NewAddressCodec is a constructor for NewAddressCodec
func NewAddressCodec(registeredMap map[string]ChainSpecificAddressCodec) AddressCodec {
	return AddressCodec{
		registeredAddressCodecMap: registeredMap,
	}
}

// AddressBytesToString converts an address from bytes to string
func (ac AddressCodec) AddressBytesToString(addr cciptypes.UnknownAddress, chainSelector cciptypes.ChainSelector) (string, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return "", fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}

	codec, exist := ac.registeredAddressCodecMap[family]
	if !exist {
		return "", fmt.Errorf("unsupported family for address decode type %s", family)
	}

	return codec.AddressBytesToString(addr)
}

// AddressStringToBytes converts an address from string to bytes
func (ac AddressCodec) AddressStringToBytes(addr string, chainSelector cciptypes.ChainSelector) (cciptypes.UnknownAddress, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}
	codec, exist := ac.registeredAddressCodecMap[family]
	if !exist {
		return nil, fmt.Errorf("unsupported family for address decode type %s", family)
	}

	return codec.AddressStringToBytes(addr)
}
