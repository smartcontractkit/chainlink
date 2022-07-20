package presenters

import (
	"fmt"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestETHKeyResource(t *testing.T) {
	var (
		now        = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		nextNonce  = int64(1)
		addressStr = "0x2aCFF2ec69aa9945Ed84f4F281eCCF6911A3B0eD"
	)
	address, err := ethkey.NewEIP55Address(addressStr)
	require.NoError(t, err)

	key := ethkey.KeyV2{
		Address: address,
	}

	state := ethkey.State{
		ID:         1,
		EVMChainID: *utils.NewBigI(42),
		Address:    address,
		CreatedAt:  now,
		UpdatedAt:  now,
		NextNonce:  nextNonce,
		IsFunding:  true,
	}

	r, err := NewETHKeyResource(key, state,
		SetETHKeyEthBalance(assets.NewEth(1)),
		SetETHKeyLinkBalance(assets.NewLinkFromJuels(1)),
		SetETHKeyMaxGasPriceWei(*utils.NewBigI(12345)),
	)
	require.NoError(t, err)

	assert.Equal(t, assets.NewEth(1), r.EthBalance)
	assert.Equal(t, assets.NewLinkFromJuels(1), r.LinkBalance)
	assert.Equal(t, *utils.NewBigI(12345), r.MaxGasPriceWei)

	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := fmt.Sprintf(`
	{
		"data":{
		   "type":"eTHKeys",
		   "id":"%s",
		   "attributes":{
			  "address":"%s",
			  "evmChainID":"42",
			  "ethBalance":"1",
			  "linkBalance":"1",
			  "isFunding":true,
			  "createdAt":"2000-01-01T00:00:00Z",
			  "updatedAt":"2000-01-01T00:00:00Z",
			  "maxGasPriceWei":"12345"
		   }
		}
	 }
	`, addressStr, addressStr)

	assert.JSONEq(t, expected, string(b))

	r, err = NewETHKeyResource(key, state,
		SetETHKeyEthBalance(assets.NewEth(1)),
		SetETHKeyLinkBalance(assets.NewLinkFromJuels(1)),
		SetETHKeyMaxGasPriceWei(*utils.NewBigI(12345)),
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
			  	"evmChainID":"42",
				"ethBalance":"1",
				"linkBalance":"1",
				"isFunding":true,
				"createdAt":"2000-01-01T00:00:00Z",
				"updatedAt":"2000-01-01T00:00:00Z",
				"maxGasPriceWei":"12345"
			}
		}
	}`,
		addressStr, addressStr,
	)

	assert.JSONEq(t, expected, string(b))
}
