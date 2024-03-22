package contracts

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_test_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_mock_ethlink_aggregator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_owner"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_load_test_with_metrics"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_wrapper_load_test_consumer"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2_consumer_wrapper"
)

// EthereumVRFCoordinatorV2 represents VRFV2 coordinator contract
type EthereumVRFCoordinatorV2 struct {
	address     *common.Address
	client      blockchain.EVMClient
	coordinator *vrf_coordinator_v2.VRFCoordinatorV2
}

type EthereumVRFOwner struct {
	address  *common.Address
	client   blockchain.EVMClient
	vrfOwner *vrf_owner.VRFOwner
}

type EthereumVRFCoordinatorTestV2 struct {
	address     *common.Address
	client      blockchain.EVMClient
	coordinator *vrf_coordinator_test_v2.VRFCoordinatorTestV2
}

// EthereumVRFConsumerV2 represents VRFv2 consumer contract
type EthereumVRFConsumerV2 struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *vrf_consumer_v2.VRFConsumerV2
}

// EthereumVRFv2Consumer represents VRFv2 consumer contract
type EthereumVRFv2Consumer struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *vrf_v2_consumer_wrapper.VRFv2Consumer
}

// EthereumVRFv2LoadTestConsumer represents VRFv2 consumer contract for performing Load Tests
type EthereumVRFv2LoadTestConsumer struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *vrf_load_test_with_metrics.VRFV2LoadTestWithMetrics
}

type EthereumVRFV2Wrapper struct {
	address *common.Address
	client  blockchain.EVMClient
	wrapper *vrfv2_wrapper.VRFV2Wrapper
}

type EthereumVRFV2WrapperLoadTestConsumer struct {
	address  *common.Address
	client   blockchain.EVMClient
	consumer *vrfv2_wrapper_load_test_consumer.VRFV2WrapperLoadTestConsumer
}

type GetRequestConfig struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	ProvingKeyHashes            [32]byte
}

type EthereumVRFMockETHLINKFeed struct {
	client  blockchain.EVMClient
	feed    *vrf_mock_ethlink_aggregator.VRFMockETHLINKAggregator
	address *common.Address
}

// DeployVRFCoordinatorV2 deploys VRFV2 coordinator contract
func (e *EthereumContractDeployer) DeployVRFCoordinatorV2(linkAddr string, bhsAddr string, linkEthFeedAddr string) (VRFCoordinatorV2, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinatorV2", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_coordinator_v2.DeployVRFCoordinatorV2(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(bhsAddr), common.HexToAddress(linkEthFeedAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinatorV2{
		client:      e.client,
		coordinator: instance.(*vrf_coordinator_v2.VRFCoordinatorV2),
		address:     address,
	}, err
}

