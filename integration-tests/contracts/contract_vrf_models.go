package contracts

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_load_test_with_metrics"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_owner"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2_consumer_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_load_test_with_metrics"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_upgraded_version"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2_wrapper_load_test_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrfv2plus_wrapper_load_test_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
)

type VRF interface {
	Fund(ethAmount *big.Float) error
	ProofLength(context.Context) (*big.Int, error)
}

type VRFCoordinator interface {
	RegisterProvingKey(
		fee *big.Int,
		oracleAddr string,
		publicProvingKey [2]*big.Int,
		jobID [32]byte,
	) error
	HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error)
	Address() string
}

type VRFCoordinatorV2 interface {
	GetRequestConfig(ctx context.Context) (GetRequestConfig, error)
	GetConfig(ctx context.Context) (vrf_coordinator_v2.GetConfig, error)
	GetFallbackWeiPerUnitLink(ctx context.Context) (*big.Int, error)
	GetFeeConfig(ctx context.Context) (vrf_coordinator_v2.GetFeeConfig, error)
	SetConfig(
		minimumRequestConfirmations uint16,
		maxGasLimit uint32,
		stalenessSeconds uint32,
		gasAfterPaymentCalculation uint32,
		fallbackWeiPerUnitLink *big.Int,
		feeConfig vrf_coordinator_v2.VRFCoordinatorV2FeeConfig,
	) error
	RegisterProvingKey(
		oracleAddr string,
		publicProvingKey [2]*big.Int,
	) error
	TransferOwnership(to common.Address) error
	HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error)
	CreateSubscription() (*types.Receipt, error)
	AddConsumer(subId uint64, consumerAddress string) error
	Address() string
	GetSubscription(ctx context.Context, subID uint64) (Subscription, error)
	GetOwner(ctx context.Context) (common.Address, error)
	PendingRequestsExist(ctx context.Context, subID uint64) (bool, error)
	OwnerCancelSubscription(subID uint64) (*seth.DecodedTransaction, *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled, error)
	CancelSubscription(subID uint64, to common.Address) (*seth.DecodedTransaction, *vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled, error)
	ParseSubscriptionCanceled(log types.Log) (*vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled, error)
	ParseRandomWordsRequested(log types.Log) (*CoordinatorRandomWordsRequested, error)
	ParseRandomWordsFulfilled(log types.Log) (*CoordinatorRandomWordsFulfilled, error)
	ParseLog(log types.Log) (generated.AbigenLog, error)
	FindSubscriptionID(subID uint64) (uint64, error)
	WaitForRandomWordsFulfilledEvent(filter RandomWordsFulfilledEventFilter) (*CoordinatorRandomWordsFulfilled, error)
	WaitForConfigSetEvent(timeout time.Duration) (*CoordinatorConfigSet, error)
	OracleWithdraw(recipient common.Address, amount *big.Int) error
}

type VRFCoordinatorV2_5 interface {
	SetLINKAndLINKNativeFeed(
		link string,
		linkNativeFeed string,
	) error
	SetConfig(
		minimumRequestConfirmations uint16,
		maxGasLimit uint32,
		stalenessSeconds uint32,
		gasAfterPaymentCalculation uint32,
		fallbackWeiPerUnitLink *big.Int,
		fulfillmentFlatFeeNativePPM uint32,
		fulfillmentFlatFeeLinkDiscountPPM uint32,
		nativePremiumPercentage uint8,
		linkPremiumPercentage uint8,
	) error
	RegisterProvingKey(
		publicProvingKey [2]*big.Int,
		gasLaneMaxGas uint64,
	) error
	HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error)
	CreateSubscription() (*types.Transaction, error)
	GetActiveSubscriptionIds(ctx context.Context, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)
	Migrate(subId *big.Int, coordinatorAddress string) (*seth.DecodedTransaction, *vrf_coordinator_v2_5.VRFCoordinatorV25MigrationCompleted, error)
	RegisterMigratableCoordinator(migratableCoordinatorAddress string) error
	AddConsumer(subId *big.Int, consumerAddress string) error
	FundSubscriptionWithNative(subId *big.Int, nativeTokenAmount *big.Int) error
	Address() string
	PendingRequestsExist(ctx context.Context, subID *big.Int) (bool, error)
	GetSubscription(ctx context.Context, subID *big.Int) (Subscription, error)
	OwnerCancelSubscription(subID *big.Int) (*seth.DecodedTransaction, *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled, error)
	CancelSubscription(subID *big.Int, to common.Address) (*seth.DecodedTransaction, *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled, error)
	Withdraw(recipient common.Address) error
	WithdrawNative(recipient common.Address) error
	GetNativeTokenTotalBalance(ctx context.Context) (*big.Int, error)
	GetLinkTotalBalance(ctx context.Context) (*big.Int, error)
	FindSubscriptionID(subID *big.Int) (*big.Int, error)
	WaitForRandomWordsFulfilledEvent(filter RandomWordsFulfilledEventFilter) (*CoordinatorRandomWordsFulfilled, error)
	ParseRandomWordsRequested(log types.Log) (*CoordinatorRandomWordsRequested, error)
	ParseRandomWordsFulfilled(log types.Log) (*CoordinatorRandomWordsFulfilled, error)
	WaitForConfigSetEvent(timeout time.Duration) (*CoordinatorConfigSet, error)
}

