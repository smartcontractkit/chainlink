package changesets

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink/deployment/tokens/internal/ops"
)

func Test_DeployEVMLinkTokens_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	var (
		csel    = chain_selectors.TEST_1000.Selector
		ethAddr = "0xeC91988D7dD84d8adE801b739172ad15c860A700"
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
					cldf_evm.Chain{
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
					chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector,
				}, // Uses an invalid Solana chain selector
			},
			wantErr: fmt.Sprintf(
				"chain selector %d is not in the evm family",
				chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector,
			),
		},
		{
			name: "error: link token contracts exists in address book",
			beforeFunc: func(t *testing.T, e *cldf.Environment) {
				t.Helper()

				e.ExistingAddresses = cldf.NewMemoryAddressBook()
				err := e.ExistingAddresses.Save(csel, ethAddr, ops.LinkTokenTypeAndVersion1)
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
					Address:       "0xeC91988D7dD84d8adE801b739172ad15c860A700",
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
				cs  = deployEVMLinkTokens{}
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

func Test_DeployEVMLinkTokens_Apply(t *testing.T) {
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

			rt, err := runtime.New(t.Context(),
				runtime.WithEnvOpts(
					environment.WithEVMSimulatedN(t, 1),
				),
			)
			require.NoError(t, err)

			err = rt.Exec(
				runtime.ChangesetTask(DeployEVMLinkTokens, tt.giveFunc(rt.Environment())),
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
