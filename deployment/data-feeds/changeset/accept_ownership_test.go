package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestAcceptOwnership(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.AllChainSelectors()[0]
	chain := env.Chains[chainSelector]

	newEnv, err := commonChangesets.Apply(t, env, nil,
		commonChangesets.Configure(
			cldf.CreateLegacyChangeSet(commonChangesets.DeployMCMSWithTimelockV2),
			map[uint64]commonTypes.MCMSWithTimelockConfigV2{
				chainSelector: proposalutils.SingleGroupTimelockConfigV2(t),
			},
		),
	)
	require.NoError(t, err)

	timeLockAddress, err := deployment.SearchAddressBook(newEnv.ExistingAddresses, chainSelector, "RBACTimelock")
	require.NoError(t, err)

	cache, _ := DeployCache(chain, []string{})
	tx, _ := cache.Contract.TransferOwnership(chain.DeployerKey, common.HexToAddress(timeLockAddress))
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	_, err = commonChangesets.Apply(t, newEnv, nil,
		commonChangesets.Configure(
			AcceptOwnershipChangeset,
			types.AcceptOwnershipConfig{
				ChainSelector:     chainSelector,
				ContractAddresses: []common.Address{cache.Contract.Address()},
				McmsConfig: &types.MCMSConfig{
					MinDelay: 1,
				},
			},
		),
	)
	require.NoError(t, err)
}
