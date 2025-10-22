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

	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
)

const flag = cre.DONTimeCapability

type DONTime struct{}

func (o *DONTime) Flag() cre.CapabilityFlag {
	return flag
}

func (o *DONTime) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	// nothing to do
	return nil, nil
}

const (
	ContractQualifier = "capability_dontime"
)

func (o *DONTime) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	_, donTimeContractAddr, timeErr := contracts.DeployOCR3Contract(testLogger, ContractQualifier, creEnv.RegistryChainSelector, creEnv.CldfEnvironment, creEnv.ContractVersions)
	if timeErr != nil {
		return fmt.Errorf("failed to deploy DONTime contract %w", timeErr)
	}

	chainID, chErr := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if chErr != nil {
		return errors.Wrapf(chErr, "failed to get chain ID from chain selector %d", creEnv.RegistryChainSelector)
	}
	jobErr := createJobs(
		ctx,
		chainID,
		donTimeContractAddr,
		creEnv.CldfEnvironment.Offchain.(*jd.JobDistributor),
		don,
		dons,
	)
	if jobErr != nil {
		return fmt.Errorf("failed to create DON Time jobs: %w", jobErr)
	}

	ocr3Config, ocr3confErr := contracts.DefaultOCR3Config()
	if ocr3confErr != nil {
		return fmt.Errorf("failed to get default OCR3 config: %w", ocr3confErr)
	}

	_, donTimeErr := operations.ExecuteOperation(
		creEnv.CldfEnvironment.OperationsBundle,
		ks_contracts_op.ConfigureOCR3Op,
		ks_contracts_op.ConfigureOCR3OpDeps{
			Env: creEnv.CldfEnvironment,
		},
		ks_contracts_op.ConfigureOCR3OpInput{
			ContractAddress: donTimeContractAddr,
			ChainSelector:   creEnv.RegistryChainSelector,
			DON:             don.KeystoneDONConfig(),
			Config:          don.ResolveORC3Config(ocr3Config),
			DryRun:          false,
		},
	)
	if donTimeErr != nil {
		return errors.Wrap(donTimeErr, "failed to configure DON Time contract")
	}

	return nil
}

func createJobs(
	ctx context.Context,
	chainID uint64,
	donTimeAddress *common.Address,
	jdClient *jd.JobDistributor,
	donTimeDON *cre.Don,
	dons *cre.Dons,
) error {
	bootstrap, isBootstrap := dons.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	workerNodes, wErr := donTimeDON.Workers()
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
		jobSpecs = append(jobSpecs, WorkerJobSpec(workerNode.JobDistributorDetails.NodeID, donTimeAddress.Hex(), evmKey.PublicAddress.Hex(), evmOCR2KeyBundle, ocrPeeringCfg, chainID))
	}

	// pass whole topology, since some jobs might need to be created on multiple DONs
	return jobs.Create(ctx, jdClient, dons, jobSpecs)
}

func WorkerJobSpec(nodeID string, ocr3CapabilityAddress, nodeEthAddress, ocr2KeyBundleID string, ocrPeeringData cre.OCRPeeringData, chainID uint64) *jobv1.ProposeJobRequest {
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
