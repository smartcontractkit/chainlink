package contracts

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_vrf_coordinator_v2plus"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_v2plus_load_test_with_metrics"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_v2plus_upgraded_version"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2plus_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2plus_wrapper_load_test_consumer"
)

// VRFSubscription holds data returned by GetSubscription.
type VRFSubscription struct {
	Balance       *big.Int
	NativeBalance *big.Int
	Owner         common.Address
	Consumers     []common.Address
}

// EthereumBlockhashStore wraps the BHS contract.
type EthereumBlockhashStore struct {
	client  *seth.Client
	store   *blockhash_store.BlockhashStore
	address common.Address
}

func (v *EthereumBlockhashStore) Address() string { return v.address.Hex() }

// EthereumBatchBlockhashStore wraps the BatchBHS contract.
type EthereumBatchBlockhashStore struct {
	client  *seth.Client
	store   *batch_blockhash_store.BatchBlockhashStore
	address common.Address
}

func (v *EthereumBatchBlockhashStore) Address() string { return v.address.Hex() }

// EthereumVRFCoordinatorV2_5 wraps the VRF Coordinator V2.5 contract.
type EthereumVRFCoordinatorV2_5 struct {
	client      *seth.Client
	coordinator *vrf_coordinator_v2_5.VRFCoordinatorV25
	address     common.Address
}

func (v *EthereumVRFCoordinatorV2_5) Address() string { return v.address.Hex() }

// EthereumBatchVRFCoordinatorV2Plus wraps the Batch VRF Coordinator V2 Plus contract.
type EthereumBatchVRFCoordinatorV2Plus struct {
	client      *seth.Client
	coordinator *batch_vrf_coordinator_v2plus.BatchVRFCoordinatorV2Plus
	address     common.Address
}

func (v *EthereumBatchVRFCoordinatorV2Plus) Address() string { return v.address.Hex() }

// EthereumVRFv2PlusLoadTestConsumer wraps the VRFv2Plus load test consumer contract.
type EthereumVRFv2PlusLoadTestConsumer struct {
	client   *seth.Client
	consumer *vrf_v2plus_load_test_with_metrics.VRFV2PlusLoadTestWithMetrics
	address  common.Address
}

func (v *EthereumVRFv2PlusLoadTestConsumer) Address() string { return v.address.Hex() }

// EthereumVRFV2PlusWrapper wraps the VRFV2PlusWrapper contract.
type EthereumVRFV2PlusWrapper struct {
	client  *seth.Client
	wrapper *vrfv2plus_wrapper.VRFV2PlusWrapper
	address common.Address
}

func (v *EthereumVRFV2PlusWrapper) Address() string { return v.address.Hex() }

// EthereumVRFV2PlusWrapperLoadTestConsumer wraps the wrapper load test consumer.
type EthereumVRFV2PlusWrapperLoadTestConsumer struct {
	client   *seth.Client
	consumer *vrfv2plus_wrapper_load_test_consumer.VRFV2PlusWrapperLoadTestConsumer
	address  common.Address
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) Address() string { return v.address.Hex() }

// EthereumVRFCoordinatorV2PlusUpgradedVersion wraps the upgraded VRF Coordinator V2 Plus contract.
type EthereumVRFCoordinatorV2PlusUpgradedVersion struct {
	client      *seth.Client
	coordinator *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersion
	address     common.Address
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) Address() string { return v.address.Hex() }

// --- Coordinator methods ---

func (v *EthereumVRFCoordinatorV2_5) SetLINKAndLINKNativeFeed(linkAddress, linkNativeFeedAddress string) error {
	_, err := v.client.Decode(v.coordinator.SetLINKAndLINKNativeFeed(
		v.client.NewTXOpts(),
		common.HexToAddress(linkAddress),
		common.HexToAddress(linkNativeFeedAddress),
	))
	return err
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
	linkPremiumPercentage uint8,
) error {
	_, err := v.client.Decode(v.coordinator.SetConfig(
		v.client.NewTXOpts(),
		minimumRequestConfirmations,
		maxGasLimit,
		stalenessSeconds,
		gasAfterPaymentCalculation,
		fallbackWeiPerUnitLink,
		fulfillmentFlatFeeNativePPM,
		fulfillmentFlatFeeLinkDiscountPPM,
		nativePremiumPercentage,
		linkPremiumPercentage,
	))
	return err
}

