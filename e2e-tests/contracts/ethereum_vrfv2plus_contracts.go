package contracts

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/montanaflynn/stats"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/e2e-tests/wrappers"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_load_test_with_metrics"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_upgraded_version"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_wrapper_load_test_consumer"
)

type EthereumVRFCoordinatorV2_5 struct {
	address     *common.Address
	client      blockchain.EVMClient
	coordinator vrf_coordinator_v2_5.VRFCoordinatorV25Interface
}

type EthereumVRFCoordinatorV2PlusUpgradedVersion struct {
	address     *common.Address
	client      blockchain.EVMClient
	coordinator *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersion
}

// EthereumVRFv2PlusLoadTestConsumer represents VRFv2Plus consumer contract for performing Load Tests
type EthereumVRFv2PlusLoadTestConsumer struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *vrf_v2plus_load_test_with_metrics.VRFV2PlusLoadTestWithMetrics
}

type EthereumVRFV2PlusWrapperLoadTestConsumer struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *vrfv2plus_wrapper_load_test_consumer.VRFV2PlusWrapperLoadTestConsumer
}

type EthereumVRFV2PlusWrapper struct {
	address *common.Address
	client  blockchain.EVMClient
	wrapper *vrfv2plus_wrapper.VRFV2PlusWrapper
}

func (v *EthereumVRFV2PlusWrapper) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFV2PlusWrapper) SetConfig(wrapperGasOverhead uint32,
	coordinatorGasOverheadNative uint32,
	coordinatorGasOverheadLink uint32,
	coordinatorGasOverheadPerWord uint16,
	wrapperNativePremiumPercentage uint8,
	wrapperLinkPremiumPercentage uint8,
	keyHash [32]byte,
	maxNumWords uint8,
	stalenessSeconds uint32,
	fallbackWeiPerUnitLink *big.Int,
	fulfillmentFlatFeeNativePPM uint32,
	fulfillmentFlatFeeLinkDiscountPPM uint32,
) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.wrapper.SetConfig(
		opts,
		wrapperGasOverhead,
		coordinatorGasOverheadNative,
		coordinatorGasOverheadLink,
		coordinatorGasOverheadPerWord,
		wrapperNativePremiumPercentage,
		wrapperLinkPremiumPercentage,
		keyHash,
		maxNumWords,
		stalenessSeconds,
		fallbackWeiPerUnitLink,
		fulfillmentFlatFeeNativePPM,
		fulfillmentFlatFeeLinkDiscountPPM,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFV2PlusWrapper) GetSubID(ctx context.Context) (*big.Int, error) {
	return v.wrapper.SUBSCRIPTIONID(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
}

func (v *EthereumVRFV2PlusWrapper) Coordinator(ctx context.Context) (common.Address, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	return v.wrapper.SVrfCoordinator(opts)
}

// DeployVRFCoordinatorV2_5 deploys VRFV2_5 coordinator contract
func (e *EthereumContractDeployer) DeployVRFCoordinatorV2_5(bhsAddr string) (VRFCoordinatorV2_5, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinatorV2Plus", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_coordinator_v2_5.DeployVRFCoordinatorV25(auth, wrappers.MustNewWrappedContractBackend(e.client, nil), common.HexToAddress(bhsAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinatorV2_5{
		client:      e.client,
		coordinator: instance.(*vrf_coordinator_v2_5.VRFCoordinatorV25),
		address:     address,
	}, err
}

func (v *EthereumVRFCoordinatorV2_5) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFCoordinatorV2_5) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	hash, err := v.coordinator.HashOfKey(opts, pubKey)
	if err != nil {
		return [32]byte{}, err
	}
	return hash, nil
}

func (v *EthereumVRFCoordinatorV2_5) GetActiveSubscriptionIds(ctx context.Context, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	activeSubscriptionIds, err := v.coordinator.GetActiveSubscriptionIds(opts, startIndex, maxCount)
	if err != nil {
		return nil, err
	}
	return activeSubscriptionIds, nil
}

