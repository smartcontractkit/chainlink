package abihelpers

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

func TestProofFlagToBits(t *testing.T) {
	genFlags := func(indexesSet []int, size int) []bool {
		bools := make([]bool, size)
		for _, indexSet := range indexesSet {
			bools[indexSet] = true
		}
		return bools
	}
	tt := []struct {
		flags    []bool
		expected *big.Int
	}{
		{
			[]bool{true, false, true},
			big.NewInt(5),
		},
		{
			[]bool{true, true, false}, // Note the bits are reversed, slightly easier to implement.
			big.NewInt(3),
		},
		{
			[]bool{false, true, true},
			big.NewInt(6),
		},
		{
			[]bool{false, false, false},
			big.NewInt(0),
		},
		{
			[]bool{true, true, true},
			big.NewInt(7),
		},
		{
			genFlags([]int{266}, 300),
			big.NewInt(0).SetBit(big.NewInt(0), 266, 1),
		},
	}
	for _, tc := range tt {
		tc := tc
		a := ProofFlagsToBits(tc.flags)
		assert.Equal(t, tc.expected.String(), a.String())
	}
}

func TestEvmWord(t *testing.T) {
	testCases := []struct {
		inp uint64
		exp common.Hash
	}{
		{inp: 1, exp: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
		{inp: math.MaxUint64, exp: common.HexToHash("0x000000000000000000000000000000000000000000000000ffffffffffffffff")},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("test %d", tc.inp), func(t *testing.T) {
			h := EvmWord(tc.inp)
			assert.Equal(t, tc.exp, h)
		})
	}
}

func TestABIEncodeDecode(t *testing.T) {
	abiStr := `[{"components": [{"name":"int1","type":"int256"},{"name":"int2","type":"int256"}], "type":"tuple"}]`
	values := []interface{}{struct {
		Int1 *big.Int `json:"int1"`
		Int2 *big.Int `json:"int2"`
	}{big.NewInt(10), big.NewInt(12)}}

	// First encoding, should call the underlying utils.ABIEncode
	encoded, err := ABIEncode(abiStr, values...)
	assert.NoError(t, err)
	assert.NotNil(t, encoded)

	// Second encoding, should retrieve from cache
	// we're just testing here that it returns same result
	encodedAgain, err := ABIEncode(abiStr, values...)

	assert.NoError(t, err)
	assert.True(t, bytes.Equal(encoded, encodedAgain))

	// Should be able to decode it back to the original values
	decoded, err := ABIDecode(abiStr, encoded)
	assert.NoError(t, err)
	assert.Equal(t, decoded, values)
}

func BenchmarkComparisonEncode(b *testing.B) {
	abiStr := `[{"components": [{"name":"int1","type":"int256"},{"name":"int2","type":"int256"}], "type":"tuple"}]`
	values := []interface{}{struct {
		Int1 *big.Int `json:"int1"`
		Int2 *big.Int `json:"int2"`
	}{big.NewInt(10), big.NewInt(12)}}

	b.Run("WithoutCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = utils.ABIEncode(abiStr, values...)
		}
	})

	// Warm up the cache
	_, _ = ABIEncode(abiStr, values...)

	b.Run("WithCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = ABIEncode(abiStr, values...)
		}
	})
}

func BenchmarkComparisonDecode(b *testing.B) {
	abiStr := `[{"components": [{"name":"int1","type":"int256"},{"name":"int2","type":"int256"}], "type":"tuple"}]`
	values := []interface{}{struct {
		Int1 *big.Int `json:"int1"`
		Int2 *big.Int `json:"int2"`
	}{big.NewInt(10), big.NewInt(12)}}
	data, _ := utils.ABIEncode(abiStr, values...)

	b.Run("WithoutCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = utils.ABIDecode(abiStr, data)
		}
	})

	// Warm up the cache
	_, _ = ABIDecode(abiStr, data)

	b.Run("WithCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = ABIDecode(abiStr, data)
		}
	})
}
