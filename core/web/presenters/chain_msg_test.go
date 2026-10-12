package presenters

import (
	"fmt"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/solanatest"
)

func TestSolanaMessageResource(t *testing.T) {
	t.Parallel()
	id := "1"
	chainID := solanatest.RandomChainID()
	r := NewSolanaMsgResource(id, chainID)
	assert.Equal(t, chainID, r.ChainID)

	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := fmt.Sprintf(`
	{
	   "data":{
		  "type":"solana_messages",
		  "id":"%s/%s",
		  "attributes":{
			 "ChainID":"%s",
			 "from":"",
			 "to":"",
			 "amount":0
		  }
	   }
	}
	`, chainID, id, chainID)

	assert.JSONEq(t, expected, string(b))
}
