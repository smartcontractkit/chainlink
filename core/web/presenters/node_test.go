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
	state := "test"
	cfg := "cfg"
	testCases := []string{"solana", "cosmos", "starknet", "tron"}
	for _, tc := range testCases {
		chainID := tc + " chain ID"
		nodeName := tc + "_node"

		nodeResource = NewNodeResource(types.NodeStatus{
			ChainID: chainID,
			Name:    nodeName,
			Config:  cfg,
			State:   state,
		})

		assert.Equal(t, chainID, nodeResource.ChainID)
		assert.Equal(t, nodeName, nodeResource.Name)
		assert.Equal(t, cfg, nodeResource.Config)
		assert.Equal(t, state, nodeResource.State)

		b, err := jsonapi.Marshal(nodeResource)
		require.NoError(t, err)

		expected := fmt.Sprintf(`
		{
		  "data":{
			  "type":"node",
			  "id":"%s/%s",
			  "attributes":{
				 "chainID":"%s",
				 "name":"%s",
				 "config":"%s",
				 "state":"%s"
			  }
		  }
		}
	`, chainID, nodeName, chainID, nodeName, cfg, state)
		assert.JSONEq(t, expected, string(b))
	}
}
