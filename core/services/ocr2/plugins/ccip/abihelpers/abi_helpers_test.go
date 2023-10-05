package abihelpers

import (
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