func (v *EthereumVRFCoordinatorV2_5) RegisterProvingKey(publicProvingKey [2]*big.Int, gasLaneMaxGas uint64) error {
	_, err := v.client.Decode(v.coordinator.RegisterProvingKey(v.client.NewTXOpts(), publicProvingKey, gasLaneMaxGas))
	return err
}

func (v *EthereumVRFCoordinatorV2_5) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
	return v.coordinator.HashOfKey(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, pubKey)
}

func (v *EthereumVRFCoordinatorV2_5) CreateSubscription() (*types.Transaction, error) {
	tx, err := v.client.Decode(v.coordinator.CreateSubscription(v.client.NewTXOpts()))
	if err != nil {
		return nil, err
	}
	return tx.Transaction, nil
}

func (v *EthereumVRFCoordinatorV2_5) AddConsumer(subID *big.Int, consumerAddress string) error {
	_, err := v.client.Decode(v.coordinator.AddConsumer(
		v.client.NewTXOpts(),
		subID,
		common.HexToAddress(consumerAddress),
	))
	return err
}

func (v *EthereumVRFCoordinatorV2_5) FundSubscriptionWithNative(subID *big.Int, nativeTokenAmount *big.Int) error {
	opts := v.client.NewTXOpts()
	opts.Value = nativeTokenAmount
	_, err := v.client.Decode(v.coordinator.FundSubscriptionWithNative(opts, subID))
	return err
}

func (v *EthereumVRFCoordinatorV2_5) GetSubscription(ctx context.Context, subID *big.Int) (VRFSubscription, error) {
	sub, err := v.coordinator.GetSubscription(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, subID)
	if err != nil {
		return VRFSubscription{}, err
	}
	return VRFSubscription{
		Balance:       sub.Balance,
		NativeBalance: sub.NativeBalance,
		Owner:         sub.SubOwner,
		Consumers:     sub.Consumers,
	}, nil
}

func (v *EthereumVRFCoordinatorV2_5) GetActiveSubscriptionIDs(ctx context.Context, startIndex, maxCount *big.Int) ([]*big.Int, error) {
	return v.coordinator.GetActiveSubscriptionIds(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, startIndex, maxCount)
}

func (v *EthereumVRFCoordinatorV2_5) PendingRequestsExist(ctx context.Context, subID *big.Int) (bool, error) {
	return v.coordinator.PendingRequestExists(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, subID)
}

func (v *EthereumVRFCoordinatorV2_5) CancelSubscription(subID *big.Int, to common.Address) error {
	_, err := v.client.Decode(v.coordinator.CancelSubscription(v.client.NewTXOpts(), subID, to))
	return err
}

func (v *EthereumVRFCoordinatorV2_5) OwnerCancelSubscription(subID *big.Int) error {
	_, err := v.client.Decode(v.coordinator.OwnerCancelSubscription(v.client.NewTXOpts(), subID))
	return err
}

func (v *EthereumVRFCoordinatorV2_5) Withdraw(recipient common.Address) error {
	_, err := v.client.Decode(v.coordinator.Withdraw(v.client.NewTXOpts(), recipient))
	return err
}

func (v *EthereumVRFCoordinatorV2_5) WithdrawNative(recipient common.Address) error {
	_, err := v.client.Decode(v.coordinator.WithdrawNative(v.client.NewTXOpts(), recipient))
	return err
}

// FilterRandomWordsFulfilled tries to find a fulfilled RandomWords event for the given requestID.
// Returns the event or an error if not yet found.
func (v *EthereumVRFCoordinatorV2_5) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestID *big.Int) (*vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled, error) {
	iter, err := v.coordinator.FilterRandomWordsFulfilled(opts, []*big.Int{requestID}, nil)
	if err != nil {
		return nil, err
	}
	if !iter.Next() {
		return nil, fmt.Errorf("no RandomWordsFulfilled event found for requestID %s", requestID)
	}
	return iter.Event, nil
}

