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
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/ocr3_1"
	ocr3_changeset "github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset"
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

	workers, err := don.Workers()
	if err != nil {
		return nil, fmt.Errorf("failed to get workers for DON %s: %w", don.Name, err)
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: []keystone_changeset.DONCapabilityWithConfig{{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   consensusQueueLabelledName,
				Version:        "1.0.0",
				CapabilityType: 2, // CONSENSUS
				ResponseType:   0, // REPORT (not used)
			},
			Config: &capabilitiespb.CapabilityConfig{
				LocalOnly: don.HasOnlyLocalCapabilities(),
			},
			UseCapRegOCRConfig: true,
		}},
		// The plugin runs on the OCR 3.1 oracle, which rejects OCR3 (version 30) offchain configs.
		CapabilityToOCR3_1Config: map[string]*ocr3_1.V3_1OracleConfig{
			consensusQueueLabelledName: contracts.DefaultOCR3_1Config(len(workers)),
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

	if err := deployContract(creEnv, testLogger); err != nil {
		return fmt.Errorf("failed to deploy consensusQueue OCR3 contract: %w", err)
	}

	err := createJobs(ctx, creEnv, don, dons)
	if err != nil {
		return fmt.Errorf("failed to create consensusQueue jobs: %w", err)
	}
	//TODO more?
	return nil
}

// deployContract deploys the OCR3 contract the job spec uses as its contractID.
// The OCR config itself comes from the Capabilities Registry.
func deployContract(creEnv *cre.Environment, testLogger zerolog.Logger) error {
	input := ocr3_changeset.DeployOCR3Input{
		ChainSelector: creEnv.RegistryChainSelector,
		Qualifier:     cre.ConsensusQueueCapability,
		Labels:        []string{cre.ConsensusQueueCapability},
	}
	if err := (ocr3_changeset.DeployOCR3{}).VerifyPreconditions(*creEnv.CldfEnvironment, input); err != nil {
		return errors.Wrap(err, "precondition verification failed")
	}
	out, err := (ocr3_changeset.DeployOCR3{}).Apply(*creEnv.CldfEnvironment, input)
	if err != nil {
		return err
	}
	contracts.MergeAllDataStores(creEnv, out)

	refKey := pkg.GetOCR3CapabilityAddressRefKey(creEnv.RegistryChainSelector, cre.ConsensusQueueCapability)
	addrRef, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(refKey)
	if err != nil {
		return errors.Wrap(err, "failed to get deployed contract address")
	}
	testLogger.Info().Msgf("Deployed consensus-queue OCR3 contract on chain %d at %s", creEnv.RegistryChainSelector, addrRef.Address)
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

	peeringCfg, err := cre.PeeringCfgs(bootstrap)
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
			"bootstrapperOCR3Urls": []string{peeringCfg.BootstrapperPeerID + "@" + peeringCfg.BootstrapperHost + ":" + strconv.Itoa(peeringCfg.Port)},
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
