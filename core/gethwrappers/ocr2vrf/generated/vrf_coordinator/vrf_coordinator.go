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

type VRFBeaconTypesCoordinatorConfig struct {
	UseReasonableGasPrice             bool
	ReentrancyLock                    bool
	Paused                            bool
	PremiumPercentage                 uint8
	UnusedGasPenaltyPercent           uint8
	StalenessSeconds                  uint32
	RedeemableRequestGasOverhead      uint32
	CallbackRequestGasOverhead        uint32
	ReasonableGasPriceStalenessBlocks uint32
	FallbackWeiPerUnitLink            *big.Int
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
	ShouldStore       bool
}

type VRFCoordinatorCallbackConfig struct {
	MaxCallbackGasLimit        uint32
	MaxCallbackArgumentsLength uint32
}

var VRFCoordinatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocksArg\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BeaconPeriodMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"earliestAllowed\",\"type\":\"uint256\"}],\"name\":\"BlockTooRecent\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confirmationDelays\",\"type\":\"uint16[10]\"},{\"internalType\":\"uint8\",\"name\":\"violatingIndex\",\"type\":\"uint8\"}],\"name\":\"ConfirmationDelaysNotIncreasing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ContractPaused\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasAllowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLeft\",\"type\":\"uint256\"}],\"name\":\"GasAllowanceExceedsGasLeft\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reportHeight\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"separatorHeight\",\"type\":\"uint64\"}],\"name\":\"HistoryDomainSeparatorTooOld\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"actualBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requiredBalance\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"expectedLength\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"actualLength\",\"type\":\"uint256\"}],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCoordinatorConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidJuelsConversion\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numRecipients\",\"type\":\"uint256\"}],\"name\":\"InvalidNumberOfRecipients\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestedSubID\",\"type\":\"uint256\"}],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"requestedVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"coordinatorVersion\",\"type\":\"uint8\"}],\"name\":\"MigrationVersionMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProducer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativePaymentGiven\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoWordsRequested\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confDelays\",\"type\":\"uint16[10]\"}],\"name\":\"NonZeroDelayAfterZeroDelay\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnMigrationNotSupported\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"producer\",\"type\":\"address\"}],\"name\":\"ProducerAlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"}],\"name\":\"RandomnessNotAvailable\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"confDelay\",\"type\":\"uint256\"}],\"name\":\"RandomnessSeedNotFoundForCallbacks\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numRecipients\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numPayments\",\"type\":\"uint256\"}],\"name\":\"RecipientsPaymentsMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"actual\",\"type\":\"address\"}],\"name\":\"ResponseMustBeRetrievedByRequester\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyRequestsReplaceContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManySlotsReplaceContract\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"max\",\"type\":\"uint256\"}],\"name\":\"TooManyWords\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockHeight\",\"type\":\"uint256\"}],\"name\":\"UniverseHasEndedBangBangBang\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCallbackArgumentsLength\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structVRFCoordinator.CallbackConfig\",\"name\":\"newConfig\",\"type\":\"tuple\"}],\"name\":\"CallbackConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"useReasonableGasPrice\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reentrancyLock\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"premiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"unusedGasPenaltyPercent\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"redeemableRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callbackRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceStalenessBlocks\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"uint96\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.CoordinatorConfig\",\"name\":\"coordinatorConfig\",\"type\":\"tuple\"}],\"name\":\"CoordinatorConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"newVersion\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"uint24\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputsServed\",\"type\":\"tuple[]\"}],\"name\":\"OutputsServed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"name\":\"PauseFlagChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"requestIDs\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"successfulFulfillment\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"truncatedErrorData\",\"type\":\"bytes[]\"},{\"indexed\":false,\"internalType\":\"uint96[]\",\"name\":\"subBalances\",\"type\":\"uint96[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"subIDs\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAllowance\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"costJuels\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newSubBalance\",\"type\":\"uint256\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"}],\"name\":\"RandomnessRedeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"costJuels\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newSubBalance\",\"type\":\"uint256\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"paymentsInJuels\",\"type\":\"uint256[]\"}],\"name\":\"batchTransferLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"deregisterMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"getCallbackMemo\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfirmationDelays\",\"outputs\":[{\"internalType\":\"uint24[8]\",\"name\":\"\",\"type\":\"uint24[8]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"getFulfillmentFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"pendingFulfillments\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"totalRequests\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"totalSuccessfulFulfillments\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSubscriptionLinkBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_link\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVRFMigration\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"encodedRequest\",\"type\":\"bytes\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrationVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onMigration\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"uint24\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"}],\"internalType\":\"structECCArithmetic.G1Point\",\"name\":\"vrfOutput\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"}],\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"},{\"internalType\":\"uint96\",\"name\":\"price\",\"type\":\"uint96\"}],\"internalType\":\"structVRFBeaconTypes.CostedCallback[]\",\"name\":\"callbacks\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"shouldStore\",\"type\":\"bool\"}],\"internalType\":\"structVRFBeaconTypes.VRFOutput[]\",\"name\":\"vrfOutputs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"}],\"name\":\"processVRFOutputs\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"containsNewOutputs\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"redeemRandomness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"randomness\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"registerMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint24\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"uint24\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"requestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_callbackConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCallbackArgumentsLength\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_coordinatorConfig\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"useReasonableGasPrice\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reentrancyLock\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"premiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"unusedGasPenaltyPercent\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"redeemableRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callbackRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceStalenessBlocks\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_pendingRequests\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_producer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"maxCallbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxCallbackArgumentsLength\",\"type\":\"uint32\"}],\"internalType\":\"structVRFCoordinator.CallbackConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"setCallbackConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint24[8]\",\"name\":\"confDelays\",\"type\":\"uint24[8]\"}],\"name\":\"setConfirmationDelays\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"useReasonableGasPrice\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reentrancyLock\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"premiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"unusedGasPenaltyPercent\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"redeemableRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callbackRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceStalenessBlocks\",\"type\":\"uint32\"},{\"internalType\":\"uint96\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"uint96\"}],\"internalType\":\"structVRFBeaconTypes.CoordinatorConfig\",\"name\":\"coordinatorConfig\",\"type\":\"tuple\"}],\"name\":\"setCoordinatorConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"pause\",\"type\":\"bool\"}],\"name\":\"setPauseFlag\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"producer\",\"type\":\"address\"}],\"name\":\"setProducer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"juelsAmount\",\"type\":\"uint256\"}],\"name\":\"transferLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b50604051620064e8380380620064e8833981016040819052620000349162000239565b8033806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf816200018e565b5050506001600160a01b03166080908152604080519182018152600080835260208301819052908201819052662386f26fc10000606090920191909152642386f26fc160b01b6006556004805463ffffffff60281b191668ffffffff00000000001790558290036200014457604051632abc297960e01b815260040160405180910390fd5b60a0829052600e805465ffffffffffff16906000620001638362000278565b91906101000a81548165ffffffffffff021916908365ffffffffffff160217905550505050620002ac565b336001600160a01b03821603620001e85760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080604083850312156200024d57600080fd5b825160208401519092506001600160a01b03811681146200026d57600080fd5b809150509250929050565b600065ffffffffffff808316818103620002a257634e487b7160e01b600052601160045260246000fd5b6001019392505050565b60805160a0516161d0620003186000396000818161075f01528181611d1701528181613d9e01528181613dcd01528181613e05015261448401526000818161047501528181610bd801528181611ab5015281816123ee01528181612df00152612e8501526161d06000f3fe6080604052600436106101d65760003560e01c806304104edb146101db5780630ae09540146101fd57806316f6ee9a1461021d578063294daa491461025d5780632b38bafc1461027f5780632f7527cc1461029f5780633e79167f146102b457806340d6bb82146102d457806347c3e2cb146102ea5780634ffac83a14610385578063597d2f3c146103985780635d06b4ab146103b657806364d51a2a146103d657806373433a2f146103fe57806376f2e3f41461041e57806379ba50971461044e5780637d253aff1461046357806385c64e11146104a45780638c7cba66146104c65780638da5cb5b146104e65780638da92e71146105045780638eef585f146105245780639e20103614610544578063a21a23e414610564578063a4c0ed3614610579578063acfc6cdd14610599578063b2a7cac5146105c6578063b79fa6f7146105e6578063bd58017f146106cd578063bec4c08c146106ed578063c3fbb6fd1461070d578063cb6317971461072d578063cd0593df1461074d578063ce3f471914610781578063dac83d29146107a1578063db972c8b146107c1578063dc311dd3146107d4578063e30afa4a14610806578063f2fde38b1461084b578063f99b1d681461086b578063f9c45ced1461088b575b600080fd5b3480156101e757600080fd5b506101fb6101f6366004614c7e565b6108ab565b005b34801561020957600080fd5b506101fb610218366004614ca2565b610a60565b34801561022957600080fd5b5061024a610238366004614cd2565b6000908152600c602052604090205490565b6040519081526020015b60405180910390f35b34801561026957600080fd5b5060015b60405160ff9091168152602001610254565b34801561028b57600080fd5b506101fb61029a366004614c7e565b610ce3565b3480156102ab57600080fd5b5061026d600881565b3480156102c057600080fd5b506101fb6102cf366004614ceb565b610d44565b3480156102e057600080fd5b5061024a6103e881565b3480156102f657600080fd5b50610348610305366004614cd2565b60106020526000908152604090205463ffffffff811690600160201b810462ffffff1690600160381b810461ffff1690600160481b90046001600160a01b031684565b6040805163ffffffff909516855262ffffff909316602085015261ffff909116918301919091526001600160a01b03166060820152608001610254565b61024a610393366004614e6d565b610df5565b3480156103a457600080fd5b506002546001600160601b031661024a565b3480156103c257600080fd5b506101fb6103d1366004614c7e565b610fa1565b3480156103e257600080fd5b506103eb606481565b60405161ffff9091168152602001610254565b34801561040a57600080fd5b506101fb610419366004614f18565b61104d565b34801561042a57600080fd5b5061043e610439366004614f9a565b611138565b6040519015158152602001610254565b34801561045a57600080fd5b506101fb611476565b34801561046f57600080fd5b506104977f000000000000000000000000000000000000000000000000000000000000000081565b604051610254919061501c565b3480156104b057600080fd5b506104b9611520565b6040516102549190615030565b3480156104d257600080fd5b506101fb6104e1366004615084565b611585565b3480156104f257600080fd5b506000546001600160a01b0316610497565b34801561051057600080fd5b506101fb61051f3660046150de565b6115f9565b34801561053057600080fd5b506101fb61053f3660046150fb565b611663565b34801561055057600080fd5b5061024a61055f366004615126565b61169f565b34801561057057600080fd5b5061024a6117c0565b34801561058557600080fd5b506101fb6105943660046151da565b611a57565b3480156105a557600080fd5b506105b96105b4366004615229565b611c42565b60405161025491906152b3565b3480156105d257600080fd5b506101fb6105e1366004614cd2565b611e46565b3480156105f257600080fd5b506004546005546106629160ff80821692610100830482169262010000810483169263010000008204811692600160201b83049091169163ffffffff600160281b8204811692600160481b8304821692600160681b8104831692600160881b90910416906001600160601b03168a565b604080519a15158b5298151560208b01529615159789019790975260ff948516606089015292909316608087015263ffffffff90811660a087015291821660c0860152811660e08501529091166101008301526001600160601b031661012082015261014001610254565b3480156106d957600080fd5b50600a54610497906001600160a01b031681565b3480156106f957600080fd5b506101fb610708366004614ca2565b611f77565b34801561071957600080fd5b506101fb6107283660046152c6565b612133565b34801561073957600080fd5b506101fb610748366004614ca2565b61261f565b34801561075957600080fd5b5061024a7f000000000000000000000000000000000000000000000000000000000000000081565b34801561078d57600080fd5b506101fb61079c36600461531a565b612909565b3480156107ad57600080fd5b506101fb6107bc366004614ca2565b612922565b61024a6107cf36600461535b565b612a33565b3480156107e057600080fd5b506107f46107ef366004614cd2565b612c91565b60405161025496959493929190615433565b34801561081257600080fd5b50600b5461082e9063ffffffff80821691600160201b90041682565b6040805163ffffffff938416815292909116602083015201610254565b34801561085757600080fd5b506101fb610866366004614c7e565b612d9d565b34801561087757600080fd5b506101fb610886366004615489565b612dae565b34801561089757600080fd5b5061024a6108a63660046154b5565b612f17565b6108b361302e565b60095460005b81811015610a3857826001600160a01b0316600982815481106108de576108de6154fb565b6000918252602090912001546001600160a01b031603610a26576009610905600184615527565b81548110610915576109156154fb565b600091825260209091200154600980546001600160a01b039092169183908110610941576109416154fb565b600091825260209091200180546001600160a01b0319166001600160a01b0392909216919091179055826009610978600185615527565b81548110610988576109886154fb565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060098054806109c7576109c761553a565b600082815260209020810160001990810180546001600160a01b03191690550190556040517ff80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af3790610a1990859061501c565b60405180910390a1505050565b80610a3081615550565b9150506108b9565b5081604051635428d44960e01b8152600401610a54919061501c565b60405180910390fd5b50565b60008281526007602052604090205482906001600160a01b031680610a9b5760405163c5171ee960e01b815260048101839052602401610a54565b336001600160a01b03821614610ac65780604051636c51fda960e11b8152600401610a54919061501c565b600454610100900460ff1615610aef5760405163769dd35360e11b815260040160405180910390fd5b600084815260086020526040902054600160601b900463ffffffff1615610b2957604051631685ecdd60e31b815260040160405180910390fd5b600084815260086020908152604091829020825160808101845290546001600160601b03811680835263ffffffff600160601b830416938301939093526001600160401b03600160801b8204811694830194909452600160c01b90049092166060830152610b9686613083565b600280546001600160601b03169082906000610bb28385615569565b92506101000a8154816001600160601b0302191690836001600160601b031602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb87846001600160601b03166040518363ffffffff1660e01b8152600401610c2d929190615590565b6020604051808303816000875af1158015610c4c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c7091906155a9565b610ca05760405163cf47918160e01b81526001600160601b03808316600483015283166024820152604401610a54565b867f3784f77e8e883de95b5d47cd713ced01229fa74d118c0a462224bcb0516d43f18784604051610cd29291906155c6565b60405180910390a250505050505050565b610ceb61302e565b600a546001600160a01b031615610d2257600a5460405163ea6d390560e01b8152610a54916001600160a01b03169060040161501c565b600a80546001600160a01b0319166001600160a01b0392909216919091179055565b610d4c61302e565b6064610d5e60a0830160808401615602565b60ff161180610d785750610d7860408201602083016150de565b80610d8e5750610d8e60608201604083016150de565b15610dac5760405163b0e7bd8360e01b815260040160405180910390fd5b806004610db98282615668565b9050507e28d3a46e95e67def989d41c66eb331add9809460b95b5fb4eb006157728fc581604051610dea919061583b565b60405180910390a150565b60045460009062010000900460ff1615610e225760405163ab35696f60e01b815260040160405180910390fd5b600454610100900460ff1615610e4b5760405163769dd35360e11b815260040160405180910390fd5b3415610e6c57604051630b829bad60e21b8152346004820152602401610a54565b6000806000610e7d883389896131c5565b925092509250600080610e90338b613310565b600087815260106020908152604091829020885181548a8401518b8601516060808e015163ffffffff90951666ffffffffffffff1990941693909317600160201b62ffffff9384160217600160381b600160e81b031916600160381b61ffff92831602600160481b600160e81b03191617600160481b6001600160a01b03909516949094029390931790935584513381526001600160401b038b1694810194909452918e169383019390935281018e9052908c16608082015260a081018390526001600160601b03821660c0820152919350915085907fb7933fba96b6b452eb44f99fdc08052a45dff82363d59abaff0456931c3d24599060e00160405180910390a2509298975050505050505050565b610fa961302e565b610fb28161351f565b15610fd2578060405163ac8a27ef60e01b8152600401610a54919061501c565b600980546001810182556000919091527f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af0180546001600160a01b0319166001600160a01b0383161790556040517fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af0162590610dea90839061501c565b600a546001600160a01b0316331461107857604051634bea32db60e11b815260040160405180910390fd5b828015806110865750601f81115b156110a757604051634ecc4fef60e01b815260048101829052602401610a54565b8082146110cb5760405163339f8a9d60e01b8152610a549082908490600401615923565b60005b818110156111305761111e8686838181106110eb576110eb6154fb565b90506020020160208101906111009190614c7e565b858584818110611112576111126154fb565b90506020020135612dae565b8061112881615550565b9150506110ce565b505050505050565b600a546000906001600160a01b0316331461116657604051634bea32db60e11b815260040160405180910390fd5b60045462010000900460ff16156111905760405163ab35696f60e01b815260040160405180910390fd5b6001600160c01b038416156111cd57600680546001600160601b038616600160a01b02600160201b600160a01b0390911663ffffffff4216171790555b6001600160401b038316156112225760068054436001600160401b03908116600160201b02600160201b600160601b0319918716600160601b0291909116600160201b600160a01b0319909216919091171790555b600080866001600160401b0381111561123d5761123d614d29565b60405190808252806020026020018201604052801561127657816020015b611263614ab7565b81526020019060019003908161125b5790505b50905060005b87811015611378576000898983818110611298576112986154fb565b90506020028101906112aa9190615931565b6112b390615ab5565b905060006112c282888b613588565b905085806112cd5750805b604083015151519096501515806112ec57506040820151516020015115155b15611363576040805160808101825283516001600160401b0316815260208085015162ffffff168183015284830180515151938301939093529151519091015160608201528451859061ffff8816908110611349576113496154fb565b6020026020010181905250848061135f90615b9b565b9550505b5050808061137090615550565b91505061127c565b5060008261ffff166001600160401b0381111561139757611397614d29565b6040519080825280602002602001820160405280156113d057816020015b6113bd614ab7565b8152602001906001900390816113b55790505b50905060005b8361ffff1681101561142c578281815181106113f4576113f46154fb565b602002602001015182828151811061140e5761140e6154fb565b6020026020010181905250808061142490615550565b9150506113d6565b507ff10ea936d00579b4c52035ee33bf46929646b3aa87554c565d8fb2c7aa549c44858888846040516114629493929190615bbc565b60405180910390a150505095945050505050565b6001546001600160a01b031633146114c95760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b6044820152606401610a54565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611528614aed565b6040805161010081019182905290600f90600890826000855b82829054906101000a900462ffffff1662ffffff16815260200190600301906020826002010492830192600103820291508084116115415790505050505050905090565b61158d61302e565b8051600b80546020808501805163ffffffff908116600160201b026001600160401b031990941695811695861793909317909355604080519485529251909116908301527f0cc54509a45ab33cd67614d4a2892c083ecf8fb43b9d29f6ea8130b9023e51df9101610dea565b61160161302e565b60045460ff6201000090910416151581151514610a5d5760048054821515620100000262ff0000199091161790556040517f49ba7c1de2d8853088b6270e43df2118516b217f38b917dd2b80dea360860fbe90610dea90831515815260200190565b600a546001600160a01b0316331461168e57604051634bea32db60e11b815260040160405180910390fd5b61169b600f826008614b0c565b5050565b604080516101408101825260045460ff80821615158352610100808304821615156020808601919091526201000084048316151585870152630100000084048316606080870191909152600160201b808604909416608080880191909152600160281b860463ffffffff90811660a0890152600160481b8704811660c0890152600160681b8704811660e0890152600160881b9096048616938701939093526005546001600160601b039081166101208801528751938401885260065480871685526001600160401b03958104861693850193909352600160601b830490941696830196909652600160a01b90049091169381019390935260009283926117ab92881691879190613720565b50506001600160601b03169695505050505050565b600454600090610100900460ff16156117ec5760405163769dd35360e11b815260040160405180910390fd5b60045462010000900460ff16156118165760405163ab35696f60e01b815260040160405180910390fd5b600033611824600143615527565b6001546040516001600160601b0319606094851b81166020830152924060348201523090931b90911660548301526001600160c01b0319600160a01b90910460c01b16606882015260700160408051808303601f19018152919052805160209091012060018054919250600160a01b9091046001600160401b03169060146118ab83615c51565b91906101000a8154816001600160401b0302191690836001600160401b03160217905550506000806001600160401b038111156118ea576118ea614d29565b604051908082528060200260200182016040528015611913578160200160208202803683370190505b50604080516080810182526000808252602080830182815283850183815260608086018581528a865260088552878620965187549451935191516001600160601b039091166001600160801b031990951694909417600160601b63ffffffff90941693909302929092176001600160801b0316600160801b6001600160401b03938416026001600160c01b031617600160c01b929093169190910291909117909355835192830184523383528281018281528385018681528884526007835294909220835181546001600160a01b03199081166001600160a01b0392831617835593516001830180549095169116179092559251805194955091939092611a21926002850192910190614baa565b505060405133915083907f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d90600090a350905090565b600454610100900460ff1615611a805760405163769dd35360e11b815260040160405180910390fd5b60045462010000900460ff1615611aaa5760405163ab35696f60e01b815260040160405180910390fd5b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614611af3576040516344b0e3c360e01b815260040160405180910390fd5b60208114611b1957604051636865567560e01b8152610a54906020908390600401615c75565b6000611b2782840184614cd2565b6000818152600760205260409020549091506001600160a01b0316611b625760405163c5171ee960e01b815260048101829052602401610a54565b600081815260086020526040812080546001600160601b031691869190611b898385615c89565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600260008282829054906101000a90046001600160601b0316611bd19190615c89565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a828784611c249190615ca9565b604051611c32929190615923565b60405180910390a2505050505050565b600454606090610100900460ff1615611c6e5760405163769dd35360e11b815260040160405180910390fd5b60008381526010602081815260408084208151608081018352815463ffffffff8116825262ffffff600160201b8204168286015261ffff600160381b820416938201939093526001600160a01b03600160481b8404811660608301908152968a9052949093526001600160e81b031990911690559151163314611d0c57806060015133604051638e30e82360e01b8152600401610a54929190615cbc565b8051600090611d42907f00000000000000000000000000000000000000000000000000000000000000009063ffffffff16615cd6565b90506000611d4e6137ca565b90506000836020015162ffffff1682611d679190615527565b9050808310611dac5782846020015162ffffff1684611d869190615ca9565b611d91906001615ca9565b6040516315ad27c360e01b8152600401610a54929190615923565b6001600160401b03831115611dd7576040516302c6ef8160e11b815260048101849052602401610a54565b604051888152339088907f16f3f633197fafab10a5df69e6f3f2f7f20092f08d8d47de0a91c0f4b96a1a259060200160405180910390a3611e3a8785600d6000611e25888a60200151613854565b81526020019081526020016000205486613863565b98975050505050505050565b600454610100900460ff1615611e6f5760405163769dd35360e11b815260040160405180910390fd5b6000818152600760205260409020546001600160a01b0316611ea75760405163c5171ee960e01b815260048101829052602401610a54565b6000818152600760205260409020600101546001600160a01b03163314611efe576000818152600760205260409081902060010154905163d084e97560e01b8152610a54916001600160a01b03169060040161501c565b6000818152600760205260409081902080546001600160a01b031980821633908117845560019093018054909116905591516001600160a01b039092169183917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c938691611f6b918591615cbc565b60405180910390a25050565b60008281526007602052604090205482906001600160a01b031680611fb25760405163c5171ee960e01b815260048101839052602401610a54565b336001600160a01b03821614611fdd5780604051636c51fda960e11b8152600401610a54919061501c565b600454610100900460ff16156120065760405163769dd35360e11b815260040160405180910390fd5b60045462010000900460ff16156120305760405163ab35696f60e01b815260040160405180910390fd5b60008481526007602052604090206002015460631901612063576040516305a48e0f60e01b815260040160405180910390fd5b600360006120718587613a1c565b815260208101919091526040016000205460ff1661212d576001600360006120998688613a1c565b815260208082019290925260409081016000908120805460ff191694151594909417909355868352600782528083206002018054600181018255908452919092200180546001600160a01b0319166001600160a01b0386161790555184907f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e19061212490869061501c565b60405180910390a25b50505050565b600454610100900460ff161561215c5760405163769dd35360e11b815260040160405180910390fd5b6121658361351f565b6121845782604051635428d44960e01b8152600401610a54919061501c565b604081146121a95760408051636865567560e01b8152610a5491908390600401615c75565b60006121b782840184615ced565b90506000806000806121cc8560200151612c91565b95509550505093509350816001600160a01b0316336001600160a01b03161461220a5781604051636c51fda960e11b8152600401610a54919061501c565b876001600160a01b031663294daa496040518163ffffffff1660e01b8152600401602060405180830381865afa158015612248573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061226c9190615d27565b60ff16856000015160ff1614612309578460000151886001600160a01b031663294daa496040518163ffffffff1660e01b8152600401602060405180830381865afa1580156122bf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122e39190615d27565b60405163e7aada9560e01b815260ff928316600482015291166024820152604401610a54565b63ffffffff83161561232e57604051631685ecdd60e31b815260040160405180910390fd5b60006040518060a00160405280612343600190565b60ff16815260200187602001518152602001846001600160a01b03168152602001838152602001866001600160601b0316815250905060008160405160200161238c9190615d44565b60405160208183030381529060405290506123aa8760200151613083565b600280548791906000906123c89084906001600160601b0316615569565b92506101000a8154816001600160601b0302191690836001600160601b031602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb8b886040518363ffffffff1660e01b815260040161243a9291906155c6565b6020604051808303816000875af1158015612459573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061247d91906155a9565b6124be5760405162461bcd60e51b8152602060048201526012602482015271696e73756666696369656e742066756e647360701b6044820152606401610a54565b60405163ce3f471960e01b81526001600160a01b038b169063ce3f4719906124ea908490600401615df0565b600060405180830381600087803b15801561250457600080fd5b505af1158015612518573d6000803e3d6000fd5b50506004805461ff00191661010017905550600090505b83518110156125c25783818151811061254a5761254a6154fb565b60200260200101516001600160a01b0316638ea981178c6040518263ffffffff1660e01b815260040161257d919061501c565b600060405180830381600087803b15801561259757600080fd5b505af11580156125ab573d6000803e3d6000fd5b5050505080806125ba90615550565b91505061252f565b506004805461ff00191690556020870151875160405160ff909116907fbd89b747474d3fc04664dfbd1d56ae7ffbe46ee097cdb9979c13916bb76269ce9061260b908e9061501c565b60405180910390a350505050505050505050565b60008281526007602052604090205482906001600160a01b03168061265a5760405163c5171ee960e01b815260048101839052602401610a54565b336001600160a01b038216146126855780604051636c51fda960e11b8152600401610a54919061501c565b600454610100900460ff16156126ae5760405163769dd35360e11b815260040160405180910390fd5b600084815260086020526040902054600160601b900463ffffffff16156126e857604051631685ecdd60e31b815260040160405180910390fd5b600360006126f68587613a1c565b815260208101919091526040016000205460ff1661272b5783836040516379bfd40160e01b8152600401610a54929190615e03565b60008481526007602090815260408083206002018054825181850281018501909352808352919290919083018282801561278e57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612770575b505050505090506000600182516127a59190615527565b905060005b82518110156128b057856001600160a01b03168382815181106127cf576127cf6154fb565b60200260200101516001600160a01b03160361289e5760008383815181106127f9576127f96154fb565b6020026020010151905080600760008a8152602001908152602001600020600201838154811061282b5761282b6154fb565b600091825260208083209190910180546001600160a01b0319166001600160a01b0394909416939093179092558981526007909152604090206002018054806128765761287661553a565b600082815260209020810160001990810180546001600160a01b0319169055019055506128b0565b806128a881615550565b9150506127aa565b50600360006128bf8789613a1c565b815260208101919091526040908101600020805460ff191690555186907f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a790611c3290889061501c565b604051632cb6686f60e01b815260040160405180910390fd5b60008281526007602052604090205482906001600160a01b03168061295d5760405163c5171ee960e01b815260048101839052602401610a54565b336001600160a01b038216146129885780604051636c51fda960e11b8152600401610a54919061501c565b600454610100900460ff16156129b15760405163769dd35360e11b815260040160405180910390fd5b6000848152600760205260409020600101546001600160a01b0384811691161461212d576000848152600760205260409081902060010180546001600160a01b0319166001600160a01b0386161790555184907f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a1906121249033908790615cbc565b60045460009062010000900460ff1615612a605760405163ab35696f60e01b815260040160405180910390fd5b600454610100900460ff1615612a895760405163769dd35360e11b815260040160405180910390fd5b3415612aaa57604051630b829bad60e21b8152346004820152602401610a54565b600080612ab989338a8a6131c5565b925050915060006040518061010001604052808481526020018a61ffff168152602001336001600160a01b031681526020018781526020018863ffffffff166001600160601b031681526020018b81526020016000815260200160008152509050600080612b2683613a32565b60c087019190915260e08601919091526040519193509150612b529085908c908f908790602001615e1a565b60405160208183030381529060405280519060200120600c6000878152602001908152602001600020819055506000604051806101600160405280878152602001336001600160a01b03168152602001866001600160401b031681526020018c62ffffff1681526020018e81526020018d61ffff1681526020018b63ffffffff1681526020018581526020018a8152602001848152602001836001600160601b0316815250905080600001517f01872fb9c7d6d68af06a17347935e04412da302a377224c205e672c26e18c37f82602001518360400151846060015185608001518660a001518760c001518860e0015160c001518960e0015160e001518a61010001518b61012001518c6101400151604051612c789b9a99989796959493929190615ecd565b60405180910390a250939b9a5050505050505050505050565b60008181526007602052604081205481908190819081906060906001600160a01b0316612cd45760405163c5171ee960e01b815260048101889052602401610a54565b60008781526008602090815260408083205460078352928190208054600290910180548351818602810186019094528084526001600160601b0386169563ffffffff600160601b820416956001600160401b03600160801b8304811696600160c01b90930416946001600160a01b031693928391830182828015612d8157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612d63575b5050505050905095509550955095509550955091939550919395565b612da561302e565b610a5d81613c9f565b600a546001600160a01b03163314612dd957604051634bea32db60e11b815260040160405180910390fd5b60405163a9059cbb60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb90612e279085908590600401615590565b6020604051808303816000875af1158015612e46573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e6a91906155a9565b61169b576040516370a0823160e01b81526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906370a0823190612eba90309060040161501c565b602060405180830381865afa158015612ed7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612efb9190615f63565b8160405163cf47918160e01b8152600401610a54929190615923565b604080516101408101825260045460ff80821615158352610100808304821615156020808601919091526201000084048316151585870152630100000084048316606080870191909152600160201b80860490941660808088019190915263ffffffff600160281b8704811660a0890152600160481b8704811660c0890152600160681b8704811660e0890152600160881b9096048616938701939093526005546001600160601b039081166101208801528751938401885260065495861684526001600160401b03948604851692840192909252600160601b850490931695820195909552600160a01b9092049093169281019290925260009161301c9190613d42565b6001600160601b031690505b92915050565b6000546001600160a01b031633146130815760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610a54565b565b6000818152600760209081526040808320815160608101835281546001600160a01b0390811682526001830154168185015260028201805484518187028101870186528181529295939486019383018282801561310957602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116130eb575b505050505081525050905060005b816040015151811015613179576003600061314f84604001518481518110613141576131416154fb565b602002602001015186613a1c565b81526020810191909152604001600020805460ff191690558061317181615550565b915050613117565b50600082815260076020526040812080546001600160a01b031990811682556001820180549091169055906131b16002830182614bff565b505050600090815260086020526040812055565b60006131cf614c1d565b60006103e88561ffff1611156131fe57846103e8604051634a90778560e01b8152600401610a54929190615c75565b8461ffff16600003613223576040516308fad2a760e01b815260040160405180910390fd5b60008061322e613d88565b600e54919350915065ffffffffffff1660006132998b8b84604080513060208201529081018490526001600160a01b038316606082015265ffffffffffff8216608082015260009060a00160408051601f198184030181529190528051602090910120949350505050565b90506132a6826001615f7c565b600e805465ffffffffffff9290921665ffffffffffff199092169190911790556040805160808101825263ffffffff909416845262ffffff8916602085015261ffff8a16908401526001600160a01b038a1660608401529550909350909150509450945094915050565b604080516080808201835260065463ffffffff8082168452600160201b8083046001600160401b03908116602080880191909152600160601b850490911686880152600160a01b9093046001600160601b0390811660608088019190915287516101408101895260045460ff808216151583526101008083048216151598840198909852620100008204811615159a83019a909a52630100000081048a169282019290925292810490971694820194909452600160281b8604821660a0820152600160481b8604821660c0820152600160681b8604821660e0820152600160881b9095041690840152600554166101208301526000918291906003836134168888613a1c565b815260208101919091526040016000205460ff1661344b5784866040516379bfd40160e01b8152600401610a54929190615e03565b60006134578284613d42565b600087815260086020526040902080546001600160601b039283169350909116828110156134a657815460405163cf47918160e01b8152610a54916001600160601b0316908590600401615f9b565b81546001600160601b03918490038281166001600160601b031960016001600160401b03600160801b80870482169290920116028116600163ffffffff60601b01600160c01b031990941693909317179093556002805480841686900390931692909116919091179055909450925050505b9250929050565b6000805b60095481101561357f57826001600160a01b03166009828154811061354a5761354a6154fb565b6000918252602090912001546001600160a01b03160361356d5750600192915050565b8061357781615550565b915050613523565b50600092915050565b6000826001600160401b031684600001516001600160401b031611156135d757835160405163012d824d60e01b81526001600160401b0380861660048301529091166024820152604401610a54565b60608401515160408086015190516000916135f491602001615fb4565b604051602081830303815290604052805190602001209050856040015160000151600060028110613627576136276154fb565b6020020151158015613640575060408601515160200151155b1561367a57600d600061366488600001516001600160401b03168960200151613854565b81526020019081526020016000205490506136fd565b8560800151156136fd576000600d60006136a589600001516001600160401b03168a60200151613854565b81526020810191909152604001600020549050806136f75781600d60006136dd8a600001516001600160401b03168b60200151613854565b8152602081019190915260400160002055600193506136fb565b8091505b505b600061370a838389613e5b565b905083806137155750805b979650505050505050565b60008060008061373086866142ff565b6001600160401b031690506000613748826010615cd6565b905060006014613759836015615cd6565b6137639190615ff4565b895161376f9190615cd6565b838960e0015163ffffffff168c6137869190615c89565b6001600160601b03166137999190615cd6565b6137a39190615ca9565b90506000806137b58360008c8c614378565b909d909c50949a509398505050505050505050565b60004661a4b18114806137df575062066eed81145b1561384d5760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613823573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138479190615f63565b91505090565b4391505090565b62ffffff1660189190911b1790565b6060826138955760405163220a34e960e11b8152600481018690526001600160401b0383166024820152604401610a54565b604080516020808201889052865163ffffffff168284015286015162ffffff166060808301919091529186015161ffff166080820152908501516001600160a01b031660a082015260c0810184905260009060e0016040516020818303038152906040528051906020012090506103e8856040015161ffff1611156139375784604001516103e8604051634a90778560e01b8152600401610a54929190615c75565b6000856040015161ffff166001600160401b0381111561395957613959614d29565b604051908082528060200260200182016040528015613982578160200160208202803683370190505b50905060005b866040015161ffff168161ffff161015613a115782816040516020016139c592919091825260f01b6001600160f01b031916602082015260220190565b6040516020818303038152906040528051906020012060001c828261ffff16815181106139f4576139f46154fb565b602090810291909101015280613a0981615b9b565b915050613988565b509695505050505050565b60a081901b6001600160a01b0383161792915050565b60008060008060036000613a4e87604001518860a00151613a1c565b815260208101919091526040016000205460ff16613a8b578460a0015185604001516040516379bfd40160e01b8152600401610a54929190615e03565b604080516080808201835260065463ffffffff80821684526001600160401b03600160201b8084048216602080880191909152600160601b8504909216868801526001600160601b03600160a01b909404841660608088019190915287516101408101895260045460ff808216151583526101008083048216151596840196909652620100008204811615159a83019a909a52630100000081048a168284015292830490981688870152600160281b8204841660a0890152600160481b8204841660c0890152600160681b8204841660e0890152600160881b90910490921690860152600554909116610120850152908801519088015191929160009182918291613b97918688613720565b60a08d0151600090815260086020526040902080546001600160601b0394851697509295509093509116841115613bef57805460405163cf47918160e01b8152610a54916001600160601b0316908690600401615f9b565b80546001600160601b0363ffffffff600160601b60016001600160401b03600160801b8087048216830190911602600160801b600160c01b03198616811783900484169091019092160263ffffffff60601b19909116600160601b600160c01b031990931692909217919091178181168690038083166001600160601b03199283161790935560028054808416889003909316929091169190911790559298509096509450925050509193509193565b336001600160a01b03821603613cf15760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610a54565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080613d4f84846142ff565b8460c0015163ffffffff16613d649190616008565b6001600160401b031690506000613d7e8260008787614378565b5095945050505050565b6000806000613d956137ca565b90506000613dc37f00000000000000000000000000000000000000000000000000000000000000008361602b565b9050600081613df27f000000000000000000000000000000000000000000000000000000000000000085615ca9565b613dfc9190615527565b90506000613e2a7f000000000000000000000000000000000000000000000000000000000000000083615ff4565b905063ffffffff8110613e50576040516307b2a52360e41b815260040160405180910390fd5b909590945092505050565b6000806040518060e00160405280866001600160401b03811115613e8157613e81614d29565b604051908082528060200260200182016040528015613eaa578160200160208202803683370190505b508152602001866001600160401b03811115613ec857613ec8614d29565b6040519080825280601f01601f191660200182016040528015613ef2576020820181803683370190505b508152602001866001600160401b03811115613f1057613f10614d29565b604051908082528060200260200182016040528015613f4357816020015b6060815260200190600190039081613f2e5790505b50815260006020820152604001866001600160401b03811115613f6857613f68614d29565b604051908082528060200260200182016040528015613f91578160200160208202803683370190505b508152602001866001600160401b03811115613faf57613faf614d29565b604051908082528060200260200182016040528015613fd8578160200160208202803683370190505b508152602001866001600160401b03811115613ff657613ff6614d29565b60405190808252806020026020018201604052801561401f578160200160208202803683370190505b509052905060005b858110156141da57600084606001518281518110614047576140476154fb565b6020026020010151905060008060008061406b89600001518a602001518c886143c3565b935093509350935083156140bf57828760400151886060015161ffff1681518110614098576140986154fb565b6020908102919091010152606087018051906140b382615b9b565b61ffff169052506140f2565b600160f81b876020015187815181106140da576140da6154fb565b60200101906001600160f81b031916908160001a9053505b87806140fc575080155b85515188518051929a50909188908110614118576141186154fb565b602002602001018181525050818760800151878151811061413b5761413b6154fb565b60200260200101906001600160601b031690816001600160601b031681525050846000015160a001518760a00151878151811061417a5761417a6154fb565b602090810291909101015284516040015160c08801518051889081106141a2576141a26154fb565b60200260200101906001600160a01b031690816001600160a01b031681525050505050505080806141d290615550565b915050614027565b50606083015151156142f7576000816060015161ffff166001600160401b0381111561420857614208614d29565b60405190808252806020026020018201604052801561423b57816020015b60608152602001906001900390816142265790505b50905060005b826060015161ffff1681101561429f5782604001518181518110614267576142676154fb565b6020026020010151828281518110614281576142816154fb565b6020026020010181905250808061429790615550565b915050614241565b507fbf0f908fceae532c9263244a7809d80917f8411f33f2746d7b0ebf835fc48de6826000015183602001518385608001518660a001518760c001516040516142ed96959493929190616078565b60405180910390a1505b509392505050565b8151600090801561431c575060408201516001600160401b031615155b156143705761010083015163ffffffff164310808061435d575061010084015161434c9063ffffffff1643615527565b83602001516001600160401b031610155b1561436e5750506040810151613028565b505b503a92915050565b6000806000606485606001516064614390919061613c565b61439d9060ff1689615cd6565b6143a79190615ff4565b90506143b58187878761477d565b925092505094509492505050565b805160a090810151600090815260086020908152604080832085519485015191519394606094869485948592614401928e928e929091879101615e1a565b60408051601f19818403018152918152815160209283012084516000908152600c90935291205490915081146144745750505460408051808201909152601081526f756e6b6e6f776e2063616c6c6261636b60801b60208201526001955093506001600160601b03169150839050614772565b5061447d614c1d565b60006144b27f00000000000000000000000000000000000000000000000000000000000000006001600160401b038e16615ff4565b6040805160808101825263ffffffff909216825262ffffff8d1660208084019190915285015161ffff16828201528401516001600160a01b031660608201529150899050614546575050604080518082019091526016815275756e617661696c61626c652072616e646f6d6e65737360501b60208201529054600195509093506001600160601b0316915060009050614772565b60006145588360000151838c8f613863565b606080840151855191860151604051939450909260009263d21ea8fd60e01b9261458792879190602401616155565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091526004805461ff00191661010017905590506000805a905060006145f28e60000151608001516001600160601b03168960400151866147f1565b909350905080614627578d516080015160405163aad1598360e01b8152610a54916001600160601b0316908490600401615923565b506000610bb85a6146389190615ca9565b6004805461ff00191690559050818110156146615761466161465a8284615527565b8f51614830565b8854600160601b900463ffffffff1689600c61467c8361618a565b825463ffffffff9182166101009390930a92830291909202199091161790555087516000908152600c60205260408120558261470157505060408051808201909152601081526f195e1958dd5d1a5bdb8819985a5b195960821b6020820152965460019b50969950506001600160601b03909516965060009550614772945050505050565b8854600160c01b90046001600160401b031689601861471f83615c51565b82546001600160401b039182166101009390930a92830291909202199091161790555050604080516020810190915260008082529854989c509a50506001600160601b0390961697508996505050505050505b945094509450949050565b60008080851561478d5785614797565b6147978585614a5c565b90506000816147ae89670de0b6b3a7640000615cd6565b6147b89190615ff4565b9050676765c793fa10079d601b1b8111156147e55760405162de437160e81b815260040160405180910390fd5b97909650945050505050565b6000805a610bb8811061482757610bb881039050856040820482031115614827576000808551602087016000898bf19250600191505b50935093915050565b80608001516001600160601b0316821115614849575050565b60045460009060649061486690600160201b900460ff16826161aa565b60ff168360c001518585608001516001600160601b03166148879190615527565b6148919190615cd6565b61489b9190615cd6565b6148a59190615ff4565b60e080840151604080516101408101825260045460ff80821615158352610100808304821615156020808601919091526201000084048316151585870152630100000084048316606080870191909152600160201b80860490941660808088019190915263ffffffff600160281b8704811660a0890152600160481b8704811660c0890152600160681b870481169a88019a909a52600160881b9095048916928601929092526005546001600160601b039081166101208701528651948501875260065498891685526001600160401b03938904841691850191909152600160601b880490921694830194909452600160a01b909504909416918401919091529293506000926149b8928592919061477d565b5060a08401516000908152600860205260408120805492935083929091906149ea9084906001600160601b0316615c89565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600260008282829054906101000a90046001600160601b0316614a329190615c89565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555050505050565b60a0820151606082015160009190600163ffffffff83161190818015614a9857508451614a8f9063ffffffff1642615527565b8363ffffffff16105b15614aa557506101208501515b6001600160601b031695945050505050565b604051806080016040528060006001600160401b03168152602001600062ffffff16815260200160008152602001600081525090565b6040518061010001604052806008906020820280368337509192915050565b600183019183908215614b9a5791602002820160005b83821115614b6957833562ffffff1683826101000a81548162ffffff021916908362ffffff1602179055509260200192600301602081600201049283019260010302614b22565b8015614b985782816101000a81549062ffffff0219169055600301602081600201049283019260010302614b69565b505b50614ba6929150614c44565b5090565b828054828255906000526020600020908101928215614b9a579160200282015b82811115614b9a57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190614bca565b5080546000825590600052602060002090810190610a5d9190614c44565b60408051608081018252600080825260208201819052918101829052606081019190915290565b5b80821115614ba65760008155600101614c45565b6001600160a01b0381168114610a5d57600080fd5b8035614c7981614c59565b919050565b600060208284031215614c9057600080fd5b8135614c9b81614c59565b9392505050565b60008060408385031215614cb557600080fd5b823591506020830135614cc781614c59565b809150509250929050565b600060208284031215614ce457600080fd5b5035919050565b60006101408284031215614cfe57600080fd5b50919050565b803561ffff81168114614c7957600080fd5b803562ffffff81168114614c7957600080fd5b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715614d6157614d61614d29565b60405290565b60405161010081016001600160401b0381118282101715614d6157614d61614d29565b60405160a081016001600160401b0381118282101715614d6157614d61614d29565b604051602081016001600160401b0381118282101715614d6157614d61614d29565b604051601f8201601f191681016001600160401b0381118282101715614df657614df6614d29565b604052919050565b600082601f830112614e0f57600080fd5b81356001600160401b03811115614e2857614e28614d29565b614e3b601f8201601f1916602001614dce565b818152846020838601011115614e5057600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060808587031215614e8357600080fd5b84359350614e9360208601614d04565b9250614ea160408601614d16565b915060608501356001600160401b03811115614ebc57600080fd5b614ec887828801614dfe565b91505092959194509250565b60008083601f840112614ee657600080fd5b5081356001600160401b03811115614efd57600080fd5b6020830191508360208260051b850101111561351857600080fd5b60008060008060408587031215614f2e57600080fd5b84356001600160401b0380821115614f4557600080fd5b614f5188838901614ed4565b90965094506020870135915080821115614f6a57600080fd5b50614f7787828801614ed4565b95989497509550505050565b80356001600160401b0381168114614c7957600080fd5b600080600080600060808688031215614fb257600080fd5b85356001600160401b03811115614fc857600080fd5b614fd488828901614ed4565b90965094505060208601356001600160c01b0381168114614ff457600080fd5b925061500260408701614f83565b915061501060608701614f83565b90509295509295909350565b6001600160a01b0391909116815260200190565b6101008101818360005b600881101561505e57815162ffffff1683526020928301929091019060010161503a565b50505092915050565b63ffffffff81168114610a5d57600080fd5b8035614c7981615067565b60006040828403121561509657600080fd5b61509e614d3f565b82356150a981615067565b815260208301356150b981615067565b60208201529392505050565b8015158114610a5d57600080fd5b8035614c79816150c5565b6000602082840312156150f057600080fd5b8135614c9b816150c5565b600061010080838503121561510f57600080fd5b83818401111561511e57600080fd5b509092915050565b6000806000806080858703121561513c57600080fd5b84359350602085013561514e81615067565b925060408501356001600160401b038082111561516a57600080fd5b61517688838901614dfe565b9350606087013591508082111561518c57600080fd5b50614ec887828801614dfe565b60008083601f8401126151ab57600080fd5b5081356001600160401b038111156151c257600080fd5b60208301915083602082850101111561351857600080fd5b600080600080606085870312156151f057600080fd5b84356151fb81614c59565b93506020850135925060408501356001600160401b0381111561521d57600080fd5b614f7787828801615199565b60008060006060848603121561523e57600080fd5b833592506020840135915060408401356001600160401b0381111561526257600080fd5b61526e86828701614dfe565b9150509250925092565b600081518084526020808501945080840160005b838110156152a85781518752958201959082019060010161528c565b509495945050505050565b602081526000614c9b6020830184615278565b6000806000604084860312156152db57600080fd5b83356152e681614c59565b925060208401356001600160401b0381111561530157600080fd5b61530d86828701615199565b9497909650939450505050565b6000806020838503121561532d57600080fd5b82356001600160401b0381111561534357600080fd5b61534f85828601615199565b90969095509350505050565b60008060008060008060c0878903121561537457600080fd5b8635955061538460208801614d04565b945061539260408801614d16565b935060608701356153a281615067565b925060808701356001600160401b03808211156153be57600080fd5b6153ca8a838b01614dfe565b935060a08901359150808211156153e057600080fd5b506153ed89828a01614dfe565b9150509295509295509295565b600081518084526020808501945080840160005b838110156152a85781516001600160a01b03168752958201959082019060010161540e565b6001600160601b038716815263ffffffff861660208201526001600160401b038581166040830152841660608201526001600160a01b038316608082015260c060a08201819052600090611e3a908301846153fa565b6000806040838503121561549c57600080fd5b82356154a781614c59565b946020939093013593505050565b600080604083850312156154c857600080fd5b8235915060208301356001600160401b038111156154e557600080fd5b6154f185828601614dfe565b9150509250929050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b8181038181111561302857613028615511565b634e487b7160e01b600052603160045260246000fd5b60006001820161556257615562615511565b5060010190565b6001600160601b0382811682821603908082111561558957615589615511565b5092915050565b6001600160a01b03929092168252602082015260400190565b6000602082840312156155bb57600080fd5b8151614c9b816150c5565b6001600160a01b039290921682526001600160601b0316602082015260400190565b60ff81168114610a5d57600080fd5b8035614c79816155e8565b60006020828403121561561457600080fd5b8135614c9b816155e8565b60008135613028816150c5565b60008135613028816155e8565b6000813561302881615067565b6001600160601b0381168114610a5d57600080fd5b6000813561302881615646565b8135615673816150c5565b815490151560ff1660ff19919091161781556156ae6156946020840161561f565b82805461ff00191691151560081b61ff0016919091179055565b6156d96156bd6040840161561f565b82805462ff0000191691151560101b62ff000016919091179055565b6157026156e86060840161562c565b825463ff000000191660189190911b63ff00000016178255565b61572f6157116080840161562c565b82805460ff60201b191660209290921b60ff60201b16919091179055565b61576261573e60a08401615639565b82805463ffffffff60281b191660289290921b63ffffffff60281b16919091179055565b61579561577160c08401615639565b82805463ffffffff60481b191660489290921b63ffffffff60481b16919091179055565b6157c86157a460e08401615639565b82805463ffffffff60681b191660689290921b63ffffffff60681b16919091179055565b6157fc6157d86101008401615639565b82805463ffffffff60881b191660889290921b63ffffffff60881b16919091179055565b61169b61580c610120840161565b565b6001830180546001600160601b0319166001600160601b0392909216919091179055565b8035614c7981615646565b61014081016158538261584d856150d3565b15159052565b61585f602084016150d3565b15156020830152615872604084016150d3565b15156040830152615885606084016155f7565b60ff166060830152615899608084016155f7565b60ff1660808301526158ad60a08401615079565b63ffffffff1660a08301526158c460c08401615079565b63ffffffff1660c08301526158db60e08401615079565b63ffffffff1660e08301526101006158f4848201615079565b63ffffffff169083015261012061590c848201615830565b6001600160601b038116848301525b505092915050565b918252602082015260400190565b6000823560be1983360301811261594757600080fd5b9190910192915050565b600082601f83011261596257600080fd5b813560206001600160401b038083111561597e5761597e614d29565b8260051b61598d838201614dce565b93845285810183019383810190888611156159a757600080fd5b84880192505b85831015611e3a578235848111156159c457600080fd5b8801601f196040828c03820112156159db57600080fd5b6159e3614d3f565b87830135878111156159f457600080fd5b8301610100818e0384011215615a0957600080fd5b615a11614d67565b9250888101358352615a2560408201614d04565b89840152615a3560608201614c6e565b6040840152608081013588811115615a4c57600080fd5b615a5a8e8b83850101614dfe565b606085015250615a6c60a08201615830565b608084015260c081013560a084015260e081013560c084015261010081013560e084015250818152615aa060408401615830565b818901528452505091840191908401906159ad565b600081360360c0811215615ac857600080fd5b615ad0614d8a565b615ad984614f83565b81526020615ae8818601614d16565b828201526040603f1984011215615afe57600080fd5b615b06614dac565b925036605f860112615b1757600080fd5b615b1f614d3f565b806080870136811115615b3157600080fd5b604088015b81811015615b4d5780358452928401928401615b36565b50908552604084019490945250509035906001600160401b03821115615b7257600080fd5b615b7e36838601615951565b6060820152615b8f60a085016150d3565b60808201529392505050565b600061ffff808316818103615bb257615bb2615511565b6001019392505050565b6000608080830160018060401b038089168552602060018060c01b038916818701526040828916818801526060858189015284895180875260a08a019150848b01965060005b81811015615c3e5787518051881684528681015162ffffff16878501528581015186850152840151848401529685019691880191600101615c02565b50909d9c50505050505050505050505050565b60006001600160401b038281166002600160401b03198101615bb257615bb2615511565b61ffff929092168252602082015260400190565b6001600160601b0381811683821601908082111561558957615589615511565b8082018082111561302857613028615511565b6001600160a01b0392831681529116602082015260400190565b808202811582820484141761302857613028615511565b600060408284031215615cff57600080fd5b615d07614d3f565b8235615d12816155e8565b81526020928301359281019290925250919050565b600060208284031215615d3957600080fd5b8151614c9b816155e8565b6020815260ff82511660208201526020820151604082015260018060a01b0360408301511660608201526000606083015160a06080840152615d8960c08401826153fa565b608094909401516001600160601b031660a093909301929092525090919050565b6000815180845260005b81811015615dd057602081850181015186830182015201615db4565b506000602082860101526020601f19601f83011685010191505092915050565b602081526000614c9b6020830184615daa565b9182526001600160a01b0316602082015260400190565b60018060401b038516815262ffffff84166020820152826040820152608060608201528151608082015261ffff60208301511660a082015260018060a01b0360408301511660c0820152600060608301516101008060e0850152615e82610180850183615daa565b91506080850151615e9d828601826001600160601b03169052565b505060a084015161012084015260c084015161014084015260e08401516101608401528091505095945050505050565b6001600160a01b038c1681526001600160401b038b16602082015262ffffff8a1660408201526060810189905261ffff8816608082015263ffffffff871660a082015260c0810186905260e081018590526101606101008201819052600090615f3883820187615daa565b61012084019590955250506001600160601b0391909116610140909101529998505050505050505050565b600060208284031215615f7557600080fd5b5051919050565b65ffffffffffff81811683821601908082111561558957615589615511565b6001600160601b03929092168252602082015260400190565b815160408201908260005b600281101561505e578251825260209283019290910190600101615fbf565b634e487b7160e01b600052601260045260246000fd5b60008261600357616003615fde565b500490565b6001600160401b0381811683821602808216919082811461591b5761591b615511565b60008261603a5761603a615fde565b500690565b600081518084526020808501945080840160005b838110156152a85781516001600160601b031687529582019590820190600101616053565b60c08152600061608b60c0830189615278565b60208382038185015261609e828a615daa565b915083820360408501528188518084528284019150828160051b850101838b0160005b838110156160ef57601f198784030185526160dd838351615daa565b948601949250908501906001016160c1565b50508681036060880152616103818b61603f565b945050505050828103608084015261611b8186615278565b905082810360a084015261612f81856153fa565b9998505050505050505050565b60ff818116838216019081111561302857613028615511565b83815260606020820152600061616e6060830185615278565b82810360408401526161808185615daa565b9695505050505050565b600063ffffffff8216806161a0576161a0615511565b6000190192915050565b60ff82811682821603908111156130285761302861551156fea164736f6c6343000813000a",
}