// CountRandomWordsFulfilledLogsInTx counts RandomWordsFulfilled events in a single tx receipt.
func (v *EthereumVRFCoordinatorV2_5) CountRandomWordsFulfilledLogsInTx(ctx context.Context, txHash common.Hash) (int, error) {
	receipt, err := v.client.Client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, l := range receipt.Logs {
		if _, pErr := v.coordinator.ParseRandomWordsFulfilled(*l); pErr == nil {
			count++
		}
	}
	return count, nil
}

func (v *EthereumVRFCoordinatorV2_5) RegisterMigratableCoordinator(newCoordAddr string) error {
	_, err := v.client.Decode(v.coordinator.RegisterMigratableCoordinator(
		v.client.NewTXOpts(),
		common.HexToAddress(newCoordAddr),
	))
	return err
}

func (v *EthereumVRFCoordinatorV2_5) Migrate(subID *big.Int, newCoordAddr string) (*vrf_coordinator_v2_5.VRFCoordinatorV25MigrationCompleted, error) {
	tx, err := v.client.Decode(v.coordinator.Migrate(
		v.client.NewTXOpts(),
		subID,
		common.HexToAddress(newCoordAddr),
	))
	if err != nil {
		return nil, err
	}
	for _, l := range tx.Receipt.Logs {
		event, pErr := v.coordinator.ParseMigrationCompleted(*l)
		if pErr == nil {
			return event, nil
		}
	}
	return nil, errors.New("no MigrationCompleted event found in Migrate receipt")
}

