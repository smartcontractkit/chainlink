package fluxmonitorv2

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

//go:generate mockery --quiet --name KeyStoreInterface --output ./mocks/ --case=underscore

// KeyStoreInterface defines an interface to interact with the keystore
type KeyStoreInterface interface {
	EnabledKeysForChain(chainID *big.Int) ([]ethkey.KeyV2, error)
	GetRoundRobinAddress(chainID *big.Int, addrs ...common.Address) (common.Address, error)
}

// KeyStore implements KeyStoreInterface
type KeyStore struct {
	keystore.Eth
}

// NewKeyStore initializes a new keystore
func NewKeyStore(ks keystore.Eth) *KeyStore {
	return &KeyStore{ks}
}