func (v *EthereumVRFCoordinatorV2_5) PendingRequestsExist(ctx context.Context, subID *big.Int) (bool, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	pendingRequestExists, err := v.coordinator.PendingRequestExists(opts, subID)
	if err != nil {
		return false, err
	}
	return pendingRequestExists, nil
}

func (v *EthereumVRFCoordinatorV2_5) ParseRandomWordsRequested(log types.Log) (*CoordinatorRandomWordsRequested, error) {
	randomWordsRequested, err := v.coordinator.ParseRandomWordsRequested(log)
	if err != nil {
		return nil, fmt.Errorf("parse RandomWordsRequested log failed, err: %w", err)
	}
	coordinatorRandomWordsRequested := &CoordinatorRandomWordsRequested{
		KeyHash:                     randomWordsRequested.KeyHash,
		RequestId:                   randomWordsRequested.RequestId,
		PreSeed:                     randomWordsRequested.PreSeed,
		SubId:                       randomWordsRequested.SubId.String(),
		MinimumRequestConfirmations: randomWordsRequested.MinimumRequestConfirmations,
		CallbackGasLimit:            randomWordsRequested.CallbackGasLimit,
		NumWords:                    randomWordsRequested.NumWords,
		ExtraArgs:                   randomWordsRequested.ExtraArgs,
		Sender:                      randomWordsRequested.Sender,
		Raw:                         randomWordsRequested.Raw,
	}
	return coordinatorRandomWordsRequested, nil
}

func (v *EthereumVRFCoordinatorV2_5) GetSubscription(ctx context.Context, subID *big.Int) (Subscription, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	subscription, err := v.coordinator.GetSubscription(opts, subID)
	if err != nil {
		return Subscription{}, err
	}
	return Subscription{
		Balance:       subscription.Balance,
		NativeBalance: subscription.NativeBalance,
		SubOwner:      subscription.SubOwner,
		Consumers:     subscription.Consumers,
		ReqCount:      subscription.ReqCount,
	}, nil
}

func (v *EthereumVRFCoordinatorV2_5) GetLinkTotalBalance(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	totalBalance, err := v.coordinator.STotalBalance(opts)
	if err != nil {
		return nil, err
	}
	return totalBalance, nil
}
func (v *EthereumVRFCoordinatorV2_5) GetNativeTokenTotalBalance(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	totalBalance, err := v.coordinator.STotalNativeBalance(opts)
	if err != nil {
		return nil, err
	}
	return totalBalance, nil
}

// OwnerCancelSubscription cancels subscription by Coordinator owner
// return funds to sub owner,
// does not check if pending requests for a sub exist
func (v *EthereumVRFCoordinatorV2_5) OwnerCancelSubscription(subID *big.Int) (*types.Transaction, error) {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := v.coordinator.OwnerCancelSubscription(
		opts,
		subID,
	)
	if err != nil {
		return nil, err
	}
	return tx, v.client.ProcessTransaction(tx)
}

