package store

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
)

// KeyStore manages a key storage directory on disk.
type KeyStore struct {
	*keystore.KeyStore
}

// NewKeyStore creates a keystore for the given directory.
func NewKeyStore(keyDir string) *KeyStore {
	ks := keystore.NewKeyStore(
		keyDir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	return &KeyStore{ks}
}

// HasAccounts returns true if there are accounts located at the keystore
// directory.
func (ks *KeyStore) HasAccounts() bool {
	return len(ks.Accounts()) > 0
}

// Unlock uses the given password to try to unlock accounts located in the
// keystore directory.
func (ks *KeyStore) Unlock(phrase string) error {
	for _, account := range ks.Accounts() {
		err := ks.KeyStore.Unlock(account, phrase)
		if err != nil {
			return fmt.Errorf("Invalid password for account: %s\n\nPlease try again...\n", account.Address.Hex())
		}
	}
	return nil
}

// SignTx uses the unlocked account to sign the given transaction.
func (ks *KeyStore) SignTx(tx *types.Transaction, chainID uint64) (*types.Transaction, error) {
	account, err := ks.GetAccount()
	if err != nil {
		return nil, err
	}

	return ks.KeyStore.SignTx(
		account,
		tx, big.NewInt(int64(chainID)),
	)
}

// GetAccount returns the unlocked account in the KeyStore object. The client
// ensures that an account exists during authentication.
func (ks *KeyStore) GetAccount() (accounts.Account, error) {
	if len(ks.Accounts()) == 0 {
		return accounts.Account{}, errors.New("No Ethereum Accounts configured")
	}
	return ks.Accounts()[0], nil
}
