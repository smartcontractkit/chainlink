package solana_test

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestSaveExistingCCIP(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		SolChains:  1,
		Nodes:      4,
	})
	solChain := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilySolana))[0]
	solAddr1 := solana.NewWallet().PublicKey().String()
	solAddr2 := solana.NewWallet().PublicKey().String()
	cfg := commonchangeset.ExistingContractsConfig{
		ExistingContracts: []commonchangeset.Contract{
			{
				Address:        solAddr1,
				TypeAndVersion: cldf.NewTypeAndVersion(shared.Router, deployment.Version1_0_0),
				ChainSelector:  solChain,
			},
			{
				Address:        solAddr2,
				TypeAndVersion: cldf.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_0_0),
				ChainSelector:  solChain,
			},
		},
	}

	output, err := commonchangeset.SaveExistingContractsChangeset(e, cfg)
	require.NoError(t, err)
	err = e.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err)
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)
	require.Equal(t, state.SolChains[solChain].Router.String(), solAddr1)
	require.Equal(t, state.SolChains[solChain].LinkToken.String(), solAddr2)
}

func TestSaveExisting(t *testing.T) {
	t.Parallel()
	dummyEnv := cldf.Environment{
		Name:              "dummy",
		Logger:            logger.TestLogger(t),
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		BlockChains: chain.NewBlockChains(
			map[uint64]chain.BlockChain{
				chainsel.SOLANA_DEVNET.Selector: cldf_solana.Chain{},
			}),
	}
	ExistingContracts := commonchangeset.ExistingContractsConfig{
		ExistingContracts: []commonchangeset.Contract{
			{
				Address: solana.NewWallet().PublicKey().String(),
				TypeAndVersion: cldf.TypeAndVersion{
					Type:    "dummy3",
					Version: deployment.Version1_1_0,
				},
				ChainSelector: chainsel.SOLANA_DEVNET.Selector,
			},
		},
	}

	output, err := commonchangeset.SaveExistingContractsChangeset(dummyEnv, ExistingContracts)
	require.NoError(t, err)
	require.NoError(t, dummyEnv.ExistingAddresses.Merge(output.AddressBook))
	addresses, err := dummyEnv.ExistingAddresses.Addresses()
	require.NoError(t, err)
	require.Len(t, addresses, 1)
	addressForSolana, exists := addresses[chainsel.SOLANA_DEVNET.Selector]
	require.True(t, exists)
	require.Len(t, addressForSolana, 1)
}
