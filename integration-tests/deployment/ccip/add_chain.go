package ccipdeployment

import (
	"math/big"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/timelock"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
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
		/*
			enablePriceRegDest, err := state.Chains[source].FeeQuoter.ApplyDestChainConfigUpdates(
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
		*/
		batches = append(batches, timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(chain.Selector),
			Batch: []mcms.Operation{
				{
					// Enable the source in on ramp
					To:    state.Chains[source].OnRamp.Address(),
					Data:  enableOnRampDest.Data(),
					Value: big.NewInt(0),
				},
				//{
				//	// Set initial dest prices to unblock testing.
				//	To:    state.Chains[source].FeeQuoter.Address(),
				//	Data:  initialPrices.Data(),
				//	Value: big.NewInt(0),
				//},
				//{
				//	// Set initial dest prices to unblock testing.
				//	To:    state.Chains[source].FeeQuoter.Address(),
				//	Data:  enablePriceRegDest.Data(),
				//	Value: big.NewInt(0),
				//},
			},
		})
		metaDataPerChain[mcms.ChainIdentifier(chain.Selector)] = timelock.MCMSWithTimelockChainMetadata{
			ChainMetadata: mcms.ChainMetadata{
				NonceOffset: 0,
				MCMAddress:  state.Chains[source].McmAddr,
			},
			TimelockAddress: state.Chains[source].TimelockAddr,
		}
	}

	// Home chain new don.
	// - Add new DONs for destination to home chain
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return nil, err
	}
	newDONArgs, err := BuildAddDONArgs(e.Logger, state.Chains[newChainSel].OffRamp, e.Chains[newChainSel], nodes)
	if err != nil {
		return nil, err
	}
	addDON, err := state.Chains[homeChainSel].CapabilityRegistry.AddDON(SimTransactOpts(),
		nodes.PeerIDs(newChainSel), []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: CCIPCapabilityId,
				Config:       newDONArgs,
			},
		}, false, false, nodes.DefaultF())
	if err != nil {
		return nil, err
	}
	homeChain, _ := chainsel.ChainBySelector(homeChainSel)
	metaDataPerChain[mcms.ChainIdentifier(homeChain.Selector)] = timelock.MCMSWithTimelockChainMetadata{
		ChainMetadata: mcms.ChainMetadata{
			NonceOffset: 0,
			MCMAddress:  state.Chains[homeChainSel].McmAddr,
		},
		TimelockAddress: state.Chains[homeChainSel].TimelockAddr,
	}
	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(homeChain.Selector),
		Batch: []mcms.Operation{
			{
				// Enable the source in on ramp
				To:    state.Chains[homeChainSel].CapabilityRegistry.Address(),
				Data:  addDON.Data(),
				Value: big.NewInt(0),
			},
		},
	})
	return timelock.NewMCMSWithTimelockProposal(
		"1",
		2004259681,
		[]mcms.Signature{},
		false,
		metaDataPerChain,
		"blah",
		batches,
		timelock.Schedule, "0s")
	// We won't actually be able to setOCR3Config on the remote until the first proposal goes through.
	// TODO: Outbound
}

//func ApplyInboundChainProposal(
//	e deployment.Environment,
//	ab deployment.AddressBook,
//	proposal managed.MCMSWithTimelockProposal,
//) (deployment.AddressBook, error) {
//	state, err := LoadOnchainState(e, ab)
//	if err != nil {
//		return ab, err
//	}
//
//	// Apply the proposal.
//})

// 1. Deploy contracts
// 2. Proposal 1 (allow for inbound testing)
// - Enables new destination in onramps using test router
// - Enables the sources in the offramp and real router.
// - Sets initial prices for destination in price reg.
// - Add new DONs for destination to home chain
// - SetOCR3Config(s) on destination offramp.
// 3. At this point should be able to test from all sources
// and ensure that its writing those source prices to the new chain.
// 4. Proposal 2 (allow for outbound testing)
// -  Add new destinations on onramp/price reg can use real router.
// No initial prices needed because DON updating them.
// - Add new sources to the remote offramps (test router).
// - Add ChainConfig to home chain so existing OCR instances become aware of the source.
// 5. Now we can test the other direction.
// 6 . Proposal 3 move onramp/offramps on existing chains to real router.
