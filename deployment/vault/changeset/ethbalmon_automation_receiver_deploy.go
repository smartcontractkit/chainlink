package changeset

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/mcms"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	proposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	vaulttypes "github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

var DeployEthBalMonWithReceiverChangeSet cldf.ChangeSetV2[vaulttypes.DeployEthBalMonWithReceiverInput] = deployEthBalMonWithReceiver{}

type deployEthBalMonWithReceiver struct{}

func (d deployEthBalMonWithReceiver) VerifyPreconditions(env cldf.Environment, config vaulttypes.DeployEthBalMonWithReceiverInput) error {
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
		if (chainCfg.ExpectedAuthor != "") != (chainCfg.ExpectedWorkflowName != "") {
			return fmt.Errorf("chain %d: expectedAuthor and expectedWorkflowName must be set together", chainSelector)
		}
		if chainCfg.ExpectedAuthor != "" {
			if err := validateEthAddress("expectedAuthor", chainCfg.ExpectedAuthor); err != nil {
				return fmt.Errorf("chain %d: %w", chainSelector, err)
			}
		}
		if err := validateDeployEthBalMonMCMSInDatastore(env, chainSelector, config.MCMSConfig); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
	}
	return nil
}

func (d deployEthBalMonWithReceiver) Apply(e cldf.Environment, config vaulttypes.DeployEthBalMonWithReceiverInput) (cldf.ChangesetOutput, error) {
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
		DeployEthBalMonWithReceiverSequence,
		deps,
		DeployEthBalMonWithReceiverSequenceInput{
			Chains:     config.Chains,
			MCMSConfig: config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("deploy EthBalMon + AutomationReceiver sequence failed: %w", err)
	}

	memoryDataStore := ds.NewMemoryDataStore()
	contractsByChain := make(map[uint64]string)

	for _, chainOut := range seqReport.Output.Chains {
		contractsByChain[chainOut.ChainSelector] = chainOut.EthBalMonAddress

		arRef := ds.AddressRef{
			ChainSelector: chainOut.ChainSelector,
			Address:       chainOut.AutomationReceiverAddress,
			Type:          ds.ContractType(vaulttypes.AutomationReceiverContractType),
			Version:       semver.MustParse("1.0.0"),
			Labels:        ds.NewLabelSet(),
		}
		ebmRef := ds.AddressRef{
			ChainSelector: chainOut.ChainSelector,
			Address:       chainOut.EthBalMonAddress,
			Type:          ds.ContractType(vaulttypes.EthBalMonContractType),
			Version:       semver.MustParse("1.0.0"),
			Labels:        ds.NewLabelSet(vaulttypes.EthBalMonContractType, "EthBalMonV1_0_0"),
		}
		contractMetadata := ds.ContractMetadata{
			ChainSelector: chainOut.ChainSelector,
			Address:       chainOut.EthBalMonAddress,
			Metadata: map[string]any{
				"deployTxHash":              chainOut.EthBalMonDeployTxHash,
				"deployBlockNumber":         chainOut.EthBalMonDeployBlockNumber,
				"keeperRegistryAddress":     chainOut.AutomationReceiverAddress,
				"minWaitPeriodSeconds":      chainOut.MinWaitPeriodSeconds,
				"automationReceiverAddress": chainOut.AutomationReceiverAddress,
				"timelockAddress":           chainOut.TimelockAddress,
				"mcmsAddress":               chainOut.MCMSAddress,
				"transferOwnershipTxHash":   chainOut.EthBalMonTransferOwnershipTxHash,
			},
		}

		if err := memoryDataStore.Addresses().Add(arRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: failed to add AR address ref: %w", chainOut.ChainSelector, err)
		}
		if err := memoryDataStore.Addresses().Add(ebmRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: failed to add EthBalMon address ref: %w", chainOut.ChainSelector, err)
		}
		if err := memoryDataStore.ContractMetadata().Add(contractMetadata); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: failed to add contract metadata: %w", chainOut.ChainSelector, err)
		}
	}

	proposal, err := BuildAcceptOwnershipTimelockProposal(
		e,
		AcceptOwnershipProposalInput{
			ContractsByChain: contractsByChain,
			Description:      "Accept ownership of EthBalanceMonitor across chains",
			MCMSConfig:       deployEthBalMonAcceptOwnershipTimelockConfig(config.MCMSConfig),
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build accept ownership proposal: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore:             memoryDataStore,
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
	}, nil
}

