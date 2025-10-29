package v2

import (
	"context"
	"fmt"
	"strings"

	"dario.cat/mergo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/consensus"
)

const flag = cre.ConsensusCapabilityV2

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
			LabelledName:   "consensus",
			Version:        "1.0.0-alpha",
			CapabilityType: 2, // CONSENSUS
			ResponseType:   0, // REPORT
		},
		Config: &capabilitiespb.CapabilityConfig{},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

const ContractQualifier = "capability_consensus"

func (c *Consensus) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	_, ocr3ContractAddr, ocrErr := contracts.DeployOCR3Contract(testLogger, ContractQualifier, creEnv.RegistryChainSelector, creEnv.CldfEnvironment, creEnv.ContractVersions)
	if ocrErr != nil {
		return fmt.Errorf("failed to deploy OCR3 (consensus v2) contract %w", ocrErr)
	}

	jobsErr := createJobs(
		ctx,
		don,
		dons,
		creEnv,
	)
	if jobsErr != nil {
		return fmt.Errorf("failed to create OCR3 jobs: %w", jobsErr)
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
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	specs := make(map[string][]string)

	capabilityConfig, ok := creEnv.CapabilityConfigs[flag]
	if !ok {
		return fmt.Errorf("%s config not found in capabilities config: %v", flag, creEnv.CapabilityConfigs)
	}

	command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryPath, creEnv.Provider)
	if cErr != nil {
		return errors.Wrap(cErr, "failed to get command for cron capability")
	}

	var nodeSet cre.NodeSetWithCapabilityConfigs
	for _, ns := range dons.AsNodeSetWithChainCapabilities() {
		if ns.GetName() == don.Name {
			nodeSet = ns
			break
		}
	}
	if nodeSet == nil {
		return fmt.Errorf("could not find node set for Don named '%s'", don.Name)
	}

	bootstrap, isBootstrap := dons.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	chainID, cErr := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if cErr != nil {
		return fmt.Errorf("failed to get chain ID from selector %d: %w", creEnv.RegistryChainSelector, cErr)
	}

	bootInput := cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: cre.EnvironmentName,
		DONName:     don.Name,
		JobName:     "consensus-v2-bootstrap",
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: don.Name},
		},
		Template: job_types.BootstrapOCR3,
		Inputs: job_types.JobSpecInput{
			"chainSelector":     creEnv.RegistryChainSelector,
			"contractQualifier": ContractQualifier,
		},
	}

	bootVerErr := cre_jobs.ProposeJobSpec{}.VerifyPreconditions(*creEnv.CldfEnvironment, bootInput)
	if bootVerErr != nil {
		return fmt.Errorf("precondition verification failed for Consensus v2 bootstrap job for chainID %d: %w", chainID, bootVerErr)
	}

	bootReport, bootErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, bootInput)
	if bootErr != nil {
		return fmt.Errorf("failed to propose Consensus v2 bootstrap job spec for chainID %d: %w", chainID, bootErr)
	}

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

	bootstrapPeers := []string{fmt.Sprintf("%s@%s:%d", strings.TrimPrefix(bootstrap.Keys.PeerID(), "p2p_"), bootstrap.Host, cre.OCRPeeringPort)}

	workerInput := cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: cre.EnvironmentName,
		DONName:     don.Name,
		JobName:     "consensus-v2-worker",
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: don.Name},
		},
		Template: job_types.Consensus,
		Inputs: job_types.JobSpecInput{
			"command":           command,
			"chainSelectorEVM":  creEnv.RegistryChainSelector,
			"contractQualifier": ContractQualifier,
			"bootstrapPeers":    bootstrapPeers,
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
		return fmt.Errorf("precondition verification failed for Consensus v2 worker job: %w", workerVerErr)
	}

	workerReport, workerErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, workerInput)
	if workerErr != nil {
		return fmt.Errorf("failed to propose Consensus v2 worker job spec: %w", workerErr)
	}

	for _, r := range workerReport.Reports {
		out, ok := r.Output.(cre_jobs_ops.ProposeStandardCapabilityJobOutput)
		if !ok {
			return fmt.Errorf("unable to cast to ProposeStandardCapabilityJobOutput, actual type: %T", r.Output)
		}
		mErr := mergo.Merge(&specs, out.Specs)
		if mErr != nil {
			return fmt.Errorf("failed to merge worker job specs: %w", mErr)
		}
	}

	approveErr := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs)
	if approveErr != nil {
		return fmt.Errorf("failed to approve Consensus v2 jobs: %w", approveErr)
	}

	return nil
}
