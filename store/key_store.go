package store

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
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

func (self *KeyStore) Unlock(phrase string) error {
	for _, account := range self.Accounts() {
		err := self.KeyStore.Unlock(account, phrase)
		if err != nil {
			return fmt.Errorf("Invalid password for account: %s\n\nPlease try again...\n", account.Address.Hex())
		}
	}
	return nil
}

func (self *KeyStore) SignTx(tx *types.Transaction, chainID int64) (*types.Transaction, error) {
	return self.KeyStore.SignTx(
		self.GetAccount().Account,
		tx, big.NewInt(chainID),
	)
}

func (self *KeyStore) GetAccount() *Account {
	return &Account{self.Accounts()[0]}
}

type Account struct {
	accounts.Account
}

func (self *Account) GetNonce(config Config) (uint64, error) {
	eth, err := rpc.Dial(config.EthereumURL)
	if err != nil {
		return 0, err
	}
	var result string
	err = eth.Call(&result, "eth_getTransactionCount", self.Address.Hex())
	if err != nil {
		return 0, err
	}
	if strings.ToLower(result[0:2]) == "0x" {
		result = result[2:]
	}
	return strconv.ParseUint(result, 16, 64)
}
