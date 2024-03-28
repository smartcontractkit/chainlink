package vrfv2

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
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
	bhs, err := env.ContractDeployer.DeployBlockhashStore()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrDeployBlockHashStore, err)
	}

	evmClient, err := env.GetEVMClient(chainID)
	if err != nil {
		return nil, err
	}

	err = evmClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	var coordinatorAddress string
	if useTestCoordinator {
		testCoordinator, err := env.ContractDeployer.DeployVRFCoordinatorTestV2(linkTokenContract.Address(), bhs.Address(), linkEthFeedContract.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrDeployCoordinator, err)
		}
		err = evmClient.WaitForEvents()
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
		}
		coordinatorAddress = testCoordinator.Address()
	} else {
		coordinator, err := env.ContractDeployer.DeployVRFCoordinatorV2(linkTokenContract.Address(), bhs.Address(), linkEthFeedContract.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrDeployCoordinator, err)
		}
		err = evmClient.WaitForEvents()
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
		}
		coordinatorAddress = coordinator.Address()
	}

	coordinator, err := env.ContractLoader.LoadVRFCoordinatorV2(coordinatorAddress)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrLoadingCoordinator, err)
	}
	if useVRFOwner {
		vrfOwner, err := env.ContractDeployer.DeployVRFOwner(coordinatorAddress)
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrDeployCoordinator, err)
		}
		err = evmClient.WaitForEvents()
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
		}
		return &vrfcommon.VRFContracts{
			CoordinatorV2:   coordinator,
			VRFOwner:        vrfOwner,
			BHS:             bhs,
			VRFV2Consumers:  nil,
			LinkToken:       linkTokenContract,
			MockETHLINKFeed: linkEthFeedContract,
		}, nil
	}
	return &vrfcommon.VRFContracts{
		CoordinatorV2:   coordinator,
		VRFOwner:        nil,
		BHS:             bhs,
		VRFV2Consumers:  nil,
		LinkToken:       linkTokenContract,
		MockETHLINKFeed: linkEthFeedContract,
	}, nil
}

