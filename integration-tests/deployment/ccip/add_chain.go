package ccipdeployment

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
)

func NewChainOutboundProposal(
	state CCIPOnChainState,
	newChainSel uint64,
	dests []uint64,
) (*timelock.MCMSWithTimelockProposal, error) {
	// Generate proposal which enables new source (from test router) on all existing destination chains.
	var batches []timelock.BatchChainOperation
	var onRampDestChainConfigArgs []onramp.OnRampDestChainConfigArgs
	var fqDestChainConfigArgs []fee_quoter.FeeQuoterDestChainConfigArgs
	for _, dest := range dests {
		onRampDestChainConfigArgs = append(onRampDestChainConfigArgs, onramp.OnRampDestChainConfigArgs{
			DestChainSelector: dest,
			Router:            state.Chains[newChainSel].TestRouter.Address(),
		})
		fqDestChainConfigArgs = append(fqDestChainConfigArgs, fee_quoter.FeeQuoterDestChainConfigArgs{
			DestChainSelector: dest,
			DestChainConfig:   defaultFeeQuoterDestChainConfig(),
		})
	}

	enableOnRampDest, err := state.Chains[newChainSel].OnRamp.ApplyDestChainConfigUpdates(
		deployment.SimTransactOpts(),
		onRampDestChainConfigArgs,
	)
	if err != nil {
		return nil, err
	}
	enableFeeQuoterDest, err := state.Chains[newChainSel].FeeQuoter.ApplyDestChainConfigUpdates(
		deployment.SimTransactOpts(),
		fqDestChainConfigArgs,
	)
	if err != nil {
		return nil, err
	}

	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(newChainSel),
		Batch: []mcms.Operation{
			{
				// Enable the source in on ramp
				To:    state.Chains[newChainSel].OnRamp.Address(),
				Data:  enableOnRampDest.Data(),
				Value: big.NewInt(0),
			},
			{
				To:    state.Chains[newChainSel].FeeQuoter.Address(),
				Data:  enableFeeQuoterDest.Data(),
				Value: big.NewInt(0),
			},
		},
	})

	timelockAddresses, metaDataPerChain, err := BuildProposalMetadata(state, []uint64{newChainSel})
	if err != nil {
		return nil, err
	}
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
	var chains []uint64
	for _, source := range sources {
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
			ChainIdentifier: mcms.ChainIdentifier(source),
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
		chains = append(chains, source)
	}

	addChainOp, err := ApplyChainConfigUpdatesOp(e, state, homeChainSel, []uint64{newChainSel})

	timelockAddresses, metaDataPerChain, err := BuildProposalMetadata(state, append(chains, homeChainSel))
	if err != nil {
		return nil, err
	}
	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(homeChainSel),
		Batch: []mcms.Operation{
			addChainOp,
		},
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

// AddDonAndSetCandidateForCommitProposal adds new DON for destination to home chain
// and sets the commit plugin config as candidateConfig for the don.
func AddDonAndSetCandidateForCommitProposal(
	state CCIPOnChainState,
	e deployment.Environment,
	nodes deployment.Nodes,
	ocrSecrets deployment.OCRSecrets,
	homeChainSel, feedChainSel, newChainSel uint64,
	tokenConfig TokenConfig,
	rmnHomeAddress common.Address,
) (*timelock.MCMSWithTimelockProposal, error) {
	newDONArgs, err := BuildOCR3ConfigForCCIPHome(
		e.Logger,
		ocrSecrets,
		state.Chains[newChainSel].OffRamp,
		e.Chains[newChainSel],
		feedChainSel,
		tokenConfig.GetLinkInfo(e.Logger, state.Chains[newChainSel].LinkToken, state.Chains[newChainSel].Weth9),
		nodes.NonBootstraps(),
		rmnHomeAddress,
	)
	if err != nil {
		return nil, err
	}
	latestDon, err := LatestCCIPDON(state.Chains[homeChainSel].CapabilityRegistry)
	if err != nil {
		return nil, err
	}
	commitConfig, ok := newDONArgs[types.PluginTypeCCIPCommit]
	if !ok {
		return nil, fmt.Errorf("missing commit plugin in ocr3Configs")
	}
	donID := latestDon.Id + 1
	addDonOp, err := SetCandidateCommitPluginWithAddDonOps(
		donID, commitConfig,
		state.Chains[homeChainSel].CapabilityRegistry,
		nodes.NonBootstraps(),
	)
	if err != nil {
		return nil, err
	}
	timelockAddresses, metaDataPerChain, err := BuildProposalMetadata(state, []uint64{homeChainSel})
	if err != nil {
		return nil, err
	}
	return timelock.NewMCMSWithTimelockProposal(
		"1",
		2004259681, // TODO: should be parameterized and based on current block timestamp.
		[]mcms.Signature{},
		false,
		metaDataPerChain,
		timelockAddresses,
		"SetCandidate for commit And AddDon for new chain",
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(homeChainSel),
			Batch:           []mcms.Operation{addDonOp},
		}},
		timelock.Schedule,
		"0s", // TODO: Should be parameterized.
	)
}
