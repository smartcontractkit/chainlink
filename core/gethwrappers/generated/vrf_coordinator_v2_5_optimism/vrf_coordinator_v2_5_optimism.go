// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_coordinator_v2_5_optimism

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

type VRFProof struct {
	Pk            [2]*big.Int
	Gamma         [2]*big.Int
	C             *big.Int
	S             *big.Int
	Seed          *big.Int
	UWitness      common.Address
	CGammaWitness [2]*big.Int
	SHashWitness  [2]*big.Int
	ZInv          *big.Int
}

type VRFTypesRequestCommitmentV2Plus struct {
	BlockNum         uint64
	SubId            *big.Int
	CallbackGasLimit uint32
	NumWords         uint32
	Sender           common.Address
	ExtraArgs        []byte
}

type VRFV2PlusClientRandomWordsRequest struct {
	KeyHash              [32]byte
	SubId                *big.Int
	RequestConfirmations uint16
	CallbackGasLimit     uint32
	NumWords             uint32
	ExtraArgs            []byte
}

var VRFCoordinatorV25OptimismMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendNative\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToTransferLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxGas\",\"type\":\"uint256\"}],\"name\":\"GasPriceExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"premiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"max\",\"type\":\"uint8\"}],\"name\":\"InvalidPremiumPercentage\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"flatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeNativePPM\",\"type\":\"uint32\"}],\"name\":\"LinkDiscountTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"max\",\"type\":\"uint32\"}],\"name\":\"MsgDataTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"name\":\"FallbackWeiPerUnitLinkUsed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativeFundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"onlyPremium\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountNative\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldNativeBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newNativeBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFundedWithNative\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_NATIVE_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"deregisterMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFTypes.RequestCommitmentV2Plus\",\"name\":\"rc\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"onlyPremium\",\"type\":\"bool\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"fundSubscriptionWithNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveSubscriptionIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"subOwner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverNativeFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"registerMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"reentrancyLock\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_currentSubNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_optimismGasModule\",\"outputs\":[{\"internalType\":\"contractIGasModule\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_provingKeyHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_provingKeys\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestCommitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalNativeBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkNativeFeed\",\"type\":\"address\"}],\"name\":\"setLINKAndLINKNativeFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"module\",\"type\":\"address\"}],\"name\":\"setOptimismGasModule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162005deb38038062005deb8339810160408190526200003491620001b7565b8133806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620000ef565b5050506001600160a01b03908116608052601280546001600160a01b0319169290911691909117905550620001ef565b336001600160a01b03821603620001495760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001b257600080fd5b919050565b60008060408385031215620001cb57600080fd5b620001d6836200019a565b9150620001e6602084016200019a565b90509250929050565b608051615bd9620002126000396000818161062801526135640152615bd96000f3fe6080604052600436106102c25760003560e01c80638402595e1161017f578063b2a7cac5116100e1578063dac83d291161008a578063e72f6e3011610064578063e72f6e301461097a578063ee9d2d381461099a578063f2fde38b146109c757600080fd5b8063dac83d2914610909578063dc311dd314610929578063e38af0a61461095a57600080fd5b8063cb631797116100bb578063cb6317971461086a578063d98e620e1461088a578063da2f2610146108aa57600080fd5b8063b2a7cac51461080a578063bec4c08c1461082a578063caf70c4a1461084a57600080fd5b80639d40a6fd11610143578063a63e0bfb1161011d578063a63e0bfb1461079d578063aa433aff146107bd578063aefb212f146107dd57600080fd5b80639d40a6fd14610730578063a21a23e414610768578063a4c0ed361461077d57600080fd5b80638402595e1461069f57806386fe91c7146106bf5780638da5cb5b146106df57806395b55cfc146106fd5780639b1c385e1461071057600080fd5b8063405b84fa1161022857806364d51a2a116101ec57806372e9d565116101c657806372e9d5651461064a57806379ba50971461066a5780637a5a2aef1461067f57600080fd5b806364d51a2a146105e157806365982744146105f6578063689c45171461061657600080fd5b8063405b84fa1461052657806340d6bb821461054657806341af6c871461057157806351cff8d9146105a15780635d06b4ab146105c157600080fd5b806315c48b841161028a5780632f622e6b116102645780632f622e6b146104c6578063301f42e9146104e6578063383d19a31461050657600080fd5b806315c48b841461042757806318e3dd271461044f5780631b6b6d231461048e57600080fd5b806304104edb146102c7578063043bd6ae146102e9578063088070f51461031257806308821d58146103e75780630ae0954014610407575b600080fd5b3480156102d357600080fd5b506102e76102e2366004614ebc565b6109e7565b005b3480156102f557600080fd5b506102ff60105481565b6040519081526020015b60405180910390f35b34801561031e57600080fd5b50600c5461038a9061ffff81169063ffffffff62010000820481169160ff660100000000000082048116926701000000000000008304811692600160581b8104821692600160781b8204831692600160981b83041691600160b81b8104821691600160c01b9091041689565b6040805161ffff909a168a5263ffffffff98891660208b01529615159689019690965293861660608801529185166080870152841660a08601529290921660c084015260ff91821660e08401521661010082015261012001610309565b3480156103f357600080fd5b506102e7610402366004614eea565b610b60565b34801561041357600080fd5b506102e7610422366004614f06565b610d1d565b34801561043357600080fd5b5061043c60c881565b60405161ffff9091168152602001610309565b34801561045b57600080fd5b50600a5461047690600160601b90046001600160601b031681565b6040516001600160601b039091168152602001610309565b34801561049a57600080fd5b506002546104ae906001600160a01b031681565b6040516001600160a01b039091168152602001610309565b3480156104d257600080fd5b506102e76104e1366004614ebc565b610d65565b3480156104f257600080fd5b50610476610501366004615168565b610eb4565b34801561051257600080fd5b506102e7610521366004614ebc565b6111ca565b34801561053257600080fd5b506102e7610541366004614f06565b6111f4565b34801561055257600080fd5b5061055c6101f481565b60405163ffffffff9091168152602001610309565b34801561057d57600080fd5b5061059161058c366004615256565b6115d6565b6040519015158152602001610309565b3480156105ad57600080fd5b506102e76105bc366004614ebc565b61168a565b3480156105cd57600080fd5b506102e76105dc366004614ebc565b61180c565b3480156105ed57600080fd5b5061043c606481565b34801561060257600080fd5b506102e761061136600461526f565b6118ca565b34801561062257600080fd5b506104ae7f000000000000000000000000000000000000000000000000000000000000000081565b34801561065657600080fd5b506003546104ae906001600160a01b031681565b34801561067657600080fd5b506102e761192a565b34801561068b57600080fd5b506102e761069a36600461529d565b6119db565b3480156106ab57600080fd5b506102e76106ba366004614ebc565b611b0f565b3480156106cb57600080fd5b50600a54610476906001600160601b031681565b3480156106eb57600080fd5b506000546001600160a01b03166104ae565b6102e761070b366004615256565b611c2a565b34801561071c57600080fd5b506102ff61072b3660046152d1565b611d4e565b34801561073c57600080fd5b50600754610750906001600160401b031681565b6040516001600160401b039091168152602001610309565b34801561077457600080fd5b506102ff61218d565b34801561078957600080fd5b506102e761079836600461530b565b612374565b3480156107a957600080fd5b506102e76107b83660046153b6565b6124f0565b3480156107c957600080fd5b506102e76107d8366004615256565b6127d7565b3480156107e957600080fd5b506107fd6107f8366004615457565b61281f565b60405161030991906154b4565b34801561081657600080fd5b506102e7610825366004615256565b612921565b34801561083657600080fd5b506102e7610845366004614f06565b612a25565b34801561085657600080fd5b506102ff6108653660046154c7565b612b18565b34801561087657600080fd5b506102e7610885366004614f06565b612b48565b34801561089657600080fd5b506102ff6108a5366004615256565b612db6565b3480156108b657600080fd5b506108ea6108c5366004615256565b600d6020526000908152604090205460ff81169061010090046001600160401b031682565b6040805192151583526001600160401b03909116602083015201610309565b34801561091557600080fd5b506102e7610924366004614f06565b612dd7565b34801561093557600080fd5b50610949610944366004615256565b612e71565b60405161030995949392919061551c565b34801561096657600080fd5b506012546104ae906001600160a01b031681565b34801561098657600080fd5b506102e7610995366004614ebc565b612f5f565b3480156109a657600080fd5b506102ff6109b5366004615256565b600f6020526000908152604090205481565b3480156109d357600080fd5b506102e76109e2366004614ebc565b613120565b6109ef613131565b60115460005b81811015610b3357826001600160a01b031660118281548110610a1a57610a1a615571565b6000918252602090912001546001600160a01b031603610b23576011610a4160018461559d565b81548110610a5157610a51615571565b600091825260209091200154601180546001600160a01b039092169183908110610a7d57610a7d615571565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b031602179055506011805480610abc57610abc6155b0565b6000828152602090819020600019908301810180546001600160a01b03191690559091019091556040516001600160a01b03851681527ff80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af3791015b60405180910390a1505050565b610b2c816155c6565b90506109f5565b50604051635428d44960e01b81526001600160a01b03831660048201526024015b60405180910390fd5b50565b610b68613131565b604080518082018252600091610b97919084906002908390839080828437600092019190915250612b18915050565b6000818152600d602090815260409182902082518084019093525460ff811615158084526101009091046001600160401b03169183019190915291925090610bf557604051631dfd6e1360e21b815260048101839052602401610b54565b6000828152600d60205260408120805468ffffffffffffffffff19169055600e54905b81811015610cc75783600e8281548110610c3457610c34615571565b906000526020600020015403610cb757600e610c5160018461559d565b81548110610c6157610c61615571565b9060005260206000200154600e8281548110610c7f57610c7f615571565b600091825260209091200155600e805480610c9c57610c9c6155b0565b60019003818190600052602060002001600090559055610cc7565b610cc0816155c6565b9050610c18565b507f9b6868e0eb737bcd72205360baa6bfd0ba4e4819a33ade2db384e8a8025639a5838360200151604051610d0f9291909182526001600160401b0316602082015260400190565b60405180910390a150505050565b81610d278161318d565b610d2f6131f7565b610d38836115d6565b15610d5657604051631685ecdd60e31b815260040160405180910390fd5b610d608383613225565b505050565b610d6d6131f7565b610d75613131565b600b54600160601b90046001600160601b0316600003610da857604051631e9acf1760e31b815260040160405180910390fd5b600b8054600160601b90046001600160601b0316908190600c610dcb83806155df565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600a600c8282829054906101000a90046001600160601b0316610e1391906155df565b92506101000a8154816001600160601b0302191690836001600160601b031602179055506000826001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d8060008114610e8d576040519150601f19603f3d011682016040523d82523d6000602084013e610e92565b606091505b5050905080610d605760405163950b247960e01b815260040160405180910390fd5b6000610ebe6131f7565b60005a9050610324361115610ef057604051630f28961b60e01b81523660048201526103246024820152604401610b54565b6000610efc86866133cb565b90506000610f1285836000015160200151613677565b60408301516060888101519293509163ffffffff16806001600160401b03811115610f3f57610f3f614f36565b604051908082528060200260200182016040528015610f68578160200160208202803683370190505b50925060005b81811015610fd05760408051602081018590529081018290526060016040516020818303038152906040528051906020012060001c848281518110610fb557610fb5615571565b6020908102919091010152610fc9816155c6565b9050610f6e565b5050602080850180516000908152600f9092526040822082905551610ff6908a856136d2565b60208a810151600090815260069091526040902080549192509060189061102c90600160c01b90046001600160401b03166155ff565b82546101009290920a6001600160401b0381810219909316918316021790915560808a01516001600160a01b03166000908152600460209081526040808320828e0151845290915290208054909160099161108f91600160481b90910416615625565b91906101000a8154816001600160401b0302191690836001600160401b0316021790555060008960a0015160018b60a00151516110cc919061559d565b815181106110dc576110dc615571565b60209101015160f81c600114905060006110f88887848d613776565b909950905080156111435760208088015160105460408051928352928201527f6ca648a381f22ead7e37773d934e64885dcf861fbfbb26c40354cbf0c4662d1a910160405180910390a15b5061115388828c602001516137ae565b6020808b015187820151604080518781526001600160601b038d16948101949094528415159084015284151560608401528b1515608084015290917faeb4b4786571e184246d39587f659abf0e26f41f6a3358692250382c0cdb47b79060a00160405180910390a3505050505050505b9392505050565b6111d2613131565b601280546001600160a01b0319166001600160a01b0392909216919091179055565b6111fc6131f7565b6112058161391b565b61122d57604051635428d44960e01b81526001600160a01b0382166004820152602401610b54565b60008060008061123c86612e71565b945094505093509350336001600160a01b0316826001600160a01b0316146112a65760405162461bcd60e51b815260206004820152601660248201527f4e6f7420737562736372697074696f6e206f776e6572000000000000000000006044820152606401610b54565b6112af866115d6565b156112fc5760405162461bcd60e51b815260206004820152601660248201527f50656e64696e67207265717565737420657869737473000000000000000000006044820152606401610b54565b6040805160c0810182526001815260208082018990526001600160a01b03851682840152606082018490526001600160601b038088166080840152861660a08301529151909160009161135191849101615648565b604051602081830303815290604052905061136b88613986565b505060405163ce3f471960e01b81526001600160a01b0388169063ce3f4719906001600160601b038816906113a4908590600401615703565b6000604051808303818588803b1580156113bd57600080fd5b505af11580156113d1573d6000803e3d6000fd5b50506002546001600160a01b0316158015935091506113fa905057506001600160601b03861615155b156114ca5760025460405163a9059cbb60e01b81526001600160a01b0389811660048301526001600160601b03891660248301529091169063a9059cbb906044016020604051808303816000875af115801561145a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061147e9190615716565b6114ca5760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610b54565b600c805466ff0000000000001916660100000000000017905560005b83518110156115795783818151811061150157611501615571565b6020908102919091010151604051638ea9811760e01b81526001600160a01b038a8116600483015290911690638ea9811790602401600060405180830381600087803b15801561155057600080fd5b505af1158015611564573d6000803e3d6000fd5b5050505080611572906155c6565b90506114e6565b50600c805466ff00000000000019169055604080516001600160a01b0389168152602081018a90527fd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be4187910160405180910390a15050505050505050565b600081815260056020526040812060020180548083036115fa575060009392505050565b60005b8181101561167f5760006004600085848154811061161d5761161d615571565b60009182526020808320909101546001600160a01b0316835282810193909352604091820181208982529092529020546001600160401b03600160481b90910416111561166f57506001949350505050565b611678816155c6565b90506115fd565b506000949350505050565b6116926131f7565b61169a613131565b6002546001600160a01b03166116c35760405163c1f0c0a160e01b815260040160405180910390fd5b600b546001600160601b03166000036116ef57604051631e9acf1760e31b815260040160405180910390fd5b600b80546001600160601b0316908190600061170b83806155df565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600a60008282829054906101000a90046001600160601b031661175391906155df565b82546101009290920a6001600160601b0381810219909316918316021790915560025460405163a9059cbb60e01b81526001600160a01b03868116600483015292851660248201529116915063a9059cbb906044016020604051808303816000875af11580156117c7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117eb9190615716565b61180857604051631e9acf1760e31b815260040160405180910390fd5b5050565b611814613131565b61181d8161391b565b156118465760405163ac8a27ef60e01b81526001600160a01b0382166004820152602401610b54565b601180546001810182556000919091527f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c680180546001600160a01b0319166001600160a01b0383169081179091556040519081527fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af016259060200160405180910390a150565b6118d2613131565b6002546001600160a01b0316156118fc57604051631688c53760e11b815260040160405180910390fd5b600280546001600160a01b039384166001600160a01b03199182161790915560038054929093169116179055565b6001546001600160a01b031633146119845760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610b54565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6119e3613131565b604080518082018252600091611a12919085906002908390839080828437600092019190915250612b18915050565b6000818152600d602052604090205490915060ff1615611a4857604051634a0b8fa760e01b815260048101829052602401610b54565b60408051808201825260018082526001600160401b0385811660208085018281526000888152600d835287812096518754925168ffffffffffffffffff1990931690151568ffffffffffffffff00191617610100929095169190910293909317909455600e805493840181559091527fbb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fd9091018490558251848152918201527f9b911b2c240bfbef3b6a8f7ed6ee321d1258bb2a3fe6becab52ac1cd3210afd39101610b16565b611b17613131565b600a544790600160601b90046001600160601b031681811115611b57576040516354ced18160e11b81526004810182905260248101839052604401610b54565b81811015610d60576000611b6b828461559d565b90506000846001600160a01b03168260405160006040518083038185875af1925050503d8060008114611bba576040519150601f19603f3d011682016040523d82523d6000602084013e611bbf565b606091505b5050905080611be15760405163950b247960e01b815260040160405180910390fd5b604080516001600160a01b0387168152602081018490527f4aed7c8eed0496c8c19ea2681fcca25741c1602342e38b045d9f1e8e905d2e9c910160405180910390a15050505050565b611c326131f7565b6000818152600560205260409020546001600160a01b0316611c6757604051630fb532db60e11b815260040160405180910390fd5b60008181526006602052604090208054600160601b90046001600160601b0316903490600c611c968385615733565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555034600a600c8282829054906101000a90046001600160601b0316611cde9190615733565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f7603b205d03651ee812f803fccde89f1012e545a9c99f0abfea9cedd0fd8e902823484611d319190615753565b604080519283526020830191909152015b60405180910390a25050565b6000611d586131f7565b602080830135600081815260059092526040909120546001600160a01b0316611d9457604051630fb532db60e11b815260040160405180910390fd5b336000908152600460209081526040808320848452808352928190208151606081018352905460ff811615158083526001600160401b036101008304811695840195909552600160481b9091049093169181019190915290611e12576040516379bfd40160e01b815260048101849052336024820152604401610b54565b600c5461ffff16611e296060870160408801615766565b61ffff161080611e4c575060c8611e466060870160408801615766565b61ffff16115b15611e9257611e616060860160408701615766565b600c5460405163539c34bb60e11b815261ffff92831660048201529116602482015260c86044820152606401610b54565b600c5462010000900463ffffffff16611eb16080870160608801615781565b63ffffffff161115611f0157611ecd6080860160608701615781565b600c54604051637aebf00f60e11b815263ffffffff9283166004820152620100009091049091166024820152604401610b54565b6101f4611f1460a0870160808801615781565b63ffffffff161115611f5a57611f3060a0860160808701615781565b6040516311ce1afb60e21b815263ffffffff90911660048201526101f46024820152604401610b54565b806020018051611f69906155ff565b6001600160401b03169052604081018051611f83906155ff565b6001600160401b03908116909152602082810151604080518935818501819052338284015260608201899052929094166080808601919091528151808603909101815260a08501825280519084012060c085019290925260e0808501839052815180860390910181526101009094019052825192909101919091206000919095509050600061202561202061201b60a08a018a61579c565b613b38565b613bb9565b905085438661203a60808b0160608c01615781565b61204a60a08c0160808d01615781565b338660405160200161206297969594939291906157e9565b60405160208183030381529060405280519060200120600f600088815260200190815260200160002081905550336001600160a01b03168588600001357feb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e89868c60400160208101906120d59190615766565b8d60600160208101906120e89190615781565b8e60800160208101906120fb9190615781565b8960405161210e96959493929190615840565b60405180910390a45050600092835260209182526040928390208151815493830151929094015168ffffffffffffffffff1990931693151568ffffffffffffffff001916939093176101006001600160401b03928316021770ffffffffffffffff0000000000000000001916600160481b91909216021790555b919050565b60006121976131f7565b6007546001600160401b0316336121af60014361559d565b6040516bffffffffffffffffffffffff19606093841b81166020830152914060348201523090921b1660548201526001600160c01b031960c083901b16606882015260700160408051601f198184030181529190528051602090910120915061221981600161587f565b6007805467ffffffffffffffff19166001600160401b03928316179055604080516000808252608082018352602080830182815283850183815260608086018581528a86526006855287862093518454935191516001600160601b039182166001600160c01b031990951694909417600160601b91909216021777ffffffffffffffffffffffffffffffffffffffffffffffff16600160c01b9290981691909102969096179055835194850184523385528481018281528585018481528884526005835294909220855181546001600160a01b03199081166001600160a01b0392831617835593516001830180549095169116179092559251805192949391926123299260028501920190614daa565b5061233991506008905084613c2a565b5060405133815283907f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d9060200160405180910390a2505090565b61237c6131f7565b6002546001600160a01b031633146123a7576040516344b0e3c360e01b815260040160405180910390fd5b602081146123c857604051638129bbcd60e01b815260040160405180910390fd5b60006123d682840184615256565b6000818152600560205260409020549091506001600160a01b031661240e57604051630fb532db60e11b815260040160405180910390fd5b600081815260066020526040812080546001600160601b0316918691906124358385615733565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600a60008282829054906101000a90046001600160601b031661247d9190615733565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a8287846124d09190615753565b6040805192835260208301919091520160405180910390a2505050505050565b6124f8613131565b60c861ffff8a1611156125325760405163539c34bb60e11b815261ffff8a1660048201819052602482015260c86044820152606401610b54565b60008513612556576040516321ea67b360e11b815260048101869052602401610b54565b8363ffffffff168363ffffffff161115612593576040516313c06e5960e11b815263ffffffff808516600483015285166024820152604401610b54565b609b60ff831611156125c457604051631d66288d60e11b815260ff83166004820152609b6024820152604401610b54565b609b60ff821611156125f557604051631d66288d60e11b815260ff82166004820152609b6024820152604401610b54565b604080516101208101825261ffff8b1680825263ffffffff808c16602084018190526000848601528b8216606085018190528b8316608086018190528a841660a08701819052938a1660c0870181905260ff808b1660e08901819052908a16610100909801889052600c8054600160c01b90990260ff60c01b19600160b81b9093029290921661ffff60b81b19600160981b90940263ffffffff60981b19600160781b9099029890981676ffffffffffffffff00000000000000000000000000000019600160581b9096026effffffff000000000000000000000019670100000000000000909802979097166effffffffffffffffff000000000000196201000090990265ffffffffffff19909c16909a179a909a1796909616979097179390931791909116959095179290921793909316929092179190911790556010869055517f2c6b6b12413678366b05b145c5f00745bdd00e739131ab5de82484a50c9d78b6906127c4908b908b908b908b908b908b908b908b908b9061ffff99909916895263ffffffff97881660208a0152958716604089015293861660608801526080870192909252841660a086015290921660c084015260ff91821660e0840152166101008201526101200190565b60405180910390a1505050505050505050565b6127df613131565b6000818152600560205260409020546001600160a01b03168061281557604051630fb532db60e11b815260040160405180910390fd5b6118088282613225565b6060600061282d6008613c36565b905080841061284f57604051631390f2a160e01b815260040160405180910390fd5b600061285b8486615753565b905081811180612869575083155b6128735780612875565b815b90506000612883868361559d565b9050806001600160401b0381111561289d5761289d614f36565b6040519080825280602002602001820160405280156128c6578160200160208202803683370190505b50935060005b81811015612916576128e96128e18883615753565b600890613c40565b8582815181106128fb576128fb615571565b602090810291909101015261290f816155c6565b90506128cc565b505050505b92915050565b6129296131f7565b6000818152600560205260409020546001600160a01b03168061295f57604051630fb532db60e11b815260040160405180910390fd5b6000828152600560205260409020600101546001600160a01b031633146129b8576000828152600560205260409081902060010154905163d084e97560e01b81526001600160a01b039091166004820152602401610b54565b6000828152600560209081526040918290208054336001600160a01b03199182168117835560019092018054909116905582516001600160a01b03851681529182015283917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c93869101611d42565b81612a2f8161318d565b612a376131f7565b6001600160a01b03821660009081526004602090815260408083208684529091529020805460ff1615612a6a5750505050565b6000848152600560205260409020600201805460631901612a9e576040516305a48e0f60e01b815260040160405180910390fd5b8154600160ff199091168117835581549081018255600082815260209081902090910180546001600160a01b0319166001600160a01b03871690811790915560405190815286917f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e191015b60405180910390a25050505050565b600081604051602001612b2b91906158c2565b604051602081830303815290604052805190602001209050919050565b81612b528161318d565b612b5a6131f7565b612b63836115d6565b15612b8157604051631685ecdd60e31b815260040160405180910390fd5b6001600160a01b038216600090815260046020908152604080832086845290915290205460ff16612bd7576040516379bfd40160e01b8152600481018490526001600160a01b0383166024820152604401610b54565b600083815260056020908152604080832060020180548251818502810185019093528083529192909190830182828015612c3a57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612c1c575b50505050509050600060018251612c51919061559d565b905060005b8251811015612d5a57846001600160a01b0316838281518110612c7b57612c7b615571565b60200260200101516001600160a01b031603612d4a576000838381518110612ca557612ca5615571565b6020026020010151905080600560008981526020019081526020016000206002018381548110612cd757612cd7615571565b600091825260208083209190910180546001600160a01b0319166001600160a01b039490941693909317909255888152600590915260409020600201805480612d2257612d226155b0565b600082815260209020810160001990810180546001600160a01b031916905501905550612d5a565b612d53816155c6565b9050612c56565b506001600160a01b0384166000818152600460209081526040808320898452825291829020805460ff19169055905191825286917f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a79101612b09565b600e8181548110612dc657600080fd5b600091825260209091200154905081565b81612de18161318d565b612de96131f7565b600083815260056020526040902060018101546001600160a01b03848116911614612e6b576001810180546001600160a01b0319166001600160a01b03851690811790915560408051338152602081019290925285917f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a1910160405180910390a25b50505050565b600081815260056020526040812054819081906001600160a01b0316606081612ead57604051630fb532db60e11b815260040160405180910390fd5b600086815260066020908152604080832054600583529281902060020180548251818502810185019093528083526001600160601b0380861695600160601b810490911694600160c01b9091046001600160401b0316938893929091839190830182828015612f4557602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612f27575b505050505090509450945094509450945091939590929450565b612f67613131565b6002546001600160a01b0316612f905760405163c1f0c0a160e01b815260040160405180910390fd5b6002546040516370a0823160e01b81523060048201526000916001600160a01b0316906370a0823190602401602060405180830381865afa158015612fd9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ffd91906158d0565b600a549091506001600160601b031681811115613037576040516354ced18160e11b81526004810182905260248101839052604401610b54565b81811015610d6057600061304b828461559d565b60025460405163a9059cbb60e01b81526001600160a01b0387811660048301526024820184905292935091169063a9059cbb906044016020604051808303816000875af11580156130a0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130c49190615716565b6130e157604051631f01ff1360e21b815260040160405180910390fd5b604080516001600160a01b0386168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b4366009101610d0f565b613128613131565b610b5d81613c4c565b6000546001600160a01b0316331461318b5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610b54565b565b6000818152600560205260409020546001600160a01b0316806131c357604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b0382161461180857604051636c51fda960e11b81526001600160a01b0382166004820152602401610b54565b600c546601000000000000900460ff161561318b5760405163769dd35360e11b815260040160405180910390fd5b60008061323184613986565b60025491935091506001600160a01b03161580159061325857506001600160601b03821615155b156132f95760025460405163a9059cbb60e01b81526001600160a01b0385811660048301526001600160601b03851660248301529091169063a9059cbb906044016020604051808303816000875af11580156132b8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906132dc9190615716565b6132f957604051631e9acf1760e31b815260040160405180910390fd5b6000836001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d806000811461334f576040519150601f19603f3d011682016040523d82523d6000602084013e613354565b606091505b50509050806133765760405163950b247960e01b815260040160405180910390fd5b604080516001600160a01b03861681526001600160601b03808616602083015284169181019190915285907f8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c490606001612b09565b6040805160a081018252600060608201818152608083018290528252602082018190529181019190915260006134048460000151612b18565b6000818152600d602090815260409182902082518084019093525460ff811615158084526101009091046001600160401b0316918301919091529192509061346257604051631dfd6e1360e21b815260048101839052602401610b54565b6000828660800151604051602001613484929190918252602082015260400190565b60408051601f1981840301815291815281516020928301206000818152600f90935290822054909250908190036134ce57604051631b44092560e11b815260040160405180910390fd5b85516020808801516040808a015160608b015160808c015160a08d015193516134fd978a9790969591016158e9565b6040516020818303038152906040528051906020012081146135325760405163354a450b60e21b815260040160405180910390fd5b85516001600160401b03164080613605578651604051631d2827a760e31b81526001600160401b0390911660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063e9413d3890602401602060405180830381865afa1580156135b3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135d791906158d0565b90508061360557865160405163175dadad60e01b81526001600160401b039091166004820152602401610b54565b6000886080015182604051602001613627929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c9050600061364e8a83613cf5565b604080516060810182529788526020880196909652948601949094525092979650505050505050565b6000816001600160401b03163a11156136ca5782156136a057506001600160401b03811661291b565b60405163435e532d60e11b81523a60048201526001600160401b0383166024820152604401610b54565b503a92915050565b6000806000631fe543e360e01b86856040516024016136f292919061593c565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252600c805466ff0000000000001916660100000000000017905590860151608087015191925061375c9163ffffffff9091169083613d60565b600c805466ff000000000000191690559695505050505050565b60008083156137955761378a868685613dac565b6000915091506137a5565b6137a0868685613e88565b915091505b94509492505050565b600081815260066020526040902082156138825780546001600160601b03600160601b90910481169085168110156137f957604051631e9acf1760e31b815260040160405180910390fd5b61380385826155df565b82547fffffffffffffffff000000000000000000000000ffffffffffffffffffffffff16600160601b6001600160601b039283168102919091178455600b805488939192600c92613858928692900416615733565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555050612e6b565b80546001600160601b039081169085168110156138b257604051631e9acf1760e31b815260040160405180910390fd5b6138bc85826155df565b82546bffffffffffffffffffffffff19166001600160601b03918216178355600b805487926000916138f091859116615733565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505050505050565b601154600090815b8181101561397c57836001600160a01b03166011828154811061394857613948615571565b6000918252602090912001546001600160a01b03160361396c575060019392505050565b613975816155c6565b9050613923565b5060009392505050565b60008181526005602090815260408083206006909252822054600290910180546001600160601b0380841694600160601b90940416925b81811015613a3257600460008483815481106139db576139db615571565b60009182526020808320909101546001600160a01b0316835282810193909352604091820181208982529092529020805470ffffffffffffffffffffffffffffffffff19169055613a2b816155c6565b90506139bd565b50600085815260056020526040812080546001600160a01b03199081168255600182018054909116905590613a6a6002830182614e0f565b5050600085815260066020526040812055613a86600886614045565b506001600160601b03841615613ad957600a8054859190600090613ab49084906001600160601b03166155df565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b6001600160601b03831615613b315782600a600c8282829054906101000a90046001600160601b0316613b0c91906155df565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b5050915091565b6040805160208101909152600081526000829003613b65575060408051602081019091526000815261291b565b63125fa26760e31b613b77838561595d565b6001600160e01b03191614613b9f57604051632923fee760e11b815260040160405180910390fd5b613bac826004818661598d565b8101906111c391906159b7565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa82604051602401613bf291511515815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b03199093169290921790915292915050565b60006111c38383614051565b600061291b825490565b60006111c383836140a0565b336001600160a01b03821603613ca45760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610b54565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000613d298360000151846020015185604001518660600151868860a001518960c001518a60e001518b61010001516140ca565b60038360200151604051602001613d41929190615a02565b60408051601f1981840301815291905280516020909101209392505050565b60005a611388811015613d7257600080fd5b611388810390508460408204820311613d8a57600080fd5b50823b613d9657600080fd5b60008083516020850160008789f1949350505050565b600080613dba6000366142f5565b905060005a600c54613dda908890600160581b900463ffffffff16615753565b613de4919061559d565b613dee9086615a16565b600c54909150600090613e1390600160781b900463ffffffff1664e8d4a51000615a16565b90508415613e5f57600c548190606490600160b81b900460ff16613e378587615753565b613e419190615a16565b613e4b9190615a43565b613e559190615753565b93505050506111c3565b600c548190606490613e7b90600160b81b900460ff1682615a57565b60ff16613e378587615753565b600080600080613e96614369565b9150915060008213613ebe576040516321ea67b360e11b815260048101839052602401610b54565b6000613ecb6000366142f5565b9050600083825a600c54613eed908d90600160581b900463ffffffff16615753565b613ef7919061559d565b613f01908b615a16565b613f0b9190615753565b613f1d90670de0b6b3a7640000615a16565b613f279190615a43565b600c54909150600090613f509063ffffffff600160981b8204811691600160781b900416615a70565b613f659063ffffffff1664e8d4a51000615a16565b9050600085613f7c83670de0b6b3a7640000615a16565b613f869190615a43565b905060008915613fc757600c548290606490613fac90600160c01b900460ff1687615a16565b613fb69190615a43565b613fc09190615753565b9050614007565b600c548290606490613fe390600160c01b900460ff1682615a57565b613ff09060ff1687615a16565b613ffa9190615a43565b6140049190615753565b90505b6b033b2e3c9fd0803ce80000008111156140345760405163e80fa38160e01b815260040160405180910390fd5b9b949a509398505050505050505050565b60006111c38383614434565b60008181526001830160205260408120546140985750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561291b565b50600061291b565b60008260000182815481106140b7576140b7615571565b9060005260206000200154905092915050565b6140d38961452e565b61411f5760405162461bcd60e51b815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e2063757276650000000000006044820152606401610b54565b6141288861452e565b6141745760405162461bcd60e51b815260206004820152601560248201527f67616d6d61206973206e6f74206f6e20637572766500000000000000000000006044820152606401610b54565b61417d8361452e565b6141c95760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610b54565b6141d28261452e565b61421e5760405162461bcd60e51b815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e206375727665000000006044820152606401610b54565b61422a878a8887614607565b6142765760405162461bcd60e51b815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e657373000000000000006044820152606401610b54565b60006142828a8761472a565b90506000614295898b878b86898961478e565b905060006142a6838d8d8a866148ba565b9050808a146142e75760405162461bcd60e51b815260206004820152600d60248201526c34b73b30b634b210383937b7b360991b6044820152606401610b54565b505050505050505050505050565b6012546040516302c6933760e01b81526000916001600160a01b0316906302c69337906143289086908690600401615a8d565b602060405180830381865afa158015614345573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111c391906158d0565b600c5460035460408051633fabe5a360e21b81529051600093849367010000000000000090910463ffffffff169284926001600160a01b039092169163feaf968c9160048082019260a0929091908290030181865afa1580156143d0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906143f49190615ad6565b50919650909250505063ffffffff8216158015906144205750614417814261559d565b8263ffffffff16105b9250821561442e5760105493505b50509091565b6000818152600183016020526040812054801561451d57600061445860018361559d565b855490915060009061446c9060019061559d565b90508181146144d157600086600001828154811061448c5761448c615571565b90600052602060002001549050808760000184815481106144af576144af615571565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806144e2576144e26155b0565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061291b565b600091505061291b565b5092915050565b80516000906401000003d019116145875760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420782d6f7264696e61746500000000000000000000000000006044820152606401610b54565b60208201516401000003d019116145e05760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420792d6f7264696e61746500000000000000000000000000006044820152606401610b54565b60208201516401000003d0199080096146008360005b60200201516148fa565b1492915050565b60006001600160a01b03821661464d5760405162461bcd60e51b815260206004820152600b60248201526a626164207769746e65737360a81b6044820152606401610b54565b60208401516000906001161561466457601c614667565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe19918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa158015614702573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b614732614e2d565b61475f6001848460405160200161474b93929190615b26565b60405160208183030381529060405261491e565b90505b61476b8161452e565b61291b578051604080516020810192909252614787910161474b565b9050614762565b614796614e2d565b825186516401000003d01991829006919006036147f55760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610b54565b61480087898861496b565b61484c5760405162461bcd60e51b815260206004820152601660248201527f4669727374206d756c20636865636b206661696c6564000000000000000000006044820152606401610b54565b61485784868561496b565b6148a35760405162461bcd60e51b815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c65640000000000000000006044820152606401610b54565b6148ae868484614a96565b98975050505050505050565b6000600286868685876040516020016148d896959493929190615b47565b60408051601f1981840301815291905280516020909101209695505050505050565b6000806401000003d01980848509840990506401000003d019600782089392505050565b614926614e2d565b61492f82614b5d565b815261494461493f8260006145f6565b614b98565b6020820181905260029006600103612188576020810180516401000003d019039052919050565b6000826000036149ab5760405162461bcd60e51b815260206004820152600b60248201526a3d32b9379039b1b0b630b960a91b6044820152606401610b54565b835160208501516000906149c190600290615ba6565b156149cd57601c6149d0565b601b5b9050600070014551231950b75fc4402da1732fc9bebe198387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa158015614a42573d6000803e3d6000fd5b505050602060405103519050600086604051602001614a619190615bba565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b614a9e614e2d565b835160208086015185519186015160009384938493614abf93909190614bb8565b919450925090506401000003d019858209600114614b1f5760405162461bcd60e51b815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a000000000000006044820152606401610b54565b60405180604001604052806401000003d01980614b3e57614b3e615a2d565b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d019811061218857604080516020808201939093528151808203840181529082019091528051910120614b65565b600061291b826002614bb16401000003d0196001615753565b901c614c98565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000614bf883838585614d3d565b9098509050614c0988828e88614d61565b9098509050614c1a88828c87614d61565b90985090506000614c2d8d878b85614d61565b9098509050614c3e88828686614d3d565b9098509050614c4f88828e89614d61565b9098509050818114614c84576401000003d019818a0998506401000003d01982890997506401000003d0198183099650614c88565b8196505b5050505050509450945094915050565b600080614ca3614e4b565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152614cd5614e69565b60208160c0846005600019fa925082600003614d335760405162461bcd60e51b815260206004820152601260248201527f6269674d6f64457870206661696c7572652100000000000000000000000000006044820152606401610b54565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215614dff579160200282015b82811115614dff57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190614dca565b50614e0b929150614e87565b5090565b5080546000825590600052602060002090810190610b5d9190614e87565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b80821115614e0b5760008155600101614e88565b6001600160a01b0381168114610b5d57600080fd5b803561218881614e9c565b600060208284031215614ece57600080fd5b81356111c381614e9c565b806040810183101561291b57600080fd5b600060408284031215614efc57600080fd5b6111c38383614ed9565b60008060408385031215614f1957600080fd5b823591506020830135614f2b81614e9c565b809150509250929050565b634e487b7160e01b600052604160045260246000fd5b60405160c081016001600160401b0381118282101715614f6e57614f6e614f36565b60405290565b60405161012081016001600160401b0381118282101715614f6e57614f6e614f36565b604051601f8201601f191681016001600160401b0381118282101715614fbf57614fbf614f36565b604052919050565b600082601f830112614fd857600080fd5b604051604081018181106001600160401b0382111715614ffa57614ffa614f36565b806040525080604084018581111561501157600080fd5b845b8181101561502b578035835260209283019201615013565b509195945050505050565b80356001600160401b038116811461218857600080fd5b803563ffffffff8116811461218857600080fd5b600060c0828403121561507357600080fd5b61507b614f4c565b905061508682615036565b81526020808301358183015261509e6040840161504d565b60408301526150af6060840161504d565b606083015260808301356150c281614e9c565b608083015260a08301356001600160401b03808211156150e157600080fd5b818501915085601f8301126150f557600080fd5b81358181111561510757615107614f36565b615119601f8201601f19168501614f97565b9150808252868482850101111561512f57600080fd5b80848401858401376000848284010152508060a085015250505092915050565b8015158114610b5d57600080fd5b80356121888161514f565b60008060008385036101e081121561517f57600080fd5b6101a08082121561518f57600080fd5b615197614f74565b91506151a38787614fc7565b82526151b28760408801614fc7565b60208301526080860135604083015260a0860135606083015260c086013560808301526151e160e08701614eb1565b60a08301526101006151f588828901614fc7565b60c0840152615208886101408901614fc7565b60e0840152610180870135908301529093508401356001600160401b0381111561523157600080fd5b61523d86828701615061565b92505061524d6101c0850161515d565b90509250925092565b60006020828403121561526857600080fd5b5035919050565b6000806040838503121561528257600080fd5b823561528d81614e9c565b91506020830135614f2b81614e9c565b600080606083850312156152b057600080fd5b6152ba8484614ed9565b91506152c860408401615036565b90509250929050565b6000602082840312156152e357600080fd5b81356001600160401b038111156152f957600080fd5b820160c081850312156111c357600080fd5b6000806000806060858703121561532157600080fd5b843561532c81614e9c565b93506020850135925060408501356001600160401b038082111561534f57600080fd5b818701915087601f83011261536357600080fd5b81358181111561537257600080fd5b88602082850101111561538457600080fd5b95989497505060200194505050565b803561ffff8116811461218857600080fd5b803560ff8116811461218857600080fd5b60008060008060008060008060006101208a8c0312156153d557600080fd5b6153de8a615393565b98506153ec60208b0161504d565b97506153fa60408b0161504d565b965061540860608b0161504d565b955060808a0135945061541d60a08b0161504d565b935061542b60c08b0161504d565b925061543960e08b016153a5565b91506154486101008b016153a5565b90509295985092959850929598565b6000806040838503121561546a57600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b838110156154a95781518752958201959082019060010161548d565b509495945050505050565b6020815260006111c36020830184615479565b6000604082840312156154d957600080fd5b6111c38383614fc7565b600081518084526020808501945080840160005b838110156154a95781516001600160a01b0316875295820195908201906001016154f7565b60006001600160601b0380881683528087166020840152506001600160401b03851660408301526001600160a01b038416606083015260a0608083015261556660a08301846154e3565b979650505050505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b8181038181111561291b5761291b615587565b634e487b7160e01b600052603160045260246000fd5b6000600182016155d8576155d8615587565b5060010190565b6001600160601b0382811682821603908082111561452757614527615587565b60006001600160401b0380831681810361561b5761561b615587565b6001019392505050565b60006001600160401b0382168061563e5761563e615587565b6000190192915050565b6020815260ff8251166020820152602082015160408201526001600160a01b0360408301511660608201526000606083015160c0608084015261568e60e08401826154e3565b905060808401516001600160601b0380821660a08601528060a08701511660c086015250508091505092915050565b6000815180845260005b818110156156e3576020818501810151868301820152016156c7565b506000602082860101526020601f19601f83011685010191505092915050565b6020815260006111c360208301846156bd565b60006020828403121561572857600080fd5b81516111c38161514f565b6001600160601b0381811683821601908082111561452757614527615587565b8082018082111561291b5761291b615587565b60006020828403121561577857600080fd5b6111c382615393565b60006020828403121561579357600080fd5b6111c38261504d565b6000808335601e198436030181126157b357600080fd5b8301803591506001600160401b038211156157cd57600080fd5b6020019150368190038213156157e257600080fd5b9250929050565b878152866020820152856040820152600063ffffffff80871660608401528086166080840152506001600160a01b03841660a083015260e060c083015261583360e08301846156bd565b9998505050505050505050565b86815285602082015261ffff85166040820152600063ffffffff808616606084015280851660808401525060c060a08301526148ae60c08301846156bd565b6001600160401b0381811683821601908082111561452757614527615587565b8060005b6002811015612e6b5781518452602093840193909101906001016158a3565b6040810161291b828461589f565b6000602082840312156158e257600080fd5b5051919050565b8781526001600160401b0387166020820152856040820152600063ffffffff80871660608401528086166080840152506001600160a01b03841660a083015260e060c083015261583360e08301846156bd565b8281526040602082015260006159556040830184615479565b949350505050565b6001600160e01b031981358181169160048510156159855780818660040360031b1b83161692505b505092915050565b6000808585111561599d57600080fd5b838611156159aa57600080fd5b5050820193919092039150565b6000602082840312156159c957600080fd5b604051602081018181106001600160401b03821117156159eb576159eb614f36565b60405282356159f98161514f565b81529392505050565b828152606081016111c3602083018461589f565b808202811582820484141761291b5761291b615587565b634e487b7160e01b600052601260045260246000fd5b600082615a5257615a52615a2d565b500490565b60ff818116838216019081111561291b5761291b615587565b63ffffffff82811682821603908082111561452757614527615587565b60208152816020820152818360408301376000818301604090810191909152601f909201601f19160101919050565b805169ffffffffffffffffffff8116811461218857600080fd5b600080600080600060a08688031215615aee57600080fd5b615af786615abc565b9450602086015193506040860151925060608601519150615b1a60808701615abc565b90509295509295909350565b838152615b36602082018461589f565b606081019190915260800192915050565b868152615b57602082018761589f565b615b64606082018661589f565b615b7160a082018561589f565b615b7e60e082018461589f565b60609190911b6bffffffffffffffffffffffff19166101208201526101340195945050505050565b600082615bb557615bb5615a2d565b500690565b615bc4818361589f565b60400191905056fea164736f6c6343000813000a",
}

