package solana

import (
	"fmt"

	// "strconv"

	"github.com/gagliardetto/solana-go"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"

	"github.com/smartcontractkit/chainlink/deployment"
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
)

// SET OCR3 CONFIG
func btoi(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// SetOCR3OffRamp will set the OCR3 offramp for the given chain.
// to the active configuration on CCIPHome. This
// is used to complete the candidate->active promotion cycle, it's
// run after the candidate is confirmed to be working correctly.
// Multichain is especially helpful for NOP rotations where we have
// to touch all the chain to change signers.
func SetOCR3ConfigSolana(e deployment.Environment, cfg cs.SetOCR3OffRampConfig) (deployment.ChangesetOutput, error) {
	state, err := cs.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	if err := cfg.Validate(e, state); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	solChains := state.SolChains

	// cfg.RemoteChainSels will be a bunch of solana chains
	// can add this in validate
	for _, remote := range cfg.RemoteChainSels {
		donID, err := internal.DonIDForChain(
			state.Chains[cfg.HomeChainSel].CapabilityRegistry,
			state.Chains[cfg.HomeChainSel].CCIPHome,
			remote)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get don id for chain %d: %w", remote, err)
		}
		args, err := internal.BuildSetOCR3ConfigArgsSolana(donID, state.Chains[cfg.HomeChainSel].CCIPHome, remote)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build set ocr3 config args: %w", err)
		}
		// TODO: check if ocr3 has already been set
		// set, err := isOCR3ConfigSetSolana(e.Logger, e.Chains[remote], state.Chains[remote].OffRamp, args)
		var instructions []solana.Instruction
		offRampConfigPDA := solChains[remote].OffRampConfigPDA
		offRampStatePDA := solChains[remote].OffRampStatePDA
		solOffRamp.SetProgramID(solChains[remote].OffRamp)
		for _, arg := range args {
			instruction, err := solOffRamp.NewSetOcrConfigInstruction(
				arg.OCRPluginType,
				solOffRamp.Ocr3ConfigInfo{
					ConfigDigest:                   arg.ConfigDigest,
					F:                              arg.F,
					IsSignatureVerificationEnabled: btoi(arg.IsSignatureVerificationEnabled),
				},
				arg.Signers,
				arg.Transmitters,
				offRampConfigPDA,
				offRampStatePDA,
				e.SolChains[remote].DeployerKey.PublicKey(),
			).ValidateAndBuild()
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %w", err)
			}
			instructions = append(instructions, instruction)
		}
		if cfg.MCMS == nil {
			if err := e.SolChains[remote].Confirm(instructions); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}
	return deployment.ChangesetOutput{}, nil
}
