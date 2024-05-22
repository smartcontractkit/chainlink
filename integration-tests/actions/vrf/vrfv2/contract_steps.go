package vrfv2

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	testconfig "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_owner"
)

func DeployVRFV2Contracts(
	env *test_env.CLClusterTestEnv,
	chainID int64,
	linkTokenContract contracts.LinkToken,
	linkEthFeedContract contracts.VRFMockETHLINKFeed,
	useVRFOwner bool,
	useTestCoordinator bool,
) (*vrfcommon.VRFContracts, error) {
	sethClient, err := env.GetSethClient(chainID)
	if err != nil {
		return nil, err
	}

	bhs, err := contracts.DeployBlockhashStore(sethClient)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrDeployBlockHashStore, err)
	}

	var coordinatorAddress string
	if useTestCoordinator {
		testCoordinator, err := contracts.DeployVRFCoordinatorTestV2(sethClient, linkTokenContract.Address(), bhs.Address(), linkEthFeedContract.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrDeployCoordinatorV2, err)
		}
		coordinatorAddress = testCoordinator.Address()
	} else {
		coordinator, err := contracts.DeployVRFCoordinatorV2(sethClient, linkTokenContract.Address(), bhs.Address(), linkEthFeedContract.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrDeployCoordinatorV2, err)
		}
		coordinatorAddress = coordinator.Address()
	}

	coordinator, err := contracts.LoadVRFCoordinatorV2(sethClient, coordinatorAddress)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrLoadingCoordinator, err)
	}

	batchCoordinator, err := contracts.DeployBatchVRFCoordinatorV2(sethClient, coordinator.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployBatchCoordinatorV2, err)
	}

	if useVRFOwner {
		vrfOwner, err := contracts.DeployVRFOwner(sethClient, coordinatorAddress)
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrDeployCoordinatorV2, err)
		}
		return &vrfcommon.VRFContracts{
			CoordinatorV2:      coordinator,
			BatchCoordinatorV2: batchCoordinator,
			VRFOwner:           vrfOwner,
			BHS:                bhs,
			VRFV2Consumers:     nil,
			LinkToken:          linkTokenContract,
			MockETHLINKFeed:    linkEthFeedContract,
		}, nil
	}
	return &vrfcommon.VRFContracts{
		CoordinatorV2:      coordinator,
		BatchCoordinatorV2: batchCoordinator,
		VRFOwner:           nil,
		BHS:                bhs,
		VRFV2Consumers:     nil,
		LinkToken:          linkTokenContract,
		MockETHLINKFeed:    linkEthFeedContract,
	}, nil
}