var VRFCoordinatorABI = VRFCoordinatorMetaData.ABI

var VRFCoordinatorBin = VRFCoordinatorMetaData.Bin

func DeployVRFCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, beaconPeriodBlocksArg *big.Int, linkToken common.Address) (common.Address, *types.Transaction, *VRFCoordinator, error) {
	parsed, err := VRFCoordinatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorBin), backend, beaconPeriodBlocksArg, linkToken)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinator{address: address, abi: *parsed, VRFCoordinatorCaller: VRFCoordinatorCaller{contract: contract}, VRFCoordinatorTransactor: VRFCoordinatorTransactor{contract: contract}, VRFCoordinatorFilterer: VRFCoordinatorFilterer{contract: contract}}, nil
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
	outstruct.PendingFulfillments = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.TotalRequests = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.TotalSuccessfulFulfillments = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.Owner = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[5], new([]common.Address)).(*[]common.Address)

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

func (_VRFCoordinator *VRFCoordinatorCaller) SCallbackConfig(opts *bind.CallOpts) (SCallbackConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "s_callbackConfig")

	outstruct := new(SCallbackConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaxCallbackGasLimit = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.MaxCallbackArgumentsLength = *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_VRFCoordinator *VRFCoordinatorSession) SCallbackConfig() (SCallbackConfig,

	error) {
	return _VRFCoordinator.Contract.SCallbackConfig(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) SCallbackConfig() (SCallbackConfig,

	error) {
	return _VRFCoordinator.Contract.SCallbackConfig(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) SCoordinatorConfig(opts *bind.CallOpts) (SCoordinatorConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "s_coordinatorConfig")

	outstruct := new(SCoordinatorConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.UseReasonableGasPrice = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.ReentrancyLock = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Paused = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.PremiumPercentage = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.UnusedGasPenaltyPercent = *abi.ConvertType(out[4], new(uint8)).(*uint8)
	outstruct.StalenessSeconds = *abi.ConvertType(out[5], new(uint32)).(*uint32)
	outstruct.RedeemableRequestGasOverhead = *abi.ConvertType(out[6], new(uint32)).(*uint32)
	outstruct.CallbackRequestGasOverhead = *abi.ConvertType(out[7], new(uint32)).(*uint32)
	outstruct.ReasonableGasPriceStalenessBlocks = *abi.ConvertType(out[8], new(uint32)).(*uint32)
	outstruct.FallbackWeiPerUnitLink = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFCoordinator *VRFCoordinatorSession) SCoordinatorConfig() (SCoordinatorConfig,

	error) {
	return _VRFCoordinator.Contract.SCoordinatorConfig(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) SCoordinatorConfig() (SCoordinatorConfig,

	error) {
	return _VRFCoordinator.Contract.SCoordinatorConfig(&_VRFCoordinator.CallOpts)
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

func (_VRFCoordinator *VRFCoordinatorTransactor) ProcessVRFOutputs(opts *bind.TransactOpts, vrfOutputs []VRFBeaconTypesVRFOutput, juelsPerFeeCoin *big.Int, reasonableGasPrice uint64, blockHeight uint64) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "processVRFOutputs", vrfOutputs, juelsPerFeeCoin, reasonableGasPrice, blockHeight)
}

func (_VRFCoordinator *VRFCoordinatorSession) ProcessVRFOutputs(vrfOutputs []VRFBeaconTypesVRFOutput, juelsPerFeeCoin *big.Int, reasonableGasPrice uint64, blockHeight uint64) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.ProcessVRFOutputs(&_VRFCoordinator.TransactOpts, vrfOutputs, juelsPerFeeCoin, reasonableGasPrice, blockHeight)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) ProcessVRFOutputs(vrfOutputs []VRFBeaconTypesVRFOutput, juelsPerFeeCoin *big.Int, reasonableGasPrice uint64, blockHeight uint64) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.ProcessVRFOutputs(&_VRFCoordinator.TransactOpts, vrfOutputs, juelsPerFeeCoin, reasonableGasPrice, blockHeight)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) RedeemRandomness(opts *bind.TransactOpts, subID *big.Int, requestID *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "redeemRandomness", subID, requestID, arg2)
}

