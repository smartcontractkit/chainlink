package contracts

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_load_test_with_metrics"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_upgraded_version"
	"math/big"
	"time"
)

type EthereumVRFCoordinatorV2Plus struct {
	address     *common.Address
	client      blockchain.EVMClient
	coordinator *vrf_coordinator_v2plus.VRFCoordinatorV2Plus
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

// DeployVRFCoordinatorV2Plus deploys VRFV2Plus coordinator contract
func (e *EthereumContractDeployer) DeployVRFCoordinatorV2Plus(bhsAddr string) (VRFCoordinatorV2Plus, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinatorV2Plus", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_coordinator_v2plus.DeployVRFCoordinatorV2Plus(auth, backend, common.HexToAddress(bhsAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinatorV2Plus{
		client:      e.client,
		coordinator: instance.(*vrf_coordinator_v2plus.VRFCoordinatorV2Plus),
		address:     address,
	}, err
}

func (v *EthereumVRFCoordinatorV2Plus) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFCoordinatorV2Plus) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
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

func (v *EthereumVRFCoordinatorV2Plus) GetSubscription(ctx context.Context, subID *big.Int) (vrf_coordinator_v2plus.GetSubscription, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	subscription, err := v.coordinator.GetSubscription(opts, subID)
	if err != nil {
		return vrf_coordinator_v2plus.GetSubscription{}, err
	}
	return subscription, nil
}

func (v *EthereumVRFCoordinatorV2Plus) GetLinkTotalBalance(ctx context.Context) (*big.Int, error) {
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
func (v *EthereumVRFCoordinatorV2Plus) GetNativeTokenTotalBalance(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	totalBalance, err := v.coordinator.STotalEthBalance(opts)
	if err != nil {
		return nil, err
	}
	return totalBalance, nil
}

func (v *EthereumVRFCoordinatorV2Plus) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig vrf_coordinator_v2plus.VRFCoordinatorV2PlusFeeConfig) error {
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
		feeConfig,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2Plus) SetLINKAndLINKETHFeed(linkAddress string, linkEthFeedAddress string) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.SetLINKAndLINKETHFeed(
		opts,
		common.HexToAddress(linkAddress),
		common.HexToAddress(linkEthFeedAddress),
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2Plus) RegisterProvingKey(
	oracleAddr string,
	publicProvingKey [2]*big.Int,
) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.RegisterProvingKey(opts, common.HexToAddress(oracleAddr), publicProvingKey)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2Plus) CreateSubscription() error {
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

func (v *EthereumVRFCoordinatorV2Plus) Migrate(subId *big.Int, coordinatorAddress string) error {
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

func (v *EthereumVRFCoordinatorV2Plus) RegisterMigratableCoordinator(migratableCoordinatorAddress string) error {
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

func (v *EthereumVRFCoordinatorV2Plus) AddConsumer(subId *big.Int, consumerAddress string) error {
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

func (v *EthereumVRFCoordinatorV2Plus) FundSubscriptionWithEth(subId *big.Int, nativeTokenAmount *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	opts.Value = nativeTokenAmount
	tx, err := v.coordinator.FundSubscriptionWithEth(
		opts,
		subId,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2Plus) FindSubscriptionID() (*big.Int, error) {
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

func (v *EthereumVRFCoordinatorV2Plus) WaitForRandomWordsFulfilledEvent(subID []*big.Int, requestID []*big.Int, timeout time.Duration) (*vrf_coordinator_v2plus.VRFCoordinatorV2PlusRandomWordsFulfilled, error) {
	randomWordsFulfilledEventsChannel := make(chan *vrf_coordinator_v2plus.VRFCoordinatorV2PlusRandomWordsFulfilled)
	subscription, err := v.coordinator.WatchRandomWordsFulfilled(nil, randomWordsFulfilledEventsChannel, requestID, subID)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for RandomWordsFulfilled event")
		case randomWordsFulfilledEvent := <-randomWordsFulfilledEventsChannel:
			return randomWordsFulfilledEvent, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2Plus) WaitForRandomWordsRequestedEvent(keyHash [][32]byte, subID []*big.Int, sender []common.Address, timeout time.Duration) (*vrf_coordinator_v2plus.VRFCoordinatorV2PlusRandomWordsRequested, error) {
	randomWordsFulfilledEventsChannel := make(chan *vrf_coordinator_v2plus.VRFCoordinatorV2PlusRandomWordsRequested)
	subscription, err := v.coordinator.WatchRandomWordsRequested(nil, randomWordsFulfilledEventsChannel, keyHash, subID, sender)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for RandomWordsFulfilled event")
		case randomWordsFulfilledEvent := <-randomWordsFulfilledEventsChannel:
			return randomWordsFulfilledEvent, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2Plus) WaitForMigrationCompletedEvent(timeout time.Duration) (*vrf_coordinator_v2plus.VRFCoordinatorV2PlusMigrationCompleted, error) {
	eventsChannel := make(chan *vrf_coordinator_v2plus.VRFCoordinatorV2PlusMigrationCompleted)
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
func (v *EthereumVRFv2PlusLoadTestConsumer) RequestRandomness(keyHash [32]byte, subID *big.Int, requestConfirmations uint16, callbackGasLimit uint32, nativePayment bool, numWords uint32, requestCount uint16) (*types.Transaction, error) {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := v.consumer.RequestRandomWords(opts, subID, requestConfirmations, keyHash, callbackGasLimit, nativePayment, numWords, requestCount)
	if err != nil {
		return nil, err
	}

	return tx, v.client.ProcessTransaction(tx)
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
		requestCount,
		fulfilmentCount,
		averageFulfillmentInMillions,
		slowestFulfillment,
		fastestFulfillment,
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

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionFeeConfig) error {
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
		feeConfig,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) SetLINKAndLINKETHFeed(linkAddress string, linkEthFeedAddress string) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.SetLINKAndLINKETHFeed(
		opts,
		common.HexToAddress(linkAddress),
		common.HexToAddress(linkEthFeedAddress),
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) RegisterProvingKey(
	oracleAddr string,
	publicProvingKey [2]*big.Int,
) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.RegisterProvingKey(opts, common.HexToAddress(oracleAddr), publicProvingKey)
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
	totalBalance, err := v.coordinator.STotalEthBalance(opts)
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

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) FundSubscriptionWithEth(subId *big.Int, nativeTokenAmount *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	opts.Value = nativeTokenAmount
	tx, err := v.coordinator.FundSubscriptionWithEth(
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

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) WaitForRandomWordsFulfilledEvent(subID []*big.Int, requestID []*big.Int, timeout time.Duration) (*vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled, error) {
	randomWordsFulfilledEventsChannel := make(chan *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled)
	subscription, err := v.coordinator.WatchRandomWordsFulfilled(nil, randomWordsFulfilledEventsChannel, requestID, subID)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for RandomWordsFulfilled event")
		case randomWordsFulfilledEvent := <-randomWordsFulfilledEventsChannel:
			return randomWordsFulfilledEvent, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) WaitForMigrationCompletedEvent(timeout time.Duration) (*vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted, error) {
	eventsChannel := make(chan *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted)
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

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) WaitForRandomWordsRequestedEvent(keyHash [][32]byte, subID []*big.Int, sender []common.Address, timeout time.Duration) (*vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested, error) {
	eventsChannel := make(chan *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested)
	subscription, err := v.coordinator.WatchRandomWordsRequested(nil, eventsChannel, keyHash, subID, sender)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for RandomWordsRequested event")
		case randomWordsRequestedEvent := <-eventsChannel:
			return randomWordsRequestedEvent, nil
		}
	}
}
