package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"

	migrate_seq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/migration"
	ccipseqs "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	_ cldf.ChangeSet[TranslateEVM2EVMOnRampsToFeeQuoterConfig] = TranslateEVM2EVMOnRampsToFeeQuoterChangeset
)

type TranslateEVM2EVMOnRampsToFeeQuoterConfig struct {
	SourceChainSelectors []uint64
	MCMS                 *proposalutils.TimelockConfig
}

func (cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) Validate(e cldf.Environment) error {
	for _, cs := range cfg.SourceChainSelectors {
		if err := cldf.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
	}

	return nil
}

func ValidatePreReqContractsInState(srcChain cldf_evm.Chain, srcChainState evm.CCIPChainState) error {
	if len(srcChainState.EVM2EVMOnRamp) == 0 {
		return fmt.Errorf("No 1.5.0 OnRamps found, skipping (chainSelector: %d, chainName: %s)", srcChain.Selector, srcChain.Name())
	}
	if srcChainState.PriceRegistry == nil {
		return fmt.Errorf("PriceRegistry not found for source chain %d (%s), cannot process 1.5.0 OnRamps", srcChain.Selector, srcChain.Name())
	}
	if srcChainState.TokenAdminRegistry == nil {
		return fmt.Errorf("TokenAdminRegistry not found for source chain %d (%s), cannot process 1.5.0 OnRamps", srcChain.Selector, srcChain.Name())
	}

	return nil
}

func TranslateEVM2EVMOnRampsToFeeQuoterChangeset(e cldf.Environment, cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) (cldf.ChangesetOutput, error) {
	output := cldf.ChangesetOutput{}
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid config: %w", err)
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	if len(cfg.SourceChainSelectors) > 0 {
		for _, sel := range cfg.SourceChainSelectors {
			if chain, ok := e.BlockChains.EVMChains()[sel]; ok {
				srcChainState := state.MustGetEVMChainState(sel)
				err := ValidatePreReqContractsInState(chain, srcChainState)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate pre-requisite contracts in state for source chain %d: %w", sel, err)
				}
			} else {
				return cldf.ChangesetOutput{}, fmt.Errorf("Chain selector not found in environment, skipping (chainSelector: %d)", sel)
			}
		}
	} else {
		return cldf.ChangesetOutput{}, fmt.Errorf("No chains to process TranslateEVM2EVMOnRampsToFeeQuoter")
	}

	seqInput := cfg.ToSequenceInput(state)
	// Translate the 1.5.0 OnRamp to the FeeQuoterDestChainConfig
	translateDynamicCfgReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		migrate_seq.SeqTranslateOnRampToFeeQDestConfig,
		e.BlockChains.EVMChains(),
		seqInput,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to translate 1.5.0 OnRamp dynamic config: %w", err)
	}

	// Update the FeeQuoterDestChainConfig on the FeeQuoter contract
	updateFeeQuoterDestsConfig := UpdateFeeQuoterDestsConfig{
		UpdatesByChain: translateDynamicCfgReport.Output.FeeQuoterUpdates,
		MCMS:           cfg.MCMS,
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.FeeQuoterApplyDestChainConfigUpdatesSequence,
		e.BlockChains.EVMChains(),
		updateFeeQuoterDestsConfig.ToSequenceInput(state),
	)

	csOutput, err := opsutil.AddEVMCallSequenceToCSOutput(e, state, output, report, err, cfg.MCMS, "Call ApplyDestChainConfigUpdates on FeeQuoters")
	if err != nil {
		return csOutput, fmt.Errorf("failed to apply FeeQuoter dest chain config updates: %w", err)
	}
	return csOutput, nil
}

