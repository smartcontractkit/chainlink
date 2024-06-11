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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendNative\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToTransferLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxGas\",\"type\":\"uint256\"}],\"name\":\"GasPriceExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"mode\",\"type\":\"uint8\"}],\"name\":\"InvalidL1FeeCalculationMode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"coefficient\",\"type\":\"uint8\"}],\"name\":\"InvalidL1FeeCoefficient\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"premiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"max\",\"type\":\"uint8\"}],\"name\":\"InvalidPremiumPercentage\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"flatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeNativePPM\",\"type\":\"uint32\"}],\"name\":\"LinkDiscountTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"max\",\"type\":\"uint32\"}],\"name\":\"MsgDataTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"name\":\"FallbackWeiPerUnitLinkUsed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"mode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"coefficient\",\"type\":\"uint8\"}],\"name\":\"L1FeeCalculationSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"L1GasFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativeFundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"onlyPremium\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountNative\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldNativeBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newNativeBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFundedWithNative\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_NATIVE_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"deregisterMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFTypes.RequestCommitmentV2Plus\",\"name\":\"rc\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"onlyPremium\",\"type\":\"bool\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"fundSubscriptionWithNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveSubscriptionIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"subOwner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverNativeFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"registerMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"reentrancyLock\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_currentSubNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_l1FeeCalculationMode\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_l1FeeCoefficient\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_provingKeyHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_provingKeys\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestCommitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalNativeBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"mode\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"coefficient\",\"type\":\"uint8\"}],\"name\":\"setL1FeeCalculation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkNativeFeed\",\"type\":\"address\"}],\"name\":\"setLINKAndLINKNativeFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040526012805461ffff19166164001790553480156200002057600080fd5b50604051620061773803806200617783398101604081905262000043916200018f565b8033806000816200009b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000ce57620000ce81620000e4565b5050506001600160a01b031660805250620001c1565b336001600160a01b038216036200013e5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000092565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001a257600080fd5b81516001600160a01b0381168114620001ba57600080fd5b9392505050565b608051615f93620001e46000396000818161066f015261331b0152615f936000f3fe6080604052600436106102dd5760003560e01c80638402595e1161017f578063aefb212f116100e1578063da2f26101161008a578063e72f6e3011610064578063e72f6e30146109c0578063ee9d2d38146109e0578063f2fde38b14610a0d57600080fd5b8063da2f261014610910578063dac83d291461096f578063dc311dd31461098f57600080fd5b8063caf70c4a116100bb578063caf70c4a146108b0578063cb631797146108d0578063d98e620e146108f057600080fd5b8063aefb212f14610843578063b2a7cac514610870578063bec4c08c1461089057600080fd5b80639b1c385e11610143578063a4c0ed361161011d578063a4c0ed36146107e3578063a63e0bfb14610803578063aa433aff1461082357600080fd5b80639b1c385e146107765780639d40a6fd14610796578063a21a23e4146107ce57600080fd5b80638402595e146106e657806386fe91c7146107065780638da5cb5b1461072657806390bd5c741461074457806395b55cfc1461076357600080fd5b8063405b84fa1161024357806364d51a2a116101ec57806372e9d565116101c657806372e9d5651461069157806379ba5097146106b15780637a5a2aef146106c657600080fd5b806364d51a2a14610628578063659827441461063d578063689c45171461065d57600080fd5b806341af6c871161021d57806341af6c87146105b857806351cff8d9146105e85780635d06b4ab1461060857600080fd5b8063405b84fa1461054157806340d6bb821461056157806340e3290f1461058c57600080fd5b806314530741116102a55780631b6b6d231161027f5780631b6b6d23146104c95780632f622e6b14610501578063301f42e91461052157600080fd5b8063145307411461044257806315c48b841461046257806318e3dd271461048a57600080fd5b806304104edb146102e2578063043bd6ae14610304578063088070f51461032d57806308821d58146104025780630ae0954014610422575b600080fd5b3480156102ee57600080fd5b506103026102fd3660046150bf565b610a2d565b005b34801561031057600080fd5b5061031a60105481565b6040519081526020015b60405180910390f35b34801561033957600080fd5b50600c546103a59061ffff81169063ffffffff62010000820481169160ff660100000000000082048116926701000000000000008304811692600160581b8104821692600160781b8204831692600160981b83041691600160b81b8104821691600160c01b9091041689565b6040805161ffff909a168a5263ffffffff98891660208b01529615159689019690965293861660608801529185166080870152841660a08601529290921660c084015260ff91821660e08401521661010082015261012001610324565b34801561040e57600080fd5b5061030261041d3660046150ed565b610ba6565b34801561042e57600080fd5b5061030261043d366004615109565b610d63565b34801561044e57600080fd5b5061030261045d36600461514a565b610dab565b34801561046e57600080fd5b5061047760c881565b60405161ffff9091168152602001610324565b34801561049657600080fd5b50600a546104b190600160601b90046001600160601b031681565b6040516001600160601b039091168152602001610324565b3480156104d557600080fd5b506002546104e9906001600160a01b031681565b6040516001600160a01b039091168152602001610324565b34801561050d57600080fd5b5061030261051c3660046150bf565b610dc1565b34801561052d57600080fd5b506104b161053c3660046153ba565b610e63565b34801561054d57600080fd5b5061030261055c366004615109565b611179565b34801561056d57600080fd5b506105776101f481565b60405163ffffffff9091168152602001610324565b34801561059857600080fd5b506012546105a69060ff1681565b60405160ff9091168152602001610324565b3480156105c457600080fd5b506105d86105d33660046154a8565b6114c7565b6040519015158152602001610324565b3480156105f457600080fd5b506103026106033660046150bf565b611569565b34801561061457600080fd5b506103026106233660046150bf565b61164a565b34801561063457600080fd5b50610477606481565b34801561064957600080fd5b506103026106583660046154c1565b611708565b34801561066957600080fd5b506104e97f000000000000000000000000000000000000000000000000000000000000000081565b34801561069d57600080fd5b506003546104e9906001600160a01b031681565b3480156106bd57600080fd5b50610302611768565b3480156106d257600080fd5b506103026106e13660046154ef565b611819565b3480156106f257600080fd5b506103026107013660046150bf565b61194d565b34801561071257600080fd5b50600a546104b1906001600160601b031681565b34801561073257600080fd5b506000546001600160a01b03166104e9565b34801561075057600080fd5b506012546105a690610100900460ff1681565b6103026107713660046154a8565b611a68565b34801561078257600080fd5b5061031a61079136600461551a565b611b78565b3480156107a257600080fd5b506007546107b6906001600160401b031681565b6040516001600160401b039091168152602001610324565b3480156107da57600080fd5b5061031a611fa3565b3480156107ef57600080fd5b506103026107fe366004615554565b61218a565b34801561080f57600080fd5b5061030261081e3660046155ee565b6122f2565b34801561082f57600080fd5b5061030261083e3660046154a8565b6125c6565b34801561084f57600080fd5b5061086361085e366004615699565b6125f9565b60405161032491906156f6565b34801561087c57600080fd5b5061030261088b3660046154a8565b6126fb565b34801561089c57600080fd5b506103026108ab366004615109565b6127ea565b3480156108bc57600080fd5b5061031a6108cb366004615709565b6128dd565b3480156108dc57600080fd5b506103026108eb366004615109565b61290d565b3480156108fc57600080fd5b5061031a61090b3660046154a8565b612b0d565b34801561091c57600080fd5b5061095061092b3660046154a8565b600d6020526000908152604090205460ff81169061010090046001600160401b031682565b6040805192151583526001600160401b03909116602083015201610324565b34801561097b57600080fd5b5061030261098a366004615109565b612b2e565b34801561099b57600080fd5b506109af6109aa3660046154a8565b612bc9565b60405161032495949392919061575e565b3480156109cc57600080fd5b506103026109db3660046150bf565b612ca2565b3480156109ec57600080fd5b5061031a6109fb3660046154a8565b600f6020526000908152604090205481565b348015610a1957600080fd5b50610302610a283660046150bf565b612e63565b610a35612e74565b60115460005b81811015610b7957826001600160a01b031660118281548110610a6057610a606157b3565b6000918252602090912001546001600160a01b031603610b69576011610a876001846157df565b81548110610a9757610a976157b3565b600091825260209091200154601180546001600160a01b039092169183908110610ac357610ac36157b3565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b031602179055506011805480610b0257610b026157f2565b6000828152602090819020600019908301810180546001600160a01b03191690559091019091556040516001600160a01b03851681527ff80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af3791015b60405180910390a1505050565b610b7281615808565b9050610a3b565b50604051635428d44960e01b81526001600160a01b03831660048201526024015b60405180910390fd5b50565b610bae612e74565b604080518082018252600091610bdd9190849060029083908390808284376000920191909152506128dd915050565b6000818152600d602090815260409182902082518084019093525460ff811615158084526101009091046001600160401b03169183019190915291925090610c3b57604051631dfd6e1360e21b815260048101839052602401610b9a565b6000828152600d60205260408120805468ffffffffffffffffff19169055600e54905b81811015610d0d5783600e8281548110610c7a57610c7a6157b3565b906000526020600020015403610cfd57600e610c976001846157df565b81548110610ca757610ca76157b3565b9060005260206000200154600e8281548110610cc557610cc56157b3565b600091825260209091200155600e805480610ce257610ce26157f2565b60019003818190600052602060002001600090559055610d0d565b610d0681615808565b9050610c5e565b507f9b6868e0eb737bcd72205360baa6bfd0ba4e4819a33ade2db384e8a8025639a5838360200151604051610d559291909182526001600160401b0316602082015260400190565b60405180910390a150505050565b81610d6d81612ed0565b610d75612f25565b610d7e836114c7565b15610d9c57604051631685ecdd60e31b815260040160405180910390fd5b610da68383612f53565b505050565b610db3612e74565b610dbd8282613036565b5050565b610dc9612f25565b610dd1612e74565b600b54600160601b90046001600160601b0316610def8115156130f0565b600b80546bffffffffffffffffffffffff60601b19169055600a8054829190600c90610e2c908490600160601b90046001600160601b0316615821565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550610dbd82826001600160601b031661310e565b6000610e6d612f25565b60005a9050610324361115610e9f57604051630f28961b60e01b81523660048201526103246024820152604401610b9a565b6000610eab8686613182565b90506000610ec18583600001516020015161342e565b60408301516060888101519293509163ffffffff16806001600160401b03811115610eee57610eee61517d565b604051908082528060200260200182016040528015610f17578160200160208202803683370190505b50925060005b81811015610f7f5760408051602081018590529081018290526060016040516020818303038152906040528051906020012060001c848281518110610f6457610f646157b3565b6020908102919091010152610f7881615808565b9050610f1d565b5050602080850180516000908152600f9092526040822082905551610fa5908a85613489565b60208a8101516000908152600690915260409020805491925090601890610fdb90600160c01b90046001600160401b0316615841565b82546101009290920a6001600160401b0381810219909316918316021790915560808a01516001600160a01b03166000908152600460209081526040808320828e0151845290915290208054909160099161103e91600160481b90910416615867565b91906101000a8154816001600160401b0302191690836001600160401b0316021790555060008960a0015160018b60a001515161107b91906157df565b8151811061108b5761108b6157b3565b60209101015160f81c600114905060006110a78887848d61352d565b909950905080156110f25760208088015160105460408051928352928201527f6ca648a381f22ead7e37773d934e64885dcf861fbfbb26c40354cbf0c4662d1a910160405180910390a15b5061110288828c60200151613565565b6020808b015187820151604080518781526001600160601b038d16948101949094528415159084015284151560608401528b1515608084015290917faeb4b4786571e184246d39587f659abf0e26f41f6a3358692250382c0cdb47b79060a00160405180910390a3505050505050505b9392505050565b611181612f25565b61118a8161369a565b6111b257604051635428d44960e01b81526001600160a01b0382166004820152602401610b9a565b6000806000806111c186612bc9565b945094505093509350336001600160a01b0316826001600160a01b03161461120757604051636c51fda960e11b81526001600160a01b0383166004820152602401610b9a565b611210866114c7565b1561122e57604051631685ecdd60e31b815260040160405180910390fd5b6040805160c0810182526001815260208082018990526001600160a01b03851682840152606082018490526001600160601b038088166080840152861660a0830152915190916000916112839184910161588a565b604051602081830303815290604052905061129d88613705565b505060405163ce3f471960e01b81526001600160a01b0388169063ce3f4719906001600160601b038816906112d690859060040161594f565b6000604051808303818588803b1580156112ef57600080fd5b505af1158015611303573d6000803e3d6000fd5b50506002546001600160a01b03161580159350915061132c905057506001600160601b03861615155b156113b85760025460405163a9059cbb60e01b81526001600160a01b0389811660048301526001600160601b03891660248301526113b892169063a9059cbb906044015b6020604051808303816000875af115801561138f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113b39190615962565b6130f0565b600c805466ff00000000000019166601000000000000179055825160005b81811015611468578481815181106113f0576113f06157b3565b6020908102919091010151604051638ea9811760e01b81526001600160a01b038b8116600483015290911690638ea9811790602401600060405180830381600087803b15801561143f57600080fd5b505af1158015611453573d6000803e3d6000fd5b505050508061146190615808565b90506113d6565b50600c805466ff00000000000019169055604080516001600160a01b038a168152602081018b90527fd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be418791015b60405180910390a1505050505050505050565b60008181526005602052604081206002018054825b8181101561155e576000600460008584815481106114fc576114fc6157b3565b60009182526020808320909101546001600160a01b0316835282810193909352604091820181208982529092529020546001600160401b03600160481b90910416111561154e57506001949350505050565b61155781615808565b90506114dc565b506000949350505050565b611571612f25565b611579612e74565b6002546001600160a01b03166115a25760405163c1f0c0a160e01b815260040160405180910390fd5b600b546001600160601b03166115b98115156130f0565b600b80546bffffffffffffffffffffffff19169055600a80548291906000906115ec9084906001600160601b0316615821565b82546101009290920a6001600160601b0381810219909316918316021790915560025460405163a9059cbb60e01b81526001600160a01b0386811660048301529285166024820152610dbd935091169063a9059cbb90604401611370565b611652612e74565b61165b8161369a565b156116845760405163ac8a27ef60e01b81526001600160a01b0382166004820152602401610b9a565b601180546001810182556000919091527f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c680180546001600160a01b0319166001600160a01b0383169081179091556040519081527fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af016259060200160405180910390a150565b611710612e74565b6002546001600160a01b03161561173a57604051631688c53760e11b815260040160405180910390fd5b600280546001600160a01b039384166001600160a01b03199182161790915560038054929093169116179055565b6001546001600160a01b031633146117c25760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610b9a565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611821612e74565b6040805180820182526000916118509190859060029083908390808284376000920191909152506128dd915050565b6000818152600d602052604090205490915060ff161561188657604051634a0b8fa760e01b815260048101829052602401610b9a565b60408051808201825260018082526001600160401b0385811660208085018281526000888152600d835287812096518754925168ffffffffffffffffff1990931690151568ffffffffffffffff00191617610100929095169190910293909317909455600e805493840181559091527fbb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fd9091018490558251848152918201527f9b911b2c240bfbef3b6a8f7ed6ee321d1258bb2a3fe6becab52ac1cd3210afd39101610b5c565b611955612e74565b600a544790600160601b90046001600160601b031681811115611995576040516354ced18160e11b81526004810182905260248101839052604401610b9a565b81811015610da65760006119a982846157df565b90506000846001600160a01b03168260405160006040518083038185875af1925050503d80600081146119f8576040519150601f19603f3d011682016040523d82523d6000602084013e6119fd565b606091505b5050905080611a1f5760405163950b247960e01b815260040160405180910390fd5b604080516001600160a01b0387168152602081018490527f4aed7c8eed0496c8c19ea2681fcca25741c1602342e38b045d9f1e8e905d2e9c910160405180910390a15050505050565b611a70612f25565b600081815260056020526040902054611a91906001600160a01b03166138b7565b60008181526006602052604090208054600160601b90046001600160601b0316903490600c611ac0838561597f565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555034600a600c8282829054906101000a90046001600160601b0316611b08919061597f565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f7603b205d03651ee812f803fccde89f1012e545a9c99f0abfea9cedd0fd8e902823484611b5b919061599f565b604080519283526020830191909152015b60405180910390a25050565b6000611b82612f25565b60208083013560008181526005909252604090912054611baa906001600160a01b03166138b7565b336000908152600460209081526040808320848452808352928190208151606081018352905460ff811615158083526001600160401b036101008304811695840195909552600160481b9091049093169181019190915290611c28576040516379bfd40160e01b815260048101849052336024820152604401610b9a565b600c5461ffff16611c3f60608701604088016159b2565b61ffff161080611c62575060c8611c5c60608701604088016159b2565b61ffff16115b15611ca857611c7760608601604087016159b2565b600c5460405163539c34bb60e11b815261ffff92831660048201529116602482015260c86044820152606401610b9a565b600c5462010000900463ffffffff16611cc760808701606088016159cd565b63ffffffff161115611d1757611ce360808601606087016159cd565b600c54604051637aebf00f60e11b815263ffffffff9283166004820152620100009091049091166024820152604401610b9a565b6101f4611d2a60a08701608088016159cd565b63ffffffff161115611d7057611d4660a08601608087016159cd565b6040516311ce1afb60e21b815263ffffffff90911660048201526101f46024820152604401610b9a565b806020018051611d7f90615841565b6001600160401b03169052604081018051611d9990615841565b6001600160401b03908116909152602082810151604080518935818501819052338284015260608201899052929094166080808601919091528151808603909101815260a08501825280519084012060c085019290925260e08085018390528151808603909101815261010090940190528251929091019190912060009190955090506000611e3b611e36611e3160a08a018a6159ea565b6138de565b61395f565b9050854386611e5060808b0160608c016159cd565b611e6060a08c0160808d016159cd565b3386604051602001611e789796959493929190615a37565b60405160208183030381529060405280519060200120600f600088815260200190815260200160002081905550336001600160a01b03168588600001357feb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e89868c6040016020810190611eeb91906159b2565b8d6060016020810190611efe91906159cd565b8e6080016020810190611f1191906159cd565b89604051611f2496959493929190615a8e565b60405180910390a45050600092835260209182526040928390208151815493830151929094015168ffffffffffffffffff1990931693151568ffffffffffffffff001916939093176101006001600160401b03928316021770ffffffffffffffff0000000000000000001916600160481b91909216021790555b919050565b6000611fad612f25565b6007546001600160401b031633611fc56001436157df565b6040516bffffffffffffffffffffffff19606093841b81166020830152914060348201523090921b1660548201526001600160c01b031960c083901b16606882015260700160408051601f198184030181529190528051602090910120915061202f816001615acd565b6007805467ffffffffffffffff19166001600160401b03928316179055604080516000808252608082018352602080830182815283850183815260608086018581528a86526006855287862093518454935191516001600160601b039182166001600160c01b031990951694909417600160601b91909216021777ffffffffffffffffffffffffffffffffffffffffffffffff16600160c01b9290981691909102969096179055835194850184523385528481018281528585018481528884526005835294909220855181546001600160a01b03199081166001600160a01b03928316178355935160018301805490951691161790925592518051929493919261213f9260028501920190614fad565b5061214f915060089050846139d0565b5060405133815283907f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d9060200160405180910390a2505090565b612192612f25565b6002546001600160a01b031633146121bd576040516344b0e3c360e01b815260040160405180910390fd5b602081146121de57604051638129bbcd60e01b815260040160405180910390fd5b60006121ec828401846154a8565b600081815260056020526040902054909150612210906001600160a01b03166138b7565b600081815260066020526040812080546001600160601b031691869190612237838561597f565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600a60008282829054906101000a90046001600160601b031661227f919061597f565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a8287846122d2919061599f565b6040805192835260208301919091520160405180910390a2505050505050565b6122fa612e74565b60c861ffff8a1611156123345760405163539c34bb60e11b815261ffff8a1660048201819052602482015260c86044820152606401610b9a565b60008513612358576040516321ea67b360e11b815260048101869052602401610b9a565b8363ffffffff168363ffffffff161115612395576040516313c06e5960e11b815263ffffffff808516600483015285166024820152604401610b9a565b609b60ff831611156123c657604051631d66288d60e11b815260ff83166004820152609b6024820152604401610b9a565b609b60ff821611156123f757604051631d66288d60e11b815260ff82166004820152609b6024820152604401610b9a565b604080516101208101825261ffff8b1680825263ffffffff808c16602084018190526000848601528b8216606085018190528b8316608086018190528a841660a08701819052938a1660c0870181905260ff808b1660e08901819052908a16610100909801889052600c8054600160c01b90990260ff60c01b19600160b81b9093029290921661ffff60b81b19600160981b90940263ffffffff60981b19600160781b9099029890981676ffffffffffffffff00000000000000000000000000000019600160581b9096026effffffff000000000000000000000019670100000000000000909802979097166effffffffffffffffff000000000000196201000090990265ffffffffffff19909c16909a179a909a1796909616979097179390931791909116959095179290921793909316929092179190911790556010869055517f2c6b6b12413678366b05b145c5f00745bdd00e739131ab5de82484a50c9d78b6906114b4908b908b908b908b908b908b908b908b908b9061ffff99909916895263ffffffff97881660208a0152958716604089015293861660608801526080870192909252841660a086015290921660c084015260ff91821660e0840152166101008201526101200190565b6125ce612e74565b6000818152600560205260409020546001600160a01b03166125ef816138b7565b610dbd8282612f53565b6060600061260760086139dc565b905080841061262957604051631390f2a160e01b815260040160405180910390fd5b6000612635848661599f565b905081811180612643575083155b61264d578061264f565b815b9050600061265d86836157df565b9050806001600160401b038111156126775761267761517d565b6040519080825280602002602001820160405280156126a0578160200160208202803683370190505b50935060005b818110156126f0576126c36126bb888361599f565b6008906139e6565b8582815181106126d5576126d56157b3565b60209081029190910101526126e981615808565b90506126a6565b505050505b92915050565b612703612f25565b6000818152600560205260409020546001600160a01b0316612724816138b7565b6000828152600560205260409020600101546001600160a01b0316331461277d576000828152600560205260409081902060010154905163d084e97560e01b81526001600160a01b039091166004820152602401610b9a565b6000828152600560209081526040918290208054336001600160a01b03199182168117835560019092018054909116905582516001600160a01b03851681529182015283917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c93869101611b6c565b816127f481612ed0565b6127fc612f25565b6001600160a01b03821660009081526004602090815260408083208684529091529020805460ff161561282f5750505050565b6000848152600560205260409020600201805460631901612863576040516305a48e0f60e01b815260040160405180910390fd5b8154600160ff199091168117835581549081018255600082815260209081902090910180546001600160a01b0319166001600160a01b03871690811790915560405190815286917f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e191015b60405180910390a25050505050565b6000816040516020016128f09190615b10565b604051602081830303815290604052805190602001209050919050565b8161291781612ed0565b61291f612f25565b612928836114c7565b1561294657604051631685ecdd60e31b815260040160405180910390fd5b6001600160a01b038216600090815260046020908152604080832086845290915290205460ff1661299c576040516379bfd40160e01b8152600481018490526001600160a01b0383166024820152604401610b9a565b6000838152600560205260408120600201805490915b81811015612ab157846001600160a01b03168382815481106129d6576129d66157b3565b6000918252602090912001546001600160a01b031603612aa157826129fc6001846157df565b81548110612a0c57612a0c6157b3565b9060005260206000200160009054906101000a90046001600160a01b0316838281548110612a3c57612a3c6157b3565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555082805480612a7a57612a7a6157f2565b600082815260209020810160001990810180546001600160a01b0319169055019055612ab1565b612aaa81615808565b90506129b2565b506001600160a01b0384166000818152600460209081526040808320898452825291829020805460ff19169055905191825286917f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a791016128ce565b600e8181548110612b1d57600080fd5b600091825260209091200154905081565b81612b3881612ed0565b612b40612f25565b600083815260056020526040902060018101546001600160a01b03848116911614612bc3576001810180546001600160a01b0319166001600160a01b03851690811790915560408051338152602081019290925285917f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a191015b60405180910390a25b50505050565b600081815260056020526040812054819081906001600160a01b03166060612bf0826138b7565b600086815260066020908152604080832054600583529281902060020180548251818502810185019093528083526001600160601b0380861695600160601b810490911694600160c01b9091046001600160401b0316938893929091839190830182828015612c8857602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612c6a575b505050505090509450945094509450945091939590929450565b612caa612e74565b6002546001600160a01b0316612cd35760405163c1f0c0a160e01b815260040160405180910390fd5b6002546040516370a0823160e01b81523060048201526000916001600160a01b0316906370a0823190602401602060405180830381865afa158015612d1c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d409190615b1e565b600a549091506001600160601b031681811115612d7a576040516354ced18160e11b81526004810182905260248101839052604401610b9a565b81811015610da6576000612d8e82846157df565b60025460405163a9059cbb60e01b81526001600160a01b0387811660048301526024820184905292935091169063a9059cbb906044016020604051808303816000875af1158015612de3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e079190615962565b612e2457604051631f01ff1360e21b815260040160405180910390fd5b604080516001600160a01b0386168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b4366009101610d55565b612e6b612e74565b610ba3816139f2565b6000546001600160a01b03163314612ece5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610b9a565b565b6000818152600560205260409020546001600160a01b0316612ef1816138b7565b336001600160a01b03821614610dbd57604051636c51fda960e11b81526001600160a01b0382166004820152602401610b9a565b600c546601000000000000900460ff1615612ece5760405163769dd35360e11b815260040160405180910390fd5b600080612f5f84613705565b60025491935091506001600160a01b031615801590612f8657506001600160601b03821615155b15612fce5760025460405163a9059cbb60e01b81526001600160a01b0385811660048301526001600160601b0385166024830152612fce92169063a9059cbb90604401611370565b612fe183826001600160601b031661310e565b604080516001600160a01b03851681526001600160601b03808516602083015283169181019190915284907f8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c490606001612bba565b60038260ff161061305e57604051621c300f60e81b815260ff83166004820152602401610b9a565b60ff81161580613071575060648160ff16115b156130935760405162d4503560e51b815260ff82166004820152602401610b9a565b6012805460ff84811661ffff199092168217610100918516918202179092556040805191825260208201929092527f8e63dc2f2e669ce73bebd2580bb9dd9a5d17fa2d046ac02057d8349fc0b0c2f3910160405180910390a15050565b80610ba357604051631e9acf1760e31b815260040160405180910390fd5b6000826001600160a01b03168260405160006040518083038185875af1925050503d806000811461315b576040519150601f19603f3d011682016040523d82523d6000602084013e613160565b606091505b5050905080610da65760405163950b247960e01b815260040160405180910390fd5b6040805160a081018252600060608201818152608083018290528252602082018190529181019190915260006131bb84600001516128dd565b6000818152600d602090815260409182902082518084019093525460ff811615158084526101009091046001600160401b0316918301919091529192509061321957604051631dfd6e1360e21b815260048101839052602401610b9a565b600082866080015160405160200161323b929190918252602082015260400190565b60408051601f1981840301815291815281516020928301206000818152600f909352908220549092509081900361328557604051631b44092560e11b815260040160405180910390fd5b85516020808801516040808a015160608b015160808c015160a08d015193516132b4978a979096959101615b37565b6040516020818303038152906040528051906020012081146132e95760405163354a450b60e21b815260040160405180910390fd5b85516001600160401b031640806133bc578651604051631d2827a760e31b81526001600160401b0390911660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063e9413d3890602401602060405180830381865afa15801561336a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061338e9190615b1e565b9050806133bc57865160405163175dadad60e01b81526001600160401b039091166004820152602401610b9a565b60008860800151826040516020016133de929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c905060006134058a83613a9b565b604080516060810182529788526020880196909652948601949094525092979650505050505050565b6000816001600160401b03163a111561348157821561345757506001600160401b0381166126f5565b60405163435e532d60e11b81523a60048201526001600160401b0383166024820152604401610b9a565b503a92915050565b6000806000631fe543e360e01b86856040516024016134a9929190615b8a565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252600c805466ff000000000000191666010000000000001790559086015160808701519192506135139163ffffffff9091169083613b06565b600c805466ff000000000000191690559695505050505050565b600080831561354c57613541868685613b52565b60009150915061355c565b613557868685613c68565b915091505b94509492505050565b600081815260066020526040902082156136155780546001600160601b03600160601b90910481169061359c9086168210156130f0565b6135a68582615821565b82546bffffffffffffffffffffffff60601b1916600160601b6001600160601b039283168102919091178455600b805488939192600c926135eb92869290041661597f565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555050612bc3565b80546001600160601b03908116906136319086168210156130f0565b61363b8582615821565b82546bffffffffffffffffffffffff19166001600160601b03918216178355600b8054879260009161366f9185911661597f565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505050505050565b601154600090815b818110156136fb57836001600160a01b0316601182815481106136c7576136c76157b3565b6000918252602090912001546001600160a01b0316036136eb575060019392505050565b6136f481615808565b90506136a2565b5060009392505050565b60008181526005602090815260408083206006909252822054600290910180546001600160601b0380841694600160601b90940416925b818110156137b1576004600084838154811061375a5761375a6157b3565b60009182526020808320909101546001600160a01b0316835282810193909352604091820181208982529092529020805470ffffffffffffffffffffffffffffffffff191690556137aa81615808565b905061373c565b50600085815260056020526040812080546001600160a01b031990811682556001820180549091169055906137e96002830182615012565b5050600085815260066020526040812055613805600886613e5f565b506001600160601b0384161561385857600a80548591906000906138339084906001600160601b0316615821565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b6001600160601b038316156138b05782600a600c8282829054906101000a90046001600160601b031661388b9190615821565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b5050915091565b6001600160a01b038116610ba357604051630fb532db60e11b815260040160405180910390fd5b604080516020810190915260008152600082900361390b57506040805160208101909152600081526126f5565b63125fa26760e31b61391d8385615bab565b6001600160e01b0319161461394557604051632923fee760e11b815260040160405180910390fd5b6139528260048186615bdb565b8101906111729190615c05565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa8260405160240161399891511515815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b03199093169290921790915292915050565b60006111728383613e6b565b60006126f5825490565b60006111728383613eba565b336001600160a01b03821603613a4a5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610b9a565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000613acf8360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151613ee4565b60038360200151604051602001613ae7929190615c50565b60408051601f1981840301815291905280516020909101209392505050565b60005a611388811015613b1857600080fd5b611388810390508460408204820311613b3057600080fd5b50823b613b3c57600080fd5b60008083516020850160008789f1949350505050565b600080613b6060003661410f565b905060005a600c54613b80908890600160581b900463ffffffff1661599f565b613b8a91906157df565b613b949086615c64565b600c54909150600090613bb990600160781b900463ffffffff1664e8d4a51000615c64565b90508215613bf5576040518381527f56296f7beae05a0db815737fdb4cd298897b1e517614d62468081531ae14d0999060200160405180910390a15b8415613c3f57600c548190606490600160b81b900460ff16613c17858761599f565b613c219190615c64565b613c2b9190615c91565b613c35919061599f565b9350505050611172565b600c548190606490613c5b90600160b81b900460ff1682615ca5565b60ff16613c17858761599f565b600080600080613c7661411b565b9150915060008213613c9e576040516321ea67b360e11b815260048101839052602401610b9a565b6000613cab60003661410f565b9050600083825a600c54613ccd908d90600160581b900463ffffffff1661599f565b613cd791906157df565b613ce1908b615c64565b613ceb919061599f565b613cfd90670de0b6b3a7640000615c64565b613d079190615c91565b600c54909150600090613d309063ffffffff600160981b8204811691600160781b900416615cbe565b613d459063ffffffff1664e8d4a51000615c64565b9050600085613d5c83670de0b6b3a7640000615c64565b613d669190615c91565b90508315613da2576040518481527f56296f7beae05a0db815737fdb4cd298897b1e517614d62468081531ae14d0999060200160405180910390a15b60008915613de157600c548290606490613dc690600160c01b900460ff1687615c64565b613dd09190615c91565b613dda919061599f565b9050613e21565b600c548290606490613dfd90600160c01b900460ff1682615ca5565b613e0a9060ff1687615c64565b613e149190615c91565b613e1e919061599f565b90505b6b033b2e3c9fd0803ce8000000811115613e4e5760405163e80fa38160e01b815260040160405180910390fd5b9b949a509398505050505050505050565b600061117283836141e6565b6000818152600183016020526040812054613eb2575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556126f5565b5060006126f5565b6000826000018281548110613ed157613ed16157b3565b9060005260206000200154905092915050565b613eed896142e0565b613f395760405162461bcd60e51b815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e2063757276650000000000006044820152606401610b9a565b613f42886142e0565b613f8e5760405162461bcd60e51b815260206004820152601560248201527f67616d6d61206973206e6f74206f6e20637572766500000000000000000000006044820152606401610b9a565b613f97836142e0565b613fe35760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610b9a565b613fec826142e0565b6140385760405162461bcd60e51b815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e206375727665000000006044820152606401610b9a565b614044878a88876143b9565b6140905760405162461bcd60e51b815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e657373000000000000006044820152606401610b9a565b600061409c8a876144dc565b905060006140af898b878b868989614540565b905060006140c0838d8d8a8661466c565b9050808a146141015760405162461bcd60e51b815260206004820152600d60248201526c34b73b30b634b210383937b7b360991b6044820152606401610b9a565b505050505050505050505050565b600061117283836146ac565b600c5460035460408051633fabe5a360e21b81529051600093849367010000000000000090910463ffffffff169284926001600160a01b039092169163feaf968c9160048082019260a0929091908290030181865afa158015614182573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906141a69190615cf5565b50919650909250505063ffffffff8216158015906141d257506141c981426157df565b8263ffffffff16105b925082156141e05760105493505b50509091565b600081815260018301602052604081205480156142cf57600061420a6001836157df565b855490915060009061421e906001906157df565b905081811461428357600086600001828154811061423e5761423e6157b3565b9060005260206000200154905080876000018481548110614261576142616157b3565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080614294576142946157f2565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506126f5565b60009150506126f5565b5092915050565b80516000906401000003d019116143395760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420782d6f7264696e61746500000000000000000000000000006044820152606401610b9a565b60208201516401000003d019116143925760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420792d6f7264696e61746500000000000000000000000000006044820152606401610b9a565b60208201516401000003d0199080096143b28360005b6020020151614778565b1492915050565b60006001600160a01b0382166143ff5760405162461bcd60e51b815260206004820152600b60248201526a626164207769746e65737360a81b6044820152606401610b9a565b60208401516000906001161561441657601c614419565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe19918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa1580156144b4573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b6144e4615030565b614511600184846040516020016144fd93929190615d45565b60405160208183030381529060405261479c565b90505b61451d816142e0565b6126f557805160408051602081019290925261453991016144fd565b9050614514565b614548615030565b825186516401000003d01991829006919006036145a75760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610b9a565b6145b28789886147e9565b6145fe5760405162461bcd60e51b815260206004820152601660248201527f4669727374206d756c20636865636b206661696c6564000000000000000000006044820152606401610b9a565b6146098486856147e9565b6146555760405162461bcd60e51b815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c65640000000000000000006044820152606401610b9a565b614660868484614914565b98975050505050505050565b60006002868686858760405160200161468a96959493929190615d66565b60408051601f1981840301815291905280516020909101209695505050505050565b60125460009060ff1661476f57600f602160991b016001600160a01b03166349948e0e8484604051806080016040528060478152602001615f40604791396040516020016146fc93929190615dc5565b6040516020818303038152906040526040518263ffffffff1660e01b8152600401614727919061594f565b602060405180830381865afa158015614744573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906147689190615b1e565b90506126f5565b611172826149db565b6000806401000003d01980848509840990506401000003d019600782089392505050565b6147a4615030565b6147ad82614ac3565b81526147c26147bd8260006143a8565b614afe565b6020820181905260029006600103611f9e576020810180516401000003d019039052919050565b6000826000036148295760405162461bcd60e51b815260206004820152600b60248201526a3d32b9379039b1b0b630b960a91b6044820152606401610b9a565b8351602085015160009061483f90600290615dec565b1561484b57601c61484e565b601b5b9050600070014551231950b75fc4402da1732fc9bebe198387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa1580156148c0573d6000803e3d6000fd5b5050506020604051035190506000866040516020016148df9190615e00565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b61491c615030565b83516020808601518551918601516000938493849361493d93909190614b1e565b919450925090506401000003d01985820960011461499d5760405162461bcd60e51b815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a000000000000006044820152606401610b9a565b60405180604001604052806401000003d019806149bc576149bc615c7b565b87860981526020016401000003d0198785099052979650505050505050565b60125460009060ff166000198101614a21576064614a026149fd60478661599f565b614bfe565b601254614a179190610100900460ff16615c64565b6111729190615c91565b60011960ff821601614aa6576064600f602160991b0163f1c7a58b614a4760478761599f565b6040518263ffffffff1660e01b8152600401614a6591815260200190565b602060405180830381865afa158015614a82573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614a029190615b1e565b604051621c300f60e81b815260ff82166004820152602401610b9a565b805160208201205b6401000003d0198110611f9e57604080516020808201939093528151808203840181529082019091528051910120614acb565b60006126f5826002614b176401000003d019600161599f565b901c614e9b565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000614b5e83838585614f40565b9098509050614b6f88828e88614f64565b9098509050614b8088828c87614f64565b90985090506000614b938d878b85614f64565b9098509050614ba488828686614f40565b9098509050614bb588828e89614f64565b9098509050818114614bea576401000003d019818a0998506401000003d01982890997506401000003d0198183099650614bee565b8196505b5050505050509450945094915050565b600080614c0c60448461599f565b614c17906010615c64565b90506000600f602160991b016001600160a01b031663519b4bd36040518163ffffffff1660e01b8152600401602060405180830381865afa158015614c60573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614c849190615b1e565b600f602160991b016001600160a01b031663c59859186040518163ffffffff1660e01b8152600401602060405180830381865afa158015614cc9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614ced9190615e12565b614cf8906010615e2f565b63ffffffff16614d089190615c64565b90506000600f602160991b016001600160a01b031663f82061406040518163ffffffff1660e01b8152600401602060405180830381865afa158015614d51573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614d759190615b1e565b600f602160991b016001600160a01b03166368d5dca66040518163ffffffff1660e01b8152600401602060405180830381865afa158015614dba573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614dde9190615e12565b63ffffffff16614dee9190615c64565b90506000614dfc828461599f565b614e069085615c64565b9050600f602160991b016001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa158015614e4d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614e719190615b1e565b614e7c90600a615f33565b614e87906010615c64565b614e919082615c91565b9695505050505050565b600080614ea661504e565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152614ed861506c565b60208160c0846005600019fa925082600003614f365760405162461bcd60e51b815260206004820152601260248201527f6269674d6f64457870206661696c7572652100000000000000000000000000006044820152606401610b9a565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215615002579160200282015b8281111561500257825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190614fcd565b5061500e92915061508a565b5090565b5080546000825590600052602060002090810190610ba3919061508a565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b8082111561500e576000815560010161508b565b6001600160a01b0381168114610ba357600080fd5b8035611f9e8161509f565b6000602082840312156150d157600080fd5b81356111728161509f565b80604081018310156126f557600080fd5b6000604082840312156150ff57600080fd5b61117283836150dc565b6000806040838503121561511c57600080fd5b82359150602083013561512e8161509f565b809150509250929050565b803560ff81168114611f9e57600080fd5b6000806040838503121561515d57600080fd5b61516683615139565b915061517460208401615139565b90509250929050565b634e487b7160e01b600052604160045260246000fd5b60405161012081016001600160401b03811182821017156151b6576151b661517d565b60405290565b600082601f8301126151cd57600080fd5b604051604081018181106001600160401b03821117156151ef576151ef61517d565b806040525080604084018581111561520657600080fd5b845b81811015615220578035835260209283019201615208565b509195945050505050565b80356001600160401b0381168114611f9e57600080fd5b63ffffffff81168114610ba357600080fd5b600082601f83011261526557600080fd5b81356001600160401b038082111561527f5761527f61517d565b604051601f8301601f19908116603f011681019082821181831017156152a7576152a761517d565b816040528381528660208588010111156152c057600080fd5b836020870160208301376000602085830101528094505050505092915050565b600060c082840312156152f257600080fd5b60405160c081016001600160401b0382821081831117156153155761531561517d565b816040528293506153258561522b565b8352602085013560208401526040850135915061534182615242565b8160408401526060850135915061535782615242565b8160608401526080850135915061536d8261509f565b81608084015260a085013591508082111561538757600080fd5b5061539485828601615254565b60a0830152505092915050565b8015158114610ba357600080fd5b8035611f9e816153a1565b60008060008385036101e08112156153d157600080fd5b6101a0808212156153e157600080fd5b6153e9615193565b91506153f587876151bc565b825261540487604088016151bc565b60208301526080860135604083015260a0860135606083015260c0860135608083015261543360e087016150b4565b60a0830152610100615447888289016151bc565b60c084015261545a8861014089016151bc565b60e0840152610180870135908301529093508401356001600160401b0381111561548357600080fd5b61548f868287016152e0565b92505061549f6101c085016153af565b90509250925092565b6000602082840312156154ba57600080fd5b5035919050565b600080604083850312156154d457600080fd5b82356154df8161509f565b9150602083013561512e8161509f565b6000806060838503121561550257600080fd5b61550c84846150dc565b91506151746040840161522b565b60006020828403121561552c57600080fd5b81356001600160401b0381111561554257600080fd5b820160c0818503121561117257600080fd5b6000806000806060858703121561556a57600080fd5b84356155758161509f565b93506020850135925060408501356001600160401b038082111561559857600080fd5b818701915087601f8301126155ac57600080fd5b8135818111156155bb57600080fd5b8860208285010111156155cd57600080fd5b95989497505060200194505050565b803561ffff81168114611f9e57600080fd5b60008060008060008060008060006101208a8c03121561560d57600080fd5b6156168a6155dc565b985060208a013561562681615242565b975060408a013561563681615242565b965060608a013561564681615242565b955060808a0135945060a08a013561565d81615242565b935060c08a013561566d81615242565b925061567b60e08b01615139565b915061568a6101008b01615139565b90509295985092959850929598565b600080604083850312156156ac57600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b838110156156eb578151875295820195908201906001016156cf565b509495945050505050565b60208152600061117260208301846156bb565b60006040828403121561571b57600080fd5b61117283836151bc565b600081518084526020808501945080840160005b838110156156eb5781516001600160a01b031687529582019590820190600101615739565b60006001600160601b0380881683528087166020840152506001600160401b03851660408301526001600160a01b038416606083015260a060808301526157a860a0830184615725565b979650505050505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b818103818111156126f5576126f56157c9565b634e487b7160e01b600052603160045260246000fd5b60006001820161581a5761581a6157c9565b5060010190565b6001600160601b038281168282160390808211156142d9576142d96157c9565b60006001600160401b0380831681810361585d5761585d6157c9565b6001019392505050565b60006001600160401b03821680615880576158806157c9565b6000190192915050565b6020815260ff8251166020820152602082015160408201526001600160a01b0360408301511660608201526000606083015160c060808401526158d060e0840182615725565b905060808401516001600160601b0380821660a08601528060a08701511660c086015250508091505092915050565b60005b8381101561591a578181015183820152602001615902565b50506000910152565b6000815180845261593b8160208601602086016158ff565b601f01601f19169290920160200192915050565b6020815260006111726020830184615923565b60006020828403121561597457600080fd5b8151611172816153a1565b6001600160601b038181168382160190808211156142d9576142d96157c9565b808201808211156126f5576126f56157c9565b6000602082840312156159c457600080fd5b611172826155dc565b6000602082840312156159df57600080fd5b813561117281615242565b6000808335601e19843603018112615a0157600080fd5b8301803591506001600160401b03821115615a1b57600080fd5b602001915036819003821315615a3057600080fd5b9250929050565b878152866020820152856040820152600063ffffffff80871660608401528086166080840152506001600160a01b03841660a083015260e060c0830152615a8160e0830184615923565b9998505050505050505050565b86815285602082015261ffff85166040820152600063ffffffff808616606084015280851660808401525060c060a083015261466060c0830184615923565b6001600160401b038181168382160190808211156142d9576142d96157c9565b8060005b6002811015612bc3578151845260209384019390910190600101615af1565b604081016126f58284615aed565b600060208284031215615b3057600080fd5b5051919050565b8781526001600160401b0387166020820152856040820152600063ffffffff80871660608401528086166080840152506001600160a01b03841660a083015260e060c0830152615a8160e0830184615923565b828152604060208201526000615ba360408301846156bb565b949350505050565b6001600160e01b03198135818116916004851015615bd35780818660040360031b1b83161692505b505092915050565b60008085851115615beb57600080fd5b83861115615bf857600080fd5b5050820193919092039150565b600060208284031215615c1757600080fd5b604051602081018181106001600160401b0382111715615c3957615c3961517d565b6040528235615c47816153a1565b81529392505050565b828152606081016111726020830184615aed565b80820281158282048414176126f5576126f56157c9565b634e487b7160e01b600052601260045260246000fd5b600082615ca057615ca0615c7b565b500490565b60ff81811683821601908111156126f5576126f56157c9565b63ffffffff8281168282160390808211156142d9576142d96157c9565b805169ffffffffffffffffffff81168114611f9e57600080fd5b600080600080600060a08688031215615d0d57600080fd5b615d1686615cdb565b9450602086015193506040860151925060608601519150615d3960808701615cdb565b90509295509295909350565b838152615d556020820184615aed565b606081019190915260800192915050565b868152615d766020820187615aed565b615d836060820186615aed565b615d9060a0820185615aed565b615d9d60e0820184615aed565b60609190911b6bffffffffffffffffffffffff19166101208201526101340195945050505050565b828482376000838201600081528351615de28183602088016158ff565b0195945050505050565b600082615dfb57615dfb615c7b565b500690565b615e0a8183615aed565b604001919050565b600060208284031215615e2457600080fd5b815161117281615242565b63ffffffff818116838216028082169190828114615bd357615bd36157c9565b600181815b80851115615e8a578160001904821115615e7057615e706157c9565b80851615615e7d57918102915b93841c9390800290615e54565b509250929050565b600082615ea1575060016126f5565b81615eae575060006126f5565b8160018114615ec45760028114615ece57615eea565b60019150506126f5565b60ff841115615edf57615edf6157c9565b50506001821b6126f5565b5060208310610133831016604e8410600b8410161715615f0d575081810a6126f5565b615f178383615e4f565b8060001904821115615f2b57615f2b6157c9565b029392505050565b60006111728383615e9256feffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa164736f6c6343000813000a",
}

