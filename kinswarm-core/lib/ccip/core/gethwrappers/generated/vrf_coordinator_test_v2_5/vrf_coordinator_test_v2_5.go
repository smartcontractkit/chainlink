// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_coordinator_test_v2_5

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

type VRFOldProof struct {
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

var VRFCoordinatorTestV25MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendNative\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToTransferLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxGas\",\"type\":\"uint256\"}],\"name\":\"GasPriceExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"premiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"max\",\"type\":\"uint8\"}],\"name\":\"InvalidPremiumPercentage\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"flatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeNativePPM\",\"type\":\"uint32\"}],\"name\":\"LinkDiscountTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"max\",\"type\":\"uint32\"}],\"name\":\"MsgDataTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"name\":\"FallbackWeiPerUnitLinkUsed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativeFundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"onlyPremium\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountNative\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldNativeBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newNativeBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFundedWithNative\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_NATIVE_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"deregisterMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRFOld.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFTypes.RequestCommitmentV2Plus\",\"name\":\"rc\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"onlyPremium\",\"type\":\"bool\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"fundSubscriptionWithNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveSubscriptionIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"subOwner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverNativeFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"registerMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"reentrancyLock\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_currentSubNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_provingKeyHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_provingKeys\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestCommitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalNativeBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkNativeFeed\",\"type\":\"address\"}],\"name\":\"setLINKAndLINKNativeFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162005e6338038062005e6383398101604081905262000034916200017e565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000d3565b5050506001600160a01b0316608052620001b0565b336001600160a01b038216036200012d5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019157600080fd5b81516001600160a01b0381168114620001a957600080fd5b9392505050565b608051615c90620001d3600039600081816105d2015261331e0152615c906000f3fe60806040526004361061028c5760003560e01c80638402595e11610164578063b2a7cac5116100c6578063da2f26101161008a578063e72f6e3011610064578063e72f6e3014610904578063ee9d2d3814610924578063f2fde38b1461095157600080fd5b8063da2f261014610854578063dac83d29146108b3578063dc311dd3146108d357600080fd5b8063b2a7cac5146107b4578063bec4c08c146107d4578063caf70c4a146107f4578063cb63179714610814578063d98e620e1461083457600080fd5b80639d40a6fd11610128578063a63e0bfb11610102578063a63e0bfb14610747578063aa433aff14610767578063aefb212f1461078757600080fd5b80639d40a6fd146106da578063a21a23e414610712578063a4c0ed361461072757600080fd5b80638402595e1461064957806386fe91c7146106695780638da5cb5b1461068957806395b55cfc146106a75780639b1c385e146106ba57600080fd5b8063405b84fa1161020d57806364d51a2a116101d157806372e9d565116101ab57806372e9d565146105f457806379ba5097146106145780637a5a2aef1461062957600080fd5b806364d51a2a1461058b57806365982744146105a0578063689c4517146105c057600080fd5b8063405b84fa146104d057806340d6bb82146104f057806341af6c871461051b57806351cff8d91461054b5780635d06b4ab1461056b57600080fd5b806315c48b841161025457806315c48b84146103f157806318e3dd27146104195780631b6b6d23146104585780632f622e6b14610490578063301f42e9146104b057600080fd5b806304104edb14610291578063043bd6ae146102b3578063088070f5146102dc57806308821d58146103b15780630ae09540146103d1575b600080fd5b34801561029d57600080fd5b506102b16102ac366004614f21565b610971565b005b3480156102bf57600080fd5b506102c960105481565b6040519081526020015b60405180910390f35b3480156102e857600080fd5b50600c546103549061ffff81169063ffffffff62010000820481169160ff660100000000000082048116926701000000000000008304811692600160581b8104821692600160781b8204831692600160981b83041691600160b81b8104821691600160c01b9091041689565b6040805161ffff909a168a5263ffffffff98891660208b01529615159689019690965293861660608801529185166080870152841660a08601529290921660c084015260ff91821660e084015216610100820152610120016102d3565b3480156103bd57600080fd5b506102b16103cc366004614f4f565b610aea565b3480156103dd57600080fd5b506102b16103ec366004614f6b565b610ca7565b3480156103fd57600080fd5b5061040660c881565b60405161ffff90911681526020016102d3565b34801561042557600080fd5b50600a5461044090600160601b90046001600160601b031681565b6040516001600160601b0390911681526020016102d3565b34801561046457600080fd5b50600254610478906001600160a01b031681565b6040516001600160a01b0390911681526020016102d3565b34801561049c57600080fd5b506102b16104ab366004614f21565b610cef565b3480156104bc57600080fd5b506104406104cb3660046151cd565b610d95565b3480156104dc57600080fd5b506102b16104eb366004614f6b565b6110ab565b3480156104fc57600080fd5b506105066101f481565b60405163ffffffff90911681526020016102d3565b34801561052757600080fd5b5061053b6105363660046152bb565b61148d565b60405190151581526020016102d3565b34801561055757600080fd5b506102b1610566366004614f21565b611541565b34801561057757600080fd5b506102b1610586366004614f21565b611666565b34801561059757600080fd5b50610406606481565b3480156105ac57600080fd5b506102b16105bb3660046152d4565b611724565b3480156105cc57600080fd5b506104787f000000000000000000000000000000000000000000000000000000000000000081565b34801561060057600080fd5b50600354610478906001600160a01b031681565b34801561062057600080fd5b506102b1611784565b34801561063557600080fd5b506102b1610644366004615302565b611835565b34801561065557600080fd5b506102b1610664366004614f21565b611969565b34801561067557600080fd5b50600a54610440906001600160601b031681565b34801561069557600080fd5b506000546001600160a01b0316610478565b6102b16106b53660046152bb565b611a84565b3480156106c657600080fd5b506102c96106d5366004615336565b611b94565b3480156106e657600080fd5b506007546106fa906001600160401b031681565b6040516001600160401b0390911681526020016102d3565b34801561071e57600080fd5b506102c9611fda565b34801561073357600080fd5b506102b1610742366004615370565b6121c1565b34801561075357600080fd5b506102b161076236600461541b565b612329565b34801561077357600080fd5b506102b16107823660046152bb565b612610565b34801561079357600080fd5b506107a76107a23660046154bc565b612643565b6040516102d39190615519565b3480156107c057600080fd5b506102b16107cf3660046152bb565b612745565b3480156107e057600080fd5b506102b16107ef366004614f6b565b612834565b34801561080057600080fd5b506102c961080f36600461552c565b612927565b34801561082057600080fd5b506102b161082f366004614f6b565b612957565b34801561084057600080fd5b506102c961084f3660046152bb565b612bc5565b34801561086057600080fd5b5061089461086f3660046152bb565b600d6020526000908152604090205460ff81169061010090046001600160401b031682565b6040805192151583526001600160401b039091166020830152016102d3565b3480156108bf57600080fd5b506102b16108ce366004614f6b565b612be6565b3480156108df57600080fd5b506108f36108ee3660046152bb565b612c81565b6040516102d3959493929190615581565b34801561091057600080fd5b506102b161091f366004614f21565b612d5a565b34801561093057600080fd5b506102c961093f3660046152bb565b600f6020526000908152604090205481565b34801561095d57600080fd5b506102b161096c366004614f21565b612f1b565b610979612f2c565b60115460005b81811015610abd57826001600160a01b0316601182815481106109a4576109a46155d6565b6000918252602090912001546001600160a01b031603610aad5760116109cb600184615602565b815481106109db576109db6155d6565b600091825260209091200154601180546001600160a01b039092169183908110610a0757610a076155d6565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b031602179055506011805480610a4657610a46615615565b6000828152602090819020600019908301810180546001600160a01b03191690559091019091556040516001600160a01b03851681527ff80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af3791015b60405180910390a1505050565b610ab68161562b565b905061097f565b50604051635428d44960e01b81526001600160a01b03831660048201526024015b60405180910390fd5b50565b610af2612f2c565b604080518082018252600091610b21919084906002908390839080828437600092019190915250612927915050565b6000818152600d602090815260409182902082518084019093525460ff811615158084526101009091046001600160401b03169183019190915291925090610b7f57604051631dfd6e1360e21b815260048101839052602401610ade565b6000828152600d60205260408120805468ffffffffffffffffff19169055600e54905b81811015610c515783600e8281548110610bbe57610bbe6155d6565b906000526020600020015403610c4157600e610bdb600184615602565b81548110610beb57610beb6155d6565b9060005260206000200154600e8281548110610c0957610c096155d6565b600091825260209091200155600e805480610c2657610c26615615565b60019003818190600052602060002001600090559055610c51565b610c4a8161562b565b9050610ba2565b507f9b6868e0eb737bcd72205360baa6bfd0ba4e4819a33ade2db384e8a8025639a5838360200151604051610c999291909182526001600160401b0316602082015260400190565b60405180910390a150505050565b81610cb181612f88565b610cb9612fdd565b610cc28361148d565b15610ce057604051631685ecdd60e31b815260040160405180910390fd5b610cea838361300b565b505050565b610cf7612fdd565b610cff612f2c565b600b54600160601b90046001600160601b0316610d1d8115156130ee565b600b80546bffffffffffffffffffffffff60601b19169055600a8054829190600c90610d5a908490600160601b90046001600160601b0316615644565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550610d9182826001600160601b031661310c565b5050565b6000610d9f612fdd565b60005a9050610324361115610dd157604051630f28961b60e01b81523660048201526103246024820152604401610ade565b6000610ddd8686613180565b90506000610df385836000015160200151613431565b60408301516060888101519293509163ffffffff16806001600160401b03811115610e2057610e20614f9b565b604051908082528060200260200182016040528015610e49578160200160208202803683370190505b50925060005b81811015610eb15760408051602081018590529081018290526060016040516020818303038152906040528051906020012060001c848281518110610e9657610e966155d6565b6020908102919091010152610eaa8161562b565b9050610e4f565b5050602080850180516000908152600f9092526040822082905551610ed7908a8561348c565b60208a8101516000908152600690915260409020805491925090601890610f0d90600160c01b90046001600160401b0316615664565b82546101009290920a6001600160401b0381810219909316918316021790915560808a01516001600160a01b03166000908152600460209081526040808320828e01518452909152902080549091600991610f7091600160481b9091041661568a565b91906101000a8154816001600160401b0302191690836001600160401b0316021790555060008960a0015160018b60a0015151610fad9190615602565b81518110610fbd57610fbd6155d6565b60209101015160f81c60011490506000610fd98887848d613530565b909950905080156110245760208088015160105460408051928352928201527f6ca648a381f22ead7e37773d934e64885dcf861fbfbb26c40354cbf0c4662d1a910160405180910390a15b5061103488828c60200151613568565b6020808b015187820151604080518781526001600160601b038d16948101949094528415159084015284151560608401528b1515608084015290917faeb4b4786571e184246d39587f659abf0e26f41f6a3358692250382c0cdb47b79060a00160405180910390a3505050505050505b9392505050565b6110b3612fdd565b6110bc816136c5565b6110e457604051635428d44960e01b81526001600160a01b0382166004820152602401610ade565b6000806000806110f386612c81565b945094505093509350336001600160a01b0316826001600160a01b03161461115d5760405162461bcd60e51b815260206004820152601660248201527f4e6f7420737562736372697074696f6e206f776e6572000000000000000000006044820152606401610ade565b6111668661148d565b156111b35760405162461bcd60e51b815260206004820152601660248201527f50656e64696e67207265717565737420657869737473000000000000000000006044820152606401610ade565b6040805160c0810182526001815260208082018990526001600160a01b03851682840152606082018490526001600160601b038088166080840152861660a083015291519091600091611208918491016156ad565b604051602081830303815290604052905061122288613730565b505060405163ce3f471960e01b81526001600160a01b0388169063ce3f4719906001600160601b0388169061125b908590600401615772565b6000604051808303818588803b15801561127457600080fd5b505af1158015611288573d6000803e3d6000fd5b50506002546001600160a01b0316158015935091506112b1905057506001600160601b03861615155b156113815760025460405163a9059cbb60e01b81526001600160a01b0389811660048301526001600160601b03891660248301529091169063a9059cbb906044016020604051808303816000875af1158015611311573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113359190615785565b6113815760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e647300000000000000000000000000006044820152606401610ade565b600c805466ff0000000000001916660100000000000017905560005b8351811015611430578381815181106113b8576113b86155d6565b6020908102919091010151604051638ea9811760e01b81526001600160a01b038a8116600483015290911690638ea9811790602401600060405180830381600087803b15801561140757600080fd5b505af115801561141b573d6000803e3d6000fd5b50505050806114299061562b565b905061139d565b50600c805466ff00000000000019169055604080516001600160a01b0389168152602081018a90527fd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be4187910160405180910390a15050505050505050565b600081815260056020526040812060020180548083036114b1575060009392505050565b60005b81811015611536576000600460008584815481106114d4576114d46155d6565b60009182526020808320909101546001600160a01b0316835282810193909352604091820181208982529092529020546001600160401b03600160481b90910416111561152657506001949350505050565b61152f8161562b565b90506114b4565b506000949350505050565b611549612fdd565b611551612f2c565b6002546001600160a01b031661157a5760405163c1f0c0a160e01b815260040160405180910390fd5b600b546001600160601b03166115918115156130ee565b600b80546bffffffffffffffffffffffff19169055600a80548291906000906115c49084906001600160601b0316615644565b82546101009290920a6001600160601b0381810219909316918316021790915560025460405163a9059cbb60e01b81526001600160a01b0386811660048301529285166024820152610d91935091169063a9059cbb906044015b6020604051808303816000875af115801561163d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116619190615785565b6130ee565b61166e612f2c565b611677816136c5565b156116a05760405163ac8a27ef60e01b81526001600160a01b0382166004820152602401610ade565b601180546001810182556000919091527f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c680180546001600160a01b0319166001600160a01b0383169081179091556040519081527fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af016259060200160405180910390a150565b61172c612f2c565b6002546001600160a01b03161561175657604051631688c53760e11b815260040160405180910390fd5b600280546001600160a01b039384166001600160a01b03199182161790915560038054929093169116179055565b6001546001600160a01b031633146117de5760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610ade565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61183d612f2c565b60408051808201825260009161186c919085906002908390839080828437600092019190915250612927915050565b6000818152600d602052604090205490915060ff16156118a257604051634a0b8fa760e01b815260048101829052602401610ade565b60408051808201825260018082526001600160401b0385811660208085018281526000888152600d835287812096518754925168ffffffffffffffffff1990931690151568ffffffffffffffff00191617610100929095169190910293909317909455600e805493840181559091527fbb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fd9091018490558251848152918201527f9b911b2c240bfbef3b6a8f7ed6ee321d1258bb2a3fe6becab52ac1cd3210afd39101610aa0565b611971612f2c565b600a544790600160601b90046001600160601b0316818111156119b1576040516354ced18160e11b81526004810182905260248101839052604401610ade565b81811015610cea5760006119c58284615602565b90506000846001600160a01b03168260405160006040518083038185875af1925050503d8060008114611a14576040519150601f19603f3d011682016040523d82523d6000602084013e611a19565b606091505b5050905080611a3b5760405163950b247960e01b815260040160405180910390fd5b604080516001600160a01b0387168152602081018490527f4aed7c8eed0496c8c19ea2681fcca25741c1602342e38b045d9f1e8e905d2e9c910160405180910390a15050505050565b611a8c612fdd565b600081815260056020526040902054611aad906001600160a01b03166138e2565b60008181526006602052604090208054600160601b90046001600160601b0316903490600c611adc83856157a2565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555034600a600c8282829054906101000a90046001600160601b0316611b2491906157a2565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f7603b205d03651ee812f803fccde89f1012e545a9c99f0abfea9cedd0fd8e902823484611b7791906157c2565b604080519283526020830191909152015b60405180910390a25050565b6000611b9e612fdd565b602080830135600081815260059092526040909120546001600160a01b0316611bda57604051630fb532db60e11b815260040160405180910390fd5b336000908152600460209081526040808320848452808352928190208151606081018352905460ff811615158083526001600160401b036101008304811695840195909552600160481b9091049093169181019190915290611c58576040516379bfd40160e01b815260048101849052336024820152604401610ade565b600c5461ffff16611c6f60608701604088016157d5565b61ffff161080611c92575060c8611c8c60608701604088016157d5565b61ffff16115b15611cd857611ca760608601604087016157d5565b600c5460405163539c34bb60e11b815261ffff92831660048201529116602482015260c86044820152606401610ade565b600c5462010000900463ffffffff16611cf760808701606088016157f0565b63ffffffff161115611d4757611d1360808601606087016157f0565b600c54604051637aebf00f60e11b815263ffffffff9283166004820152620100009091049091166024820152604401610ade565b6101f4611d5a60a08701608088016157f0565b63ffffffff161115611da057611d7660a08601608087016157f0565b6040516311ce1afb60e21b815263ffffffff90911660048201526101f46024820152604401610ade565b806020018051611daf90615664565b6001600160401b03169052604081018051611dc990615664565b6001600160401b03908116909152602082810151604080518935818501819052338284015260608201899052929094166080808601919091528151808603909101815260a08501825280519084012060c085019290925260e08085018390528151808603909101815261010090940190528251929091019190912060009190955090506000611e6b611e66611e6160a08a018a61580b565b613909565b61398a565b905085611e766139fb565b86611e8760808b0160608c016157f0565b611e9760a08c0160808d016157f0565b3386604051602001611eaf9796959493929190615858565b60405160208183030381529060405280519060200120600f600088815260200190815260200160002081905550336001600160a01b03168588600001357feb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e89868c6040016020810190611f2291906157d5565b8d6060016020810190611f3591906157f0565b8e6080016020810190611f4891906157f0565b89604051611f5b969594939291906158af565b60405180910390a45050600092835260209182526040928390208151815493830151929094015168ffffffffffffffffff1990931693151568ffffffffffffffff001916939093176101006001600160401b03928316021770ffffffffffffffff0000000000000000001916600160481b91909216021790555b919050565b6000611fe4612fdd565b6007546001600160401b031633611ffc600143615602565b6040516bffffffffffffffffffffffff19606093841b81166020830152914060348201523090921b1660548201526001600160c01b031960c083901b16606882015260700160408051601f19818403018152919052805160209091012091506120668160016158ee565b6007805467ffffffffffffffff19166001600160401b03928316179055604080516000808252608082018352602080830182815283850183815260608086018581528a86526006855287862093518454935191516001600160601b039182166001600160c01b031990951694909417600160601b91909216021777ffffffffffffffffffffffffffffffffffffffffffffffff16600160c01b9290981691909102969096179055835194850184523385528481018281528585018481528884526005835294909220855181546001600160a01b03199081166001600160a01b0392831617835593516001830180549095169116179092559251805192949391926121769260028501920190614e0f565b5061218691506008905084613a7c565b5060405133815283907f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d9060200160405180910390a2505090565b6121c9612fdd565b6002546001600160a01b031633146121f4576040516344b0e3c360e01b815260040160405180910390fd5b6020811461221557604051638129bbcd60e01b815260040160405180910390fd5b6000612223828401846152bb565b600081815260056020526040902054909150612247906001600160a01b03166138e2565b600081815260066020526040812080546001600160601b03169186919061226e83856157a2565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600a60008282829054906101000a90046001600160601b03166122b691906157a2565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a82878461230991906157c2565b6040805192835260208301919091520160405180910390a2505050505050565b612331612f2c565b60c861ffff8a16111561236b5760405163539c34bb60e11b815261ffff8a1660048201819052602482015260c86044820152606401610ade565b6000851361238f576040516321ea67b360e11b815260048101869052602401610ade565b8363ffffffff168363ffffffff1611156123cc576040516313c06e5960e11b815263ffffffff808516600483015285166024820152604401610ade565b609b60ff831611156123fd57604051631d66288d60e11b815260ff83166004820152609b6024820152604401610ade565b609b60ff8216111561242e57604051631d66288d60e11b815260ff82166004820152609b6024820152604401610ade565b604080516101208101825261ffff8b1680825263ffffffff808c16602084018190526000848601528b8216606085018190528b8316608086018190528a841660a08701819052938a1660c0870181905260ff808b1660e08901819052908a16610100909801889052600c8054600160c01b90990260ff60c01b19600160b81b9093029290921661ffff60b81b19600160981b90940263ffffffff60981b19600160781b9099029890981676ffffffffffffffff00000000000000000000000000000019600160581b9096026effffffff000000000000000000000019670100000000000000909802979097166effffffffffffffffff000000000000196201000090990265ffffffffffff19909c16909a179a909a1796909616979097179390931791909116959095179290921793909316929092179190911790556010869055517f2c6b6b12413678366b05b145c5f00745bdd00e739131ab5de82484a50c9d78b6906125fd908b908b908b908b908b908b908b908b908b9061ffff99909916895263ffffffff97881660208a0152958716604089015293861660608801526080870192909252841660a086015290921660c084015260ff91821660e0840152166101008201526101200190565b60405180910390a1505050505050505050565b612618612f2c565b6000818152600560205260409020546001600160a01b0316612639816138e2565b610d91828261300b565b606060006126516008613a88565b905080841061267357604051631390f2a160e01b815260040160405180910390fd5b600061267f84866157c2565b90508181118061268d575083155b6126975780612699565b815b905060006126a78683615602565b9050806001600160401b038111156126c1576126c1614f9b565b6040519080825280602002602001820160405280156126ea578160200160208202803683370190505b50935060005b8181101561273a5761270d61270588836157c2565b600890613a92565b85828151811061271f5761271f6155d6565b60209081029190910101526127338161562b565b90506126f0565b505050505b92915050565b61274d612fdd565b6000818152600560205260409020546001600160a01b031661276e816138e2565b6000828152600560205260409020600101546001600160a01b031633146127c7576000828152600560205260409081902060010154905163d084e97560e01b81526001600160a01b039091166004820152602401610ade565b6000828152600560209081526040918290208054336001600160a01b03199182168117835560019092018054909116905582516001600160a01b03851681529182015283917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c93869101611b88565b8161283e81612f88565b612846612fdd565b6001600160a01b03821660009081526004602090815260408083208684529091529020805460ff16156128795750505050565b60008481526005602052604090206002018054606319016128ad576040516305a48e0f60e01b815260040160405180910390fd5b8154600160ff199091168117835581549081018255600082815260209081902090910180546001600160a01b0319166001600160a01b03871690811790915560405190815286917f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e191015b60405180910390a25050505050565b60008160405160200161293a9190615931565b604051602081830303815290604052805190602001209050919050565b8161296181612f88565b612969612fdd565b6129728361148d565b1561299057604051631685ecdd60e31b815260040160405180910390fd5b6001600160a01b038216600090815260046020908152604080832086845290915290205460ff166129e6576040516379bfd40160e01b8152600481018490526001600160a01b0383166024820152604401610ade565b600083815260056020908152604080832060020180548251818502810185019093528083529192909190830182828015612a4957602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612a2b575b50505050509050600060018251612a609190615602565b905060005b8251811015612b6957846001600160a01b0316838281518110612a8a57612a8a6155d6565b60200260200101516001600160a01b031603612b59576000838381518110612ab457612ab46155d6565b6020026020010151905080600560008981526020019081526020016000206002018381548110612ae657612ae66155d6565b600091825260208083209190910180546001600160a01b0319166001600160a01b039490941693909317909255888152600590915260409020600201805480612b3157612b31615615565b600082815260209020810160001990810180546001600160a01b031916905501905550612b69565b612b628161562b565b9050612a65565b506001600160a01b0384166000818152600460209081526040808320898452825291829020805460ff19169055905191825286917f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a79101612918565b600e8181548110612bd557600080fd5b600091825260209091200154905081565b81612bf081612f88565b612bf8612fdd565b600083815260056020526040902060018101546001600160a01b03848116911614612c7b576001810180546001600160a01b0319166001600160a01b03851690811790915560408051338152602081019290925285917f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a191015b60405180910390a25b50505050565b600081815260056020526040812054819081906001600160a01b03166060612ca8826138e2565b600086815260066020908152604080832054600583529281902060020180548251818502810185019093528083526001600160601b0380861695600160601b810490911694600160c01b9091046001600160401b0316938893929091839190830182828015612d4057602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612d22575b505050505090509450945094509450945091939590929450565b612d62612f2c565b6002546001600160a01b0316612d8b5760405163c1f0c0a160e01b815260040160405180910390fd5b6002546040516370a0823160e01b81523060048201526000916001600160a01b0316906370a0823190602401602060405180830381865afa158015612dd4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612df8919061593f565b600a549091506001600160601b031681811115612e32576040516354ced18160e11b81526004810182905260248101839052604401610ade565b81811015610cea576000612e468284615602565b60025460405163a9059cbb60e01b81526001600160a01b0387811660048301526024820184905292935091169063a9059cbb906044016020604051808303816000875af1158015612e9b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ebf9190615785565b612edc57604051631f01ff1360e21b815260040160405180910390fd5b604080516001600160a01b0386168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b4366009101610c99565b612f23612f2c565b610ae781613a9e565b6000546001600160a01b03163314612f865760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610ade565b565b6000818152600560205260409020546001600160a01b0316612fa9816138e2565b336001600160a01b03821614610d9157604051636c51fda960e11b81526001600160a01b0382166004820152602401610ade565b600c546601000000000000900460ff1615612f865760405163769dd35360e11b815260040160405180910390fd5b60008061301784613730565b60025491935091506001600160a01b03161580159061303e57506001600160601b03821615155b156130865760025460405163a9059cbb60e01b81526001600160a01b0385811660048301526001600160601b038516602483015261308692169063a9059cbb9060440161161e565b61309983826001600160601b031661310c565b604080516001600160a01b03851681526001600160601b03808516602083015283169181019190915284907f8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c490606001612c72565b80610ae757604051631e9acf1760e31b815260040160405180910390fd5b6000826001600160a01b03168260405160006040518083038185875af1925050503d8060008114613159576040519150601f19603f3d011682016040523d82523d6000602084013e61315e565b606091505b5050905080610cea5760405163950b247960e01b815260040160405180910390fd5b6040805160a081018252600060608201818152608083018290528252602082018190529181019190915260006131b98460000151612927565b6000818152600d602090815260409182902082518084019093525460ff811615158084526101009091046001600160401b0316918301919091529192509061321757604051631dfd6e1360e21b815260048101839052602401610ade565b6000828660800151604051602001613239929190918252602082015260400190565b60408051601f1981840301815291815281516020928301206000818152600f909352908220549092509081900361328357604051631b44092560e11b815260040160405180910390fd5b85516020808801516040808a015160608b015160808c015160a08d015193516132b2978a979096959101615958565b6040516020818303038152906040528051906020012081146132e75760405163354a450b60e21b815260040160405180910390fd5b60006132f68760000151613b47565b9050806133bf578651604051631d2827a760e31b81526001600160401b0390911660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063e9413d3890602401602060405180830381865afa15801561336d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613391919061593f565b9050806133bf57865160405163175dadad60e01b81526001600160401b039091166004820152602401610ade565b60008860800151826040516020016133e1929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c905060006134088a83613c1a565b604080516060810182529788526020880196909652948601949094525092979650505050505050565b6000816001600160401b03163a111561348457821561345a57506001600160401b03811661273f565b60405163435e532d60e11b81523a60048201526001600160401b0383166024820152604401610ade565b503a92915050565b6000806000631fe543e360e01b86856040516024016134ac9291906159ab565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252600c805466ff000000000000191666010000000000001790559086015160808701519192506135169163ffffffff9091169083613c85565b600c805466ff000000000000191690559695505050505050565b600080831561354f57613544868685613cd1565b60009150915061355f565b61355a868685613de2565b915091505b94509492505050565b6000818152600660205260409020821561362c5780546001600160601b03600160601b90910481169085168110156135b357604051631e9acf1760e31b815260040160405180910390fd5b6135bd8582615644565b82546bffffffffffffffffffffffff60601b1916600160601b6001600160601b039283168102919091178455600b805488939192600c926136029286929004166157a2565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555050612c7b565b80546001600160601b0390811690851681101561365c57604051631e9acf1760e31b815260040160405180910390fd5b6136668582615644565b82546bffffffffffffffffffffffff19166001600160601b03918216178355600b8054879260009161369a918591166157a2565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505050505050565b601154600090815b8181101561372657836001600160a01b0316601182815481106136f2576136f26155d6565b6000918252602090912001546001600160a01b031603613716575060019392505050565b61371f8161562b565b90506136cd565b5060009392505050565b60008181526005602090815260408083206006909252822054600290910180546001600160601b0380841694600160601b90940416925b818110156137dc5760046000848381548110613785576137856155d6565b60009182526020808320909101546001600160a01b0316835282810193909352604091820181208982529092529020805470ffffffffffffffffffffffffffffffffff191690556137d58161562b565b9050613767565b50600085815260056020526040812080546001600160a01b031990811682556001820180549091169055906138146002830182614e74565b5050600085815260066020526040812055613830600886613fd4565b506001600160601b0384161561388357600a805485919060009061385e9084906001600160601b0316615644565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b6001600160601b038316156138db5782600a600c8282829054906101000a90046001600160601b03166138b69190615644565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b5050915091565b6001600160a01b038116610ae757604051630fb532db60e11b815260040160405180910390fd5b6040805160208101909152600081526000829003613936575060408051602081019091526000815261273f565b63125fa26760e31b61394883856159cc565b6001600160e01b0319161461397057604051632923fee760e11b815260040160405180910390fd5b61397d82600481866159fc565b8101906110a49190615a26565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa826040516024016139c391511515815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b03199093169290921790915292915050565b600046613a0781613fe0565b15613a755760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613a4b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a6f919061593f565b91505090565b4391505090565b60006110a48383614003565b600061273f825490565b60006110a48383614052565b336001600160a01b03821603613af65760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610ade565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600046613b5381613fe0565b15613c0b57610100836001600160401b0316613b6d6139fb565b613b779190615602565b1180613b935750613b866139fb565b836001600160401b031610155b15613ba15750600092915050565b6040516315a03d4160e11b81526001600160401b0384166004820152606490632b407a82906024015b602060405180830381865afa158015613be7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110a4919061593f565b50506001600160401b03164090565b6000613c4e8360000151846020015185604001518660600151868860a001518960c001518a60e001518b610100015161407c565b60038360200151604051602001613c66929190615a71565b60408051601f1981840301815291905280516020909101209392505050565b60005a611388811015613c9757600080fd5b611388810390508460408204820311613caf57600080fd5b50823b613cbb57600080fd5b60008083516020850160008789f1949350505050565b600080613d146000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506142a792505050565b905060005a600c54613d34908890600160581b900463ffffffff166157c2565b613d3e9190615602565b613d489086615a85565b600c54909150600090613d6d90600160781b900463ffffffff1664e8d4a51000615a85565b90508415613db957600c548190606490600160b81b900460ff16613d9185876157c2565b613d9b9190615a85565b613da59190615ab2565b613daf91906157c2565b93505050506110a4565b600c548190606490613dd590600160b81b900460ff1682615ac6565b60ff16613d9185876157c2565b600080600080613df0614387565b9150915060008213613e18576040516321ea67b360e11b815260048101839052602401610ade565b6000613e5a6000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506142a792505050565b9050600083825a600c54613e7c908d90600160581b900463ffffffff166157c2565b613e869190615602565b613e90908b615a85565b613e9a91906157c2565b613eac90670de0b6b3a7640000615a85565b613eb69190615ab2565b600c54909150600090613edf9063ffffffff600160981b8204811691600160781b900416615adf565b613ef49063ffffffff1664e8d4a51000615a85565b9050600085613f0b83670de0b6b3a7640000615a85565b613f159190615ab2565b905060008915613f5657600c548290606490613f3b90600160c01b900460ff1687615a85565b613f459190615ab2565b613f4f91906157c2565b9050613f96565b600c548290606490613f7290600160c01b900460ff1682615ac6565b613f7f9060ff1687615a85565b613f899190615ab2565b613f9391906157c2565b90505b6b033b2e3c9fd0803ce8000000811115613fc35760405163e80fa38160e01b815260040160405180910390fd5b9b949a509398505050505050505050565b60006110a48383614452565b600061a4b1821480613ff4575062066eed82145b8061273f57505062066eee1490565b600081815260018301602052604081205461404a5750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561273f565b50600061273f565b6000826000018281548110614069576140696155d6565b9060005260206000200154905092915050565b6140858961454c565b6140d15760405162461bcd60e51b815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e2063757276650000000000006044820152606401610ade565b6140da8861454c565b6141265760405162461bcd60e51b815260206004820152601560248201527f67616d6d61206973206e6f74206f6e20637572766500000000000000000000006044820152606401610ade565b61412f8361454c565b61417b5760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610ade565b6141848261454c565b6141d05760405162461bcd60e51b815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e206375727665000000006044820152606401610ade565b6141dc878a8887614625565b6142285760405162461bcd60e51b815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e657373000000000000006044820152606401610ade565b60006142348a87614748565b90506000614247898b878b8689896147ac565b90506000614258838d8d8a866148d8565b9050808a146142995760405162461bcd60e51b815260206004820152600d60248201526c34b73b30b634b210383937b7b360991b6044820152606401610ade565b505050505050505050505050565b6000466142b381613fe0565b156142f757606c6001600160a01b031663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613be7573d6000803e3d6000fd5b61430081614918565b1561437e5773420000000000000000000000000000000000000f6001600160a01b03166349948e0e84604051806080016040528060488152602001615c3c60489139604051602001614353929190615afc565b6040516020818303038152906040526040518263ffffffff1660e01b8152600401613bca9190615772565b50600092915050565b600c5460035460408051633fabe5a360e21b81529051600093849367010000000000000090910463ffffffff169284926001600160a01b039092169163feaf968c9160048082019260a0929091908290030181865afa1580156143ee573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906144129190615b45565b50919650909250505063ffffffff82161580159061443e57506144358142615602565b8263ffffffff16105b9250821561444c5760105493505b50509091565b6000818152600183016020526040812054801561453b576000614476600183615602565b855490915060009061448a90600190615602565b90508181146144ef5760008660000182815481106144aa576144aa6155d6565b90600052602060002001549050808760000184815481106144cd576144cd6155d6565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061450057614500615615565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061273f565b600091505061273f565b5092915050565b80516000906401000003d019116145a55760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420782d6f7264696e61746500000000000000000000000000006044820152606401610ade565b60208201516401000003d019116145fe5760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420792d6f7264696e61746500000000000000000000000000006044820152606401610ade565b60208201516401000003d01990800961461e8360005b602002015161495f565b1492915050565b60006001600160a01b03821661466b5760405162461bcd60e51b815260206004820152600b60248201526a626164207769746e65737360a81b6044820152606401610ade565b60208401516000906001161561468257601c614685565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe19918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa158015614720573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b614750614e92565b61477d6001848460405160200161476993929190615b95565b604051602081830303815290604052614983565b90505b6147898161454c565b61273f5780516040805160208101929092526147a59101614769565b9050614780565b6147b4614e92565b825186516401000003d01991829006919006036148135760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610ade565b61481e8789886149d0565b61486a5760405162461bcd60e51b815260206004820152601660248201527f4669727374206d756c20636865636b206661696c6564000000000000000000006044820152606401610ade565b6148758486856149d0565b6148c15760405162461bcd60e51b815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c65640000000000000000006044820152606401610ade565b6148cc868484614afb565b98975050505050505050565b6000600286868685876040516020016148f696959493929190615bb6565b60408051601f1981840301815291905280516020909101209695505050505050565b6000600a82148061492a57506101a482145b80614937575062aa37dc82145b80614943575061210582145b80614950575062014a3382145b8061273f57505062014a341490565b6000806401000003d01980848509840990506401000003d019600782089392505050565b61498b614e92565b61499482614bc2565b81526149a96149a4826000614614565b614bfd565b6020820181905260029006600103611fd5576020810180516401000003d019039052919050565b600082600003614a105760405162461bcd60e51b815260206004820152600b60248201526a3d32b9379039b1b0b630b960a91b6044820152606401610ade565b83516020850151600090614a2690600290615c15565b15614a3257601c614a35565b601b5b9050600070014551231950b75fc4402da1732fc9bebe198387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa158015614aa7573d6000803e3d6000fd5b505050602060405103519050600086604051602001614ac69190615c29565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b614b03614e92565b835160208086015185519186015160009384938493614b2493909190614c1d565b919450925090506401000003d019858209600114614b845760405162461bcd60e51b815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a000000000000006044820152606401610ade565b60405180604001604052806401000003d01980614ba357614ba3615a9c565b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d0198110611fd557604080516020808201939093528151808203840181529082019091528051910120614bca565b600061273f826002614c166401000003d01960016157c2565b901c614cfd565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000614c5d83838585614da2565b9098509050614c6e88828e88614dc6565b9098509050614c7f88828c87614dc6565b90985090506000614c928d878b85614dc6565b9098509050614ca388828686614da2565b9098509050614cb488828e89614dc6565b9098509050818114614ce9576401000003d019818a0998506401000003d01982890997506401000003d0198183099650614ced565b8196505b5050505050509450945094915050565b600080614d08614eb0565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152614d3a614ece565b60208160c0846005600019fa925082600003614d985760405162461bcd60e51b815260206004820152601260248201527f6269674d6f64457870206661696c7572652100000000000000000000000000006044820152606401610ade565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215614e64579160200282015b82811115614e6457825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190614e2f565b50614e70929150614eec565b5090565b5080546000825590600052602060002090810190610ae79190614eec565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b80821115614e705760008155600101614eed565b6001600160a01b0381168114610ae757600080fd5b8035611fd581614f01565b600060208284031215614f3357600080fd5b81356110a481614f01565b806040810183101561273f57600080fd5b600060408284031215614f6157600080fd5b6110a48383614f3e565b60008060408385031215614f7e57600080fd5b823591506020830135614f9081614f01565b809150509250929050565b634e487b7160e01b600052604160045260246000fd5b60405160c081016001600160401b0381118282101715614fd357614fd3614f9b565b60405290565b60405161012081016001600160401b0381118282101715614fd357614fd3614f9b565b604051601f8201601f191681016001600160401b038111828210171561502457615024614f9b565b604052919050565b600082601f83011261503d57600080fd5b604051604081018181106001600160401b038211171561505f5761505f614f9b565b806040525080604084018581111561507657600080fd5b845b81811015615090578035835260209283019201615078565b509195945050505050565b80356001600160401b0381168114611fd557600080fd5b803563ffffffff81168114611fd557600080fd5b600060c082840312156150d857600080fd5b6150e0614fb1565b90506150eb8261509b565b815260208083013581830152615103604084016150b2565b6040830152615114606084016150b2565b6060830152608083013561512781614f01565b608083015260a08301356001600160401b038082111561514657600080fd5b818501915085601f83011261515a57600080fd5b81358181111561516c5761516c614f9b565b61517e601f8201601f19168501614ffc565b9150808252868482850101111561519457600080fd5b80848401858401376000848284010152508060a085015250505092915050565b8015158114610ae757600080fd5b8035611fd5816151b4565b60008060008385036101e08112156151e457600080fd5b6101a0808212156151f457600080fd5b6151fc614fd9565b9150615208878761502c565b8252615217876040880161502c565b60208301526080860135604083015260a0860135606083015260c0860135608083015261524660e08701614f16565b60a083015261010061525a8882890161502c565b60c084015261526d88610140890161502c565b60e0840152610180870135908301529093508401356001600160401b0381111561529657600080fd5b6152a2868287016150c6565b9250506152b26101c085016151c2565b90509250925092565b6000602082840312156152cd57600080fd5b5035919050565b600080604083850312156152e757600080fd5b82356152f281614f01565b91506020830135614f9081614f01565b6000806060838503121561531557600080fd5b61531f8484614f3e565b915061532d6040840161509b565b90509250929050565b60006020828403121561534857600080fd5b81356001600160401b0381111561535e57600080fd5b820160c081850312156110a457600080fd5b6000806000806060858703121561538657600080fd5b843561539181614f01565b93506020850135925060408501356001600160401b03808211156153b457600080fd5b818701915087601f8301126153c857600080fd5b8135818111156153d757600080fd5b8860208285010111156153e957600080fd5b95989497505060200194505050565b803561ffff81168114611fd557600080fd5b803560ff81168114611fd557600080fd5b60008060008060008060008060006101208a8c03121561543a57600080fd5b6154438a6153f8565b985061545160208b016150b2565b975061545f60408b016150b2565b965061546d60608b016150b2565b955060808a0135945061548260a08b016150b2565b935061549060c08b016150b2565b925061549e60e08b0161540a565b91506154ad6101008b0161540a565b90509295985092959850929598565b600080604083850312156154cf57600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b8381101561550e578151875295820195908201906001016154f2565b509495945050505050565b6020815260006110a460208301846154de565b60006040828403121561553e57600080fd5b6110a4838361502c565b600081518084526020808501945080840160005b8381101561550e5781516001600160a01b03168752958201959082019060010161555c565b60006001600160601b0380881683528087166020840152506001600160401b03851660408301526001600160a01b038416606083015260a060808301526155cb60a0830184615548565b979650505050505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b8181038181111561273f5761273f6155ec565b634e487b7160e01b600052603160045260246000fd5b60006001820161563d5761563d6155ec565b5060010190565b6001600160601b03828116828216039080821115614545576145456155ec565b60006001600160401b03808316818103615680576156806155ec565b6001019392505050565b60006001600160401b038216806156a3576156a36155ec565b6000190192915050565b6020815260ff8251166020820152602082015160408201526001600160a01b0360408301511660608201526000606083015160c060808401526156f360e0840182615548565b905060808401516001600160601b0380821660a08601528060a08701511660c086015250508091505092915050565b60005b8381101561573d578181015183820152602001615725565b50506000910152565b6000815180845261575e816020860160208601615722565b601f01601f19169290920160200192915050565b6020815260006110a46020830184615746565b60006020828403121561579757600080fd5b81516110a4816151b4565b6001600160601b03818116838216019080821115614545576145456155ec565b8082018082111561273f5761273f6155ec565b6000602082840312156157e757600080fd5b6110a4826153f8565b60006020828403121561580257600080fd5b6110a4826150b2565b6000808335601e1984360301811261582257600080fd5b8301803591506001600160401b0382111561583c57600080fd5b60200191503681900382131561585157600080fd5b9250929050565b878152866020820152856040820152600063ffffffff80871660608401528086166080840152506001600160a01b03841660a083015260e060c08301526158a260e0830184615746565b9998505050505050505050565b86815285602082015261ffff85166040820152600063ffffffff808616606084015280851660808401525060c060a08301526148cc60c0830184615746565b6001600160401b03818116838216019080821115614545576145456155ec565b8060005b6002811015612c7b578151845260209384019390910190600101615912565b6040810161273f828461590e565b60006020828403121561595157600080fd5b5051919050565b8781526001600160401b0387166020820152856040820152600063ffffffff80871660608401528086166080840152506001600160a01b03841660a083015260e060c08301526158a260e0830184615746565b8281526040602082015260006159c460408301846154de565b949350505050565b6001600160e01b031981358181169160048510156159f45780818660040360031b1b83161692505b505092915050565b60008085851115615a0c57600080fd5b83861115615a1957600080fd5b5050820193919092039150565b600060208284031215615a3857600080fd5b604051602081018181106001600160401b0382111715615a5a57615a5a614f9b565b6040528235615a68816151b4565b81529392505050565b828152606081016110a4602083018461590e565b808202811582820484141761273f5761273f6155ec565b634e487b7160e01b600052601260045260246000fd5b600082615ac157615ac1615a9c565b500490565b60ff818116838216019081111561273f5761273f6155ec565b63ffffffff828116828216039080821115614545576145456155ec565b60008351615b0e818460208801615722565b835190830190615b22818360208801615722565b01949350505050565b805169ffffffffffffffffffff81168114611fd557600080fd5b600080600080600060a08688031215615b5d57600080fd5b615b6686615b2b565b9450602086015193506040860151925060608601519150615b8960808701615b2b565b90509295509295909350565b838152615ba5602082018461590e565b606081019190915260800192915050565b868152615bc6602082018761590e565b615bd3606082018661590e565b615be060a082018561590e565b615bed60e082018461590e565b60609190911b6bffffffffffffffffffffffff19166101208201526101340195945050505050565b600082615c2457615c24615a9c565b500690565b615c33818361590e565b60400191905056fe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000813000a",
}

