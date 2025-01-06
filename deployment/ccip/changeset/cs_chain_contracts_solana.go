package changeset

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
)

// TODO:
// func (cfg UpdateOnRampDestsConfig) Validate(e deployment.Environment) error {
// 	state, err := LoadOnchainState(e)
// 	if err != nil {
// 		return err
// 	}
// 	supportedChains := state.SupportedChains()
// 	for chainSel, updates := range cfg.UpdatesByChain {
// 		chainState, ok := state.Chains[chainSel]
// 		if !ok {
// 			return fmt.Errorf("chain %d not found in onchain state", chainSel)
// 		}
// 		if chainState.TestRouter == nil {
// 			return fmt.Errorf("missing test router for chain %d", chainSel)
// 		}
// 		if chainState.Router == nil {
// 			return fmt.Errorf("missing router for chain %d", chainSel)
// 		}
// 		if chainState.OnRamp == nil {
// 			return fmt.Errorf("missing onramp onramp for chain %d", chainSel)
// 		}
// 		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.Chains[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.OnRamp); err != nil {
// 			return err
// 		}

// 		for destination := range updates {
// 			// Destination cannot be an unknown destination.
// 			if _, ok := supportedChains[destination]; !ok {
// 				return fmt.Errorf("destination chain %d is not a supported %s", destination, chainState.OnRamp.Address())
// 			}
// 			sc, err := chainState.OnRamp.GetStaticConfig(&bind.CallOpts{Context: e.GetContext()})
// 			if err != nil {
// 				return fmt.Errorf("failed to get onramp static config %s: %w", chainState.OnRamp.Address(), err)
// 			}
// 			if destination == sc.ChainSelector {
// 				return fmt.Errorf("cannot update onramp destination to the same chain")
// 			}
// 		}
// 	}
// 	return nil
// }

// UpdateOnRampsDests updates the onramp destinations for each onramp
// in the chains specified. Multichain support is important - consider when we add a new chain
// and need to update the onramp destinations for all chains to support the new chain.
func UpdateOnRampsDestsSolana(e deployment.Environment, cfg UpdateOnRampDestsConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	s, err := LoadOnchainStateSolana(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	for chainSel, updates := range cfg.UpdatesByChain {
		e.Logger.Infow("Updating onramp destinations", "chain", chainSel, "updates", updates)
		chain := e.SolChains[chainSel]

		validSourceChainConfig := ccip_router.SourceChainConfig{
			OnRamp:    []byte{1, 2, 3},
			IsEnabled: true,
		}

		ccipRouterID := s.SolChains[chainSel].CcipRouter
		// ccip_router.SetProgramID(ccipRouterId) //cannot set this again

		for destination, update := range updates {
			EvmSourceChainStatePDA := GetEvmSourceChainStatePDA(ccipRouterID, destination)
			e.Logger.Infow("EvmSourceChainStatePDA", "EvmSourceChainStatePDA", EvmSourceChainStatePDA)
			EvmDestChainStatePDA := GetEvmDestChainStatePDA(ccipRouterID, destination)
			validDestChainConfig := ccip_router.DestChainConfig{
				IsEnabled: update.IsEnabled,

				// minimal valid config
				DefaultTxGasLimit:       1,
				MaxPerMsgGasLimit:       100,
				MaxDataBytes:            32,
				MaxNumberOfTokensPerMsg: 1,
				// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
				ChainFamilySelector: [4]uint8{40, 18, 213, 44},
			}

			instruction, err := ccip_router.NewAddChainSelectorInstruction(
				destination,
				validSourceChainConfig,
				validDestChainConfig,
				EvmSourceChainStatePDA,
				EvmDestChainStatePDA,
				GetRouterConfigPDA(ccipRouterID),
				chain.DeployerKey.PublicKey(),
				solana.SystemProgramID,
			).ValidateAndBuild()

			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate instructions: %v", err)
			}

			err = chain.Confirm([]solana.Instruction{instruction})

			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %v", err)
			} else {
				e.Logger.Infow("Confirmed instruction", "instruction", instruction)
			}
		}
	}

	return deployment.ChangesetOutput{}, nil
}

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
func SetOCR3ConfigSolana(e deployment.Environment, cfg SetOCR3OffRampConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	state, err := LoadOnchainState(e)
	if err != nil {
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
			return deployment.ChangesetOutput{}, err
		}
		args, err := internal.BuildSetOCR3ConfigArgsSolana(donID, state.Chains[cfg.HomeChainSel].CCIPHome, remote)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		// set, err := isOCR3ConfigSetSolana(e.Logger, e.Chains[remote], state.Chains[remote].OffRamp, args)
		// if err != nil {
		// 	return deployment.ChangesetOutput{}, err
		// }
		// if set {
		// 	e.Logger.Infof("OCR3 config already set on offramp for chain %d", remote)
		// 	continue
		// }
		var instructions []solana.Instruction
		ccipRouterId := solChains[remote].CcipRouter
		for _, arg := range args {
			instruction, err := ccip_router.NewSetOcrConfigInstruction(
				uint8(arg.OcrPluginType),
				ccip_router.Ocr3ConfigInfo{
					ConfigDigest:                   arg.ConfigDigest,
					F:                              arg.F,
					IsSignatureVerificationEnabled: uint8(btoi(arg.IsSignatureVerificationEnabled)),
				},
				arg.Signers,
				arg.Transmitters,
				GetRouterConfigPDA(ccipRouterId),
				GetRouterStatePDA(ccipRouterId),
				e.SolChains[remote].DeployerKey.PublicKey(),
			).ValidateAndBuild()
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			instructions = append(instructions, instruction)
		}
		if cfg.MCMS == nil {
			err := e.SolChains[remote].Confirm(instructions)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
		}
	}

	return deployment.ChangesetOutput{}, nil

	// var batches []timelock.BatchChainOperation
	// timelocks := make(map[uint64]common.Address)
	// proposers := make(map[uint64]*mcm.MCM)
	// else {
	// 	batches = append(batches, timelock.BatchChainOperation{
	// 		ChainIdentifier: mcms.ChainIdentifier(remote),
	// 		Batch: []mcms.Operation{
	// 			{
	// 				To:    offRamp.Address(),
	// 				Data:  tx.Data(),
	// 				Value: big.NewInt(0),
	// 			},
	// 		},
	// 	})
	// 	timelocks[remote] = state.Chains[remote].Timelock.Address()
	// 	proposers[remote] = state.Chains[remote].ProposerMcm
	// }
	// p, err := proposalutils.BuildProposalFromBatches(
	// 	timelocks,
	// 	proposers,
	// 	batches,
	// 	"Update OCR3 config",
	// 	cfg.MCMS.MinDelay,
	// )
	// if err != nil {
	// 	return deployment.ChangesetOutput{}, err
	// }
	// e.Logger.Infof("Proposing OCR3 config update for", cfg.RemoteChainSels)
	// return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{
	// 	*p,
	// }}, nil

}
