package contracts

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_load_test_with_metrics"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2_wrapper_load_test_consumer"
)

// CoordinatorV2Log is implemented by EthereumVRFCoordinatorV2 for log parsing and waits.
type CoordinatorV2Log interface {
	ParseRandomWordsRequested(log types.Log) (*CoordinatorRandomWordsRequested, error)
	ParseRandomWordsFulfilled(log types.Log) (*CoordinatorRandomWordsFulfilled, error)
	Address() string
	WaitForRandomWordsFulfilledEvent(filter RandomWordsFulfilledEventFilter) (*CoordinatorRandomWordsFulfilled, error)
	WaitForConfigSetEvent(timeout time.Duration) (*CoordinatorConfigSet, error)
	FilterRandomWordsFulfilledEvent(opts *bind.FilterOpts, requestID *big.Int) (*CoordinatorRandomWordsFulfilled, error)
}

// CoordinatorRandomWordsFulfilled mirrors integration-tests coordinator fulfilled shape.
type CoordinatorRandomWordsFulfilled struct {
	RequestID     *big.Int
	OutputSeed    *big.Int
	SubID         string
	Payment       *big.Int
	NativePayment bool
	Success       bool
	OnlyPremium   bool
	Raw           types.Log
}

// CoordinatorRandomWordsRequested mirrors integration-tests coordinator requested shape.
type CoordinatorRandomWordsRequested struct {
	KeyHash                     [32]byte
	RequestID                   *big.Int
	PreSeed                     *big.Int
	SubID                       string
	MinimumRequestConfirmations uint16
	CallbackGasLimit            uint32
	NumWords                    uint32
	ExtraArgs                   []byte
	Sender                      common.Address
	Raw                         types.Log
}

// CoordinatorConfigSet is emitted when coordinator config changes.
type CoordinatorConfigSet struct {
	MinimumRequestConfirmations       uint16
	MaxGasLimit                       uint32
	StalenessSeconds                  uint32
	GasAfterPaymentCalculation        uint32
	FallbackWeiPerUnitLink            *big.Int
	FulfillmentFlatFeeNativePPM       uint32
	FulfillmentFlatFeeLinkDiscountPPM uint32
	NativePremiumPercentage           uint8
	LinkPremiumPercentage             uint8
	FeeConfig                         VRFCoordinatorV2OnChainFeeConfig
	Raw                               types.Log
}

// VRFCoordinatorV2OnChainFeeConfig is the coordinator fee config struct.
type VRFCoordinatorV2OnChainFeeConfig struct {
	FulfillmentFlatFeeLinkPPMTier1 uint32
	FulfillmentFlatFeeLinkPPMTier2 uint32
	FulfillmentFlatFeeLinkPPMTier3 uint32
	FulfillmentFlatFeeLinkPPMTier4 uint32
	FulfillmentFlatFeeLinkPPMTier5 uint32
	ReqsForTier2                   *big.Int
	ReqsForTier3                   *big.Int
	ReqsForTier4                   *big.Int
	ReqsForTier5                   *big.Int
}

// RandomWordsFulfilledEventFilter filters fulfilled events.
type RandomWordsFulfilledEventFilter struct {
	RequestIDs []*big.Int
	SubIDs     []*big.Int
	Timeout    time.Duration
}

// VRFv2Subscription is returned by GetSubscription on the V2 coordinator.
type VRFv2Subscription struct {
	Balance       *big.Int
	NativeBalance *big.Int
	ReqCount      uint64
	SubOwner      common.Address
	Consumers     []common.Address
}

// VRFLoadTestMetrics holds consumer load-test counters.
type VRFLoadTestMetrics struct {
	RequestCount                 *big.Int
	FulfilmentCount              *big.Int
	AverageFulfillmentInMillions *big.Int
	SlowestFulfillment           *big.Int
	FastestFulfillment           *big.Int
}

// --- Coordinator V2 ---

