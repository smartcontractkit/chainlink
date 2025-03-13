package v1_6_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/stretchr/testify/require"
)

func mustHaveOwner(t *testing.T, ownable commonchangeset.Ownable, expectedOwner string) {
	owner, err := ownable.Owner(nil)
	require.NoError(t, err, "must get owner")
	require.Equal(t, expectedOwner, owner.Hex(), fmt.Sprintf("owner must be %s", expectedOwner))
}

func TestConnectNewChain(t *testing.T) {
	t.Parallel()

	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
		testCfg.Chains = 3
	})
	e := deployedEnvironment.Env

	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err, "must load onchain state")

	selectors := e.AllChainSelectors()
	var newSelector uint64
	remoteChainSelectors := make([]uint64, 0, len(selectors)-1)
	for _, selector := range selectors {
		if selector != deployedEnvironment.HomeChainSel && newSelector == 0 {
			newSelector = selector // Just take any non-home chain selector
			continue
		}
		remoteChainSelectors = append(remoteChainSelectors, selector)
	}

	timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts, len(selectors))
	for _, selector := range selectors {
		// Assemble map of addresses required for Timelock scheduling & execution
		timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.Chains[selector].Timelock,
			CallProxy: state.Chains[selector].CallProxy,
		}
	}

	// onRamp, offRamp, and router on non-new chains are assumed to be owned by the timelock
	contractsToTransfer := make(map[uint64][]common.Address, len(remoteChainSelectors))
	for _, selector := range remoteChainSelectors {
		contractsToTransfer[selector] = []common.Address{
			state.Chains[selector].OnRamp.Address(),
			state.Chains[selector].OffRamp.Address(),
			state.Chains[selector].Router.Address(),
		}
	}
	e, err = commonchangeset.Apply(t, e, timelockContracts,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commoncs.TransferToMCMSWithTimelock),
			commoncs.TransferToMCMSWithTimelockConfig{
				ContractsByChain: contractsToTransfer,
				MinDelay:         0 * time.Second,
			},
		),
	)
	require.NoError(t, err, "must apply TransferToMCMSWithTimelock")

	testRouter := false
	e, err = commonchangeset.Apply(t, e, timelockContracts,
		commonchangeset.Configure(
			v1_6.ConnectNewChainChangeset,
			v1_6.ConnectNewChainConfig{
				TestRouter:           &testRouter,
				RemoteChainSelectors: remoteChainSelectors,
				NewChainSelector:     newSelector,
				MCMSConfig: &changeset.MCMSConfig{
					MinDelay:   0 * time.Second,
					MCMSAction: timelock.Schedule,
				},
			},
		),
	)
	require.NoError(t, err, "must apply ActivateChainOnProd")

	state, err = changeset.LoadOnchainState(e)
	require.NoError(t, err, "must load onchain state")

	for _, selector := range selectors {
		if selector == newSelector {
			// New chain must have all contracts owned by timelock
			mustHaveOwner(t, state.Chains[selector].OnRamp, state.Chains[selector].Timelock.Address().Hex())
			mustHaveOwner(t, state.Chains[selector].OffRamp, state.Chains[selector].Timelock.Address().Hex())
			mustHaveOwner(t, state.Chains[selector].FeeQuoter, state.Chains[selector].Timelock.Address().Hex())
			mustHaveOwner(t, state.Chains[selector].RMNProxy, state.Chains[selector].Timelock.Address().Hex())
			mustHaveOwner(t, state.Chains[selector].NonceManager, state.Chains[selector].Timelock.Address().Hex())
			mustHaveOwner(t, state.Chains[selector].TokenAdminRegistry, state.Chains[selector].Timelock.Address().Hex())
			mustHaveOwner(t, state.Chains[selector].Router, state.Chains[selector].Timelock.Address().Hex())
			mustHaveOwner(t, state.Chains[selector].RMNRemote, state.Chains[selector].Timelock.Address().Hex())

			// New chain must have other chains supported
			for _, remoteChainSelector := range remoteChainSelectors {
				destChainConfig, err := state.Chains[selector].OnRamp.GetDestChainConfig(nil, remoteChainSelector)
				require.NoError(t, err, "must get dest chain config from onRamp")
				require.Equal(t, state.Chains[selector].Router.Address().Hex(), destChainConfig.Router.Hex(), "router must be the prod router")

				srcChainConfig, err := state.Chains[selector].OffRamp.GetSourceChainConfig(nil, remoteChainSelector)
				require.NoError(t, err, "must get src chain config from offRamp")
				require.True(t, srcChainConfig.IsEnabled, "src chain config must be enabled")
				require.True(t, srcChainConfig.IsRMNVerificationDisabled, "src chain config must have rmn verification disabled") // TODO: Fix this, should be able to accept as input
				require.Equal(t, common.LeftPadBytes(state.Chains[remoteChainSelector].OnRamp.Address().Bytes(), 32), srcChainConfig.OnRamp, "remote onRamp must be set on offRamp")
				require.Equal(t, state.Chains[selector].Router.Address().Hex(), srcChainConfig.Router.Hex(), "router must be the prod router")

				// TODO: Router checks
			}
		} else {
			// Supported remote chains
		}
	}
}