func (e *EthereumContractDeployer) DeployVRFOwner(coordinatorAddr string) (VRFOwner, error) {
	address, _, instance, err := e.client.DeployContract("VRFOwner", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_owner.DeployVRFOwner(auth, backend, common.HexToAddress(coordinatorAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFOwner{
		client:   e.client,
		vrfOwner: instance.(*vrf_owner.VRFOwner),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployVRFCoordinatorTestV2(linkAddr string, bhsAddr string, linkEthFeedAddr string) (*EthereumVRFCoordinatorTestV2, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinatorTestV2", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_coordinator_test_v2.DeployVRFCoordinatorTestV2(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(bhsAddr), common.HexToAddress(linkEthFeedAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinatorTestV2{
		client:      e.client,
		coordinator: instance.(*vrf_coordinator_test_v2.VRFCoordinatorTestV2),
		address:     address,
	}, err
}

// DeployVRFConsumerV2 deploys VRFv@ consumer contract
func (e *EthereumContractDeployer) DeployVRFConsumerV2(linkAddr string, coordinatorAddr string) (VRFConsumerV2, error) {
	address, _, instance, err := e.client.DeployContract("VRFConsumerV2", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_consumer_v2.DeployVRFConsumerV2(auth, backend, common.HexToAddress(coordinatorAddr), common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFConsumerV2{
		client:   e.client,
		consumer: instance.(*vrf_consumer_v2.VRFConsumerV2),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployVRFv2Consumer(coordinatorAddr string) (VRFv2Consumer, error) {
	address, _, instance, err := e.client.DeployContract("VRFv2Consumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_v2_consumer_wrapper.DeployVRFv2Consumer(auth, backend, common.HexToAddress(coordinatorAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFv2Consumer{
		client:   e.client,
		consumer: instance.(*vrf_v2_consumer_wrapper.VRFv2Consumer),
		address:  address,
	}, err
}

// DeployVRFv2LoadTestConsumer(coordinatorAddr string) (VRFv2Consumer, error)
func (e *EthereumContractDeployer) DeployVRFv2LoadTestConsumer(coordinatorAddr string) (VRFv2LoadTestConsumer, error) {
	address, _, instance, err := e.client.DeployContract("VRFV2LoadTestWithMetrics", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrf_load_test_with_metrics.DeployVRFV2LoadTestWithMetrics(auth, backend, common.HexToAddress(coordinatorAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFv2LoadTestConsumer{
		client:   e.client,
		consumer: instance.(*vrf_load_test_with_metrics.VRFV2LoadTestWithMetrics),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployVRFV2Wrapper(linkAddr string, linkEthFeedAddr string, coordinatorAddr string) (VRFV2Wrapper, error) {
	address, _, instance, err := e.client.DeployContract("VRFV2Wrapper", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrfv2_wrapper.DeployVRFV2Wrapper(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(linkEthFeedAddr), common.HexToAddress(coordinatorAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFV2Wrapper{
		address: address,
		client:  e.client,
		wrapper: instance.(*vrfv2_wrapper.VRFV2Wrapper),
	}, err
}

func (e *EthereumContractDeployer) DeployVRFV2WrapperLoadTestConsumer(linkAddr string, vrfV2WrapperAddr string) (VRFv2WrapperLoadTestConsumer, error) {
	address, _, instance, err := e.client.DeployContract("VRFV2WrapperLoadTestConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return vrfv2_wrapper_load_test_consumer.DeployVRFV2WrapperLoadTestConsumer(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(vrfV2WrapperAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFV2WrapperLoadTestConsumer{
		address:  address,
		client:   e.client,
		consumer: instance.(*vrfv2_wrapper_load_test_consumer.VRFV2WrapperLoadTestConsumer),
	}, err
}

func (v *EthereumVRFCoordinatorV2) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFCoordinatorV2) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
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

func (v *EthereumVRFCoordinatorV2) GetSubscription(ctx context.Context, subID uint64) (vrf_coordinator_v2.GetSubscription, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	subscription, err := v.coordinator.GetSubscription(opts, subID)
	if err != nil {
		return vrf_coordinator_v2.GetSubscription{}, err
	}
	return subscription, nil
}

func (v *EthereumVRFCoordinatorV2) GetOwner(ctx context.Context) (common.Address, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	coordinatorOwnerAddress, err := v.coordinator.Owner(opts)
	if err != nil {
		return common.Address{}, err
	}
	return coordinatorOwnerAddress, nil
}

func (v *EthereumVRFCoordinatorV2) GetRequestConfig(ctx context.Context) (GetRequestConfig, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	minConfirmations, maxGas, keyHashes, err := v.coordinator.GetRequestConfig(opts)
	if err != nil {
		return GetRequestConfig{}, err
	}
	requestConfig := GetRequestConfig{
		MinimumRequestConfirmations: minConfirmations,
		MaxGasLimit:                 maxGas,
		ProvingKeyHashes:            keyHashes[0],
	}

	return requestConfig, nil
}

func (v *EthereumVRFCoordinatorV2) GetConfig(ctx context.Context) (vrf_coordinator_v2.GetConfig, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	config, err := v.coordinator.GetConfig(opts)
	if err != nil {
		return vrf_coordinator_v2.GetConfig{}, err
	}
	return config, nil
}

func (v *EthereumVRFCoordinatorV2) GetFallbackWeiPerUnitLink(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	fallbackWeiPerUnitLink, err := v.coordinator.GetFallbackWeiPerUnitLink(opts)
	if err != nil {
		return nil, err
	}
	return fallbackWeiPerUnitLink, nil
}

func (v *EthereumVRFCoordinatorV2) GetFeeConfig(ctx context.Context) (vrf_coordinator_v2.GetFeeConfig, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	config, err := v.coordinator.GetFeeConfig(opts)
	if err != nil {
		return vrf_coordinator_v2.GetFeeConfig{}, err
	}
	return config, nil
}

func (v *EthereumVRFCoordinatorV2) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig vrf_coordinator_v2.VRFCoordinatorV2FeeConfig) error {
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

func (v *EthereumVRFCoordinatorV2) RegisterProvingKey(
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

func (v *EthereumVRFCoordinatorV2) TransferOwnership(to common.Address) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.TransferOwnership(opts, to)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFCoordinatorV2) CreateSubscription() (*types.Transaction, error) {
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

func (v *EthereumVRFCoordinatorV2) AddConsumer(subId uint64, consumerAddress string) error {
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

func (v *EthereumVRFCoordinatorV2) PendingRequestsExist(ctx context.Context, subID uint64) (bool, error) {
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

func (v *EthereumVRFCoordinatorV2) OracleWithdraw(recipient common.Address, amount *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.coordinator.OracleWithdraw(opts, recipient, amount)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// OwnerCancelSubscription cancels subscription,
// return funds to the subscription owner,
// down not check if pending requests for a sub exist,
// outstanding requests may fail onchain
func (v *EthereumVRFCoordinatorV2) OwnerCancelSubscription(subID uint64) (*types.Transaction, error) {
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
func (v *EthereumVRFCoordinatorV2) CancelSubscription(subID uint64, to common.Address) (*types.Transaction, error) {
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

func (v *EthereumVRFCoordinatorV2) FindSubscriptionID(subID uint64) (uint64, error) {
	owner := v.client.GetDefaultWallet().Address()
	subscriptionIterator, err := v.coordinator.FilterSubscriptionCreated(
		nil,
		[]uint64{subID},
	)
	if err != nil {
		return 0, err
	}

	if !subscriptionIterator.Next() {
		return 0, fmt.Errorf("expected at least 1 subID for the given owner %s", owner)
	}

	return subscriptionIterator.Event.SubId, nil
}

func (v *EthereumVRFCoordinatorV2) WaitForRandomWordsFulfilledEvent(requestID []*big.Int, timeout time.Duration) (*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled, error) {
	randomWordsFulfilledEventsChannel := make(chan *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled)
	subscription, err := v.coordinator.WatchRandomWordsFulfilled(nil, randomWordsFulfilledEventsChannel, requestID)
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

func (v *EthereumVRFCoordinatorV2) WaitForRandomWordsRequestedEvent(keyHash [][32]byte, subID []uint64, sender []common.Address, timeout time.Duration) (*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested, error) {
	randomWordsFulfilledEventsChannel := make(chan *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested)
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
			return nil, fmt.Errorf("timeout waiting for RandomWordsRequested event")
		case randomWordsFulfilledEvent := <-randomWordsFulfilledEventsChannel:
			return randomWordsFulfilledEvent, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2) WaitForSubscriptionFunded(subID []uint64, timeout time.Duration) (*vrf_coordinator_v2.VRFCoordinatorV2SubscriptionFunded, error) {
	eventsChannel := make(chan *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionFunded)
	subscription, err := v.coordinator.WatchSubscriptionFunded(nil, eventsChannel, subID)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for SubscriptionFunded event")
		case event := <-eventsChannel:
			return event, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2) WaitForSubscriptionCanceledEvent(subID []uint64, timeout time.Duration) (*vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled, error) {
	eventsChannel := make(chan *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled)
	subscription, err := v.coordinator.WatchSubscriptionCanceled(nil, eventsChannel, subID)
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

func (v *EthereumVRFCoordinatorV2) WaitForSubscriptionCreatedEvent(subID []uint64, timeout time.Duration) (*vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCreated, error) {
	eventsChannel := make(chan *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCreated)
	subscription, err := v.coordinator.WatchSubscriptionCreated(nil, eventsChannel, subID)
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
		case event := <-eventsChannel:
			return event, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2) WaitForSubscriptionConsumerAdded(subID []uint64, timeout time.Duration) (*vrf_coordinator_v2.VRFCoordinatorV2SubscriptionConsumerAdded, error) {
	eventsChannel := make(chan *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionConsumerAdded)
	subscription, err := v.coordinator.WatchSubscriptionConsumerAdded(nil, eventsChannel, subID)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for SubscriptionConsumerAdded event")
		case event := <-eventsChannel:
			return event, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2) WaitForSubscriptionConsumerRemoved(subID []uint64, timeout time.Duration) (*vrf_coordinator_v2.VRFCoordinatorV2SubscriptionConsumerRemoved, error) {
	eventsChannel := make(chan *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionConsumerRemoved)
	subscription, err := v.coordinator.WatchSubscriptionConsumerRemoved(nil, eventsChannel, subID)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for SubscriptionConsumerRemoved event")
		case event := <-eventsChannel:
			return event, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2) WaitForConfigSetEvent(timeout time.Duration) (*vrf_coordinator_v2.VRFCoordinatorV2ConfigSet, error) {
	eventsChannel := make(chan *vrf_coordinator_v2.VRFCoordinatorV2ConfigSet)
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
			return event, nil
		}
	}
}

// GetAllRandomWords get all VRFv2 randomness output words
func (v *EthereumVRFConsumerV2) GetAllRandomWords(ctx context.Context, num int) ([]*big.Int, error) {
	words := make([]*big.Int, 0)
	for i := 0; i < num; i++ {
		word, err := v.consumer.SRandomWords(&bind.CallOpts{
			From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
			Context: ctx,
		}, big.NewInt(int64(i)))
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, nil
}

// LoadExistingConsumer loads an EthereumVRFConsumerV2 with a specified address
func (v *EthereumVRFConsumerV2) LoadExistingConsumer(address string, client blockchain.EVMClient) error {
	a := common.HexToAddress(address)
	consumer, err := vrf_consumer_v2.NewVRFConsumerV2(a, client.(*blockchain.EthereumClient).Client)
	if err != nil {
		return err
	}
	v.client = client
	v.consumer = consumer
	v.address = &a
	return nil
}

// CreateFundedSubscription create funded subscription for VRFv2 randomness
func (v *EthereumVRFConsumerV2) CreateFundedSubscription(funds *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.CreateSubscriptionAndFund(opts, funds)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// TopUpSubscriptionFunds add funds to a VRFv2 subscription
func (v *EthereumVRFConsumerV2) TopUpSubscriptionFunds(funds *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.TopUpSubscription(opts, funds)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFConsumerV2) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFv2Consumer) Address() string {
	return v.address.Hex()
}

// CurrentSubscription get current VRFv2 subscription
func (v *EthereumVRFConsumerV2) CurrentSubscription() (uint64, error) {
	return v.consumer.SSubId(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
}

// GasAvailable get available gas after randomness fulfilled
func (v *EthereumVRFConsumerV2) GasAvailable() (*big.Int, error) {
	return v.consumer.SGasAvailable(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
}

func (v *EthereumVRFConsumerV2) Fund(ethAmount *big.Float) error {
	gasEstimates, err := v.client.EstimateGas(ethereum.CallMsg{
		To: v.address,
	})
	if err != nil {
		return err
	}
	return v.client.Fund(v.address.Hex(), ethAmount, gasEstimates)
}

// RequestRandomness request VRFv2 random words
func (v *EthereumVRFConsumerV2) RequestRandomness(hash [32]byte, subID uint64, confs uint16, gasLimit uint32, numWords uint32) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.RequestRandomness(opts, hash, subID, confs, gasLimit, numWords)
	if err != nil {
		return err
	}
	log.Info().Interface("Sub ID", subID).
		Interface("Number of Words", numWords).
		Interface("Number of Confirmations", confs).
		Interface("Callback Gas Limit", gasLimit).
		Interface("KeyHash", hex.EncodeToString(hash[:])).
		Interface("Consumer Contract", v.address).
		Msg("RequestRandomness called")
	return v.client.ProcessTransaction(tx)
}

// RequestRandomness request VRFv2 random words
func (v *EthereumVRFv2Consumer) RequestRandomness(hash [32]byte, subID uint64, confs uint16, gasLimit uint32, numWords uint32) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.consumer.RequestRandomWords(opts, subID, gasLimit, confs, numWords, hash)
	if err != nil {
		return err
	}
	log.Info().Interface("Sub ID", subID).
		Interface("Number of Words", numWords).
		Interface("Number of Confirmations", confs).
		Interface("Callback Gas Limit", gasLimit).
		Interface("KeyHash", hex.EncodeToString(hash[:])).
		Interface("Consumer Contract", v.address).
		Msg("RequestRandomness called")
	return v.client.ProcessTransaction(tx)
}

// RandomnessOutput get VRFv2 randomness output (word)
func (v *EthereumVRFConsumerV2) RandomnessOutput(ctx context.Context, arg0 *big.Int) (*big.Int, error) {
	return v.consumer.SRandomWords(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}, arg0)
}

func (v *EthereumVRFv2LoadTestConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFv2LoadTestConsumer) RequestRandomness(
	keyHash [32]byte,
	subID uint64,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	numWords uint32,
	requestCount uint16,
) (*types.Transaction, error) {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}

	tx, err := v.consumer.RequestRandomWords(opts, subID, requestConfirmations, keyHash, callbackGasLimit, numWords, requestCount)
	if err != nil {
		return nil, err
	}
	return tx, v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFv2LoadTestConsumer) RequestRandomWordsWithForceFulfill(
	keyHash [32]byte,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	numWords uint32,
	requestCount uint16,
	subTopUpAmount *big.Int,
	linkAddress common.Address,
) (*types.Transaction, error) {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := v.consumer.RequestRandomWordsWithForceFulfill(
		opts,
		requestConfirmations,
		keyHash,
		callbackGasLimit,
		numWords,
		requestCount,
		subTopUpAmount,
		linkAddress,
	)
	if err != nil {
		return nil, err
	}
	return tx, v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFv2Consumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_v2_consumer_wrapper.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}, requestID)
}

func (v *EthereumVRFv2Consumer) GetLastRequestId(ctx context.Context) (*big.Int, error) {
	return v.consumer.LastRequestId(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
}

func (v *EthereumVRFv2LoadTestConsumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_load_test_with_metrics.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}, requestID)
}

func (v *EthereumVRFv2LoadTestConsumer) GetLastRequestId(ctx context.Context) (*big.Int, error) {
	return v.consumer.SLastRequestId(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
}

func (v *EthereumVRFv2LoadTestConsumer) ResetMetrics() error {
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

func (v *EthereumVRFv2LoadTestConsumer) GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error) {
	requestCount, err := v.consumer.SRequestCount(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}
	fulfilmentCount, err := v.consumer.SResponseCount(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})

	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}
	averageFulfillmentInMillions, err := v.consumer.SAverageFulfillmentInMillions(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}
	slowestFulfillment, err := v.consumer.SSlowestFulfillment(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})

	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}
	fastestFulfillment, err := v.consumer.SFastestFulfillment(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}

	return &VRFLoadTestMetrics{
		requestCount,
		fulfilmentCount,
		averageFulfillmentInMillions,
		slowestFulfillment,
		fastestFulfillment,
	}, nil
}

func (v *EthereumVRFV2Wrapper) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFV2Wrapper) SetConfig(wrapperGasOverhead uint32, coordinatorGasOverhead uint32, wrapperPremiumPercentage uint8, keyHash [32]byte, maxNumWords uint8) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.wrapper.SetConfig(
		opts,
		wrapperGasOverhead,
		coordinatorGasOverhead,
		wrapperPremiumPercentage,
		keyHash,
		maxNumWords,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFV2Wrapper) GetSubID(ctx context.Context) (uint64, error) {
	return v.wrapper.SUBSCRIPTIONID(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) Fund(ethAmount *big.Float) error {
	gasEstimates, err := v.client.EstimateGas(ethereum.CallMsg{})
	if err != nil {
		return err
	}
	return v.client.Fund(v.address.Hex(), ethAmount, gasEstimates)
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) RequestRandomness(requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, requestCount uint16) (*types.Transaction, error) {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := v.consumer.MakeRequests(opts, callbackGasLimit, requestConfirmations, numWords, requestCount)
	if err != nil {
		return nil, err
	}
	return tx, v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrfv2_wrapper_load_test_consumer.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	}, requestID)
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) GetLastRequestId(ctx context.Context) (*big.Int, error) {
	return v.consumer.SLastRequestId(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) GetWrapper(ctx context.Context) (common.Address, error) {
	return v.consumer.IVrfV2Wrapper(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: ctx,
	})
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error) {
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

func (v *EthereumVRFOwner) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFOwner) SetAuthorizedSenders(senders []common.Address) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.vrfOwner.SetAuthorizedSenders(
		opts,
		senders,
	)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFOwner) AcceptVRFOwnership() error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.vrfOwner.AcceptVRFOwnership(opts)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumVRFOwner) WaitForRandomWordsForcedEvent(requestIDs []*big.Int, subIds []uint64, senders []common.Address, timeout time.Duration) (*vrf_owner.VRFOwnerRandomWordsForced, error) {
	eventsChannel := make(chan *vrf_owner.VRFOwnerRandomWordsForced)
	subscription, err := v.vrfOwner.WatchRandomWordsForced(nil, eventsChannel, requestIDs, subIds, senders)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for RandomWordsForced event")
		case event := <-eventsChannel:
			return event, nil
		}
	}
}

func (v *EthereumVRFCoordinatorTestV2) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFMockETHLINKFeed) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFMockETHLINKFeed) LatestRoundData() (*big.Int, error) {
	data, err := v.feed.LatestRoundData(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.Ans, nil
}

func (v *EthereumVRFMockETHLINKFeed) LatestRoundDataUpdatedAt() (*big.Int, error) {
	data, err := v.feed.LatestRoundData(&bind.CallOpts{
		From:    common.HexToAddress(v.client.GetDefaultWallet().Address()),
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.UpdatedAt, nil
}

func (v *EthereumVRFMockETHLINKFeed) SetBlockTimestampDeduction(blockTimestampDeduction *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := v.feed.SetBlockTimestampDeduction(opts, blockTimestampDeduction)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}
