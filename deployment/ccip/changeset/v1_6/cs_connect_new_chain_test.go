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

type test struct {
	Msg                        string
	TransferRemoteChainsToMCMS bool
	TestRouter                 bool
	MCMS                       *changeset.MCMSConfig
	ErrStr                     string
}

func TestConnectNewChain(t *testing.T) {
	t.Parallel()

	mcmsConfig := &changeset.MCMSConfig{
		MinDelay:   0 * time.Second,
		MCMSAction: timelock.Schedule,
	}

	tests := []test{
		{
			Msg:                        "MCMS config should be defined",
			TransferRemoteChainsToMCMS: true,
			MCMS:                       nil,
			ErrStr:                     "not owned by deployer key",
		},
		{
			Msg:                        "MCMS config should be undefined",
			TransferRemoteChainsToMCMS: false,
			MCMS:                       mcmsConfig,
			ErrStr:                     "not owned by timelock",
		},
		{
			Msg:                        "Use production router (with MCMS)",
			TransferRemoteChainsToMCMS: true,
			TestRouter:                 false,
			MCMS:                       mcmsConfig,
		},
		{
			Msg:                        "Use production router (without MCMS)",
			TransferRemoteChainsToMCMS: false,
			TestRouter:                 false,
			MCMS:                       nil,
		},
		{
			Msg:                        "Use test router (with MCMS)",
			TransferRemoteChainsToMCMS: true,
			TestRouter:                 true,
			MCMS:                       mcmsConfig,
		},
		{
			Msg:                        "Use test router (without MCMS)",
			TransferRemoteChainsToMCMS: false,
			TestRouter:                 true,
			MCMS:                       nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
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

			if test.TransferRemoteChainsToMCMS {
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
			}

			remoteChains := make(map[uint64]v1_6.ConnectionConfig, len(remoteChainSelectors))
			for _, selector := range remoteChainSelectors {
				remoteChains[selector] = v1_6.ConnectionConfig{
					RMNVerificationDisabled: false,
					AllowListEnabled:        false,
				}
			}

			e, err = commonchangeset.Apply(t, e, timelockContracts,
				commonchangeset.Configure(
					v1_6.ConnectNewChainChangeset,
					v1_6.ConnectNewChainConfig{
						TestRouter:       &test.TestRouter,
						RemoteChains:     remoteChains,
						NewChainSelector: newSelector,
						NewChainConnectionConfig: v1_6.ConnectionConfig{
							RMNVerificationDisabled: true,
							AllowListEnabled:        true,
						},
						MCMSConfig: test.MCMS,
					},
				),
			)
			if test.ErrStr != "" {
				require.ErrorContains(t, err, test.ErrStr, "expected ConnectNewChainChangeset error")
				return
			}
			require.NoError(t, err, "must apply ConnectNewChainChangeset")

			for _, selector := range selectors {
				expectedAllowListEnabled := true
				expectedRMNVerificationDisabled := true
				remoteSelectors := []uint64{newSelector}
				if selector == newSelector {
					expectedAllowListEnabled = false
					expectedRMNVerificationDisabled = false
					remoteSelectors = remoteChainSelectors
					if !test.TestRouter && test.MCMS != nil {
						// New chain must have all contracts owned by timelock
						mustHaveOwner(t, state.Chains[selector].OnRamp, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].OffRamp, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].FeeQuoter, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].RMNProxy, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].NonceManager, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].TokenAdminRegistry, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].Router, state.Chains[selector].Timelock.Address().Hex())
						mustHaveOwner(t, state.Chains[selector].RMNRemote, state.Chains[selector].Timelock.Address().Hex())

						// Admin role for deployer key should be revoked
						adminRole, err := state.Chains[selector].Timelock.ADMINROLE(nil)
						require.NoError(t, err, "must get admin role")
						hasRole, err := state.Chains[selector].Timelock.HasRole(nil, adminRole, e.Chains[selector].DeployerKey.From)
						require.NoError(t, err, "must get admin role")
						require.False(t, hasRole, "deployer key must not have admin role")
					} else {
						// onRamp, offRamp, and router should still be owned by deployer key
						mustHaveOwner(t, state.Chains[selector].OnRamp, e.Chains[selector].DeployerKey.From.Hex())
						mustHaveOwner(t, state.Chains[selector].OffRamp, e.Chains[selector].DeployerKey.From.Hex())
						mustHaveOwner(t, state.Chains[selector].Router, e.Chains[selector].DeployerKey.From.Hex())
					}
				}

				for _, remoteChainSelector := range remoteSelectors {
					expectedRouter := state.Chains[selector].Router
					if test.TestRouter {
						expectedRouter = state.Chains[selector].TestRouter
					}

					// onRamp checks
					destChainConfig, err := state.Chains[selector].OnRamp.GetDestChainConfig(nil, remoteChainSelector)
					require.NoError(t, err, "must get dest chain config from onRamp")
					require.Equal(t, expectedRouter.Address().Hex(), destChainConfig.Router.Hex(), "router must equal expected")
					require.Equal(t, expectedAllowListEnabled, destChainConfig.AllowlistEnabled, "allowListEnabled must equal expected")

					// offRamp checks
					srcChainConfig, err := state.Chains[selector].OffRamp.GetSourceChainConfig(nil, remoteChainSelector)
					require.NoError(t, err, "must get src chain config from offRamp")
					require.True(t, srcChainConfig.IsEnabled, "src chain config must be enabled")
					require.Equal(t, expectedRMNVerificationDisabled, srcChainConfig.IsRMNVerificationDisabled, "rmnVerificationDisabled must equal expected")
					require.Equal(t, common.LeftPadBytes(state.Chains[remoteChainSelector].OnRamp.Address().Bytes(), 32), srcChainConfig.OnRamp, "remote onRamp must be set on offRamp")
					require.Equal(t, expectedRouter.Address().Hex(), srcChainConfig.Router.Hex(), "router must equal expected")

					// router checks
					isOffRamp, err := expectedRouter.IsOffRamp(nil, remoteChainSelector, state.Chains[selector].OffRamp.Address())
					require.NoError(t, err, "must check if router has offRamp")
					require.True(t, isOffRamp, "router must have offRamp")
					onRamp, err := expectedRouter.GetOnRamp(nil, remoteChainSelector)
					require.NoError(t, err, "must get onRamp from router")
					require.Equal(t, state.Chains[selector].OnRamp.Address().Hex(), onRamp.Hex(), "onRamp must equal expected")
				}
			}
		})
	}
}
