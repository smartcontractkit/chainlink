package solana

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

func TestFromFallsBackToInputChainIDWhenOutputMissing(t *testing.T) {
	t.Setenv("SOLANA_PRIVATE_KEY", DefaultSolanaPrivateKey.String())

	contractsDir := t.TempDir()
	input := &blockchain.Input{
		ChainID:      "22222222222222222222222222222222222222222222",
		ContractsDir: contractsDir,
	}
	out := &blockchain.Output{
		Type:    blockchain.TypeSolana,
		ChainID: "",
		Family:  blockchain.FamilySolana,
		Nodes: []*blockchain.Node{
			{ExternalHTTPUrl: "http://localhost:8550"},
		},
	}

	got, err := From(input, out)
	require.NoError(t, err, "expected reconstruction to use input chain id fallback")
	require.Equal(t, input.ChainID, got.SolanaChainID, "expected fallback chain id to be retained")
}

