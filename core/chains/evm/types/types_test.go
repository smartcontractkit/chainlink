package types_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	testGethLog1 = &gethTypes.Log{
		Address: common.HexToAddress("0x11111111"),
		Topics: []common.Hash{
			common.HexToHash("0xaaaaaaaa"),
			common.HexToHash("0xbbbbbbbb"),
		},
		Data:        []byte{1, 2, 3, 4, 5},
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0xffffffff"),
		TxHash:      common.HexToHash("0xcccccccc"),
		TxIndex:     100,
		Index:       200,
		Removed:     false,
	}

	testGethLog2 = &gethTypes.Log{
		Address: common.HexToAddress("0x11111112"),
		Topics: []common.Hash{
			common.HexToHash("0xaaaaaaab"),
			common.HexToHash("0xbbbbbbbc"),
		},
		Data:        []byte{2, 3, 4, 5, 6},
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0xfffffff0"),
		TxHash:      common.HexToHash("0xcccccccd"),
		TxIndex:     101,
		Index:       201,
		Removed:     true,
	}

	testGethReceipt = &gethTypes.Receipt{
		PostState:         []byte{1, 2, 3, 4, 5},
		Status:            1,
		CumulativeGasUsed: 100,
		Bloom:             gethTypes.BytesToBloom([]byte{1, 3, 4}),
		TxHash:            common.HexToHash("0x1020304050"),
		ContractAddress:   common.HexToAddress("0x1122334455"),
		GasUsed:           123,
		BlockHash:         common.HexToHash("0x11111111111111"),
		BlockNumber:       big.NewInt(555),
		TransactionIndex:  777,
		Logs: []*gethTypes.Log{
			testGethLog1,
			testGethLog2,
		},
	}
)

func Test_PersistsReadsChain(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	val := assets.NewWeiI(rand.Int63())
	addr := testutils.NewAddress()
	ks := make(map[string]types.ChainCfg)
	ks[addr.Hex()] = types.ChainCfg{EvmMaxGasPriceWei: val}
	chain := types.DBChain{
		ID: *utils.NewBigI(rand.Int63()),
		Cfg: &types.ChainCfg{
			KeySpecific: ks,
		},
	}

	evmtest.MustInsertChain(t, db, &chain)

	var loadedChain types.DBChain
	require.NoError(t, db.Get(&loadedChain, "SELECT * FROM evm_chains WHERE id = $1", chain.ID))

	loadedVal := loadedChain.Cfg.KeySpecific[addr.Hex()].EvmMaxGasPriceWei
	assert.Equal(t, loadedVal, val)
}

func TestFromGethReceipt(t *testing.T) {
	t.Parallel()

	receipt := types.FromGethReceipt(testGethReceipt)

	assert.NotNil(t, receipt)
	assert.Equal(t, testGethReceipt.PostState, receipt.PostState)
	assert.Equal(t, testGethReceipt.Status, receipt.Status)
	assert.Equal(t, testGethReceipt.CumulativeGasUsed, receipt.CumulativeGasUsed)
	assert.Equal(t, testGethReceipt.Bloom, receipt.Bloom)
	assert.Equal(t, testGethReceipt.TxHash, receipt.TxHash)
	assert.Equal(t, testGethReceipt.ContractAddress, receipt.ContractAddress)
	assert.Equal(t, testGethReceipt.GasUsed, receipt.GasUsed)
	assert.Equal(t, testGethReceipt.BlockHash, receipt.BlockHash)
	assert.Equal(t, testGethReceipt.BlockNumber, receipt.BlockNumber)
	assert.Equal(t, testGethReceipt.TransactionIndex, receipt.TransactionIndex)
	assert.Len(t, receipt.Logs, len(testGethReceipt.Logs))

	for i, log := range receipt.Logs {
		expectedLog := testGethReceipt.Logs[i]
		assert.Equal(t, expectedLog.Address, log.Address)
		assert.Equal(t, expectedLog.Topics, log.Topics)
		assert.Equal(t, expectedLog.Data, log.Data)
		assert.Equal(t, expectedLog.BlockHash, log.BlockHash)
		assert.Equal(t, expectedLog.BlockNumber, log.BlockNumber)
		assert.Equal(t, expectedLog.TxHash, log.TxHash)
		assert.Equal(t, expectedLog.TxIndex, log.TxIndex)
		assert.Equal(t, expectedLog.Index, log.Index)
		assert.Equal(t, expectedLog.Removed, log.Removed)
	}
}

func TestReceipt_IsZero(t *testing.T) {
	t.Parallel()

	receipt := types.FromGethReceipt(testGethReceipt)
	assert.False(t, receipt.IsZero())

	zeroTxHash := *testGethReceipt
	zeroTxHash.TxHash = common.HexToHash("0x0")
	receipt = types.FromGethReceipt(&zeroTxHash)
	assert.True(t, receipt.IsZero())
}

func TestReceipt_IsUnmined(t *testing.T) {
	t.Parallel()

	receipt := types.FromGethReceipt(testGethReceipt)
	assert.False(t, receipt.IsUnmined())

	zeroBlockHash := *testGethReceipt
	zeroBlockHash.BlockHash = common.HexToHash("0x0")
	receipt = types.FromGethReceipt(&zeroBlockHash)
	assert.True(t, receipt.IsUnmined())
}

func TestReceipt_MarshalUnmarshalJson(t *testing.T) {
	t.Parallel()

	receipt := types.FromGethReceipt(testGethReceipt)
	json, err := receipt.MarshalJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, json)

	parsedReceipt := &types.Receipt{}
	err = parsedReceipt.UnmarshalJSON(json)
	assert.NoError(t, err)

	assert.Equal(t, receipt, parsedReceipt)
}

func TestLog_MarshalUnmarshalJson(t *testing.T) {
	t.Parallel()

	log := types.FromGethLog(testGethLog1)
	json, err := log.MarshalJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, json)

	parsedLog := &types.Log{}
	err = parsedLog.UnmarshalJSON(json)
	assert.NoError(t, err)

	assert.Equal(t, log, parsedLog)
}