type VRFCoordinatorV2PlusUpgradedVersion interface {
	SetLINKAndLINKNativeFeed(
		link string,
		linkNativeFeed string,
	) error
	SetConfig(
		minimumRequestConfirmations uint16,
		maxGasLimit uint32,
		stalenessSeconds uint32,
		gasAfterPaymentCalculation uint32,
		fallbackWeiPerUnitLink *big.Int,
		fulfillmentFlatFeeNativePPM uint32,
		fulfillmentFlatFeeLinkDiscountPPM uint32,
		nativePremiumPercentage uint8,
		linkPremiumPercentage uint8,
	) error
	RegisterProvingKey(
		publicProvingKey [2]*big.Int,
		gasLaneMaxGasPrice uint64,
	) error
	HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error)
	CreateSubscription() error
	GetNativeTokenTotalBalance(ctx context.Context) (*big.Int, error)
	GetLinkTotalBalance(ctx context.Context) (*big.Int, error)
	Migrate(subId *big.Int, coordinatorAddress string) error
	RegisterMigratableCoordinator(migratableCoordinatorAddress string) error
	AddConsumer(subId *big.Int, consumerAddress string) error
	FundSubscriptionWithNative(subId *big.Int, nativeTokenAmount *big.Int) error
	Address() string
	GetSubscription(ctx context.Context, subID *big.Int) (vrf_v2plus_upgraded_version.GetSubscription, error)
	GetActiveSubscriptionIds(ctx context.Context, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)
	FindSubscriptionID() (*big.Int, error)
	WaitForRandomWordsFulfilledEvent(filter RandomWordsFulfilledEventFilter) (*CoordinatorRandomWordsFulfilled, error)
	ParseRandomWordsRequested(log types.Log) (*CoordinatorRandomWordsRequested, error)
	ParseRandomWordsFulfilled(log types.Log) (*CoordinatorRandomWordsFulfilled, error)
	WaitForConfigSetEvent(timeout time.Duration) (*CoordinatorConfigSet, error)
}

type VRFV2Wrapper interface {
	Address() string
	SetConfig(wrapperGasOverhead uint32, coordinatorGasOverhead uint32, wrapperPremiumPercentage uint8, keyHash [32]byte, maxNumWords uint8) error
	GetSubID(ctx context.Context) (uint64, error)
}

type VRFV2PlusWrapper interface {
	Address() string
	SetConfig(
		wrapperGasOverhead uint32,
		coordinatorGasOverheadNative uint32,
		coordinatorGasOverheadLink uint32,
		coordinatorGasOverheadPerWord uint16,
		wrapperNativePremiumPercentage uint8,
		wrapperLinkPremiumPercentage uint8,
		keyHash [32]byte, maxNumWords uint8,
		stalenessSeconds uint32,
		fallbackWeiPerUnitLink *big.Int,
		fulfillmentFlatFeeNativePPM uint32,
		fulfillmentFlatFeeLinkDiscountPPM uint32,
	) error
	GetSubID(ctx context.Context) (*big.Int, error)
	Coordinator(ctx context.Context) (common.Address, error)
}

type VRFOwner interface {
	Address() string
	SetAuthorizedSenders(senders []common.Address) error
	AcceptVRFOwnership() error
	WaitForRandomWordsForcedEvent(requestIDs []*big.Int, subIds []uint64, senders []common.Address, timeout time.Duration) (*vrf_owner.VRFOwnerRandomWordsForced, error)
	OwnerCancelSubscription(subID uint64) (*types.Transaction, error)
}