func DeployVRFV2Consumers(client *seth.Client, coordinatorAddress string, consumerContractsAmount int) ([]contracts.VRFv2LoadTestConsumer, error) {
	var consumers []contracts.VRFv2LoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contracts.DeployVRFv2LoadTestConsumer(client, coordinatorAddress)
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func DeployVRFV2WrapperConsumers(client *seth.Client, linkTokenAddress string, vrfV2Wrapper contracts.VRFV2Wrapper, consumerContractsAmount int) ([]contracts.VRFv2WrapperLoadTestConsumer, error) {
	var consumers []contracts.VRFv2WrapperLoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contracts.DeployVRFV2WrapperLoadTestConsumer(client, linkTokenAddress, vrfV2Wrapper.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func DeployVRFV2DirectFundingContracts(
	client *seth.Client,
	linkTokenAddress string,
	linkEthFeedAddress string,
	coordinator contracts.VRFCoordinatorV2,
	consumerContractsAmount int,
) (*VRFV2WrapperContracts, error) {
	vrfv2Wrapper, err := contracts.DeployVRFV2Wrapper(client, linkTokenAddress, linkEthFeedAddress, coordinator.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployVRFV2Wrapper, err)
	}
	consumers, err := DeployVRFV2WrapperConsumers(client, linkTokenAddress, vrfv2Wrapper, consumerContractsAmount)
	if err != nil {
		return nil, err
	}
	return &VRFV2WrapperContracts{vrfv2Wrapper, consumers}, nil
}

func VRFV2RegisterProvingKey(
	vrfKey *client.VRFKey,
	oracleAddress string,
	coordinator contracts.VRFCoordinatorV2,
) (vrfcommon.VRFEncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return vrfcommon.VRFEncodedProvingKey{}, fmt.Errorf("%s, err %w", vrfcommon.ErrEncodingProvingKey, err)
	}
	err = coordinator.RegisterProvingKey(
		oracleAddress,
		provingKey,
	)
	if err != nil {
		return vrfcommon.VRFEncodedProvingKey{}, fmt.Errorf("%s, err %w", vrfcommon.ErrRegisterProvingKey, err)
	}
	return provingKey, nil
}

func SetupVRFV2Contracts(
	env *test_env.CLClusterTestEnv,
	chainID int64,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.VRFMockETHLINKFeed,
	useVRFOwner bool,
	useTestCoordinator bool,
	vrfv2Config *testconfig.General,
	l zerolog.Logger,
) (*vrfcommon.VRFContracts, error) {
	l.Info().Msg("Deploying VRFV2 contracts")
	vrfContracts, err := DeployVRFV2Contracts(
		env,
		chainID,
		linkToken,
		mockNativeLINKFeed,
		useVRFOwner,
		useTestCoordinator,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployVRFV2Contracts, err)
	}

	vrfCoordinatorV2FeeConfig := vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
		FulfillmentFlatFeeLinkPPMTier1: *vrfv2Config.FulfillmentFlatFeeLinkPPMTier1,
		FulfillmentFlatFeeLinkPPMTier2: *vrfv2Config.FulfillmentFlatFeeLinkPPMTier2,
		FulfillmentFlatFeeLinkPPMTier3: *vrfv2Config.FulfillmentFlatFeeLinkPPMTier3,
		FulfillmentFlatFeeLinkPPMTier4: *vrfv2Config.FulfillmentFlatFeeLinkPPMTier4,
		FulfillmentFlatFeeLinkPPMTier5: *vrfv2Config.FulfillmentFlatFeeLinkPPMTier5,
		ReqsForTier2:                   big.NewInt(*vrfv2Config.ReqsForTier2),
		ReqsForTier3:                   big.NewInt(*vrfv2Config.ReqsForTier3),
		ReqsForTier4:                   big.NewInt(*vrfv2Config.ReqsForTier4),
		ReqsForTier5:                   big.NewInt(*vrfv2Config.ReqsForTier5)}

	l.Info().Str("Coordinator", vrfContracts.CoordinatorV2.Address()).Msg("Setting Coordinator Config")
	err = vrfContracts.CoordinatorV2.SetConfig(
		*vrfv2Config.MinimumConfirmations,
		*vrfv2Config.MaxGasLimitCoordinatorConfig,
		*vrfv2Config.StalenessSeconds,
		*vrfv2Config.GasAfterPaymentCalculation,
		decimal.RequireFromString(*vrfv2Config.FallbackWeiPerUnitLink).BigInt(),
		vrfCoordinatorV2FeeConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrSetVRFCoordinatorConfig, err)
	}
	return vrfContracts, nil
}

func setupVRFOwnerContract(contracts *vrfcommon.VRFContracts, allNativeTokenKeyAddressStrings []string, allNativeTokenKeyAddresses []common.Address, l zerolog.Logger) error {
	l.Info().Msg("Setting up VRFOwner contract")
	l.Info().
		Str("Coordinator", contracts.CoordinatorV2.Address()).
		Str("VRFOwner", contracts.VRFOwner.Address()).
		Msg("Transferring ownership of Coordinator to VRFOwner")
	err := contracts.CoordinatorV2.TransferOwnership(common.HexToAddress(contracts.VRFOwner.Address()))
	if err != nil {
		return nil
	}
	l.Info().
		Str("VRFOwner", contracts.VRFOwner.Address()).
		Msg("Accepting VRF Ownership")
	err = contracts.VRFOwner.AcceptVRFOwnership()
	if err != nil {
		return nil
	}
	l.Info().
		Strs("Authorized Senders", allNativeTokenKeyAddressStrings).
		Str("VRFOwner", contracts.VRFOwner.Address()).
		Msg("Setting authorized senders for VRFOwner contract")
	err = contracts.VRFOwner.SetAuthorizedSenders(allNativeTokenKeyAddresses)
	if err != nil {
		return nil
	}
	return err
}

