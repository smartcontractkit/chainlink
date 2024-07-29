// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_coordinator_v2_5_arbitrum

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

var VRFCoordinatorV25ArbitrumMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendNative\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToTransferLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxGas\",\"type\":\"uint256\"}],\"name\":\"GasPriceExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"premiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"max\",\"type\":\"uint8\"}],\"name\":\"InvalidPremiumPercentage\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"flatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeNativePPM\",\"type\":\"uint32\"}],\"name\":\"LinkDiscountTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"max\",\"type\":\"uint32\"}],\"name\":\"MsgDataTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"name\":\"FallbackWeiPerUnitLinkUsed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"L1GasFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativeFundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"onlyPremium\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountNative\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldNativeBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newNativeBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFundedWithNative\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_NATIVE_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"deregisterMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFTypes.RequestCommitmentV2Plus\",\"name\":\"rc\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"onlyPremium\",\"type\":\"bool\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"fundSubscriptionWithNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveSubscriptionIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"subOwner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverNativeFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"registerMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"reentrancyLock\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_currentSubNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_provingKeyHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_provingKeys\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestCommitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalNativeBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkNativeFeed\",\"type\":\"address\"}],\"name\":\"setLINKAndLINKNativeFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620059e4380380620059e4833981016040819052620000349162000180565b8033806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf81620000d5565b5050506001600160a01b031660805250620001b2565b336001600160a01b038216036200012f5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019357600080fd5b81516001600160a01b0381168114620001ab57600080fd5b9392505050565b60805161580f620001d5600039600081816105d201526131b5015261580f6000f3fe60806040526004361061028c5760003560e01c80638402595e11610164578063b2a7cac5116100c6578063da2f26101161008a578063e72f6e3011610064578063e72f6e3014610904578063ee9d2d3814610924578063f2fde38b1461095157600080fd5b8063da2f261014610854578063dac83d29146108b3578063dc311dd3146108d357600080fd5b8063b2a7cac5146107b4578063bec4c08c146107d4578063caf70c4a146107f4578063cb63179714610814578063d98e620e1461083457600080fd5b80639d40a6fd11610128578063a63e0bfb11610102578063a63e0bfb14610747578063aa433aff14610767578063aefb212f1461078757600080fd5b80639d40a6fd146106da578063a21a23e414610712578063a4c0ed361461072757600080fd5b80638402595e1461064957806386fe91c7146106695780638da5cb5b1461068957806395b55cfc146106a75780639b1c385e146106ba57600080fd5b8063405b84fa1161020d57806364d51a2a116101d157806372e9d565116101ab57806372e9d565146105f457806379ba5097146106145780637a5a2aef1461062957600080fd5b806364d51a2a1461058b57806365982744146105a0578063689c4517146105c057600080fd5b8063405b84fa146104d057806340d6bb82146104f057806341af6c871461051b57806351cff8d91461054b5780635d06b4ab1461056b57600080fd5b806315c48b841161025457806315c48b84146103f157806318e3dd27146104195780631b6b6d23146104585780632f622e6b14610490578063301f42e9146104b057600080fd5b806304104edb14610291578063043bd6ae146102b3578063088070f5146102dc57806308821d58146103b15780630ae09540146103d1575b600080fd5b34801561029d57600080fd5b506102b16102ac366004614d1f565b610971565b005b3480156102bf57600080fd5b506102c960105481565b6040519081526020015b60405180910390f35b3480156102e857600080fd5b50600c546103549061ffff81169063ffffffff62010000820481169160ff660100000000000082048116926701000000000000008304811692600160581b8104821692600160781b8204831692600160981b83041691600160b81b8104821691600160c01b9091041689565b6040805161ffff909a168a5263ffffffff98891660208b01529615159689019690965293861660608801529185166080870152841660a08601529290921660c084015260ff91821660e084015216610100820152610120016102d3565b3480156103bd57600080fd5b506102b16103cc366004614d4d565b610aea565b3480156103dd57600080fd5b506102b16103ec366004614d69565b610c83565b3480156103fd57600080fd5b5061040660c881565b60405161ffff90911681526020016102d3565b34801561042557600080fd5b50600a5461044090600160601b90046001600160601b031681565b6040516001600160601b0390911681526020016102d3565b34801561046457600080fd5b50600254610478906001600160a01b031681565b6040516001600160a01b0390911681526020016102d3565b34801561049c57600080fd5b506102b16104ab366004614d1f565b610ccb565b3480156104bc57600080fd5b506104406104cb366004614dbf565b610d71565b3480156104dc57600080fd5b506102b16104eb366004614d69565b6110c2565b3480156104fc57600080fd5b506105066101f481565b60405163ffffffff90911681526020016102d3565b34801561052757600080fd5b5061053b610536366004614e2b565b611410565b60405190151581526020016102d3565b34801561055757600080fd5b506102b1610566366004614d1f565b6114b2565b34801561057757600080fd5b506102b1610586366004614d1f565b611593565b34801561059757600080fd5b50610406606481565b3480156105ac57600080fd5b506102b16105bb366004614e44565b611651565b3480156105cc57600080fd5b506104787f000000000000000000000000000000000000000000000000000000000000000081565b34801561060057600080fd5b50600354610478906001600160a01b031681565b34801561062057600080fd5b506102b16116b1565b34801561063557600080fd5b506102b1610644366004614e89565b611762565b34801561065557600080fd5b506102b1610664366004614d1f565b611872565b34801561067557600080fd5b50600a54610440906001600160601b031681565b34801561069557600080fd5b506000546001600160a01b0316610478565b6102b16106b5366004614e2b565b61198d565b3480156106c657600080fd5b506102c96106d5366004614ebd565b611a9d565b3480156106e657600080fd5b506007546106fa906001600160401b031681565b6040516001600160401b0390911681526020016102d3565b34801561071e57600080fd5b506102c9611ecf565b34801561073357600080fd5b506102b1610742366004614ef9565b6120b6565b34801561075357600080fd5b506102b1610762366004614fb8565b61221e565b34801561077357600080fd5b506102b1610782366004614e2b565b6124f2565b34801561079357600080fd5b506107a76107a2366004615059565b612525565b6040516102d391906150b6565b3480156107c057600080fd5b506102b16107cf366004614e2b565b612627565b3480156107e057600080fd5b506102b16107ef366004614d69565b612716565b34801561080057600080fd5b506102c961080f366004614d4d565b612809565b34801561082057600080fd5b506102b161082f366004614d69565b612839565b34801561084057600080fd5b506102c961084f366004614e2b565b612a39565b34801561086057600080fd5b5061089461086f366004614e2b565b600d6020526000908152604090205460ff81169061010090046001600160401b031682565b6040805192151583526001600160401b039091166020830152016102d3565b3480156108bf57600080fd5b506102b16108ce366004614d69565b612a5a565b3480156108df57600080fd5b506108f36108ee366004614e2b565b612af5565b6040516102d3959493929190615102565b34801561091057600080fd5b506102b161091f366004614d1f565b612bce565b34801561093057600080fd5b506102c961093f366004614e2b565b600f6020526000908152604090205481565b34801561095d57600080fd5b506102b161096c366004614d1f565b612d8f565b610979612da0565b60115460005b81811015610abd57826001600160a01b0316601182815481106109a4576109a4615157565b6000918252602090912001546001600160a01b031603610aad5760116109cb600184615183565b815481106109db576109db615157565b600091825260209091200154601180546001600160a01b039092169183908110610a0757610a07615157565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b031602179055506011805480610a4657610a46615196565b6000828152602090819020600019908301810180546001600160a01b03191690559091019091556040516001600160a01b03851681527ff80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af3791015b60405180910390a1505050565b610ab6816151ac565b905061097f565b50604051635428d44960e01b81526001600160a01b03831660048201526024015b60405180910390fd5b50565b610af2612da0565b6000610afd82612809565b6000818152600d602090815260409182902082518084019093525460ff811615158084526101009091046001600160401b03169183019190915291925090610b5b57604051631dfd6e1360e21b815260048101839052602401610ade565b6000828152600d60205260408120805468ffffffffffffffffff19169055600e54905b81811015610c2d5783600e8281548110610b9a57610b9a615157565b906000526020600020015403610c1d57600e610bb7600184615183565b81548110610bc757610bc7615157565b9060005260206000200154600e8281548110610be557610be5615157565b600091825260209091200155600e805480610c0257610c02615196565b60019003818190600052602060002001600090559055610c2d565b610c26816151ac565b9050610b7e565b507f9b6868e0eb737bcd72205360baa6bfd0ba4e4819a33ade2db384e8a8025639a5838360200151604051610c759291909182526001600160401b0316602082015260400190565b60405180910390a150505050565b81610c8d81612dfc565b610c95612e51565b610c9e83611410565b15610cbc57604051631685ecdd60e31b815260040160405180910390fd5b610cc68383612e7f565b505050565b610cd3612e51565b610cdb612da0565b600b54600160601b90046001600160601b0316610cf9811515612f62565b600b80546bffffffffffffffffffffffff60601b19169055600a8054829190600c90610d36908490600160601b90046001600160601b03166151c5565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550610d6d82826001600160601b0316612f80565b5050565b6000610d7b612e51565b60005a9050610324361115610dad57604051630f28961b60e01b81523660048201526103246024820152604401610ade565b6000610db98686612ff4565b90506000610dcf858360000151602001516132eb565b60408301519091506060906000610deb60808a018a85016151e5565b63ffffffff169050806001600160401b03811115610e0b57610e0b615200565b604051908082528060200260200182016040528015610e34578160200160208202803683370190505b50925060005b81811015610e9c5760408051602081018590529081018290526060016040516020818303038152906040528051906020012060001c848281518110610e8157610e81615157565b6020908102919091010152610e95816151ac565b9050610e3a565b5050602080850180516000908152600f9092526040822082905551610ec2908a85613346565b60208a8101356000908152600690915260409020805491925090601890610ef890600160c01b90046001600160401b0316615216565b91906101000a8154816001600160401b0302191690836001600160401b03160217905550600460008a6080016020810190610f339190614d1f565b6001600160a01b03168152602080820192909252604090810160009081208c840135825290925290208054600990610f7a90600160481b90046001600160401b031661523c565b91906101000a8154816001600160401b0302191690836001600160401b031602179055506000898060a00190610fb0919061525f565b6001610fbf60a08e018e61525f565b610fca929150615183565b818110610fd957610fd9615157565b9091013560f81c600114915060009050610ff58887848d6133fd565b909950905080156110405760208088015160105460408051928352928201527f6ca648a381f22ead7e37773d934e64885dcf861fbfbb26c40354cbf0c4662d1a910160405180910390a15b5061105088828c60200135613435565b602086810151604080518681526001600160601b038c16818501528415158183015285151560608201528c151560808201529051928d0135927faeb4b4786571e184246d39587f659abf0e26f41f6a3358692250382c0cdb47b79181900360a00190a3505050505050505b9392505050565b6110ca612e51565b6110d38161356a565b6110fb57604051635428d44960e01b81526001600160a01b0382166004820152602401610ade565b60008060008061110a86612af5565b945094505093509350336001600160a01b0316826001600160a01b03161461115057604051636c51fda960e11b81526001600160a01b0383166004820152602401610ade565b61115986611410565b1561117757604051631685ecdd60e31b815260040160405180910390fd5b6040805160c0810182526001815260208082018990526001600160a01b03851682840152606082018490526001600160601b038088166080840152861660a0830152915190916000916111cc918491016152ac565b60405160208183030381529060405290506111e6886135d5565b505060405163ce3f471960e01b81526001600160a01b0388169063ce3f4719906001600160601b0388169061121f908590600401615367565b6000604051808303818588803b15801561123857600080fd5b505af115801561124c573d6000803e3d6000fd5b50506002546001600160a01b031615801593509150611275905057506001600160601b03861615155b156113015760025460405163a9059cbb60e01b81526001600160a01b0389811660048301526001600160601b038916602483015261130192169063a9059cbb906044015b6020604051808303816000875af11580156112d8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112fc919061537a565b612f62565b600c805466ff00000000000019166601000000000000179055825160005b818110156113b15784818151811061133957611339615157565b6020908102919091010151604051638ea9811760e01b81526001600160a01b038b8116600483015290911690638ea9811790602401600060405180830381600087803b15801561138857600080fd5b505af115801561139c573d6000803e3d6000fd5b50505050806113aa906151ac565b905061131f565b50600c805466ff00000000000019169055604080516001600160a01b038a168152602081018b90527fd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be418791015b60405180910390a1505050505050505050565b60008181526005602052604081206002018054825b818110156114a75760006004600085848154811061144557611445615157565b60009182526020808320909101546001600160a01b0316835282810193909352604091820181208982529092529020546001600160401b03600160481b90910416111561149757506001949350505050565b6114a0816151ac565b9050611425565b506000949350505050565b6114ba612e51565b6114c2612da0565b6002546001600160a01b03166114eb5760405163c1f0c0a160e01b815260040160405180910390fd5b600b546001600160601b0316611502811515612f62565b600b80546bffffffffffffffffffffffff19169055600a80548291906000906115359084906001600160601b03166151c5565b82546101009290920a6001600160601b0381810219909316918316021790915560025460405163a9059cbb60e01b81526001600160a01b0386811660048301529285166024820152610d6d935091169063a9059cbb906044016112b9565b61159b612da0565b6115a48161356a565b156115cd5760405163ac8a27ef60e01b81526001600160a01b0382166004820152602401610ade565b601180546001810182556000919091527f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c680180546001600160a01b0319166001600160a01b0383169081179091556040519081527fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af016259060200160405180910390a150565b611659612da0565b6002546001600160a01b03161561168357604051631688c53760e11b815260040160405180910390fd5b600280546001600160a01b039384166001600160a01b03199182161790915560038054929093169116179055565b6001546001600160a01b0316331461170b5760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610ade565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61176a612da0565b600061177583612809565b6000818152600d602052604090205490915060ff16156117ab57604051634a0b8fa760e01b815260048101829052602401610ade565b60408051808201825260018082526001600160401b0385811660208085018281526000888152600d835287812096518754925168ffffffffffffffffff1990931690151568ffffffffffffffff00191617610100929095169190910293909317909455600e805493840181559091527fbb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fd9091018490558251848152918201527f9b911b2c240bfbef3b6a8f7ed6ee321d1258bb2a3fe6becab52ac1cd3210afd39101610aa0565b61187a612da0565b600a544790600160601b90046001600160601b0316818111156118ba576040516354ced18160e11b81526004810182905260248101839052604401610ade565b81811015610cc65760006118ce8284615183565b90506000846001600160a01b03168260405160006040518083038185875af1925050503d806000811461191d576040519150601f19603f3d011682016040523d82523d6000602084013e611922565b606091505b50509050806119445760405163950b247960e01b815260040160405180910390fd5b604080516001600160a01b0387168152602081018490527f4aed7c8eed0496c8c19ea2681fcca25741c1602342e38b045d9f1e8e905d2e9c910160405180910390a15050505050565b611995612e51565b6000818152600560205260409020546119b6906001600160a01b0316613787565b60008181526006602052604090208054600160601b90046001600160601b0316903490600c6119e58385615397565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555034600a600c8282829054906101000a90046001600160601b0316611a2d9190615397565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f7603b205d03651ee812f803fccde89f1012e545a9c99f0abfea9cedd0fd8e902823484611a8091906153b7565b604080519283526020830191909152015b60405180910390a25050565b6000611aa7612e51565b60208083013560008181526005909252604090912054611acf906001600160a01b0316613787565b336000908152600460209081526040808320848452808352928190208151606081018352905460ff811615158083526001600160401b036101008304811695840195909552600160481b9091049093169181019190915290611b4d576040516379bfd40160e01b815260048101849052336024820152604401610ade565b600c5461ffff16611b6460608701604088016153ca565b61ffff161080611b87575060c8611b8160608701604088016153ca565b61ffff16115b15611bcd57611b9c60608601604087016153ca565b600c5460405163539c34bb60e11b815261ffff92831660048201529116602482015260c86044820152606401610ade565b600c5462010000900463ffffffff16611bec60808701606088016151e5565b63ffffffff161115611c3c57611c0860808601606087016151e5565b600c54604051637aebf00f60e11b815263ffffffff9283166004820152620100009091049091166024820152604401610ade565b6101f4611c4f60a08701608088016151e5565b63ffffffff161115611c9557611c6b60a08601608087016151e5565b6040516311ce1afb60e21b815263ffffffff90911660048201526101f46024820152604401610ade565b806020018051611ca490615216565b6001600160401b03169052604081018051611cbe90615216565b6001600160401b03908116909152602082810151604080518935818501819052338284015260608201899052929094166080808601919091528151808603909101815260a08501825280519084012060c085019290925260e08085018390528151808603909101815261010090940190528251929091019190912060009190955090506000611d60611d5b611d5660a08a018a61525f565b6137ae565b61382f565b905085611d6b6138a0565b86611d7c60808b0160608c016151e5565b611d8c60a08c0160808d016151e5565b3386604051602001611da497969594939291906153e5565b60405160208183030381529060405280519060200120600f600088815260200190815260200160002081905550336001600160a01b03168588600001357feb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e89868c6040016020810190611e1791906153ca565b8d6060016020810190611e2a91906151e5565b8e6080016020810190611e3d91906151e5565b89604051611e509695949392919061543c565b60405180910390a45050600092835260209182526040928390208151815493830151929094015168ffffffffffffffffff1990931693151568ffffffffffffffff001916939093176101006001600160401b03928316021770ffffffffffffffff0000000000000000001916600160481b91909216021790555b919050565b6000611ed9612e51565b6007546001600160401b031633611ef1600143615183565b6040516bffffffffffffffffffffffff19606093841b81166020830152914060348201523090921b1660548201526001600160c01b031960c083901b16606882015260700160408051601f1981840301815291905280516020909101209150611f5b81600161547b565b6007805467ffffffffffffffff19166001600160401b03928316179055604080516000808252608082018352602080830182815283850183815260608086018581528a86526006855287862093518454935191516001600160601b039182166001600160c01b031990951694909417600160601b91909216021777ffffffffffffffffffffffffffffffffffffffffffffffff16600160c01b9290981691909102969096179055835194850184523385528481018281528585018481528884526005835294909220855181546001600160a01b03199081166001600160a01b03928316178355935160018301805490951691161790925592518051929493919261206b9260028501920190614c18565b5061207b9150600890508461390a565b5060405133815283907f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d9060200160405180910390a2505090565b6120be612e51565b6002546001600160a01b031633146120e9576040516344b0e3c360e01b815260040160405180910390fd5b6020811461210a57604051638129bbcd60e01b815260040160405180910390fd5b600061211882840184614e2b565b60008181526005602052604090205490915061213c906001600160a01b0316613787565b600081815260066020526040812080546001600160601b0316918691906121638385615397565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600a60008282829054906101000a90046001600160601b03166121ab9190615397565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a8287846121fe91906153b7565b6040805192835260208301919091520160405180910390a2505050505050565b612226612da0565b60c861ffff8a1611156122605760405163539c34bb60e11b815261ffff8a1660048201819052602482015260c86044820152606401610ade565b60008513612284576040516321ea67b360e11b815260048101869052602401610ade565b8363ffffffff168363ffffffff1611156122c1576040516313c06e5960e11b815263ffffffff808516600483015285166024820152604401610ade565b609b60ff831611156122f257604051631d66288d60e11b815260ff83166004820152609b6024820152604401610ade565b609b60ff8216111561232357604051631d66288d60e11b815260ff82166004820152609b6024820152604401610ade565b604080516101208101825261ffff8b1680825263ffffffff808c16602084018190526000848601528b8216606085018190528b8316608086018190528a841660a08701819052938a1660c0870181905260ff808b1660e08901819052908a16610100909801889052600c8054600160c01b90990260ff60c01b19600160b81b9093029290921661ffff60b81b19600160981b90940263ffffffff60981b19600160781b9099029890981676ffffffffffffffff00000000000000000000000000000019600160581b9096026effffffff000000000000000000000019670100000000000000909802979097166effffffffffffffffff000000000000196201000090990265ffffffffffff19909c16909a179a909a1796909616979097179390931791909116959095179290921793909316929092179190911790556010869055517f2c6b6b12413678366b05b145c5f00745bdd00e739131ab5de82484a50c9d78b6906113fd908b908b908b908b908b908b908b908b908b9061ffff99909916895263ffffffff97881660208a0152958716604089015293861660608801526080870192909252841660a086015290921660c084015260ff91821660e0840152166101008201526101200190565b6124fa612da0565b6000818152600560205260409020546001600160a01b031661251b81613787565b610d6d8282612e7f565b606060006125336008613916565b905080841061255557604051631390f2a160e01b815260040160405180910390fd5b600061256184866153b7565b90508181118061256f575083155b612579578061257b565b815b905060006125898683615183565b9050806001600160401b038111156125a3576125a3615200565b6040519080825280602002602001820160405280156125cc578160200160208202803683370190505b50935060005b8181101561261c576125ef6125e788836153b7565b600890613920565b85828151811061260157612601615157565b6020908102919091010152612615816151ac565b90506125d2565b505050505b92915050565b61262f612e51565b6000818152600560205260409020546001600160a01b031661265081613787565b6000828152600560205260409020600101546001600160a01b031633146126a9576000828152600560205260409081902060010154905163d084e97560e01b81526001600160a01b039091166004820152602401610ade565b6000828152600560209081526040918290208054336001600160a01b03199182168117835560019092018054909116905582516001600160a01b03851681529182015283917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c93869101611a91565b8161272081612dfc565b612728612e51565b6001600160a01b03821660009081526004602090815260408083208684529091529020805460ff161561275b5750505050565b600084815260056020526040902060020180546063190161278f576040516305a48e0f60e01b815260040160405180910390fd5b8154600160ff199091168117835581549081018255600082815260209081902090910180546001600160a01b0319166001600160a01b03871690811790915560405190815286917f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e191015b60405180910390a25050505050565b60008160405160200161281c919061549b565b604051602081830303815290604052805190602001209050919050565b8161284381612dfc565b61284b612e51565b61285483611410565b1561287257604051631685ecdd60e31b815260040160405180910390fd5b6001600160a01b038216600090815260046020908152604080832086845290915290205460ff166128c8576040516379bfd40160e01b8152600481018490526001600160a01b0383166024820152604401610ade565b6000838152600560205260408120600201805490915b818110156129dd57846001600160a01b031683828154811061290257612902615157565b6000918252602090912001546001600160a01b0316036129cd5782612928600184615183565b8154811061293857612938615157565b9060005260206000200160009054906101000a90046001600160a01b031683828154811061296857612968615157565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b03160217905550828054806129a6576129a6615196565b600082815260209020810160001990810180546001600160a01b03191690550190556129dd565b6129d6816151ac565b90506128de565b506001600160a01b0384166000818152600460209081526040808320898452825291829020805460ff19169055905191825286917f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a791016127fa565b600e8181548110612a4957600080fd5b600091825260209091200154905081565b81612a6481612dfc565b612a6c612e51565b600083815260056020526040902060018101546001600160a01b03848116911614612aef576001810180546001600160a01b0319166001600160a01b03851690811790915560408051338152602081019290925285917f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a191015b60405180910390a25b50505050565b600081815260056020526040812054819081906001600160a01b03166060612b1c82613787565b600086815260066020908152604080832054600583529281902060020180548251818502810185019093528083526001600160601b0380861695600160601b810490911694600160c01b9091046001600160401b0316938893929091839190830182828015612bb457602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612b96575b505050505090509450945094509450945091939590929450565b612bd6612da0565b6002546001600160a01b0316612bff5760405163c1f0c0a160e01b815260040160405180910390fd5b6002546040516370a0823160e01b81523060048201526000916001600160a01b0316906370a0823190602401602060405180830381865afa158015612c48573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c6c91906154aa565b600a549091506001600160601b031681811115612ca6576040516354ced18160e11b81526004810182905260248101839052604401610ade565b81811015610cc6576000612cba8284615183565b60025460405163a9059cbb60e01b81526001600160a01b0387811660048301526024820184905292935091169063a9059cbb906044016020604051808303816000875af1158015612d0f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d33919061537a565b612d5057604051631f01ff1360e21b815260040160405180910390fd5b604080516001600160a01b0386168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b4366009101610c75565b612d97612da0565b610ae78161392c565b6000546001600160a01b03163314612dfa5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610ade565b565b6000818152600560205260409020546001600160a01b0316612e1d81613787565b336001600160a01b03821614610d6d57604051636c51fda960e11b81526001600160a01b0382166004820152602401610ade565b600c546601000000000000900460ff1615612dfa5760405163769dd35360e11b815260040160405180910390fd5b600080612e8b846135d5565b60025491935091506001600160a01b031615801590612eb257506001600160601b03821615155b15612efa5760025460405163a9059cbb60e01b81526001600160a01b0385811660048301526001600160601b0385166024830152612efa92169063a9059cbb906044016112b9565b612f0d83826001600160601b0316612f80565b604080516001600160a01b03851681526001600160601b03808516602083015283169181019190915284907f8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c490606001612ae6565b80610ae757604051631e9acf1760e31b815260040160405180910390fd5b6000826001600160a01b03168260405160006040518083038185875af1925050503d8060008114612fcd576040519150601f19603f3d011682016040523d82523d6000602084013e612fd2565b606091505b5050905080610cc65760405163950b247960e01b815260040160405180910390fd5b6040805160a081018252600060608201818152608083018290528252602082018190529181018290529061302784612809565b6000818152600d602090815260409182902082518084019093525460ff811615158084526101009091046001600160401b0316918301919091529192509061308557604051631dfd6e1360e21b815260048101839052602401610ade565b6000828660c001356040516020016130a7929190918252602082015260400190565b60408051601f1981840301815291815281516020928301206000818152600f90935290822054909250908190036130f157604051631b44092560e11b815260040160405180910390fd5b816130ff60208801886154c3565b602088013561311460608a0160408b016151e5565b61312460808b0160608c016151e5565b61313460a08c0160808d01614d1f565b61314160a08d018d61525f565b6040516020016131589897969594939291906154de565b60405160208183030381529060405280519060200120811461318d5760405163354a450b60e21b815260040160405180910390fd5b60006131a461319f60208901896154c3565b6139d5565b905080613288576001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001663e9413d386131e760208a018a6154c3565b6040516001600160e01b031960e084901b1681526001600160401b039091166004820152602401602060405180830381865afa15801561322b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061324f91906154aa565b9050806132885761326360208801886154c3565b60405163175dadad60e01b81526001600160401b039091166004820152602401610ade565b6040805160c08a01356020808301919091528183018490528251808303840181526060909201909252805191012060006132c28a83613a8d565b604080516060810182529788526020880196909652948601949094525092979650505050505050565b6000816001600160401b03163a111561333e57821561331457506001600160401b038116612621565b60405163435e532d60e11b81523a60048201526001600160401b0383166024820152604401610ade565b503a92915050565b6000806000631fe543e360e01b8685604051602401613366929190615556565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252600c805466ff0000000000001916660100000000000017905591506133e3906133c790606088019088016151e5565b63ffffffff166133dd60a0880160808901614d1f565b83613b83565b600c805466ff000000000000191690559695505050505050565b600080831561341c57613411868685613bcf565b60009150915061342c565b613427868685613ce5565b915091505b94509492505050565b600081815260066020526040902082156134e55780546001600160601b03600160601b90910481169061346c908616821015612f62565b61347685826151c5565b82546bffffffffffffffffffffffff60601b1916600160601b6001600160601b039283168102919091178455600b805488939192600c926134bb928692900416615397565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555050612aef565b80546001600160601b0390811690613501908616821015612f62565b61350b85826151c5565b82546bffffffffffffffffffffffff19166001600160601b03918216178355600b8054879260009161353f91859116615397565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505050505050565b601154600090815b818110156135cb57836001600160a01b03166011828154811061359757613597615157565b6000918252602090912001546001600160a01b0316036135bb575060019392505050565b6135c4816151ac565b9050613572565b5060009392505050565b60008181526005602090815260408083206006909252822054600290910180546001600160601b0380841694600160601b90940416925b81811015613681576004600084838154811061362a5761362a615157565b60009182526020808320909101546001600160a01b0316835282810193909352604091820181208982529092529020805470ffffffffffffffffffffffffffffffffff1916905561367a816151ac565b905061360c565b50600085815260056020526040812080546001600160a01b031990811682556001820180549091169055906136b96002830182614c7d565b50506000858152600660205260408120556136d5600886613edc565b506001600160601b0384161561372857600a80548591906000906137039084906001600160601b03166151c5565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b6001600160601b038316156137805782600a600c8282829054906101000a90046001600160601b031661375b91906151c5565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b5050915091565b6001600160a01b038116610ae757604051630fb532db60e11b815260040160405180910390fd5b60408051602081019091526000815260008290036137db5750604080516020810190915260008152612621565b63125fa26760e31b6137ed838561556f565b6001600160e01b0319161461381557604051632923fee760e11b815260040160405180910390fd5b613822826004818661559f565b8101906110bb91906155c9565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa8260405160240161386891511515815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b03199093169290921790915292915050565b600060646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156138e1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061390591906154aa565b905090565b60006110bb8383613ee8565b6000612621825490565b60006110bb8383613f37565b336001600160a01b038216036139845760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610ade565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000806139e06138a0565b9050806001600160401b0316836001600160401b0316101580613a165750610100613a0b8483615622565b6001600160401b0316115b15613a245750600092915050565b6040516315a03d4160e11b81526001600160401b0384166004820152606490632b407a8290602401602060405180830381865afa158015613a69573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110bb91906154aa565b604080518082018252600091613b4d9190859060029083908390808284376000920191909152505060408051808201825291508087019060029083908390808284376000920191909152505050608086013560a087013586613af66101008a0160e08b01614d1f565b604080518082018252906101008c019060029083908390808284376000920191909152505060408051808201825291506101408d0190600290839083908082843760009201919091525050506101808c0135613f61565b600383604001604051602001613b64929190615642565b60408051601f1981840301815291905280516020909101209392505050565b60005a611388811015613b9557600080fd5b611388810390508460408204820311613bad57600080fd5b50823b613bb957600080fd5b60008083516020850160008789f1949350505050565b600080613bdd60003661418c565b905060005a600c54613bfd908890600160581b900463ffffffff166153b7565b613c079190615183565b613c119086615658565b600c54909150600090613c3690600160781b900463ffffffff1664e8d4a51000615658565b90508215613c72576040518381527f56296f7beae05a0db815737fdb4cd298897b1e517614d62468081531ae14d0999060200160405180910390a15b8415613cbc57600c548190606490600160b81b900460ff16613c9485876153b7565b613c9e9190615658565b613ca89190615685565b613cb291906153b7565b93505050506110bb565b600c548190606490613cd890600160b81b900460ff1682615699565b60ff16613c9485876153b7565b600080600080613cf3614196565b9150915060008213613d1b576040516321ea67b360e11b815260048101839052602401610ade565b6000613d2860003661418c565b9050600083825a600c54613d4a908d90600160581b900463ffffffff166153b7565b613d549190615183565b613d5e908b615658565b613d6891906153b7565b613d7a90670de0b6b3a7640000615658565b613d849190615685565b600c54909150600090613dad9063ffffffff600160981b8204811691600160781b9004166156b2565b613dc29063ffffffff1664e8d4a51000615658565b9050600085613dd983670de0b6b3a7640000615658565b613de39190615685565b90508315613e1f576040518481527f56296f7beae05a0db815737fdb4cd298897b1e517614d62468081531ae14d0999060200160405180910390a15b60008915613e5e57600c548290606490613e4390600160c01b900460ff1687615658565b613e4d9190615685565b613e5791906153b7565b9050613e9e565b600c548290606490613e7a90600160c01b900460ff1682615699565b613e879060ff1687615658565b613e919190615685565b613e9b91906153b7565b90505b6b033b2e3c9fd0803ce8000000811115613ecb5760405163e80fa38160e01b815260040160405180910390fd5b9b949a509398505050505050505050565b60006110bb8383614261565b6000818152600183016020526040812054613f2f57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155612621565b506000612621565b6000826000018281548110613f4e57613f4e615157565b9060005260206000200154905092915050565b613f6a8961435b565b613fb65760405162461bcd60e51b815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e2063757276650000000000006044820152606401610ade565b613fbf8861435b565b61400b5760405162461bcd60e51b815260206004820152601560248201527f67616d6d61206973206e6f74206f6e20637572766500000000000000000000006044820152606401610ade565b6140148361435b565b6140605760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610ade565b6140698261435b565b6140b55760405162461bcd60e51b815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e206375727665000000006044820152606401610ade565b6140c1878a8887614434565b61410d5760405162461bcd60e51b815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e657373000000000000006044820152606401610ade565b60006141198a87614557565b9050600061412c898b878b8689896145bb565b9050600061413d838d8d8a866146e7565b9050808a1461417e5760405162461bcd60e51b815260206004820152600d60248201526c34b73b30b634b210383937b7b360991b6044820152606401610ade565b505050505050505050505050565b60006110bb614727565b600c5460035460408051633fabe5a360e21b81529051600093849367010000000000000090910463ffffffff169284926001600160a01b039092169163feaf968c9160048082019260a0929091908290030181865afa1580156141fd573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061422191906156e9565b50919650909250505063ffffffff82161580159061424d57506142448142615183565b8263ffffffff16105b9250821561425b5760105493505b50509091565b6000818152600183016020526040812054801561434a576000614285600183615183565b855490915060009061429990600190615183565b90508181146142fe5760008660000182815481106142b9576142b9615157565b90600052602060002001549050808760000184815481106142dc576142dc615157565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061430f5761430f615196565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050612621565b6000915050612621565b5092915050565b80516000906401000003d019116143b45760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420782d6f7264696e61746500000000000000000000000000006044820152606401610ade565b60208201516401000003d0191161440d5760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420792d6f7264696e61746500000000000000000000000000006044820152606401610ade565b60208201516401000003d01990800961442d8360005b6020020151614768565b1492915050565b60006001600160a01b03821661447a5760405162461bcd60e51b815260206004820152600b60248201526a626164207769746e65737360a81b6044820152606401610ade565b60208401516000906001161561449157601c614494565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe19918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa15801561452f573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b61455f614c9b565b61458c600184846040516020016145789392919061575c565b60405160208183030381529060405261478c565b90505b6145988161435b565b6126215780516040805160208101929092526145b49101614578565b905061458f565b6145c3614c9b565b825186516401000003d01991829006919006036146225760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610ade565b61462d8789886147d9565b6146795760405162461bcd60e51b815260206004820152601660248201527f4669727374206d756c20636865636b206661696c6564000000000000000000006044820152606401610ade565b6146848486856147d9565b6146d05760405162461bcd60e51b815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c65640000000000000000006044820152606401610ade565b6146db868484614904565b98975050505050505050565b6000600286868685876040516020016147059695949392919061577d565b60408051601f1981840301815291905280516020909101209695505050505050565b6000606c6001600160a01b031663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156138e1573d6000803e3d6000fd5b6000806401000003d01980848509840990506401000003d019600782089392505050565b614794614c9b565b61479d826149cb565b81526147b26147ad826000614423565b614a06565b6020820181905260029006600103611eca576020810180516401000003d019039052919050565b6000826000036148195760405162461bcd60e51b815260206004820152600b60248201526a3d32b9379039b1b0b630b960a91b6044820152606401610ade565b8351602085015160009061482f906002906157dc565b1561483b57601c61483e565b601b5b9050600070014551231950b75fc4402da1732fc9bebe198387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa1580156148b0573d6000803e3d6000fd5b5050506020604051035190506000866040516020016148cf91906157f0565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b61490c614c9b565b83516020808601518551918601516000938493849361492d93909190614a26565b919450925090506401000003d01985820960011461498d5760405162461bcd60e51b815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a000000000000006044820152606401610ade565b60405180604001604052806401000003d019806149ac576149ac61566f565b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d0198110611eca576040805160208082019390935281518082038401815290820190915280519101206149d3565b6000612621826002614a1f6401000003d01960016153b7565b901c614b06565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000614a6683838585614bab565b9098509050614a7788828e88614bcf565b9098509050614a8888828c87614bcf565b90985090506000614a9b8d878b85614bcf565b9098509050614aac88828686614bab565b9098509050614abd88828e89614bcf565b9098509050818114614af2576401000003d019818a0998506401000003d01982890997506401000003d0198183099650614af6565b8196505b5050505050509450945094915050565b600080614b11614cb9565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152614b43614cd7565b60208160c0846005600019fa925082600003614ba15760405162461bcd60e51b815260206004820152601260248201527f6269674d6f64457870206661696c7572652100000000000000000000000000006044820152606401610ade565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215614c6d579160200282015b82811115614c6d57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190614c38565b50614c79929150614cf5565b5090565b5080546000825590600052602060002090810190610ae79190614cf5565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b80821115614c795760008155600101614cf6565b6001600160a01b0381168114610ae757600080fd5b600060208284031215614d3157600080fd5b81356110bb81614d0a565b806040810183101561262157600080fd5b600060408284031215614d5f57600080fd5b6110bb8383614d3c565b60008060408385031215614d7c57600080fd5b823591506020830135614d8e81614d0a565b809150509250929050565b600060c08284031215614dab57600080fd5b50919050565b8015158114610ae757600080fd5b60008060008385036101e0811215614dd657600080fd5b6101a080821215614de657600080fd5b85945084013590506001600160401b03811115614e0257600080fd5b614e0e86828701614d99565b9250506101c0840135614e2081614db1565b809150509250925092565b600060208284031215614e3d57600080fd5b5035919050565b60008060408385031215614e5757600080fd5b8235614e6281614d0a565b91506020830135614d8e81614d0a565b80356001600160401b0381168114611eca57600080fd5b60008060608385031215614e9c57600080fd5b614ea68484614d3c565b9150614eb460408401614e72565b90509250929050565b600060208284031215614ecf57600080fd5b81356001600160401b03811115614ee557600080fd5b614ef184828501614d99565b949350505050565b60008060008060608587031215614f0f57600080fd5b8435614f1a81614d0a565b93506020850135925060408501356001600160401b0380821115614f3d57600080fd5b818701915087601f830112614f5157600080fd5b813581811115614f6057600080fd5b886020828501011115614f7257600080fd5b95989497505060200194505050565b803561ffff81168114611eca57600080fd5b803563ffffffff81168114611eca57600080fd5b803560ff81168114611eca57600080fd5b60008060008060008060008060006101208a8c031215614fd757600080fd5b614fe08a614f81565b9850614fee60208b01614f93565b9750614ffc60408b01614f93565b965061500a60608b01614f93565b955060808a0135945061501f60a08b01614f93565b935061502d60c08b01614f93565b925061503b60e08b01614fa7565b915061504a6101008b01614fa7565b90509295985092959850929598565b6000806040838503121561506c57600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b838110156150ab5781518752958201959082019060010161508f565b509495945050505050565b6020815260006110bb602083018461507b565b600081518084526020808501945080840160005b838110156150ab5781516001600160a01b0316875295820195908201906001016150dd565b60006001600160601b0380881683528087166020840152506001600160401b03851660408301526001600160a01b038416606083015260a0608083015261514c60a08301846150c9565b979650505050505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b818103818111156126215761262161516d565b634e487b7160e01b600052603160045260246000fd5b6000600182016151be576151be61516d565b5060010190565b6001600160601b038281168282160390808211156143545761435461516d565b6000602082840312156151f757600080fd5b6110bb82614f93565b634e487b7160e01b600052604160045260246000fd5b60006001600160401b038083168181036152325761523261516d565b6001019392505050565b60006001600160401b038216806152555761525561516d565b6000190192915050565b6000808335601e1984360301811261527657600080fd5b8301803591506001600160401b0382111561529057600080fd5b6020019150368190038213156152a557600080fd5b9250929050565b6020815260ff8251166020820152602082015160408201526001600160a01b0360408301511660608201526000606083015160c060808401526152f260e08401826150c9565b905060808401516001600160601b0380821660a08601528060a08701511660c086015250508091505092915050565b6000815180845260005b818110156153475760208185018101518683018201520161532b565b506000602082860101526020601f19601f83011685010191505092915050565b6020815260006110bb6020830184615321565b60006020828403121561538c57600080fd5b81516110bb81614db1565b6001600160601b038181168382160190808211156143545761435461516d565b808201808211156126215761262161516d565b6000602082840312156153dc57600080fd5b6110bb82614f81565b878152866020820152856040820152600063ffffffff80871660608401528086166080840152506001600160a01b03841660a083015260e060c083015261542f60e0830184615321565b9998505050505050505050565b86815285602082015261ffff85166040820152600063ffffffff808616606084015280851660808401525060c060a08301526146db60c0830184615321565b6001600160401b038181168382160190808211156143545761435461516d565b60408181019083833792915050565b6000602082840312156154bc57600080fd5b5051919050565b6000602082840312156154d557600080fd5b6110bb82614e72565b8881526001600160401b0388166020820152866040820152600063ffffffff80881660608401528087166080840152506001600160a01b03851660a083015260e060c08301528260e08301526101008385828501376000838501820152601f909301601f191690910190910198975050505050505050565b828152604060208201526000614ef1604083018461507b565b6001600160e01b031981358181169160048510156155975780818660040360031b1b83161692505b505092915050565b600080858511156155af57600080fd5b838611156155bc57600080fd5b5050820193919092039150565b6000602082840312156155db57600080fd5b604051602081018181106001600160401b038211171561560b57634e487b7160e01b600052604160045260246000fd5b604052823561561981614db1565b81529392505050565b6001600160401b038281168282160390808211156143545761435461516d565b8281526060810160408360208401379392505050565b80820281158282048414176126215761262161516d565b634e487b7160e01b600052601260045260246000fd5b6000826156945761569461566f565b500490565b60ff81811683821601908111156126215761262161516d565b63ffffffff8281168282160390808211156143545761435461516d565b805169ffffffffffffffffffff81168114611eca57600080fd5b600080600080600060a0868803121561570157600080fd5b61570a866156cf565b945060208601519350604086015192506060860151915061572d608087016156cf565b90509295509295909350565b8060005b6002811015612aef57815184526020938401939091019060010161573d565b83815261576c6020820184615739565b606081019190915260800192915050565b86815261578d6020820187615739565b61579a6060820186615739565b6157a760a0820185615739565b6157b460e0820184615739565b60609190911b6bffffffffffffffffffffffff19166101208201526101340195945050505050565b6000826157eb576157eb61566f565b500690565b6157fa8183615739565b60400191905056fea164736f6c6343000813000a",
}

