package contracts

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink/integration-tests/wrappers"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_mock_ethlink_aggregator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_owner"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_load_test_with_metrics"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_wrapper_load_test_consumer"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2_consumer_wrapper"
)

type EthereumVRFOwner struct {
	address  common.Address
	client   *seth.Client
	vrfOwner *vrf_owner.VRFOwner
}

type EthereumVRFCoordinatorV2 struct {
	address     common.Address
	client      *seth.Client
	coordinator *vrf_coordinator_v2.VRFCoordinatorV2
}

type EthereumBatchVRFCoordinatorV2 struct {
	address          common.Address
	client           *seth.Client
	batchCoordinator *batch_vrf_coordinator_v2.BatchVRFCoordinatorV2
}

// EthereumVRFConsumerV2 represents VRFv2 consumer contract
type EthereumVRFConsumerV2 struct {
	address  common.Address
	client   *seth.Client
	consumer *vrf_consumer_v2.VRFConsumerV2
}

// EthereumVRFv2Consumer represents VRFv2 consumer contract
type EthereumVRFv2Consumer struct {
	address  common.Address
	client   *seth.Client
	consumer *vrf_v2_consumer_wrapper.VRFv2Consumer
}

// EthereumVRFv2LoadTestConsumer represents VRFv2 consumer contract for performing Load Tests
type EthereumVRFv2LoadTestConsumer struct {
	address  common.Address
	client   *seth.Client
	consumer *vrf_load_test_with_metrics.VRFV2LoadTestWithMetrics
}

type EthereumVRFV2Wrapper struct {
	address common.Address
	client  *seth.Client
	wrapper *vrfv2_wrapper.VRFV2Wrapper
}

type EthereumVRFV2WrapperLoadTestConsumer struct {
	address  common.Address
	client   *seth.Client
	consumer *vrfv2_wrapper_load_test_consumer.VRFV2WrapperLoadTestConsumer
}

type GetRequestConfig struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	ProvingKeyHashes            [32]byte
}

type EthereumVRFMockETHLINKFeed struct {
	client  *seth.Client
	feed    *vrf_mock_ethlink_aggregator.VRFMockETHLINKAggregator
	address common.Address
}

func DeployVRFCoordinatorV2(seth *seth.Client, linkAddr, bhsAddr, linkEthFeedAddr string) (VRFCoordinatorV2, error) {
	abi, err := vrf_coordinator_v2.VRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return &EthereumVRFCoordinatorV2{}, fmt.Errorf("failed to get VRFCoordinatorV2 ABI: %w", err)
	}

	coordinatorDeploymentData, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFCoordinatorV2",
		*abi,
		common.FromHex(vrf_coordinator_v2.VRFCoordinatorV2MetaData.Bin),
		common.HexToAddress(linkAddr),
		common.HexToAddress(bhsAddr),
		common.HexToAddress(linkEthFeedAddr))
	if err != nil {
		return &EthereumVRFCoordinatorV2{}, fmt.Errorf("VRFCoordinatorV2 instance deployment have failed: %w", err)
	}

	coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(coordinatorDeploymentData.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFCoordinatorV2{}, fmt.Errorf("failed to instantiate VRFCoordinatorV2 instance: %w", err)
	}

	return &EthereumVRFCoordinatorV2{
		client:      seth,
		coordinator: coordinator,
		address:     coordinatorDeploymentData.Address,
	}, err
}

func LoadVRFCoordinatorV2(seth *seth.Client, address string) (*EthereumVRFCoordinatorV2, error) {
	abi, err := vrf_coordinator_v2.VRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return &EthereumVRFCoordinatorV2{}, fmt.Errorf("failed to get VRFCoordinatorV2 ABI: %w", err)
	}
	seth.ContractStore.AddABI("VRFCoordinatorV2", *abi)
	seth.ContractStore.AddBIN("VRFCoordinatorV2", common.FromHex(vrf_coordinator_v2.VRFCoordinatorV2MetaData.Bin))

	contract, err := vrf_coordinator_v2.NewVRFCoordinatorV2(common.HexToAddress(address), wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFCoordinatorV2{}, fmt.Errorf("failed to instantiate VRFCoordinatorV2 instance: %w", err)
	}

	return &EthereumVRFCoordinatorV2{
		client:      seth,
		address:     common.HexToAddress(address),
		coordinator: contract,
	}, nil
}

