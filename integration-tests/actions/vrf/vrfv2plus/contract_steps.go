package vrfv2plus

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	vrfv2plus_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2plus"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

func DeployVRFV2_5Contracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
) (*vrfcommon.VRFContracts, error) {
	bhs, err := contractDeployer.DeployBlockhashStore()
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrDeployBlockHashStore, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrWaitTXsComplete, err)
	}
	batchBHS, err := contractDeployer.DeployBatchBlockhashStore(bhs.Address())
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrDeployBatchBlockHashStore, err)
	}
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2_5(bhs.Address())
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrDeployCoordinatorV2Plus, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrWaitTXsComplete, err)
	}
	batchCoordinator, err := contractDeployer.DeployBatchVRFCoordinatorV2Plus(coordinator.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployBatchCoordinatorV2Plus, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrWaitTXsComplete, err)
	}
	return &vrfcommon.VRFContracts{
		CoordinatorV2Plus:      coordinator,
		BatchCoordinatorV2Plus: batchCoordinator,
		BHS:                    bhs,
		BatchBHS:               batchBHS,
		VRFV2PlusConsumer:      nil,
	}, nil
}

func DeployVRFV2PlusConsumers(contractDeployer contracts.ContractDeployer, coordinator contracts.VRFCoordinatorV2_5, consumerContractsAmount int) ([]contracts.VRFv2PlusLoadTestConsumer, error) {
	var consumers []contracts.VRFv2PlusLoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contractDeployer.DeployVRFv2PlusLoadTestConsumer(coordinator.Address())
		if err != nil {
			return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func VRFV2_5RegisterProvingKey(
	vrfKey *client.VRFKey,
	coordinator contracts.VRFCoordinatorV2_5,
	gasLaneMaxGas uint64,
) (vrfcommon.VRFEncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return vrfcommon.VRFEncodedProvingKey{}, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrEncodingProvingKey, err)
	}
	err = coordinator.RegisterProvingKey(
		provingKey,
		gasLaneMaxGas,
	)
	if err != nil {
		return vrfcommon.VRFEncodedProvingKey{}, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrRegisterProvingKey, err)
	}
	return provingKey, nil
}

func VRFV2PlusUpgradedVersionRegisterProvingKey(
	vrfKey *client.VRFKey,
	coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion,
	gasLaneMaxGasPrice uint64,
) (vrfcommon.VRFEncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return vrfcommon.VRFEncodedProvingKey{}, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrEncodingProvingKey, err)
	}
	err = coordinator.RegisterProvingKey(
		provingKey,
		gasLaneMaxGasPrice,
	)
	if err != nil {
		return vrfcommon.VRFEncodedProvingKey{}, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrRegisterProvingKey, err)
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
		return fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrABIEncodingFunding, err)
	}
	_, err = linkToken.TransferAndCall(coordinator.Address(), linkFundingAmountJuels, encodedSubId)
	if err != nil {
		return fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrSendingLinkToken, err)
	}
	return chainClient.WaitForEvents()
}

func CreateFundSubsAndAddConsumers(
	env *test_env.CLClusterTestEnv,
	chainID int64,
	subscriptionFundingAmountNative *big.Float,
	subscriptionFundingAmountLink *big.Float,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2_5,
	consumers []contracts.VRFv2PlusLoadTestConsumer,
	numberOfSubToCreate int,
) ([]*big.Int, error) {
	subIDs, err := CreateSubsAndFund(
		env,
		chainID,
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

	evmClient, err := env.GetEVMClient(chainID)
	if err != nil {
		return nil, err
	}

	err = evmClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrWaitTXsComplete, err)
	}
	return subIDs, nil
}

