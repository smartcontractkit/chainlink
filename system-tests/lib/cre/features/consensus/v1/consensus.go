package v1

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/consensus"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const flag = cre.ConsensusCapability

type Consensus struct{}

func (o *Consensus) Flag() cre.CapabilityFlag {
	return flag
}

func (o *Consensus) PreEnvStartup(
	testLogger zerolog.Logger,
	registryChainSelector uint64,
	cldfEnv *cldf.Environment,
	provider infra.Provider,
	topology *cre.Topology,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	capabilityConfigs cre.CapabilityConfigs,
	contractVersions map[string]string,
	gatewayConfigs map[cre.NodeUUID]*config.GatewayConfig,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := make(map[uint64][]keystone_changeset.DONCapabilityWithConfig)
	for _, donMetadata := range topology.DonsMetadataWithFlag(flag) {
		if capabilities[donMetadata.ID] == nil {
			capabilities[donMetadata.ID] = []keystone_changeset.DONCapabilityWithConfig{}
		}
		capabilities[donMetadata.ID] = append(capabilities[donMetadata.ID], keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "offchain_reporting",
				Version:        "1.0.0",
				CapabilityType: 2, // CONSENSUS
				ResponseType:   0, // REPORT
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfigs: capabilities,
	}, nil
}

const (
	OCR3ContractQualifier = "capability_ocr3"
)

func (o *Consensus) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	nodeSetOutput []*cre.WrappedNodeOutput,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	contractVersions map[string]string,
	provider infra.Provider,
	capabilityConfigs map[string]cre.CapabilityConfig,
) error {
	// should we support more than one DON with OCR3 capability? Could there be 0? I guess as long as there's 1 with consensus v2?
	dons := creEnv.DonTopology.DonsWithFlag(flag)
	if len(dons) == 0 {
		return nil
	}
	if len(dons) > 1 {
		return fmt.Errorf("more than one DON with consensus v1 capability is not supported yet")
	}
	consensusDON := dons[0]

	_, ocr3ContractAddr, ocrErr := contracts.DeployOCR3Contract(testLogger, OCR3ContractQualifier, creEnv.DonTopology.HomeChainSelector, creEnv.CldfEnvironment, contractVersions)
	if ocrErr != nil {
		return fmt.Errorf("failed to deploy OCR3 contract %w", ocrErr)
	}

	chainID, chErr := chainselectors.ChainIdFromSelector(creEnv.DonTopology.HomeChainSelector)
	if chErr != nil {
		return errors.Wrapf(chErr, "failed to get chain ID from chain selector %d", creEnv.DonTopology.HomeChainSelector)
	}

	bootstrap, isBootstrap := creEnv.DonTopology.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	jobErr := createJobs(
		ctx,
		chainID,
		ocr3ContractAddr,
		creEnv.CldfEnvironment.Offchain.(*jd.JobDistributor),
		consensusDON,
		creEnv.DonTopology,
		bootstrap,
	)
	if jobErr != nil {
		return fmt.Errorf("failed to create OCR3 jobs: %w", jobErr)
	}

	// wait for LP to be started (otherwise it won't pick up contract's configuration events)
	if err := consensus.WaitForLogPollerToBeHealthy(consensusDON); err != nil {
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
			ChainSelector:   creEnv.DonTopology.HomeChainSelector,
			DON:             consensusDON.KeystoneDONConfig(),
			Config:          consensusDON.ResolveORC3Config(ocr3Config),
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
	consensusDON *cre.DON,
	donTopology *cre.DonTopology,
	bootstrap *cre.Node,
) error {
	workerNodes, wErr := consensusDON.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	_, ocrPeeringCfg, err := cre.PeeringCfgs(bootstrap)
	if err != nil {
		return errors.Wrap(err, "failed to get peering configs")
	}

	jobSpecs := []*job.ProposeJobRequest{}
	jobSpecs = append(jobSpecs, jobs.BootstrapOCR3(bootstrap.JobDistributorDetails.NodeID, "ocr3-capability", ocr3ContractAddr.Hex(), chainID))

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
		jobSpecs = append(jobSpecs, jobs.WorkerOCR3(workerNode.JobDistributorDetails.NodeID, ocr3ContractAddr.Hex(), evmKey.PublicAddress.Hex(), evmOCR2KeyBundle, workerNode.Keys.OCR2BundleIDs, ocrPeeringCfg, chainID))
	}

	// pass whole topology, since some jobs might need to be created on multiple DONs
	return jobs.Create(ctx, jdClient, donTopology, jobSpecs)
}
