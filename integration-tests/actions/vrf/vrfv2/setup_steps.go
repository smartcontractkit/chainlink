package vrfv2

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
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
	if vrfJobSpecConfig.BatchFulfillmentEnabled {
		spec.BatchCoordinatorAddress = vrfJobSpecConfig.BatchCoordinatorAddress
	}
	job, err := chainlinkNode.MustCreateJob(spec)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrCreatingVRFv2Job, err)
	}
	return job, nil
}

// SetupVRFV2Environment will create specified number of subscriptions and add the same conumer/s to each of them
func SetupVRFV2Environment(
	ctx context.Context,
	env *test_env.CLClusterTestEnv,
	chainID int64,
	nodesToCreate []vrfcommon.VRFNodeType,
	vrfv2TestConfig types.VRFv2TestConfig,
	useVRFOwner bool,
	useTestCoordinator bool,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.VRFMockETHLINKFeed,
	registerProvingKeyAgainstAddress string,
	numberOfTxKeysToCreate int,
	l zerolog.Logger,
) (*vrfcommon.VRFContracts, *vrfcommon.VRFKeyData, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) {
	l.Info().Msg("Starting VRFV2 environment setup")
	configGeneral := vrfv2TestConfig.GetVRFv2Config().General
	vrfContracts, err := SetupVRFV2Contracts(
		env,
		chainID,
		linkToken,
		mockNativeLINKFeed,
		useVRFOwner,
		useTestCoordinator,
		configGeneral,
		l,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	nodeTypeToNodeMap, err := vrfcommon.CreateNodeTypeToNodeMap(env.ClCluster, nodesToCreate)
	if err != nil {
		return nil, nil, nil, err
	}
	vrfKey, pubKeyCompressed, err := vrfcommon.CreateVRFKeyOnVRFNode(nodeTypeToNodeMap[vrfcommon.VRF], l)
	if err != nil {
		return nil, nil, nil, err
	}

	l.Info().Str("Coordinator", vrfContracts.CoordinatorV2.Address()).Msg("Registering Proving Key")
	provingKey, err := VRFV2RegisterProvingKey(vrfKey, registerProvingKeyAgainstAddress, vrfContracts.CoordinatorV2)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRegisteringProvingKey, err)
	}
	keyHash, err := vrfContracts.CoordinatorV2.HashOfKey(ctx, provingKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrCreatingProvingKeyHash, err)
	}

	sethClient, err := env.GetSethClient(chainID)
	if err != nil {
		return nil, nil, nil, err
	}

	vrfTXKeyAddressStrings, vrfTXKeyAddresses, err := vrfcommon.CreateFundAndGetSendingKeys(
		l,
		sethClient,
		nodeTypeToNodeMap[vrfcommon.VRF],
		*vrfv2TestConfig.GetCommonConfig().ChainlinkNodeFunding,
		numberOfTxKeysToCreate,
		big.NewInt(chainID),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	nodeTypeToNodeMap[vrfcommon.VRF].TXKeyAddressStrings = vrfTXKeyAddressStrings

	vrfOwnerConfig, err := SetupVRFOwnerContractIfNeeded(useVRFOwner, vrfContracts, vrfTXKeyAddressStrings, vrfTXKeyAddresses, l)
	if err != nil {
		return nil, nil, nil, err
	}

	g := errgroup.Group{}
	if vrfNode, exists := nodeTypeToNodeMap[vrfcommon.VRF]; exists {
		g.Go(func() error {
			err := setupVRFNode(vrfContracts, big.NewInt(chainID), configGeneral, pubKeyCompressed, vrfOwnerConfig, l, vrfNode)
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
				big.NewInt(chainID),
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
		return nil, nil, nil, fmt.Errorf("VRF node setup ended up with an error: %w", err)
	}

	vrfKeyData := vrfcommon.VRFKeyData{
		VRFKey:            vrfKey,
		EncodedProvingKey: provingKey,
		KeyHash:           keyHash,
		PubKeyCompressed:  pubKeyCompressed,
	}

	l.Info().Msg("VRFV2 environment setup is finished")
	return vrfContracts, &vrfKeyData, nodeTypeToNodeMap, nil
}

func setupVRFNode(contracts *vrfcommon.VRFContracts, chainID *big.Int, vrfv2Config *testconfig.General, pubKeyCompressed string, vrfOwnerConfig *vrfcommon.VRFOwnerConfig, l zerolog.Logger, vrfNode *vrfcommon.VRFNode) error {
	vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
		ForwardingAllowed:             *vrfv2Config.VRFJobForwardingAllowed,
		CoordinatorAddress:            contracts.CoordinatorV2.Address(),
		BatchCoordinatorAddress:       contracts.BatchCoordinatorV2.Address(),
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
	ctx context.Context,
	env *test_env.CLClusterTestEnv,
	chainID int64,
	vrfv2TestConfig tc.VRFv2TestConfig,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.VRFMockETHLINKFeed,
	coordinator contracts.VRFCoordinatorV2,
	keyHash [32]byte,
	wrapperConsumerContractsAmount int,
) (*VRFV2WrapperContracts, *uint64, error) {
	sethClient, err := env.GetSethClient(chainID)
	if err != nil {
		return nil, nil, err
	}

	// Deploy VRF v2 direct funding contracts
	wrapperContracts, err := DeployVRFV2DirectFundingContracts(
		sethClient,
		linkToken.Address(),
		mockNativeLINKFeed.Address(),
		coordinator,
		wrapperConsumerContractsAmount,
	)
	if err != nil {
		return nil, nil, err
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

	// Fetch wrapper subscription ID
	wrapperSubID, err := wrapperContracts.VRFV2Wrapper.GetSubID(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Fund wrapper subscription
	err = FundSubscriptions(big.NewFloat(*vrfv2Config.General.SubscriptionFundingAmountLink), linkToken, coordinator, []uint64{wrapperSubID})
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

	return wrapperContracts, &wrapperSubID, nil
}

func SetupVRFV2Universe(
	ctx context.Context,
	t *testing.T,
	testConfig tc.TestConfig,
	chainID int64,
	cleanupFn func(),
	newEnvConfig vrfcommon.NewEnvConfig,
	l zerolog.Logger,
	chainlinkNodeLogScannerSettings test_env.ChainlinkNodeLogScannerSettings,
) (*test_env.CLClusterTestEnv, *vrfcommon.VRFContracts, *vrfcommon.VRFKeyData, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) {
	var (
		env               *test_env.CLClusterTestEnv
		vrfContracts      *vrfcommon.VRFContracts
		vrfKey            *vrfcommon.VRFKeyData
		nodeTypeToNodeMap map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		err               error
	)
	if *testConfig.VRFv2.General.UseExistingEnv {
		vrfContracts, vrfKey, env, err = SetupVRFV2ForExistingEnv(t, testConfig, chainID, cleanupFn, l)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "Error setting up VRF V2 for Existing env", err)
		}
	} else {
		vrfContracts, vrfKey, env, nodeTypeToNodeMap, err = SetupVRFV2ForNewEnv(ctx, t, testConfig, chainID, cleanupFn, newEnvConfig, l, chainlinkNodeLogScannerSettings)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "Error setting up VRF V2 for New env", err)
		}
	}
	return env, vrfContracts, vrfKey, nodeTypeToNodeMap, nil
}