func CreateFundSubsAndAddConsumers(
	subscriptionFundingAmountLink *big.Float,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	consumers []contracts.VRFv2LoadTestConsumer,
	numberOfSubToCreate int,
) ([]uint64, error) {
	subIDs, err := CreateSubsAndFund(subscriptionFundingAmountLink, linkToken, coordinator, numberOfSubToCreate)
	if err != nil {
		return nil, err
	}
	subToConsumersMap := map[uint64][]contracts.VRFv2LoadTestConsumer{}

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
	return subIDs, nil
}

func CreateSubsAndFund(
	subscriptionFundingAmountLink *big.Float,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	subAmountToCreate int,
) ([]uint64, error) {
	subs, err := CreateSubs(coordinator, subAmountToCreate)
	if err != nil {
		return nil, err
	}
	err = FundSubscriptions(subscriptionFundingAmountLink, linkToken, coordinator, subs)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func CreateSubs(
	coordinator contracts.VRFCoordinatorV2,
	subAmountToCreate int,
) ([]uint64, error) {
	var subIDArr []uint64

	for i := 0; i < subAmountToCreate; i++ {
		subID, err := CreateSubAndFindSubID(coordinator)
		if err != nil {
			return nil, err
		}
		subIDArr = append(subIDArr, subID)
	}
	return subIDArr, nil
}

func AddConsumersToSubs(
	subToConsumerMap map[uint64][]contracts.VRFv2LoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2,
) error {
	for subID, consumers := range subToConsumerMap {
		for _, consumer := range consumers {
			err := coordinator.AddConsumer(subID, consumer.Address())
			if err != nil {
				return fmt.Errorf("%s, err %w", vrfcommon.ErrAddConsumerToSub, err)
			}
		}
	}
	return nil
}

func CreateSubAndFindSubID(coordinator contracts.VRFCoordinatorV2) (uint64, error) {
	receipt, err := coordinator.CreateSubscription()
	if err != nil {
		return 0, fmt.Errorf("%s, err %w", vrfcommon.ErrCreateVRFSubscription, err)
	}
	//SubscriptionsCreated Log should be emitted with the subscription ID
	subID := receipt.Logs[0].Topics[1].Big().Uint64()

	return subID, nil
}

func FundSubscriptions(
	subscriptionFundingAmountLink *big.Float,
	linkAddress contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	subIDs []uint64,
) error {
	for _, subID := range subIDs {
		//Link Billing
		amountJuels := conversions.EtherToWei(subscriptionFundingAmountLink)
		err := FundVRFCoordinatorV2Subscription(linkAddress, coordinator, subID, amountJuels)
		if err != nil {
			return fmt.Errorf("%s, err %w", vrfcommon.ErrFundSubWithLinkToken, err)
		}
	}
	return nil
}

func FundVRFCoordinatorV2Subscription(
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	subscriptionID uint64,
	linkFundingAmountJuels *big.Int,
) error {
	encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint64"}]`, subscriptionID)
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrABIEncodingFunding, err)
	}
	_, err = linkToken.TransferAndCall(coordinator.Address(), linkFundingAmountJuels, encodedSubId)
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrSendingLinkToken, err)
	}
	return nil
}

