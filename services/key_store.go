package services

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type KeyStore struct {
	*keystore.KeyStore
}

func NewKeyStore(keyDir string) *KeyStore {
	ks := keystore.NewKeyStore(
		keyDir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	return &KeyStore{ks}
}

func (self *KeyStore) HasAccounts() bool {
	return len(self.Accounts()) > 0
}
