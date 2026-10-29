package queue

import (
	"context"
	"fmt"
	"strconv"

	"dario.cat/mergo"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
)

const flag = cre.ConsensusQueueCapability

const consensusQueueLabelledName = cre.ConsensusQueueCapability

type ConsensusQueue struct{}

func (o *ConsensusQueue) Flag() cre.CapabilityFlag { return flag }
func (o *ConsensusQueue) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	if !don.HasFlag(flag) {
		return nil, nil
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: []keystone_changeset.DONCapabilityWithConfig{{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName: consensusQueueLabelledName,
				Version:      "1.0.0",
			},
			Config: &capabilitiespb.CapabilityConfig{
				LocalOnly: don.HasOnlyLocalCapabilities(),
			},
			UseCapRegOCRConfig: true,
		}},
		CapabilityToOCR3Config: map[string]*ocr3.OracleConfig{
			//TODO don't we need DefaultOCR3_1Config?
			consensusQueueLabelledName: contracts.DefaultOCR3Config(),
		},
	}, nil
}

func (o *ConsensusQueue) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	if !don.HasFlag(flag) {
		return nil
	}

	err := createJobs(ctx, creEnv, don, dons)
	if err != nil {
		return fmt.Errorf("failed to create consensusQueue jobs: %w", err)
	}
	//TODO more?
	return nil
}

func createJobs(
	ctx context.Context,
	creEnv *cre.Environment,
	don *cre.Don,
	dons *cre.Dons,
) error {
	bootstrap, isBootstrap := dons.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	_, ocrPeeringCfg, err := cre.PeeringCfgs(bootstrap)
	if err != nil {
		return errors.Wrap(err, "failed to get peering configs")
	}

	input := cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: creEnv.CldfEnvironment.Name,
		DONName:     don.Name,
		JobName:     cre.ConsensusQueueCapability + "-worker",
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: cre.ConsensusQueueCapability},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: don.Name},
		},
		Template: job_types.OCR3,
		Inputs: job_types.JobSpecInput{
			"chainSelectorEVM":     creEnv.RegistryChainSelector,
			"contractQualifier":    cre.ConsensusQueueCapability,
			"templateName":         cre.ConsensusQueueCapability,
			"bootstrapperOCR3Urls": []string{ocrPeeringCfg.OCRBootstrapperPeerID + "@" + ocrPeeringCfg.OCRBootstrapperHost + ":" + strconv.Itoa(ocrPeeringCfg.Port)},
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
		out, ok := r.Output.(cre_jobs_ops.ProposeOCR3JobOutput)
		if !ok {
			return fmt.Errorf("unable to cast to ProposeOCR3JobOutput, actual type: %T", r.Output)
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
