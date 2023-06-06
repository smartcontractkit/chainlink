// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_coordinator

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

type ECCArithmeticG1Point struct {
	P [2]*big.Int
}

type VRFBeaconTypesBillingConfig struct {
	UseReasonableGasPrice             bool
	UnusedGasPenaltyPercent           uint8
	StalenessSeconds                  uint32
	RedeemableRequestGasOverhead      uint32
	CallbackRequestGasOverhead        uint32
	PremiumPercentage                 uint32
	ReasonableGasPriceStalenessBlocks uint32
	FallbackWeiPerUnitLink            *big.Int
}

type VRFBeaconTypesCallback struct {
	RequestID      *big.Int
	NumWords       uint16
	Requester      common.Address
	Arguments      []byte
	GasAllowance   *big.Int
	SubID          *big.Int
	GasPrice       *big.Int
	WeiPerUnitLink *big.Int
}

type VRFBeaconTypesCostedCallback struct {
	Callback VRFBeaconTypesCallback
	Price    *big.Int
}

type VRFBeaconTypesOutputServed struct {
	Height            uint64
	ConfirmationDelay *big.Int
	ProofG1X          *big.Int
	ProofG1Y          *big.Int
}

type VRFBeaconTypesVRFOutput struct {
	BlockHeight       uint64
	ConfirmationDelay *big.Int
	VrfOutput         ECCArithmeticG1Point
	Callbacks         []VRFBeaconTypesCostedCallback
}

type VRFCoordinatorConfig struct {
	MaxCallbackGasLimit        uint32
	MaxCallbackArgumentsLength uint32
}

var VRFCoordinatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocksArg\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BeaconPeriodMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"earliestAllowed\",\"type\":\"uint256\"}],\"name\":\"BlockTooRecent\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"providedLength\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxLength\",\"type\":\"uint32\"}],\"name\":\"CallbackArgumentsLengthTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confirmationDelays\",\"type\":\"uint16[10]\"},{\"internalType\":\"uint8\",\"name\":\"violatingIndex\",\"type\":\"uint8\"}],\"name\":\"ConfirmationDelaysNotIncreasing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasAllowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLeft\",\"type\":\"uint256\"}],\"name\":\"GasAllowanceExceedsGasLeft\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"providedLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxLimit\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reportHeight\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"separatorHeight\",\"type\":\"uint64\"}],\"name\":\"HistoryDomainSeparatorTooOld\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"actualBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requiredBalance\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidBillingConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"expectedLength\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"actualLength\",\"type\":\"uint256\"}],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidJuelsConversion\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numRecipients\",\"type\":\"uint256\"}],\"name\":\"InvalidNumberOfRecipients\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestedSubID\",\"type\":\"uint256\"}],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"requestedVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"coordinatorVersion\",\"type\":\"uint8\"}],\"name\":\"MigrationVersionMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProducer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativePaymentGiven\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoWordsRequested\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confDelays\",\"type\":\"uint16[10]\"}],\"name\":\"NonZeroDelayAfterZeroDelay\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnMigrationNotSupported\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"producer\",\"type\":\"address\"}],\"name\":\"ProducerAlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"}],\"name\":\"RandomnessNotAvailable\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numRecipients\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numPayments\",\"type\":\"uint256\"}],\"name\":\"RecipientsPaymentsMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"actual\",\"type\":\"address\"}],\"name\":\"ResponseMustBeRetrievedByRequester\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyRequestsReplaceContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManySlotsReplaceContract\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"max\",\"type\":\"uint256\"}],\"name\":\"TooManyWords\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockHeight\",\"type\":\"uint256\"}],\"name\":\"UniverseHasEndedBangBangBang\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"givenDelay\",\"type\":\"uint24\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"knownDelays\",\"type\":\"uint24[8]\"}],\"name\":\"UnknownConfirmationDelay\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"useReasonableGasPrice\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"unusedGasPenaltyPercent\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"redeemableRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callbackRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"premiumPercentage\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceStalenessBlocks\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.BillingConfig\",\"name\":\"billingConfig\",\"type\":\"tuple\"}],\"name\":\"BillingConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCallbackArgumentsLength\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structVRFCoordinator.Config\",\"name\":\"newConfig\",\"type\":\"tuple\"}],\"name\":\"CoordinatorConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"newVersion\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputsServed\",\"type\":\"tuple[]\"}],\"name\":\"OutputsServed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.RequestID[]\",\"name\":\"requestIDs\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"successfulFulfillment\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"truncatedErrorData\",\"type\":\"bytes[]\"},{\"indexed\":false,\"internalType\":\"uint96[]\",\"name\":\"subBalances\",\"type\":\"uint96[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"subIDs\",\"type\":\"uint256[]\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAllowance\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"costJuels\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newSubBalance\",\"type\":\"uint256\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"}],\"name\":\"RandomnessRedeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"costJuels\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newSubBalance\",\"type\":\"uint256\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_JUELS_SUPPLY\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"paymentsInJuels\",\"type\":\"uint256[]\"}],\"name\":\"batchTransferLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"calculateRequestPriceCallbackJuels\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateRequestPriceJuels\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"deregisterMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"getCallbackMemo\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfirmationDelays\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"\",\"type\":\"uint24[8]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"getFulfillmentFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"pendingFulfillments\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSubscriptionLinkBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_link\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_link_eth_feed\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_startSlot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVRFMigration\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"encodedRequest\",\"type\":\"bytes\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrationVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onMigration\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"}],\"internalType\":\"structECCArithmetic.G1Point\",\"name\":\"vrfOutput\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"}],\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"},{\"internalType\":\"uint96\",\"name\":\"price\",\"type\":\"uint96\"}],\"internalType\":\"structVRFBeaconTypes.CostedCallback[]\",\"name\":\"callbacks\",\"type\":\"tuple[]\"}],\"internalType\":\"structVRFBeaconTypes.VRFOutput[]\",\"name\":\"vrfOutputs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"processVRFOutputs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestIDArg\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"redeemRandomness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"randomness\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"registerMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint24\",\"name\":\"confDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint24\",\"name\":\"confDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"requestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_billingConfig\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"useReasonableGasPrice\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"unusedGasPenaltyPercent\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"redeemableRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callbackRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"premiumPercentage\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceStalenessBlocks\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCallbackArgumentsLength\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_pendingRequests\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.SlotNumber\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_producer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"useReasonableGasPrice\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"unusedGasPenaltyPercent\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"redeemableRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callbackRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"premiumPercentage\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceStalenessBlocks\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"internalType\":\"structVRFBeaconTypes.BillingConfig\",\"name\":\"billingConfig\",\"type\":\"tuple\"}],\"name\":\"setBillingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCallbackArgumentsLength\",\"type\":\"uint32\"}],\"internalType\":\"structVRFCoordinator.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"confDelays\",\"type\":\"uint24[8]\"}],\"name\":\"setConfirmationDelays\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"producer\",\"type\":\"address\"}],\"name\":\"setProducer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"gasPrice\",\"type\":\"uint64\"}],\"name\":\"setReasonableGasPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"juelsAmount\",\"type\":\"uint256\"}],\"name\":\"transferLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162006206380380620062068339810160408190526200003591620002f0565b818133806000816200008e5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c157620000c18162000199565b50506001805460ff60a01b19169055506001600160a01b039182166080521660a05260008390036200010657604051632abc297960e01b815260040160405180910390fd5b60c083905260006200011762000244565b9050600060c051826200012b919062000331565b905060008160c0516200013f91906200036a565b90506200014d818462000386565b60e0526010805465ffffffffffff169060006200016a836200039c565b91906101000a81548165ffffffffffff021916908365ffffffffffff16021790555050505050505050620003de565b336001600160a01b03821603620001f35760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000085565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60004661a4b18114806200025a575062066eed81145b15620002cc5760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620002a0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002c69190620003c4565b91505090565b4391505090565b80516001600160a01b0381168114620002eb57600080fd5b919050565b6000806000606084860312156200030657600080fd5b835192506200031860208501620002d3565b91506200032860408501620002d3565b90509250925092565b6000826200034f57634e487b7160e01b600052601260045260246000fd5b500690565b634e487b7160e01b600052601160045260246000fd5b8181038181111562000380576200038062000354565b92915050565b8082018082111562000380576200038062000354565b600065ffffffffffff808316818103620003ba57620003ba62000354565b6001019392505050565b600060208284031215620003d757600080fd5b5051919050565b60805160a05160c05160e051615da162000465600039600061039b01526000818161086801528181611b8e015281816140090152818161403801528181614070015261412f0152600081816109cb015261449401526000818161067001528181610f350152818161192e015281816122d501528181612d1f0152612db40152615da16000f3fe6080604052600436106102385760003560e01c80625bd5241461023d57806304104edb1461027357806305f4acc614610295578063088070f5146102b55780630ae09540146102eb57806316f6ee9a1461030b578063181f5a77146103465780631a961f2514610389578063294daa49146103bd5780632b38bafc146103df5780632d9297b0146103ff5780632f7527cc1461042c5780633f4ba83a1461044157806340d6bb821461045657806346942d181461046c57806347c3e2cb1461048c5780634ffac83a14610527578063597d2f3c1461053a5780635c975abb146105585780635d06b4ab1461057d57806364d51a2a1461059d57806369a3164e146105c557806373433a2f1461060a578063747c91f71461062a57806379ba5097146106495780637d253aff1461065e5780638456cb591461069f57806385c64e11146106b45780638da5cb5b146106d65780638eef585f146106f45780639e20103614610714578063a21a23e414610734578063a4c0ed3614610749578063acfc6cdd14610769578063b072419514610796578063b2a7cac5146107b6578063bd58017f146107d6578063bec4c08c146107f6578063c3fbb6fd14610816578063cb63179714610836578063cd0593df14610856578063ce3f47191461088a578063daa28f5e146108aa578063dac83d2914610956578063db972c8b14610976578063dc311dd314610989578063f27fcfb8146109b9578063f2fde38b146109ed578063f99b1d6814610a0d578063f9c45ced14610a2d575b600080fd5b34801561024957600080fd5b5061025d61025836600461487c565b610a4d565b60405161026a919061496f565b60405180910390f35b34801561027f57600080fd5b5061029361028e3660046149a2565b610cec565b005b3480156102a157600080fd5b506102936102b03660046149bf565b610d64565b3480156102c157600080fd5b50600d546102dd9063ffffffff80821691600160201b90041682565b60405161026a9291906149da565b3480156102f757600080fd5b506102936103063660046149f1565b610de1565b34801561031757600080fd5b50610338610326366004614a21565b6000908152600e602052604090205490565b60405190815260200161026a565b34801561035257600080fd5b5060408051808201825260148152730565246436f6f7264696e61746f7220312e302e360641b6020820152905161026a9190614a80565b34801561039557600080fd5b506103387f000000000000000000000000000000000000000000000000000000000000000081565b3480156103c957600080fd5b5060015b60405160ff909116815260200161026a565b3480156103eb57600080fd5b506102936103fa3660046149a2565b611041565b34801561040b57600080fd5b506104146110a2565b6040516001600160601b03909116815260200161026a565b34801561043857600080fd5b506103cd600881565b34801561044d57600080fd5b506102936110e7565b34801561046257600080fd5b506103386103e881565b34801561047857600080fd5b50610293610487366004614a93565b6110f9565b34801561049857600080fd5b506104ea6104a7366004614a21565b60126020526000908152604090205463ffffffff811690600160201b810462ffffff1690600160381b810461ffff1690600160481b90046001600160a01b031684565b6040805163ffffffff909516855262ffffff909316602085015261ffff909116918301919091526001600160a01b0316606082015260800161026a565b610338610535366004614c13565b611174565b34801561054657600080fd5b506002546001600160601b0316610338565b34801561056457600080fd5b5061056d611316565b604051901515815260200161026a565b34801561058957600080fd5b506102936105983660046149a2565b611326565b3480156105a957600080fd5b506105b2606481565b60405161ffff909116815260200161026a565b3480156105d157600080fd5b506105e56105e0366004614c93565b611387565b604080516001600160601b03909416845260208401929092529082015260600161026a565b34801561061657600080fd5b50610293610625366004614ce0565b611430565b34801561063657600080fd5b50610414676765c793fa10079d601b1b81565b34801561065557600080fd5b5061029361151b565b34801561066a57600080fd5b506106927f000000000000000000000000000000000000000000000000000000000000000081565b60405161026a9190614d4b565b3480156106ab57600080fd5b506102936115c5565b3480156106c057600080fd5b506106c96115d5565b60405161026a9190614d87565b3480156106e257600080fd5b506000546001600160a01b0316610692565b34801561070057600080fd5b5061029361070f366004614d96565b61163a565b34801561072057600080fd5b5061033861072f366004614dde565b611676565b34801561074057600080fd5b5061033861169e565b34801561075557600080fd5b50610293610764366004614e92565b6118da565b34801561077557600080fd5b50610789610784366004614ee1565b611abb565b60405161026a9190614f60565b3480156107a257600080fd5b506102936107b1366004614f73565b611cca565b3480156107c257600080fd5b506102936107d1366004614a21565b611d3e565b3480156107e257600080fd5b50600c54610692906001600160a01b031681565b34801561080257600080fd5b506102936108113660046149f1565b611e6a565b34801561082257600080fd5b50610293610831366004614fb4565b61201c565b34801561084257600080fd5b506102936108513660046149f1565b6124ea565b34801561086257600080fd5b506103387f000000000000000000000000000000000000000000000000000000000000000081565b34801561089657600080fd5b506102936108a5366004615008565b6127e5565b3480156108b657600080fd5b506004546005546109099160ff8082169261010083049091169163ffffffff620100008204811692600160301b8304821692600160501b8104831692600160701b8204811692600160901b909204169088565b60408051981515895260ff909716602089015263ffffffff9586169688019690965292841660608701529083166080860152821660a08501521660c083015260e08201526101000161026a565b34801561096257600080fd5b506102936109713660046149f1565b6127fe565b610338610984366004615049565b61290a565b34801561099557600080fd5b506109a96109a4366004614a21565b612bdf565b60405161026a9493929190615123565b3480156109c557600080fd5b506106927f000000000000000000000000000000000000000000000000000000000000000081565b3480156109f957600080fd5b50610293610a083660046149a2565b612ccc565b348015610a1957600080fd5b50610293610a2836600461516f565b612cdd565b348015610a3957600080fd5b50610338610a4836600461519b565b612e46565b600c546060906001600160a01b03163314610a7b57604051634bea32db60e11b815260040160405180910390fd5b610a83611316565b15610aa95760405162461bcd60e51b8152600401610aa0906151cb565b60405180910390fd5b600080876001600160401b03811115610ac457610ac4614acf565b604051908082528060200260200182016040528015610afd57816020015b610aea6146a5565b815260200190600190039081610ae25790505b50905060005b88811015610bed5760008a8a83818110610b1f57610b1f6151f5565b9050602002810190610b31919061520b565b610b3a9061539b565b9050610b4781888b612e62565b60408101515151151580610b6357506040810151516020015115155b15610bda576040805160808101825282516001600160401b0316815260208084015162ffffff168183015283830180515151938301939093529151519091015160608201528351849061ffff8716908110610bc057610bc06151f5565b60200260200101819052508380610bd690615488565b9450505b5080610be5816154a9565b915050610b03565b5060008261ffff166001600160401b03811115610c0c57610c0c614acf565b604051908082528060200260200182016040528015610c4557816020015b610c326146a5565b815260200190600190039081610c2a5790505b50905060005b8361ffff16811015610ca157828181518110610c6957610c696151f5565b6020026020010151828281518110610c8357610c836151f5565b60200260200101819052508080610c99906154a9565b915050610c4b565b507ff10ea936d00579b4c52035ee33bf46929646b3aa87554c565d8fb2c7aa549c4486898984604051610cd794939291906154c2565b60405180910390a19998505050505050505050565b610cf4612f89565b610cff600a82612fdc565b610d1e5780604051635428d44960e01b8152600401610aa09190614d4b565b610d29600a82612ff1565b507ff80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af3781604051610d599190614d4b565b60405180910390a150565b600c546001600160a01b03163314610d8f57604051634bea32db60e11b815260040160405180910390fd5b60045460ff1615610dde5760408051808201909152436001600160401b039081168083529083166020909201829052600780546001600160801b031916909117600160401b9092029190911790555b50565b60008281526008602052604090205482906001600160a01b031680610e1c5760405163c5171ee960e01b815260048101839052602401610aa0565b336001600160a01b03821614610e475780604051636c51fda960e11b8152600401610aa09190614d4b565b6000848152600960205260409020548490600160601b90046001600160401b031615610e8657604051631685ecdd60e31b815260040160405180910390fd5b60065460ff1615610eaa5760405163769dd35360e11b815260040160405180910390fd5b6000858152600960209081526040918290208251808401909352546001600160601b038116808452600160601b9091046001600160401b031691830191909152610ef387613006565b600280546001600160601b03169082906000610f0f83856154fd565b92506101000a8154816001600160601b0302191690836001600160601b031602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb88846001600160601b03166040518363ffffffff1660e01b8152600401610f8a92919061551d565b6020604051808303816000875af1158015610fa9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fcd9190615544565b610ffd5760405163cf47918160e01b81526001600160601b03808316600483015283166024820152604401610aa0565b877f3784f77e8e883de95b5d47cd713ced01229fa74d118c0a462224bcb0516d43f1888460405161102f929190615561565b60405180910390a25050505050505050565b611049612f89565b600c546001600160a01b03161561108057600c5460405163ea6d390560e01b8152610aa0916001600160a01b031690600401614d4b565b600c80546001600160a01b0319166001600160a01b0392909216919091179055565b6000806110ad613155565b6004546110c79190600160301b900463ffffffff16615583565b6001600160401b0316905060006110df8260006131e6565b509392505050565b6110ef612f89565b6110f761323d565b565b611101612f89565b606461111360408301602084016155bd565b60ff1611156111355760405163015f505560e31b815260040160405180910390fd5b80600461114282826155e7565b9050507f6dff48329afe669cb5cd0d9af619ec435fe9550e84871839e9df73a1fa34c9bf81604051610d59919061572f565b600061117e611316565b1561119b5760405162461bcd60e51b8152600401610aa0906151cb565b60065460ff16156111bf5760405163769dd35360e11b815260040160405180910390fd5b34156111e057604051630b829bad60e21b8152346004820152602401610aa0565b82600080806111f189338a876132cf565b925092509250600080611204338c613525565b600087815260126020908152604091829020885181548a8401518b8601516060808e015163ffffffff90951666ffffffffffffff1990941693909317600160201b62ffffff9384160217600160381b600160e81b031916600160381b61ffff92831602600160481b600160e81b03191617600160481b6001600160a01b03909516949094029390931790935584513381526001600160401b038b1694810194909452918c169383019390935281018f9052908d16608082015260a081018390526001600160601b03821660c0820152919350915085907fb7933fba96b6b452eb44f99fdc08052a45dff82363d59abaff0456931c3d24599060e00160405180910390a250929998505050505050505050565b600154600160a01b900460ff1690565b61132e612f89565b611339600a82613606565b611358578060405163ac8a27ef60e01b8152600401610aa09190614d4b565b7fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af0162581604051610d599190614d4b565b600080600080611395613155565b6001600160401b0316905060006113ad8260106157de565b9050600060146113be8360156157de565b6113c8919061580b565b87516113d491906157de565b60045484906113f090600160501b900463ffffffff168b61581f565b6001600160601b031661140391906157de565b61140d919061583f565b905060008061141d8360006131e6565b9098509650939450505050509250925092565b600c546001600160a01b0316331461145b57604051634bea32db60e11b815260040160405180910390fd5b828015806114695750601f81115b1561148a57604051634ecc4fef60e01b815260048101829052602401610aa0565b8082146114ae5760405163339f8a9d60e01b8152610aa09082908490600401615852565b60005b81811015611513576115018686838181106114ce576114ce6151f5565b90506020020160208101906114e391906149a2565b8585848181106114f5576114f56151f5565b90506020020135612cdd565b8061150b816154a9565b9150506114b1565b505050505050565b6001546001600160a01b0316331461156e5760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b6044820152606401610aa0565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6115cd612f89565b6110f761361b565b6115dd6146db565b6040805161010081019182905290601190600890826000855b82829054906101000a900462ffffff1662ffffff16815260200190600301906020826002010492830192600103820291508084116115f65790505050505050905090565b600c546001600160a01b0316331461166557604051634bea32db60e11b815260040160405180910390fd5b61167260118260086146fa565b5050565b6000806116898563ffffffff1685611387565b50506001600160601b03169695505050505050565b60065460009060ff16156116c55760405163769dd35360e11b815260040160405180910390fd5b6116cd611316565b156116ea5760405162461bcd60e51b8152600401610aa0906151cb565b6000336116f8600143615860565b6001546040516001600160601b0319606094851b81166020830152924060348201523090931b90911660548301526001600160c01b0319600160a81b90910460c01b16606882015260700160408051808303601f19018152919052805160209091012060018054919250600160a81b9091046001600160401b031690601561177f83615873565b91906101000a8154816001600160401b0302191690836001600160401b03160217905550506000806001600160401b038111156117be576117be614acf565b6040519080825280602002602001820160405280156117e7578160200160208202803683370190505b5060408051808201825260008082526020808301828152878352600982528483209351845491516001600160601b039091166001600160a01b031992831617600160601b6001600160401b039092169190910217909355835160608101855233815280820183815281860187815289855260088452959093208151815486166001600160a01b03918216178255935160018201805490961694169390931790935592518051949550919390926118a4926002850192910190614798565b505060405133915083907f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d90600090a350905090565b60065460ff16156118fe5760405163769dd35360e11b815260040160405180910390fd5b611906611316565b156119235760405162461bcd60e51b8152600401610aa0906151cb565b336001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461196c576040516344b0e3c360e01b815260040160405180910390fd5b6020811461199257604051636865567560e01b8152610aa0906020908390600401615897565b60006119a082840184614a21565b6000818152600860205260409020549091506001600160a01b03166119db5760405163c5171ee960e01b815260048101829052602401610aa0565b600081815260096020526040812080546001600160601b031691869190611a02838561581f565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600260008282829054906101000a90046001600160601b0316611a4a919061581f565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a828784611a9d919061583f565b604051611aab929190615852565b60405180910390a2505050505050565b60065460609060ff1615611ae25760405163769dd35360e11b815260040160405180910390fd5b60008381526012602081815260408084208151608081018352815463ffffffff8116825262ffffff600160201b8204168286015261ffff600160381b820416938201939093526001600160a01b03600160481b8404811660608301908152968a9052949093526001600160e81b031990911690559151859291163314611b8357806060015133604051638e30e82360e01b8152600401610aa09291906158ab565b8051600090611bb9907f00000000000000000000000000000000000000000000000000000000000000009063ffffffff166157de565b90506000611bc561367b565b90506000836020015162ffffff1682611bde9190615860565b9050808310611c235782846020015162ffffff1684611bfd919061583f565b611c0890600161583f565b6040516315ad27c360e01b8152600401610aa0929190615852565b6001600160401b03831115611c4e576040516302c6ef8160e11b815260048101849052602401610aa0565b604051898152339086907f16f3f633197fafab10a5df69e6f3f2f7f20092f08d8d47de0a91c0f4b96a1a259060200160405180910390a3611cbb8585600f6000611ca6888a6020015162ffffff1660189190911b1790565b81526020019081526020016000205486613705565b955050505050505b9392505050565b611cd2612f89565b8051600d80546020808501805163ffffffff908116600160201b026001600160401b031990941695811695861793909317909355604080519485529251909116908301527ffffe83c0c6d543712480c43dcc77636fc5671d01e3199dfb237b1bbf29d971e49101610d59565b60065460ff1615611d625760405163769dd35360e11b815260040160405180910390fd5b6000818152600860205260409020546001600160a01b0316611d9a5760405163c5171ee960e01b815260048101829052602401610aa0565b6000818152600860205260409020600101546001600160a01b03163314611df1576000818152600860205260409081902060010154905163d084e97560e01b8152610aa0916001600160a01b031690600401614d4b565b6000818152600860205260409081902080546001600160a01b031980821633908117845560019093018054909116905591516001600160a01b039092169183917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c938691611e5e9185916158ab565b60405180910390a25050565b60008281526008602052604090205482906001600160a01b031680611ea55760405163c5171ee960e01b815260048101839052602401610aa0565b336001600160a01b03821614611ed05780604051636c51fda960e11b8152600401610aa09190614d4b565b60065460ff1615611ef45760405163769dd35360e11b815260040160405180910390fd5b611efc611316565b15611f195760405162461bcd60e51b8152600401610aa0906151cb565b60008481526008602052604090206002015460631901611f4c576040516305a48e0f60e01b815260040160405180910390fd5b60036000611f5a85876138be565b815260208101919091526040016000205460ff1661201657600160036000611f8286886138be565b815260208082019290925260409081016000908120805460ff191694151594909417909355868352600882528083206002018054600181018255908452919092200180546001600160a01b0319166001600160a01b0386161790555184907f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e19061200d908690614d4b565b60405180910390a25b50505050565b60065460ff16156120405760405163769dd35360e11b815260040160405180910390fd5b61204b600a84612fdc565b61206a5782604051635428d44960e01b8152600401610aa09190614d4b565b6040811461208f5760408051636865567560e01b8152610aa091908390600401615897565b600061209d828401846158c5565b90506000806000806120b28560200151612bdf565b9350935093509350816001600160a01b0316336001600160a01b0316146120ee5781604051636c51fda960e11b8152600401610aa09190614d4b565b876001600160a01b031663294daa496040518163ffffffff1660e01b8152600401602060405180830381865afa15801561212c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061215091906158ff565b60ff16856000015160ff16146121ed578460000151886001600160a01b031663294daa496040518163ffffffff1660e01b8152600401602060405180830381865afa1580156121a3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121c791906158ff565b60405163e7aada9560e01b815260ff928316600482015291166024820152604401610aa0565b6001600160401b0383161561221557604051631685ecdd60e31b815260040160405180910390fd5b60006040518060a0016040528061222a600190565b60ff16815260200187602001518152602001846001600160a01b03168152602001838152602001866001600160601b03168152509050600081604051602001612273919061591c565b60405160208183030381529060405290506122918760200151613006565b600280548791906000906122af9084906001600160601b03166154fd565b92506101000a8154816001600160601b0302191690836001600160601b031602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb8b886040518363ffffffff1660e01b8152600401612321929190615561565b6020604051808303816000875af1158015612340573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123649190615544565b6123a55760405162461bcd60e51b8152602060048201526012602482015271696e73756666696369656e742066756e647360701b6044820152606401610aa0565b60405163ce3f471960e01b81526001600160a01b038b169063ce3f4719906123d1908490600401614a80565b600060405180830381600087803b1580156123eb57600080fd5b505af11580156123ff573d6000803e3d6000fd5b5050505060005b835181101561249957838181518110612421576124216151f5565b60200260200101516001600160a01b0316638ea981178c6040518263ffffffff1660e01b81526004016124549190614d4b565b600060405180830381600087803b15801561246e57600080fd5b505af1158015612482573d6000803e3d6000fd5b505050508080612491906154a9565b915050612406565b508660200151876000015160ff167fbd89b747474d3fc04664dfbd1d56ae7ffbe46ee097cdb9979c13916bb76269ce8c6040516124d69190614d4b565b60405180910390a350505050505050505050565b60008281526008602052604090205482906001600160a01b0316806125255760405163c5171ee960e01b815260048101839052602401610aa0565b336001600160a01b038216146125505780604051636c51fda960e11b8152600401610aa09190614d4b565b6000848152600960205260409020548490600160601b90046001600160401b03161561258f57604051631685ecdd60e31b815260040160405180910390fd5b60065460ff16156125b35760405163769dd35360e11b815260040160405180910390fd5b600360006125c186886138be565b815260208101919091526040016000205460ff166125f65784846040516379bfd40160e01b8152600401610aa0929190615982565b60008581526008602090815260408083206002018054825181850281018501909352808352919290919083018282801561265957602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161263b575b505050505090506000600182516126709190615860565b905060005b825181101561277b57866001600160a01b031683828151811061269a5761269a6151f5565b60200260200101516001600160a01b0316036127695760008383815181106126c4576126c46151f5565b6020026020010151905080600860008b815260200190815260200160002060020183815481106126f6576126f66151f5565b600091825260208083209190910180546001600160a01b0319166001600160a01b0394909416939093179092558a815260089091526040902060020180548061274157612741615999565b600082815260209020810160001990810180546001600160a01b03191690550190555061277b565b80612773816154a9565b915050612675565b506003600061278a888a6138be565b815260208101919091526040908101600020805460ff191690555187907f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a7906127d4908990614d4b565b60405180910390a250505050505050565b604051632cb6686f60e01b815260040160405180910390fd5b60008281526008602052604090205482906001600160a01b0316806128395760405163c5171ee960e01b815260048101839052602401610aa0565b336001600160a01b038216146128645780604051636c51fda960e11b8152600401610aa09190614d4b565b60065460ff16156128885760405163769dd35360e11b815260040160405180910390fd5b6000848152600860205260409020600101546001600160a01b03848116911614612016576000848152600860205260409081902060010180546001600160a01b0319166001600160a01b0386161790555184907f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a19061200d90339087906158ab565b6000612914611316565b156129315760405162461bcd60e51b8152600401610aa0906151cb565b60065460ff16156129555760405163769dd35360e11b815260040160405180910390fd5b341561297657604051630b829bad60e21b8152346004820152602401610aa0565b60408051808201909152600d5463ffffffff808216808452600160201b90920481166020840152861611156129c2578051604051637aebf00f60e11b8152610aa09187916004016149da565b806020015163ffffffff16845111156129f75783518160200151604051631961a6f960e31b8152600401610aa09291906149da565b50600080612a0789338a8a6132cf565b925050915060006040518061010001604052808481526020018a61ffff168152602001336001600160a01b031681526020018781526020018863ffffffff166001600160601b031681526020018b81526020016000815260200160008152509050600080612a74836138d4565b60c087019190915260e08601919091526040519193509150612aa09085908c908f908790602001615a2c565b60405160208183030381529060405280519060200120600e6000878152602001908152602001600020819055506000604051806101600160405280878152602001336001600160a01b03168152602001866001600160401b031681526020018c62ffffff1681526020018e81526020018d61ffff1681526020018b63ffffffff1681526020018581526020018a8152602001848152602001836001600160601b0316815250905080600001517f01872fb9c7d6d68af06a17347935e04412da302a377224c205e672c26e18c37f82602001518360400151846060015185608001518660a001518760c001518860e0015160c001518960e0015160e001518a61010001518b61012001518c6101400151604051612bc69b9a99989796959493929190615a5e565b60405180910390a250939b9a5050505050505050505050565b600081815260086020526040812054819081906060906001600160a01b0316612c1e5760405163c5171ee960e01b815260048101869052602401610aa0565b60008581526009602090815260408083205460088352928190208054600290910180548351818602810186019094528084526001600160601b03861695600160601b90046001600160401b0316946001600160a01b03909316939192839190830182828015612cb657602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612c98575b5050505050905093509350935093509193509193565b612cd4612f89565b610dde81613a14565b600c546001600160a01b03163314612d0857604051634bea32db60e11b815260040160405180910390fd5b60405163a9059cbb60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90612d56908590859060040161551d565b6020604051808303816000875af1158015612d75573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d999190615544565b611672576040516370a0823160e01b81526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906370a0823190612de9903090600401614d4b565b602060405180830381865afa158015612e06573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e2a9190615af4565b8160405163cf47918160e01b8152600401610aa0929190615852565b6000612e506110a2565b6001600160601b031690505b92915050565b82516001600160401b0380841691161115612ea657825160405163012d824d60e01b81526001600160401b0380851660048301529091166024820152604401610aa0565b60408301515151600090158015612ec4575060408401515160200151155b15612f0857600f6000612ef286600001516001600160401b0316876020015162ffffff1660189190911b1790565b8152602001908152602001600020549050612f71565b8360400151604051602001612f1d9190615b0d565b60405160208183030381529060405280519060200120905080600f6000612f5f87600001516001600160401b0316886020015162ffffff1660189190911b1790565b81526020810191909152604001600020555b606084015151612f82818387613ab7565b5050505050565b6000546001600160a01b031633146110f75760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610aa0565b6000611cc3836001600160a01b038416613e71565b6000611cc3836001600160a01b038416613e89565b6000818152600860209081526040808320815160608101835281546001600160a01b0390811682526001830154168185015260028201805484518187028101870186528181529295939486019383018282801561308c57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161306e575b505050505081525050905060005b8160400151518110156130fc57600360006130d2846040015184815181106130c4576130c46151f5565b6020026020010151866138be565b81526020810191909152604001600020805460ff19169055806130f4816154a9565b91505061309a565b50600082815260086020526040812080546001600160a01b0319908116825560018201805490911690559061313460028301826147ed565b505050600090815260096020526040902080546001600160a01b0319169055565b60045460009060ff16801561317b5750600754600160401b90046001600160401b031615155b156131e157600454600160901b900463ffffffff16431080806131c257506004546131b390600160901b900463ffffffff1643615860565b6007546001600160401b031610155b156131df575050600754600160401b90046001600160401b031690565b505b503a90565b6004546000908190819060649061320a90600160701b900463ffffffff1682615b40565b61321a9063ffffffff16876157de565b613224919061580b565b90506132308185613f83565b92509250505b9250929050565b613245611316565b6132885760405162461bcd60e51b815260206004820152601460248201527314185d5cd8589b194e881b9bdd081c185d5cd95960621b6044820152606401610aa0565b6001805460ff60a01b191690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b6040516132c59190614d4b565b60405180910390a1565b604080516080810182526000808252602082018190529181018290526060810182905260006103e88561ffff16111561332157846103e8604051634a90778560e01b8152600401610aa0929190615897565b8461ffff16600003613346576040516308fad2a760e01b815260040160405180910390fd5b600080613351613ff3565b601054919350915065ffffffffffff1660006133bc8b8b84604080513060208201529081018490526001600160a01b038316606082015265ffffffffffff8216608082015260009060a00160408051601f198184030181529190528051602090910120949350505050565b6040805180820182526010805465ffffffffffff16825282516101008101938490529394506000939192909160208401916011906008908288855b82829054906101000a900462ffffff1662ffffff16815260200190600301906020826002010492830192600103820291508084116133f757905050505050508152505090508260016134499190615b5d565b6010805465ffffffffffff191665ffffffffffff9290921691909117905560005b60088110156134b0578962ffffff168260200151826008811061348f5761348f6151f5565b602002015162ffffff16146134b057806134a8816154a9565b91505061346a565b600881106134d8576020820151604051630c4f769b60e41b8152610aa0918c91600401615b7c565b50506040805160808101825263ffffffff909416845262ffffff8916602085015261ffff8a16908401526001600160a01b038a1660608401529550909350909150505b9450945094915050565b6000806003600061353686866138be565b815260208101919091526040016000205460ff1661356b5782846040516379bfd40160e01b8152600401610aa0929190615982565b60006135756110a2565b600085815260096020526040902080546001600160601b039283169350909116828110156135c457815460405163cf47918160e01b8152610aa0916001600160601b0316908590600401615b96565b81546001600160601b0319908116918490036001600160601b038181169390931790935560028054918216918316859003909216179055909590945092505050565b6000611cc3836001600160a01b0384166140c6565b613623611316565b156136405760405162461bcd60e51b8152600401610aa0906151cb565b6001805460ff60a01b1916600160a01b1790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586132b83390565b60004661a4b1811480613690575062066eed81145b156136fe5760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156136d4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136f89190615af4565b91505090565b4391505090565b6060826137375760405163220a34e960e11b8152600481018690526001600160401b0383166024820152604401610aa0565b604080516020808201889052865163ffffffff168284015286015162ffffff166060808301919091529186015161ffff166080820152908501516001600160a01b031660a082015260c0810184905260009060e0016040516020818303038152906040528051906020012090506103e8856040015161ffff1611156137d95784604001516103e8604051634a90778560e01b8152600401610aa0929190615897565b6000856040015161ffff166001600160401b038111156137fb576137fb614acf565b604051908082528060200260200182016040528015613824578160200160208202803683370190505b50905060005b866040015161ffff168161ffff1610156138b357828160405160200161386792919091825260f01b6001600160f01b031916602082015260220190565b6040516020818303038152906040528051906020012060001c828261ffff1681518110613896576138966151f5565b6020908102919091010152806138ab81615488565b91505061382a565b509695505050505050565b60a081901b6001600160a01b0383161792915050565b600080600080600360006138f087604001518860a001516138be565b815260208101919091526040016000205460ff1661392d578460a0015185604001516040516379bfd40160e01b8152600401610aa0929190615982565b600080600061394488608001518960600151611387565b60a08b0151600090815260096020526040902080546001600160601b039485169750929550909350911684111561399c57805460405163cf47918160e01b8152610aa0916001600160601b0316908690600401615b96565b80546001600160601b0360016001600160401b03600160601b8085048216929092011602818116828416178790038083166001600160601b03199283166001600160a01b0319909516949094179390931790935560028054808316889003909216919093161790915592989197509550909350915050565b336001600160a01b03821603613a665760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610aa0565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000836001600160401b03811115613ad157613ad1614acf565b604051908082528060200260200182016040528015613afa578160200160208202803683370190505b5090506000846001600160401b03811115613b1757613b17614acf565b6040519080825280601f01601f191660200182016040528015613b41576020820181803683370190505b5090506000856001600160401b03811115613b5e57613b5e614acf565b604051908082528060200260200182016040528015613b9157816020015b6060815260200190600190039081613b7c5790505b509050600080876001600160401b03811115613baf57613baf614acf565b604051908082528060200260200182016040528015613bd8578160200160208202803683370190505b5090506000886001600160401b03811115613bf557613bf5614acf565b604051908082528060200260200182016040528015613c1e578160200160208202803683370190505b50905060005b89811015613d6b57600088606001518281518110613c4457613c446151f5565b602002602001015190506000806000613c678c600001518d602001518f87614110565b9250925092508215613ca85781898961ffff1681518110613c8a57613c8a6151f5565b60200260200101819052508780613ca090615488565b985050613cd7565b600160f81b8a8681518110613cbf57613cbf6151f5565b60200101906001600160f81b031916908160001a9053505b8351518b518c9087908110613cee57613cee6151f5565b60200260200101818152505080878681518110613d0d57613d0d6151f5565b60200260200101906001600160601b031690816001600160601b031681525050836000015160a00151868681518110613d4857613d486151f5565b602002602001018181525050505050508080613d63906154a9565b915050613c24565b5060608701515115613e665760008361ffff166001600160401b03811115613d9557613d95614acf565b604051908082528060200260200182016040528015613dc857816020015b6060815260200190600190039081613db35790505b50905060005b8461ffff16811015613e2457858181518110613dec57613dec6151f5565b6020026020010151828281518110613e0657613e066151f5565b60200260200101819052508080613e1c906154a9565b915050613dce565b507f8f79f730779e875ce76c428039cc2052b5b5918c2a55c598fab251c1198aec548787838686604051613e5c959493929190615be8565b60405180910390a1505b505050505050505050565b60009081526001919091016020526040902054151590565b60008181526001830160205260408120548015613f72576000613ead600183615860565b8554909150600090613ec190600190615860565b9050818114613f26576000866000018281548110613ee157613ee16151f5565b9060005260206000200154905080876000018481548110613f0457613f046151f5565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080613f3757613f37615999565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050612e5c565b6000915050612e5c565b5092915050565b600080808315613f935783613f9b565b613f9b614465565b9050600081613fb287670de0b6b3a76400006157de565b613fbc919061580b565b9050676765c793fa10079d601b1b811115613fe95760405162de437160e81b815260040160405180910390fd5b9590945092505050565b600080600061400061367b565b9050600061402e7f000000000000000000000000000000000000000000000000000000000000000083615cb2565b905060008161405d7f00000000000000000000000000000000000000000000000000000000000000008561583f565b6140679190615860565b905060006140957f00000000000000000000000000000000000000000000000000000000000000008361580b565b905063ffffffff81106140bb576040516307b2a52360e41b815260040160405180910390fd5b909590945092505050565b60006140d28383613e71565b61410857508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155612e5c565b506000612e5c565b805160a00151600090815260096020526040812060609082908161415d7f00000000000000000000000000000000000000000000000000000000000000006001600160401b038b1661580b565b865160a08101516040519293509091600091614181918d918d918690602001615a2c565b60408051601f19818403018152918152815160209283012084516000908152600e90935291205490915081146141f4575050905460408051808201909152601081526f756e6b6e6f776e2063616c6c6261636b60801b60208201526001955093506001600160601b0316915061351b9050565b506040805160808101825263ffffffff8416815262ffffff8b1660208083019190915283015161ffff1681830152908201516001600160a01b0316606082015288614282575050604080518082019091526016815275756e617661696c61626c652072616e646f6d6e65737360501b6020820152915460019550919350506001600160601b0316905061351b565b60006142948360000151838c8f613705565b606080840151855191860151604051939450909260009263d21ea8fd60e01b926142c392879190602401615cc6565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091526006805460ff1916600117905590506000805a9050600061432c8e60000151608001516001600160601b031689604001518661453d565b909350905080614361578d516080015160405163aad1598360e01b8152610aa0916001600160601b0316908490600401615852565b506000610bb85a614372919061583f565b6006805460ff1916905590508181101561439a5761439a6143938284615860565b8f5161457c565b8954600160601b90046001600160401b03168a600c6143b883615cf1565b82546001600160401b039182166101009390930a92830291909202199091161790555087516000908152600e60205260408120558261442c5760408051808201909152601081526f195e1958dd5d1a5bdb8819985a5b195960821b60208201528a54600191906001600160601b031661444b565b604080516020810190915260008082528b549091906001600160601b03165b9c509c509c50505050505050505050509450945094915050565b6004805460408051633fabe5a360e21b815290516000936201000090930463ffffffff169283151592859283927f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169263feaf968c928183019260a0928290030181865afa1580156144e3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906145079190615d2b565b50945090925084915050801561452b57506145228242615860565b8463ffffffff16105b1561453557506005545b949350505050565b6000805a610bb8811061457357610bb881039050856040820482031115614573576000808551602087016000898bf19250600191505b50935093915050565b80608001516001600160601b0316821115614595575050565b6004546000906064906145b090610100900460ff1682615d7b565b60ff168360c001518585608001516001600160601b03166145d19190615860565b6145db91906157de565b6145e591906157de565b6145ef919061580b565b90506000614601828460e00151613f83565b5060a08401516000908152600960205260408120805492935083929091906146339084906001600160601b031661581f565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600260008282829054906101000a90046001600160601b031661467b919061581f565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555050505050565b604051806080016040528060006001600160401b03168152602001600062ffffff16815260200160008152602001600081525090565b6040518061010001604052806008906020820280368337509192915050565b6001830191839082156147885791602002820160005b8382111561475757833562ffffff1683826101000a81548162ffffff021916908362ffffff1602179055509260200192600301602081600201049283019260010302614710565b80156147865782816101000a81549062ffffff0219169055600301602081600201049283019260010302614757565b505b50614794929150614807565b5090565b828054828255906000526020600020908101928215614788579160200282015b8281111561478857825182546001600160a01b0319166001600160a01b039091161782556020909201916001909101906147b8565b5080546000825590600052602060002090810190610dde91905b5b808211156147945760008155600101614808565b60008083601f84011261482e57600080fd5b5081356001600160401b0381111561484557600080fd5b6020830191508360208260051b850101111561323657600080fd5b80356001600160401b038116811461487757600080fd5b919050565b60008060008060008060a0878903121561489557600080fd5b86356001600160401b038111156148ab57600080fd5b6148b789828a0161481c565b90975095505060208701356001600160c01b03811681146148d757600080fd5b93506148e560408801614860565b92506148f360608801614860565b9150608087013590509295509295509295565b600081518084526020808501945080840160005b8381101561496457815180516001600160401b031688528381015162ffffff168489015260408082015190890152606090810151908801526080909601959082019060010161491a565b509495945050505050565b602081526000611cc36020830184614906565b6001600160a01b0381168114610dde57600080fd5b803561487781614982565b6000602082840312156149b457600080fd5b8135611cc381614982565b6000602082840312156149d157600080fd5b611cc382614860565b63ffffffff92831681529116602082015260400190565b60008060408385031215614a0457600080fd5b823591506020830135614a1681614982565b809150509250929050565b600060208284031215614a3357600080fd5b5035919050565b6000815180845260005b81811015614a6057602081850181015186830182015201614a44565b506000602082860101526020601f19601f83011685010191505092915050565b602081526000611cc36020830184614a3a565b60006101008284031215614aa657600080fd5b50919050565b803561ffff8116811461487757600080fd5b62ffffff81168114610dde57600080fd5b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715614b0757614b07614acf565b60405290565b60405161010081016001600160401b0381118282101715614b0757614b07614acf565b604051608081016001600160401b0381118282101715614b0757614b07614acf565b604051602081016001600160401b0381118282101715614b0757614b07614acf565b604051601f8201601f191681016001600160401b0381118282101715614b9c57614b9c614acf565b604052919050565b600082601f830112614bb557600080fd5b81356001600160401b03811115614bce57614bce614acf565b614be1601f8201601f1916602001614b74565b818152846020838601011115614bf657600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060808587031215614c2957600080fd5b84359350614c3960208601614aac565b92506040850135614c4981614abe565b915060608501356001600160401b03811115614c6457600080fd5b614c7087828801614ba4565b91505092959194509250565b80356001600160601b038116811461487757600080fd5b60008060408385031215614ca657600080fd5b614caf83614c7c565b915060208301356001600160401b03811115614cca57600080fd5b614cd685828601614ba4565b9150509250929050565b60008060008060408587031215614cf657600080fd5b84356001600160401b0380821115614d0d57600080fd5b614d198883890161481c565b90965094506020870135915080821115614d3257600080fd5b50614d3f8782880161481c565b95989497509550505050565b6001600160a01b0391909116815260200190565b8060005b600881101561201657815162ffffff16845260209384019390910190600101614d63565b6101008101612e5c8284614d5f565b6000610100808385031215614daa57600080fd5b838184011115614db957600080fd5b509092915050565b63ffffffff81168114610dde57600080fd5b803561487781614dc1565b60008060008060808587031215614df457600080fd5b843593506020850135614e0681614dc1565b925060408501356001600160401b0380821115614e2257600080fd5b614e2e88838901614ba4565b93506060870135915080821115614e4457600080fd5b50614c7087828801614ba4565b60008083601f840112614e6357600080fd5b5081356001600160401b03811115614e7a57600080fd5b60208301915083602082850101111561323657600080fd5b60008060008060608587031215614ea857600080fd5b8435614eb381614982565b93506020850135925060408501356001600160401b03811115614ed557600080fd5b614d3f87828801614e51565b600080600060608486031215614ef657600080fd5b833592506020840135915060408401356001600160401b03811115614f1a57600080fd5b614f2686828701614ba4565b9150509250925092565b600081518084526020808501945080840160005b8381101561496457815187529582019590820190600101614f44565b602081526000611cc36020830184614f30565b600060408284031215614f8557600080fd5b614f8d614ae5565b8235614f9881614dc1565b81526020830135614fa881614dc1565b60208201529392505050565b600080600060408486031215614fc957600080fd5b8335614fd481614982565b925060208401356001600160401b03811115614fef57600080fd5b614ffb86828701614e51565b9497909650939450505050565b6000806020838503121561501b57600080fd5b82356001600160401b0381111561503157600080fd5b61503d85828601614e51565b90969095509350505050565b60008060008060008060c0878903121561506257600080fd5b8635955061507260208801614aac565b9450604087013561508281614abe565b9350606087013561509281614dc1565b925060808701356001600160401b03808211156150ae57600080fd5b6150ba8a838b01614ba4565b935060a08901359150808211156150d057600080fd5b506150dd89828a01614ba4565b9150509295509295509295565b600081518084526020808501945080840160005b838110156149645781516001600160a01b0316875295820195908201906001016150fe565b6001600160601b03851681526001600160401b03841660208201526001600160a01b0383166040820152608060608201819052600090615165908301846150ea565b9695505050505050565b6000806040838503121561518257600080fd5b823561518d81614982565b946020939093013593505050565b600080604083850312156151ae57600080fd5b8235915060208301356001600160401b03811115614cca57600080fd5b60208082526010908201526f14185d5cd8589b194e881c185d5cd95960821b604082015260600190565b634e487b7160e01b600052603260045260246000fd5b60008235609e1983360301811261522157600080fd5b9190910192915050565b600082601f83011261523c57600080fd5b813560206001600160401b038083111561525857615258614acf565b8260051b615267838201614b74565b938452858101830193838101908886111561528157600080fd5b84880192505b8583101561538f5782358481111561529e57600080fd5b8801601f196040828c03820112156152b557600080fd5b6152bd614ae5565b87830135878111156152ce57600080fd5b8301610100818e03840112156152e357600080fd5b6152eb614b0d565b92508881013583526152ff60408201614aac565b8984015261530f60608201614997565b604084015260808101358881111561532657600080fd5b6153348e8b83850101614ba4565b60608501525061534660a08201614c7c565b608084015260c081013560a084015260e081013560c084015261010081013560e08401525081815261537a60408401614c7c565b81890152845250509184019190840190615287565b98975050505050505050565b600081360360a08112156153ae57600080fd5b6153b6614b30565b6153bf84614860565b81526020808501356153d081614abe565b828201526040603f19840112156153e657600080fd5b6153ee614b52565b925036605f8601126153ff57600080fd5b615407614ae5565b80608087013681111561541957600080fd5b604088015b81811015615435578035845292840192840161541e565b50908552604084019490945250509035906001600160401b0382111561545a57600080fd5b6154663683860161522b565b60608201529392505050565b634e487b7160e01b600052601160045260246000fd5b600061ffff80831681810361549f5761549f615472565b6001019392505050565b6000600182016154bb576154bb615472565b5060010190565b6001600160401b0385811682526001600160c01b03851660208301528316604082015260806060820181905260009061516590830184614906565b6001600160601b03828116828216039080821115613f7c57613f7c615472565b6001600160a01b03929092168252602082015260400190565b8015158114610dde57600080fd5b60006020828403121561555657600080fd5b8151611cc381615536565b6001600160a01b039290921682526001600160601b0316602082015260400190565b6001600160401b038181168382160280821691908281146155a6576155a6615472565b505092915050565b60ff81168114610dde57600080fd5b6000602082840312156155cf57600080fd5b8135611cc3816155ae565b60008135612e5c81614dc1565b81356155f281615536565b815460ff19811691151560ff1691821783556020840135615612816155ae565b61ff008160081b168361ffff19841617178455505050615655615637604084016155da565b825465ffffffff0000191660109190911b65ffffffff000016178255565b615688615664606084016155da565b82805463ffffffff60301b191660309290921b63ffffffff60301b16919091179055565b6156bb615697608084016155da565b82805463ffffffff60501b191660509290921b63ffffffff60501b16919091179055565b6156ee6156ca60a084016155da565b82805463ffffffff60701b191660709290921b63ffffffff60701b16919091179055565b6157216156fd60c084016155da565b82805463ffffffff60901b191660909290921b63ffffffff60901b16919091179055565b60e082013560018201555050565b6101008101823561573f81615536565b151582526020830135615751816155ae565b60ff166020830152604083013561576781614dc1565b63ffffffff16604083015261577e60608401614dd3565b63ffffffff16606083015261579560808401614dd3565b63ffffffff1660808301526157ac60a08401614dd3565b63ffffffff1660a08301526157c360c08401614dd3565b63ffffffff811660c08401525060e092830135919092015290565b8082028115828204841417612e5c57612e5c615472565b634e487b7160e01b600052601260045260246000fd5b60008261581a5761581a6157f5565b500490565b6001600160601b03818116838216019080821115613f7c57613f7c615472565b80820180821115612e5c57612e5c615472565b918252602082015260400190565b81810381811115612e5c57612e5c615472565b60006001600160401b038281166002600160401b0319810161549f5761549f615472565b61ffff929092168252602082015260400190565b6001600160a01b0392831681529116602082015260400190565b6000604082840312156158d757600080fd5b6158df614ae5565b82356158ea816155ae565b81526020928301359281019290925250919050565b60006020828403121561591157600080fd5b8151611cc3816155ae565b6020815260ff82511660208201526020820151604082015260018060a01b0360408301511660608201526000606083015160a0608084015261596160c08401826150ea565b608094909401516001600160601b031660a093909301929092525090919050565b9182526001600160a01b0316602082015260400190565b634e487b7160e01b600052603160045260246000fd5b60006101008251845261ffff602084015116602085015260018060a01b03604084015116604085015260608301518160608601526159ef82860182614a3a565b6080858101516001600160601b03169087015260a0808601519087015260c0808601519087015260e0948501519490950193909352509192915050565b60018060401b038516815262ffffff8416602082015282604082015260806060820152600061516560808301846159af565b6001600160a01b038c1681526001600160401b038b16602082015262ffffff8a1660408201526060810189905261ffff8816608082015263ffffffff871660a082015260c0810186905260e081018590526101606101008201819052600090615ac983820187614a3a565b61012084019590955250506001600160601b0391909116610140909101529998505050505050505050565b600060208284031215615b0657600080fd5b5051919050565b815160408201908260005b6002811015615b37578251825260209283019290910190600101615b18565b50505092915050565b63ffffffff818116838216019080821115613f7c57613f7c615472565b65ffffffffffff818116838216019080821115613f7c57613f7c615472565b62ffffff831681526101208101611cc36020830184614d5f565b6001600160601b03929092168252602082015260400190565b600081518084526020808501945080840160005b838110156149645781516001600160601b031687529582019590820190600101615bc3565b60a0808252865190820181905260009060209060c0840190828a01845b82811015615c2157815184529284019290840190600101615c05565b50505083810382850152615c358189614a3a565b905083810360408501528087518083528383019150838160051b840101848a0160005b83811015615c8657601f19868403018552615c74838351614a3a565b94870194925090860190600101615c58565b50508681036060880152615c9a818a615baf565b945050505050828103608084015261538f8185614f30565b600082615cc157615cc16157f5565b500690565b838152606060208201526000615cdf6060830185614f30565b82810360408401526151658185614a3a565b60006001600160401b03821680615d0a57615d0a615472565b6000190192915050565b80516001600160501b038116811461487757600080fd5b600080600080600060a08688031215615d4357600080fd5b615d4c86615d14565b9450602086015193506040860151925060608601519150615d6f60808701615d14565b90509295509295909350565b60ff8281168282160390811115612e5c57612e5c61547256fea164736f6c6343000813000a",
}