func DeployVRFV2Consumers(contractDeployer contracts.ContractDeployer, coordinatorAddress string, consumerContractsAmount int) ([]contracts.VRFv2LoadTestConsumer, error) {
	var consumers []contracts.VRFv2LoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contractDeployer.DeployVRFv2LoadTestConsumer(coordinatorAddress)
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func DeployVRFV2WrapperConsumers(contractDeployer contracts.ContractDeployer, linkTokenAddress string, vrfV2Wrapper contracts.VRFV2Wrapper, consumerContractsAmount int) ([]contracts.VRFv2WrapperLoadTestConsumer, error) {
	var consumers []contracts.VRFv2WrapperLoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contractDeployer.DeployVRFV2WrapperLoadTestConsumer(linkTokenAddress, vrfV2Wrapper.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func DeployVRFV2DirectFundingContracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
	linkTokenAddress string,
	linkEthFeedAddress string,
	coordinator contracts.VRFCoordinatorV2,
	consumerContractsAmount int,
) (*VRFV2WrapperContracts, error) {
	vrfv2Wrapper, err := contractDeployer.DeployVRFV2Wrapper(linkTokenAddress, linkEthFeedAddress, coordinator.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployVRFV2Wrapper, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	consumers, err := DeployVRFV2WrapperConsumers(contractDeployer, linkTokenAddress, vrfv2Wrapper, consumerContractsAmount)
	if err != nil {
		return nil, err
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
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
		big.NewInt(*vrfv2Config.FallbackWeiPerUnitLink),
		vrfCoordinatorV2FeeConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrSetVRFCoordinatorConfig, err)
	}

	evmClient, err := env.GetEVMClient(chainID)
	if err != nil {
		return nil, err
	}

	err = evmClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return vrfContracts, nil
}

func setupVRFOwnerContract(env *test_env.CLClusterTestEnv, chainID int64, contracts *vrfcommon.VRFContracts, allNativeTokenKeyAddressStrings []string, allNativeTokenKeyAddresses []common.Address, l zerolog.Logger) error {
	l.Info().Msg("Setting up VRFOwner contract")
	l.Info().
		Str("Coordinator", contracts.CoordinatorV2.Address()).
		Str("VRFOwner", contracts.VRFOwner.Address()).
		Msg("Transferring ownership of Coordinator to VRFOwner")
	err := contracts.CoordinatorV2.TransferOwnership(common.HexToAddress(contracts.VRFOwner.Address()))
	if err != nil {
		return nil
	}
	evmClient, err := env.GetEVMClient(chainID)
	if err != nil {
		return err
	}

	err = evmClient.WaitForEvents()
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
	err = evmClient.WaitForEvents()
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
	err = evmClient.WaitForEvents()
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return err
}

func CreateFundSubsAndAddConsumers(
	env *test_env.CLClusterTestEnv,
	chainID int64,
	subscriptionFundingAmountLink *big.Float,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	consumers []contracts.VRFv2LoadTestConsumer,
	numberOfSubToCreate int,
) ([]uint64, error) {
	subIDs, err := CreateSubsAndFund(env, chainID, subscriptionFundingAmountLink, linkToken, coordinator, numberOfSubToCreate)
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

	evmClient, err := env.GetEVMClient(chainID)
	if err != nil {
		return nil, err
	}

	err = evmClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return subIDs, nil
}

func CreateSubsAndFund(
	env *test_env.CLClusterTestEnv,
	chainID int64,
	subscriptionFundingAmountLink *big.Float,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	subAmountToCreate int,
) ([]uint64, error) {
	subs, err := CreateSubs(env, chainID, coordinator, subAmountToCreate)
	if err != nil {
		return nil, err
	}
	evmClient, err := env.GetEVMClient(chainID)
	if err != nil {
		return nil, err
	}

	err = evmClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	err = FundSubscriptions(env, chainID, subscriptionFundingAmountLink, linkToken, coordinator, subs)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func CreateSubs(
	env *test_env.CLClusterTestEnv,
	chainID int64,
	coordinator contracts.VRFCoordinatorV2,
	subAmountToCreate int,
) ([]uint64, error) {
	var subIDArr []uint64

	for i := 0; i < subAmountToCreate; i++ {
		subID, err := CreateSubAndFindSubID(env, chainID, coordinator)
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

func CreateSubAndFindSubID(env *test_env.CLClusterTestEnv, chainID int64, coordinator contracts.VRFCoordinatorV2) (uint64, error) {
	tx, err := coordinator.CreateSubscription()
	if err != nil {
		return 0, fmt.Errorf("%s, err %w", vrfcommon.ErrCreateVRFSubscription, err)
	}
	evmClient, err := env.GetEVMClient(chainID)
	if err != nil {
		return 0, err
	}

	err = evmClient.WaitForEvents()
	if err != nil {
		return 0, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	receipt, err := evmClient.GetTxReceipt(tx.Hash())
	if err != nil {
		return 0, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	//SubscriptionsCreated Log should be emitted with the subscription ID
	subID := receipt.Logs[0].Topics[1].Big().Uint64()

	return subID, nil
}

func FundSubscriptions(
	env *test_env.CLClusterTestEnv,
	chainID int64,
	subscriptionFundingAmountLink *big.Float,
	linkAddress contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	subIDs []uint64,
) error {
	evmClient, err := env.GetEVMClient(chainID)
	if err != nil {
		return err
	}

	for _, subID := range subIDs {
		//Link Billing
		amountJuels := conversions.EtherToWei(subscriptionFundingAmountLink)
		err := FundVRFCoordinatorV2Subscription(linkAddress, coordinator, evmClient, subID, amountJuels)
		if err != nil {
			return fmt.Errorf("%s, err %w", vrfcommon.ErrFundSubWithLinkToken, err)
		}
	}
	err = evmClient.WaitForEvents()
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return nil
}

func FundVRFCoordinatorV2Subscription(
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	chainClient blockchain.EVMClient,
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
	return chainClient.WaitForEvents()
}

func DirectFundingRequestRandomnessAndWaitForFulfillment(
	l zerolog.Logger,
	consumer contracts.VRFv2WrapperLoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2,
	subID uint64,
	vrfv2KeyData *vrfcommon.VRFKeyData,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	randomWordsFulfilledEventTimeout time.Duration,
) (*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled, error) {
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
		randomWordsFulfilledEventTimeout,
		l,
	)
	return fulfillmentEvents, err
}

func RequestRandomnessAndWaitForFulfillment(
	l zerolog.Logger,
	consumer contracts.VRFv2LoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2,
	subID uint64,
	vrfKeyData *vrfcommon.VRFKeyData,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	randomWordsFulfilledEventTimeout time.Duration,
) (*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled, error) {
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
	)
	if err != nil {
		return nil, err
	}
	fulfillmentEvents, err := WaitRandomWordsFulfilledEvent(
		coordinator,
		randomWordsRequestedEvent.RequestId,
		randomWordsFulfilledEventTimeout,
		l,
	)
	if err != nil {
		return nil, err
	}
	return fulfillmentEvents, nil
}

func RequestRandomness(
	l zerolog.Logger,
	consumer contracts.VRFv2LoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2,
	subID uint64,
	vrfKeyData *vrfcommon.VRFKeyData,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
) (*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested, error) {
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
	)
	randomWordsRequestedEvent, err := consumer.RequestRandomness(
		coordinator,
		vrfKeyData.KeyHash,
		subID,
		minimumConfirmations,
		callbackGasLimit,
		numberOfWords,
		randomnessRequestCountPerRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRequestRandomness, err)
	}
	LogRandomnessRequestedEvent(l, coordinator, randomWordsRequestedEvent)

	return randomWordsRequestedEvent, err
}

func RequestRandomnessWithForceFulfillAndWaitForFulfillment(
	l zerolog.Logger,
	consumer contracts.VRFv2LoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2,
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
) (*vrf_coordinator_v2.VRFCoordinatorV2ConfigSet, *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled, *vrf_owner.VRFOwnerRandomWordsForced, error) {
	logRandRequest(l, consumer.Address(), coordinator.Address(), 0, minimumConfirmations, callbackGasLimit, numberOfWords, randomnessRequestCountPerRequest, randomnessRequestCountPerRequestDeviation, vrfv2KeyData.KeyHash)
	_, err := consumer.RequestRandomWordsWithForceFulfill(
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

	randomWordsRequestedEvent, err := coordinator.WaitForRandomWordsRequestedEvent(
		[][32]byte{vrfv2KeyData.KeyHash},
		nil,
		[]common.Address{common.HexToAddress(consumer.Address())},
		time.Minute*1,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitRandomWordsRequestedEvent, err)
	}
	LogRandomnessRequestedEvent(l, coordinator, randomWordsRequestedEvent)

	errorChannel := make(chan error)
	configSetEventChannel := make(chan *vrf_coordinator_v2.VRFCoordinatorV2ConfigSet)
	randWordsFulfilledEventChannel := make(chan *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled)
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
			[]*big.Int{randomWordsRequestedEvent.RequestId},
			randomWordsFulfilledEventTimeout,
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

	var configSetEvent *vrf_coordinator_v2.VRFCoordinatorV2ConfigSet
	var randomWordsFulfilledEvent *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled
	var randomWordsForcedEvent *vrf_owner.VRFOwnerRandomWordsForced
	for i := 0; i < 3; i++ {
		select {
		case err = <-errorChannel:
			return nil, nil, nil, err
		case configSetEvent = <-configSetEventChannel:
		case randomWordsFulfilledEvent = <-randWordsFulfilledEventChannel:
			LogRandomWordsFulfilledEvent(l, coordinator, randomWordsFulfilledEvent)
		case randomWordsForcedEvent = <-randWordsForcedEventChannel:
			LogRandomWordsForcedEvent(l, vrfOwner, randomWordsForcedEvent)
		case <-time.After(randomWordsFulfilledEventTimeout):
			err = fmt.Errorf("timeout waiting for ConfigSet, RandomWordsFulfilled and RandomWordsForced events")
		}
	}
	return configSetEvent, randomWordsFulfilledEvent, randomWordsForcedEvent, err
}

func WaitRandomWordsFulfilledEvent(
	coordinator contracts.VRFCoordinatorV2,
	requestId *big.Int,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled, error) {
	randomWordsFulfilledEvent, err := coordinator.WaitForRandomWordsFulfilledEvent(
		[]*big.Int{requestId},
		randomWordsFulfilledEventTimeout,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitRandomWordsFulfilledEvent, err)
	}
	LogRandomWordsFulfilledEvent(l, coordinator, randomWordsFulfilledEvent)
	return randomWordsFulfilledEvent, err
}

func SetupVRFOwnerContractIfNeeded(useVRFOwner bool, env *test_env.CLClusterTestEnv, chainID int64, vrfContracts *vrfcommon.VRFContracts, vrfTXKeyAddressStrings []string, vrfTXKeyAddresses []common.Address, l zerolog.Logger) (*vrfcommon.VRFOwnerConfig, error) {
	var vrfOwnerConfig *vrfcommon.VRFOwnerConfig
	if useVRFOwner {
		err := setupVRFOwnerContract(env, chainID, vrfContracts, vrfTXKeyAddressStrings, vrfTXKeyAddresses, l)
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
	consumers, err := DeployVRFV2Consumers(env.ContractDeployer, coordinator.Address(), numberOfConsumerContractsToDeployAndAddToSub)
	if err != nil {
		return nil, nil, fmt.Errorf("err: %w", err)
	}
	evmClient, err := env.GetEVMClient(chainID)
	if err != nil {
		return nil, []uint64{}, err
	}
	err = evmClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err: %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	l.Info().
		Str("Coordinator", *testConfig.VRFv2.ExistingEnvConfig.ExistingEnvConfig.CoordinatorAddress).
		Int("Number of Subs to create", numberOfSubToCreate).
		Msg("Creating and funding subscriptions, deploying and adding consumers to subs")
	subIDs, err := CreateFundSubsAndAddConsumers(
		env,
		chainID,
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
			_, err := vrfContracts.CoordinatorV2.CancelSubscription(subID, common.HexToAddress(eoaWalletAddress))
			if err != nil {
				l.Error().Err(err).Msg("Error canceling subscription")
			}
		} else {
			l.Error().Uint64("Sub ID", subID).Msg("Pending requests exist for subscription, cannot cancel subscription and return funds")
		}
	}
}
