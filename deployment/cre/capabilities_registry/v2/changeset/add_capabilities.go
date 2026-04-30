package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

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
)

var _ cldf.ChangeSetV2[AddCapabilitiesInput] = AddCapabilities{}

type AddCapabilitiesInput struct {
	RegistryChainSel  uint64 `json:"registryChainSel" yaml:"registryChainSel"`
	RegistryQualifier string `json:"registryQualifier" yaml:"registryQualifier"`

	MCMSConfig *crecontracts.MCMSConfig `json:"mcmsConfig" yaml:"mcmsConfig"`

	// DonCapabilityConfigs maps DON name to the list of capability configs for that DON.
	DonCapabilityConfigs map[string][]contracts.CapabilityConfig `json:"donCapabilityConfigs" yaml:"donCapabilityConfigs"`

	// Force indicates whether to force the update even if we cannot validate that all forwarder contracts are ready to accept the new configure version.
	// This is very dangerous, and could break the whole platform if the forwarders are not ready. Be very careful with this option.
	Force bool `json:"force" yaml:"force"`

	// FirstOCR3ConfigCapabilities maps DON name to capability IDs for which this
	// is the first OCR3 config (no existing config on-chain). Without listing a
	// capability here, the changeset will fail if it cannot read the current config
	// count from the registry, preventing accidental config count collisions.
	FirstOCR3ConfigCapabilities map[string][]string `json:"firstOCR3ConfigCapabilities" yaml:"firstOCR3ConfigCapabilities"`
}

type AddCapabilities struct{}

func (u AddCapabilities) VerifyPreconditions(_ cldf.Environment, config AddCapabilitiesInput) error {
	if len(config.DonCapabilityConfigs) == 0 {
		return errors.New("donCapabilityConfigs must contain at least one DON entry")
	}
	for donName, configs := range config.DonCapabilityConfigs {
		if donName == "" {
			return errors.New("donCapabilityConfigs keys cannot be empty strings")
		}
		if len(configs) == 0 {
			return fmt.Errorf("donCapabilityConfigs[%q] must contain at least one capability config", donName)
		}
	}
	return nil
}

func (u AddCapabilities) Apply(e cldf.Environment, config AddCapabilitiesInput) (cldf.ChangesetOutput, error) {
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

	for donName, donCapConfigs := range config.DonCapabilityConfigs {
		_, nodes, donErr := sequences.GetDonNodes(donName, capReg)
		if donErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("DON %q: failed to get nodes: %w", donName, donErr)
		}

		ocr3CapConfigs := make([]ocr3CapConfig, len(donCapConfigs))
		for i, cc := range donCapConfigs {
			ocr3CapConfigs[i] = ocr3CapConfig{CapabilityID: cc.Capability.CapabilityID, Config: cc.Config}
		}

		firstCaps := config.FirstOCR3ConfigCapabilities[donName]

		configCountFn := func(capID, ocrConfigKey string) (uint64, error) {
			isFirst := isFirstOCR3Config(firstCaps, capID)

			currentCount, countErr := ocr3.GetCurrentOCR3ConfigCount(capReg, donName, capID, ocrConfigKey)
			if countErr != nil {
				if !isFirst {
					return 0, fmt.Errorf(
						"failed to read current OCR3 config count for capability %q[%q] in DON %q: %w. "+
							"Add %q to firstOCR3ConfigCapabilities[%q] if this is the initial OCR3 config",
						capID, ocrConfigKey, donName, countErr, capID, donName)
				}
				currentCount = 0
			}
			if currentCount == 0 && !isFirst {
				return 0, fmt.Errorf(
					"OCR3 config count is 0 for capability %q[%q] in DON %q, which suggests no prior config exists. "+
						"Add %q to firstOCR3ConfigCapabilities[%q] to confirm this is the initial OCR3 config",
					capID, ocrConfigKey, donName, capID, donName)
			}
			return currentCount + 1, nil
		}

		if expandErr := expandOCR3Configs(e, config.RegistryChainSel, nodes, ocr3CapConfigs, configCountFn); expandErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("DON %q: failed to expand OCR3 configs: %w", donName, expandErr)
		}
	}

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		sequences.AddCapabilities,
		sequences.AddCapabilitiesDeps{Env: &e, MCMSContracts: mcmsContracts},
		sequences.AddCapabilitiesInput{
			RegistryRef:          registryRef,
			DonCapabilityConfigs: config.DonCapabilityConfigs,
			Force:                config.Force,
			MCMSConfig:           config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		Reports:               seqReport.ExecutionReports,
		MCMSTimelockProposals: seqReport.Output.Proposals,
	}, nil
}
