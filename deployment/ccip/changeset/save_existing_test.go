package changeset_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

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
		Nodes:      4,
	})
	chains := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain1 := chains[0]
	chain2 := chains[1]
	cfg := commonchangeset.ExistingContractsConfig{
		ExistingContracts: []commonchangeset.Contract{
			{
				Address:        common.BigToAddress(big.NewInt(1)).String(),
				TypeAndVersion: cldf.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_0_0),
				ChainSelector:  chain1,
			},
			{
				Address:        common.BigToAddress(big.NewInt(2)).String(),
				TypeAndVersion: cldf.NewTypeAndVersion(shared.WETH9, deployment.Version1_0_0),
				ChainSelector:  chain1,
			},
			{
				Address:        common.BigToAddress(big.NewInt(3)).String(),
				TypeAndVersion: cldf.NewTypeAndVersion(shared.TokenAdminRegistry, deployment.Version1_5_0),
				ChainSelector:  chain1,
			},
			{
				Address:        common.BigToAddress(big.NewInt(4)).String(),
				TypeAndVersion: cldf.NewTypeAndVersion(shared.RegistryModule, deployment.Version1_6_0),
				ChainSelector:  chain2,
			},
			{
				Address:        common.BigToAddress(big.NewInt(5)).String(),
				TypeAndVersion: cldf.NewTypeAndVersion(shared.Router, deployment.Version1_2_0),
				ChainSelector:  chain2,
			},
		},
	}

	output, err := commonchangeset.SaveExistingContractsChangeset(e, cfg)
	require.NoError(t, err)
	err = e.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err)
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)
	chainState, _ := state.EVMChainState(chain1)
	require.Equal(t, chainState.LinkToken.Address(), common.BigToAddress(big.NewInt(1)))
	require.Equal(t, chainState.Weth9.Address(), common.BigToAddress(big.NewInt(2)))
	require.Equal(t, chainState.TokenAdminRegistry.Address(), common.BigToAddress(big.NewInt(3)))
	require.NotEmpty(t, state.MustGetEVMChainState(chain2).RegistryModules1_6)
	require.Equal(t, state.MustGetEVMChainState(chain2).RegistryModules1_6[0].Address(), common.BigToAddress(big.NewInt(4)))
	require.Equal(t, state.MustGetEVMChainState(chain2).Router.Address(), common.BigToAddress(big.NewInt(5)))
}
