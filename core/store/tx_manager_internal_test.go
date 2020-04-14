package store

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func (a *ManagedAccount) PublicLastSafeNonce() uint64 {
	return a.lastSafeNonce
}

func (a *ManagedAccount) SetLastSafeNonce(n uint64) {
	a.lastSafeNonce = n
}

func TestManagedAccount_updateLastSafeNonce(t *testing.T) {
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
			ma := &ManagedAccount{lastSafeNonce: test.last}
			ma.updateLastSafeNonce(test.submitted)

			assert.Equal(t, test.want, ma.lastSafeNonce)
		})
	}
}

func TestTxManager_updateLastSafeNonce_success(t *testing.T) {
	t.Parallel()

	from := common.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca754")
	nonce := uint64(234)
	a := accounts.Account{Address: from}
	ma := &ManagedAccount{Account: a}
	txm := &EthTxManager{availableAccounts: []*ManagedAccount{ma}}
	tx := &models.Tx{From: from, Nonce: nonce}

	txm.updateLastSafeNonce(tx)

	assert.Equal(t, nonce, ma.lastSafeNonce)
}

func TestTxManager_updateLastSafeNonce_noMatchingAccount(t *testing.T) {
	t.Parallel()

	from := common.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca754")
	nonce := uint64(234)
	a := accounts.Account{Address: common.HexToAddress("0xffffffffffff666546e30d74d50d173d20bca754")}
	ma := &ManagedAccount{Account: a}
	txm := &EthTxManager{availableAccounts: []*ManagedAccount{ma}}
	tx := &models.Tx{From: from, Nonce: nonce}

	txm.updateLastSafeNonce(tx)

	assert.NotEqual(t, nonce, ma.lastSafeNonce)
}