func (v *EthereumVRFCoordinatorV2_5) GetLinkTotalBalance(ctx context.Context) (*big.Int, error) {
	return v.coordinator.STotalBalance(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFCoordinatorV2_5) GetNativeTotalBalance(ctx context.Context) (*big.Int, error) {
	return v.coordinator.STotalNativeBalance(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

// --- Upgraded coordinator methods ---

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) RegisterProvingKey(publicProvingKey [2]*big.Int, gasLaneMaxGas uint64) error {
	_, err := v.client.Decode(v.coordinator.RegisterProvingKey(v.client.NewTXOpts(), publicProvingKey, gasLaneMaxGas))
	return err
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
	_, err := v.client.Decode(v.coordinator.SetConfig(
		v.client.NewTXOpts(),
		minimumRequestConfirmations,
		maxGasLimit,
		stalenessSeconds,
		gasAfterPaymentCalculation,
		fallbackWeiPerUnitLink,
		fulfillmentFlatFeeNativePPM,
		fulfillmentFlatFeeLinkDiscountPPM,
		nativePremiumPercentage,
		linkPremiumPercentage,
	))
	return err
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) SetLINKAndLINKNativeFeed(linkAddress, linkNativeFeedAddress string) error {
	_, err := v.client.Decode(v.coordinator.SetLINKAndLINKNativeFeed(
		v.client.NewTXOpts(),
		common.HexToAddress(linkAddress),
		common.HexToAddress(linkNativeFeedAddress),
	))
	return err
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) GetSubscription(ctx context.Context, subID *big.Int) (VRFSubscription, error) {
	sub, err := v.coordinator.GetSubscription(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, subID)
	if err != nil {
		return VRFSubscription{}, err
	}
	return VRFSubscription{
		Balance:       sub.Balance,
		NativeBalance: sub.NativeBalance,
		Owner:         sub.SubOwner,
		Consumers:     sub.Consumers,
	}, nil
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) GetActiveSubscriptionIDs(ctx context.Context, startIndex, maxCount *big.Int) ([]*big.Int, error) {
	return v.coordinator.GetActiveSubscriptionIds(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, startIndex, maxCount)
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) GetLinkTotalBalance(ctx context.Context) (*big.Int, error) {
	return v.coordinator.STotalBalance(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) GetNativeTotalBalance(ctx context.Context) (*big.Int, error) {
	return v.coordinator.STotalNativeBalance(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestID *big.Int) (*vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled, error) {
	iter, err := v.coordinator.FilterRandomWordsFulfilled(opts, []*big.Int{requestID}, nil)
	if err != nil {
		return nil, err
	}
	if !iter.Next() {
		return nil, fmt.Errorf("no RandomWordsFulfilled event found for requestID %s", requestID)
	}
	return iter.Event, nil
}

// FindSubscriptionID parses the sub ID from the first log of a CreateSubscription receipt.
func FindSubscriptionID(receipt *types.Receipt) (*big.Int, error) {
	if len(receipt.Logs) == 0 {
		return nil, errors.New("no logs in CreateSubscription receipt")
	}
	if len(receipt.Logs[0].Topics) < 2 {
		return nil, errors.New("not enough topics in SubscriptionCreated log")
	}
	return receipt.Logs[0].Topics[1].Big(), nil
}

// --- Consumer methods ---

// RequestRandomness sends a VRF request from the load test consumer and returns the request ID.
func (v *EthereumVRFv2PlusLoadTestConsumer) RequestRandomness(
	keyHash [32]byte,
	subID *big.Int,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	nativePayment bool,
	numWords uint32,
	requestCount uint16,
) (*big.Int, error) {
	tx, err := v.client.Decode(v.consumer.RequestRandomWords(
		v.client.NewTXOpts(),
		subID,
		requestConfirmations,
		keyHash,
		callbackGasLimit,
		nativePayment,
		numWords,
		requestCount,
	))
	if err != nil {
		return nil, err
	}
	return parseRequestIDFromLogs(tx.Receipt.Logs)
}

// RequestRandomnessWithEvent sends a VRF request and returns the full RandomWordsRequested event.
func (v *EthereumVRFv2PlusLoadTestConsumer) RequestRandomnessWithEvent(
	keyHash [32]byte,
	subID *big.Int,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	nativePayment bool,
	numWords uint32,
	requestCount uint16,
) (*vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested, error) {
	tx, err := v.client.Decode(v.consumer.RequestRandomWords(
		v.client.NewTXOpts(),
		subID,
		requestConfirmations,
		keyHash,
		callbackGasLimit,
		nativePayment,
		numWords,
		requestCount,
	))
	if err != nil {
		return nil, err
	}
	return parseRandomWordsRequestedFromLogs(tx.Receipt.Logs)
}

func (v *EthereumVRFv2PlusLoadTestConsumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_v2plus_load_test_with_metrics.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, requestID)
}

func (v *EthereumVRFv2PlusLoadTestConsumer) RequestCount(ctx context.Context) (*big.Int, error) {
	return v.consumer.SRequestCount(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFv2PlusLoadTestConsumer) ResponseCount(ctx context.Context) (*big.Int, error) {
	return v.consumer.SResponseCount(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFv2PlusLoadTestConsumer) GetCoordinator(ctx context.Context) (common.Address, error) {
	return v.consumer.SVrfCoordinator(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

// --- Wrapper methods ---

func (v *EthereumVRFV2PlusWrapper) SetConfig(
	wrapperGasOverhead uint32,
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
	_, err := v.client.Decode(v.wrapper.SetConfig(
		v.client.NewTXOpts(),
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
	))
	return err
}

func (v *EthereumVRFV2PlusWrapper) GetSubID(ctx context.Context) (*big.Int, error) {
	return v.wrapper.SUBSCRIPTIONID(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFV2PlusWrapper) Coordinator(ctx context.Context) (common.Address, error) {
	return v.wrapper.SVrfCoordinator(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

// --- Wrapper consumer methods ---

// RequestRandomWords sends a LINK-funded wrapper request and returns the request ID.
func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) RequestRandomWords(
	requestConfirmations uint16,
	callbackGasLimit uint32,
	numWords uint32,
	requestCount uint16,
) (*big.Int, error) {
	tx, err := v.client.Decode(v.consumer.MakeRequests(
		v.client.NewTXOpts(),
		callbackGasLimit,
		requestConfirmations,
		numWords,
		requestCount,
	))
	if err != nil {
		return nil, err
	}
	return parseRequestIDFromLogs(tx.Receipt.Logs)
}

// RequestRandomWordsNative sends a native-funded wrapper request and returns the request ID.
func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) RequestRandomWordsNative(
	requestConfirmations uint16,
	callbackGasLimit uint32,
	numWords uint32,
	requestCount uint16,
) (*big.Int, error) {
	tx, err := v.client.Decode(v.consumer.MakeRequestsNative(
		v.client.NewTXOpts(),
		callbackGasLimit,
		requestConfirmations,
		numWords,
		requestCount,
	))
	if err != nil {
		return nil, err
	}
	return parseRequestIDFromLogs(tx.Receipt.Logs)
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrfv2plus_wrapper_load_test_consumer.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, requestID)
}

// --- Deploy functions ---

func DeployBlockhashStore(client *seth.Client) (*EthereumBlockhashStore, error) {
	abi, err := blockhash_store.BlockhashStoreMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get BlockhashStore ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "BlockhashStore", *abi, common.FromHex(blockhash_store.BlockhashStoreMetaData.Bin))
	if err != nil {
		return nil, fmt.Errorf("BlockhashStore deployment failed: %w", err)
	}
	instance, err := blockhash_store.NewBlockhashStore(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate BlockhashStore: %w", err)
	}
	return &EthereumBlockhashStore{client: client, store: instance, address: data.Address}, nil
}

func DeployBatchBlockhashStore(client *seth.Client, bhsAddr string) (*EthereumBatchBlockhashStore, error) {
	abi, err := batch_blockhash_store.BatchBlockhashStoreMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get BatchBlockhashStore ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "BatchBlockhashStore", *abi,
		common.FromHex(batch_blockhash_store.BatchBlockhashStoreMetaData.Bin),
		common.HexToAddress(bhsAddr))
	if err != nil {
		return nil, fmt.Errorf("BatchBlockhashStore deployment failed: %w", err)
	}
	instance, err := batch_blockhash_store.NewBatchBlockhashStore(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate BatchBlockhashStore: %w", err)
	}
	return &EthereumBatchBlockhashStore{client: client, store: instance, address: data.Address}, nil
}

func DeployVRFCoordinatorV2_5(client *seth.Client, bhsAddr string) (*EthereumVRFCoordinatorV2_5, error) {
	abi, err := vrf_coordinator_v2_5.VRFCoordinatorV25MetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFCoordinatorV2_5 ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "VRFCoordinatorV2_5", *abi,
		common.FromHex(vrf_coordinator_v2_5.VRFCoordinatorV25MetaData.Bin),
		common.HexToAddress(bhsAddr))
	if err != nil {
		return nil, fmt.Errorf("VRFCoordinatorV2_5 deployment failed: %w", err)
	}
	instance, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFCoordinatorV2_5: %w", err)
	}
	return &EthereumVRFCoordinatorV2_5{client: client, coordinator: instance, address: data.Address}, nil
}

func DeployBatchVRFCoordinatorV2Plus(client *seth.Client, coordAddr string) (*EthereumBatchVRFCoordinatorV2Plus, error) {
	abi, err := batch_vrf_coordinator_v2plus.BatchVRFCoordinatorV2PlusMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get BatchVRFCoordinatorV2Plus ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "BatchVRFCoordinatorV2Plus", *abi,
		common.FromHex(batch_vrf_coordinator_v2plus.BatchVRFCoordinatorV2PlusMetaData.Bin),
		common.HexToAddress(coordAddr))
	if err != nil {
		return nil, fmt.Errorf("BatchVRFCoordinatorV2Plus deployment failed: %w", err)
	}
	instance, err := batch_vrf_coordinator_v2plus.NewBatchVRFCoordinatorV2Plus(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate BatchVRFCoordinatorV2Plus: %w", err)
	}
	return &EthereumBatchVRFCoordinatorV2Plus{client: client, coordinator: instance, address: data.Address}, nil
}

func DeployVRFv2PlusLoadTestConsumer(client *seth.Client, coordAddr string) (*EthereumVRFv2PlusLoadTestConsumer, error) {
	abi, err := vrf_v2plus_load_test_with_metrics.VRFV2PlusLoadTestWithMetricsMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFV2PlusLoadTestWithMetrics ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "VRFV2PlusLoadTestWithMetrics", *abi,
		common.FromHex(vrf_v2plus_load_test_with_metrics.VRFV2PlusLoadTestWithMetricsMetaData.Bin),
		common.HexToAddress(coordAddr))
	if err != nil {
		return nil, fmt.Errorf("VRFV2PlusLoadTestWithMetrics deployment failed: %w", err)
	}
	instance, err := vrf_v2plus_load_test_with_metrics.NewVRFV2PlusLoadTestWithMetrics(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFV2PlusLoadTestWithMetrics: %w", err)
	}
	return &EthereumVRFv2PlusLoadTestConsumer{client: client, consumer: instance, address: data.Address}, nil
}

func DeployVRFV2PlusWrapper(client *seth.Client, linkAddr, feedAddr, coordAddr string, subID *big.Int) (*EthereumVRFV2PlusWrapper, error) {
	abi, err := vrfv2plus_wrapper.VRFV2PlusWrapperMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFV2PlusWrapper ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "VRFV2PlusWrapper", *abi,
		common.FromHex(vrfv2plus_wrapper.VRFV2PlusWrapperMetaData.Bin),
		common.HexToAddress(linkAddr),
		common.HexToAddress(feedAddr),
		common.HexToAddress(coordAddr),
		subID)
	if err != nil {
		return nil, fmt.Errorf("VRFV2PlusWrapper deployment failed: %w", err)
	}
	instance, err := vrfv2plus_wrapper.NewVRFV2PlusWrapper(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFV2PlusWrapper: %w", err)
	}
	return &EthereumVRFV2PlusWrapper{client: client, wrapper: instance, address: data.Address}, nil
}

