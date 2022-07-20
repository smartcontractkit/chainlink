package ethkey_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
)

func TestEIP55Address(t *testing.T) {
	t.Parallel()

	address := ethkey.EIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")

	assert.Equal(t, []byte{
		0xa0, 0x78, 0x8f, 0xc1, 0x7b, 0x1d, 0xee, 0x36,
		0xf0, 0x57, 0xc4, 0x2b, 0x6f, 0x37, 0x3a, 0x34,
		0xb0, 0x14, 0x68, 0x7e,
	}, address.Bytes())

	bi, _ := (new(big.Int)).SetString("a0788FC17B1dEe36f057c42B6F373A34B014687e", 16)
	assert.Equal(t, bi, address.Big())

	assert.Equal(t, "0xa0788FC17B1dEe36f057c42B6F373A34B014687e", address.String())

	assert.Equal(t, common.BytesToHash([]byte{
		0xa0, 0x78, 0x8f, 0xc1, 0x7b, 0x1d, 0xee, 0x36,
		0xf0, 0x57, 0xc4, 0x2b, 0x6f, 0x37, 0x3a, 0x34,
		0xb0, 0x14, 0x68, 0x7e,
	}), address.Hash())

	assert.Equal(t, "0xa0788FC17B1dEe36f057c42B6F373A34B014687e", address.String())

	zeroAddress := ethkey.EIP55Address("")
	err := json.Unmarshal([]byte(`"0xa0788FC17B1dEe36f057c42B6F373A34B014687e"`), &zeroAddress)
	assert.NoError(t, err)
	assert.Equal(t, "0xa0788FC17B1dEe36f057c42B6F373A34B014687e", zeroAddress.String())

	zeroAddress = ethkey.EIP55Address("")
	err = zeroAddress.UnmarshalText([]byte("0xa0788FC17B1dEe36f057c42B6F373A34B014687e"))
	assert.NoError(t, err)
	assert.Equal(t, "0xa0788FC17B1dEe36f057c42B6F373A34B014687e", zeroAddress.String())
}

func TestValidateEIP55Address(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid address", "0xa0788FC17B1dEe36f057c42B6F373A34B014687e", true},
		{"lowercase address", "0xa0788fc17b1dee36f057c42b6f373a34b014687e", false},
		{"invalid checksum", "0xA0788FC17B1dEe36f057c42B6F373A34B014687e", false},
		{"no leading 0x", "A0788FC17B1dEe36f057c42B6F373A34B014687e", false},
		{"non hex character", "0xz0788FC17B1dEe36f057c42B6F373A34B014687e", false},
		{"wrong length", "0xa0788FC17B1dEe36f057c42B6F373A34B014687", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ethkey.NewEIP55Address(test.input)
			valid := err == nil
			assert.Equal(t, test.valid, valid)
		})
	}
}

func TestEIP55AddressFromAddress(t *testing.T) {
	t.Parallel()

	addr := common.HexToAddress("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")
	eip55 := ethkey.EIP55AddressFromAddress(addr)
	assert.Equal(t, addr, eip55.Address())
}

func TestEIP55Address_Scan_Value(t *testing.T) {
	t.Parallel()

	eip55, err := ethkey.NewEIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")
	assert.NoError(t, err)

	val, err := eip55.Value()
	assert.NoError(t, err)

	var eip55New ethkey.EIP55Address
	err = eip55New.Scan(val)
	assert.NoError(t, err)

	assert.Equal(t, eip55, eip55New)
}

func TestEIP55AddressCollection_Scan_Value(t *testing.T) {
	t.Parallel()

	collection := ethkey.EIP55AddressCollection{
		ethkey.EIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B0146111"),
		ethkey.EIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B0146222"),
	}

	val, err := collection.Value()
	assert.NoError(t, err)

	var collectionNew ethkey.EIP55AddressCollection
	err = collectionNew.Scan(val)
	assert.NoError(t, err)

	assert.Equal(t, collection, collectionNew)
}

func TestEIP55Address_IsZero(t *testing.T) {
	t.Parallel()

	eip55 := ethkey.EIP55AddressFromAddress(common.HexToAddress("0x0"))
	assert.True(t, eip55.IsZero())

	eip55 = ethkey.EIP55AddressFromAddress(common.HexToAddress("0x1"))
	assert.False(t, eip55.IsZero())
}
