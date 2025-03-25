package ccipevm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressBytesToString(t *testing.T) {
	addressCodec := AddressCodec{}
	addr := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13}
	want := "0x000102030405060708090a0b0c0d0e0f10111213"
	got, err := addressCodec.AddressBytesToString(addr)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestAddressStringToBytes(t *testing.T) {
	addressCodec := AddressCodec{}
	addr := "0x000102030405060708090a0b0c0d0e0f10111213"
	want := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13}
	got, err := addressCodec.AddressStringToBytes(addr)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

// we allow various sizes since some contracts store the 20-byte address as 32-byte
// func TestInvalidAddressBytesToString(t *testing.T) {
// 	addressCodec := AddressCodec{}
// 	addr := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12}
// 	_, err := addressCodec.AddressBytesToString(addr)
// 	require.Error(t, err)
// }

func TestInvalidAddressStringToBytes(t *testing.T) {
	addressCodec := AddressCodec{}
	addr := "0x000102030405060708090a0b0c0d0e0f1011121"
	_, err := addressCodec.AddressStringToBytes(addr)
	require.Error(t, err)
}

func TestValidEVMAddress(t *testing.T) {
	addressCodec := AddressCodec{}
	addr := []byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}
	want := "0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF"
	got, err := addressCodec.AddressBytesToString(addr)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestInvalidHexString(t *testing.T) {
	addressCodec := AddressCodec{}
	addr := "0xZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ"
	_, err := addressCodec.AddressStringToBytes(addr)
	require.Error(t, err)
}
