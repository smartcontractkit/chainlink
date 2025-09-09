package seqs

import (
	"fmt"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func Test_SeqDeploySolTokens(t *testing.T) {
	t.Parallel()

	var (
		chainSelector = chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	)

	// Boots up a Solana testing chain in a container. This is done outside of the tests to
	// avoid booting up the container for each test.
	chains := cldf_chain.NewBlockChainsFromSlice(memory.NewMemoryChainsSol(t, 1, "")).SolanaChains()

	tests := []struct {
		name               string
		skipGenerateChains bool // Used to avoid generating chains which boots up containers
		give               SeqDeploySolTokensInput
		wantErr            string
	}{
		{
			name: "valid input",
			give: SeqDeploySolTokensInput{
				ChainSelectors: []uint64{chainSelector},
			},
		},
		{
			name:               "error: failed to get family",
			skipGenerateChains: true,
			give: SeqDeploySolTokensInput{
				ChainSelectors: []uint64{1},
			},
			wantErr: "unknown chain selector 1",
		},
		{
			name:               "error: not a Solana chain",
			skipGenerateChains: true,
			give: SeqDeploySolTokensInput{
				ChainSelectors: []uint64{
					chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector,
				},
			},
			wantErr: fmt.Sprintf(
				"chain selector %d is not a Solana chain",
				chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				ab = cldf.NewMemoryAddressBook()
				ds = datastore.NewMemoryDataStore()

				b    = optest.NewBundle(t)
				deps = SeqDeploySolTokensDeps{
					SolChains: chains,
					AddrBook:  ab,
					Datastore: ds,
				}
			)

			got, err := operations.ExecuteSequence(b, SeqDeploySolTokens, deps, tt.give)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)

				// Check that the output has the address
				require.Len(t, got.Output.Addresses, len(tt.give.ChainSelectors))

				// Check that the address book has the link token contract for each chain
				for _, csel := range tt.give.ChainSelectors {
					addrBookByChain, err := ab.AddressesForChain(csel)
					require.NoError(t, err)
					require.NotEmpty(t, addrBookByChain)
					require.Len(t, addrBookByChain, 1)
				}

				// Check the address book has the link token contract for each chain
				addrRefs, err := ds.Addresses().Fetch()
				require.NoError(t, err)
				require.Len(t, addrRefs, 1)
			}
		})
	}
}
