package common

import (
	"fmt"
	"sync"

	chainsel "github.com/smartcontractkit/chain-selectors"
	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// Ensure AddressCodec implements the gRPC-compatible AddressCodec interface from chainlink-common.
var _ ccipocr3common.AddressCodec = &AddressCodec{}

// AddressCodec is a concurrency-safe registry that manages gRPC-compatible
// ChainSpecificAddressCodec instances for different chain families. It implements
// the AddressCodec interface from chainlink-common by delegating to the
// registered chain family specific codec.
//
// Terminology:
//   - "AddressCodec": refers to the registry instance. It maintains the map of
//     chain family to codec and provides thread-safe access to it.
//   - "AddressCodec interface": is the chain-agnostic interface defined in
//     chainlink-common that can be called over gRPC.
//   - "ChainSpecificAddressCodec": is the chain family specific interface that
//     matches the gRPC-compatible interface defined in
//     chainlink-common/pkg/types/ccipocr3/plugincodec.go and returned by
//     CCIPProvider.Codec().
//
// AddressCodec values constructed with NewAddressCodec may be copied; copies
// share the same underlying registry.
type AddressCodec struct {
	registry *addressCodecRegistry
}

type addressCodecRegistry struct {
	registeredAddressCodecMap map[string]ChainSpecificAddressCodec
	mu                        sync.RWMutex
}

// NewAddressCodec constructs an AddressCodec registry.
func NewAddressCodec(registeredMap map[string]ChainSpecificAddressCodec) AddressCodec {
	return AddressCodec{
		registry: &addressCodecRegistry{
			registeredAddressCodecMap: cloneAddressCodecMap(registeredMap),
		},
	}
}

// RegisterCodec registers a chain family specific address codec.
func (ac *AddressCodec) RegisterCodec(chainFamily string, codec ChainSpecificAddressCodec) {
	registry := ac.getRegistry()
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if registry.registeredAddressCodecMap == nil {
		registry.registeredAddressCodecMap = make(map[string]ChainSpecificAddressCodec)
	}
	registry.registeredAddressCodecMap[chainFamily] = codec
}

// HasCodec returns true if a chain family specific address codec is registered.
func (ac *AddressCodec) HasCodec(chainFamily string) bool {
	registry := ac.registry
	if registry == nil {
		return false
	}

	registry.mu.RLock()
	defer registry.mu.RUnlock()

	_, ok := registry.registeredAddressCodecMap[chainFamily]
	return ok
}

// ============ gRPC-compatible implementation of AddressCodec interface ============

// AddressBytesToString converts an address from bytes to string
func (ac *AddressCodec) AddressBytesToString(addr ccipocr3common.UnknownAddress, chainSelector ccipocr3common.ChainSelector) (string, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return "", fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}

	codec, err := ac.getCodec(family, "address decode")
	if err != nil {
		return "", err
	}

	return codec.AddressBytesToString(addr)
}

// AddressStringToBytes converts an address from string to bytes
func (ac *AddressCodec) AddressStringToBytes(addr string, chainSelector ccipocr3common.ChainSelector) (ccipocr3common.UnknownAddress, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}
	codec, err := ac.getCodec(family, "address decode")
	if err != nil {
		return nil, err
	}

	return codec.AddressStringToBytes(addr)
}

// ============ Superset of ChainSpecificAddressCodec for core OCR config handling ============

// TransmitterBytesToString converts a transmitter account from bytes to string.
func (ac *AddressCodec) TransmitterBytesToString(addr ccipocr3common.UnknownAddress, chainSelector ccipocr3common.ChainSelector) (string, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return "", fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}

	codec, err := ac.getCodec(family, "transmitter decode")
	if err != nil {
		return "", err
	}

	return codec.TransmitterBytesToString(addr)
}

// OracleIDAsAddressBytes returns valid address bytes for a given chain selector and oracle ID.
// Used for making nil transmitters in the OCR config valid, it just means that this oracle does not support the destination chain.
func (ac *AddressCodec) OracleIDAsAddressBytes(oracleID uint8, chainSelector ccipocr3common.ChainSelector) ([]byte, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}
	codec, err := ac.getCodec(family, "address decode")
	if err != nil {
		return nil, err
	}

	return codec.OracleIDAsAddressBytes(oracleID)
}

func (ac *AddressCodec) getRegistry() *addressCodecRegistry {
	if ac.registry == nil {
		ac.registry = &addressCodecRegistry{}
	}
	return ac.registry
}

func (ac *AddressCodec) getCodec(family string, decodeType string) (ChainSpecificAddressCodec, error) {
	registry := ac.registry
	if registry == nil {
		return nil, fmt.Errorf("unsupported family for %s type %s", decodeType, family)
	}

	registry.mu.RLock()
	defer registry.mu.RUnlock()

	codec, exist := registry.registeredAddressCodecMap[family]
	if !exist {
		return nil, fmt.Errorf("unsupported family for %s type %s", decodeType, family)
	}

	return codec, nil
}

func cloneAddressCodecMap(registeredMap map[string]ChainSpecificAddressCodec) map[string]ChainSpecificAddressCodec {
	if registeredMap == nil {
		return nil
	}

	clone := make(map[string]ChainSpecificAddressCodec, len(registeredMap))
	for family, codec := range registeredMap {
		clone[family] = codec
	}
	return clone
}
