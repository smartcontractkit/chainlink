package ccipdeployment

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/timelock"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
)

// AddChain deploys chain contracts for a new chain
// and generates 3 proposals to connect that new chain to all existing chains.
// We testing in between each proposal.
func NewChainInbound(
	e deployment.Environment,
	ab deployment.AddressBook,
	homeChainSel uint64,
	newChainSel uint64,
	sources []uint64,
) ([]timelock.MCMSWithTimelockProposal, deployment.AddressBook, error) {
	// 1. Deploy contracts to new chain and wire them.
	newAddresses, err := DeployChainContracts(e, e.Chains[newChainSel], deployment.NewMemoryAddressBook())
	if err != nil {
		return nil, ab, err
	}
	if err := ab.Merge(newAddresses); err != nil {
		return nil, ab, err
	}
	state, err := LoadOnchainState(e, ab)
	if err != nil {
		return nil, ab, err
	}

	// 2. Generate proposal which enables new destination (from test router) on all source chains.
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
			return nil, ab, err
		}
		enablePriceRegDest, err := state.Chains[source].FeeQuoter.ApplyDestChainConfigUpdates(
			SimTransactOpts(),
			[]fee_quoter.FeeQuoterDestChainConfigArgs{
				{
					DestChainSelector: newChainSel,
					DestChainConfig:   defaultFeeQuoterDestChainConfig(),
				},
			})
		if err != nil {
			return nil, ab, err
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
			return nil, ab, err
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
					// Set initial dest prices to unblock testing.
					To:    state.Chains[source].FeeQuoter.Address(),
					Data:  enablePriceRegDest.Data(),
					Value: big.NewInt(0),
				},
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
		return nil, ab, err
	}
	newDONArgs, err := BuildAddDONArgs(e.Logger, state.Chains[newChainSel].OffRamp, e.Chains[newChainSel], nodes)
	if err != nil {
		return nil, ab, err
	}
	addDON, err := state.Chains[homeChainSel].CapabilityRegistry.AddDON(SimTransactOpts(),
		nodes.PeerIDs(newChainSel), []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: CCIPCapabilityId,
				Config:       newDONArgs,
			},
		}, false, false, nodes.DefaultF())
	if err != nil {
		return nil, ab, err
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
	newDestProposal, err := timelock.NewMCMSWithTimelockProposal(
		"1",
		uint32(time.Now().Add(1*time.Hour).Unix()),
		[]mcms.Signature{},
		false,
		metaDataPerChain,
		"blah",
		batches,
		timelock.Schedule, "1h")
	if err != nil {
		return nil, ab, err
	}

	// New chain we can configure directly with deployer key first.
	var offRampEnables []offramp.OffRampSourceChainConfigArgs
	for _, source := range sources {
		offRampEnables = append(offRampEnables, offramp.OffRampSourceChainConfigArgs{
			Router:              state.Chains[newChainSel].Router.Address(),
			SourceChainSelector: source,
			IsEnabled:           true,
			OnRamp:              common.LeftPadBytes(state.Chains[source].OnRamp.Address().Bytes(), 32),
		})
	}
	tx, err := state.Chains[newChainSel].OffRamp.ApplySourceChainConfigUpdates(e.Chains[newChainSel].DeployerKey, offRampEnables)
	if _, err := deployment.ConfirmIfNoError(e.Chains[newChainSel], tx, err); err != nil {
		return nil, ab, err
	}

	// We won't actually be able to setOCR3Config on the remote until the first proposal goes through.
	// TODO: Outbound
	return []timelock.MCMSWithTimelockProposal{*newDestProposal}, ab, nil
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
