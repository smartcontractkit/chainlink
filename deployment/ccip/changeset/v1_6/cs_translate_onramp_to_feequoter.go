package v1_6

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"

	migrate_seq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/migration"
	ccipseqs "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	TranslateEVM2EVMOnRampsToFQDestConfig          = cldf.CreateChangeSet(TranslateEVM2EVMOnRampsToFeeQuoterChangeset, ValidatePreReqContractsInState)
	TranslateEVM2EVMOnRampsToFQTokenTransferConfig = cldf.CreateChangeSet(TranslateEVM2EVMOnRampsToFeeQTokenTransferFeeConfigChangeset, ValidatePreReqContractsInState)
)

type TranslateEVM2EVMOnRampsToFeeQuoterConfig struct {
	DestChainSelector           uint64
	NewFeeQuoterParamsPerSource map[uint64]migrate_seq.NewFeeQuoterDestChainConfigParams
	MCMS                        *proposalutils.TimelockConfig
}

func (cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) Validate(e cldf.Environment) error {
	if err := cldf.IsValidChainSelector(cfg.DestChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", cfg.DestChainSelector, err)
	}

	return nil
}

func ValidatePreReqContractsInState(e cldf.Environment, cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) error {
	if err := cfg.Validate(e); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	state, err := stateview.LoadOnchainState(e, stateview.WithLoadLegacyContracts(true))
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for sourceChainSel, sourceChain := range state.Chains {
		if sourceChainSel == cfg.DestChainSelector {
			continue // Skip the destination chain, we are not processing OnRamps for it.
		}

		var hasRampToDest bool
		for destChainOnSrc := range sourceChain.EVM2EVMOnRamp {
			if _, ok := cfg.NewFeeQuoterParamsPerSource[sourceChainSel]; destChainOnSrc == cfg.DestChainSelector && !ok {
				return fmt.Errorf("no new FeeQuoter params found for source chain %d to destination chain %d", sourceChainSel, cfg.DestChainSelector)
			} else if destChainOnSrc == cfg.DestChainSelector {
				hasRampToDest = true
			}
		}

		if hasRampToDest {
			if sourceChain.PriceRegistry == nil {
				return fmt.Errorf("priceRegistry not found for source chain %d, cannot process 1.5.0 OnRamps", sourceChainSel)
			}
			if sourceChain.TokenAdminRegistry == nil {
				return fmt.Errorf("tokenAdminRegistry not found for source chain %d, cannot process 1.5.0 OnRamps", sourceChainSel)
			}
		}
	}

	return nil
}

