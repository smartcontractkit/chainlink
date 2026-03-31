package changeset

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

var _ cldf.ChangeSetV2[ConfigureCapabilitiesRegistryInput] = ConfigureCapabilitiesRegistry{}

// ConfigureCapabilitiesRegistryInput must be JSON and YAML Serializable with no private fields
type ConfigureCapabilitiesRegistryInput struct {
	ChainSelector uint64 `json:"chainSelector" yaml:"chainSelector"`
	// Deprecated: Use Qualifier instead
	// TODO(PRODCRE-1030): Remove support for CapabilitiesRegistryAddress
	CapabilitiesRegistryAddress string                             `json:"capabilitiesRegistryAddress" yaml:"capabilitiesRegistryAddress"`
	MCMSConfig                  *crecontracts.MCMSConfig           `json:"mcmsConfig,omitempty" yaml:"mcmsConfig,omitempty"`
	Nops                        []CapabilitiesRegistryNodeOperator `json:"nops,omitempty" yaml:"nops,omitempty"`
	Capabilities                []CapabilitiesRegistryCapability   `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
	Nodes                       []CapabilitiesRegistryNodeParams   `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	DONs                        []CapabilitiesRegistryNewDONParams `json:"dons,omitempty" yaml:"dons,omitempty"`
	Qualifier                   string                             `json:"qualifier,omitempty" yaml:"qualifier,omitempty"`
}

type ConfigureCapabilitiesRegistryDeps struct {
	Env           *cldf.Environment
	MCMSContracts *commonchangeset.MCMSWithTimelockState // Required if MCMSConfig input is not nil
}

type ConfigureCapabilitiesRegistry struct{}

func (l ConfigureCapabilitiesRegistry) VerifyPreconditions(e cldf.Environment, config ConfigureCapabilitiesRegistryInput) error {
	if config.CapabilitiesRegistryAddress == "" && config.Qualifier == "" {
		return fmt.Errorf("must set either contract address or qualifier (address: %s, qualifier: %s)", config.CapabilitiesRegistryAddress, config.Qualifier)
	}
	if _, ok := e.BlockChains.EVMChains()[config.ChainSelector]; !ok {
		return fmt.Errorf("chain %d not found in environment", config.ChainSelector)
	}

	return nil
}

func (l ConfigureCapabilitiesRegistry) Apply(e cldf.Environment, config ConfigureCapabilitiesRegistryInput) (cldf.ChangesetOutput, error) {
	// Get MCMS contracts if needed
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if config.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.ChainSelector, *config.MCMSConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	nops := make([]capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams, len(config.Nops))
	for i, nop := range config.Nops {
		nops[i] = nop.ToWrapper()
	}

	capabilities := make([]contracts.RegisterableCapability, len(config.Capabilities))
	for i, cp := range config.Capabilities {
		c, err := cp.ToWrapper()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert capability %d: %w", i, err)
		}
		capabilities[i] = c
	}

	nodes := make([]contracts.NodesInput, len(config.Nodes))
	for i, node := range config.Nodes {
		n, err := node.ToWrapper()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert node %d: %w", i, err)
		}
		nodes[i] = n
	}

	// Generate OCR3 configs for capability entries in new DONs that need them.
	for donIdx := range config.DONs {
		don := &config.DONs[donIdx]

		donNodes, err := p2pIDsToNodeInfo(don.Nodes)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("DON %q: failed to convert P2P IDs to node info: %w", don.Name, err)
		}

		capConfigs := make([]ocr3CapConfig, len(don.CapabilityConfigurations))
		for i, cc := range don.CapabilityConfigurations {
			capConfigs[i] = ocr3CapConfig(cc)
		}

		firstConfigCount := func(_, _ string) (uint64, error) { return 1, nil }

		if err := expandOCR3Configs(e, config.ChainSelector, donNodes, capConfigs, firstConfigCount); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("DON %q: failed to process OCR3 configs: %w", don.Name, err)
		}
	}

	dons := make([]capabilities_registry_v2.CapabilitiesRegistryNewDONParams, len(config.DONs))
	for i, don := range config.DONs {
		d, err := don.ToWrapper(e)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert DON %d: %w", i, err)
		}
		dons[i] = d
	}

	var (
		registryRef  datastore.AddressRefKey
		contractAddr = config.CapabilitiesRegistryAddress
	)

	if config.Qualifier != "" {
		registryRef = pkg.GetCapRegV2AddressRefKey(config.ChainSelector, config.Qualifier)
		contractAddr = ""
	}

	capabilitiesRegistryConfigurationReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		sequences.ConfigureCapabilitiesRegistry,
		sequences.ConfigureCapabilitiesRegistryDeps{
			Env:           &e,
			MCMSContracts: mcmsContracts,
		},
		sequences.ConfigureCapabilitiesRegistryInput{
			RegistryChainSel: config.ChainSelector,
			MCMSConfig:       config.MCMSConfig,
			ContractAddress:  contractAddr,
			RegistryRef:      registryRef,
			Nops:             nops,
			Capabilities:     capabilities,
			Nodes:            nodes,
			DONs:             dons,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure capabilities registry: %w", err)
	}

	reports := make([]operations.Report[any, any], 0)
	reports = append(reports, capabilitiesRegistryConfigurationReport.ToGenericReport())

	return cldf.ChangesetOutput{
		Reports:               reports,
		MCMSTimelockProposals: capabilitiesRegistryConfigurationReport.Output.MCMSTimelockProposals,
	}, nil
}

