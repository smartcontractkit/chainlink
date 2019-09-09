package models

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBig_UnmarshalText(t *testing.T) {
	t.Parallel()

	i := &Big{}
	tests := []struct {
		name  string
		input string
		want  *big.Int
	}{
		{"number", `1234`, big.NewInt(1234)},
		{"string", `"1234"`, big.NewInt(1234)},
		{"hex number", `0x1234`, big.NewInt(4660)},
		{"hex string", `"0x1234"`, big.NewInt(4660)},
		{"single quoted", `'1234'`, big.NewInt(1234)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := i.UnmarshalText([]byte(test.input))
			require.NoError(t, err)
			assert.Equal(t, test.want, i.ToInt())
		})
	}
}

func TestBig_UnmarshalTextErrors(t *testing.T) {
	t.Parallel()

	i := &Big{}
	tests := []struct {
		name  string
		input string
		want  *big.Int
	}{
		{"quoted word", `"word"`, big.NewInt(0)},
		{"word", `word`, big.NewInt(0)},
		{"empty", ``, big.NewInt(0)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := i.UnmarshalText([]byte(test.input))
			require.Error(t, err)
		})
	}
}

func TestBig_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *big.Int
		want  string
	}{
		{"number", big.NewInt(1234), `1234`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := (*Big)(test.input)
			b, err := json.Marshal(&i)
			assert.NoError(t, err)
			assert.Equal(t, test.want, string(b))
		})
	}
}

func TestBig_Scan(t *testing.T) {
	t.Parallel()

	uint256Max, ok := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	require.True(t, ok)

	tests := []struct {
		name  string
		input interface{}
		want  *Big
	}{
		{"zero string", "0", NewBig(big.NewInt(0))},
		{"one string", "1", NewBig(big.NewInt(1))},
		{
			"large string",
			"115792089237316195423570985008687907853269984665640564039457584007913129639935",
			NewBig(uint256Max),
		},
		{"zero as bytes", []uint8{48}, NewBig(big.NewInt(0))},
		{"small number as bytes", []uint8{49, 52}, NewBig(big.NewInt(14))},
		{
			"max number as bytes",
			[]uint8{
				49, 49, 53, 55, 57, 50, 48, 56, 57, 50, 51, 55, 51, 49, 54, 49, 57, 53,
				52, 50, 51, 53, 55, 48, 57, 56, 53, 48, 48, 56, 54, 56, 55, 57, 48, 55,
				56, 53, 51, 50, 54, 57, 57, 56, 52, 54, 54, 53, 54, 52, 48, 53, 54, 52,
				48, 51, 57, 52, 53, 55, 53, 56, 52, 48, 48, 55, 57, 49, 51, 49, 50, 57,
				54, 51, 57, 57, 51, 53,
			},
			NewBig(uint256Max),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			big := &Big{}
			err := big.Scan(test.input)
			require.NoError(t, err)
			assert.Equal(t, test.want, big)
		})
	}
}

func TestBig_ScanErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input interface{}
	}{
		{"zero integer", 0},
		{"one integer", 1},
		{"zero wrapped string", `"0"`},
		{"one wrapped string", `"1"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			big := &Big{}
			err := big.Scan(test.input)
			require.Error(t, err)
		})
	}
}
