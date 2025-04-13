package aptos

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/operations"
	"github.com/smartcontractkit/mcms"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

var _ deployment.ChangeSetV2[config.DeployAptosChainConfig] = DeployAptosChain{}

// DeployAptosChain deploys Aptos chain packages and modules
type DeployAptosChain struct{}

func (cs DeployAptosChain) VerifyPreconditions(env deployment.Environment, config config.DeployAptosChainConfig) error {
	// Validate env and prerequisite contracts
	state, err := changeset.LoadOnchainStateAptos(env)
	if err != nil {
		return fmt.Errorf("failed to load existing Aptos onchain state: %w", err)
	}
	var errs []error
	for chainSel := range config.ContractParamsPerChain {
		if err := config.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("invalid config for Aptos chain %d: %w", chainSel, err))
			continue
		}
		if _, ok := env.AptosChains[chainSel]; !ok {
			errs = append(errs, fmt.Errorf("aptos chain %d not found in env", chainSel))
		}
		chainState, ok := state[chainSel]
		if !ok {
			errs = append(errs, fmt.Errorf("aptos chain %d not found in state", chainSel))
			continue
		}
		if chainState.MCMSAddress == (aptos.AccountAddress{}) {
			mcmsConfig := config.MCMSConfigPerChain[chainSel]
			for _, cfg := range []mcmstypes.Config{mcmsConfig.Bypasser, mcmsConfig.Canceller, mcmsConfig.Proposer} {
				if err := cfg.Validate(); err != nil {
					errs = append(errs, fmt.Errorf("invalid mcms configs for Aptos chain %d: %w", chainSel, err))
				}
			}
		}
	}

	return errors.Join(errs...)
}

func (cs DeployAptosChain) Apply(env deployment.Environment, config config.DeployAptosChainConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainStateAptos(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	ab := deployment.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)
	proposals := make([]mcms.TimelockProposal, 0)

	// Deploy CCIP on each Aptos chain in config
	for chainSel := range config.ContractParamsPerChain {
		mcmsOperations := []mcmstypes.BatchOperation{}
		chainState := state[chainSel]
		aptosChain := env.AptosChains[chainSel]

		deps := operation.AptosDeps{
			AB:           ab,
			AptosChain:   aptosChain,
			OnChainState: chainState,
		}

		// MCMS Deploy operations
		mcmsSeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployMCMSSequence, deps, config.MCMSConfigPerChain[chainSel])
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		seqReports = append(seqReports, mcmsSeqReport.ExecutionReports...)
		mcmsOperations = append(mcmsOperations, mcmsSeqReport.Output.MCMSOperation)

		// CCIP Deploy operations
		ccipSeqInput := seq.DeployCCIPSeqInput{
			MCMSAddress: mcmsSeqReport.Output.MCMSAddress,
			CCIPConfig:  config.ContractParamsPerChain[chainSel],
		}
		ccipSeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployCCIPSequence, deps, ccipSeqInput)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Aptos chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, ccipSeqReport.ExecutionReports...)
		mcmsOperations = append(mcmsOperations, ccipSeqReport.Output.MCMSOperations...)

		// Generate MCMS proposals
		proposal, err := utils.GenerateProposal(
			aptosChain.Client,
			mcmsSeqReport.Output.MCMSAddress,
			chainSel,
			mcmsOperations,
			"Deploy Aptos MCMS and CCIP",
			aptosmcms.TimelockRoleProposer,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", chainSel, err)
		}
		proposals = append(proposals, *proposal)
	}
	return deployment.ChangesetOutput{
		AddressBook:           ab,
		MCMSTimelockProposals: proposals,
		Reports:               seqReports,
	}, nil
}
