package changeset

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	proposeutils "github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"
	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	proposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/automation-cre/generated/latest/automation_receiver"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	vaulttypes "github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

// SetExpectedWorkflowIdentityChangeSet builds a timelock proposal to configure the
// AutomationReceiver's inbound identity guard (setExpectedAuthor + setExpectedWorkflowName) on
// already-deployed receivers whose ownership is the Timelock. The AutomationReceiver reverts
// inbound reports with WorkflowIdentityNotConfigured until the identity is configured, so this is
// required for a deployed receiver to accept reports. Author + name are stable across workflow
// redeploys, so this only needs to be set once (unlike pinning the workflow id).
var SetExpectedWorkflowIdentityChangeSet cldf.ChangeSetV2[vaulttypes.SetExpectedWorkflowIdentityInput] = setExpectedWorkflowIdentity{}

type setExpectedWorkflowIdentity struct{}

func (s setExpectedWorkflowIdentity) VerifyPreconditions(env cldf.Environment, config vaulttypes.SetExpectedWorkflowIdentityInput) error {
	if len(config.Chains) == 0 {
		return errors.New("chains must not be empty")
	}
	if config.MCMSConfig != nil && config.MCMSConfig.MinDelay < 0 {
		return fmt.Errorf("MCMS minimum delay cannot be negative: %d", config.MCMSConfig.MinDelay)
	}
	for chainSelector, chainCfg := range config.Chains {
		if err := validateChainSelector(chainSelector, env); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if chainCfg.ExpectedAuthor == "" || chainCfg.ExpectedWorkflowName == "" {
			return fmt.Errorf(
				"chain %d: both expectedAuthor and expectedWorkflowName are required "+
					"(the AutomationReceiver reverts with WorkflowIdentityNotConfigured otherwise)",
				chainSelector,
			)
		}
		if err := validateEthAddress("expectedAuthor", chainCfg.ExpectedAuthor); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
	}
	return nil
}

func (s setExpectedWorkflowIdentity) Apply(e cldf.Environment, config vaulttypes.SetExpectedWorkflowIdentityInput) (cldf.ChangesetOutput, error) {
	logger := e.Logger
	logger.Infow("Generating setExpectedWorkflowIdentity proposal", "numChains", len(config.Chains))

	evmChains := e.BlockChains.EVMChains()

	// deps only carries the deployer key/chain used to build the proposal (the batch itself is
	// generated per-chain), so any chain works — pick deterministically to keep output reproducible.
	primaryChainSelector, _ := lowestChainSelector(config.Chains)
	primaryChain := evmChains[primaryChainSelector]

	deps := VaultDeps{
		Auth:        primaryChain.DeployerKey,
		Chain:       primaryChain,
		Environment: e,
		DataStore:   e.DataStore,
	}

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		SetExpectedWorkflowIdentitySequence,
		deps,
		SetExpectedWorkflowIdentitySequenceInput{
			Chains:     config.Chains,
			MCMSConfig: config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("setExpectedWorkflowIdentity sequence failed: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: seqReport.Output.MCMSTimelockProposals,
	}, nil
}

// ================================================
// ================================================
// SetExpectedWorkflowIdentity SEQUENCE
// ================================================
// ================================================

type SetExpectedWorkflowIdentitySequenceInput struct {
	Chains     map[uint64]vaulttypes.SetExpectedWorkflowIdentityChainConfig `json:"chains"`
	MCMSConfig *proposalutils.TimelockConfig                                `json:"mcms_config,omitempty"`
}

type SetExpectedWorkflowIdentitySequenceOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal `json:"mcms_timelock_proposals"`
}

var SetExpectedWorkflowIdentitySequence = operations.NewSequence(
	"set-expected-workflow-identity-sequence",
	semver.MustParse("1.0.0"),
	"Build a timelock proposal to set the expected workflow identity on AutomationReceiver for each chain",
	func(b operations.Bundle, deps VaultDeps, input SetExpectedWorkflowIdentitySequenceInput) (SetExpectedWorkflowIdentitySequenceOutput, error) {
		b.Logger.Infow("Starting setExpectedWorkflowIdentity sequence", "chains", len(input.Chains))

		var batches []mcmstypes.BatchOperation
		timelockAddresses := make(map[uint64]string)
		mcmAddressByChain := make(map[uint64]string)

		for chainSelector, chainCfg := range input.Chains {
			opReport, err := operations.ExecuteOperation(b, SetExpectedWorkflowIdentityProposalOperation, deps, SetExpectedWorkflowIdentityProposalOpInput{
				ChainSelector:        chainSelector,
				ExpectedAuthor:       chainCfg.ExpectedAuthor,
				ExpectedWorkflowName: chainCfg.ExpectedWorkflowName,
				MCMSConfig:           input.MCMSConfig,
			})
			if err != nil {
				return SetExpectedWorkflowIdentitySequenceOutput{}, fmt.Errorf("chain %d: failed to generate setExpectedWorkflowIdentity batch: %w", chainSelector, err)
			}
			opOutput := opReport.Output

			batches = append(batches, opOutput.BatchOperation)
			timelockAddresses[chainSelector] = opOutput.TimelockAddress
			mcmAddressByChain[chainSelector] = opOutput.MCMSAddress
		}

		proposal, err := proposeutils.BuildProposalFromBatchesV2(deps.Environment, timelockAddresses, mcmAddressByChain, nil, batches, "AutomationReceiver SetExpectedWorkflowIdentity", ethBalMonProposalTimelockConfig(input.MCMSConfig))
		if err != nil {
			return SetExpectedWorkflowIdentitySequenceOutput{}, fmt.Errorf("failed to build timelock proposal: %w", err)
		}

		b.Logger.Infow("Generated setExpectedWorkflowIdentity proposal",
			"chains", len(input.Chains), "operations", len(batches))

		return SetExpectedWorkflowIdentitySequenceOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	},
)

