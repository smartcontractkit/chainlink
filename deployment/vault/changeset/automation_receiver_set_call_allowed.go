package changeset

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	proposeutils "github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"
	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	proposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/automation-cre/generated/latest/automation_receiver"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	vaulttypes "github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

var SetCallAllowedChangeSet cldf.ChangeSetV2[vaulttypes.SetCallAllowedInput] = setCallAllowed{}

type setCallAllowed struct{}

func (s setCallAllowed) VerifyPreconditions(env cldf.Environment, config vaulttypes.SetCallAllowedInput) error {
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
		if err := validateEthAddress("targetAddress", chainCfg.TargetAddress); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if _, err := parseSelectorHex(chainCfg.Selector); err != nil {
			return fmt.Errorf("chain %d: invalid selector %q: %w", chainSelector, chainCfg.Selector, err)
		}
	}
	return nil
}

func (s setCallAllowed) Apply(e cldf.Environment, config vaulttypes.SetCallAllowedInput) (cldf.ChangesetOutput, error) {
	logger := e.Logger
	logger.Infow("Generating setCallAllowed proposal", "numChains", len(config.Chains))

	evmChains := e.BlockChains.EVMChains()

	var primaryChain cldf_evm.Chain
	for chainSelector := range config.Chains {
		primaryChain = evmChains[chainSelector]
		break
	}

	deps := VaultDeps{
		Auth:        primaryChain.DeployerKey,
		Chain:       primaryChain,
		Environment: e,
		DataStore:   e.DataStore,
	}

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		SetCallAllowedSequence,
		deps,
		SetCallAllowedSequenceInput{
			Chains:     config.Chains,
			MCMSConfig: config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("setCallAllowed sequence failed: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: seqReport.Output.MCMSTimelockProposals,
	}, nil
}

// ================================================
// ================================================
// SetCallAllowed SEQUENCE
// ================================================
// ================================================

type SetCallAllowedSequenceInput struct {
	Chains     map[uint64]vaulttypes.SetCallAllowedChainConfig `json:"chains"`
	MCMSConfig *proposalutils.TimelockConfig                   `json:"mcms_config,omitempty"`
}

type SetCallAllowedSequenceOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal `json:"mcms_timelock_proposals"`
}

var SetCallAllowedSequence = operations.NewSequence(
	"set-call-allowed-sequence",
	semver.MustParse("1.0.0"),
	"Build a timelock proposal to call setCallAllowed on AutomationReceiver for each chain",
	func(b operations.Bundle, deps VaultDeps, input SetCallAllowedSequenceInput) (SetCallAllowedSequenceOutput, error) {
		b.Logger.Infow("Starting setCallAllowed sequence", "chains", len(input.Chains))

		var batches []mcmstypes.BatchOperation
		timelockAddresses := make(map[uint64]string)
		mcmAddressByChain := make(map[uint64]string)

		for chainSelector, chainCfg := range input.Chains {
			selector, err := parseSelectorHex(chainCfg.Selector)
			if err != nil {
				return SetCallAllowedSequenceOutput{}, fmt.Errorf("chain %d: invalid selector: %w", chainSelector, err)
			}

			opReport, err := operations.ExecuteOperation(b, SetCallAllowedProposalOperation, deps, SetCallAllowedProposalOpInput{
				ChainSelector: chainSelector,
				TargetAddress: chainCfg.TargetAddress,
				Selector:      selector,
				Allowed:       chainCfg.Allowed,
				MCMSConfig:    input.MCMSConfig,
			})
			if err != nil {
				return SetCallAllowedSequenceOutput{}, fmt.Errorf("chain %d: failed to generate setCallAllowed batch: %w", chainSelector, err)
			}
			opOutput := opReport.Output

			batches = append(batches, opOutput.BatchOperation)
			timelockAddresses[chainSelector] = opOutput.TimelockAddress
			mcmAddressByChain[chainSelector] = opOutput.MCMSAddress
		}

		proposal, err := proposeutils.BuildProposalFromBatchesV2(deps.Environment, timelockAddresses, mcmAddressByChain, nil, batches, "AutomationReceiver SetCallAllowed", ethBalMonProposalTimelockConfig(input.MCMSConfig))
		if err != nil {
			return SetCallAllowedSequenceOutput{}, fmt.Errorf("failed to build timelock proposal: %w", err)
		}

		b.Logger.Infow("Generated setCallAllowed proposal",
			"chains", len(input.Chains), "operations", len(batches))

		return SetCallAllowedSequenceOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	},
)