var VRFCoordinatorV25OptimismABI = VRFCoordinatorV25OptimismMetaData.ABI

var VRFCoordinatorV25OptimismBin = VRFCoordinatorV25OptimismMetaData.Bin

func DeployVRFCoordinatorV25Optimism(auth *bind.TransactOpts, backend bind.ContractBackend, blockhashStore common.Address, module common.Address) (common.Address, *types.Transaction, *VRFCoordinatorV25Optimism, error) {
	parsed, err := VRFCoordinatorV25OptimismMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorV25OptimismBin), backend, blockhashStore, module)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorV25Optimism{address: address, abi: *parsed, VRFCoordinatorV25OptimismCaller: VRFCoordinatorV25OptimismCaller{contract: contract}, VRFCoordinatorV25OptimismTransactor: VRFCoordinatorV25OptimismTransactor{contract: contract}, VRFCoordinatorV25OptimismFilterer: VRFCoordinatorV25OptimismFilterer{contract: contract}}, nil
}

type VRFCoordinatorV25Optimism struct {
	address common.Address
	abi     abi.ABI
	VRFCoordinatorV25OptimismCaller
	VRFCoordinatorV25OptimismTransactor
	VRFCoordinatorV25OptimismFilterer
}

type VRFCoordinatorV25OptimismCaller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV25OptimismTransactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV25OptimismFilterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV25OptimismSession struct {
	Contract     *VRFCoordinatorV25Optimism
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV25OptimismCallerSession struct {
	Contract *VRFCoordinatorV25OptimismCaller
	CallOpts bind.CallOpts
}

type VRFCoordinatorV25OptimismTransactorSession struct {
	Contract     *VRFCoordinatorV25OptimismTransactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV25OptimismRaw struct {
	Contract *VRFCoordinatorV25Optimism
}

type VRFCoordinatorV25OptimismCallerRaw struct {
	Contract *VRFCoordinatorV25OptimismCaller
}

type VRFCoordinatorV25OptimismTransactorRaw struct {
	Contract *VRFCoordinatorV25OptimismTransactor
}

func NewVRFCoordinatorV25Optimism(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorV25Optimism, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorV25OptimismABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFCoordinatorV25Optimism(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25Optimism{address: address, abi: abi, VRFCoordinatorV25OptimismCaller: VRFCoordinatorV25OptimismCaller{contract: contract}, VRFCoordinatorV25OptimismTransactor: VRFCoordinatorV25OptimismTransactor{contract: contract}, VRFCoordinatorV25OptimismFilterer: VRFCoordinatorV25OptimismFilterer{contract: contract}}, nil
}

func NewVRFCoordinatorV25OptimismCaller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorV25OptimismCaller, error) {
	contract, err := bindVRFCoordinatorV25Optimism(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismCaller{contract: contract}, nil
}

func NewVRFCoordinatorV25OptimismTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorV25OptimismTransactor, error) {
	contract, err := bindVRFCoordinatorV25Optimism(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismTransactor{contract: contract}, nil
}

func NewVRFCoordinatorV25OptimismFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorV25OptimismFilterer, error) {
	contract, err := bindVRFCoordinatorV25Optimism(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismFilterer{contract: contract}, nil
}

func bindVRFCoordinatorV25Optimism(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFCoordinatorV25OptimismMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV25Optimism.Contract.VRFCoordinatorV25OptimismCaller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.VRFCoordinatorV25OptimismTransactor.contract.Transfer(opts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.VRFCoordinatorV25OptimismTransactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV25Optimism.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.contract.Transfer(opts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "BLOCKHASH_STORE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV25Optimism.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV25Optimism.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) LINK() (common.Address, error) {
	return _VRFCoordinatorV25Optimism.Contract.LINK(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) LINK() (common.Address, error) {
	return _VRFCoordinatorV25Optimism.Contract.LINK(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) LINKNATIVEFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "LINK_NATIVE_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) LINKNATIVEFEED() (common.Address, error) {
	return _VRFCoordinatorV25Optimism.Contract.LINKNATIVEFEED(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) LINKNATIVEFEED() (common.Address, error) {
	return _VRFCoordinatorV25Optimism.Contract.LINKNATIVEFEED(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV25Optimism.Contract.MAXCONSUMERS(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV25Optimism.Contract.MAXCONSUMERS(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) MAXNUMWORDS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "MAX_NUM_WORDS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV25Optimism.Contract.MAXNUMWORDS(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV25Optimism.Contract.MAXNUMWORDS(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "MAX_REQUEST_CONFIRMATIONS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV25Optimism.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV25Optimism.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) GetActiveSubscriptionIds(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "getActiveSubscriptionIds", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) GetActiveSubscriptionIds(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VRFCoordinatorV25Optimism.Contract.GetActiveSubscriptionIds(&_VRFCoordinatorV25Optimism.CallOpts, startIndex, maxCount)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) GetActiveSubscriptionIds(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VRFCoordinatorV25Optimism.Contract.GetActiveSubscriptionIds(&_VRFCoordinatorV25Optimism.CallOpts, startIndex, maxCount)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "getSubscription", subId)

	outstruct := new(GetSubscription)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Balance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.NativeBalance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ReqCount = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.SubOwner = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[4], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinatorV25Optimism.Contract.GetSubscription(&_VRFCoordinatorV25Optimism.CallOpts, subId)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinatorV25Optimism.Contract.GetSubscription(&_VRFCoordinatorV25Optimism.CallOpts, subId)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "hashOfKey", publicKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25Optimism.Contract.HashOfKey(&_VRFCoordinatorV25Optimism.CallOpts, publicKey)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25Optimism.Contract.HashOfKey(&_VRFCoordinatorV25Optimism.CallOpts, publicKey)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) Owner() (common.Address, error) {
	return _VRFCoordinatorV25Optimism.Contract.Owner(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) Owner() (common.Address, error) {
	return _VRFCoordinatorV25Optimism.Contract.Owner(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) PendingRequestExists(opts *bind.CallOpts, subId *big.Int) (bool, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "pendingRequestExists", subId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) PendingRequestExists(subId *big.Int) (bool, error) {
	return _VRFCoordinatorV25Optimism.Contract.PendingRequestExists(&_VRFCoordinatorV25Optimism.CallOpts, subId)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) PendingRequestExists(subId *big.Int) (bool, error) {
	return _VRFCoordinatorV25Optimism.Contract.PendingRequestExists(&_VRFCoordinatorV25Optimism.CallOpts, subId)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) SConfig(opts *bind.CallOpts) (SConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "s_config")

	outstruct := new(SConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MinimumRequestConfirmations = *abi.ConvertType(out[0], new(uint16)).(*uint16)
	outstruct.MaxGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.ReentrancyLock = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.StalenessSeconds = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.GasAfterPaymentCalculation = *abi.ConvertType(out[4], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeNativePPM = *abi.ConvertType(out[5], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeLinkDiscountPPM = *abi.ConvertType(out[6], new(uint32)).(*uint32)
	outstruct.NativePremiumPercentage = *abi.ConvertType(out[7], new(uint8)).(*uint8)
	outstruct.LinkPremiumPercentage = *abi.ConvertType(out[8], new(uint8)).(*uint8)

	return *outstruct, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) SConfig() (SConfig,

	error) {
	return _VRFCoordinatorV25Optimism.Contract.SConfig(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) SConfig() (SConfig,

	error) {
	return _VRFCoordinatorV25Optimism.Contract.SConfig(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) SCurrentSubNonce(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "s_currentSubNonce")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) SCurrentSubNonce() (uint64, error) {
	return _VRFCoordinatorV25Optimism.Contract.SCurrentSubNonce(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) SCurrentSubNonce() (uint64, error) {
	return _VRFCoordinatorV25Optimism.Contract.SCurrentSubNonce(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) SFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "s_fallbackWeiPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV25Optimism.Contract.SFallbackWeiPerUnitLink(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV25Optimism.Contract.SFallbackWeiPerUnitLink(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) SOptimismGasModule(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "s_optimismGasModule")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) SOptimismGasModule() (common.Address, error) {
	return _VRFCoordinatorV25Optimism.Contract.SOptimismGasModule(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) SOptimismGasModule() (common.Address, error) {
	return _VRFCoordinatorV25Optimism.Contract.SOptimismGasModule(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) SProvingKeyHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "s_provingKeyHashes", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) SProvingKeyHashes(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25Optimism.Contract.SProvingKeyHashes(&_VRFCoordinatorV25Optimism.CallOpts, arg0)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) SProvingKeyHashes(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25Optimism.Contract.SProvingKeyHashes(&_VRFCoordinatorV25Optimism.CallOpts, arg0)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) SProvingKeys(opts *bind.CallOpts, arg0 [32]byte) (SProvingKeys,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "s_provingKeys", arg0)

	outstruct := new(SProvingKeys)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.MaxGas = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) SProvingKeys(arg0 [32]byte) (SProvingKeys,

	error) {
	return _VRFCoordinatorV25Optimism.Contract.SProvingKeys(&_VRFCoordinatorV25Optimism.CallOpts, arg0)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) SProvingKeys(arg0 [32]byte) (SProvingKeys,

	error) {
	return _VRFCoordinatorV25Optimism.Contract.SProvingKeys(&_VRFCoordinatorV25Optimism.CallOpts, arg0)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) SRequestCommitments(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "s_requestCommitments", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) SRequestCommitments(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25Optimism.Contract.SRequestCommitments(&_VRFCoordinatorV25Optimism.CallOpts, arg0)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) SRequestCommitments(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25Optimism.Contract.SRequestCommitments(&_VRFCoordinatorV25Optimism.CallOpts, arg0)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) STotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "s_totalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) STotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV25Optimism.Contract.STotalBalance(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) STotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV25Optimism.Contract.STotalBalance(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) STotalNativeBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "s_totalNativeBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) STotalNativeBalance() (*big.Int, error) {
	return _VRFCoordinatorV25Optimism.Contract.STotalNativeBalance(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) STotalNativeBalance() (*big.Int, error) {
	return _VRFCoordinatorV25Optimism.Contract.STotalNativeBalance(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "acceptOwnership")
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.AcceptOwnership(&_VRFCoordinatorV25Optimism.TransactOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.AcceptOwnership(&_VRFCoordinatorV25Optimism.TransactOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV25Optimism.TransactOpts, subId)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV25Optimism.TransactOpts, subId)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "addConsumer", subId, consumer)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.AddConsumer(&_VRFCoordinatorV25Optimism.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.AddConsumer(&_VRFCoordinatorV25Optimism.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "cancelSubscription", subId, to)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.CancelSubscription(&_VRFCoordinatorV25Optimism.TransactOpts, subId, to)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.CancelSubscription(&_VRFCoordinatorV25Optimism.TransactOpts, subId, to)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "createSubscription")
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.CreateSubscription(&_VRFCoordinatorV25Optimism.TransactOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.CreateSubscription(&_VRFCoordinatorV25Optimism.TransactOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) DeregisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "deregisterMigratableCoordinator", target)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) DeregisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.DeregisterMigratableCoordinator(&_VRFCoordinatorV25Optimism.TransactOpts, target)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) DeregisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.DeregisterMigratableCoordinator(&_VRFCoordinatorV25Optimism.TransactOpts, target)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "deregisterProvingKey", publicProvingKey)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.DeregisterProvingKey(&_VRFCoordinatorV25Optimism.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.DeregisterProvingKey(&_VRFCoordinatorV25Optimism.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "fulfillRandomWords", proof, rc, onlyPremium)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) FulfillRandomWords(proof VRFProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.FulfillRandomWords(&_VRFCoordinatorV25Optimism.TransactOpts, proof, rc, onlyPremium)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) FulfillRandomWords(proof VRFProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.FulfillRandomWords(&_VRFCoordinatorV25Optimism.TransactOpts, proof, rc, onlyPremium)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) FundSubscriptionWithNative(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "fundSubscriptionWithNative", subId)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) FundSubscriptionWithNative(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.FundSubscriptionWithNative(&_VRFCoordinatorV25Optimism.TransactOpts, subId)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) FundSubscriptionWithNative(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.FundSubscriptionWithNative(&_VRFCoordinatorV25Optimism.TransactOpts, subId)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) Migrate(opts *bind.TransactOpts, subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "migrate", subId, newCoordinator)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) Migrate(subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.Migrate(&_VRFCoordinatorV25Optimism.TransactOpts, subId, newCoordinator)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) Migrate(subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.Migrate(&_VRFCoordinatorV25Optimism.TransactOpts, subId, newCoordinator)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.OnTokenTransfer(&_VRFCoordinatorV25Optimism.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.OnTokenTransfer(&_VRFCoordinatorV25Optimism.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) OwnerCancelSubscription(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "ownerCancelSubscription", subId)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) OwnerCancelSubscription(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.OwnerCancelSubscription(&_VRFCoordinatorV25Optimism.TransactOpts, subId)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) OwnerCancelSubscription(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.OwnerCancelSubscription(&_VRFCoordinatorV25Optimism.TransactOpts, subId)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "recoverFunds", to)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RecoverFunds(&_VRFCoordinatorV25Optimism.TransactOpts, to)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RecoverFunds(&_VRFCoordinatorV25Optimism.TransactOpts, to)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) RecoverNativeFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "recoverNativeFunds", to)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) RecoverNativeFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RecoverNativeFunds(&_VRFCoordinatorV25Optimism.TransactOpts, to)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) RecoverNativeFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RecoverNativeFunds(&_VRFCoordinatorV25Optimism.TransactOpts, to)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) RegisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "registerMigratableCoordinator", target)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) RegisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RegisterMigratableCoordinator(&_VRFCoordinatorV25Optimism.TransactOpts, target)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) RegisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RegisterMigratableCoordinator(&_VRFCoordinatorV25Optimism.TransactOpts, target)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) RegisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int, maxGas uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "registerProvingKey", publicProvingKey, maxGas)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) RegisterProvingKey(publicProvingKey [2]*big.Int, maxGas uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RegisterProvingKey(&_VRFCoordinatorV25Optimism.TransactOpts, publicProvingKey, maxGas)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) RegisterProvingKey(publicProvingKey [2]*big.Int, maxGas uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RegisterProvingKey(&_VRFCoordinatorV25Optimism.TransactOpts, publicProvingKey, maxGas)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "removeConsumer", subId, consumer)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RemoveConsumer(&_VRFCoordinatorV25Optimism.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RemoveConsumer(&_VRFCoordinatorV25Optimism.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) RequestRandomWords(opts *bind.TransactOpts, req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "requestRandomWords", req)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RequestRandomWords(&_VRFCoordinatorV25Optimism.TransactOpts, req)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RequestRandomWords(&_VRFCoordinatorV25Optimism.TransactOpts, req)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV25Optimism.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV25Optimism.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeNativePPM uint32, fulfillmentFlatFeeLinkDiscountPPM uint32, nativePremiumPercentage uint8, linkPremiumPercentage uint8) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fulfillmentFlatFeeNativePPM, fulfillmentFlatFeeLinkDiscountPPM, nativePremiumPercentage, linkPremiumPercentage)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeNativePPM uint32, fulfillmentFlatFeeLinkDiscountPPM uint32, nativePremiumPercentage uint8, linkPremiumPercentage uint8) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.SetConfig(&_VRFCoordinatorV25Optimism.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fulfillmentFlatFeeNativePPM, fulfillmentFlatFeeLinkDiscountPPM, nativePremiumPercentage, linkPremiumPercentage)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeNativePPM uint32, fulfillmentFlatFeeLinkDiscountPPM uint32, nativePremiumPercentage uint8, linkPremiumPercentage uint8) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.SetConfig(&_VRFCoordinatorV25Optimism.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fulfillmentFlatFeeNativePPM, fulfillmentFlatFeeLinkDiscountPPM, nativePremiumPercentage, linkPremiumPercentage)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) SetLINKAndLINKNativeFeed(opts *bind.TransactOpts, link common.Address, linkNativeFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "setLINKAndLINKNativeFeed", link, linkNativeFeed)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) SetLINKAndLINKNativeFeed(link common.Address, linkNativeFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.SetLINKAndLINKNativeFeed(&_VRFCoordinatorV25Optimism.TransactOpts, link, linkNativeFeed)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) SetLINKAndLINKNativeFeed(link common.Address, linkNativeFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.SetLINKAndLINKNativeFeed(&_VRFCoordinatorV25Optimism.TransactOpts, link, linkNativeFeed)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) SetOptimismGasModule(opts *bind.TransactOpts, module common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "setOptimismGasModule", module)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) SetOptimismGasModule(module common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.SetOptimismGasModule(&_VRFCoordinatorV25Optimism.TransactOpts, module)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) SetOptimismGasModule(module common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.SetOptimismGasModule(&_VRFCoordinatorV25Optimism.TransactOpts, module)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.TransferOwnership(&_VRFCoordinatorV25Optimism.TransactOpts, to)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.TransferOwnership(&_VRFCoordinatorV25Optimism.TransactOpts, to)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) Withdraw(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "withdraw", recipient)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) Withdraw(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.Withdraw(&_VRFCoordinatorV25Optimism.TransactOpts, recipient)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) Withdraw(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.Withdraw(&_VRFCoordinatorV25Optimism.TransactOpts, recipient)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) WithdrawNative(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "withdrawNative", recipient)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) WithdrawNative(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.WithdrawNative(&_VRFCoordinatorV25Optimism.TransactOpts, recipient)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) WithdrawNative(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.WithdrawNative(&_VRFCoordinatorV25Optimism.TransactOpts, recipient)
}