var VRFCoordinatorABI = VRFCoordinatorMetaData.ABI

var VRFCoordinatorBin = VRFCoordinatorMetaData.Bin

func DeployVRFCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, beaconPeriodBlocksArg *big.Int, linkToken common.Address, linkEthFeed common.Address) (common.Address, *types.Transaction, *VRFCoordinator, error) {
	parsed, err := VRFCoordinatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorBin), backend, beaconPeriodBlocksArg, linkToken, linkEthFeed)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinator{VRFCoordinatorCaller: VRFCoordinatorCaller{contract: contract}, VRFCoordinatorTransactor: VRFCoordinatorTransactor{contract: contract}, VRFCoordinatorFilterer: VRFCoordinatorFilterer{contract: contract}}, nil
}

type VRFCoordinator struct {
	address common.Address
	abi     abi.ABI
	VRFCoordinatorCaller
	VRFCoordinatorTransactor
	VRFCoordinatorFilterer
}

type VRFCoordinatorCaller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorTransactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorFilterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorSession struct {
	Contract     *VRFCoordinator
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorCallerSession struct {
	Contract *VRFCoordinatorCaller
	CallOpts bind.CallOpts
}

type VRFCoordinatorTransactorSession struct {
	Contract     *VRFCoordinatorTransactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorRaw struct {
	Contract *VRFCoordinator
}

type VRFCoordinatorCallerRaw struct {
	Contract *VRFCoordinatorCaller
}

type VRFCoordinatorTransactorRaw struct {
	Contract *VRFCoordinatorTransactor
}

func NewVRFCoordinator(address common.Address, backend bind.ContractBackend) (*VRFCoordinator, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFCoordinator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinator{address: address, abi: abi, VRFCoordinatorCaller: VRFCoordinatorCaller{contract: contract}, VRFCoordinatorTransactor: VRFCoordinatorTransactor{contract: contract}, VRFCoordinatorFilterer: VRFCoordinatorFilterer{contract: contract}}, nil
}

func NewVRFCoordinatorCaller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorCaller, error) {
	contract, err := bindVRFCoordinator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorCaller{contract: contract}, nil
}

func NewVRFCoordinatorTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorTransactor, error) {
	contract, err := bindVRFCoordinator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTransactor{contract: contract}, nil
}

func NewVRFCoordinatorFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorFilterer, error) {
	contract, err := bindVRFCoordinator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorFilterer{contract: contract}, nil
}

func bindVRFCoordinator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFCoordinatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinator *VRFCoordinatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinator.Contract.VRFCoordinatorCaller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinator *VRFCoordinatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.VRFCoordinatorTransactor.contract.Transfer(opts)
}

func (_VRFCoordinator *VRFCoordinatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.VRFCoordinatorTransactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinator *VRFCoordinatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinator.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinator *VRFCoordinatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.contract.Transfer(opts)
}

func (_VRFCoordinator *VRFCoordinatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinator *VRFCoordinatorCaller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinator.Contract.MAXCONSUMERS(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinator.Contract.MAXCONSUMERS(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) MAXJUELSSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "MAX_JUELS_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) MAXJUELSSUPPLY() (*big.Int, error) {
	return _VRFCoordinator.Contract.MAXJUELSSUPPLY(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) MAXJUELSSUPPLY() (*big.Int, error) {
	return _VRFCoordinator.Contract.MAXJUELSSUPPLY(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) MAXNUMWORDS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "MAX_NUM_WORDS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) MAXNUMWORDS() (*big.Int, error) {
	return _VRFCoordinator.Contract.MAXNUMWORDS(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) MAXNUMWORDS() (*big.Int, error) {
	return _VRFCoordinator.Contract.MAXNUMWORDS(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "NUM_CONF_DELAYS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) NUMCONFDELAYS() (uint8, error) {
	return _VRFCoordinator.Contract.NUMCONFDELAYS(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) NUMCONFDELAYS() (uint8, error) {
	return _VRFCoordinator.Contract.NUMCONFDELAYS(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) CalculateRequestPriceCallbackJuels(opts *bind.CallOpts, gasAllowance *big.Int, arguments []byte) (*big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "calculateRequestPriceCallbackJuels", gasAllowance, arguments)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

func (_VRFCoordinator *VRFCoordinatorSession) CalculateRequestPriceCallbackJuels(gasAllowance *big.Int, arguments []byte) (*big.Int, *big.Int, *big.Int, error) {
	return _VRFCoordinator.Contract.CalculateRequestPriceCallbackJuels(&_VRFCoordinator.CallOpts, gasAllowance, arguments)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) CalculateRequestPriceCallbackJuels(gasAllowance *big.Int, arguments []byte) (*big.Int, *big.Int, *big.Int, error) {
	return _VRFCoordinator.Contract.CalculateRequestPriceCallbackJuels(&_VRFCoordinator.CallOpts, gasAllowance, arguments)
}

func (_VRFCoordinator *VRFCoordinatorCaller) CalculateRequestPriceJuels(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "calculateRequestPriceJuels")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) CalculateRequestPriceJuels() (*big.Int, error) {
	return _VRFCoordinator.Contract.CalculateRequestPriceJuels(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) CalculateRequestPriceJuels() (*big.Int, error) {
	return _VRFCoordinator.Contract.CalculateRequestPriceJuels(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) GetCallbackMemo(opts *bind.CallOpts, requestId *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "getCallbackMemo", requestId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) GetCallbackMemo(requestId *big.Int) ([32]byte, error) {
	return _VRFCoordinator.Contract.GetCallbackMemo(&_VRFCoordinator.CallOpts, requestId)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) GetCallbackMemo(requestId *big.Int) ([32]byte, error) {
	return _VRFCoordinator.Contract.GetCallbackMemo(&_VRFCoordinator.CallOpts, requestId)
}

func (_VRFCoordinator *VRFCoordinatorCaller) GetConfirmationDelays(opts *bind.CallOpts) ([8]*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "getConfirmationDelays")

	if err != nil {
		return *new([8]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([8]*big.Int)).(*[8]*big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) GetConfirmationDelays() ([8]*big.Int, error) {
	return _VRFCoordinator.Contract.GetConfirmationDelays(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) GetConfirmationDelays() ([8]*big.Int, error) {
	return _VRFCoordinator.Contract.GetConfirmationDelays(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) GetFee(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "getFee", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) GetFee(arg0 *big.Int, arg1 []byte) (*big.Int, error) {
	return _VRFCoordinator.Contract.GetFee(&_VRFCoordinator.CallOpts, arg0, arg1)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) GetFee(arg0 *big.Int, arg1 []byte) (*big.Int, error) {
	return _VRFCoordinator.Contract.GetFee(&_VRFCoordinator.CallOpts, arg0, arg1)
}

func (_VRFCoordinator *VRFCoordinatorCaller) GetFulfillmentFee(opts *bind.CallOpts, arg0 *big.Int, callbackGasLimit uint32, arguments []byte, arg3 []byte) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "getFulfillmentFee", arg0, callbackGasLimit, arguments, arg3)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) GetFulfillmentFee(arg0 *big.Int, callbackGasLimit uint32, arguments []byte, arg3 []byte) (*big.Int, error) {
	return _VRFCoordinator.Contract.GetFulfillmentFee(&_VRFCoordinator.CallOpts, arg0, callbackGasLimit, arguments, arg3)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) GetFulfillmentFee(arg0 *big.Int, callbackGasLimit uint32, arguments []byte, arg3 []byte) (*big.Int, error) {
	return _VRFCoordinator.Contract.GetFulfillmentFee(&_VRFCoordinator.CallOpts, arg0, callbackGasLimit, arguments, arg3)
}

func (_VRFCoordinator *VRFCoordinatorCaller) GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

	error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "getSubscription", subId)

	outstruct := new(GetSubscription)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Balance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.PendingFulfillments = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.Owner = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_VRFCoordinator *VRFCoordinatorSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinator.Contract.GetSubscription(&_VRFCoordinator.CallOpts, subId)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinator.Contract.GetSubscription(&_VRFCoordinator.CallOpts, subId)
}

func (_VRFCoordinator *VRFCoordinatorCaller) GetSubscriptionLinkBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "getSubscriptionLinkBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) GetSubscriptionLinkBalance() (*big.Int, error) {
	return _VRFCoordinator.Contract.GetSubscriptionLinkBalance(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) GetSubscriptionLinkBalance() (*big.Int, error) {
	return _VRFCoordinator.Contract.GetSubscriptionLinkBalance(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "i_beaconPeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) IBeaconPeriodBlocks() (*big.Int, error) {
	return _VRFCoordinator.Contract.IBeaconPeriodBlocks(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) IBeaconPeriodBlocks() (*big.Int, error) {
	return _VRFCoordinator.Contract.IBeaconPeriodBlocks(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) ILink(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "i_link")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) ILink() (common.Address, error) {
	return _VRFCoordinator.Contract.ILink(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) ILink() (common.Address, error) {
	return _VRFCoordinator.Contract.ILink(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) ILinkEthFeed(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "i_link_eth_feed")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) ILinkEthFeed() (common.Address, error) {
	return _VRFCoordinator.Contract.ILinkEthFeed(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) ILinkEthFeed() (common.Address, error) {
	return _VRFCoordinator.Contract.ILinkEthFeed(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) IStartSlot(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "i_startSlot")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) IStartSlot() (*big.Int, error) {
	return _VRFCoordinator.Contract.IStartSlot(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) IStartSlot() (*big.Int, error) {
	return _VRFCoordinator.Contract.IStartSlot(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) MigrationVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "migrationVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) MigrationVersion() (uint8, error) {
	return _VRFCoordinator.Contract.MigrationVersion(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) MigrationVersion() (uint8, error) {
	return _VRFCoordinator.Contract.MigrationVersion(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) OnMigration(opts *bind.CallOpts, arg0 []byte) error {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "onMigration", arg0)

	if err != nil {
		return err
	}

	return err

}

func (_VRFCoordinator *VRFCoordinatorSession) OnMigration(arg0 []byte) error {
	return _VRFCoordinator.Contract.OnMigration(&_VRFCoordinator.CallOpts, arg0)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) OnMigration(arg0 []byte) error {
	return _VRFCoordinator.Contract.OnMigration(&_VRFCoordinator.CallOpts, arg0)
}

func (_VRFCoordinator *VRFCoordinatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) Owner() (common.Address, error) {
	return _VRFCoordinator.Contract.Owner(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) Owner() (common.Address, error) {
	return _VRFCoordinator.Contract.Owner(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) Paused() (bool, error) {
	return _VRFCoordinator.Contract.Paused(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) Paused() (bool, error) {
	return _VRFCoordinator.Contract.Paused(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) SBillingConfig(opts *bind.CallOpts) (SBillingConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "s_billingConfig")

	outstruct := new(SBillingConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UseReasonableGasPrice = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.UnusedGasPenaltyPercent = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.StalenessSeconds = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.RedeemableRequestGasOverhead = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.CallbackRequestGasOverhead = *abi.ConvertType(out[4], new(uint32)).(*uint32)
	outstruct.PremiumPercentage = *abi.ConvertType(out[5], new(uint32)).(*uint32)
	outstruct.ReasonableGasPriceStalenessBlocks = *abi.ConvertType(out[6], new(uint32)).(*uint32)
	outstruct.FallbackWeiPerUnitLink = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFCoordinator *VRFCoordinatorSession) SBillingConfig() (SBillingConfig,

	error) {
	return _VRFCoordinator.Contract.SBillingConfig(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) SBillingConfig() (SBillingConfig,

	error) {
	return _VRFCoordinator.Contract.SBillingConfig(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) SConfig(opts *bind.CallOpts) (SConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "s_config")

	outstruct := new(SConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaxCallbackGasLimit = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.MaxCallbackArgumentsLength = *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_VRFCoordinator *VRFCoordinatorSession) SConfig() (SConfig,

	error) {
	return _VRFCoordinator.Contract.SConfig(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) SConfig() (SConfig,

	error) {
	return _VRFCoordinator.Contract.SConfig(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) SPendingRequests(opts *bind.CallOpts, arg0 *big.Int) (SPendingRequests,

	error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "s_pendingRequests", arg0)

	outstruct := new(SPendingRequests)
	if err != nil {
		return *outstruct, err
	}

	outstruct.SlotNumber = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.ConfirmationDelay = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.NumWords = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.Requester = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)

	return *outstruct, err

}

func (_VRFCoordinator *VRFCoordinatorSession) SPendingRequests(arg0 *big.Int) (SPendingRequests,

	error) {
	return _VRFCoordinator.Contract.SPendingRequests(&_VRFCoordinator.CallOpts, arg0)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) SPendingRequests(arg0 *big.Int) (SPendingRequests,

	error) {
	return _VRFCoordinator.Contract.SPendingRequests(&_VRFCoordinator.CallOpts, arg0)
}

func (_VRFCoordinator *VRFCoordinatorCaller) SProducer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "s_producer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) SProducer() (common.Address, error) {
	return _VRFCoordinator.Contract.SProducer(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) SProducer() (common.Address, error) {
	return _VRFCoordinator.Contract.SProducer(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) TypeAndVersion() (string, error) {
	return _VRFCoordinator.Contract.TypeAndVersion(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) TypeAndVersion() (string, error) {
	return _VRFCoordinator.Contract.TypeAndVersion(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "acceptOwnership")
}

func (_VRFCoordinator *VRFCoordinatorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinator.Contract.AcceptOwnership(&_VRFCoordinator.TransactOpts)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinator.Contract.AcceptOwnership(&_VRFCoordinator.TransactOpts)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_VRFCoordinator *VRFCoordinatorSession) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinator.TransactOpts, subId)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinator.TransactOpts, subId)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "addConsumer", subId, consumer)
}

func (_VRFCoordinator *VRFCoordinatorSession) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.AddConsumer(&_VRFCoordinator.TransactOpts, subId, consumer)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.AddConsumer(&_VRFCoordinator.TransactOpts, subId, consumer)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) BatchTransferLink(opts *bind.TransactOpts, recipients []common.Address, paymentsInJuels []*big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "batchTransferLink", recipients, paymentsInJuels)
}

func (_VRFCoordinator *VRFCoordinatorSession) BatchTransferLink(recipients []common.Address, paymentsInJuels []*big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.BatchTransferLink(&_VRFCoordinator.TransactOpts, recipients, paymentsInJuels)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) BatchTransferLink(recipients []common.Address, paymentsInJuels []*big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.BatchTransferLink(&_VRFCoordinator.TransactOpts, recipients, paymentsInJuels)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "cancelSubscription", subId, to)
}

func (_VRFCoordinator *VRFCoordinatorSession) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.CancelSubscription(&_VRFCoordinator.TransactOpts, subId, to)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.CancelSubscription(&_VRFCoordinator.TransactOpts, subId, to)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "createSubscription")
}

func (_VRFCoordinator *VRFCoordinatorSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinator.Contract.CreateSubscription(&_VRFCoordinator.TransactOpts)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinator.Contract.CreateSubscription(&_VRFCoordinator.TransactOpts)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) DeregisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "deregisterMigratableCoordinator", target)
}

func (_VRFCoordinator *VRFCoordinatorSession) DeregisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.DeregisterMigratableCoordinator(&_VRFCoordinator.TransactOpts, target)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) DeregisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.DeregisterMigratableCoordinator(&_VRFCoordinator.TransactOpts, target)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) Migrate(opts *bind.TransactOpts, newCoordinator common.Address, encodedRequest []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "migrate", newCoordinator, encodedRequest)
}

func (_VRFCoordinator *VRFCoordinatorSession) Migrate(newCoordinator common.Address, encodedRequest []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Migrate(&_VRFCoordinator.TransactOpts, newCoordinator, encodedRequest)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) Migrate(newCoordinator common.Address, encodedRequest []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Migrate(&_VRFCoordinator.TransactOpts, newCoordinator, encodedRequest)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_VRFCoordinator *VRFCoordinatorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.OnTokenTransfer(&_VRFCoordinator.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.OnTokenTransfer(&_VRFCoordinator.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "pause")
}

func (_VRFCoordinator *VRFCoordinatorSession) Pause() (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Pause(&_VRFCoordinator.TransactOpts)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) Pause() (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Pause(&_VRFCoordinator.TransactOpts)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) ProcessVRFOutputs(opts *bind.TransactOpts, vrfOutputs []VRFBeaconTypesVRFOutput, juelsPerFeeCoin *big.Int, reasonableGasPrice uint64, blockHeight uint64, arg4 [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "processVRFOutputs", vrfOutputs, juelsPerFeeCoin, reasonableGasPrice, blockHeight, arg4)
}

func (_VRFCoordinator *VRFCoordinatorSession) ProcessVRFOutputs(vrfOutputs []VRFBeaconTypesVRFOutput, juelsPerFeeCoin *big.Int, reasonableGasPrice uint64, blockHeight uint64, arg4 [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.ProcessVRFOutputs(&_VRFCoordinator.TransactOpts, vrfOutputs, juelsPerFeeCoin, reasonableGasPrice, blockHeight, arg4)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) ProcessVRFOutputs(vrfOutputs []VRFBeaconTypesVRFOutput, juelsPerFeeCoin *big.Int, reasonableGasPrice uint64, blockHeight uint64, arg4 [32]byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.ProcessVRFOutputs(&_VRFCoordinator.TransactOpts, vrfOutputs, juelsPerFeeCoin, reasonableGasPrice, blockHeight, arg4)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) RedeemRandomness(opts *bind.TransactOpts, subID *big.Int, requestIDArg *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "redeemRandomness", subID, requestIDArg, arg2)
}

func (_VRFCoordinator *VRFCoordinatorSession) RedeemRandomness(subID *big.Int, requestIDArg *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RedeemRandomness(&_VRFCoordinator.TransactOpts, subID, requestIDArg, arg2)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RedeemRandomness(subID *big.Int, requestIDArg *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RedeemRandomness(&_VRFCoordinator.TransactOpts, subID, requestIDArg, arg2)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) RegisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "registerMigratableCoordinator", target)
}

func (_VRFCoordinator *VRFCoordinatorSession) RegisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RegisterMigratableCoordinator(&_VRFCoordinator.TransactOpts, target)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RegisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RegisterMigratableCoordinator(&_VRFCoordinator.TransactOpts, target)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "removeConsumer", subId, consumer)
}

func (_VRFCoordinator *VRFCoordinatorSession) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RemoveConsumer(&_VRFCoordinator.TransactOpts, subId, consumer)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RemoveConsumer(&_VRFCoordinator.TransactOpts, subId, consumer)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) RequestRandomness(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelayArg *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "requestRandomness", subID, numWords, confDelayArg, arg3)
}

func (_VRFCoordinator *VRFCoordinatorSession) RequestRandomness(subID *big.Int, numWords uint16, confDelayArg *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomness(&_VRFCoordinator.TransactOpts, subID, numWords, confDelayArg, arg3)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RequestRandomness(subID *big.Int, numWords uint16, confDelayArg *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomness(&_VRFCoordinator.TransactOpts, subID, numWords, confDelayArg, arg3)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) RequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, arg5 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "requestRandomnessFulfillment", subID, numWords, confDelayArg, callbackGasLimit, arguments, arg5)
}

func (_VRFCoordinator *VRFCoordinatorSession) RequestRandomnessFulfillment(subID *big.Int, numWords uint16, confDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, arg5 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomnessFulfillment(&_VRFCoordinator.TransactOpts, subID, numWords, confDelayArg, callbackGasLimit, arguments, arg5)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RequestRandomnessFulfillment(subID *big.Int, numWords uint16, confDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, arg5 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomnessFulfillment(&_VRFCoordinator.TransactOpts, subID, numWords, confDelayArg, callbackGasLimit, arguments, arg5)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

func (_VRFCoordinator *VRFCoordinatorSession) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinator.TransactOpts, subId, newOwner)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinator.TransactOpts, subId, newOwner)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) SetBillingConfig(opts *bind.TransactOpts, billingConfig VRFBeaconTypesBillingConfig) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "setBillingConfig", billingConfig)
}