func DeployBatchVRFCoordinatorV2(seth *seth.Client, coordinatorAddress string) (BatchVRFCoordinatorV2, error) {
	abi, err := batch_vrf_coordinator_v2.BatchVRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return &EthereumBatchVRFCoordinatorV2{}, fmt.Errorf("failed to get BatchVRFCoordinatorV2 ABI: %w", err)
	}

	coordinatorDeploymentData, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFCoordinatorV2Plus",
		*abi,
		common.FromHex(batch_vrf_coordinator_v2.BatchVRFCoordinatorV2MetaData.Bin),
		common.HexToAddress(coordinatorAddress))
	if err != nil {
		return &EthereumBatchVRFCoordinatorV2{}, fmt.Errorf("BatchVRFCoordinatorV2 instance deployment have failed: %w", err)
	}

	contract, err := batch_vrf_coordinator_v2.NewBatchVRFCoordinatorV2(coordinatorDeploymentData.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumBatchVRFCoordinatorV2{}, fmt.Errorf("failed to instantiate BatchVRFCoordinatorV2 instance: %w", err)
	}

	return &EthereumBatchVRFCoordinatorV2{
		client:           seth,
		batchCoordinator: contract,
		address:          coordinatorDeploymentData.Address,
	}, err
}

func (v *EthereumBatchVRFCoordinatorV2) Address() string {
	return v.address.Hex()
}

func DeployVRFOwner(seth *seth.Client, coordinatorAddress string) (VRFOwner, error) {
	abi, err := vrf_owner.VRFOwnerMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFOwner{}, fmt.Errorf("failed to get VRFOwner ABI: %w", err)
	}

	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFOwner",
		*abi,
		common.FromHex(vrf_owner.VRFOwnerMetaData.Bin),
		common.HexToAddress(coordinatorAddress))
	if err != nil {
		return &EthereumVRFOwner{}, fmt.Errorf("VRFOwner instance deployment have failed: %w", err)
	}

	contract, err := vrf_owner.NewVRFOwner(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFOwner{}, fmt.Errorf("failed to instantiate VRFOwner instance: %w", err)
	}

	return &EthereumVRFOwner{
		client:   seth,
		vrfOwner: contract,
		address:  data.Address,
	}, err
}

// DeployVRFConsumerV2 deploys VRFv@ consumer contract
func DeployVRFConsumerV2(seth *seth.Client, linkAddr, coordinatorAddr common.Address) (VRFConsumerV2, error) {
	abi, err := vrf_consumer_v2.VRFConsumerV2MetaData.GetAbi()
	if err != nil {
		return &EthereumVRFConsumerV2{}, fmt.Errorf("failed to get VRFConsumerV2 ABI: %w", err)
	}

	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFConsumerV2",
		*abi,
		common.FromHex(vrf_consumer_v2.VRFConsumerV2MetaData.Bin),
		coordinatorAddr,
		linkAddr)
	if err != nil {
		return &EthereumVRFConsumerV2{}, fmt.Errorf("VRFConsumerV2 instance deployment have failed: %w", err)
	}

	contract, err := vrf_consumer_v2.NewVRFConsumerV2(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFConsumerV2{}, fmt.Errorf("failed to instantiate VRFConsumerV2 instance: %w", err)
	}

	return &EthereumVRFConsumerV2{
		client:   seth,
		consumer: contract,
		address:  data.Address,
	}, err
}

func DeployVRFv2Consumer(seth *seth.Client, coordinatorAddr common.Address) (VRFv2Consumer, error) {
	abi, err := vrf_v2_consumer_wrapper.VRFv2ConsumerMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFv2Consumer{}, fmt.Errorf("failed to get VRFv2Consumer ABI: %w", err)
	}

	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFv2Consumer",
		*abi,
		common.FromHex(vrf_v2_consumer_wrapper.VRFv2ConsumerMetaData.Bin),
		coordinatorAddr)
	if err != nil {
		return &EthereumVRFv2Consumer{}, fmt.Errorf("VRFv2Consumer instance deployment have failed: %w", err)
	}

	contract, err := vrf_v2_consumer_wrapper.NewVRFv2Consumer(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFv2Consumer{}, fmt.Errorf("failed to instantiate VRFv2Consumer instance: %w", err)
	}

	return &EthereumVRFv2Consumer{
		client:   seth,
		consumer: contract,
		address:  data.Address,
	}, err
}

