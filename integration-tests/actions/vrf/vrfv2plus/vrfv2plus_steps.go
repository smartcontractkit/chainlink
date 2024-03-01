package vrfv2plus

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"golang.org/x/sync/errgroup"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	testconfig "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_wrapper_load_test_consumer"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	vrfv2plus_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/types"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_upgraded_version"
)

func DeployVRFV2_5Contracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
	consumerContractsAmount int,
) (*vrfcommon.VRFContracts, error) {
	bhs, err := contractDeployer.DeployBlockhashStore()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrDeployBlockHashStore, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2_5(bhs.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrDeployCoordinator, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	consumers, err := DeployVRFV2PlusConsumers(contractDeployer, coordinator, consumerContractsAmount)
	if err != nil {
		return nil, err
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return &vrfcommon.VRFContracts{
		CoordinatorV2Plus: coordinator,
		BHS:               bhs,
		VRFV2PlusConsumer: consumers,
	}, nil
}

func DeployVRFV2PlusConsumers(contractDeployer contracts.ContractDeployer, coordinator contracts.VRFCoordinatorV2_5, consumerContractsAmount int) ([]contracts.VRFv2PlusLoadTestConsumer, error) {
	var consumers []contracts.VRFv2PlusLoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contractDeployer.DeployVRFv2PlusLoadTestConsumer(coordinator.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func CreateVRFV2PlusJob(
	chainlinkNode *client.ChainlinkClient,
	vrfJobSpecConfig vrfcommon.VRFJobSpecConfig,
) (*client.Job, error) {
	jobUUID := uuid.New()
	os := &client.VRFV2PlusTxPipelineSpec{
		Address:               vrfJobSpecConfig.CoordinatorAddress,
		EstimateGasMultiplier: vrfJobSpecConfig.EstimateGasMultiplier,
		FromAddress:           vrfJobSpecConfig.FromAddresses[0],
		SimulationBlock:       vrfJobSpecConfig.SimulationBlock,
	}
	ost, err := os.String()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrParseJob, err)
	}

	job, err := chainlinkNode.MustCreateJob(&client.VRFV2PlusJobSpec{
		Name:                          fmt.Sprintf("vrf-v2-plus-%s", jobUUID),
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
	})
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrCreatingVRFv2PlusJob, err)
	}
	return job, nil
}

func VRFV2_5RegisterProvingKey(
	vrfKey *client.VRFKey,
	coordinator contracts.VRFCoordinatorV2_5,
	gasLaneMaxGas uint64,
) (vrfcommon.VRFEncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return vrfcommon.VRFEncodedProvingKey{}, fmt.Errorf("%s, err %w", vrfcommon.ErrEncodingProvingKey, err)
	}
	err = coordinator.RegisterProvingKey(
		provingKey,
		gasLaneMaxGas,
	)
	if err != nil {
		return vrfcommon.VRFEncodedProvingKey{}, fmt.Errorf("%s, err %w", vrfcommon.ErrRegisterProvingKey, err)
	}
	return provingKey, nil
}

func VRFV2PlusUpgradedVersionRegisterProvingKey(
	vrfKey *client.VRFKey,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
) (vrfcommon.VRFEncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return vrfcommon.VRFEncodedProvingKey{}, fmt.Errorf("%s, err %w", vrfcommon.ErrEncodingProvingKey, err)
	}
	err = coordinator.RegisterProvingKey(
		provingKey,
	)
	if err != nil {
		return vrfcommon.VRFEncodedProvingKey{}, fmt.Errorf("%s, err %w", vrfcommon.ErrRegisterProvingKey, err)
	}
	return provingKey, nil
}

func FundVRFCoordinatorV2_5Subscription(
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2_5,
	chainClient blockchain.EVMClient,
	subscriptionID *big.Int,
	linkFundingAmountJuels *big.Int,
) error {
	encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint256"}]`, subscriptionID)
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrABIEncodingFunding, err)
	}
	_, err = linkToken.TransferAndCall(coordinator.Address(), linkFundingAmountJuels, encodedSubId)
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrSendingLinkToken, err)
	}
	return chainClient.WaitForEvents()
}