// ================================================
// ================================================
// Deploy EthBalMon + AutomationReceiver SEQUENCE
// ================================================
// ================================================

type DeployEthBalMonWithReceiverSequenceInput struct {
	Chains     map[uint64]vaulttypes.DeployEthBalMonWithReceiverChainConfig
	MCMSConfig *proposalutils.TimelockConfig
}

type DeployEthBalMonWithReceiverPerChainOutput struct {
	ChainSelector                    uint64
	AutomationReceiverAddress        string
	EthBalMonAddress                 string
	EthBalMonDeployTxHash            string
	EthBalMonDeployBlockNumber       uint64
	MinWaitPeriodSeconds             uint64
	TimelockAddress                  string
	MCMSAddress                      string
	EthBalMonTransferOwnershipTxHash string
}

type DeployEthBalMonWithReceiverSequenceOutput struct {
	Chains []DeployEthBalMonWithReceiverPerChainOutput
}

var DeployEthBalMonWithReceiverSequence = operations.NewSequence(
	"deploy-ethbalmon-with-receiver-sequence",
	semver.MustParse("1.0.0"),
	"Deploy AutomationReceiver + EthBalMon, wire them, and transfer ownerships",
	func(b operations.Bundle, deps VaultDeps, input DeployEthBalMonWithReceiverSequenceInput) (DeployEthBalMonWithReceiverSequenceOutput, error) {
		out := DeployEthBalMonWithReceiverSequenceOutput{
			Chains: []DeployEthBalMonWithReceiverPerChainOutput{},
		}

		performUpkeepSig := crypto.Keccak256([]byte("performUpkeep(bytes)"))
		performUpkeepSelector := [4]byte{performUpkeepSig[0], performUpkeepSig[1], performUpkeepSig[2], performUpkeepSig[3]}

		for chainSelector, chainCfg := range input.Chains {
			var rawMinWait uint64
			if chainCfg.SetMinWaitPeriodSeconds != nil {
				rawMinWait = *chainCfg.SetMinWaitPeriodSeconds
			}
			minWait := effectiveMinWaitPeriodSeconds(rawMinWait)

			timelockAddr, err := getRequiredContractAddress(deps.DataStore, chainSelector, commontypes.RBACTimelock)
			if err != nil {
				return DeployEthBalMonWithReceiverSequenceOutput{}, fmt.Errorf("chain %d: failed to get timelock address: %w", chainSelector, err)
			}
			mcmsAddr, err := getRequiredContractAddress(
				deps.DataStore,
				chainSelector,
				ethBalMonMCMSContractTypeForAction(deployEthBalMonAcceptOwnershipMCMSAction(input.MCMSConfig)),
			)
			if err != nil {
				return DeployEthBalMonWithReceiverSequenceOutput{}, fmt.Errorf("chain %d: failed to get mcms address: %w", chainSelector, err)
			}

			// 1. Deploy AutomationReceiver
			arReport, err := operations.ExecuteOperation(b, DeployAutomationReceiverOperation, deps, DeployAutomationReceiverInput{
				ChainSelector:    chainSelector,
				ForwarderAddress: chainCfg.ForwarderAddress,
			})
			if err != nil {
				return DeployEthBalMonWithReceiverSequenceOutput{}, fmt.Errorf("chain %d: deploy AutomationReceiver failed: %w", chainSelector, err)
			}
			arAddress := arReport.Output.AutomationReceiverAddress

			// 2. Deploy EthBalMon with AR as keeper registry
			ebmReport, err := operations.ExecuteOperation(b, DeployEthBalMonContractOperation, deps, DeployEthBalMonContractInput{
				ChainSelector:         chainSelector,
				KeeperRegistryAddress: arAddress,
				MinWaitPeriodSeconds:  minWait,
			})
			if err != nil {
				return DeployEthBalMonWithReceiverSequenceOutput{}, fmt.Errorf("chain %d: deploy EthBalMon failed: %w", chainSelector, err)
			}
			ebmAddress := ebmReport.Output.ContractAddress

			// 3. setCallAllowed on AR → EthBalMon.performUpkeep
			_, err = operations.ExecuteOperation(b, SetCallAllowedOperation, deps, SetCallAllowedOperationInput{
				ChainSelector:             chainSelector,
				AutomationReceiverAddress: arAddress,
				TargetAddress:             ebmAddress,
				Selector:                  performUpkeepSelector,
				Allowed:                   true,
			})
			if err != nil {
				return DeployEthBalMonWithReceiverSequenceOutput{}, fmt.Errorf("chain %d: setCallAllowed failed: %w", chainSelector, err)
			}

			// 3b. Optionally lock the AR's inbound identity guard (deployer still owns the AR).
			if chainCfg.ExpectedAuthor != "" && chainCfg.ExpectedWorkflowName != "" {
				_, err = operations.ExecuteOperation(b, SetExpectedWorkflowIdentityOperation, deps, SetExpectedWorkflowIdentityOperationInput{
					ChainSelector:             chainSelector,
					AutomationReceiverAddress: arAddress,
					ExpectedAuthor:            chainCfg.ExpectedAuthor,
					ExpectedWorkflowName:      chainCfg.ExpectedWorkflowName,
				})
				if err != nil {
					return DeployEthBalMonWithReceiverSequenceOutput{}, fmt.Errorf("chain %d: setExpectedWorkflowIdentity failed: %w", chainSelector, err)
				}
			}

			// 4. Transfer AR ownership to Timelock (OZ Ownable — single-step, immediate)
			_, err = operations.ExecuteOperation(b, TransferAutomationReceiverOwnershipOperation, deps, TransferAutomationReceiverOwnershipInput{
				ChainSelector:             chainSelector,
				AutomationReceiverAddress: arAddress,
				TimelockAddress:           timelockAddr,
			})
			if err != nil {
				return DeployEthBalMonWithReceiverSequenceOutput{}, fmt.Errorf("chain %d: transfer AR ownership failed: %w", chainSelector, err)
			}

			// 5. Transfer EthBalMon ownership to Timelock (ConfirmedOwner — two-step, needs MCMS accept)
			transferReport, err := operations.ExecuteOperation(b, TransferOwnershipOperation, deps, TransferEthBalMonOwnershipInput{
				ChainSelector:   chainSelector,
				ContractAddress: ebmAddress,
				TimelockAddress: timelockAddr,
			})
			if err != nil {
				return DeployEthBalMonWithReceiverSequenceOutput{}, fmt.Errorf("chain %d: transfer EthBalMon ownership failed: %w", chainSelector, err)
			}

			out.Chains = append(out.Chains, DeployEthBalMonWithReceiverPerChainOutput{
				ChainSelector:                    chainSelector,
				AutomationReceiverAddress:        arAddress,
				EthBalMonAddress:                 ebmAddress,
				EthBalMonDeployTxHash:            ebmReport.Output.TxHash,
				EthBalMonDeployBlockNumber:       ebmReport.Output.BlockNumber,
				MinWaitPeriodSeconds:             minWait,
				TimelockAddress:                  timelockAddr,
				MCMSAddress:                      mcmsAddr,
				EthBalMonTransferOwnershipTxHash: transferReport.Output.TxHash,
			})
		}

		return out, nil
	},
)
