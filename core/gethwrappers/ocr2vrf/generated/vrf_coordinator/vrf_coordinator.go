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

var VRFCoordinatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocksArg\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BeaconPeriodMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"earliestAllowed\",\"type\":\"uint256\"}],\"name\":\"BlockTooRecent\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"firstDelay\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"minDelay\",\"type\":\"uint16\"}],\"name\":\"ConfirmationDelayBlocksTooShort\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confirmationDelays\",\"type\":\"uint16[10]\"},{\"internalType\":\"uint8\",\"name\":\"violatingIndex\",\"type\":\"uint8\"}],\"name\":\"ConfirmationDelaysNotIncreasing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasAllowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLeft\",\"type\":\"uint256\"}],\"name\":\"GasAllowanceExceedsGasLeft\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reportHeight\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"separatorHeight\",\"type\":\"uint64\"}],\"name\":\"HistoryDomainSeparatorTooOld\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"actualBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requiredBalance\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidBillingConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"expectedLength\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"actualLength\",\"type\":\"uint256\"}],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidJuelsConversion\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numRecipients\",\"type\":\"uint256\"}],\"name\":\"InvalidNumberOfRecipients\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestedSubID\",\"type\":\"uint256\"}],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"requestedVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"coordinatorVersion\",\"type\":\"uint8\"}],\"name\":\"MigrationVersionMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProducer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeRouter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoWordsRequested\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confDelays\",\"type\":\"uint16[10]\"}],\"name\":\"NonZeroDelayAfterZeroDelay\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnMigrationNotSupported\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"producer\",\"type\":\"address\"}],\"name\":\"ProducerAlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"}],\"name\":\"RandomnessNotAvailable\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numRecipients\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numPayments\",\"type\":\"uint256\"}],\"name\":\"RecipientsPaymentsMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"actual\",\"type\":\"address\"}],\"name\":\"ResponseMustBeRetrievedByRequester\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyRequestsReplaceContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManySlotsReplaceContract\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"max\",\"type\":\"uint256\"}],\"name\":\"TooManyWords\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockHeight\",\"type\":\"uint256\"}],\"name\":\"UniverseHasEndedBangBangBang\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"givenDelay\",\"type\":\"uint24\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"knownDelays\",\"type\":\"uint24[8]\"}],\"name\":\"UnknownConfirmationDelay\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"newVersion\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputsServed\",\"type\":\"tuple[]\"}],\"name\":\"OutputsServed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.RequestID[]\",\"name\":\"requestIDs\",\"type\":\"uint48[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"successfulFulfillment\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"truncatedErrorData\",\"type\":\"bytes[]\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAllowance\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_JUELS_SUPPLY\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"paymentsInJuels\",\"type\":\"uint256[]\"}],\"name\":\"batchTransferLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"}],\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"}],\"name\":\"calculateRequestPriceCallbackJuels\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateRequestPriceJuels\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestId\",\"type\":\"uint48\"}],\"name\":\"getCallbackMemo\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfirmationDelays\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"\",\"type\":\"uint24[8]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"pendingFulfillments\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalLinkBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_link\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_link_eth_feed\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_router\",\"outputs\":[{\"internalType\":\"contractVRFRouter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_startSlot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVRFMigration\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"encodedRequest\",\"type\":\"bytes\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrationVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minDelay\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onMigration\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"}],\"internalType\":\"structECCArithmetic.G1Point\",\"name\":\"vrfOutput\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"}],\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"},{\"internalType\":\"uint96\",\"name\":\"price\",\"type\":\"uint96\"}],\"internalType\":\"structVRFBeaconTypes.CostedCallback[]\",\"name\":\"callbacks\",\"type\":\"tuple[]\"}],\"internalType\":\"structVRFBeaconTypes.VRFOutput[]\",\"name\":\"vrfOutputs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"processVRFOutputs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"redeemRandomness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"randomness\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"requestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"name\":\"s_pendingRequests\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.SlotNumber\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_producer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"useReasonableGasPrice\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"unusedGasPenaltyPercent\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"redeemableRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callbackRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"premiumPercentage\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceStalenessBlocks\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"internalType\":\"structVRFBeaconTypes.BillingConfig\",\"name\":\"billingConfig\",\"type\":\"tuple\"}],\"name\":\"setBillingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"confDelays\",\"type\":\"uint24[8]\"}],\"name\":\"setConfirmationDelays\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"producer\",\"type\":\"address\"}],\"name\":\"setProducer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"gasPrice\",\"type\":\"uint64\"}],\"name\":\"setReasonableGasPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"juelsAmount\",\"type\":\"uint256\"}],\"name\":\"transferLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b5060405162006360380380620063608339810160408190526200003591620002cb565b82828233806000816200008f5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c257620000c28162000174565b50506001805460ff60a01b19169055506001600160a01b0392831660805290821660a0521660c05260008490036200010d57604051632abc297960e01b815260040160405180910390fd5b8360e0818152505060006200012c6200021f60201b6200317b1760201c565b9050600060e051826200014091906200031f565b905060008160e05162000154919062000358565b905062000162818462000372565b6101005250620003a795505050505050565b336001600160a01b03821603620001ce5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000086565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60004661a4b181148062000235575062066eed81145b15620002a75760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200027b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002a191906200038d565b91505090565b4391505090565b80516001600160a01b0381168114620002c657600080fd5b919050565b60008060008060808587031215620002e257600080fd5b84519350620002f460208601620002ae565b92506200030460408601620002ae565b91506200031460608601620002ae565b905092959194509250565b6000826200033d57634e487b7160e01b600052601260045260246000fd5b500690565b634e487b7160e01b600052601160045260246000fd5b6000828210156200036d576200036d62000342565b500390565b6000821982111562000388576200038862000342565b500190565b600060208284031215620003a057600080fd5b5051919050565b60805160a05160c05160e05161010051615ee86200047860003960006103ab01526000818161068101528181611e1c01528181613a4501528181613a7401528181613aac015261442f01526000818161046d01528181610b00015281816110580152818161153d015281816119ec01528181611cba015281816122f301528181612411015261469101526000818161079e0152613ea701526000818161053c015281816109f501528181611b000152818161271601528181612e8b01528181612f8e01526130e70152615ee86000f3fe608060405234801561001057600080fd5b506004361061030a5760003560e01c806385c64e111161019c578063c3fbb6fd116100ee578063dac83d2911610097578063f27fcfb811610071578063f27fcfb814610799578063f2fde38b146107c0578063f99b1d68146107d357600080fd5b8063dac83d291461074f578063dc311dd314610762578063e72f6e301461078657600080fd5b8063cd0593df116100c8578063cd0593df1461067c578063ce3f4719146106a3578063d0b925fd146106b657600080fd5b8063c3fbb6fd1461064e578063c63c4e9b14610661578063cb6317971461066957600080fd5b8063a21a23e411610150578063b2a7cac51161012a578063b2a7cac514610615578063bd58017f14610628578063bec4c08c1461063b57600080fd5b8063a21a23e4146105da578063a4c0ed36146105e2578063abbf1c9b146105f557600080fd5b80638da5cb5b116101815780638da5cb5b146105a35780638eef585f146105b457806395009f08146105c757600080fd5b806385c64e11146105665780638a8a13901461057b57600080fd5b806345ccbb8b1161026057806364d51a2a1161020957806379ba5097116101e357806379ba50971461052f5780637d253aff146105375780638456cb591461055e57600080fd5b806364d51a2a146104ee57806373433a2f14610509578063747c91f71461051c57600080fd5b8063581bdd161161023a578063581bdd16146104685780635c975abb146104a757806362f8b620146104c457600080fd5b806345ccbb8b1461042457806346942d1814610435578063483af7031461044857600080fd5b8063294daa49116102c25780632f7527cc1161029c5780632f7527cc1461040b5780633f4ba83a1461041357806340d6bb821461041b57600080fd5b8063294daa49146103db5780632b38bafc146103f05780632d9297b01461040357600080fd5b8063170b317a116102f3578063170b317a14610337578063181f5a77146103675780631a961f25146103a657600080fd5b806305f4acc61461030f5780630ae0954014610324575b600080fd5b61032261031d366004614a86565b6107e6565b005b610322610332366004614ac1565b610880565b61034a610345366004614d1e565b610bba565b6040516001600160601b0390911681526020015b60405180910390f35b604080518082018252601481527f565246436f6f7264696e61746f7220312e302e300000000000000000000000006020820152905161035e9190614d9f565b6103cd7f000000000000000000000000000000000000000000000000000000000000000081565b60405190815260200161035e565b60015b60405160ff909116815260200161035e565b6103226103fe366004614db2565b610c22565b61034a610c9e565b6103de600881565b610322610ce2565b6103cd6103e881565b6002546001600160601b03166103cd565b610322610443366004614dcf565b610cf4565b61045b610456366004614e33565b610d5b565b60405161035e9190614f37565b61048f7f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200161035e565b600154600160a01b900460ff16604051901515815260200161035e565b6104d76104d2366004614f5d565b61104b565b60405165ffffffffffff909116815260200161035e565b6104f6606481565b60405161ffff909116815260200161035e565b610322610517366004614fd7565b6112ab565b61034a6b033b2e3c9fd0803ce800000081565b6103226113ce565b61048f7f000000000000000000000000000000000000000000000000000000000000000081565b61032261147f565b61056e61148f565b60405161035e919061506a565b6103cd610589366004615079565b65ffffffffffff166000908152600b602052604090205490565b6000546001600160a01b031661048f565b6103226105c2366004615094565b6114f4565b6104d76105d53660046150d1565b611530565b6103cd61173d565b6103226105f03660046151c5565b611a84565b610608610603366004615214565b611cad565b60405161035e91906152ad565b6103226106233660046152c0565b611f52565b600a5461048f906001600160a01b031681565b610322610649366004614ac1565b61209c565b61032261065c3660046152d9565b612297565b6104f6600381565b610322610677366004614ac1565b6128e5565b6103cd7f000000000000000000000000000000000000000000000000000000000000000081565b6103226106b136600461532d565b612c09565b6107126106c4366004615079565b60106020526000908152604090205463ffffffff811690640100000000810462ffffff1690670100000000000000810461ffff1690690100000000000000000090046001600160a01b031684565b6040805163ffffffff909516855262ffffff909316602085015261ffff909116918301919091526001600160a01b0316606082015260800161035e565b61032261075d366004614ac1565b612c3b565b6107756107703660046152c0565b612d56565b60405161035e9594939291906153a7565b610322610794366004614db2565b612e52565b61048f7f000000000000000000000000000000000000000000000000000000000000000081565b6103226107ce366004614db2565b613046565b6103226107e13660046153fc565b613057565b600a546001600160a01b0316331461081157604051634bea32db60e11b815260040160405180910390fd5b60045460ff161561087d5760408051808201909152436001600160401b039081168083529083166020909201829052600780547fffffffffffffffffffffffffffffffff0000000000000000000000000000000016909117680100000000000000009092029190911790555b50565b60008281526008602052604090205482906001600160a01b0316806108c05760405163c5171ee960e01b8152600481018390526024015b60405180910390fd5b336001600160a01b038216146108f457604051636c51fda960e11b81526001600160a01b03821660048201526024016108b7565b6000848152600960205260409020548490600160a01b90046001600160401b03161561093357604051631685ecdd60e31b815260040160405180910390fd5b60065460ff16156109575760405163769dd35360e11b815260040160405180910390fd5b600085815260096020908152604091829020825160608101845290546001600160601b0381168083526001600160401b03600160601b8304811694840194909452600160a01b909104909216928101929092526109b3876131ff565b600280546001600160601b031690829060006109cf838561543e565b92506101000a8154816001600160601b0302191690836001600160601b031602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb88846001600160601b03166040518363ffffffff1660e01b8152600401610a5e9291906001600160a01b03929092168252602082015260400190565b6020604051808303816000875af1158015610a7d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610aa19190615474565b610ad15760405163cf47918160e01b81526001600160601b038083166004830152831660248201526044016108b7565b6040517fbb9b2b34000000000000000000000000000000000000000000000000000000008152600481018990527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063bb9b2b3490602401600060405180830381600087803b158015610b4c57600080fd5b505af1158015610b60573d6000803e3d6000fd5b5050604080516001600160a01b038b1681526001600160601b03861660208201528b93507f3784f77e8e883de95b5d47cd713ced01229fa74d118c0a462224bcb0516d43f192500160405180910390a25050505050505050565b600080610bc561337c565b60045460808501516001600160401b039290921691610bf8916a0100000000000000000000900463ffffffff1690615491565b610c0291906154bc565b6001600160601b03169050610c1b818460e00151613435565b9392505050565b610c2a613492565b600a546001600160a01b031615610c7c57600a546040517fea6d39050000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024016108b7565b600a80546001600160a01b0319166001600160a01b0392909216919091179055565b600080610ca961337c565b600454610cc691906601000000000000900463ffffffff166154eb565b6001600160401b03169050610cdc816000613435565b91505090565b610cea613492565b610cf26134ec565b565b610cfc613492565b6064610d0e6040830160208401615520565b60ff161115610d49576040517f0afa82a800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806004610d56828261554a565b505050565b600a546060906001600160a01b03163314610d8957604051634bea32db60e11b815260040160405180910390fd5b600154600160a01b900460ff1615610dd65760405162461bcd60e51b815260206004820152601060248201526f14185d5cd8589b194e881c185d5cd95960821b60448201526064016108b7565b600080876001600160401b03811115610df157610df1614af1565b604051908082528060200260200182016040528015610e4357816020015b604080516080810182526000808252602080830182905292820181905260608201528252600019909201910181610e0f5790505b50905060005b88811015610f335760008a8a83818110610e6557610e656156e0565b9050602002810190610e7791906156f6565b610e80906157f9565b9050610e8d81888b613592565b60408101515151151580610ea957506040810151516020015115155b15610f20576040805160808101825282516001600160401b0316815260208084015162ffffff168183015283830180515151938301939093529151519091015160608201528351849061ffff8716908110610f0657610f066156e0565b60200260200101819052508380610f1c906158ce565b9450505b5080610f2b816158ef565b915050610e49565b5060008261ffff166001600160401b03811115610f5257610f52614af1565b604051908082528060200260200182016040528015610fa457816020015b604080516080810182526000808252602080830182905292820181905260608201528252600019909201910181610f705790505b50905060005b8361ffff1681101561100057828181518110610fc857610fc86156e0565b6020026020010151828281518110610fe257610fe26156e0565b60200260200101819052508080610ff8906158ef565b915050610faa565b507ff10ea936d00579b4c52035ee33bf46929646b3aa87554c565d8fb2c7aa549c44868989846040516110369493929190615908565b60405180910390a19998505050505050505050565b6000336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146110965760405163f74c318f60e01b815260040160405180910390fd5b600154600160a01b900460ff16156110e35760405162461bcd60e51b815260206004820152601060248201526f14185d5cd8589b194e881c185d5cd95960821b60448201526064016108b7565b60065460ff16156111075760405163769dd35360e11b815260040160405180910390fd5b600080600061111789888861399a565b9250925092506111278989613cd3565b65ffffffffffff8316600090815260106020908152604080832085518154848801518489015160608a01516001600160a01b03166901000000000000000000027fffffff0000000000000000000000000000000000000000ffffffffffffffffff61ffff90921667010000000000000002919091167fffffff00000000000000000000000000000000000000000000ffffffffffffff62ffffff9093166401000000000266ffffffffffffff1990941663ffffffff909516949094179290921716919091171790558a83526009909152902080546001600160401b03600160601b9091041690600c6112188361595e565b82546101009290920a6001600160401b0381810219909316918316021790915560408051918416825262ffffff8916602083015281018a905261ffff891660608201526001600160a01b038b16915065ffffffffffff8516907fc3b31df4232b05afd212fc28027dae6fd6a81618c2a3116182cb57c7f0a3fd0a9060800160405180910390a35090979650505050505050565b600a546001600160a01b031633146112d657604051634bea32db60e11b815260040160405180910390fd5b828015806112e45750601f81115b1561131e576040517f4ecc4fef000000000000000000000000000000000000000000000000000000008152600481018290526024016108b7565b808214611361576040517f339f8a9d00000000000000000000000000000000000000000000000000000000815260048101829052602481018390526044016108b7565b60005b818110156113c6576113b4868683818110611381576113816156e0565b90506020020160208101906113969190614db2565b8585848181106113a8576113a86156e0565b90506020020135613057565b806113be816158ef565b915050611364565b505050505050565b6001546001600160a01b031633146114285760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016108b7565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611487613492565b610cf2613dd7565b611497614929565b6040805161010081019182905290600f90600890826000855b82829054906101000a900462ffffff1662ffffff16815260200190600301906020826002010492830192600103820291508084116114b05790505050505050905090565b600a546001600160a01b0316331461151f57604051634bea32db60e11b815260040160405180910390fd5b61152c600f826008614948565b5050565b6000336001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461157b5760405163f74c318f60e01b815260040160405180910390fd5b600154600160a01b900460ff16156115c85760405162461bcd60e51b815260206004820152601060248201526f14185d5cd8589b194e881c185d5cd95960821b60448201526064016108b7565b60065460ff16156115ec5760405163769dd35360e11b815260040160405180910390fd5b6000806115fa8a898961399a565b925050915060006040518061010001604052808465ffffffffffff1681526020018a61ffff1681526020018c6001600160a01b031681526020018781526020018863ffffffff166001600160601b031681526020018b815260200161165d61337c565b6001600160401b03168152602001611673613e5f565b9052905061168081613f50565b81888b83604051602001611697949392919061597a565b60408051601f19818403018152918152815160209283012065ffffffffffff86166000908152600b9093529120556116ce8a61406a565b8a6001600160a01b03168365ffffffffffff167f24f0e469e0097d1e8d9975137f9f4dd17d2c1481b3a2f25f2382f51287eda1dc848b8e8e8d8860c001518960e001518f604051611726989796959493929190615a40565b60405180910390a350909998505050505050505050565b60065460009060ff16156117645760405163769dd35360e11b815260040160405180910390fd5b600154600160a01b900460ff16156117b15760405162461bcd60e51b815260206004820152601060248201526f14185d5cd8589b194e881c185d5cd95960821b60448201526064016108b7565b6000336117bf600143615aa5565b6001546040516bffffffffffffffffffffffff19606094851b81166020830152924060348201523090931b90911660548301527fffffffffffffffff000000000000000000000000000000000000000000000000600160a81b90910460c01b16606882015260700160408051808303601f19018152919052805160209091012060018054919250600160a81b9091046001600160401b03169060156118638361595e565b91906101000a8154816001600160401b0302191690836001600160401b03160217905550506000806001600160401b038111156118a2576118a2614af1565b6040519080825280602002602001820160405280156118cb578160200160208202803683370190505b506040805160608082018352600080835260208084018281528486018381528984526009835286842095518654925191516001600160601b039091166001600160a01b031993841617600160601b6001600160401b0393841602177fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff16600160a01b9290911691909102179094558451928301855233835282810182815283860187815289845260088352959092208351815486166001600160a01b0391821617825592516001820180549096169316929092179093559251805194955090936119bb92600285019201906149e6565b50506040517f9479b74e000000000000000000000000000000000000000000000000000000008152600481018490527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169150639479b74e90602401600060405180830381600087803b158015611a3957600080fd5b505af1158015611a4d573d6000803e3d6000fd5b50506040513392508491507f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d90600090a350905090565b60065460ff1615611aa85760405163769dd35360e11b815260040160405180910390fd5b600154600160a01b900460ff1615611af55760405162461bcd60e51b815260206004820152601060248201526f14185d5cd8589b194e881c185d5cd95960821b60448201526064016108b7565b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614611b57576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114611b8257604051636865567560e01b815260206004820152602481018290526044016108b7565b6000611b90828401846152c0565b6000818152600860205260409020549091506001600160a01b0316611bcb5760405163c5171ee960e01b8152600481018290526024016108b7565b600081815260096020526040812080546001600160601b031691869190611bf28385615491565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600260008282829054906101000a90046001600160601b0316611c3a9190615491565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a828784611c8d9190615abc565b6040805192835260208301919091520160405180910390a2505050505050565b6060336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614611cf85760405163f74c318f60e01b815260040160405180910390fd5b60065460ff1615611d1c5760405163769dd35360e11b815260040160405180910390fd5b65ffffffffffff831660008181526010602081815260408084208151608081018352815463ffffffff8116825262ffffff6401000000008204168286015261ffff670100000000000000820416938201939093526001600160a01b036901000000000000000000840481166060830190815297909652939092527fffffff00000000000000000000000000000000000000000000000000000000001690559151811690871614611e115760608101516040517f8e30e8230000000000000000000000000000000000000000000000000000000081526001600160a01b03918216600482015290871660248201526044016108b7565b8051600090611e47907f00000000000000000000000000000000000000000000000000000000000000009063ffffffff16615ad4565b90506000611e5361317b565b90506000836020015162ffffff1682611e6c9190615aa5565b9050808310611ed15782846020015162ffffff1684611e8b9190615abc565b611e96906001615abc565b6040517f15ad27c3000000000000000000000000000000000000000000000000000000008152600481019290925260248201526044016108b7565b6001600160401b03831115611f15576040517f058ddf02000000000000000000000000000000000000000000000000000000008152600481018490526024016108b7565b6000838152600d602090815260408083208288015162ffffff168452909152902054611f459088908690866140f7565b9998505050505050505050565b60065460ff1615611f765760405163769dd35360e11b815260040160405180910390fd5b6000818152600860205260409020546001600160a01b0316611fae5760405163c5171ee960e01b8152600481018290526024016108b7565b6000818152600860205260409020600101546001600160a01b0316331461202057600081815260086020526040908190206001015490517fd084e9750000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024016108b7565b6000818152600860209081526040918290208054336001600160a01b0319808316821784556001909301805490931690925583516001600160a01b0390911680825292810191909152909183917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c9386910160405180910390a25050565b60008281526008602052604090205482906001600160a01b0316806120d75760405163c5171ee960e01b8152600481018390526024016108b7565b336001600160a01b0382161461210b57604051636c51fda960e11b81526001600160a01b03821660048201526024016108b7565b60065460ff161561212f5760405163769dd35360e11b815260040160405180910390fd5b600154600160a01b900460ff161561217c5760405162461bcd60e51b815260206004820152601060248201526f14185d5cd8589b194e881c185d5cd95960821b60448201526064016108b7565b600084815260086020526040902060020154606319016121c8576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b03831660009081526003602090815260408083208784529091528120546001600160401b03169003612291576001600160a01b03831660008181526003602090815260408083208884528252808320805467ffffffffffffffff19166001908117909155600883528184206002018054918201815584529282902090920180546001600160a01b03191684179055905191825285917f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e191015b60405180910390a25b50505050565b60065460ff16156122bb5760405163769dd35360e11b815260040160405180910390fd5b6040517fcef40aa60000000000000000000000000000000000000000000000000000000081526001600160a01b0384811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063cef40aa690602401602060405180830381865afa15801561233a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061235e9190615474565b612394576040517f3fae142000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604081146123c15760408051636865567560e01b81526004810191909152602481018290526044016108b7565b60006123cf82840184615af3565b60208101516040517ffa66358a000000000000000000000000000000000000000000000000000000008152600481019190915290915030906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063fa66358a90602401602060405180830381865afa158015612458573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061247c9190615b2d565b6001600160a01b0316146124ab57806020015160405163c5171ee960e01b81526004016108b791815260200190565b60008060008060006124c08660200151612d56565b94509450945094509450816001600160a01b0316336001600160a01b03161461250757604051636c51fda960e11b81526001600160a01b03831660048201526024016108b7565b886001600160a01b031663294daa496040518163ffffffff1660e01b8152600401602060405180830381865afa158015612545573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125699190615b4a565b60ff16866000015160ff161461261f578560000151896001600160a01b031663294daa496040518163ffffffff1660e01b8152600401602060405180830381865afa1580156125bc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125e09190615b4a565b6040517fe7aada9500000000000000000000000000000000000000000000000000000000815260ff9283166004820152911660248201526044016108b7565b6001600160401b0383161561264757604051631685ecdd60e31b815260040160405180910390fd5b60006040518060c0016040528061265c600190565b60ff16815260200188602001518152602001846001600160a01b03168152602001838152602001876001600160601b03168152602001866001600160401b031681525090506000816040516020016126b49190615b67565b60405160208183030381529060405290506126d288602001516131ff565b600280548891906000906126f09084906001600160601b031661543e565b92506101000a8154816001600160601b0302191690836001600160601b031602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb8c896040518363ffffffff1660e01b815260040161277f9291906001600160a01b039290921682526001600160601b0316602082015260400190565b6020604051808303816000875af115801561279e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127c29190615474565b61280e5760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e6473000000000000000000000000000060448201526064016108b7565b6040517fce3f47190000000000000000000000000000000000000000000000000000000081526001600160a01b038c169063ce3f471990612853908490600401614d9f565b600060405180830381600087803b15801561286d57600080fd5b505af1158015612881573d6000803e3d6000fd5b505050508760200151886000015160ff167fbd89b747474d3fc04664dfbd1d56ae7ffbe46ee097cdb9979c13916bb76269ce8d6040516128d091906001600160a01b0391909116815260200190565b60405180910390a35050505050505050505050565b60008281526008602052604090205482906001600160a01b0316806129205760405163c5171ee960e01b8152600481018390526024016108b7565b336001600160a01b0382161461295457604051636c51fda960e11b81526001600160a01b03821660048201526024016108b7565b6000848152600960205260409020548490600160a01b90046001600160401b03161561299357604051631685ecdd60e31b815260040160405180910390fd5b60065460ff16156129b75760405163769dd35360e11b815260040160405180910390fd5b6001600160a01b03841660009081526003602090815260408083208884529091528120546001600160401b03169003612a15576040516379bfd40160e01b8152600481018690526001600160a01b03851660248201526044016108b7565b600085815260086020908152604080832060020180548251818502810185019093528083529192909190830182828015612a7857602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612a5a575b50505050509050600060018251612a8f9190615aa5565b905060005b8251811015612b9a57866001600160a01b0316838281518110612ab957612ab96156e0565b60200260200101516001600160a01b031603612b88576000838381518110612ae357612ae36156e0565b6020026020010151905080600860008b81526020019081526020016000206002018381548110612b1557612b156156e0565b600091825260208083209190910180546001600160a01b0319166001600160a01b0394909416939093179092558a8152600890915260409020600201805480612b6057612b60615bdf565b600082815260209020810160001990810180546001600160a01b031916905501905550612b9a565b80612b92816158ef565b915050612a94565b506001600160a01b03861660008181526003602090815260408083208b8452825291829020805467ffffffffffffffff19169055905191825288917f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a7910160405180910390a250505050505050565b6040517f2cb6686f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526008602052604090205482906001600160a01b031680612c765760405163c5171ee960e01b8152600481018390526024016108b7565b336001600160a01b03821614612caa57604051636c51fda960e11b81526001600160a01b03821660048201526024016108b7565b60065460ff1615612cce5760405163769dd35360e11b815260040160405180910390fd5b6000848152600860205260409020600101546001600160a01b038481169116146122915760008481526008602090815260409182902060010180546001600160a01b0319166001600160a01b03871690811790915582513381529182015285917f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a19101612288565b6000818152600860205260408120548190819081906060906001600160a01b0316612d975760405163c5171ee960e01b8152600481018790526024016108b7565b60008681526009602090815260408083205460088352928190208054600290910180548351818602810186019094528084526001600160601b038616956001600160401b03600160601b8204811696600160a01b90920416946001600160a01b0390941693918391830182828015612e3857602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612e1a575b505050505090509450945094509450945091939590929450565b612e5a613492565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906370a0823190602401602060405180830381865afa158015612eda573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612efe9190615bf5565b6002549091506001600160601b031681811115612f51576040517fa99da30200000000000000000000000000000000000000000000000000000000815260048101829052602481018390526044016108b7565b81811015610d56576000612f658284615aa5565b60405163a9059cbb60e01b81526001600160a01b038681166004830152602482018390529192507f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb906044016020604051808303816000875af1158015612fd9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ffd9190615474565b50604080516001600160a01b0386168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a150505050565b61304e613492565b61087d816142fc565b600a546001600160a01b0316331461308257604051634bea32db60e11b815260040160405180910390fd5b600280546001600160601b0316908290600061309e838561543e565b82546001600160601b039182166101009390930a92830291909202199091161790555060405163a9059cbb60e01b81526001600160a01b038481166004830152602482018490527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af1158015613130573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906131549190615474565b610d565760405163cf47918160e01b815260048101829052602481018390526044016108b7565b60004661a4b1811480613190575062066eed81145b156131f85760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156131d4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cdc9190615bf5565b4391505090565b6000818152600860209081526040808320815160608101835281546001600160a01b0390811682526001830154168185015260028201805484518187028101870186528181529295939486019383018282801561328557602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311613267575b505050505081525050905060005b81604001515181101561330b5760036000836040015183815181106132ba576132ba6156e0565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081208682529092529020805467ffffffffffffffff1916905580613303816158ef565b915050613293565b50600082815260086020526040812080546001600160a01b031990811682556001820180549091169055906133436002830182614a3b565b505050600090815260096020526040902080547fffffffff00000000000000000000000000000000000000000000000000000000169055565b60045460009060ff1680156133a757506007546801000000000000000090046001600160401b031615155b15613430576004547201000000000000000000000000000000000000900463ffffffff164310808061340c57506004546133fd907201000000000000000000000000000000000000900463ffffffff1643615aa5565b6007546001600160401b031610155b1561342e5750506007546801000000000000000090046001600160401b031690565b505b503a90565b6004546000908190606490613462906e010000000000000000000000000000900463ffffffff1682615c0e565b6134729063ffffffff1686615ad4565b61347c9190615c43565b905061348881846143a5565b9150505b92915050565b6000546001600160a01b03163314610cf25760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016108b7565b600154600160a01b900460ff166135455760405162461bcd60e51b815260206004820152601460248201527f5061757361626c653a206e6f742070617573656400000000000000000000000060448201526064016108b7565b6001805460ff60a01b191690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b6040516001600160a01b03909116815260200160405180910390a1565b82516001600160401b03808416911611156135ef5782516040517f012d824d0000000000000000000000000000000000000000000000000000000081526001600160401b03808516600483015290911660248201526044016108b7565b6040830151515160009015801561360d575060408401515160200151155b15613645575082516001600160401b03166000908152600d602090815260408083208287015162ffffff16845290915290205461369f565b836040015160405160200161365a9190615c57565b60408051601f19818403018152918152815160209283012086516001600160401b03166000908152600d84528281208885015162ffffff168252909352912081905590505b6060840151516000816001600160401b038111156136bf576136bf614af1565b6040519080825280602002602001820160405280156136e8578160200160208202803683370190505b5090506000826001600160401b0381111561370557613705614af1565b6040519080825280601f01601f19166020018201604052801561372f576020820181803683370190505b5090506000836001600160401b0381111561374c5761374c614af1565b60405190808252806020026020018201604052801561377f57816020015b606081526020019060019003908161376a5790505b5090506000805b858110156138985760008a6060015182815181106137a6576137a66156e0565b602002602001015190506000806137c78d600001518e602001518c86614425565b9150915081156138065780868661ffff16815181106137e8576137e86156e0565b602002602001018190525084806137fe906158ce565b95505061384d565b600160f81b87858151811061381d5761381d6156e0565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505b8251518851899086908110613864576138646156e0565b602002602001019065ffffffffffff16908165ffffffffffff16815250505050508080613890906158ef565b915050613786565b506060890151511561398f5760008161ffff166001600160401b038111156138c2576138c2614af1565b6040519080825280602002602001820160405280156138f557816020015b60608152602001906001900390816138e05790505b50905060005b8261ffff1681101561395157838181518110613919576139196156e0565b6020026020010151828281518110613933576139336156e0565b60200260200101819052508080613949906158ef565b9150506138fb565b507f47ddf7bb0cbd94c1b43c5097f1352a80db0ceb3696f029d32b24f32cd631d2b785858360405161398593929190615c8a565b60405180910390a1505b505050505050505050565b604080516080810182526000808252602082018190529181018290526060810182905260006103e88561ffff1611156139f457604051634a90778560e01b815261ffff861660048201526103e860248201526044016108b7565b8461ffff16600003613a32576040517f08fad2a700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000613a3c61317b565b90506000613a6a7f000000000000000000000000000000000000000000000000000000000000000083615d40565b9050600081613a997f000000000000000000000000000000000000000000000000000000000000000085615abc565b613aa39190615aa5565b90506000613ad17f000000000000000000000000000000000000000000000000000000000000000083615c43565b905063ffffffff8110613b10576040517f7b2a523000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080518082018252600e805465ffffffffffff1682528251610100810193849052849360009392916020840191600f906008908288855b82829054906101000a900462ffffff1662ffffff1681526020019060030190602082600201049283019260010382029150808411613b4857905050505091909252505081519192505065ffffffffffff90811610613bd2576040517f2b4655b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051613bdf906001615d54565b600e805465ffffffffffff191665ffffffffffff9290921691909117905560005b6008811015613c46578a62ffffff1682602001518260088110613c2557613c256156e0565b602002015162ffffff1614613c465780613c3e816158ef565b915050613c00565b60088110613c875760208201516040517fc4f769b00000000000000000000000000000000000000000000000000000000081526108b7918d91600401615d75565b506040805160808101825263ffffffff909316835262ffffff8b16602084015261ffff8c16908301526001600160a01b038c166060830152519750955090935050505093509350939050565b6001600160a01b03821660009081526003602090815260408083208484529091528120546001600160401b03169003613d31576040516379bfd40160e01b8152600481018290526001600160a01b03831660248201526044016108b7565b6000613d3b610c9e565b600083815260096020526040902080546001600160601b039283169350909116821115613d9057805460405163cf47918160e01b81526001600160601b039091166004820152602481018390526044016108b7565b805482908290600090613dad9084906001600160601b031661543e565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555050505050565b600154600160a01b900460ff1615613e245760405162461bcd60e51b815260206004820152601060248201526f14185d5cd8589b194e881c185d5cd95960821b60448201526064016108b7565b6001805460ff60a01b1916600160a01b1790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586135753390565b60048054604080517ffeaf968c00000000000000000000000000000000000000000000000000000000815290516000936201000090930463ffffffff169283151592859283927f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169263feaf968c928183019260a0928290030181865afa158015613ef6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613f1a9190615da9565b509450909250849150508015613f3e5750613f358242615aa5565b8463ffffffff16105b15613f4857506005545b949350505050565b6040808201516001600160a01b031660009081526003602090815282822060a085015183529052908120546001600160401b03169003613fbf5760a081015160408083015190516379bfd40160e01b815260048101929092526001600160a01b031660248201526044016108b7565b6000613fca82610bba565b60a0830151600090815260096020526040902080546001600160601b03928316935090911682111561402457805460405163cf47918160e01b81526001600160601b039091166004820152602481018390526044016108b7565b8054829082906000906140419084906001600160601b031661543e565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550505050565b60008181526009602052604090208054600160601b90046001600160401b031681600c6140968361595e565b82546101009290920a6001600160401b038181021990931691831602179091558254600160a01b90041690508160146140ce8361595e565b91906101000a8154816001600160401b0302191690836001600160401b03160217905550505050565b606082614149576040517fc7d41b1b00000000000000000000000000000000000000000000000000000000815265ffffffffffff861660048201526001600160401b03831660248201526044016108b7565b6040805165ffffffffffff8716602080830191909152865163ffffffff168284015286015162ffffff166060808301919091529186015161ffff166080820152908501516001600160a01b031660a082015260c0810184905260009060e0016040516020818303038152906040528051906020012090506103e8856040015161ffff1611156141ff576040808601519051634a90778560e01b815261ffff90911660048201526103e860248201526044016108b7565b6000856040015161ffff166001600160401b0381111561422157614221614af1565b60405190808252806020026020018201604052801561424a578160200160208202803683370190505b50905060005b866040015161ffff168161ffff1610156142f15782816040516020016142a592919091825260f01b7fffff00000000000000000000000000000000000000000000000000000000000016602082015260220190565b6040516020818303038152906040528051906020012060001c828261ffff16815181106142d4576142d46156e0565b6020908102919091010152806142e9816158ce565b915050614250565b509695505050505050565b336001600160a01b038216036143545760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016108b7565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008082156143b457826143bc565b6143bc613e5f565b90506000816143d386670de0b6b3a7640000615ad4565b6143dd9190615c43565b90506b033b2e3c9fd0803ce8000000811115613488576040517fde43710000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060608161445d7f00000000000000000000000000000000000000000000000000000000000000006001600160401b038916615c43565b845160a08101516040519293509091600091614481918b918b91869060200161597a565b60408051601f198184030181529181528151602092830120845165ffffffffffff166000908152600b90935291205490915081146144fd5760016040518060400160405280601081526020017f756e6b6e6f776e2063616c6c6261636b000000000000000000000000000000008152509450945050505061486a565b6040805160808101825263ffffffff8516815262ffffff8a1660208083019190915284015161ffff1681830152908301516001600160a01b03166060820152876145865760016040518060400160405280601681526020017f756e617661696c61626c652072616e646f6d6e65737300000000000000000000815250955095505050505061486a565b60006145988460000151838b8e6140f7565b60608084015186519187015160405193945090926000927f5a47dd7100000000000000000000000000000000000000000000000000000000926145e092879190602401615df9565b60408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526006805460ff191660011790559050600080805a8d51608001516040808c015190517f55fe97630000000000000000000000000000000000000000000000000000000081529293506001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016926355fe9763926146c89290918990600401615e2c565b60408051808303816000875af11580156146e6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061470a9190615e66565b90935091508161475f578c51608001516040517faad159830000000000000000000000000000000000000000000000000000000081526001600160601b039091166004820152602481018290526044016108b7565b60006113885a61476f9190615abc565b6006805460ff19169055905081811015614797576147976147908284615aa5565b8f51614873565b8d5160a0015160009081526009602052604090208054600160a01b90046001600160401b03169060146147c983615e95565b82546001600160401b039182166101009390930a928302919092021990911617905550895165ffffffffffff166000908152600b6020526040812055836148475760016040518060400160405280601081526020017f657865637574696f6e206661696c65640000000000000000000000000000000081525061485a565b6000604051806020016040528060008152505b9c509c5050505050505050505050505b94509492505050565b80608001516001600160601b031682111561488c575050565b6004546000906064906148a790610100900460ff1682615eb8565b60ff168360c001518585608001516001600160601b03166148c89190615aa5565b6148d29190615ad4565b6148dc9190615ad4565b6148e69190615c43565b905060006148f8828460e001516143a5565b60a0840151600090815260096020526040812080549293508392909190613dad9084906001600160601b0316615491565b6040518061010001604052806008906020820280368337509192915050565b6001830191839082156149d65791602002820160005b838211156149a557833562ffffff1683826101000a81548162ffffff021916908362ffffff160217905550926020019260030160208160020104928301926001030261495e565b80156149d45782816101000a81549062ffffff02191690556003016020816002010492830192600103026149a5565b505b506149e2929150614a55565b5090565b8280548282559060005260206000209081019282156149d6579160200282015b828111156149d657825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190614a06565b508054600082559060005260206000209081019061087d91905b5b808211156149e25760008155600101614a56565b80356001600160401b0381168114614a8157600080fd5b919050565b600060208284031215614a9857600080fd5b610c1b82614a6a565b6001600160a01b038116811461087d57600080fd5b8035614a8181614aa1565b60008060408385031215614ad457600080fd5b823591506020830135614ae681614aa1565b809150509250929050565b634e487b7160e01b600052604160045260246000fd5b60405161010081016001600160401b0381118282101715614b2a57614b2a614af1565b60405290565b604080519081016001600160401b0381118282101715614b2a57614b2a614af1565b604051608081016001600160401b0381118282101715614b2a57614b2a614af1565b604051602081016001600160401b0381118282101715614b2a57614b2a614af1565b604051601f8201601f191681016001600160401b0381118282101715614bbe57614bbe614af1565b604052919050565b803565ffffffffffff81168114614a8157600080fd5b803561ffff81168114614a8157600080fd5b600082601f830112614bff57600080fd5b81356001600160401b03811115614c1857614c18614af1565b614c2b601f8201601f1916602001614b96565b818152846020838601011115614c4057600080fd5b816020850160208301376000918101602001919091529392505050565b80356001600160601b0381168114614a8157600080fd5b60006101008284031215614c8757600080fd5b614c8f614b07565b9050614c9a82614bc6565b8152614ca860208301614bdc565b6020820152614cb960408301614ab6565b604082015260608201356001600160401b03811115614cd757600080fd5b614ce384828501614bee565b606083015250614cf560808301614c5d565b608082015260a082013560a082015260c082013560c082015260e082013560e082015292915050565b600060208284031215614d3057600080fd5b81356001600160401b03811115614d4657600080fd5b61348884828501614c74565b6000815180845260005b81811015614d7857602081850181015186830182015201614d5c565b81811115614d8a576000602083870101525b50601f01601f19169290920160200192915050565b602081526000610c1b6020830184614d52565b600060208284031215614dc457600080fd5b8135610c1b81614aa1565b60006101008284031215614de257600080fd5b50919050565b60008083601f840112614dfa57600080fd5b5081356001600160401b03811115614e1157600080fd5b6020830191508360208260051b8501011115614e2c57600080fd5b9250929050565b60008060008060008060a08789031215614e4c57600080fd5b86356001600160401b03811115614e6257600080fd5b614e6e89828a01614de8565b909750955050602087013577ffffffffffffffffffffffffffffffffffffffffffffffff81168114614e9f57600080fd5b9350614ead60408801614a6a565b9250614ebb60608801614a6a565b9150608087013590509295509295509295565b600081518084526020808501945080840160005b83811015614f2c57815180516001600160401b031688528381015162ffffff1684890152604080820151908901526060908101519088015260809096019590820190600101614ee2565b509495945050505050565b602081526000610c1b6020830184614ece565b803562ffffff81168114614a8157600080fd5b600080600080600060a08688031215614f7557600080fd5b8535614f8081614aa1565b945060208601359350614f9560408701614bdc565b9250614fa360608701614f4a565b915060808601356001600160401b03811115614fbe57600080fd5b614fca88828901614bee565b9150509295509295909350565b60008060008060408587031215614fed57600080fd5b84356001600160401b038082111561500457600080fd5b61501088838901614de8565b9096509450602087013591508082111561502957600080fd5b5061503687828801614de8565b95989497509550505050565b8060005b600881101561229157815162ffffff16845260209384019390910190600101615046565b610100810161348c8284615042565b60006020828403121561508b57600080fd5b610c1b82614bc6565b60006101008083850312156150a857600080fd5b8381840111156150b757600080fd5b509092915050565b63ffffffff8116811461087d57600080fd5b600080600080600080600060e0888a0312156150ec57600080fd5b87356150f781614aa1565b96506020880135955061510c60408901614bdc565b945061511a60608901614f4a565b9350608088013561512a816150bf565b925060a08801356001600160401b038082111561514657600080fd5b6151528b838c01614bee565b935060c08a013591508082111561516857600080fd5b506151758a828b01614bee565b91505092959891949750929550565b60008083601f84011261519657600080fd5b5081356001600160401b038111156151ad57600080fd5b602083019150836020828501011115614e2c57600080fd5b600080600080606085870312156151db57600080fd5b84356151e681614aa1565b93506020850135925060408501356001600160401b0381111561520857600080fd5b61503687828801615184565b6000806000806080858703121561522a57600080fd5b843561523581614aa1565b93506020850135925061524a60408601614bc6565b915060608501356001600160401b0381111561526557600080fd5b61527187828801614bee565b91505092959194509250565b600081518084526020808501945080840160005b83811015614f2c57815187529582019590820190600101615291565b602081526000610c1b602083018461527d565b6000602082840312156152d257600080fd5b5035919050565b6000806000604084860312156152ee57600080fd5b83356152f981614aa1565b925060208401356001600160401b0381111561531457600080fd5b61532086828701615184565b9497909650939450505050565b6000806020838503121561534057600080fd5b82356001600160401b0381111561535657600080fd5b61536285828601615184565b90969095509350505050565b600081518084526020808501945080840160005b83811015614f2c5781516001600160a01b031687529582019590820190600101615382565b6001600160601b038616815260006001600160401b0380871660208401528086166040840152506001600160a01b038416606083015260a060808301526153f160a083018461536e565b979650505050505050565b6000806040838503121561540f57600080fd5b823561541a81614aa1565b946020939093013593505050565b634e487b7160e01b600052601160045260246000fd5b60006001600160601b038381169083168181101561545e5761545e615428565b039392505050565b801515811461087d57600080fd5b60006020828403121561548657600080fd5b8151610c1b81615466565b60006001600160601b038083168185168083038211156154b3576154b3615428565b01949350505050565b60006001600160601b03808316818516818304811182151516156154e2576154e2615428565b02949350505050565b60006001600160401b03808316818516818304811182151516156154e2576154e2615428565b60ff8116811461087d57600080fd5b60006020828403121561553257600080fd5b8135610c1b81615511565b6000813561348c816150bf565b813561555581615466565b815460ff19811691151560ff169182178355602084013561557581615511565b61ff008160081b168361ffff198416171784555050506155b861559a6040840161553d565b825465ffffffff0000191660109190911b65ffffffff000016178255565b6155ed6155c76060840161553d565b825469ffffffff000000000000191660309190911b69ffffffff00000000000016178255565b61562a6155fc6080840161553d565b82546dffffffff00000000000000000000191660509190911b6dffffffff0000000000000000000016178255565b61567c61563960a0840161553d565b82547fffffffffffffffffffffffffffff00000000ffffffffffffffffffffffffffff1660709190911b71ffffffff000000000000000000000000000016178255565b6156d261568b60c0840161553d565b82547fffffffffffffffffffff00000000ffffffffffffffffffffffffffffffffffff1660909190911b75ffffffff00000000000000000000000000000000000016178255565b60e082013560018201555050565b634e487b7160e01b600052603260045260246000fd5b60008235609e1983360301811261570c57600080fd5b9190910192915050565b600082601f83011261572757600080fd5b813560206001600160401b038083111561574357615743614af1565b8260051b615752838201614b96565b938452858101830193838101908886111561576c57600080fd5b84880192505b858310156157ed5782358481111561578a5760008081fd5b88016040818b03601f19018113156157a25760008081fd5b6157aa614b30565b87830135878111156157bc5760008081fd5b6157ca8d8a83870101614c74565b8252506157d8828401614c5d565b81890152845250509184019190840190615772565b98975050505050505050565b600081360360a081121561580c57600080fd5b615814614b52565b61581d84614a6a565b8152602061582c818601614f4a565b828201526040603f198401121561584257600080fd5b61584a614b74565b925036605f86011261585b57600080fd5b615863614b30565b80608087013681111561587557600080fd5b604088015b81811015615891578035845292840192840161587a565b50908552604084019490945250509035906001600160401b038211156158b657600080fd5b6158c236838601615716565b60608201529392505050565b600061ffff8083168181036158e5576158e5615428565b6001019392505050565b60006001820161590157615901615428565b5060010190565b60006001600160401b03808716835277ffffffffffffffffffffffffffffffffffffffffffffffff86166020840152808516604084015250608060608301526159546080830184614ece565b9695505050505050565b60006001600160401b038083168181036158e5576158e5615428565b6001600160401b038516815262ffffff841660208201528260408201526080606082015265ffffffffffff8251166080820152600060208301516159c460a084018261ffff169052565b5060408301516001600160a01b031660c0830152606083015161010060e084018190526159f5610180850183614d52565b91506080850151615a10828601826001600160601b03169052565b505060a084015161012084015260c084015161014084015260e08401516101608401528091505095945050505050565b60006101006001600160401b038b16835262ffffff8a16602084015288604084015261ffff8816606084015263ffffffff871660808401528560a08401528460c08401528060e0840152615a9681840185614d52565b9b9a5050505050505050505050565b600082821015615ab757615ab7615428565b500390565b60008219821115615acf57615acf615428565b500190565b6000816000190483118215151615615aee57615aee615428565b500290565b600060408284031215615b0557600080fd5b615b0d614b30565b8235615b1881615511565b81526020928301359281019290925250919050565b600060208284031215615b3f57600080fd5b8151610c1b81614aa1565b600060208284031215615b5c57600080fd5b8151610c1b81615511565b6020815260ff8251166020820152602082015160408201526001600160a01b0360408301511660608201526000606083015160c06080840152615bad60e084018261536e565b90506001600160601b0360808501511660a08401526001600160401b0360a08501511660c08401528091505092915050565b634e487b7160e01b600052603160045260246000fd5b600060208284031215615c0757600080fd5b5051919050565b600063ffffffff8083168185168083038211156154b3576154b3615428565b634e487b7160e01b600052601260045260246000fd5b600082615c5257615c52615c2d565b500490565b815160408201908260005b6002811015615c81578251825260209283019290910190600101615c62565b50505092915050565b606080825284519082018190526000906020906080840190828801845b82811015615ccb57815165ffffffffffff1684529284019290840190600101615ca7565b50505083810382850152615cdf8187614d52565b905083810360408501528085518083528383019150838160051b84010184880160005b83811015615d3057601f19868403018552615d1e838351614d52565b94870194925090860190600101615d02565b50909a9950505050505050505050565b600082615d4f57615d4f615c2d565b500690565b600065ffffffffffff8083168185168083038211156154b3576154b3615428565b62ffffff831681526101208101610c1b6020830184615042565b805169ffffffffffffffffffff81168114614a8157600080fd5b600080600080600060a08688031215615dc157600080fd5b615dca86615d8f565b9450602086015193506040860151925060608601519150615ded60808701615d8f565b90509295509295909350565b65ffffffffffff84168152606060208201526000615e1a606083018561527d565b82810360408401526159548185614d52565b6001600160601b03841681526001600160a01b0383166020820152606060408201526000615e5d6060830184614d52565b95945050505050565b60008060408385031215615e7957600080fd5b8251615e8481615466565b6020840151909250614ae681615466565b60006001600160401b03821680615eae57615eae615428565b6000190192915050565b600060ff821660ff841680821015615ed257615ed2615428565b9003939250505056fea164736f6c634300080f000a",
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