func DeployVRFv2LoadTestConsumer(client *seth.Client, coordinatorAddr string) (VRFv2LoadTestConsumer, error) {
	abi, err := vrf_load_test_with_metrics.VRFV2LoadTestWithMetricsMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFv2LoadTestConsumer{}, fmt.Errorf("failed to get VRFV2LoadTestWithMetrics ABI: %w", err)
	}

	data, err := client.DeployContract(
		client.NewTXOpts(),
		"VRFV2LoadTestWithMetrics",
		*abi,
		common.FromHex(vrf_load_test_with_metrics.VRFV2LoadTestWithMetricsMetaData.Bin),
		common.HexToAddress(coordinatorAddr))
	if err != nil {
		return &EthereumVRFv2LoadTestConsumer{}, fmt.Errorf("VRFV2LoadTestWithMetrics instance deployment have failed: %w", err)
	}

	contract, err := vrf_load_test_with_metrics.NewVRFV2LoadTestWithMetrics(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumVRFv2LoadTestConsumer{}, fmt.Errorf("failed to instantiate VRFV2LoadTestWithMetrics instance: %w", err)
	}

	return &EthereumVRFv2LoadTestConsumer{
		client:   client,
		consumer: contract,
		address:  data.Address,
	}, err
}

func LoadVRFv2LoadTestConsumer(seth *seth.Client, addr common.Address) (VRFv2LoadTestConsumer, error) {
	abi, err := vrf_load_test_with_metrics.VRFV2LoadTestWithMetricsMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFv2LoadTestConsumer{}, fmt.Errorf("failed to get VRFV2LoadTestWithMetrics ABI: %w", err)
	}
	seth.ContractStore.AddABI("VRFV2LoadTestWithMetrics", *abi)
	seth.ContractStore.AddBIN("VRFV2LoadTestWithMetrics", common.FromHex(vrf_load_test_with_metrics.VRFV2LoadTestWithMetricsMetaData.Bin))

	contract, err := vrf_load_test_with_metrics.NewVRFV2LoadTestWithMetrics(addr, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFv2LoadTestConsumer{}, fmt.Errorf("failed to instantiate VRFV2LoadTestWithMetrics instance: %w", err)
	}

	return &EthereumVRFv2LoadTestConsumer{
		client:   seth,
		address:  addr,
		consumer: contract,
	}, nil
}

func DeployVRFV2Wrapper(client *seth.Client, linkAddr string, linkEthFeedAddr string, coordinatorAddr string) (VRFV2Wrapper, error) {
	abi, err := vrfv2_wrapper.VRFV2WrapperMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFV2Wrapper{}, fmt.Errorf("failed to get VRFV2Wrapper ABI: %w", err)
	}

	data, err := client.DeployContract(
		client.NewTXOpts(),
		"VRFV2Wrapper",
		*abi,
		common.FromHex(vrfv2_wrapper.VRFV2WrapperMetaData.Bin),
		common.HexToAddress(linkAddr),
		common.HexToAddress(linkEthFeedAddr),
		common.HexToAddress(coordinatorAddr))
	if err != nil {
		return &EthereumVRFV2Wrapper{}, fmt.Errorf("VRFV2Wrapper instance deployment have failed: %w", err)
	}

	contract, err := vrfv2_wrapper.NewVRFV2Wrapper(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumVRFV2Wrapper{}, fmt.Errorf("failed to instantiate VRFV2Wrapper instance: %w", err)
	}

	return &EthereumVRFV2Wrapper{
		client:  client,
		wrapper: contract,
		address: data.Address,
	}, err
}