// SetupVRFV2_5Environment will create specified number of subscriptions and add the same conumer/s to each of them
func SetupVRFV2_5Environment(
	env *test_env.CLClusterTestEnv,
	nodesToCreate []vrfcommon.VRFNodeType,
	vrfv2PlusTestConfig types.VRFv2PlusTestConfig,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	numberOfTxKeysToCreate int,
	numberOfConsumers int,
	numberOfSubToCreate int,
	l zerolog.Logger,
) (*vrfcommon.VRFContracts, []*big.Int, *vrfcommon.VRFKeyData, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) {
	l.Info().Msg("Starting VRFV2 Plus environment setup")
	configGeneral := vrfv2PlusTestConfig.GetVRFv2PlusConfig().General
	vrfContracts, subIDs, err := SetupVRFV2PlusContracts(
		env,
		linkToken,
		mockNativeLINKFeed,
		configGeneral,
		numberOfSubToCreate,
		numberOfConsumers,
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

	l.Info().Str("Coordinator", vrfContracts.CoordinatorV2Plus.Address()).Msg("Registering Proving Key")
	provingKey, err := VRFV2_5RegisterProvingKey(vrfKey, vrfContracts.CoordinatorV2Plus, uint64(assets.GWei(*configGeneral.CLNodeMaxGasPriceGWei).Int64()))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRegisteringProvingKey, err)
	}
	keyHash, err := vrfContracts.CoordinatorV2Plus.HashOfKey(context.Background(), provingKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrCreatingProvingKeyHash, err)
	}

	chainID := env.EVMClient.GetChainID()
	vrfTXKeyAddressStrings, _, err := vrfcommon.CreateFundAndGetSendingKeys(
		env.EVMClient,
		nodeTypeToNodeMap[vrfcommon.VRF],
		*vrfv2PlusTestConfig.GetCommonConfig().ChainlinkNodeFunding,
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

	g := errgroup.Group{}
	if vrfNode, exists := nodeTypeToNodeMap[vrfcommon.VRF]; exists {
		g.Go(func() error {
			err := setupVRFNode(vrfContracts, chainID, configGeneral, pubKeyCompressed, l, vrfNode)
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
				vrfContracts.CoordinatorV2Plus.Address(),
				vrfContracts.BHS.Address(),
				*vrfv2PlusTestConfig.GetCommonConfig().ChainlinkNodeFunding,
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

	l.Info().Msg("VRFV2 Plus environment setup is finished")
	return vrfContracts, subIDs, &vrfKeyData, nodeTypeToNodeMap, nil
}

func SetupVRFV2PlusContracts(
	env *test_env.CLClusterTestEnv,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	configGeneral *testconfig.General,
	numberOfSubToCreate int,
	numberOfConsumers int,
	l zerolog.Logger,
) (*vrfcommon.VRFContracts, []*big.Int, error) {
	l.Info().Msg("Deploying VRFV2 Plus contracts")
	vrfContracts, err := DeployVRFV2_5Contracts(env.ContractDeployer, env.EVMClient, numberOfConsumers)
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrDeployVRFV2_5Contracts, err)
	}

	l.Info().Str("Coordinator", vrfContracts.CoordinatorV2Plus.Address()).Msg("Setting Coordinator Config")
	err = vrfContracts.CoordinatorV2Plus.SetConfig(
		*configGeneral.MinimumConfirmations,
		*configGeneral.MaxGasLimitCoordinatorConfig,
		*configGeneral.StalenessSeconds,
		*configGeneral.GasAfterPaymentCalculation,
		big.NewInt(*configGeneral.FallbackWeiPerUnitLink),
		*configGeneral.FulfillmentFlatFeeNativePPM,
		*configGeneral.FulfillmentFlatFeeLinkDiscountPPM,
		*configGeneral.NativePremiumPercentage,
		*configGeneral.LinkPremiumPercentage,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrSetVRFCoordinatorConfig, err)
	}

	l.Info().Str("Coordinator", vrfContracts.CoordinatorV2Plus.Address()).Msg("Setting Link and ETH/LINK feed")
	err = vrfContracts.CoordinatorV2Plus.SetLINKAndLINKNativeFeed(linkToken.Address(), mockNativeLINKFeed.Address())
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrSetLinkNativeLinkFeed, err)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	l.Info().
		Str("Coordinator", vrfContracts.CoordinatorV2Plus.Address()).
		Int("Number of Subs to create", numberOfSubToCreate).
		Msg("Creating and funding subscriptions, adding consumers")
	subIDs, err := CreateFundSubsAndAddConsumers(
		env,
		big.NewFloat(*configGeneral.SubscriptionFundingAmountNative),
		big.NewFloat(*configGeneral.SubscriptionFundingAmountLink),
		linkToken,
		vrfContracts.CoordinatorV2Plus, vrfContracts.VRFV2PlusConsumer,
		numberOfSubToCreate,
	)
	if err != nil {
		return nil, nil, err
	}
	return vrfContracts, subIDs, nil
}

