package contracts

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/montanaflynn/stats"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_blockhash_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_vrf_coordinator_v2plus"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_test_v2_5"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2_5_arbitrum"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2_5_optimism"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_v2plus_load_test_with_metrics"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_v2plus_upgraded_version"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2plus_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2plus_wrapper_arbitrum"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2plus_wrapper_load_test_consumer"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2plus_wrapper_optimism"
	"github.com/smartcontractkit/chainlink/integration-tests/wrappers"
)

type EthereumVRFCoordinatorV2_5 struct {
	address     common.Address
	client      *seth.Client
	coordinator vrf_coordinator_v2_5.VRFCoordinatorV25Interface
}

type EthereumVRFCoordinatorV2_5_Optimism struct {
	Address     common.Address
	client      *seth.Client
	coordinator vrf_coordinator_v2_5_optimism.VRFCoordinatorV25Optimism
}

type EthereumVRFCoordinatorV2_5_Arbitrum struct {
	Address     common.Address
	client      *seth.Client
	coordinator vrf_coordinator_v2_5_arbitrum.VRFCoordinatorV25Arbitrum
}

type EthereumVRFCoordinatorTestV2_5 struct {
	Address     common.Address
	client      *seth.Client
	coordinator vrf_coordinator_test_v2_5.VRFCoordinatorTestV25
}

type EthereumBatchVRFCoordinatorV2Plus struct {
	address          common.Address
	client           *seth.Client
	batchCoordinator *batch_vrf_coordinator_v2plus.BatchVRFCoordinatorV2Plus
}

type EthereumVRFCoordinatorV2PlusUpgradedVersion struct {
	address     common.Address
	client      *seth.Client
	coordinator *vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersion
}

// EthereumVRFv2PlusLoadTestConsumer represents VRFv2Plus consumer contract for performing Load Tests
type EthereumVRFv2PlusLoadTestConsumer struct {
	address  common.Address
	client   *seth.Client
	consumer *vrf_v2plus_load_test_with_metrics.VRFV2PlusLoadTestWithMetrics
}

type EthereumVRFV2PlusWrapperLoadTestConsumer struct {
	address  common.Address
	client   *seth.Client
	consumer *vrfv2plus_wrapper_load_test_consumer.VRFV2PlusWrapperLoadTestConsumer
}

type EthereumVRFV2PlusWrapper struct {
	address common.Address
	client  *seth.Client
	wrapper *vrfv2plus_wrapper.VRFV2PlusWrapper
}

type EthereumVRFV2PlusWrapperOptimism struct {
	Address common.Address
	client  *seth.Client
	wrapper *vrfv2plus_wrapper_optimism.VRFV2PlusWrapperOptimism
}

type EthereumVRFV2PlusWrapperArbitrum struct {
	Address common.Address
	client  *seth.Client
	wrapper *vrfv2plus_wrapper_arbitrum.VRFV2PlusWrapperArbitrum
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
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	return v.wrapper.SVrfCoordinator(opts)
}

// DeployVRFCoordinatorV2_5 deploys VRFV2_5 coordinator contract
func DeployVRFCoordinatorV2_5(seth *seth.Client, bhsAddr string) (VRFCoordinatorV2_5, error) {
	abi, err := vrf_coordinator_v2_5.VRFCoordinatorV25MetaData.GetAbi()
	if err != nil {
		return &EthereumVRFCoordinatorV2_5{}, fmt.Errorf("failed to get VRFCoordinatorV2_5 ABI: %w", err)
	}

	coordinatorDeploymentData, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFCoordinatorV2_5",
		*abi,
		common.FromHex(vrf_coordinator_v2_5.VRFCoordinatorV25MetaData.Bin),
		common.HexToAddress(bhsAddr))
	if err != nil {
		return &EthereumVRFCoordinatorV2_5{}, fmt.Errorf("VRFCoordinatorV2_5 instance deployment have failed: %w", err)
	}

	contract, err := vrf_coordinator_v2_5.NewVRFCoordinatorV25(coordinatorDeploymentData.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFCoordinatorV2_5{}, fmt.Errorf("failed to instantiate VRFCoordinatorV2_5 instance: %w", err)
	}

	return &EthereumVRFCoordinatorV2_5{
		client:      seth,
		coordinator: contract,
		address:     coordinatorDeploymentData.Address,
	}, err
}

