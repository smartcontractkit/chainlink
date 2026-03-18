package cre

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAptosAccountsForNode_ReturnsCachedAddresses(t *testing.T) {
	t.Parallel()

	node := &Node{
		Name: "node-1",
		Addresses: Addresses{
			AptosAddresses: []string{"0x1", "0x2"},
		},
	}

	addresses, err := AptosAccountsForNode(context.Background(), node)
	require.NoError(t, err)
	require.Equal(t, []string{"0x1", "0x2"}, addresses)

	addresses[0] = "0xdead"
	require.Equal(t, []string{"0x1", "0x2"}, node.Addresses.AptosAddresses)
}

func TestAptosAccountsForNode_RequiresClientWhenCacheMissing(t *testing.T) {
	t.Parallel()

	node := &Node{
		Name: "node-1",
	}

	_, err := AptosAccountsForNode(context.Background(), node)
	require.Error(t, err)
	require.ErrorContains(t, err, "missing cached aptos addresses for node node-1 and node has no rest client")
}
