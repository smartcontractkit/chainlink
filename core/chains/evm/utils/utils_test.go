package utils_test

import (
	"context"
	"math/big"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

func TestKeccak256(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic", "0xf00b", "0x2433bb36d5f9b14e4fea87c2d32d79abfe34e56808b891e471f4400fca2a336c"},
		{"long input", "0xf00b2433bb36d5f9b14e4fea87c2d32d79abfe34e56808b891e471f4400fca2a336c", "0x6b917c56ad7bea7d09132b9e1e29bb5d9aa7d32d067c638dfa886bbbf6874cdf"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input, err := hexutil.Decode(test.input)
			assert.NoError(t, err)
			result, err := utils.Keccak256(input)
			assert.NoError(t, err)

			assert.Equal(t, test.want, hexutil.Encode(result))
		})
	}
}

func TestUtils_IsEmptyAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		addr common.Address
		want bool
	}{
		{"zero address", common.Address{}, true},
		{"non-zero address", testutils.NewAddress(), false},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := utils.IsEmptyAddress(test.addr)
			assert.Equal(t, test.want, actual)
		})
	}
}

// From https://github.com/ethereum/EIPs/blob/master/EIPS/eip-55.md#test-cases
var testAddresses = []string{
	"0x52908400098527886E0F7030069857D2E4169EE7",
	"0x8617E340B3D01FA5F11F306F4090FD50E238070D",
	"0xde709f2102306220921060314715629080e2fb77",
	"0x27b1fdb04752bbc536007a920d24acb045561c26",
	"0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed",
	"0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359",
	"0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB",
	"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb",
}

func TestClient_EIP55CapitalizedAddress(t *testing.T) {
	t.Parallel()

	valid := utils.EIP55CapitalizedAddress
	for _, address := range testAddresses {
		assert.True(t, valid(address))
		assert.False(t, valid(strings.ToLower(address)) &&
			valid(strings.ToUpper(address)))
	}
}

func TestClient_ParseEthereumAddress(t *testing.T) {
	t.Parallel()

	parse := utils.ParseEthereumAddress
	for _, address := range testAddresses {
		a1, err := parse(address)
		assert.NoError(t, err)
		no0xPrefix := address[2:]
		a2, err := parse(no0xPrefix)
		assert.NoError(t, err)
		assert.True(t, a1 == a2)
		_, lowerErr := parse(strings.ToLower(address))
		_, upperErr := parse(strings.ToUpper(address))
		shouldBeError := multierr.Combine(lowerErr, upperErr)
		assert.Error(t, shouldBeError)
		assert.True(t, strings.Contains(shouldBeError.Error(), no0xPrefix))
	}
	_, notHexErr := parse("0xCeci n'est pas une chaîne hexadécimale")
	assert.Error(t, notHexErr)
	_, tooLongErr := parse("0x0123456789abcdef0123456789abcdef0123456789abcdef")
	assert.Error(t, tooLongErr)
}

func TestUint256ToBytes(t *testing.T) {
	t.Parallel()

	v := big.NewInt(0).Sub(utils.MaxUint256, big.NewInt(1))
	uint256, err := utils.Uint256ToBytes(v)
	assert.NoError(t, err)

	b32 := utils.Uint256ToBytes32(v)
	assert.Equal(t, uint256, b32)

	large := big.NewInt(0).Add(utils.MaxUint256, big.NewInt(1))
	_, err = utils.Uint256ToBytes(large)
	assert.Error(t, err, "too large to convert to uint256")

	negative := big.NewInt(-1)
	assert.Panics(t, func() {
		_, _ = utils.Uint256ToBytes(negative)
	}, "failed to round-trip uint256 back to source big.Int")
}

func TestCheckUint256(t *testing.T) {
	t.Parallel()

	large := big.NewInt(0).Add(utils.MaxUint256, big.NewInt(1))
	err := utils.CheckUint256(large)
	assert.Error(t, err, "number out of range for uint256")

	negative := big.NewInt(-123)
	err = utils.CheckUint256(negative)
	assert.Error(t, err, "number out of range for uint256")

	err = utils.CheckUint256(big.NewInt(123))
	assert.NoError(t, err)
}

func TestRandUint256(t *testing.T) {
	t.Parallel()

	for i := 0; i < 1000; i++ {
		uint256 := utils.RandUint256()
		assert.NoError(t, utils.CheckUint256(uint256))
	}
}

func TestHexToUint256(t *testing.T) {
	t.Parallel()

	b, err := utils.HexToUint256("0x00")
	assert.NoError(t, err)
	assert.Zero(t, b.Cmp(big.NewInt(0)))

	b, err = utils.HexToUint256("0xFFFFFFFF")
	assert.NoError(t, err)
	assert.Zero(t, b.Cmp(big.NewInt(4294967295)))
}

func TestNewHash(t *testing.T) {
	t.Parallel()

	h1 := utils.NewHash()
	h2 := utils.NewHash()
	assert.NotEqual(t, h1, h2)
	assert.NotEqual(t, h1, common.HexToHash("0x0"))
	assert.NotEqual(t, h2, common.HexToHash("0x0"))
}

func TestPadByteToHash(t *testing.T) {
	t.Parallel()

	h := utils.PadByteToHash(1)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000001", h.String())
}

func TestUtils_BackoffSleeper(t *testing.T) {
	t.Parallel()

	bs := utils.NewBackoffSleeper()
	assert.Equal(t, time.Duration(0), bs.Duration(), "should initially return immediately")
	bs.Sleep()

	d := 1 * time.Nanosecond
	bs.Min = d
	bs.Factor = 2
	assert.Equal(t, d, bs.Duration())
	bs.Sleep()

	d2 := 2 * time.Nanosecond
	assert.Equal(t, d2, bs.Duration())

	bs.Reset()
	assert.Equal(t, time.Duration(0), bs.Duration(), "should initially return immediately")
}

func TestRetryWithBackoff(t *testing.T) {
	t.Parallel()

	var counter atomic.Int32
	ctx, cancel := context.WithCancel(tests.Context(t))

	utils.RetryWithBackoff(ctx, func() bool {
		return false
	})

	retry := func() bool {
		return counter.Add(1) < 3
	}

	go utils.RetryWithBackoff(ctx, retry)

	assert.Eventually(t, func() bool {
		return counter.Load() == 3
	}, tests.WaitTimeout(t), tests.TestInterval)

	cancel()

	utils.RetryWithBackoff(ctx, retry)
	assert.Equal(t, int32(4), counter.Load())
}