func setupVRFNode(contracts *vrfcommon.VRFContracts, chainID *big.Int, config *vrfv2plus_config.General, pubKeyCompressed string, l zerolog.Logger, vrfNode *vrfcommon.VRFNode) error {
	vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
		ForwardingAllowed:             *config.VRFJobForwardingAllowed,
		CoordinatorAddress:            contracts.CoordinatorV2Plus.Address(),
		FromAddresses:                 vrfNode.TXKeyAddressStrings,
		EVMChainID:                    chainID.String(),
		MinIncomingConfirmations:      int(*config.MinimumConfirmations),
		PublicKey:                     pubKeyCompressed,
		EstimateGasMultiplier:         *config.VRFJobEstimateGasMultiplier,
		BatchFulfillmentEnabled:       *config.VRFJobBatchFulfillmentEnabled,
		BatchFulfillmentGasMultiplier: *config.VRFJobBatchFulfillmentGasMultiplier,
		PollPeriod:                    config.VRFJobPollPeriod.Duration,
		RequestTimeout:                config.VRFJobRequestTimeout.Duration,
		SimulationBlock:               config.VRFJobSimulationBlock,
		VRFOwnerConfig:                nil,
	}

	l.Info().Msg("Creating VRFV2 Plus Job")
	job, err := CreateVRFV2PlusJob(
		vrfNode.CLNode.API,
		vrfJobSpecConfig,
	)
	if err != nil {
		return fmt.Errorf("%s, err %w", ErrCreateVRFV2PlusJobs, err)
	}
	vrfNode.Job = job

	// this part is here because VRFv2 can work with only a specific key
	// [[EVM.KeySpecific]]
	//	Key = '...'
	nodeConfig := node.NewConfig(vrfNode.CLNode.NodeConfig,
		node.WithLogPollInterval(1*time.Second),
		node.WithVRFv2EVMEstimator(vrfNode.TXKeyAddressStrings, *config.CLNodeMaxGasPriceGWei),
	)
	l.Info().Msg("Restarting Node with new sending key PriceMax configuration")
	err = vrfNode.CLNode.Restart(nodeConfig)
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrRestartCLNode, err)
	}
	return nil
}

func CreateFundSubsAndAddConsumers(
	env *test_env.CLClusterTestEnv,
	subscriptionFundingAmountNative *big.Float,
	subscriptionFundingAmountLink *big.Float,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2_5,
	consumers []contracts.VRFv2PlusLoadTestConsumer,
	numberOfSubToCreate int,
) ([]*big.Int, error) {
	subIDs, err := CreateSubsAndFund(
		env,
		subscriptionFundingAmountNative,
		subscriptionFundingAmountLink,
		linkToken,
		coordinator,
		numberOfSubToCreate,
	)
	if err != nil {
		return nil, err
	}
	subToConsumersMap := map[*big.Int][]contracts.VRFv2PlusLoadTestConsumer{}

	//each subscription will have the same consumers
	for _, subID := range subIDs {
		subToConsumersMap[subID] = consumers
	}

	err = AddConsumersToSubs(
		subToConsumersMap,
		coordinator,
	)
	if err != nil {
		return nil, err
	}

	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return subIDs, nil
}

func CreateSubsAndFund(
	env *test_env.CLClusterTestEnv,
	subscriptionFundingAmountNative *big.Float,
	subscriptionFundingAmountLink *big.Float,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2_5,
	subAmountToCreate int,
) ([]*big.Int, error) {
	subs, err := CreateSubs(env, coordinator, subAmountToCreate)
	if err != nil {
		return nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	err = FundSubscriptions(
		env,
		subscriptionFundingAmountNative,
		subscriptionFundingAmountLink,
		linkToken,
		coordinator,
		subs,
	)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func CreateSubs(
	env *test_env.CLClusterTestEnv,
	coordinator contracts.VRFCoordinatorV2_5,
	subAmountToCreate int,
) ([]*big.Int, error) {
	var subIDArr []*big.Int

	for i := 0; i < subAmountToCreate; i++ {
		subID, err := CreateSubAndFindSubID(env, coordinator)
		if err != nil {
			return nil, err
		}
		subIDArr = append(subIDArr, subID)
	}
	return subIDArr, nil
}

func AddConsumersToSubs(
	subToConsumerMap map[*big.Int][]contracts.VRFv2PlusLoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2_5,
) error {
	for subID, consumers := range subToConsumerMap {
		for _, consumer := range consumers {
			err := coordinator.AddConsumer(subID, consumer.Address())
			if err != nil {
				return fmt.Errorf("%s, err %w", ErrAddConsumerToSub, err)
			}
		}
	}
	return nil
}

func CreateSubAndFindSubID(env *test_env.CLClusterTestEnv, coordinator contracts.VRFCoordinatorV2_5) (*big.Int, error) {
	tx, err := coordinator.CreateSubscription()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrCreateVRFSubscription, err)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	receipt, err := env.EVMClient.GetTxReceipt(tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	//SubscriptionsCreated Log should be emitted with the subscription ID
	subID := receipt.Logs[0].Topics[1].Big()

	return subID, nil
}

func FundSubscriptions(
	env *test_env.CLClusterTestEnv,
	subscriptionFundingAmountNative *big.Float,
	subscriptionFundingAmountLink *big.Float,
	linkAddress contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2_5,
	subIDs []*big.Int,
) error {
	for _, subID := range subIDs {
		//Native Billing
		amountWei := conversions.EtherToWei(subscriptionFundingAmountNative)
		err := coordinator.FundSubscriptionWithNative(
			subID,
			amountWei,
		)
		if err != nil {
			return fmt.Errorf("%s, err %w", ErrFundSubWithNativeToken, err)
		}
		//Link Billing
		amountJuels := conversions.EtherToWei(subscriptionFundingAmountLink)
		err = FundVRFCoordinatorV2_5Subscription(linkAddress, coordinator, env.EVMClient, subID, amountJuels)
		if err != nil {
			return fmt.Errorf("%s, err %w", vrfcommon.ErrFundSubWithLinkToken, err)
		}
	}
	err := env.EVMClient.WaitForEvents()
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return nil
}

func GetUpgradedCoordinatorTotalBalance(coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion) (linkTotalBalance *big.Int, nativeTokenTotalBalance *big.Int, err error) {
	linkTotalBalance, err = coordinator.GetLinkTotalBalance(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrLinkTotalBalance, err)
	}
	nativeTokenTotalBalance, err = coordinator.GetNativeTokenTotalBalance(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrNativeTokenBalance, err)
	}
	return
}

