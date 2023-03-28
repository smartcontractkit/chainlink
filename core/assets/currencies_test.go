package assets_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
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
	assert.Equal(t, assets.ErrNoQuotesForCurrency, err)
}

func TestAssets_LinkToInt(t *testing.T) {
	t.Parallel()

	link := assets.NewLinkFromJuels(0)
	assert.Equal(t, big.NewInt(0), link.ToInt())

	link = assets.NewLinkFromJuels(123)
	assert.Equal(t, big.NewInt(123), link.ToInt())
}

func TestAssets_LinkToHash(t *testing.T) {
	t.Parallel()

	link := assets.NewLinkFromJuels(123)
	expected := common.BigToHash((*big.Int)(link))
	assert.Equal(t, expected, link.ToHash())
}

func TestAssets_LinkSetLink(t *testing.T) {
	t.Parallel()

	link1 := assets.NewLinkFromJuels(123)
	link2 := assets.NewLinkFromJuels(321)
	link3 := link1.Set(link2)
	assert.Equal(t, link3, link2)
}

func TestAssets_LinkCmpLink(t *testing.T) {
	t.Parallel()

	link1 := assets.NewLinkFromJuels(123)
	link2 := assets.NewLinkFromJuels(321)
	assert.NotZero(t, link1.Cmp(link2))

	link3 := assets.NewLinkFromJuels(321)
	assert.Zero(t, link3.Cmp(link2))
}

func TestAssets_LinkAddLink(t *testing.T) {
	t.Parallel()

	link1 := assets.NewLinkFromJuels(123)
	link2 := assets.NewLinkFromJuels(321)
	sum := assets.NewLinkFromJuels(123 + 321)
	assert.Equal(t, sum, link1.Add(link1, link2))
}

func TestAssets_LinkText(t *testing.T) {
	t.Parallel()

	link := assets.NewLinkFromJuels(123)
	assert.Equal(t, "123", link.Text(10))
	assert.Equal(t, "7b", link.Text(16))
}

func TestAssets_LinkIsZero(t *testing.T) {
	t.Parallel()

	link := assets.NewLinkFromJuels(123)
	assert.False(t, link.IsZero())

	link = assets.NewLinkFromJuels(0)
	assert.True(t, link.IsZero())
}

func TestAssets_LinkSymbol(t *testing.T) {
	t.Parallel()

	link := assets.NewLinkFromJuels(123)
	assert.Equal(t, "LINK", link.Symbol())
}

func TestAssets_LinkScanValue(t *testing.T) {
	t.Parallel()

	link := assets.NewLinkFromJuels(123)
	v, err := link.Value()
	assert.NoError(t, err)

	link2 := assets.NewLinkFromJuels(0)
	err = link2.Scan(v)
	assert.NoError(t, err)
	assert.Equal(t, link2, link)

	err = link2.Scan("123")
	assert.NoError(t, err)
	assert.Equal(t, link2, link)

	err = link2.Scan([]uint8{'1', '2', '3'})
	assert.NoError(t, err)
	assert.Equal(t, link2, link)

	assert.ErrorContains(t, link2.Scan([]uint8{'x'}), "unable to set string")
	assert.ErrorContains(t, link2.Scan("123.56"), "unable to set string")
	assert.ErrorContains(t, link2.Scan(1.5), "unable to convert")
	assert.ErrorContains(t, link2.Scan(int64(123)), "unable to convert")
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

func TestLink(t *testing.T) {
	for _, tt := range []struct {
		input string
		exp   string
	}{
		{"0", "0"},
		{"1", "1"},
		{"1 juels", "1"},
		{"100000000000", "100000000000"},
		{"0.0000001 link", "100000000000"},
		{"1000000000000", "0.000001 link"},
		{"100000000000000", "0.0001 link"},
		{"0.0001 link", "0.0001 link"},
		{"10000000000000000", "0.01 link"},
		{"0.01 link", "0.01 link"},
		{"100000000000000000", "0.1 link"},
		{"0.1 link", "0.1 link"},
		{"1.0 link", "1 link"},
		{"1000000000000000000", "1 link"},
		{"1000000000000000000 juels", "1 link"},
		{"1100000000000000000", "1.1 link"},
		{"1.1link", "1.1 link"},
		{"1.1 link", "1.1 link"},
	} {
		t.Run(tt.input, func(t *testing.T) {
			var l assets.Link
			err := l.UnmarshalText([]byte(tt.input))
			require.NoError(t, err)
			b, err := l.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tt.exp, string(b))
		})
	}
}

func FuzzLink(f *testing.F) {
	f.Add("1")
	f.Add("1 link")
	f.Add("1.1link")
	f.Add("2.3")
	f.Add("2.3 link")
	f.Add("00005 link")
	f.Add("0.0005link")
	f.Add("1100000000000000000000000000000")
	f.Add("1100000000000000000000000000000 juels")
	f.Fuzz(func(t *testing.T, v string) {
		if len(v) > 1_000 {
			t.Skip()
		}
		var l assets.Link
		err := l.UnmarshalText([]byte(v))
		if err != nil {
			t.Skip()
		}

		b, err := l.MarshalText()
		require.NoErrorf(t, err, "failed to marshal %v after unmarshaling from %q", l, v)

		var l2 assets.Link
		err = l2.UnmarshalText(b)
		require.NoErrorf(t, err, "failed to unmarshal %s after marshaling from %v", string(b), l)
		require.Equal(t, l, l2, "unequal values after marshal/unmarshal")
	})
}