func (_VRFCoordinator *VRFCoordinatorSession) SetBillingConfig(billingConfig VRFBeaconTypesBillingConfig) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetBillingConfig(&_VRFCoordinator.TransactOpts, billingConfig)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) SetBillingConfig(billingConfig VRFBeaconTypesBillingConfig) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetBillingConfig(&_VRFCoordinator.TransactOpts, billingConfig)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) SetConfig(opts *bind.TransactOpts, config VRFCoordinatorConfig) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "setConfig", config)
}

func (_VRFCoordinator *VRFCoordinatorSession) SetConfig(config VRFCoordinatorConfig) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetConfig(&_VRFCoordinator.TransactOpts, config)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) SetConfig(config VRFCoordinatorConfig) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetConfig(&_VRFCoordinator.TransactOpts, config)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) SetConfirmationDelays(opts *bind.TransactOpts, confDelays [8]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "setConfirmationDelays", confDelays)
}

func (_VRFCoordinator *VRFCoordinatorSession) SetConfirmationDelays(confDelays [8]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetConfirmationDelays(&_VRFCoordinator.TransactOpts, confDelays)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) SetConfirmationDelays(confDelays [8]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetConfirmationDelays(&_VRFCoordinator.TransactOpts, confDelays)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) SetProducer(opts *bind.TransactOpts, producer common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "setProducer", producer)
}

