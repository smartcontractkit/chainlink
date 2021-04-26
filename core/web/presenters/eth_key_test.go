package presenters

import (
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestETHKeyResource(t *testing.T) {
	var (
		now       = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		nextNonce = int64(1)
	)
	address, err := models.NewEIP55Address("0x2aCFF2ec69aa9945Ed84f4F281eCCF6911A3B0eD")
	require.NoError(t, err)

	key := models.Key{
		ID:        1,
		Address:   address,
		CreatedAt: now,
		UpdatedAt: now,
		NextNonce: nextNonce,
		LastUsed:  &now,
		IsFunding: true,
	}

	r, err := NewETHKeyResource(key,
		SetETHKeyEthBalance(assets.NewEth(1)),
		SetETHKeyLinkBalance(assets.NewLink(1)),
	)
	require.NoError(t, err)

	assert.Equal(t, assets.NewEth(1), r.EthBalance)
	assert.Equal(t, assets.NewLink(1), r.LinkBalance)

	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := `
	{
		"data":{
		   "type":"eTHKeys",
		   "id":"0x2aCFF2ec69aa9945Ed84f4F281eCCF6911A3B0eD",
		   "attributes":{
			  "address":"0x2aCFF2ec69aa9945Ed84f4F281eCCF6911A3B0eD",
			  "ethBalance":"1",
			  "linkBalance":"1",
			  "nextNonce":1,
			  "lastUsed":"2000-01-01T00:00:00Z",
			  "isFunding":true,
			  "createdAt":"2000-01-01T00:00:00Z",
			  "updatedAt":"2000-01-01T00:00:00Z",
			  "deletedAt":null
		   }
		}
	 }
	`

	assert.JSONEq(t, expected, string(b))
}