func (_VRFCoordinator *VRFCoordinatorCaller) CalculateRequestPriceCallbackJuels(opts *bind.CallOpts, callback VRFBeaconTypesCallback) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "calculateRequestPriceCallbackJuels", callback)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) CalculateRequestPriceCallbackJuels(callback VRFBeaconTypesCallback) (*big.Int, error) {
	return _VRFCoordinator.Contract.CalculateRequestPriceCallbackJuels(&_VRFCoordinator.CallOpts, callback)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) CalculateRequestPriceCallbackJuels(callback VRFBeaconTypesCallback) (*big.Int, error) {
	return _VRFCoordinator.Contract.CalculateRequestPriceCallbackJuels(&_VRFCoordinator.CallOpts, callback)
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

func (_VRFCoordinator *VRFCoordinatorCaller) GetTotalLinkBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "getTotalLinkBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) GetTotalLinkBalance() (*big.Int, error) {
	return _VRFCoordinator.Contract.GetTotalLinkBalance(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) GetTotalLinkBalance() (*big.Int, error) {
	return _VRFCoordinator.Contract.GetTotalLinkBalance(&_VRFCoordinator.CallOpts)
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

func (_VRFCoordinator *VRFCoordinatorTransactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "recoverFunds", to)
}

func (_VRFCoordinator *VRFCoordinatorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RecoverFunds(&_VRFCoordinator.TransactOpts, to)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RecoverFunds(&_VRFCoordinator.TransactOpts, to)
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

type VRFCoordinatorFundsRecoveredIterator struct {
	Event *VRFCoordinatorFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorFundsRecovered)
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
		it.Event = new(VRFCoordinatorFundsRecovered)
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

func (it *VRFCoordinatorFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorFundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorFundsRecoveredIterator{contract: _VRFCoordinator.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorFundsRecovered)
				if err := _VRFCoordinator.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseFundsRecovered(log types.Log) (*VRFCoordinatorFundsRecovered, error) {
	event := new(VRFCoordinatorFundsRecovered)
	if err := _VRFCoordinator.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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
type SPendingRequests struct {
	SlotNumber        uint32
	ConfirmationDelay *big.Int
	NumWords          uint16
	Requester         common.Address
}

func (_VRFCoordinator *VRFCoordinator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinator.abi.Events["ConfigSet"].ID:
		return _VRFCoordinator.ParseConfigSet(log)
	case _VRFCoordinator.abi.Events["FundsRecovered"].ID:
		return _VRFCoordinator.ParseFundsRecovered(log)
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

func (VRFCoordinatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
}

func (VRFCoordinatorFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
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

	CalculateRequestPriceCallbackJuels(opts *bind.CallOpts, callback VRFBeaconTypesCallback) (*big.Int, error)

	CalculateRequestPriceJuels(opts *bind.CallOpts) (*big.Int, error)

	GetCallbackMemo(opts *bind.CallOpts, requestId *big.Int) ([32]byte, error)

	GetConfirmationDelays(opts *bind.CallOpts) ([8]*big.Int, error)

	GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

		error)

	GetTotalLinkBalance(opts *bind.CallOpts) (*big.Int, error)

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

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RedeemRandomness(opts *bind.TransactOpts, sender common.Address, arg1 *big.Int, requestID *big.Int, arg3 []byte) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, requester common.Address, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, arg4 []byte) (*types.Transaction, error)

	RequestRandomnessFulfillment(opts *bind.TransactOpts, requester common.Address, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, arg6 []byte) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error)

	SetBillingConfig(opts *bind.TransactOpts, billingConfig VRFBeaconTypesBillingConfig) (*types.Transaction, error)

	SetConfirmationDelays(opts *bind.TransactOpts, confDelays [8]*big.Int) (*types.Transaction, error)

	SetProducer(opts *bind.TransactOpts, producer common.Address) (*types.Transaction, error)

	SetReasonableGasPrice(opts *bind.TransactOpts, gasPrice uint64) (*types.Transaction, error)

	TransferLink(opts *bind.TransactOpts, recipient common.Address, juelsAmount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFCoordinatorConfigSet, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorFundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorFundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*VRFCoordinatorFundsRecovered, error)

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