func CreateSubsAndFund(
	env *test_env.CLClusterTestEnv,
	chainID int64,
	subscriptionFundingAmountNative *big.Float,
	subscriptionFundingAmountLink *big.Float,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2_5,
	subAmountToCreate int,
) ([]*big.Int, error) {
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
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrWaitTXsComplete, err)
	}
	err = FundSubscriptions(
		env,
		chainID,
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
	chainID int64,
	coordinator contracts.VRFCoordinatorV2_5,
	subAmountToCreate int,
) ([]*big.Int, error) {
	var subIDArr []*big.Int

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
	subToConsumerMap map[*big.Int][]contracts.VRFv2PlusLoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2_5,
) error {
	for subID, consumers := range subToConsumerMap {
		for _, consumer := range consumers {
			err := coordinator.AddConsumer(subID, consumer.Address())
			if err != nil {
				return fmt.Errorf(vrfcommon.ErrGenericFormat, ErrAddConsumerToSub, err)
			}
		}
	}
	return nil
}

func CreateSubAndFindSubID(env *test_env.CLClusterTestEnv, chainID int64, coordinator contracts.VRFCoordinatorV2_5) (*big.Int, error) {
	tx, err := coordinator.CreateSubscription()
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrCreateVRFSubscription, err)
	}
	evmClient, err := env.GetEVMClient(chainID)
	if err != nil {
		return nil, err
	}
	err = evmClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrWaitTXsComplete, err)
	}

	receipt, err := evmClient.GetTxReceipt(tx.Hash())
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrWaitTXsComplete, err)
	}

	//SubscriptionsCreated Log should be emitted with the subscription ID
	subID := receipt.Logs[0].Topics[1].Big()

	return subID, nil
}

func FundSubscriptions(
	env *test_env.CLClusterTestEnv,
	chainID int64,
	subscriptionFundingAmountNative *big.Float,
	subscriptionFundingAmountLink *big.Float,
	linkAddress contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2_5,
	subIDs []*big.Int,
) error {
	evmClient, err := env.GetEVMClient(chainID)
	if err != nil {
		return err
	}

	for _, subID := range subIDs {
		//Native Billing
		amountWei := conversions.EtherToWei(subscriptionFundingAmountNative)
		err := coordinator.FundSubscriptionWithNative(
			subID,
			amountWei,
		)
		if err != nil {
			return fmt.Errorf(vrfcommon.ErrGenericFormat, ErrFundSubWithNativeToken, err)
		}
		//Link Billing
		amountJuels := conversions.EtherToWei(subscriptionFundingAmountLink)
		err = FundVRFCoordinatorV2_5Subscription(linkAddress, coordinator, evmClient, subID, amountJuels)
		if err != nil {
			return fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrFundSubWithLinkToken, err)
		}
	}
	err = evmClient.WaitForEvents()
	if err != nil {
		return fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrWaitTXsComplete, err)
	}
	return nil
}

func GetUpgradedCoordinatorTotalBalance(coordinator contracts.VRFCoordinatorV2PlusUpgradedVersion) (linkTotalBalance *big.Int, nativeTokenTotalBalance *big.Int, err error) {
	linkTotalBalance, err = coordinator.GetLinkTotalBalance(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrLinkTotalBalance, err)
	}
	nativeTokenTotalBalance, err = coordinator.GetNativeTokenTotalBalance(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrNativeTokenBalance, err)
	}
	return
}

func GetCoordinatorTotalBalance(coordinator contracts.VRFCoordinatorV2_5) (linkTotalBalance *big.Int, nativeTokenTotalBalance *big.Int, err error) {
	linkTotalBalance, err = coordinator.GetLinkTotalBalance(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrLinkTotalBalance, err)
	}
	nativeTokenTotalBalance, err = coordinator.GetNativeTokenTotalBalance(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrNativeTokenBalance, err)
	}
	return
}