func (_VRFCoordinator *VRFCoordinatorSession) SetProducer(producer common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetProducer(&_VRFCoordinator.TransactOpts, producer)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) SetProducer(producer common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetProducer(&_VRFCoordinator.TransactOpts, producer)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) SetReasonableGasPrice(opts *bind.TransactOpts, gasPrice uint64) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "setReasonableGasPrice", gasPrice)
}

func (_VRFCoordinator *VRFCoordinatorSession) SetReasonableGasPrice(gasPrice uint64) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetReasonableGasPrice(&_VRFCoordinator.TransactOpts, gasPrice)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) SetReasonableGasPrice(gasPrice uint64) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetReasonableGasPrice(&_VRFCoordinator.TransactOpts, gasPrice)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) TransferLink(opts *bind.TransactOpts, recipient common.Address, juelsAmount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "transferLink", recipient, juelsAmount)
}

func (_VRFCoordinator *VRFCoordinatorSession) TransferLink(recipient common.Address, juelsAmount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.TransferLink(&_VRFCoordinator.TransactOpts, recipient, juelsAmount)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) TransferLink(recipient common.Address, juelsAmount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.TransferLink(&_VRFCoordinator.TransactOpts, recipient, juelsAmount)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFCoordinator *VRFCoordinatorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.TransferOwnership(&_VRFCoordinator.TransactOpts, to)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.TransferOwnership(&_VRFCoordinator.TransactOpts, to)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "unpause")
}

func (_VRFCoordinator *VRFCoordinatorSession) Unpause() (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Unpause(&_VRFCoordinator.TransactOpts)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) Unpause() (*types.Transaction, error) {
	return _VRFCoordinator.Contract.Unpause(&_VRFCoordinator.TransactOpts)
}

