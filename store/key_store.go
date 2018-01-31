package store

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	NodeWarningString       = "WARNING: Chainlink node may not be fully functional. "
	NoETHConnectivityString = "Unable to query the provided ETH wallet's balance. " +
		"Are you connected to the Ethereum network? " + NodeWarningString
	MissingWalletString  = "No ETH wallet found. " + NodeWarningString
	ZeroETHBalanceString = "Zero balance in the provided ETH wallet. " + NodeWarningString
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
	return ks.KeyStore.SignTx(
		ks.GetAccount(),
		tx, big.NewInt(int64(chainID)),
	)
}

// GetAccount returns the unlocked account in the KeyStore object.
func (ks *KeyStore) GetAccount() accounts.Account {
	return ks.Accounts()[0]
}

func (ks *KeyStore) ShowEthBalance(txm *TxManager) string {
	result := ""
	if ks.HasAccounts() {
		account := ks.GetAccount()
		balance, err := txm.GetEthBalance(account.Address)
		if err != nil {
			result = NoETHConnectivityString + err.Error()
		} else {
			address := account.Address.Hex()
			result += fmt.Sprintf("ETH Balance for %v: %v. ", address, balance)
			if balance == 0 {
				result += ZeroETHBalanceString
			}
		}
	} else {
		result += MissingWalletString
	}
	return result
}
