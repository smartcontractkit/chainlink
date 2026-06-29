package changeset

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

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
		return fmt.Errorf("chains must not be empty")
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

	memoryDataStore := ds.NewMemoryDataStore()

	for chainSelector, chainCfg := range config.Chains {
		timelockAddr, err := getRequiredContractAddress(e.DataStore, chainSelector, commontypes.RBACTimelock)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: failed to get timelock address: %w", chainSelector, err)
		}

		deployReport, err := operations.ExecuteOperation(
			e.OperationsBundle,
			DeployAutomationReceiverOperation,
			deps,
			DeployAutomationReceiverInput{
				ChainSelector:    chainSelector,
				ForwarderAddress: chainCfg.ForwarderAddress,
			},
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: deploy AutomationReceiver failed: %w", chainSelector, err)
		}
		arAddress := deployReport.Output.AutomationReceiverAddress

		_, err = operations.ExecuteOperation(
			e.OperationsBundle,
			TransferAutomationReceiverOwnershipOperation,
			deps,
			TransferAutomationReceiverOwnershipInput{
				ChainSelector:             chainSelector,
				AutomationReceiverAddress: arAddress,
				TimelockAddress:           timelockAddr,
			},
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: transfer AutomationReceiver ownership failed: %w", chainSelector, err)
		}

		ref := ds.AddressRef{
			ChainSelector: chainSelector,
			Address:       arAddress,
			Type:          ds.ContractType(vaulttypes.AutomationReceiverContractType),
			Version:       semver.MustParse("1.0.0"),
			Qualifier:     "",
			Labels:        ds.NewLabelSet(),
		}
		if err := memoryDataStore.Addresses().Add(ref); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: failed to add AutomationReceiver address ref: %w", chainSelector, err)
		}
	}

	return cldf.ChangesetOutput{DataStore: memoryDataStore}, nil
}
