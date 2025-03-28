package v1_6_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type laneDefinition struct {
	Source v1_6.ChainDefinition
	Dest   v1_6.ChainDefinition
}

func getAllPossibleLanes(chains []v1_6.ChainDefinition) []v1_6.BidirectionalLaneDefinition {
	lanes := make([]v1_6.BidirectionalLaneDefinition, 0)
	paired := make(map[uint64]map[uint64]bool)

	for i, chainA := range chains {
		for j, chainB := range chains {
			if i == j {
				continue
			}
			if paired[chainA.Selector] != nil && paired[chainA.Selector][chainB.Selector] {
				continue
			}
			if paired[chainB.Selector] != nil && paired[chainB.Selector][chainA.Selector] {
				continue
			}

			lanes = append(lanes, v1_6.BidirectionalLaneDefinition{chainA, chainB})
			if paired[chainA.Selector] == nil {
				paired[chainA.Selector] = make(map[uint64]bool)
			}
			paired[chainA.Selector][chainB.Selector] = true
			if paired[chainB.Selector] == nil {
				paired[chainB.Selector] = make(map[uint64]bool)
			}
			paired[chainB.Selector][chainA.Selector] = true
		}
	}

	return lanes
}

func getRemoteChains(chains []v1_6.ChainDefinition, currentIndex int) []v1_6.ChainDefinition {
	remoteChains := make([]v1_6.ChainDefinition, 0, len(chains)-1)
	for i, chain := range chains {
		if i == currentIndex {
			continue
		}
		remoteChains = append(remoteChains, chain)
	}
	return remoteChains
}

func checkBidirectionalLaneConnectivity(
	t *testing.T,
	e deployment.Environment,
	state changeset.CCIPOnChainState,
	chainOne v1_6.ChainDefinition,
	chainTwo v1_6.ChainDefinition,
	testRouter bool,
) {
	lanes := []laneDefinition{
		{
			Source: chainOne,
			Dest:   chainTwo,
		},
		{
			Source: chainTwo,
			Dest:   chainOne,
		},
	}
	for _, lane := range lanes {
		onRamp := state.Chains[lane.Source.Selector].OnRamp
		offRamp := state.Chains[lane.Dest.Selector].OffRamp
		feeQuoterOnSrc := state.Chains[lane.Source.Selector].FeeQuoter
		routerOnSrc := state.Chains[lane.Source.Selector].Router
		routerOnDest := state.Chains[lane.Dest.Selector].Router
		if testRouter {
			routerOnSrc = state.Chains[lane.Source.Selector].TestRouter
			routerOnDest = state.Chains[lane.Dest.Selector].TestRouter
		}

		destChainConfig, err := onRamp.GetDestChainConfig(nil, lane.Dest.Selector)
		require.NoError(t, err, "must get dest chain config from onRamp")
		require.Equal(t, routerOnSrc.Address().Hex(), destChainConfig.Router.Hex(), "router must equal expected")
		require.Equal(t, lane.Dest.AllowListEnabled, destChainConfig.AllowlistEnabled, "allowListEnabled must equal expected")

		srcChainConfig, err := offRamp.GetSourceChainConfig(nil, lane.Source.Selector)
		require.NoError(t, err, "must get src chain config from offRamp")
		require.True(t, srcChainConfig.IsEnabled, "src chain config must be enabled")
		require.Equal(t, lane.Source.RMNVerificationDisabled, srcChainConfig.IsRMNVerificationDisabled, "rmnVerificationDisabled must equal expected")
		require.Equal(t, common.LeftPadBytes(state.Chains[lane.Source.Selector].OnRamp.Address().Bytes(), 32), srcChainConfig.OnRamp, "remote onRamp must be set on offRamp")
		require.Equal(t, routerOnDest.Address().Hex(), srcChainConfig.Router.Hex(), "router must equal expected")

		isOffRamp, err := routerOnSrc.IsOffRamp(nil, lane.Dest.Selector, state.Chains[lane.Source.Selector].OffRamp.Address())
		require.NoError(t, err, "must check if router has offRamp")
		require.True(t, isOffRamp, "router must have offRamp")
		onRampOnRouter, err := routerOnSrc.GetOnRamp(nil, lane.Dest.Selector)
		require.NoError(t, err, "must get onRamp from router")
		require.Equal(t, state.Chains[lane.Source.Selector].OnRamp.Address().Hex(), onRampOnRouter.Hex(), "onRamp must equal expected")

		isOffRamp, err = routerOnDest.IsOffRamp(nil, lane.Source.Selector, state.Chains[lane.Dest.Selector].OffRamp.Address())
		require.NoError(t, err, "must check if router has offRamp")
		require.True(t, isOffRamp, "router must have offRamp")
		onRampOnRouter, err = routerOnDest.GetOnRamp(nil, lane.Source.Selector)
		require.NoError(t, err, "must get onRamp from router")
		require.Equal(t, state.Chains[lane.Dest.Selector].OnRamp.Address().Hex(), onRampOnRouter.Hex(), "onRamp must equal expected")

		feeQuoterDestConfig, err := feeQuoterOnSrc.GetDestChainConfig(nil, lane.Dest.Selector)
		require.NoError(t, err, "must get dest chain config from feeQuoter")
		require.Equal(t, lane.Dest.FeeQuoterDestChainConfig, feeQuoterDestConfig, "feeQuoter dest chain config must equal expected")

		price, err := feeQuoterOnSrc.GetDestinationChainGasPrice(nil, lane.Dest.Selector)
		require.NoError(t, err, "must get price from feeQuoter")
		require.Equal(t, lane.Dest.GasPrice, price.Value, "price must equal expected")
	}
}

