package assets_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
)

func TestAssets_NewEthAndString(t *testing.T) {
	t.Parallel()

	eth := assets.NewEth(0)

	assert.Equal(t, "0.000000000000000000", eth.String())

	eth.SetInt64(1)
	assert.Equal(t, "0.000000000000000001", eth.String())

	eth.SetString("900000000000000000", 10)
	assert.Equal(t, "0.900000000000000000", eth.String())

	eth.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	assert.Equal(t, "115792089237316195423570985008687907853269984665640564039457.584007913129639935", eth.String())

	eth.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639936", 10)
	assert.Equal(t, "115792089237316195423570985008687907853269984665640564039457.584007913129639936", eth.String())
}

func TestAssets_Eth_IsZero(t *testing.T) {
	t.Parallel()

	zeroEth := assets.NewEth(0)
	assert.True(t, zeroEth.IsZero())

	oneLink := assets.NewEth(1)
	assert.False(t, oneLink.IsZero())
}

func TestAssets_Eth_MarshalJson(t *testing.T) {
	t.Parallel()

	eth := assets.NewEth(1)

	b, err := json.Marshal(eth)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"1"`), b)
}

func TestAssets_Eth_UnmarshalJsonOk(t *testing.T) {
	t.Parallel()

	eth := assets.Eth{}

	err := json.Unmarshal([]byte(`"1"`), &eth)
	assert.NoError(t, err)
	assert.Equal(t, "0.000000000000000001", eth.String())
}

func TestAssets_Eth_UnmarshalJsonError(t *testing.T) {
	t.Parallel()

	eth := assets.Eth{}

	err := json.Unmarshal([]byte(`"x"`), &eth)
	assert.EqualError(t, err, "assets: cannot unmarshal \"x\" into a *assets.Eth")

	err = json.Unmarshal([]byte(`1`), &eth)
	assert.Equal(t, commonassets.ErrNoQuotesForCurrency, err)
}

func TestAssets_NewEth(t *testing.T) {
	t.Parallel()

	ethRef := assets.NewEth(123)
	ethVal := assets.NewEthValue(123)
	ethStr, err := assets.NewEthValueS(ethRef.String())
	assert.NoError(t, err)
	assert.Equal(t, *ethRef, ethVal)
	assert.Equal(t, *ethRef, ethStr)
}

func TestAssets_EthSymbol(t *testing.T) {
	t.Parallel()

	eth := assets.NewEth(123)
	assert.Equal(t, "ETH", eth.Symbol())
}

func TestAssets_EthScanValue(t *testing.T) {
	t.Parallel()

	eth := assets.NewEth(123)
	v, err := eth.Value()
	assert.NoError(t, err)

	eth2 := assets.NewEth(0)
	err = eth2.Scan(v)
	assert.NoError(t, err)

	assert.Equal(t, eth, eth2)
}

func TestAssets_EthCmpEth(t *testing.T) {
	t.Parallel()

	eth1 := assets.NewEth(123)
	eth2 := assets.NewEth(321)
	assert.NotZero(t, eth1.Cmp(eth2))

	eth3 := assets.NewEth(321)
	assert.Zero(t, eth3.Cmp(eth2))
}