var VRFCoordinatorTestV25ABI = VRFCoordinatorTestV25MetaData.ABI

var VRFCoordinatorTestV25Bin = VRFCoordinatorTestV25MetaData.Bin

func DeployVRFCoordinatorTestV25(auth *bind.TransactOpts, backend bind.ContractBackend, blockhashStore common.Address) (common.Address, *types.Transaction, *VRFCoordinatorTestV25, error) {
	parsed, err := VRFCoordinatorTestV25MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorTestV25Bin), backend, blockhashStore)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorTestV25{address: address, abi: *parsed, VRFCoordinatorTestV25Caller: VRFCoordinatorTestV25Caller{contract: contract}, VRFCoordinatorTestV25Transactor: VRFCoordinatorTestV25Transactor{contract: contract}, VRFCoordinatorTestV25Filterer: VRFCoordinatorTestV25Filterer{contract: contract}}, nil
}

type VRFCoordinatorTestV25 struct {
	address common.Address
	abi     abi.ABI
	VRFCoordinatorTestV25Caller
	VRFCoordinatorTestV25Transactor
	VRFCoordinatorTestV25Filterer
}

type VRFCoordinatorTestV25Caller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorTestV25Transactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorTestV25Filterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorTestV25Session struct {
	Contract     *VRFCoordinatorTestV25
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorTestV25CallerSession struct {
	Contract *VRFCoordinatorTestV25Caller
	CallOpts bind.CallOpts
}

type VRFCoordinatorTestV25TransactorSession struct {
	Contract     *VRFCoordinatorTestV25Transactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorTestV25Raw struct {
	Contract *VRFCoordinatorTestV25
}

type VRFCoordinatorTestV25CallerRaw struct {
	Contract *VRFCoordinatorTestV25Caller
}

type VRFCoordinatorTestV25TransactorRaw struct {
	Contract *VRFCoordinatorTestV25Transactor
}

func NewVRFCoordinatorTestV25(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorTestV25, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorTestV25ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFCoordinatorTestV25(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25{address: address, abi: abi, VRFCoordinatorTestV25Caller: VRFCoordinatorTestV25Caller{contract: contract}, VRFCoordinatorTestV25Transactor: VRFCoordinatorTestV25Transactor{contract: contract}, VRFCoordinatorTestV25Filterer: VRFCoordinatorTestV25Filterer{contract: contract}}, nil
}

func NewVRFCoordinatorTestV25Caller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorTestV25Caller, error) {
	contract, err := bindVRFCoordinatorTestV25(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25Caller{contract: contract}, nil
}

func NewVRFCoordinatorTestV25Transactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorTestV25Transactor, error) {
	contract, err := bindVRFCoordinatorTestV25(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25Transactor{contract: contract}, nil
}

func NewVRFCoordinatorTestV25Filterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorTestV25Filterer, error) {
	contract, err := bindVRFCoordinatorTestV25(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25Filterer{contract: contract}, nil
}

func bindVRFCoordinatorTestV25(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFCoordinatorTestV25MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorTestV25.Contract.VRFCoordinatorTestV25Caller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.VRFCoordinatorTestV25Transactor.contract.Transfer(opts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.VRFCoordinatorTestV25Transactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorTestV25.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.contract.Transfer(opts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "BLOCKHASH_STORE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorTestV25.Contract.BLOCKHASHSTORE(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorTestV25.Contract.BLOCKHASHSTORE(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) LINK() (common.Address, error) {
	return _VRFCoordinatorTestV25.Contract.LINK(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) LINK() (common.Address, error) {
	return _VRFCoordinatorTestV25.Contract.LINK(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) LINKNATIVEFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "LINK_NATIVE_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) LINKNATIVEFEED() (common.Address, error) {
	return _VRFCoordinatorTestV25.Contract.LINKNATIVEFEED(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) LINKNATIVEFEED() (common.Address, error) {
	return _VRFCoordinatorTestV25.Contract.LINKNATIVEFEED(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorTestV25.Contract.MAXCONSUMERS(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorTestV25.Contract.MAXCONSUMERS(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) MAXNUMWORDS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "MAX_NUM_WORDS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorTestV25.Contract.MAXNUMWORDS(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorTestV25.Contract.MAXNUMWORDS(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "MAX_REQUEST_CONFIRMATIONS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorTestV25.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorTestV25.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) GetActiveSubscriptionIds(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "getActiveSubscriptionIds", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) GetActiveSubscriptionIds(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VRFCoordinatorTestV25.Contract.GetActiveSubscriptionIds(&_VRFCoordinatorTestV25.CallOpts, startIndex, maxCount)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) GetActiveSubscriptionIds(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VRFCoordinatorTestV25.Contract.GetActiveSubscriptionIds(&_VRFCoordinatorTestV25.CallOpts, startIndex, maxCount)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

	error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "getSubscription", subId)

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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinatorTestV25.Contract.GetSubscription(&_VRFCoordinatorTestV25.CallOpts, subId)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinatorTestV25.Contract.GetSubscription(&_VRFCoordinatorTestV25.CallOpts, subId)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "hashOfKey", publicKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorTestV25.Contract.HashOfKey(&_VRFCoordinatorTestV25.CallOpts, publicKey)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorTestV25.Contract.HashOfKey(&_VRFCoordinatorTestV25.CallOpts, publicKey)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) Owner() (common.Address, error) {
	return _VRFCoordinatorTestV25.Contract.Owner(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) Owner() (common.Address, error) {
	return _VRFCoordinatorTestV25.Contract.Owner(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) PendingRequestExists(opts *bind.CallOpts, subId *big.Int) (bool, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "pendingRequestExists", subId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) PendingRequestExists(subId *big.Int) (bool, error) {
	return _VRFCoordinatorTestV25.Contract.PendingRequestExists(&_VRFCoordinatorTestV25.CallOpts, subId)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) PendingRequestExists(subId *big.Int) (bool, error) {
	return _VRFCoordinatorTestV25.Contract.PendingRequestExists(&_VRFCoordinatorTestV25.CallOpts, subId)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) SConfig(opts *bind.CallOpts) (SConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "s_config")

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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) SConfig() (SConfig,

	error) {
	return _VRFCoordinatorTestV25.Contract.SConfig(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) SConfig() (SConfig,

	error) {
	return _VRFCoordinatorTestV25.Contract.SConfig(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) SCurrentSubNonce(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "s_currentSubNonce")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) SCurrentSubNonce() (uint64, error) {
	return _VRFCoordinatorTestV25.Contract.SCurrentSubNonce(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) SCurrentSubNonce() (uint64, error) {
	return _VRFCoordinatorTestV25.Contract.SCurrentSubNonce(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) SFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "s_fallbackWeiPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorTestV25.Contract.SFallbackWeiPerUnitLink(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorTestV25.Contract.SFallbackWeiPerUnitLink(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) SProvingKeyHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "s_provingKeyHashes", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) SProvingKeyHashes(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorTestV25.Contract.SProvingKeyHashes(&_VRFCoordinatorTestV25.CallOpts, arg0)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) SProvingKeyHashes(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorTestV25.Contract.SProvingKeyHashes(&_VRFCoordinatorTestV25.CallOpts, arg0)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) SProvingKeys(opts *bind.CallOpts, arg0 [32]byte) (SProvingKeys,

	error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "s_provingKeys", arg0)

	outstruct := new(SProvingKeys)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.MaxGas = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) SProvingKeys(arg0 [32]byte) (SProvingKeys,

	error) {
	return _VRFCoordinatorTestV25.Contract.SProvingKeys(&_VRFCoordinatorTestV25.CallOpts, arg0)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) SProvingKeys(arg0 [32]byte) (SProvingKeys,

	error) {
	return _VRFCoordinatorTestV25.Contract.SProvingKeys(&_VRFCoordinatorTestV25.CallOpts, arg0)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) SRequestCommitments(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "s_requestCommitments", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) SRequestCommitments(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorTestV25.Contract.SRequestCommitments(&_VRFCoordinatorTestV25.CallOpts, arg0)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) SRequestCommitments(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorTestV25.Contract.SRequestCommitments(&_VRFCoordinatorTestV25.CallOpts, arg0)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) STotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "s_totalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) STotalBalance() (*big.Int, error) {
	return _VRFCoordinatorTestV25.Contract.STotalBalance(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) STotalBalance() (*big.Int, error) {
	return _VRFCoordinatorTestV25.Contract.STotalBalance(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Caller) STotalNativeBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV25.contract.Call(opts, &out, "s_totalNativeBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) STotalNativeBalance() (*big.Int, error) {
	return _VRFCoordinatorTestV25.Contract.STotalNativeBalance(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25CallerSession) STotalNativeBalance() (*big.Int, error) {
	return _VRFCoordinatorTestV25.Contract.STotalNativeBalance(&_VRFCoordinatorTestV25.CallOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "acceptOwnership")
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.AcceptOwnership(&_VRFCoordinatorTestV25.TransactOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.AcceptOwnership(&_VRFCoordinatorTestV25.TransactOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorTestV25.TransactOpts, subId)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorTestV25.TransactOpts, subId)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "addConsumer", subId, consumer)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.AddConsumer(&_VRFCoordinatorTestV25.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.AddConsumer(&_VRFCoordinatorTestV25.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "cancelSubscription", subId, to)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.CancelSubscription(&_VRFCoordinatorTestV25.TransactOpts, subId, to)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.CancelSubscription(&_VRFCoordinatorTestV25.TransactOpts, subId, to)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "createSubscription")
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.CreateSubscription(&_VRFCoordinatorTestV25.TransactOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.CreateSubscription(&_VRFCoordinatorTestV25.TransactOpts)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) DeregisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "deregisterMigratableCoordinator", target)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) DeregisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.DeregisterMigratableCoordinator(&_VRFCoordinatorTestV25.TransactOpts, target)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) DeregisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.DeregisterMigratableCoordinator(&_VRFCoordinatorTestV25.TransactOpts, target)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "deregisterProvingKey", publicProvingKey)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.DeregisterProvingKey(&_VRFCoordinatorTestV25.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.DeregisterProvingKey(&_VRFCoordinatorTestV25.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFOldProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "fulfillRandomWords", proof, rc, onlyPremium)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) FulfillRandomWords(proof VRFOldProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.FulfillRandomWords(&_VRFCoordinatorTestV25.TransactOpts, proof, rc, onlyPremium)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) FulfillRandomWords(proof VRFOldProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.FulfillRandomWords(&_VRFCoordinatorTestV25.TransactOpts, proof, rc, onlyPremium)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) FundSubscriptionWithNative(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "fundSubscriptionWithNative", subId)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) FundSubscriptionWithNative(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.FundSubscriptionWithNative(&_VRFCoordinatorTestV25.TransactOpts, subId)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) FundSubscriptionWithNative(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.FundSubscriptionWithNative(&_VRFCoordinatorTestV25.TransactOpts, subId)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) Migrate(opts *bind.TransactOpts, subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "migrate", subId, newCoordinator)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) Migrate(subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.Migrate(&_VRFCoordinatorTestV25.TransactOpts, subId, newCoordinator)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) Migrate(subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.Migrate(&_VRFCoordinatorTestV25.TransactOpts, subId, newCoordinator)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.OnTokenTransfer(&_VRFCoordinatorTestV25.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.OnTokenTransfer(&_VRFCoordinatorTestV25.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) OwnerCancelSubscription(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "ownerCancelSubscription", subId)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) OwnerCancelSubscription(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.OwnerCancelSubscription(&_VRFCoordinatorTestV25.TransactOpts, subId)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) OwnerCancelSubscription(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.OwnerCancelSubscription(&_VRFCoordinatorTestV25.TransactOpts, subId)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "recoverFunds", to)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RecoverFunds(&_VRFCoordinatorTestV25.TransactOpts, to)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RecoverFunds(&_VRFCoordinatorTestV25.TransactOpts, to)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) RecoverNativeFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "recoverNativeFunds", to)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) RecoverNativeFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RecoverNativeFunds(&_VRFCoordinatorTestV25.TransactOpts, to)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) RecoverNativeFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RecoverNativeFunds(&_VRFCoordinatorTestV25.TransactOpts, to)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) RegisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "registerMigratableCoordinator", target)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) RegisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RegisterMigratableCoordinator(&_VRFCoordinatorTestV25.TransactOpts, target)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) RegisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RegisterMigratableCoordinator(&_VRFCoordinatorTestV25.TransactOpts, target)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) RegisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int, maxGas uint64) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "registerProvingKey", publicProvingKey, maxGas)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) RegisterProvingKey(publicProvingKey [2]*big.Int, maxGas uint64) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RegisterProvingKey(&_VRFCoordinatorTestV25.TransactOpts, publicProvingKey, maxGas)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) RegisterProvingKey(publicProvingKey [2]*big.Int, maxGas uint64) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RegisterProvingKey(&_VRFCoordinatorTestV25.TransactOpts, publicProvingKey, maxGas)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "removeConsumer", subId, consumer)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RemoveConsumer(&_VRFCoordinatorTestV25.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RemoveConsumer(&_VRFCoordinatorTestV25.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) RequestRandomWords(opts *bind.TransactOpts, req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "requestRandomWords", req)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RequestRandomWords(&_VRFCoordinatorTestV25.TransactOpts, req)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RequestRandomWords(&_VRFCoordinatorTestV25.TransactOpts, req)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorTestV25.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorTestV25.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeNativePPM uint32, fulfillmentFlatFeeLinkDiscountPPM uint32, nativePremiumPercentage uint8, linkPremiumPercentage uint8) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fulfillmentFlatFeeNativePPM, fulfillmentFlatFeeLinkDiscountPPM, nativePremiumPercentage, linkPremiumPercentage)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeNativePPM uint32, fulfillmentFlatFeeLinkDiscountPPM uint32, nativePremiumPercentage uint8, linkPremiumPercentage uint8) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.SetConfig(&_VRFCoordinatorTestV25.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fulfillmentFlatFeeNativePPM, fulfillmentFlatFeeLinkDiscountPPM, nativePremiumPercentage, linkPremiumPercentage)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeNativePPM uint32, fulfillmentFlatFeeLinkDiscountPPM uint32, nativePremiumPercentage uint8, linkPremiumPercentage uint8) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.SetConfig(&_VRFCoordinatorTestV25.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fulfillmentFlatFeeNativePPM, fulfillmentFlatFeeLinkDiscountPPM, nativePremiumPercentage, linkPremiumPercentage)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) SetLINKAndLINKNativeFeed(opts *bind.TransactOpts, link common.Address, linkNativeFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "setLINKAndLINKNativeFeed", link, linkNativeFeed)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) SetLINKAndLINKNativeFeed(link common.Address, linkNativeFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.SetLINKAndLINKNativeFeed(&_VRFCoordinatorTestV25.TransactOpts, link, linkNativeFeed)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) SetLINKAndLINKNativeFeed(link common.Address, linkNativeFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.SetLINKAndLINKNativeFeed(&_VRFCoordinatorTestV25.TransactOpts, link, linkNativeFeed)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.TransferOwnership(&_VRFCoordinatorTestV25.TransactOpts, to)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.TransferOwnership(&_VRFCoordinatorTestV25.TransactOpts, to)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) Withdraw(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "withdraw", recipient)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) Withdraw(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.Withdraw(&_VRFCoordinatorTestV25.TransactOpts, recipient)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) Withdraw(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.Withdraw(&_VRFCoordinatorTestV25.TransactOpts, recipient)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Transactor) WithdrawNative(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.contract.Transact(opts, "withdrawNative", recipient)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Session) WithdrawNative(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.WithdrawNative(&_VRFCoordinatorTestV25.TransactOpts, recipient)
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25TransactorSession) WithdrawNative(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV25.Contract.WithdrawNative(&_VRFCoordinatorTestV25.TransactOpts, recipient)
}