var VRFCoordinatorV25ArbitrumABI = VRFCoordinatorV25ArbitrumMetaData.ABI

var VRFCoordinatorV25ArbitrumBin = VRFCoordinatorV25ArbitrumMetaData.Bin

func DeployVRFCoordinatorV25Arbitrum(auth *bind.TransactOpts, backend bind.ContractBackend, blockhashStore common.Address) (common.Address, *types.Transaction, *VRFCoordinatorV25Arbitrum, error) {
	parsed, err := VRFCoordinatorV25ArbitrumMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorV25ArbitrumBin), backend, blockhashStore)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorV25Arbitrum{address: address, abi: *parsed, VRFCoordinatorV25ArbitrumCaller: VRFCoordinatorV25ArbitrumCaller{contract: contract}, VRFCoordinatorV25ArbitrumTransactor: VRFCoordinatorV25ArbitrumTransactor{contract: contract}, VRFCoordinatorV25ArbitrumFilterer: VRFCoordinatorV25ArbitrumFilterer{contract: contract}}, nil
}

type VRFCoordinatorV25Arbitrum struct {
	address common.Address
	abi     abi.ABI
	VRFCoordinatorV25ArbitrumCaller
	VRFCoordinatorV25ArbitrumTransactor
	VRFCoordinatorV25ArbitrumFilterer
}