func GetCoordinatorTotalBalance(coordinator contracts.VRFCoordinatorV2_5) (linkTotalBalance *big.Int, nativeTokenTotalBalance *big.Int, err error) {
	linkTotalBalance, err = coordinator.GetLinkTotalBalance(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrLinkTotalBalance, err)
	}
	nativeTokenTotalBalance, err = coordinator.GetNativeTokenTotalBalance(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrNativeTokenBalance, err)
	}
	return
}

func RequestRandomnessAndWaitForFulfillment(
	consumer contracts.VRFv2PlusLoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2_5,
	vrfKeyData *vrfcommon.VRFKeyData,
	subID *big.Int,
	isNativeBilling bool,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled, error) {
	logRandRequest(
		l,
		consumer.Address(),
		coordinator.Address(),
		subID,
		isNativeBilling,
		minimumConfirmations,
		callbackGasLimit,
		numberOfWords,
		vrfKeyData.KeyHash,
		randomnessRequestCountPerRequest,
		randomnessRequestCountPerRequestDeviation,
	)
	_, err := consumer.RequestRandomness(
		vrfKeyData.KeyHash,
		subID,
		minimumConfirmations,
		callbackGasLimit,
		isNativeBilling,
		numberOfWords,
		randomnessRequestCountPerRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRequestRandomness, err)
	}

	return WaitForRequestAndFulfillmentEvents(
		consumer.Address(),
		coordinator,
		vrfKeyData,
		subID,
		isNativeBilling,
		randomWordsFulfilledEventTimeout,
		l,
	)
}

func RequestRandomnessAndWaitForFulfillmentUpgraded(
	consumer contracts.VRFv2PlusLoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
	vrfKeyData *vrfcommon.VRFKeyData,
	subID *big.Int,
	isNativeBilling bool,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled, error) {
	logRandRequest(
		l,
		consumer.Address(),
		coordinator.Address(),
		subID,
		isNativeBilling,
		minimumConfirmations,
		callbackGasLimit,
		numberOfWords,
		vrfKeyData.KeyHash,
		randomnessRequestCountPerRequest,
		randomnessRequestCountPerRequestDeviation,
	)
	_, err := consumer.RequestRandomness(
		vrfKeyData.KeyHash,
		subID,
		minimumConfirmations,
		callbackGasLimit,
		isNativeBilling,
		numberOfWords,
		randomnessRequestCountPerRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRequestRandomness, err)
	}

	return WaitForRequestAndFulfillmentEventsUpgraded(
		consumer.Address(),
		coordinator,
		vrfKeyData,
		subID,
		isNativeBilling,
		randomWordsFulfilledEventTimeout,
		l,
	)
}