func DeployVRFCoordinatorV2PlusUpgradedVersion(client *seth.Client, bhsAddr string) (*EthereumVRFCoordinatorV2PlusUpgradedVersion, error) {
	abi, err := vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFCoordinatorV2PlusUpgradedVersion ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "VRFCoordinatorV2PlusUpgradedVersion", *abi,
		common.FromHex(vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionMetaData.Bin),
		common.HexToAddress(bhsAddr))
	if err != nil {
		return nil, fmt.Errorf("VRFCoordinatorV2PlusUpgradedVersion deployment failed: %w", err)
	}
	instance, err := vrf_v2plus_upgraded_version.NewVRFCoordinatorV2PlusUpgradedVersion(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFCoordinatorV2PlusUpgradedVersion: %w", err)
	}
	return &EthereumVRFCoordinatorV2PlusUpgradedVersion{client: client, coordinator: instance, address: data.Address}, nil
}

func DeployVRFV2PlusWrapperLoadTestConsumer(client *seth.Client, wrapperAddr string) (*EthereumVRFV2PlusWrapperLoadTestConsumer, error) {
	abi, err := vrfv2plus_wrapper_load_test_consumer.VRFV2PlusWrapperLoadTestConsumerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFV2PlusWrapperLoadTestConsumer ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "VRFV2PlusWrapperLoadTestConsumer", *abi,
		common.FromHex(vrfv2plus_wrapper_load_test_consumer.VRFV2PlusWrapperLoadTestConsumerMetaData.Bin),
		common.HexToAddress(wrapperAddr))
	if err != nil {
		return nil, fmt.Errorf("VRFV2PlusWrapperLoadTestConsumer deployment failed: %w", err)
	}
	instance, err := vrfv2plus_wrapper_load_test_consumer.NewVRFV2PlusWrapperLoadTestConsumer(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFV2PlusWrapperLoadTestConsumer: %w", err)
	}
	return &EthereumVRFV2PlusWrapperLoadTestConsumer{client: client, consumer: instance, address: data.Address}, nil
}

