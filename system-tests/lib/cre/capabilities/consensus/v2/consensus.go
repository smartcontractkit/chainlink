package consensus

import (
	"bytes"
	"html/template"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr/donlevel"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const flag = cre.ConsensusCapabilityV2
const configTemplate = `'{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}","nodeAddress":"{{.NodeAddress}}"}'`

func New() (*capabilities.Capability, error) {
	return capabilities.New(
		flag,
		capabilities.WithJobSpecFn(jobSpec),
		capabilities.WithCapabilityRegistryV1ConfigFn(registerWithV1),
	)
}

func registerWithV1(donFlags []string, _ *cre.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	if flags.HasFlag(donFlags, flag) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "consensus",
				Version:        "1.0.0-alpha",
				CapabilityType: 2, // CONSENSUS
				ResponseType:   0, // REPORT
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities, nil
}

func buildRuntimeValues(chainID uint64, networkFamily, nodeAddress string) map[string]any {
	return map[string]any{
		"ChainID":       chainID,
		"NetworkFamily": networkFamily,
		"NodeAddress":   nodeAddress,
	}
}

func jobSpec(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
	var generateJobSpec = func(logger zerolog.Logger, chainID uint64, nodeAddress string, mergedConfig map[string]any) (string, error) {
		runtimeFallbacks := buildRuntimeValues(chainID, "evm", nodeAddress)

		templateData, aErr := don.ApplyRuntimeValues(mergedConfig, runtimeFallbacks)
		if aErr != nil {
			return "", errors.Wrap(aErr, "failed to apply runtime values")
		}

		tmpl, err := template.New("consensusConfig").Parse(configTemplate)
		if err != nil {
			return "", errors.Wrap(err, "failed to parse consensus config template")
		}

		var configBuffer bytes.Buffer
		if err := tmpl.Execute(&configBuffer, templateData); err != nil {
			return "", errors.Wrap(err, "failed to execute consensus config template")
		}

		return configBuffer.String(), nil
	}

	var dataStoreOCR3ContractKeyProvider = func(contractName string, chainSelector uint64) datastore.AddressRefKey {
		return datastore.NewAddressRefKey(
			chainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			contractName,
		)
	}

	return ocr.GenerateJobSpecsForStandardCapabilityWithOCR(
		input.DonTopology,
		input.CldEnvironment.DataStore,
		input.CapabilitiesAwareNodeSets,
		input.InfraInput,
		flag,
		func(_ uint64) string {
			return contracts.ConsensusV2ContractQualifier
		},
		dataStoreOCR3ContractKeyProvider,
		donlevel.CapabilityEnabler,
		donlevel.EnabledChainsProvider,
		generateJobSpec,
		donlevel.ConfigMerger,
		input.CapabilityConfigs,
	)
}