func RequestRandomness(
	consumer contracts.VRFv2PlusLoadTestConsumer,
	coordinator contracts.Coordinator,
	vrfKeyData *vrfcommon.VRFKeyData,
	subID *big.Int,
	isNativeBilling bool,
	config *vrfv2plus_config.General,
	l zerolog.Logger,
) (*contracts.CoordinatorRandomWordsRequested, error) {
	LogRandRequest(
		l,
		consumer.Address(),
		coordinator.Address(),
		subID,
		isNativeBilling,
		vrfKeyData.KeyHash,
		config,
	)
	randomWordsRequestedEvent, err := consumer.RequestRandomness(
		coordinator,
		vrfKeyData.KeyHash,
		subID,
		*config.MinimumConfirmations,
		*config.CallbackGasLimit,
		isNativeBilling,
		*config.NumberOfWords,
		*config.RandomnessRequestCountPerRequest,
	)
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrRequestRandomness, err)
	}
	vrfcommon.LogRandomnessRequestedEvent(l, coordinator, randomWordsRequestedEvent, isNativeBilling)

	return randomWordsRequestedEvent, err
}

func RequestRandomnessAndWaitForFulfillment(
	consumer contracts.VRFv2PlusLoadTestConsumer,
	coordinator contracts.Coordinator,
	vrfKeyData *vrfcommon.VRFKeyData,
	subID *big.Int,
	isNativeBilling bool,
	config *vrfv2plus_config.General,
	l zerolog.Logger,
) (*contracts.CoordinatorRandomWordsRequested, *contracts.CoordinatorRandomWordsFulfilled, error) {
	randomWordsRequestedEvent, err := RequestRandomness(
		consumer,
		coordinator,
		vrfKeyData,
		subID,
		isNativeBilling,
		config,
		l,
	)
	if err != nil {
		return nil, nil, err
	}

	randomWordsFulfilledEvent, err := WaitRandomWordsFulfilledEvent(
		coordinator,
		randomWordsRequestedEvent.RequestId,
		subID,
		randomWordsRequestedEvent.Raw.BlockNumber,
		isNativeBilling,
		config.RandomWordsFulfilledEventTimeout.Duration,
		l,
	)
	if err != nil {
		return nil, nil, err
	}
	return randomWordsRequestedEvent, randomWordsFulfilledEvent, nil

}

func DeployVRFV2PlusDirectFundingContracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
	linkTokenAddress string,
	linkEthFeedAddress string,
	coordinator contracts.VRFCoordinatorV2_5,
	consumerContractsAmount int,
	wrapperSubId *big.Int,
) (*VRFV2PlusWrapperContracts, error) {

	vrfv2PlusWrapper, err := contractDeployer.DeployVRFV2PlusWrapper(linkTokenAddress, linkEthFeedAddress, coordinator.Address(), wrapperSubId)
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrDeployWrapper, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrWaitTXsComplete, err)
	}

	consumers, err := DeployVRFV2PlusWrapperConsumers(contractDeployer, vrfv2PlusWrapper, consumerContractsAmount)
	if err != nil {
		return nil, err
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrWaitTXsComplete, err)
	}
	return &VRFV2PlusWrapperContracts{vrfv2PlusWrapper, consumers}, nil
}

func WrapperRequestRandomness(consumer contracts.VRFv2PlusWrapperLoadTestConsumer, coordinator contracts.Coordinator, vrfKeyData *vrfcommon.VRFKeyData, subID *big.Int, isNativeBilling bool, config *vrfv2plus_config.General, l zerolog.Logger) (*contracts.CoordinatorRandomWordsRequested, string, error) {
	LogRandRequest(
		l,
		consumer.Address(),
		coordinator.Address(),
		subID,
		isNativeBilling,
		vrfKeyData.KeyHash,
		config,
	)
	var randomWordsRequestedEvent *contracts.CoordinatorRandomWordsRequested
	var err error
	if isNativeBilling {
		randomWordsRequestedEvent, err = consumer.RequestRandomnessNative(
			coordinator,
			*config.MinimumConfirmations,
			*config.CallbackGasLimit,
			*config.NumberOfWords,
			*config.RandomnessRequestCountPerRequest,
		)
		if err != nil {
			return nil, "", fmt.Errorf(vrfcommon.ErrGenericFormat, ErrRequestRandomnessDirectFundingNativePayment, err)
		}
	} else {
		randomWordsRequestedEvent, err = consumer.RequestRandomness(
			coordinator,
			*config.MinimumConfirmations,
			*config.CallbackGasLimit,
			*config.NumberOfWords,
			*config.RandomnessRequestCountPerRequest,
		)
		if err != nil {
			return nil, "", fmt.Errorf(vrfcommon.ErrGenericFormat, ErrRequestRandomnessDirectFundingLinkPayment, err)
		}
	}
	vrfcommon.LogRandomnessRequestedEvent(l, coordinator, randomWordsRequestedEvent, isNativeBilling)
	wrapperAddress, err := consumer.GetWrapper(context.Background())
	if err != nil {
		return nil, "", fmt.Errorf("error getting wrapper address, err: %w", err)
	}
	return randomWordsRequestedEvent, wrapperAddress.Hex(), nil
}

