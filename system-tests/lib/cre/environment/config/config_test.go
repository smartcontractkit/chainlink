package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func TestBlockchainNormalizeAndValidate(t *testing.T) {
	b := &Blockchain{Input: blockchain.Input{Type: blockchain.TypeAnvil}}
	b.Normalize()
	require.Equal(t, PlacementLocal, b.Placement)
	require.Equal(t, RemoteStartPolicyReuseIfIdentical, b.RemoteStartPolicy)
	require.NoError(t, b.Validate())

	b = &Blockchain{Input: blockchain.Input{Type: blockchain.TypeAnvil}, Placement: ComponentPlacement("invalid")}
	err := b.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid blockchain placement")
}

func TestJobDistributorNormalizeAndValidate(t *testing.T) {
	j := &JobDistributor{Input: jd.Input{}}
	j.Normalize()
	require.Equal(t, PlacementLocal, j.Placement)
	require.Equal(t, RemoteStartPolicyReuseIfIdentical, j.RemoteStartPolicy)
	require.NoError(t, j.Validate())

	j = &JobDistributor{Input: jd.Input{}, Placement: PlacementRemote, RemoteStartPolicy: RemoteStartPolicy("bad")}
	err := j.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid jd remote_start_policy")
}

func TestNodeSetPlacementNormalizeAndValidate(t *testing.T) {
	nodeSet := &cre.NodeSet{}
	normalizeNodeSetPlacement(nodeSet)
	require.Equal(t, string(PlacementLocal), nodeSet.Placement)
	require.Equal(t, string(RemoteStartPolicyReuseIfIdentical), nodeSet.RemoteStartPolicy)
	require.NoError(t, validateNodeSetPlacement(nodeSet))

	nodeSet.Placement = "bad"
	err := validateNodeSetPlacement(nodeSet)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid nodeset placement")
}

func TestResolveBlockchainInputs(t *testing.T) {
	_, err := ResolveBlockchainInputs(nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "at least one blockchain")

	out, err := ResolveBlockchainInputs([]*Blockchain{
		{Input: blockchain.Input{Type: blockchain.TypeAnvil}, Placement: PlacementRemote},
	})
	require.NoError(t, err)
	require.Len(t, out, 1)
	require.Equal(t, blockchain.TypeAnvil, out[0].Type)
}

func TestRemoveChainIDFromFlag(t *testing.T) {
	require.Equal(t, "write-evm", removeChainIDFromFlag("write-evm-1337"))
	require.Equal(t, "write-evm-mainnet", removeChainIDFromFlag("write-evm-mainnet"))
	require.Equal(t, "cron", removeChainIDFromFlag("cron"))
}

func TestNormalizeComponentPlacement(t *testing.T) {
	require.Equal(t, PlacementLocal, normalizeComponentPlacement(ComponentPlacement(" LOCAL ")))
	require.Equal(t, PlacementRemote, normalizeComponentPlacement(ComponentPlacement("REMOTE")))
	require.Equal(t, ComponentPlacement("weird"), normalizeComponentPlacement(ComponentPlacement("weird")))
}