func DeployVRFCoordinatorV2_5_Optimism(seth *seth.Client, bhsAddr string) (*EthereumVRFCoordinatorV2_5_Optimism, error) {
	abi, err := vrf_coordinator_v2_5_optimism.VRFCoordinatorV25OptimismMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFCoordinatorV2_5_Optimism ABI: %w", err)
	}
	coordinatorDeploymentData, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFCoordinatorV2_5_Optimism",
		*abi,
		common.FromHex(vrf_coordinator_v2_5_optimism.VRFCoordinatorV25OptimismMetaData.Bin),
		common.HexToAddress(bhsAddr))
	if err != nil {
		return nil, fmt.Errorf("VRFCoordinatorV2_5_Optimism instance deployment have failed: %w", err)
	}
	contract, err := vrf_coordinator_v2_5_optimism.NewVRFCoordinatorV25Optimism(coordinatorDeploymentData.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFCoordinatorV2_5_Optimism instance: %w", err)
	}
	return &EthereumVRFCoordinatorV2_5_Optimism{
		client:      seth,
		coordinator: *contract,
		Address:     coordinatorDeploymentData.Address,
	}, err
}

func DeployVRFCoordinatorV2_5_Arbitrum(seth *seth.Client, bhsAddr string) (*EthereumVRFCoordinatorV2_5_Arbitrum, error) {
	abi, err := vrf_coordinator_v2_5_arbitrum.VRFCoordinatorV25ArbitrumMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFCoordinatorV2_5_Arbitrum ABI: %w", err)
	}
	coordinatorDeploymentData, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFCoordinatorV2_5_Arbitrum",
		*abi,
		common.FromHex(vrf_coordinator_v2_5_arbitrum.VRFCoordinatorV25ArbitrumMetaData.Bin),
		common.HexToAddress(bhsAddr))
	if err != nil {
		return nil, fmt.Errorf("VRFCoordinatorV2_5_Arbitrum instance deployment have failed: %w", err)
	}
	contract, err := vrf_coordinator_v2_5_arbitrum.NewVRFCoordinatorV25Arbitrum(coordinatorDeploymentData.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFCoordinatorV2_5_Arbitrum instance: %w", err)
	}
	return &EthereumVRFCoordinatorV2_5_Arbitrum{
		client:      seth,
		coordinator: *contract,
		Address:     coordinatorDeploymentData.Address,
	}, err
}

func DeployVRFCoordinatorTestV2_5(seth *seth.Client, bhsAddr string) (*EthereumVRFCoordinatorTestV2_5, error) {
	abi, err := vrf_coordinator_test_v2_5.VRFCoordinatorTestV25MetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFCoordinatorTestV2_5 ABI: %w", err)
	}
	coordinatorDeploymentData, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFCoordinatorTestV2_5",
		*abi,
		common.FromHex(vrf_coordinator_test_v2_5.VRFCoordinatorTestV25MetaData.Bin),
		common.HexToAddress(bhsAddr))
	if err != nil {
		return nil, fmt.Errorf("VRFCoordinatorTestV2_5 instance deployment have failed: %w", err)
	}
	contract, err := vrf_coordinator_test_v2_5.NewVRFCoordinatorTestV25(coordinatorDeploymentData.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFCoordinatorTestV2_5 instance: %w", err)
	}
	return &EthereumVRFCoordinatorTestV2_5{
		client:      seth,
		coordinator: *contract,
		Address:     coordinatorDeploymentData.Address,
	}, err
}

func DeployBatchVRFCoordinatorV2Plus(seth *seth.Client, coordinatorAddress string) (BatchVRFCoordinatorV2Plus, error) {
	abi, err := batch_vrf_coordinator_v2plus.BatchVRFCoordinatorV2PlusMetaData.GetAbi()
	if err != nil {
		return &EthereumBatchVRFCoordinatorV2Plus{}, fmt.Errorf("failed to get BatchVRFCoordinatorV2Plus ABI: %w", err)
	}

	coordinatorDeploymentData, err := seth.DeployContract(
		seth.NewTXOpts(),
		"BatchVRFCoordinatorV2Plus",
		*abi,
		common.FromHex(batch_vrf_coordinator_v2plus.BatchVRFCoordinatorV2PlusMetaData.Bin),
		common.HexToAddress(coordinatorAddress))
	if err != nil {
		return &EthereumBatchVRFCoordinatorV2Plus{}, fmt.Errorf("BatchVRFCoordinatorV2Plus instance deployment have failed: %w", err)
	}

	contract, err := batch_vrf_coordinator_v2plus.NewBatchVRFCoordinatorV2Plus(coordinatorDeploymentData.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumBatchVRFCoordinatorV2Plus{}, fmt.Errorf("failed to instantiate BatchVRFCoordinatorV2Plus instance: %w", err)
	}

	return &EthereumBatchVRFCoordinatorV2Plus{
		client:           seth,
		batchCoordinator: contract,
		address:          coordinatorDeploymentData.Address,
	}, err
}

func (v *EthereumVRFCoordinatorV2_5) Address() string {
	return v.address.Hex()
}

