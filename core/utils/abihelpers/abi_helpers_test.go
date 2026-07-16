package abihelpers

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
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
			[]bool{true, true, false},
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
	values := []any{struct {
		Int1 *big.Int `json:"int1"`
		Int2 *big.Int `json:"int2"`
	}{big.NewInt(10), big.NewInt(12)}}

	encoded, err := ABIEncode(abiStr, values...)
	assert.NoError(t, err)
	assert.NotNil(t, encoded)

	encodedAgain, err := ABIEncode(abiStr, values...)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(encoded, encodedAgain))

	decoded, err := ABIDecode(abiStr, encoded)
	assert.NoError(t, err)
	assert.Equal(t, decoded, values)
}