func SetupVRFV2PlusWrapperEnvironment(
	env *test_env.CLClusterTestEnv,
	vrfv2PlusTestConfig types.VRFv2PlusTestConfig,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	coordinator contracts.VRFCoordinatorV2_5,
	keyHash [32]byte,
	wrapperConsumerContractsAmount int,
) (*VRFV2PlusWrapperContracts, *big.Int, error) {
	vrfv2PlusConfig := vrfv2PlusTestConfig.GetVRFv2PlusConfig().General
	wrapperContracts, err := DeployVRFV2PlusDirectFundingContracts(
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
	err = wrapperContracts.VRFV2PlusWrapper.SetConfig(
		*vrfv2PlusConfig.WrapperGasOverhead,
		*vrfv2PlusConfig.CoordinatorGasOverhead,
		*vrfv2PlusConfig.WrapperPremiumPercentage,
		keyHash,
		*vrfv2PlusConfig.WrapperMaxNumberOfWords,
		*vrfv2PlusConfig.StalenessSeconds,
		big.NewInt(*vrfv2PlusConfig.FallbackWeiPerUnitLink),
		*vrfv2PlusConfig.FulfillmentFlatFeeLinkPPM,
		*vrfv2PlusConfig.FulfillmentFlatFeeNativePPM,
	)
	if err != nil {
		return nil, nil, err
	}

	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	//fund sub
	wrapperSubID, err := wrapperContracts.VRFV2PlusWrapper.GetSubID(context.Background())
	if err != nil {
		return nil, nil, err
	}

	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	err = FundSubscriptions(
		env,
		big.NewFloat(*vrfv2PlusTestConfig.GetVRFv2PlusConfig().General.SubscriptionFundingAmountNative),
		big.NewFloat(*vrfv2PlusTestConfig.GetVRFv2PlusConfig().General.SubscriptionFundingAmountLink),
		linkToken,
		coordinator,
		[]*big.Int{wrapperSubID},
	)
	if err != nil {
		return nil, nil, err
	}

	//fund consumer with Link
	err = linkToken.Transfer(
		wrapperContracts.LoadTestConsumers[0].Address(),
		big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(*vrfv2PlusConfig.WrapperConsumerFundingAmountLink)),
	)
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	//fund consumer with Eth
	err = wrapperContracts.LoadTestConsumers[0].Fund(big.NewFloat(*vrfv2PlusConfig.WrapperConsumerFundingAmountNativeToken))
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return wrapperContracts, wrapperSubID, nil
}