func DirectFundingRequestRandomnessAndWaitForFulfillment(
	l zerolog.Logger,
	consumer contracts.VRFv2WrapperLoadTestConsumer,
	coordinator contracts.Coordinator,
	subID uint64,
	vrfv2KeyData *vrfcommon.VRFKeyData,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	randomWordsFulfilledEventTimeout time.Duration,
) (*contracts.CoordinatorRandomWordsFulfilled, error) {
	logRandRequest(
		l,
		consumer.Address(),
		coordinator.Address(),
		subID,
		minimumConfirmations,
		callbackGasLimit,
		numberOfWords,
		randomnessRequestCountPerRequest,
		randomnessRequestCountPerRequestDeviation,
		vrfv2KeyData.KeyHash,
		0,
	)
	randomWordsRequestedEvent, err := consumer.RequestRandomness(
		coordinator,
		minimumConfirmations,
		callbackGasLimit,
		numberOfWords,
		randomnessRequestCountPerRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRequestRandomness, err)
	}
	fulfillmentEvents, err := WaitRandomWordsFulfilledEvent(
		coordinator,
		randomWordsRequestedEvent.RequestId,
		randomWordsRequestedEvent.Raw.BlockNumber,
		randomWordsFulfilledEventTimeout,
		l,
	)
	return fulfillmentEvents, err
}

func RequestRandomnessAndWaitForFulfillment(
	l zerolog.Logger,
	consumer contracts.VRFv2LoadTestConsumer,
	coordinator contracts.Coordinator,
	subID uint64,
	vrfKeyData *vrfcommon.VRFKeyData,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	randomWordsFulfilledEventTimeout time.Duration,
	keyNum int,
) (*contracts.CoordinatorRandomWordsRequested, *contracts.CoordinatorRandomWordsFulfilled, error) {
	randomWordsRequestedEvent, err := RequestRandomness(
		l,
		consumer,
		coordinator,
		subID,
		vrfKeyData,
		minimumConfirmations,
		callbackGasLimit,
		numberOfWords,
		randomnessRequestCountPerRequest,
		randomnessRequestCountPerRequestDeviation,
		keyNum,
	)
	if err != nil {
		return nil, nil, err
	}
	randomWordsFulfilledEvent, err := WaitRandomWordsFulfilledEvent(
		coordinator,
		randomWordsRequestedEvent.RequestId,
		randomWordsRequestedEvent.Raw.BlockNumber,
		randomWordsFulfilledEventTimeout,
		l,
	)
	if err != nil {
		return nil, nil, err
	}
	return randomWordsRequestedEvent, randomWordsFulfilledEvent, nil
}

func RequestRandomness(
	l zerolog.Logger,
	consumer contracts.VRFv2LoadTestConsumer,
	coordinator contracts.Coordinator,
	subID uint64,
	vrfKeyData *vrfcommon.VRFKeyData,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	keyNum int,
) (*contracts.CoordinatorRandomWordsRequested, error) {
	logRandRequest(
		l,
		consumer.Address(),
		coordinator.Address(),
		subID,
		minimumConfirmations,
		callbackGasLimit,
		numberOfWords,
		randomnessRequestCountPerRequest,
		randomnessRequestCountPerRequestDeviation,
		vrfKeyData.KeyHash,
		keyNum,
	)
	randomWordsRequestedEvent, err := consumer.RequestRandomnessFromKey(
		coordinator,
		vrfKeyData.KeyHash,
		subID,
		minimumConfirmations,
		callbackGasLimit,
		numberOfWords,
		randomnessRequestCountPerRequest,
		keyNum,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRequestRandomness, err)
	}
	vrfcommon.LogRandomnessRequestedEvent(l, coordinator, randomWordsRequestedEvent, false, keyNum)

	return randomWordsRequestedEvent, err
}

