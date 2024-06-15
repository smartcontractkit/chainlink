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

	"github.com/smartcontractkit/seth"

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
	chainClient *seth.Client,
) (*vrfcommon.VRFContracts, error) {
	bhs, err := contracts.DeployBlockhashStore(chainClient)
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrDeployBlockHashStore, err)
	}
	batchBHS, err := contracts.DeployBatchBlockhashStore(chainClient, bhs.Address())
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrDeployBatchBlockHashStore, err)
	}
	coordinator, err := contracts.DeployVRFCoordinatorV2_5(chainClient, bhs.Address())
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrDeployCoordinatorV2Plus, err)
	}
	batchCoordinator, err := contracts.DeployBatchVRFCoordinatorV2Plus(chainClient, coordinator.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployBatchCoordinatorV2Plus, err)
	}
	return &vrfcommon.VRFContracts{
		CoordinatorV2Plus:      coordinator,
		BatchCoordinatorV2Plus: batchCoordinator,
		BHS:                    bhs,
		BatchBHS:               batchBHS,
		VRFV2PlusConsumer:      nil,
	}, nil
}

func DeployVRFV2PlusConsumers(client *seth.Client, coordinator contracts.VRFCoordinatorV2_5, consumerContractsAmount int) ([]contracts.VRFv2PlusLoadTestConsumer, error) {
	var consumers []contracts.VRFv2PlusLoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contracts.DeployVRFv2PlusLoadTestConsumer(client, coordinator.Address())
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
	return nil
}

func CreateFundSubsAndAddConsumers(
	ctx context.Context,
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
		ctx,
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

	return subIDs, err
}

func CreateSubsAndFund(
	ctx context.Context,
	env *test_env.CLClusterTestEnv,
	chainID int64,
	subscriptionFundingAmountNative *big.Float,
	subscriptionFundingAmountLink *big.Float,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2_5,
	subAmountToCreate int,
) ([]*big.Int, error) {
	subs, err := CreateSubs(ctx, env, chainID, coordinator, subAmountToCreate)
	if err != nil {
		return nil, err
	}
	err = FundSubscriptions(
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
	ctx context.Context,
	env *test_env.CLClusterTestEnv,
	chainID int64,
	coordinator contracts.VRFCoordinatorV2_5,
	subAmountToCreate int,
) ([]*big.Int, error) {
	var subIDArr []*big.Int

	for i := 0; i < subAmountToCreate; i++ {
		subID, err := CreateSubAndFindSubID(ctx, env, chainID, coordinator)
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

func CreateSubAndFindSubID(ctx context.Context, env *test_env.CLClusterTestEnv, chainID int64, coordinator contracts.VRFCoordinatorV2_5) (*big.Int, error) {
	tx, err := coordinator.CreateSubscription()
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrCreateVRFSubscription, err)
	}
	sethClient, err := env.GetSethClient(chainID)
	if err != nil {
		return nil, err
	}
	receipt, err := sethClient.Client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrWaitTXsComplete, err)
	}

	//SubscriptionsCreated Log should be emitted with the subscription ID
	subID := receipt.Logs[0].Topics[1].Big()

	return subID, nil
}

func FundSubscriptions(
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
			return fmt.Errorf(vrfcommon.ErrGenericFormat, ErrFundSubWithNativeToken, err)
		}
		//Link Billing
		amountJuels := conversions.EtherToWei(subscriptionFundingAmountLink)
		err = FundVRFCoordinatorV2_5Subscription(linkAddress, coordinator, subID, amountJuels)
		if err != nil {
			return fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrFundSubWithLinkToken, err)
		}
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
	keyNum int,
) (*contracts.CoordinatorRandomWordsRequested, error) {
	LogRandRequest(
		l,
		consumer.Address(),
		coordinator.Address(),
		subID,
		isNativeBilling,
		vrfKeyData.KeyHash,
		config,
		keyNum,
	)
	randomWordsRequestedEvent, err := consumer.RequestRandomnessFromKey(
		coordinator,
		vrfKeyData.KeyHash,
		subID,
		*config.MinimumConfirmations,
		*config.CallbackGasLimit,
		isNativeBilling,
		*config.NumberOfWords,
		*config.RandomnessRequestCountPerRequest,
		keyNum,
	)
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrRequestRandomness, err)
	}
	vrfcommon.LogRandomnessRequestedEvent(l, coordinator, randomWordsRequestedEvent, isNativeBilling, keyNum)

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
	keyNum int,
) (*contracts.CoordinatorRandomWordsRequested, *contracts.CoordinatorRandomWordsFulfilled, error) {
	randomWordsRequestedEvent, err := RequestRandomness(
		consumer,
		coordinator,
		vrfKeyData,
		subID,
		isNativeBilling,
		config,
		l,
		keyNum,
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
		keyNum,
	)
	if err != nil {
		return nil, nil, err
	}
	return randomWordsRequestedEvent, randomWordsFulfilledEvent, nil

}