// CancelSubscription cancels subscription by Sub owner,
// return funds to specified address,
// checks if pending requests for a sub exist
func (v *EthereumVRFCoordinatorV2_5) CancelSubscription(subID *big.Int, to common.Address) (*types.Transaction, error) {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := v.coordinator.CancelSubscription(
		opts,
		subID,
		to,
	)
	if err != nil {
		return nil, err
	}
	return tx, v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2_5) Withdraw(recipient common.Address) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.Withdraw(
		opts,
		recipient,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2_5) WithdrawNative(recipient common.Address) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.WithdrawNative(
		opts,
		recipient,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2_5) SetConfig(
	minimumRequestConfirmations uint16,
	maxGasLimit uint32,
	stalenessSeconds uint32,
	gasAfterPaymentCalculation uint32,
	fallbackWeiPerUnitLink *big.Int,
	fulfillmentFlatFeeNativePPM uint32,
	fulfillmentFlatFeeLinkDiscountPPM uint32,
	nativePremiumPercentage uint8,
	linkPremiumPercentage uint8) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.SetConfig(
		opts,
		minimumRequestConfirmations,
		maxGasLimit,
		stalenessSeconds,
		gasAfterPaymentCalculation,
		fallbackWeiPerUnitLink,
		fulfillmentFlatFeeNativePPM,
		fulfillmentFlatFeeLinkDiscountPPM,
		nativePremiumPercentage,
		linkPremiumPercentage,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2_5) SetLINKAndLINKNativeFeed(linkAddress string, linkNativeFeedAddress string) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.SetLINKAndLINKNativeFeed(
		opts,
		common.HexToAddress(linkAddress),
		common.HexToAddress(linkNativeFeedAddress),
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2_5) RegisterProvingKey(
	publicProvingKey [2]*big.Int,
	gasLaneMaxGas uint64,
) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.RegisterProvingKey(opts, publicProvingKey, gasLaneMaxGas)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2_5) CreateSubscription() (*types.Transaction, error) {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := v.coordinator.CreateSubscription(opts)
	if err != nil {
		return nil, err
	}
	return tx, v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2_5) Migrate(subId *big.Int, coordinatorAddress string) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.Migrate(opts, subId, common.HexToAddress(coordinatorAddress))
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2_5) RegisterMigratableCoordinator(migratableCoordinatorAddress string) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.RegisterMigratableCoordinator(opts, common.HexToAddress(migratableCoordinatorAddress))
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2_5) AddConsumer(subId *big.Int, consumerAddress string) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.AddConsumer(
		opts,
		subId,
		common.HexToAddress(consumerAddress),
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2_5) FundSubscriptionWithNative(subId *big.Int, nativeTokenAmount *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	opts.Value = nativeTokenAmount
	tx, err := v.coordinator.FundSubscriptionWithNative(
		opts,
		subId,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2_5) FindSubscriptionID(subID *big.Int) (*big.Int, error) {
	owner := v.client.GetDefaultWallet().Address()
	subscriptionIterator, err := v.coordinator.FilterSubscriptionCreated(
		nil,
		[]*big.Int{subID},
	)
	if err != nil {
		return nil, err
	}

	if !subscriptionIterator.Next() {
		return nil, fmt.Errorf("expected at least 1 subID for the given owner %s", owner)
	}

	return subscriptionIterator.Event.SubId, nil
}

func (v *EthereumVRFCoordinatorV2_5) WaitForSubscriptionCreatedEvent(timeout time.Duration) (*vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCreated, error) {
	eventsChannel := make(chan *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCreated)
	subscription, err := v.coordinator.WatchSubscriptionCreated(nil, eventsChannel, nil)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for SubscriptionCreated event")
		case sub := <-eventsChannel:
			return sub, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2_5) WaitForSubscriptionCanceledEvent(subID *big.Int, timeout time.Duration) (*vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled, error) {
	eventsChannel := make(chan *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled)
	subscription, err := v.coordinator.WatchSubscriptionCanceled(nil, eventsChannel, []*big.Int{subID})
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for SubscriptionCanceled event")
		case sub := <-eventsChannel:
			return sub, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2_5) WaitForRandomWordsFulfilledEvent(filter RandomWordsFulfilledEventFilter) (*CoordinatorRandomWordsFulfilled, error) {
	randomWordsFulfilledEventsChannel := make(chan *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled)
	subscription, err := v.coordinator.WatchRandomWordsFulfilled(nil, randomWordsFulfilledEventsChannel, filter.RequestIds, filter.SubIDs)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(filter.Timeout):
			return nil, fmt.Errorf("timeout waiting for RandomWordsFulfilled event")
		case randomWordsFulfilledEvent := <-randomWordsFulfilledEventsChannel:
			return &CoordinatorRandomWordsFulfilled{
				RequestId:     randomWordsFulfilledEvent.RequestId,
				OutputSeed:    randomWordsFulfilledEvent.OutputSeed,
				SubId:         randomWordsFulfilledEvent.SubId.String(),
				Payment:       randomWordsFulfilledEvent.Payment,
				NativePayment: randomWordsFulfilledEvent.NativePayment,
				Success:       randomWordsFulfilledEvent.Success,
				OnlyPremium:   randomWordsFulfilledEvent.OnlyPremium,
				Raw:           randomWordsFulfilledEvent.Raw,
			}, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2_5) WaitForConfigSetEvent(timeout time.Duration) (*CoordinatorConfigSet, error) {
	eventsChannel := make(chan *vrf_coordinator_v2_5.VRFCoordinatorV25ConfigSet)
	subscription, err := v.coordinator.WatchConfigSet(nil, eventsChannel)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()
	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for ConfigSet event")
		case event := <-eventsChannel:
			return &CoordinatorConfigSet{
				MinimumRequestConfirmations:       event.MinimumRequestConfirmations,
				MaxGasLimit:                       event.MaxGasLimit,
				StalenessSeconds:                  event.StalenessSeconds,
				GasAfterPaymentCalculation:        event.GasAfterPaymentCalculation,
				FallbackWeiPerUnitLink:            event.FallbackWeiPerUnitLink,
				FulfillmentFlatFeeNativePPM:       event.FulfillmentFlatFeeNativePPM,
				FulfillmentFlatFeeLinkDiscountPPM: event.FulfillmentFlatFeeLinkDiscountPPM,
				NativePremiumPercentage:           event.NativePremiumPercentage,
				LinkPremiumPercentage:             event.LinkPremiumPercentage,
			}, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2_5) WaitForMigrationCompletedEvent(timeout time.Duration) (*vrf_coordinator_v2_5.VRFCoordinatorV25MigrationCompleted, error) {
	eventsChannel := make(chan *vrf_coordinator_v2_5.VRFCoordinatorV25MigrationCompleted)
	subscription, err := v.coordinator.WatchMigrationCompleted(nil, eventsChannel)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for MigrationCompleted event")
		case migrationCompletedEvent := <-eventsChannel:
			return migrationCompletedEvent, nil
		}
	}
}

