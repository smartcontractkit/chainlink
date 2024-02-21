package presenters

import (
	"fmt"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/cosmostest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/solanatest"
)

func TestSolanaMessageResource(t *testing.T) {
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

func TestCosmosMessageResource(t *testing.T) {
	id := "1"
	chainID := cosmostest.RandomChainID()
	contractID := "cosmos1p3ucd3ptpw902fluyjzkq3fflq4btddac9sa3s"
	r := NewCosmosMsgResource(id, chainID, contractID)
	assert.Equal(t, chainID, r.ChainID)
	assert.Equal(t, contractID, r.ContractID)

	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := fmt.Sprintf(`
	{
	   "data":{
		  "type":"cosmos_messages",
		  "id":"%s/%s",
		  "attributes":{
			 "ChainID":"%s",
			 "ContractID":"%s",
			 "State":"",
			 "TxHash":null
		  }
	   }
	}
	`, chainID, id, chainID, contractID)

	assert.JSONEq(t, expected, string(b))
}
