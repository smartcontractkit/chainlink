package changeset

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/crypto"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	vaulttypes "github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

var DeployAutomationReceiverChangeSet cldf.ChangeSetV2[vaulttypes.DeployAutomationReceiverInput] = deployAutomationReceiver{}

type deployAutomationReceiver struct{}

func (d deployAutomationReceiver) VerifyPreconditions(env cldf.Environment, config vaulttypes.DeployAutomationReceiverInput) error {
	if len(config.Chains) == 0 {
		return errors.New("chains must not be empty")
	}
	evmChains := env.BlockChains.EVMChains()
	for chainSelector, chainCfg := range config.Chains {
		if err := validateChainSelector(chainSelector, env); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if _, ok := evmChains[chainSelector]; !ok {
			return fmt.Errorf("chain %d: not found in environment", chainSelector)
		}
		if err := validateEthAddress("forwarderAddress", chainCfg.ForwarderAddress); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if err := validateEthAddress("targetAddress", chainCfg.TargetAddress); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if chainCfg.Selector != "" {
			if _, err := parseSelectorHex(chainCfg.Selector); err != nil {
				return fmt.Errorf("chain %d: invalid selector %q: %w", chainSelector, chainCfg.Selector, err)
			}
		}
		if (chainCfg.ExpectedAuthor != "") != (chainCfg.ExpectedWorkflowName != "") {
			return fmt.Errorf("chain %d: expectedAuthor and expectedWorkflowName must be set together", chainSelector)
		}
		if chainCfg.ExpectedAuthor != "" {
			if err := validateEthAddress("expectedAuthor", chainCfg.ExpectedAuthor); err != nil {
				return fmt.Errorf("chain %d: %w", chainSelector, err)
			}
		}
	}
	return nil
}

func (d deployAutomationReceiver) Apply(e cldf.Environment, config vaulttypes.DeployAutomationReceiverInput) (cldf.ChangesetOutput, error) {
	evmChains := e.BlockChains.EVMChains()

	var primaryChainSelector uint64
	for sel := range config.Chains {
		if primaryChainSelector == 0 || sel < primaryChainSelector {
			primaryChainSelector = sel
		}
	}

	primaryChain := evmChains[primaryChainSelector]
	deps := VaultDeps{
		Auth:        primaryChain.DeployerKey,
		Chain:       primaryChain,
		Environment: e,
		DataStore:   e.DataStore,
	}

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		DeployAutomationReceiverSequence,
		deps,
		DeployAutomationReceiverSequenceInput{
			Chains: config.Chains,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("deploy AutomationReceiver sequence failed: %w", err)
	}

	memoryDataStore := ds.NewMemoryDataStore()
	for _, chainOut := range seqReport.Output.Chains {
		ref := ds.AddressRef{
			ChainSelector: chainOut.ChainSelector,
			Address:       chainOut.AutomationReceiverAddress,
			Type:          ds.ContractType(vaulttypes.AutomationReceiverContractType),
			Version:       semver.MustParse("1.0.0"),
			Labels:        ds.NewLabelSet(),
		}
		if err := memoryDataStore.Addresses().Add(ref); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: failed to add AutomationReceiver address ref: %w", chainOut.ChainSelector, err)
		}
	}

	return cldf.ChangesetOutput{DataStore: memoryDataStore}, nil
}

// ================================================
// ================================================
// Deploy AutomationReceiver SEQUENCE
// ================================================
// ================================================

type DeployAutomationReceiverSequenceInput struct {
	Chains map[uint64]vaulttypes.AutomationReceiverChainConfig
}

type DeployAutomationReceiverPerChainOutput struct {
	ChainSelector             uint64
	AutomationReceiverAddress string
}

type DeployAutomationReceiverSequenceOutput struct {
	Chains []DeployAutomationReceiverPerChainOutput
}

