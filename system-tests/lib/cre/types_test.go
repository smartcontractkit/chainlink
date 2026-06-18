package cre

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
)

func TestNewNodesPreservesInputOrder(t *testing.T) {
	t.Parallel()

	cfgs := []NodeMetadataConfig{
		{Host: "node-0", Roles: []string{WorkerNode}, Index: 0, Keys: NodeKeyInput{Password: "password"}},
		{Host: "node-1", Roles: []string{BootstrapNode}, Index: 1, Keys: NodeKeyInput{Password: "password"}},
		{Host: "node-2", Roles: []string{GatewayNode}, Index: 2, Keys: NodeKeyInput{Password: "password"}},
	}

	nodes, err := newNodes(cfgs)

	require.NoError(t, err)
	require.Len(t, nodes, len(cfgs))
	for i, node := range nodes {
		require.Equal(t, cfgs[i].Host, node.Host)
		require.Equal(t, cfgs[i].Index, node.Index)
		require.Equal(t, cfgs[i].Roles, node.Roles)
	}
}

func TestNewNodesImportedSecretsBypassGeneration(t *testing.T) {
	t.Parallel()

	importedKeys, err := NewNodeKeys(NodeKeyInput{
		EVMChainIDs: []uint64{111},
		Password:    "password",
	})
	require.NoError(t, err)

	importedSecrets, err := importedKeys.ToNodeSecretsTOML()
	require.NoError(t, err)

	nodes, err := newNodes([]NodeMetadataConfig{{
		Host:  "node-imported",
		Roles: []string{WorkerNode},
		Index: 0,
		Keys: NodeKeyInput{
			ImportedSecrets: importedSecrets,
		},
	}})

	require.NoError(t, err)
	require.Len(t, nodes, 1)
	require.Equal(t, importedKeys.PeerID(), nodes[0].Keys.PeerID())
	require.Equal(t, importedKeys.EVM[111].PublicAddress, nodes[0].Keys.EVM[111].PublicAddress)
}

func TestNewNodesReturnsFailingIndex(t *testing.T) {
	t.Parallel()

	importedKeys, err := NewNodeKeys(NodeKeyInput{Password: "password"})
	require.NoError(t, err)

	importedSecrets, err := importedKeys.ToNodeSecretsTOML()
	require.NoError(t, err)

	_, err = newNodes([]NodeMetadataConfig{
		{
			Host:  "node-0",
			Roles: []string{WorkerNode},
			Index: 0,
			Keys:  NodeKeyInput{ImportedSecrets: importedSecrets},
		},
		{
			Host:  "node-1",
			Roles: []string{WorkerNode},
			Index: 1,
			Keys:  NodeKeyInput{ImportedSecrets: "not valid toml"},
		},
		{
			Host:  "node-2",
			Roles: []string{WorkerNode},
			Index: 2,
			Keys:  NodeKeyInput{ImportedSecrets: importedSecrets},
		},
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "index: 1")
}

var _ = secrets.NodeKeys{}