var VRFCoordinatorV25OptimismABI = VRFCoordinatorV25OptimismMetaData.ABI

var VRFCoordinatorV25OptimismBin = VRFCoordinatorV25OptimismMetaData.Bin

func DeployVRFCoordinatorV25Optimism(auth *bind.TransactOpts, backend bind.ContractBackend, blockhashStore common.Address) (common.Address, *types.Transaction, *VRFCoordinatorV25Optimism, error) {
	parsed, err := VRFCoordinatorV25OptimismMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorV25OptimismBin), backend, blockhashStore)
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) SL1FeeCalculationMode(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "s_l1FeeCalculationMode")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) SL1FeeCalculationMode() (uint8, error) {
	return _VRFCoordinatorV25Optimism.Contract.SL1FeeCalculationMode(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) SL1FeeCalculationMode() (uint8, error) {
	return _VRFCoordinatorV25Optimism.Contract.SL1FeeCalculationMode(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCaller) SL1FeeCoefficient(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Optimism.contract.Call(opts, &out, "s_l1FeeCoefficient")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) SL1FeeCoefficient() (uint8, error) {
	return _VRFCoordinatorV25Optimism.Contract.SL1FeeCoefficient(&_VRFCoordinatorV25Optimism.CallOpts)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismCallerSession) SL1FeeCoefficient() (uint8, error) {
	return _VRFCoordinatorV25Optimism.Contract.SL1FeeCoefficient(&_VRFCoordinatorV25Optimism.CallOpts)
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactor) SetL1FeeCalculation(opts *bind.TransactOpts, mode uint8, coefficient uint8) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.contract.Transact(opts, "setL1FeeCalculation", mode, coefficient)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismSession) SetL1FeeCalculation(mode uint8, coefficient uint8) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.SetL1FeeCalculation(&_VRFCoordinatorV25Optimism.TransactOpts, mode, coefficient)
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismTransactorSession) SetL1FeeCalculation(mode uint8, coefficient uint8) (*types.Transaction, error) {
	return _VRFCoordinatorV25Optimism.Contract.SetL1FeeCalculation(&_VRFCoordinatorV25Optimism.TransactOpts, mode, coefficient)
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

type VRFCoordinatorV25OptimismL1FeeCalculationSetIterator struct {
	Event *VRFCoordinatorV25OptimismL1FeeCalculationSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismL1FeeCalculationSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismL1FeeCalculationSet)
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
		it.Event = new(VRFCoordinatorV25OptimismL1FeeCalculationSet)
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

func (it *VRFCoordinatorV25OptimismL1FeeCalculationSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismL1FeeCalculationSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismL1FeeCalculationSet struct {
	Mode        uint8
	Coefficient uint8
	Raw         types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterL1FeeCalculationSet(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismL1FeeCalculationSetIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "L1FeeCalculationSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismL1FeeCalculationSetIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "L1FeeCalculationSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchL1FeeCalculationSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismL1FeeCalculationSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "L1FeeCalculationSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismL1FeeCalculationSet)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "L1FeeCalculationSet", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseL1FeeCalculationSet(log types.Log) (*VRFCoordinatorV25OptimismL1FeeCalculationSet, error) {
	event := new(VRFCoordinatorV25OptimismL1FeeCalculationSet)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "L1FeeCalculationSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OptimismL1GasFeeIterator struct {
	Event *VRFCoordinatorV25OptimismL1GasFee

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OptimismL1GasFeeIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OptimismL1GasFee)
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
		it.Event = new(VRFCoordinatorV25OptimismL1GasFee)
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

func (it *VRFCoordinatorV25OptimismL1GasFeeIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OptimismL1GasFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OptimismL1GasFee struct {
	Fee *big.Int
	Raw types.Log
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) FilterL1GasFee(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismL1GasFeeIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.FilterLogs(opts, "L1GasFee")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OptimismL1GasFeeIterator{contract: _VRFCoordinatorV25Optimism.contract, event: "L1GasFee", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) WatchL1GasFee(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismL1GasFee) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Optimism.contract.WatchLogs(opts, "L1GasFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OptimismL1GasFee)
				if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "L1GasFee", log); err != nil {
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

func (_VRFCoordinatorV25Optimism *VRFCoordinatorV25OptimismFilterer) ParseL1GasFee(log types.Log) (*VRFCoordinatorV25OptimismL1GasFee, error) {
	event := new(VRFCoordinatorV25OptimismL1GasFee)
	if err := _VRFCoordinatorV25Optimism.contract.UnpackLog(event, "L1GasFee", log); err != nil {
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
	case _VRFCoordinatorV25Optimism.abi.Events["L1FeeCalculationSet"].ID:
		return _VRFCoordinatorV25Optimism.ParseL1FeeCalculationSet(log)
	case _VRFCoordinatorV25Optimism.abi.Events["L1GasFee"].ID:
		return _VRFCoordinatorV25Optimism.ParseL1GasFee(log)
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

func (VRFCoordinatorV25OptimismL1FeeCalculationSet) Topic() common.Hash {
	return common.HexToHash("0x8e63dc2f2e669ce73bebd2580bb9dd9a5d17fa2d046ac02057d8349fc0b0c2f3")
}

func (VRFCoordinatorV25OptimismL1GasFee) Topic() common.Hash {
	return common.HexToHash("0x56296f7beae05a0db815737fdb4cd298897b1e517614d62468081531ae14d099")
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

	SL1FeeCalculationMode(opts *bind.CallOpts) (uint8, error)

	SL1FeeCoefficient(opts *bind.CallOpts) (uint8, error)

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

	SetL1FeeCalculation(opts *bind.TransactOpts, mode uint8, coefficient uint8) (*types.Transaction, error)

	SetLINKAndLINKNativeFeed(opts *bind.TransactOpts, link common.Address, linkNativeFeed common.Address) (*types.Transaction, error)

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

	FilterL1FeeCalculationSet(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismL1FeeCalculationSetIterator, error)

	WatchL1FeeCalculationSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismL1FeeCalculationSet) (event.Subscription, error)

	ParseL1FeeCalculationSet(log types.Log) (*VRFCoordinatorV25OptimismL1FeeCalculationSet, error)

	FilterL1GasFee(opts *bind.FilterOpts) (*VRFCoordinatorV25OptimismL1GasFeeIterator, error)

	WatchL1GasFee(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OptimismL1GasFee) (event.Subscription, error)

	ParseL1GasFee(log types.Log) (*VRFCoordinatorV25OptimismL1GasFee, error)

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
