package store

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/multierr"
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
	var merr error
	for _, account := range ks.Accounts() {
		err := ks.KeyStore.Unlock(account, phrase)
		if err != nil {
			merr = multierr.Combine(merr, fmt.Errorf("invalid password for account %s", account.Address.Hex()), err)
		} else {
			logger.Infow(fmt.Sprint("Unlocked account ", account.Address.Hex()), "address", account.Address.Hex())
		}
	}
	return merr
}

// NewAccount adds an account to the keystore
func (ks *KeyStore) NewAccount(passphrase string) (accounts.Account, error) {
	account, err := ks.KeyStore.NewAccount(passphrase)
	if err != nil {
		return accounts.Account{}, err
	}

	err = ks.KeyStore.Unlock(account, passphrase)
	if err != nil {
		return accounts.Account{}, err
	}

	return account, nil
}

// SignTx uses the unlocked account to sign the given transaction.
func (ks *KeyStore) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return ks.KeyStore.SignTx(account, tx, chainID)
}

// Sign creates an HMAC from some input data using the account's private key
func (ks *KeyStore) Sign(input []byte) (models.Signature, error) {
	account, err := ks.GetFirstAccount()
	if err != nil {
		return models.Signature{}, err
	}
	hash, err := utils.Keccak256(input)
	if err != nil {
		return models.Signature{}, err
	}

	output, err := ks.KeyStore.SignHash(account, hash)
	if err != nil {
		return models.Signature{}, err
	}
	var signature models.Signature
	signature.SetBytes(output)
	return signature, nil
}

// GetFirstAccount returns the unlocked account in the KeyStore object. The client
// ensures that an account exists during authentication.
func (ks *KeyStore) GetFirstAccount() (accounts.Account, error) {
	if len(ks.Accounts()) == 0 {
		return accounts.Account{}, errors.New("No Ethereum Accounts configured")
	}
	return ks.Accounts()[0], nil
}

// GetAccounts returns all accounts
func (ks *KeyStore) GetAccounts() []accounts.Account {
	return ks.Accounts()
}
