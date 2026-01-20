package sharding

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	deployment_contracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	ring_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	ocr3_changeset "github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset"
	ocr3_contracts "github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	shard_config_changeset "github.com/smartcontractkit/chainlink/deployment/cre/shard_config/v1/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/consensus"
)

const (
	// RingContractQualifier is the qualifier used for the Ring OCR3 contract
	RingContractQualifier = "ring"
)

type SetupShardingInput struct {
	Logger   zerolog.Logger
	CreEnv   *cre.Environment
	Topology *cre.Topology
	Dons     *cre.Dons
}

func SetupSharding(ctx context.Context, input SetupShardingInput) error {
	// Get the shard leader DON
	shardLeaderDON, err := getShardLeaderDON(input.Dons)
	if err != nil {
		return fmt.Errorf("failed to get shard leader DON: %w", err)
	}

	input.Logger.Info().Msgf("Setting up Ring for shard leader DON '%s'", shardLeaderDON.Name)

	// 1. Deploy ShardConfig contract
	shardConfigAddr, err := deployShardConfigContract(input.CreEnv, input.Logger, input.Dons)
	if err != nil {
		return fmt.Errorf("failed to deploy ShardConfig contract: %w", err)
	}

	// 2. Deploy Ring OCR3 contract
	ringOCR3Addr, err := deployRingOCR3Contract(input.CreEnv, input.Logger)
	if err != nil {
		return fmt.Errorf("failed to deploy Ring OCR3 contract: %w", err)
	}

	// 3. Get bootstrap URLs for Ring P2P
	bootstrapURLs, err := getBootstrapURLs(input.Dons)
	if err != nil {
		return fmt.Errorf("failed to get bootstrap URLs: %w", err)
	}

	// 4. Create Ring jobs on the shard leader DON
	err = createRingJobs(ctx, input.CreEnv, shardLeaderDON, input.Dons, ringOCR3Addr, shardConfigAddr, bootstrapURLs)
	if err != nil {
		return fmt.Errorf("failed to create Ring jobs: %w", err)
	}

	time.Sleep(60 * time.Second)
	// 5. Wait for LogPoller to be healthy before configuring OCR3
	if lpErr := consensus.WaitForLogPollerToBeHealthy(shardLeaderDON); lpErr != nil {
		return errors.Wrap(lpErr, "failed while waiting for Log Poller to become healthy")
	}

	// 6. Configure OCR3 contract
	err = configureRingOCR3(input.CreEnv, ringOCR3Addr, shardLeaderDON, input.Logger)
	if err != nil {
		return fmt.Errorf("failed to configure Ring OCR3: %w", err)
	}

	input.Logger.Info().Msgf("Ring setup completed for shard leader DON '%s'", shardLeaderDON.Name)
	return nil
}

// getShardLeaderDON finds the shard leader DON (ShardIndex == 0)
func getShardLeaderDON(dons *cre.Dons) (*cre.Don, error) {
	shardDONs := dons.DonsWithFlag(cre.ShardDON)
	for _, don := range shardDONs {
		if don.Metadata().IsShardLeader() {
			return don, nil
		}
	}
	return nil, errors.New("no shard leader DON found")
}

// deployShardConfigContract deploys the ShardConfig contract
func deployShardConfigContract(creEnv *cre.Environment, logger zerolog.Logger, dons *cre.Dons) (common.Address, error) {
	// Count shard DONs for initial shard count
	shardDONs := dons.DonsWithFlag(cre.ShardDON)
	initialShardCount := uint64(len(shardDONs))

	deployInput := shard_config_changeset.DeployShardConfigInput{
		ChainSelector:     creEnv.RegistryChainSelector,
		InitialShardCount: initialShardCount,
	}

	vErr := shard_config_changeset.DeployShardConfig{}.VerifyPreconditions(*creEnv.CldfEnvironment, deployInput)
	if vErr != nil {
		return common.Address{}, fmt.Errorf("preconditions verification for Shard Config contract failed: %w", vErr)
	}

	out, dErr := shard_config_changeset.DeployShardConfig{}.Apply(*creEnv.CldfEnvironment, deployInput)
	if dErr != nil {
		return common.Address{}, fmt.Errorf("failed to deploy Shard Config contract: %w", dErr)
	}

	crecontracts.MergeAllDataStores(creEnv, out)

	shardConfigAddrStr := crecontracts.MustGetAddressFromDataStore(creEnv.CldfEnvironment.DataStore, creEnv.RegistryChainSelector, deployment_contracts.ShardConfig.String(), semver.MustParse("1"), "")
	shardConfigAddr := common.HexToAddress(shardConfigAddrStr)
	logger.Info().Msgf("Deployed Shard Config v1 contract on chain %d at %s", creEnv.RegistryChainSelector, shardConfigAddr.Hex())

	return shardConfigAddr, nil
}

