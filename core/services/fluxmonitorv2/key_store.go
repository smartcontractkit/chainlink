package fluxmonitorv2

import (
	"github.com/ethereum/go-ethereum/common"
	corestore "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name KeyStoreInterface --output ./mocks/ --case=underscore

// KeyStoreInterface defines an interface to interact with the keystore
type KeyStoreInterface interface {
	SendingKeys() ([]models.Key, error)
	GetRoundRobinAddress(...common.Address) (common.Address, error)
}

// KeyStore implements KeyStoreInterface
type KeyStore struct {
	corestore.KeyStoreInterface
}

// NewKeyStore initializes a new keystore
func NewKeyStore(ks corestore.KeyStoreInterface) *KeyStore {
	return &KeyStore{ks}
}