func (cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) ToSequenceInput(state stateview.CCIPOnChainState) migrate_seq.FeeQuoterUpdateTokenTransferConfig {
	input := make(map[uint64]opsutil.EVMCallInput[migrate_seq.OnRampToFeeQuoterDestChainConfigInput], len(cfg.SourceChainSelectors))
	for _, sel := range cfg.SourceChainSelectors {
		srcChainState := state.Chains[sel]
		onRamps := make(map[uint64]common.Address, len(srcChainState.EVM2EVMOnRamp))
		// Iterate over all onRamps found for the Source Chain
		for destChainSel, onRamp1_5 := range srcChainState.EVM2EVMOnRamp {
			if onRamp1_5 != nil {
				onRamps[destChainSel] = onRamp1_5.Address()
			}
		}
		input[sel] = opsutil.EVMCallInput[migrate_seq.OnRampToFeeQuoterDestChainConfigInput]{
			ChainSelector: sel,
			Address:       state.Chains[sel].FeeQuoter.Address(),
			CallInput: migrate_seq.OnRampToFeeQuoterDestChainConfigInput{
				OnRamps:            onRamps,
				TokenAdminRegistry: srcChainState.TokenAdminRegistry.Address(),
			},
			NoSend: cfg.MCMS != nil,
		}
	}
	return migrate_seq.FeeQuoterUpdateTokenTransferConfig{
		UpdatesByChain: input,
	}
}

/* // Translate the 1.5.0 OnRamp to the FeeQuoterDestChainConfig
	destChainEVM2EVMDynamicCfg, err2 := onRamp1_5.GetDynamicConfig(nil)
	if err2 != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("Cannot GetDynamicConfig on source chain %d, to destination Chain: %d, for 1.5.0 OnRamp %d: %w", srcChainSel, destChainSel, onRamp1_5.Address(), err2)
	}

	priceRegistry, err := price_registry.NewPriceRegistry(destChainEVM2EVMDynamicCfg.PriceRegistry, srcChain.Client)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create PriceRegistry contract binding on source chain %d: %w", srcChainSel, err)
	}

	allFeeTokens, err := priceRegistry.GetFeeTokens(nil)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get fee tokens from PriceRegistry on source chain %d: %w", srcChainSel, err)
	}

	for _, ft := range allFeeTokens {
		ftCfg, err := onRamp1_5.GetFeeTokenConfig(nil, ft)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get fee token config for %s on source chain %d: %w", ft.Hex(), srcChainSel, err)
		}

	}

	// Port token transfer fee config args from all 1.5.0 OnRamps into FeeQuoter
	// get tokenAdminregistry from state
	//	-> get alltokens from tokenAdminReg.getAllConfiguredTokens
	// 	-> for each token
	// 		-->  get tokenpool on onRamp.getPoolBySourceToken(selector, token)
	//		--> supportedChains := tokenPool.getSupportedChains()
	//		--> if targetSelector not in supportedChains then continue
	// 		--> add this token to validTokens to process
	tokenar := srcChainState.TokenAdminRegistry
	// Example: create a slice of custom struct inline
	type tokenAddress struct {
		Address common.Address
	}
	allValidTokensAndCfgs := make([]fee_quoter.FeeQuoterTokenTransferFeeConfigSingleTokenArgs, 0)
	allTokens, err := tokenar.GetAllConfiguredTokens(nil, 0, 1000) // TODO: Adjust the pagination as needed
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get all configured tokens from TokenAdminRegistry on source chain %d: %w", srcChainSel, err)
	}
	for _, token := range allTokens {
		tokenPoolAddress, err := onRamp1_5.GetPoolBySourceToken(nil, destChainSel, token)
		if err != nil {
			lggr.Warnw("Failed to get pool for token on 1.5.0 OnRamp", "sourceChainSelector", srcChainSel, "destinationChainSelector", destChainSel, "token", token.Hex(), "error", err)
			continue // continue or exit?
		}

		tokenPool, err := token_pool.NewTokenPool(tokenPoolAddress, srcChain.Client)
		if err != nil {
			lggr.Warnw("Failed to create tokenpool contract binding", "sourceChainSelector", srcChainSel, "destinationChainSelector", destChainSel, "token", token.Hex(), "error", err)
			continue // continue or exit?
		}

		supportedChains, err := tokenPool.GetSupportedChains(nil)
		if err != nil {
			lggr.Warnw("Failed to get supported chains from token pool", "sourceChainSelector", srcChainSel, "destinationChainSelector", destChainSel, "token", token.Hex(), "error", err)
			continue // continue or exit?
		}

		// Check if the destination chain selector is in the supported chains of token pool
		found := false
		for _, chainSel := range supportedChains {
			if chainSel == destChainSel {
				found = true
				break
			}
		}
		if !found {
			continue // skip this token if the destination chain is not supported
		}
		tokenTransferFeeCfg, err2 := onRamp1_5.GetTokenTransferFeeConfig(nil, token)
		if err2 != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("Cannot GetTokenTransferFeeConfig on source chain %d, to destination Chain: %d, for 1.5.0 OnRamp %d and token %s: %w", srcChainSel, destChainSel, onRamp1_5.Address(), token.Hex(), err2)
		}
		allValidTokensAndCfgs = append(allValidTokensAndCfgs,
			migrateOnRamp.TranslateOnrampToFeequoterTokenTransferFeeConfig(destChainSel, token, tokenTransferFeeCfg),
		)
	}

	// fix this

	_ = fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{
		DestChainSelector:       destChainSel,
		TokenTransferFeeConfigs: allValidTokensAndCfgs,
	}

}

// This needs to be a seq
/* tx, err := v1_6.ApplyFeeTokensUpdatesFeeQuoterChangeset(e, v1_6.ApplyFeeTokensUpdatesConfig{
	UpdatesByChain: tokenTransferFeeConfigArgs,
	MCMS:           cfg.MCMS,
}) */
/*
	if len(feeQuoterUpdates) == 0 {
		lggr.Info("No FeeQuoter updates found from 1.5.0 OnRamps, skipping changeset")
		return cldf.ChangesetOutput{}, nil
	}

	// this needs to be a seq
	lggr.Infow("Applying translated configurations to FeeQuoters via v1_6.UpdateFeeQuoterDestsChangeset", "updatesCount", len(feeQuoterUpdates))
	return v1_6.UpdateFeeQuoterDestsChangeset(e, v1_6.UpdateFeeQuoterDestsConfig{
		UpdatesByChain: feeQuoterUpdates,
		MCMS:           cfg.MCMS,
	}) */
