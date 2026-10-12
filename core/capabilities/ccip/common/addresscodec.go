package common

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

var _ ccipocr3.AddressCodec = &AddressCodec{}

// AddressCodec is a struct that holds the chain specific address codecs and
// implements a superset of the ccipocr3.AddressCodec interface.
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
func (ac AddressCodec) AddressBytesToString(addr ccipocr3.UnknownAddress, chainSelector ccipocr3.ChainSelector) (string, error) {
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

// TransmitterBytesToString converts a transmitter account from bytes to string
func (ac AddressCodec) TransmitterBytesToString(addr ccipocr3.UnknownAddress, chainSelector ccipocr3.ChainSelector) (string, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return "", fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}

	codec, exist := ac.registeredAddressCodecMap[family]
	if !exist {
		return "", fmt.Errorf("unsupported family for transmitter decode type %s", family)
	}

	return codec.TransmitterBytesToString(addr)
}

// AddressStringToBytes converts an address from string to bytes
func (ac AddressCodec) AddressStringToBytes(addr string, chainSelector ccipocr3.ChainSelector) (ccipocr3.UnknownAddress, error) {
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

// OracleIDAsAddressBytes returns valid address bytes for a given chain selector and oracle ID.
// Used for making nil transmitters in the OCR config valid, it just means that this oracle does not support the destination chain.
func (ac AddressCodec) OracleIDAsAddressBytes(oracleID uint8, chainSelector ccipocr3.ChainSelector) ([]byte, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}
	codec, exist := ac.registeredAddressCodecMap[family]
	if !exist {
		return nil, fmt.Errorf("unsupported family for address decode type %s", family)
	}

	return codec.OracleIDAsAddressBytes(oracleID)
}