// ================================================
// ================================================
// SetExpectedWorkflowIdentity proposal OPERATION
// ================================================
// ================================================

type SetExpectedWorkflowIdentityProposalOpInput struct {
	ChainSelector        uint64                        `json:"chain_selector"`
	ExpectedAuthor       string                        `json:"expected_author"`
	ExpectedWorkflowName string                        `json:"expected_workflow_name"`
	MCMSConfig           *proposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

type SetExpectedWorkflowIdentityProposalOpOutput struct {
	ChainSelector   uint64                   `json:"chain_selector"`
	BatchOperation  mcmstypes.BatchOperation `json:"batch_operation"`
	TimelockAddress string                   `json:"timelock_address"`
	MCMSAddress     string                   `json:"mcms_address"`
}

var SetExpectedWorkflowIdentityProposalOperation = operations.NewOperation(
	"set-expected-workflow-identity-proposal-operation",
	semver.MustParse("1.0.0"),
	"Operation to create a transaction batch for AutomationReceiver setExpectedAuthor + setExpectedWorkflowName",
	func(b operations.Bundle, deps VaultDeps, input SetExpectedWorkflowIdentityProposalOpInput) (SetExpectedWorkflowIdentityProposalOpOutput, error) {
		b.Logger.Infow("Starting setExpectedWorkflowIdentity proposal operation",
			"chainSelector", input.ChainSelector,
			"expectedAuthor", input.ExpectedAuthor,
			"expectedWorkflowName", input.ExpectedWorkflowName,
		)

		chain, ok := deps.Environment.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return SetExpectedWorkflowIdentityProposalOpOutput{}, fmt.Errorf("chain not found in environment: %d", input.ChainSelector)
		}

		automationReceiverAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			cldf.ContractType(vaulttypes.AutomationReceiverContractType),
		)
		if err != nil {
			return SetExpectedWorkflowIdentityProposalOpOutput{}, fmt.Errorf("failed to get AutomationReceiver address: %w", err)
		}

		timelockAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			commontypes.RBACTimelock,
		)
		if err != nil {
			return SetExpectedWorkflowIdentityProposalOpOutput{}, fmt.Errorf("failed to get timelock address: %w", err)
		}
		mcmsAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			ethBalMonMCMSContractTypeForProposal(input.MCMSConfig),
		)
		if err != nil {
			return SetExpectedWorkflowIdentityProposalOpOutput{}, fmt.Errorf("failed to get MCMS address: %w", err)
		}

		receiver, err := automation_receiver.NewAutomationReceiver(
			common.HexToAddress(automationReceiverAddr),
			chain.Client,
		)
		if err != nil {
			return SetExpectedWorkflowIdentityProposalOpOutput{}, fmt.Errorf("failed to instantiate AutomationReceiver at %s: %w", automationReceiverAddr, err)
		}

		// setExpectedAuthor and setExpectedWorkflowName, executed atomically in one timelock batch.
		authorTx, err := receiver.SetExpectedAuthor(cldf.SimTransactOpts(), common.HexToAddress(input.ExpectedAuthor))
		if err != nil {
			return SetExpectedWorkflowIdentityProposalOpOutput{}, fmt.Errorf("failed to generate setExpectedAuthor calldata on chain %d: %w", input.ChainSelector, err)
		}
		nameTx, err := receiver.SetExpectedWorkflowName(cldf.SimTransactOpts(), input.ExpectedWorkflowName)
		if err != nil {
			return SetExpectedWorkflowIdentityProposalOpOutput{}, fmt.Errorf("failed to generate setExpectedWorkflowName calldata on chain %d: %w", input.ChainSelector, err)
		}

		batch := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(input.ChainSelector),
			Transactions: []mcmstypes.Transaction{
				{
					OperationMetadata: mcmstypes.OperationMetadata{
						ContractType: vaulttypes.AutomationReceiverContractType,
						Tags:         []string{"setExpectedAuthor"},
					},
					To:               automationReceiverAddr,
					Data:             authorTx.Data(),
					AdditionalFields: json.RawMessage(`{"value": 0}`),
				},
				{
					OperationMetadata: mcmstypes.OperationMetadata{
						ContractType: vaulttypes.AutomationReceiverContractType,
						Tags:         []string{"setExpectedWorkflowName"},
					},
					To:               automationReceiverAddr,
					Data:             nameTx.Data(),
					AdditionalFields: json.RawMessage(`{"value": 0}`),
				},
			},
		}

		b.Logger.Infow("Generated setExpectedWorkflowIdentity batch",
			"chainSelector", input.ChainSelector,
			"automationReceiver", automationReceiverAddr,
		)

		return SetExpectedWorkflowIdentityProposalOpOutput{
			ChainSelector:   input.ChainSelector,
			BatchOperation:  batch,
			TimelockAddress: timelockAddr,
			MCMSAddress:     mcmsAddr,
		}, nil
	},
)