func DeployVRFV2WrapperLoadTestConsumer(client *seth.Client, linkAddr string, vrfV2WrapperAddr string) (VRFv2WrapperLoadTestConsumer, error) {
	abi, err := vrfv2_wrapper_load_test_consumer.VRFV2WrapperLoadTestConsumerMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFV2WrapperLoadTestConsumer{}, fmt.Errorf("failed to get VRFV2WrapperLoadTestConsumer ABI: %w", err)
	}

	data, err := client.DeployContract(
		client.NewTXOpts(),
		"VRFV2WrapperLoadTestConsumer",
		*abi,
		common.FromHex(vrfv2_wrapper_load_test_consumer.VRFV2WrapperLoadTestConsumerMetaData.Bin),
		common.HexToAddress(linkAddr), common.HexToAddress(vrfV2WrapperAddr))
	if err != nil {
		return &EthereumVRFV2WrapperLoadTestConsumer{}, fmt.Errorf("VRFV2WrapperLoadTestConsumer instance deployment have failed: %w", err)
	}

	contract, err := vrfv2_wrapper_load_test_consumer.NewVRFV2WrapperLoadTestConsumer(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumVRFV2WrapperLoadTestConsumer{}, fmt.Errorf("failed to instantiate VRFV2WrapperLoadTestConsumer instance: %w", err)
	}

	return &EthereumVRFV2WrapperLoadTestConsumer{
		client:   client,
		consumer: contract,
		address:  data.Address,
	}, err
}

func (v *EthereumVRFCoordinatorV2) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFCoordinatorV2) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	hash, err := v.coordinator.HashOfKey(opts, pubKey)
	if err != nil {
		return [32]byte{}, err
	}
	return hash, nil
}

func (v *EthereumVRFCoordinatorV2) GetSubscription(ctx context.Context, subID uint64) (Subscription, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	subscription, err := v.coordinator.GetSubscription(opts, subID)
	if err != nil {
		return Subscription{}, err
	}
	return Subscription{
		Balance:       subscription.Balance,
		NativeBalance: nil,
		SubOwner:      subscription.Owner,
		Consumers:     subscription.Consumers,
		ReqCount:      subscription.ReqCount,
	}, nil
}

func (v *EthereumVRFCoordinatorV2) GetOwner(ctx context.Context) (common.Address, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
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
		From:    v.client.MustGetRootKeyAddress(),
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
		From:    v.client.MustGetRootKeyAddress(),
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
		From:    v.client.MustGetRootKeyAddress(),
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
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	config, err := v.coordinator.GetFeeConfig(opts)
	if err != nil {
		return vrf_coordinator_v2.GetFeeConfig{}, err
	}
	return config, nil
}

func (v *EthereumVRFCoordinatorV2) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig vrf_coordinator_v2.VRFCoordinatorV2FeeConfig) error {
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

func (v *EthereumVRFCoordinatorV2) RegisterProvingKey(
	oracleAddr string,
	publicProvingKey [2]*big.Int,
) error {
	_, err := v.client.Decode(v.coordinator.RegisterProvingKey(v.client.NewTXOpts(),
		common.HexToAddress(oracleAddr), publicProvingKey))
	return err
}

func (v *EthereumVRFCoordinatorV2) TransferOwnership(to common.Address) error {
	_, err := v.client.Decode(v.coordinator.TransferOwnership(v.client.NewTXOpts(), to))
	return err
}

func (v *EthereumVRFCoordinatorV2) CreateSubscription() (*types.Receipt, error) {
	tx, err := v.client.Decode(v.coordinator.CreateSubscription(v.client.NewTXOpts()))
	return tx.Receipt, err
}

func (v *EthereumVRFCoordinatorV2) AddConsumer(subId uint64, consumerAddress string) error {
	_, err := v.client.Decode(v.coordinator.AddConsumer(
		v.client.NewTXOpts(),
		subId,
		common.HexToAddress(consumerAddress),
	))
	return err
}

func (v *EthereumVRFCoordinatorV2) PendingRequestsExist(ctx context.Context, subID uint64) (bool, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	pendingRequestExists, err := v.coordinator.PendingRequestExists(opts, subID)
	if err != nil {
		return false, err
	}
	return pendingRequestExists, nil
}

func (v *EthereumVRFCoordinatorV2) OracleWithdraw(recipient common.Address, amount *big.Int) error {
	_, err := v.client.Decode(v.coordinator.OracleWithdraw(v.client.NewTXOpts(), recipient, amount))
	return err
}

// OwnerCancelSubscription cancels subscription,
// return funds to the subscription owner,
// down not check if pending requests for a sub exist,
// outstanding requests may fail onchain
func (v *EthereumVRFCoordinatorV2) OwnerCancelSubscription(subID uint64) (*seth.DecodedTransaction, *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled, error) {
	tx, err := v.client.Decode(v.coordinator.OwnerCancelSubscription(
		v.client.NewTXOpts(),
		subID,
	))
	if err != nil {
		return nil, nil, err
	}
	var subCanceledEvent *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled
	for _, log := range tx.Receipt.Logs {
		for _, topic := range log.Topics {
			if topic.Cmp(vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled{}.Topic()) == 0 {
				subCanceledEvent, err = v.coordinator.ParseSubscriptionCanceled(*log)
				if err != nil {
					return nil, nil, fmt.Errorf("parsing SubscriptionCanceled log failed, err: %w", err)
				}
			}
		}
	}
	return tx, subCanceledEvent, err
}

func (v *EthereumVRFCoordinatorV2) ParseSubscriptionCanceled(log types.Log) (*vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled, error) {
	return v.coordinator.ParseSubscriptionCanceled(log)
}

func (v *EthereumVRFCoordinatorV2) ParseRandomWordsRequested(log types.Log) (*CoordinatorRandomWordsRequested, error) {
	requested, err := v.coordinator.ParseRandomWordsRequested(log)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RandomWordsRequested event: %w", err)
	}

	return &CoordinatorRandomWordsRequested{
		KeyHash:                     requested.KeyHash,
		RequestId:                   requested.RequestId,
		PreSeed:                     requested.PreSeed,
		SubId:                       strconv.FormatUint(requested.SubId, 10),
		MinimumRequestConfirmations: requested.MinimumRequestConfirmations,
		CallbackGasLimit:            requested.CallbackGasLimit,
		NumWords:                    requested.NumWords,
		Sender:                      requested.Sender,
		ExtraArgs:                   nil,
		Raw:                         requested.Raw,
	}, nil
}

