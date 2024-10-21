package ccipdeployment

import (
	"fmt"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"math/big"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
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
	var chains []uint64
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
		chains = append(chains, source)
	}

	homeChain, _ := chainsel.ChainBySelector(homeChainSel)
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

	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(homeChain.Selector),
		Batch: []mcms.Operation{
			{
				// Add the chain first, don needs it to be there.
				To:    state.Chains[homeChainSel].CCIPHome.Address(),
				Data:  addChain.Data(),
				Value: big.NewInt(0),
			},
		},
	})

	return BuildProposalFromBatches(state, batches, "proposal to set new chains", "0s")
}

// AddDonAndSetCandidateProposal adds new DON for destination to home chain
// and sets the commit plugin config as candidateConfig for the don.
func AddDonAndSetCandidateProposal(
	state CCIPOnChainState,
	e deployment.Environment,
	nodes deployment.Nodes,
	ocrSecrets deployment.OCRSecrets,
	homeChainSel, feedChainSel, newChainSel uint64,
	tokenConfig TokenConfig,
	pluginType types.PluginType,
) (*timelock.MCMSWithTimelockProposal, error) {
	newDONArgs, err := BuildOCR3ConfigForCCIPHome(
		e.Logger,
		ocrSecrets,
		state.Chains[newChainSel].OffRamp,
		e.Chains[newChainSel],
		feedChainSel,
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[newChainSel].LinkToken),
		nodes.NonBootstraps(),
		state.Chains[homeChainSel].RMNHome.Address(),
	)
	if err != nil {
		return nil, err
	}
	latestDon, err := LatestCCIPDON(state.Chains[homeChainSel].CapabilityRegistry)
	if err != nil {
		return nil, err
	}
	commitConfig, ok := newDONArgs[pluginType]
	if !ok {
		return nil, fmt.Errorf("missing commit plugin in ocr3Configs")
	}
	donID := latestDon.Id + 1
	addDonOp, err := NewDonWithCandidateOp(
		donID, commitConfig,
		state.Chains[homeChainSel].CapabilityRegistry,
		nodes.NonBootstraps(),
	)
	if err != nil {
		return nil, err
	}

	return BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(homeChainSel),
		Batch:           []mcms.Operation{addDonOp},
	}}, "setCandidate for commit and AddDon on new Chain", "0s")
}