func RequestRandomnessWithForceFulfillAndWaitForFulfillment(
	l zerolog.Logger,
	consumer contracts.VRFv2LoadTestConsumer,
	coordinator contracts.Coordinator,
	vrfOwner contracts.VRFOwner,
	vrfv2KeyData *vrfcommon.VRFKeyData,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	subTopUpAmount *big.Int,
	linkAddress common.Address,
	randomWordsFulfilledEventTimeout time.Duration,
) (*contracts.CoordinatorConfigSet, *contracts.CoordinatorRandomWordsFulfilled, *vrf_owner.VRFOwnerRandomWordsForced, error) {
	logRandRequest(l, consumer.Address(), coordinator.Address(), 0, minimumConfirmations, callbackGasLimit, numberOfWords, randomnessRequestCountPerRequest, randomnessRequestCountPerRequestDeviation, vrfv2KeyData.KeyHash, 0)
	randomWordsRequestedEvent, err := consumer.RequestRandomWordsWithForceFulfill(
		coordinator,
		vrfv2KeyData.KeyHash,
		minimumConfirmations,
		callbackGasLimit,
		numberOfWords,
		randomnessRequestCountPerRequest,
		subTopUpAmount,
		linkAddress,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRequestRandomness, err)
	}

	vrfcommon.LogRandomnessRequestedEvent(l, coordinator, randomWordsRequestedEvent, false, 0)

	errorChannel := make(chan error)
	configSetEventChannel := make(chan *contracts.CoordinatorConfigSet)
	randWordsFulfilledEventChannel := make(chan *contracts.CoordinatorRandomWordsFulfilled)
	randWordsForcedEventChannel := make(chan *vrf_owner.VRFOwnerRandomWordsForced)

	go func() {
		configSetEvent, err := coordinator.WaitForConfigSetEvent(
			randomWordsFulfilledEventTimeout,
		)
		if err != nil {
			l.Error().Err(err).Msg("error waiting for ConfigSetEvent")
			errorChannel <- err
		}
		configSetEventChannel <- configSetEvent
	}()

	go func() {
		randomWordsFulfilledEvent, err := coordinator.WaitForRandomWordsFulfilledEvent(
			contracts.RandomWordsFulfilledEventFilter{
				RequestIds: []*big.Int{randomWordsRequestedEvent.RequestId},
				Timeout:    randomWordsFulfilledEventTimeout,
			},
		)
		if err != nil {
			l.Error().Err(err).Msg("error waiting for RandomWordsFulfilledEvent")
			errorChannel <- err
		}
		randWordsFulfilledEventChannel <- randomWordsFulfilledEvent
	}()

	go func() {
		randomWordsForcedEvent, err := vrfOwner.WaitForRandomWordsForcedEvent(
			[]*big.Int{randomWordsRequestedEvent.RequestId},
			nil,
			nil,
			randomWordsFulfilledEventTimeout,
		)
		if err != nil {
			l.Error().Err(err).Msg("error waiting for RandomWordsForcedEvent")
			errorChannel <- err
		}
		randWordsForcedEventChannel <- randomWordsForcedEvent
	}()

	var configSetEvent *contracts.CoordinatorConfigSet
	var randomWordsFulfilledEvent *contracts.CoordinatorRandomWordsFulfilled
	var randomWordsForcedEvent *vrf_owner.VRFOwnerRandomWordsForced
	for i := 0; i < 3; i++ {
		select {
		case err = <-errorChannel:
			return nil, nil, nil, err
		case configSetEvent = <-configSetEventChannel:
		case randomWordsFulfilledEvent = <-randWordsFulfilledEventChannel:
			vrfcommon.LogRandomWordsFulfilledEvent(l, coordinator, randomWordsFulfilledEvent, false, 0)
		case randomWordsForcedEvent = <-randWordsForcedEventChannel:
			vrfcommon.LogRandomWordsForcedEvent(l, vrfOwner, randomWordsForcedEvent)
		case <-time.After(randomWordsFulfilledEventTimeout):
			err = fmt.Errorf("timeout waiting for ConfigSet, RandomWordsFulfilled and RandomWordsForced events")
		}
	}
	return configSetEvent, randomWordsFulfilledEvent, randomWordsForcedEvent, err
}

func WaitRandomWordsFulfilledEvent(
	coordinator contracts.Coordinator,
	requestId *big.Int,
	randomWordsRequestedEventBlockNumber uint64,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*contracts.CoordinatorRandomWordsFulfilled, error) {
	randomWordsFulfilledEvent, err := coordinator.WaitForRandomWordsFulfilledEvent(
		contracts.RandomWordsFulfilledEventFilter{
			RequestIds: []*big.Int{requestId},
			Timeout:    randomWordsFulfilledEventTimeout,
		},
	)
	if err != nil {
		l.Warn().
			Str("requestID", requestId.String()).
			Err(err).Msg("Error waiting for random words fulfilled event, trying to filter for the event")
		randomWordsFulfilledEvent, err = coordinator.FilterRandomWordsFulfilledEvent(
			&bind.FilterOpts{
				Start: randomWordsRequestedEventBlockNumber,
			},
			requestId,
		)
		if err != nil {
			return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrFilterRandomWordsFulfilledEvent, err)
		}
	}
	vrfcommon.LogRandomWordsFulfilledEvent(l, coordinator, randomWordsFulfilledEvent, false, 0)
	return randomWordsFulfilledEvent, err
}

