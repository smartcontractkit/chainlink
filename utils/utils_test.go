package utils_test

import (
	"testing"
	"time"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestUtils_NewBytes32ID(t *testing.T) {
	t.Parallel()
	id := utils.NewBytes32ID()
	assert.NotContains(t, id, "-")
}

func TestUtils_WeiToEth(t *testing.T) {
	t.Parallel()
	var numWei *big.Int = new(big.Int).SetInt64(1)
	var expectedNumEth float64 = 1e-18
	actualNumEth := utils.WeiToEth(numWei)
	assert.Equal(t, expectedNumEth, actualNumEth)
}

func TestUtils_EthToWei(t *testing.T) {
	t.Parallel()
	var numEth float64 = 1.0
	var expectedNumWei *big.Int = new(big.Int).SetInt64(1e18)
	actualNumWei := utils.EthToWei(numEth)
	assert.Equal(t, actualNumWei, expectedNumWei)
}

func TestUtils_IsEmptyAddress(t *testing.T) {
	tests := []struct {
		name string
		addr common.Address
		want bool
	}{
		{"zero address", common.Address{}, true},
		{"non-zero address", cltest.NewAddress(), false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual := utils.IsEmptyAddress(test.addr)
			assert.Equal(t, test.want, actual)
		})
	}
}

func TestUtils_StringToHex(t *testing.T) {
	tests := []struct {
		utf8 string
		hex  string
	}{
		{"abc", "0x616263"},
		{"Hi Mom!", "0x4869204d6f6d21"},
		{"", "0x"},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.utf8, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.hex, utils.StringToHex(test.utf8))
		})
	}
}

func TestUtils_HexToString(t *testing.T) {
	tests := []struct {
		hex     string
		utf8    string
		errored bool
	}{
		{"0x616263", "abc", false},
		{"616263", "abc", false},
		{"0x4869204d6f6d21", "Hi Mom!", false},
		{"0x", "", false},
		{"uh oh", "", true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.hex, func(t *testing.T) {
			t.Parallel()
			actualUtf8, err := utils.HexToString(test.hex)
			assert.Equal(t, test.errored, err != nil)
			assert.Equal(t, test.utf8, actualUtf8)
		})
	}
}

func TestUtils_BackoffSleeper(t *testing.T) {
	bs := utils.NewBackoffSleeper()
	d := 1 * time.Nanosecond
	bs.Min = d
	bs.Factor = 2
	assert.Equal(t, d, bs.Duration())
	bs.Sleep()
	d2 := 2 * time.Nanosecond
	assert.Equal(t, d2, bs.Duration())
}
