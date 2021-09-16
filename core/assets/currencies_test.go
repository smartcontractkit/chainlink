package assets_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"

	"github.com/stretchr/testify/assert"
)

func TestAssets_NewLinkAndString(t *testing.T) {
	t.Parallel()

	link := assets.NewLinkFromJuels(0)

	assert.Equal(t, "0", link.String())

	link.SetInt64(1)
	assert.Equal(t, "1", link.String())

	link.SetString("900000000000000000", 10)
	assert.Equal(t, "900000000000000000", link.String())

	link.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	assert.Equal(t, "115792089237316195423570985008687907853269984665640564039457584007913129639935", link.String())

	var nilLink *assets.Link
	assert.Equal(t, "0", nilLink.String())
}

func TestAssets_NewLinkAndLink(t *testing.T) {
	t.Parallel()

	link := assets.NewLinkFromJuels(0)

	assert.Equal(t, "0.000000000000000000", link.Link())

	link.SetInt64(1)
	assert.Equal(t, "0.000000000000000001", link.Link())

	link.SetString("900000000000000000", 10)
	assert.Equal(t, "0.900000000000000000", link.Link())

	link.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	assert.Equal(t, "115792089237316195423570985008687907853269984665640564039457.584007913129639935", link.Link())

	link.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639936", 10)
	assert.Equal(t, "115792089237316195423570985008687907853269984665640564039457.584007913129639936", link.Link())

	var nilLink *assets.Link
	assert.Equal(t, "0", nilLink.Link())
}

func TestAssets_Link_MarshalJson(t *testing.T) {
	t.Parallel()

	link := assets.NewLinkFromJuels(1)

	b, err := json.Marshal(link)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"1"`), b)
}

func TestAssets_Link_UnmarshalJsonOk(t *testing.T) {
	t.Parallel()

	link := assets.Link{}

	err := json.Unmarshal([]byte(`"1"`), &link)
	assert.NoError(t, err)
	assert.Equal(t, "0.000000000000000001", link.Link())
}

func TestAssets_Link_UnmarshalJsonError(t *testing.T) {
	t.Parallel()

	link := assets.Link{}

	err := json.Unmarshal([]byte(`"x"`), &link)
	assert.EqualError(t, err, "assets: cannot unmarshal \"x\" into a *assets.Link")

	err = json.Unmarshal([]byte(`1`), &link)
	assert.Equal(t, assets.ErrNoQuotesForCurrency, err)
}

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

	oneWei := assets.NewEth(1)
	assert.False(t, oneWei.IsZero())
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
	assert.Equal(t, assets.ErrNoQuotesForCurrency, err)
}