// --- Load functions ---

func LoadBlockhashStore(client *seth.Client, addr string) (*EthereumBlockhashStore, error) {
	address := common.HexToAddress(addr)
	abi, err := blockhash_store.BlockhashStoreMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get BlockhashStore ABI: %w", err)
	}
	client.ContractStore.AddABI("BlockhashStore", *abi)
	client.ContractStore.AddBIN("BlockhashStore", common.FromHex(blockhash_store.BlockhashStoreMetaData.Bin))
	instance, err := blockhash_store.NewBlockhashStore(address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to load BlockhashStore: %w", err)
	}
	return &EthereumBlockhashStore{client: client, store: instance, address: address}, nil
}

func (v *EthereumBlockhashStore) GetBlockhash(ctx context.Context, blockNumber uint64) ([32]byte, error) {
	return v.store.GetBlockhash(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, new(big.Int).SetUint64(blockNumber))
}

func LoadBatchVRFCoordinatorV2Plus(client *seth.Client, addr string) (*EthereumBatchVRFCoordinatorV2Plus, error) {
	address := common.HexToAddress(addr)
	abi, err := batch_vrf_coordinator_v2plus.BatchVRFCoordinatorV2PlusMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get BatchVRFCoordinatorV2Plus ABI: %w", err)
	}
	client.ContractStore.AddABI("BatchVRFCoordinatorV2Plus", *abi)
	client.ContractStore.AddBIN("BatchVRFCoordinatorV2Plus", common.FromHex(batch_vrf_coordinator_v2plus.BatchVRFCoordinatorV2PlusMetaData.Bin))
	instance, err := batch_vrf_coordinator_v2plus.NewBatchVRFCoordinatorV2Plus(address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to load BatchVRFCoordinatorV2Plus: %w", err)
	}
	return &EthereumBatchVRFCoordinatorV2Plus{client: client, coordinator: instance, address: address}, nil
}