func DirectFundingRequestRandomnessAndWaitForFulfillment(
	consumer contracts.VRFv2PlusWrapperLoadTestConsumer,
	coordinator contracts.Coordinator,
	vrfKeyData *vrfcommon.VRFKeyData,
	subID *big.Int,
	isNativeBilling bool,
	config *vrfv2plus_config.General,
	l zerolog.Logger,
) (*contracts.CoordinatorRandomWordsFulfilled, error) {
	randomWordsRequestedEvent, _, err := WrapperRequestRandomness(consumer, coordinator, vrfKeyData, subID,
		isNativeBilling, config, l)
	if err != nil {
		return nil, fmt.Errorf("error getting wrapper address, err: %w", err)
	}
	return WaitRandomWordsFulfilledEvent(
		coordinator,
		randomWordsRequestedEvent.RequestId,
		subID,
		randomWordsRequestedEvent.Raw.BlockNumber,
		isNativeBilling,
		config.RandomWordsFulfilledEventTimeout.Duration,
		l,
	)
}

func WaitRandomWordsFulfilledEvent(
	coordinator contracts.Coordinator,
	requestId *big.Int,
	subID *big.Int,
	randomWordsRequestedEventBlockNumber uint64,
	isNativeBilling bool,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*contracts.CoordinatorRandomWordsFulfilled, error) {
	randomWordsFulfilledEvent, err := coordinator.WaitForRandomWordsFulfilledEvent(
		contracts.RandomWordsFulfilledEventFilter{
			SubIDs:     []*big.Int{subID},
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
	vrfcommon.LogRandomWordsFulfilledEvent(l, coordinator, randomWordsFulfilledEvent, isNativeBilling)
	return randomWordsFulfilledEvent, err
}

func DeployVRFV2PlusWrapperConsumers(contractDeployer contracts.ContractDeployer, vrfV2PlusWrapper contracts.VRFV2PlusWrapper, consumerContractsAmount int) ([]contracts.VRFv2PlusWrapperLoadTestConsumer, error) {
	var consumers []contracts.VRFv2PlusWrapperLoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contractDeployer.DeployVRFV2PlusWrapperLoadTestConsumer(vrfV2PlusWrapper.Address())
		if err != nil {
			return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func SetupVRFV2PlusContracts(
	env *test_env.CLClusterTestEnv,
	chainID int64,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.VRFMockETHLINKFeed,
	configGeneral *vrfv2plus_config.General,
	l zerolog.Logger,
) (*vrfcommon.VRFContracts, error) {
	l.Info().Msg("Deploying VRFV2 Plus contracts")
	evmClient, err := env.GetEVMClient(chainID)
	if err != nil {
		return nil, err
	}
	vrfContracts, err := DeployVRFV2_5Contracts(env.ContractDeployer, evmClient)
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrDeployVRFV2_5Contracts, err)
	}
	vrfContracts.LinkToken = linkToken
	vrfContracts.MockETHLINKFeed = mockNativeLINKFeed

	l.Info().Str("Coordinator", vrfContracts.CoordinatorV2Plus.Address()).Msg("Setting Coordinator Config")
	err = vrfContracts.CoordinatorV2Plus.SetConfig(
		*configGeneral.MinimumConfirmations,
		*configGeneral.MaxGasLimitCoordinatorConfig,
		*configGeneral.StalenessSeconds,
		*configGeneral.GasAfterPaymentCalculation,
		decimal.RequireFromString(*configGeneral.FallbackWeiPerUnitLink).BigInt(),
		*configGeneral.FulfillmentFlatFeeNativePPM,
		*configGeneral.FulfillmentFlatFeeLinkDiscountPPM,
		*configGeneral.NativePremiumPercentage,
		*configGeneral.LinkPremiumPercentage,
	)
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrSetVRFCoordinatorConfig, err)
	}

	l.Info().Str("Coordinator", vrfContracts.CoordinatorV2Plus.Address()).Msg("Setting Link and ETH/LINK feed")
	err = vrfContracts.CoordinatorV2Plus.SetLINKAndLINKNativeFeed(linkToken.Address(), mockNativeLINKFeed.Address())
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrSetLinkNativeLinkFeed, err)
	}
	err = evmClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrWaitTXsComplete, err)
	}

	return vrfContracts, nil
}