func DeployVRFV2PlusWrapperConsumers(contractDeployer contracts.ContractDeployer, linkTokenAddress string, vrfV2PlusWrapper contracts.VRFV2PlusWrapper, consumerContractsAmount int) ([]contracts.VRFv2PlusWrapperLoadTestConsumer, error) {
	var consumers []contracts.VRFv2PlusWrapperLoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contractDeployer.DeployVRFV2PlusWrapperLoadTestConsumer(linkTokenAddress, vrfV2PlusWrapper.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func DeployVRFV2PlusDirectFundingContracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
	linkTokenAddress string,
	linkEthFeedAddress string,
	coordinator contracts.VRFCoordinatorV2_5,
	consumerContractsAmount int,
) (*VRFV2PlusWrapperContracts, error) {

	vrfv2PlusWrapper, err := contractDeployer.DeployVRFV2PlusWrapper(linkTokenAddress, linkEthFeedAddress, coordinator.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployWrapper, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	consumers, err := DeployVRFV2PlusWrapperConsumers(contractDeployer, linkTokenAddress, vrfv2PlusWrapper, consumerContractsAmount)
	if err != nil {
		return nil, err
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return &VRFV2PlusWrapperContracts{vrfv2PlusWrapper, consumers}, nil
}

func WrapperRequestRandomness(
	consumer contracts.VRFv2PlusWrapperLoadTestConsumer,
	coordinatorAddress string,
	vrfKeyData *vrfcommon.VRFKeyData,
	subID *big.Int,
	isNativeBilling bool,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	l zerolog.Logger) (string, error) {
	logRandRequest(
		l,
		consumer.Address(),
		coordinatorAddress,
		subID,
		isNativeBilling,
		minimumConfirmations,
		callbackGasLimit,
		numberOfWords,
		vrfKeyData.KeyHash,
		randomnessRequestCountPerRequest,
		randomnessRequestCountPerRequestDeviation,
	)
	if isNativeBilling {
		_, err := consumer.RequestRandomnessNative(
			minimumConfirmations,
			callbackGasLimit,
			numberOfWords,
			randomnessRequestCountPerRequest,
		)
		if err != nil {
			return "", fmt.Errorf("%s, err %w", ErrRequestRandomnessDirectFundingNativePayment, err)
		}
	} else {
		_, err := consumer.RequestRandomness(
			minimumConfirmations,
			callbackGasLimit,
			numberOfWords,
			randomnessRequestCountPerRequest,
		)
		if err != nil {
			return "", fmt.Errorf("%s, err %w", ErrRequestRandomnessDirectFundingLinkPayment, err)
		}
	}
	wrapperAddress, err := consumer.GetWrapper(context.Background())
	if err != nil {
		return "", fmt.Errorf("error getting wrapper address, err: %w", err)
	}
	return wrapperAddress.Hex(), nil
}

func DirectFundingRequestRandomnessAndWaitForFulfillment(
	consumer contracts.VRFv2PlusWrapperLoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2_5,
	vrfKeyData *vrfcommon.VRFKeyData,
	subID *big.Int,
	isNativeBilling bool,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled, error) {
	wrapperAddress, err := WrapperRequestRandomness(consumer, coordinator.Address(), vrfKeyData, subID,
		isNativeBilling, minimumConfirmations, callbackGasLimit, numberOfWords,
		randomnessRequestCountPerRequest, randomnessRequestCountPerRequestDeviation,
		l)
	if err != nil {
		return nil, fmt.Errorf("error getting wrapper address, err: %w", err)
	}
	return WaitForRequestAndFulfillmentEvents(
		wrapperAddress,
		coordinator,
		vrfKeyData,
		subID,
		isNativeBilling,
		randomWordsFulfilledEventTimeout,
		l,
	)
}

func DirectFundingRequestRandomnessAndWaitForFulfillmentUpgraded(
	consumer contracts.VRFv2PlusWrapperLoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
	vrfKeyData *vrfcommon.VRFKeyData,
	subID *big.Int,
	isNativeBilling bool,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled, error) {
	wrapperAddress, err := WrapperRequestRandomness(consumer, coordinator.Address(), vrfKeyData, subID,
		isNativeBilling, minimumConfirmations, callbackGasLimit, numberOfWords,
		randomnessRequestCountPerRequest, randomnessRequestCountPerRequestDeviation,
		l)
	if err != nil {
		return nil, fmt.Errorf("error getting wrapper address, err: %w", err)
	}
	return WaitForRequestAndFulfillmentEventsUpgraded(
		wrapperAddress,
		coordinator,
		vrfKeyData,
		subID,
		isNativeBilling,
		randomWordsFulfilledEventTimeout,
		l,
	)
}

func WaitForRequestAndFulfillmentEvents(
	consumerAddress string,
	coordinator contracts.VRFCoordinatorV2_5,
	vrfKeyData *vrfcommon.VRFKeyData,
	subID *big.Int,
	isNativeBilling bool,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled, error) {
	randomWordsRequestedEvent, err := coordinator.WaitForRandomWordsRequestedEvent(
		[][32]byte{vrfKeyData.KeyHash},
		[]*big.Int{subID},
		[]common.Address{common.HexToAddress(consumerAddress)},
		time.Minute*1,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitRandomWordsRequestedEvent, err)
	}

	LogRandomnessRequestedEvent(l, coordinator, randomWordsRequestedEvent, isNativeBilling)

	randomWordsFulfilledEvent, err := coordinator.WaitForRandomWordsFulfilledEvent(
		[]*big.Int{subID},
		[]*big.Int{randomWordsRequestedEvent.RequestId},
		randomWordsFulfilledEventTimeout,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitRandomWordsFulfilledEvent, err)
	}

	LogRandomWordsFulfilledEvent(l, coordinator, randomWordsFulfilledEvent, isNativeBilling)
	return randomWordsFulfilledEvent, err
}

func WaitForRequestAndFulfillmentEventsUpgraded(
	consumerAddress string,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
	vrfKeyData *vrfcommon.VRFKeyData,
	subID *big.Int,
	isNativeBilling bool,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled, error) {
	randomWordsRequestedEvent, err := coordinator.WaitForRandomWordsRequestedEvent(
		[][32]byte{vrfKeyData.KeyHash},
		[]*big.Int{subID},
		[]common.Address{common.HexToAddress(consumerAddress)},
		time.Minute*1,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitRandomWordsRequestedEvent, err)
	}

	LogRandomnessRequestedEventUpgraded(l, coordinator, randomWordsRequestedEvent, isNativeBilling)

	randomWordsFulfilledEvent, err := coordinator.WaitForRandomWordsFulfilledEvent(
		[]*big.Int{subID},
		[]*big.Int{randomWordsRequestedEvent.RequestId},
		randomWordsFulfilledEventTimeout,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitRandomWordsFulfilledEvent, err)
	}
	LogRandomWordsFulfilledEventUpgraded(l, coordinator, randomWordsFulfilledEvent, isNativeBilling)
	return randomWordsFulfilledEvent, err
}

func ReturnFundsForFulfilledRequests(client blockchain.EVMClient, coordinator contracts.VRFCoordinatorV2_5, l zerolog.Logger) error {
	linkTotalBalance, err := coordinator.GetLinkTotalBalance(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting LINK total balance, err: %w", err)
	}
	defaultWallet := client.GetDefaultWallet().Address()
	l.Info().
		Str("LINK amount", linkTotalBalance.String()).
		Str("Returning to", defaultWallet).
		Msg("Returning LINK for fulfilled requests")
	err = coordinator.Withdraw(
		common.HexToAddress(defaultWallet),
	)
	if err != nil {
		return fmt.Errorf("Error withdrawing LINK from coordinator to default wallet, err: %w", err)
	}
	nativeTotalBalance, err := coordinator.GetNativeTokenTotalBalance(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting NATIVE total balance, err: %w", err)
	}
	l.Info().
		Str("Native Token amount", nativeTotalBalance.String()).
		Str("Returning to", defaultWallet).
		Msg("Returning Native Token for fulfilled requests")
	err = coordinator.WithdrawNative(
		common.HexToAddress(defaultWallet),
	)
	if err != nil {
		return fmt.Errorf("Error withdrawing NATIVE from coordinator to default wallet, err: %w", err)
	}
	return nil
}