type VRFCoordinatorTestV25ConfigSetIterator struct {
	Event *VRFCoordinatorTestV25ConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25ConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25ConfigSet)
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
		it.Event = new(VRFCoordinatorTestV25ConfigSet)
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

func (it *VRFCoordinatorTestV25ConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25ConfigSet struct {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorTestV25ConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25ConfigSetIterator{contract: _VRFCoordinatorTestV25.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25ConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25ConfigSet)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseConfigSet(log types.Log) (*VRFCoordinatorTestV25ConfigSet, error) {
	event := new(VRFCoordinatorTestV25ConfigSet)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25CoordinatorDeregisteredIterator struct {
	Event *VRFCoordinatorTestV25CoordinatorDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25CoordinatorDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25CoordinatorDeregistered)
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
		it.Event = new(VRFCoordinatorTestV25CoordinatorDeregistered)
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

func (it *VRFCoordinatorTestV25CoordinatorDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25CoordinatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25CoordinatorDeregistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorTestV25CoordinatorDeregisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25CoordinatorDeregisteredIterator{contract: _VRFCoordinatorTestV25.contract, event: "CoordinatorDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25CoordinatorDeregistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25CoordinatorDeregistered)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseCoordinatorDeregistered(log types.Log) (*VRFCoordinatorTestV25CoordinatorDeregistered, error) {
	event := new(VRFCoordinatorTestV25CoordinatorDeregistered)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25CoordinatorRegisteredIterator struct {
	Event *VRFCoordinatorTestV25CoordinatorRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25CoordinatorRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25CoordinatorRegistered)
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
		it.Event = new(VRFCoordinatorTestV25CoordinatorRegistered)
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

func (it *VRFCoordinatorTestV25CoordinatorRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25CoordinatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25CoordinatorRegistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFCoordinatorTestV25CoordinatorRegisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25CoordinatorRegisteredIterator{contract: _VRFCoordinatorTestV25.contract, event: "CoordinatorRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25CoordinatorRegistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25CoordinatorRegistered)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseCoordinatorRegistered(log types.Log) (*VRFCoordinatorTestV25CoordinatorRegistered, error) {
	event := new(VRFCoordinatorTestV25CoordinatorRegistered)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsedIterator struct {
	Event *VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsed)
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
		it.Event = new(VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsed)
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

func (it *VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsed struct {
	RequestId              *big.Int
	FallbackWeiPerUnitLink *big.Int
	Raw                    types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterFallbackWeiPerUnitLinkUsed(opts *bind.FilterOpts) (*VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsedIterator, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "FallbackWeiPerUnitLinkUsed")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsedIterator{contract: _VRFCoordinatorTestV25.contract, event: "FallbackWeiPerUnitLinkUsed", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchFallbackWeiPerUnitLinkUsed(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsed) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "FallbackWeiPerUnitLinkUsed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsed)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "FallbackWeiPerUnitLinkUsed", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseFallbackWeiPerUnitLinkUsed(log types.Log) (*VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsed, error) {
	event := new(VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsed)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "FallbackWeiPerUnitLinkUsed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25FundsRecoveredIterator struct {
	Event *VRFCoordinatorTestV25FundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25FundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25FundsRecovered)
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
		it.Event = new(VRFCoordinatorTestV25FundsRecovered)
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

func (it *VRFCoordinatorTestV25FundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25FundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25FundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorTestV25FundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25FundsRecoveredIterator{contract: _VRFCoordinatorTestV25.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25FundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25FundsRecovered)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseFundsRecovered(log types.Log) (*VRFCoordinatorTestV25FundsRecovered, error) {
	event := new(VRFCoordinatorTestV25FundsRecovered)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25MigrationCompletedIterator struct {
	Event *VRFCoordinatorTestV25MigrationCompleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25MigrationCompletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25MigrationCompleted)
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
		it.Event = new(VRFCoordinatorTestV25MigrationCompleted)
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

func (it *VRFCoordinatorTestV25MigrationCompletedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25MigrationCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25MigrationCompleted struct {
	NewCoordinator common.Address
	SubId          *big.Int
	Raw            types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterMigrationCompleted(opts *bind.FilterOpts) (*VRFCoordinatorTestV25MigrationCompletedIterator, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "MigrationCompleted")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25MigrationCompletedIterator{contract: _VRFCoordinatorTestV25.contract, event: "MigrationCompleted", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchMigrationCompleted(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25MigrationCompleted) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "MigrationCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25MigrationCompleted)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseMigrationCompleted(log types.Log) (*VRFCoordinatorTestV25MigrationCompleted, error) {
	event := new(VRFCoordinatorTestV25MigrationCompleted)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25NativeFundsRecoveredIterator struct {
	Event *VRFCoordinatorTestV25NativeFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25NativeFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25NativeFundsRecovered)
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
		it.Event = new(VRFCoordinatorTestV25NativeFundsRecovered)
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

func (it *VRFCoordinatorTestV25NativeFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25NativeFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25NativeFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterNativeFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorTestV25NativeFundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "NativeFundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25NativeFundsRecoveredIterator{contract: _VRFCoordinatorTestV25.contract, event: "NativeFundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchNativeFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25NativeFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "NativeFundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25NativeFundsRecovered)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "NativeFundsRecovered", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseNativeFundsRecovered(log types.Log) (*VRFCoordinatorTestV25NativeFundsRecovered, error) {
	event := new(VRFCoordinatorTestV25NativeFundsRecovered)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "NativeFundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25OwnershipTransferRequestedIterator struct {
	Event *VRFCoordinatorTestV25OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25OwnershipTransferRequested)
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
		it.Event = new(VRFCoordinatorTestV25OwnershipTransferRequested)
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

func (it *VRFCoordinatorTestV25OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorTestV25OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25OwnershipTransferRequestedIterator{contract: _VRFCoordinatorTestV25.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25OwnershipTransferRequested)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorTestV25OwnershipTransferRequested, error) {
	event := new(VRFCoordinatorTestV25OwnershipTransferRequested)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25OwnershipTransferredIterator struct {
	Event *VRFCoordinatorTestV25OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25OwnershipTransferred)
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
		it.Event = new(VRFCoordinatorTestV25OwnershipTransferred)
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

func (it *VRFCoordinatorTestV25OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorTestV25OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25OwnershipTransferredIterator{contract: _VRFCoordinatorTestV25.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25OwnershipTransferred)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorTestV25OwnershipTransferred, error) {
	event := new(VRFCoordinatorTestV25OwnershipTransferred)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25ProvingKeyDeregisteredIterator struct {
	Event *VRFCoordinatorTestV25ProvingKeyDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25ProvingKeyDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25ProvingKeyDeregistered)
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
		it.Event = new(VRFCoordinatorTestV25ProvingKeyDeregistered)
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

func (it *VRFCoordinatorTestV25ProvingKeyDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25ProvingKeyDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25ProvingKeyDeregistered struct {
	KeyHash [32]byte
	MaxGas  uint64
	Raw     types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterProvingKeyDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorTestV25ProvingKeyDeregisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "ProvingKeyDeregistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25ProvingKeyDeregisteredIterator{contract: _VRFCoordinatorTestV25.contract, event: "ProvingKeyDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25ProvingKeyDeregistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "ProvingKeyDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25ProvingKeyDeregistered)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorTestV25ProvingKeyDeregistered, error) {
	event := new(VRFCoordinatorTestV25ProvingKeyDeregistered)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25ProvingKeyRegisteredIterator struct {
	Event *VRFCoordinatorTestV25ProvingKeyRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25ProvingKeyRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25ProvingKeyRegistered)
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
		it.Event = new(VRFCoordinatorTestV25ProvingKeyRegistered)
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

func (it *VRFCoordinatorTestV25ProvingKeyRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25ProvingKeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25ProvingKeyRegistered struct {
	KeyHash [32]byte
	MaxGas  uint64
	Raw     types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterProvingKeyRegistered(opts *bind.FilterOpts) (*VRFCoordinatorTestV25ProvingKeyRegisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "ProvingKeyRegistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25ProvingKeyRegisteredIterator{contract: _VRFCoordinatorTestV25.contract, event: "ProvingKeyRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25ProvingKeyRegistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "ProvingKeyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25ProvingKeyRegistered)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorTestV25ProvingKeyRegistered, error) {
	event := new(VRFCoordinatorTestV25ProvingKeyRegistered)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25RandomWordsFulfilledIterator struct {
	Event *VRFCoordinatorTestV25RandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25RandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25RandomWordsFulfilled)
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
		it.Event = new(VRFCoordinatorTestV25RandomWordsFulfilled)
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

func (it *VRFCoordinatorTestV25RandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25RandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25RandomWordsFulfilled struct {
	RequestId     *big.Int
	OutputSeed    *big.Int
	SubId         *big.Int
	Payment       *big.Int
	NativePayment bool
	Success       bool
	OnlyPremium   bool
	Raw           types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subId []*big.Int) (*VRFCoordinatorTestV25RandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule, subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25RandomWordsFulfilledIterator{contract: _VRFCoordinatorTestV25.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25RandomWordsFulfilled, requestId []*big.Int, subId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule, subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25RandomWordsFulfilled)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorTestV25RandomWordsFulfilled, error) {
	event := new(VRFCoordinatorTestV25RandomWordsFulfilled)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25RandomWordsRequestedIterator struct {
	Event *VRFCoordinatorTestV25RandomWordsRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25RandomWordsRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25RandomWordsRequested)
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
		it.Event = new(VRFCoordinatorTestV25RandomWordsRequested)
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

func (it *VRFCoordinatorTestV25RandomWordsRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25RandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25RandomWordsRequested struct {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (*VRFCoordinatorTestV25RandomWordsRequestedIterator, error) {

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

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25RandomWordsRequestedIterator{contract: _VRFCoordinatorTestV25.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25RandomWordsRequested, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25RandomWordsRequested)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorTestV25RandomWordsRequested, error) {
	event := new(VRFCoordinatorTestV25RandomWordsRequested)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25SubscriptionCanceledIterator struct {
	Event *VRFCoordinatorTestV25SubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25SubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25SubscriptionCanceled)
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
		it.Event = new(VRFCoordinatorTestV25SubscriptionCanceled)
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

func (it *VRFCoordinatorTestV25SubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25SubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25SubscriptionCanceled struct {
	SubId        *big.Int
	To           common.Address
	AmountLink   *big.Int
	AmountNative *big.Int
	Raw          types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25SubscriptionCanceledIterator{contract: _VRFCoordinatorTestV25.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionCanceled, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25SubscriptionCanceled)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorTestV25SubscriptionCanceled, error) {
	event := new(VRFCoordinatorTestV25SubscriptionCanceled)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25SubscriptionConsumerAddedIterator struct {
	Event *VRFCoordinatorTestV25SubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25SubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25SubscriptionConsumerAdded)
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
		it.Event = new(VRFCoordinatorTestV25SubscriptionConsumerAdded)
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

func (it *VRFCoordinatorTestV25SubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25SubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25SubscriptionConsumerAdded struct {
	SubId    *big.Int
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25SubscriptionConsumerAddedIterator{contract: _VRFCoordinatorTestV25.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionConsumerAdded, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25SubscriptionConsumerAdded)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorTestV25SubscriptionConsumerAdded, error) {
	event := new(VRFCoordinatorTestV25SubscriptionConsumerAdded)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25SubscriptionConsumerRemovedIterator struct {
	Event *VRFCoordinatorTestV25SubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25SubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25SubscriptionConsumerRemoved)
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
		it.Event = new(VRFCoordinatorTestV25SubscriptionConsumerRemoved)
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

func (it *VRFCoordinatorTestV25SubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25SubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25SubscriptionConsumerRemoved struct {
	SubId    *big.Int
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25SubscriptionConsumerRemovedIterator{contract: _VRFCoordinatorTestV25.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionConsumerRemoved, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25SubscriptionConsumerRemoved)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorTestV25SubscriptionConsumerRemoved, error) {
	event := new(VRFCoordinatorTestV25SubscriptionConsumerRemoved)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25SubscriptionCreatedIterator struct {
	Event *VRFCoordinatorTestV25SubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25SubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25SubscriptionCreated)
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
		it.Event = new(VRFCoordinatorTestV25SubscriptionCreated)
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

func (it *VRFCoordinatorTestV25SubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25SubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25SubscriptionCreated struct {
	SubId *big.Int
	Owner common.Address
	Raw   types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25SubscriptionCreatedIterator{contract: _VRFCoordinatorTestV25.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionCreated, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25SubscriptionCreated)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorTestV25SubscriptionCreated, error) {
	event := new(VRFCoordinatorTestV25SubscriptionCreated)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25SubscriptionFundedIterator struct {
	Event *VRFCoordinatorTestV25SubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25SubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25SubscriptionFunded)
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
		it.Event = new(VRFCoordinatorTestV25SubscriptionFunded)
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

func (it *VRFCoordinatorTestV25SubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25SubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25SubscriptionFunded struct {
	SubId      *big.Int
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25SubscriptionFundedIterator{contract: _VRFCoordinatorTestV25.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionFunded, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25SubscriptionFunded)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorTestV25SubscriptionFunded, error) {
	event := new(VRFCoordinatorTestV25SubscriptionFunded)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25SubscriptionFundedWithNativeIterator struct {
	Event *VRFCoordinatorTestV25SubscriptionFundedWithNative

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25SubscriptionFundedWithNativeIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25SubscriptionFundedWithNative)
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
		it.Event = new(VRFCoordinatorTestV25SubscriptionFundedWithNative)
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

func (it *VRFCoordinatorTestV25SubscriptionFundedWithNativeIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25SubscriptionFundedWithNativeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25SubscriptionFundedWithNative struct {
	SubId            *big.Int
	OldNativeBalance *big.Int
	NewNativeBalance *big.Int
	Raw              types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterSubscriptionFundedWithNative(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionFundedWithNativeIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "SubscriptionFundedWithNative", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25SubscriptionFundedWithNativeIterator{contract: _VRFCoordinatorTestV25.contract, event: "SubscriptionFundedWithNative", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchSubscriptionFundedWithNative(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionFundedWithNative, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "SubscriptionFundedWithNative", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25SubscriptionFundedWithNative)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionFundedWithNative", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseSubscriptionFundedWithNative(log types.Log) (*VRFCoordinatorTestV25SubscriptionFundedWithNative, error) {
	event := new(VRFCoordinatorTestV25SubscriptionFundedWithNative)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionFundedWithNative", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25SubscriptionOwnerTransferRequestedIterator struct {
	Event *VRFCoordinatorTestV25SubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25SubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25SubscriptionOwnerTransferRequested)
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
		it.Event = new(VRFCoordinatorTestV25SubscriptionOwnerTransferRequested)
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

func (it *VRFCoordinatorTestV25SubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25SubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25SubscriptionOwnerTransferRequested struct {
	SubId *big.Int
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25SubscriptionOwnerTransferRequestedIterator{contract: _VRFCoordinatorTestV25.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionOwnerTransferRequested, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25SubscriptionOwnerTransferRequested)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorTestV25SubscriptionOwnerTransferRequested, error) {
	event := new(VRFCoordinatorTestV25SubscriptionOwnerTransferRequested)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV25SubscriptionOwnerTransferredIterator struct {
	Event *VRFCoordinatorTestV25SubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV25SubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV25SubscriptionOwnerTransferred)
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
		it.Event = new(VRFCoordinatorTestV25SubscriptionOwnerTransferred)
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

func (it *VRFCoordinatorTestV25SubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV25SubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV25SubscriptionOwnerTransferred struct {
	SubId *big.Int
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV25SubscriptionOwnerTransferredIterator{contract: _VRFCoordinatorTestV25.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionOwnerTransferred, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV25.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV25SubscriptionOwnerTransferred)
				if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25Filterer) ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorTestV25SubscriptionOwnerTransferred, error) {
	event := new(VRFCoordinatorTestV25SubscriptionOwnerTransferred)
	if err := _VRFCoordinatorTestV25.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinatorTestV25.abi.Events["ConfigSet"].ID:
		return _VRFCoordinatorTestV25.ParseConfigSet(log)
	case _VRFCoordinatorTestV25.abi.Events["CoordinatorDeregistered"].ID:
		return _VRFCoordinatorTestV25.ParseCoordinatorDeregistered(log)
	case _VRFCoordinatorTestV25.abi.Events["CoordinatorRegistered"].ID:
		return _VRFCoordinatorTestV25.ParseCoordinatorRegistered(log)
	case _VRFCoordinatorTestV25.abi.Events["FallbackWeiPerUnitLinkUsed"].ID:
		return _VRFCoordinatorTestV25.ParseFallbackWeiPerUnitLinkUsed(log)
	case _VRFCoordinatorTestV25.abi.Events["FundsRecovered"].ID:
		return _VRFCoordinatorTestV25.ParseFundsRecovered(log)
	case _VRFCoordinatorTestV25.abi.Events["MigrationCompleted"].ID:
		return _VRFCoordinatorTestV25.ParseMigrationCompleted(log)
	case _VRFCoordinatorTestV25.abi.Events["NativeFundsRecovered"].ID:
		return _VRFCoordinatorTestV25.ParseNativeFundsRecovered(log)
	case _VRFCoordinatorTestV25.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFCoordinatorTestV25.ParseOwnershipTransferRequested(log)
	case _VRFCoordinatorTestV25.abi.Events["OwnershipTransferred"].ID:
		return _VRFCoordinatorTestV25.ParseOwnershipTransferred(log)
	case _VRFCoordinatorTestV25.abi.Events["ProvingKeyDeregistered"].ID:
		return _VRFCoordinatorTestV25.ParseProvingKeyDeregistered(log)
	case _VRFCoordinatorTestV25.abi.Events["ProvingKeyRegistered"].ID:
		return _VRFCoordinatorTestV25.ParseProvingKeyRegistered(log)
	case _VRFCoordinatorTestV25.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFCoordinatorTestV25.ParseRandomWordsFulfilled(log)
	case _VRFCoordinatorTestV25.abi.Events["RandomWordsRequested"].ID:
		return _VRFCoordinatorTestV25.ParseRandomWordsRequested(log)
	case _VRFCoordinatorTestV25.abi.Events["SubscriptionCanceled"].ID:
		return _VRFCoordinatorTestV25.ParseSubscriptionCanceled(log)
	case _VRFCoordinatorTestV25.abi.Events["SubscriptionConsumerAdded"].ID:
		return _VRFCoordinatorTestV25.ParseSubscriptionConsumerAdded(log)
	case _VRFCoordinatorTestV25.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _VRFCoordinatorTestV25.ParseSubscriptionConsumerRemoved(log)
	case _VRFCoordinatorTestV25.abi.Events["SubscriptionCreated"].ID:
		return _VRFCoordinatorTestV25.ParseSubscriptionCreated(log)
	case _VRFCoordinatorTestV25.abi.Events["SubscriptionFunded"].ID:
		return _VRFCoordinatorTestV25.ParseSubscriptionFunded(log)
	case _VRFCoordinatorTestV25.abi.Events["SubscriptionFundedWithNative"].ID:
		return _VRFCoordinatorTestV25.ParseSubscriptionFundedWithNative(log)
	case _VRFCoordinatorTestV25.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _VRFCoordinatorTestV25.ParseSubscriptionOwnerTransferRequested(log)
	case _VRFCoordinatorTestV25.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _VRFCoordinatorTestV25.ParseSubscriptionOwnerTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFCoordinatorTestV25ConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2c6b6b12413678366b05b145c5f00745bdd00e739131ab5de82484a50c9d78b6")
}

func (VRFCoordinatorTestV25CoordinatorDeregistered) Topic() common.Hash {
	return common.HexToHash("0xf80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af37")
}

func (VRFCoordinatorTestV25CoordinatorRegistered) Topic() common.Hash {
	return common.HexToHash("0xb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af01625")
}

func (VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsed) Topic() common.Hash {
	return common.HexToHash("0x6ca648a381f22ead7e37773d934e64885dcf861fbfbb26c40354cbf0c4662d1a")
}

func (VRFCoordinatorTestV25FundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (VRFCoordinatorTestV25MigrationCompleted) Topic() common.Hash {
	return common.HexToHash("0xd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be4187")
}

func (VRFCoordinatorTestV25NativeFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x4aed7c8eed0496c8c19ea2681fcca25741c1602342e38b045d9f1e8e905d2e9c")
}

func (VRFCoordinatorTestV25OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFCoordinatorTestV25OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFCoordinatorTestV25ProvingKeyDeregistered) Topic() common.Hash {
	return common.HexToHash("0x9b6868e0eb737bcd72205360baa6bfd0ba4e4819a33ade2db384e8a8025639a5")
}

func (VRFCoordinatorTestV25ProvingKeyRegistered) Topic() common.Hash {
	return common.HexToHash("0x9b911b2c240bfbef3b6a8f7ed6ee321d1258bb2a3fe6becab52ac1cd3210afd3")
}

func (VRFCoordinatorTestV25RandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0xaeb4b4786571e184246d39587f659abf0e26f41f6a3358692250382c0cdb47b7")
}

func (VRFCoordinatorTestV25RandomWordsRequested) Topic() common.Hash {
	return common.HexToHash("0xeb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e")
}

func (VRFCoordinatorTestV25SubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0x8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c4")
}

func (VRFCoordinatorTestV25SubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e1")
}

func (VRFCoordinatorTestV25SubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a7")
}

func (VRFCoordinatorTestV25SubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d")
}

func (VRFCoordinatorTestV25SubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0x1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a")
}

func (VRFCoordinatorTestV25SubscriptionFundedWithNative) Topic() common.Hash {
	return common.HexToHash("0x7603b205d03651ee812f803fccde89f1012e545a9c99f0abfea9cedd0fd8e902")
}

func (VRFCoordinatorTestV25SubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a1")
}

func (VRFCoordinatorTestV25SubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0xd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c9386")
}

func (_VRFCoordinatorTestV25 *VRFCoordinatorTestV25) Address() common.Address {
	return _VRFCoordinatorTestV25.address
}

type VRFCoordinatorTestV25Interface interface {
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

	FulfillRandomWords(opts *bind.TransactOpts, proof VRFOldProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error)

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

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error)

	WithdrawNative(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorTestV25ConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25ConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFCoordinatorTestV25ConfigSet, error)

	FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorTestV25CoordinatorDeregisteredIterator, error)

	WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25CoordinatorDeregistered) (event.Subscription, error)

	ParseCoordinatorDeregistered(log types.Log) (*VRFCoordinatorTestV25CoordinatorDeregistered, error)

	FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFCoordinatorTestV25CoordinatorRegisteredIterator, error)

	WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25CoordinatorRegistered) (event.Subscription, error)

	ParseCoordinatorRegistered(log types.Log) (*VRFCoordinatorTestV25CoordinatorRegistered, error)

	FilterFallbackWeiPerUnitLinkUsed(opts *bind.FilterOpts) (*VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsedIterator, error)

	WatchFallbackWeiPerUnitLinkUsed(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsed) (event.Subscription, error)

	ParseFallbackWeiPerUnitLinkUsed(log types.Log) (*VRFCoordinatorTestV25FallbackWeiPerUnitLinkUsed, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorTestV25FundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25FundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*VRFCoordinatorTestV25FundsRecovered, error)

	FilterMigrationCompleted(opts *bind.FilterOpts) (*VRFCoordinatorTestV25MigrationCompletedIterator, error)

	WatchMigrationCompleted(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25MigrationCompleted) (event.Subscription, error)

	ParseMigrationCompleted(log types.Log) (*VRFCoordinatorTestV25MigrationCompleted, error)

	FilterNativeFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorTestV25NativeFundsRecoveredIterator, error)

	WatchNativeFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25NativeFundsRecovered) (event.Subscription, error)

	ParseNativeFundsRecovered(log types.Log) (*VRFCoordinatorTestV25NativeFundsRecovered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorTestV25OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorTestV25OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorTestV25OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorTestV25OwnershipTransferred, error)

	FilterProvingKeyDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorTestV25ProvingKeyDeregisteredIterator, error)

	WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25ProvingKeyDeregistered) (event.Subscription, error)

	ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorTestV25ProvingKeyDeregistered, error)

	FilterProvingKeyRegistered(opts *bind.FilterOpts) (*VRFCoordinatorTestV25ProvingKeyRegisteredIterator, error)

	WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25ProvingKeyRegistered) (event.Subscription, error)

	ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorTestV25ProvingKeyRegistered, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subId []*big.Int) (*VRFCoordinatorTestV25RandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25RandomWordsFulfilled, requestId []*big.Int, subId []*big.Int) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorTestV25RandomWordsFulfilled, error)

	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (*VRFCoordinatorTestV25RandomWordsRequestedIterator, error)

	WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25RandomWordsRequested, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorTestV25RandomWordsRequested, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionCanceled, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorTestV25SubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionConsumerAdded, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorTestV25SubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionConsumerRemoved, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorTestV25SubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionCreated, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorTestV25SubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionFunded, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorTestV25SubscriptionFunded, error)

	FilterSubscriptionFundedWithNative(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionFundedWithNativeIterator, error)

	WatchSubscriptionFundedWithNative(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionFundedWithNative, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionFundedWithNative(log types.Log) (*VRFCoordinatorTestV25SubscriptionFundedWithNative, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionOwnerTransferRequested, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorTestV25SubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorTestV25SubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV25SubscriptionOwnerTransferred, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorTestV25SubscriptionOwnerTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
