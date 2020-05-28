package models_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHead_NewHead(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input *big.Int
		want  string
	}{
		{big.NewInt(0), "0"},
		{big.NewInt(0xf), "f"},
		{big.NewInt(0x10), "10"},
	}
	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			num := models.NewHead(test.input, cltest.NewHash(), cltest.NewHash(), big.NewInt(0))
			assert.Equal(t, test.want, fmt.Sprintf("%x", num.ToInt()))
		})
	}
}

func TestHead_GreaterThan(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		left    *models.Head
		right   *models.Head
		greater bool
	}{
		{"nil nil", nil, nil, false},
		{"present nil", cltest.Head(1), nil, true},
		{"nil present", nil, cltest.Head(1), false},
		{"less", cltest.Head(1), cltest.Head(2), false},
		{"equal", cltest.Head(2), cltest.Head(2), false},
		{"greater", cltest.Head(2), cltest.Head(1), true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.greater, test.left.GreaterThan(test.right))
		})
	}
}

func TestHead_NextInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		bn   *models.Head
		want *big.Int
	}{
		{"nil", nil, nil},
		{"one", cltest.Head(1), big.NewInt(2)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.bn.NextInt())
		})
	}
}

func TestTx_PresenterMatchesHex(t *testing.T) {
	t.Parallel()

	createdTx := models.Tx{
		From:        common.HexToAddress("0xf208"),
		To:          common.HexToAddress("0x70"),
		Data:        []byte(`{"data": "is wilding out"}`),
		Nonce:       0x8008,
		Value:       utils.NewBig(big.NewInt(777)),
		GasLimit:    1999,
		Hash:        common.HexToHash("0x0"),
		GasPrice:    utils.NewBig(big.NewInt(333)),
		Confirmed:   true,
		SentAt:      1745,
		SignedRawTx: hexutil.MustDecode("0xcafe"),
	}

	ptx := presenters.NewTx(&createdTx)
	bytes, err := json.Marshal(ptx)
	require.NoError(t, err)
	assert.JSONEq(t, `{`+
		`"confirmed":true,`+
		`"data":"0x7b2264617461223a202269732077696c64696e67206f7574227d",`+
		`"from":"0x000000000000000000000000000000000000f208",`+
		`"gasLimit":"1999",`+
		`"gasPrice":"333",`+
		`"hash":"0x0000000000000000000000000000000000000000000000000000000000000000",`+
		`"rawHex":"0xcafe",`+
		`"nonce":"32776",`+
		`"sentAt":"1745",`+
		`"to":"0x0000000000000000000000000000000000000070",`+
		`"value":"777"`+
		`}`, string(bytes))
}

func TestHighestPricedTxAttemptPerTx(t *testing.T) {
	items := []models.TxAttempt{
		{TxID: 1, GasPrice: utils.NewBig(big.NewInt(5555))},
		{TxID: 1, GasPrice: utils.NewBig(big.NewInt(444))},
		{TxID: 1, GasPrice: utils.NewBig(big.NewInt(2))},
		{TxID: 1, GasPrice: utils.NewBig(big.NewInt(33333))},
		{TxID: 2, GasPrice: utils.NewBig(big.NewInt(4444))},
		{TxID: 2, GasPrice: utils.NewBig(big.NewInt(999))},
		{TxID: 2, GasPrice: utils.NewBig(big.NewInt(12211))},
	}

	items = models.HighestPricedTxAttemptPerTx(items)

	sort.Slice(items, func(i, j int) bool { return items[i].TxID < items[j].TxID })

	assert.Len(t, items, 2)
	assert.True(t, items[0].GasPrice.ToInt().Cmp(big.NewInt(33333)) == 0)
	assert.True(t, items[1].GasPrice.ToInt().Cmp(big.NewInt(12211)) == 0)
}
