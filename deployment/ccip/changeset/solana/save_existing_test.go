package solana_test

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
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
	solChain := e.AllChainSelectorsSolana()[0]
	solAddr1 := solana.NewWallet().PublicKey().String()
	solAddr2 := solana.NewWallet().PublicKey().String()
	cfg := commonchangeset.ExistingContractsConfig{
		ExistingContracts: []commonchangeset.Contract{
			{
				Address:        solAddr1,
				TypeAndVersion: deployment.NewTypeAndVersion(changeset.Router, deployment.Version1_0_0),
				ChainSelector:  solChain,
			},
			{
				Address:        solAddr2,
				TypeAndVersion: deployment.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_0_0),
				ChainSelector:  solChain,
			},
		},
	}

	output, err := commonchangeset.SaveExistingContractsChangeset(e, cfg)
	require.NoError(t, err)
	err = e.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err)
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)
	require.Equal(t, state.SolChains[solChain].Router.String(), solAddr1)
	require.Equal(t, state.SolChains[solChain].LinkToken.String(), solAddr2)
}

func TestSaveExisting(t *testing.T) {
	t.Parallel()
	dummyEnv := deployment.Environment{
		Name:              "dummy",
		Logger:            logger.TestLogger(t),
		ExistingAddresses: deployment.NewMemoryAddressBook(),
		SolChains: map[uint64]deployment.SolChain{
			chainsel.SOLANA_DEVNET.Selector: {},
		},
	}
	ExistingContracts := commonchangeset.ExistingContractsConfig{
		ExistingContracts: []commonchangeset.Contract{
			{
				Address: solana.NewWallet().PublicKey().String(),
				TypeAndVersion: deployment.TypeAndVersion{
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
