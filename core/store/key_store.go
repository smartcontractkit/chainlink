package store

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/multierr"
)

// EthereumMessageHashPrefix is a Geth-originating message prefix that seeks to
// prevent arbitrary message data to be representable as a valid Ethereum transaction
// For more information, see: https://github.com/ethereum/go-ethereum/issues/3731
const EthereumMessageHashPrefix = "\x19Ethereum Signed Message:\n32"

//go:generate mockery --name KeyStoreInterface --output ../internal/mocks/ --case=underscore
type KeyStoreInterface interface {
	Accounts() []accounts.Account
	Wallets() []accounts.Wallet
	GetFirstAccount() (accounts.Account, error)
	HasAccounts() bool
	Unlock(phrase string) error
	NewAccount(passphrase string) (accounts.Account, error)
	SignHash(hash common.Hash) (models.Signature, error)
	Import(keyJSON []byte, passphrase, newPassphrase string) (accounts.Account, error)
	Export(a accounts.Account, passphrase, newPassphrase string) ([]byte, error)
	GetAccounts() []accounts.Account
	GetAccountByAddress(common.Address) (accounts.Account, error)

	SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
}

// KeyStore manages a key storage directory on disk.
type KeyStore struct {
	*keystore.KeyStore
	scryptParams utils.ScryptParams
}

// NewKeyStore creates a keystore for the given directory.
func NewKeyStore(keyDir string, scryptParams utils.ScryptParams) *KeyStore {
	ks := keystore.NewKeyStore(keyDir, scryptParams.N, scryptParams.P)
	return &KeyStore{ks, scryptParams}
}

// NewInsecureKeyStore creates an *INSECURE* keystore for the given directory.
// NOTE: Should only be used for testing!
func NewInsecureKeyStore(keyDir string) *KeyStore {
	return NewKeyStore(keyDir, utils.FastScryptParams)
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

// SignHash signs a precomputed digest, using the first account's private key
// This method adds an ethereum message prefix to the message before signing it,
// invalidating any would-be valid Ethereum transactions
func (ks *KeyStore) SignHash(hash common.Hash) (models.Signature, error) {
	prefixedMessageBytes, err := utils.Keccak256(append([]byte(EthereumMessageHashPrefix), hash.Bytes()...))
	if err != nil {
		return models.Signature{}, err
	}

	signature, err := ks.unsafeSignHash(common.BytesToHash(prefixedMessageBytes))
	if err != nil {
		return models.Signature{}, err
	}

	return signature, nil
}

// unsafeSignHash signs a precomputed digest, using the first account's private
// key
// NOTE: Do not use this method to sign arbitrary message hashes, it may be an
// Ethereum transaction in disguise! Use SignHashSafe instead unless this is
// strictly needed
func (ks *KeyStore) unsafeSignHash(hash common.Hash) (models.Signature, error) {
	account, err := ks.GetFirstAccount()
	if err != nil {
		return models.Signature{}, err
	}
	output, err := ks.KeyStore.SignHash(account, hash.Bytes())
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
		return accounts.Account{}, errors.New("no Ethereum Accounts configured")
	}
	return ks.Accounts()[0], nil
}

// GetAccounts returns all accounts
func (ks *KeyStore) GetAccounts() []accounts.Account {
	return ks.Accounts()
}

// GetAccountByAddress returns the account matching the address provided, or an error if it is missing
func (ks *KeyStore) GetAccountByAddress(address common.Address) (accounts.Account, error) {
	for _, account := range ks.Accounts() {
		if account.Address == address {
			return account, nil
		}
	}
	return accounts.Account{}, errors.New("no account found with that address")
}