func SetupVRFOwnerContractIfNeeded(useVRFOwner bool, vrfContracts *vrfcommon.VRFContracts, vrfTXKeyAddressStrings []string, vrfTXKeyAddresses []common.Address, l zerolog.Logger) (*vrfcommon.VRFOwnerConfig, error) {
	var vrfOwnerConfig *vrfcommon.VRFOwnerConfig
	if useVRFOwner {
		err := setupVRFOwnerContract(vrfContracts, vrfTXKeyAddressStrings, vrfTXKeyAddresses, l)
		if err != nil {
			return nil, err
		}
		vrfOwnerConfig = &vrfcommon.VRFOwnerConfig{
			OwnerAddress: vrfContracts.VRFOwner.Address(),
			UseVRFOwner:  useVRFOwner,
		}
	} else {
		vrfOwnerConfig = &vrfcommon.VRFOwnerConfig{
			OwnerAddress: "",
			UseVRFOwner:  useVRFOwner,
		}
	}
	return vrfOwnerConfig, nil
}

func SetupNewConsumersAndSubs(
	env *test_env.CLClusterTestEnv,
	chainID int64,
	coordinator contracts.VRFCoordinatorV2,
	testConfig tc.TestConfig,
	linkToken contracts.LinkToken,
	numberOfConsumerContractsToDeployAndAddToSub int,
	numberOfSubToCreate int,
	l zerolog.Logger,
) ([]contracts.VRFv2LoadTestConsumer, []uint64, error) {
	sethClient, err := env.GetSethClient(chainID)
	if err != nil {
		return nil, nil, err
	}
	consumers, err := DeployVRFV2Consumers(sethClient, coordinator.Address(), numberOfConsumerContractsToDeployAndAddToSub)
	if err != nil {
		return nil, nil, fmt.Errorf("err: %w", err)
	}
	l.Info().
		Str("Coordinator", *testConfig.VRFv2.ExistingEnvConfig.ExistingEnvConfig.CoordinatorAddress).
		Int("Number of Subs to create", numberOfSubToCreate).
		Msg("Creating and funding subscriptions, deploying and adding consumers to subs")
	subIDs, err := CreateFundSubsAndAddConsumers(
		big.NewFloat(*testConfig.VRFv2.General.SubscriptionFundingAmountLink),
		linkToken,
		coordinator,
		consumers,
		numberOfSubToCreate,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("err: %w", err)
	}
	return consumers, subIDs, nil
}

func CancelSubsAndReturnFunds(ctx context.Context, vrfContracts *vrfcommon.VRFContracts, eoaWalletAddress string, subIDs []uint64, l zerolog.Logger) {
	for _, subID := range subIDs {
		l.Info().
			Uint64("Returning funds from SubID", subID).
			Str("Returning funds to", eoaWalletAddress).
			Msg("Canceling subscription and returning funds to subscription owner")
		pendingRequestsExist, err := vrfContracts.CoordinatorV2.PendingRequestsExist(ctx, subID)
		if err != nil {
			l.Error().Err(err).Msg("Error checking if pending requests exist")
		}
		if !pendingRequestsExist {
			_, _, err := vrfContracts.CoordinatorV2.CancelSubscription(subID, common.HexToAddress(eoaWalletAddress))
			if err != nil {
				l.Error().Err(err).Msg("Error canceling subscription")
			}
		} else {
			l.Error().Uint64("Sub ID", subID).Msg("Pending requests exist for subscription, cannot cancel subscription and return funds")
		}
	}
}
