package presenters

import (
	"net/url"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func TestBridgeResource(t *testing.T) {
	t.Parallel()

	timestamp := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	url, err := url.Parse("https://bridge.example.com/api")
	require.NoError(t, err)

	bridge := bridges.BridgeType{
		Name:                   "test",
		URL:                    models.WebURL(*url),
		Confirmations:          1,
		OutgoingToken:          "vjNL7X8Ea6GFJoa6PBsvK2ECzNK3b8IZ",
		MinimumContractPayment: assets.NewLinkFromJuels(1),
		CreatedAt:              timestamp,
	}

	r := NewBridgeResource(bridge)

	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := `
{
	"data": {
		"type":"bridges",
		"id":"test",
		"attributes":{
			"name":"test",
			"url":"https://bridge.example.com/api",
			"confirmations":1,
			"outgoingToken":"vjNL7X8Ea6GFJoa6PBsvK2ECzNK3b8IZ",
			"minimumContractPayment":"1",
			"createdAt":"2000-01-01T00:00:00Z"
		}
	}
}
`

	assert.JSONEq(t, expected, string(b))

	// Test insertion of IncomingToken
	r.IncomingToken = "cd+OfGXy3UHEDAlD0y27F6/rJE14X1UI"
	b, err = jsonapi.Marshal(r)
	require.NoError(t, err)

	expected = `
{
	"data": {
		"type":"bridges",
		"id":"test",
		"attributes":{
			"name":"test",
			"url":"https://bridge.example.com/api",
			"confirmations":1,
			"incomingToken": "cd+OfGXy3UHEDAlD0y27F6/rJE14X1UI",
			"outgoingToken":"vjNL7X8Ea6GFJoa6PBsvK2ECzNK3b8IZ",
			"minimumContractPayment":"1",
			"createdAt":"2000-01-01T00:00:00Z"
		}
	}
}
`

	assert.JSONEq(t, expected, string(b))
}