func LoadVRFCoordinatorV2_5(client *seth.Client, addr string) (*EthereumVRFCoordinatorV2_5, error) {
	address := common.HexToAddress(addr)
	abi, err := vrf_coordinator_v2_5.VRFCoordinatorV25MetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFCoordinatorV2_5 ABI: %w", err)
	}
	client.ContractStore.AddABI("VRFCoordinatorV2_5", *abi)
	client.ContractStore.AddBIN("VRFCoordinatorV2_5", common.FromHex(vrf_coordinator_v2_5.VRFCoordinatorV25MetaData.Bin))
	instance, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to load VRFCoordinatorV2_5: %w", err)
	}
	return &EthereumVRFCoordinatorV2_5{client: client, coordinator: instance, address: address}, nil
}

func LoadVRFv2PlusLoadTestConsumer(client *seth.Client, addr string) (*EthereumVRFv2PlusLoadTestConsumer, error) {
	address := common.HexToAddress(addr)
	abi, err := vrf_v2plus_load_test_with_metrics.VRFV2PlusLoadTestWithMetricsMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFV2PlusLoadTestWithMetrics ABI: %w", err)
	}
	client.ContractStore.AddABI("VRFV2PlusLoadTestWithMetrics", *abi)
	client.ContractStore.AddBIN("VRFV2PlusLoadTestWithMetrics", common.FromHex(vrf_v2plus_load_test_with_metrics.VRFV2PlusLoadTestWithMetricsMetaData.Bin))
	instance, err := vrf_v2plus_load_test_with_metrics.NewVRFV2PlusLoadTestWithMetrics(address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to load VRFV2PlusLoadTestWithMetrics: %w", err)
	}
	return &EthereumVRFv2PlusLoadTestConsumer{client: client, consumer: instance, address: address}, nil
}

