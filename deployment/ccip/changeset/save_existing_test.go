package changeset_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

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
		Nodes:      4,
	})
	chains := e.AllChainSelectors()
	chain1 := chains[0]
	chain2 := chains[1]
	cfg := commonchangeset.ExistingContractsConfig{
		ExistingContracts: []commonchangeset.Contract{
			{
				Address:        common.BigToAddress(big.NewInt(1)).String(),
				TypeAndVersion: deployment.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_0_0),
				ChainSelector:  chain1,
			},
			{
				Address:        common.BigToAddress(big.NewInt(2)).String(),
				TypeAndVersion: deployment.NewTypeAndVersion(changeset.WETH9, deployment.Version1_0_0),
				ChainSelector:  chain1,
			},
			{
				Address:        common.BigToAddress(big.NewInt(3)).String(),
				TypeAndVersion: deployment.NewTypeAndVersion(changeset.TokenAdminRegistry, deployment.Version1_5_0),
				ChainSelector:  chain1,
			},
			{
				Address:        common.BigToAddress(big.NewInt(4)).String(),
				TypeAndVersion: deployment.NewTypeAndVersion(changeset.RegistryModule, deployment.Version1_6_0),
				ChainSelector:  chain2,
			},
			{
				Address:        common.BigToAddress(big.NewInt(5)).String(),
				TypeAndVersion: deployment.NewTypeAndVersion(changeset.Router, deployment.Version1_2_0),
				ChainSelector:  chain2,
			},
		},
	}

	output, err := commonchangeset.SaveExistingContractsChangeset(e, cfg)
	require.NoError(t, err)
	err = e.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err)
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)
	require.Equal(t, state.Chains[chain1].LinkToken.Address(), common.BigToAddress(big.NewInt(1)))
	require.Equal(t, state.Chains[chain1].Weth9.Address(), common.BigToAddress(big.NewInt(2)))
	require.Equal(t, state.Chains[chain1].TokenAdminRegistry.Address(), common.BigToAddress(big.NewInt(3)))
	require.NotEmpty(t, state.Chains[chain2].RegistryModules1_6)
	require.Equal(t, state.Chains[chain2].RegistryModules1_6[0].Address(), common.BigToAddress(big.NewInt(4)))
	require.Equal(t, state.Chains[chain2].Router.Address(), common.BigToAddress(big.NewInt(5)))
}