func (_VRFCoordinator *VRFCoordinatorSession) RedeemRandomness(subID *big.Int, requestID *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RedeemRandomness(&_VRFCoordinator.TransactOpts, subID, requestID, arg2)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RedeemRandomness(subID *big.Int, requestID *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RedeemRandomness(&_VRFCoordinator.TransactOpts, subID, requestID, arg2)
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

func (_VRFCoordinator *VRFCoordinatorTransactor) RequestRandomness(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "requestRandomness", subID, numWords, confDelay, arg3)
}

func (_VRFCoordinator *VRFCoordinatorSession) RequestRandomness(subID *big.Int, numWords uint16, confDelay *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomness(&_VRFCoordinator.TransactOpts, subID, numWords, confDelay, arg3)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RequestRandomness(subID *big.Int, numWords uint16, confDelay *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomness(&_VRFCoordinator.TransactOpts, subID, numWords, confDelay, arg3)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) RequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte, arg5 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "requestRandomnessFulfillment", subID, numWords, confDelay, callbackGasLimit, arguments, arg5)
}

func (_VRFCoordinator *VRFCoordinatorSession) RequestRandomnessFulfillment(subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte, arg5 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomnessFulfillment(&_VRFCoordinator.TransactOpts, subID, numWords, confDelay, callbackGasLimit, arguments, arg5)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) RequestRandomnessFulfillment(subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte, arg5 []byte) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.RequestRandomnessFulfillment(&_VRFCoordinator.TransactOpts, subID, numWords, confDelay, callbackGasLimit, arguments, arg5)
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

