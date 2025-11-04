package solana_test

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestSaveExistingCCIP(t *testing.T) {
	t.Parallel()

	e, err := environment.New(t.Context())
	require.NoError(t, err)

	// Insert a fake Solana chain into the environment since we are not making any onchain calls.
	selector := chainsel.TEST_22222222222222222222222222222222222222222222.Selector
	e.BlockChains = cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{
		cldf_solana.Chain{Selector: selector},
	})

	solAddr1 := solana.NewWallet().PublicKey().String()
	solAddr2 := solana.NewWallet().PublicKey().String()
	cfg := commonchangeset.ExistingContractsConfig{
		ExistingContracts: []commonchangeset.Contract{
			{
				Address:        solAddr1,
				TypeAndVersion: cldf.NewTypeAndVersion(shared.Router, deployment.Version1_0_0),
				ChainSelector:  selector,
			},
			{
				Address:        solAddr2,
				TypeAndVersion: cldf.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_0_0),
				ChainSelector:  selector,
			},
		},
	}

	output, err := commonchangeset.SaveExistingContractsChangeset(*e, cfg)
	require.NoError(t, err)
	err = e.ExistingAddresses.Merge(output.AddressBook) //nolint:staticcheck // AddressBook is deprecated but still in use for this changeset
	require.NoError(t, err)
	state, err := stateview.LoadOnchainState(*e)
	require.NoError(t, err)
	require.Equal(t, state.SolChains[selector].Router.String(), solAddr1)
	require.Equal(t, state.SolChains[selector].LinkToken.String(), solAddr2)
}

func TestSaveExisting(t *testing.T) {
	t.Parallel()

	e, err := environment.New(t.Context())
	require.NoError(t, err)

	// Insert a fake Solana chain into the environment since we are not making any onchain calls.
	selector := chainsel.TEST_22222222222222222222222222222222222222222222.Selector
	e.BlockChains = cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{
		cldf_solana.Chain{Selector: selector},
	})

	ExistingContracts := commonchangeset.ExistingContractsConfig{
		ExistingContracts: []commonchangeset.Contract{
			{
				Address: solana.NewWallet().PublicKey().String(),
				TypeAndVersion: cldf.TypeAndVersion{
					Type:    "dummy3",
					Version: deployment.Version1_1_0,
				},
				ChainSelector: selector,
			},
		},
	}

	output, err := commonchangeset.SaveExistingContractsChangeset(*e, ExistingContracts)
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook)) //nolint:staticcheck // AddressBook is deprecated but still in use for this changeset
	addresses, err := e.ExistingAddresses.Addresses()
	require.NoError(t, err)
	require.Len(t, addresses, 1)
	addressForSolana, exists := addresses[selector]
	require.True(t, exists)
	require.Len(t, addressForSolana, 1)
}
