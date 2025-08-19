package factory

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	creregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
)

// Type aliases for cleaner function signatures

// RuntimeValuesExtractorFn extracts runtime values from node metadata for template substitution.
// chainID is 0 for DON-level capabilities that don't operate on specific chains.
type RuntimeValuesExtractorFn func(chainID uint64, nodeMetadata *cre.NodeMetadata) map[string]any

// CommandBuilderFn constructs the command string for executing a capability binary or built-in capability.
type CommandBuilderFn func(input *cre.JobSpecInput, capabilityConfig cre.CapabilityConfig) (string, error)

type JobNameFn func(chainID uint64, flag cre.CapabilityFlag) string
type IsEnabledFn func(donWithMetadata *cre.DonWithMetadata, nodeSet *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) bool
type EnabledChainsFn func(donTopology *cre.DonTopology, nodeSetInput *cre.CapabilitiesAwareNodeSet, flag cre.CapabilityFlag) []uint64
type ConfigResolverFn func(nodeSetInput *cre.CapabilitiesAwareNodeSet, capabilityConfig cre.CapabilityConfig, chainID uint64, flag cre.CapabilityFlag) (bool, map[string]any, error)

// NoOpExtractor is a no-operation runtime values extractor for DON-level capabilities
// that don't need runtime values extraction from node metadata
var NoOpExtractor RuntimeValuesExtractorFn = func(chainID uint64, nodeMetadata *cre.NodeMetadata) map[string]any {
	return map[string]any{} // Return empty map - DON-level capabilities typically don't need runtime values
}

// BinaryPathBuilder constructs the container path for capability binaries by combining
// the default container directory with the base name of the capability's binary path
var BinaryPathBuilder CommandBuilderFn = func(input *cre.JobSpecInput, capabilityConfig cre.CapabilityConfig) (string, error) {
	containerPath, pathErr := creregistry.DefaultContainerDirectory(input.InfraInput.Type)
	if pathErr != nil {
		return "", errors.Wrapf(pathErr, "failed to get default container directory for infra type %s", input.InfraInput.Type)
	}

	return filepath.Join(containerPath, filepath.Base(capabilityConfig.BinaryPath)), nil
}

// NewCapabilityJobSpecFactory creates a job spec factory for capabilities that operate
// at the DON level without chain-specific configuration (e.g., cron, mock, custom-compute, web-api-*).
// These capabilities use the home chain selector and can have per-DON configuration overrides.
func NewCapabilityJobSpecFactory(
	isEnabledFn IsEnabledFn,
	enabledChainsFn EnabledChainsFn,
	configResolverFn ConfigResolverFn,
	jobNameFn JobNameFn,
) *CapabilityJobSpecFactory {
	return &CapabilityJobSpecFactory{
		isEnabledFn:      isEnabledFn,
		enabledChainsFn:  enabledChainsFn,
		configResolverFn: configResolverFn,
		jobNameFn:        jobNameFn,
	}
}

func (f *CapabilityJobSpecFactory) BuildJobSpecFn(
	capabilityFlag cre.CapabilityFlag,
	configTemplate string,
	runtimeValuesExtractorFn RuntimeValuesExtractorFn,
	commandBuilderFn CommandBuilderFn,
) func(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
	return func(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
		if input.DonTopology == nil {
			return nil, errors.New("topology is nil")
		}

		donToJobSpecs := make(cre.DonsToJobSpecs)

		for donIdx, donWithMetadata := range input.DonTopology.DonsWithMetadata {
			if donIdx >= len(input.CapabilitiesAwareNodeSets) || input.CapabilitiesAwareNodeSets[donIdx] == nil {
				continue
			}

			if f.isEnabledFn != nil && !f.isEnabledFn(donWithMetadata, input.CapabilitiesAwareNodeSets[donIdx], capabilityFlag) {
				continue
			}

			capabilityConfig, ok := input.CapabilityConfigs[capabilityFlag]
			if !ok {
				return nil, errors.Errorf("%s config not found in capabilities config", capabilityFlag)
			}

			command, cmdErr := commandBuilderFn(input, capabilityConfig)
			if cmdErr != nil {
				return nil, errors.Wrap(cmdErr, "failed to get capability command")
			}

			workflowNodeSet, setErr := crenode.FindManyWithLabel(
				donWithMetadata.NodesMetadata,
				&cre.Label{Key: crenode.NodeTypeKey, Value: cre.WorkerNode},
				crenode.EqualLabels,
			)

			if setErr != nil {
				return nil, errors.Wrap(setErr, "failed to find worker nodes")
			}

			// Generate job specs for each enabled chain
			for _, chainIDUint64 := range f.enabledChainsFn(input.DonTopology, input.CapabilitiesAwareNodeSets[donIdx], capabilityFlag) {
				enabled, mergedConfig, rErr := f.configResolverFn(input.CapabilitiesAwareNodeSets[donIdx], capabilityConfig, chainIDUint64, capabilityFlag)
				if rErr != nil {
					return nil, errors.Wrap(rErr, "failed to resolve capability config for chain")
				}
				if !enabled {
					continue
				}

				// Create job specs for each worker node
				for _, workerNode := range workflowNodeSet {
					nodeID, nodeIDErr := crenode.FindLabelValue(workerNode, crenode.NodeIDKey)
					if nodeIDErr != nil {
						return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
					}

					// Apply runtime values to merged config using the runtime value builder
					templateData, aErr := don.ApplyRuntimeValues(mergedConfig, runtimeValuesExtractorFn(chainIDUint64, workerNode))
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

					if err := don.ValidateTemplateSubstitution(configStr, capabilityFlag); err != nil {
						return nil, errors.Wrapf(err, "%s template validation failed", capabilityFlag)
					}

					jobSpec := jobs.WorkerStandardCapability(nodeID, f.jobNameFn(chainIDUint64, capabilityFlag), command, configStr, "")
					donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
				}
			}
		}

		return donToJobSpecs, nil
	}
}

// CapabilityJobSpecFactory is a unified factory that uses strategy functions to handle
// both DON-level and chain-specific capabilities through composition.
type CapabilityJobSpecFactory struct {
	// capabilityFlag           cre.CapabilityFlag
	// configTemplate           string
	// runtimeValuesExtractorFn RuntimeValuesExtractorFn
	// commandBuilderFn         CommandBuilderFn

	// Strategy functions that differ between DON-level and chain-specific capabilities
	jobNameFn        JobNameFn
	isEnabledFn      IsEnabledFn
	enabledChainsFn  EnabledChainsFn
	configResolverFn ConfigResolverFn
}
