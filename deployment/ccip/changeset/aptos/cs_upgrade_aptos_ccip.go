package aptos

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSetV2[config.UpgradeAptosChainConfig] = UpgradeAptosChain{}

// UpgradeAptosChain upgrades Aptos chain packages and modules
type UpgradeAptosChain struct{}

func (cs UpgradeAptosChain) VerifyPreconditions(env cldf.Environment, cfg config.UpgradeAptosChainConfig) error {
	var errs []error
	if !cfg.UpgradeCCIP && !cfg.UpgradeOffRamp && !cfg.UpgradeOnRamp && !cfg.UpgradeRouter {
		errs = append(errs, errors.New("no upgrades selected"))
	}
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to load Aptos onchain state: %w", err))
	}
	// Validate supported chain
	supportedChains := state.SupportedChains()
	if _, ok := supportedChains[cfg.ChainSelector]; !ok {
		errs = append(errs, fmt.Errorf("unsupported chain: %d", cfg.ChainSelector))
	}
	// Validate MCMS config
	if cfg.MCMS == nil {
		errs = append(errs, errors.New("MCMS config is required for UpgradeAptosChain changeset"))
	}
	// Check CCIP Package
	ccipAddress := state.AptosChains[cfg.ChainSelector].CCIPAddress
	client := env.BlockChains.AptosChains()[cfg.ChainSelector].Client
	if ccipAddress == (aptos.AccountAddress{}) {
		errs = append(errs, fmt.Errorf("package CCIP is not deployed on Aptos chain %d", cfg.ChainSelector))
	}
	// Check OnRamp module
	hasOnramp, err := utils.IsModuleDeployed(client, ccipAddress, "onramp")
	if err != nil || !hasOnramp {
		errs = append(errs, fmt.Errorf("onRamp module is not deployed on Aptos chain %d: %w", cfg.ChainSelector, err))
	}
	// Check OffRamp module
	hasOfframp, err := utils.IsModuleDeployed(client, ccipAddress, "offramp")
	if err != nil || !hasOfframp {
		errs = append(errs, fmt.Errorf("offRamp module is not deployed on Aptos chain %d: %w", cfg.ChainSelector, err))
	}
	// Check Router module
	hasRouter, err := utils.IsModuleDeployed(client, ccipAddress, "router")
	if err != nil || !hasRouter {
		errs = append(errs, fmt.Errorf("router module is not deployed on Aptos chain %d: %w", cfg.ChainSelector, err))
	}

	return errors.Join(errs...)
}

func (cs UpgradeAptosChain) Apply(env cldf.Environment, cfg config.UpgradeAptosChainConfig) (cldf.ChangesetOutput, error) {
	timeLockProposals := []mcms.TimelockProposal{}
	mcmsOperations := []types.BatchOperation{}
	seqReports := make([]operations.Report[any, any], 0)

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	deps := dependency.AptosDeps{
		AptosChain:       env.BlockChains.AptosChains()[cfg.ChainSelector],
		CCIPOnChainState: state,
	}

	// Execute the sequence
	upgradeSeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.UpgradeCCIPSequence, deps, cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	seqReports = append(seqReports, upgradeSeqReport.ExecutionReports...)
	mcmsOperations = append(mcmsOperations, upgradeSeqReport.Output...)

	// Generate MCMS proposals
	proposal, err := utils.GenerateProposal(
		env,
		state.AptosChains[cfg.ChainSelector].MCMSAddress,
		cfg.ChainSelector,
		mcmsOperations,
		"Upgrade chain contracts on Aptos chain",
		*cfg.MCMS,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", cfg.ChainSelector, err)
	}
	timeLockProposals = append(timeLockProposals, *proposal)

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: timeLockProposals,
		Reports:               seqReports,
	}, nil
}