// ================================================
// ================================================
// SetCallAllowed proposal OPERATION
// ================================================
// ================================================

type SetCallAllowedProposalOpInput struct {
	ChainSelector uint64                        `json:"chain_selector"`
	TargetAddress string                        `json:"target_address"`
	Selector      [4]byte                       `json:"selector"`
	Allowed       bool                          `json:"allowed"`
	MCMSConfig    *proposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

type SetCallAllowedProposalOpOutput struct {
	ChainSelector   uint64                   `json:"chain_selector"`
	BatchOperation  mcmstypes.BatchOperation `json:"batch_operation"`
	TimelockAddress string                   `json:"timelock_address"`
	MCMSAddress     string                   `json:"mcms_address"`
}

var SetCallAllowedProposalOperation = operations.NewOperation(
	"set-call-allowed-proposal-operation",
	semver.MustParse("1.0.0"),
	"Operation to create transaction batch for AutomationReceiver setCallAllowed",
	func(b operations.Bundle, deps VaultDeps, input SetCallAllowedProposalOpInput) (SetCallAllowedProposalOpOutput, error) {
		b.Logger.Infow("Starting setCallAllowed proposal operation",
			"chainSelector", input.ChainSelector,
			"target", input.TargetAddress,
			"selector", fmt.Sprintf("0x%08x", input.Selector),
			"allowed", input.Allowed,
		)

		chain, ok := deps.Environment.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return SetCallAllowedProposalOpOutput{}, fmt.Errorf("chain not found in environment: %d", input.ChainSelector)
		}

		automationReceiverAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			cldf.ContractType(vaulttypes.AutomationReceiverContractType),
		)
		if err != nil {
			return SetCallAllowedProposalOpOutput{}, fmt.Errorf("failed to get AutomationReceiver address: %w", err)
		}

		timelockAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			commontypes.RBACTimelock,
		)
		if err != nil {
			return SetCallAllowedProposalOpOutput{}, fmt.Errorf("failed to get timelock address: %w", err)
		}
		mcmsAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			ethBalMonMCMSContractTypeForProposal(input.MCMSConfig),
		)
		if err != nil {
			return SetCallAllowedProposalOpOutput{}, fmt.Errorf("failed to get MCMS address: %w", err)
		}

		receiver, err := automation_receiver.NewAutomationReceiver(
			common.HexToAddress(automationReceiverAddr),
			chain.Client,
		)
		if err != nil {
			return SetCallAllowedProposalOpOutput{}, fmt.Errorf("failed to instantiate AutomationReceiver at %s: %w", automationReceiverAddr, err)
		}

		setCallAllowedTx, err := receiver.SetCallAllowed(
			cldf.SimTransactOpts(),
			common.HexToAddress(input.TargetAddress),
			input.Selector,
			input.Allowed,
		)
		if err != nil {
			return SetCallAllowedProposalOpOutput{}, fmt.Errorf("failed to generate setCallAllowed calldata on chain %d: %w", input.ChainSelector, err)
		}

		batch := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(input.ChainSelector),
			Transactions: []mcmstypes.Transaction{
				{
					OperationMetadata: mcmstypes.OperationMetadata{
						ContractType: vaulttypes.AutomationReceiverContractType,
						Tags: []string{
							"setCallAllowed",
						},
					},
					To:               automationReceiverAddr,
					Data:             setCallAllowedTx.Data(),
					AdditionalFields: json.RawMessage(`{"value": 0}`),
				},
			},
		}

		b.Logger.Infow("Generated setCallAllowed batch",
			"chainSelector", input.ChainSelector,
			"automationReceiver", automationReceiverAddr,
			"target", input.TargetAddress,
		)

		return SetCallAllowedProposalOpOutput{
			ChainSelector:   input.ChainSelector,
			BatchOperation:  batch,
			TimelockAddress: timelockAddr,
			MCMSAddress:     mcmsAddr,
		}, nil
	},
)

// parseSelectorHex parses a hex string like "0x4b9f5c20" into a [4]byte selector.
func parseSelectorHex(s string) ([4]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	b, err := hex.DecodeString(s)
	if err != nil {
		return [4]byte{}, fmt.Errorf("not valid hex: %w", err)
	}
	if len(b) != 4 {
		return [4]byte{}, fmt.Errorf("selector must be exactly 4 bytes, got %d", len(b))
	}
	return [4]byte{b[0], b[1], b[2], b[3]}, nil
}
