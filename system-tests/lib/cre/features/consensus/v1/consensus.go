package v1

import (
	"context"
	"fmt"
	"strconv"

	"dario.cat/mergo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/consensus"
)

const flag = cre.ConsensusCapability

type Consensus struct{}

func (c *Consensus) Flag() cre.CapabilityFlag {
	return flag
}

func (c *Consensus) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "offchain_reporting",
			Version:        "1.0.0",
			CapabilityType: 2, // CONSENSUS
			ResponseType:   0, // REPORT
		},
		Config: &capabilitiespb.CapabilityConfig{},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

const (
	ContractQualifier = "capability_ocr3"
)

func (c *Consensus) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	// should we support more than one DON with OCR3 capability? Could there be 0? I guess as long as there's 1 with consensus v2?
	_, ocr3ContractAddr, ocrErr := contracts.DeployOCR3Contract(testLogger, ContractQualifier, creEnv.RegistryChainSelector, creEnv.CldfEnvironment, creEnv.ContractVersions)
	if ocrErr != nil {
		return fmt.Errorf("failed to deploy OCR3 contract %w", ocrErr)
	}

	jobErr := createJobs(
		ctx,
		creEnv,
		don,
		dons,
	)
	if jobErr != nil {
		return fmt.Errorf("failed to create OCR3 jobs: %w", jobErr)
	}

	// wait for LP to be started (otherwise it won't pick up contract's configuration events)
	if err := consensus.WaitForLogPollerToBeHealthy(don); err != nil {
		return errors.Wrap(err, "failed while waiting for Log Poller to become healthy")
	}

	ocr3Config, ocr3confErr := contracts.DefaultOCR3Config()
	if ocr3confErr != nil {
		return fmt.Errorf("failed to get default OCR3 config: %w", ocr3confErr)
	}

	_, ocr3Err := operations.ExecuteOperation(
		creEnv.CldfEnvironment.OperationsBundle,
		ks_contracts_op.ConfigureOCR3Op,
		ks_contracts_op.ConfigureOCR3OpDeps{
			Env: creEnv.CldfEnvironment,
		},
		ks_contracts_op.ConfigureOCR3OpInput{
			ContractAddress: ocr3ContractAddr,
			ChainSelector:   creEnv.RegistryChainSelector,
			DON:             don.KeystoneDONConfig(),
			Config:          don.ResolveORC3Config(ocr3Config),
			DryRun:          false,
		},
	)

	if ocr3Err != nil {
		return errors.Wrap(ocr3Err, "failed to configure OCR3 contract")
	}

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

	bootInput := cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: cre.EnvironmentName,
		DONName:     bootstrap.DON.Name,
		JobName:     "consensus-v1-bootstrap",
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: bootstrap.DON.Name},
		},
		Template: job_types.BootstrapOCR3,
		Inputs: job_types.JobSpecInput{
			"chainSelector":     creEnv.RegistryChainSelector,
			"contractQualifier": ContractQualifier,
		},
	}

	bootVerErr := cre_jobs.ProposeJobSpec{}.VerifyPreconditions(*creEnv.CldfEnvironment, bootInput)
	if bootVerErr != nil {
		return fmt.Errorf("precondition verification failed for Consensus v1 bootstrap job: %w", bootVerErr)
	}

	bootReport, bootErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, bootInput)
	if bootErr != nil {
		return fmt.Errorf("failed to propose Consensus v1 bootstrap job spec: %w", bootErr)
	}

	specs := make(map[string][]string)
	for _, r := range bootReport.Reports {
		out, ok := r.Output.(cre_jobs_ops.ProposeOCR3BootstrapJobOutput)
		if !ok {
			return fmt.Errorf("unable to cast to ProposeOCR3BootstrapJobOutput, actual type: %T", r.Output)
		}
		mErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice)
		if mErr != nil {
			return fmt.Errorf("failed to merge bootstrap job specs: %w", mErr)
		}
	}

	_, ocrPeeringCfg, err := cre.PeeringCfgs(bootstrap)
	if err != nil {
		return errors.Wrap(err, "failed to get peering configs")
	}

	workerInput := cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: cre.EnvironmentName,
		DONName:     don.Name,
		JobName:     "consensus-v1-worker",
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: don.Name},
		},
		Template: job_types.OCR3,
		Inputs: job_types.JobSpecInput{
			"chainSelectorEVM":     creEnv.RegistryChainSelector,
			"contractQualifier":    ContractQualifier,
			"templateName":         "worker-ocr3",
			"bootstrapperOCR3Urls": []string{ocrPeeringCfg.OCRBootstraperPeerID + "@" + ocrPeeringCfg.OCRBootstraperHost + ":" + strconv.Itoa(ocrPeeringCfg.Port)},
		},
	}

	for _, blockchain := range creEnv.Blockchains {
		if blockchain.IsFamily(chainselectors.FamilySolana) {
			workerInput.Inputs["chainSelectorSolana"] = blockchain.ChainSelector()
			break
		}
	}

	workerVerErr := cre_jobs.ProposeJobSpec{}.VerifyPreconditions(*creEnv.CldfEnvironment, workerInput)
	if workerVerErr != nil {
		return fmt.Errorf("precondition verification failed for Consensus v1 worker job: %w", workerVerErr)
	}

	workerReport, workerErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, workerInput)
	if workerErr != nil {
		return fmt.Errorf("failed to propose Consensus v1 worker job spec: %w", workerErr)
	}

	for _, r := range workerReport.Reports {
		out, ok := r.Output.(cre_jobs_ops.ProposeOCR3JobOutput)
		if !ok {
			return fmt.Errorf("unable to cast to ProposeOCR3JobOutput, actual type: %T", r.Output)
		}
		mErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice)
		if mErr != nil {
			return fmt.Errorf("failed to merge worker job specs: %w", mErr)
		}
	}

	approveErr := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs)
	if approveErr != nil {
		return fmt.Errorf("failed to approve Consensus v1 jobs: %w", approveErr)
	}

	return nil
}
