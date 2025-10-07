package dontime

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const flag = cre.DONTimeCapability

type DONTime struct{}

func (o *DONTime) Flag() cre.CapabilityFlag {
	return flag
}

func (o *DONTime) PreDONStartup(
	testLogger zerolog.Logger,
	registryChainSelector uint64,
	cldfEnv *cldf.Environment,
	provider infra.Provider,
	topology *cre.Topology,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	capabilityConfigs cre.CapabilityConfigs,
	contractVersions map[string]string,
) error {
	// nothing to do
	return nil
}

const (
	DONTimeContractQualifier = "capability_dontime"
)

func (o *DONTime) PostDONStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	nodeSetOutput []*cre.WrappedNodeOutput,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	contractVersions map[string]string,
) (*cre.PostDONStartupOutput, error) {
	// should we support more than one DON with DON Time capability?
	donTimeDON, oneErr := creEnv.DonTopology.OneWithFlag(flag)
	if oneErr != nil {
		return nil, oneErr
	}

	_, donTimeContractAddr, timeErr := contracts.DeployOCR3Contract(testLogger, DONTimeContractQualifier, creEnv.DonTopology.HomeChainSelector, creEnv.CldfEnvironment, contractVersions)
	if timeErr != nil {
		return nil, fmt.Errorf("failed to deploy DONTime contract %w", timeErr)
	}

	chainID, chErr := chainselectors.ChainIdFromSelector(creEnv.DonTopology.HomeChainSelector)
	if chErr != nil {
		return nil, errors.Wrapf(chErr, "failed to get chain ID from chain selector %d", creEnv.DonTopology.HomeChainSelector)
	}

	bootstrap, isBootstrap := creEnv.DonTopology.Bootstrap()
	if !isBootstrap {
		return nil, errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	jobErr := createJobs(
		ctx,
		chainID,
		donTimeContractAddr,
		creEnv.CldfEnvironment.Offchain.(*jd.JobDistributor),
		donTimeDON,
		bootstrap,
	)
	if jobErr != nil {
		return nil, fmt.Errorf("failed to create DON Time jobs: %w", jobErr)
	}

	ocr3Config, ocr3confErr := contracts.DefaultOCR3Config()
	if ocr3confErr != nil {
		return nil, fmt.Errorf("failed to get default OCR3 config: %w", ocr3confErr)
	}

	_, donTimeErr := operations.ExecuteOperation(
		creEnv.CldfEnvironment.OperationsBundle,
		ks_contracts_op.ConfigureOCR3Op,
		ks_contracts_op.ConfigureOCR3OpDeps{
			Env: creEnv.CldfEnvironment,
		},
		ks_contracts_op.ConfigureOCR3OpInput{
			ContractAddress: donTimeContractAddr,
			ChainSelector:   creEnv.DonTopology.HomeChainSelector,
			DON:             donTimeDON.KeystoneDONConfig(),
			Config:          donTimeDON.ResolveORC3Config(ocr3Config),
			DryRun:          false,
		},
	)
	if donTimeErr != nil {
		return nil, errors.Wrap(donTimeErr, "failed to configure DON Time contract")
	}

	return &cre.PostDONStartupOutput{}, nil
}

func createJobs(
	ctx context.Context,
	chainID uint64,
	donTimeAddress *common.Address,
	jdClient *jd.JobDistributor,
	don *cre.DON,
	bootstrap *cre.Node,
) error {
	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	_, ocrPeeringCfg, err := cre.PeeringCfgs(bootstrap)
	if err != nil {
		return errors.Wrap(err, "failed to get peering configs")
	}

	jobSpecs := []*jobv1.ProposeJobRequest{}
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
		jobSpecs = append(jobSpecs, donTimeJobSpec(workerNode.JobDistributorDetails.NodeID, donTimeAddress.Hex(), evmKey.PublicAddress.Hex(), evmOCR2KeyBundle, ocrPeeringCfg, chainID))
	}

	return jobs.CreateForDON(ctx, jdClient, don, jobSpecs)
}

func donTimeJobSpec(nodeID string, ocr3CapabilityAddress, nodeEthAddress, ocr2KeyBundleID string, ocrPeeringData cre.OCRPeeringData, chainID uint64) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()
	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "offchainreporting2"
	schemaVersion = 1
	externalJobID = "%s"
	name = "dontime"
	forwardingAllowed = false
	maxTaskDuration = "0s"
	contractID = "%s"
	relay = "evm"
	pluginType = "dontime"
	ocrKeyBundleID = "%s"
	p2pv2Bootstrappers = [
		"%s@%s",
	]
	transmitterID = "%s"

	[relayConfig]
	chainID = "%d"
	providerType = "dontime"

	[pluginConfig]
	pluginName = "dontime"
	ocrVersion = 3
	telemetryType = "plugin"

	[onchainSigningStrategy]
	strategyName = 'multi-chain'
	[onchainSigningStrategy.config]
	evm = "%s"
`,
			uuid,
			ocr3CapabilityAddress, // re-use OCR3Capability contract
			ocr2KeyBundleID,
			ocrPeeringData.OCRBootstraperPeerID,
			fmt.Sprintf("%s:%d", ocrPeeringData.OCRBootstraperHost, ocrPeeringData.Port),
			nodeEthAddress, // transmitterID (although this shouldn't be used for this plugin?)
			chainID,
			ocr2KeyBundleID,
		),
	}
}
