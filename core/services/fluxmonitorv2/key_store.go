package fluxmonitorv2

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	corestore "github.com/smartcontractkit/chainlink/core/store"
)

//go:generate mockery --name KeyStoreInterface --output ./mocks/ --case=underscore

// KeyStoreInterface defines an interface to interact with the keystore
type KeyStoreInterface interface {
	Accounts() []accounts.Account
	GetRoundRobinAddress() (common.Address, error)
}

// KeyStore implements KeyStoreInterface
type KeyStore struct {
	store *corestore.Store
}

// NewKeyStore initializes a new keystore
func NewKeyStore(store *corestore.Store) *KeyStore {
	return &KeyStore{store: store}
}

// Accounts gets the node's accounts from the keystore
func (ks *KeyStore) Accounts() []accounts.Account {
	return ks.store.KeyStore.Accounts()
}

// GetRoundRobinAddress queries the database for the address of a random
// ethereum key derived from the id.
func (ks *KeyStore) GetRoundRobinAddress() (common.Address, error) {
	return ks.store.GetRoundRobinAddress()
}
