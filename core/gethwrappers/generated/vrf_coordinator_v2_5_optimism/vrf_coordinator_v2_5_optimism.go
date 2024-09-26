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
	Bin: "0x60a06040526012805461ffff19166164001790553480156200002057600080fd5b50604051620060253803806200602583398101604081905262000043916200018f565b8033806000816200009b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000ce57620000ce81620000e4565b5050506001600160a01b031660805250620001c1565b336001600160a01b038216036200013e5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000092565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208284031215620001a257600080fd5b81516001600160a01b0381168114620001ba57600080fd5b9392505050565b608051615e41620001e46000396000818161066f015261333e0152615e416000f3fe6080604052600436106102dd5760003560e01c80638402595e1161017f578063aefb212f116100e1578063da2f26101161008a578063e72f6e3011610064578063e72f6e30146109c0578063ee9d2d38146109e0578063f2fde38b14610a0d57600080fd5b8063da2f261014610910578063dac83d291461096f578063dc311dd31461098f57600080fd5b8063caf70c4a116100bb578063caf70c4a146108b0578063cb631797146108d0578063d98e620e146108f057600080fd5b8063aefb212f14610843578063b2a7cac514610870578063bec4c08c1461089057600080fd5b80639b1c385e11610143578063a4c0ed361161011d578063a4c0ed36146107e3578063a63e0bfb14610803578063aa433aff1461082357600080fd5b80639b1c385e146107765780639d40a6fd14610796578063a21a23e4146107ce57600080fd5b80638402595e146106e657806386fe91c7146107065780638da5cb5b1461072657806390bd5c741461074457806395b55cfc1461076357600080fd5b8063405b84fa1161024357806364d51a2a116101ec57806372e9d565116101c657806372e9d5651461069157806379ba5097146106b15780637a5a2aef146106c657600080fd5b806364d51a2a14610628578063659827441461063d578063689c45171461065d57600080fd5b806341af6c871161021d57806341af6c87146105b857806351cff8d9146105e85780635d06b4ab1461060857600080fd5b8063405b84fa1461054157806340d6bb821461056157806340e3290f1461058c57600080fd5b806314530741116102a55780631b6b6d231161027f5780631b6b6d23146104c95780632f622e6b14610501578063301f42e91461052157600080fd5b8063145307411461044257806315c48b841461046257806318e3dd271461048a57600080fd5b806304104edb146102e2578063043bd6ae14610304578063088070f51461032d57806308821d58146104025780630ae0954014610422575b600080fd5b3480156102ee57600080fd5b506103026102fd366004615198565b610a2d565b005b34801561031057600080fd5b5061031a60105481565b6040519081526020015b60405180910390f35b34801561033957600080fd5b50600c546103a59061ffff81169063ffffffff62010000820481169160ff660100000000000082048116926701000000000000008304811692600160581b8104821692600160781b8204831692600160981b83041691600160b81b8104821691600160c01b9091041689565b6040805161ffff909a168a5263ffffffff98891660208b01529615159689019690965293861660608801529185166080870152841660a08601529290921660c084015260ff91821660e08401521661010082015261012001610324565b34801561040e57600080fd5b5061030261041d3660046151c6565b610ba6565b34801561042e57600080fd5b5061030261043d3660046151e2565b610d3f565b34801561044e57600080fd5b5061030261045d366004615223565b610d87565b34801561046e57600080fd5b5061047760c881565b60405161ffff9091168152602001610324565b34801561049657600080fd5b50600a546104b190600160601b90046001600160601b031681565b6040516001600160601b039091168152602001610324565b3480156104d557600080fd5b506002546104e9906001600160a01b031681565b6040516001600160a01b039091168152602001610324565b34801561050d57600080fd5b5061030261051c366004615198565b610d9d565b34801561052d57600080fd5b506104b161053c36600461527c565b610e3f565b34801561054d57600080fd5b5061030261055c3660046151e2565b611190565b34801561056d57600080fd5b506105776101f481565b60405163ffffffff9091168152602001610324565b34801561059857600080fd5b506012546105a69060ff1681565b60405160ff9091168152602001610324565b3480156105c457600080fd5b506105d86105d33660046152e8565b6114de565b6040519015158152602001610324565b3480156105f457600080fd5b50610302610603366004615198565b611580565b34801561061457600080fd5b50610302610623366004615198565b611661565b34801561063457600080fd5b50610477606481565b34801561064957600080fd5b50610302610658366004615301565b61171f565b34801561066957600080fd5b506104e97f000000000000000000000000000000000000000000000000000000000000000081565b34801561069d57600080fd5b506003546104e9906001600160a01b031681565b3480156106bd57600080fd5b5061030261177f565b3480156106d257600080fd5b506103026106e1366004615346565b611830565b3480156106f257600080fd5b50610302610701366004615198565b611940565b34801561071257600080fd5b50600a546104b1906001600160601b031681565b34801561073257600080fd5b506000546001600160a01b03166104e9565b34801561075057600080fd5b506012546105a690610100900460ff1681565b6103026107713660046152e8565b611a5b565b34801561078257600080fd5b5061031a610791366004615371565b611b6b565b3480156107a257600080fd5b506007546107b6906001600160401b031681565b6040516001600160401b039091168152602001610324565b3480156107da57600080fd5b5061031a611f96565b3480156107ef57600080fd5b506103026107fe3660046153ad565b61217d565b34801561080f57600080fd5b5061030261081e366004615459565b6122e5565b34801561082f57600080fd5b5061030261083e3660046152e8565b6125b9565b34801561084f57600080fd5b5061086361085e366004615504565b6125ec565b6040516103249190615561565b34801561087c57600080fd5b5061030261088b3660046152e8565b6126ee565b34801561089c57600080fd5b506103026108ab3660046151e2565b6127dd565b3480156108bc57600080fd5b5061031a6108cb3660046151c6565b6128d0565b3480156108dc57600080fd5b506103026108eb3660046151e2565b612900565b3480156108fc57600080fd5b5061031a61090b3660046152e8565b612b00565b34801561091c57600080fd5b5061095061092b3660046152e8565b600d6020526000908152604090205460ff81169061010090046001600160401b031682565b6040805192151583526001600160401b03909116602083015201610324565b34801561097b57600080fd5b5061030261098a3660046151e2565b612b21565b34801561099b57600080fd5b506109af6109aa3660046152e8565b612bbc565b6040516103249594939291906155ad565b3480156109cc57600080fd5b506103026109db366004615198565b612c95565b3480156109ec57600080fd5b5061031a6109fb3660046152e8565b600f6020526000908152604090205481565b348015610a1957600080fd5b50610302610a28366004615198565b612e56565b610a35612e67565b60115460005b81811015610b7957826001600160a01b031660118281548110610a6057610a60615602565b6000918252602090912001546001600160a01b031603610b69576011610a8760018461562e565b81548110610a9757610a97615602565b600091825260209091200154601180546001600160a01b039092169183908110610ac357610ac3615602565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b031602179055506011805480610b0257610b02615641565b6000828152602090819020600019908301810180546001600160a01b03191690559091019091556040516001600160a01b03851681527ff80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af3791015b60405180910390a1505050565b610b7281615657565b9050610a3b565b50604051635428d44960e01b81526001600160a01b03831660048201526024015b60405180910390fd5b50565b610bae612e67565b6000610bb9826128d0565b6000818152600d602090815260409182902082518084019093525460ff811615158084526101009091046001600160401b03169183019190915291925090610c1757604051631dfd6e1360e21b815260048101839052602401610b9a565b6000828152600d60205260408120805468ffffffffffffffffff19169055600e54905b81811015610ce95783600e8281548110610c5657610c56615602565b906000526020600020015403610cd957600e610c7360018461562e565b81548110610c8357610c83615602565b9060005260206000200154600e8281548110610ca157610ca1615602565b600091825260209091200155600e805480610cbe57610cbe615641565b60019003818190600052602060002001600090559055610ce9565b610ce281615657565b9050610c3a565b507f9b6868e0eb737bcd72205360baa6bfd0ba4e4819a33ade2db384e8a8025639a5838360200151604051610d319291909182526001600160401b0316602082015260400190565b60405180910390a150505050565b81610d4981612ec3565b610d51612f18565b610d5a836114de565b15610d7857604051631685ecdd60e31b815260040160405180910390fd5b610d828383612f46565b505050565b610d8f612e67565b610d998282613029565b5050565b610da5612f18565b610dad612e67565b600b54600160601b90046001600160601b0316610dcb8115156130e3565b600b80546bffffffffffffffffffffffff60601b19169055600a8054829190600c90610e08908490600160601b90046001600160601b0316615670565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550610d9982826001600160601b0316613101565b6000610e49612f18565b60005a9050610324361115610e7b57604051630f28961b60e01b81523660048201526103246024820152604401610b9a565b6000610e878686613175565b90506000610e9d85836000015160200151613474565b60408301519091506060906000610eb960808a018a8501615690565b63ffffffff169050806001600160401b03811115610ed957610ed96156ad565b604051908082528060200260200182016040528015610f02578160200160208202803683370190505b50925060005b81811015610f6a5760408051602081018590529081018290526060016040516020818303038152906040528051906020012060001c848281518110610f4f57610f4f615602565b6020908102919091010152610f6381615657565b9050610f08565b5050602080850180516000908152600f9092526040822082905551610f90908a856134cf565b60208a8101356000908152600690915260409020805491925090601890610fc690600160c01b90046001600160401b03166156c3565b91906101000a8154816001600160401b0302191690836001600160401b03160217905550600460008a60800160208101906110019190615198565b6001600160a01b03168152602080820192909252604090810160009081208c84013582529092529020805460099061104890600160481b90046001600160401b03166156e9565b91906101000a8154816001600160401b0302191690836001600160401b031602179055506000898060a0019061107e919061570c565b600161108d60a08e018e61570c565b61109892915061562e565b8181106110a7576110a7615602565b9091013560f81c6001149150600090506110c38887848d613586565b9099509050801561110e5760208088015160105460408051928352928201527f6ca648a381f22ead7e37773d934e64885dcf861fbfbb26c40354cbf0c4662d1a910160405180910390a15b5061111e88828c602001356135be565b602086810151604080518681526001600160601b038c16818501528415158183015285151560608201528c151560808201529051928d0135927faeb4b4786571e184246d39587f659abf0e26f41f6a3358692250382c0cdb47b79181900360a00190a3505050505050505b9392505050565b611198612f18565b6111a1816136f3565b6111c957604051635428d44960e01b81526001600160a01b0382166004820152602401610b9a565b6000806000806111d886612bbc565b945094505093509350336001600160a01b0316826001600160a01b03161461121e57604051636c51fda960e11b81526001600160a01b0383166004820152602401610b9a565b611227866114de565b1561124557604051631685ecdd60e31b815260040160405180910390fd5b6040805160c0810182526001815260208082018990526001600160a01b03851682840152606082018490526001600160601b038088166080840152861660a08301529151909160009161129a91849101615759565b60405160208183030381529060405290506112b48861375e565b505060405163ce3f471960e01b81526001600160a01b0388169063ce3f4719906001600160601b038816906112ed90859060040161581e565b6000604051808303818588803b15801561130657600080fd5b505af115801561131a573d6000803e3d6000fd5b50506002546001600160a01b031615801593509150611343905057506001600160601b03861615155b156113cf5760025460405163a9059cbb60e01b81526001600160a01b0389811660048301526001600160601b03891660248301526113cf92169063a9059cbb906044015b6020604051808303816000875af11580156113a6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113ca9190615831565b6130e3565b600c805466ff00000000000019166601000000000000179055825160005b8181101561147f5784818151811061140757611407615602565b6020908102919091010151604051638ea9811760e01b81526001600160a01b038b8116600483015290911690638ea9811790602401600060405180830381600087803b15801561145657600080fd5b505af115801561146a573d6000803e3d6000fd5b505050508061147890615657565b90506113ed565b50600c805466ff00000000000019169055604080516001600160a01b038a168152602081018b90527fd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be418791015b60405180910390a1505050505050505050565b60008181526005602052604081206002018054825b818110156115755760006004600085848154811061151357611513615602565b60009182526020808320909101546001600160a01b0316835282810193909352604091820181208982529092529020546001600160401b03600160481b90910416111561156557506001949350505050565b61156e81615657565b90506114f3565b506000949350505050565b611588612f18565b611590612e67565b6002546001600160a01b03166115b95760405163c1f0c0a160e01b815260040160405180910390fd5b600b546001600160601b03166115d08115156130e3565b600b80546bffffffffffffffffffffffff19169055600a80548291906000906116039084906001600160601b0316615670565b82546101009290920a6001600160601b0381810219909316918316021790915560025460405163a9059cbb60e01b81526001600160a01b0386811660048301529285166024820152610d99935091169063a9059cbb90604401611387565b611669612e67565b611672816136f3565b1561169b5760405163ac8a27ef60e01b81526001600160a01b0382166004820152602401610b9a565b601180546001810182556000919091527f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c680180546001600160a01b0319166001600160a01b0383169081179091556040519081527fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af016259060200160405180910390a150565b611727612e67565b6002546001600160a01b03161561175157604051631688c53760e11b815260040160405180910390fd5b600280546001600160a01b039384166001600160a01b03199182161790915560038054929093169116179055565b6001546001600160a01b031633146117d95760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610b9a565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611838612e67565b6000611843836128d0565b6000818152600d602052604090205490915060ff161561187957604051634a0b8fa760e01b815260048101829052602401610b9a565b60408051808201825260018082526001600160401b0385811660208085018281526000888152600d835287812096518754925168ffffffffffffffffff1990931690151568ffffffffffffffff00191617610100929095169190910293909317909455600e805493840181559091527fbb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fd9091018490558251848152918201527f9b911b2c240bfbef3b6a8f7ed6ee321d1258bb2a3fe6becab52ac1cd3210afd39101610b5c565b611948612e67565b600a544790600160601b90046001600160601b031681811115611988576040516354ced18160e11b81526004810182905260248101839052604401610b9a565b81811015610d8257600061199c828461562e565b90506000846001600160a01b03168260405160006040518083038185875af1925050503d80600081146119eb576040519150601f19603f3d011682016040523d82523d6000602084013e6119f0565b606091505b5050905080611a125760405163950b247960e01b815260040160405180910390fd5b604080516001600160a01b0387168152602081018490527f4aed7c8eed0496c8c19ea2681fcca25741c1602342e38b045d9f1e8e905d2e9c910160405180910390a15050505050565b611a63612f18565b600081815260056020526040902054611a84906001600160a01b0316613910565b60008181526006602052604090208054600160601b90046001600160601b0316903490600c611ab3838561584e565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555034600a600c8282829054906101000a90046001600160601b0316611afb919061584e565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f7603b205d03651ee812f803fccde89f1012e545a9c99f0abfea9cedd0fd8e902823484611b4e919061586e565b604080519283526020830191909152015b60405180910390a25050565b6000611b75612f18565b60208083013560008181526005909252604090912054611b9d906001600160a01b0316613910565b336000908152600460209081526040808320848452808352928190208151606081018352905460ff811615158083526001600160401b036101008304811695840195909552600160481b9091049093169181019190915290611c1b576040516379bfd40160e01b815260048101849052336024820152604401610b9a565b600c5461ffff16611c326060870160408801615881565b61ffff161080611c55575060c8611c4f6060870160408801615881565b61ffff16115b15611c9b57611c6a6060860160408701615881565b600c5460405163539c34bb60e11b815261ffff92831660048201529116602482015260c86044820152606401610b9a565b600c5462010000900463ffffffff16611cba6080870160608801615690565b63ffffffff161115611d0a57611cd66080860160608701615690565b600c54604051637aebf00f60e11b815263ffffffff9283166004820152620100009091049091166024820152604401610b9a565b6101f4611d1d60a0870160808801615690565b63ffffffff161115611d6357611d3960a0860160808701615690565b6040516311ce1afb60e21b815263ffffffff90911660048201526101f46024820152604401610b9a565b806020018051611d72906156c3565b6001600160401b03169052604081018051611d8c906156c3565b6001600160401b03908116909152602082810151604080518935818501819052338284015260608201899052929094166080808601919091528151808603909101815260a08501825280519084012060c085019290925260e08085018390528151808603909101815261010090940190528251929091019190912060009190955090506000611e2e611e29611e2460a08a018a61570c565b613937565b6139b8565b9050854386611e4360808b0160608c01615690565b611e5360a08c0160808d01615690565b3386604051602001611e6b979695949392919061589c565b60405160208183030381529060405280519060200120600f600088815260200190815260200160002081905550336001600160a01b03168588600001357feb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e89868c6040016020810190611ede9190615881565b8d6060016020810190611ef19190615690565b8e6080016020810190611f049190615690565b89604051611f17969594939291906158f3565b60405180910390a45050600092835260209182526040928390208151815493830151929094015168ffffffffffffffffff1990931693151568ffffffffffffffff001916939093176101006001600160401b03928316021770ffffffffffffffff0000000000000000001916600160481b91909216021790555b919050565b6000611fa0612f18565b6007546001600160401b031633611fb860014361562e565b6040516bffffffffffffffffffffffff19606093841b81166020830152914060348201523090921b1660548201526001600160c01b031960c083901b16606882015260700160408051601f1981840301815291905280516020909101209150612022816001615932565b6007805467ffffffffffffffff19166001600160401b03928316179055604080516000808252608082018352602080830182815283850183815260608086018581528a86526006855287862093518454935191516001600160601b039182166001600160c01b031990951694909417600160601b91909216021777ffffffffffffffffffffffffffffffffffffffffffffffff16600160c01b9290981691909102969096179055835194850184523385528481018281528585018481528884526005835294909220855181546001600160a01b03199081166001600160a01b0392831617835593516001830180549095169116179092559251805192949391926121329260028501920190615091565b5061214291506008905084613a29565b5060405133815283907f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d9060200160405180910390a2505090565b612185612f18565b6002546001600160a01b031633146121b0576040516344b0e3c360e01b815260040160405180910390fd5b602081146121d157604051638129bbcd60e01b815260040160405180910390fd5b60006121df828401846152e8565b600081815260056020526040902054909150612203906001600160a01b0316613910565b600081815260066020526040812080546001600160601b03169186919061222a838561584e565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600a60008282829054906101000a90046001600160601b0316612272919061584e565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a8287846122c5919061586e565b6040805192835260208301919091520160405180910390a2505050505050565b6122ed612e67565b60c861ffff8a1611156123275760405163539c34bb60e11b815261ffff8a1660048201819052602482015260c86044820152606401610b9a565b6000851361234b576040516321ea67b360e11b815260048101869052602401610b9a565b8363ffffffff168363ffffffff161115612388576040516313c06e5960e11b815263ffffffff808516600483015285166024820152604401610b9a565b609b60ff831611156123b957604051631d66288d60e11b815260ff83166004820152609b6024820152604401610b9a565b609b60ff821611156123ea57604051631d66288d60e11b815260ff82166004820152609b6024820152604401610b9a565b604080516101208101825261ffff8b1680825263ffffffff808c16602084018190526000848601528b8216606085018190528b8316608086018190528a841660a08701819052938a1660c0870181905260ff808b1660e08901819052908a16610100909801889052600c8054600160c01b90990260ff60c01b19600160b81b9093029290921661ffff60b81b19600160981b90940263ffffffff60981b19600160781b9099029890981676ffffffffffffffff00000000000000000000000000000019600160581b9096026effffffff000000000000000000000019670100000000000000909802979097166effffffffffffffffff000000000000196201000090990265ffffffffffff19909c16909a179a909a1796909616979097179390931791909116959095179290921793909316929092179190911790556010869055517f2c6b6b12413678366b05b145c5f00745bdd00e739131ab5de82484a50c9d78b6906114cb908b908b908b908b908b908b908b908b908b9061ffff99909916895263ffffffff97881660208a0152958716604089015293861660608801526080870192909252841660a086015290921660c084015260ff91821660e0840152166101008201526101200190565b6125c1612e67565b6000818152600560205260409020546001600160a01b03166125e281613910565b610d998282612f46565b606060006125fa6008613a35565b905080841061261c57604051631390f2a160e01b815260040160405180910390fd5b6000612628848661586e565b905081811180612636575083155b6126405780612642565b815b90506000612650868361562e565b9050806001600160401b0381111561266a5761266a6156ad565b604051908082528060200260200182016040528015612693578160200160208202803683370190505b50935060005b818110156126e3576126b66126ae888361586e565b600890613a3f565b8582815181106126c8576126c8615602565b60209081029190910101526126dc81615657565b9050612699565b505050505b92915050565b6126f6612f18565b6000818152600560205260409020546001600160a01b031661271781613910565b6000828152600560205260409020600101546001600160a01b03163314612770576000828152600560205260409081902060010154905163d084e97560e01b81526001600160a01b039091166004820152602401610b9a565b6000828152600560209081526040918290208054336001600160a01b03199182168117835560019092018054909116905582516001600160a01b03851681529182015283917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c93869101611b5f565b816127e781612ec3565b6127ef612f18565b6001600160a01b03821660009081526004602090815260408083208684529091529020805460ff16156128225750505050565b6000848152600560205260409020600201805460631901612856576040516305a48e0f60e01b815260040160405180910390fd5b8154600160ff199091168117835581549081018255600082815260209081902090910180546001600160a01b0319166001600160a01b03871690811790915560405190815286917f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e191015b60405180910390a25050505050565b6000816040516020016128e39190615952565b604051602081830303815290604052805190602001209050919050565b8161290a81612ec3565b612912612f18565b61291b836114de565b1561293957604051631685ecdd60e31b815260040160405180910390fd5b6001600160a01b038216600090815260046020908152604080832086845290915290205460ff1661298f576040516379bfd40160e01b8152600481018490526001600160a01b0383166024820152604401610b9a565b6000838152600560205260408120600201805490915b81811015612aa457846001600160a01b03168382815481106129c9576129c9615602565b6000918252602090912001546001600160a01b031603612a9457826129ef60018461562e565b815481106129ff576129ff615602565b9060005260206000200160009054906101000a90046001600160a01b0316838281548110612a2f57612a2f615602565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555082805480612a6d57612a6d615641565b600082815260209020810160001990810180546001600160a01b0319169055019055612aa4565b612a9d81615657565b90506129a5565b506001600160a01b0384166000818152600460209081526040808320898452825291829020805460ff19169055905191825286917f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a791016128c1565b600e8181548110612b1057600080fd5b600091825260209091200154905081565b81612b2b81612ec3565b612b33612f18565b600083815260056020526040902060018101546001600160a01b03848116911614612bb6576001810180546001600160a01b0319166001600160a01b03851690811790915560408051338152602081019290925285917f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a191015b60405180910390a25b50505050565b600081815260056020526040812054819081906001600160a01b03166060612be382613910565b600086815260066020908152604080832054600583529281902060020180548251818502810185019093528083526001600160601b0380861695600160601b810490911694600160c01b9091046001600160401b0316938893929091839190830182828015612c7b57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612c5d575b505050505090509450945094509450945091939590929450565b612c9d612e67565b6002546001600160a01b0316612cc65760405163c1f0c0a160e01b815260040160405180910390fd5b6002546040516370a0823160e01b81523060048201526000916001600160a01b0316906370a0823190602401602060405180830381865afa158015612d0f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d339190615961565b600a549091506001600160601b031681811115612d6d576040516354ced18160e11b81526004810182905260248101839052604401610b9a565b81811015610d82576000612d81828461562e565b60025460405163a9059cbb60e01b81526001600160a01b0387811660048301526024820184905292935091169063a9059cbb906044016020604051808303816000875af1158015612dd6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612dfa9190615831565b612e1757604051631f01ff1360e21b815260040160405180910390fd5b604080516001600160a01b0386168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b4366009101610d31565b612e5e612e67565b610ba381613a4b565b6000546001600160a01b03163314612ec15760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610b9a565b565b6000818152600560205260409020546001600160a01b0316612ee481613910565b336001600160a01b03821614610d9957604051636c51fda960e11b81526001600160a01b0382166004820152602401610b9a565b600c546601000000000000900460ff1615612ec15760405163769dd35360e11b815260040160405180910390fd5b600080612f528461375e565b60025491935091506001600160a01b031615801590612f7957506001600160601b03821615155b15612fc15760025460405163a9059cbb60e01b81526001600160a01b0385811660048301526001600160601b0385166024830152612fc192169063a9059cbb90604401611387565b612fd483826001600160601b0316613101565b604080516001600160a01b03851681526001600160601b03808516602083015283169181019190915284907f8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c490606001612bad565b60038260ff161061305157604051621c300f60e81b815260ff83166004820152602401610b9a565b60ff81161580613064575060648160ff16115b156130865760405162d4503560e51b815260ff82166004820152602401610b9a565b6012805460ff84811661ffff199092168217610100918516918202179092556040805191825260208201929092527f8e63dc2f2e669ce73bebd2580bb9dd9a5d17fa2d046ac02057d8349fc0b0c2f3910160405180910390a15050565b80610ba357604051631e9acf1760e31b815260040160405180910390fd5b6000826001600160a01b03168260405160006040518083038185875af1925050503d806000811461314e576040519150601f19603f3d011682016040523d82523d6000602084013e613153565b606091505b5050905080610d825760405163950b247960e01b815260040160405180910390fd5b6040805160a08101825260006060820181815260808301829052825260208201819052918101829052906131a8846128d0565b6000818152600d602090815260409182902082518084019093525460ff811615158084526101009091046001600160401b0316918301919091529192509061320657604051631dfd6e1360e21b815260048101839052602401610b9a565b6000828660c00135604051602001613228929190918252602082015260400190565b60408051601f1981840301815291815281516020928301206000818152600f909352908220549092509081900361327257604051631b44092560e11b815260040160405180910390fd5b81613280602088018861597a565b602088013561329560608a0160408b01615690565b6132a560808b0160608c01615690565b6132b560a08c0160808d01615198565b6132c260a08d018d61570c565b6040516020016132d9989796959493929190615995565b60405160208183030381529060405280519060200120811461330e5760405163354a450b60e21b815260040160405180910390fd5b600061332d613320602089018961597a565b6001600160401b03164090565b905080613411576001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001663e9413d3861337060208a018a61597a565b6040516001600160e01b031960e084901b1681526001600160401b039091166004820152602401602060405180830381865afa1580156133b4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133d89190615961565b905080613411576133ec602088018861597a565b60405163175dadad60e01b81526001600160401b039091166004820152602401610b9a565b6040805160c08a013560208083019190915281830184905282518083038401815260609092019092528051910120600061344b8a83613af4565b604080516060810182529788526020880196909652948601949094525092979650505050505050565b6000816001600160401b03163a11156134c757821561349d57506001600160401b0381166126e8565b60405163435e532d60e11b81523a60048201526001600160401b0383166024820152604401610b9a565b503a92915050565b6000806000631fe543e360e01b86856040516024016134ef929190615a0d565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252600c805466ff00000000000019166601000000000000179055915061356c906135509060608801908801615690565b63ffffffff1661356660a0880160808901615198565b83613bea565b600c805466ff000000000000191690559695505050505050565b60008083156135a55761359a868685613c36565b6000915091506135b5565b6135b0868685613d4c565b915091505b94509492505050565b6000818152600660205260409020821561366e5780546001600160601b03600160601b9091048116906135f59086168210156130e3565b6135ff8582615670565b82546bffffffffffffffffffffffff60601b1916600160601b6001600160601b039283168102919091178455600b805488939192600c9261364492869290041661584e565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555050612bb6565b80546001600160601b039081169061368a9086168210156130e3565b6136948582615670565b82546bffffffffffffffffffffffff19166001600160601b03918216178355600b805487926000916136c89185911661584e565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505050505050565b601154600090815b8181101561375457836001600160a01b03166011828154811061372057613720615602565b6000918252602090912001546001600160a01b031603613744575060019392505050565b61374d81615657565b90506136fb565b5060009392505050565b60008181526005602090815260408083206006909252822054600290910180546001600160601b0380841694600160601b90940416925b8181101561380a57600460008483815481106137b3576137b3615602565b60009182526020808320909101546001600160a01b0316835282810193909352604091820181208982529092529020805470ffffffffffffffffffffffffffffffffff1916905561380381615657565b9050613795565b50600085815260056020526040812080546001600160a01b0319908116825560018201805490911690559061384260028301826150f6565b505060008581526006602052604081205561385e600886613f43565b506001600160601b038416156138b157600a805485919060009061388c9084906001600160601b0316615670565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b6001600160601b038316156139095782600a600c8282829054906101000a90046001600160601b03166138e49190615670565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b5050915091565b6001600160a01b038116610ba357604051630fb532db60e11b815260040160405180910390fd5b604080516020810190915260008152600082900361396457506040805160208101909152600081526126e8565b63125fa26760e31b6139768385615a26565b6001600160e01b0319161461399e57604051632923fee760e11b815260040160405180910390fd5b6139ab8260048186615a56565b8101906111899190615a80565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa826040516024016139f191511515815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b03199093169290921790915292915050565b60006111898383613f4f565b60006126e8825490565b60006111898383613f9e565b336001600160a01b03821603613aa35760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610b9a565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b604080518082018252600091613bb49190859060029083908390808284376000920191909152505060408051808201825291508087019060029083908390808284376000920191909152505050608086013560a087013586613b5d6101008a0160e08b01615198565b604080518082018252906101008c019060029083908390808284376000920191909152505060408051808201825291506101408d0190600290839083908082843760009201919091525050506101808c0135613fc8565b600383604001604051602001613bcb929190615ad9565b60408051601f1981840301815291905280516020909101209392505050565b60005a611388811015613bfc57600080fd5b611388810390508460408204820311613c1457600080fd5b50823b613c2057600080fd5b60008083516020850160008789f1949350505050565b600080613c446000366141f3565b905060005a600c54613c64908890600160581b900463ffffffff1661586e565b613c6e919061562e565b613c789086615aef565b600c54909150600090613c9d90600160781b900463ffffffff1664e8d4a51000615aef565b90508215613cd9576040518381527f56296f7beae05a0db815737fdb4cd298897b1e517614d62468081531ae14d0999060200160405180910390a15b8415613d2357600c548190606490600160b81b900460ff16613cfb858761586e565b613d059190615aef565b613d0f9190615b1c565b613d19919061586e565b9350505050611189565b600c548190606490613d3f90600160b81b900460ff1682615b30565b60ff16613cfb858761586e565b600080600080613d5a6141ff565b9150915060008213613d82576040516321ea67b360e11b815260048101839052602401610b9a565b6000613d8f6000366141f3565b9050600083825a600c54613db1908d90600160581b900463ffffffff1661586e565b613dbb919061562e565b613dc5908b615aef565b613dcf919061586e565b613de190670de0b6b3a7640000615aef565b613deb9190615b1c565b600c54909150600090613e149063ffffffff600160981b8204811691600160781b900416615b49565b613e299063ffffffff1664e8d4a51000615aef565b9050600085613e4083670de0b6b3a7640000615aef565b613e4a9190615b1c565b90508315613e86576040518481527f56296f7beae05a0db815737fdb4cd298897b1e517614d62468081531ae14d0999060200160405180910390a15b60008915613ec557600c548290606490613eaa90600160c01b900460ff1687615aef565b613eb49190615b1c565b613ebe919061586e565b9050613f05565b600c548290606490613ee190600160c01b900460ff1682615b30565b613eee9060ff1687615aef565b613ef89190615b1c565b613f02919061586e565b90505b6b033b2e3c9fd0803ce8000000811115613f325760405163e80fa38160e01b815260040160405180910390fd5b9b949a509398505050505050505050565b600061118983836142ca565b6000818152600183016020526040812054613f96575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556126e8565b5060006126e8565b6000826000018281548110613fb557613fb5615602565b9060005260206000200154905092915050565b613fd1896143c4565b61401d5760405162461bcd60e51b815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e2063757276650000000000006044820152606401610b9a565b614026886143c4565b6140725760405162461bcd60e51b815260206004820152601560248201527f67616d6d61206973206e6f74206f6e20637572766500000000000000000000006044820152606401610b9a565b61407b836143c4565b6140c75760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610b9a565b6140d0826143c4565b61411c5760405162461bcd60e51b815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e206375727665000000006044820152606401610b9a565b614128878a888761449d565b6141745760405162461bcd60e51b815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e657373000000000000006044820152606401610b9a565b60006141808a876145c0565b90506000614193898b878b868989614624565b905060006141a4838d8d8a86614750565b9050808a146141e55760405162461bcd60e51b815260206004820152600d60248201526c34b73b30b634b210383937b7b360991b6044820152606401610b9a565b505050505050505050505050565b60006111898383614790565b600c5460035460408051633fabe5a360e21b81529051600093849367010000000000000090910463ffffffff169284926001600160a01b039092169163feaf968c9160048082019260a0929091908290030181865afa158015614266573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061428a9190615b80565b50919650909250505063ffffffff8216158015906142b657506142ad814261562e565b8263ffffffff16105b925082156142c45760105493505b50509091565b600081815260018301602052604081205480156143b35760006142ee60018361562e565b85549091506000906143029060019061562e565b905081811461436757600086600001828154811061432257614322615602565b906000526020600020015490508087600001848154811061434557614345615602565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061437857614378615641565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506126e8565b60009150506126e8565b5092915050565b80516000906401000003d0191161441d5760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420782d6f7264696e61746500000000000000000000000000006044820152606401610b9a565b60208201516401000003d019116144765760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420792d6f7264696e61746500000000000000000000000000006044820152606401610b9a565b60208201516401000003d0199080096144968360005b602002015161485c565b1492915050565b60006001600160a01b0382166144e35760405162461bcd60e51b815260206004820152600b60248201526a626164207769746e65737360a81b6044820152606401610b9a565b6020840151600090600116156144fa57601c6144fd565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe19918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa158015614598573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b6145c8615114565b6145f5600184846040516020016145e193929190615bf3565b604051602081830303815290604052614880565b90505b614601816143c4565b6126e857805160408051602081019290925261461d91016145e1565b90506145f8565b61462c615114565b825186516401000003d019918290069190060361468b5760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610b9a565b6146968789886148cd565b6146e25760405162461bcd60e51b815260206004820152601660248201527f4669727374206d756c20636865636b206661696c6564000000000000000000006044820152606401610b9a565b6146ed8486856148cd565b6147395760405162461bcd60e51b815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c65640000000000000000006044820152606401610b9a565b6147448684846149f8565b98975050505050505050565b60006002868686858760405160200161476e96959493929190615c14565b60408051601f1981840301815291905280516020909101209695505050505050565b60125460009060ff1661485357600f602160991b016001600160a01b03166349948e0e8484604051806080016040528060478152602001615dee604791396040516020016147e093929190615c73565b6040516020818303038152906040526040518263ffffffff1660e01b815260040161480b919061581e565b602060405180830381865afa158015614828573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061484c9190615961565b90506126e8565b61118982614abf565b6000806401000003d01980848509840990506401000003d019600782089392505050565b614888615114565b61489182614ba7565b81526148a66148a182600061448c565b614be2565b6020820181905260029006600103611f91576020810180516401000003d019039052919050565b60008260000361490d5760405162461bcd60e51b815260206004820152600b60248201526a3d32b9379039b1b0b630b960a91b6044820152606401610b9a565b8351602085015160009061492390600290615c9a565b1561492f57601c614932565b601b5b9050600070014551231950b75fc4402da1732fc9bebe198387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa1580156149a4573d6000803e3d6000fd5b5050506020604051035190506000866040516020016149c39190615cae565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b614a00615114565b835160208086015185519186015160009384938493614a2193909190614c02565b919450925090506401000003d019858209600114614a815760405162461bcd60e51b815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a000000000000006044820152606401610b9a565b60405180604001604052806401000003d01980614aa057614aa0615b06565b87860981526020016401000003d0198785099052979650505050505050565b60125460009060ff166000198101614b05576064614ae6614ae160478661586e565b614ce2565b601254614afb9190610100900460ff16615aef565b6111899190615b1c565b60011960ff821601614b8a576064600f602160991b0163f1c7a58b614b2b60478761586e565b6040518263ffffffff1660e01b8152600401614b4991815260200190565b602060405180830381865afa158015614b66573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614ae69190615961565b604051621c300f60e81b815260ff82166004820152602401610b9a565b805160208201205b6401000003d0198110611f9157604080516020808201939093528151808203840181529082019091528051910120614baf565b60006126e8826002614bfb6401000003d019600161586e565b901c614f7f565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000614c4283838585615024565b9098509050614c5388828e88615048565b9098509050614c6488828c87615048565b90985090506000614c778d878b85615048565b9098509050614c8888828686615024565b9098509050614c9988828e89615048565b9098509050818114614cce576401000003d019818a0998506401000003d01982890997506401000003d0198183099650614cd2565b8196505b5050505050509450945094915050565b600080614cf060448461586e565b614cfb906010615aef565b90506000600f602160991b016001600160a01b031663519b4bd36040518163ffffffff1660e01b8152600401602060405180830381865afa158015614d44573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614d689190615961565b600f602160991b016001600160a01b031663c59859186040518163ffffffff1660e01b8152600401602060405180830381865afa158015614dad573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614dd19190615cc0565b614ddc906010615cdd565b63ffffffff16614dec9190615aef565b90506000600f602160991b016001600160a01b031663f82061406040518163ffffffff1660e01b8152600401602060405180830381865afa158015614e35573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614e599190615961565b600f602160991b016001600160a01b03166368d5dca66040518163ffffffff1660e01b8152600401602060405180830381865afa158015614e9e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614ec29190615cc0565b63ffffffff16614ed29190615aef565b90506000614ee0828461586e565b614eea9085615aef565b9050600f602160991b016001600160a01b031663313ce5676040518163ffffffff1660e01b8152600401602060405180830381865afa158015614f31573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614f559190615961565b614f6090600a615de1565b614f6b906010615aef565b614f759082615b1c565b9695505050505050565b600080614f8a615132565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152614fbc615150565b60208160c0846005600019fa92508260000361501a5760405162461bcd60e51b815260206004820152601260248201527f6269674d6f64457870206661696c7572652100000000000000000000000000006044820152606401610b9a565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b8280548282559060005260206000209081019282156150e6579160200282015b828111156150e657825182546001600160a01b0319166001600160a01b039091161782556020909201916001909101906150b1565b506150f292915061516e565b5090565b5080546000825590600052602060002090810190610ba3919061516e565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b808211156150f2576000815560010161516f565b6001600160a01b0381168114610ba357600080fd5b6000602082840312156151aa57600080fd5b813561118981615183565b80604081018310156126e857600080fd5b6000604082840312156151d857600080fd5b61118983836151b5565b600080604083850312156151f557600080fd5b82359150602083013561520781615183565b809150509250929050565b803560ff81168114611f9157600080fd5b6000806040838503121561523657600080fd5b61523f83615212565b915061524d60208401615212565b90509250929050565b600060c0828403121561526857600080fd5b50919050565b8015158114610ba357600080fd5b60008060008385036101e081121561529357600080fd5b6101a0808212156152a357600080fd5b85945084013590506001600160401b038111156152bf57600080fd5b6152cb86828701615256565b9250506101c08401356152dd8161526e565b809150509250925092565b6000602082840312156152fa57600080fd5b5035919050565b6000806040838503121561531457600080fd5b823561531f81615183565b9150602083013561520781615183565b80356001600160401b0381168114611f9157600080fd5b6000806060838503121561535957600080fd5b61536384846151b5565b915061524d6040840161532f565b60006020828403121561538357600080fd5b81356001600160401b0381111561539957600080fd5b6153a584828501615256565b949350505050565b600080600080606085870312156153c357600080fd5b84356153ce81615183565b93506020850135925060408501356001600160401b03808211156153f157600080fd5b818701915087601f83011261540557600080fd5b81358181111561541457600080fd5b88602082850101111561542657600080fd5b95989497505060200194505050565b803561ffff81168114611f9157600080fd5b63ffffffff81168114610ba357600080fd5b60008060008060008060008060006101208a8c03121561547857600080fd5b6154818a615435565b985060208a013561549181615447565b975060408a01356154a181615447565b965060608a01356154b181615447565b955060808a0135945060a08a01356154c881615447565b935060c08a01356154d881615447565b92506154e660e08b01615212565b91506154f56101008b01615212565b90509295985092959850929598565b6000806040838503121561551757600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b838110156155565781518752958201959082019060010161553a565b509495945050505050565b6020815260006111896020830184615526565b600081518084526020808501945080840160005b838110156155565781516001600160a01b031687529582019590820190600101615588565b60006001600160601b0380881683528087166020840152506001600160401b03851660408301526001600160a01b038416606083015260a060808301526155f760a0830184615574565b979650505050505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b818103818111156126e8576126e8615618565b634e487b7160e01b600052603160045260246000fd5b60006001820161566957615669615618565b5060010190565b6001600160601b038281168282160390808211156143bd576143bd615618565b6000602082840312156156a257600080fd5b813561118981615447565b634e487b7160e01b600052604160045260246000fd5b60006001600160401b038083168181036156df576156df615618565b6001019392505050565b60006001600160401b0382168061570257615702615618565b6000190192915050565b6000808335601e1984360301811261572357600080fd5b8301803591506001600160401b0382111561573d57600080fd5b60200191503681900382131561575257600080fd5b9250929050565b6020815260ff8251166020820152602082015160408201526001600160a01b0360408301511660608201526000606083015160c0608084015261579f60e0840182615574565b905060808401516001600160601b0380821660a08601528060a08701511660c086015250508091505092915050565b60005b838110156157e95781810151838201526020016157d1565b50506000910152565b6000815180845261580a8160208601602086016157ce565b601f01601f19169290920160200192915050565b60208152600061118960208301846157f2565b60006020828403121561584357600080fd5b81516111898161526e565b6001600160601b038181168382160190808211156143bd576143bd615618565b808201808211156126e8576126e8615618565b60006020828403121561589357600080fd5b61118982615435565b878152866020820152856040820152600063ffffffff80871660608401528086166080840152506001600160a01b03841660a083015260e060c08301526158e660e08301846157f2565b9998505050505050505050565b86815285602082015261ffff85166040820152600063ffffffff808616606084015280851660808401525060c060a083015261474460c08301846157f2565b6001600160401b038181168382160190808211156143bd576143bd615618565b60408181019083833792915050565b60006020828403121561597357600080fd5b5051919050565b60006020828403121561598c57600080fd5b6111898261532f565b8881526001600160401b0388166020820152866040820152600063ffffffff80881660608401528087166080840152506001600160a01b03851660a083015260e060c08301528260e08301526101008385828501376000838501820152601f909301601f191690910190910198975050505050505050565b8281526040602082015260006153a56040830184615526565b6001600160e01b03198135818116916004851015615a4e5780818660040360031b1b83161692505b505092915050565b60008085851115615a6657600080fd5b83861115615a7357600080fd5b5050820193919092039150565b600060208284031215615a9257600080fd5b604051602081018181106001600160401b0382111715615ac257634e487b7160e01b600052604160045260246000fd5b6040528235615ad08161526e565b81529392505050565b8281526060810160408360208401379392505050565b80820281158282048414176126e8576126e8615618565b634e487b7160e01b600052601260045260246000fd5b600082615b2b57615b2b615b06565b500490565b60ff81811683821601908111156126e8576126e8615618565b63ffffffff8281168282160390808211156143bd576143bd615618565b805169ffffffffffffffffffff81168114611f9157600080fd5b600080600080600060a08688031215615b9857600080fd5b615ba186615b66565b9450602086015193506040860151925060608601519150615bc460808701615b66565b90509295509295909350565b8060005b6002811015612bb6578151845260209384019390910190600101615bd4565b838152615c036020820184615bd0565b606081019190915260800192915050565b868152615c246020820187615bd0565b615c316060820186615bd0565b615c3e60a0820185615bd0565b615c4b60e0820184615bd0565b60609190911b6bffffffffffffffffffffffff19166101208201526101340195945050505050565b828482376000838201600081528351615c908183602088016157ce565b0195945050505050565b600082615ca957615ca9615b06565b500690565b615cb88183615bd0565b604001919050565b600060208284031215615cd257600080fd5b815161118981615447565b63ffffffff818116838216028082169190828114615a4e57615a4e615618565b600181815b80851115615d38578160001904821115615d1e57615d1e615618565b80851615615d2b57918102915b93841c9390800290615d02565b509250929050565b600082615d4f575060016126e8565b81615d5c575060006126e8565b8160018114615d725760028114615d7c57615d98565b60019150506126e8565b60ff841115615d8d57615d8d615618565b50506001821b6126e8565b5060208310610133831016604e8410600b8410161715615dbb575081810a6126e8565b615dc58383615cfd565b8060001904821115615dd957615dd9615618565b029392505050565b60006111898383615d4056feffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa164736f6c6343000813000a",
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
