package testutils

import (
	"encoding/hex"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
)

func AssertLiquidityTransferredEventSlicesEqual(
	t *testing.T,
	expected,
	actual []*liquiditymanager.LiquidityManagerLiquidityTransferred,
	sortComparator func(a, b *liquiditymanager.LiquidityManagerLiquidityTransferred) bool,
) {
	require.Equal(t, len(expected), len(actual))
	sort.Slice(expected, func(i, j int) bool {
		return sortComparator(expected[i], expected[j])
	})
	sort.Slice(actual, func(i, j int) bool {
		return sortComparator(actual[i], actual[j])
	})
	for i := range expected {
		assert.Equal(t, expected[i].OcrSeqNum, actual[i].OcrSeqNum)
		assert.Equal(t, expected[i].FromChainSelector, actual[i].FromChainSelector)
		assert.Equal(t, expected[i].ToChainSelector, actual[i].ToChainSelector)
		assert.Equal(t, expected[i].To, actual[i].To)
		assert.Equal(t, expected[i].Amount, actual[i].Amount)
		assert.Equal(t, expected[i].BridgeSpecificData, actual[i].BridgeSpecificData)
		assert.Equal(t, expected[i].BridgeReturnData, actual[i].BridgeReturnData)
	}
}

func SortByBridgeReturnData(a, b *liquiditymanager.LiquidityManagerLiquidityTransferred) bool {
	return hex.EncodeToString(a.BridgeReturnData) < hex.EncodeToString(b.BridgeReturnData)
}

func MustPackBridgeData(t *testing.T, bridgeDataHex string) []byte {
	packed, err := hex.DecodeString(bridgeDataHex[2:])
	require.NoError(t, err)
	return packed
}

func MustConvertHexBridgeDataToBytes(t *testing.T, hexData string) []byte {
	packed, err := hex.DecodeString(hexData[2:])
	require.NoError(t, err)
	return packed
}