// EthereumVRFCoordinatorV2 wraps vrf_coordinator_v2.VRFCoordinatorV2.
type EthereumVRFCoordinatorV2 struct {
	client      *seth.Client
	coordinator *vrf_coordinator_v2.VRFCoordinatorV2
	address     common.Address
}

func (v *EthereumVRFCoordinatorV2) Address() string { return v.address.Hex() }

func (v *EthereumVRFCoordinatorV2) Coordinator() *vrf_coordinator_v2.VRFCoordinatorV2 {
	return v.coordinator
}

// DeployVRFCoordinatorV2 deploys the VRF Coordinator V2 contract.
func DeployVRFCoordinatorV2(client *seth.Client, linkAddr, bhsAddr, linkEthFeedAddr string) (*EthereumVRFCoordinatorV2, error) {
	abi, err := vrf_coordinator_v2.VRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFCoordinatorV2 ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "VRFCoordinatorV2", *abi,
		common.FromHex(vrf_coordinator_v2.VRFCoordinatorV2MetaData.Bin),
		common.HexToAddress(linkAddr),
		common.HexToAddress(bhsAddr),
		common.HexToAddress(linkEthFeedAddr))
	if err != nil {
		return nil, fmt.Errorf("VRFCoordinatorV2 deployment failed: %w", err)
	}
	instance, err := vrf_coordinator_v2.NewVRFCoordinatorV2(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFCoordinatorV2: %w", err)
	}
	return &EthereumVRFCoordinatorV2{client: client, coordinator: instance, address: data.Address}, nil
}

// LoadVRFCoordinatorV2 binds an existing VRF Coordinator V2.
func LoadVRFCoordinatorV2(client *seth.Client, addr string) (*EthereumVRFCoordinatorV2, error) {
	address := common.HexToAddress(addr)
	abi, err := vrf_coordinator_v2.VRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFCoordinatorV2 ABI: %w", err)
	}
	client.ContractStore.AddABI("VRFCoordinatorV2", *abi)
	client.ContractStore.AddBIN("VRFCoordinatorV2", common.FromHex(vrf_coordinator_v2.VRFCoordinatorV2MetaData.Bin))
	instance, err := vrf_coordinator_v2.NewVRFCoordinatorV2(address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to load VRFCoordinatorV2: %w", err)
	}
	return &EthereumVRFCoordinatorV2{client: client, coordinator: instance, address: address}, nil
}

func (v *EthereumVRFCoordinatorV2) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
	return v.coordinator.HashOfKey(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx}, pubKey)
}

func (v *EthereumVRFCoordinatorV2) GetSubscription(ctx context.Context, subID uint64) (VRFv2Subscription, error) {
	sub, err := v.coordinator.GetSubscription(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx}, subID)
	if err != nil {
		return VRFv2Subscription{}, err
	}
	return VRFv2Subscription{
		Balance:       sub.Balance,
		NativeBalance: nil,
		ReqCount:      sub.ReqCount,
		SubOwner:      sub.Owner,
		Consumers:     sub.Consumers,
	}, nil
}

func (v *EthereumVRFCoordinatorV2) SetConfig(
	minimumRequestConfirmations uint16,
	maxGasLimit uint32,
	stalenessSeconds uint32,
	gasAfterPaymentCalculation uint32,
	fallbackWeiPerUnitLink *big.Int,
	feeConfig vrf_coordinator_v2.VRFCoordinatorV2FeeConfig,
) error {
	_, err := v.client.Decode(v.coordinator.SetConfig(
		v.client.NewTXOpts(),
		minimumRequestConfirmations,
		maxGasLimit,
		stalenessSeconds,
		gasAfterPaymentCalculation,
		fallbackWeiPerUnitLink,
		feeConfig,
	))
	return err
}

func (v *EthereumVRFCoordinatorV2) RegisterProvingKey(oracleAddr string, publicProvingKey [2]*big.Int) error {
	_, err := v.client.Decode(v.coordinator.RegisterProvingKey(
		v.client.NewTXOpts(),
		common.HexToAddress(oracleAddr),
		publicProvingKey,
	))
	return err
}

