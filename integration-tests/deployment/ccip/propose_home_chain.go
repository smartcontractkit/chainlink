package ccipdeployment

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

// SetCandidateExecPluginProposal calls setCandidate on the CCIPHome for setting up OCR3 exec Plugin config for the new chain.
func SetCandidateExecPluginProposal(
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
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[newChainSel].LinkToken, state.Chains[newChainSel].Weth9),
		nodes.NonBootstraps(),
		rmnHomeAddress,
	)
	if err != nil {
		return nil, err
	}

	execConfig, ok := newDONArgs[types.PluginTypeCCIPExec]
	if !ok {
		return nil, fmt.Errorf("missing exec plugin in ocr3Configs")
	}

	setCandidateMCMSOps, err := SetCandidateExecPluginOps(
		execConfig,
		state.Chains[homeChainSel].CapabilityRegistry,
		state.Chains[homeChainSel].CCIPHome,
		newChainSel,
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
		"SetCandidate for execution",
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(homeChainSel),
			Batch:           setCandidateMCMSOps,
		}},
		timelock.Schedule,
		"0s", // TODO: Should be parameterized.
	)
}

// PromoteCandidateProposal generates a proposal to call promoteCandidate on the CCIPHome through CapReg.
// This needs to be called after SetCandidateProposal is executed.
func PromoteCandidateProposal(
	state CCIPOnChainState,
	homeChainSel, newChainSel uint64,
	nodes deployment.Nodes,
) (*timelock.MCMSWithTimelockProposal, error) {
	promoteCandidateOps, err := PromoteCandidateOps(
		state.Chains[homeChainSel].CapabilityRegistry,
		state.Chains[homeChainSel].CCIPHome,
		newChainSel,
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
		"promoteCandidate for commit and execution",
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(homeChainSel),
			Batch:           promoteCandidateOps,
		}},
		timelock.Schedule,
		"0s", // TODO: Should be parameterized.
	)
}