func TranslateEVM2EVMOnRampsToFeeQuoterChangeset(e cldf.Environment, cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) (cldf.ChangesetOutput, error) {
	csOutput := cldf.ChangesetOutput{}
	state, err := stateview.LoadOnchainState(e, stateview.WithLoadLegacyContracts(true))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Translate the 1.5.0 OnRamp to the FeeQuoterDestChainConfig
	translateDynamicCfgReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		migrate_seq.SeqTranslateOnRampToFeeQDestConfig,
		e.BlockChains.EVMChains(),
		cfg.toSequenceInput(e, state),
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to translate 1.5.0 OnRamp dynamic config: %w", err)
	}

	// Update the FeeQuoterDestChainConfig on the FeeQuoter contract with translated configs
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
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute FeeQuoterApplyDestChainConfigUpdatesSequence: %w", err)
	}
	csOutput, err = opsutil.AddEVMCallSequenceToCSOutput(e, csOutput, report, err, state.EVMMCMSStateByChain(), cfg.MCMS, "Call ApplyDestChainConfigUpdates on FeeQuoter")
	if err != nil {
		return csOutput, fmt.Errorf("failed to apply FeeQuoter dest chain config updates: %w", err)
	}

	// ApplyFeeTokensUpdates to add fee tokens on FeeQ with translated configs
	feeTokensReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.FeeQuoterApplyFeeTokensUpdatesSeq,
		e.BlockChains.EVMChains(),
		cfg.toFeeTokenApplySeqInput(state, translateDynamicCfgReport.Output.FeeTokens),
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute FeeQuoterApplyFeeTokensUpdatesSeq: %w", err)
	}
	csOutput, err = opsutil.AddEVMCallSequenceToCSOutput(e, csOutput, feeTokensReport, err, state.EVMMCMSStateByChain(), cfg.MCMS, "Call ApplyFeeTokensUpdatesConfig on FeeQuoter")
	if err != nil {
		return csOutput, fmt.Errorf("failed to apply FeeQuoter fee tokens updates: %w", err)
	}

	// applyPremiumMultiplierWeiPerEthUpdates to add premiumMultiplier Cfg on FeeQ with translated configs
	premiumMultiplierSqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.FeeQApplyPremiumMultiplierWeiPerEthUpdatesSeq,
		e.BlockChains.EVMChains(),
		cfg.toPremiumMultiplierCfgSeqInput(state, translateDynamicCfgReport.Output.FeeTokenPremiumMultipliers),
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute FeeQApplyPremiumMultiplierWeiPerEthUpdatesSeq: %w", err)
	}
	csOutput, err = opsutil.AddEVMCallSequenceToCSOutput(e, csOutput, premiumMultiplierSqReport, err, state.EVMMCMSStateByChain(), cfg.MCMS, "Call ApplyPremiumMultiplierWeiPerEthUpdates on FeeQuoter")
	if err != nil {
		return csOutput, fmt.Errorf("failed to apply FeeQuoter premium Multiplier config updates: %w", err)
	}
	return csOutput, nil
}