func p2pIDsToNodeInfo(p2pIDs []string) ([]capabilities_registry_v2.INodeInfoProviderNodeInfo, error) {
	nodes := make([]capabilities_registry_v2.INodeInfoProviderNodeInfo, len(p2pIDs))
	for i, id := range p2pIDs {
		peerID, err := p2pkey.MakePeerID(id)
		if err != nil {
			return nil, fmt.Errorf("invalid P2P ID %q: %w", id, err)
		}
		nodes[i] = capabilities_registry_v2.INodeInfoProviderNodeInfo{
			P2pId: peerID,
		}
	}
	return nodes, nil
}

// reference to a capability config entry for OCR3 - mutations propagate to the caller.
type ocr3CapConfig struct {
	CapabilityID string
	Config       map[string]any
}

type configCountFunc func(capID, ocrConfigKey string) (uint64, error)

// expandOCR3Configs scans capability configs for ocr3Configs entries whose
// offchainConfig sub-key is a map of oracle config parameters, and expands
// them into full OCR3Configs (signers, transmitters, offchain config, etc.)
// using node info from the Job Distributor. Entries where offchainConfig is
// already a base64 string (i.e. already processed) are left untouched.
//
// Example YAML input:
//
//	capabilityConfigurations:
//	  - capabilityID: "consensus@1.0.0"
//	    config:
//	      ocr3Configs:
//	        __default__:
//	          offchainConfig:
//	            deltaProgressMillis: 5000
//	            maxFaultyOracles: 2
//	            transmissionSchedule: [10]
//	          extraSignerFamilies:
//	            - aptos
//
// Example output (the __default__ entry is replaced in-place):
//
//	ocr3Configs:
//	  __default__:
//	    signers: ["base64...", ...]
//	    transmitters: ["base64...", ...]
//	    f: 1
//	    offchainConfig: "base64..."
//	    offchainConfigVersion: 30
//	    configCount: 1
func expandOCR3Configs(
	e cldf.Environment,
	chainSel uint64,
	nodes []capabilities_registry_v2.INodeInfoProviderNodeInfo,
	capConfigs []ocr3CapConfig,
	configCountFn configCountFunc,
) error {
	type genEntry struct {
		capIdx              int
		capID               string
		ocrConfigKey        string
		oracleConfig        *ocr3.OracleConfig
		extraSignerFamilies []string
	}
	var entries []genEntry

	for i, capCfg := range capConfigs {
		ocr3Configs := extractOCR3Configs(capCfg.Config)
		if ocr3Configs == nil {
			continue
		}

		for key, entryCfg := range ocr3Configs {
			entryMap, ok := entryCfg.(map[string]any)
			if !ok {
				continue
			}

			offchainCfg, ok := entryMap["offchainConfig"].(map[string]any)
			if !ok {
				// offchainConfig is absent or already a base64 string (processed); skip.
				continue
			}

			oc, err := parseOracleConfig(offchainCfg)
			if err != nil {
				return fmt.Errorf("capability %q, ocr3Configs[%q].offchainConfig: failed to parse oracle config: %w",
					capCfg.CapabilityID, key, err)
			}

			extraFamilies, err := parseExtraSignerFamilies(entryMap)
			if err != nil {
				return fmt.Errorf("capability %q, ocr3Configs[%q].extraSignerFamilies: %w",
					capCfg.CapabilityID, key, err)
			}

			entries = append(entries, genEntry{
				capIdx:              i,
				capID:               capCfg.CapabilityID,
				ocrConfigKey:        key,
				oracleConfig:        oc,
				extraSignerFamilies: extraFamilies,
			})
		}
	}

	if len(entries) == 0 {
		return nil
	}

	configCounts := make([]uint64, len(entries))
	for i, entry := range entries {
		count, err := configCountFn(entry.capID, entry.ocrConfigKey)
		if err != nil {
			return err
		}
		configCounts[i] = count
	}

	for i, entry := range entries {
		ocrConfig, err := ocr3.ComputeOCR3Config(
			e, chainSel, nodes, *entry.oracleConfig, nil, entry.extraSignerFamilies,
		)
		if err != nil {
			return fmt.Errorf("failed to generate OCR3 config for %q[%q]: %w",
				entry.capID, entry.ocrConfigKey, err)
		}

		if err := ocr3.ValidateOCR2OracleConfig(ocrConfig); err != nil {
			return fmt.Errorf("OCR3 config validation failed for %q[%q]: %w",
				entry.capID, entry.ocrConfigKey, err)
		}

		newConfigCount := configCounts[i]

		e.Logger.Infof("Generated OCR3 config for capability %q (key=%q, configCount=%d, signers=%d, f=%d)",
			entry.capID, entry.ocrConfigKey, newConfigCount, len(ocrConfig.Signers), ocrConfig.F)

		ocr3ConfigMap := ocr3.OCR2OracleConfigToMap(ocrConfig, newConfigCount)

		ocr3Configs := extractOCR3Configs(capConfigs[entry.capIdx].Config)
		ocr3Configs[entry.ocrConfigKey] = ocr3ConfigMap
		capConfigs[entry.capIdx].Config["ocr3Configs"] = ocr3Configs
	}

	return nil
}

