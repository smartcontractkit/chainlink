package store

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func (a *ManagedAccount) PublicLastConfirmedNonce() uint64 {
	return a.lastConfirmedNonce
}

func (txm *EthTxManager) GetAvailableAccount(from common.Address) *ManagedAccount {
	for _, a := range txm.availableAccounts {
		if a.Address == from {
			return a
		}
	}
	return nil
}

func TestManagedAccount_updateLastConfirmedNonce(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		last      uint64
		submitted uint64
		want      uint64
	}{
		{"greater", 100, 101, 101},
		{"less", 100, 99, 100},
		{"equal", 100, 100, 100},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ma := &ManagedAccount{lastConfirmedNonce: test.last}
			ma.updateLastConfirmedNonce(test.submitted)

			assert.Equal(t, test.want, ma.lastConfirmedNonce)
		})
	}
}

func TestTxManager_updateLastConfirmedTx_success(t *testing.T) {
	t.Parallel()

	from := common.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca754")
	nonce := uint64(234)
	a := accounts.Account{Address: from}
	ma := &ManagedAccount{Account: a}
	txm := &EthTxManager{availableAccounts: []*ManagedAccount{ma}}
	tx := &models.Tx{From: from, Nonce: nonce}

	txm.updateManagedAccounts(tx)

	assert.Equal(t, nonce, ma.lastConfirmedNonce)
}

func TestTxManager_updateLastConfirmedTx_noMatchingAccount(t *testing.T) {
	t.Parallel()

	from := common.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca754")
	nonce := uint64(234)
	a := accounts.Account{Address: common.HexToAddress("0xffffffffffff666546e30d74d50d173d20bca754")}
	ma := &ManagedAccount{Account: a}
	txm := &EthTxManager{availableAccounts: []*ManagedAccount{ma}}
	tx := &models.Tx{From: from, Nonce: nonce}

	txm.updateManagedAccounts(tx)

	assert.NotEqual(t, nonce, ma.lastConfirmedNonce)
}
