package ccipdeployment

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

type CreateDonArgs struct {
	CommitConfig *ccip_home.CCIPHomeOCR3Config
	ExecConfig   *ccip_home.CCIPHomeOCR3Config
	DonId        uint32
	CCIPHomeAbi  *abi.ABI
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
		}),
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

func NewSetCandidateProposal(
	state CCIPOnChainState,
	e deployment.Environment,
	homeChainSel, feedChainSel, newChainSel uint64,
	tokenConfig TokenConfig,
	rmnHomeAddress []byte,
) (*timelock.MCMSWithTimelockProposal, *CreateDonArgs, error) {
	// Home chain new don.
	// - Add new DONs for destination to home chain
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return nil, nil, err
	}

	newDONArgs, err := BuildAddDONArgs(
		e.Logger,
		state.Chains[newChainSel].OffRamp,
		e.Chains[newChainSel],
		feedChainSel,
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[newChainSel].LinkToken),
		nodes.NonBootstraps(),
		rmnHomeAddress,
	)
	if err != nil {
		return nil, nil, err
	}
	tabi, err := ccip_home.CCIPHomeMetaData.GetAbi()
	if err != nil {
		return nil, nil, err
	}
	commitConfig, execConfig, donId, err := FormCreateDonArgs(state.Chains[homeChainSel].CapabilityRegistry, newDONArgs)
	createDonArgs := &CreateDonArgs{
		CommitConfig: &commitConfig,
		ExecConfig:   &execConfig,
		DonId:        donId,
		CCIPHomeAbi:  tabi,
	}
	setCandidateMCMSOps, err := SetCandidateOps(
		tabi,
		state.Chains[homeChainSel].CapabilityRegistry,
		commitConfig,
		execConfig,
		donId,
		nodes.NonBootstraps(),
	)

	if err != nil {
		return nil, nil, err
	}
	tl := state.Chains[homeChainSel].Timelock
	mcm := state.Chains[homeChainSel].ProposerMcm
	proposal, err := CreateSingleChainMCMSOps(setCandidateMCMSOps, homeChainSel, "Set Candidate", tl, mcm)
	if err != nil {
		return nil, nil, err
	}
	return proposal, createDonArgs, nil
}

func NewPromoteCandidateProposal(
	state CCIPOnChainState, homeChainSel uint64,
	createDonArgs *CreateDonArgs,
	e deployment.Environment,
) (*timelock.MCMSWithTimelockProposal, error) {
	if createDonArgs == nil {
		return nil, fmt.Errorf("createDonArgs is nil")
	}
	if createDonArgs.CommitConfig == nil {
		return nil, fmt.Errorf("commitConfig is nil")
	}
	if createDonArgs.ExecConfig == nil {
		return nil, fmt.Errorf("execConfig is nil")
	}
	if createDonArgs.CCIPHomeAbi == nil {
		return nil, fmt.Errorf("ccipHomeAbi is nil")
	}
	if createDonArgs.DonId == 0 {
		return nil, fmt.Errorf("donId is 0")
	}
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return nil, err
	}
	promoteCandidateOps, err := PromoteCandidateOps(
		createDonArgs.CCIPHomeAbi,
		state.Chains[homeChainSel].CapabilityRegistry,
		state.Chains[homeChainSel].CCIPHome,
		*createDonArgs.CommitConfig,
		*createDonArgs.ExecConfig,
		createDonArgs.DonId,
		nodes.NonBootstraps(),
	)
	if err != nil {
		return nil, err
	}
	tl := state.Chains[homeChainSel].Timelock
	mcm := state.Chains[homeChainSel].ProposerMcm
	return CreateSingleChainMCMSOps(promoteCandidateOps, homeChainSel, "Promote Candidate", tl, mcm)
}