func (v *EthereumBatchVRFCoordinatorV2Plus) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFCoordinatorV2_5) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
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

func (v *EthereumVRFCoordinatorV2_5) GetActiveSubscriptionIds(ctx context.Context, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
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
		From:    v.client.MustGetRootKeyAddress(),
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

func (v *EthereumVRFCoordinatorV2_5) ParseRandomWordsFulfilled(log types.Log) (*CoordinatorRandomWordsFulfilled, error) {
	fulfilled, err := v.coordinator.ParseRandomWordsFulfilled(log)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RandomWordsFulfilled event: %w", err)
	}
	return &CoordinatorRandomWordsFulfilled{
		RequestId:     fulfilled.RequestId,
		OutputSeed:    fulfilled.OutputSeed,
		Payment:       fulfilled.Payment,
		SubId:         fulfilled.SubId.String(),
		NativePayment: fulfilled.NativePayment,
		OnlyPremium:   fulfilled.OnlyPremium,
		Success:       fulfilled.Success,
		Raw:           fulfilled.Raw,
	}, nil
}

func (v *EthereumVRFCoordinatorV2_5) GetSubscription(ctx context.Context, subID *big.Int) (Subscription, error) {
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
		NativeBalance: subscription.NativeBalance,
		SubOwner:      subscription.SubOwner,
		Consumers:     subscription.Consumers,
		ReqCount:      subscription.ReqCount,
	}, nil
}

func (v *EthereumVRFCoordinatorV2_5) GetLinkTotalBalance(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
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
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	totalBalance, err := v.coordinator.STotalNativeBalance(opts)
	if err != nil {
		return nil, err
	}
	return totalBalance, nil
}

func (v *EthereumVRFCoordinatorV2_5) GetBlockHashStoreAddress(ctx context.Context) (common.Address, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	blockHashStoreAddress, err := v.coordinator.BLOCKHASHSTORE(opts)
	if err != nil {
		return common.Address{}, err
	}
	return blockHashStoreAddress, nil
}

func (v *EthereumVRFCoordinatorV2_5) GetLinkAddress(ctx context.Context) (common.Address, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	address, err := v.coordinator.LINK(opts)
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}

func (v *EthereumVRFCoordinatorV2_5) GetLinkNativeFeed(ctx context.Context) (common.Address, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	address, err := v.coordinator.LINKNATIVEFEED(opts)
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}

func (v *EthereumVRFCoordinatorV2_5) GetConfig(ctx context.Context) (vrf_coordinator_v2_5.SConfig, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	config, err := v.coordinator.SConfig(opts)
	if err != nil {
		return vrf_coordinator_v2_5.SConfig{}, err
	}
	return config, nil
}

// OwnerCancelSubscription cancels subscription by Coordinator owner
// return funds to sub owner,
// does not check if pending requests for a sub exist
func (v *EthereumVRFCoordinatorV2_5) OwnerCancelSubscription(subID *big.Int) (*seth.DecodedTransaction, *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled, error) {
	tx, err := v.client.Decode(v.coordinator.OwnerCancelSubscription(
		v.client.NewTXOpts(),
		subID,
	))
	if err != nil {
		return nil, nil, err
	}
	var cancelEvent *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled
	for _, log := range tx.Receipt.Logs {
		for _, topic := range log.Topics {
			if topic.Cmp(vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled{}.Topic()) == 0 {
				cancelEvent, err = v.coordinator.ParseSubscriptionCanceled(*log)
				if err != nil {
					return nil, nil, fmt.Errorf("parsing SubscriptionCanceled log failed, err: %w", err)
				}
			}
		}
	}
	return tx, cancelEvent, err
}

// CancelSubscription cancels subscription by Sub owner,
// return funds to specified address,
// checks if pending requests for a sub exist
func (v *EthereumVRFCoordinatorV2_5) CancelSubscription(subID *big.Int, to common.Address) (*seth.DecodedTransaction, *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled, error) {
	tx, err := v.client.Decode(v.coordinator.CancelSubscription(
		v.client.NewTXOpts(),
		subID,
		to,
	))
	if err != nil {
		return nil, nil, err
	}
	var cancelEvent *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled
	for _, log := range tx.Receipt.Logs {
		for _, topic := range log.Topics {
			if topic.Cmp(vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled{}.Topic()) == 0 {
				cancelEvent, err = v.coordinator.ParseSubscriptionCanceled(*log)
				if err != nil {
					return nil, nil, fmt.Errorf("parsing SubscriptionCanceled log failed, err: %w", err)
				}
			}
		}
	}
	return tx, cancelEvent, err
}

