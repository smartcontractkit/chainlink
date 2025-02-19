package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"

	"github.com/smartcontractkit/chainlink/deployment"
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
)

const (
	OcrCommitPlugin uint8 = iota
	OcrExecutePlugin
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

	for _, remote := range cfg.RemoteChainSels {
		chainFamily, _ := chain_selectors.GetSelectorFamily(remote)
		if chainFamily != chain_selectors.FamilySolana {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain %d is not a solana chain", remote)
		}
	}

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
		set, err := isOCR3ConfigSetOnOffRampSolana(e, e.SolChains[remote], state.SolChains[remote], args)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to check if ocr3 config is set on offramp: %w", err)
		}
		if set {
			e.Logger.Infof("OCR3 config already set on offramp for chain %d", remote)
			continue
		}

		var instructions []solana.Instruction
		offRampConfigPDA := state.SolChains[remote].OffRampConfigPDA
		offRampStatePDA := state.SolChains[remote].OffRampStatePDA
		solOffRamp.SetProgramID(state.SolChains[remote].OffRamp)
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

func isOCR3ConfigSetOnOffRampSolana(
	e deployment.Environment,
	chain deployment.SolChain,
	chainState cs.SolCCIPChainState,
	args []internal.MultiOCR3BaseOCRConfigArgsSolana,
) (bool, error) {
	var configAccount solOffRamp.Config
	err := chain.GetAccountDataBorshInto(e.GetContext(), chainState.OffRampConfigPDA, &configAccount)
	if err != nil {
		return false, fmt.Errorf("failed to get account info: %w", err)
	}
	for _, newState := range args {
		existingState := configAccount.Ocr3[newState.OCRPluginType]
		if existingState.ConfigInfo.ConfigDigest != newState.ConfigDigest {
			e.Logger.Infof("OCR3 config digest mismatch")
			return false, nil
		}
		if existingState.ConfigInfo.F != newState.F {
			e.Logger.Infof("OCR3 config F mismatch")
			return false, nil
		}
		if existingState.ConfigInfo.IsSignatureVerificationEnabled != btoi(newState.IsSignatureVerificationEnabled) {
			e.Logger.Infof("OCR3 config signature verification mismatch")
			return false, nil
		}
		if newState.OCRPluginType == OcrCommitPlugin {
			// only commit will set signers, exec doesn't need them.
			if len(existingState.Signers) != len(newState.Signers) {
				e.Logger.Infof("OCR3 config signers length mismatch")
				return false, nil
			}
			for i := 0; i < len(existingState.Signers); i++ {
				if existingState.Signers[i] != newState.Signers[i] {
					e.Logger.Infof("OCR3 config signers mismatch")
					return false, nil
				}
			}
		}
		if len(existingState.Transmitters) != len(newState.Transmitters) {
			e.Logger.Infof("OCR3 config transmitters length mismatch")
			return false, nil
		}
		for i := 0; i < len(existingState.Transmitters); i++ {
			if existingState.Transmitters[i] != newState.Transmitters[i] {
				e.Logger.Infof("OCR3 config transmitters mismatch")
				return false, nil
			}
		}
	}
	return true, nil
}
