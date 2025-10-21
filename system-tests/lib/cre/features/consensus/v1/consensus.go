package v1

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr"
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

	chainID, chErr := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if chErr != nil {
		return errors.Wrapf(chErr, "failed to get chain ID from chain selector %d", creEnv.RegistryChainSelector)
	}

	jobErr := createJobs(
		ctx,
		chainID,
		ocr3ContractAddr,
		creEnv.CldfEnvironment.Offchain.(*jd.JobDistributor),
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
	chainID uint64,
	ocr3ContractAddr *common.Address,
	jdClient *jd.JobDistributor,
	consensusDON *cre.Don,
	dons *cre.Dons,
) error {
	bootstrap, isBootstrap := dons.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	workerNodes, wErr := consensusDON.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	_, ocrPeeringCfg, err := cre.PeeringCfgs(bootstrap)
	if err != nil {
		return errors.Wrap(err, "failed to get peering configs")
	}

	jobSpecs := []*jobv1.ProposeJobRequest{}
	jobSpecs = append(jobSpecs, ocr.BootstrapOCR3(bootstrap.JobDistributorDetails.NodeID, "ocr3-capability", ocr3ContractAddr.Hex(), chainID))

	for _, workerNode := range workerNodes {
		evmKey, ok := workerNode.Keys.EVM[chainID]
		if !ok {
			return fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
		}

		// we need the OCR2 key bundle for the EVM chain, because OCR jobs currently run only on EVM chains
		evmOCR2KeyBundle, ok := workerNode.Keys.OCR2BundleIDs[chainselectors.FamilyEVM]
		if !ok {
			return fmt.Errorf("node %s does not have OCR2 key bundle for evm", workerNode.Name)
		}

		// we pass here bundles for all chains to enable multi-chain signing
		jobSpecs = append(jobSpecs, WorkerJobSpec(workerNode.JobDistributorDetails.NodeID, ocr3ContractAddr.Hex(), evmKey.PublicAddress.Hex(), evmOCR2KeyBundle, workerNode.Keys.OCR2BundleIDs, ocrPeeringCfg, chainID))
	}

	// pass whole topology, since some jobs might need to be created on multiple DONs
	return jobs.Create(ctx, jdClient, dons, jobSpecs)
}

func WorkerJobSpec(nodeID string, ocr3CapabilityAddress, nodeEthAddress, offchainBundleID string, ocr2KeyBundles map[string]string, ocrPeeringData cre.OCRPeeringData, chainID uint64) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()

	spec := fmt.Sprintf(`
	type = "offchainreporting2"
	schemaVersion = 1
	externalJobID = "%s"
	name = "%s"
	contractID = "%s"
	ocrKeyBundleID = "%s"
	p2pv2Bootstrappers = [
		"%s@%s",
	]
	relay = "evm"
	pluginType = "plugin"
	transmitterID = "%s"
	[relayConfig]
	chainID = "%d"
	[pluginConfig]
	command = "/usr/local/bin/chainlink-ocr3-capability"
	ocrVersion = 3
	pluginName = "ocr-capability"
	providerType = "ocr3-capability"
	telemetryType = "plugin"
	[onchainSigningStrategy]
	strategyName = "multi-chain"
	[onchainSigningStrategy.config]
`,
		uuid,
		cre.ConsensusCapability,
		ocr3CapabilityAddress,
		offchainBundleID,
		ocrPeeringData.OCRBootstraperPeerID,
		fmt.Sprintf("%s:%d", ocrPeeringData.OCRBootstraperHost, ocrPeeringData.Port),
		nodeEthAddress,
		chainID,
	)
	for family, key := range ocr2KeyBundles {
		spec += fmt.Sprintf(`
        %s = "%s"`, family, key)
		spec += "\n"
	}

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec:   spec,
	}
}
