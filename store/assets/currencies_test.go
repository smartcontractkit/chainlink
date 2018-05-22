package assets_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/stretchr/testify/assert"
)

func TestAssets_NewLinkAndString(t *testing.T) {
	t.Parallel()

	link := assets.NewLink(0)

	assert.Equal(t, "0.000000000000000000", link.String())

	link.SetInt64(1)
	assert.Equal(t, "0.000000000000000001", link.String())

	link.SetString("900000000000000000", 10)
	assert.Equal(t, "0.900000000000000000", link.String())

	link.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	assert.Equal(t, "115792089237316195423570985008687907853269984665640564039457.584007913129639935", link.String())

	link.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639936", 10)
	assert.Equal(t, "115792089237316195423570985008687907853269984665640564039457.584007913129639936", link.String())
}

func TestAssets_Link_MarshalJson(t *testing.T) {
	t.Parallel()

	link := assets.NewLink(1)

	b, err := json.Marshal(link)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"1"`), b)
}

func TestAssets_Link_UnmarshalJsonOk(t *testing.T) {
	t.Parallel()

	link := assets.Link{}

	err := json.Unmarshal([]byte(`"1"`), &link)
	assert.NoError(t, err)
	assert.Equal(t, "0.000000000000000001", link.String())
}

func TestAssets_Link_UnmarshalJsonError(t *testing.T) {
	t.Parallel()

	link := assets.Link{}

	err := json.Unmarshal([]byte(`"a"`), &link)
	assert.EqualError(t, err, "assets: cannot unmarshal \"a\" into a *assets.Link")

	err = json.Unmarshal([]byte(`1`), &link)
	assert.EqualError(t, err, "json: cannot unmarshal number into Go value of type *assets.Link")
}
