package evm

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr/chainlevel"
)

const flag = cre.EVMCapability
const configTemplate = `'{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}","logTriggerPollInterval":{{.LogTriggerPollInterval}}, "creForwarderAddress":"{{.CreForwarderAddress}}","receiverGasMinimum":{{.ReceiverGasMinimum}},"nodeAddress":"{{.NodeAddress}}"}'`
const registrationRefresh = 20 * time.Second
const registrationExpiry = 60 * time.Second

func New() (*capabilities.Capability, error) {
	return capabilities.New(
		flag,
		capabilities.WithJobSpecFn(jobSpec),
		capabilities.WithCapabilityRegistryV1ConfigFn(registerWithV1),
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
		selector, selectorErr := chainselectors.SelectorFromChainId(chainID)
		if selectorErr != nil {
			return nil, errors.Wrapf(selectorErr, "failed to get selector from chainID: %d", chainID)
		}
		faultyNodes, faultyErr := nodeSetInput.MaxFaultyNodes()
		if faultyErr != nil {
			return nil, errors.Wrap(faultyErr, "failed to get faulty nodes")
		}
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "evm" + ":ChainSelector:" + strconv.FormatUint(selector, 10),
				Version:        "1.0.0",
				CapabilityType: 0, // TRIGGER
			},
			Config: &capabilitiespb.CapabilityConfig{
				RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
					RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
						// needed for message_cache.go#Ready(), without these events from the capability will never be accepted
						RegistrationRefresh:     durationpb.New(registrationRefresh),
						RegistrationExpiry:      durationpb.New(registrationExpiry),
						MinResponsesToAggregate: faultyNodes + 1,
					},
				},
			},
		})
	}

	return capabilities, nil
}

// buildRuntimeValues creates runtime-generated  values for any keys not specified in TOML
func buildRuntimeValues(chainID uint64, networkFamily, creForwarderAddress, nodeAddress string) map[string]any {
	return map[string]any{
		"ChainID":             chainID,
		"NetworkFamily":       networkFamily,
		"CreForwarderAddress": creForwarderAddress,
		"NodeAddress":         nodeAddress,
	}
}

func jobSpec(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
	var generateJobSpec = func(logger zerolog.Logger, chainID uint64, nodeAddress string, mergedConfig map[string]any) (string, error) {
		cs, ok := chainselectors.EvmChainIdToChainSelector()[chainID]
		if !ok {
			return "", fmt.Errorf("chain selector not found for chainID: %d", chainID)
		}

		creForwarderKey := datastore.NewAddressRefKey(
			cs,
			datastore.ContractType(keystone_changeset.KeystoneForwarder.String()),
			semver.MustParse("1.0.0"),
			"",
		)
		creForwarderAddress, err := input.CldEnvironment.DataStore.Addresses().Get(creForwarderKey)
		if err != nil {
			return "", errors.Wrap(err, "failed to get CRE Forwarder address")
		}

		logger.Debug().Msgf("Found CRE Forwarder contract on chain %d at %s", chainID, creForwarderAddress.Address)

		runtimeFallbacks := buildRuntimeValues(chainID, "evm", creForwarderAddress.Address, nodeAddress)

		templateData, aErr := don.ApplyRuntimeValues(mergedConfig, runtimeFallbacks)
		if aErr != nil {
			return "", errors.Wrap(aErr, "failed to apply runtime values")
		}

		tmpl, err := template.New("evmConfig").Parse(configTemplate)
		if err != nil {
			return "", errors.Wrapf(err, "failed to parse %s config template", flag)
		}

		var configBuffer bytes.Buffer
		if err := tmpl.Execute(&configBuffer, templateData); err != nil {
			return "", errors.Wrapf(err, "failed to execute %s config template", flag)
		}

		configStr := configBuffer.String()

		if err := don.ValidateTemplateSubstitution(configStr, flag); err != nil {
			return "", errors.Wrapf(err, "%s template validation failed", flag)
		}

		return configStr, nil
	}

	return ocr.GenerateJobSpecsForStandardCapabilityWithOCR(
		input.DonTopology,
		input.CldEnvironment.DataStore,
		input.CapabilitiesAwareNodeSets,
		input.InfraInput,
		flag,
		contracts.GetCapabilityContractIdentifier,
		chainlevel.CapabilityEnabler,
		chainlevel.EnabledChainsProvider,
		generateJobSpec,
		chainlevel.ConfigMerger,
		input.CapabilityConfigs,
	)
}
