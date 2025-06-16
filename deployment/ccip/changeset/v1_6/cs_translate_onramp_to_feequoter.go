package v1_6

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
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
	csOutput := cldf.ChangesetOutput{}
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid config: %w", err)
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	if len(cfg.SourceChainSelectors) == 0 {
		return cldf.ChangesetOutput{}, errors.New("No chains to process TranslateEVM2EVMOnRampsToFeeQuoter")
	}
	for _, sel := range cfg.SourceChainSelectors {
		chain, ok := e.BlockChains.EVMChains()[sel]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("Chain selector not found in environment, skipping (chainSelector: %d)", sel)
		}
		srcChainState := state.MustGetEVMChainState(sel)
		err := ValidatePreReqContractsInState(chain, srcChainState)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate pre-requisite contracts in state for source chain %d: %w", sel, err)
		}
	}

	// Translate the 1.5.0 OnRamp to the FeeQuoterDestChainConfig
	translateDynamicCfgReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		migrate_seq.SeqTranslateOnRampToFeeQDestConfig,
		e.BlockChains.EVMChains(),
		cfg.toSequenceInput(state),
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
		return cldf.ChangesetOutput{}, fmt.Errorf("Failed to execute FeeQuoterApplyDestChainConfigUpdatesSequence: %w", err)
	}
	csOutput, err = opsutil.AddEVMCallSequenceToCSOutput(e, state, csOutput, report, err, cfg.MCMS, "Call ApplyDestChainConfigUpdates on FeeQuoters")
	if err != nil {
		return csOutput, fmt.Errorf("failed to apply FeeQuoter dest chain config updates: %w", err)
	}

	// ApplyFeeTokensUpdates to add fee tokens on FeeQ with translated configs
	feeTokensReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.FeeQuoterApplyFeeTokensUpdatesSeq,
		e.BlockChains.EVMChains(),
		cfg.toFeeTokenApplySeqInput(state, translateDynamicCfgReport.Output.Feetokens),
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("Failed to execute FeeQuoterApplyFeeTokensUpdatesSeq: %w", err)
	}
	csOutput, err = opsutil.AddEVMCallSequenceToCSOutput(e, state, csOutput, feeTokensReport, err, cfg.MCMS, "Call ApplyFeeTokensUpdatesConfig on FeeQuoter")
	if err != nil {
		return csOutput, fmt.Errorf("failed to apply FeeQuoter fee tokens updates: %w", err)
	}

	// applyPremiumMultiplierWeiPerEthUpdates to add premiumMultiplier Cfg on FeeQ with translated configs
	premiumMultiplierSqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.FeeQApplyPremiumMultiplierWeiPerEthUpdatesSeq,
		e.BlockChains.EVMChains(),
		cfg.toPremiumMultiplierCfgSequencInput(state, translateDynamicCfgReport.Output.FeeTokenPremiumMultipliers),
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("Failed to execute FeeQApplyPremiumMultiplierWeiPerEthUpdatesSeq: %w", err)
	}
	csOutput, err = opsutil.AddEVMCallSequenceToCSOutput(e, state, csOutput, premiumMultiplierSqReport, err, cfg.MCMS, "Call ApplyFeeTokensUpdatesConfig on FeeQuoter")
	if err != nil {
		return csOutput, fmt.Errorf("failed to apply FeeQuoter premium Multiplier config updates: %w", err)
	}
	return csOutput, nil
}

func (cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) toSequenceInput(state stateview.CCIPOnChainState) migrate_seq.FeeQuoterUpdateTokenTransferConfig {
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

func (cfg TranslateEVM2EVMOnRampsToFeeQuoterConfig) toPremiumMultiplierCfgSequencInput(state stateview.CCIPOnChainState, tokenPremiumCfgs map[uint64][]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs) ccipseqs.FeeQuoterUpdatePremiumMultiplierWeiPerEthConfig {
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
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid config: %w", err)
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	if len(cfg.SourceChainSelectors) == 0 {
		return cldf.ChangesetOutput{}, errors.New("No chains to process TranslateEVM2EVMOnRampsToFeeQuoter")
	}
	for _, sel := range cfg.SourceChainSelectors {
		chain, ok := e.BlockChains.EVMChains()[sel]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("Chain selector not found in environment, skipping (chainSelector: %d)", sel)
		}
		srcChainState := state.MustGetEVMChainState(sel)
		err := ValidatePreReqContractsInState(chain, srcChainState)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate pre-requisite contracts in state for source chain %d: %w", sel, err)
		}
	}

	// Translate the 1.5.0 OnRamp token transfer fee configs to FeeQuoterTokenTransferFeeConfig
	translateTokenTransferFeeCfgReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		migrate_seq.SeqTranslateOnRampToFeeQTokenTransferFeeCfg,
		e.BlockChains.EVMChains(),
		cfg.toSequenceInput(state),
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
		return cldf.ChangesetOutput{}, fmt.Errorf("Failed to execute FeeQUpdateTransferTokenFeeCfgSeq: %w", err)
	}

	// TODO: check the output: its empty
	csOutput, err = opsutil.AddEVMCallSequenceToCSOutput(e, state, csOutput, report, err, cfg.MCMS, "Call ApplyTokenTransferFeeConfigUpdates on FeeQuoter")
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
