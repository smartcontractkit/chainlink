package ccipsolana

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
)

func TestPublicKeyFromBytes(t *testing.T) {
	tests := []struct {
		name     string
		inHex    string
		isErr    bool
		expected string
	}{
		{
			"empty",
			"",
			true,
			solana.PublicKey{}.String(),
		},
		{
			"smaller than required",
			"010203040506",
			true,
			solana.PublicKey{}.String(),
		},
		{
			"equal to 32 bytes",
			"0102030405060102030405060102030405060102030405060102030405060101",
			false,
			solana.MustPublicKeyFromBase58("4wBqpZM9msxygzsdeLPq6Zw3LoiAxJk3GjtKPpqkcsi").String(),
		},
		{
			"longer than required",
			"0102030405060102030405060102030405060102030405060102030405060101FFFFFFFFFF",
			true,
			solana.PublicKey{}.String(),
		},
	}

	codec := AddressCodec{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bytes, err := hex.DecodeString(test.inHex)
			require.NoError(t, err)

			if test.isErr {
				_, err := codec.AddressBytesToString(bytes)
				require.Error(t, err)
			} else {
				actual, err := codec.AddressBytesToString(bytes)
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestPublicKeyFromBase58(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		expected    []byte
		expectedErr error
	}{
		{
			"hand crafted",
			"SerumkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
			solana.MustPublicKeyFromBase58("SerumkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA").Bytes(),
			nil,
		},
		{
			"hand crafted error",
			"SerkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
			solana.PublicKey{}.Bytes(),
			errors.New("invalid length, expected 32, got 30"),
		},
	}

	codec := AddressCodec{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := codec.AddressStringToBytes(test.in)
			if test.expectedErr == nil {
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			} else {
				require.Error(t, err)
			}
		})
	}
}