func (v *EthereumVRFCoordinatorV2) CreateSubscription() (*types.Receipt, error) {
	tx, err := v.client.Decode(v.coordinator.CreateSubscription(v.client.NewTXOpts()))
	if err != nil {
		return nil, err
	}
	return tx.Receipt, err
}

func (v *EthereumVRFCoordinatorV2) AddConsumer(subID uint64, consumerAddress string) error {
	_, err := v.client.Decode(v.coordinator.AddConsumer(
		v.client.NewTXOpts(),
		subID,
		common.HexToAddress(consumerAddress),
	))
	return err
}

func (v *EthereumVRFCoordinatorV2) PendingRequestsExist(ctx context.Context, subID uint64) (bool, error) {
	return v.coordinator.PendingRequestExists(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx}, subID)
}

func (v *EthereumVRFCoordinatorV2) OracleWithdraw(recipient common.Address, amount *big.Int) error {
	_, err := v.client.Decode(v.coordinator.OracleWithdraw(v.client.NewTXOpts(), recipient, amount))
	return err
}

func (v *EthereumVRFCoordinatorV2) CancelSubscription(subID uint64, to common.Address) (*seth.DecodedTransaction, *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled, error) {
	tx, err := v.client.Decode(v.coordinator.CancelSubscription(v.client.NewTXOpts(), subID, to))
	if err != nil {
		return nil, nil, err
	}
	var canceled *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled
	for _, lg := range tx.Receipt.Logs {
		for _, topic := range lg.Topics {
			if topic.Cmp(vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled{}.Topic()) == 0 {
				canceled, err = v.coordinator.ParseSubscriptionCanceled(*lg)
				if err != nil {
					return nil, nil, fmt.Errorf("parse SubscriptionCanceled: %w", err)
				}
			}
		}
	}
	if canceled == nil {
		return tx, nil, errors.New("no SubscriptionCanceled event in transaction receipt logs")
	}
	return tx, canceled, nil
}

func (v *EthereumVRFCoordinatorV2) OwnerCancelSubscription(subID uint64) (*seth.DecodedTransaction, *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled, error) {
	tx, err := v.client.Decode(v.coordinator.OwnerCancelSubscription(v.client.NewTXOpts(), subID))
	if err != nil {
		return nil, nil, err
	}
	var canceled *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled
	for _, lg := range tx.Receipt.Logs {
		for _, topic := range lg.Topics {
			if topic.Cmp(vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled{}.Topic()) == 0 {
				canceled, err = v.coordinator.ParseSubscriptionCanceled(*lg)
				if err != nil {
					return nil, nil, fmt.Errorf("parse SubscriptionCanceled: %w", err)
				}
			}
		}
	}
	if canceled == nil {
		return tx, nil, errors.New("no SubscriptionCanceled event in transaction receipt logs")
	}
	return tx, canceled, nil
}

func (v *EthereumVRFCoordinatorV2) GetConfig(ctx context.Context) (vrf_coordinator_v2.GetConfig, error) {
	return v.coordinator.GetConfig(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx})
}

func (v *EthereumVRFCoordinatorV2) GetFeeConfig(ctx context.Context) (vrf_coordinator_v2.GetFeeConfig, error) {
	return v.coordinator.GetFeeConfig(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx})
}

func (v *EthereumVRFCoordinatorV2) GetFallbackWeiPerUnitLink(ctx context.Context) (*big.Int, error) {
	return v.coordinator.GetFallbackWeiPerUnitLink(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx})
}

func (v *EthereumVRFCoordinatorV2) ParseRandomWordsRequested(log types.Log) (*CoordinatorRandomWordsRequested, error) {
	ev, err := v.coordinator.ParseRandomWordsRequested(log)
	if err != nil {
		return nil, err
	}
	return &CoordinatorRandomWordsRequested{
		KeyHash:                     ev.KeyHash,
		RequestID:                   ev.RequestId,
		PreSeed:                     ev.PreSeed,
		SubID:                       strconv.FormatUint(ev.SubId, 10),
		MinimumRequestConfirmations: ev.MinimumRequestConfirmations,
		CallbackGasLimit:            ev.CallbackGasLimit,
		NumWords:                    ev.NumWords,
		Sender:                      ev.Sender,
		Raw:                         ev.Raw,
	}, nil
}

