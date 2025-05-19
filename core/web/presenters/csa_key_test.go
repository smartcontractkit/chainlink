package presenters

import (
	"fmt"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
)

func TestCSAKeyResource(t *testing.T) {
	keyV2, err := csakey.NewV2()
	require.NoError(t, err)

	r := NewCSAKeyResource(keyV2)
	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := fmt.Sprintf(`
	{
		"data":{
			"type":"csaKeys",
			"id":"%[1]s",
			"attributes":{
				"publicKey": "csa_%[1]s",
				"version": 1
			}
		}
	}`, keyV2.PublicKeyString())

	assert.JSONEq(t, expected, string(b))
}
