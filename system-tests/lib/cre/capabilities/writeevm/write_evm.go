package writeevm

import (
	"bytes"
	"fmt"
	"math/big"
	"text/template"

	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	evmworkflow "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	chainlinkbig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
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

	workflowNodeSet, wErr := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
	if wErr != nil {
		return nil, errors.Wrap(wErr, "failed to find worker nodes")
	}

	for i := range workflowNodeSet {
		var nodeIndex int
		for _, label := range workflowNodeSet[i].Labels {
			if label.Key == node.IndexKey {
				var nErr error
				nodeIndex, nErr = strconv.Atoi(label.Value)
				if nErr != nil {
					return nil, errors.Wrap(nErr, "failed to convert node index to int")
				}
			}
		}

		// // get all the forwarders and add workflow config (FromAddress + Forwarder) for chains that have write-evm enabled
		data := []writeEVMData{}
		for idx, chainID := range input.NodeSet.ChainCapabilities[flag].EnabledChains {
			chain, exists := chain_selectors.ChainByEvmChainID(chainID)
			if !exists {
				return nil, errors.Errorf("failed to find selector for chain ID %d", chainID)
			}

			addrsForChains, addErr := input.AddressBook.AddressesForChain(chain.Selector)
			if addErr != nil {
				return nil, errors.Wrap(addErr, "failed to get addresses from address book")
			}

			for addr, addrValue := range addrsForChains {
				if addrValue.Type == keystone_changeset.KeystoneForwarder {
					input := writeEVMData{}
					input.ForwarderAddress = addr
					input.ChainID = chainID
					input.ChainSelector = chain.Selector

					expectedAddressKey := node.AddressKeyFromSelector(input.ChainSelector)
					for _, label := range workflowNodeSet[i].Labels {
						if label.Key == expectedAddressKey {
							if label.Value == "" {
								return nil, errors.Errorf("%s label value is empty", expectedAddressKey)
							}
							input.FromAddress = common.HexToAddress(label.Value)
							break
						}
					}
					if input.FromAddress == (common.Address{}) {
						return nil, errors.Errorf("failed to get from address for chain %d", input.ChainSelector)
					}

					data = append(data, input)
				}
			}

			if input.CapabilityConfigs == nil {
				return nil, errors.New("additional capabilities configs are nil, but are required to configure the write-evm capability")
			}

			if writeEvmConfig, ok := input.CapabilityConfigs[cre.WriteEVMCapability]; ok {
				enabled, mergedConfig, rErr := envconfig.ResolveCapabilityForChain(
					cre.WriteEVMCapability,
					input.NodeSet.ChainCapabilities,
					writeEvmConfig.Config,
					data[idx].ChainID,
				)
				if rErr != nil {
					return nil, errors.Wrapf(rErr, "failed to resolve write-evm config for chain %d", data[idx].ChainID)
				}

				if !enabled {
					// This should never happen, but guard anyway. We have already checked that the capability is enabled in the chain capabilities, when we generated the workerEVMInputs.
					continue
				}

				runtimeValues := map[string]any{
					"FromAddress":      data[idx].FromAddress.Hex(),
					"ForwarderAddress": data[idx].ForwarderAddress,
				}

				var mErr error
				data[idx].WorkflowConfig, mErr = don.ApplyRuntimeValues(mergedConfig, runtimeValues)
				if mErr != nil {
					return nil, errors.Wrap(mErr, "failed to apply runtime values")
				}
			}
		}

		if len(existingConfigs) < nodeIndex+1 {
			return nil, errors.Errorf("missing config for node index %d", nodeIndex)
		}

		currentConfig := existingConfigs[nodeIndex]

		var typedConfig corechainlink.Config
		unmarshallErr := toml.Unmarshal([]byte(currentConfig), &typedConfig)
		if unmarshallErr != nil {
			return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal config for node index %d", nodeIndex)
		}

		if len(typedConfig.EVM) < len(data) {
			return nil, fmt.Errorf("not enough EVM chains configured in node index %d to add write-evm config. Expected at least %d chains, but found %d", nodeIndex, len(data), len(typedConfig.EVM))
		}

		for _, writeEVMInput := range data {
			found := false
		INNER:
			for idx, evmChain := range typedConfig.EVM {
				if evmChain.ChainID.Cmp(chainlinkbig.New(big.NewInt(libc.MustSafeInt64(writeEVMInput.ChainID)))) == 0 {
					var evmWorkflow evmworkflow.Workflow

					// Execute template with chain's workflow configuration
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

					typedConfig.EVM[idx].Workflow = evmWorkflow
					typedConfig.EVM[idx].Transactions.ForwardersEnabled = ptr.Ptr(true)

					found = true
					break INNER
				}
			}

			if !found {
				return nil, fmt.Errorf("failed to find EVM chain with ID %d in the config of node index %d to add write-evm config", writeEVMInput.ChainID, nodeIndex)
			}
		}

		marshalledConfig, mErr := toml.Marshal(typedConfig)
		if mErr != nil {
			return nil, errors.Wrapf(mErr, "failed to marshal config for node index %d", nodeIndex)
		}

		existingConfigs[nodeIndex] = string(marshalledConfig)
	}

	return existingConfigs, nil
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