type VRFConsumer interface {
	Address() string
	RequestRandomness(hash [32]byte, fee *big.Int) error
	CurrentRoundID(ctx context.Context) (*big.Int, error)
	RandomnessOutput(ctx context.Context) (*big.Int, error)
	Fund(ethAmount *big.Float) error
}

type VRFConsumerV2 interface {
	Address() string
	CurrentSubscription() (uint64, error)
	CreateFundedSubscription(funds *big.Int) error
	TopUpSubscriptionFunds(funds *big.Int) error
	RequestRandomness(hash [32]byte, subID uint64, confs uint16, gasLimit uint32, numWords uint32) error
	RandomnessOutput(ctx context.Context, arg0 *big.Int) (*big.Int, error)
	GetAllRandomWords(ctx context.Context, num int) ([]*big.Int, error)
	GasAvailable() (*big.Int, error)
	Fund(ethAmount *big.Float) error
}

type VRFv2Consumer interface {
	Address() string
	RequestRandomness(hash [32]byte, subID uint64, confs uint16, gasLimit uint32, numWords uint32) error
	GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_v2_consumer_wrapper.GetRequestStatus, error)
	GetLastRequestId(ctx context.Context) (*big.Int, error)
}

type VRFv2LoadTestConsumer interface {
	Address() string
	RequestRandomness(
		coordinator Coordinator,
		keyHash [32]byte,
		subID uint64,
		requestConfirmations uint16,
		callbackGasLimit uint32,
		numWords uint32,
		requestCount uint16,
	) (*CoordinatorRandomWordsRequested, error)
	RequestRandomnessFromKey(
		coordinator Coordinator,
		keyHash [32]byte,
		subID uint64,
		requestConfirmations uint16,
		callbackGasLimit uint32,
		numWords uint32,
		requestCount uint16,
		keyNum int,
	) (*CoordinatorRandomWordsRequested, error)
	RequestRandomWordsWithForceFulfill(
		coordinator Coordinator,
		keyHash [32]byte,
		requestConfirmations uint16,
		callbackGasLimit uint32,
		numWords uint32,
		requestCount uint16,
		subTopUpAmount *big.Int,
		linkAddress common.Address,
	) (*CoordinatorRandomWordsRequested, error)
	GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_load_test_with_metrics.GetRequestStatus, error)
	GetLastRequestId(ctx context.Context) (*big.Int, error)
	GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error)
	ResetMetrics() error
}

type VRFv2WrapperLoadTestConsumer interface {
	Address() string
	Fund(ethAmount *big.Float) error
	RequestRandomness(coordinator Coordinator, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, requestCount uint16) (*CoordinatorRandomWordsRequested, error)
	GetRequestStatus(ctx context.Context, requestID *big.Int) (vrfv2_wrapper_load_test_consumer.GetRequestStatus, error)
	GetLastRequestId(ctx context.Context) (*big.Int, error)
	GetWrapper(ctx context.Context) (common.Address, error)
	GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error)
}

type VRFv2PlusLoadTestConsumer interface {
	Address() string
	RequestRandomness(
		coordinator Coordinator,
		keyHash [32]byte, subID *big.Int,
		requestConfirmations uint16,
		callbackGasLimit uint32,
		nativePayment bool,
		numWords uint32,
		requestCount uint16,
	) (*CoordinatorRandomWordsRequested, error)
	RequestRandomnessFromKey(
		coordinator Coordinator,
		keyHash [32]byte, subID *big.Int,
		requestConfirmations uint16,
		callbackGasLimit uint32,
		nativePayment bool,
		numWords uint32,
		requestCount uint16,
		keyNum int,
	) (*CoordinatorRandomWordsRequested, error)
	GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_v2plus_load_test_with_metrics.GetRequestStatus, error)
	GetLastRequestId(ctx context.Context) (*big.Int, error)
	GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error)
	GetCoordinator(ctx context.Context) (common.Address, error)
	ResetMetrics() error
}