func (_VRFCoordinator *VRFCoordinatorTransactor) SetCallbackConfig(opts *bind.TransactOpts, config VRFCoordinatorCallbackConfig) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "setCallbackConfig", config)
}

func (_VRFCoordinator *VRFCoordinatorSession) SetCallbackConfig(config VRFCoordinatorCallbackConfig) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetCallbackConfig(&_VRFCoordinator.TransactOpts, config)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) SetCallbackConfig(config VRFCoordinatorCallbackConfig) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetCallbackConfig(&_VRFCoordinator.TransactOpts, config)
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

func (_VRFCoordinator *VRFCoordinatorTransactor) SetCoordinatorConfig(opts *bind.TransactOpts, coordinatorConfig VRFBeaconTypesCoordinatorConfig) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "setCoordinatorConfig", coordinatorConfig)
}

func (_VRFCoordinator *VRFCoordinatorSession) SetCoordinatorConfig(coordinatorConfig VRFBeaconTypesCoordinatorConfig) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetCoordinatorConfig(&_VRFCoordinator.TransactOpts, coordinatorConfig)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) SetCoordinatorConfig(coordinatorConfig VRFBeaconTypesCoordinatorConfig) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetCoordinatorConfig(&_VRFCoordinator.TransactOpts, coordinatorConfig)
}

