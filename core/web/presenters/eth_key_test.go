package presenters

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestETHKeyResource(t *testing.T) {
	var (
		now        = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		nextNonce  = int64(1)
		addressStr = "0x2aCFF2ec69aa9945Ed84f4F281eCCF6911A3B0eD"
	)
	address, err := ethkey.NewEIP55Address(addressStr)
	require.NoError(t, err)

	key := ethkey.Key{
		ID:        1,
		Address:   address,
		CreatedAt: now,
		UpdatedAt: now,
		NextNonce: nextNonce,
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

	expected := fmt.Sprintf(`
	{
		"data":{
		   "type":"eTHKeys",
		   "id":"%s",
		   "attributes":{
			  "address":"%s",
			  "ethBalance":"1",
			  "linkBalance":"1",
			  "nextNonce":1,
			  "isFunding":true,
			  "createdAt":"2000-01-01T00:00:00Z",
			  "updatedAt":"2000-01-01T00:00:00Z",
			  "deletedAt":null
		   }
		}
	 }
	`, addressStr, addressStr)

	assert.JSONEq(t, expected, string(b))

	// With a deleted field
	key.DeletedAt = gorm.DeletedAt(sql.NullTime{Time: now, Valid: true})

	r, err = NewETHKeyResource(key,
		SetETHKeyEthBalance(assets.NewEth(1)),
		SetETHKeyLinkBalance(assets.NewLink(1)),
	)
	require.NoError(t, err)
	b, err = jsonapi.Marshal(r)
	require.NoError(t, err)

	expected = fmt.Sprintf(`
	{
		"data": {
			"type":"eTHKeys",
			"id":"%s",
			"attributes":{
				"address":"%s",
				"ethBalance":"1",
				"linkBalance":"1",
				"nextNonce":1,
				"isFunding":true,
				"createdAt":"2000-01-01T00:00:00Z",
				"updatedAt":"2000-01-01T00:00:00Z",
				"deletedAt":"2000-01-01T00:00:00Z"
			}
		}
	}`,
		addressStr, addressStr,
	)

	assert.JSONEq(t, expected, string(b))
}
