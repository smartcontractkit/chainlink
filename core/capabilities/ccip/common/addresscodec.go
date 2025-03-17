package common

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// ChainSpecificAddressCodec is an interface that defines the methods for encoding and decoding addresses
type ChainSpecificAddressCodec interface {
	AddressBytesToString([]byte) (string, error)
	AddressStringToBytes(string) ([]byte, error)
}

// AddressCodec is a struct that holds the chain specific address codecs
type AddressCodec struct {
	EVMAddressCodec    ChainSpecificAddressCodec
	SolanaAddressCodec ChainSpecificAddressCodec
}

// AddressCodecParams is a struct that holds the parameters for creating a AddressCodec
type AddressCodecParams struct {
	evmAddressCodec    ChainSpecificAddressCodec
	solanaAddressCodec ChainSpecificAddressCodec
}

// NewAddressCodecParams is a constructor for AddressCodecParams
func NewAddressCodecParams(evmAddressCodec ChainSpecificAddressCodec, solanaAddressCodec ChainSpecificAddressCodec) AddressCodecParams {
	return AddressCodecParams{
		evmAddressCodec:    evmAddressCodec,
		solanaAddressCodec: solanaAddressCodec,
	}
}

// NewAddressCodec is a constructor for AddressCodec
func NewAddressCodec(params AddressCodecParams) AddressCodec {
	return AddressCodec{
		EVMAddressCodec:    params.evmAddressCodec,
		SolanaAddressCodec: params.solanaAddressCodec,
	}
}

// AddressBytesToString converts an address from bytes to string
func (ac AddressCodec) AddressBytesToString(addr cciptypes.UnknownAddress, chainSelector cciptypes.ChainSelector) (string, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return "", fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}

	switch family {
	case chainsel.FamilyEVM:
		return ac.EVMAddressCodec.AddressBytesToString(addr)

	case chainsel.FamilySolana:
		return ac.SolanaAddressCodec.AddressBytesToString(addr)

	default:
		return "", fmt.Errorf("unsupported family for address encode type %s", family)
	}
}

// AddressStringToBytes converts an address from string to bytes
func (ac AddressCodec) AddressStringToBytes(addr string, chainSelector cciptypes.ChainSelector) (cciptypes.UnknownAddress, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}

	switch family {
	case chainsel.FamilyEVM:
		return ac.EVMAddressCodec.AddressStringToBytes(addr)

	case chainsel.FamilySolana:
		return ac.SolanaAddressCodec.AddressStringToBytes(addr)

	default:
		return nil, fmt.Errorf("unsupported family for address decode type %s", family)
	}
}