func (v *EthereumVRFCoordinatorV2) ParseRandomWordsFulfilled(log types.Log) (*CoordinatorRandomWordsFulfilled, error) {
	fulfilled, err := v.coordinator.ParseRandomWordsFulfilled(log)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RandomWordsFulfilled event: %w", err)
	}

	return &CoordinatorRandomWordsFulfilled{
		RequestId:  fulfilled.RequestId,
		OutputSeed: fulfilled.OutputSeed,
		Payment:    fulfilled.Payment,
		Success:    fulfilled.Success,
		Raw:        fulfilled.Raw,
	}, nil
}

func (v *EthereumVRFCoordinatorV2) ParseLog(log types.Log) (generated.AbigenLog, error) {
	return v.coordinator.ParseLog(log)
}

// CancelSubscription cancels subscription by Sub owner,
// return funds to specified address,
// checks if pending requests for a sub exist
func (v *EthereumVRFCoordinatorV2) CancelSubscription(subID uint64, to common.Address) (*seth.DecodedTransaction, *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled, error) {
	tx, err := v.client.Decode(v.coordinator.CancelSubscription(
		v.client.NewTXOpts(),
		subID,
		to,
	))
	if err != nil {
		return nil, nil, err
	}
	var subCanceledEvent *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled
	for _, log := range tx.Receipt.Logs {
		for _, topic := range log.Topics {
			if topic.Cmp(vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled{}.Topic()) == 0 {
				subCanceledEvent, err = v.coordinator.ParseSubscriptionCanceled(*log)
				if err != nil {
					return nil, nil, fmt.Errorf("parsing SubscriptionCanceled log failed, err: %w", err)
				}
			}
		}
	}
	return tx, subCanceledEvent, err
}

func (v *EthereumVRFCoordinatorV2) FindSubscriptionID(subID uint64) (uint64, error) {
	owner := v.client.MustGetRootKeyAddress()
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

func (v *EthereumVRFCoordinatorV2) WaitForRandomWordsFulfilledEvent(filter RandomWordsFulfilledEventFilter) (*CoordinatorRandomWordsFulfilled, error) {
	randomWordsFulfilledEventsChannel := make(chan *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled)
	subscription, err := v.coordinator.WatchRandomWordsFulfilled(nil, randomWordsFulfilledEventsChannel, filter.RequestIds)
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
				RequestId:  randomWordsFulfilledEvent.RequestId,
				OutputSeed: randomWordsFulfilledEvent.OutputSeed,
				Payment:    randomWordsFulfilledEvent.Payment,
				Success:    randomWordsFulfilledEvent.Success,
				Raw:        randomWordsFulfilledEvent.Raw,
			}, nil
		}
	}
}