type VRFv2PlusWrapperLoadTestConsumer interface {
	Address() string
	Fund(ethAmount *big.Float) error
	RequestRandomness(
		coordinator Coordinator,
		requestConfirmations uint16,
		callbackGasLimit uint32,
		numWords uint32,
		requestCount uint16,
	) (*CoordinatorRandomWordsRequested, error)
	RequestRandomnessNative(
		coordinator Coordinator,
		requestConfirmations uint16,
		callbackGasLimit uint32,
		numWords uint32,
		requestCount uint16,
	) (*CoordinatorRandomWordsRequested, error)
	GetRequestStatus(ctx context.Context, requestID *big.Int) (vrfv2plus_wrapper_load_test_consumer.GetRequestStatus, error)
	GetLastRequestId(ctx context.Context) (*big.Int, error)
	GetWrapper(ctx context.Context) (common.Address, error)
	GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error)
}

type DKG interface {
	Address() string
	AddClient(keyID string, clientAddress string) error
	SetConfig(
		signerAddresses []common.Address,
		transmitterAddresses []common.Address,
		f uint8,
		onchainConfig []byte,
		offchainConfigVersion uint64,
		offchainConfig []byte,
	) error
	WaitForConfigSetEvent(timeout time.Duration) (*dkg.DKGConfigSet, error)
	WaitForTransmittedEvent(timeout time.Duration) (*dkg.DKGTransmitted, error)
}

type VRFCoordinatorV3 interface {
	Address() string
	SetProducer(producerAddress string) error
	CreateSubscription() error
	FindSubscriptionID() (*big.Int, error)
	AddConsumer(subId *big.Int, consumerAddress string) error
	SetConfig(maxCallbackGasLimit, maxCallbackArgumentsLength uint32) error
}

type VRFBeacon interface {
	Address() string
	SetPayees(transmitterAddresses []common.Address, payeesAddresses []common.Address) error
	SetConfig(
		signerAddresses []common.Address,
		transmitterAddresses []common.Address,
		f uint8,
		onchainConfig []byte,
		offchainConfigVersion uint64,
		offchainConfig []byte,
	) error
	WaitForConfigSetEvent(timeout time.Duration) (*vrf_beacon.VRFBeaconConfigSet, error)
	WaitForNewTransmissionEvent(timeout time.Duration) (*vrf_beacon.VRFBeaconNewTransmission, error)
	LatestConfigDigestAndEpoch(ctx context.Context) (vrf_beacon.LatestConfigDigestAndEpoch,
		error)
}

type VRFBeaconConsumer interface {
	Address() string
	RequestRandomness(
		numWords uint16,
		subID, confirmationDelayArg *big.Int,
	) (*types.Receipt, error)
	RedeemRandomness(subID, requestID *big.Int) error
	RequestRandomnessFulfillment(
		numWords uint16,
		subID, confirmationDelayArg *big.Int,
		requestGasLimit,
		callbackGasLimit uint32,
		arguments []byte,
	) (*types.Receipt, error)
	IBeaconPeriodBlocks(ctx context.Context) (*big.Int, error)
	GetRequestIdsBy(ctx context.Context, nextBeaconOutputHeight *big.Int, confDelay *big.Int) (*big.Int, error)
	GetRandomnessByRequestId(ctx context.Context, requestID *big.Int, numWordIndex *big.Int) (*big.Int, error)
}

type BatchBlockhashStore interface {
	Address() string
}

type BatchVRFCoordinatorV2 interface {
	Address() string
}
type BatchVRFCoordinatorV2Plus interface {
	Address() string
}

type VRFMockETHLINKFeed interface {
	Address() string
	LatestRoundData() (*big.Int, error)
	LatestRoundDataUpdatedAt() (*big.Int, error)
	SetBlockTimestampDeduction(blockTimestampDeduction *big.Int) error
}

type RequestStatus struct {
	Fulfilled   bool
	RandomWords []*big.Int
}

type LoadTestRequestStatus struct {
	Fulfilled   bool
	RandomWords []*big.Int
	// Currently Unused November 8, 2023, Mignt be used in near future, will remove if not.
	// requestTimestamp      *big.Int
	// fulfilmentTimestamp   *big.Int
	// requestBlockNumber    *big.Int
	// fulfilmentBlockNumber *big.Int
}

type VRFLoadTestMetrics struct {
	RequestCount                         *big.Int
	FulfilmentCount                      *big.Int
	AverageFulfillmentInMillions         *big.Int
	SlowestFulfillment                   *big.Int
	FastestFulfillment                   *big.Int
	P90FulfillmentBlockTime              float64
	P95FulfillmentBlockTime              float64
	AverageResponseTimeInSecondsMillions *big.Int
	SlowestResponseTimeInSeconds         *big.Int
	FastestResponseTimeInSeconds         *big.Int
}
