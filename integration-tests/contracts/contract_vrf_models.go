package contracts

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_load_test_with_metrics"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2_consumer_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_load_test_with_metrics"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_upgraded_version"
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
	HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error)
	CreateSubscription() error
	AddConsumer(subId uint64, consumerAddress string) error
	Address() string
	GetSubscription(ctx context.Context, subID uint64) (vrf_coordinator_v2.GetSubscription, error)
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
		feeConfig vrf_coordinator_v2_5.VRFCoordinatorV25FeeConfig,
	) error
	RegisterProvingKey(
		oracleAddr string,
		publicProvingKey [2]*big.Int,
	) error
	HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error)
	CreateSubscription() error
	GetActiveSubscriptionIds(ctx context.Context, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)
	Migrate(subId *big.Int, coordinatorAddress string) error
	RegisterMigratableCoordinator(migratableCoordinatorAddress string) error
	AddConsumer(subId *big.Int, consumerAddress string) error
	FundSubscriptionWithNative(subId *big.Int, nativeTokenAmount *big.Int) error
	Address() string
	GetSubscription(ctx context.Context, subID *big.Int) (vrf_coordinator_v2_5.GetSubscription, error)
	GetNativeTokenTotalBalance(ctx context.Context) (*big.Int, error)
	GetLinkTotalBalance(ctx context.Context) (*big.Int, error)
	FindSubscriptionID() (*big.Int, error)
	WaitForSubscriptionCreatedEvent(timeout time.Duration) (*vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCreated, error)
	WaitForRandomWordsFulfilledEvent(subID []*big.Int, requestID []*big.Int, timeout time.Duration) (*vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled, error)
	WaitForRandomWordsRequestedEvent(keyHash [][32]byte, subID []*big.Int, sender []common.Address, timeout time.Duration) (*vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsRequested, error)
	WaitForMigrationCompletedEvent(timeout time.Duration) (*vrf_coordinator_v2_5.VRFCoordinatorV25MigrationCompleted, error)
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
		feeConfig vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionFeeConfig,
	) error
	RegisterProvingKey(
		oracleAddr string,
		publicProvingKey [2]*big.Int,
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
	WaitForRandomWordsFulfilledEvent(subID []*big.Int, requestID []*big.Int, timeout time.Duration) (*vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled, error)
	WaitForMigrationCompletedEvent(timeout time.Duration) (*vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted, error)
	WaitForRandomWordsRequestedEvent(keyHash [][32]byte, subID []*big.Int, sender []common.Address, timeout time.Duration) (*vrf_v2plus_upgraded_version.VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested, error)
}

type VRFV2PlusWrapper interface {
	Address() string
	SetConfig(wrapperGasOverhead uint32, coordinatorGasOverhead uint32, wrapperPremiumPercentage uint8, keyHash [32]byte, maxNumWords uint8, stalenessSeconds uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeLinkPPM uint32, fulfillmentFlatFeeNativePPM uint32) error
	GetSubID(ctx context.Context) (*big.Int, error)
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
	RequestRandomness(hash [32]byte, subID uint64, confs uint16, gasLimit uint32, numWords uint32, requestCount uint16) error
	GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_load_test_with_metrics.GetRequestStatus, error)
	GetLastRequestId(ctx context.Context) (*big.Int, error)
	GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error)
}

type VRFv2PlusLoadTestConsumer interface {
	Address() string
	RequestRandomness(keyHash [32]byte, subID *big.Int, requestConfirmations uint16, callbackGasLimit uint32, nativePayment bool, numWords uint32, requestCount uint16) (*types.Transaction, error)
	GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_v2plus_load_test_with_metrics.GetRequestStatus, error)
	GetLastRequestId(ctx context.Context) (*big.Int, error)
	GetLoadTestMetrics(ctx context.Context) (*VRFLoadTestMetrics, error)
	GetCoordinator(ctx context.Context) (common.Address, error)
	ResetMetrics() error
}

type VRFv2PlusWrapperLoadTestConsumer interface {
	Address() string
	Fund(ethAmount *big.Float) error
	RequestRandomness(requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, requestCount uint16) (*types.Transaction, error)
	RequestRandomnessNative(requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, requestCount uint16) (*types.Transaction, error)
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

type RequestStatus struct {
	Fulfilled   bool
	RandomWords []*big.Int
}

type LoadTestRequestStatus struct {
	Fulfilled             bool
	RandomWords           []*big.Int
	requestTimestamp      *big.Int
	fulfilmentTimestamp   *big.Int
	requestBlockNumber    *big.Int
	fulfilmentBlockNumber *big.Int
}

type VRFLoadTestMetrics struct {
	RequestCount                 *big.Int
	FulfilmentCount              *big.Int
	AverageFulfillmentInMillions *big.Int
	SlowestFulfillment           *big.Int
	FastestFulfillment           *big.Int
}
