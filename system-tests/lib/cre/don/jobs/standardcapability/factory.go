package factory

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"

	ptypes "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
)

// Type aliases for cleaner function signatures

// RuntimeValuesExtractor extracts runtime values from node metadata for template substitution.
// chainID is 0 for DON-level capabilities that don't operate on specific chains.
type RuntimeValuesExtractor func(chainID uint64, node *cre.Node) map[string]any

// CommandBuilder constructs the command string for executing a capability binary or built-in capability.
type CommandBuilder func(input *cre.JobSpecInput, capabilityConfig cre.CapabilityConfig) (string, error)

// JobNamer constructs the job name for a capability.
type JobNamer func(chainID uint64, flag cre.CapabilityFlag) string

// CapabilityEnabler determines if a capability is enabled for a given DON.
type CapabilityEnabler func(capabilities []string, nodeSet *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) bool

// EnabledChainsProvider provides the list of enabled chains for a given capability.
type EnabledChainsProvider func(donTopology *cre.DonTopology, nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) []uint64

// ConfigResolver resolves the capability config for a given chain.
type ConfigResolver func(nodeSetInput *cre.CapabilitiesAwareNodeSet, capabilityConfig cre.CapabilityConfig, chainID uint64, flag cre.CapabilityFlag) (bool, map[string]any, error)

// NoOpExtractor is a no-operation runtime values extractor for DON-level capabilities
// that don't need runtime values extraction from node metadata
var NoOpExtractor RuntimeValuesExtractor = func(_ uint64, _ *cre.Node) map[string]any {
	return map[string]any{} // Return empty map - DON-level capabilities typically don't need runtime values
}

// BinaryPathBuilder constructs the container path for capability binaries by combining
// the default container directory with the base name of the capability's binary path
var BinaryPathBuilder CommandBuilder = func(input *cre.JobSpecInput, capabilityConfig cre.CapabilityConfig) (string, error) {
	containerPath, pathErr := crecapabilities.DefaultContainerDirectory(input.InfraInput.Type)
	if pathErr != nil {
		return "", errors.Wrapf(pathErr, "failed to get default container directory for infra type %s", input.InfraInput.Type)
	}

	return filepath.Join(containerPath, filepath.Base(capabilityConfig.BinaryPath)), nil
}

// CapabilityJobSpecFactory is a unified factory that uses strategy functions to handle
// both DON-level and chain-specific capabilities through composition.
type CapabilityJobSpecFactory struct {
	// Strategy functions that differ between DON-level and chain-specific capabilities
	jobNamer              JobNamer
	capabilityEnabler     CapabilityEnabler
	enabledChainsProvider EnabledChainsProvider
	configResolver        ConfigResolver
}

// NewCapabilityJobSpecFactory creates a job spec factory for capabilities that operate
// at the DON level without chain-specific configuration (e.g., cron, mock, custom-compute, web-api-*).
// These capabilities use the home chain selector and can have per-DON configuration overrides.
func NewCapabilityJobSpecFactory(
	capabilityEnabler CapabilityEnabler,
	enabledChainsProvider EnabledChainsProvider,
	configResolver ConfigResolver,
	jobNamer JobNamer,
) (*CapabilityJobSpecFactory, error) {
	if capabilityEnabler == nil {
		return nil, errors.New("capability enabler is nil")
	}
	if enabledChainsProvider == nil {
		return nil, errors.New("enabled chains provider is nil")
	}
	if configResolver == nil {
		return nil, errors.New("config resolver is nil")
	}
	if jobNamer == nil {
		return nil, errors.New("job namer is nil")
	}

	return &CapabilityJobSpecFactory{
		capabilityEnabler:     capabilityEnabler,
		enabledChainsProvider: enabledChainsProvider,
		configResolver:        configResolver,
		jobNamer:              jobNamer,
	}, nil
}

func (f *CapabilityJobSpecFactory) BuildJobSpec(
	capabilityFlag cre.CapabilityFlag,
	configTemplate string,
	runtimeValuesExtractor RuntimeValuesExtractor,
	commandBuilder CommandBuilder,
) func(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
	return func(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
		if input.DonTopology == nil {
			return nil, errors.New("topology is nil")
		}
		if runtimeValuesExtractor == nil {
			return nil, errors.New("runtime values extractor is nil")
		}
		if commandBuilder == nil {
			return nil, errors.New("command builder is nil")
		}

		donToJobSpecs := make(cre.DonsToJobSpecs)

		for donIdx, don := range input.DonTopology.Dons.List() {
			if donIdx >= len(input.CapabilitiesAwareNodeSets) || input.CapabilitiesAwareNodeSets[donIdx] == nil {
				continue
			}

			if f.capabilityEnabler != nil && !f.capabilityEnabler(don.Flags, input.CapabilitiesAwareNodeSets[donIdx], capabilityFlag) {
				continue
			}

			capabilityConfig, ok := input.CapabilityConfigs[capabilityFlag]
			if !ok {
				return nil, errors.Errorf("%s config not found in capabilities config. Make sure you have set it in the TOML config", capabilityFlag)
			}

			command, cmdErr := commandBuilder(input, capabilityConfig)
			if cmdErr != nil {
				return nil, errors.Wrap(cmdErr, "failed to get capability command")
			}

			workerNodes, wErr := don.Workers()
			if wErr != nil {
				return nil, errors.Wrap(wErr, "failed to find worker nodes")
			}

			// Generate job specs for each enabled chain
			for _, chainID := range f.enabledChainsProvider(input.DonTopology, input.CapabilitiesAwareNodeSets[donIdx], capabilityFlag) {
				enabled, mergedConfig, rErr := f.configResolver(input.CapabilitiesAwareNodeSets[donIdx], capabilityConfig, chainID, capabilityFlag)
				if rErr != nil {
					return nil, errors.Wrap(rErr, "failed to resolve capability config for chain")
				}
				if !enabled {
					continue
				}

				// Create job specs for each worker node
				for _, workerNode := range workerNodes {
					// Apply runtime values to merged config using the runtime value builder
					templateData, aErr := credon.ApplyRuntimeValues(mergedConfig, runtimeValuesExtractor(chainID, workerNode))
					if aErr != nil {
						return nil, errors.Wrap(aErr, "failed to apply runtime values")
					}

					// Parse and execute template
					tmpl, tmplErr := template.New(capabilityFlag + "-config").Parse(configTemplate)
					if tmplErr != nil {
						return nil, errors.Wrapf(tmplErr, "failed to parse %s config template", capabilityFlag)
					}

					var configBuffer bytes.Buffer
					if err := tmpl.Execute(&configBuffer, templateData); err != nil {
						return nil, errors.Wrapf(err, "failed to execute %s config template", capabilityFlag)
					}
					configStr := configBuffer.String()

					if err := credon.ValidateTemplateSubstitution(configStr, capabilityFlag); err != nil {
						return nil, errors.Wrapf(err, "%s template validation failed", capabilityFlag)
					}

					jobSpec := jobs.WorkerStandardCapability(workerNode.JobDistributorDetails.NodeID, f.jobNamer(chainID, capabilityFlag), command, configStr, "")
					jobSpec.Labels = []*ptypes.Label{{Key: cre.CapabilityLabelKey, Value: &capabilityFlag}}
					donToJobSpecs[don.ID] = append(donToJobSpecs[don.ID], jobSpec)
				}
			}
		}

		return donToJobSpecs, nil
	}
}