func (v *EthereumVRFv2PlusLoadTestConsumer) Address() string {
	return v.address.Hex()
}
func (v *EthereumVRFv2PlusLoadTestConsumer) RequestRandomness(
	coordinator Coordinator,
	keyHash [32]byte, subID *big.Int,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	nativePayment bool,
	numWords uint32,
	requestCount uint16,
) (*CoordinatorRandomWordsRequested, error) {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := v.consumer.RequestRandomWords(opts, subID, requestConfirmations, keyHash, callbackGasLimit, nativePayment, numWords, requestCount)
	if err != nil {
		return nil, err
	}
	randomWordsRequestedEvent, err := RetrieveRequestRandomnessLogs(coordinator, v.client, tx)
	if err != nil {
		return nil, err
	}
	return randomWordsRequestedEvent, nil
}

func (v *EthereumVRFv2PlusLoadTestConsumer) ResetMetrics() error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.Reset(opts)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFv2PlusLoadTestConsumer) GetCoordinator(ctx context.Context) (common.Address, error) {
	return v.consumer.SVrfCoordinator(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
}
func (v *EthereumVRFv2PlusLoadTestConsumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_v2plus_load_test_with_metrics.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}, requestID)
}

func (v *EthereumVRFv2PlusLoadTestConsumer) GetLastRequestId(ctx context.Context) (*big.Int, error) {
	return v.consumer.SLastRequestId(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
}

func (v *EthereumVRFv2PlusLoadTestConsumer) GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error) {
	requestCount, err := v.consumer.SRequestCount(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	fulfilmentCount, err := v.consumer.SResponseCount(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})

	if err != nil {
		return nil, err
	}
	averageFulfillmentInMillions, err := v.consumer.SAverageResponseTimeInBlocksMillions(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	slowestFulfillment, err := v.consumer.SSlowestResponseTimeInBlocks(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})

	if err != nil {
		return nil, err
	}
	fastestFulfillment, err := v.consumer.SFastestResponseTimeInBlocks(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	averageResponseTimeInSeconds, err := v.consumer.SAverageResponseTimeInSecondsMillions(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	slowestResponseTimeInSeconds, err := v.consumer.SSlowestResponseTimeInSeconds(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	fastestResponseTimeInSeconds, err := v.consumer.SFastestResponseTimeInSeconds(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	var responseTimesInBlocks []uint32
	for {
		currentResponseTimesInBlocks, err := v.consumer.GetRequestBlockTimes(&bind.CallOpts{
			From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
			Context: ctx,
		}, big.NewInt(int64(len(responseTimesInBlocks))), big.NewInt(1000))
		if err != nil {
			return nil, err
		}
		if len(currentResponseTimesInBlocks) == 0 {
			break
		}
		responseTimesInBlocks = append(responseTimesInBlocks, currentResponseTimesInBlocks...)
	}
	var p90FulfillmentBlockTime, p95FulfillmentBlockTime float64
	if len(responseTimesInBlocks) == 0 {
		p90FulfillmentBlockTime = 0
		p95FulfillmentBlockTime = 0
	} else {
		responseTimesInBlocksFloat64 := make([]float64, len(responseTimesInBlocks))
		for i, value := range responseTimesInBlocks {
			responseTimesInBlocksFloat64[i] = float64(value)
		}
		p90FulfillmentBlockTime, err = stats.Percentile(responseTimesInBlocksFloat64, 90)
		if err != nil {
			return nil, err
		}
		p95FulfillmentBlockTime, err = stats.Percentile(responseTimesInBlocksFloat64, 95)
		if err != nil {
			return nil, err
		}
	}
	return &VRFLoadTestMetrics{
		RequestCount:                         requestCount,
		FulfilmentCount:                      fulfilmentCount,
		AverageFulfillmentInMillions:         averageFulfillmentInMillions,
		SlowestFulfillment:                   slowestFulfillment,
		FastestFulfillment:                   fastestFulfillment,
		P90FulfillmentBlockTime:              p90FulfillmentBlockTime,
		P95FulfillmentBlockTime:              p95FulfillmentBlockTime,
		AverageResponseTimeInSecondsMillions: averageResponseTimeInSeconds,
		SlowestResponseTimeInSeconds:         slowestResponseTimeInSeconds,
		FastestResponseTimeInSeconds:         fastestResponseTimeInSeconds,
	}, nil
}

func (e *EthereumContractDeployer) DeployVRFCoordinatorV2PlusUpgradedVersion(bhsAddr string) (VRFCoordinatorV2PlusUpgradedVersion, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinatorV2PlusUpgradedVersion", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_v2plus_upgraded_version.DeployVRFCoordinatorV2PlusUpgradedVersion(auth, backend, common.HexToAddress(bhsAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinatorV2PlusUpgradedVersion{
		client:      e.client,
		coordinator: instance.(*vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersion),
		address:     address,
	}, err
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	hash, err := v.coordinator.HashOfKey(opts, pubKey)
	if err != nil {
		return [32]byte{}, err
	}
	return hash, nil
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) GetActiveSubscriptionIds(ctx context.Context, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	activeSubscriptionIds, err := v.coordinator.GetActiveSubscriptionIds(opts, startIndex, maxCount)
	if err != nil {
		return nil, err
	}
	return activeSubscriptionIds, nil
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) GetSubscription(ctx context.Context, subID *big.Int) (vrf_v2plus_upgraded_version.GetSubscription, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	subscription, err := v.coordinator.GetSubscription(opts, subID)
	if err != nil {
		return vrf_v2plus_upgraded_version.GetSubscription{}, err
	}
	return subscription, nil
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) SetConfig(
	minimumRequestConfirmations uint16,
	maxGasLimit uint32,
	stalenessSeconds uint32,
	gasAfterPaymentCalculation uint32,
	fallbackWeiPerUnitLink *big.Int,
	fulfillmentFlatFeeNativePPM uint32,
	fulfillmentFlatFeeLinkDiscountPPM uint32,
	nativePremiumPercentage uint8,
	linkPremiumPercentage uint8,
) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.SetConfig(
		opts,
		minimumRequestConfirmations,
		maxGasLimit,
		stalenessSeconds,
		gasAfterPaymentCalculation,
		fallbackWeiPerUnitLink,
		fulfillmentFlatFeeNativePPM,
		fulfillmentFlatFeeLinkDiscountPPM,
		nativePremiumPercentage,
		linkPremiumPercentage,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) SetLINKAndLINKNativeFeed(linkAddress string, linkNativeFeedAddress string) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.SetLINKAndLINKNativeFeed(
		opts,
		common.HexToAddress(linkAddress),
		common.HexToAddress(linkNativeFeedAddress),
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) RegisterProvingKey(
	publicProvingKey [2]*big.Int,
	gasLaneMaxGas uint64,
) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.RegisterProvingKey(opts, publicProvingKey, gasLaneMaxGas)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) CreateSubscription() error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.CreateSubscription(opts)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) GetLinkTotalBalance(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	totalBalance, err := v.coordinator.STotalBalance(opts)
	if err != nil {
		return nil, err
	}
	return totalBalance, nil
}
func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) GetNativeTokenTotalBalance(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	totalBalance, err := v.coordinator.STotalNativeBalance(opts)
	if err != nil {
		return nil, err
	}
	return totalBalance, nil
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) Migrate(subId *big.Int, coordinatorAddress string) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.Migrate(opts, subId, common.HexToAddress(coordinatorAddress))
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) RegisterMigratableCoordinator(migratableCoordinatorAddress string) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.RegisterMigratableCoordinator(opts, common.HexToAddress(migratableCoordinatorAddress))
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) AddConsumer(subId *big.Int, consumerAddress string) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.AddConsumer(
		opts,
		subId,
		common.HexToAddress(consumerAddress),
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) FundSubscriptionWithNative(subId *big.Int, nativeTokenAmount *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	opts.Value = nativeTokenAmount
	tx, err := v.coordinator.FundSubscriptionWithNative(
		opts,
		subId,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) FindSubscriptionID() (*big.Int, error) {
	owner := v.client.GetDefaultWallet().Address()
	subscriptionIterator, err := v.coordinator.FilterSubscriptionCreated(
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	if !subscriptionIterator.Next() {
		return nil, fmt.Errorf("expected at least 1 subID for the given owner %s", owner)
	}

	return subscriptionIterator.Event.SubId, nil
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) WaitForRandomWordsFulfilledEvent(filter RandomWordsFulfilledEventFilter) (*CoordinatorRandomWordsFulfilled, error) {
	randomWordsFulfilledEventsChannel := make(chan *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled)
	subscription, err := v.coordinator.WatchRandomWordsFulfilled(nil, randomWordsFulfilledEventsChannel, filter.RequestIds, filter.SubIDs)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(filter.Timeout):
			return nil, fmt.Errorf("timeout waiting for RandomWordsFulfilled event")
		case randomWordsFulfilledEvent := <-randomWordsFulfilledEventsChannel:
			return &CoordinatorRandomWordsFulfilled{
				RequestId:     randomWordsFulfilledEvent.RequestId,
				OutputSeed:    randomWordsFulfilledEvent.OutputSeed,
				SubId:         randomWordsFulfilledEvent.SubId.String(),
				Payment:       randomWordsFulfilledEvent.Payment,
				NativePayment: randomWordsFulfilledEvent.NativePayment,
				Success:       randomWordsFulfilledEvent.Success,
				OnlyPremium:   randomWordsFulfilledEvent.OnlyPremium,
				Raw:           randomWordsFulfilledEvent.Raw,
			}, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) WaitForMigrationCompletedEvent(timeout time.Duration) (*vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted, error) {
	eventsChannel := make(chan *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted)
	subscription, err := v.coordinator.WatchMigrationCompleted(nil, eventsChannel)
	if err != nil {
		return nil, fmt.Errorf("parse RandomWordsRequested log failed, err: %w", err)
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for MigrationCompleted event")
		case migrationCompletedEvent := <-eventsChannel:
			return migrationCompletedEvent, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) ParseRandomWordsRequested(log types.Log) (*CoordinatorRandomWordsRequested, error) {
	randomWordsRequested, err := v.coordinator.ParseRandomWordsRequested(log)
	if err != nil {
		return nil, err
	}
	coordinatorRandomWordsRequested := &CoordinatorRandomWordsRequested{
		KeyHash:                     randomWordsRequested.KeyHash,
		RequestId:                   randomWordsRequested.RequestId,
		PreSeed:                     randomWordsRequested.PreSeed,
		SubId:                       randomWordsRequested.SubId.String(),
		MinimumRequestConfirmations: randomWordsRequested.MinimumRequestConfirmations,
		CallbackGasLimit:            randomWordsRequested.CallbackGasLimit,
		NumWords:                    randomWordsRequested.NumWords,
		ExtraArgs:                   randomWordsRequested.ExtraArgs,
		Sender:                      randomWordsRequested.Sender,
		Raw:                         randomWordsRequested.Raw,
	}
	return coordinatorRandomWordsRequested, nil
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) WaitForConfigSetEvent(timeout time.Duration) (*CoordinatorConfigSet, error) {
	eventsChannel := make(chan *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionConfigSet)
	subscription, err := v.coordinator.WatchConfigSet(nil, eventsChannel)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()
	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for ConfigSet event")
		case event := <-eventsChannel:
			return &CoordinatorConfigSet{
				MinimumRequestConfirmations:       event.MinimumRequestConfirmations,
				MaxGasLimit:                       event.MaxGasLimit,
				StalenessSeconds:                  event.StalenessSeconds,
				GasAfterPaymentCalculation:        event.GasAfterPaymentCalculation,
				FallbackWeiPerUnitLink:            event.FallbackWeiPerUnitLink,
				FulfillmentFlatFeeNativePPM:       event.FulfillmentFlatFeeNativePPM,
				FulfillmentFlatFeeLinkDiscountPPM: event.FulfillmentFlatFeeLinkDiscountPPM,
				NativePremiumPercentage:           event.NativePremiumPercentage,
				LinkPremiumPercentage:             event.LinkPremiumPercentage,
			}, nil
		}
	}
}

func (e *EthereumContractDeployer) DeployVRFv2PlusLoadTestConsumer(coordinatorAddr string) (VRFv2PlusLoadTestConsumer, error) {
	address, _, instance, err := e.client.DeployContract("VRFV2PlusLoadTestWithMetrics", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_v2plus_load_test_with_metrics.DeployVRFV2PlusLoadTestWithMetrics(auth, wrappers.MustNewWrappedContractBackend(e.client, nil), common.HexToAddress(coordinatorAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFv2PlusLoadTestConsumer{
		client:   e.client,
		consumer: instance.(*vrf_v2plus_load_test_with_metrics.VRFV2PlusLoadTestWithMetrics),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployVRFV2PlusWrapper(linkAddr string, linkEthFeedAddr string, coordinatorAddr string, subId *big.Int) (VRFV2PlusWrapper, error) {
	address, _, instance, err := e.client.DeployContract("VRFV2PlusWrapper", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrfv2plus_wrapper.DeployVRFV2PlusWrapper(auth, wrappers.MustNewWrappedContractBackend(e.client, nil), common.HexToAddress(linkAddr), common.HexToAddress(linkEthFeedAddr), common.HexToAddress(coordinatorAddr), subId)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFV2PlusWrapper{
		client:  e.client,
		wrapper: instance.(*vrfv2plus_wrapper.VRFV2PlusWrapper),
		address: address,
	}, err
}

func (e *EthereumContractDeployer) DeployVRFV2PlusWrapperLoadTestConsumer(vrfV2PlusWrapperAddr string) (VRFv2PlusWrapperLoadTestConsumer, error) {
	address, _, instance, err := e.client.DeployContract("VRFV2PlusWrapperLoadTestConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrfv2plus_wrapper_load_test_consumer.DeployVRFV2PlusWrapperLoadTestConsumer(auth, wrappers.MustNewWrappedContractBackend(e.client, nil), common.HexToAddress(vrfV2PlusWrapperAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFV2PlusWrapperLoadTestConsumer{
		client:   e.client,
		consumer: instance.(*vrfv2plus_wrapper_load_test_consumer.VRFV2PlusWrapperLoadTestConsumer),
		address:  address,
	}, err
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) Fund(ethAmount *big.Float) error {
	gasEstimates, err := v.client.EstimateGas(ethereum.CallMsg{
		To: v.address,
	})
	if err != nil {
		return err
	}
	return v.client.Fund(v.address.Hex(), ethAmount, gasEstimates)
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) RequestRandomness(
	coordinator Coordinator,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	numWords uint32,
	requestCount uint16,
) (*CoordinatorRandomWordsRequested, error) {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := v.consumer.MakeRequests(opts, callbackGasLimit, requestConfirmations, numWords, requestCount)
	if err != nil {
		return nil, err
	}
	randomWordsRequestedEvent, err := RetrieveRequestRandomnessLogs(coordinator, v.client, tx)
	if err != nil {
		return nil, err
	}
	return randomWordsRequestedEvent, nil
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) RequestRandomnessNative(
	coordinator Coordinator,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	numWords uint32,
	requestCount uint16,
) (*CoordinatorRandomWordsRequested, error) {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := v.consumer.MakeRequestsNative(opts, callbackGasLimit, requestConfirmations, numWords, requestCount)
	if err != nil {
		return nil, err
	}
	randomWordsRequestedEvent, err := RetrieveRequestRandomnessLogs(coordinator, v.client, tx)
	if err != nil {
		return nil, err
	}
	return randomWordsRequestedEvent, nil
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrfv2plus_wrapper_load_test_consumer.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}, requestID)
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) GetLastRequestId(ctx context.Context) (*big.Int, error) {
	return v.consumer.SLastRequestId(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) GetWrapper(ctx context.Context) (common.Address, error) {
	return v.consumer.IVrfV2PlusWrapper(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error) {
	requestCount, err := v.consumer.SRequestCount(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	fulfilmentCount, err := v.consumer.SResponseCount(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})

	if err != nil {
		return nil, err
	}
	averageFulfillmentInMillions, err := v.consumer.SAverageFulfillmentInMillions(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	slowestFulfillment, err := v.consumer.SSlowestFulfillment(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})

	if err != nil {
		return nil, err
	}
	fastestFulfillment, err := v.consumer.SFastestFulfillment(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}

	return &VRFLoadTestMetrics{
		RequestCount:                         requestCount,
		FulfilmentCount:                      fulfilmentCount,
		AverageFulfillmentInMillions:         averageFulfillmentInMillions,
		SlowestFulfillment:                   slowestFulfillment,
		FastestFulfillment:                   fastestFulfillment,
		AverageResponseTimeInSecondsMillions: nil,
		SlowestResponseTimeInSeconds:         nil,
		FastestResponseTimeInSeconds:         nil,
	}, nil
}