func LoadVRFV2PlusWrapper(client *seth.Client, addr string) (*EthereumVRFV2PlusWrapper, error) {
	address := common.HexToAddress(addr)
	abi, err := vrfv2plus_wrapper.VRFV2PlusWrapperMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFV2PlusWrapper ABI: %w", err)
	}
	client.ContractStore.AddABI("VRFV2PlusWrapper", *abi)
	client.ContractStore.AddBIN("VRFV2PlusWrapper", common.FromHex(vrfv2plus_wrapper.VRFV2PlusWrapperMetaData.Bin))
	instance, err := vrfv2plus_wrapper.NewVRFV2PlusWrapper(address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to load VRFV2PlusWrapper: %w", err)
	}
	return &EthereumVRFV2PlusWrapper{client: client, wrapper: instance, address: address}, nil
}

func LoadVRFV2PlusWrapperLoadTestConsumer(client *seth.Client, addr string) (*EthereumVRFV2PlusWrapperLoadTestConsumer, error) {
	address := common.HexToAddress(addr)
	abi, err := vrfv2plus_wrapper_load_test_consumer.VRFV2PlusWrapperLoadTestConsumerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFV2PlusWrapperLoadTestConsumer ABI: %w", err)
	}
	client.ContractStore.AddABI("VRFV2PlusWrapperLoadTestConsumer", *abi)
	client.ContractStore.AddBIN("VRFV2PlusWrapperLoadTestConsumer", common.FromHex(vrfv2plus_wrapper_load_test_consumer.VRFV2PlusWrapperLoadTestConsumerMetaData.Bin))
	instance, err := vrfv2plus_wrapper_load_test_consumer.NewVRFV2PlusWrapperLoadTestConsumer(address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to load VRFV2PlusWrapperLoadTestConsumer: %w", err)
	}
	return &EthereumVRFV2PlusWrapperLoadTestConsumer{client: client, consumer: instance, address: address}, nil
}

// parseRequestIDFromLogs parses the RandomWordsRequested event from receipt logs and returns the requestID.
func parseRequestIDFromLogs(logs []*types.Log) (*big.Int, error) {
	event, err := parseRandomWordsRequestedFromLogs(logs)
	if err != nil {
		return nil, err
	}
	return event.RequestId, nil
}

func parseRandomWordsRequestedFromLogs(logs []*types.Log) (*vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested, error) {
	coordABI, err := vrf_coordinator_v2_5.VRFCoordinatorV25MetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get coordinator ABI: %w", err)
	}
	bc := bind.NewBoundContract(common.Address{}, *coordABI, nil, nil, nil)
	for _, log := range logs {
		event := new(vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested)
		if pErr := bc.UnpackLog(event, "RandomWordsRequested", *log); pErr != nil {
			continue
		}
		event.Raw = *log
		return event, nil
	}
	return nil, errors.New("no RandomWordsRequested event found in receipt")
}

// EncodeOnChainVRFProvingKey encodes an uncompressed VRF public key (0x + 128 hex chars) to [2]*big.Int.
func EncodeOnChainVRFProvingKey(uncompressed string) ([2]*big.Int, error) {
	// strip 0x prefix
	raw := uncompressed
	if len(raw) > 1 && raw[:2] == "0x" {
		raw = raw[2:]
	}
	if len(raw) != 128 {
		// also accept 130 char uncompressed point (04 prefix + 128 hex)
		if len(raw) != 130 || raw[:2] != "04" {
			return [2]*big.Int{}, fmt.Errorf("unexpected uncompressed key length %d (expected 128 hex chars after 0x)", len(raw))
		}
		raw = raw[2:]
	}
	if _, err := hex.DecodeString(raw); err != nil {
		return [2]*big.Int{}, fmt.Errorf("invalid hex in VRF key: %w", err)
	}
	provingKey := [2]*big.Int{}
	var set1, set2 bool
	provingKey[0], set1 = new(big.Int).SetString(raw[:64], 16)
	if !set1 {
		return [2]*big.Int{}, errors.New("cannot convert VRF key X coordinate to *big.Int")
	}
	provingKey[1], set2 = new(big.Int).SetString(raw[64:], 16)
	if !set2 {
		return [2]*big.Int{}, errors.New("cannot convert VRF key Y coordinate to *big.Int")
	}
	return provingKey, nil
}