func SetupNewConsumersAndSubs(
	env *test_env.CLClusterTestEnv,
	chainID int64,
	coordinator contracts.VRFCoordinatorV2_5,
	testConfig tc.TestConfig,
	linkToken contracts.LinkToken,
	consumerContractsAmount int,
	numberOfSubToCreate int,
	l zerolog.Logger,
) ([]contracts.VRFv2PlusLoadTestConsumer, []*big.Int, error) {
	consumers, err := DeployVRFV2PlusConsumers(env.ContractDeployer, coordinator, consumerContractsAmount)
	if err != nil {
		return nil, nil, fmt.Errorf("err: %w", err)
	}
	evmClient, err := env.GetEVMClient(chainID)
	if err != nil {
		return nil, nil, err
	}
	err = evmClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err: %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	l.Info().
		Str("Coordinator", *testConfig.VRFv2Plus.ExistingEnvConfig.ExistingEnvConfig.CoordinatorAddress).
		Int("Number of Subs to create", numberOfSubToCreate).
		Msg("Creating and funding subscriptions, deploying and adding consumers to subs")
	subIDs, err := CreateFundSubsAndAddConsumers(
		env,
		chainID,
		big.NewFloat(*testConfig.VRFv2Plus.General.SubscriptionFundingAmountNative),
		big.NewFloat(*testConfig.VRFv2Plus.General.SubscriptionFundingAmountLink),
		linkToken,
		coordinator,
		consumers,
		*testConfig.VRFv2Plus.General.NumberOfSubToCreate,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("err: %w", err)
	}
	return consumers, subIDs, nil
}

func CancelSubsAndReturnFunds(ctx context.Context, vrfContracts *vrfcommon.VRFContracts, eoaWalletAddress string, subIDs []*big.Int, l zerolog.Logger) {
	for _, subID := range subIDs {
		l.Info().
			Str("Returning funds from SubID", subID.String()).
			Str("Returning funds to", eoaWalletAddress).
			Msg("Canceling subscription and returning funds to subscription owner")
		pendingRequestsExist, err := vrfContracts.CoordinatorV2Plus.PendingRequestsExist(ctx, subID)
		if err != nil {
			l.Error().Err(err).Msg("Error checking if pending requests exist")
		}
		if !pendingRequestsExist {
			_, err := vrfContracts.CoordinatorV2Plus.CancelSubscription(subID, common.HexToAddress(eoaWalletAddress))
			if err != nil {
				l.Error().Err(err).Msg("Error canceling subscription")
			}
		} else {
			l.Error().Str("Sub ID", subID.String()).Msg("Pending requests exist for subscription, cannot cancel subscription and return funds")
		}
	}
}
