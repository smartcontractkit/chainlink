package testutils

import (
	"context"
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ethkey"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// SigningKeystore is a mock keystore that actually signs messages
type SigningKeystore struct {
	addressToKey map[common.Address]*ecdsa.PrivateKey
	addresses    []common.Address
}

// NewSigningKeystore creates a new SigningKeystore from a map of address strings to private keys
// This allows callers to avoid importing go-ethereum/common
func NewSigningKeystore(addressToKey map[string]*ecdsa.PrivateKey, addresses []ethkey.KeyV2) *SigningKeystore {
	addrToKey := make(map[common.Address]*ecdsa.PrivateKey)
	addrList := make([]common.Address, 0, len(addresses))

	for addrStr, key := range addressToKey {
		addr := common.HexToAddress(addrStr)
		addrToKey[addr] = key
	}

	for _, key := range addresses {
		addrList = append(addrList, key.Address)
	}

	return &SigningKeystore{
		addressToKey: addrToKey,
		addresses:    addrList,
	}
}

func (s *SigningKeystore) CheckEnabled(ctx context.Context, address common.Address) error {
	for _, addr := range s.addresses {
		if addr == address {
			return nil
		}
	}
	return errors.New("not enabled")
}

func (s *SigningKeystore) EnabledAddresses(ctx context.Context) ([]common.Address, error) {
	return s.addresses, nil
}

func (s *SigningKeystore) SignMessage(ctx context.Context, address common.Address, message []byte) ([]byte, error) {
	key, ok := s.addressToKey[address]
	if !ok {
		return nil, errors.New("address not found")
	}
	// SignMessage uses accounts.TextHash which applies Ethereum message prefix
	// Use utils.GenerateEthSignature which does the same thing
	return utils.GenerateEthSignature(key, message)
}

func (s *SigningKeystore) Sign(ctx context.Context, address common.Address, bytes []byte) ([]byte, error) {
	key, ok := s.addressToKey[address]
	if !ok {
		return nil, errors.New("address not found")
	}
	// Sign signs raw bytes without prefix - use utils.GenerateEthSignature which applies prefix
	// This is not perfect but works for our test since we only test SignMessage
	return utils.GenerateEthSignature(key, bytes)
}

func (s *SigningKeystore) GetNextAddress(ctx context.Context, addresses ...common.Address) (common.Address, error) {
	if len(s.addresses) == 0 {
		return common.Address{}, errors.New("no addresses available")
	}
	return s.addresses[0], nil
}

func (s *SigningKeystore) GetMutex(address common.Address) *keys.Mutex {
	return &keys.Mutex{}
}