type VRFCoordinatorV25OptimismConfigSetIterator struct {
	Event *VRFCoordinatorV25OptimismConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismConfigSet)
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
		it.Event = new(VRFCoordinatorV25OptimismConfigSet)
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

func (it *VRFCoordinatorV25OptimismConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismConfigSet struct {
	MinimumRequestConfirmations       uint16
	MaxGasLimit                       uint32
	StalenessSeconds                  uint32
	GasAfterPaymentCalculation        uint32
	FallbackWeiPerUnitLink            *big.Int
	FulfillmentFlatFeeNativePPM       uint32
	FulfillmentFlatFeeLinkDiscountPPM uint32
	NativePremiumPercentage           uint8
	LinkPremiumPercentage             uint8
	Raw                               types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismConfigSetIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismConfigSet)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseConfigSet(log types.Log) (*VRFCoordinatorV25OptimismConfigSet, error) {
	event := new(VRFCoordinatorV25OptimismConfigSet)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismCoordinatorDeregisteredIterator struct {
	Event *VRFCoordinatorV25OptimismCoordinatorDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismCoordinatorDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismCoordinatorDeregistered)
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
		it.Event = new(VRFCoordinatorV25OptimismCoordinatorDeregistered)
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

func (it *VRFCoordinatorV25OptimismCoordinatorDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismCoordinatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismCoordinatorDeregistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismCoordinatorDeregisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismCoordinatorDeregisteredIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "CoordinatorDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismCoordinatorDeregistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismCoordinatorDeregistered)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseCoordinatorDeregistered(log types.Log) (*VRFCoordinatorV25OptimismCoordinatorDeregistered, error) {
	event := new(VRFCoordinatorV25OptimismCoordinatorDeregistered)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismCoordinatorRegisteredIterator struct {
	Event *VRFCoordinatorV25OptimismCoordinatorRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismCoordinatorRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismCoordinatorRegistered)
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
		it.Event = new(VRFCoordinatorV25OptimismCoordinatorRegistered)
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

func (it *VRFCoordinatorV25OptimismCoordinatorRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismCoordinatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismCoordinatorRegistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismCoordinatorRegisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismCoordinatorRegisteredIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "CoordinatorRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismCoordinatorRegistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismCoordinatorRegistered)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseCoordinatorRegistered(log types.Log) (*VRFCoordinatorV25OptimismCoordinatorRegistered, error) {
	event := new(VRFCoordinatorV25OptimismCoordinatorRegistered)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsedIterator struct {
	Event *VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsed)
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
		it.Event = new(VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsed)
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

func (it *VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsed struct {
	RequestId              *big.Int
	FallbackWeiPerUnitLink *big.Int
	Raw                    types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterFallbackWeiPerUnitLinkUsed(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsedIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "FallbackWeiPerUnitLinkUsed")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsedIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "FallbackWeiPerUnitLinkUsed", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchFallbackWeiPerUnitLinkUsed(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsed) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "FallbackWeiPerUnitLinkUsed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsed)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "FallbackWeiPerUnitLinkUsed", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseFallbackWeiPerUnitLinkUsed(log types.Log) (*VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsed, error) {
	event := new(VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsed)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "FallbackWeiPerUnitLinkUsed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismFundsRecoveredIterator struct {
	Event *VRFCoordinatorV25OptimismFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismFundsRecovered)
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
		it.Event = new(VRFCoordinatorV25OptimismFundsRecovered)
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

func (it *VRFCoordinatorV25OptimismFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismFundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismFundsRecoveredIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismFundsRecovered)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseFundsRecovered(log types.Log) (*VRFCoordinatorV25OptimismFundsRecovered, error) {
	event := new(VRFCoordinatorV25OptimismFundsRecovered)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismMigrationCompletedIterator struct {
	Event *VRFCoordinatorV25OptimismMigrationCompleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismMigrationCompletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismMigrationCompleted)
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
		it.Event = new(VRFCoordinatorV25OptimismMigrationCompleted)
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

func (it *VRFCoordinatorV25OptimismMigrationCompletedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismMigrationCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismMigrationCompleted struct {
	NewCoordinator common.Address
	SubId          *big.Int
	Raw            types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterMigrationCompleted(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismMigrationCompletedIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "MigrationCompleted")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismMigrationCompletedIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "MigrationCompleted", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchMigrationCompleted(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismMigrationCompleted) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "MigrationCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismMigrationCompleted)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseMigrationCompleted(log types.Log) (*VRFCoordinatorV25OptimismMigrationCompleted, error) {
	event := new(VRFCoordinatorV25OptimismMigrationCompleted)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismNativeFundsRecoveredIterator struct {
	Event *VRFCoordinatorV25OptimismNativeFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismNativeFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismNativeFundsRecovered)
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
		it.Event = new(VRFCoordinatorV25OptimismNativeFundsRecovered)
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

func (it *VRFCoordinatorV25OptimismNativeFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismNativeFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismNativeFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterNativeFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismNativeFundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "NativeFundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismNativeFundsRecoveredIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "NativeFundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchNativeFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismNativeFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "NativeFundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismNativeFundsRecovered)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "NativeFundsRecovered", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseNativeFundsRecovered(log types.Log) (*VRFCoordinatorV25OptimismNativeFundsRecovered, error) {
	event := new(VRFCoordinatorV25OptimismNativeFundsRecovered)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "NativeFundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismOwnershipTransferRequestedIterator struct {
	Event *VRFCoordinatorV25OptimismOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismOwnershipTransferRequested)
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
		it.Event = new(VRFCoordinatorV25OptimismOwnershipTransferRequested)
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

func (it *VRFCoordinatorV25OptimismOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV25OptimismOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismOwnershipTransferRequestedIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismOwnershipTransferRequested)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV25OptimismOwnershipTransferRequested, error) {
	event := new(VRFCoordinatorV25OptimismOwnershipTransferRequested)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismOwnershipTransferredIterator struct {
	Event *VRFCoordinatorV25OptimismOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismOwnershipTransferred)
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
		it.Event = new(VRFCoordinatorV25OptimismOwnershipTransferred)
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

func (it *VRFCoordinatorV25OptimismOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV25OptimismOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismOwnershipTransferredIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismOwnershipTransferred)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV25OptimismOwnershipTransferred, error) {
	event := new(VRFCoordinatorV25OptimismOwnershipTransferred)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismProvingKeyDeregisteredIterator struct {
	Event *VRFCoordinatorV25OptimismProvingKeyDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismProvingKeyDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismProvingKeyDeregistered)
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
		it.Event = new(VRFCoordinatorV25OptimismProvingKeyDeregistered)
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

func (it *VRFCoordinatorV25OptimismProvingKeyDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismProvingKeyDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismProvingKeyDeregistered struct {
	KeyHash [32]byte
	MaxGas  uint64
	Raw     types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterProvingKeyDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismProvingKeyDeregisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "ProvingKeyDeregistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismProvingKeyDeregisteredIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "ProvingKeyDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismProvingKeyDeregistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "ProvingKeyDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismProvingKeyDeregistered)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorV25OptimismProvingKeyDeregistered, error) {
	event := new(VRFCoordinatorV25OptimismProvingKeyDeregistered)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismProvingKeyRegisteredIterator struct {
	Event *VRFCoordinatorV25OptimismProvingKeyRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismProvingKeyRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismProvingKeyRegistered)
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
		it.Event = new(VRFCoordinatorV25OptimismProvingKeyRegistered)
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

func (it *VRFCoordinatorV25OptimismProvingKeyRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismProvingKeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismProvingKeyRegistered struct {
	KeyHash [32]byte
	MaxGas  uint64
	Raw     types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterProvingKeyRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismProvingKeyRegisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "ProvingKeyRegistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismProvingKeyRegisteredIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "ProvingKeyRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismProvingKeyRegistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "ProvingKeyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismProvingKeyRegistered)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV25OptimismProvingKeyRegistered, error) {
	event := new(VRFCoordinatorV25OptimismProvingKeyRegistered)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismRandomWordsFulfilledIterator struct {
	Event *VRFCoordinatorV25OptimismRandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismRandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismRandomWordsFulfilled)
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
		it.Event = new(VRFCoordinatorV25OptimismRandomWordsFulfilled)
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

func (it *VRFCoordinatorV25OptimismRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismRandomWordsFulfilled struct {
	RequestId     *big.Int
	OutputSeed    *big.Int
	SubId         *big.Int
	Payment       *big.Int
	NativePayment bool
	Success       bool
	OnlyPremium   bool
	Raw           types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subId []*big.Int) (*VRFCoordinatorV25OptimismRandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule, subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismRandomWordsFulfilledIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismRandomWordsFulfilled, requestId []*big.Int, subId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule, subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismRandomWordsFulfilled)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV25OptimismRandomWordsFulfilled, error) {
	event := new(VRFCoordinatorV25OptimismRandomWordsFulfilled)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismRandomWordsRequestedIterator struct {
	Event *VRFCoordinatorV25OptimismRandomWordsRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismRandomWordsRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismRandomWordsRequested)
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
		it.Event = new(VRFCoordinatorV25OptimismRandomWordsRequested)
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

func (it *VRFCoordinatorV25OptimismRandomWordsRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismRandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismRandomWordsRequested struct {
	KeyHash                     [32]byte
	RequestId                   *big.Int
	PreSeed                     *big.Int
	SubId                       *big.Int
	MinimumRequestConfirmations uint16
	CallbackGasLimit            uint32
	NumWords                    uint32
	ExtraArgs                   []byte
	Sender                      common.Address
	Raw                         types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (*VRFCoordinatorV25OptimismRandomWordsRequestedIterator, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismRandomWordsRequestedIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismRandomWordsRequested, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (event.Subscription, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismRandomWordsRequested)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV25OptimismRandomWordsRequested, error) {
	event := new(VRFCoordinatorV25OptimismRandomWordsRequested)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismSubscriptionCanceledIterator struct {
	Event *VRFCoordinatorV25OptimismSubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismSubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismSubscriptionCanceled)
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
		it.Event = new(VRFCoordinatorV25OptimismSubscriptionCanceled)
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

func (it *VRFCoordinatorV25OptimismSubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismSubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismSubscriptionCanceled struct {
	SubId        *big.Int
	To           common.Address
	AmountLink   *big.Int
	AmountNative *big.Int
	Raw          types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismSubscriptionCanceledIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionCanceled, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismSubscriptionCanceled)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionCanceled, error) {
	event := new(VRFCoordinatorV25OptimismSubscriptionCanceled)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismSubscriptionConsumerAddedIterator struct {
	Event *VRFCoordinatorV25OptimismSubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismSubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismSubscriptionConsumerAdded)
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
		it.Event = new(VRFCoordinatorV25OptimismSubscriptionConsumerAdded)
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

func (it *VRFCoordinatorV25OptimismSubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismSubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismSubscriptionConsumerAdded struct {
	SubId    *big.Int
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismSubscriptionConsumerAddedIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionConsumerAdded, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismSubscriptionConsumerAdded)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionConsumerAdded, error) {
	event := new(VRFCoordinatorV25OptimismSubscriptionConsumerAdded)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismSubscriptionConsumerRemovedIterator struct {
	Event *VRFCoordinatorV25OptimismSubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismSubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismSubscriptionConsumerRemoved)
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
		it.Event = new(VRFCoordinatorV25OptimismSubscriptionConsumerRemoved)
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

func (it *VRFCoordinatorV25OptimismSubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismSubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismSubscriptionConsumerRemoved struct {
	SubId    *big.Int
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismSubscriptionConsumerRemovedIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionConsumerRemoved, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismSubscriptionConsumerRemoved)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionConsumerRemoved, error) {
	event := new(VRFCoordinatorV25OptimismSubscriptionConsumerRemoved)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismSubscriptionCreatedIterator struct {
	Event *VRFCoordinatorV25OptimismSubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismSubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismSubscriptionCreated)
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
		it.Event = new(VRFCoordinatorV25OptimismSubscriptionCreated)
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

func (it *VRFCoordinatorV25OptimismSubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismSubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismSubscriptionCreated struct {
	SubId *big.Int
	Owner common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismSubscriptionCreatedIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionCreated, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismSubscriptionCreated)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionCreated, error) {
	event := new(VRFCoordinatorV25OptimismSubscriptionCreated)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismSubscriptionFundedIterator struct {
	Event *VRFCoordinatorV25OptimismSubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismSubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismSubscriptionFunded)
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
		it.Event = new(VRFCoordinatorV25OptimismSubscriptionFunded)
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

func (it *VRFCoordinatorV25OptimismSubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismSubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismSubscriptionFunded struct {
	SubId      *big.Int
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismSubscriptionFundedIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionFunded, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismSubscriptionFunded)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionFunded, error) {
	event := new(VRFCoordinatorV25OptimismSubscriptionFunded)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismSubscriptionFundedWithNativeIterator struct {
	Event *VRFCoordinatorV25OptimismSubscriptionFundedWithNative

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismSubscriptionFundedWithNativeIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismSubscriptionFundedWithNative)
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
		it.Event = new(VRFCoordinatorV25OptimismSubscriptionFundedWithNative)
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

func (it *VRFCoordinatorV25OptimismSubscriptionFundedWithNativeIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismSubscriptionFundedWithNativeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismSubscriptionFundedWithNative struct {
	SubId            *big.Int
	OldNativeBalance *big.Int
	NewNativeBalance *big.Int
	Raw              types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterSubscriptionFundedWithNative(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionFundedWithNativeIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "SubscriptionFundedWithNative", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismSubscriptionFundedWithNativeIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "SubscriptionFundedWithNative", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchSubscriptionFundedWithNative(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionFundedWithNative, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "SubscriptionFundedWithNative", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismSubscriptionFundedWithNative)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionFundedWithNative", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseSubscriptionFundedWithNative(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionFundedWithNative, error) {
	event := new(VRFCoordinatorV25OptimismSubscriptionFundedWithNative)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionFundedWithNative", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequestedIterator struct {
	Event *VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequested)
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
		it.Event = new(VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequested)
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

func (it *VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequested struct {
	SubId *big.Int
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequestedIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequested, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequested)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequested, error) {
	event := new(VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequested)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismSubscriptionOwnerTransferredIterator struct {
	Event *VRFCoordinatorV25OptimismSubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismSubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismSubscriptionOwnerTransferred)
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
		it.Event = new(VRFCoordinatorV25OptimismSubscriptionOwnerTransferred)
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

func (it *VRFCoordinatorV25OptimismSubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismSubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismSubscriptionOwnerTransferred struct {
	SubId *big.Int
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismSubscriptionOwnerTransferredIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionOwnerTransferred, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismSubscriptionOwnerTransferred)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionOwnerTransferred, error) {
	event := new(VRFCoordinatorV25OptimismSubscriptionOwnerTransferred)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetSubscription struct {
	Balance       *big.Int
	NativeBalance *big.Int
	ReqCount      uint64
	SubOwner      common.Address
	Consumers     []common.Address
}
type SConfig struct {
	MinimumRequestConfirmations       uint16
	MaxGasLimit                       uint32
	ReentrancyLock                    bool
	StalenessSeconds                  uint32
	GasAfterPaymentCalculation        uint32
	FulfillmentFlatFeeNativePPM       uint32
	FulfillmentFlatFeeLinkDiscountPPM uint32
	NativePremiumPercentage           uint8
	LinkPremiumPercentage             uint8
}
type SProvingKeys struct {
	Exists bool
	MaxGas uint64
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25Optimism) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinatorV25Optimism.abi.Events["ConfigSet"].ID:
		return _VRFCoordinatorV25Optimism.ParseConfigSet(log)
	case _VRFCoordinatorV25Optimism.abi.Events["CoordinatorDeregistered"].ID:
		return _VRFCoordinatorV25Optimism.ParseCoordinatorDeregistered(log)
	case _VRFCoordinatorV25Optimism.abi.Events["CoordinatorRegistered"].ID:
		return _VRFCoordinatorV25Optimism.ParseCoordinatorRegistered(log)
	case _VRFCoordinatorV25Optimism.abi.Events["FallbackWeiPerUnitLinkUsed"].ID:
		return _VRFCoordinatorV25Optimism.ParseFallbackWeiPerUnitLinkUsed(log)
	case _VRFCoordinatorV25Optimism.abi.Events["FundsRecovered"].ID:
		return _VRFCoordinatorV25Optimism.ParseFundsRecovered(log)
	case _VRFCoordinatorV25Optimism.abi.Events["MigrationCompleted"].ID:
		return _VRFCoordinatorV25Optimism.ParseMigrationCompleted(log)
	case _VRFCoordinatorV25Optimism.abi.Events["NativeFundsRecovered"].ID:
		return _VRFCoordinatorV25Optimism.ParseNativeFundsRecovered(log)
	case _VRFCoordinatorV25Optimism.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFCoordinatorV25Optimism.ParseOwnershipTransferRequested(log)
	case _VRFCoordinatorV25Optimism.abi.Events["OwnershipTransferred"].ID:
		return _VRFCoordinatorV25Optimism.ParseOwnershipTransferred(log)
	case _VRFCoordinatorV25Optimism.abi.Events["ProvingKeyDeregistered"].ID:
		return _VRFCoordinatorV25Optimism.ParseProvingKeyDeregistered(log)
	case _VRFCoordinatorV25Optimism.abi.Events["ProvingKeyRegistered"].ID:
		return _VRFCoordinatorV25Optimism.ParseProvingKeyRegistered(log)
	case _VRFCoordinatorV25Optimism.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFCoordinatorV25Optimism.ParseRandomWordsFulfilled(log)
	case _VRFCoordinatorV25Optimism.abi.Events["RandomWordsRequested"].ID:
		return _VRFCoordinatorV25Optimism.ParseRandomWordsRequested(log)
	case _VRFCoordinatorV25Optimism.abi.Events["SubscriptionCanceled"].ID:
		return _VRFCoordinatorV25Optimism.ParseSubscriptionCanceled(log)
	case _VRFCoordinatorV25Optimism.abi.Events["SubscriptionConsumerAdded"].ID:
		return _VRFCoordinatorV25Optimism.ParseSubscriptionConsumerAdded(log)
	case _VRFCoordinatorV25Optimism.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _VRFCoordinatorV25Optimism.ParseSubscriptionConsumerRemoved(log)
	case _VRFCoordinatorV25Optimism.abi.Events["SubscriptionCreated"].ID:
		return _VRFCoordinatorV25Optimism.ParseSubscriptionCreated(log)
	case _VRFCoordinatorV25Optimism.abi.Events["SubscriptionFunded"].ID:
		return _VRFCoordinatorV25Optimism.ParseSubscriptionFunded(log)
	case _VRFCoordinatorV25Optimism.abi.Events["SubscriptionFundedWithNative"].ID:
		return _VRFCoordinatorV25Optimism.ParseSubscriptionFundedWithNative(log)
	case _VRFCoordinatorV25Optimism.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _VRFCoordinatorV25Optimism.ParseSubscriptionOwnerTransferRequested(log)
	case _VRFCoordinatorV25Optimism.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _VRFCoordinatorV25Optimism.ParseSubscriptionOwnerTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFCoordinatorV25OptimismConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2c6b6b12413678366b05b145c5f00745bdd00e739131ab5de82484a50c9d78b6")
}

func (VRFCoordinatorV25OptimismCoordinatorDeregistered) Topic() common.Hash {
	return common.HexToHash("0xf80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af37")
}

func (VRFCoordinatorV25OptimismCoordinatorRegistered) Topic() common.Hash {
	return common.HexToHash("0xb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af01625")
}

func (VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsed) Topic() common.Hash {
	return common.HexToHash("0x6ca648a381f22ead7e37773d934e64885dcf861fbfbb26c40354cbf0c4662d1a")
}

func (VRFCoordinatorV25OptimismFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (VRFCoordinatorV25OptimismMigrationCompleted) Topic() common.Hash {
	return common.HexToHash("0xd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be4187")
}

func (VRFCoordinatorV25OptimismNativeFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x4aed7c8eed0496c8c19ea2681fcca25741c1602342e38b045d9f1e8e905d2e9c")
}

func (VRFCoordinatorV25OptimismOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFCoordinatorV25OptimismOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFCoordinatorV25OptimismProvingKeyDeregistered) Topic() common.Hash {
	return common.HexToHash("0x9b6868e0eb737bcd72205360baa6bfd0ba4e4819a33ade2db384e8a8025639a5")
}

func (VRFCoordinatorV25OptimismProvingKeyRegistered) Topic() common.Hash {
	return common.HexToHash("0x9b911b2c240bfbef3b6a8f7ed6ee321d1258bb2a3fe6becab52ac1cd3210afd3")
}

func (VRFCoordinatorV25OptimismRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0xaeb4b4786571e184246d39587f659abf0e26f41f6a3358692250382c0cdb47b7")
}

func (VRFCoordinatorV25OptimismRandomWordsRequested) Topic() common.Hash {
	return common.HexToHash("0xeb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e")
}

func (VRFCoordinatorV25OptimismSubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0x8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c4")
}

func (VRFCoordinatorV25OptimismSubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e1")
}

func (VRFCoordinatorV25OptimismSubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a7")
}

func (VRFCoordinatorV25OptimismSubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d")
}

func (VRFCoordinatorV25OptimismSubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0x1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a")
}

func (VRFCoordinatorV25OptimismSubscriptionFundedWithNative) Topic() common.Hash {
	return common.HexToHash("0x7603b205d03651ee812f803fccde89f1012e545a9c99f0abfea9cedd0fd8e902")
}

func (VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a1")
}

func (VRFCoordinatorV25OptimismSubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0xd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c9386")
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25Optimism) Address() common.Address {
	return _VRFCoordinatorV25Optimism.address
}

type VRFCoordinatorV25OptimismInterface interface {
	BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	LINKNATIVEFEED(opts *bind.CallOpts) (common.Address, error)

	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	MAXNUMWORDS(opts *bind.CallOpts) (uint32, error)

	MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error)

	GetActiveSubscriptionIds(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

		error)

	HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	PendingRequestExists(opts *bind.CallOpts, subId *big.Int) (bool, error)

	SConfig(opts *bind.CallOpts) (SConfig,

		error)

	SCurrentSubNonce(opts *bind.CallOpts) (uint64, error)

	SFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error)

	SOptimismGasModule(opts *bind.CallOpts) (common.Address, error)

	SProvingKeyHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	SProvingKeys(opts *bind.CallOpts, arg0 [32]byte) (SProvingKeys,

		error)

	SRequestCommitments(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	STotalBalance(opts *bind.CallOpts) (*big.Int, error)

	STotalNativeBalance(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	DeregisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error)

	DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error)

	FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error)

	FundSubscriptionWithNative(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	Migrate(opts *bind.TransactOpts, subId *big.Int, newCoordinator common.Address) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RecoverNativeFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RegisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error)

	RegisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int, maxGas uint64) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeNativePPM uint32, fulfillmentFlatFeeLinkDiscountPPM uint32, nativePremiumPercentage uint8, linkPremiumPercentage uint8) (*types.Transaction, error)

	SetLINKAndLINKNativeFeed(opts *bind.TransactOpts, link common.Address, linkNativeFeed common.Address) (*types.Transaction, error)

	SetOptimismGasModule(opts *bind.TransactOpts, module common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error)

	WithdrawNative(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFCoordinatorV25OptimismConfigSet, error)

	FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismCoordinatorDeregisteredIterator, error)

	WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismCoordinatorDeregistered) (event.Subscription, error)

	ParseCoordinatorDeregistered(log types.Log) (*VRFCoordinatorV25OptimismCoordinatorDeregistered, error)

	FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismCoordinatorRegisteredIterator, error)

	WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismCoordinatorRegistered) (event.Subscription, error)

	ParseCoordinatorRegistered(log types.Log) (*VRFCoordinatorV25OptimismCoordinatorRegistered, error)

	FilterFallbackWeiPerUnitLinkUsed(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsedIterator, error)

	WatchFallbackWeiPerUnitLinkUsed(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsed) (event.Subscription, error)

	ParseFallbackWeiPerUnitLinkUsed(log types.Log) (*VRFCoordinatorV25OptimismFallbackWeiPerUnitLinkUsed, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismFundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismFundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*VRFCoordinatorV25OptimismFundsRecovered, error)

	FilterMigrationCompleted(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismMigrationCompletedIterator, error)

	WatchMigrationCompleted(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismMigrationCompleted) (event.Subscription, error)

	ParseMigrationCompleted(log types.Log) (*VRFCoordinatorV25OptimismMigrationCompleted, error)

	FilterNativeFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismNativeFundsRecoveredIterator, error)

	WatchNativeFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismNativeFundsRecovered) (event.Subscription, error)

	ParseNativeFundsRecovered(log types.Log) (*VRFCoordinatorV25OptimismNativeFundsRecovered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV25OptimismOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV25OptimismOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV25OptimismOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV25OptimismOwnershipTransferred, error)

	FilterProvingKeyDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismProvingKeyDeregisteredIterator, error)

	WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismProvingKeyDeregistered) (event.Subscription, error)

	ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorV25OptimismProvingKeyDeregistered, error)

	FilterProvingKeyRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismProvingKeyRegisteredIterator, error)

	WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismProvingKeyRegistered) (event.Subscription, error)

	ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV25OptimismProvingKeyRegistered, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subId []*big.Int) (*VRFCoordinatorV25OptimismRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismRandomWordsFulfilled, requestId []*big.Int, subId []*big.Int) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV25OptimismRandomWordsFulfilled, error)

	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (*VRFCoordinatorV25OptimismRandomWordsRequestedIterator, error)

	WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismRandomWordsRequested, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV25OptimismRandomWordsRequested, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionCanceled, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionConsumerAdded, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionConsumerRemoved, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionCreated, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionFunded, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionFunded, error)

	FilterSubscriptionFundedWithNative(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionFundedWithNativeIterator, error)

	WatchSubscriptionFundedWithNative(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionFundedWithNative, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionFundedWithNative(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionFundedWithNative, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequested, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25OptimismSubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismSubscriptionOwnerTransferred, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV25OptimismSubscriptionOwnerTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