var DeployAutomationReceiverSequence = operations.NewSequence(
	"deploy-automation-receiver-sequence",
	semver.MustParse("1.0.0"),
	"Deploy AutomationReceiver, call setCallAllowed, and transfer ownership to Timelock",
	func(b operations.Bundle, deps VaultDeps, input DeployAutomationReceiverSequenceInput) (DeployAutomationReceiverSequenceOutput, error) {
		out := DeployAutomationReceiverSequenceOutput{
			Chains: []DeployAutomationReceiverPerChainOutput{},
		}

		for chainSelector, chainCfg := range input.Chains {
			timelockAddr, err := getRequiredContractAddress(deps.DataStore, chainSelector, commontypes.RBACTimelock)
			if err != nil {
				return DeployAutomationReceiverSequenceOutput{}, fmt.Errorf("chain %d: failed to get timelock address: %w", chainSelector, err)
			}

			deployReport, err := operations.ExecuteOperation(b, DeployAutomationReceiverOperation, deps, DeployAutomationReceiverInput{
				ChainSelector:    chainSelector,
				ForwarderAddress: chainCfg.ForwarderAddress,
			})
			if err != nil {
				return DeployAutomationReceiverSequenceOutput{}, fmt.Errorf("chain %d: deploy AutomationReceiver failed: %w", chainSelector, err)
			}
			arAddress := deployReport.Output.AutomationReceiverAddress

			selector, err := resolveSelector(chainCfg.Selector)
			if err != nil {
				return DeployAutomationReceiverSequenceOutput{}, fmt.Errorf("chain %d: %w", chainSelector, err)
			}

			_, err = operations.ExecuteOperation(b, SetCallAllowedOperation, deps, SetCallAllowedOperationInput{
				ChainSelector:             chainSelector,
				AutomationReceiverAddress: arAddress,
				TargetAddress:             chainCfg.TargetAddress,
				Selector:                  selector,
				Allowed:                   true,
			})
			if err != nil {
				return DeployAutomationReceiverSequenceOutput{}, fmt.Errorf("chain %d: setCallAllowed failed: %w", chainSelector, err)
			}

			// Optionally lock the receiver's inbound identity guard while the deployer still owns
			// the contract (before ownership is transferred to the Timelock below).
			if chainCfg.ExpectedAuthor != "" && chainCfg.ExpectedWorkflowName != "" {
				_, err = operations.ExecuteOperation(b, SetExpectedWorkflowIdentityOperation, deps, SetExpectedWorkflowIdentityOperationInput{
					ChainSelector:             chainSelector,
					AutomationReceiverAddress: arAddress,
					ExpectedAuthor:            chainCfg.ExpectedAuthor,
					ExpectedWorkflowName:      chainCfg.ExpectedWorkflowName,
				})
				if err != nil {
					return DeployAutomationReceiverSequenceOutput{}, fmt.Errorf("chain %d: setExpectedWorkflowIdentity failed: %w", chainSelector, err)
				}
			}

			_, err = operations.ExecuteOperation(b, TransferAutomationReceiverOwnershipOperation, deps, TransferAutomationReceiverOwnershipInput{
				ChainSelector:             chainSelector,
				AutomationReceiverAddress: arAddress,
				TimelockAddress:           timelockAddr,
			})
			if err != nil {
				return DeployAutomationReceiverSequenceOutput{}, fmt.Errorf("chain %d: transfer AutomationReceiver ownership failed: %w", chainSelector, err)
			}

			out.Chains = append(out.Chains, DeployAutomationReceiverPerChainOutput{
				ChainSelector:             chainSelector,
				AutomationReceiverAddress: arAddress,
			})
		}

		return out, nil
	},
)

// resolveSelector returns the performUpkeep(bytes) selector when s is empty,
// otherwise parses the provided hex string.
func resolveSelector(s string) ([4]byte, error) {
	if s == "" {
		sig := crypto.Keccak256([]byte("performUpkeep(bytes)"))
		return [4]byte{sig[0], sig[1], sig[2], sig[3]}, nil
	}
	return parseSelectorHex(s)
}