func LogSubDetails(l zerolog.Logger, subscription vrf_coordinator_v2_5.GetSubscription, subID *big.Int, coordinator contracts.VRFCoordinatorV2_5) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Str("Link Balance", (*commonassets.Link)(subscription.Balance).Link()).
		Str("Native Token Balance", assets.FormatWei(subscription.NativeBalance)).
		Str("Subscription ID", subID.String()).
		Str("Subscription Owner", subscription.Owner.String()).
		Interface("Subscription Consumers", subscription.Consumers).
		Msg("Subscription Data")
}

func LogRandomnessRequestedEventUpgraded(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
	randomWordsRequestedEvent *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested,
	isNativeBilling bool,
) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Bool("Native Billing", isNativeBilling).
		Str("Request ID", randomWordsRequestedEvent.RequestId.String()).
		Str("Subscription ID", randomWordsRequestedEvent.SubId.String()).
		Str("Sender Address", randomWordsRequestedEvent.Sender.String()).
		Str("Keyhash", fmt.Sprintf("0x%x", randomWordsRequestedEvent.KeyHash)).
		Uint32("Callback Gas Limit", randomWordsRequestedEvent.CallbackGasLimit).
		Uint32("Number of Words", randomWordsRequestedEvent.NumWords).
		Uint16("Minimum Request Confirmations", randomWordsRequestedEvent.MinimumRequestConfirmations).
		Msg("RandomnessRequested Event")
}

func LogRandomWordsFulfilledEventUpgraded(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
	randomWordsFulfilledEvent *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled,
	isNativeBilling bool,
) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Bool("Native Billing", isNativeBilling).
		Str("Total Payment in Juels", randomWordsFulfilledEvent.Payment.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Str("Subscription ID", randomWordsFulfilledEvent.SubID.String()).
		Str("Request ID", randomWordsFulfilledEvent.RequestId.String()).
		Bool("Success", randomWordsFulfilledEvent.Success).
		Msg("RandomWordsFulfilled Event (TX metadata)")
}

func LogRandomnessRequestedEvent(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2_5,
	randomWordsRequestedEvent *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested,
	isNativeBilling bool,
) {
	l.Info().
		Str("Coordinator", coordinator.Address()).
		Bool("Native Billing", isNativeBilling).
		Str("Request ID", randomWordsRequestedEvent.RequestId.String()).
		Str("Subscription ID", randomWordsRequestedEvent.SubId.String()).
		Str("Sender Address", randomWordsRequestedEvent.Sender.String()).
		Str("Keyhash", fmt.Sprintf("0x%x", randomWordsRequestedEvent.KeyHash)).
		Uint32("Callback Gas Limit", randomWordsRequestedEvent.CallbackGasLimit).
		Uint32("Number of Words", randomWordsRequestedEvent.NumWords).
		Uint16("Minimum Request Confirmations", randomWordsRequestedEvent.MinimumRequestConfirmations).
		Str("TX Hash", randomWordsRequestedEvent.Raw.TxHash.String()).
		Uint64("BlockNumber", randomWordsRequestedEvent.Raw.BlockNumber).
		Str("BlockHash", randomWordsRequestedEvent.Raw.BlockHash.String()).
		Msg("RandomnessRequested Event")
}

