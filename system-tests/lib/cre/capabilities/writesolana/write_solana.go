package writesolana

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/gagliardetto/solana-go"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

const flag = cre.WriteSolanaCapability

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

	if slices.Contains(nodeSetInput.Capabilities, flag) {
		// TODO PLEX-296
		// fullName := solana.GenerateName()
		fullName := "write_solana_devnet@1.0.0"
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
	if !flags.HasFlagForAnyChain(input.NodeSet.ComputedCapabilities, flag) {
		return existingConfigs, nil
	}

	data := solanaInput{}
	for _, bcOut := range input.BlockchainOutput {
		if bcOut.SolChain == nil {
			continue
		}
		data.ChainSelector = bcOut.SolChain.ChainSelector
		// find Solana forwarder address
		forwarders := input.Datastore.Addresses().Filter(datastore.AddressRefByChainSelector(data.ChainSelector))
		for _, addr := range forwarders {
			if addr.Type == ks_sol.ForwarderState {
				data.ForwarderState = addr.Address
				continue
			}
			data.ForwarderAddress = addr.Address
		}

		break
	}

	workerNodes, wErr := input.DonMetadata.Workers()
	if wErr != nil {
		return nil, errors.Wrap(wErr, "failed to find worker nodes")
	}

	for _, workerNode := range workerNodes {
		chainID, chErr := chainselectors.SolanaChainIdFromSelector(data.ChainSelector)
		if chErr != nil {
			return nil, errors.Wrapf(chErr, "failed to get Solana chain ID from selector %d", data.ChainSelector)
		}

		key, ok := workerNode.Keys.Solana[chainID]
		if !ok {
			return nil, errors.Errorf("missing Solana key for chainID %s on node index %d", chainID, workerNode.Index)
		}
		data.FromAddress = key.PublicAddress

		if input.CapabilityConfigs == nil {
			return nil, errors.New("additional capabilities configs are nil, but are required to configure the write-solana capability")
		}

		if writeSolConfig, ok := input.CapabilityConfigs[cre.WriteSolanaCapability]; ok {
			mergedConfig := envconfig.ResolveCapabilityConfigForDON(
				cre.WriteSolanaCapability,
				writeSolConfig.Config,
				nil,
			)

			runtimeValues := map[string]any{
				"FromAddress":      data.FromAddress.String(),
				"ForwarderAddress": data.ForwarderAddress,
				"ForwarderState":   data.ForwarderState,
			}

			var mErr error
			data.WorkflowConfig, mErr = don.ApplyRuntimeValues(mergedConfig, runtimeValues)
			if mErr != nil {
				return nil, errors.Wrap(mErr, "failed to apply runtime values")
			}
		} else {
			fmt.Println("sol config not found")
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

		if len(typedConfig.Solana) != 1 {
			return nil, fmt.Errorf("only 1 Solana chain is supported, but found %d for node at index %d", len(typedConfig.Solana), workerNode.Index)
		}

		if typedConfig.Solana[0].ChainID == nil {
			return nil, fmt.Errorf("solana chainID is nil for node at index %d", workerNode.Index)
		}

		var solCfg solcfg.WorkflowConfig

		// Execute template with chain's workflow configuration
		tmpl, err := template.New("solanaWorkflowConfig").Parse(solWorkflowConfigTemplate)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse Solana workflow config template")
		}
		var configBuffer bytes.Buffer
		if executeErr := tmpl.Execute(&configBuffer, data.WorkflowConfig); executeErr != nil {
			return nil, errors.Wrap(executeErr, "failed to execute Solana workflow config template")
		}

		configStr := configBuffer.String()

		if err := don.ValidateTemplateSubstitution(configStr, flag); err != nil {
			return nil, errors.Wrapf(err, "%s template validation failed", flag)
		}

		unmarshallErr = toml.Unmarshal([]byte(configStr), &solCfg)
		if unmarshallErr != nil {
			return nil, errors.Wrap(unmarshallErr, "failed to unmarshal Solana.Workflow config")
		}

		typedConfig.Solana[0].Workflow = solCfg

		marshalledConfig, mErr := toml.Marshal(typedConfig)
		if mErr != nil {
			return nil, errors.Wrapf(mErr, "failed to marshal config for node index %d", workerNode.Index)
		}

		existingConfigs[workerNode.Index] = string(marshalledConfig)
	}

	return existingConfigs, nil
}

type solanaInput struct {
	ChainSelector    uint64
	FromAddress      solana.PublicKey
	ForwarderAddress string
	ForwarderState   string
	HasWrite         bool
	WorkflowConfig   map[string]any // Configuration for Solana.Workflow section
}

const solWorkflowConfigTemplate = `
		ForwarderAddress = '{{.ForwarderAddress}}'
		FromAddress      = '{{.FromAddress}}'
		ForwarderState   = '{{.ForwarderState}}'
		PollPeriod = '{{.PollPeriod}}'
		AcceptanceTimeout = '{{.AcceptanceTimeout}}'
		TxAcceptanceState = {{.TxAcceptanceState}}
		Local = {{.Local}}
	`
