package stellar

import (
	"testing"

	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	types2 "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/deployment"
)

func testStellarKey(seed byte) []byte {
	k := make([]byte, 32)
	for i := range k {
		k[i] = seed + byte(i)
	}
	return k
}

func testStellarNode(id string, cfgs map[uint64][]byte) deployment.Node {
	n := deployment.Node{
		NodeID:         id,
		SelToOCRConfig: map[chainselectors.ChainDetails]deployment.OCRConfig{},
	}
	for sel, key := range cfgs {
		n.SelToOCRConfig[chainselectors.ChainDetails{ChainSelector: sel}] = deployment.OCRConfig{
			OnchainPublicKey: types2.OnchainPublicKey(key),
		}
	}
	return n
}

func TestStellarSigners(t *testing.T) {
	t.Parallel()

	testnet := chainselectors.STELLAR_TESTNET.Selector
	mainnet := chainselectors.STELLAR_MAINNET.Selector
	ethMainnet := chainselectors.ETHEREUM_MAINNET.Selector

	keyA := testStellarKey(1)
	keyB := testStellarKey(2)

	var wantA, wantB [32]byte
	copy(wantA[:], keyA)
	copy(wantB[:], keyB)

	t.Run("prefers the config registered for the target chain", func(t *testing.T) {
		t.Parallel()

		nodes := deployment.Nodes{testStellarNode("n1", map[uint64][]byte{testnet: keyA, mainnet: keyB})}

		got, err := stellarSigners(nodes, mainnet)
		require.NoError(t, err)
		require.Equal(t, [][32]byte{wantB}, got)

		got, err = stellarSigners(nodes, testnet)
		require.NoError(t, err)
		require.Equal(t, [][32]byte{wantA}, got)
	})

	t.Run("falls back to another stellar chain config with the same key", func(t *testing.T) {
		t.Parallel()

		nodes := deployment.Nodes{testStellarNode("n1", map[uint64][]byte{mainnet: keyA, ethMainnet: testStellarKey(9)[:20]})}

		got, err := stellarSigners(nodes, testnet)
		require.NoError(t, err)
		require.Equal(t, [][32]byte{wantA}, got)
	})

	t.Run("rejects ambiguous fallback across stellar chains", func(t *testing.T) {
		t.Parallel()

		localnet := chainselectors.STELLAR_LOCALNET.Selector
		nodes := deployment.Nodes{testStellarNode("n1", map[uint64][]byte{mainnet: keyA, localnet: keyB})}

		_, err := stellarSigners(nodes, testnet)
		require.ErrorContains(t, err, "different key bundles")
	})

	t.Run("errors when a node has no stellar config", func(t *testing.T) {
		t.Parallel()

		nodes := deployment.Nodes{
			testStellarNode("n1", map[uint64][]byte{testnet: keyA}),
			testStellarNode("n2", map[uint64][]byte{ethMainnet: testStellarKey(9)[:20]}),
		}

		_, err := stellarSigners(nodes, testnet)
		require.ErrorContains(t, err, "no stellar OCR2 config for node n2")
	})

	t.Run("rejects a key that is not 32 bytes", func(t *testing.T) {
		t.Parallel()

		nodes := deployment.Nodes{testStellarNode("n1", map[uint64][]byte{testnet: keyA[:31]})}

		_, err := stellarSigners(nodes, testnet)
		require.ErrorContains(t, err, "expected 32-byte stellar onchain public key")
	})

	t.Run("skips bootstrap nodes and keeps slice order", func(t *testing.T) {
		t.Parallel()

		bt := testStellarNode("bt", nil)
		bt.IsBootstrap = true
		nodes := deployment.Nodes{
			testStellarNode("n1", map[uint64][]byte{testnet: keyB}),
			bt,
			testStellarNode("n2", map[uint64][]byte{testnet: keyA}),
		}

		got, err := stellarSigners(nodes, testnet)
		require.NoError(t, err)
		require.Equal(t, [][32]byte{wantB, wantA}, got)
	})

	t.Run("errors when every node is a bootstrapper", func(t *testing.T) {
		t.Parallel()

		bt := testStellarNode("bt", nil)
		bt.IsBootstrap = true

		_, err := stellarSigners(deployment.Nodes{bt}, testnet)
		require.ErrorContains(t, err, "no stellar signers resolved")
	})
}