type VRFCoordinatorV25ArbitrumCaller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV25ArbitrumTransactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV25ArbitrumFilterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV25ArbitrumSession struct {
	Contract     *VRFCoordinatorV25Arbitrum
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV25ArbitrumCallerSession struct {
	Contract *VRFCoordinatorV25ArbitrumCaller
	CallOpts bind.CallOpts
}

type VRFCoordinatorV25ArbitrumTransactorSession struct {
	Contract     *VRFCoordinatorV25ArbitrumTransactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV25ArbitrumRaw struct {
	Contract *VRFCoordinatorV25Arbitrum
}

type VRFCoordinatorV25ArbitrumCallerRaw struct {
	Contract *VRFCoordinatorV25ArbitrumCaller
}

type VRFCoordinatorV25ArbitrumTransactorRaw struct {
	Contract *VRFCoordinatorV25ArbitrumTransactor
}

func NewVRFCoordinatorV25Arbitrum(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorV25Arbitrum, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorV25ArbitrumABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFCoordinatorV25Arbitrum(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25Arbitrum{address: address, abi: abi, VRFCoordinatorV25ArbitrumCaller: VRFCoordinatorV25ArbitrumCaller{contract: contract}, VRFCoordinatorV25ArbitrumTransactor: VRFCoordinatorV25ArbitrumTransactor{contract: contract}, VRFCoordinatorV25ArbitrumFilterer: VRFCoordinatorV25ArbitrumFilterer{contract: contract}}, nil
}

func NewVRFCoordinatorV25ArbitrumCaller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorV25ArbitrumCaller, error) {
	contract, err := bindVRFCoordinatorV25Arbitrum(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumCaller{contract: contract}, nil
}

func NewVRFCoordinatorV25ArbitrumTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorV25ArbitrumTransactor, error) {
	contract, err := bindVRFCoordinatorV25Arbitrum(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumTransactor{contract: contract}, nil
}

func NewVRFCoordinatorV25ArbitrumFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorV25ArbitrumFilterer, error) {
	contract, err := bindVRFCoordinatorV25Arbitrum(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumFilterer{contract: contract}, nil
}

func bindVRFCoordinatorV25Arbitrum(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFCoordinatorV25ArbitrumMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV25Arbitrum.Contract.VRFCoordinatorV25ArbitrumCaller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.VRFCoordinatorV25ArbitrumTransactor.contract.Transfer(opts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.VRFCoordinatorV25ArbitrumTransactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV25Arbitrum.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.contract.Transfer(opts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "BLOCKHASH_STORE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) LINK() (common.Address, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.LINK(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) LINK() (common.Address, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.LINK(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) LINKNATIVEFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "LINK_NATIVE_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) LINKNATIVEFEED() (common.Address, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.LINKNATIVEFEED(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) LINKNATIVEFEED() (common.Address, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.LINKNATIVEFEED(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.MAXCONSUMERS(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.MAXCONSUMERS(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) MAXNUMWORDS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "MAX_NUM_WORDS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.MAXNUMWORDS(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.MAXNUMWORDS(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "MAX_REQUEST_CONFIRMATIONS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) GetActiveSubscriptionIds(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "getActiveSubscriptionIds", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) GetActiveSubscriptionIds(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.GetActiveSubscriptionIds(&_VRFCoordinatorV25Arbitrum.CallOpts, startIndex, maxCount)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) GetActiveSubscriptionIds(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.GetActiveSubscriptionIds(&_VRFCoordinatorV25Arbitrum.CallOpts, startIndex, maxCount)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "getSubscription", subId)

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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinatorV25Arbitrum.Contract.GetSubscription(&_VRFCoordinatorV25Arbitrum.CallOpts, subId)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinatorV25Arbitrum.Contract.GetSubscription(&_VRFCoordinatorV25Arbitrum.CallOpts, subId)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "hashOfKey", publicKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.HashOfKey(&_VRFCoordinatorV25Arbitrum.CallOpts, publicKey)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.HashOfKey(&_VRFCoordinatorV25Arbitrum.CallOpts, publicKey)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) Owner() (common.Address, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.Owner(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) Owner() (common.Address, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.Owner(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) PendingRequestExists(opts *bind.CallOpts, subId *big.Int) (bool, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "pendingRequestExists", subId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) PendingRequestExists(subId *big.Int) (bool, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.PendingRequestExists(&_VRFCoordinatorV25Arbitrum.CallOpts, subId)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) PendingRequestExists(subId *big.Int) (bool, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.PendingRequestExists(&_VRFCoordinatorV25Arbitrum.CallOpts, subId)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) SConfig(opts *bind.CallOpts) (SConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "s_config")

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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) SConfig() (SConfig,

	error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SConfig(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) SConfig() (SConfig,

	error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SConfig(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) SCurrentSubNonce(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "s_currentSubNonce")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) SCurrentSubNonce() (uint64, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SCurrentSubNonce(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) SCurrentSubNonce() (uint64, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SCurrentSubNonce(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) SFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "s_fallbackWeiPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SFallbackWeiPerUnitLink(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SFallbackWeiPerUnitLink(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) SProvingKeyHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "s_provingKeyHashes", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) SProvingKeyHashes(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SProvingKeyHashes(&_VRFCoordinatorV25Arbitrum.CallOpts, arg0)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) SProvingKeyHashes(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SProvingKeyHashes(&_VRFCoordinatorV25Arbitrum.CallOpts, arg0)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) SProvingKeys(opts *bind.CallOpts, arg0 [32]byte) (SProvingKeys,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "s_provingKeys", arg0)

	outstruct := new(SProvingKeys)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.MaxGas = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) SProvingKeys(arg0 [32]byte) (SProvingKeys,

	error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SProvingKeys(&_VRFCoordinatorV25Arbitrum.CallOpts, arg0)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) SProvingKeys(arg0 [32]byte) (SProvingKeys,

	error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SProvingKeys(&_VRFCoordinatorV25Arbitrum.CallOpts, arg0)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) SRequestCommitments(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "s_requestCommitments", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) SRequestCommitments(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SRequestCommitments(&_VRFCoordinatorV25Arbitrum.CallOpts, arg0)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) SRequestCommitments(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SRequestCommitments(&_VRFCoordinatorV25Arbitrum.CallOpts, arg0)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) STotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "s_totalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) STotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.STotalBalance(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) STotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.STotalBalance(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCaller) STotalNativeBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV25Arbitrum.contract.Call(opts, &out, "s_totalNativeBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) STotalNativeBalance() (*big.Int, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.STotalNativeBalance(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumCallerSession) STotalNativeBalance() (*big.Int, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.STotalNativeBalance(&_VRFCoordinatorV25Arbitrum.CallOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "acceptOwnership")
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.AcceptOwnership(&_VRFCoordinatorV25Arbitrum.TransactOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.AcceptOwnership(&_VRFCoordinatorV25Arbitrum.TransactOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "addConsumer", subId, consumer)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.AddConsumer(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.AddConsumer(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "cancelSubscription", subId, to)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.CancelSubscription(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId, to)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.CancelSubscription(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId, to)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "createSubscription")
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.CreateSubscription(&_VRFCoordinatorV25Arbitrum.TransactOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.CreateSubscription(&_VRFCoordinatorV25Arbitrum.TransactOpts)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) DeregisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "deregisterMigratableCoordinator", target)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) DeregisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.DeregisterMigratableCoordinator(&_VRFCoordinatorV25Arbitrum.TransactOpts, target)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) DeregisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.DeregisterMigratableCoordinator(&_VRFCoordinatorV25Arbitrum.TransactOpts, target)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "deregisterProvingKey", publicProvingKey)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.DeregisterProvingKey(&_VRFCoordinatorV25Arbitrum.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.DeregisterProvingKey(&_VRFCoordinatorV25Arbitrum.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "fulfillRandomWords", proof, rc, onlyPremium)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) FulfillRandomWords(proof VRFProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.FulfillRandomWords(&_VRFCoordinatorV25Arbitrum.TransactOpts, proof, rc, onlyPremium)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) FulfillRandomWords(proof VRFProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.FulfillRandomWords(&_VRFCoordinatorV25Arbitrum.TransactOpts, proof, rc, onlyPremium)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) FundSubscriptionWithNative(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "fundSubscriptionWithNative", subId)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) FundSubscriptionWithNative(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.FundSubscriptionWithNative(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) FundSubscriptionWithNative(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.FundSubscriptionWithNative(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) Migrate(opts *bind.TransactOpts, subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "migrate", subId, newCoordinator)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) Migrate(subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.Migrate(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId, newCoordinator)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) Migrate(subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.Migrate(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId, newCoordinator)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.OnTokenTransfer(&_VRFCoordinatorV25Arbitrum.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.OnTokenTransfer(&_VRFCoordinatorV25Arbitrum.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) OwnerCancelSubscription(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "ownerCancelSubscription", subId)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) OwnerCancelSubscription(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.OwnerCancelSubscription(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) OwnerCancelSubscription(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.OwnerCancelSubscription(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "recoverFunds", to)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RecoverFunds(&_VRFCoordinatorV25Arbitrum.TransactOpts, to)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RecoverFunds(&_VRFCoordinatorV25Arbitrum.TransactOpts, to)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) RecoverNativeFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "recoverNativeFunds", to)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) RecoverNativeFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RecoverNativeFunds(&_VRFCoordinatorV25Arbitrum.TransactOpts, to)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) RecoverNativeFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RecoverNativeFunds(&_VRFCoordinatorV25Arbitrum.TransactOpts, to)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) RegisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "registerMigratableCoordinator", target)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) RegisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RegisterMigratableCoordinator(&_VRFCoordinatorV25Arbitrum.TransactOpts, target)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) RegisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RegisterMigratableCoordinator(&_VRFCoordinatorV25Arbitrum.TransactOpts, target)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) RegisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int, maxGas uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "registerProvingKey", publicProvingKey, maxGas)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) RegisterProvingKey(publicProvingKey [2]*big.Int, maxGas uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RegisterProvingKey(&_VRFCoordinatorV25Arbitrum.TransactOpts, publicProvingKey, maxGas)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) RegisterProvingKey(publicProvingKey [2]*big.Int, maxGas uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RegisterProvingKey(&_VRFCoordinatorV25Arbitrum.TransactOpts, publicProvingKey, maxGas)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "removeConsumer", subId, consumer)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RemoveConsumer(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RemoveConsumer(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) RequestRandomWords(opts *bind.TransactOpts, req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "requestRandomWords", req)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RequestRandomWords(&_VRFCoordinatorV25Arbitrum.TransactOpts, req)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RequestRandomWords(&_VRFCoordinatorV25Arbitrum.TransactOpts, req)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV25Arbitrum.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeNativePPM uint32, fulfillmentFlatFeeLinkDiscountPPM uint32, nativePremiumPercentage uint8, linkPremiumPercentage uint8) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fulfillmentFlatFeeNativePPM, fulfillmentFlatFeeLinkDiscountPPM, nativePremiumPercentage, linkPremiumPercentage)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeNativePPM uint32, fulfillmentFlatFeeLinkDiscountPPM uint32, nativePremiumPercentage uint8, linkPremiumPercentage uint8) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SetConfig(&_VRFCoordinatorV25Arbitrum.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fulfillmentFlatFeeNativePPM, fulfillmentFlatFeeLinkDiscountPPM, nativePremiumPercentage, linkPremiumPercentage)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeNativePPM uint32, fulfillmentFlatFeeLinkDiscountPPM uint32, nativePremiumPercentage uint8, linkPremiumPercentage uint8) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SetConfig(&_VRFCoordinatorV25Arbitrum.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fulfillmentFlatFeeNativePPM, fulfillmentFlatFeeLinkDiscountPPM, nativePremiumPercentage, linkPremiumPercentage)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) SetLINKAndLINKNativeFeed(opts *bind.TransactOpts, link common.Address, linkNativeFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "setLINKAndLINKNativeFeed", link, linkNativeFeed)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) SetLINKAndLINKNativeFeed(link common.Address, linkNativeFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SetLINKAndLINKNativeFeed(&_VRFCoordinatorV25Arbitrum.TransactOpts, link, linkNativeFeed)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) SetLINKAndLINKNativeFeed(link common.Address, linkNativeFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.SetLINKAndLINKNativeFeed(&_VRFCoordinatorV25Arbitrum.TransactOpts, link, linkNativeFeed)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.TransferOwnership(&_VRFCoordinatorV25Arbitrum.TransactOpts, to)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.TransferOwnership(&_VRFCoordinatorV25Arbitrum.TransactOpts, to)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) Withdraw(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "withdraw", recipient)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) Withdraw(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.Withdraw(&_VRFCoordinatorV25Arbitrum.TransactOpts, recipient)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) Withdraw(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.Withdraw(&_VRFCoordinatorV25Arbitrum.TransactOpts, recipient)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactor) WithdrawNative(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.contract.Transact(opts, "withdrawNative", recipient)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumSession) WithdrawNative(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.WithdrawNative(&_VRFCoordinatorV25Arbitrum.TransactOpts, recipient)
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumTransactorSession) WithdrawNative(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25Arbitrum.Contract.WithdrawNative(&_VRFCoordinatorV25Arbitrum.TransactOpts, recipient)
}