func (v *EthereumVRFCoordinatorV2) ParseRandomWordsFulfilled(log types.Log) (*CoordinatorRandomWordsFulfilled, error) {
	ev, err := v.coordinator.ParseRandomWordsFulfilled(log)
	if err != nil {
		return nil, err
	}
	return &CoordinatorRandomWordsFulfilled{
		RequestID:  ev.RequestId,
		OutputSeed: ev.OutputSeed,
		Payment:    ev.Payment,
		Success:    ev.Success,
		Raw:        ev.Raw,
	}, nil
}

func (v *EthereumVRFCoordinatorV2) FilterRandomWordsFulfilledEvent(opts *bind.FilterOpts, requestID *big.Int) (*CoordinatorRandomWordsFulfilled, error) {
	it, err := v.coordinator.FilterRandomWordsFulfilled(opts, []*big.Int{requestID})
	if err != nil {
		return nil, err
	}
	if !it.Next() {
		return nil, fmt.Errorf("no RandomWordsFulfilled for request %s", requestID.String())
	}
	ev := it.Event
	return &CoordinatorRandomWordsFulfilled{
		RequestID:  ev.RequestId,
		OutputSeed: ev.OutputSeed,
		Payment:    ev.Payment,
		Success:    ev.Success,
		Raw:        ev.Raw,
	}, nil
}

func (v *EthereumVRFCoordinatorV2) WaitForRandomWordsFulfilledEvent(filter RandomWordsFulfilledEventFilter) (*CoordinatorRandomWordsFulfilled, error) {
	ch := make(chan *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled)
	sub, err := v.coordinator.WatchRandomWordsFulfilled(nil, ch, filter.RequestIDs)
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			return nil, err
		case <-time.After(filter.Timeout):
			return nil, errors.New("timeout waiting for RandomWordsFulfilled event")
		case ev := <-ch:
			return &CoordinatorRandomWordsFulfilled{
				RequestID:  ev.RequestId,
				OutputSeed: ev.OutputSeed,
				Payment:    ev.Payment,
				Success:    ev.Success,
				Raw:        ev.Raw,
			}, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2) WaitForConfigSetEvent(timeout time.Duration) (*CoordinatorConfigSet, error) {
	ch := make(chan *vrf_coordinator_v2.VRFCoordinatorV2ConfigSet)
	sub, err := v.coordinator.WatchConfigSet(nil, ch)
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			return nil, err
		case <-time.After(timeout):
			return nil, errors.New("timeout waiting for ConfigSet event")
		case ev := <-ch:
			return &CoordinatorConfigSet{
				MinimumRequestConfirmations: ev.MinimumRequestConfirmations,
				MaxGasLimit:                 ev.MaxGasLimit,
				StalenessSeconds:            ev.StalenessSeconds,
				GasAfterPaymentCalculation:  ev.GasAfterPaymentCalculation,
				FallbackWeiPerUnitLink:      ev.FallbackWeiPerUnitLink,
				FeeConfig: VRFCoordinatorV2OnChainFeeConfig{
					FulfillmentFlatFeeLinkPPMTier1: ev.FeeConfig.FulfillmentFlatFeeLinkPPMTier1,
					FulfillmentFlatFeeLinkPPMTier2: ev.FeeConfig.FulfillmentFlatFeeLinkPPMTier2,
					FulfillmentFlatFeeLinkPPMTier3: ev.FeeConfig.FulfillmentFlatFeeLinkPPMTier3,
					FulfillmentFlatFeeLinkPPMTier4: ev.FeeConfig.FulfillmentFlatFeeLinkPPMTier4,
					FulfillmentFlatFeeLinkPPMTier5: ev.FeeConfig.FulfillmentFlatFeeLinkPPMTier5,
					ReqsForTier2:                   ev.FeeConfig.ReqsForTier2,
					ReqsForTier3:                   ev.FeeConfig.ReqsForTier3,
					ReqsForTier4:                   ev.FeeConfig.ReqsForTier4,
					ReqsForTier5:                   ev.FeeConfig.ReqsForTier5,
				},
			}, nil
		}
	}
}

