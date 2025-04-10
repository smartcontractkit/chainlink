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
	"github.com/smartcontractkit/chainlink/deployment/operations"
	"github.com/smartcontractkit/mcms"
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
		if chainState.MCMSAddress == aptos.AccountZero {
			mcmsConfig := config.MCMSConfigPerChain[chainSel]
			if err := mcmsConfig.Validate(); err != nil {
				errs = append(errs, fmt.Errorf("invalid mcms configs for Aptos chain %d: %w", chainSel, err))
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
	proposals := []mcms.Proposal{}
	seqReports := make([]operations.Report[any, any], 0)

	// Deploy CCIP on each Aptos chain in config
	for chainSel := range config.ContractParamsPerChain {
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
		proposals = append(proposals, *mcmsSeqReport.Output.MCMSProposal)

		// CCIP Deploy operations
		ccipSeqInput := seq.DeployCCIPSeqInput{
			MCMSAddress: mcmsSeqReport.Output.MCMSAddress,
			MCMSOpCount: mcmsSeqReport.Output.NextOpCount,
			CCIPConfig:  config.ContractParamsPerChain[chainSel],
		}
		ccipSeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployCCIPSequence, deps, ccipSeqInput)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Aptos chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, ccipSeqReport.ExecutionReports...)
		for _, proposal := range ccipSeqReport.Output.MCMSProposals {
			proposals = append(proposals, *proposal)
		}
	}

	return deployment.ChangesetOutput{
		AddressBook:   ab,
		MCMSProposals: proposals,
		Reports:       seqReports,
	}, nil
}