func SetupVRFV2ForNewEnv(
	ctx context.Context,
	t *testing.T,
	testConfig tc.TestConfig,
	chainID int64,
	cleanupFn func(),
	newEnvConfig vrfcommon.NewEnvConfig,
	l zerolog.Logger,
	chainlinkNodeLogScannerSettings test_env.ChainlinkNodeLogScannerSettings,
) (*vrfcommon.VRFContracts, *vrfcommon.VRFKeyData, *test_env.CLClusterTestEnv, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) {
	network, err := actions.EthereumNetworkConfigFromConfig(l, &testConfig)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "Error building ethereum network config", err)
	}

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&testConfig).
		WithPrivateEthereumNetwork(network.EthereumNetworkConfig).
		WithCLNodes(len(newEnvConfig.NodesToCreate)).
		WithFunding(big.NewFloat(*testConfig.Common.ChainlinkNodeFunding)).
		WithChainlinkNodeLogScanner(chainlinkNodeLogScannerSettings).
		WithCustomCleanup(cleanupFn).
		WithSeth().
		Build()

	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error creating test env", err)
	}

	sethClient, err := env.GetSethClientForSelectedNetwork()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	mockETHLinkFeed, err := contracts.DeployVRFMockETHLINKFeed(sethClient, big.NewInt(*testConfig.VRFv2.General.LinkNativeFeedResponse))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error deploying mock ETH/LINK feed", err)
	}

	linkToken, err := contracts.DeployLinkTokenContract(l, sethClient)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error deploying LINK contract", err)
	}

	vrfContracts, vrfKey, nodeTypeToNode, err := SetupVRFV2Environment(
		ctx,
		env,
		chainID,
		newEnvConfig.NodesToCreate,
		&testConfig,
		newEnvConfig.UseVRFOwner,
		newEnvConfig.UseTestCoordinator,
		linkToken,
		mockETHLinkFeed,
		//register proving key against EOA address in order to return funds to this address
		sethClient.MustGetRootKeyAddress().Hex(),
		newEnvConfig.NumberOfTxKeysToCreate,
		l,
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error setting up VRF v2 env", err)
	}
	return vrfContracts, vrfKey, env, nodeTypeToNode, nil
}