func LogRandomWordsFulfilledEvent(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2_5,
	randomWordsFulfilledEvent *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled,
	isNativeBilling bool,
) {
	l.Info().
		Bool("Native Billing", isNativeBilling).
		Str("Coordinator", coordinator.Address()).
		Str("Total Payment", randomWordsFulfilledEvent.Payment.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Str("Subscription ID", randomWordsFulfilledEvent.SubId.String()).
		Str("Request ID", randomWordsFulfilledEvent.RequestId.String()).
		Bool("Success", randomWordsFulfilledEvent.Success).
		Uint64("BlockNumber", randomWordsFulfilledEvent.Raw.BlockNumber).
		Str("BlockHash", randomWordsFulfilledEvent.Raw.BlockHash.String()).
		Msg("RandomWordsFulfilled Event (TX metadata)")
}

func LogMigrationCompletedEvent(l zerolog.Logger, migrationCompletedEvent *vrf_coordinator_v2_5.VRFCoordinatorV25MigrationCompleted, vrfv2PlusContracts *vrfcommon.VRFContracts) {
	l.Info().
		Str("Subscription ID", migrationCompletedEvent.SubId.String()).
		Str("Migrated From Coordinator", vrfv2PlusContracts.CoordinatorV2Plus.Address()).
		Str("Migrated To Coordinator", migrationCompletedEvent.NewCoordinator.String()).
		Msg("MigrationCompleted Event")
}

func LogSubDetailsAfterMigration(l zerolog.Logger, newCoordinator contracts.VRFCoordinatorV2PlusUpgradedVersion, subID *big.Int, migratedSubscription vrf_v2plus_upgraded_version.GetSubscription) {
	l.Info().
		Str("New Coordinator", newCoordinator.Address()).
		Str("Subscription ID", subID.String()).
		Str("Juels Balance", migratedSubscription.Balance.String()).
		Str("Native Token Balance", migratedSubscription.NativeBalance.String()).
		Str("Subscription Owner", migratedSubscription.Owner.String()).
		Interface("Subscription Consumers", migratedSubscription.Consumers).
		Msg("Subscription Data After Migration to New Coordinator")
}

func LogFulfillmentDetailsLinkBilling(
	l zerolog.Logger,
	wrapperConsumerJuelsBalanceBeforeRequest *big.Int,
	wrapperConsumerJuelsBalanceAfterRequest *big.Int,
	consumerStatus vrfv2plus_wrapper_load_test_consumer.GetRequestStatus,
	randomWordsFulfilledEvent *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled,
) {
	l.Info().
		Str("Consumer Balance Before Request (Link)", (*commonassets.Link)(wrapperConsumerJuelsBalanceBeforeRequest).Link()).
		Str("Consumer Balance After Request (Link)", (*commonassets.Link)(wrapperConsumerJuelsBalanceAfterRequest).Link()).
		Bool("Fulfilment Status", consumerStatus.Fulfilled).
		Str("Paid by Consumer Contract (Link)", (*commonassets.Link)(consumerStatus.Paid).Link()).
		Str("Paid by Coordinator Sub (Link)", (*commonassets.Link)(randomWordsFulfilledEvent.Payment).Link()).
		Str("RequestTimestamp", consumerStatus.RequestTimestamp.String()).
		Str("FulfilmentTimestamp", consumerStatus.FulfilmentTimestamp.String()).
		Str("RequestBlockNumber", consumerStatus.RequestBlockNumber.String()).
		Str("FulfilmentBlockNumber", consumerStatus.FulfilmentBlockNumber.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Msg("Random Words Fulfilment Details For Link Billing")
}

func LogFulfillmentDetailsNativeBilling(
	l zerolog.Logger,
	wrapperConsumerBalanceBeforeRequestWei *big.Int,
	wrapperConsumerBalanceAfterRequestWei *big.Int,
	consumerStatus vrfv2plus_wrapper_load_test_consumer.GetRequestStatus,
	randomWordsFulfilledEvent *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled,
) {
	l.Info().
		Str("Consumer Balance Before Request", assets.FormatWei(wrapperConsumerBalanceBeforeRequestWei)).
		Str("Consumer Balance After Request", assets.FormatWei(wrapperConsumerBalanceAfterRequestWei)).
		Bool("Fulfilment Status", consumerStatus.Fulfilled).
		Str("Paid by Consumer Contract", assets.FormatWei(consumerStatus.Paid)).
		Str("Paid by Coordinator Sub", assets.FormatWei(randomWordsFulfilledEvent.Payment)).
		Str("RequestTimestamp", consumerStatus.RequestTimestamp.String()).
		Str("FulfilmentTimestamp", consumerStatus.FulfilmentTimestamp.String()).
		Str("RequestBlockNumber", consumerStatus.RequestBlockNumber.String()).
		Str("FulfilmentBlockNumber", consumerStatus.FulfilmentBlockNumber.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Msg("Random Words Request Fulfilment Details For Native Billing")
}

func logRandRequest(
	l zerolog.Logger,
	consumer string,
	coordinator string,
	subID *big.Int,
	isNativeBilling bool,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	keyHash [32]byte,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16) {
	l.Info().
		Str("Consumer", consumer).
		Str("Coordinator", coordinator).
		Str("SubID", subID.String()).
		Bool("IsNativePayment", isNativeBilling).
		Uint16("MinimumConfirmations", minimumConfirmations).
		Uint32("CallbackGasLimit", callbackGasLimit).
		Uint32("NumberOfWords", numberOfWords).
		Str("KeyHash", fmt.Sprintf("0x%x", keyHash)).
		Uint16("RandomnessRequestCountPerRequest", randomnessRequestCountPerRequest).
		Uint16("RandomnessRequestCountPerRequestDeviation", randomnessRequestCountPerRequestDeviation).
		Msg("Requesting randomness")
}
