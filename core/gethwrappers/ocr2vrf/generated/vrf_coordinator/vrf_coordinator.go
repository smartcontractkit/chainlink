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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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
)

type ECCArithmeticG1Point struct {
	P [2]*big.Int
}

type VRFBeaconTypesBeaconRequest struct {
	SlotNumber        uint32
	ConfirmationDelay *big.Int
	NumWords          uint16
	Requester         common.Address
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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"beaconPeriodBlocksArg\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BeaconPeriodMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"earliestAllowed\",\"type\":\"uint256\"}],\"name\":\"BlockTooRecent\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"firstDelay\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"minDelay\",\"type\":\"uint16\"}],\"name\":\"ConfirmationDelayBlocksTooShort\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confirmationDelays\",\"type\":\"uint16[10]\"},{\"internalType\":\"uint8\",\"name\":\"violatingIndex\",\"type\":\"uint8\"}],\"name\":\"ConfirmationDelaysNotIncreasing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reportHeight\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"separatorHeight\",\"type\":\"uint64\"}],\"name\":\"HistoryDomainSeparatorTooOld\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidBillingConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidJuelsConversion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"requestedVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"coordinatorVersion\",\"type\":\"uint8\"}],\"name\":\"MigrationVersionMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoWordsRequested\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16[10]\",\"name\":\"confDelays\",\"type\":\"uint16[10]\"}],\"name\":\"NonZeroDelayAfterZeroDelay\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnMigrationNotSupported\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint256\",\"name\":\"requestHeight\",\"type\":\"uint256\"}],\"name\":\"RandomnessNotAvailable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"actual\",\"type\":\"address\"}],\"name\":\"ResponseMustBeRetrievedByRequester\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyRequestsReplaceContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManySlotsReplaceContract\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"max\",\"type\":\"uint256\"}],\"name\":\"TooManyWords\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockHeight\",\"type\":\"uint256\"}],\"name\":\"UniverseHasEndedBangBangBang\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"givenDelay\",\"type\":\"uint24\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"knownDelays\",\"type\":\"uint24[8]\"}],\"name\":\"UnknownConfirmationDelay\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"f\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"onchainConfig\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"offchainConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"offchainConfig\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"newVersion\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint40\",\"name\":\"epochAndRound\",\"type\":\"uint40\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"recentBlockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputsServed\",\"type\":\"tuple[]\"}],\"name\":\"OutputsServed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.RequestID[]\",\"name\":\"requestIDs\",\"type\":\"uint48[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"successfulFulfillment\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"truncatedErrorData\",\"type\":\"bytes[]\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAllowance\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"}],\"name\":\"RandomnessFulfillmentRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nextBeaconOutputHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confDelay\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NUM_CONF_DELAYS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ROUTER\",\"outputs\":[{\"internalType\":\"contractVRFRouter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"paymentsInJuels\",\"type\":\"uint256[]\"}],\"name\":\"batchTransferLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"}],\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"}],\"name\":\"calculateRequestPriceCallbackJuels\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateRequestPriceJuels\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"name\":\"forgetConsumerSubscriptionID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestId\",\"type\":\"uint48\"}],\"name\":\"getBeaconRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"VRFBeaconTypes.SlotNumber\",\"name\":\"slotNumber\",\"type\":\"uint32\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"internalType\":\"structVRFBeaconTypes.BeaconRequest\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestId\",\"type\":\"uint48\"}],\"name\":\"getCallbackMemo\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfirmationDelays\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"\",\"type\":\"uint24[8]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"pendingFulfillments\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalLinkBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_StartSlot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_beaconPeriodBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxNumWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractVRFMigrationInterface\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"encodedRequest\",\"type\":\"bytes\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrationVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minDelay\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onMigration\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"}],\"internalType\":\"structECCArithmetic.G1Point\",\"name\":\"vrfOutput\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"uint96\",\"name\":\"gasAllowance\",\"type\":\"uint96\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weiPerUnitLink\",\"type\":\"uint256\"}],\"internalType\":\"structVRFBeaconTypes.Callback\",\"name\":\"callback\",\"type\":\"tuple\"},{\"internalType\":\"uint96\",\"name\":\"price\",\"type\":\"uint96\"}],\"internalType\":\"structVRFBeaconTypes.CostedCallback[]\",\"name\":\"callbacks\",\"type\":\"tuple[]\"}],\"internalType\":\"structVRFBeaconTypes.VRFOutput[]\",\"name\":\"vrfOutputs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint192\",\"name\":\"juelsPerFeeCoin\",\"type\":\"uint192\"},{\"internalType\":\"uint64\",\"name\":\"reasonableGasPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"processVRFOutputs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelay\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"proofG1X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"proofG1Y\",\"type\":\"uint256\"}],\"internalType\":\"structVRFBeaconTypes.OutputServed[]\",\"name\":\"outputs\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"producer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"requestID\",\"type\":\"uint48\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"redeemRandomness\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"randomness\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"requestRandomness\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"numWords\",\"type\":\"uint16\"},{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay\",\"name\":\"confirmationDelayArg\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"arguments\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"requestRandomnessFulfillment\",\"outputs\":[{\"internalType\":\"VRFBeaconTypes.RequestID\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"useReasonableGasPrice\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"unusedGasPenaltyPercent\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"redeemableRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callbackRequestGasOverhead\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"premiumPercentage\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPriceStalenessBlocks\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"internalType\":\"structVRFBeaconTypes.BillingConfig\",\"name\":\"billingConfig\",\"type\":\"tuple\"}],\"name\":\"setBillingConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"VRFBeaconTypes.ConfirmationDelay[8]\",\"name\":\"confDelays\",\"type\":\"uint24[8]\"}],\"name\":\"setConfirmationDelays\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setProducer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"gasPrice\",\"type\":\"uint64\"}],\"name\":\"setReasonableGasPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"juelsAmount\",\"type\":\"uint256\"}],\"name\":\"transferLink\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b50604051620060c9380380620060c98339810160408190526200003591620002be565b82828233806000816200008f5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c257620000c28162000167565b5050506001600160a01b0392831660805290821660a0521660c05260008490036200010057604051632abc297960e01b815260040160405180910390fd5b8360e0818152505060006200011f6200021260201b62002ff31760201c565b9050600060e0518262000133919062000312565b905060008160e0516200014791906200034b565b905062000155818462000365565b61010052506200039a95505050505050565b336001600160a01b03821603620001c15760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000086565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60004661a4b181148062000228575062066eed81145b156200029a5760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200026e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000294919062000380565b91505090565b4391505090565b80516001600160a01b0381168114620002b957600080fd5b919050565b60008060008060808587031215620002d557600080fd5b84519350620002e760208601620002a1565b9250620002f760408601620002a1565b91506200030760608601620002a1565b905092959194509250565b6000826200033057634e487b7160e01b600052601260045260246000fd5b500690565b634e487b7160e01b600052601160045260246000fd5b60008282101562000360576200036062000335565b500390565b600082198211156200037b576200037b62000335565b500190565b6000602082840312156200039357600080fd5b5051919050565b60805160a05160c05160e05161010051615c5e6200046b60003960006106d401526000818161069a01528181611ceb015281816138190152818161384801528181613880015261410901526000818161046c01528181610a4101528181610f080152818161148e0152818161190a01528181611b7e0152818161214b0152818161225d01526143200152600081816106160152613bde01526000818161034001528181610949015281816119d10152818161255a01528181612cba01528181612dbd0152612f3a0152615c5e6000f3fe608060405234801561001057600080fd5b50600436106102de5760003560e01c80638eef585f11610186578063c3fbb6fd116100e3578063cf7e754a11610097578063e72f6e3011610071578063e72f6e301461072d578063f2fde38b14610740578063f99b1d681461075357600080fd5b8063cf7e754a146106cf578063dac83d29146106f6578063dc311dd31461070957600080fd5b8063cb631797116100c8578063cb63179714610682578063cd0593df14610695578063ce3f4719146106bc57600080fd5b8063c3fbb6fd14610667578063c63c4e9b1461067a57600080fd5b8063abbf1c9b1161013a578063b2a7cac51161011f578063b2a7cac514610638578063bbcdd0d81461064b578063bec4c08c1461065457600080fd5b8063abbf1c9b146105f1578063ad1783611461061157600080fd5b80639e3616f41161016b5780639e3616f4146105c3578063a21a23e4146105d6578063a4c0ed36146105de57600080fd5b80638eef585f1461059d57806395009f08146105b057600080fd5b8063376126721161023f57806364d51a2a116101f357806385c64e11116101cd57806385c64e111461054f5780638a8a1390146105645780638da5cb5b1461058c57600080fd5b806364d51a2a1461051957806373433a2f1461053457806379ba50971461054757600080fd5b806346942d181161022457806346942d18146104bc578063483af703146104cf57806362f8b620146104ef57600080fd5b8063376126721461048e57806345ccbb8b146104a157600080fd5b80632a5e19f6116102965780632d9297b01161027b5780632d9297b0146104575780632f7527cc1461045f57806332fe7b261461046757600080fd5b80632a5e19f61461038f5780632b38bafc1461044457600080fd5b8063170b317a116102c7578063170b317a1461030b5780631b6b6d231461033b578063294daa491461037a57600080fd5b806305f4acc6146102e35780630ae09540146102f8575b600080fd5b6102f66102f1366004614712565b610766565b005b6102f661030636600461474d565b6107d5565b61031e6103193660046149aa565b610afb565b6040516001600160601b0390911681526020015b60405180910390f35b6103627f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610332565b60015b60405160ff9091168152602001610332565b61043761039d3660046149de565b6040805160808101825260008082526020820181905291810182905260608101919091525065ffffffffffff166000908152601160209081526040918290208251608081018452905463ffffffff81168252640100000000810462ffffff1692820192909252670100000000000000820461ffff1692810192909252690100000000000000000090046001600160a01b0316606082015290565b60405161033291906149f9565b6102f6610452366004614a3c565b610b63565b61031e610b8d565b61037d600881565b6103627f000000000000000000000000000000000000000000000000000000000000000081565b600b54610362906001600160a01b031681565b6002546001600160601b03165b604051908152602001610332565b6102f66104ca366004614a59565b610bd1565b6104e26104dd366004614abd565b610c38565b6040516103329190614bc1565b6105026104fd366004614be7565b610efb565b60405165ffffffffffff9091168152602001610332565b610521606481565b60405161ffff9091168152602001610332565b6102f6610542366004614c61565b611119565b6102f661130d565b6105576113be565b6040516103329190614cf4565b6104ae6105723660046149de565b65ffffffffffff166000908152600c602052604090205490565b6000546001600160a01b0316610362565b6102f66105ab366004614d03565b611423565b6105026105be366004614d40565b611481565b6102f66105d1366004614df3565b61164c565b6104ae6116b8565b6102f66105ec366004614e75565b6119a2565b6106046105ff366004614ec4565b611b71565b6040516103329190614f5d565b6103627f000000000000000000000000000000000000000000000000000000000000000081565b6102f6610646366004614f70565b611e21565b6104ae6103e881565b6102f661066236600461474d565b611f68565b6102f6610675366004614f89565b612113565b610521600381565b6102f661069036600461474d565b612729565b6104ae7f000000000000000000000000000000000000000000000000000000000000000081565b6102f66106ca366004614fdd565b612a3e565b6104ae7f000000000000000000000000000000000000000000000000000000000000000081565b6102f661070436600461474d565b612a70565b61071c610717366004614f70565b612b88565b60405161033295949392919061504b565b6102f661073b366004614a3c565b612c81565b6102f661074e366004614a3c565b612e75565b6102f66107613660046150a0565b612e86565b60045460ff16156107d25760408051808201909152436001600160401b039081168083529083166020909201829052600780547fffffffffffffffffffffffffffffffff0000000000000000000000000000000016909117680100000000000000009092029190911790555b50565b60008281526008602052604090205482906001600160a01b03168061080d57604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b0382161461084657604051636c51fda960e11b81526001600160a01b03821660048201526024015b60405180910390fd5b6000848152600960205260409020548490600160a01b90046001600160401b03161561088557604051631685ecdd60e31b815260040160405180910390fd5b60065460ff16156108a95760405163769dd35360e11b815260040160405180910390fd5b600085815260096020908152604091829020825160608101845290546001600160601b0381168083526001600160401b03600160601b8304811694840194909452600160a01b9091049092169281019290925261090587613077565b600280548291906000906109239084906001600160601b03166150e2565b92506101000a8154816001600160601b0302191690836001600160601b031602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb87836001600160601b03166040518363ffffffff1660e01b81526004016109b29291906001600160a01b03929092168252602082015260400190565b6020604051808303816000875af11580156109d1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109f59190615118565b610a1257604051631e9acf1760e31b815260040160405180910390fd5b6040517fbb9b2b34000000000000000000000000000000000000000000000000000000008152600481018890527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063bb9b2b3490602401600060405180830381600087803b158015610a8d57600080fd5b505af1158015610aa1573d6000803e3d6000fd5b5050604080516001600160a01b038a1681526001600160601b03851660208201528a93507f3784f77e8e883de95b5d47cd713ced01229fa74d118c0a462224bcb0516d43f19250015b60405180910390a250505050505050565b600080610b066131f4565b60045460808501516001600160401b039290921691610b39916a0100000000000000000000900463ffffffff1690615135565b610b439190615160565b6001600160601b03169050610b5c818460e001516132ad565b9392505050565b610b6b61330a565b600b80546001600160a01b0319166001600160a01b0392909216919091179055565b600080610b986131f4565b600454610bb591906601000000000000900463ffffffff1661518f565b6001600160401b03169050610bcb8160006132ad565b91505090565b610bd961330a565b6064610beb60408301602084016151c4565b60ff161115610c26576040517f0afa82a800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806004610c3382826151ee565b505050565b600b546060906001600160a01b03163314610c885760405162461bcd60e51b815260206004820152601060248201526f31b0b6361039b2ba283937b23ab1b2b960811b604482015260640161083d565b600080876001600160401b03811115610ca357610ca361477d565b604051908082528060200260200182016040528015610cf557816020015b604080516080810182526000808252602080830182905292820181905260608201528252600019909201910181610cc15790505b50905060005b88811015610de15760008a8a83818110610d1757610d17615384565b9050602002810190610d29919061539a565b610d329061549d565b9050610d3f81888b613366565b60408101515151151580610d5b57506040810151516020015115155b15610dce576040805160808101825282516001600160401b0316815260208084015162ffffff168183015283830180515151938301939093529151519091015160608201528351849084908110610db457610db4615384565b60200260200101819052508380610dca90615572565b9450505b5080610dd981615593565b915050610cfb565b5060008261ffff166001600160401b03811115610e0057610e0061477d565b604051908082528060200260200182016040528015610e5257816020015b604080516080810182526000808252602080830182905292820181905260608201528252600019909201910181610e1e5790505b50905060005b8361ffff16811015610eae57828181518110610e7657610e76615384565b6020026020010151828281518110610e9057610e90615384565b60200260200101819052508080610ea690615593565b915050610e58565b507fe1d18855b43b829b66a7f664301b0733507d67e5d8e163e3f0b778717e884ee086338a8a85604051610ee69594939291906155ac565b60405180910390a19998505050505050505050565b6000336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610f755760405162461bcd60e51b815260206004820152601560248201527f63616c6c6572206d75737420626520726f757465720000000000000000000000604482015260640161083d565b6000806000610f8589888861376e565b925092509250610f958989613aa7565b65ffffffffffff8316600090815260116020908152604080832085518154848801518489015160608a01516001600160a01b03166901000000000000000000027fffffff0000000000000000000000000000000000000000ffffffffffffffffff61ffff90921667010000000000000002919091167fffffff00000000000000000000000000000000000000000000ffffffffffffff62ffffff9093166401000000000266ffffffffffffff1990941663ffffffff909516949094179290921716919091171790558a83526009909152902080546001600160401b03600160601b9091041690600c61108683615607565b82546101009290920a6001600160401b0381810219909316918316021790915560408051918416825262ffffff8916602083015281018a905261ffff891660608201526001600160a01b038b16915065ffffffffffff8516907fc3b31df4232b05afd212fc28027dae6fd6a81618c2a3116182cb57c7f0a3fd0a9060800160405180910390a35090979650505050505050565b600b546001600160a01b031633146111665760405162461bcd60e51b815260206004820152601060248201526f31b0b6361039b2ba283937b23ab1b2b960811b604482015260640161083d565b82806111da5760405162461bcd60e51b815260206004820152602b60248201527f6e756d626572206f6620726563697069656e7473206d7573742062652067726560448201527f61746572207468616e2030000000000000000000000000000000000000000000606482015260840161083d565b601f81111561122b5760405162461bcd60e51b815260206004820152601360248201527f746f6f206d616e7920726563697069656e747300000000000000000000000000604482015260640161083d565b8082146112a05760405162461bcd60e51b815260206004820152603660248201527f6c656e677468206f6620726563697069656e747320616e64207061796d656e7460448201527f73496e4a75656c7320646964206e6f74206d6174636800000000000000000000606482015260840161083d565b60005b81811015611305576112f38686838181106112c0576112c0615384565b90506020020160208101906112d59190614a3c565b8585848181106112e7576112e7615384565b90506020020135612e86565b806112fd81615593565b9150506112a3565b505050505050565b6001546001600160a01b031633146113675760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161083d565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6113c66145b5565b6040805161010081019182905290601090600890826000855b82829054906101000a900462ffffff1662ffffff16815260200190600301906020826002010492830192600103820291508084116113df5790505050505050905090565b600b546001600160a01b031633146114705760405162461bcd60e51b815260206004820152601060248201526f31b0b6361039b2ba283937b23ab1b2b960811b604482015260640161083d565b61147d60108260086145d4565b5050565b6000336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146114fb5760405162461bcd60e51b815260206004820152601560248201527f63616c6c6572206d75737420626520726f757465720000000000000000000000604482015260640161083d565b6000806115098a898961376e565b925050915060006040518061010001604052808465ffffffffffff1681526020018a61ffff1681526020018c6001600160a01b031681526020018781526020018863ffffffff166001600160601b031681526020018b815260200161156c6131f4565b6001600160401b03168152602001611582613b96565b9052905061158f81613c87565b81888b836040516020016115a69493929190615670565b60408051601f19818403018152918152815160209283012065ffffffffffff86166000908152600c9093529120556115dd8a613d8c565b8a6001600160a01b03168365ffffffffffff167f24f0e469e0097d1e8d9975137f9f4dd17d2c1481b3a2f25f2382f51287eda1dc848b8e8e8d8860c001518960e001518f604051611635989796959493929190615736565b60405180910390a350909998505050505050505050565b61165461330a565b60005b81811015610c33576000600a600085858581811061167757611677615384565b905060200201602081019061168c9190614a3c565b6001600160a01b03168152602081019190915260400160002055806116b081615593565b915050611657565b60065460009060ff16156116df5760405163769dd35360e11b815260040160405180910390fd5b60006116ec60014361579b565b60015460408051924060208401523060601b6bffffffffffffffffffffffff191690830152600160a01b900460c01b7fffffffffffffffff000000000000000000000000000000000000000000000000166054820152605c0160408051808303601f19018152919052805160209091012060018054919250600160a01b9091046001600160401b031690601461178183615607565b91906101000a8154816001600160401b0302191690836001600160401b03160217905550506000806001600160401b038111156117c0576117c061477d565b6040519080825280602002602001820160405280156117e9578160200160208202803683370190505b506040805160608082018352600080835260208084018281528486018381528984526009835286842095518654925191516001600160601b039091166001600160a01b031993841617600160601b6001600160401b0393841602177fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff16600160a01b9290911691909102179094558451928301855233835282810182815283860187815289845260088352959092208351815486166001600160a01b0391821617825592516001820180549096169316929092179093559251805194955090936118d99260028501920190614672565b50506040517f9479b74e000000000000000000000000000000000000000000000000000000008152600481018490527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169150639479b74e90602401600060405180830381600087803b15801561195757600080fd5b505af115801561196b573d6000803e3d6000fd5b50506040513392508491507f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d90600090a350905090565b60065460ff16156119c65760405163769dd35360e11b815260040160405180910390fd5b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614611a28576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114611a4957604051638129bbcd60e01b815260040160405180910390fd5b6000611a5782840184614f70565b6000818152600860205260409020549091506001600160a01b0316611a8f57604051630fb532db60e11b815260040160405180910390fd5b600081815260096020526040812080546001600160601b031691869190611ab68385615135565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600260008282829054906101000a90046001600160601b0316611afe9190615135565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a828784611b5191906157b2565b6040805192835260208301919091520160405180910390a2505050505050565b6060336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614611beb5760405162461bcd60e51b815260206004820152601560248201527f63616c6c6572206d75737420626520726f757465720000000000000000000000604482015260640161083d565b65ffffffffffff831660008181526011602081815260408084208151608081018352815463ffffffff8116825262ffffff6401000000008204168286015261ffff670100000000000000820416938201939093526001600160a01b036901000000000000000000840481166060830190815297909652939092527fffffff00000000000000000000000000000000000000000000000000000000001690559151811690871614611ce05760608101516040517f8e30e8230000000000000000000000000000000000000000000000000000000081526001600160a01b039182166004820152908716602482015260440161083d565b8051600090611d16907f00000000000000000000000000000000000000000000000000000000000000009063ffffffff166157ca565b90506000611d22612ff3565b90506000836020015162ffffff1682611d3b919061579b565b9050808310611da05782846020015162ffffff1684611d5a91906157b2565b611d659060016157b2565b6040517f15ad27c30000000000000000000000000000000000000000000000000000000081526004810192909252602482015260440161083d565b6001600160401b03831115611de4576040517f058ddf020000000000000000000000000000000000000000000000000000000081526004810184905260240161083d565b6000838152600e602090815260408083208288015162ffffff168452909152902054611e14908890869086613e19565b9998505050505050505050565b60065460ff1615611e455760405163769dd35360e11b815260040160405180910390fd5b6000818152600860205260409020546001600160a01b0316611e7a57604051630fb532db60e11b815260040160405180910390fd5b6000818152600860205260409020600101546001600160a01b03163314611eec57600081815260086020526040908190206001015490517fd084e9750000000000000000000000000000000000000000000000000000000081526001600160a01b03909116600482015260240161083d565b6000818152600860209081526040918290208054336001600160a01b0319808316821784556001909301805490931690925583516001600160a01b0390911680825292810191909152909183917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c9386910160405180910390a25050565b60008281526008602052604090205482906001600160a01b031680611fa057604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b03821614611fd457604051636c51fda960e11b81526001600160a01b038216600482015260240161083d565b60065460ff1615611ff85760405163769dd35360e11b815260040160405180910390fd5b60008481526008602052604090206002015460631901612044576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b03831660009081526003602090815260408083208784529091528120546001600160401b0316900361210d576001600160a01b03831660008181526003602090815260408083208884528252808320805467ffffffffffffffff19166001908117909155600883528184206002018054918201815584529282902090920180546001600160a01b03191684179055905191825285917f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e191015b60405180910390a25b50505050565b6040517fcef40aa60000000000000000000000000000000000000000000000000000000081526001600160a01b0384811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063cef40aa690602401602060405180830381865afa158015612192573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121b69190615118565b6121ec576040517f3fae142000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040811461220d57604051638129bbcd60e01b815260040160405180910390fd5b600061221b828401846157e9565b60208101516040517ffa66358a000000000000000000000000000000000000000000000000000000008152600481019190915290915030906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063fa66358a90602401602060405180830381865afa1580156122a4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122c89190615823565b6001600160a01b0316146122ef57604051630fb532db60e11b815260040160405180910390fd5b60008060008060006123048660200151612b88565b94509450945094509450816001600160a01b0316336001600160a01b03161461234b57604051636c51fda960e11b81526001600160a01b038316600482015260240161083d565b886001600160a01b031663294daa496040518163ffffffff1660e01b8152600401602060405180830381865afa158015612389573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123ad9190615840565b60ff16866000015160ff1614612463578560000151896001600160a01b031663294daa496040518163ffffffff1660e01b8152600401602060405180830381865afa158015612400573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124249190615840565b6040517fe7aada9500000000000000000000000000000000000000000000000000000000815260ff92831660048201529116602482015260440161083d565b6001600160401b0383161561248b57604051631685ecdd60e31b815260040160405180910390fd5b60006040518060c001604052806124a0600190565b60ff16815260200188602001518152602001846001600160a01b03168152602001838152602001876001600160601b03168152602001866001600160401b031681525090506000816040516020016124f8919061585d565b60405160208183030381529060405290506125168860200151613077565b600280548891906000906125349084906001600160601b03166150e2565b92506101000a8154816001600160601b0302191690836001600160601b031602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb8c896040518363ffffffff1660e01b81526004016125c39291906001600160a01b039290921682526001600160601b0316602082015260400190565b6020604051808303816000875af11580156125e2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126069190615118565b6126525760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015260640161083d565b6040517fce3f47190000000000000000000000000000000000000000000000000000000081526001600160a01b038c169063ce3f4719906126979084906004016158d5565b600060405180830381600087803b1580156126b157600080fd5b505af11580156126c5573d6000803e3d6000fd5b505050508760200151886000015160ff167fbd89b747474d3fc04664dfbd1d56ae7ffbe46ee097cdb9979c13916bb76269ce8d60405161271491906001600160a01b0391909116815260200190565b60405180910390a35050505050505050505050565b60008281526008602052604090205482906001600160a01b03168061276157604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b0382161461279557604051636c51fda960e11b81526001600160a01b038216600482015260240161083d565b6000848152600960205260409020548490600160a01b90046001600160401b0316156127d457604051631685ecdd60e31b815260040160405180910390fd5b60065460ff16156127f85760405163769dd35360e11b815260040160405180910390fd5b6001600160a01b03841660009081526003602090815260408083208884529091528120546001600160401b03169003612856576040516379bfd40160e01b8152600481018690526001600160a01b038516602482015260440161083d565b6000858152600860209081526040808320600201805482518185028101850190935280835291929091908301828280156128b957602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161289b575b505050505090506000600182516128d0919061579b565b905060005b82518110156129db57866001600160a01b03168382815181106128fa576128fa615384565b60200260200101516001600160a01b0316036129c957600083838151811061292457612924615384565b6020026020010151905080600860008b8152602001908152602001600020600201838154811061295657612956615384565b600091825260208083209190910180546001600160a01b0319166001600160a01b0394909416939093179092558a81526008909152604090206002018054806129a1576129a16158e8565b600082815260209020810160001990810180546001600160a01b0319169055019055506129db565b806129d381615593565b9150506128d5565b506001600160a01b03861660008181526003602090815260408083208b8452825291829020805467ffffffffffffffff19169055905191825288917f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a79101610aea565b6040517f2cb6686f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281526008602052604090205482906001600160a01b031680612aa857604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b03821614612adc57604051636c51fda960e11b81526001600160a01b038216600482015260240161083d565b60065460ff1615612b005760405163769dd35360e11b815260040160405180910390fd5b6000848152600860205260409020600101546001600160a01b0384811691161461210d5760008481526008602090815260409182902060010180546001600160a01b0319166001600160a01b03871690811790915582513381529182015285917f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a19101612104565b6000818152600860205260408120548190819081906060906001600160a01b0316612bc657604051630fb532db60e11b815260040160405180910390fd5b60008681526009602090815260408083205460088352928190208054600290910180548351818602810186019094528084526001600160601b038616956001600160401b03600160601b8204811696600160a01b90920416946001600160a01b0390941693918391830182828015612c6757602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612c49575b505050505090509450945094509450945091939590929450565b612c8961330a565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906370a0823190602401602060405180830381865afa158015612d09573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d2d91906158fe565b6002549091506001600160601b031681811115612d80576040517fa99da302000000000000000000000000000000000000000000000000000000008152600481018290526024810183905260440161083d565b81811015610c33576000612d94828461579b565b60405163a9059cbb60e01b81526001600160a01b038681166004830152602482018390529192507f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb906044016020604051808303816000875af1158015612e08573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e2c9190615118565b50604080516001600160a01b0386168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a150505050565b612e7d61330a565b6107d281613fd6565b600b546001600160a01b03163314612ed35760405162461bcd60e51b815260206004820152601060248201526f31b0b6361039b2ba283937b23ab1b2b960811b604482015260640161083d565b60028054829190600090612ef19084906001600160601b03166150e2565b82546001600160601b039182166101009390930a92830291909202199091161790555060405163a9059cbb60e01b81526001600160a01b038381166004830152602482018390527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af1158015612f83573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612fa79190615118565b61147d5760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015260640161083d565b60004661a4b1811480613008575062066eed81145b156130705760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561304c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bcb91906158fe565b4391505090565b6000818152600860209081526040808320815160608101835281546001600160a01b039081168252600183015416818501526002820180548451818702810187018652818152929593948601938301828280156130fd57602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116130df575b505050505081525050905060005b81604001515181101561318357600360008360400151838151811061313257613132615384565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081208682529092529020805467ffffffffffffffff191690558061317b81615593565b91505061310b565b50600082815260086020526040812080546001600160a01b031990811682556001820180549091169055906131bb60028301826146c7565b505050600090815260096020526040902080547fffffffff00000000000000000000000000000000000000000000000000000000169055565b60045460009060ff16801561321f57506007546801000000000000000090046001600160401b031615155b156132a8576004547201000000000000000000000000000000000000900463ffffffff16431080806132845750600454613275907201000000000000000000000000000000000000900463ffffffff164361579b565b6007546001600160401b031610155b156132a65750506007546801000000000000000090046001600160401b031690565b505b503a90565b60045460009081906064906132da906e010000000000000000000000000000900463ffffffff1682615917565b6132ea9063ffffffff16866157ca565b6132f4919061594c565b9050613300818461407f565b9150505b92915050565b6000546001600160a01b031633146133645760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161083d565b565b82516001600160401b03808416911611156133c35782516040517f012d824d0000000000000000000000000000000000000000000000000000000081526001600160401b038085166004830152909116602482015260440161083d565b604083015151516000901580156133e1575060408401515160200151155b15613419575082516001600160401b03166000908152600e602090815260408083208287015162ffffff168452909152902054613473565b836040015160405160200161342e9190615960565b60408051601f19818403018152918152815160209283012086516001600160401b03166000908152600e84528281208885015162ffffff168252909352912081905590505b6060840151516000816001600160401b038111156134935761349361477d565b6040519080825280602002602001820160405280156134bc578160200160208202803683370190505b5090506000826001600160401b038111156134d9576134d961477d565b6040519080825280601f01601f191660200182016040528015613503576020820181803683370190505b5090506000836001600160401b038111156135205761352061477d565b60405190808252806020026020018201604052801561355357816020015b606081526020019060019003908161353e5790505b5090506000805b8581101561366c5760008a60600151828151811061357a5761357a615384565b6020026020010151905060008061359b8d600001518e602001518c866140ff565b9150915081156135da5780868661ffff16815181106135bc576135bc615384565b602002602001018190525084806135d290615572565b955050613621565b600160f81b8785815181106135f1576135f1615384565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505b825151885189908690811061363857613638615384565b602002602001019065ffffffffffff16908165ffffffffffff1681525050505050808061366490615593565b91505061355a565b50606089015151156137635760008161ffff166001600160401b038111156136965761369661477d565b6040519080825280602002602001820160405280156136c957816020015b60608152602001906001900390816136b45790505b50905060005b8261ffff16811015613725578381815181106136ed576136ed615384565b602002602001015182828151811061370757613707615384565b6020026020010181905250808061371d90615593565b9150506136cf565b507f47ddf7bb0cbd94c1b43c5097f1352a80db0ceb3696f029d32b24f32cd631d2b785858360405161375993929190615993565b60405180910390a1505b505050505050505050565b604080516080810182526000808252602082018190529181018290526060810182905260006103e88561ffff1611156137c857604051634a90778560e01b815261ffff861660048201526103e8602482015260440161083d565b8461ffff16600003613806576040517f08fad2a700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000613810612ff3565b9050600061383e7f000000000000000000000000000000000000000000000000000000000000000083615a49565b905060008161386d7f0000000000000000000000000000000000000000000000000000000000000000856157b2565b613877919061579b565b905060006138a57f00000000000000000000000000000000000000000000000000000000000000008361594c565b905063ffffffff81106138e4576040517f7b2a523000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080518082018252600f805465ffffffffffff16825282516101008101938490528493600093929160208401916010906008908288855b82829054906101000a900462ffffff1662ffffff168152602001906003019060208260020104928301926001038202915080841161391c57905050505091909252505081519192505065ffffffffffff908116106139a6576040517f2b4655b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80516139b3906001615a5d565b600f805465ffffffffffff191665ffffffffffff9290921691909117905560005b6008811015613a1a578a62ffffff16826020015182600881106139f9576139f9615384565b602002015162ffffff1614613a1a5780613a1281615593565b9150506139d4565b60088110613a5b5760208201516040517fc4f769b000000000000000000000000000000000000000000000000000000000815261083d918d91600401615a7e565b506040805160808101825263ffffffff909316835262ffffff8b16602084015261ffff8c16908301526001600160a01b038c166060830152519750955090935050505093509350939050565b6001600160a01b03821660009081526003602090815260408083208484529091528120546001600160401b03169003613b05576040516379bfd40160e01b8152600481018290526001600160a01b038316602482015260440161083d565b6000613b0f610b8d565b600083815260096020526040902080546001600160601b039283169350909116821115613b4f57604051631e9acf1760e31b815260040160405180910390fd5b805482908290600090613b6c9084906001600160601b03166150e2565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555050505050565b60048054604080517ffeaf968c00000000000000000000000000000000000000000000000000000000815290516000936201000090930463ffffffff169283151592859283927f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169263feaf968c928183019260a0928290030181865afa158015613c2d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613c519190615ab2565b509450909250849150508015613c755750613c6c824261579b565b8463ffffffff16105b15613c7f57506005545b949350505050565b6040808201516001600160a01b031660009081526003602090815282822060a085015183529052908120546001600160401b03169003613cf65760a081015160408083015190516379bfd40160e01b815260048101929092526001600160a01b0316602482015260440161083d565b6000613d0182610afb565b60a0830151600090815260096020526040902080546001600160601b039283169350909116821115613d4657604051631e9acf1760e31b815260040160405180910390fd5b805482908290600090613d639084906001600160601b03166150e2565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550505050565b60008181526009602052604090208054600160601b90046001600160401b031681600c613db883615607565b82546101009290920a6001600160401b038181021990931691831602179091558254600160a01b9004169050816014613df083615607565b91906101000a8154816001600160401b0302191690836001600160401b03160217905550505050565b606082613e6b576040517fc7d41b1b00000000000000000000000000000000000000000000000000000000815265ffffffffffff861660048201526001600160401b038316602482015260440161083d565b6000858585604051602001613e8293929190615b02565b6040516020818303038152906040528051906020012090506103e8856040015161ffff161115613ed9576040808601519051634a90778560e01b815261ffff90911660048201526103e8602482015260440161083d565b6000856040015161ffff166001600160401b03811115613efb57613efb61477d565b604051908082528060200260200182016040528015613f24578160200160208202803683370190505b50905060005b866040015161ffff168161ffff161015613fcb578281604051602001613f7f92919091825260f01b7fffff00000000000000000000000000000000000000000000000000000000000016602082015260220190565b6040516020818303038152906040528051906020012060001c828261ffff1681518110613fae57613fae615384565b602090810291909101015280613fc381615572565b915050613f2a565b509695505050505050565b336001600160a01b0382160361402e5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161083d565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600080821561408e5782614096565b614096613b96565b90506000816140ad86670de0b6b3a76400006157ca565b6140b7919061594c565b90506b033b2e3c9fd0803ce8000000811115613300576040517fde43710000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006060816141377f00000000000000000000000000000000000000000000000000000000000000006001600160401b03891661594c565b845160a0810151604051929350909160009161415b918b918b918690602001615670565b60408051601f198184030181529181528151602092830120845165ffffffffffff166000908152600c90935291205490915081146141d75760016040518060400160405280601081526020017f756e6b6e6f776e2063616c6c6261636b00000000000000000000000000000000815250945094505050506144f6565b6040805160808101825263ffffffff8516815262ffffff8a1660208083019190915284015161ffff1681830152908301516001600160a01b03166060820152825160009061422790838b8e613e19565b60608084015186519187015160405193945090926000927f5a47dd71000000000000000000000000000000000000000000000000000000009261426f92879190602401615b65565b60408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526006805460ff191660011790559050600080805a8d51608001516040808c015190517f55fe97630000000000000000000000000000000000000000000000000000000081529293506001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016926355fe9763926143579290918990600401615ba2565b60408051808303816000875af1158015614375573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906143999190615bdc565b9093509150816143eb5760405162461bcd60e51b815260206004820152601060248201527f696e73756666696369656e742067617300000000000000000000000000000000604482015260640161083d565b60006113885a6143fb91906157b2565b6006805460ff191690559050818110156144235761442361441c828461579b565b8f516144ff565b8d5160a0015160009081526009602052604090208054600160a01b90046001600160401b031690601461445583615c0b565b82546001600160401b039182166101009390930a928302919092021990911617905550895165ffffffffffff166000908152600c6020526040812055836144d35760016040518060400160405280601081526020017f657865637574696f6e206661696c6564000000000000000000000000000000008152506144e6565b6000604051806020016040528060008152505b9c509c5050505050505050505050505b94509492505050565b80608001516001600160601b0316821115614518575050565b60045460009060649061453390610100900460ff1682615c2e565b60ff168360c001518585608001516001600160601b0316614554919061579b565b61455e91906157ca565b61456891906157ca565b614572919061594c565b90506000614584828460e0015161407f565b60a0840151600090815260096020526040812080549293508392909190613b6c9084906001600160601b0316615135565b6040518061010001604052806008906020820280368337509192915050565b6001830191839082156146625791602002820160005b8382111561463157833562ffffff1683826101000a81548162ffffff021916908362ffffff16021790555092602001926003016020816002010492830192600103026145ea565b80156146605782816101000a81549062ffffff0219169055600301602081600201049283019260010302614631565b505b5061466e9291506146e1565b5090565b828054828255906000526020600020908101928215614662579160200282015b8281111561466257825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190614692565b50805460008255906000526020600020908101906107d291905b5b8082111561466e57600081556001016146e2565b80356001600160401b038116811461470d57600080fd5b919050565b60006020828403121561472457600080fd5b610b5c826146f6565b6001600160a01b03811681146107d257600080fd5b803561470d8161472d565b6000806040838503121561476057600080fd5b8235915060208301356147728161472d565b809150509250929050565b634e487b7160e01b600052604160045260246000fd5b60405161010081016001600160401b03811182821017156147b6576147b661477d565b60405290565b604080519081016001600160401b03811182821017156147b6576147b661477d565b604051608081016001600160401b03811182821017156147b6576147b661477d565b604051602081016001600160401b03811182821017156147b6576147b661477d565b604051601f8201601f191681016001600160401b038111828210171561484a5761484a61477d565b604052919050565b803565ffffffffffff8116811461470d57600080fd5b803561ffff8116811461470d57600080fd5b600082601f83011261488b57600080fd5b81356001600160401b038111156148a4576148a461477d565b6148b7601f8201601f1916602001614822565b8181528460208386010111156148cc57600080fd5b816020850160208301376000918101602001919091529392505050565b80356001600160601b038116811461470d57600080fd5b6000610100828403121561491357600080fd5b61491b614793565b905061492682614852565b815261493460208301614868565b602082015261494560408301614742565b604082015260608201356001600160401b0381111561496357600080fd5b61496f8482850161487a565b606083015250614981608083016148e9565b608082015260a082013560a082015260c082013560c082015260e082013560e082015292915050565b6000602082840312156149bc57600080fd5b81356001600160401b038111156149d257600080fd5b61330084828501614900565b6000602082840312156149f057600080fd5b610b5c82614852565b815163ffffffff16815260208083015162ffffff169082015260408083015161ffff16908201526060808301516001600160a01b03169082015260808101613304565b600060208284031215614a4e57600080fd5b8135610b5c8161472d565b60006101008284031215614a6c57600080fd5b50919050565b60008083601f840112614a8457600080fd5b5081356001600160401b03811115614a9b57600080fd5b6020830191508360208260051b8501011115614ab657600080fd5b9250929050565b60008060008060008060a08789031215614ad657600080fd5b86356001600160401b03811115614aec57600080fd5b614af889828a01614a72565b909750955050602087013577ffffffffffffffffffffffffffffffffffffffffffffffff81168114614b2957600080fd5b9350614b37604088016146f6565b9250614b45606088016146f6565b9150608087013590509295509295509295565b600081518084526020808501945080840160005b83811015614bb657815180516001600160401b031688528381015162ffffff1684890152604080820151908901526060908101519088015260809096019590820190600101614b6c565b509495945050505050565b602081526000610b5c6020830184614b58565b803562ffffff8116811461470d57600080fd5b600080600080600060a08688031215614bff57600080fd5b8535614c0a8161472d565b945060208601359350614c1f60408701614868565b9250614c2d60608701614bd4565b915060808601356001600160401b03811115614c4857600080fd5b614c548882890161487a565b9150509295509295909350565b60008060008060408587031215614c7757600080fd5b84356001600160401b0380821115614c8e57600080fd5b614c9a88838901614a72565b90965094506020870135915080821115614cb357600080fd5b50614cc087828801614a72565b95989497509550505050565b8060005b600881101561210d57815162ffffff16845260209384019390910190600101614cd0565b61010081016133048284614ccc565b6000610100808385031215614d1757600080fd5b838184011115614d2657600080fd5b509092915050565b63ffffffff811681146107d257600080fd5b600080600080600080600060e0888a031215614d5b57600080fd5b8735614d668161472d565b965060208801359550614d7b60408901614868565b9450614d8960608901614bd4565b93506080880135614d9981614d2e565b925060a08801356001600160401b0380821115614db557600080fd5b614dc18b838c0161487a565b935060c08a0135915080821115614dd757600080fd5b50614de48a828b0161487a565b91505092959891949750929550565b60008060208385031215614e0657600080fd5b82356001600160401b03811115614e1c57600080fd5b614e2885828601614a72565b90969095509350505050565b60008083601f840112614e4657600080fd5b5081356001600160401b03811115614e5d57600080fd5b602083019150836020828501011115614ab657600080fd5b60008060008060608587031215614e8b57600080fd5b8435614e968161472d565b93506020850135925060408501356001600160401b03811115614eb857600080fd5b614cc087828801614e34565b60008060008060808587031215614eda57600080fd5b8435614ee58161472d565b935060208501359250614efa60408601614852565b915060608501356001600160401b03811115614f1557600080fd5b614f218782880161487a565b91505092959194509250565b600081518084526020808501945080840160005b83811015614bb657815187529582019590820190600101614f41565b602081526000610b5c6020830184614f2d565b600060208284031215614f8257600080fd5b5035919050565b600080600060408486031215614f9e57600080fd5b8335614fa98161472d565b925060208401356001600160401b03811115614fc457600080fd5b614fd086828701614e34565b9497909650939450505050565b60008060208385031215614ff057600080fd5b82356001600160401b0381111561500657600080fd5b614e2885828601614e34565b600081518084526020808501945080840160005b83811015614bb65781516001600160a01b031687529582019590820190600101615026565b6001600160601b038616815260006001600160401b0380871660208401528086166040840152506001600160a01b038416606083015260a0608083015261509560a0830184615012565b979650505050505050565b600080604083850312156150b357600080fd5b82356150be8161472d565b946020939093013593505050565b634e487b7160e01b600052601160045260246000fd5b60006001600160601b0383811690831681811015615102576151026150cc565b039392505050565b80151581146107d257600080fd5b60006020828403121561512a57600080fd5b8151610b5c8161510a565b60006001600160601b03808316818516808303821115615157576151576150cc565b01949350505050565b60006001600160601b0380831681851681830481118215151615615186576151866150cc565b02949350505050565b60006001600160401b0380831681851681830481118215151615615186576151866150cc565b60ff811681146107d257600080fd5b6000602082840312156151d657600080fd5b8135610b5c816151b5565b6000813561330481614d2e565b81356151f98161510a565b815460ff19811691151560ff1691821783556020840135615219816151b5565b61ff008160081b168361ffff1984161717845550505061525c61523e604084016151e1565b825465ffffffff0000191660109190911b65ffffffff000016178255565b61529161526b606084016151e1565b825469ffffffff000000000000191660309190911b69ffffffff00000000000016178255565b6152ce6152a0608084016151e1565b82546dffffffff00000000000000000000191660509190911b6dffffffff0000000000000000000016178255565b6153206152dd60a084016151e1565b82547fffffffffffffffffffffffffffff00000000ffffffffffffffffffffffffffff1660709190911b71ffffffff000000000000000000000000000016178255565b61537661532f60c084016151e1565b82547fffffffffffffffffffff00000000ffffffffffffffffffffffffffffffffffff1660909190911b75ffffffff00000000000000000000000000000000000016178255565b60e082013560018201555050565b634e487b7160e01b600052603260045260246000fd5b60008235609e198336030181126153b057600080fd5b9190910192915050565b600082601f8301126153cb57600080fd5b813560206001600160401b03808311156153e7576153e761477d565b8260051b6153f6838201614822565b938452858101830193838101908886111561541057600080fd5b84880192505b858310156154915782358481111561542e5760008081fd5b88016040818b03601f19018113156154465760008081fd5b61544e6147bc565b87830135878111156154605760008081fd5b61546e8d8a83870101614900565b82525061547c8284016148e9565b81890152845250509184019190840190615416565b98975050505050505050565b600081360360a08112156154b057600080fd5b6154b86147de565b6154c1846146f6565b815260206154d0818601614bd4565b828201526040603f19840112156154e657600080fd5b6154ee614800565b925036605f8601126154ff57600080fd5b6155076147bc565b80608087013681111561551957600080fd5b604088015b81811015615535578035845292840192840161551e565b50908552604084019490945250509035906001600160401b0382111561555a57600080fd5b615566368386016153ba565b60608201529392505050565b600061ffff808316818103615589576155896150cc565b6001019392505050565b6000600182016155a5576155a56150cc565b5060010190565b60006001600160401b0380881683526001600160a01b038716602084015277ffffffffffffffffffffffffffffffffffffffffffffffff8616604084015280851660608401525060a0608083015261509560a0830184614b58565b60006001600160401b03808316818103615589576155896150cc565b6000815180845260005b818110156156495760208185018101518683018201520161562d565b8181111561565b576000602083870101525b50601f01601f19169290920160200192915050565b6001600160401b038516815262ffffff841660208201528260408201526080606082015265ffffffffffff8251166080820152600060208301516156ba60a084018261ffff169052565b5060408301516001600160a01b031660c0830152606083015161010060e084018190526156eb610180850183615623565b91506080850151615706828601826001600160601b03169052565b505060a084015161012084015260c084015161014084015260e08401516101608401528091505095945050505050565b60006101006001600160401b038b16835262ffffff8a16602084015288604084015261ffff8816606084015263ffffffff871660808401528560a08401528460c08401528060e084015261578c81840185615623565b9b9a5050505050505050505050565b6000828210156157ad576157ad6150cc565b500390565b600082198211156157c5576157c56150cc565b500190565b60008160001904831182151516156157e4576157e46150cc565b500290565b6000604082840312156157fb57600080fd5b6158036147bc565b823561580e816151b5565b81526020928301359281019290925250919050565b60006020828403121561583557600080fd5b8151610b5c8161472d565b60006020828403121561585257600080fd5b8151610b5c816151b5565b6020815260ff8251166020820152602082015160408201526001600160a01b0360408301511660608201526000606083015160c060808401526158a360e0840182615012565b90506001600160601b0360808501511660a08401526001600160401b0360a08501511660c08401528091505092915050565b602081526000610b5c6020830184615623565b634e487b7160e01b600052603160045260246000fd5b60006020828403121561591057600080fd5b5051919050565b600063ffffffff808316818516808303821115615157576151576150cc565b634e487b7160e01b600052601260045260246000fd5b60008261595b5761595b615936565b500490565b815160408201908260005b600281101561598a57825182526020928301929091019060010161596b565b50505092915050565b606080825284519082018190526000906020906080840190828801845b828110156159d457815165ffffffffffff16845292840192908401906001016159b0565b505050838103828501526159e88187615623565b905083810360408501528085518083528383019150838160051b84010184880160005b83811015615a3957601f19868403018552615a27838351615623565b94870194925090860190600101615a0b565b50909a9950505050505050505050565b600082615a5857615a58615936565b500690565b600065ffffffffffff808316818516808303821115615157576151576150cc565b62ffffff831681526101208101610b5c6020830184614ccc565b805169ffffffffffffffffffff8116811461470d57600080fd5b600080600080600060a08688031215615aca57600080fd5b615ad386615a98565b9450602086015193506040860151925060608601519150615af660808701615a98565b90509295509295909350565b65ffffffffffff8416815260c08101615b57602083018563ffffffff815116825262ffffff602082015116602083015261ffff60408201511660408301526001600160a01b0360608201511660608301525050565b8260a0830152949350505050565b65ffffffffffff84168152606060208201526000615b866060830185614f2d565b8281036040840152615b988185615623565b9695505050505050565b6001600160601b03841681526001600160a01b0383166020820152606060408201526000615bd36060830184615623565b95945050505050565b60008060408385031215615bef57600080fd5b8251615bfa8161510a565b60208401519092506147728161510a565b60006001600160401b03821680615c2457615c246150cc565b6000190192915050565b600060ff821660ff841680821015615c4857615c486150cc565b9003939250505056fea164736f6c634300080f000a",
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
	parsed, err := abi.JSON(strings.NewReader(VRFCoordinatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

func (_VRFCoordinator *VRFCoordinatorCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) LINK() (common.Address, error) {
	return _VRFCoordinator.Contract.LINK(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) LINK() (common.Address, error) {
	return _VRFCoordinator.Contract.LINK(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCaller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) LINKETHFEED() (common.Address, error) {
	return _VRFCoordinator.Contract.LINKETHFEED(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) LINKETHFEED() (common.Address, error) {
	return _VRFCoordinator.Contract.LINKETHFEED(&_VRFCoordinator.CallOpts)
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

func (_VRFCoordinator *VRFCoordinatorCaller) ROUTER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "ROUTER")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) ROUTER() (common.Address, error) {
	return _VRFCoordinator.Contract.ROUTER(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) ROUTER() (common.Address, error) {
	return _VRFCoordinator.Contract.ROUTER(&_VRFCoordinator.CallOpts)
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

func (_VRFCoordinator *VRFCoordinatorCaller) GetBeaconRequest(opts *bind.CallOpts, requestId *big.Int) (VRFBeaconTypesBeaconRequest, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "getBeaconRequest", requestId)

	if err != nil {
		return *new(VRFBeaconTypesBeaconRequest), err
	}

	out0 := *abi.ConvertType(out[0], new(VRFBeaconTypesBeaconRequest)).(*VRFBeaconTypesBeaconRequest)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) GetBeaconRequest(requestId *big.Int) (VRFBeaconTypesBeaconRequest, error) {
	return _VRFCoordinator.Contract.GetBeaconRequest(&_VRFCoordinator.CallOpts, requestId)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) GetBeaconRequest(requestId *big.Int) (VRFBeaconTypesBeaconRequest, error) {
	return _VRFCoordinator.Contract.GetBeaconRequest(&_VRFCoordinator.CallOpts, requestId)
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

func (_VRFCoordinator *VRFCoordinatorCaller) IStartSlot(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "i_StartSlot")

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

func (_VRFCoordinator *VRFCoordinatorCaller) MaxNumWords(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "maxNumWords")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) MaxNumWords() (*big.Int, error) {
	return _VRFCoordinator.Contract.MaxNumWords(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) MaxNumWords() (*big.Int, error) {
	return _VRFCoordinator.Contract.MaxNumWords(&_VRFCoordinator.CallOpts)
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

func (_VRFCoordinator *VRFCoordinatorCaller) Producer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinator.contract.Call(opts, &out, "producer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinator *VRFCoordinatorSession) Producer() (common.Address, error) {
	return _VRFCoordinator.Contract.Producer(&_VRFCoordinator.CallOpts)
}

func (_VRFCoordinator *VRFCoordinatorCallerSession) Producer() (common.Address, error) {
	return _VRFCoordinator.Contract.Producer(&_VRFCoordinator.CallOpts)
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

func (_VRFCoordinator *VRFCoordinatorTransactor) ForgetConsumerSubscriptionID(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "forgetConsumerSubscriptionID", consumers)
}

func (_VRFCoordinator *VRFCoordinatorSession) ForgetConsumerSubscriptionID(consumers []common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.ForgetConsumerSubscriptionID(&_VRFCoordinator.TransactOpts, consumers)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) ForgetConsumerSubscriptionID(consumers []common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.ForgetConsumerSubscriptionID(&_VRFCoordinator.TransactOpts, consumers)
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

func (_VRFCoordinator *VRFCoordinatorTransactor) SetProducer(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.contract.Transact(opts, "setProducer", addr)
}

func (_VRFCoordinator *VRFCoordinatorSession) SetProducer(addr common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetProducer(&_VRFCoordinator.TransactOpts, addr)
}

func (_VRFCoordinator *VRFCoordinatorTransactorSession) SetProducer(addr common.Address) (*types.Transaction, error) {
	return _VRFCoordinator.Contract.SetProducer(&_VRFCoordinator.TransactOpts, addr)
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
	Transmitter        common.Address
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

type GetSubscription struct {
	Balance             *big.Int
	ReqCount            uint64
	PendingFulfillments uint64
	Owner               common.Address
	Consumers           []common.Address
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
	return common.HexToHash("0xe1d18855b43b829b66a7f664301b0733507d67e5d8e163e3f0b778717e884ee0")
}

func (VRFCoordinatorOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFCoordinatorOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
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

func (_VRFCoordinator *VRFCoordinator) Address() common.Address {
	return _VRFCoordinator.address
}

type VRFCoordinatorInterface interface {
	LINK(opts *bind.CallOpts) (common.Address, error)

	LINKETHFEED(opts *bind.CallOpts) (common.Address, error)

	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	NUMCONFDELAYS(opts *bind.CallOpts) (uint8, error)

	ROUTER(opts *bind.CallOpts) (common.Address, error)

	CalculateRequestPriceCallbackJuels(opts *bind.CallOpts, callback VRFBeaconTypesCallback) (*big.Int, error)

	CalculateRequestPriceJuels(opts *bind.CallOpts) (*big.Int, error)

	GetBeaconRequest(opts *bind.CallOpts, requestId *big.Int) (VRFBeaconTypesBeaconRequest, error)

	GetCallbackMemo(opts *bind.CallOpts, requestId *big.Int) ([32]byte, error)

	GetConfirmationDelays(opts *bind.CallOpts) ([8]*big.Int, error)

	GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

		error)

	GetTotalLinkBalance(opts *bind.CallOpts) (*big.Int, error)

	IStartSlot(opts *bind.CallOpts) (*big.Int, error)

	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)

	MaxNumWords(opts *bind.CallOpts) (*big.Int, error)

	MigrationVersion(opts *bind.CallOpts) (uint8, error)

	MinDelay(opts *bind.CallOpts) (uint16, error)

	OnMigration(opts *bind.CallOpts, arg0 []byte) error

	Owner(opts *bind.CallOpts) (common.Address, error)

	Producer(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	BatchTransferLink(opts *bind.TransactOpts, recipients []common.Address, paymentsInJuels []*big.Int) (*types.Transaction, error)

	CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	ForgetConsumerSubscriptionID(opts *bind.TransactOpts, consumers []common.Address) (*types.Transaction, error)

	Migrate(opts *bind.TransactOpts, newCoordinator common.Address, encodedRequest []byte) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	ProcessVRFOutputs(opts *bind.TransactOpts, vrfOutputs []VRFBeaconTypesVRFOutput, juelsPerFeeCoin *big.Int, reasonableGasPrice uint64, blockHeight uint64, arg4 [32]byte) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RedeemRandomness(opts *bind.TransactOpts, sender common.Address, arg1 *big.Int, requestID *big.Int, arg3 []byte) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	RequestRandomness(opts *bind.TransactOpts, requester common.Address, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, arg4 []byte) (*types.Transaction, error)

	RequestRandomnessFulfillment(opts *bind.TransactOpts, requester common.Address, subID *big.Int, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte, arg6 []byte) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error)

	SetBillingConfig(opts *bind.TransactOpts, billingConfig VRFBeaconTypesBillingConfig) (*types.Transaction, error)

	SetConfirmationDelays(opts *bind.TransactOpts, confDelays [8]*big.Int) (*types.Transaction, error)

	SetProducer(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error)

	SetReasonableGasPrice(opts *bind.TransactOpts, gasPrice uint64) (*types.Transaction, error)

	TransferLink(opts *bind.TransactOpts, recipient common.Address, juelsAmount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

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

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
