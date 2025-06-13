package ccipaptos

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func generateAddressBytes32() [32]byte {
	var result [32]byte
	_, err := rand.Read(result[:])
	if err != nil {
		panic(fmt.Sprintf("failed to generate random address bytes: %v", err))
	}
	return result
}

func generateAddressBytes() []byte {
	a := generateAddressBytes32()
	return a[:]
}

func generateAddressString() string {
	addressBytes := generateAddressBytes()
	addressString, err := addressBytesToString(addressBytes)
	if err != nil {
		panic(fmt.Sprintf("failed to generate random address string: %v", err))
	}
	return addressString
}

func TestAddressBytesToString(t *testing.T) {
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
			"",
		},
		{
			"equal to 32 bytes",
			"0102030405060102030405060102030405060102030405060102030405060101",
			false,
			"0x0102030405060102030405060102030405060102030405060102030405060101",
		},
		{
			"longer than required",
			"0102030405060102030405060102030405060102030405060102030405060101FFFFFFFFFF",
			true,
			"",
		},
		{
			"shorter than required",
			"010203040506",
			false,
			"0x0000000000000000000000000000000000000000000000000000010203040506",
		},
	}

	codec := AddressCodec{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bytes, err := hex.DecodeString(test.inHex)
			require.NoError(t, err)

			if test.isErr {
				_, err := codec.AddressBytesToString(bytes)
				require.Error(t, err, "expected error for %s, input %s", test.name, test.inHex)
			} else {
				actual, err := codec.AddressBytesToString(bytes)
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestAddressStringToBytes(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		expected    []byte
		expectedErr error
	}{
		{
			"hand crafted",
			"0x0102030405060102030405060102030405060102030405060102030405060101",
			[]byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
				0x01, 0x01,
			},
			nil,
		},
		{
			"hand crafted error",
			"invalidAddress",
			nil,
			errors.New("failed to decode Aptos address 'invalidAddress': invalid address"),
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