func (v *EthereumVRFCoordinatorV2) WaitForConfigSetEvent(timeout time.Duration) (*CoordinatorConfigSet, error) {
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
			return &CoordinatorConfigSet{
				MinimumRequestConfirmations: event.MinimumRequestConfirmations,
				MaxGasLimit:                 event.MaxGasLimit,
				StalenessSeconds:            event.StalenessSeconds,
				GasAfterPaymentCalculation:  event.GasAfterPaymentCalculation,
				FallbackWeiPerUnitLink:      event.FallbackWeiPerUnitLink,
				FeeConfig: VRFCoordinatorV2FeeConfig{
					FulfillmentFlatFeeLinkPPMTier1: event.FeeConfig.FulfillmentFlatFeeLinkPPMTier1,
					FulfillmentFlatFeeLinkPPMTier2: event.FeeConfig.FulfillmentFlatFeeLinkPPMTier2,
					FulfillmentFlatFeeLinkPPMTier3: event.FeeConfig.FulfillmentFlatFeeLinkPPMTier3,
					FulfillmentFlatFeeLinkPPMTier4: event.FeeConfig.FulfillmentFlatFeeLinkPPMTier4,
					FulfillmentFlatFeeLinkPPMTier5: event.FeeConfig.FulfillmentFlatFeeLinkPPMTier5,
					ReqsForTier2:                   event.FeeConfig.ReqsForTier2,
					ReqsForTier3:                   event.FeeConfig.ReqsForTier3,
					ReqsForTier4:                   event.FeeConfig.ReqsForTier4,
					ReqsForTier5:                   event.FeeConfig.ReqsForTier5,
				},
			}, nil
		}
	}
}

// GetAllRandomWords get all VRFv2 randomness output words
func (v *EthereumVRFConsumerV2) GetAllRandomWords(ctx context.Context, num int) ([]*big.Int, error) {
	words := make([]*big.Int, 0)
	for i := 0; i < num; i++ {
		word, err := v.consumer.SRandomWords(&bind.CallOpts{
			From:    v.client.MustGetRootKeyAddress(),
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
func (v *EthereumVRFConsumerV2) LoadExistingConsumer(seth *seth.Client, address common.Address) error {
	abi, err := vrf_consumer_v2.VRFConsumerV2MetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("failed to get VRFConsumerV2 ABI: %w", err)
	}
	seth.ContractStore.AddABI("VRFConsumerV2", *abi)
	seth.ContractStore.AddBIN("VRFConsumerV2", common.FromHex(vrf_consumer_v2.VRFConsumerV2MetaData.Bin))

	contract, err := vrf_consumer_v2.NewVRFConsumerV2(address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return fmt.Errorf("failed to instantiate VRFConsumerV2 instance: %w", err)
	}

	v.client = seth
	v.consumer = contract
	v.address = address

	return nil
}

// CreateFundedSubscription create funded subscription for VRFv2 randomness
func (v *EthereumVRFConsumerV2) CreateFundedSubscription(funds *big.Int) error {
	_, err := v.client.Decode(v.consumer.CreateSubscriptionAndFund(v.client.NewTXOpts(), funds))
	return err
}

// TopUpSubscriptionFunds add funds to a VRFv2 subscription
func (v *EthereumVRFConsumerV2) TopUpSubscriptionFunds(funds *big.Int) error {
	_, err := v.client.Decode(v.consumer.TopUpSubscription(v.client.NewTXOpts(), funds))
	return err
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
		From:    v.client.MustGetRootKeyAddress(),
		Context: context.Background(),
	})
}

// GasAvailable get available gas after randomness fulfilled
func (v *EthereumVRFConsumerV2) GasAvailable() (*big.Int, error) {
	return v.consumer.SGasAvailable(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: context.Background(),
	})
}

func (v *EthereumVRFConsumerV2) Fund(_ *big.Float) error {
	panic("do not use this function, use actions_seth.SendFunds() instead, otherwise we will have to deal with circular dependencies")
}

// RequestRandomness request VRFv2 random words
func (v *EthereumVRFConsumerV2) RequestRandomness(hash [32]byte, subID uint64, confs uint16, gasLimit uint32, numWords uint32) error {
	_, err := v.client.Decode(v.consumer.RequestRandomness(v.client.NewTXOpts(), hash, subID, confs, gasLimit, numWords))
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
	return nil
}

