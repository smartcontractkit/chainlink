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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocksArg\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BeaconPeriodMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"earliestAllowed\",\"type\":\"uint256\"}],\"name\":\"BlockTooRecent\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"providedLength\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxLength\",\"type\":\"uint32\"}],\"name\":\"CallbackArgumentsLengthTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"firstDelay\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"minDelay\",\"type\":\"uint16\"}],\"name\":\"ConfirmationDelayBlocksTooShort\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confirmationDelays\",\"type\":\"uint16[10]\"},{\"internalType\":\"uint8\",\"name\":\"violatingIndex\",\"type\":\"uint8\"}],\"name\":\"ConfirmationDelaysNotIncreasing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasAllowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLeft\",\"type\":\"uint256\"}],\"name\":\"GasAllowanceExceedsGasLeft\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"providedLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxLimit\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reportHeight\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"separatorHeight\",\"type\":\"uint64\"}],\"name\":\"HistoryDomainSeparatorTooOld\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"actualBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requiredBalance\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidBillingConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"expectedLength\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"actualLength\",\"type\":\"uint256\"}],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidJuelsConversion\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numRecipients\",\"type\":\"uint256\"}],\"name\":\"InvalidNumberOfRecipients\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestedSubID\",\"type\":\"uint256\"}],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"requestedVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"coordinatorVersion\",\"type\":\"uint8\"}],\"name\":\"MigrationVersionMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProducer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeRouter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoWordsRequested\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confDelays\",\"type\":\"uint16[10]\"}],\"name\":\"NonZeroDelayAfterZeroDelay\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnMigrationNotSupported\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"producer\",\"type\":\"address\"}],\"name\":\"ProducerAlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"}],\"name\":\"RandomnessNotAvailable\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numRecipients\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numPayments\",\"type\":\"uint256\"}],\"name\":\"RecipientsPaymentsMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"actual\",\"type\":\"address\"}],\"name\":\"ResponseMustBeRetrievedByRequester\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyRequestsReplaceContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManySlotsReplaceContract\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"max\",\"type\":\"uint256\"}],\"name\":\"TooManyWords\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockHeight\",\"type\":\"uint256\"}],\"name\":\"UniverseHasEndedBangBangBang\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"givenDelay\",\"type\":\"uint24\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"knownDelays\",\"type\":\"uint24[8]\"}],\"name\":\"UnknownConfirmationDelay\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"useReasonableGasPrice\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"unusedGasPenaltyPercent\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"redeemableRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callbackRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"premiumPercentage\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceStalenessBlocks\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.BillingConfig\",\"name\":\"billingConfig\",\"type\":\"tuple\"}],\"name\":\"BillingConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCallbackArgumentsLength\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structVRFCoordinator.Config\",\"name\":\"newConfig\",\"type\":\"tuple\"}],\"name\":\"CoordinatorConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"newVersion\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputsServed\",\"type\":\"tuple[]\"}],\"name\":\"OutputsServed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.RequestID[]\",\"name\":\"requestIDs\",\"type\":\"uint48[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"successfulFulfillment\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"truncatedErrorData\",\"type\":\"bytes[]\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAllowance\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_JUELS_SUPPLY\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"paymentsInJuels\",\"type\":\"uint256[]\"}],\"name\":\"batchTransferLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"calculateRequestPriceCallbackJuels\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateRequestPriceJuels\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBillingConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"useReasonableGasPrice\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"unusedGasPenaltyPercent\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"redeemableRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callbackRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"premiumPercentage\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceStalenessBlocks\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"internalType\":\"structVRFBeaconTypes.BillingConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestId\",\"type\":\"uint48\"}],\"name\":\"getCallbackMemo\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfirmationDelays\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"\",\"type\":\"uint24[8]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"getFulfillmentFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"pendingFulfillments\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSubscriptionLinkBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_link\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_link_eth_feed\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_router\",\"outputs\":[{\"internalType\":\"contractVRFRouter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_startSlot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVRFMigration\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"encodedRequest\",\"type\":\"bytes\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrationVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minDelay\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onMigration\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"}],\"internalType\":\"structECCArithmetic.G1Point\",\"name\":\"vrfOutput\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"}],\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"},{\"internalType\":\"uint96\",\"name\":\"price\",\"type\":\"uint96\"}],\"internalType\":\"structVRFBeaconTypes.CostedCallback[]\",\"name\":\"callbacks\",\"type\":\"tuple[]\"}],\"internalType\":\"structVRFBeaconTypes.VRFOutput[]\",\"name\":\"vrfOutputs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"processVRFOutputs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"redeemRandomness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"randomness\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"requestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCallbackArgumentsLength\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"name\":\"s_pendingRequests\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.SlotNumber\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_producer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"useReasonableGasPrice\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"unusedGasPenaltyPercent\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"redeemableRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callbackRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"premiumPercentage\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceStalenessBlocks\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"internalType\":\"structVRFBeaconTypes.BillingConfig\",\"name\":\"billingConfig\",\"type\":\"tuple\"}],\"name\":\"setBillingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCallbackArgumentsLength\",\"type\":\"uint32\"}],\"internalType\":\"structVRFCoordinator.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"confDelays\",\"type\":\"uint24[8]\"}],\"name\":\"setConfirmationDelays\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"producer\",\"type\":\"address\"}],\"name\":\"setProducer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"gasPrice\",\"type\":\"uint64\"}],\"name\":\"setReasonableGasPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"juelsAmount\",\"type\":\"uint256\"}],\"name\":\"transferLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b5060405162006186380380620061868339810160408190526200003591620002bd565b82828233806000816200008f5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c257620000c28162000166565b50506001805460ff60a01b19169055506001600160a01b0392831660805290821660a0521660c05260008490036200010d57604051632abc297960e01b815260040160405180910390fd5b60e084905260006200011e62000211565b9050600060e0518262000132919062000311565b905060008160e0516200014691906200034a565b905062000154818462000366565b61010052506200039695505050505050565b336001600160a01b03821603620001c05760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000086565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60004661a4b181148062000227575062066eed81145b15620002995760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200026d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200029391906200037c565b91505090565b4391505090565b80516001600160a01b0381168114620002b857600080fd5b919050565b60008060008060808587031215620002d457600080fd5b84519350620002e660208601620002a0565b9250620002f660408601620002a0565b91506200030660608601620002a0565b905092959194509250565b6000826200032f57634e487b7160e01b600052601260045260246000fd5b500690565b634e487b7160e01b600052601160045260246000fd5b8181038181111562000360576200036062000334565b92915050565b8082018082111562000360576200036062000334565b6000602082840312156200038f57600080fd5b5051919050565b60805160a05160c05160e05161010051615d1f62000467600039600061035b01526000818161073f01528181611e0601528181613637015281816136660152818161369e015261403201526000818161042401528181610b7501528181610c060152818161109c0152818161159501528181611a6401528181611cf10152818161228d01528181612384015261422101526000818161083e0152613a4101526000818161050701528181610a8201528181611b500152818161266701528181612d6e0152612df60152615d1f6000f3fe608060405234801561001057600080fd5b50600436106102b65760003560e01c80638a8a139011610172578063bec4c08c116100d9578063d0b925fd11610092578063d0b925fd14610774578063dac83d2914610802578063dc311dd314610815578063f27fcfb814610839578063f2fde38b14610860578063f99b1d6814610873578063f9c45ced1461088657600080fd5b8063bec4c08c146106f9578063c3fbb6fd1461070c578063c63c4e9b1461071f578063cb63179714610727578063cd0593df1461073a578063ce3f47191461076157600080fd5b8063a21a23e41161012b578063a21a23e414610685578063a4c0ed361461068d578063abbf1c9b146106a0578063b0724195146106c0578063b2a7cac5146106d3578063bd58017f146106e657600080fd5b80638a8a1390146105465780638d907c621461056e5780638da5cb5b1461063b5780638eef585f1461064c57806395009f081461065f5780639e2010361461067257600080fd5b8063483af7031161022157806369a3164e116101da57806369a3164e146104c157806373433a2f146104d4578063747c91f7146104e757806379ba5097146104fa5780637d253aff146105025780638456cb591461052957806385c64e111461053157600080fd5b8063483af703146103ff578063581bdd161461041f578063597d2f3c146104535780635c975abb1461046457806362f8b6201461047c57806364d51a2a146104a657600080fd5b80632b38bafc116102735780632b38bafc146103a05780632d9297b0146103b35780632f7527cc146103d35780633f4ba83a146103db57806340d6bb82146103e357806346942d18146103ec57600080fd5b806305f4acc6146102bb578063088070f5146102d05780630ae095401461030d578063181f5a77146103205780631a961f2514610356578063294daa491461038b575b600080fd5b6102ce6102c9366004614698565b610899565b005b600b546102eb9063ffffffff80821691600160201b90041682565b6040805163ffffffff9384168152929091166020830152015b60405180910390f35b6102ce61031b3660046146da565b610916565b60408051808201825260148152730565246436f6f7264696e61746f7220312e302e360641b602082015290516103049190614750565b61037d7f000000000000000000000000000000000000000000000000000000000000000081565b604051908152602001610304565b60015b60405160ff9091168152602001610304565b6102ce6103ae366004614763565b610caf565b6103bb610d10565b6040516001600160601b039091168152602001610304565b61038e600881565b6102ce610d51565b61037d6103e881565b6102ce6103fa366004614780565b610d63565b61041261040d3660046147e4565b610de9565b60405161030491906148d7565b6104467f000000000000000000000000000000000000000000000000000000000000000081565b60405161030491906148ea565b6002546001600160601b031661037d565b61046c61107f565b6040519015158152602001610304565b61048f61048a366004614a67565b61108f565b60405165ffffffffffff9091168152602001610304565b6104ae606481565b60405161ffff9091168152602001610304565b6103bb6104cf366004614af8565b611292565b6102ce6104e2366004614b45565b61133c565b6103bb6b033b2e3c9fd0803ce800000081565b6102ce61142d565b6104467f000000000000000000000000000000000000000000000000000000000000000081565b6102ce6114d7565b6105396114e7565b6040516103049190614bd8565b61037d610554366004614bfd565b65ffffffffffff166000908152600c602052604090205490565b61062e6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101919091525060408051610100808201835260045460ff80821615158452918104909116602083015263ffffffff620100008204811693830193909352600160301b810483166060830152600160501b810483166080830152600160701b8104831660a0830152600160901b900490911660c082015260055460e082015290565b6040516103049190614c18565b6000546001600160a01b0316610446565b6102ce61065a366004614caa565b61154c565b61048f61066d366004614cf2565b611588565b61037d610680366004614da5565b6117fe565b61037d611822565b6102ce61069b366004614e65565b611afc565b6106b36106ae366004614eb4565b611ce4565b6040516103049190614f41565b6102ce6106ce366004614f54565b611f0a565b6102ce6106e1366004614f95565b611f7f565b600a54610446906001600160a01b031681565b6102ce6107073660046146da565b6120ab565b6102ce61071a366004614fae565b612252565b6104ae600381565b6102ce6107353660046146da565b6127e6565b61037d7f000000000000000000000000000000000000000000000000000000000000000081565b6102ce61076f366004615002565b612aeb565b6107c5610782366004614bfd565b60106020526000908152604090205463ffffffff811690600160201b810462ffffff1690600160381b810461ffff1690600160481b90046001600160a01b031684565b6040805163ffffffff909516855262ffffff909316602085015261ffff909116918301919091526001600160a01b03166060820152608001610304565b6102ce6108103660046146da565b612b04565b610828610823366004614f95565b612c10565b60405161030495949392919061507c565b6104467f000000000000000000000000000000000000000000000000000000000000000081565b6102ce61086e366004614763565b612d0c565b6102ce6108813660046150d1565b612d1d565b61037d6108943660046150fd565b612e90565b600a546001600160a01b031633146108c457604051634bea32db60e11b815260040160405180910390fd5b60045460ff16156109135760408051808201909152436001600160401b039081168083529083166020909201829052600780546001600160801b031916909117600160401b9092029190911790555b50565b60008281526008602052604090205482906001600160a01b0316806109565760405163c5171ee960e01b8152600481018390526024015b60405180910390fd5b336001600160a01b038216146109815780604051636c51fda960e11b815260040161094d91906148ea565b6000848152600960205260409020548490600160a01b90046001600160401b0316156109c057604051631685ecdd60e31b815260040160405180910390fd5b60065460ff16156109e45760405163769dd35360e11b815260040160405180910390fd5b600085815260096020908152604091829020825160608101845290546001600160601b0381168083526001600160401b03600160601b8304811694840194909452600160a01b90910490921692810192909252610a4087612eaa565b600280546001600160601b03169082906000610a5c8385615143565b92506101000a8154816001600160601b0302191690836001600160601b031602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb88846001600160601b03166040518363ffffffff1660e01b8152600401610aeb9291906001600160a01b03929092168252602082015260400190565b6020604051808303816000875af1158015610b0a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b2e9190615178565b610b5e5760405163cf47918160e01b81526001600160601b0380831660048301528316602482015260440161094d565b60405163677a055360e11b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063cef40aa690610baa9030906004016148ea565b602060405180830381865afa158015610bc7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610beb9190615178565b15610c6b57604051632ee6cacd60e21b8152600481018990527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063bb9b2b3490602401600060405180830381600087803b158015610c5257600080fd5b505af1158015610c66573d6000803e3d6000fd5b505050505b877f3784f77e8e883de95b5d47cd713ced01229fa74d118c0a462224bcb0516d43f18884604051610c9d929190615195565b60405180910390a25050505050505050565b610cb7613008565b600a546001600160a01b031615610cee57600a5460405163ea6d390560e01b815261094d916001600160a01b0316906004016148ea565b600a80546001600160a01b0319166001600160a01b0392909216919091179055565b600080610d1b61305b565b600454610d359190600160301b900463ffffffff166151b7565b6001600160401b03169050610d4b8160006130ec565b91505090565b610d59613008565b610d6161313c565b565b610d6b613008565b6064610d7d60408301602084016151f1565b60ff161115610d9f5760405163015f505560e31b815260040160405180910390fd5b806004610dac828261521b565b9050507f6dff48329afe669cb5cd0d9af619ec435fe9550e84871839e9df73a1fa34c9bf81604051610dde9190615365565b60405180910390a150565b600a546060906001600160a01b03163314610e1757604051634bea32db60e11b815260040160405180910390fd5b610e1f61107f565b15610e3c5760405162461bcd60e51b815260040161094d90615414565b600080876001600160401b03811115610e5757610e57614923565b604051908082528060200260200182016040528015610e9057816020015b610e7d614505565b815260200190600190039081610e755790505b50905060005b88811015610f805760008a8a83818110610eb257610eb261543e565b9050602002810190610ec49190615454565b610ecd906155eb565b9050610eda81888b6131ce565b60408101515151151580610ef657506040810151516020015115155b15610f6d576040805160808101825282516001600160401b0316815260208084015162ffffff168183015283830180515151938301939093529151519091015160608201528351849061ffff8716908110610f5357610f5361543e565b60200260200101819052508380610f69906156c0565b9450505b5080610f78816156e1565b915050610e96565b5060008261ffff166001600160401b03811115610f9f57610f9f614923565b604051908082528060200260200182016040528015610fd857816020015b610fc5614505565b815260200190600190039081610fbd5790505b50905060005b8361ffff1681101561103457828181518110610ffc57610ffc61543e565b60200260200101518282815181106110165761101661543e565b6020026020010181905250808061102c906156e1565b915050610fde565b507ff10ea936d00579b4c52035ee33bf46929646b3aa87554c565d8fb2c7aa549c448689898460405161106a94939291906156fa565b60405180910390a19998505050505050505050565b600154600160a01b900460ff1690565b6000336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146110da5760405163f74c318f60e01b815260040160405180910390fd5b6110e261107f565b156110ff5760405162461bcd60e51b815260040161094d90615414565b60065460ff16156111235760405163769dd35360e11b815260040160405180910390fd5b60008060006111338988886135a5565b925092509250611143898961387a565b65ffffffffffff8316600090815260106020908152604080832085518154848801518489015160608a01516001600160a01b0316600160481b02600160481b600160e81b031961ffff909216600160381b0291909116600160381b600160e81b031962ffffff909316600160201b0266ffffffffffffff1990941663ffffffff909516949094179290921716919091171790558a83526009909152902080546001600160401b03600160601b9091041690600c6111ff8361573f565b82546101009290920a6001600160401b0381810219909316918316021790915560408051918416825262ffffff8916602083015281018a905261ffff891660608201526001600160a01b038b16915065ffffffffffff8516907fc3b31df4232b05afd212fc28027dae6fd6a81618c2a3116182cb57c7f0a3fd0a9060800160405180910390a35090979650505050505050565b60008061129d61305b565b6112a89060106151b7565b6001600160401b03169050600060146112c283601561575b565b6112cc9190615788565b84516112d8919061575b565b6112e061305b565b6004546001600160401b03919091169061130790600160501b900463ffffffff168861579c565b61131191906157bc565b6001600160601b031661132491906157df565b90506113318160006130ec565b925050505b92915050565b600a546001600160a01b0316331461136757604051634bea32db60e11b815260040160405180910390fd5b828015806113755750601f81115b1561139657604051634ecc4fef60e01b81526004810182905260240161094d565b8082146113c05760405163339f8a9d60e01b8152600481018290526024810183905260440161094d565b60005b81811015611425576114138686838181106113e0576113e061543e565b90506020020160208101906113f59190614763565b8585848181106114075761140761543e565b90506020020135612d1d565b8061141d816156e1565b9150506113c3565b505050505050565b6001546001600160a01b031633146114805760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b604482015260640161094d565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6114df613008565b610d616139b2565b6114ef61453b565b6040805161010081019182905290600f90600890826000855b82829054906101000a900462ffffff1662ffffff16815260200190600301906020826002010492830192600103820291508084116115085790505050505050905090565b600a546001600160a01b0316331461157757604051634bea32db60e11b815260040160405180910390fd5b611584600f82600861455a565b5050565b6000336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146115d35760405163f74c318f60e01b815260040160405180910390fd5b6115db61107f565b156115f85760405162461bcd60e51b815260040161094d90615414565b60065460ff161561161c5760405163769dd35360e11b815260040160405180910390fd5b600b5463ffffffff908116908516111561165d57600b54604051637aebf00f60e11b815263ffffffff8087166004830152909116602482015260440161094d565b600b548351600160201b90910463ffffffff1610156116ad578251600b54604051631961a6f960e31b815263ffffffff9283166004820152600160201b909104909116602482015260440161094d565b6000806116bb8a89896135a5565b925050915060006040518061010001604052808465ffffffffffff1681526020018a61ffff1681526020018c6001600160a01b031681526020018781526020018863ffffffff166001600160601b031681526020018b815260200161171e61305b565b6001600160401b03168152602001611734613a12565b9052905061174181613ae9565b81888b8360405160200161175894939291906157f2565b60408051601f19818403018152918152815160209283012065ffffffffffff86166000908152600c90935291205561178f8a613c3a565b8a6001600160a01b03168365ffffffffffff167f24f0e469e0097d1e8d9975137f9f4dd17d2c1481b3a2f25f2382f51287eda1dc848b8e8e8d8860c001518960e001518f6040516117e79897969594939291906158b8565b60405180910390a350909998505050505050505050565b60006118108463ffffffff1684611292565b6001600160601b031695945050505050565b60065460009060ff16156118495760405163769dd35360e11b815260040160405180910390fd5b61185161107f565b1561186e5760405162461bcd60e51b815260040161094d90615414565b60003361187c60014361591d565b6001546040516bffffffffffffffffffffffff19606094851b81166020830152924060348201523090931b90911660548301526001600160c01b0319600160a81b90910460c01b16606882015260700160408051808303601f19018152919052805160209091012060018054919250600160a81b9091046001600160401b03169060156119088361573f565b91906101000a8154816001600160401b0302191690836001600160401b03160217905550506000806001600160401b0381111561194757611947614923565b604051908082528060200260200182016040528015611970578160200160208202803683370190505b506040805160608082018352600080835260208084018281528486018381528984526009835286842095518654925191516001600160601b039091166001600160a01b031993841617600160601b6001600160401b03938416021767ffffffffffffffff60a01b1916600160a01b9290911691909102179094558451928301855233835282810182815283860187815289845260088352959092208351815486166001600160a01b039182161782559251600182018054909616931692909217909355925180519495509093611a4c92600285019201906145f8565b5050604051634a3cdba760e11b8152600481018490527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169150639479b74e90602401600060405180830381600087803b158015611ab157600080fd5b505af1158015611ac5573d6000803e3d6000fd5b50506040513392508491507f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d90600090a350905090565b60065460ff1615611b205760405163769dd35360e11b815260040160405180910390fd5b611b2861107f565b15611b455760405162461bcd60e51b815260040161094d90615414565b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614611b8e576040516344b0e3c360e01b815260040160405180910390fd5b60208114611bb957604051636865567560e01b8152602060048201526024810182905260440161094d565b6000611bc782840184614f95565b6000818152600860205260409020549091506001600160a01b0316611c025760405163c5171ee960e01b81526004810182905260240161094d565b600081815260096020526040812080546001600160601b031691869190611c29838561579c565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600260008282829054906101000a90046001600160601b0316611c71919061579c565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a828784611cc491906157df565b6040805192835260208301919091520160405180910390a2505050505050565b6060336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614611d2f5760405163f74c318f60e01b815260040160405180910390fd5b60065460ff1615611d535760405163769dd35360e11b815260040160405180910390fd5b65ffffffffffff831660008181526010602081815260408084208151608081018352815463ffffffff8116825262ffffff600160201b8204168286015261ffff600160381b820416938201939093526001600160a01b03600160481b840481166060830190815297909652939092526001600160e81b03191690559151811690871614611dfb57806060015186604051638e30e82360e01b815260040161094d929190615930565b8051600090611e31907f00000000000000000000000000000000000000000000000000000000000000009063ffffffff1661575b565b90506000611e3d613cc7565b90506000836020015162ffffff1682611e56919061591d565b9050808310611ea25782846020015162ffffff1684611e7591906157df565b611e809060016157df565b6040516315ad27c360e01b81526004810192909252602482015260440161094d565b6001600160401b03831115611ecd576040516302c6ef8160e11b81526004810184905260240161094d565b6000838152600d602090815260408083208288015162ffffff168452909152902054611efd908890869086613d4b565b9998505050505050505050565b611f12613008565b8051600b80546020808501805163ffffffff908116600160201b0267ffffffffffffffff1990941695811695861793909317909355604080519485529251909116908301527ffffe83c0c6d543712480c43dcc77636fc5671d01e3199dfb237b1bbf29d971e49101610dde565b60065460ff1615611fa35760405163769dd35360e11b815260040160405180910390fd5b6000818152600860205260409020546001600160a01b0316611fdb5760405163c5171ee960e01b81526004810182905260240161094d565b6000818152600860205260409020600101546001600160a01b03163314612032576000818152600860205260409081902060010154905163d084e97560e01b815261094d916001600160a01b0316906004016148ea565b6000818152600860205260409081902080546001600160a01b031980821633908117845560019093018054909116905591516001600160a01b039092169183917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c93869161209f918591615930565b60405180910390a25050565b60008281526008602052604090205482906001600160a01b0316806120e65760405163c5171ee960e01b81526004810183905260240161094d565b336001600160a01b038216146121115780604051636c51fda960e11b815260040161094d91906148ea565b60065460ff16156121355760405163769dd35360e11b815260040160405180910390fd5b61213d61107f565b1561215a5760405162461bcd60e51b815260040161094d90615414565b6000848152600860205260409020600201546063190161218d576040516305a48e0f60e01b815260040160405180910390fd5b6001600160a01b038316600090815260036020908152604080832087845290915290205460ff1661224c576001600160a01b03831660008181526003602090815260408083208884528252808320805460ff1916600190811790915560088352818420600201805491820181558452919092200180546001600160a01b0319169092179091555184907f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e1906122439086906148ea565b60405180910390a25b50505050565b60065460ff16156122765760405163769dd35360e11b815260040160405180910390fd5b60405163677a055360e11b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063cef40aa6906122c29086906004016148ea565b602060405180830381865afa1580156122df573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123039190615178565b612320576040516301fd70a160e51b815260040160405180910390fd5b6040811461234d5760408051636865567560e01b815260048101919091526024810182905260440161094d565b600061235b8284018461594a565b6020810151604051637d331ac560e11b8152600481019190915290915030906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063fa66358a90602401602060405180830381865afa1580156123cb573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123ef9190615984565b6001600160a01b03161461241e57806020015160405163c5171ee960e01b815260040161094d91815260200190565b60008060008060006124338660200151612c10565b94509450945094509450816001600160a01b0316336001600160a01b0316146124715781604051636c51fda960e11b815260040161094d91906148ea565b886001600160a01b031663294daa496040518163ffffffff1660e01b8152600401602060405180830381865afa1580156124af573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124d391906159a1565b60ff16866000015160ff1614612570578560000151896001600160a01b031663294daa496040518163ffffffff1660e01b8152600401602060405180830381865afa158015612526573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061254a91906159a1565b60405163e7aada9560e01b815260ff92831660048201529116602482015260440161094d565b6001600160401b0383161561259857604051631685ecdd60e31b815260040160405180910390fd5b60006040518060c001604052806125ad600190565b60ff16815260200188602001518152602001846001600160a01b03168152602001838152602001876001600160601b03168152602001866001600160401b0316815250905060008160405160200161260591906159be565b60405160208183030381529060405290506126238860200151612eaa565b600280548891906000906126419084906001600160601b0316615143565b92506101000a8154816001600160601b0302191690836001600160601b031602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb8c896040518363ffffffff1660e01b81526004016126b3929190615195565b6020604051808303816000875af11580156126d2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126f69190615178565b6127375760405162461bcd60e51b8152602060048201526012602482015271696e73756666696369656e742066756e647360701b604482015260640161094d565b60405163ce3f471960e01b81526001600160a01b038c169063ce3f471990612763908490600401614750565b600060405180830381600087803b15801561277d57600080fd5b505af1158015612791573d6000803e3d6000fd5b505050508760200151886000015160ff167fbd89b747474d3fc04664dfbd1d56ae7ffbe46ee097cdb9979c13916bb76269ce8d6040516127d191906148ea565b60405180910390a35050505050505050505050565b60008281526008602052604090205482906001600160a01b0316806128215760405163c5171ee960e01b81526004810183905260240161094d565b336001600160a01b0382161461284c5780604051636c51fda960e11b815260040161094d91906148ea565b6000848152600960205260409020548490600160a01b90046001600160401b03161561288b57604051631685ecdd60e31b815260040160405180910390fd5b60065460ff16156128af5760405163769dd35360e11b815260040160405180910390fd5b6001600160a01b038416600090815260036020908152604080832088845290915290205460ff166128f75784846040516379bfd40160e01b815260040161094d929190615a35565b60008581526008602090815260408083206002018054825181850281018501909352808352919290919083018282801561295a57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161293c575b50505050509050600060018251612971919061591d565b905060005b8251811015612a7c57866001600160a01b031683828151811061299b5761299b61543e565b60200260200101516001600160a01b031603612a6a5760008383815181106129c5576129c561543e565b6020026020010151905080600860008b815260200190815260200160002060020183815481106129f7576129f761543e565b600091825260208083209190910180546001600160a01b0319166001600160a01b0394909416939093179092558a8152600890915260409020600201805480612a4257612a42615a4c565b600082815260209020810160001990810180546001600160a01b031916905501905550612a7c565b80612a74816156e1565b915050612976565b506001600160a01b03861660009081526003602090815260408083208a845290915290819020805460ff191690555187907f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a790612ada9089906148ea565b60405180910390a250505050505050565b604051632cb6686f60e01b815260040160405180910390fd5b60008281526008602052604090205482906001600160a01b031680612b3f5760405163c5171ee960e01b81526004810183905260240161094d565b336001600160a01b03821614612b6a5780604051636c51fda960e11b815260040161094d91906148ea565b60065460ff1615612b8e5760405163769dd35360e11b815260040160405180910390fd5b6000848152600860205260409020600101546001600160a01b0384811691161461224c576000848152600860205260409081902060010180546001600160a01b0319166001600160a01b0386161790555184907f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a1906122439033908790615930565b6000818152600860205260408120548190819081906060906001600160a01b0316612c515760405163c5171ee960e01b81526004810187905260240161094d565b60008681526009602090815260408083205460088352928190208054600290910180548351818602810186019094528084526001600160601b038616956001600160401b03600160601b8204811696600160a01b90920416946001600160a01b0390941693918391830182828015612cf257602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612cd4575b505050505090509450945094509450945091939590929450565b612d14613008565b61091381613f1f565b600a546001600160a01b03163314612d4857604051634bea32db60e11b815260040160405180910390fd5b60405163a9059cbb60e01b81526001600160a01b038381166004830152602482018390527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af1158015612db7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ddb9190615178565b611584576040516370a0823160e01b81526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906370a0823190612e2b9030906004016148ea565b602060405180830381865afa158015612e48573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e6c9190615a62565b60405163cf47918160e01b815260048101919091526024810182905260440161094d565b6000612e9a610d10565b6001600160601b03169392505050565b6000818152600860209081526040808320815160608101835281546001600160a01b03908116825260018301541681850152600282018054845181870281018701865281815292959394860193830182828015612f3057602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612f12575b505050505081525050905060005b816040015151811015612faf576003600083604001518381518110612f6557612f6561543e565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081208682529092529020805460ff1916905580612fa7816156e1565b915050612f3e565b50600082815260086020526040812080546001600160a01b03199081168255600182018054909116905590612fe7600283018261464d565b505050600090815260096020526040902080546001600160e01b0319169055565b6000546001600160a01b03163314610d615760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b604482015260640161094d565b60045460009060ff1680156130815750600754600160401b90046001600160401b031615155b156130e757600454600160901b900463ffffffff16431080806130c857506004546130b990600160901b900463ffffffff164361591d565b6007546001600160401b031610155b156130e5575050600754600160401b90046001600160401b031690565b505b503a90565b600454600090819060649061310e90600160701b900463ffffffff1682615a7b565b61311e9063ffffffff168661575b565b6131289190615788565b90506131348184613fc2565b949350505050565b61314461107f565b6131875760405162461bcd60e51b815260206004820152601460248201527314185d5cd8589b194e881b9bdd081c185d5cd95960621b604482015260640161094d565b6001805460ff60a01b191690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b6040516131c491906148ea565b60405180910390a1565b82516001600160401b038084169116111561321257825160405163012d824d60e01b81526001600160401b038085166004830152909116602482015260440161094d565b60408301515151600090158015613230575060408401515160200151155b15613268575082516001600160401b03166000908152600d602090815260408083208287015162ffffff1684529091529020546132c2565b836040015160405160200161327d9190615a98565b60408051601f19818403018152918152815160209283012086516001600160401b03166000908152600d84528281208885015162ffffff168252909352912081905590505b6060840151516000816001600160401b038111156132e2576132e2614923565b60405190808252806020026020018201604052801561330b578160200160208202803683370190505b5090506000826001600160401b0381111561332857613328614923565b6040519080825280601f01601f191660200182016040528015613352576020820181803683370190505b5090506000836001600160401b0381111561336f5761336f614923565b6040519080825280602002602001820160405280156133a257816020015b606081526020019060019003908161338d5790505b5090506000805b858110156134a35760008a6060015182815181106133c9576133c961543e565b602002602001015190506000806133ea8d600001518e602001518c86614028565b9150915081156134295780868661ffff168151811061340b5761340b61543e565b60200260200101819052508480613421906156c0565b955050613458565b600160f81b8785815181106134405761344061543e565b60200101906001600160f81b031916908160001a9053505b825151885189908690811061346f5761346f61543e565b602002602001019065ffffffffffff16908165ffffffffffff1681525050505050808061349b906156e1565b9150506133a9565b506060890151511561359a5760008161ffff166001600160401b038111156134cd576134cd614923565b60405190808252806020026020018201604052801561350057816020015b60608152602001906001900390816134eb5790505b50905060005b8261ffff1681101561355c578381815181106135245761352461543e565b602002602001015182828151811061353e5761353e61543e565b60200260200101819052508080613554906156e1565b915050613506565b507f47ddf7bb0cbd94c1b43c5097f1352a80db0ceb3696f029d32b24f32cd631d2b785858360405161359093929190615acb565b60405180910390a1505b505050505050505050565b604080516080810182526000808252602082018190529181018290526060810182905260006103e88561ffff1611156135ff57604051634a90778560e01b815261ffff861660048201526103e8602482015260440161094d565b8461ffff16600003613624576040516308fad2a760e01b815260040160405180910390fd5b600061362e613cc7565b9050600061365c7f000000000000000000000000000000000000000000000000000000000000000083615b81565b905060008161368b7f0000000000000000000000000000000000000000000000000000000000000000856157df565b613695919061591d565b905060006136c37f000000000000000000000000000000000000000000000000000000000000000083615788565b905063ffffffff81106136e9576040516307b2a52360e41b815260040160405180910390fd5b604080518082018252600e805465ffffffffffff1682528251610100810193849052849360009392916020840191600f906008908288855b82829054906101000a900462ffffff1662ffffff168152602001906003019060208260020104928301926001038202915080841161372157905050505091909252505081519192505065ffffffffffff9081161061379257604051630568cab760e31b815260040160405180910390fd5b805161379f906001615b95565b600e805465ffffffffffff191665ffffffffffff9290921691909117905560005b6008811015613806578a62ffffff16826020015182600881106137e5576137e561543e565b602002015162ffffff161461380657806137fe816156e1565b9150506137c0565b6008811061382e576020820151604051630c4f769b60e41b815261094d918d91600401615bb4565b506040805160808101825263ffffffff909316835262ffffff8b16602084015261ffff8c16908301526001600160a01b038c166060830152519750955090935050505093509350939050565b6001600160a01b038216600090815260036020908152604080832084845290915290205460ff166138c25780826040516379bfd40160e01b815260040161094d929190615a35565b60006138cc610d10565b600083815260096020526040902080546001600160601b0392831693509091168281101561392257815460405163cf47918160e01b81526001600160601b0390911660048201526024810184905260440161094d565b81548390839060009061393f9084906001600160601b0316615143565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555082600260008282829054906101000a90046001600160601b03166139879190615143565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505050505050565b6139ba61107f565b156139d75760405162461bcd60e51b815260040161094d90615414565b6001805460ff60a01b1916600160a01b1790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586131b73390565b6004805460408051633fabe5a360e21b815290516000936201000090930463ffffffff169283151592859283927f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169263feaf968c928183019260a0928290030181865afa158015613a90573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ab49190615be8565b509450909250849150508015613ad85750613acf824261591d565b8463ffffffff16105b156131345750600554949350505050565b6040808201516001600160a01b031660009081526003602090815282822060a085015183529052205460ff16613b3e578060a0015181604001516040516379bfd40160e01b815260040161094d929190615a35565b6000613b5282608001518360600151611292565b60a0830151600090815260096020526040902080546001600160601b039283169350909116821115613bac57805460405163cf47918160e01b81526001600160601b0390911660048201526024810183905260440161094d565b805482908290600090613bc99084906001600160601b0316615143565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555081600260008282829054906101000a90046001600160601b0316613c119190615143565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550505050565b60008181526009602052604090208054600160601b90046001600160401b031681600c613c668361573f565b82546101009290920a6001600160401b038181021990931691831602179091558254600160a01b9004169050816014613c9e8361573f565b91906101000a8154816001600160401b0302191690836001600160401b03160217905550505050565b60004661a4b1811480613cdc575062066eed81145b15613d445760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613d20573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d4b9190615a62565b4391505090565b606082613d845760405163c7d41b1b60e01b815265ffffffffffff861660048201526001600160401b038316602482015260440161094d565b6040805165ffffffffffff8716602080830191909152865163ffffffff168284015286015162ffffff166060808301919091529186015161ffff166080820152908501516001600160a01b031660a082015260c0810184905260009060e0016040516020818303038152906040528051906020012090506103e8856040015161ffff161115613e3a576040808601519051634a90778560e01b815261ffff90911660048201526103e8602482015260440161094d565b6000856040015161ffff166001600160401b03811115613e5c57613e5c614923565b604051908082528060200260200182016040528015613e85578160200160208202803683370190505b50905060005b866040015161ffff168161ffff161015613f14578281604051602001613ec892919091825260f01b6001600160f01b031916602082015260220190565b6040516020818303038152906040528051906020012060001c828261ffff1681518110613ef757613ef761543e565b602090810291909101015280613f0c816156c0565b915050613e8b565b509695505050505050565b336001600160a01b03821603613f715760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b604482015260640161094d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000808215613fd15782613fd9565b613fd9613a12565b9050600081613ff086670de0b6b3a764000061575b565b613ffa9190615788565b90506b033b2e3c9fd0803ce80000008111156131345760405162de437160e81b815260040160405180910390fd5b60006060816140607f00000000000000000000000000000000000000000000000000000000000000006001600160401b038916615788565b845160a08101516040519293509091600091614084918b918b9186906020016157f2565b60408051601f198184030181529181528151602092830120845165ffffffffffff166000908152600c90935291205490915081146140f35760016040518060400160405280601081526020016f756e6b6e6f776e2063616c6c6261636b60801b815250945094505050506143d4565b6040805160808101825263ffffffff8516815262ffffff8a1660208083019190915284015161ffff1681830152908301516001600160a01b031660608201528761417557600160405180604001604052806016815260200175756e617661696c61626c652072616e646f6d6e65737360501b81525095509550505050506143d4565b60006141878460000151838b8e613d4b565b6060808401518651918701516040519394509092600092635a47dd7160e01b926141b692879190602401615c38565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091526006805460ff191660011790559050600080805a8d51608001516040808c015190516355fe976360e01b81529293506001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016926355fe9763926142589290918990600401615c6b565b60408051808303816000875af1158015614276573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061429a9190615ca7565b9093509150816142d6578c516080015160405163aad1598360e01b81526001600160601b0390911660048201526024810182905260440161094d565b60006113885a6142e691906157df565b6006805460ff1916905590508181101561430e5761430e614307828461591d565b8f516143dd565b8d5160a0015160009081526009602052604090208054600160a01b90046001600160401b031690601461434083615cd6565b82546001600160401b039182166101009390930a928302919092021990911617905550895165ffffffffffff166000908152600c6020526040812055836143b15760016040518060400160405280601081526020016f195e1958dd5d1a5bdb8819985a5b195960821b8152506143c4565b6000604051806020016040528060008152505b9c509c5050505050505050505050505b94509492505050565b80608001516001600160601b03168211156143f6575050565b60045460009060649061441190610100900460ff1682615cf9565b60ff168360c001518585608001516001600160601b0316614432919061591d565b61443c919061575b565b614446919061575b565b6144509190615788565b90506000614462828460e00151613fc2565b60a08401516000908152600960205260408120805492935083929091906144939084906001600160601b031661579c565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600260008282829054906101000a90046001600160601b03166144db919061579c565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555050505050565b604051806080016040528060006001600160401b03168152602001600062ffffff16815260200160008152602001600081525090565b6040518061010001604052806008906020820280368337509192915050565b6001830191839082156145e85791602002820160005b838211156145b757833562ffffff1683826101000a81548162ffffff021916908362ffffff1602179055509260200192600301602081600201049283019260010302614570565b80156145e65782816101000a81549062ffffff02191690556003016020816002010492830192600103026145b7565b505b506145f4929150614667565b5090565b8280548282559060005260206000209081019282156145e8579160200282015b828111156145e857825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190614618565b508054600082559060005260206000209081019061091391905b5b808211156145f45760008155600101614668565b80356001600160401b038116811461469357600080fd5b919050565b6000602082840312156146aa57600080fd5b6146b38261467c565b9392505050565b6001600160a01b038116811461091357600080fd5b8035614693816146ba565b600080604083850312156146ed57600080fd5b8235915060208301356146ff816146ba565b809150509250929050565b6000815180845260005b8181101561473057602081850181015186830182015201614714565b506000602082860101526020601f19601f83011685010191505092915050565b6020815260006146b3602083018461470a565b60006020828403121561477557600080fd5b81356146b3816146ba565b6000610100828403121561479357600080fd5b50919050565b60008083601f8401126147ab57600080fd5b5081356001600160401b038111156147c257600080fd5b6020830191508360208260051b85010111156147dd57600080fd5b9250929050565b60008060008060008060a087890312156147fd57600080fd5b86356001600160401b0381111561481357600080fd5b61481f89828a01614799565b90975095505060208701356001600160c01b038116811461483f57600080fd5b935061484d6040880161467c565b925061485b6060880161467c565b9150608087013590509295509295509295565b600081518084526020808501945080840160005b838110156148cc57815180516001600160401b031688528381015162ffffff1684890152604080820151908901526060908101519088015260809096019590820190600101614882565b509495945050505050565b6020815260006146b3602083018461486e565b6001600160a01b0391909116815260200190565b803561ffff8116811461469357600080fd5b803562ffffff8116811461469357600080fd5b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b038111828210171561495b5761495b614923565b60405290565b60405161010081016001600160401b038111828210171561495b5761495b614923565b604051608081016001600160401b038111828210171561495b5761495b614923565b604051602081016001600160401b038111828210171561495b5761495b614923565b604051601f8201601f191681016001600160401b03811182821017156149f0576149f0614923565b604052919050565b600082601f830112614a0957600080fd5b81356001600160401b03811115614a2257614a22614923565b614a35601f8201601f19166020016149c8565b818152846020838601011115614a4a57600080fd5b816020850160208301376000918101602001919091529392505050565b600080600080600060a08688031215614a7f57600080fd5b8535614a8a816146ba565b945060208601359350614a9f604087016148fe565b9250614aad60608701614910565b915060808601356001600160401b03811115614ac857600080fd5b614ad4888289016149f8565b9150509295509295909350565b80356001600160601b038116811461469357600080fd5b60008060408385031215614b0b57600080fd5b614b1483614ae1565b915060208301356001600160401b03811115614b2f57600080fd5b614b3b858286016149f8565b9150509250929050565b60008060008060408587031215614b5b57600080fd5b84356001600160401b0380821115614b7257600080fd5b614b7e88838901614799565b90965094506020870135915080821115614b9757600080fd5b50614ba487828801614799565b95989497509550505050565b8060005b600881101561224c57815162ffffff16845260209384019390910190600101614bb4565b61010081016113368284614bb0565b803565ffffffffffff8116811461469357600080fd5b600060208284031215614c0f57600080fd5b6146b382614be7565b60006101008201905082511515825260ff6020840151166020830152604083015163ffffffff808216604085015280606086015116606085015250506080830151614c6b608084018263ffffffff169052565b5060a0830151614c8360a084018263ffffffff169052565b5060c0830151614c9b60c084018263ffffffff169052565b5060e092830151919092015290565b6000610100808385031215614cbe57600080fd5b838184011115614ccd57600080fd5b509092915050565b63ffffffff8116811461091357600080fd5b803561469381614cd5565b600080600080600080600060e0888a031215614d0d57600080fd5b8735614d18816146ba565b965060208801359550614d2d604089016148fe565b9450614d3b60608901614910565b93506080880135614d4b81614cd5565b925060a08801356001600160401b0380821115614d6757600080fd5b614d738b838c016149f8565b935060c08a0135915080821115614d8957600080fd5b50614d968a828b016149f8565b91505092959891949750929550565b60008060008060808587031215614dbb57600080fd5b843593506020850135614dcd81614cd5565b925060408501356001600160401b0380821115614de957600080fd5b614df5888389016149f8565b93506060870135915080821115614e0b57600080fd5b50614e18878288016149f8565b91505092959194509250565b60008083601f840112614e3657600080fd5b5081356001600160401b03811115614e4d57600080fd5b6020830191508360208285010111156147dd57600080fd5b60008060008060608587031215614e7b57600080fd5b8435614e86816146ba565b93506020850135925060408501356001600160401b03811115614ea857600080fd5b614ba487828801614e24565b60008060008060808587031215614eca57600080fd5b8435614ed5816146ba565b935060208501359250614eea60408601614be7565b915060608501356001600160401b03811115614f0557600080fd5b614e18878288016149f8565b600081518084526020808501945080840160005b838110156148cc57815187529582019590820190600101614f25565b6020815260006146b36020830184614f11565b600060408284031215614f6657600080fd5b614f6e614939565b8235614f7981614cd5565b81526020830135614f8981614cd5565b60208201529392505050565b600060208284031215614fa757600080fd5b5035919050565b600080600060408486031215614fc357600080fd5b8335614fce816146ba565b925060208401356001600160401b03811115614fe957600080fd5b614ff586828701614e24565b9497909650939450505050565b6000806020838503121561501557600080fd5b82356001600160401b0381111561502b57600080fd5b61503785828601614e24565b90969095509350505050565b600081518084526020808501945080840160005b838110156148cc5781516001600160a01b031687529582019590820190600101615057565b6001600160601b03861681526001600160401b038581166020830152841660408201526001600160a01b038316606082015260a0608082018190526000906150c690830184615043565b979650505050505050565b600080604083850312156150e457600080fd5b82356150ef816146ba565b946020939093013593505050565b6000806040838503121561511057600080fd5b8235915060208301356001600160401b03811115614b2f57600080fd5b634e487b7160e01b600052601160045260246000fd5b6001600160601b038281168282160390808211156151635761516361512d565b5092915050565b801515811461091357600080fd5b60006020828403121561518a57600080fd5b81516146b38161516a565b6001600160a01b039290921682526001600160601b0316602082015260400190565b6001600160401b038181168382160280821691908281146151da576151da61512d565b505092915050565b60ff8116811461091357600080fd5b60006020828403121561520357600080fd5b81356146b3816151e2565b6000813561133681614cd5565b81356152268161516a565b815460ff19811691151560ff1691821783556020840135615246816151e2565b61ff008160081b168361ffff1984161717845550505061528961526b6040840161520e565b825465ffffffff0000191660109190911b65ffffffff000016178255565b6152be6152986060840161520e565b825469ffffffff000000000000191660309190911b69ffffffff00000000000016178255565b6152f16152cd6080840161520e565b82805463ffffffff60501b191660509290921b63ffffffff60501b16919091179055565b61532461530060a0840161520e565b82805463ffffffff60701b191660709290921b63ffffffff60701b16919091179055565b61535761533360c0840161520e565b82805463ffffffff60901b191660909290921b63ffffffff60901b16919091179055565b60e082013560018201555050565b610100810182356153758161516a565b151582526020830135615387816151e2565b60ff166020830152604083013561539d81614cd5565b63ffffffff1660408301526153b460608401614ce7565b63ffffffff1660608301526153cb60808401614ce7565b63ffffffff1660808301526153e260a08401614ce7565b63ffffffff1660a08301526153f960c08401614ce7565b63ffffffff811660c08401525060e092830135919092015290565b60208082526010908201526f14185d5cd8589b194e881c185d5cd95960821b604082015260600190565b634e487b7160e01b600052603260045260246000fd5b60008235609e1983360301811261546a57600080fd5b9190910192915050565b600082601f83011261548557600080fd5b813560206001600160401b03808311156154a1576154a1614923565b8260051b6154b08382016149c8565b93845285810183019383810190888611156154ca57600080fd5b84880192505b858310156155df578235848111156154e757600080fd5b8801601f196040828c03820112156154fe57600080fd5b615506614939565b878301358781111561551757600080fd5b8301610100818e038401121561552c57600080fd5b615534614961565b9250615541898201614be7565b835261554f604082016148fe565b8984015261555f606082016146cf565b604084015260808101358881111561557657600080fd5b6155848e8b838501016149f8565b60608501525061559660a08201614ae1565b608084015260c081013560a084015260e081013560c084015261010081013560e0840152508181526155ca60408401614ae1565b818901528452505091840191908401906154d0565b98975050505050505050565b600081360360a08112156155fe57600080fd5b615606614984565b61560f8461467c565b8152602061561e818601614910565b828201526040603f198401121561563457600080fd5b61563c6149a6565b925036605f86011261564d57600080fd5b615655614939565b80608087013681111561566757600080fd5b604088015b81811015615683578035845292840192840161566c565b50908552604084019490945250509035906001600160401b038211156156a857600080fd5b6156b436838601615474565b60608201529392505050565b600061ffff8083168181036156d7576156d761512d565b6001019392505050565b6000600182016156f3576156f361512d565b5060010190565b6001600160401b0385811682526001600160c01b0385166020830152831660408201526080606082018190526000906157359083018461486e565b9695505050505050565b60006001600160401b038083168181036156d7576156d761512d565b80820281158282048414176113365761133661512d565b634e487b7160e01b600052601260045260246000fd5b60008261579757615797615772565b500490565b6001600160601b038181168382160190808211156151635761516361512d565b6001600160601b038181168382160280821691908281146151da576151da61512d565b808201808211156113365761133661512d565b6001600160401b038516815262ffffff841660208201528260408201526080606082015265ffffffffffff82511660808201526000602083015161583c60a084018261ffff169052565b5060408301516001600160a01b031660c0830152606083015161010060e0840181905261586d61018085018361470a565b91506080850151615888828601826001600160601b03169052565b505060a084015161012084015260c084015161014084015260e08401516101608401528091505095945050505050565b60006101006001600160401b038b16835262ffffff8a16602084015288604084015261ffff8816606084015263ffffffff871660808401528560a08401528460c08401528060e084015261590e8184018561470a565b9b9a5050505050505050505050565b818103818111156113365761133661512d565b6001600160a01b0392831681529116602082015260400190565b60006040828403121561595c57600080fd5b615964614939565b823561596f816151e2565b81526020928301359281019290925250919050565b60006020828403121561599657600080fd5b81516146b3816146ba565b6000602082840312156159b357600080fd5b81516146b3816151e2565b6020815260ff82511660208201526020820151604082015260018060a01b0360408301511660608201526000606083015160c06080840152615a0360e0840182615043565b90506001600160601b0360808501511660a08401526001600160401b0360a08501511660c08401528091505092915050565b9182526001600160a01b0316602082015260400190565b634e487b7160e01b600052603160045260246000fd5b600060208284031215615a7457600080fd5b5051919050565b63ffffffff8181168382160190808211156151635761516361512d565b815160408201908260005b6002811015615ac2578251825260209283019290910190600101615aa3565b50505092915050565b606080825284519082018190526000906020906080840190828801845b82811015615b0c57815165ffffffffffff1684529284019290840190600101615ae8565b50505083810382850152615b20818761470a565b905083810360408501528085518083528383019150838160051b84010184880160005b83811015615b7157601f19868403018552615b5f83835161470a565b94870194925090860190600101615b43565b50909a9950505050505050505050565b600082615b9057615b90615772565b500690565b65ffffffffffff8181168382160190808211156151635761516361512d565b62ffffff8316815261012081016146b36020830184614bb0565b805169ffffffffffffffffffff8116811461469357600080fd5b600080600080600060a08688031215615c0057600080fd5b615c0986615bce565b9450602086015193506040860151925060608601519150615c2c60808701615bce565b90509295509295909350565b65ffffffffffff84168152606060208201526000615c596060830185614f11565b8281036040840152615735818561470a565b6001600160601b03841681526001600160a01b0383166020820152606060408201819052600090615c9e9083018461470a565b95945050505050565b60008060408385031215615cba57600080fd5b8251615cc58161516a565b60208401519092506146ff8161516a565b60006001600160401b03821680615cef57615cef61512d565b6000190192915050565b60ff82811682821603908111156113365761133661512d56fea164736f6c6343000813000a",
}