func SetupVRFV2ForExistingEnv(t *testing.T, testConfig tc.TestConfig, chainID int64, cleanupFn func(), l zerolog.Logger) (*vrfcommon.VRFContracts, *vrfcommon.VRFKeyData, *test_env.CLClusterTestEnv, error) {
	commonExistingEnvConfig := testConfig.VRFv2.ExistingEnvConfig.ExistingEnvConfig
	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&testConfig).
		WithCustomCleanup(cleanupFn).
		WithSeth().
		Build()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err: %w", "error creating test env", err)
	}
	client, err := env.GetSethClientForSelectedNetwork()
	if err != nil {
		return nil, nil, nil, err
	}
	coordinator, err := contracts.LoadVRFCoordinatorV2(client, *commonExistingEnvConfig.ConsumerAddress)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err: %w", "error loading VRFCoordinator2", err)
	}
	linkAddr := common.HexToAddress(*commonExistingEnvConfig.LinkAddress)
	linkToken, err := contracts.LoadLinkTokenContract(l, client, linkAddr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err: %w", "error loading LinkToken", err)
	}

	sethClient, err := env.GetSethClient(chainID)
	if err != nil {
		return nil, nil, nil, err
	}

	err = vrfcommon.FundNodesIfNeeded(testcontext.Get(t), commonExistingEnvConfig, sethClient, l)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("err: %w", err)
	}
	vrfContracts := &vrfcommon.VRFContracts{
		CoordinatorV2:  coordinator,
		VRFV2Consumers: nil,
		LinkToken:      linkToken,
		BHS:            nil,
	}
	vrfKey := &vrfcommon.VRFKeyData{
		VRFKey:            nil,
		EncodedProvingKey: [2]*big.Int{},
		KeyHash:           common.HexToHash(*commonExistingEnvConfig.KeyHash),
	}
	return vrfContracts, vrfKey, env, nil
}

func SetupSubsAndConsumersForExistingEnv(
	env *test_env.CLClusterTestEnv,
	chainID int64,
	coordinator contracts.VRFCoordinatorV2,
	linkToken contracts.LinkToken,
	numberOfConsumerContractsToDeployAndAddToSub int,
	numberOfSubToCreate int,
	testConfig tc.TestConfig,
	l zerolog.Logger,
) ([]uint64, []contracts.VRFv2LoadTestConsumer, error) {
	var (
		subIDs    []uint64
		consumers []contracts.VRFv2LoadTestConsumer
		err       error
	)
	if *testConfig.VRFv2.General.UseExistingEnv {
		commonExistingEnvConfig := testConfig.VRFv2.ExistingEnvConfig.ExistingEnvConfig
		if *commonExistingEnvConfig.CreateFundSubsAndAddConsumers {
			consumers, subIDs, err = SetupNewConsumersAndSubs(
				env,
				chainID,
				coordinator,
				testConfig,
				linkToken,
				numberOfConsumerContractsToDeployAndAddToSub,
				numberOfSubToCreate,
				l,
			)
			if err != nil {
				return nil, nil, fmt.Errorf("err: %w", err)
			}
		} else {
			client, err := env.GetSethClient(chainID)
			if err != nil {
				return nil, nil, err
			}
			addr := common.HexToAddress(*commonExistingEnvConfig.ConsumerAddress)
			consumer, err := contracts.LoadVRFv2LoadTestConsumer(client, addr)
			if err != nil {
				return nil, nil, fmt.Errorf("err: %w", err)
			}
			consumers = append(consumers, consumer)
			subIDs = append(subIDs, *testConfig.VRFv2.ExistingEnvConfig.SubID)
		}
	} else {
		consumers, subIDs, err = SetupNewConsumersAndSubs(
			env,
			chainID,
			coordinator,
			testConfig,
			linkToken,
			numberOfConsumerContractsToDeployAndAddToSub,
			numberOfSubToCreate,
			l,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("err: %w", err)
		}
	}
	return subIDs, consumers, nil
}
