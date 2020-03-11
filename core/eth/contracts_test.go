package eth

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetContractCodec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		contract  string
		expectErr bool
	}{
		{"Get Oracle contract", "Oracle", false},
		{"Get non-existent contract", "not-a-contract", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			contract, err := GetContractCodec(test.contract)
			if test.expectErr {
				assert.Error(t, err)
				assert.Nil(t, contract)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, contract)
			}
		})
	}
}

var address common.Address = common.HexToAddress(
	"0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
)

// NB: This test needs a compiled oracle contract, which can be built with
// `yarn workspace chainlink run setup` in the base project directory.
func TestContractCodec_EncodeMessageCall(t *testing.T) {
	t.Parallel()

	// Test with the Oracle contract
	oracle, err := GetContractCodec("Oracle")
	require.NoError(t, err)
	require.NotNil(t, oracle)

	data, err := oracle.EncodeMessageCall("withdraw", address, big.NewInt(10))
	assert.NoError(t, err)
	assert.NotNil(t, data)
}

// NB: This test needs a compiled oracle contract, which can be built with
// `yarn workspace chainlink run setup` in the base project directory.
func TestContractCodec_EncodeMessageCall_errors(t *testing.T) {
	t.Parallel()

	// Test with the Oracle contract
	oracle, err := GetContractCodec("Oracle")
	require.NoError(t, err)
	require.NotNil(t, oracle)

	tenLINK := big.NewInt(10)

	tests := []struct {
		name   string
		method string
		args   []interface{}
	}{
		{"Non-existent method", "not-a-method", []interface{}{address, tenLINK}},
		{"Too few arguments", "withdraw", []interface{}{address}},
		{"Too many arguments", "withdraw", []interface{}{address, tenLINK, tenLINK}},
		{"Incorrect argument types", "withdraw", []interface{}{tenLINK, address}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := oracle.EncodeMessageCall(test.method, test.args...)
			assert.Error(t, err)
			assert.Nil(t, data)
		})
	}
}
