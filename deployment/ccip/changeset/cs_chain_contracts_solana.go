package changeset

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink/deployment"
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

		ccipRouterId := s.SolChains[chainSel].CcipRouter
		// ccip_router.SetProgramID(ccipRouterId) //cannot set this again

		for destination, update := range updates {
			EvmSourceChainStatePDA := GetEvmSourceChainStatePDA(ccipRouterId, destination)
			e.Logger.Infow("EvmSourceChainStatePDA", "EvmSourceChainStatePDA", EvmSourceChainStatePDA)
			EvmDestChainStatePDA := GetEvmDestChainStatePDA(ccipRouterId, destination)
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
				GetRouterConfigPDA(ccipRouterId),
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