// ParseRandomWordsFulfilledLogs parses all RandomWordsFulfilled logs in a receipt.
func ParseRandomWordsFulfilledLogs(coordinator CoordinatorV2Log, logs []*types.Log) ([]*CoordinatorRandomWordsFulfilled, error) {
	var out []*CoordinatorRandomWordsFulfilled
	for _, lg := range logs {
		for _, topic := range lg.Topics {
			if topic.Cmp(vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled{}.Topic()) == 0 {
				ev, err := coordinator.ParseRandomWordsFulfilled(*lg)
				if err != nil {
					return nil, err
				}
				out = append(out, ev)
			}
		}
	}
	return out, nil
}

func parseRequestRandomnessLogs(coordinator CoordinatorV2Log, logs []*types.Log) (*CoordinatorRandomWordsRequested, error) {
	var requested *CoordinatorRandomWordsRequested
	var err error
	for _, lg := range logs {
		for _, topic := range lg.Topics {
			if topic.Cmp(vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{}.Topic()) == 0 {
				requested, err = coordinator.ParseRandomWordsRequested(*lg)
				if err != nil {
					return nil, fmt.Errorf("parse RandomWordsRequested: %w", err)
				}
			}
		}
	}
	if requested == nil {
		return nil, errors.New("no RandomWordsRequested event in transaction receipt logs")
	}
	return requested, nil
}

// FindVRFv2SubscriptionID returns the uint64 sub ID from a CreateSubscription receipt by locating SubscriptionCreated.
func FindVRFv2SubscriptionID(receipt *types.Receipt) (uint64, error) {
	wantTopic := vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCreated{}.Topic()
	for _, lg := range receipt.Logs {
		if len(lg.Topics) < 2 {
			continue
		}
		if lg.Topics[0].Cmp(wantTopic) == 0 {
			return lg.Topics[1].Big().Uint64(), nil
		}
	}
	return 0, errors.New("no SubscriptionCreated event in transaction receipt logs")
}

// FallbackWeiBigInt parses decimal string fallback wei per unit LINK.
func FallbackWeiBigInt(s string) (*big.Int, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}
	return d.BigInt(), nil
}

// --- Batch coordinator V2 ---

// EthereumBatchVRFCoordinatorV2 wraps batch_vrf_coordinator_v2.
type EthereumBatchVRFCoordinatorV2 struct {
	client           *seth.Client
	batchCoordinator *batch_vrf_coordinator_v2.BatchVRFCoordinatorV2
	address          common.Address
}

func (v *EthereumBatchVRFCoordinatorV2) Address() string { return v.address.Hex() }

// DeployBatchVRFCoordinatorV2 deploys the batch coordinator pointing at VRF Coordinator V2.
func DeployBatchVRFCoordinatorV2(client *seth.Client, coordinatorAddress string) (*EthereumBatchVRFCoordinatorV2, error) {
	abi, err := batch_vrf_coordinator_v2.BatchVRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get BatchVRFCoordinatorV2 ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "BatchVRFCoordinatorV2", *abi,
		common.FromHex(batch_vrf_coordinator_v2.BatchVRFCoordinatorV2MetaData.Bin),
		common.HexToAddress(coordinatorAddress))
	if err != nil {
		return nil, fmt.Errorf("BatchVRFCoordinatorV2 deployment failed: %w", err)
	}
	instance, err := batch_vrf_coordinator_v2.NewBatchVRFCoordinatorV2(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate BatchVRFCoordinatorV2: %w", err)
	}
	return &EthereumBatchVRFCoordinatorV2{client: client, batchCoordinator: instance, address: data.Address}, nil
}

// --- Load test consumer ---

