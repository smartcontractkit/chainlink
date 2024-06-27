package types_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
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

// constraint satisfied by common.Hash and common.Address
type ByteString interface {
	Bytes() []byte
}

type ScannableArrayType interface {
	Scan(src any) error
}

type HexArrayScanTestArgs struct {
	b1        []byte
	b2        []byte
	wrongsize []byte
}

func testHexArrayScan[T ScannableArrayType](t *testing.T, dest T, args HexArrayScanTestArgs) {
	b0 := "NULL"
	empty := "{}"
	b1, b2, wrongsize := args.b1, args.b2, args.wrongsize

	src1 := fmt.Sprintf("{\"\\\\x%x\"}", b1)
	src2 := fmt.Sprintf("{\"\\\\x%x\",\"\\\\x%x\"}", b2, b2)
	src3 := fmt.Sprintf("{\"\\\\x%x\"}", wrongsize)
	invalid := fmt.Sprintf("{\"\\\\x%x\", NULL}", b1)
	d2 := fmt.Sprintf("[1][1]={{\"\\\\x%x\"}}", b1)

	get := func(d T, ind int) (bs ByteString) {
		switch val := (ScannableArrayType(dest)).(type) {
		case *types.HashArray:
			bs = ([]common.Hash(*val))[ind]
		case *types.AddressArray:
			bs = ([]common.Address(*val))[ind]
		}
		return bs
	}

	length := func(d T) (l int) {
		switch val := (ScannableArrayType(dest)).(type) {
		case *types.HashArray:
			l = len([]common.Hash(*val))
		case *types.AddressArray:
			l = len([]common.Address(*val))
		}
		return l
	}

	err := dest.Scan(b0)
	require.Error(t, err)

	err = dest.Scan(empty)
	assert.NoError(t, err)

	err = dest.Scan(src1)
	require.NoError(t, err)
	require.Equal(t, length(dest), 1)
	assert.Equal(t, get(dest, 0).Bytes(), b1)

	err = dest.Scan(src2)
	require.NoError(t, err)
	require.Equal(t, length(dest), 3)
	assert.Equal(t, get(dest, 1).Bytes(), b2)
	assert.Equal(t, get(dest, 2).Bytes(), b2)

	err = dest.Scan(src3)
	require.Error(t, err)

	err = dest.Scan(invalid)
	require.Error(t, err)

	err = dest.Scan(d2)
	require.Error(t, err)
}

func Test_AddressArrayScan(t *testing.T) {
	t.Parallel()
	addr1, err := hex.DecodeString("2ab9a2dc53736b361b72d900cdf9f78f9406fbbb")
	require.NoError(t, err)
	require.Len(t, addr1, 20)
	addr2, err := hex.DecodeString("56b9a2dc53736b361b72d900cdf9f78f9406fbbb")
	require.Len(t, addr2, 20)
	toolong, err := hex.DecodeString("6b361b72d900cdf9f78f9406fbbb6b361b72d900cdf9f78f9406fbbb")
	require.Len(t, toolong, 28)

	a := types.AddressArray{}
	args := HexArrayScanTestArgs{addr1, addr2, toolong}
	testHexArrayScan[*types.AddressArray](t, &a, args)
}

func Test_HashArrayScan(t *testing.T) {
	t.Parallel()

	h1, err := hex.DecodeString("2ab9130c6b361b72d900cdf9f78f9406fbbb6b361b72d900cdf9f78f9406fbbb")
	require.NoError(t, err)
	require.Len(t, h1, 32)
	h2, err := hex.DecodeString("56b9a2dc53736b361b72d900cdf9f78f9406fbbb06fbbb6b361b7206fbbb6b36")
	require.Len(t, h2, 32)
	tooshort, err := hex.DecodeString("6b361b72d900cdf9f78f9406fbbb6b361b72d900cdf9f78f9406fbbb")
	require.Len(t, tooshort, 28)

	h := types.HashArray{}
	args := HexArrayScanTestArgs{h1, h2, tooshort}
	testHexArrayScan[*types.HashArray](t, &h, args)
}