// extractOCR3Configs extracts the "ocr3Configs" map from a capability config.
// Returns nil if absent or not a map.
func extractOCR3Configs(config map[string]any) map[string]any {
	if config == nil {
		return nil
	}
	ocr3Configs, ok := config["ocr3Configs"].(map[string]any)
	if !ok {
		return nil
	}
	return ocr3Configs
}

// parseExtraSignerFamilies extracts and validates the optional "extraSignerFamilies"
// string slice from an ocr3Config entry map.
func parseExtraSignerFamilies(entryMap map[string]any) ([]string, error) {
	raw, ok := entryMap["extraSignerFamilies"]
	if !ok {
		return nil, nil
	}
	rawSlice, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("expected []string, got %T", raw)
	}
	families := make([]string, 0, len(rawSlice))
	for i, v := range rawSlice {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("element [%d]: expected string, got %T", i, v)
		}
		families = append(families, s)
	}
	if err := ocr3.ValidateExtraSignerFamilies(families); err != nil {
		return nil, err
	}
	return families, nil
}

// ParseOracleConfig JSON-roundtrips an untyped map into an OracleConfig.
func parseOracleConfig(raw any) (*ocr3.OracleConfig, error) {
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	var oc ocr3.OracleConfig
	if err := json.Unmarshal(jsonBytes, &oc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal oracle config: %w", err)
	}

	return &oc, nil
}