// }

/*
func TranslateEVM2EVMOnRampsToFeeQTransferTokensChangeset(e cldf.Environment, cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) (cldf.ChangesetOutput, error) {
	output := cldf.ChangesetOutput{}
	seqReports := make([]operations.Report[any, any], 0)
	var allProposals []mcmslib.TimelockProposal

	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid config: %w", err)
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	lggr := e.Logger
	// todo: preallocate for the number of chains?
	_ = make(map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig)
	// tokenTransferFeeConfigUpdates := make(map[uint64]fee_quoter.FeeQuoterTokenTransferFeeConfigArgs)
	chainsToProcess := make(map[uint64]cldf_evm.Chain, len(cfg.SourceChainSelectors))

	if len(cfg.SourceChainSelectors) > 0 {
		for _, sel := range cfg.SourceChainSelectors {
			if chain, ok := e.BlockChains.EVMChains()[sel]; ok {
				chainsToProcess[sel] = chain
			} else {
				lggr.Warnw("Chain selector not found in environment, skipping", "chainSelector", sel)
			}
		}
	} else {
		return cldf.ChangesetOutput{}, fmt.Errorf("No chains to process TranslateEVM2EVMOnRampsToFeeQuoter")
	}

	for srcChainSel, srcChain := range chainsToProcess {
		srcChainState := state.MustGetEVMChainState(srcChainSel)
		err := ValidatePreReqContractsInState(srcChain, srcChainState)
		if err != nil {
			lggr.Warnw("failed to validate pre-requisite contracts in state for source chain %d: %w", srcChainSel, err)
			continue // continue or exit?
		}
		onRamps := make(map[uint64]common.Address, len(srcChainState.EVM2EVMOnRamp))
		// Iterate over all onRamps found for the Source Chain
		for destChainSel, onRamp1_5 := range srcChainState.EVM2EVMOnRamp {
			if onRamp1_5 != nil {
				onRamps[destChainSel] = onRamp1_5.Address()
			}
		}
		// Translate the 1.5.0 OnRamp to the FeeQuoterDestChainConfig
		transferTokensCfgReport, err := operations.ExecuteSequence(
			e.OperationsBundle,
			migrate_seq.SeqTranslateOnRampToFeeQTokenTransferFeeCfg,
			migrate_seq.TranslateOnRampToFeeQuoterDestChainConfigDeps{
				Chain: srcChain,
			},
			migrate_seq.TranslateOnRampToFeeQuoterDestChainConfigInput{
				SourceChainSelector: srcChainSel,
				OnRamps:             onRamps,
				TokenAdminRegistry:  srcChainState.TokenAdminRegistry.Address(),
			},
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to translate 1.5.0 OnRamp dynamic config for source chain %d: %w", srcChainSel, err)
		}

		// Update the FeeQuoterDestChainConfig on the FeeQuoter contract
		updateFeeQuoterDestsConfig := UpdateFeeQuoterDestsConfig{
			UpdatesByChain: transferTokensCfgReport.Output,
			MCMS:           cfg.MCMS,
		}

		report, err := operations.ExecuteSequence(
			e.OperationsBundle,
			ccipseqs.FeeQuoterApplyDestChainConfigUpdatesSequence,
			e.BlockChains.EVMChains(),
			updateFeeQuoterDestsConfig.ToSequenceInput(state),
		)

		csOutput, err := opsutil.AddEVMCallSequenceToCSOutput(e, state, output, report, err, cfg.MCMS, "Call ApplyDestChainConfigUpdates on FeeQuoters")
		if err != nil {
			return csOutput, fmt.Errorf("failed to apply FeeQuoter dest chain config updates for source chain %d: %w", srcChainSel, err)
		}
		seqReports = append(seqReports, csOutput.Reports...)
		allProposals = append(allProposals, csOutput.MCMSTimelockProposals...)
	}

	proposal, err := proposalutils.AggregateProposals(
		e,
		state.EVMMCMSStateByChain(),
		nil,
		allProposals,
		"FeeQuoterApplyDestChainConfigUpdatesSequence migration",
		cfg.MCMS,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	if proposal == nil {
		return cldf.ChangesetOutput{Reports: seqReports}, nil
	}
	return cldf.ChangesetOutput{Reports: seqReports, MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil */
/* // Translate the 1.5.0 OnRamp to the FeeQuoterDestChainConfig
	destChainEVM2EVMDynamicCfg, err2 := onRamp1_5.GetDynamicConfig(nil)
	if err2 != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("Cannot GetDynamicConfig on source chain %d, to destination Chain: %d, for 1.5.0 OnRamp %d: %w", srcChainSel, destChainSel, onRamp1_5.Address(), err2)
	}

	priceRegistry, err := price_registry.NewPriceRegistry(destChainEVM2EVMDynamicCfg.PriceRegistry, srcChain.Client)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create PriceRegistry contract binding on source chain %d: %w", srcChainSel, err)
	}

	allFeeTokens, err := priceRegistry.GetFeeTokens(nil)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get fee tokens from PriceRegistry on source chain %d: %w", srcChainSel, err)
	}

	for _, ft := range allFeeTokens {
		ftCfg, err := onRamp1_5.GetFeeTokenConfig(nil, ft)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get fee token config for %s on source chain %d: %w", ft.Hex(), srcChainSel, err)
		}

	}

	// Port token transfer fee config args from all 1.5.0 OnRamps into FeeQuoter
	// get tokenAdminregistry from state
	//	-> get alltokens from tokenAdminReg.getAllConfiguredTokens
	// 	-> for each token
	// 		-->  get tokenpool on onRamp.getPoolBySourceToken(selector, token)
	//		--> supportedChains := tokenPool.getSupportedChains()
	//		--> if targetSelector not in supportedChains then continue
	// 		--> add this token to validTokens to process
	tokenar := srcChainState.TokenAdminRegistry
	// Example: create a slice of custom struct inline
	type tokenAddress struct {
		Address common.Address
	}
	allValidTokensAndCfgs := make([]fee_quoter.FeeQuoterTokenTransferFeeConfigSingleTokenArgs, 0)
	allTokens, err := tokenar.GetAllConfiguredTokens(nil, 0, 1000) // TODO: Adjust the pagination as needed
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get all configured tokens from TokenAdminRegistry on source chain %d: %w", srcChainSel, err)
	}
	for _, token := range allTokens {
		tokenPoolAddress, err := onRamp1_5.GetPoolBySourceToken(nil, destChainSel, token)
		if err != nil {
			lggr.Warnw("Failed to get pool for token on 1.5.0 OnRamp", "sourceChainSelector", srcChainSel, "destinationChainSelector", destChainSel, "token", token.Hex(), "error", err)
			continue // continue or exit?
		}

		tokenPool, err := token_pool.NewTokenPool(tokenPoolAddress, srcChain.Client)
		if err != nil {
			lggr.Warnw("Failed to create tokenpool contract binding", "sourceChainSelector", srcChainSel, "destinationChainSelector", destChainSel, "token", token.Hex(), "error", err)
			continue // continue or exit?
		}

		supportedChains, err := tokenPool.GetSupportedChains(nil)
		if err != nil {
			lggr.Warnw("Failed to get supported chains from token pool", "sourceChainSelector", srcChainSel, "destinationChainSelector", destChainSel, "token", token.Hex(), "error", err)
			continue // continue or exit?
		}

		// Check if the destination chain selector is in the supported chains of token pool
		found := false
		for _, chainSel := range supportedChains {
			if chainSel == destChainSel {
				found = true
				break
			}
		}
		if !found {
			continue // skip this token if the destination chain is not supported
		}
		tokenTransferFeeCfg, err2 := onRamp1_5.GetTokenTransferFeeConfig(nil, token)
		if err2 != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("Cannot GetTokenTransferFeeConfig on source chain %d, to destination Chain: %d, for 1.5.0 OnRamp %d and token %s: %w", srcChainSel, destChainSel, onRamp1_5.Address(), token.Hex(), err2)
		}
		allValidTokensAndCfgs = append(allValidTokensAndCfgs,
			migrateOnRamp.TranslateOnrampToFeequoterTokenTransferFeeConfig(destChainSel, token, tokenTransferFeeCfg),
		)
	}

	// fix this

	_ = fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{
		DestChainSelector:       destChainSel,
		TokenTransferFeeConfigs: allValidTokensAndCfgs,
	}

}

// This needs to be a seq
/* tx, err := v1_6.ApplyFeeTokensUpdatesFeeQuoterChangeset(e, v1_6.ApplyFeeTokensUpdatesConfig{
	UpdatesByChain: tokenTransferFeeConfigArgs,
	MCMS:           cfg.MCMS,
}) */
/*
	if len(feeQuoterUpdates) == 0 {
		lggr.Info("No FeeQuoter updates found from 1.5.0 OnRamps, skipping changeset")
		return cldf.ChangesetOutput{}, nil
	}

	// this needs to be a seq
	lggr.Infow("Applying translated configurations to FeeQuoters via v1_6.UpdateFeeQuoterDestsChangeset", "updatesCount", len(feeQuoterUpdates))
	return v1_6.UpdateFeeQuoterDestsChangeset(e, v1_6.UpdateFeeQuoterDestsConfig{
		UpdatesByChain: feeQuoterUpdates,
		MCMS:           cfg.MCMS,
	}) */
// }
