package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/modifier"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
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

	// FirstOCR3ConfigCapabilities lists capability IDs (e.g. "consensus@1.0.0")
	// for which this is the first OCR3 config (no existing config on-chain).
	// Without listing a capability here, the changeset will fail if it cannot
	// read the current config count from the registry, preventing accidental
	// config count collisions.
	FirstOCR3ConfigCapabilities []string `json:"firstOCR3ConfigCapabilities" yaml:"firstOCR3ConfigCapabilities"`

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
	ocr3CapConfigs := make([]ocr3CapConfig, len(config.CapabilityConfigs))
	for i, cc := range config.CapabilityConfigs {
		ocr3CapConfigs[i] = ocr3CapConfig{CapabilityID: cc.Capability.CapabilityID, Config: cc.Config}
	}

	// fetch and increment existing count or start from 1 for FirstOCR3ConfigCapabilities
	configCountFn := func(capID, ocrConfigKey string) (uint64, error) {
		isFirst := isFirstOCR3Config(config.FirstOCR3ConfigCapabilities, capID)

		currentCount, err := ocr3.GetCurrentOCR3ConfigCount(
			capReg, config.DONName, capID, ocrConfigKey,
		)
		if err != nil {
			if !isFirst {
				return 0, fmt.Errorf(
					"failed to read current OCR3 config count for capability %q[%q]: %w. "+
						"Add %q to firstOCR3ConfigCapabilities if this is the initial OCR3 config for this capability",
					capID, ocrConfigKey, err, capID)
			}
			currentCount = 0
		}
		if currentCount == 0 && !isFirst {
			return 0, fmt.Errorf(
				"OCR3 config count is 0 for capability %q[%q], which suggests no prior config exists. "+
					"Add %q to firstOCR3ConfigCapabilities to confirm this is the initial OCR3 config",
				capID, ocrConfigKey, capID)
		}
		return currentCount + 1, nil
	}

	if err := expandOCR3Configs(e, config.RegistryChainSel, nodes, ocr3CapConfigs, configCountFn); err != nil {
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

	// apply modifiers to capability configs
	// currently we add p2pToTransmitterMap to the specConfig for Aptos capabilities
	// more modifiers can be added here as needed
	modifierParams := modifier.CapabilityConfigModifierParams{
		Env:     &e,
		DonName: config.DONName,
		P2PIDs:  p2pIDs,
		Configs: config.CapabilityConfigs,
	}
	for _, mod := range modifier.DefaultCapabilityConfigModifiers() {
		if err := mod.Modify(modifierParams); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("modify capability configs for DON %s: %w", config.DONName, err)
		}
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

// isFirstOCR3Config returns true if capID is listed in the
// firstOCR3ConfigCapabilities slice.
func isFirstOCR3Config(firstCaps []string, capID string) bool {
	for _, c := range firstCaps {
		if c == capID {
			return true
		}
	}
	return false
}
