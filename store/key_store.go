package store

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
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

func (ks *KeyStore) HasAccounts() bool {
	return len(ks.Accounts()) > 0
}

func (ks *KeyStore) Unlock(phrase string) error {
	for _, account := range ks.Accounts() {
		err := ks.KeyStore.Unlock(account, phrase)
		if err != nil {
			return fmt.Errorf("Invalid password for account: %s\n\nPlease try again...\n", account.Address.Hex())
		}
	}
	return nil
}

func (ks *KeyStore) SignTx(tx *types.Transaction, chainID uint64) (*types.Transaction, error) {
	return ks.KeyStore.SignTx(
		ks.GetAccount(),
		tx, big.NewInt(int64(chainID)),
	)
}

func (ks *KeyStore) GetAccount() accounts.Account {
	return ks.Accounts()[0]
}
