//go:build !integration

package v1_6_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type laneDefinition struct {
	Source v1_6.ChainDefinition
	Dest   v1_6.ChainDefinition
}

func getAllPossibleLanes(chains []v1_6.ChainDefinition, disable bool) []v1_6.BidirectionalLaneDefinition {
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

			lanes = append(lanes, v1_6.BidirectionalLaneDefinition{
				Chains:     [2]v1_6.ChainDefinition{chainA, chainB},
				IsDisabled: disable,
			})
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

func TestBuildConfigs(t *testing.T) {
	selectors := []uint64{1, 2}

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

	cfg := v1_6.UpdateBidirectionalLanesConfig{
		TestRouter: false,
		MCMSConfig: &proposalutils.TimelockConfig{
			MinDelay:   0 * time.Second,
			MCMSAction: types.TimelockActionSchedule,
		},
		Lanes: getAllPossibleLanes(chains, false),
	}

	configs := cfg.BuildConfigs()

	require.Equal(t, v1_6.UpdateFeeQuoterDestsConfig{
		UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
			1: {
				2: v1_6.DefaultFeeQuoterDestChainConfig(true),
			},
			2: {
				1: v1_6.DefaultFeeQuoterDestChainConfig(true),
			},
		},
		MCMS: cfg.MCMSConfig,
	}, configs.UpdateFeeQuoterDestsConfig)
	require.Equal(t, v1_6.UpdateFeeQuoterPricesConfig{
		PricesByChain: map[uint64]v1_6.FeeQuoterPriceUpdatePerSource{
			1: {
				GasPrices: map[uint64]*big.Int{
					2: big.NewInt(1e17),
				},
			},
			2: {
				GasPrices: map[uint64]*big.Int{
					1: big.NewInt(1e17),
				},
			},
		},
		MCMS: cfg.MCMSConfig,
	}, configs.UpdateFeeQuoterPricesConfig)
	require.Equal(t, v1_6.UpdateOffRampSourcesConfig{
		UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
			1: {
				2: {
					IsEnabled:                 true,
					TestRouter:                false,
					IsRMNVerificationDisabled: true,
				},
			},
			2: {
				1: {
					IsEnabled:                 true,
					TestRouter:                false,
					IsRMNVerificationDisabled: true,
				},
			},
		},
		MCMS: cfg.MCMSConfig,
	}, configs.UpdateOffRampSourcesConfig)
	require.Equal(t, v1_6.UpdateOnRampDestsConfig{
		UpdatesByChain: map[uint64]map[uint64]v1_6.OnRampDestinationUpdate{
			1: {
				2: {
					IsEnabled:        true,
					TestRouter:       false,
					AllowListEnabled: false,
				},
			},
			2: {
				1: {
					IsEnabled:        true,
					TestRouter:       false,
					AllowListEnabled: false,
				},
			},
		},
		MCMS: cfg.MCMSConfig,
	}, configs.UpdateOnRampDestsConfig)
	require.Equal(t, v1_6.UpdateRouterRampsConfig{
		UpdatesByChain: map[uint64]v1_6.RouterUpdates{
			1: {
				OnRampUpdates: map[uint64]bool{
					2: true,
				},
				OffRampUpdates: map[uint64]bool{
					2: true,
				},
			},
			2: {
				OnRampUpdates: map[uint64]bool{
					1: true,
				},
				OffRampUpdates: map[uint64]bool{
					1: true,
				},
			},
		},
		MCMS: cfg.MCMSConfig,
	}, configs.UpdateRouterRampsConfig)
}