var VRFCoordinatorABI = VRFCoordinatorMetaData.ABI

var VRFCoordinatorBin = VRFCoordinatorMetaData.Bin

func DeployVRFCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, beaconPeriodBlocksArg *big.Int, linkToken common.Address, linkEthFeed common.Address, router common.Address) (common.Address, *types.Transaction, *VRFCoordinator, error) {
	parsed, err := VRFCoordinatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorBin), backend, beaconPeriodBlocksArg, linkToken, linkEthFeed, router)
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

func (_VRFCoordinator *VRFCoordinatorCaller) CalculateRequestPriceCallbackJuels(opts *bind.CallOpts, gasAllowance *big.Int, arguments []byte) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "calculateRequestPriceCallbackJuels", gasAllowance, arguments)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) CalculateRequestPriceCallbackJuels(gasAllowance *big.Int, arguments []byte) (*big.Int, error) {
	return _VRFCoordinator.Contract.CalculateRequestPriceCallbackJuels(&_VRFCoordinator.CallOpts, gasAllowance, arguments)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) CalculateRequestPriceCallbackJuels(gasAllowance *big.Int, arguments []byte) (*big.Int, error) {
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

func (_VRFCoordinator *VRFCoordinatorCaller) GetBillingConfig(opts *bind.CallOpts) (VRFBeaconTypesBillingConfig, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "getBillingConfig")

	if err != nil {
		return *new(VRFBeaconTypesBillingConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(VRFBeaconTypesBillingConfig)).(*VRFBeaconTypesBillingConfig)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) GetBillingConfig() (VRFBeaconTypesBillingConfig, error) {
	return _VRFCoordinator.Contract.GetBillingConfig(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) GetBillingConfig() (VRFBeaconTypesBillingConfig, error) {
	return _VRFCoordinator.Contract.GetBillingConfig(&_VRFCoordinator.CallOpts)
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
	outstruct.ReqCount = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.PendingFulfillments = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Owner = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[4], new([]common.Address)).(*[]common.Address)

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

func (_VRFCoordinator *VRFCoordinatorCaller) IRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "i_router")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) IRouter() (common.Address, error) {
	return _VRFCoordinator.Contract.IRouter(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) IRouter() (common.Address, error) {
	return _VRFCoordinator.Contract.IRouter(&_VRFCoordinator.CallOpts)
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

func (_VRFCoordinator *VRFCoordinatorCaller) MinDelay(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "minDelay")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) MinDelay() (uint16, error) {
	return _VRFCoordinator.Contract.MinDelay(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) MinDelay() (uint16, error) {
	return _VRFCoordinator.Contract.MinDelay(&_VRFCoordinator.CallOpts)
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

func (_VRFCoordinator *VRFCoordinatorTransactor) RedeemRandomness(opts *bind.TransactOpts, sender common.Address, arg1 *big.Int, requestID *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "redeemRandomness", sender, arg1, requestID, arg3)
}

func (_VRFCoordinator *VRFCoordinatorSession) RedeemRandomness(sender common.Address, arg1 *big.Int, requestID *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RedeemRandomness(&_VRFCoordinator.TransactOpts, sender, arg1, requestID, arg3)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RedeemRandomness(sender common.Address, arg1 *big.Int, requestID *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RedeemRandomness(&_VRFCoordinator.TransactOpts, sender, arg1, requestID, arg3)
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

func (_VRFCoordinator *VRFCoordinatorTransactor) RequestRandomness(opts *bind.TransactOpts, requester common.Address, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "requestRandomness", requester, subID, numWords, confirmationDelayArg, arg4)
}

func (_VRFCoordinator *VRFCoordinatorSession) RequestRandomness(requester common.Address, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomness(&_VRFCoordinator.TransactOpts, requester, subID, numWords, confirmationDelayArg, arg4)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RequestRandomness(requester common.Address, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomness(&_VRFCoordinator.TransactOpts, requester, subID, numWords, confirmationDelayArg, arg4)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) RequestRandomnessFulfillment(opts *bind.TransactOpts, requester common.Address, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, arg6 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "requestRandomnessFulfillment", requester, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments, arg6)
}

func (_VRFCoordinator *VRFCoordinatorSession) RequestRandomnessFulfillment(requester common.Address, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, arg6 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomnessFulfillment(&_VRFCoordinator.TransactOpts, requester, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments, arg6)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RequestRandomnessFulfillment(requester common.Address, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, arg6 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomnessFulfillment(&_VRFCoordinator.TransactOpts, requester, subID, numWords, confirmationDelayArg, callbackGasLimit, arguments, arg6)
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
	Raw                    types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFCoordinatorRandomnessFulfillmentRequestedIterator, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "RandomnessFulfillmentRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorRandomnessFulfillmentRequestedIterator{contract: _VRFCoordinator.contract, event: "RandomnessFulfillmentRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessFulfillmentRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "RandomnessFulfillmentRequested", requestIDRule, requesterRule)
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
	Raw                    types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterRandomnessRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFCoordinatorRandomnessRequestedIterator, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "RandomnessRequested", requestIDRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorRandomnessRequestedIterator{contract: _VRFCoordinator.contract, event: "RandomnessRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIDRule []interface{}
	for _, requestIDItem := range requestID {
		requestIDRule = append(requestIDRule, requestIDItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "RandomnessRequested", requestIDRule, requesterRule)
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
	ReqCount            uint64
	PendingFulfillments uint64
	Owner               common.Address
	Consumers           []common.Address
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
	return common.HexToHash("0x47ddf7bb0cbd94c1b43c5097f1352a80db0ceb3696f029d32b24f32cd631d2b7")
}

func (VRFCoordinatorRandomnessFulfillmentRequested) Topic() common.Hash {
	return common.HexToHash("0x24f0e469e0097d1e8d9975137f9f4dd17d2c1481b3a2f25f2382f51287eda1dc")
}

func (VRFCoordinatorRandomnessRequested) Topic() common.Hash {
	return common.HexToHash("0xc3b31df4232b05afd212fc28027dae6fd6a81618c2a3116182cb57c7f0a3fd0a")
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

	CalculateRequestPriceCallbackJuels(opts *bind.CallOpts, gasAllowance *big.Int, arguments []byte) (*big.Int, error)

	CalculateRequestPriceJuels(opts *bind.CallOpts) (*big.Int, error)

	GetBillingConfig(opts *bind.CallOpts) (VRFBeaconTypesBillingConfig, error)

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

	IRouter(opts *bind.CallOpts) (common.Address, error)

	IStartSlot(opts *bind.CallOpts) (*big.Int, error)

	MigrationVersion(opts *bind.CallOpts) (uint8, error)

	MinDelay(opts *bind.CallOpts) (uint16, error)

	OnMigration(opts *bind.CallOpts, arg0 []byte) error

	Owner(opts *bind.CallOpts) (common.Address, error)

	Paused(opts *bind.CallOpts) (bool, error)

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

	Migrate(opts *bind.TransactOpts, newCoordinator common.Address, encodedRequest []byte) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	ProcessVRFOutputs(opts *bind.TransactOpts, vrfOutputs []VRFBeaconTypesVRFOutput, juelsPerFeeCoin *big.Int, reasonableGasPrice uint64, blockHeight uint64, arg4 [32]byte) (*types.Transaction, error)

	RedeemRandomness(opts *bind.TransactOpts, sender common.Address, arg1 *big.Int, requestID *big.Int, arg3 []byte) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, requester common.Address, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, arg4 []byte) (*types.Transaction, error)

	RequestRandomnessFulfillment(opts *bind.TransactOpts, requester common.Address, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, arg6 []byte) (*types.Transaction, error)

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

	FilterRandomnessFulfillmentRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFCoordinatorRandomnessFulfillmentRequestedIterator, error)

	WatchRandomnessFulfillmentRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessFulfillmentRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error)

	ParseRandomnessFulfillmentRequested(log types.Log) (*VRFCoordinatorRandomnessFulfillmentRequested, error)

	FilterRandomnessRequested(opts *bind.FilterOpts, requestID []*big.Int, requester []common.Address) (*VRFCoordinatorRandomnessRequestedIterator, error)

	WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorRandomnessRequested, requestID []*big.Int, requester []common.Address) (event.Subscription, error)

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
