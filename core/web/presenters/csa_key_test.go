package presenters

import (
	"fmt"
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSAKeyResource(t *testing.T) {
	timestamp := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	key, err := csakey.New("passphrase", utils.FastScryptParams)
	require.NoError(t, err)
	key.ID = 1
	key.CreatedAt = timestamp
	key.UpdatedAt = timestamp

	r := NewCSAKeyResource(*key)
	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := fmt.Sprintf(`
	{
		"data":{
			"type":"csaKeys",
			"id":"1",
			"attributes":{
				"publicKey": "%s",
				"createdAt":"2000-01-01T00:00:00Z",
				"updatedAt":"2000-01-01T00:00:00Z"
			}
		}
	}`, key.PublicKey.String())

	assert.JSONEq(t, expected, string(b))
}
