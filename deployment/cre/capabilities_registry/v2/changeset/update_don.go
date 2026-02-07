package changeset

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var _ cldf.ChangeSetV2[UpdateDONInput] = UpdateDON{}

type UpdateDONInput struct {
	RegistryQualifier string `json:"registryQualifier" yaml:"registryQualifier"`
	RegistryChainSel  uint64 `json:"registryChainSel" yaml:"registryChainSel"`

	// DONName to update, this is required
	DONName string `json:"donName" yaml:"donName"`
	// NewDonName is optional
	NewDonName string `json:"newDonName" yaml:"newDonName"`

	CapabilityConfigs []contracts.CapabilityConfig `json:"capabilityConfigs" yaml:"capabilityConfigs"` // if Config subfield is nil, a default config is used

	// Force indicates whether to force the update even if we cannot validate that all forwarder contracts are ready to accept the new configure version.
	// This is very dangerous, and could break the whole platform if the forwarders are not ready. Be very careful with this option.
	Force bool `json:"force" yaml:"force"`

	// FirstOCR3Config must be set to true when this is the first OCR3 config
	// for a capability (no existing config on-chain). Without this flag, the
	// changeset will fail if it cannot read the current config count from the
	// registry, preventing accidental config count collisions.
	FirstOCR3Config bool `json:"firstOCR3Config" yaml:"firstOCR3Config"`

	MCMSConfig *crecontracts.MCMSConfig `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type UpdateDON struct{}

func (u UpdateDON) VerifyPreconditions(_ cldf.Environment, config UpdateDONInput) error {
	if config.DONName == "" {
		return errors.New("must provide a non-empty DONName")
	}

	return nil
}

func (u UpdateDON) Apply(e cldf.Environment, config UpdateDONInput) (cldf.ChangesetOutput, error) {
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if config.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.RegistryChainSel, *config.MCMSConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	registryRef := pkg.GetCapRegV2AddressRefKey(config.RegistryChainSel, config.RegistryQualifier)

	chain, ok := e.BlockChains.EVMChains()[config.RegistryChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain not found for selector %d", config.RegistryChainSel)
	}

	registryAddressRef, err := e.DataStore.Addresses().Get(registryRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get registry address: %w", err)
	}

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		common.HexToAddress(registryAddressRef.Address), chain.Client,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create CapabilitiesRegistry: %w", err)
	}

	don, nodes, err := sequences.GetDonNodes(config.DONName, capReg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get DON %s nodes: %w", config.DONName, err)
	}

	// Generate OCR3 configs for capability entries that need them.
	if err := processOCR3Configs(e, &config, capReg, nodes); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to process OCR3 configs: %w", err)
	}

	p2pIDs := make([]p2pkey.PeerID, 0)
	for _, node := range nodes {
		p2pIDs = append(p2pIDs, node.P2pId)
	}

	// Create the appropriate strategy
	strategy, err := strategies.CreateStrategy(
		chain,
		e,
		config.MCMSConfig,
		mcmsContracts,
		capReg.Address(),
		contracts.UpdateDONDescription,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create strategy: %w", err)
	}

	updateDonReport, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.UpdateDON,
		contracts.UpdateDONDeps{
			Env:                  &e,
			Strategy:             strategy,
			CapabilitiesRegistry: capReg,
		},
		contracts.UpdateDONInput{
			ChainSelector:     config.RegistryChainSel,
			P2PIDs:            p2pIDs,
			CapabilityConfigs: config.CapabilityConfigs,
			DonName:           config.DONName,
			NewDonName:        config.NewDonName,
			F:                 don.F,
			IsPrivate:         !don.IsPublic,
			Force:             config.Force,
			MCMSConfig:        config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to update DON %s: %w", config.DONName, err)
	}

	var proposals []mcmslib.TimelockProposal

	if updateDonReport.Output.Operation != nil {
		proposal, mcmsErr := strategy.BuildProposal([]types.BatchOperation{*updateDonReport.Output.Operation})
		if mcmsErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal for UpdateDON on chain %d: %w", config.RegistryChainSel, mcmsErr)
		}

		proposals = append(proposals, *proposal)
	}

	return cldf.ChangesetOutput{
		Reports:               []operations.Report[any, any]{updateDonReport.ToGenericReport()},
		MCMSTimelockProposals: proposals,
	}, nil
}

// processOCR3Configs scans capability configs for ocr3Configs entries missing
// "signers". Those entries are treated as oracle config parameters and expanded
// into full OCR3Configs (signers, transmitters, offchain config, etc.) using
// node info from the Job Distributor.
//
// Entries that already contain "signers" are left as-is.
//
// Example YAML input:
//
//	capabilityConfigurations:
//	  - capability:
//	      capabilityID: "offchain_reporting@1.0.0"
//	    config:
//	      ocr3Configs:
//	        __default__:
//	          deltaProgressMillis: 5000
//	          maxFaultyOracles: 2
//	          transmissionSchedule: [10]
func processOCR3Configs(
	e cldf.Environment,
	config *UpdateDONInput,
	capReg *capabilities_registry_v2.CapabilitiesRegistry,
	nodes []capabilities_registry_v2.INodeInfoProviderNodeInfo,
) error {
	type genEntry struct {
		capIdx       int
		capID        string
		ocrConfigKey string
		oracleConfig *ocr3.OracleConfig
	}
	var entries []genEntry

	for i, capCfg := range config.CapabilityConfigs {
		ocr3Configs := ExtractOCR3Configs(capCfg.Config)
		if ocr3Configs == nil {
			continue
		}

		for key, entryCfg := range ocr3Configs {
			entryMap, ok := entryCfg.(map[string]any)
			if !ok {
				continue
			}

			if NeedsOCR3Generation(entryMap) {
				oc, err := ParseOracleConfig(entryMap)
				if err != nil {
					return fmt.Errorf("capability %q, ocr3Configs[%q]: failed to parse oracle config: %w",
						capCfg.Capability.CapabilityID, key, err)
				}

				entries = append(entries, genEntry{
					capIdx:       i,
					capID:        capCfg.Capability.CapabilityID,
					ocrConfigKey: key,
					oracleConfig: oc,
				})
			}
		}
	}

	if len(entries) == 0 {
		return nil
	}

	for _, entry := range entries {
		ocrConfig, err := ocr3.ComputeOCR3Config(
			e, config.RegistryChainSel, nodes, *entry.oracleConfig, nil,
		)
		if err != nil {
			return fmt.Errorf("failed to generate OCR3 config for %q[%q]: %w",
				entry.capID, entry.ocrConfigKey, err)
		}

		if err := ocr3.ValidateOCR2OracleConfig(ocrConfig); err != nil {
			return fmt.Errorf("OCR3 config validation failed for %q[%q]: %w",
				entry.capID, entry.ocrConfigKey, err)
		}

		currentCount, err := ocr3.GetCurrentOCR3ConfigCount(
			capReg, config.DONName, entry.capID, entry.ocrConfigKey,
		)
		if err != nil {
			if !config.FirstOCR3Config {
				return fmt.Errorf(
					"failed to read current OCR3 config count for capability %q[%q]: %w. "+
						"Set firstOCR3Config=true if this is the initial OCR3 config for this capability",
					entry.capID, entry.ocrConfigKey, err)
			}
			currentCount = 0
		}
		if currentCount == 0 && !config.FirstOCR3Config {
			return fmt.Errorf(
				"OCR3 config count is 0 for capability %q[%q], which suggests no prior config exists. "+
					"Set firstOCR3Config=true to confirm this is the initial OCR3 config",
				entry.capID, entry.ocrConfigKey)
		}
		newConfigCount := currentCount + 1

		e.Logger.Infof("Generated OCR3 config for capability %q (key=%q, configCount=%d, signers=%d, f=%d)",
			entry.capID, entry.ocrConfigKey, newConfigCount, len(ocrConfig.Signers), ocrConfig.F)

		ocr3ConfigMap := ocr3.OCR2OracleConfigToMap(ocrConfig, newConfigCount)

		ocr3Configs := ExtractOCR3Configs(config.CapabilityConfigs[entry.capIdx].Config)
		ocr3Configs[entry.ocrConfigKey] = ocr3ConfigMap
		config.CapabilityConfigs[entry.capIdx].Config["ocr3Configs"] = ocr3Configs
	}

	return nil
}

// NeedsOCR3Generation returns true if an ocr3Configs entry is missing signers,
// indicating it contains oracle config parameters that need to be expanded
// into a full OCR3Config.
func NeedsOCR3Generation(entry map[string]any) bool {
	signers, ok := entry["signers"]
	if !ok {
		return true
	}
	if signers == nil {
		return true
	}
	if arr, ok := signers.([]any); ok && len(arr) == 0 {
		return true
	}
	return false
}

// ExtractOCR3Configs extracts the "ocr3Configs" map from a capability config.
// Returns nil if absent or not a map.
func ExtractOCR3Configs(config map[string]any) map[string]any {
	if config == nil {
		return nil
	}
	ocr3Configs, ok := config["ocr3Configs"].(map[string]any)
	if !ok {
		return nil
	}
	return ocr3Configs
}

// ParseOracleConfig JSON-roundtrips an untyped map into an OracleConfig.
func ParseOracleConfig(raw any) (*ocr3.OracleConfig, error) {
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