func (v *EthereumVRFCoordinatorV2_5) Withdraw(recipient common.Address) error {
	_, err := v.client.Decode(v.coordinator.Withdraw(
		v.client.NewTXOpts(),
		recipient,
	))
	return err
}

func (v *EthereumVRFCoordinatorV2_5) WithdrawNative(recipient common.Address) error {
	_, err := v.client.Decode(v.coordinator.WithdrawNative(
		v.client.NewTXOpts(),
		recipient,
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
	linkPremiumPercentage uint8) error {
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

func (v *EthereumVRFCoordinatorV2_5) SetLINKAndLINKNativeFeed(linkAddress string, linkNativeFeedAddress string) error {
	_, err := v.client.Decode(v.coordinator.SetLINKAndLINKNativeFeed(
		v.client.NewTXOpts(),
		common.HexToAddress(linkAddress),
		common.HexToAddress(linkNativeFeedAddress),
	))
	return err
}

func (v *EthereumVRFCoordinatorV2_5) RegisterProvingKey(
	publicProvingKey [2]*big.Int,
	gasLaneMaxGas uint64,
) error {
	_, err := v.client.Decode(v.coordinator.RegisterProvingKey(v.client.NewTXOpts(), publicProvingKey, gasLaneMaxGas))
	return err
}

func (v *EthereumVRFCoordinatorV2_5) CreateSubscription() (*types.Transaction, error) {
	tx, err := v.client.Decode(v.coordinator.CreateSubscription(v.client.NewTXOpts()))
	if err != nil {
		return nil, err
	}
	return tx.Transaction, nil
}

func (v *EthereumVRFCoordinatorV2_5) Migrate(subId *big.Int, coordinatorAddress string) (*seth.DecodedTransaction, *vrf_coordinator_v2_5.VRFCoordinatorV25MigrationCompleted, error) {
	tx, err := v.client.Decode(v.coordinator.Migrate(v.client.NewTXOpts(), subId, common.HexToAddress(coordinatorAddress)))
	if err != nil {
		return nil, nil, err
	}
	var migrationCompletedEvent *vrf_coordinator_v2_5.VRFCoordinatorV25MigrationCompleted
	for _, log := range tx.Receipt.Logs {
		for _, topic := range log.Topics {
			if topic.Cmp(vrf_coordinator_v2_5.VRFCoordinatorV25MigrationCompleted{}.Topic()) == 0 {
				migrationCompletedEvent, err = v.coordinator.ParseMigrationCompleted(*log)
				if err != nil {
					return nil, nil, fmt.Errorf("parsing MigrationCompleted log failed, err: %w", err)
				}
			}
		}
	}
	return tx, migrationCompletedEvent, err
}

func (v *EthereumVRFCoordinatorV2_5) RegisterMigratableCoordinator(migratableCoordinatorAddress string) error {
	_, err := v.client.Decode(v.coordinator.RegisterMigratableCoordinator(v.client.NewTXOpts(), common.HexToAddress(migratableCoordinatorAddress)))
	return err
}

func (v *EthereumVRFCoordinatorV2_5) AddConsumer(subId *big.Int, consumerAddress string) error {
	_, err := v.client.Decode(v.coordinator.AddConsumer(
		v.client.NewTXOpts(),
		subId,
		common.HexToAddress(consumerAddress),
	))
	return err
}

func (v *EthereumVRFCoordinatorV2_5) FundSubscriptionWithNative(subId *big.Int, nativeTokenAmount *big.Int) error {
	opts := v.client.NewTXOpts()
	opts.Value = nativeTokenAmount
	_, err := v.client.Decode(v.coordinator.FundSubscriptionWithNative(
		opts,
		subId,
	))
	if err != nil {
		return err
	}
	return nil
}

func (v *EthereumVRFCoordinatorV2_5) FindSubscriptionID(subID *big.Int) (*big.Int, error) {
	owner := v.client.MustGetRootKeyAddress()
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

func (v *EthereumVRFCoordinatorV2_5) FilterRandomWordsFulfilledEvent(opts *bind.FilterOpts, requestId *big.Int) (*CoordinatorRandomWordsFulfilled, error) {
	iterator, err := v.coordinator.FilterRandomWordsFulfilled(
		opts,
		[]*big.Int{requestId},
		nil,
	)
	if err != nil {
		return nil, err
	}
	if !iterator.Next() {
		return nil, fmt.Errorf("expected at least 1 RandomWordsFulfilled event for request Id: %s", requestId.String())
	}
	return &CoordinatorRandomWordsFulfilled{
		RequestId:     iterator.Event.RequestId,
		OutputSeed:    iterator.Event.OutputSeed,
		SubId:         iterator.Event.SubId.String(),
		Payment:       iterator.Event.Payment,
		NativePayment: iterator.Event.NativePayment,
		Success:       iterator.Event.Success,
		OnlyPremium:   iterator.Event.OnlyPremium,
		Raw:           iterator.Event.Raw,
	}, nil
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

func (v *EthereumVRFCoordinatorV2_5_Optimism) SetL1FeeCalculation(
	mode uint8,
	coefficient uint8,
) error {
	_, err := v.client.Decode(v.coordinator.SetL1FeeCalculation(v.client.NewTXOpts(), mode, coefficient))
	return err
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
	return v.RequestRandomnessFromKey(coordinator, keyHash, subID, requestConfirmations, callbackGasLimit, nativePayment, numWords, requestCount, 0)
}

func (v *EthereumVRFv2PlusLoadTestConsumer) RequestRandomnessFromKey(
	coordinator Coordinator,
	keyHash [32]byte, subID *big.Int,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	nativePayment bool,
	numWords uint32,
	requestCount uint16,
	keyNum int,
) (*CoordinatorRandomWordsRequested, error) {
	tx, err := v.client.Decode(v.consumer.RequestRandomWords(v.client.NewTXKeyOpts(keyNum), subID, requestConfirmations, keyHash, callbackGasLimit, nativePayment, numWords, requestCount))
	if err != nil {
		return nil, err
	}
	randomWordsRequestedEvent, err := parseRequestRandomnessLogs(coordinator, tx.Receipt.Logs)
	if err != nil {
		return nil, err
	}
	return randomWordsRequestedEvent, nil
}

func (v *EthereumVRFv2PlusLoadTestConsumer) ResetMetrics() error {
	_, err := v.client.Decode(v.consumer.Reset(v.client.NewTXOpts()))
	return err
}

func (v *EthereumVRFv2PlusLoadTestConsumer) GetCoordinator(ctx context.Context) (common.Address, error) {
	return v.consumer.SVrfCoordinator(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}
func (v *EthereumVRFv2PlusLoadTestConsumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_v2plus_load_test_with_metrics.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, requestID)
}

func (v *EthereumVRFv2PlusLoadTestConsumer) GetLastRequestId(ctx context.Context) (*big.Int, error) {
	return v.consumer.SLastRequestId(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFv2PlusLoadTestConsumer) GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error) {
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
	averageFulfillmentInMillions, err := v.consumer.SAverageResponseTimeInBlocksMillions(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	slowestFulfillment, err := v.consumer.SSlowestResponseTimeInBlocks(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})

	if err != nil {
		return nil, err
	}
	fastestFulfillment, err := v.consumer.SFastestResponseTimeInBlocks(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	averageResponseTimeInSeconds, err := v.consumer.SAverageResponseTimeInSecondsMillions(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	slowestResponseTimeInSeconds, err := v.consumer.SSlowestResponseTimeInSeconds(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	fastestResponseTimeInSeconds, err := v.consumer.SFastestResponseTimeInSeconds(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	var responseTimesInBlocks []uint32
	for {
		currentResponseTimesInBlocks, err := v.consumer.GetRequestBlockTimes(&bind.CallOpts{
			From:    v.client.MustGetRootKeyAddress(),
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

// DeployBatchBlockhashStore deploys DeployBatchBlockhashStore contract
func DeployBatchBlockhashStore(seth *seth.Client, blockhashStoreAddr string) (BatchBlockhashStore, error) {
	abi, err := batch_blockhash_store.BatchBlockhashStoreMetaData.GetAbi()
	if err != nil {
		return &EthereumBatchBlockhashStore{}, fmt.Errorf("failed to get BatchBlockhashStore ABI: %w", err)
	}
	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"BatchBlockhashStore",
		*abi,
		common.FromHex(batch_blockhash_store.BatchBlockhashStoreMetaData.Bin),
		common.HexToAddress(blockhashStoreAddr))
	if err != nil {
		return &EthereumBatchBlockhashStore{}, fmt.Errorf("BatchBlockhashStore instance deployment have failed: %w", err)
	}

	contract, err := batch_blockhash_store.NewBatchBlockhashStore(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumBatchBlockhashStore{}, fmt.Errorf("failed to instantiate BatchBlockhashStore instance: %w", err)
	}

	return &EthereumBatchBlockhashStore{
		client:              seth,
		batchBlockhashStore: contract,
		address:             data.Address,
	}, err
}

func DeployVRFCoordinatorV2PlusUpgradedVersion(client *seth.Client, bhsAddr string) (VRFCoordinatorV2PlusUpgradedVersion, error) {
	abi, err := vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFCoordinatorV2PlusUpgradedVersion{}, fmt.Errorf("failed to get VRFCoordinatorV2PlusUpgradedVersion ABI: %w", err)
	}

	data, err := client.DeployContract(
		client.NewTXOpts(),
		"VRFCoordinatorV2PlusUpgradedVersion",
		*abi,
		common.FromHex(vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionMetaData.Bin),
		common.HexToAddress(bhsAddr))
	if err != nil {
		return &EthereumVRFCoordinatorV2PlusUpgradedVersion{}, fmt.Errorf("VRFCoordinatorV2PlusUpgradedVersion instance deployment have failed: %w", err)
	}

	contract, err := vrf_v2plus_upgraded_version.NewVRFCoordinatorV2PlusUpgradedVersion(data.Address, wrappers.MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumVRFCoordinatorV2PlusUpgradedVersion{}, fmt.Errorf("failed to instantiate VRFCoordinatorV2PlusUpgradedVersion instance: %w", err)
	}

	return &EthereumVRFCoordinatorV2PlusUpgradedVersion{
		client:      client,
		coordinator: contract,
		address:     data.Address,
	}, err
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
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

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) GetActiveSubscriptionIds(ctx context.Context, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
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
		From:    v.client.MustGetRootKeyAddress(),
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

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) SetLINKAndLINKNativeFeed(linkAddress string, linkNativeFeedAddress string) error {
	_, err := v.client.Decode(v.coordinator.SetLINKAndLINKNativeFeed(
		v.client.NewTXOpts(),
		common.HexToAddress(linkAddress),
		common.HexToAddress(linkNativeFeedAddress),
	))
	return err
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) RegisterProvingKey(
	publicProvingKey [2]*big.Int,
	gasLaneMaxGas uint64,
) error {
	_, err := v.client.Decode(v.coordinator.RegisterProvingKey(v.client.NewTXOpts(), publicProvingKey, gasLaneMaxGas))
	return err
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) CreateSubscription() error {
	_, err := v.client.Decode(v.coordinator.CreateSubscription(v.client.NewTXOpts()))
	return err
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) GetLinkTotalBalance(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
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
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}
	totalBalance, err := v.coordinator.STotalNativeBalance(opts)
	if err != nil {
		return nil, err
	}
	return totalBalance, nil
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) Migrate(subId *big.Int, coordinatorAddress string) error {
	_, err := v.client.Decode(v.coordinator.Migrate(v.client.NewTXOpts(), subId, common.HexToAddress(coordinatorAddress)))
	return err
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) RegisterMigratableCoordinator(migratableCoordinatorAddress string) error {
	_, err := v.client.Decode(v.coordinator.RegisterMigratableCoordinator(v.client.NewTXOpts(), common.HexToAddress(migratableCoordinatorAddress)))
	return err
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) AddConsumer(subId *big.Int, consumerAddress string) error {
	_, err := v.client.Decode(v.coordinator.AddConsumer(
		v.client.NewTXOpts(),
		subId,
		common.HexToAddress(consumerAddress),
	))
	return err
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) FundSubscriptionWithNative(subId *big.Int, nativeTokenAmount *big.Int) error {
	opts := v.client.NewTXOpts()
	opts.Value = nativeTokenAmount
	_, err := v.client.Decode(v.coordinator.FundSubscriptionWithNative(
		opts,
		subId,
	))
	return err
}

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) FindSubscriptionID() (*big.Int, error) {
	owner := v.client.MustGetRootKeyAddress()
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

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) FilterRandomWordsFulfilledEvent(opts *bind.FilterOpts, requestId *big.Int) (*CoordinatorRandomWordsFulfilled, error) {
	iterator, err := v.coordinator.FilterRandomWordsFulfilled(
		opts,
		[]*big.Int{requestId},
		nil,
	)
	if err != nil {
		return nil, err
	}
	if !iterator.Next() {
		return nil, fmt.Errorf("expected at least 1 RandomWordsFulfilled event for request Id: %s", requestId.String())
	}
	return &CoordinatorRandomWordsFulfilled{
		RequestId:     iterator.Event.RequestId,
		OutputSeed:    iterator.Event.OutputSeed,
		SubId:         iterator.Event.SubId.String(),
		Payment:       iterator.Event.Payment,
		NativePayment: iterator.Event.NativePayment,
		Success:       iterator.Event.Success,
		OnlyPremium:   iterator.Event.OnlyPremium,
		Raw:           iterator.Event.Raw,
	}, nil
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

func (v *EthereumVRFCoordinatorV2PlusUpgradedVersion) ParseRandomWordsFulfilled(log types.Log) (*CoordinatorRandomWordsFulfilled, error) {
	fulfilled, err := v.coordinator.ParseRandomWordsFulfilled(log)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RandomWordsFulfilled event: %w", err)
	}
	return &CoordinatorRandomWordsFulfilled{
		RequestId:     fulfilled.RequestId,
		OutputSeed:    fulfilled.OutputSeed,
		Payment:       fulfilled.Payment,
		SubId:         fulfilled.SubId.String(),
		NativePayment: fulfilled.NativePayment,
		OnlyPremium:   fulfilled.OnlyPremium,
		Success:       fulfilled.Success,
		Raw:           fulfilled.Raw,
	}, nil
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

func DeployVRFv2PlusLoadTestConsumer(seth *seth.Client, coordinatorAddr string) (VRFv2PlusLoadTestConsumer, error) {
	abi, err := vrf_v2plus_load_test_with_metrics.VRFV2PlusLoadTestWithMetricsMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFv2PlusLoadTestConsumer{}, fmt.Errorf("failed to get VRFV2PlusLoadTestWithMetrics ABI: %w", err)
	}

	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFV2PlusLoadTestWithMetrics",
		*abi,
		common.FromHex(vrf_v2plus_load_test_with_metrics.VRFV2PlusLoadTestWithMetricsMetaData.Bin),
		common.HexToAddress(coordinatorAddr))
	if err != nil {
		return &EthereumVRFv2PlusLoadTestConsumer{}, fmt.Errorf("VRFV2PlusLoadTestWithMetrics instance deployment have failed: %w", err)
	}

	contract, err := vrf_v2plus_load_test_with_metrics.NewVRFV2PlusLoadTestWithMetrics(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFv2PlusLoadTestConsumer{}, fmt.Errorf("failed to instantiate VRFV2PlusLoadTestWithMetrics instance: %w", err)
	}

	return &EthereumVRFv2PlusLoadTestConsumer{
		client:   seth,
		consumer: contract,
		address:  data.Address,
	}, err
}

func DeployVRFV2PlusWrapper(seth *seth.Client, linkAddr string, linkEthFeedAddr string, coordinatorAddr string, subId *big.Int) (VRFV2PlusWrapper, error) {
	abi, err := vrfv2plus_wrapper.VRFV2PlusWrapperMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFV2PlusWrapper{}, fmt.Errorf("failed to get VRFV2PlusWrapper ABI: %w", err)
	}
	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFV2PlusWrapper",
		*abi,
		common.FromHex(vrfv2plus_wrapper.VRFV2PlusWrapperMetaData.Bin),
		common.HexToAddress(linkAddr), common.HexToAddress(linkEthFeedAddr),
		common.HexToAddress(coordinatorAddr), subId)
	if err != nil {
		return &EthereumVRFV2PlusWrapper{}, fmt.Errorf("VRFV2PlusWrapper instance deployment have failed: %w", err)
	}
	contract, err := vrfv2plus_wrapper.NewVRFV2PlusWrapper(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFV2PlusWrapper{}, fmt.Errorf("failed to instantiate VRFV2PlusWrapper instance: %w", err)
	}
	return &EthereumVRFV2PlusWrapper{
		client:  seth,
		wrapper: contract,
		address: data.Address,
	}, err
}

func DeployVRFV2PlusWrapperArbitrum(seth *seth.Client, linkAddr string, linkEthFeedAddr string, coordinatorAddr string, subId *big.Int) (*EthereumVRFV2PlusWrapperArbitrum, error) {
	abi, err := vrfv2plus_wrapper_arbitrum.VRFV2PlusWrapperArbitrumMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFV2PlusWrapper_Arbitrum ABI: %w", err)
	}
	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFV2PlusWrapper_Arbitrum",
		*abi,
		common.FromHex(vrfv2plus_wrapper_arbitrum.VRFV2PlusWrapperArbitrumMetaData.Bin),
		common.HexToAddress(linkAddr), common.HexToAddress(linkEthFeedAddr),
		common.HexToAddress(coordinatorAddr), subId)
	if err != nil {
		return nil, fmt.Errorf("VRFV2PlusWrapper_Arbitrum instance deployment have failed: %w", err)
	}
	contract, err := vrfv2plus_wrapper_arbitrum.NewVRFV2PlusWrapperArbitrum(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFV2PlusWrapper_Arbitrum instance: %w", err)
	}
	return &EthereumVRFV2PlusWrapperArbitrum{
		client:  seth,
		wrapper: contract,
		Address: data.Address,
	}, err
}

func DeployVRFV2PlusWrapperOptimism(seth *seth.Client, linkAddr string, linkEthFeedAddr string, coordinatorAddr string, subId *big.Int) (*EthereumVRFV2PlusWrapperOptimism, error) {
	abi, err := vrfv2plus_wrapper_optimism.VRFV2PlusWrapperOptimismMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get VRFV2PlusWrapper_Optimism ABI: %w", err)
	}
	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFV2PlusWrapper_Optimism",
		*abi,
		common.FromHex(vrfv2plus_wrapper_optimism.VRFV2PlusWrapperOptimismMetaData.Bin),
		common.HexToAddress(linkAddr), common.HexToAddress(linkEthFeedAddr),
		common.HexToAddress(coordinatorAddr), subId)
	if err != nil {
		return nil, fmt.Errorf("VRFV2PlusWrapper_Optimism instance deployment have failed: %w", err)
	}
	contract, err := vrfv2plus_wrapper_optimism.NewVRFV2PlusWrapperOptimism(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate VRFV2PlusWrapper_Optimism instance: %w", err)
	}
	return &EthereumVRFV2PlusWrapperOptimism{
		client:  seth,
		wrapper: contract,
		Address: data.Address,
	}, err
}

func DeployVRFV2PlusWrapperLoadTestConsumer(seth *seth.Client, vrfV2PlusWrapperAddr string) (VRFv2PlusWrapperLoadTestConsumer, error) {
	abi, err := vrfv2plus_wrapper_load_test_consumer.VRFV2PlusWrapperLoadTestConsumerMetaData.GetAbi()
	if err != nil {
		return &EthereumVRFV2PlusWrapperLoadTestConsumer{}, fmt.Errorf("failed to get VRFV2PlusWrapperLoadTestConsumer ABI: %w", err)
	}
	data, err := seth.DeployContract(
		seth.NewTXOpts(),
		"VRFV2PlusWrapperLoadTestConsumer",
		*abi,
		common.FromHex(vrfv2plus_wrapper_load_test_consumer.VRFV2PlusWrapperLoadTestConsumerMetaData.Bin),
		common.HexToAddress(vrfV2PlusWrapperAddr))
	if err != nil {
		return &EthereumVRFV2PlusWrapperLoadTestConsumer{}, fmt.Errorf("VRFV2PlusWrapperLoadTestConsumer instance deployment have failed: %w", err)
	}

	contract, err := vrfv2plus_wrapper_load_test_consumer.NewVRFV2PlusWrapperLoadTestConsumer(data.Address, wrappers.MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumVRFV2PlusWrapperLoadTestConsumer{}, fmt.Errorf("failed to instantiate VRFV2PlusWrapperLoadTestConsumer instance: %w", err)
	}
	return &EthereumVRFV2PlusWrapperLoadTestConsumer{
		client:   seth,
		consumer: contract,
		address:  data.Address,
	}, err
}

func (v *EthereumVRFV2PlusWrapperOptimism) SetL1FeeCalculation(mode uint8, coefficient uint8) error {
	_, err := v.client.Decode(v.wrapper.SetL1FeeCalculation(v.client.NewTXOpts(), mode, coefficient))
	return err
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) Fund(_ *big.Float) error {
	panic("do not use this function, use actions.SendFunds() instead, otherwise we will have to deal with circular dependencies")
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) RequestRandomness(
	coordinator Coordinator,
	requestConfirmations uint16,
	callbackGasLimit uint32,
	numWords uint32,
	requestCount uint16,
) (*CoordinatorRandomWordsRequested, error) {
	tx, err := v.client.Decode(v.consumer.MakeRequests(v.client.NewTXOpts(), callbackGasLimit, requestConfirmations, numWords, requestCount))
	if err != nil {
		return nil, err
	}
	randomWordsRequestedEvent, err := parseRequestRandomnessLogs(coordinator, tx.Receipt.Logs)
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
	tx, err := v.client.Decode(v.consumer.MakeRequestsNative(v.client.NewTXOpts(), callbackGasLimit, requestConfirmations, numWords, requestCount))
	if err != nil {
		return nil, err
	}
	randomWordsRequestedEvent, err := parseRequestRandomnessLogs(coordinator, tx.Receipt.Logs)
	if err != nil {
		return nil, err
	}
	return randomWordsRequestedEvent, nil
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) GetRequestStatus(ctx context.Context, requestID *big.Int) (vrfv2plus_wrapper_load_test_consumer.GetRequestStatus, error) {
	return v.consumer.GetRequestStatus(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	}, requestID)
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) GetLastRequestId(ctx context.Context) (*big.Int, error) {
	return v.consumer.SLastRequestId(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) GetWrapper(ctx context.Context) (common.Address, error) {
	return v.consumer.IVrfV2PlusWrapper(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func (v *EthereumVRFV2PlusWrapperLoadTestConsumer) GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error) {
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
		AverageResponseTimeInSecondsMillions: nil,
		SlowestResponseTimeInSeconds:         nil,
		FastestResponseTimeInSeconds:         nil,
	}, nil
}