// deployRingOCR3Contract deploys the OCR3 capability contract for Ring
func deployRingOCR3Contract(creEnv *cre.Environment, logger zerolog.Logger) (common.Address, error) {
	deployInput := ocr3_changeset.DeployOCR3Input{
		ChainSelector: creEnv.RegistryChainSelector,
		Qualifier:     RingContractQualifier,
		Labels:        []string{"ring", "sharding"},
	}

	vErr := ocr3_changeset.DeployOCR3{}.VerifyPreconditions(*creEnv.CldfEnvironment, deployInput)
	if vErr != nil {
		return common.Address{}, fmt.Errorf("preconditions verification for Ring OCR3 contract failed: %w", vErr)
	}

	out, dErr := ocr3_changeset.DeployOCR3{}.Apply(*creEnv.CldfEnvironment, deployInput)
	if dErr != nil {
		return common.Address{}, fmt.Errorf("failed to deploy Ring OCR3 contract: %w", dErr)
	}

	crecontracts.MergeAllDataStores(creEnv, out)

	// Get the deployed contract address
	refKey := pkg.GetOCR3CapabilityAddressRefKey(creEnv.RegistryChainSelector, RingContractQualifier)
	addrRef, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(refKey)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get Ring OCR3 contract address: %w", err)
	}

	addr := common.HexToAddress(addrRef.Address)
	logger.Info().Msgf("Deployed Ring OCR3 contract on chain %d at %s", creEnv.RegistryChainSelector, addr.Hex())

	return addr, nil
}

// getBootstrapURLs extracts P2P bootstrap URLs from the topology's bootstrap node
func getBootstrapURLs(dons *cre.Dons) ([]string, error) {
	bootstrap, ok := dons.Bootstrap()
	if !ok {
		return nil, errors.New("no bootstrap node found in dons")
	}

	_, ocrPeeringCfg, err := cre.PeeringCfgs(bootstrap)
	if err != nil {
		return nil, fmt.Errorf("failed to get peering configs: %w", err)
	}

	bootstrapURL := ocrPeeringCfg.OCRBootstraperPeerID + "@" + ocrPeeringCfg.OCRBootstraperHost + ":" + strconv.Itoa(ocrPeeringCfg.Port)
	return []string{bootstrapURL}, nil
}

// createRingJobs creates Ring jobs on the shard leader DON
func createRingJobs(ctx context.Context, creEnv *cre.Environment, shardLeaderDON *cre.Don, dons *cre.Dons, ringOCR3Addr, shardConfigAddr common.Address, bootstrapURLs []string) error {
	ringJobInput := ring_ops.ProposeRingJobInput{
		Domain:           offchain.ProductLabel,
		EnvName:          cre.EnvironmentName,
		DONName:          shardLeaderDON.Name,
		JobName:          "ring-capability",
		ContractAddress:  ringOCR3Addr.Hex(),
		ChainSelectorEVM: creEnv.RegistryChainSelector,
		ShardConfigAddr:  shardConfigAddr.Hex(),
		BootstrapperUrls: bootstrapURLs,
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: shardLeaderDON.Name},
		},
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: "ring"},
	}

	report, err := operations.ExecuteSequence(
		creEnv.CldfEnvironment.OperationsBundle,
		ring_ops.ProposeRingJob,
		ring_ops.ProposeRingJobDeps{Env: *creEnv.CldfEnvironment},
		ringJobInput,
	)
	if err != nil {
		return fmt.Errorf("failed to propose Ring jobs: %w", err)
	}

	// Approve the proposed jobs
	if err := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, report.Output.Specs); err != nil {
		return fmt.Errorf("failed to approve Ring jobs: %w", err)
	}

	return nil
}

// configureRingOCR3 configures the Ring OCR3 contract with DON signers
func configureRingOCR3(creEnv *cre.Environment, ringOCR3Addr common.Address, shardLeaderDON *cre.Don, logger zerolog.Logger) error {
	// Get worker node IDs for the shard leader DON
	nodeIDs := shardLeaderDON.JDNodeIDs()
	if len(nodeIDs) == 0 {
		return fmt.Errorf("no nodes found in shard leader DON '%s'", shardLeaderDON.Name)
	}

	// Create default OCR3 configuration for Ring
	oracleConfig := &ocr3.OracleConfig{
		UniqueReports:                     false,
		DeltaProgressMillis:               30000,
		DeltaResendMillis:                 10000,
		DeltaInitialMillis:                20000,
		DeltaRoundMillis:                  2000,
		DeltaGraceMillis:                  500,
		DeltaCertifiedCommitRequestMillis: 500,
		DeltaStageMillis:                  60000,
		MaxRoundsPerEpoch:                 100,
		TransmissionSchedule:              []int{len(nodeIDs)},
		MaxDurationQueryMillis:            5000,
		MaxDurationObservationMillis:      10000,
		MaxDurationShouldAcceptMillis:     5000,
		MaxDurationShouldTransmitMillis:   5000,
		MaxFaultyOracles:                  int(shardLeaderDON.F),
	}

	configInput := ocr3_changeset.ConfigureOCR3Input{
		ContractChainSelector: creEnv.RegistryChainSelector,
		ContractQualifier:     RingContractQualifier,
		DON: ocr3_contracts.DonNodeSet{
			Name:    shardLeaderDON.Name,
			NodeIDs: nodeIDs,
		},
		OracleConfig: oracleConfig,
		DryRun:       false,
	}

	vErr := ocr3_changeset.ConfigureOCR3{}.VerifyPreconditions(*creEnv.CldfEnvironment, configInput)
	if vErr != nil {
		return fmt.Errorf("preconditions verification for Ring OCR3 configuration failed: %w", vErr)
	}

	_, cErr := ocr3_changeset.ConfigureOCR3{}.Apply(*creEnv.CldfEnvironment, configInput)
	if cErr != nil {
		return fmt.Errorf("failed to configure Ring OCR3 contract: %w", cErr)
	}

	logger.Info().Msgf("Configured Ring OCR3 contract at %s with DON '%s' (%d nodes)", ringOCR3Addr.Hex(), shardLeaderDON.Name, len(nodeIDs))

	return nil
}
