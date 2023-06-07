package contracts

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_load_test_with_metrics"
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
	GetRequestStatus(ctx context.Context, requestID *big.Int) (RequestStatus, error)
	GetLastRequestId(ctx context.Context) (*big.Int, error)
}

type VRFv2LoadTestConsumer interface {
	Address() string
	RequestRandomness(hash [32]byte, subID uint64, confs uint16, gasLimit uint32, numWords uint32, requestCount uint16) error
	GetRequestStatus(ctx context.Context, requestID *big.Int) (vrf_load_test_with_metrics.GetRequestStatus, error)
	GetLastRequestId(ctx context.Context) (*big.Int, error)
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