func TestAddBidirectionalLanesChangeset(t *testing.T) {
	t.Parallel()

	type test struct {
		Msg        string
		TestRouter bool
		MCMS       *proposalutils.TimelockConfig
	}

	mcmsConfig := &proposalutils.TimelockConfig{
		MinDelay:   0 * time.Second,
		MCMSAction: types.TimelockActionSchedule,
	}

	tests := []test{
		{
			Msg:        "Use production router (with MCMS)",
			TestRouter: false,
			MCMS:       mcmsConfig,
		},
		{
			Msg:        "Use test router (without MCMS)",
			TestRouter: true,
			MCMS:       nil,
		},
		{
			Msg:        "Use test router (with MCMS for other contracts)",
			TestRouter: true,
			MCMS:       mcmsConfig,
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

			timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts, len(selectors))
			for _, selector := range selectors {
				// Assemble map of addresses required for Timelock scheduling & execution
				timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
					Timelock:  state.Chains[selector].Timelock,
					CallProxy: state.Chains[selector].CallProxy,
				}
			}

			if test.MCMS != nil {
				contractsToTransfer := make(map[uint64][]common.Address, len(selectors))
				for _, selector := range selectors {
					contractsToTransfer[selector] = []common.Address{
						state.Chains[selector].OnRamp.Address(),
						state.Chains[selector].OffRamp.Address(),
						state.Chains[selector].Router.Address(),
						state.Chains[selector].FeeQuoter.Address(),
					}
				}
				e, err = commonchangeset.Apply(t, e, timelockContracts,
					commonchangeset.Configure(
						deployment.CreateLegacyChangeSet(commoncs.TransferToMCMSWithTimelock),
						commoncs.TransferToMCMSWithTimelockConfig{
							ContractsByChain: contractsToTransfer,
							MCMSConfig: proposalutils.TimelockConfig{
								MinDelay:   0 * time.Second,
								MCMSAction: types.TimelockActionSchedule,
							},
						},
					),
				)
				require.NoError(t, err, "must apply TransferToMCMSWithTimelock")
			}

			chains := make([]v1_6.ChainDefinition, len(selectors))
			for i, selector := range selectors {
				chains[i] = v1_6.ChainDefinition{
					ConnectionConfig: v1_6.ConnectionConfig{
						RMNVerificationDisabled: true,
						AllowListEnabled:        false,
					},
					Selector:                 selector,
					GasPrice:                 big.NewInt(1e17),
					FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
				}
			}

			e, err = commonchangeset.Apply(t, e, timelockContracts,
				commonchangeset.Configure(
					v1_6.AddBidirectionalLanesChangeset,
					v1_6.AddBidirectionalLanesConfig{
						TestRouter: test.TestRouter,
						MCMSConfig: test.MCMS,
						Lanes:      getAllPossibleLanes(chains),
					},
				),
			)
			require.NoError(t, err, "must apply AddBidirectionalLanesChangeset")

			for i, chain := range chains {
				remoteChains := getRemoteChains(chains, i)
				for _, remoteChain := range remoteChains {
					checkBidirectionalLaneConnectivity(t, e, state, chain, remoteChain, test.TestRouter)
				}
			}
		})
	}
}
