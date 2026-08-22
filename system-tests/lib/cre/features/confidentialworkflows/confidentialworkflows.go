// Package confidentialworkflows installs the confidential workflows capability
// as a standard CRE feature: it publishes the enclave list to the on-chain
// registry and proposes the capability job on each worker node.
//
// The enclave list comes from the DON's capability config rather than Go values,
// so a topology can declare it and `cre env start` can stand the capability up
// without a test driving the environment in-process.
package confidentialworkflows

import (
	"context"
	"fmt"

	"dario.cat/mergo"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/confidentialcompute"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
)

const flag = cre.ConfidentialWorkflowsCapability

// VersionConfigKey is the capability config value holding the version the
// capability registers under.
const VersionConfigKey = "version"

// defaultVersion is used when the topology does not declare one.
const defaultVersion = "1.0.0"

type ConfidentialWorkflows struct{}

func (o *ConfidentialWorkflows) Flag() cre.CapabilityFlag {
	return flag
}

// PreEnvStartup publishes the enclave list the topology declared, which the
// confidential relay handler reads to discover where it may route requests.
func (o *ConfidentialWorkflows) PreEnvStartup(
	_ context.Context,
	_ zerolog.Logger,
	don *cre.DonMetadata,
	_ *cre.Topology,
	_ *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	if !don.HasFlag(flag) {
		return &cre.PreEnvStartupOutput{}, nil
	}

	nodeSet := don.MustNodeSet()
	enclaves, eErr := confidentialcompute.EnclavesFromConfig(nodeSet, flag)
	if eErr != nil {
		return nil, eErr
	}

	capabilityConfig, cErr := confidentialcompute.RegistryCapabilityConfig(flag, version(don), enclaves)
	if cErr != nil {
		return nil, cErr
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: []keystone_changeset.DONCapabilityWithConfig{capabilityConfig},
	}, nil
}

// PostEnvStartup proposes the capability job on each worker node, sealing the
// capability's API key to every node's workflow key.
func (o *ConfidentialWorkflows) PostEnvStartup(
	ctx context.Context,
	_ zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	if !don.HasFlag(flag) {
		return nil
	}

	capabilityConfig, ok := don.GetCapabilityConfig(flag)
	if !ok {
		return fmt.Errorf("config for '%s' capability not found for %s DON", flag, don.GetName())
	}

	command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryName)
	if cErr != nil {
		return errors.Wrapf(cErr, "failed to get command for %s capability", flag)
	}

	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	encryptedAPIKeys, eErr := confidentialcompute.EncryptedAPIKeys(workerNodes)
	if eErr != nil {
		return eErr
	}

	configJSON, jErr := confidentialcompute.JobConfigJSON(encryptedAPIKeys)
	if jErr != nil {
		return jErr
	}

	input := cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: creEnv.CldfEnvironment.Name,
		DONName:     don.Name,
		JobName:     flag + "-worker",
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: don.Name},
		},
		Template: job_types.Cron,
		Inputs: job_types.JobSpecInput{
			"command": command,
			"config":  configJSON,
		},
	}
	if creEnv.FreshExternalJobIDs {
		input.Inputs["externalJobID"] = uuid.NewString()
	}

	if err := (cre_jobs.ProposeJobSpec{}).VerifyPreconditions(*creEnv.CldfEnvironment, input); err != nil {
		return errors.Wrapf(err, "precondition verification failed for %s worker job", flag)
	}

	report, err := (cre_jobs.ProposeJobSpec{}).Apply(*creEnv.CldfEnvironment, input)
	if err != nil {
		return errors.Wrapf(err, "failed to propose %s worker job spec", flag)
	}

	// Proposing only queues the spec in Job Distributor; without approval the
	// nodes never run it and the capability never registers locally.
	specs := make(map[string][]string)
	for _, r := range report.Reports {
		out, ok := r.Output.(cre_jobs_ops.ProposeStandardCapabilityJobOutput)
		if !ok {
			return fmt.Errorf("unable to cast to ProposeStandardCapabilityJobOutput, actual type: %T", r.Output)
		}
		if mErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice); mErr != nil {
			return errors.Wrapf(mErr, "failed to merge %s worker job specs", flag)
		}
	}

	if aErr := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs); aErr != nil {
		return errors.Wrapf(aErr, "failed to approve %s worker jobs", flag)
	}

	return nil
}

// version returns the version the capability registers under, from the DON's
// capability config when declared.
func version(don *cre.DonMetadata) string {
	cfg, ok := don.CapabilityConfigs[flag]
	if !ok || cfg.Values == nil {
		return defaultVersion
	}

	if v, found := cfg.Values[VersionConfigKey]; found {
		if s, isString := v.(string); isString && s != "" {
			return s
		}
	}

	return defaultVersion
}
