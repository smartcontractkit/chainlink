package common

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// AddressCodec is a struct that holds the chain specific address codecs
type AddressCodec struct {
	EVMAddressCodec    ChainSpecificAddressCodec
	SolanaAddressCodec ChainSpecificAddressCodec
	AptosAddressCodec  ChainSpecificAddressCodec
}

// NewAddressCodec is a constructor for NewAddressCodec
func NewAddressCodec(evmAddrCodec, solanaAddrCodec, aptosAddressCodec ChainSpecificAddressCodec) AddressCodec {
	return AddressCodec{
		EVMAddressCodec:    evmAddrCodec,
		SolanaAddressCodec: solanaAddrCodec,
		AptosAddressCodec:  aptosAddressCodec,
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

	case chainsel.FamilyAptos:
		return ac.AptosAddressCodec.AddressBytesToString(addr)

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

	case chainsel.FamilyAptos:
		return ac.AptosAddressCodec.AddressStringToBytes(addr)

	default:
		return nil, fmt.Errorf("unsupported family for address decode type %s", family)
	}
}
