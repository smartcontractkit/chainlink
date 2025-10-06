package changesets

import (
	"fmt"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink/deployment/tokens/internal/ops"
)

func Test_DeploySolLinkTokens_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	var (
		csel    = chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
		solAddr = "J6oVJ42pE6eXdTCcCidhjzHWS7Sxz6yMsXHxXphT1U7Y"
	)

	tests := []struct {
		name       string
		beforeFunc func(t *testing.T, e *cldf.Environment)
		input      DeployLinkTokensInput
		wantErr    string
	}{
		{
			name: "valid input",
			beforeFunc: func(t *testing.T, e *cldf.Environment) {
				e.BlockChains = cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{
					cldf_solana.Chain{
						Selector: csel,
					},
				})

				// Inject empty address book and datastore
				e.ExistingAddresses = cldf.NewMemoryAddressBook()
				e.DataStore = datastore.NewMemoryDataStore().Seal()
			},
			input: DeployLinkTokensInput{
				ChainSelectors: []uint64{csel},
			},
		},
		{
			name: "error: duplicate chain selectors",
			input: DeployLinkTokensInput{
				ChainSelectors: []uint64{1, 1},
			},
			wantErr: "duplicate chain selector found",
		},
		{
			name: "error: invalid chain selector family",
			input: DeployLinkTokensInput{
				ChainSelectors: []uint64{
					chain_selectors.TEST_1000.Selector,
				}, // Uses an invalid EVM chain selector
			},
			wantErr: fmt.Sprintf(
				"chain selector %d is not in the solana family",
				chain_selectors.TEST_1000.Selector,
			),
		},
		{
			name: "error: link token contracts exists in address book",
			beforeFunc: func(t *testing.T, e *cldf.Environment) {
				t.Helper()

				e.ExistingAddresses = cldf.NewMemoryAddressBook()
				err := e.ExistingAddresses.Save(csel, solAddr, ops.LinkTokenTypeAndVersion1)
				require.NoError(t, err)
			},
			input: DeployLinkTokensInput{
				ChainSelectors: []uint64{csel},
			},
			wantErr: "link token contract already exists for chain selector",
		},
		{
			name: "error: link token contract exists in datastore",
			beforeFunc: func(t *testing.T, e *cldf.Environment) {
				t.Helper()

				// Insert the selector with no addresses to pass address book check
				e.ExistingAddresses = cldf.NewMemoryAddressBookFromMap(
					map[uint64]map[string]cldf.TypeAndVersion{
						csel: {},
					},
				)

				ds := datastore.NewMemoryDataStore()
				err := ds.Addresses().Add(datastore.AddressRef{
					ChainSelector: csel,
					Address:       solAddr,
					Type:          datastore.ContractType(ops.LinkTokenTypeAndVersion1.Type.String()),
					Version:       &ops.LinkTokenTypeAndVersion1.Version,
				})
				require.NoError(t, err)

				e.DataStore = ds.Seal()
			},
			input: DeployLinkTokensInput{
				ChainSelectors: []uint64{csel},
			},
			wantErr: "link token contract already exists for chain selector",
		},
		{
			name: "error: chain selector not found in environment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				cs  = deploySolLinkTokens{}
				env = &cldf.Environment{}
			)

			if tt.beforeFunc != nil {
				tt.beforeFunc(t, env)
			}

			err := cs.VerifyPreconditions(*env, tt.input)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_DeploySolLinkTokens_Apply(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		giveFunc func(e cldf.Environment) DeployLinkTokensInput
	}{
		{
			name: "valid input",
			giveFunc: func(e cldf.Environment) DeployLinkTokensInput {
				csels := e.BlockChains.ListChainSelectors()

				return DeployLinkTokensInput{
					ChainSelectors: csels,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			programsPath := t.TempDir() // We don't need programs in here to deploy the token

			rt, err := runtime.New(t.Context(),
				runtime.WithEnvOpts(
					environment.WithSolanaContainerN(t, 1, programsPath, map[string]string{}),
				),
			)
			require.NoError(t, err)

			err = rt.Exec(
				runtime.ChangesetTask(DeploySolLinkTokens, tt.giveFunc(rt.Environment())),
			)
			require.NoError(t, err)

			// Check that the address book has the link token contract for each chain
			for _, csel := range rt.Environment().BlockChains.ListChainSelectors() {
				addrBookByChain, err := rt.State().AddressBook.AddressesForChain(csel)
				require.NoError(t, err)
				require.NotEmpty(t, addrBookByChain)
				require.Len(t, addrBookByChain, 1)
			}

			// Check the address book has the link token contract for each chain
			addrRefs, err := rt.State().DataStore.Addresses().Fetch()
			require.NoError(t, err)
			require.Len(t, addrRefs, 1)
		})
	}
}