func (cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) toSequenceInput(e cldf.Environment, state stateview.CCIPOnChainState) migrate_seq.FeeQuoterUpdateTokenTransferConfig {
	input := make(map[uint64]opsutil.EVMCallInput[migrate_seq.OnRampToFeeQuoterDestChainConfigInput])
	for sel := range e.BlockChains.EVMChains() {
		onRamps := make(map[uint64]common.Address)
		newFeeQuoterParams := make(map[uint64]migrate_seq.NewFeeQuoterDestChainConfigParams)
		srcChainState := state.Chains[sel]
		for destChainSel, onRamp1_5 := range srcChainState.EVM2EVMOnRamp {
			if destChainSel == cfg.DestChainSelector && onRamp1_5 != nil {
				onRamps[destChainSel] = onRamp1_5.Address()
				newFeeQuoterParams[destChainSel] = cfg.NewFeeQuoterParamsPerSource[sel]
			}
		}
		input[sel] = opsutil.EVMCallInput[migrate_seq.OnRampToFeeQuoterDestChainConfigInput]{
			ChainSelector: sel,
			Address:       state.Chains[sel].FeeQuoter.Address(),
			CallInput: migrate_seq.OnRampToFeeQuoterDestChainConfigInput{
				NewFeeQuoterParams: newFeeQuoterParams,
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

func (cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) toFeeTokenApplySeqInput(state stateview.CCIPOnChainState, tokens map[uint64][]common.Address) ccipseqs.FeeQuoterUpdateFeeTokensConfig {
	input := make(map[uint64]opsutil.EVMCallInput[ccipops.ApplyFeeTokensUpdatesInput], len(tokens))

	for chainSel, tokens := range tokens {
		var tokensToRemove, tokensToAdd []common.Address
		tokensToAdd = append(tokensToAdd, tokens...)
		input[chainSel] = opsutil.EVMCallInput[ccipops.ApplyFeeTokensUpdatesInput]{
			ChainSelector: chainSel,
			Address:       state.Chains[chainSel].FeeQuoter.Address(),
			CallInput: ccipops.ApplyFeeTokensUpdatesInput{
				FeeTokensToAdd:    tokensToAdd,
				FeeTokensToRemove: tokensToRemove,
			},
			NoSend: cfg.MCMS != nil,
		}
	}
	return ccipseqs.FeeQuoterUpdateFeeTokensConfig{
		UpdatesByChain: input,
	}
}

func (cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) toPremiumMultiplierCfgSeqInput(state stateview.CCIPOnChainState, tokenPremiumCfgs map[uint64][]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs) ccipseqs.FeeQuoterUpdatePremiumMultiplierWeiPerEthConfig {
	input := make(map[uint64]opsutil.EVMCallInput[[]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs], len(tokenPremiumCfgs))

	for chainSel, updates := range tokenPremiumCfgs {
		var premiumMultiplierUpdates []fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs
		for _, update := range updates {
			premiumMultiplierUpdates = append(premiumMultiplierUpdates, fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs{
				Token:                      update.Token,
				PremiumMultiplierWeiPerEth: update.PremiumMultiplierWeiPerEth,
			})
		}
		input[chainSel] = opsutil.EVMCallInput[[]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs]{
			ChainSelector: chainSel,
			Address:       state.Chains[chainSel].FeeQuoter.Address(),
			CallInput:     premiumMultiplierUpdates,
			NoSend:        cfg.MCMS != nil, // If MCMS exists, we do not want to send the transaction.
		}
	}
	return ccipseqs.FeeQuoterUpdatePremiumMultiplierWeiPerEthConfig{
		UpdatesByChain: input,
	}
}

func TranslateEVM2EVMOnRampsToFeeQTokenTransferFeeConfigChangeset(e cldf.Environment, cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) (cldf.ChangesetOutput, error) {
	csOutput := cldf.ChangesetOutput{}
	state, err := stateview.LoadOnchainState(e, stateview.WithLoadLegacyContracts(true))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	// Translate the 1.5.0 OnRamp token transfer fee configs to FeeQuoterTokenTransferFeeConfig
	translateTokenTransferFeeCfgReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		migrate_seq.SeqTranslateOnRampToFeeQTokenTransferFeeCfg,
		e.BlockChains.EVMChains(),
		cfg.toSequenceInput(e, state),
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to translate 1.5.0 OnRamp dynamic config: %w", err)
	}

	// ApplyTokenTransferFeeConfigUpdates on the FeeQuoter contract
	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.FeeQUpdateTransferTokenFeeCfgSeq,
		e.BlockChains.EVMChains(),
		cfg.tokenTransferFeeConfigArgsToSeqInput(state, translateTokenTransferFeeCfgReport.Output.FeeQuoterUpdates),
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute FeeQUpdateTransferTokenFeeCfgSeq: %w", err)
	}

	// TODO: check the output: its empty
	csOutput, err = opsutil.AddEVMCallSequenceToCSOutput(e, csOutput, report, err, state.EVMMCMSStateByChain(), cfg.MCMS, "Call ApplyTokenTransferFeeConfigUpdates on FeeQuoter")
	if err != nil {
		return csOutput, fmt.Errorf("failed to apply FeeQuoter fee tokens updates: %w", err)
	}
	return csOutput, nil
}

func (cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) tokenTransferFeeConfigArgsToSeqInput(state stateview.CCIPOnChainState, tokenTransferFeeCfgArgs map[uint64][]fee_quoter.FeeQuoterTokenTransferFeeConfigArgs) ccipseqs.FeeQuoterUpdateTokenTransferConfig {
	input := make(map[uint64]opsutil.EVMCallInput[ccipops.ApplyTokenTransferFeeConfigUpdatesConfigPerChain])
	for chainSel, tokensFeeCfgArgs := range tokenTransferFeeCfgArgs {
		input[chainSel] = opsutil.EVMCallInput[ccipops.ApplyTokenTransferFeeConfigUpdatesConfigPerChain]{
			ChainSelector: chainSel,
			Address:       state.Chains[chainSel].FeeQuoter.Address(),
			CallInput: ccipops.ApplyTokenTransferFeeConfigUpdatesConfigPerChain{
				TokenTransferFeeConfigs:       tokensFeeCfgArgs,
				TokenTransferFeeConfigsRemove: nil, // not removing any token transfer configs for now
			},
			NoSend: cfg.MCMS != nil,
		}
	}

	return ccipseqs.FeeQuoterUpdateTokenTransferConfig{
		UpdatesByChain: input,
	}
}