// EthereumVRFv2LoadTestConsumer wraps VRFV2LoadTestWithMetrics.
type EthereumVRFv2LoadTestConsumer struct {
	client   *seth.Client
	consumer *vrf_load_test_with_metrics.VRFV2LoadTestWithMetrics
	address  common.Address
}

func (v *EthereumVRFv2LoadTestConsumer) Address() string { return v.address.Hex() }

// DeployVRFv2LoadTestConsumer deploys a VRF v2 load test consumer.
func DeployVRFv2LoadTestConsumer(client *seth.Client, coordinatorAddr string) (*EthereumVRFv2LoadTestConsumer, error) {
	abi, err := vrf_load_test_with_metrics.VRFV2LoadTestWithMetricsMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFV2LoadTestWithMetrics ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "VRFV2LoadTestWithMetrics", *abi,
		common.FromHex(vrf_load_test_with_metrics.VRFV2LoadTestWithMetricsMetaData.Bin),
		common.HexToAddress(coordinatorAddr))
	if err != nil {
		return nil, fmt.Errorf("VRFV2LoadTestWithMetrics deployment failed: %w", err)
	}
	instance, err := vrf_load_test_with_metrics.NewVRFV2LoadTestWithMetrics(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFV2LoadTestWithMetrics: %w", err)
	}
	return &EthereumVRFv2LoadTestConsumer{client: client, consumer: instance, address: data.Address}, nil
}

func (v *EthereumVRFv2LoadTestConsumer) RequestRandomnessFromKey(
	coordinator CoordinatorV2Log,
	keyHash [32]byte,
	subID uint64,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	numWords uint32,
	requestCount uint16,
	keyNum int,
) (*CoordinatorRandomWordsRequested, error) {
	tx, err := v.client.Decode(v.consumer.RequestRandomWords(
		v.client.NewTXKeyOpts(keyNum),
		subID,
		requestConfirmations,
		keyHash,
		callbackGasLimit,
		numWords,
		requestCount,
	)) // matches integration-tests RequestRandomWordsFromKey argument order
	if err != nil {
		return nil, fmt.Errorf("RequestRandomWords: %w", err)
	}
	return parseRequestRandomnessLogs(coordinator, tx.Receipt.Logs)
}

func (v *EthereumVRFv2LoadTestConsumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_load_test_with_metrics.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx}, requestID)
}

func (v *EthereumVRFv2LoadTestConsumer) GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error) {
	reqCount, err := v.consumer.SRequestCount(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx})
	if err != nil {
		return nil, err
	}
	fulfillCount, err := v.consumer.SResponseCount(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx})
	if err != nil {
		return nil, err
	}
	avg, err := v.consumer.SAverageFulfillmentInMillions(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx})
	if err != nil {
		return nil, err
	}
	slow, err := v.consumer.SSlowestFulfillment(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx})
	if err != nil {
		return nil, err
	}
	fast, err := v.consumer.SFastestFulfillment(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx})
	if err != nil {
		return nil, err
	}
	return &VRFLoadTestMetrics{
		RequestCount:                 reqCount,
		FulfilmentCount:              fulfillCount,
		AverageFulfillmentInMillions: avg,
		SlowestFulfillment:           slow,
		FastestFulfillment:           fast,
	}, nil
}

// WaitRandomWordsFulfilled waits for fulfillment using watch then filter fallback.
func WaitRandomWordsFulfilled(
	coordinator *EthereumVRFCoordinatorV2,
	requestID *big.Int,
	requestBlock uint64,
	timeout time.Duration,
) (*CoordinatorRandomWordsFulfilled, error) {
	ev, err := coordinator.WaitForRandomWordsFulfilledEvent(RandomWordsFulfilledEventFilter{
		RequestIDs: []*big.Int{requestID},
		Timeout:    timeout,
	})
	if err == nil {
		return ev, nil
	}
	return coordinator.FilterRandomWordsFulfilledEvent(&bind.FilterOpts{Start: requestBlock}, requestID)
}

