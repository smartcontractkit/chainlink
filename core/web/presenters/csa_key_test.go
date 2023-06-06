package presenters

import (
	"fmt"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestCSAKeyResource(t *testing.T) {
	key, err := csakey.New("passphrase", utils.FastScryptParams)
	require.NoError(t, err)
	key.ID = 1

	r := NewCSAKeyResource(key.ToV2())
	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := fmt.Sprintf(`
	{
		"data":{
			"type":"csaKeys",
			"id":"%s",
			"attributes":{
				"publicKey": "%s",
				"version": 1
			}
		}
	}`, key.PublicKey.String(), key.PublicKey.String())

	assert.JSONEq(t, expected, string(b))
}