type VRFCoordinatorV25ArbitrumConfigSetIterator struct {
	Event *VRFCoordinatorV25ArbitrumConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumConfigSet)
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
		it.Event = new(VRFCoordinatorV25ArbitrumConfigSet)
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

func (it *VRFCoordinatorV25ArbitrumConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumConfigSet struct {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumConfigSetIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumConfigSet)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseConfigSet(log types.Log) (*VRFCoordinatorV25ArbitrumConfigSet, error) {
	event := new(VRFCoordinatorV25ArbitrumConfigSet)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumCoordinatorDeregisteredIterator struct {
	Event *VRFCoordinatorV25ArbitrumCoordinatorDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumCoordinatorDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumCoordinatorDeregistered)
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
		it.Event = new(VRFCoordinatorV25ArbitrumCoordinatorDeregistered)
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

func (it *VRFCoordinatorV25ArbitrumCoordinatorDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumCoordinatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumCoordinatorDeregistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumCoordinatorDeregisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumCoordinatorDeregisteredIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "CoordinatorDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumCoordinatorDeregistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "CoordinatorDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumCoordinatorDeregistered)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseCoordinatorDeregistered(log types.Log) (*VRFCoordinatorV25ArbitrumCoordinatorDeregistered, error) {
	event := new(VRFCoordinatorV25ArbitrumCoordinatorDeregistered)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "CoordinatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumCoordinatorRegisteredIterator struct {
	Event *VRFCoordinatorV25ArbitrumCoordinatorRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumCoordinatorRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumCoordinatorRegistered)
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
		it.Event = new(VRFCoordinatorV25ArbitrumCoordinatorRegistered)
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

func (it *VRFCoordinatorV25ArbitrumCoordinatorRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumCoordinatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumCoordinatorRegistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumCoordinatorRegisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumCoordinatorRegisteredIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "CoordinatorRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumCoordinatorRegistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumCoordinatorRegistered)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseCoordinatorRegistered(log types.Log) (*VRFCoordinatorV25ArbitrumCoordinatorRegistered, error) {
	event := new(VRFCoordinatorV25ArbitrumCoordinatorRegistered)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsedIterator struct {
	Event *VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsed)
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
		it.Event = new(VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsed)
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

func (it *VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsed struct {
	RequestId              *big.Int
	FallbackWeiPerUnitLink *big.Int
	Raw                    types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterFallbackWeiPerUnitLinkUsed(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsedIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "FallbackWeiPerUnitLinkUsed")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsedIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "FallbackWeiPerUnitLinkUsed", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchFallbackWeiPerUnitLinkUsed(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsed) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "FallbackWeiPerUnitLinkUsed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsed)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "FallbackWeiPerUnitLinkUsed", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseFallbackWeiPerUnitLinkUsed(log types.Log) (*VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsed, error) {
	event := new(VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsed)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "FallbackWeiPerUnitLinkUsed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumFundsRecoveredIterator struct {
	Event *VRFCoordinatorV25ArbitrumFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumFundsRecovered)
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
		it.Event = new(VRFCoordinatorV25ArbitrumFundsRecovered)
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

func (it *VRFCoordinatorV25ArbitrumFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumFundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumFundsRecoveredIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumFundsRecovered)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseFundsRecovered(log types.Log) (*VRFCoordinatorV25ArbitrumFundsRecovered, error) {
	event := new(VRFCoordinatorV25ArbitrumFundsRecovered)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumL1GasFeeIterator struct {
	Event *VRFCoordinatorV25ArbitrumL1GasFee

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumL1GasFeeIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumL1GasFee)
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
		it.Event = new(VRFCoordinatorV25ArbitrumL1GasFee)
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

func (it *VRFCoordinatorV25ArbitrumL1GasFeeIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumL1GasFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumL1GasFee struct {
	Fee *big.Int
	Raw types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterL1GasFee(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumL1GasFeeIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "L1GasFee")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumL1GasFeeIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "L1GasFee", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchL1GasFee(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumL1GasFee) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "L1GasFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumL1GasFee)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "L1GasFee", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseL1GasFee(log types.Log) (*VRFCoordinatorV25ArbitrumL1GasFee, error) {
	event := new(VRFCoordinatorV25ArbitrumL1GasFee)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "L1GasFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumMigrationCompletedIterator struct {
	Event *VRFCoordinatorV25ArbitrumMigrationCompleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumMigrationCompletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumMigrationCompleted)
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
		it.Event = new(VRFCoordinatorV25ArbitrumMigrationCompleted)
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

func (it *VRFCoordinatorV25ArbitrumMigrationCompletedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumMigrationCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumMigrationCompleted struct {
	NewCoordinator common.Address
	SubId          *big.Int
	Raw            types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterMigrationCompleted(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumMigrationCompletedIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "MigrationCompleted")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumMigrationCompletedIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "MigrationCompleted", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchMigrationCompleted(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumMigrationCompleted) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "MigrationCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumMigrationCompleted)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseMigrationCompleted(log types.Log) (*VRFCoordinatorV25ArbitrumMigrationCompleted, error) {
	event := new(VRFCoordinatorV25ArbitrumMigrationCompleted)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumNativeFundsRecoveredIterator struct {
	Event *VRFCoordinatorV25ArbitrumNativeFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumNativeFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumNativeFundsRecovered)
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
		it.Event = new(VRFCoordinatorV25ArbitrumNativeFundsRecovered)
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

func (it *VRFCoordinatorV25ArbitrumNativeFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumNativeFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumNativeFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterNativeFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumNativeFundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "NativeFundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumNativeFundsRecoveredIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "NativeFundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchNativeFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumNativeFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "NativeFundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumNativeFundsRecovered)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "NativeFundsRecovered", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseNativeFundsRecovered(log types.Log) (*VRFCoordinatorV25ArbitrumNativeFundsRecovered, error) {
	event := new(VRFCoordinatorV25ArbitrumNativeFundsRecovered)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "NativeFundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumOwnershipTransferRequestedIterator struct {
	Event *VRFCoordinatorV25ArbitrumOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumOwnershipTransferRequested)
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
		it.Event = new(VRFCoordinatorV25ArbitrumOwnershipTransferRequested)
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

func (it *VRFCoordinatorV25ArbitrumOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV25ArbitrumOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumOwnershipTransferRequestedIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumOwnershipTransferRequested)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV25ArbitrumOwnershipTransferRequested, error) {
	event := new(VRFCoordinatorV25ArbitrumOwnershipTransferRequested)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumOwnershipTransferredIterator struct {
	Event *VRFCoordinatorV25ArbitrumOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumOwnershipTransferred)
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
		it.Event = new(VRFCoordinatorV25ArbitrumOwnershipTransferred)
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

func (it *VRFCoordinatorV25ArbitrumOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV25ArbitrumOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumOwnershipTransferredIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumOwnershipTransferred)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV25ArbitrumOwnershipTransferred, error) {
	event := new(VRFCoordinatorV25ArbitrumOwnershipTransferred)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumProvingKeyDeregisteredIterator struct {
	Event *VRFCoordinatorV25ArbitrumProvingKeyDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumProvingKeyDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumProvingKeyDeregistered)
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
		it.Event = new(VRFCoordinatorV25ArbitrumProvingKeyDeregistered)
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

func (it *VRFCoordinatorV25ArbitrumProvingKeyDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumProvingKeyDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumProvingKeyDeregistered struct {
	KeyHash [32]byte
	MaxGas  uint64
	Raw     types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterProvingKeyDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumProvingKeyDeregisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "ProvingKeyDeregistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumProvingKeyDeregisteredIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "ProvingKeyDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumProvingKeyDeregistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "ProvingKeyDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumProvingKeyDeregistered)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorV25ArbitrumProvingKeyDeregistered, error) {
	event := new(VRFCoordinatorV25ArbitrumProvingKeyDeregistered)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumProvingKeyRegisteredIterator struct {
	Event *VRFCoordinatorV25ArbitrumProvingKeyRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumProvingKeyRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumProvingKeyRegistered)
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
		it.Event = new(VRFCoordinatorV25ArbitrumProvingKeyRegistered)
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

func (it *VRFCoordinatorV25ArbitrumProvingKeyRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumProvingKeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumProvingKeyRegistered struct {
	KeyHash [32]byte
	MaxGas  uint64
	Raw     types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterProvingKeyRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumProvingKeyRegisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "ProvingKeyRegistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumProvingKeyRegisteredIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "ProvingKeyRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumProvingKeyRegistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "ProvingKeyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumProvingKeyRegistered)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV25ArbitrumProvingKeyRegistered, error) {
	event := new(VRFCoordinatorV25ArbitrumProvingKeyRegistered)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumRandomWordsFulfilledIterator struct {
	Event *VRFCoordinatorV25ArbitrumRandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumRandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumRandomWordsFulfilled)
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
		it.Event = new(VRFCoordinatorV25ArbitrumRandomWordsFulfilled)
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

func (it *VRFCoordinatorV25ArbitrumRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumRandomWordsFulfilled struct {
	RequestId     *big.Int
	OutputSeed    *big.Int
	SubId         *big.Int
	Payment       *big.Int
	NativePayment bool
	Success       bool
	OnlyPremium   bool
	Raw           types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subId []*big.Int) (*VRFCoordinatorV25ArbitrumRandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule, subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumRandomWordsFulfilledIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumRandomWordsFulfilled, requestId []*big.Int, subId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule, subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumRandomWordsFulfilled)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV25ArbitrumRandomWordsFulfilled, error) {
	event := new(VRFCoordinatorV25ArbitrumRandomWordsFulfilled)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumRandomWordsRequestedIterator struct {
	Event *VRFCoordinatorV25ArbitrumRandomWordsRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumRandomWordsRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumRandomWordsRequested)
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
		it.Event = new(VRFCoordinatorV25ArbitrumRandomWordsRequested)
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

func (it *VRFCoordinatorV25ArbitrumRandomWordsRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumRandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumRandomWordsRequested struct {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (*VRFCoordinatorV25ArbitrumRandomWordsRequestedIterator, error) {

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

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumRandomWordsRequestedIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumRandomWordsRequested, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumRandomWordsRequested)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV25ArbitrumRandomWordsRequested, error) {
	event := new(VRFCoordinatorV25ArbitrumRandomWordsRequested)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumSubscriptionCanceledIterator struct {
	Event *VRFCoordinatorV25ArbitrumSubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionCanceled)
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
		it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionCanceled)
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

func (it *VRFCoordinatorV25ArbitrumSubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumSubscriptionCanceled struct {
	SubId        *big.Int
	To           common.Address
	AmountLink   *big.Int
	AmountNative *big.Int
	Raw          types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumSubscriptionCanceledIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionCanceled, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumSubscriptionCanceled)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionCanceled, error) {
	event := new(VRFCoordinatorV25ArbitrumSubscriptionCanceled)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumSubscriptionConsumerAddedIterator struct {
	Event *VRFCoordinatorV25ArbitrumSubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionConsumerAdded)
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
		it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionConsumerAdded)
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

func (it *VRFCoordinatorV25ArbitrumSubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumSubscriptionConsumerAdded struct {
	SubId    *big.Int
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumSubscriptionConsumerAddedIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionConsumerAdded, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumSubscriptionConsumerAdded)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionConsumerAdded, error) {
	event := new(VRFCoordinatorV25ArbitrumSubscriptionConsumerAdded)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumSubscriptionConsumerRemovedIterator struct {
	Event *VRFCoordinatorV25ArbitrumSubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionConsumerRemoved)
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
		it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionConsumerRemoved)
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

func (it *VRFCoordinatorV25ArbitrumSubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumSubscriptionConsumerRemoved struct {
	SubId    *big.Int
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumSubscriptionConsumerRemovedIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionConsumerRemoved, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumSubscriptionConsumerRemoved)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionConsumerRemoved, error) {
	event := new(VRFCoordinatorV25ArbitrumSubscriptionConsumerRemoved)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumSubscriptionCreatedIterator struct {
	Event *VRFCoordinatorV25ArbitrumSubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionCreated)
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
		it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionCreated)
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

func (it *VRFCoordinatorV25ArbitrumSubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumSubscriptionCreated struct {
	SubId *big.Int
	Owner common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumSubscriptionCreatedIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionCreated, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumSubscriptionCreated)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionCreated, error) {
	event := new(VRFCoordinatorV25ArbitrumSubscriptionCreated)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumSubscriptionFundedIterator struct {
	Event *VRFCoordinatorV25ArbitrumSubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionFunded)
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
		it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionFunded)
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

func (it *VRFCoordinatorV25ArbitrumSubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumSubscriptionFunded struct {
	SubId      *big.Int
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumSubscriptionFundedIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionFunded, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumSubscriptionFunded)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionFunded, error) {
	event := new(VRFCoordinatorV25ArbitrumSubscriptionFunded)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumSubscriptionFundedWithNativeIterator struct {
	Event *VRFCoordinatorV25ArbitrumSubscriptionFundedWithNative

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionFundedWithNativeIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionFundedWithNative)
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
		it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionFundedWithNative)
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

func (it *VRFCoordinatorV25ArbitrumSubscriptionFundedWithNativeIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionFundedWithNativeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumSubscriptionFundedWithNative struct {
	SubId            *big.Int
	OldNativeBalance *big.Int
	NewNativeBalance *big.Int
	Raw              types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterSubscriptionFundedWithNative(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionFundedWithNativeIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "SubscriptionFundedWithNative", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumSubscriptionFundedWithNativeIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "SubscriptionFundedWithNative", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchSubscriptionFundedWithNative(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionFundedWithNative, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "SubscriptionFundedWithNative", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumSubscriptionFundedWithNative)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionFundedWithNative", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseSubscriptionFundedWithNative(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionFundedWithNative, error) {
	event := new(VRFCoordinatorV25ArbitrumSubscriptionFundedWithNative)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionFundedWithNative", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequestedIterator struct {
	Event *VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequested)
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
		it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequested)
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

func (it *VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequested struct {
	SubId *big.Int
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequestedIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequested, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequested)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequested, error) {
	event := new(VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequested)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferredIterator struct {
	Event *VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferred)
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
		it.Event = new(VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferred)
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

func (it *VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferred struct {
	SubId *big.Int
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferredIterator{contract: _VRFCoordinatorV25Arbitrum.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferred, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25Arbitrum.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferred)
				if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25ArbitrumFilterer) ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferred, error) {
	event := new(VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferred)
	if err := _VRFCoordinatorV25Arbitrum.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25Arbitrum) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinatorV25Arbitrum.abi.Events["ConfigSet"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseConfigSet(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["CoordinatorDeregistered"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseCoordinatorDeregistered(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["CoordinatorRegistered"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseCoordinatorRegistered(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["FallbackWeiPerUnitLinkUsed"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseFallbackWeiPerUnitLinkUsed(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["FundsRecovered"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseFundsRecovered(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["L1GasFee"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseL1GasFee(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["MigrationCompleted"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseMigrationCompleted(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["NativeFundsRecovered"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseNativeFundsRecovered(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseOwnershipTransferRequested(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["OwnershipTransferred"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseOwnershipTransferred(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["ProvingKeyDeregistered"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseProvingKeyDeregistered(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["ProvingKeyRegistered"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseProvingKeyRegistered(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseRandomWordsFulfilled(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["RandomWordsRequested"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseRandomWordsRequested(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["SubscriptionCanceled"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseSubscriptionCanceled(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["SubscriptionConsumerAdded"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseSubscriptionConsumerAdded(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseSubscriptionConsumerRemoved(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["SubscriptionCreated"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseSubscriptionCreated(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["SubscriptionFunded"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseSubscriptionFunded(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["SubscriptionFundedWithNative"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseSubscriptionFundedWithNative(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseSubscriptionOwnerTransferRequested(log)
	case _VRFCoordinatorV25Arbitrum.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _VRFCoordinatorV25Arbitrum.ParseSubscriptionOwnerTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFCoordinatorV25ArbitrumConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2c6b6b12413678366b05b145c5f00745bdd00e739131ab5de82484a50c9d78b6")
}

func (VRFCoordinatorV25ArbitrumCoordinatorDeregistered) Topic() common.Hash {
	return common.HexToHash("0xf80a1a97fd42251f3c33cda98635e7399253033a6774fe37cd3f650b5282af37")
}

func (VRFCoordinatorV25ArbitrumCoordinatorRegistered) Topic() common.Hash {
	return common.HexToHash("0xb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af01625")
}

func (VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsed) Topic() common.Hash {
	return common.HexToHash("0x6ca648a381f22ead7e37773d934e64885dcf861fbfbb26c40354cbf0c4662d1a")
}

func (VRFCoordinatorV25ArbitrumFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (VRFCoordinatorV25ArbitrumL1GasFee) Topic() common.Hash {
	return common.HexToHash("0x56296f7beae05a0db815737fdb4cd298897b1e517614d62468081531ae14d099")
}

func (VRFCoordinatorV25ArbitrumMigrationCompleted) Topic() common.Hash {
	return common.HexToHash("0xd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be4187")
}

func (VRFCoordinatorV25ArbitrumNativeFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x4aed7c8eed0496c8c19ea2681fcca25741c1602342e38b045d9f1e8e905d2e9c")
}

func (VRFCoordinatorV25ArbitrumOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFCoordinatorV25ArbitrumOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFCoordinatorV25ArbitrumProvingKeyDeregistered) Topic() common.Hash {
	return common.HexToHash("0x9b6868e0eb737bcd72205360baa6bfd0ba4e4819a33ade2db384e8a8025639a5")
}

func (VRFCoordinatorV25ArbitrumProvingKeyRegistered) Topic() common.Hash {
	return common.HexToHash("0x9b911b2c240bfbef3b6a8f7ed6ee321d1258bb2a3fe6becab52ac1cd3210afd3")
}

func (VRFCoordinatorV25ArbitrumRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0xaeb4b4786571e184246d39587f659abf0e26f41f6a3358692250382c0cdb47b7")
}

func (VRFCoordinatorV25ArbitrumRandomWordsRequested) Topic() common.Hash {
	return common.HexToHash("0xeb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e")
}

func (VRFCoordinatorV25ArbitrumSubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0x8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c4")
}

func (VRFCoordinatorV25ArbitrumSubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e1")
}

func (VRFCoordinatorV25ArbitrumSubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a7")
}

func (VRFCoordinatorV25ArbitrumSubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d")
}

func (VRFCoordinatorV25ArbitrumSubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0x1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a")
}

func (VRFCoordinatorV25ArbitrumSubscriptionFundedWithNative) Topic() common.Hash {
	return common.HexToHash("0x7603b205d03651ee812f803fccde89f1012e545a9c99f0abfea9cedd0fd8e902")
}

func (VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a1")
}

func (VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0xd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c9386")
}

func (_VRFCoordinatorV25Arbitrum *VRFCoordinatorV25Arbitrum) Address() common.Address {
	return _VRFCoordinatorV25Arbitrum.address
}

type VRFCoordinatorV25ArbitrumInterface interface {
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

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Withdraw(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error)

	WithdrawNative(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFCoordinatorV25ArbitrumConfigSet, error)

	FilterCoordinatorDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumCoordinatorDeregisteredIterator, error)

	WatchCoordinatorDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumCoordinatorDeregistered) (event.Subscription, error)

	ParseCoordinatorDeregistered(log types.Log) (*VRFCoordinatorV25ArbitrumCoordinatorDeregistered, error)

	FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumCoordinatorRegisteredIterator, error)

	WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumCoordinatorRegistered) (event.Subscription, error)

	ParseCoordinatorRegistered(log types.Log) (*VRFCoordinatorV25ArbitrumCoordinatorRegistered, error)

	FilterFallbackWeiPerUnitLinkUsed(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsedIterator, error)

	WatchFallbackWeiPerUnitLinkUsed(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsed) (event.Subscription, error)

	ParseFallbackWeiPerUnitLinkUsed(log types.Log) (*VRFCoordinatorV25ArbitrumFallbackWeiPerUnitLinkUsed, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumFundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumFundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*VRFCoordinatorV25ArbitrumFundsRecovered, error)

	FilterL1GasFee(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumL1GasFeeIterator, error)

	WatchL1GasFee(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumL1GasFee) (event.Subscription, error)

	ParseL1GasFee(log types.Log) (*VRFCoordinatorV25ArbitrumL1GasFee, error)

	FilterMigrationCompleted(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumMigrationCompletedIterator, error)

	WatchMigrationCompleted(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumMigrationCompleted) (event.Subscription, error)

	ParseMigrationCompleted(log types.Log) (*VRFCoordinatorV25ArbitrumMigrationCompleted, error)

	FilterNativeFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumNativeFundsRecoveredIterator, error)

	WatchNativeFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumNativeFundsRecovered) (event.Subscription, error)

	ParseNativeFundsRecovered(log types.Log) (*VRFCoordinatorV25ArbitrumNativeFundsRecovered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV25ArbitrumOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV25ArbitrumOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV25ArbitrumOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV25ArbitrumOwnershipTransferred, error)

	FilterProvingKeyDeregistered(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumProvingKeyDeregisteredIterator, error)

	WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumProvingKeyDeregistered) (event.Subscription, error)

	ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorV25ArbitrumProvingKeyDeregistered, error)

	FilterProvingKeyRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV25ArbitrumProvingKeyRegisteredIterator, error)

	WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumProvingKeyRegistered) (event.Subscription, error)

	ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV25ArbitrumProvingKeyRegistered, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subId []*big.Int) (*VRFCoordinatorV25ArbitrumRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumRandomWordsFulfilled, requestId []*big.Int, subId []*big.Int) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV25ArbitrumRandomWordsFulfilled, error)

	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (*VRFCoordinatorV25ArbitrumRandomWordsRequestedIterator, error)

	WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumRandomWordsRequested, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV25ArbitrumRandomWordsRequested, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionCanceled, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionConsumerAdded, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionConsumerRemoved, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionCreated, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionFunded, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionFunded, error)

	FilterSubscriptionFundedWithNative(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionFundedWithNativeIterator, error)

	WatchSubscriptionFundedWithNative(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionFundedWithNative, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionFundedWithNative(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionFundedWithNative, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequested, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferred, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV25ArbitrumSubscriptionOwnerTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
