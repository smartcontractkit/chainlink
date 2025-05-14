package reward_manager

import (
	"testing"

	"github.com/stretchr/testify/require"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	dsTypes "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

// RewardManagerDeploy deploys RewardManager
// and returns the updated environment and the addresses of RewardManager.
func RewardManagerDeploy(
	t *testing.T,
	cfg testutil.MemoryEnv,
) (cldf.Environment, common.Address) {
	t.Helper()

	chainSelector := testutil.TestChain.Selector

	var shouldTransfer bool
	var mcmsProposalCfg *proposalutils.TimelockConfig
	if len(cfg.Timelocks) > 0 {
		shouldTransfer = true
		mcmsProposalCfg = &proposalutils.TimelockConfig{
			MinDelay: 0,
		}
	}

	// 2) Deploy RewardManager
	deployRewardManagerCfg := DeployRewardManagerConfig{
		Ownership: dsTypes.OwnershipSettings{
			ShouldTransfer:     shouldTransfer,
			MCMSProposalConfig: mcmsProposalCfg,
		},
		ChainsToDeploy: map[uint64]DeployRewardManager{
			testutil.TestChain.Selector: {LinkTokenAddress: cfg.LinkTokenState.LinkToken.Address()},
		},
	}

	env, err := commonChangesets.Apply(t, cfg.Environment, nil,
		commonChangesets.Configure(
			DeployRewardManagerChangeset,
			deployRewardManagerCfg,
		),
	)
	require.NoError(t, err, "deploying RewardManager should not fail")

	// Get the RewardManager address
	record, err := env.DataStore.Addresses().Get(ds.NewAddressRefKey(chainSelector, ds.ContractType(dsTypes.RewardManager), &deployment.Version0_5_0, ""))
	require.NoError(t, err)
	rewardManagerAddr := common.HexToAddress(record.Address)

	return env, rewardManagerAddr
}
