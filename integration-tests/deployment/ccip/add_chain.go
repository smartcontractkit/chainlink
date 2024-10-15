package ccipdeployment

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
)

// NewChainInboundProposal generates a proposal
// to connect the new chain to the existing chains.
func NewChainInboundProposal(
	e deployment.Environment,
	ocrSecrets deployment.OCRSecrets,
	state CCIPOnChainState,
	homeChainSel uint64,
	feedChainSel uint64,
	newChainSel uint64,
	sources []uint64,
	tokenConfig TokenConfig,
	rmnHomeAddress []byte,
) (*timelock.MCMSWithTimelockProposal, error) {
	// Generate proposal which enables new destination (from test router) on all source chains.
	var batches []timelock.BatchChainOperation
	metaDataPerChain := make(map[mcms.ChainIdentifier]mcms.ChainMetadata)
	timelockAddresses := make(map[mcms.ChainIdentifier]common.Address)
	for _, source := range sources {
		chain, _ := chainsel.ChainBySelector(source)
		enableOnRampDest, err := state.Chains[source].OnRamp.ApplyDestChainConfigUpdates(deployment.SimTransactOpts(), []onramp.OnRampDestChainConfigArgs{
			{
				DestChainSelector: newChainSel,
				Router:            state.Chains[source].TestRouter.Address(),
			},
		})
		if err != nil {
			return nil, err
		}
		enableFeeQuoterDest, err := state.Chains[source].FeeQuoter.ApplyDestChainConfigUpdates(
			deployment.SimTransactOpts(),
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
			deployment.SimTransactOpts(),
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
		opCount, err := state.Chains[source].ProposerMcm.GetOpCount(nil)
		if err != nil {
			return nil, err
		}
		metaDataPerChain[mcms.ChainIdentifier(chain.Selector)] = mcms.ChainMetadata{
			StartingOpCount: opCount.Uint64(),
			MCMAddress:      state.Chains[source].ProposerMcm.Address(),
		}
		timelockAddresses[mcms.ChainIdentifier(chain.Selector)] = state.Chains[source].Timelock.Address()
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
	addChain, err := state.Chains[homeChainSel].CCIPHome.ApplyChainConfigUpdates(
		deployment.SimTransactOpts(), nil, []ccip_home.CCIPHomeChainConfigArgs{
			chainConfig,
		})
	if err != nil {
		return nil, err
	}

	newDONArgs, err := BuildAddDONArgs(
		e.Logger,
		ocrSecrets,
		state.Chains[newChainSel].OffRamp,
		e.Chains[newChainSel],
		feedChainSel,
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[newChainSel].LinkToken),
		nodes.NonBootstraps(),
		rmnHomeAddress,
	)
	if err != nil {
		return nil, err
	}
	mcmsOps, err := CreateDON(
		e.Logger,
		state.Chains[homeChainSel].CapabilityRegistry,
		state.Chains[homeChainSel].CCIPHome,
		newDONArgs,
		e.Chains[homeChainSel],
		nodes,
	)
	//addDON, err := state.Chains[homeChainSel].CapabilityRegistry.AddDON(SimTransactOpts(),
	//	nodes.NonBootstraps().PeerIDs(), []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
	//		{
	//			CapabilityId: CCIPCapabilityID,
	//			Config:       newDONArgs,
	//		},
	//	}, false, false, nodes.NonBootstraps().DefaultF())
	if err != nil {
		return nil, err
	}
	homeChain, _ := chainsel.ChainBySelector(homeChainSel)
	opCount, err := state.Chains[homeChainSel].ProposerMcm.GetOpCount(nil)
	if err != nil {
		return nil, err
	}
	metaDataPerChain[mcms.ChainIdentifier(homeChain.Selector)] = mcms.ChainMetadata{
		StartingOpCount: opCount.Uint64(),
		MCMAddress:      state.Chains[homeChainSel].ProposerMcm.Address(),
	}
	timelockAddresses[mcms.ChainIdentifier(homeChain.Selector)] = state.Chains[homeChainSel].Timelock.Address()
	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(homeChain.Selector),
		Batch: append([]mcms.Operation{
			{
				// Add the chain first, don needs it to be there.
				To:    state.Chains[homeChainSel].CCIPHome.Address(),
				Data:  addChain.Data(),
				Value: big.NewInt(0),
			},
		}, mcmsOps...),
	})
	return timelock.NewMCMSWithTimelockProposal(
		"1",
		2004259681, // TODO: should be parameterized and based on current block timestamp.
		[]mcms.Signature{},
		false,
		metaDataPerChain,
		timelockAddresses,
		"blah", // TODO
		batches,
		timelock.Schedule,
		"0s", // TODO: Should be parameterized.
	)
}