// RequestRandomness request VRFv2 random words
func (v *EthereumVRFv2Consumer) RequestRandomness(hash [32]byte, subID uint64, confs uint16, gasLimit uint32, numWords uint32) error {
	_, err := v.client.Decode(v.consumer.RequestRandomWords(v.client.NewTXOpts(), subID, gasLimit, confs, numWords, hash))
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
	return nil
}

// RandomnessOutput get VRFv2 randomness output (word)
func (v *EthereumVRFConsumerV2) RandomnessOutput(ctx context.Context, arg0 *big.Int) (*big.Int, error) {
	return v.consumer.SRandomWords(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, arg0)
}

func (v *EthereumVRFv2LoadTestConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFv2LoadTestConsumer) RequestRandomness(
	coordinator Coordinator,
	keyHash [32]byte,
	subID uint64,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	numWords uint32,
	requestCount uint16,
) (*CoordinatorRandomWordsRequested, error) {
	return v.RequestRandomnessFromKey(coordinator, keyHash, subID, requestConfirmations, callbackGasLimit, numWords, requestCount, 0)
}

func (v *EthereumVRFv2LoadTestConsumer) RequestRandomnessFromKey(
	coordinator Coordinator,
	keyHash [32]byte,
	subID uint64,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	numWords uint32,
	requestCount uint16,
	keyNum int,
) (*CoordinatorRandomWordsRequested, error) {
	tx, err := v.client.Decode(v.consumer.RequestRandomWords(v.client.NewTXKeyOpts(keyNum), subID, requestConfirmations, keyHash, callbackGasLimit, numWords, requestCount))
	if err != nil {
		return nil, fmt.Errorf("RequestRandomWords failed, err: %w", err)
	}
	randomWordsRequestedEvent, err := parseRequestRandomnessLogs(coordinator, tx.Receipt.Logs)
	if err != nil {
		return nil, err
	}
	return randomWordsRequestedEvent, nil
}

func (v *EthereumVRFv2LoadTestConsumer) RequestRandomWordsWithForceFulfill(
	coordinator Coordinator,
	keyHash [32]byte,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	numWords uint32,
	requestCount uint16,
	subTopUpAmount *big.Int,
	linkAddress common.Address,
) (*CoordinatorRandomWordsRequested, error) {
	tx, err := v.client.Decode(v.consumer.RequestRandomWordsWithForceFulfill(
		v.client.NewTXOpts(),
		requestConfirmations,
		keyHash,
		callbackGasLimit,
		numWords,
		requestCount,
		subTopUpAmount,
		linkAddress,
	))
	if err != nil {
		return nil, fmt.Errorf("RequestRandomWords failed, err: %w", err)
	}
	randomWordsRequestedEvent, err := parseRequestRandomnessLogs(coordinator, tx.Receipt.Logs)
	if err != nil {
		return nil, err
	}
	return randomWordsRequestedEvent, nil
}

func (v *EthereumVRFv2Consumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_v2_consumer_wrapper.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, requestID)
}

func (v *EthereumVRFv2Consumer) GetLastRequestId(ctx context.Context) (*big.Int, error) {
	return v.consumer.LastRequestId(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFv2LoadTestConsumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_load_test_with_metrics.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, requestID)
}

func (v *EthereumVRFv2LoadTestConsumer) GetLastRequestId(ctx context.Context) (*big.Int, error) {
	return v.consumer.SLastRequestId(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFv2LoadTestConsumer) ResetMetrics() error {
	_, err := v.client.Decode(v.consumer.Reset(v.client.NewTXOpts()))
	return err
}

func (v *EthereumVRFv2LoadTestConsumer) GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error) {
	requestCount, err := v.consumer.SRequestCount(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}
	fulfilmentCount, err := v.consumer.SResponseCount(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})

	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}
	averageFulfillmentInMillions, err := v.consumer.SAverageFulfillmentInMillions(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}
	slowestFulfillment, err := v.consumer.SSlowestFulfillment(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})

	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}
	fastestFulfillment, err := v.consumer.SFastestFulfillment(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
	if err != nil {
		return &VRFLoadTestMetrics{}, err
	}

	return &VRFLoadTestMetrics{
		RequestCount:                         requestCount,
		FulfilmentCount:                      fulfilmentCount,
		AverageFulfillmentInMillions:         averageFulfillmentInMillions,
		SlowestFulfillment:                   slowestFulfillment,
		FastestFulfillment:                   fastestFulfillment,
		P90FulfillmentBlockTime:              0.0,
		P95FulfillmentBlockTime:              0.0,
		AverageResponseTimeInSecondsMillions: nil,
		SlowestResponseTimeInSeconds:         nil,
		FastestResponseTimeInSeconds:         nil,
	}, nil
}

