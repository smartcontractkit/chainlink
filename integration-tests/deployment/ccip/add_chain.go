package ccipdeployment

import (
	"math/big"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/timelock"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

// NewChainInboundProposal generates a proposal
// to connect the new chain to the existing chains.
func NewChainInboundProposal(
	e deployment.Environment,
	state CCIPOnChainState,
	homeChainSel uint64,
	newChainSel uint64,
	sources []uint64,
) (*timelock.MCMSWithTimelockProposal, error) {
	// Generate proposal which enables new destination (from test router) on all source chains.
	var batches []timelock.BatchChainOperation
	metaDataPerChain := make(map[mcms.ChainIdentifier]timelock.MCMSWithTimelockChainMetadata)
	for _, source := range sources {
		chain, _ := chainsel.ChainBySelector(source)
		enableOnRampDest, err := state.Chains[source].OnRamp.ApplyDestChainConfigUpdates(SimTransactOpts(), []onramp.OnRampDestChainConfigArgs{
			{
				DestChainSelector: newChainSel,
				Router:            state.Chains[source].TestRouter.Address(),
			},
		})
		if err != nil {
			return nil, err
		}
		enableFeeQuoterDest, err := state.Chains[source].FeeQuoter.ApplyDestChainConfigUpdates(
			SimTransactOpts(),
			[]fee_quoter.FeeQuoterDestChainConfigArgs{
				{
					DestChainSelector: newChainSel,
					DestChainConfig:   defaultFeeQuoterDestChainConfig(),
				},
			})
		if err != nil {
			return nil, err
		}
		initialPrices, err := state.Chains[source].FeeQuoter.UpdatePrices(
			SimTransactOpts(),
			fee_quoter.InternalPriceUpdates{
				TokenPriceUpdates: []fee_quoter.InternalTokenPriceUpdate{},
				GasPriceUpdates: []fee_quoter.InternalGasPriceUpdate{
					{
						DestChainSelector: newChainSel,
						// TODO: parameterize
						UsdPerUnitGas: big.NewInt(2e12),
					},
				}})
		if err != nil {
			return nil, err
		}
		batches = append(batches, timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(chain.Selector),
			Batch: []mcms.Operation{
				{
					// Enable the source in on ramp
					To:    state.Chains[source].OnRamp.Address(),
					Data:  enableOnRampDest.Data(),
					Value: big.NewInt(0),
				},
				{
					// Set initial dest prices to unblock testing.
					To:    state.Chains[source].FeeQuoter.Address(),
					Data:  initialPrices.Data(),
					Value: big.NewInt(0),
				},
				{
					To:    state.Chains[source].FeeQuoter.Address(),
					Data:  enableFeeQuoterDest.Data(),
					Value: big.NewInt(0),
				},
			},
		})
		metaDataPerChain[mcms.ChainIdentifier(chain.Selector)] = timelock.MCMSWithTimelockChainMetadata{
			ChainMetadata: mcms.ChainMetadata{
				NonceOffset: 0,
				MCMAddress:  state.Chains[source].Mcm.Address(),
			},
			TimelockAddress: state.Chains[source].Timelock.Address(),
		}
	}

	// Home chain new don.
	// - Add new DONs for destination to home chain
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return nil, err
	}
	encodedExtraChainConfig, err := chainconfig.EncodeChainConfig(chainconfig.ChainConfig{
		GasPriceDeviationPPB:    ccipocr3.NewBigIntFromInt64(1000),
		DAGasPriceDeviationPPB:  ccipocr3.NewBigIntFromInt64(0),
		OptimisticConfirmations: 1,
	})
	if err != nil {
		return nil, err
	}
	chainConfig := SetupConfigInfo(newChainSel, nodes.NonBootstraps().PeerIDs(),
		nodes.DefaultF(), encodedExtraChainConfig)
	addChain, err := state.Chains[homeChainSel].CCIPConfig.ApplyChainConfigUpdates(SimTransactOpts(), nil, []ccip_config.CCIPConfigTypesChainConfigInfo{
		chainConfig,
	})
	if err != nil {
		return nil, err
	}

	newDONArgs, err := BuildAddDONArgs(e.Logger, state.Chains[newChainSel].OffRamp, e.Chains[newChainSel], nodes.NonBootstraps())
	if err != nil {
		return nil, err
	}
	addDON, err := state.Chains[homeChainSel].CapabilityRegistry.AddDON(SimTransactOpts(),
		nodes.NonBootstraps().PeerIDs(), []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: CCIPCapabilityID,
				Config:       newDONArgs,
			},
		}, false, false, nodes.NonBootstraps().DefaultF())
	if err != nil {
		return nil, err
	}
	homeChain, _ := chainsel.ChainBySelector(homeChainSel)
	metaDataPerChain[mcms.ChainIdentifier(homeChain.Selector)] = timelock.MCMSWithTimelockChainMetadata{
		ChainMetadata: mcms.ChainMetadata{
			NonceOffset: 0,
			MCMAddress:  state.Chains[homeChainSel].Mcm.Address(),
		},
		TimelockAddress: state.Chains[homeChainSel].Timelock.Address(),
	}
	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(homeChain.Selector),
		Batch: []mcms.Operation{
			{
				// Add the chain first, don needs it to be there.
				To:    state.Chains[homeChainSel].CCIPConfig.Address(),
				Data:  addChain.Data(),
				Value: big.NewInt(0),
			},
			{
				To:    state.Chains[homeChainSel].CapabilityRegistry.Address(),
				Data:  addDON.Data(),
				Value: big.NewInt(0),
			},
		},
	})
	return timelock.NewMCMSWithTimelockProposal(
		"1",
		2004259681, // TODO: should be parameterized and based on current block timestamp.
		[]mcms.Signature{},
		false,
		metaDataPerChain,
		"blah", // TODO
		batches,
		timelock.Schedule,
		"0s", // TODO: Should be parameterized.
	)
}