func (_VRFCoordinator *VRFCoordinatorTransactor) SetPauseFlag(opts *bind.TransactOpts, pause bool) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "setPauseFlag", pause)
}

func (_VRFCoordinator *VRFCoordinatorSession) SetPauseFlag(pause bool) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetPauseFlag(&_VRFCoordinator.TransactOpts, pause)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) SetPauseFlag(pause bool) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetPauseFlag(&_VRFCoordinator.TransactOpts, pause)
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

type VRFCoordinatorCallbackConfigSetIterator struct {
	Event *VRFCoordinatorCallbackConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorCallbackConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorCallbackConfigSet)
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
		it.Event = new(VRFCoordinatorCallbackConfigSet)
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

func (it *VRFCoordinatorCallbackConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorCallbackConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorCallbackConfigSet struct {
	NewConfig VRFCoordinatorCallbackConfig
	Raw       types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterCallbackConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorCallbackConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "CallbackConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorCallbackConfigSetIterator{contract: _VRFCoordinator.contract, event: "CallbackConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchCallbackConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorCallbackConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "CallbackConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorCallbackConfigSet)
				if err := _VRFCoordinator.contract.UnpackLog(event, "CallbackConfigSet", log); err != nil {
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

func (_VRFCoordinator *VRFCoordinatorFilterer) ParseCallbackConfigSet(log types.Log) (*VRFCoordinatorCallbackConfigSet, error) {
	event := new(VRFCoordinatorCallbackConfigSet)
	if err := _VRFCoordinator.contract.UnpackLog(event, "CallbackConfigSet", log); err != nil {
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
	CoordinatorConfig VRFBeaconTypesCoordinatorConfig
	Raw               types.Log
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

type VRFCoordinatorPauseFlagChangedIterator struct {
	Event *VRFCoordinatorPauseFlagChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorPauseFlagChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorPauseFlagChanged)
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
		it.Event = new(VRFCoordinatorPauseFlagChanged)
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

func (it *VRFCoordinatorPauseFlagChangedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorPauseFlagChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorPauseFlagChanged struct {
	Paused bool
	Raw    types.Log
}

func (_VRFCoordinator *VRFCoordinatorFilterer) FilterPauseFlagChanged(opts *bind.FilterOpts) (*VRFCoordinatorPauseFlagChangedIterator, error) {

	logs, sub, err := _VRFCoordinator.contract.FilterLogs(opts, "PauseFlagChanged")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorPauseFlagChangedIterator{contract: _VRFCoordinator.contract, event: "PauseFlagChanged", logs: logs, sub: sub}, nil
}

func (_VRFCoordinator *VRFCoordinatorFilterer) WatchPauseFlagChanged(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorPauseFlagChanged) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinator.contract.WatchLogs(opts, "PauseFlagChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorPauseFlagChanged)
				if err := _VRFCoordinator.contract.UnpackLog(event, "PauseFlagChanged", log); err != nil {
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

func (_VRFCoordinator *VRFCoordinatorFilterer) ParsePauseFlagChanged(log types.Log) (*VRFCoordinatorPauseFlagChanged, error) {
	event := new(VRFCoordinatorPauseFlagChanged)
	if err := _VRFCoordinator.contract.UnpackLog(event, "PauseFlagChanged", log); err != nil {
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
	Consumers             []common.Address
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

type GetSubscription struct {
	Balance                     *big.Int
	PendingFulfillments         uint32
	TotalRequests               uint64
	TotalSuccessfulFulfillments uint64
	Owner                       common.Address
	Consumers                   []common.Address
}
type SCallbackConfig struct {
	MaxCallbackGasLimit        uint32
	MaxCallbackArgumentsLength uint32
}
type SCoordinatorConfig struct {
	UseReasonableGasPrice             bool
	ReentrancyLock                    bool
	Paused                            bool
	PremiumPercentage                 uint8
	UnusedGasPenaltyPercent           uint8
	StalenessSeconds                  uint32
	RedeemableRequestGasOverhead      uint32
	CallbackRequestGasOverhead        uint32
	ReasonableGasPriceStalenessBlocks uint32
	FallbackWeiPerUnitLink            *big.Int
}
type SPendingRequests struct {
	SlotNumber        uint32
	ConfirmationDelay *big.Int
	NumWords          uint16
	Requester         common.Address
}

func (_VRFCoordinator *VRFCoordinator) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinator.abi.Events["CallbackConfigSet"].ID:
		return _VRFCoordinator.ParseCallbackConfigSet(log)
	case _VRFCoordinator.abi.Events["CoordinatorConfigSet"].ID:
		return _VRFCoordinator.ParseCoordinatorConfigSet(log)
	case _VRFCoordinator.abi.Events["CoordinatorDeregistered"].ID:
		return _VRFCoordinator.ParseCoordinatorDeregistered(log)
	case _VRFCoordinator.abi.Events["CoordinatorRegistered"].ID:
		return _VRFCoordinator.ParseCoordinatorRegistered(log)
	case _VRFCoordinator.abi.Events["MigrationCompleted"].ID:
		return _VRFCoordinator.ParseMigrationCompleted(log)
	case _VRFCoordinator.abi.Events["OutputsServed"].ID:
		return _VRFCoordinator.ParseOutputsServed(log)
	case _VRFCoordinator.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFCoordinator.ParseOwnershipTransferRequested(log)
	case _VRFCoordinator.abi.Events["OwnershipTransferred"].ID:
		return _VRFCoordinator.ParseOwnershipTransferred(log)
	case _VRFCoordinator.abi.Events["PauseFlagChanged"].ID:
		return _VRFCoordinator.ParsePauseFlagChanged(log)
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

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFCoordinatorCallbackConfigSet) Topic() common.Hash {
	return common.HexToHash("0x0cc54509a45ab33cd67614d4a2892c083ecf8fb43b9d29f6ea8130b9023e51df")
}

func (VRFCoordinatorCoordinatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0x0028d3a46e95e67def989d41c66eb331add9809460b95b5fb4eb006157728fc5")
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

func (VRFCoordinatorOutputsServed) Topic() common.Hash {
	return common.HexToHash("0xf10ea936d00579b4c52035ee33bf46929646b3aa87554c565d8fb2c7aa549c44")
}

func (VRFCoordinatorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFCoordinatorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFCoordinatorPauseFlagChanged) Topic() common.Hash {
	return common.HexToHash("0x49ba7c1de2d8853088b6270e43df2118516b217f38b917dd2b80dea360860fbe")
}

func (VRFCoordinatorRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0xbf0f908fceae532c9263244a7809d80917f8411f33f2746d7b0ebf835fc48de6")
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

func (_VRFCoordinator *VRFCoordinator) Address() common.Address {
	return _VRFCoordinator.address
}

type VRFCoordinatorInterface interface {
	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	MAXNUMWORDS(opts *bind.CallOpts) (*big.Int, error)

	NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error)

	GetCallbackMemo(opts *bind.CallOpts, requestId *big.Int) ([32]byte, error)

	GetConfirmationDelays(opts *bind.CallOpts) ([8]*big.Int, error)

	GetFee(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) (*big.Int, error)

	GetFulfillmentFee(opts *bind.CallOpts, arg0 *big.Int, callbackGasLimit uint32, arguments []byte, arg3 []byte) (*big.Int, error)

	GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

		error)

	GetSubscriptionLinkBalance(opts *bind.CallOpts) (*big.Int, error)

	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)

	ILink(opts *bind.CallOpts) (common.Address, error)

	MigrationVersion(opts *bind.CallOpts) (uint8, error)

	OnMigration(opts *bind.CallOpts, arg0 []byte) error

	Owner(opts *bind.CallOpts) (common.Address, error)

	SCallbackConfig(opts *bind.CallOpts) (SCallbackConfig,

		error)

	SCoordinatorConfig(opts *bind.CallOpts) (SCoordinatorConfig,

		error)

	SPendingRequests(opts *bind.CallOpts, arg0 *big.Int) (SPendingRequests,

		error)

	SProducer(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	BatchTransferLink(opts *bind.TransactOpts, recipients []common.Address, paymentsInJuels []*big.Int) (*types.Transaction, error)

	CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	DeregisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error)

	Migrate(opts *bind.TransactOpts, newCoordinator common.Address, encodedRequest []byte) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	ProcessVRFOutputs(opts *bind.TransactOpts, vrfOutputs []VRFBeaconTypesVRFOutput, juelsPerFeeCoin *big.Int, reasonableGasPrice uint64, blockHeight uint64) (*types.Transaction, error)

	RedeemRandomness(opts *bind.TransactOpts, subID *big.Int, requestID *big.Int, arg2 []byte) (*types.Transaction, error)

	RegisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, arg3 []byte) (*types.Transaction, error)

	RequestRandomnessFulfillment(opts *bind.TransactOpts, subID *big.Int, numWords uint16, confDelay *big.Int, callbackGasLimit uint32, arguments []byte, arg5 []byte) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error)

	SetCallbackConfig(opts *bind.TransactOpts, config VRFCoordinatorCallbackConfig) (*types.Transaction, error)

	SetConfirmationDelays(opts *bind.TransactOpts, confDelays [8]*big.Int) (*types.Transaction, error)

	SetCoordinatorConfig(opts *bind.TransactOpts, coordinatorConfig VRFBeaconTypesCoordinatorConfig) (*types.Transaction, error)

	SetPauseFlag(opts *bind.TransactOpts, pause bool) (*types.Transaction, error)

	SetProducer(opts *bind.TransactOpts, producer common.Address) (*types.Transaction, error)

	TransferLink(opts *bind.TransactOpts, recipient common.Address, juelsAmount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterCallbackConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorCallbackConfigSetIterator, error)

	WatchCallbackConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorCallbackConfigSet) (event.Subscription, error)

	ParseCallbackConfigSet(log types.Log) (*VRFCoordinatorCallbackConfigSet, error)

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

	FilterOutputsServed(opts *bind.FilterOpts) (*VRFCoordinatorOutputsServedIterator, error)

	WatchOutputsServed(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorOutputsServed) (event.Subscription, error)

	ParseOutputsServed(log types.Log) (*VRFCoordinatorOutputsServed, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorOwnershipTransferred, error)

	FilterPauseFlagChanged(opts *bind.FilterOpts) (*VRFCoordinatorPauseFlagChangedIterator, error)

	WatchPauseFlagChanged(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorPauseFlagChanged) (event.Subscription, error)

	ParsePauseFlagChanged(log types.Log) (*VRFCoordinatorPauseFlagChanged, error)

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

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
