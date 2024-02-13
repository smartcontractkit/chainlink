package presenters

import (
	"fmt"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestNodeResource(t *testing.T) {
	var nodeResource NodeResource
	var r interface{}
	state := "test"
	cfg := "cfg"
	testCases := []string{"solana", "cosmos", "starknet"}
	for _, tc := range testCases {
		chainID := fmt.Sprintf("%s chain ID", tc)
		nodeName := fmt.Sprintf("%s_node", tc)

		switch tc {
		case "evm":
			evmNodeResource := NewEVMNodeResource(
				types.NodeStatus{
					ChainID: chainID,
					Name:    nodeName,
					Config:  cfg,
					State:   state,
				})
			r = evmNodeResource
			nodeResource = evmNodeResource.NodeResource
		case "solana":
			solanaNodeResource := NewSolanaNodeResource(
				types.NodeStatus{
					ChainID: chainID,
					Name:    nodeName,
					Config:  cfg,
					State:   state,
				})
			r = solanaNodeResource
			nodeResource = solanaNodeResource.NodeResource
		case "cosmos":
			cosmosNodeResource := NewCosmosNodeResource(
				types.NodeStatus{
					ChainID: chainID,
					Name:    nodeName,
					Config:  cfg,
					State:   state,
				})
			r = cosmosNodeResource
			nodeResource = cosmosNodeResource.NodeResource
		case "starknet":
			starknetNodeResource := NewStarkNetNodeResource(
				types.NodeStatus{
					ChainID: chainID,
					Name:    nodeName,
					Config:  cfg,
					State:   state,
				})
			r = starknetNodeResource
			nodeResource = starknetNodeResource.NodeResource
		default:
			t.Fail()
		}
		assert.Equal(t, chainID, nodeResource.ChainID)
		assert.Equal(t, nodeName, nodeResource.Name)
		assert.Equal(t, cfg, nodeResource.Config)
		assert.Equal(t, state, nodeResource.State)

		b, err := jsonapi.Marshal(r)
		require.NoError(t, err)

		expected := fmt.Sprintf(`
		{
		  "data":{
			  "type":"%s_node",
			  "id":"%s/%s",
			  "attributes":{
				 "chainID":"%s",
				 "name":"%s",
				 "config":"%s",
				 "state":"%s"
			  }
		  }
		}
	`, tc, chainID, nodeName, chainID, nodeName, cfg, state)
		assert.JSONEq(t, expected, string(b))
	}
}
