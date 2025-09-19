package common

import (
	"sync"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// Ensure AddressCodecRegistry implements the AddressCodecBundle interface
var _ cciptypes.AddressCodecBundle = (*AddressCodecRegistry)(nil)

// AddressCodecRegistry is a singleton registry that manages ChainSpecificAddressCodec instances
// for different chain families. It implements the AddressCodecBundle interface
// by delegating to the existing AddressCodec implementation.
//
// Terminology:
//   - "AddressCodecRegistry": refers to the entire singleton registry instance. It both maintains the map of
//     chain family to codec and provides thread-safe access to it.
//   - "AddressCodecBundle": is the interface that the registry implements
type AddressCodecRegistry struct {
	addressCodec cciptypes.AddressCodecMap
	mu           sync.RWMutex
}

var (
	addressRegistryInstance *AddressCodecRegistry
	addressRegistryOnce     sync.Once
)

// GetAddressCodecRegistry returns the singleton instance of AddressCodecRegistry. This is only called
// in core node.
func GetAddressCodecRegistry() *AddressCodecRegistry {
	addressRegistryOnce.Do(func() {
		addressRegistryInstance = &AddressCodecRegistry{
			addressCodec: make(cciptypes.AddressCodecMap),
		}
	})
	return addressRegistryInstance
}

// RegisterFamily registers a chain family with a no-op ChainSpecificAddressCodec if not already registered.
// This is used when we know which chain families we want to support but don't have a specific codec
// implementation initialized for it yet.
func (r *AddressCodecRegistry) RegisterFamily(chainFamily string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.addressCodec[chainFamily]; !exists {
		r.addressCodec[chainFamily] = NoOpChainSpecificAddressCodec{}
	}
}

// RegisterCodec registers a ChainSpecificAddressCodec for a specific chain family
func (r *AddressCodecRegistry) RegisterCodec(chainFamily string, codec ChainSpecificAddressCodec) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.addressCodec[chainFamily] = codec
}

// ============ Implementation of AddressCodecBundle interface ============

func (r *AddressCodecRegistry) AddressBytesToString(addr cciptypes.UnknownAddress, chainSelector cciptypes.ChainSelector) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.addressCodec.AddressBytesToString(addr, chainSelector)
}

func (r *AddressCodecRegistry) TransmitterBytesToString(addr cciptypes.UnknownAddress, chainSelector cciptypes.ChainSelector) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.addressCodec.TransmitterBytesToString(addr, chainSelector)
}

func (r *AddressCodecRegistry) AddressStringToBytes(addr string, chainSelector cciptypes.ChainSelector) (cciptypes.UnknownAddress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.addressCodec.AddressStringToBytes(addr, chainSelector)
}

func (r *AddressCodecRegistry) OracleIDAsAddressBytes(oracleID uint8, chainSelector cciptypes.ChainSelector) ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.addressCodec.OracleIDAsAddressBytes(oracleID, chainSelector)
}

type NoOpChainSpecificAddressCodec struct{}

func (n NoOpChainSpecificAddressCodec) AddressBytesToString(addr []byte) (string, error) {
	return string(addr), nil
}

func (n NoOpChainSpecificAddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	return []byte(addr), nil
}

func (n NoOpChainSpecificAddressCodec) OracleIDAsAddressBytes(oracleID uint8) ([]byte, error) {
	return []byte{oracleID}, nil
}

func (n NoOpChainSpecificAddressCodec) TransmitterBytesToString(addr []byte) (string, error) {
	return string(addr), nil
}