func (v *EthereumVRFV2Wrapper) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFV2Wrapper) SetConfig(wrapperGasOverhead uint32, coordinatorGasOverhead uint32, wrapperPremiumPercentage uint8, keyHash [32]byte, maxNumWords uint8) error {
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
	return v.wrapper.SUBSCRIPTIONID(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) Fund(_ *big.Float) error {
	panic("do not use this function, use actions_seth.SendFunds() instead, otherwise we will have to deal with circular dependencies")
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) RequestRandomness(coordinator Coordinator, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, requestCount uint16) (*CoordinatorRandomWordsRequested, error) {
	tx, err := v.client.Decode(v.consumer.MakeRequests(v.client.NewTXOpts(),
		callbackGasLimit, requestConfirmations, numWords, requestCount))
	if err != nil {
		return nil, err
	}
	randomWordsRequestedEvent, err := parseRequestRandomnessLogs(coordinator, tx.Receipt.Logs)
	if err != nil {
		return nil, err
	}
	return randomWordsRequestedEvent, nil
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrfv2_wrapper_load_test_consumer.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, requestID)
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) GetLastRequestId(ctx context.Context) (*big.Int, error) {
	return v.consumer.SLastRequestId(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) GetWrapper(ctx context.Context) (common.Address, error) {
	return v.consumer.IVrfV2Wrapper(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFV2WrapperLoadTestConsumer) GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error) {
	requestCount, err := v.consumer.SRequestCount(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	fulfilmentCount, err := v.consumer.SResponseCount(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})

	if err != nil {
		return nil, err
	}
	averageFulfillmentInMillions, err := v.consumer.SAverageFulfillmentInMillions(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	slowestFulfillment, err := v.consumer.SSlowestFulfillment(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})

	if err != nil {
		return nil, err
	}
	fastestFulfillment, err := v.consumer.SFastestFulfillment(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
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
		P90FulfillmentBlockTime:              0.0,
		P95FulfillmentBlockTime:              0.0,
		AverageResponseTimeInSecondsMillions: nil,
		SlowestResponseTimeInSeconds:         nil,
		FastestResponseTimeInSeconds:         nil,
	}, nil
}

func (v *EthereumVRFOwner) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFOwner) SetAuthorizedSenders(senders []common.Address) error {
	_, err := v.client.Decode(v.vrfOwner.SetAuthorizedSenders(
		v.client.NewTXOpts(),
		senders,
	))
	return err
}

func (v *EthereumVRFOwner) AcceptVRFOwnership() error {
	_, err := v.client.Decode(v.vrfOwner.AcceptVRFOwnership(v.client.NewTXOpts()))
	return err
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

func (v *EthereumVRFOwner) OwnerCancelSubscription(subID uint64) (*types.Transaction, error) {
	// Do not wrap in Decode() to avoid waiting until the transaction is mined
	return v.vrfOwner.OwnerCancelSubscription(
		v.client.NewTXOpts(),
		subID,
	)
}

func (v *EthereumVRFMockETHLINKFeed) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFMockETHLINKFeed) LatestRoundData() (*big.Int, error) {
	data, err := v.feed.LatestRoundData(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.Ans, nil
}

func (v *EthereumVRFMockETHLINKFeed) LatestRoundDataUpdatedAt() (*big.Int, error) {
	data, err := v.feed.LatestRoundData(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.UpdatedAt, nil
}

func (v *EthereumVRFMockETHLINKFeed) SetBlockTimestampDeduction(blockTimestampDeduction *big.Int) error {
	_, err := v.client.Decode(v.feed.SetBlockTimestampDeduction(v.client.NewTXOpts(), blockTimestampDeduction))
	return err
}