type VRFCoordinatorBillingConfigSetIterator struct {
	Event *VRFCoordinatorBillingConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorBillingConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorBillingConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorBillingConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorBillingConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorBillingConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorBillingConfigSet struct {
	BillingConfig VRFBeaconTypesBillingConfig
	Raw           types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterBillingConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorBillingConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "BillingConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorBillingConfigSetIterator{contract: _VRFCoordinator.contract, event: "BillingConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorBillingConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "BillingConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorBillingConfigSet)
				if err := _VRFCoordinator.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseBillingConfigSet(log types.Log) (*VRFCoordinatorBillingConfigSet, error) {
	event := new(VRFCoordinatorBillingConfigSet)
	if err := _VRFCoordinator.contract.UnpackLog(event, "BillingConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorConfigSetIterator struct {
	Event *VRFCoordinatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
	Raw                       types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorConfigSetIterator{contract: _VRFCoordinator.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorConfigSet)
				if err := _VRFCoordinator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseConfigSet(log types.Log) (*VRFCoordinatorConfigSet, error) {
	event := new(VRFCoordinatorConfigSet)
	if err := _VRFCoordinator.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorCoordinatorConfigSetIterator struct {
	Event *VRFCoordinatorCoordinatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorCoordinatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorCoordinatorConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorCoordinatorConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorCoordinatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorCoordinatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorCoordinatorConfigSet struct {
	NewConfig VRFCoordinatorConfig
	Raw       types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterCoordinatorConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorCoordinatorConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "CoordinatorConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorCoordinatorConfigSetIterator{contract: _VRFCoordinator.contract, event: "CoordinatorConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchCoordinatorConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorCoordinatorConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "CoordinatorConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorCoordinatorConfigSet)
				if err := _VRFCoordinator.contract.UnpackLog(event, "CoordinatorConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseCoordinatorConfigSet(log types.Log) (*VRFCoordinatorCoordinatorConfigSet, error) {
	event := new(VRFCoordinatorCoordinatorConfigSet)
	if err := _VRFCoordinator.contract.UnpackLog(event, "CoordinatorConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorCoordinatorDeregisteredIterator struct {
	Event *VRFCoordinatorCoordinatorDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorCoordinatorDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorCoordinatorDeregistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorCoordinatorDeregistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorCoordinatorDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorCoordinatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorCoordinatorDeregistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorCoordinatorDeregisteredIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorCoordinatorDeregisteredIterator{contract: _VRFCoordinator.contract, event: "CoordinatorDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorCoordinatorDeregistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorCoordinatorDeregistered)
				if err := _VRFCoordinator.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseCoordinatorDeregistered(log types.Log) (*VRFCoordinatorCoordinatorDeregistered, error) {
	event := new(VRFCoordinatorCoordinatorDeregistered)
	if err := _VRFCoordinator.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorCoordinatorRegisteredIterator struct {
	Event *VRFCoordinatorCoordinatorRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorCoordinatorRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorCoordinatorRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorCoordinatorRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorCoordinatorRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorCoordinatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorCoordinatorRegistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFCoordinatorCoordinatorRegisteredIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorCoordinatorRegisteredIterator{contract: _VRFCoordinator.contract, event: "CoordinatorRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorCoordinatorRegistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorCoordinatorRegistered)
				if err := _VRFCoordinator.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseCoordinatorRegistered(log types.Log) (*VRFCoordinatorCoordinatorRegistered, error) {
	event := new(VRFCoordinatorCoordinatorRegistered)
	if err := _VRFCoordinator.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorMigrationCompletedIterator struct {
	Event *VRFCoordinatorMigrationCompleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorMigrationCompletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorMigrationCompleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorMigrationCompleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorMigrationCompletedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorMigrationCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorMigrationCompleted struct {
	NewVersion     uint8
	NewCoordinator common.Address
	SubID          *big.Int
	Raw            types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterMigrationCompleted(opts *bind.FilterOpts, newVersion []uint8, subID []*big.Int) (*VRFCoordinatorMigrationCompletedIterator, error) {

	var newVersionRule []interface{}
	for _, newVersionItem := range newVersion {
		newVersionRule = append(newVersionRule, newVersionItem)
	}

	var subIDRule []interface{}
	for _, subIDItem := range subID {
		subIDRule = append(subIDRule, subIDItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "MigrationCompleted", newVersionRule, subIDRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorMigrationCompletedIterator{contract: _VRFCoordinator.contract, event: "MigrationCompleted", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchMigrationCompleted(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorMigrationCompleted, newVersion []uint8, subID []*big.Int) (event.Subscription, error) {

	var newVersionRule []interface{}
	for _, newVersionItem := range newVersion {
		newVersionRule = append(newVersionRule, newVersionItem)
	}

	var subIDRule []interface{}
	for _, subIDItem := range subID {
		subIDRule = append(subIDRule, subIDItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "MigrationCompleted", newVersionRule, subIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorMigrationCompleted)
				if err := _VRFCoordinator.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseMigrationCompleted(log types.Log) (*VRFCoordinatorMigrationCompleted, error) {
	event := new(VRFCoordinatorMigrationCompleted)
	if err := _VRFCoordinator.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorNewTransmissionIterator struct {
	Event *VRFCoordinatorNewTransmission

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorNewTransmissionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorNewTransmission)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorNewTransmission)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorNewTransmissionIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorNewTransmission struct {
	AggregatorRoundId  uint32
	EpochAndRound      *big.Int
	Transmitter        common.Address
	JuelsPerFeeCoin    *big.Int
	ReasonableGasPrice uint64
	ConfigDigest       [32]byte
	Raw                types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32, epochAndRound []*big.Int) (*VRFCoordinatorNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}
	var epochAndRoundRule []interface{}
	for _, epochAndRoundItem := range epochAndRound {
		epochAndRoundRule = append(epochAndRoundRule, epochAndRoundItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule, epochAndRoundRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorNewTransmissionIterator{contract: _VRFCoordinator.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorNewTransmission, aggregatorRoundId []uint32, epochAndRound []*big.Int) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}
	var epochAndRoundRule []interface{}
	for _, epochAndRoundItem := range epochAndRound {
		epochAndRoundRule = append(epochAndRoundRule, epochAndRoundItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule, epochAndRoundRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorNewTransmission)
				if err := _VRFCoordinator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseNewTransmission(log types.Log) (*VRFCoordinatorNewTransmission, error) {
	event := new(VRFCoordinatorNewTransmission)
	if err := _VRFCoordinator.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorOutputsServedIterator struct {
	Event *VRFCoordinatorOutputsServed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorOutputsServedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorOutputsServed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorOutputsServed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorOutputsServedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorOutputsServedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorOutputsServed struct {
	RecentBlockHeight  uint64
	JuelsPerFeeCoin    *big.Int
	ReasonableGasPrice uint64
	OutputsServed      []VRFBeaconTypesOutputServed
	Raw                types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterOutputsServed(opts *bind.FilterOpts) (*VRFCoordinatorOutputsServedIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "OutputsServed")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorOutputsServedIterator{contract: _VRFCoordinator.contract, event: "OutputsServed", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchOutputsServed(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorOutputsServed) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "OutputsServed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorOutputsServed)
				if err := _VRFCoordinator.contract.UnpackLog(event, "OutputsServed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseOutputsServed(log types.Log) (*VRFCoordinatorOutputsServed, error) {
	event := new(VRFCoordinatorOutputsServed)
	if err := _VRFCoordinator.contract.UnpackLog(event, "OutputsServed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorOwnershipTransferRequestedIterator struct {
	Event *VRFCoordinatorOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorOwnershipTransferRequestedIterator{contract: _VRFCoordinator.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorOwnershipTransferRequested)
				if err := _VRFCoordinator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorOwnershipTransferRequested, error) {
	event := new(VRFCoordinatorOwnershipTransferRequested)
	if err := _VRFCoordinator.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorOwnershipTransferredIterator struct {
	Event *VRFCoordinatorOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorOwnershipTransferredIterator{contract: _VRFCoordinator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorOwnershipTransferred)
				if err := _VRFCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorOwnershipTransferred, error) {
	event := new(VRFCoordinatorOwnershipTransferred)
	if err := _VRFCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorPausedIterator struct {
	Event *VRFCoordinatorPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorPausedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterPaused(opts *bind.FilterOpts) (*VRFCoordinatorPausedIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorPausedIterator{contract: _VRFCoordinator.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorPaused) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorPaused)
				if err := _VRFCoordinator.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParsePaused(log types.Log) (*VRFCoordinatorPaused, error) {
	event := new(VRFCoordinatorPaused)
	if err := _VRFCoordinator.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorRandomWordsFulfilledIterator struct {
	Event *VRFCoordinatorRandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorRandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorRandomWordsFulfilled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorRandomWordsFulfilled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorRandomWordsFulfilled struct {
	RequestIDs            []*big.Int
	SuccessfulFulfillment []byte
	TruncatedErrorData    [][]byte
	SubBalances           []*big.Int
	SubIDs                []*big.Int
	Raw                   types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts) (*VRFCoordinatorRandomWordsFulfilledIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "RandomWordsFulfilled")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorRandomWordsFulfilledIterator{contract: _VRFCoordinator.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomWordsFulfilled) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "RandomWordsFulfilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorRandomWordsFulfilled)
				if err := _VRFCoordinator.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorRandomWordsFulfilled, error) {
	event := new(VRFCoordinatorRandomWordsFulfilled)
	if err := _VRFCoordinator.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorRandomnessFulfillmentRequestedIterator struct {
	Event *VRFCoordinatorRandomnessFulfillmentRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorRandomnessFulfillmentRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorRandomnessFulfillmentRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorRandomnessFulfillmentRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorRandomnessFulfillmentRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorRandomnessFulfillmentRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorRandomnessFulfillmentRequested struct {
	RequestID              *big.Int
	Requester              common.Address
	NextBeaconOutputHeight uint64
	ConfDelay              *big.Int
	SubID                  *big.Int
	NumWords               uint16
	GasAllowance           uint32
	GasPrice               *big.Int
	WeiPerUnitLink         *big.Int
	Arguments              []byte
	CostJuels              *big.Int
	NewSubBalance          *big.Int
	Raw                    types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts, requestID []*big.Int) (*VRFCoordinatorRandomnessFulfillmentRequestedIterator, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "RandomnessFulfillmentRequested", requestIDRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorRandomnessFulfillmentRequestedIterator{contract: _VRFCoordinator.contract, event: "RandomnessFulfillmentRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessFulfillmentRequested, requestID []*big.Int) (event.Subscription, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "RandomnessFulfillmentRequested", requestIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorRandomnessFulfillmentRequested)
				if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessFulfillmentRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseRandomnessFulfillmentRequested(log types.Log) (*VRFCoordinatorRandomnessFulfillmentRequested, error) {
	event := new(VRFCoordinatorRandomnessFulfillmentRequested)
	if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessFulfillmentRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorRandomnessRedeemedIterator struct {
	Event *VRFCoordinatorRandomnessRedeemed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorRandomnessRedeemedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorRandomnessRedeemed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorRandomnessRedeemed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorRandomnessRedeemedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorRandomnessRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorRandomnessRedeemed struct {
	RequestID *big.Int
	Requester common.Address
	SubID     *big.Int
	Raw       types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterRandomnessRedeemed(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFCoordinatorRandomnessRedeemedIterator, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "RandomnessRedeemed", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorRandomnessRedeemedIterator{contract: _VRFCoordinator.contract, event: "RandomnessRedeemed", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchRandomnessRedeemed(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessRedeemed, requestID []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "RandomnessRedeemed", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorRandomnessRedeemed)
				if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRedeemed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseRandomnessRedeemed(log types.Log) (*VRFCoordinatorRandomnessRedeemed, error) {
	event := new(VRFCoordinatorRandomnessRedeemed)
	if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRedeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorRandomnessRequestedIterator struct {
	Event *VRFCoordinatorRandomnessRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorRandomnessRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorRandomnessRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorRandomnessRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorRandomnessRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorRandomnessRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorRandomnessRequested struct {
	RequestID              *big.Int
	Requester              common.Address
	NextBeaconOutputHeight uint64
	ConfDelay              *big.Int
	SubID                  *big.Int
	NumWords               uint16
	CostJuels              *big.Int
	NewSubBalance          *big.Int
	Raw                    types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterRandomnessRequested(opts *bind.FilterOpts, requestID []*big.Int) (*VRFCoordinatorRandomnessRequestedIterator, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "RandomnessRequested", requestIDRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorRandomnessRequestedIterator{contract: _VRFCoordinator.contract, event: "RandomnessRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessRequested, requestID []*big.Int) (event.Subscription, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "RandomnessRequested", requestIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorRandomnessRequested)
				if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseRandomnessRequested(log types.Log) (*VRFCoordinatorRandomnessRequested, error) {
	event := new(VRFCoordinatorRandomnessRequested)
	if err := _VRFCoordinator.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorSubscriptionCanceledIterator struct {
	Event *VRFCoordinatorSubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorSubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorSubscriptionCanceled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorSubscriptionCanceled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorSubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorSubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorSubscriptionCanceled struct {
	SubId  *big.Int
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorSubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorSubscriptionCanceledIterator{contract: _VRFCoordinator.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionCanceled, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorSubscriptionCanceled)
				if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorSubscriptionCanceled, error) {
	event := new(VRFCoordinatorSubscriptionCanceled)
	if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorSubscriptionConsumerAddedIterator struct {
	Event *VRFCoordinatorSubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorSubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorSubscriptionConsumerAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorSubscriptionConsumerAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorSubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorSubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorSubscriptionConsumerAdded struct {
	SubId    *big.Int
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorSubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorSubscriptionConsumerAddedIterator{contract: _VRFCoordinator.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionConsumerAdded, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorSubscriptionConsumerAdded)
				if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorSubscriptionConsumerAdded, error) {
	event := new(VRFCoordinatorSubscriptionConsumerAdded)
	if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorSubscriptionConsumerRemovedIterator struct {
	Event *VRFCoordinatorSubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorSubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorSubscriptionConsumerRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorSubscriptionConsumerRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorSubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorSubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorSubscriptionConsumerRemoved struct {
	SubId    *big.Int
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorSubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorSubscriptionConsumerRemovedIterator{contract: _VRFCoordinator.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionConsumerRemoved, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorSubscriptionConsumerRemoved)
				if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorSubscriptionConsumerRemoved, error) {
	event := new(VRFCoordinatorSubscriptionConsumerRemoved)
	if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorSubscriptionCreatedIterator struct {
	Event *VRFCoordinatorSubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorSubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorSubscriptionCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorSubscriptionCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorSubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorSubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorSubscriptionCreated struct {
	SubId *big.Int
	Owner common.Address
	Raw   types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []*big.Int, owner []common.Address) (*VRFCoordinatorSubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorSubscriptionCreatedIterator{contract: _VRFCoordinator.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionCreated, subId []*big.Int, owner []common.Address) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorSubscriptionCreated)
				if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorSubscriptionCreated, error) {
	event := new(VRFCoordinatorSubscriptionCreated)
	if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorSubscriptionFundedIterator struct {
	Event *VRFCoordinatorSubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorSubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorSubscriptionFunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorSubscriptionFunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorSubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorSubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorSubscriptionFunded struct {
	SubId      *big.Int
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorSubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorSubscriptionFundedIterator{contract: _VRFCoordinator.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionFunded, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorSubscriptionFunded)
				if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorSubscriptionFunded, error) {
	event := new(VRFCoordinatorSubscriptionFunded)
	if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorSubscriptionOwnerTransferRequestedIterator struct {
	Event *VRFCoordinatorSubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorSubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorSubscriptionOwnerTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorSubscriptionOwnerTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorSubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorSubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorSubscriptionOwnerTransferRequested struct {
	SubId *big.Int
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorSubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorSubscriptionOwnerTransferRequestedIterator{contract: _VRFCoordinator.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionOwnerTransferRequested, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorSubscriptionOwnerTransferRequested)
				if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorSubscriptionOwnerTransferRequested, error) {
	event := new(VRFCoordinatorSubscriptionOwnerTransferRequested)
	if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorSubscriptionOwnerTransferredIterator struct {
	Event *VRFCoordinatorSubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorSubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorSubscriptionOwnerTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorSubscriptionOwnerTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorSubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorSubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorSubscriptionOwnerTransferred struct {
	SubId *big.Int
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorSubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorSubscriptionOwnerTransferredIterator{contract: _VRFCoordinator.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionOwnerTransferred, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorSubscriptionOwnerTransferred)
				if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorSubscriptionOwnerTransferred, error) {
	event := new(VRFCoordinatorSubscriptionOwnerTransferred)
	if err := _VRFCoordinator.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorUnpausedIterator struct {
	Event *VRFCoordinatorUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *VRFCoordinatorUnpausedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterUnpaused(opts *bind.FilterOpts) (*VRFCoordinatorUnpausedIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorUnpausedIterator{contract: _VRFCoordinator.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorUnpaused) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorUnpaused)
				if err := _VRFCoordinator.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseUnpaused(log types.Log) (*VRFCoordinatorUnpaused, error) {
	event := new(VRFCoordinatorUnpaused)
	if err := _VRFCoordinator.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetSubscription struct {
	Balance             *big.Int
	PendingFulfillments uint64
	Owner               common.Address
	Consumers           []common.Address
}
type SBillingConfig struct {
	UseReasonableGasPrice             bool
	UnusedGasPenaltyPercent           uint8
	StalenessSeconds                  uint32
	RedeemableRequestGasOverhead      uint32
	CallbackRequestGasOverhead        uint32
	PremiumPercentage                 uint32
	ReasonableGasPriceStalenessBlocks uint32
	FallbackWeiPerUnitLink            *big.Int
}
type SConfig struct {
	MaxCallbackGasLimit        uint32
	MaxCallbackArgumentsLength uint32
}
type SPendingRequests struct {
	SlotNumber        uint32
	ConfirmationDelay *big.Int
	NumWords          uint16
	Requester         common.Address
}

func (_VRFCoordinator *VRFCoordinator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinator.abi.Events["BillingConfigSet"].ID:
		return _VRFCoordinator.ParseBillingConfigSet(log)
	case _VRFCoordinator.abi.Events["ConfigSet"].ID:
		return _VRFCoordinator.ParseConfigSet(log)
	case _VRFCoordinator.abi.Events["CoordinatorConfigSet"].ID:
		return _VRFCoordinator.ParseCoordinatorConfigSet(log)
	case _VRFCoordinator.abi.Events["CoordinatorDeregistered"].ID:
		return _VRFCoordinator.ParseCoordinatorDeregistered(log)
	case _VRFCoordinator.abi.Events["CoordinatorRegistered"].ID:
		return _VRFCoordinator.ParseCoordinatorRegistered(log)
	case _VRFCoordinator.abi.Events["MigrationCompleted"].ID:
		return _VRFCoordinator.ParseMigrationCompleted(log)
	case _VRFCoordinator.abi.Events["NewTransmission"].ID:
		return _VRFCoordinator.ParseNewTransmission(log)
	case _VRFCoordinator.abi.Events["OutputsServed"].ID:
		return _VRFCoordinator.ParseOutputsServed(log)
	case _VRFCoordinator.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFCoordinator.ParseOwnershipTransferRequested(log)
	case _VRFCoordinator.abi.Events["OwnershipTransferred"].ID:
		return _VRFCoordinator.ParseOwnershipTransferred(log)
	case _VRFCoordinator.abi.Events["Paused"].ID:
		return _VRFCoordinator.ParsePaused(log)
	case _VRFCoordinator.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFCoordinator.ParseRandomWordsFulfilled(log)
	case _VRFCoordinator.abi.Events["RandomnessFulfillmentRequested"].ID:
		return _VRFCoordinator.ParseRandomnessFulfillmentRequested(log)
	case _VRFCoordinator.abi.Events["RandomnessRedeemed"].ID:
		return _VRFCoordinator.ParseRandomnessRedeemed(log)
	case _VRFCoordinator.abi.Events["RandomnessRequested"].ID:
		return _VRFCoordinator.ParseRandomnessRequested(log)
	case _VRFCoordinator.abi.Events["SubscriptionCanceled"].ID:
		return _VRFCoordinator.ParseSubscriptionCanceled(log)
	case _VRFCoordinator.abi.Events["SubscriptionConsumerAdded"].ID:
		return _VRFCoordinator.ParseSubscriptionConsumerAdded(log)
	case _VRFCoordinator.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _VRFCoordinator.ParseSubscriptionConsumerRemoved(log)
	case _VRFCoordinator.abi.Events["SubscriptionCreated"].ID:
		return _VRFCoordinator.ParseSubscriptionCreated(log)
	case _VRFCoordinator.abi.Events["SubscriptionFunded"].ID:
		return _VRFCoordinator.ParseSubscriptionFunded(log)
	case _VRFCoordinator.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _VRFCoordinator.ParseSubscriptionOwnerTransferRequested(log)
	case _VRFCoordinator.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _VRFCoordinator.ParseSubscriptionOwnerTransferred(log)
	case _VRFCoordinator.abi.Events["Unpaused"].ID:
		return _VRFCoordinator.ParseUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFCoordinatorBillingConfigSet) Topic() common.Hash {
	return common.HexToHash("0x6dff48329afe669cb5cd0d9af619ec435fe9550e84871839e9df73a1fa34c9bf")
}

func (VRFCoordinatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (VRFCoordinatorCoordinatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0xfffe83c0c6d543712480c43dcc77636fc5671d01e3199dfb237b1bbf29d971e4")
}

func (VRFCoordinatorCoordinatorDeregistered) Topic() common.Hash {
	return common.HexToHash("0xf80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af37")
}

func (VRFCoordinatorCoordinatorRegistered) Topic() common.Hash {
	return common.HexToHash("0xb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af01625")
}

func (VRFCoordinatorMigrationCompleted) Topic() common.Hash {
	return common.HexToHash("0xbd89b747474d3fc04664dfbd1d56ae7ffbe46ee097cdb9979c13916bb76269ce")
}

func (VRFCoordinatorNewTransmission) Topic() common.Hash {
	return common.HexToHash("0x27bf3f1077f091da6885751ba10f5775d06657fd59e47a6ab1f7635e5a115afe")
}

func (VRFCoordinatorOutputsServed) Topic() common.Hash {
	return common.HexToHash("0xf10ea936d00579b4c52035ee33bf46929646b3aa87554c565d8fb2c7aa549c44")
}

func (VRFCoordinatorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFCoordinatorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFCoordinatorPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (VRFCoordinatorRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x8f79f730779e875ce76c428039cc2052b5b5918c2a55c598fab251c1198aec54")
}

func (VRFCoordinatorRandomnessFulfillmentRequested) Topic() common.Hash {
	return common.HexToHash("0x01872fb9c7d6d68af06a17347935e04412da302a377224c205e672c26e18c37f")
}

func (VRFCoordinatorRandomnessRedeemed) Topic() common.Hash {
	return common.HexToHash("0x16f3f633197fafab10a5df69e6f3f2f7f20092f08d8d47de0a91c0f4b96a1a25")
}

func (VRFCoordinatorRandomnessRequested) Topic() common.Hash {
	return common.HexToHash("0xb7933fba96b6b452eb44f99fdc08052a45dff82363d59abaff0456931c3d2459")
}

func (VRFCoordinatorSubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0x3784f77e8e883de95b5d47cd713ced01229fa74d118c0a462224bcb0516d43f1")
}

func (VRFCoordinatorSubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e1")
}

func (VRFCoordinatorSubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a7")
}

func (VRFCoordinatorSubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d")
}

func (VRFCoordinatorSubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0x1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a")
}

func (VRFCoordinatorSubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a1")
}

func (VRFCoordinatorSubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0xd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c9386")
}

func (VRFCoordinatorUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (_VRFCoordinator *VRFCoordinator) Address() common.Address {
	return _VRFCoordinator.address
}

type VRFCoordinatorInterface interface {
	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	MAXJUELSSUPPLY(opts *bind.CallOpts) (*big.Int, error)

	MAXNUMWORDS(opts *bind.CallOpts) (*big.Int, error)

	NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error)

	CalculateRequestPriceCallbackJuels(opts *bind.CallOpts, gasAllowance *big.Int, arguments []byte) (*big.Int, *big.Int, *big.Int, error)

	CalculateRequestPriceJuels(opts *bind.CallOpts) (*big.Int, error)

	GetCallbackMemo(opts *bind.CallOpts, requestId *big.Int) ([32]byte, error)

	GetConfirmationDelays(opts *bind.CallOpts) ([8]*big.Int, error)

	GetFee(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) (*big.Int, error)

	GetFulfillmentFee(opts *bind.CallOpts, arg0 *big.Int, callbackGasLimit uint32, arguments []byte, arg3 []byte) (*big.Int, error)

	GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

		error)

	GetSubscriptionLinkBalance(opts *bind.CallOpts) (*big.Int, error)

	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)

	ILink(opts *bind.CallOpts) (common.Address, error)

	ILinkEthFeed(opts *bind.CallOpts) (common.Address, error)

	IStartSlot(opts *bind.CallOpts) (*big.Int, error)

	MigrationVersion(opts *bind.CallOpts) (uint8, error)

	OnMigration(opts *bind.CallOpts, arg0 []byte) error

	Owner(opts *bind.CallOpts) (common.Address, error)

	Paused(opts *bind.CallOpts) (bool, error)

	SBillingConfig(opts *bind.CallOpts) (SBillingConfig,

		error)

	SConfig(opts *bind.CallOpts) (SConfig,

		error)

	SPendingRequests(opts *bind.CallOpts, arg0 *big.Int) (SPendingRequests,

		error)

	SProducer(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	BatchTransferLink(opts *bind.TransactOpts, recipients []common.Address, paymentsInJuels []*big.Int) (*types.Transaction, error)

	CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	DeregisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error)

	Migrate(opts *bind.TransactOpts, newCoordinator common.Address, encodedRequest []byte) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	ProcessVRFOutputs(opts *bind.TransactOpts, vrfOutputs []VRFBeaconTypesVRFOutput, juelsPerFeeCoin *big.Int, reasonableGasPrice uint64, blockHeight uint64, arg4 [32]byte) (*types.Transaction, error)

	RedeemRandomness(opts *bind.TransactOpts, subID *big.Int, requestIDArg *big.Int, arg2 []byte) (*types.Transaction, error)

	RegisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelayArg *big.Int, arg3 []byte) (*types.Transaction, error)

	RequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, arg5 []byte) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error)

	SetBillingConfig(opts *bind.TransactOpts, billingConfig VRFBeaconTypesBillingConfig) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, config VRFCoordinatorConfig) (*types.Transaction, error)

	SetConfirmationDelays(opts *bind.TransactOpts, confDelays [8]*big.Int) (*types.Transaction, error)

	SetProducer(opts *bind.TransactOpts, producer common.Address) (*types.Transaction, error)

	SetReasonableGasPrice(opts *bind.TransactOpts, gasPrice uint64) (*types.Transaction, error)

	TransferLink(opts *bind.TransactOpts, recipient common.Address, juelsAmount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterBillingConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorBillingConfigSetIterator, error)

	WatchBillingConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorBillingConfigSet) (event.Subscription, error)

	ParseBillingConfigSet(log types.Log) (*VRFCoordinatorBillingConfigSet, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFCoordinatorConfigSet, error)

	FilterCoordinatorConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorCoordinatorConfigSetIterator, error)

	WatchCoordinatorConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorCoordinatorConfigSet) (event.Subscription, error)

	ParseCoordinatorConfigSet(log types.Log) (*VRFCoordinatorCoordinatorConfigSet, error)

	FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorCoordinatorDeregisteredIterator, error)

	WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorCoordinatorDeregistered) (event.Subscription, error)

	ParseCoordinatorDeregistered(log types.Log) (*VRFCoordinatorCoordinatorDeregistered, error)

	FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFCoordinatorCoordinatorRegisteredIterator, error)

	WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorCoordinatorRegistered) (event.Subscription, error)

	ParseCoordinatorRegistered(log types.Log) (*VRFCoordinatorCoordinatorRegistered, error)

	FilterMigrationCompleted(opts *bind.FilterOpts, newVersion []uint8, subID []*big.Int) (*VRFCoordinatorMigrationCompletedIterator, error)

	WatchMigrationCompleted(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorMigrationCompleted, newVersion []uint8, subID []*big.Int) (event.Subscription, error)

	ParseMigrationCompleted(log types.Log) (*VRFCoordinatorMigrationCompleted, error)

	FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32, epochAndRound []*big.Int) (*VRFCoordinatorNewTransmissionIterator, error)

	WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorNewTransmission, aggregatorRoundId []uint32, epochAndRound []*big.Int) (event.Subscription, error)

	ParseNewTransmission(log types.Log) (*VRFCoordinatorNewTransmission, error)

	FilterOutputsServed(opts *bind.FilterOpts) (*VRFCoordinatorOutputsServedIterator, error)

	WatchOutputsServed(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorOutputsServed) (event.Subscription, error)

	ParseOutputsServed(log types.Log) (*VRFCoordinatorOutputsServed, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*VRFCoordinatorPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*VRFCoordinatorPaused, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts) (*VRFCoordinatorRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomWordsFulfilled) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorRandomWordsFulfilled, error)

	FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts, requestID []*big.Int) (*VRFCoordinatorRandomnessFulfillmentRequestedIterator, error)

	WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessFulfillmentRequested, requestID []*big.Int) (event.Subscription, error)

	ParseRandomnessFulfillmentRequested(log types.Log) (*VRFCoordinatorRandomnessFulfillmentRequested, error)

	FilterRandomnessRedeemed(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFCoordinatorRandomnessRedeemedIterator, error)

	WatchRandomnessRedeemed(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessRedeemed, requestID []*big.Int, requester []common.Address) (event.Subscription, error)

	ParseRandomnessRedeemed(log types.Log) (*VRFCoordinatorRandomnessRedeemed, error)

	FilterRandomnessRequested(opts *bind.FilterOpts, requestID []*big.Int) (*VRFCoordinatorRandomnessRequestedIterator, error)

	WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessRequested, requestID []*big.Int) (event.Subscription, error)

	ParseRandomnessRequested(log types.Log) (*VRFCoordinatorRandomnessRequested, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorSubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionCanceled, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorSubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorSubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionConsumerAdded, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorSubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorSubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionConsumerRemoved, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorSubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []*big.Int, owner []common.Address) (*VRFCoordinatorSubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionCreated, subId []*big.Int, owner []common.Address) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorSubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorSubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionFunded, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorSubscriptionFunded, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorSubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionOwnerTransferRequested, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorSubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorSubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorSubscriptionOwnerTransferred, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorSubscriptionOwnerTransferred, error)

	FilterUnpaused(opts *bind.FilterOpts) (*VRFCoordinatorUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*VRFCoordinatorUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
