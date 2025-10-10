package writeevm

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	evmworkflow "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	chainlinkbig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

const flag = cre.WriteEVMCapability

func New() (*capabilities.Capability, error) {
	return capabilities.New(
		flag,
		capabilities.WithCapabilityRegistryV1ConfigFn(registerWithV1),
		capabilities.WithNodeConfigTransformerFn(transformNodeConfig),
	)
}

func registerWithV1(_ []string, nodeSetInput *cre.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	capabilities := make([]keystone_changeset.DONCapabilityWithConfig, 0)

	if nodeSetInput == nil {
		return nil, errors.New("node set input is nil")
	}

	// it's fine if there are no chain capabilities
	if nodeSetInput.ChainCapabilities == nil {
		return nil, nil
	}

	if _, ok := nodeSetInput.ChainCapabilities[flag]; !ok {
		return nil, nil
	}

	for _, chainID := range nodeSetInput.ChainCapabilities[flag].EnabledChains {
		fullName := corevm.GenerateWriteTargetName(chainID)
		splitName := strings.Split(fullName, "@")

		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   splitName[0],
				Version:        splitName[1],
				CapabilityType: 3, // TARGET
				ResponseType:   1, // OBSERVATION_IDENTICAL
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities, nil
}

func transformNodeConfig(input cre.GenerateConfigsInput, existingConfigs cre.NodeIndexToConfigOverride) (cre.NodeIndexToConfigOverride, error) {
	if input.NodeSet == nil {
		return nil, errors.New("node set input is nil")
	}

	if input.NodeSet.ChainCapabilities == nil || input.NodeSet.ChainCapabilities[flag] == nil {
		return existingConfigs, nil
	}

	if input.CapabilityConfigs == nil {
		return nil, errors.New("additional capabilities configs are nil, but are required to configure the write-evm capability")
	}

	workerNodes, wErr := input.DonMetadata.Workers()
	if wErr != nil {
		return nil, errors.Wrap(wErr, "failed to find worker nodes")
	}

	for _, workerNode := range workerNodes {
		// // get all the forwarders and add workflow config (FromAddress + Forwarder) for chains that have write-evm enabled
		data := []writeEVMData{}
		for _, chainID := range input.NodeSet.ChainCapabilities[flag].EnabledChains {
			chain, exists := chain_selectors.ChainByEvmChainID(chainID)
			if !exists {
				return nil, errors.Errorf("failed to find selector for chain ID %d", chainID)
			}

			evmData := writeEVMData{
				ChainID:       chainID,
				ChainSelector: chain.Selector,
			}

			forwarderAddress, fErr := findForwarderAddress(chain, input.AddressBook)
			if fErr != nil {
				return nil, errors.Errorf("failed to find forwarder address for chain %d", chain.Selector)
			}
			evmData.ForwarderAddress = forwarderAddress.Hex()

			evmKey, ok := workerNode.Keys.EVM[chainID]
			if !ok {
				return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
			}
			evmData.FromAddress = evmKey.PublicAddress

			var mergeErr error
			evmData, mergeErr = mergeDefaultAndRuntimeConfigValues(evmData, input.CapabilityConfigs, input.NodeSet.ChainCapabilities, chainID)
			if mergeErr != nil {
				return nil, errors.Wrap(mergeErr, "failed to merge default and runtime write-evm config values")
			}

			data = append(data, evmData)
		}

		if len(existingConfigs) < workerNode.Index+1 {
			return nil, errors.Errorf("missing config for node index %d", workerNode.Index)
		}

		currentConfig := existingConfigs[workerNode.Index]

		var typedConfig corechainlink.Config
		unmarshallErr := toml.Unmarshal([]byte(currentConfig), &typedConfig)
		if unmarshallErr != nil {
			return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal config for node index %d", workerNode.Index)
		}

		if len(typedConfig.EVM) < len(data) {
			return nil, fmt.Errorf("not enough EVM chains configured in node index %d to add write-evm config. Expected at least %d chains, but found %d", workerNode.Index, len(data), len(typedConfig.EVM))
		}

		for _, writeEVMInput := range data {
			chainFound := false
			for idx, evmChain := range typedConfig.EVM {
				chainIDIsEqual := evmChain.ChainID.Cmp(chainlinkbig.New(big.NewInt(libc.MustSafeInt64(writeEVMInput.ChainID)))) == 0
				if chainIDIsEqual {
					evmWorkflow, evmErr := buildEVMWorkflowConfig(writeEVMInput)
					if evmErr != nil {
						return nil, errors.Wrap(evmErr, "failed to build EVM workflow config")
					}

					typedConfig.EVM[idx].Workflow = *evmWorkflow
					typedConfig.EVM[idx].Transactions.ForwardersEnabled = ptr.Ptr(true)

					chainFound = true
					break
				}
			}

			if !chainFound {
				return nil, fmt.Errorf("failed to find EVM chain with ID %d in the config of node index %d to add write-evm config", writeEVMInput.ChainID, workerNode.Index)
			}
		}

		stringifiedConfig, mErr := toml.Marshal(typedConfig)
		if mErr != nil {
			return nil, errors.Wrapf(mErr, "failed to marshal config for node index %d", workerNode.Index)
		}

		existingConfigs[workerNode.Index] = string(stringifiedConfig)
	}

	return existingConfigs, nil
}

func findForwarderAddress(chain chain_selectors.Chain, addressBook deployment.AddressBook) (*common.Address, error) {
	addrsForChains, addErr := addressBook.AddressesForChain(chain.Selector)
	if addErr != nil {
		return nil, errors.Wrap(addErr, "failed to get addresses from address book")
	}

	for addr, addrValue := range addrsForChains {
		if addrValue.Type == keystone_changeset.KeystoneForwarder {
			return ptr.Ptr(common.HexToAddress(addr)), nil
		}
	}

	return nil, errors.Errorf("failed to find forwarder address for chain %d", chain.Selector)
}

func mergeDefaultAndRuntimeConfigValues(data writeEVMData, defaultCapabilityConfigs cre.CapabilityConfigs, nodeSetChainCapabilities map[string]*cre.ChainCapabilityConfig, chainID uint64) (writeEVMData, error) {
	if writeEvmConfig, ok := defaultCapabilityConfigs[flag]; ok {
		_, mergedConfig, rErr := envconfig.ResolveCapabilityForChain(
			cre.WriteEVMCapability,
			nodeSetChainCapabilities,
			writeEvmConfig.Config,
			chainID,
		)
		if rErr != nil {
			return data, errors.Wrapf(rErr, "failed to resolve write-evm config for chain %d", chainID)
		}

		runtimeValues := map[string]any{
			"FromAddress":      data.FromAddress.Hex(),
			"ForwarderAddress": data.ForwarderAddress,
		}

		var mErr error
		data.WorkflowConfig, mErr = don.ApplyRuntimeValues(mergedConfig, runtimeValues)
		if mErr != nil {
			return data, errors.Wrap(mErr, "failed to apply runtime values")
		}
	}

	return data, nil
}

func buildEVMWorkflowConfig(writeEVMInput writeEVMData) (*evmworkflow.Workflow, error) {
	var evmWorkflow evmworkflow.Workflow

	tmpl, tErr := template.New("evmWorkflowConfig").Parse(evmWorkflowConfigTemplate)
	if tErr != nil {
		return nil, errors.Wrap(tErr, "failed to parse evm workflow config template")
	}
	var configBuffer bytes.Buffer
	if executeErr := tmpl.Execute(&configBuffer, writeEVMInput.WorkflowConfig); executeErr != nil {
		return nil, errors.Wrap(executeErr, "failed to execute evm workflow config template")
	}

	configStr := configBuffer.String()
	if err := don.ValidateTemplateSubstitution(configStr, flag); err != nil {
		return nil, errors.Wrapf(err, "%s template validation failed", flag)
	}

	unmarshallErr := toml.Unmarshal([]byte(configStr), &evmWorkflow)
	if unmarshallErr != nil {
		return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal EVM.Workflow config for chain %d", writeEVMInput.ChainID)
	}

	return &evmWorkflow, nil
}

type writeEVMData struct {
	ChainID          uint64
	ChainSelector    uint64
	FromAddress      common.Address
	ForwarderAddress string
	WorkflowConfig   map[string]any // Configuration for EVM.Workflow section
}

const evmWorkflowConfigTemplate = `
	FromAddress = '{{.FromAddress}}'
	ForwarderAddress = '{{.ForwarderAddress}}'
	GasLimitDefault = {{.GasLimitDefault}}
	TxAcceptanceState = {{.TxAcceptanceState}}
	PollPeriod = '{{.PollPeriod}}'
	AcceptanceTimeout = '{{.AcceptanceTimeout}}'
`