func DeployVRFV2PlusDirectFundingContracts(
	sethClient *seth.Client,
	linkTokenAddress string,
	linkEthFeedAddress string,
	coordinator contracts.VRFCoordinatorV2_5,
	consumerContractsAmount int,
	wrapperSubId *big.Int,
) (*VRFV2PlusWrapperContracts, error) {
	vrfv2PlusWrapper, err := contracts.DeployVRFV2PlusWrapper(sethClient, linkTokenAddress, linkEthFeedAddress, coordinator.Address(), wrapperSubId)
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrDeployWrapper, err)
	}
	consumers, err := DeployVRFV2PlusWrapperConsumers(sethClient, vrfv2PlusWrapper, consumerContractsAmount)
	if err != nil {
		return nil, err
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
		0,
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
	vrfcommon.LogRandomnessRequestedEvent(l, coordinator, randomWordsRequestedEvent, isNativeBilling, 0)
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
		0,
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
	keyNum int,
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

	vrfcommon.LogRandomWordsFulfilledEvent(l, coordinator, randomWordsFulfilledEvent, isNativeBilling, keyNum)
	return randomWordsFulfilledEvent, err
}

func DeployVRFV2PlusWrapperConsumers(client *seth.Client, vrfV2PlusWrapper contracts.VRFV2PlusWrapper, consumerContractsAmount int) ([]contracts.VRFv2PlusWrapperLoadTestConsumer, error) {
	var consumers []contracts.VRFv2PlusWrapperLoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contracts.DeployVRFV2PlusWrapperLoadTestConsumer(client, vrfV2PlusWrapper.Address())
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
	sethClient, err := env.GetSethClient(chainID)
	if err != nil {
		return nil, err
	}
	vrfContracts, err := DeployVRFV2_5Contracts(sethClient)
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

	return vrfContracts, nil
}

func SetupNewConsumersAndSubs(
	ctx context.Context,
	env *test_env.CLClusterTestEnv,
	chainID int64,
	coordinator contracts.VRFCoordinatorV2_5,
	testConfig tc.TestConfig,
	linkToken contracts.LinkToken,
	consumerContractsAmount int,
	numberOfSubToCreate int,
	l zerolog.Logger,
) ([]contracts.VRFv2PlusLoadTestConsumer, []*big.Int, error) {
	sethClient, err := env.GetSethClient(chainID)
	if err != nil {
		return nil, nil, err
	}
	consumers, err := DeployVRFV2PlusConsumers(sethClient, coordinator, consumerContractsAmount)
	if err != nil {
		return nil, nil, fmt.Errorf("err: %w", err)
	}
	l.Info().
		Str("Coordinator", *testConfig.VRFv2Plus.ExistingEnvConfig.ExistingEnvConfig.CoordinatorAddress).
		Int("Number of Subs to create", numberOfSubToCreate).
		Msg("Creating and funding subscriptions, deploying and adding consumers to subs")
	subIDs, err := CreateFundSubsAndAddConsumers(
		ctx,
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
			_, _, err := vrfContracts.CoordinatorV2Plus.CancelSubscription(subID, common.HexToAddress(eoaWalletAddress))
			if err != nil {
				l.Error().Err(err).Msg("Error canceling subscription")
			}
		} else {
			l.Error().Str("Sub ID", subID.String()).Msg("Pending requests exist for subscription, cannot cancel subscription and return funds")
		}
	}
}
