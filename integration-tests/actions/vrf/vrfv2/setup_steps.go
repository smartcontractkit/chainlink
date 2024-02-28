package vrfv2

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	testconfig "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"

	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/types"
)

func CreateVRFV2Job(
	chainlinkNode *client.ChainlinkClient,
	vrfJobSpecConfig vrfcommon.VRFJobSpecConfig,
) (*client.Job, error) {
	jobUUID := uuid.New()
	os := &client.VRFV2TxPipelineSpec{
		Address:               vrfJobSpecConfig.CoordinatorAddress,
		EstimateGasMultiplier: vrfJobSpecConfig.EstimateGasMultiplier,
		FromAddress:           vrfJobSpecConfig.FromAddresses[0],
		SimulationBlock:       vrfJobSpecConfig.SimulationBlock,
	}
	ost, err := os.String()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrParseJob, err)
	}

	spec := &client.VRFV2JobSpec{
		Name:                          fmt.Sprintf("vrf-v2-%s", jobUUID),
		ForwardingAllowed:             vrfJobSpecConfig.ForwardingAllowed,
		CoordinatorAddress:            vrfJobSpecConfig.CoordinatorAddress,
		FromAddresses:                 vrfJobSpecConfig.FromAddresses,
		EVMChainID:                    vrfJobSpecConfig.EVMChainID,
		MinIncomingConfirmations:      vrfJobSpecConfig.MinIncomingConfirmations,
		PublicKey:                     vrfJobSpecConfig.PublicKey,
		ExternalJobID:                 jobUUID.String(),
		ObservationSource:             ost,
		BatchFulfillmentEnabled:       vrfJobSpecConfig.BatchFulfillmentEnabled,
		BatchFulfillmentGasMultiplier: vrfJobSpecConfig.BatchFulfillmentGasMultiplier,
		PollPeriod:                    vrfJobSpecConfig.PollPeriod,
		RequestTimeout:                vrfJobSpecConfig.RequestTimeout,
	}
	if vrfJobSpecConfig.VRFOwnerConfig.UseVRFOwner {
		spec.VRFOwner = vrfJobSpecConfig.VRFOwnerConfig.OwnerAddress
		spec.UseVRFOwner = true
	}

	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrParseJob, err)

	}
	job, err := chainlinkNode.MustCreateJob(spec)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrCreatingVRFv2Job, err)
	}
	return job, nil
}