// --- VRF v2 wrapper (direct funding) ---

// EthereumVRFV2Wrapper wraps vrfv2_wrapper.VRFV2Wrapper.
type EthereumVRFV2Wrapper struct {
	client  *seth.Client
	wrapper *vrfv2_wrapper.VRFV2Wrapper
	address common.Address
}

func (v *EthereumVRFV2Wrapper) Address() string { return v.address.Hex() }

// DeployVRFV2Wrapper deploys VRFV2Wrapper.
func DeployVRFV2Wrapper(client *seth.Client, linkAddr, linkEthFeedAddr, coordinatorAddr string) (*EthereumVRFV2Wrapper, error) {
	abi, err := vrfv2_wrapper.VRFV2WrapperMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFV2Wrapper ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "VRFV2Wrapper", *abi,
		common.FromHex(vrfv2_wrapper.VRFV2WrapperMetaData.Bin),
		common.HexToAddress(linkAddr),
		common.HexToAddress(linkEthFeedAddr),
		common.HexToAddress(coordinatorAddr))
	if err != nil {
		return nil, fmt.Errorf("VRFV2Wrapper deployment failed: %w", err)
	}
	instance, err := vrfv2_wrapper.NewVRFV2Wrapper(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFV2Wrapper: %w", err)
	}
	return &EthereumVRFV2Wrapper{client: client, wrapper: instance, address: data.Address}, nil
}

func (v *EthereumVRFV2Wrapper) SetConfig(wrapperGasOverhead, coordinatorGasOverhead uint32, wrapperPremiumPercentage uint8, keyHash [32]byte, maxNumWords uint8) error {
	_, err := v.client.Decode(v.wrapper.SetConfig(
		v.client.NewTXOpts(),
		wrapperGasOverhead,
		coordinatorGasOverhead,
		wrapperPremiumPercentage,
		keyHash,
		maxNumWords,
	))
	return err
}

func (v *EthereumVRFV2Wrapper) GetSubID(ctx context.Context) (uint64, error) {
	return v.wrapper.SUBSCRIPTIONID(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx})
}

// EthereumVRFV2WrapperLoadTestConsumer wraps the wrapper load test consumer.
type EthereumVRFV2WrapperLoadTestConsumer struct {
	client   *seth.Client
	consumer *vrfv2_wrapper_load_test_consumer.VRFV2WrapperLoadTestConsumer
	address  common.Address
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) Address() string { return v.address.Hex() }

// DeployVRFV2WrapperLoadTestConsumer deploys wrapper load test consumer.
func DeployVRFV2WrapperLoadTestConsumer(client *seth.Client, linkAddr, wrapperAddr string) (*EthereumVRFV2WrapperLoadTestConsumer, error) {
	abi, err := vrfv2_wrapper_load_test_consumer.VRFV2WrapperLoadTestConsumerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFV2WrapperLoadTestConsumer ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "VRFV2WrapperLoadTestConsumer", *abi,
		common.FromHex(vrfv2_wrapper_load_test_consumer.VRFV2WrapperLoadTestConsumerMetaData.Bin),
		common.HexToAddress(linkAddr), common.HexToAddress(wrapperAddr))
	if err != nil {
		return nil, fmt.Errorf("VRFV2WrapperLoadTestConsumer deployment failed: %w", err)
	}
	instance, err := vrfv2_wrapper_load_test_consumer.NewVRFV2WrapperLoadTestConsumer(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFV2WrapperLoadTestConsumer: %w", err)
	}
	return &EthereumVRFV2WrapperLoadTestConsumer{client: client, consumer: instance, address: data.Address}, nil
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) RequestRandomness(
	coordinator CoordinatorV2Log,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	numWords uint32,
	requestCount uint16,
) (*CoordinatorRandomWordsRequested, error) {
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
	return parseRequestRandomnessLogs(coordinator, tx.Receipt.Logs)
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrfv2_wrapper_load_test_consumer.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{From: v.client.MustGetRootKeyAddress(), Context: ctx}, requestID)
}