// SetupVRFV2Environment will create specified number of subscriptions and add the same conumer/s to each of them
func SetupVRFV2Environment(
	env *test_env.CLClusterTestEnv,
	nodesToCreate []vrfcommon.VRFNodeType,
	vrfv2TestConfig types.VRFv2TestConfig,
	useVRFOwner bool,
	useTestCoordinator bool,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	registerProvingKeyAgainstAddress string,
	numberOfTxKeysToCreate int,
	numberOfConsumers int,
	numberOfSubToCreate int,
	l zerolog.Logger,
) (*vrfcommon.VRFContracts, []uint64, *vrfcommon.VRFKeyData, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) {
	l.Info().Msg("Starting VRFV2 environment setup")
	configGeneral := vrfv2TestConfig.GetVRFv2Config().General
	vrfContracts, subIDs, err := SetupVRFV2Contracts(
		env,
		linkToken,
		mockNativeLINKFeed,
		numberOfConsumers,
		useVRFOwner,
		useTestCoordinator,
		configGeneral,
		numberOfSubToCreate,
		l,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	nodeTypeToNodeMap := vrfcommon.CreateNodeTypeToNodeMap(env.ClCluster, nodesToCreate)
	vrfKey, pubKeyCompressed, err := vrfcommon.CreateVRFKeyOnVRFNode(nodeTypeToNodeMap[vrfcommon.VRF], l)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	l.Info().Str("Coordinator", vrfContracts.CoordinatorV2.Address()).Msg("Registering Proving Key")
	provingKey, err := VRFV2RegisterProvingKey(vrfKey, registerProvingKeyAgainstAddress, vrfContracts.CoordinatorV2)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRegisteringProvingKey, err)
	}
	keyHash, err := vrfContracts.CoordinatorV2.HashOfKey(context.Background(), provingKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrCreatingProvingKeyHash, err)
	}

	chainID := env.EVMClient.GetChainID()
	vrfTXKeyAddressStrings, vrfTXKeyAddresses, err := vrfcommon.CreateFundAndGetSendingKeys(
		env.EVMClient,
		nodeTypeToNodeMap[vrfcommon.VRF],
		*vrfv2TestConfig.GetCommonConfig().ChainlinkNodeFunding,
		numberOfTxKeysToCreate,
		chainID,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	nodeTypeToNodeMap[vrfcommon.VRF].TXKeyAddressStrings = vrfTXKeyAddressStrings

	vrfOwnerConfig, err := SetupVRFOwnerContractIfNeeded(useVRFOwner, env, vrfContracts, vrfTXKeyAddressStrings, vrfTXKeyAddresses, l)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	g := errgroup.Group{}
	if vrfNode, exists := nodeTypeToNodeMap[vrfcommon.VRF]; exists {
		g.Go(func() error {
			err := setupVRFNode(vrfContracts, chainID, configGeneral, pubKeyCompressed, vrfOwnerConfig, l, vrfNode)
			if err != nil {
				return err
			}
			return nil
		})
	}

	if bhsNode, exists := nodeTypeToNodeMap[vrfcommon.BHS]; exists {
		g.Go(func() error {
			err := vrfcommon.SetupBHSNode(
				env,
				configGeneral.General,
				numberOfTxKeysToCreate,
				chainID,
				vrfContracts.CoordinatorV2.Address(),
				vrfContracts.BHS.Address(),
				*vrfv2TestConfig.GetCommonConfig().ChainlinkNodeFunding,
				l,
				bhsNode,
			)
			if err != nil {
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("VRF node setup ended up with an error: %w", err)
	}

	vrfKeyData := vrfcommon.VRFKeyData{
		VRFKey:            vrfKey,
		EncodedProvingKey: provingKey,
		KeyHash:           keyHash,
	}

	l.Info().Msg("VRFV2 environment setup is finished")
	return vrfContracts, subIDs, &vrfKeyData, nodeTypeToNodeMap, nil
}

func setupVRFNode(contracts *vrfcommon.VRFContracts, chainID *big.Int, vrfv2Config *testconfig.General, pubKeyCompressed string, vrfOwnerConfig *vrfcommon.VRFOwnerConfig, l zerolog.Logger, vrfNode *vrfcommon.VRFNode) error {
	vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
		ForwardingAllowed:             *vrfv2Config.VRFJobForwardingAllowed,
		CoordinatorAddress:            contracts.CoordinatorV2.Address(),
		FromAddresses:                 vrfNode.TXKeyAddressStrings,
		EVMChainID:                    chainID.String(),
		MinIncomingConfirmations:      int(*vrfv2Config.MinimumConfirmations),
		PublicKey:                     pubKeyCompressed,
		EstimateGasMultiplier:         *vrfv2Config.VRFJobEstimateGasMultiplier,
		BatchFulfillmentEnabled:       *vrfv2Config.VRFJobBatchFulfillmentEnabled,
		BatchFulfillmentGasMultiplier: *vrfv2Config.VRFJobBatchFulfillmentGasMultiplier,
		PollPeriod:                    vrfv2Config.VRFJobPollPeriod.Duration,
		RequestTimeout:                vrfv2Config.VRFJobRequestTimeout.Duration,
		SimulationBlock:               vrfv2Config.VRFJobSimulationBlock,
		VRFOwnerConfig:                vrfOwnerConfig,
	}

	l.Info().Msg("Creating VRFV2 Job")
	vrfV2job, err := CreateVRFV2Job(
		vrfNode.CLNode.API,
		vrfJobSpecConfig,
	)
	if err != nil {
		return fmt.Errorf("%s, err %w", ErrCreateVRFV2Jobs, err)
	}
	vrfNode.Job = vrfV2job

	// this part is here because VRFv2 can work with only a specific key
	// [[EVM.KeySpecific]]
	//	Key = '...'
	nodeConfig := node.NewConfig(vrfNode.CLNode.NodeConfig,
		node.WithLogPollInterval(1*time.Second),
		node.WithVRFv2EVMEstimator(vrfNode.TXKeyAddressStrings, *vrfv2Config.CLNodeMaxGasPriceGWei),
	)
	l.Info().Msg("Restarting Node with new sending key PriceMax configuration")
	err = vrfNode.CLNode.Restart(nodeConfig)
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrRestartCLNode, err)
	}
	return nil
}

func SetupVRFV2WrapperEnvironment(
	env *test_env.CLClusterTestEnv,
	vrfv2TestConfig tc.VRFv2TestConfig,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	coordinator contracts.VRFCoordinatorV2,
	keyHash [32]byte,
	wrapperConsumerContractsAmount int,
) (*VRFV2WrapperContracts, *uint64, error) {
	// Deploy VRF v2 direct funding contracts
	wrapperContracts, err := DeployVRFV2DirectFundingContracts(
		env.ContractDeployer,
		env.EVMClient,
		linkToken.Address(),
		mockNativeLINKFeed.Address(),
		coordinator,
		wrapperConsumerContractsAmount,
	)
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	vrfv2Config := vrfv2TestConfig.GetVRFv2Config()

	// Configure VRF v2 wrapper contract
	err = wrapperContracts.VRFV2Wrapper.SetConfig(
		*vrfv2Config.General.WrapperGasOverhead,
		*vrfv2Config.General.CoordinatorGasOverhead,
		*vrfv2Config.General.WrapperPremiumPercentage,
		keyHash,
		*vrfv2Config.General.WrapperMaxNumberOfWords,
	)
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	// Fetch wrapper subscription ID
	wrapperSubID, err := wrapperContracts.VRFV2Wrapper.GetSubID(context.Background())
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	// Fund wrapper subscription
	err = FundSubscriptions(env, big.NewFloat(*vrfv2Config.General.SubscriptionFundingAmountLink), linkToken, coordinator, []uint64{wrapperSubID})
	if err != nil {
		return nil, nil, err
	}

	// Fund consumer with LINK
	err = linkToken.Transfer(
		wrapperContracts.LoadTestConsumers[0].Address(),
		big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(*vrfv2Config.General.WrapperConsumerFundingAmountLink)),
	)
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	return wrapperContracts, &wrapperSubID, nil
}
