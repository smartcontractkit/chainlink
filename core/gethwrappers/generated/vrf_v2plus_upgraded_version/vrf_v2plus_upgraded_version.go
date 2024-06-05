// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_v2plus_upgraded_version

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

var VRFCoordinatorV2PlusUpgradedVersionMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendNative\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToTransferLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxGas\",\"type\":\"uint256\"}],\"name\":\"GasPriceExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"transferredValue\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"expectedValue\",\"type\":\"uint96\"}],\"name\":\"InvalidNativeBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"premiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"max\",\"type\":\"uint8\"}],\"name\":\"InvalidPremiumPercentage\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"requestVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"expectedVersion\",\"type\":\"uint8\"}],\"name\":\"InvalidVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"flatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"flatFeeNativePPM\",\"type\":\"uint32\"}],\"name\":\"LinkDiscountTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"max\",\"type\":\"uint32\"}],\"name\":\"MsgDataTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SubscriptionIDCollisionFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"}],\"name\":\"FallbackWeiPerUnitLinkUsed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativeFundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"onlyPremium\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountNative\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldNativeBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newNativeBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFundedWithNative\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_NATIVE_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFTypes.RequestCommitmentV2Plus\",\"name\":\"rc\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"onlyPremium\",\"type\":\"bool\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"fundSubscriptionWithNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveSubscriptionIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"subOwner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrationVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedData\",\"type\":\"bytes\"}],\"name\":\"onMigration\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverNativeFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"registerMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"reentrancyLock\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_currentSubNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_provingKeyHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_provingKeys\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"maxGas\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestCommitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalNativeBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkNativeFeed\",\"type\":\"address\"}],\"name\":\"setLINKAndLINKNativeFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162005fad38038062005fad83398101604081905262000034916200017e565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000d3565b5050506001600160a01b0316608052620001b0565b336001600160a01b038216036200012d5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019157600080fd5b81516001600160a01b0381168114620001a957600080fd5b9392505050565b608051615dda620001d36000396000818161052e015261319c0152615dda6000f3fe6080604052600436106101f65760003560e01c8062012291146101fb578063043bd6ae14610228578063088070f51461024c5780630ae095401461031a57806315c48b841461033c57806318e3dd27146103645780631b6b6d23146103a3578063294daa49146103d05780632f622e6b146103ec578063301f42e91461040c578063405b84fa1461042c57806340d6bb821461044c57806341af6c871461047757806351cff8d9146104a75780635d06b4ab146104c757806364d51a2a146104e757806365982744146104fc578063689c45171461051c57806372e9d5651461055057806379ba5097146105705780637a5a2aef146105855780638402595e146105a557806386fe91c7146105c55780638da5cb5b146105e557806395b55cfc146106035780639b1c385e146106165780639d40a6fd14610636578063a21a23e414610663578063a4c0ed3614610678578063a63e0bfb14610698578063aa433aff146106b8578063aefb212f146106d8578063b2a7cac514610705578063bec4c08c14610725578063caf70c4a14610745578063cb63179714610765578063ce3f471914610785578063d98e620e14610798578063da2f2610146107b8578063dac83d2914610817578063dc311dd314610837578063e72f6e3014610868578063ee9d2d3814610888578063f2fde38b146108b5575b600080fd5b34801561020757600080fd5b506102106108d5565b60405161021f93929190614d5a565b60405180910390f35b34801561023457600080fd5b5061023e60105481565b60405190815260200161021f565b34801561025857600080fd5b50600c546102bd9061ffff81169063ffffffff62010000820481169160ff600160301b8204811692600160381b8304811692600160581b8104821692600160781b8204831692600160981b83041691600160b81b8104821691600160c01b9091041689565b6040805161ffff909a168a5263ffffffff98891660208b01529615159689019690965293861660608801529185166080870152841660a08601529290921660c084015260ff91821660e0840152166101008201526101200161021f565b34801561032657600080fd5b5061033a610335366004614dd9565b610951565b005b34801561034857600080fd5b5061035160c881565b60405161ffff909116815260200161021f565b34801561037057600080fd5b50600a5461038b90600160601b90046001600160601b031681565b6040516001600160601b03909116815260200161021f565b3480156103af57600080fd5b506002546103c3906001600160a01b031681565b60405161021f9190614e09565b3480156103dc57600080fd5b506040516002815260200161021f565b3480156103f857600080fd5b5061033a610407366004614e1d565b610999565b34801561041857600080fd5b5061038b61042736600461506c565b610a3a565b34801561043857600080fd5b5061033a610447366004614dd9565b610d54565b34801561045857600080fd5b506104626101f481565b60405163ffffffff909116815260200161021f565b34801561048357600080fd5b5061049761049236600461515a565b6110f7565b604051901515815260200161021f565b3480156104b357600080fd5b5061033a6104c2366004614e1d565b6111ab565b3480156104d357600080fd5b5061033a6104e2366004614e1d565b6112c8565b3480156104f357600080fd5b50610351606481565b34801561050857600080fd5b5061033a610517366004615173565b61137f565b34801561052857600080fd5b506103c37f000000000000000000000000000000000000000000000000000000000000000081565b34801561055c57600080fd5b506003546103c3906001600160a01b031681565b34801561057c57600080fd5b5061033a6113df565b34801561059157600080fd5b5061033a6105a03660046151a1565b611489565b3480156105b157600080fd5b5061033a6105c0366004614e1d565b6115c5565b3480156105d157600080fd5b50600a5461038b906001600160601b031681565b3480156105f157600080fd5b506000546001600160a01b03166103c3565b61033a61061136600461515a565b6116d1565b34801561062257600080fd5b5061023e6106313660046151db565b6117de565b34801561064257600080fd5b50600754610656906001600160401b031681565b60405161021f9190615215565b34801561066f57600080fd5b5061023e611b98565b34801561068457600080fd5b5061033a610693366004615271565b611d6b565b3480156106a457600080fd5b5061033a6106b33660046152ef565b611ed1565b3480156106c457600080fd5b5061033a6106d336600461515a565b612173565b3480156106e457600080fd5b506106f86106f3366004615390565b6121a6565b60405161021f91906153ed565b34801561071157600080fd5b5061033a61072036600461515a565b6122a8565b34801561073157600080fd5b5061033a610740366004614dd9565b612388565b34801561075157600080fd5b5061023e610760366004615400565b61247a565b34801561077157600080fd5b5061033a610780366004614dd9565b6124aa565b61033a61079336600461541c565b61270c565b3480156107a457600080fd5b5061023e6107b336600461515a565b612a73565b3480156107c457600080fd5b506107f86107d336600461515a565b600d6020526000908152604090205460ff81169061010090046001600160401b031682565b6040805192151583526001600160401b0390911660208301520161021f565b34801561082357600080fd5b5061033a610832366004614dd9565b612a94565b34801561084357600080fd5b5061085761085236600461515a565b612b2a565b60405161021f959493929190615496565b34801561087457600080fd5b5061033a610883366004614e1d565b612c03565b34801561089457600080fd5b5061023e6108a336600461515a565b600f6020526000908152604090205481565b3480156108c157600080fd5b5061033a6108d0366004614e1d565b612dc0565b600c54600e805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff1693919283919083018282801561093f57602002820191906000526020600020905b81548152602001906001019080831161092b575b50505050509050925092509250909192565b8161095b81612dd4565b610963612e20565b61096c836110f7565b1561098a57604051631685ecdd60e31b815260040160405180910390fd5b6109948383612e4d565b505050565b6109a1612e20565b6109a9612f2a565b600b54600160601b90046001600160601b03166109c7811515612f7d565b600b8054600160601b600160c01b0319169055600a8054829190600c906109ff908490600160601b90046001600160601b0316615501565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550610a3682826001600160601b0316612f9b565b5050565b6000610a44612e20565b60005a9050610324361115610a7b57604051630f28961b60e01b815236600482015261032460248201526044015b60405180910390fd5b6000610a87868661300f565b90506000610a9d858360000151602001516132a6565b60408301516060888101519293509163ffffffff16806001600160401b03811115610aca57610aca614e3a565b604051908082528060200260200182016040528015610af3578160200160208202803683370190505b50925060005b81811015610b5a578281604051602001610b14929190615521565b6040516020818303038152906040528051906020012060001c848281518110610b3f57610b3f61552f565b6020908102919091010152610b5381615545565b9050610af9565b5050602080850180516000908152600f9092526040822082905551610b80908a856132f4565b60208a8101516000908152600690915260409020805491925090601890610bb690600160c01b90046001600160401b031661555e565b82546101009290920a6001600160401b0381810219909316918316021790915560808a01516001600160a01b03166000908152600460209081526040808320828e01518452909152902080549091600991610c1991600160481b9091041661558c565b91906101000a8154816001600160401b0302191690836001600160401b0316021790555060008960a0015160018b60a0015151610c5691906155af565b81518110610c6657610c6661552f565b60209101015160f81c60011490506000610c828887848d61338f565b90995090508015610ccd577f6ca648a381f22ead7e37773d934e64885dcf861fbfbb26c40354cbf0c4662d1a8760200151601054604051610cc4929190615521565b60405180910390a15b50610cdd88828c602001516133c7565b6020808b015187820151604080518781526001600160601b038d16948101949094528415159084015284151560608401528b1515608084015290917faeb4b4786571e184246d39587f659abf0e26f41f6a3358692250382c0cdb47b79060a00160405180910390a3505050505050505b9392505050565b610d5c612e20565b610d658161351a565b610d845780604051635428d44960e01b8152600401610a729190614e09565b600080600080610d9386612b2a565b945094505093509350336001600160a01b0316826001600160a01b031614610df65760405162461bcd60e51b81526020600482015260166024820152752737ba1039bab139b1b934b83a34b7b71037bbb732b960511b6044820152606401610a72565b610dff866110f7565b15610e455760405162461bcd60e51b815260206004820152601660248201527550656e64696e6720726571756573742065786973747360501b6044820152606401610a72565b6040805160c0810182526001815260208082018990526001600160a01b03851682840152606082018490526001600160601b038088166080840152861660a083015291519091600091610e9a918491016155c2565b6040516020818303038152906040529050610eb488613585565b505060405163ce3f471960e01b81526001600160a01b0388169063ce3f4719906001600160601b03881690610eed908590600401615687565b6000604051808303818588803b158015610f0657600080fd5b505af1158015610f1a573d6000803e3d6000fd5b50506002546001600160a01b031615801593509150610f43905057506001600160601b03861615155b15610ffe5760025460405163a9059cbb60e01b81526001600160a01b039091169063a9059cbb90610f7a908a908a9060040161569a565b6020604051808303816000875af1158015610f99573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fbd91906156bc565b610ffe5760405162461bcd60e51b8152602060048201526012602482015271696e73756666696369656e742066756e647360701b6044820152606401610a72565b600c805460ff60301b1916600160301b17905560005b83518110156110a55783818151811061102f5761102f61552f565b60200260200101516001600160a01b0316638ea98117896040518263ffffffff1660e01b81526004016110629190614e09565b600060405180830381600087803b15801561107c57600080fd5b505af1158015611090573d6000803e3d6000fd5b505050508061109e90615545565b9050611014565b50600c805460ff60301b191690556040517fd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be4187906110e59089908b906156d9565b60405180910390a15050505050505050565b6000818152600560205260408120600201805480830361111b575060009392505050565b60005b818110156111a05760006004600085848154811061113e5761113e61552f565b60009182526020808320909101546001600160a01b0316835282810193909352604091820181208982529092529020546001600160401b03600160481b90910416111561119057506001949350505050565b61119981615545565b905061111e565b506000949350505050565b6111b3612e20565b6111bb612f2a565b6002546001600160a01b03166111e45760405163c1f0c0a160e01b815260040160405180910390fd5b600b546001600160601b03166111fb811515612f7d565b600b80546001600160601b0319169055600a80548291906000906112299084906001600160601b0316615501565b82546001600160601b039182166101009390930a92830291909202199091161790555060025460405163a9059cbb60e01b8152610a36916001600160a01b03169063a9059cbb90611280908690869060040161569a565b6020604051808303816000875af115801561129f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112c391906156bc565b612f7d565b6112d0612f2a565b6112d98161351a565b156112f9578060405163ac8a27ef60e01b8152600401610a729190614e09565b601180546001810182556000919091527f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c680180546001600160a01b0319166001600160a01b0383161790556040517fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af0162590611374908390614e09565b60405180910390a150565b611387612f2a565b6002546001600160a01b0316156113b157604051631688c53760e11b815260040160405180910390fd5b600280546001600160a01b039384166001600160a01b03199182161790915560038054929093169116179055565b6001546001600160a01b031633146114325760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b6044820152606401610a72565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b611491612f2a565b6040805180820182526000916114c091908590600290839083908082843760009201919091525061247a915050565b6000818152600d602052604090205490915060ff16156114f657604051634a0b8fa760e01b815260048101829052602401610a72565b60408051808201825260018082526001600160401b0385811660208085019182526000878152600d9091528581209451855492516001600160481b0319909316901515610100600160481b03191617610100929093169190910291909117909255600e805491820181559091527fbb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fd01829055517f9b911b2c240bfbef3b6a8f7ed6ee321d1258bb2a3fe6becab52ac1cd3210afd3906115b890839085906156f2565b60405180910390a1505050565b6115cd612f2a565b600a544790600160601b90046001600160601b0316818111156116075780826040516354ced18160e11b8152600401610a72929190615521565b8181101561099457600061161b82846155af565b90506000846001600160a01b03168260405160006040518083038185875af1925050503d806000811461166a576040519150601f19603f3d011682016040523d82523d6000602084013e61166f565b606091505b50509050806116915760405163950b247960e01b815260040160405180910390fd5b7f4aed7c8eed0496c8c19ea2681fcca25741c1602342e38b045d9f1e8e905d2e9c85836040516116c29291906156d9565b60405180910390a15050505050565b6116d9612e20565b6000818152600560205260409020546116fa906001600160a01b031661372d565b60008181526006602052604090208054600160601b90046001600160601b0316903490600c6117298385615709565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555034600a600c8282829054906101000a90046001600160601b03166117719190615709565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f7603b205d03651ee812f803fccde89f1012e545a9c99f0abfea9cedd0fd8e9028234846117c49190615729565b6040516117d2929190615521565b60405180910390a25050565b60006117e8612e20565b602080830135600081815260059092526040909120546001600160a01b031661182457604051630fb532db60e11b815260040160405180910390fd5b336000908152600460209081526040808320848452808352928190208151606081018352905460ff811615158083526001600160401b036101008304811695840195909552600160481b909104909316918101919091529061189d5782336040516379bfd40160e01b8152600401610a7292919061573c565b600c5461ffff166118b46060870160408801615753565b61ffff1610806118d7575060c86118d16060870160408801615753565b61ffff16115b15611911576118ec6060860160408701615753565b600c5460405163539c34bb60e11b8152610a72929161ffff169060c89060040161576e565b600c5462010000900463ffffffff16611930608087016060880161578c565b63ffffffff1611156119765761194c608086016060870161578c565b600c54604051637aebf00f60e11b8152610a72929162010000900463ffffffff16906004016157a7565b6101f461198960a087016080880161578c565b63ffffffff1611156119c3576119a560a086016080870161578c565b6101f46040516311ce1afb60e21b8152600401610a729291906157a7565b8060200180516119d29061555e565b6001600160401b031690526040810180516119ec9061555e565b6001600160401b031690526020810151600090611a0f9087359033908790613754565b90955090506000611a33611a2e611a2960a08a018a6157be565b6137dd565b61385e565b905085611a3e6138cf565b86611a4f60808b0160608c0161578c565b611a5f60a08c0160808d0161578c565b3386604051602001611a779796959493929190615804565b60405160208183030381529060405280519060200120600f600088815260200190815260200160002081905550336001600160a01b03168588600001357feb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e89868c6040016020810190611aea9190615753565b8d6060016020810190611afd919061578c565b8e6080016020810190611b10919061578c565b89604051611b239695949392919061585d565b60405180910390a4505060009283526020918252604092839020815181549383015192909401516001600160481b0319909316931515610100600160481b031916939093176101006001600160401b039283160217600160481b600160881b031916600160481b91909216021790555b919050565b6000611ba2612e20565b6007546001600160401b031633611bba6001436155af565b6040516001600160601b0319606093841b81166020830152914060348201523090921b1660548201526001600160c01b031960c083901b16606882015260700160408051601f1981840301815291905280516020909101209150611c1f81600161589c565b600780546001600160401b0319166001600160401b03928316179055604080516000808252608082018352602080830182815283850183815260608086018581528a86526006855287862093518454935191516001600160601b039182166001600160c01b031990951694909417600160601b9190921602176001600160c01b0316600160c01b9290981691909102969096179055835194850184523385528481018281528585018481528884526005835294909220855181546001600160a01b03199081166001600160a01b039283161783559351600183018054909516911617909255925180519294939192611d1d9260028501920190614c68565b50611d2d91506008905084613950565b50827f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d33604051611d5e9190614e09565b60405180910390a2505090565b611d73612e20565b6002546001600160a01b03163314611d9e576040516344b0e3c360e01b815260040160405180910390fd5b60208114611dbf57604051638129bbcd60e01b815260040160405180910390fd5b6000611dcd8284018461515a565b600081815260056020526040902054909150611df1906001600160a01b031661372d565b600081815260066020526040812080546001600160601b031691869190611e188385615709565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600a60008282829054906101000a90046001600160601b0316611e609190615709565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a828784611eb39190615729565b604051611ec1929190615521565b60405180910390a2505050505050565b611ed9612f2a565b60c861ffff8a161115611f0657888960c860405163539c34bb60e11b8152600401610a729392919061576e565b60008513611f2a576040516321ea67b360e11b815260048101869052602401610a72565b8363ffffffff168363ffffffff161115611f5b5782846040516313c06e5960e11b8152600401610a729291906157a7565b609b60ff83161115611f855781609b604051631d66288d60e11b8152600401610a729291906158bc565b609b60ff82161115611faf5780609b604051631d66288d60e11b8152600401610a729291906158bc565b604080516101208101825261ffff8b1680825263ffffffff808c16602084018190526000848601528b8216606085018190528b8316608086018190528a841660a08701819052938a1660c0870181905260ff808b1660e08901819052908a16610100909801889052600c8054600160c01b90990260ff60c01b19600160b81b9093029290921661ffff60b81b19600160981b90940263ffffffff60981b19600160781b90990298909816600160781b600160b81b0319600160581b90960263ffffffff60581b19600160381b90980297909716600160301b600160781b03196201000090990265ffffffffffff19909c16909a179a909a1796909616979097179390931791909116959095179290921793909316929092179190911790556010869055517f2c6b6b12413678366b05b145c5f00745bdd00e739131ab5de82484a50c9d78b690612160908b908b908b908b908b908b908b908b908b9061ffff99909916895263ffffffff97881660208a0152958716604089015293861660608801526080870192909252841660a086015290921660c084015260ff91821660e0840152166101008201526101200190565b60405180910390a1505050505050505050565b61217b612f2a565b6000818152600560205260409020546001600160a01b031661219c8161372d565b610a368282612e4d565b606060006121b4600861395c565b90508084106121d657604051631390f2a160e01b815260040160405180910390fd5b60006121e28486615729565b9050818111806121f0575083155b6121fa57806121fc565b815b9050600061220a86836155af565b9050806001600160401b0381111561222457612224614e3a565b60405190808252806020026020018201604052801561224d578160200160208202803683370190505b50935060005b8181101561229d576122706122688883615729565b600890613966565b8582815181106122825761228261552f565b602090810291909101015261229681615545565b9050612253565b505050505b92915050565b6122b0612e20565b6000818152600560205260409020546001600160a01b03166122d18161372d565b6000828152600560205260409020600101546001600160a01b03163314612328576000828152600560205260409081902060010154905163d084e97560e01b8152610a72916001600160a01b031690600401614e09565b600082815260056020526040908190208054336001600160a01b031991821681178355600190920180549091169055905183917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c9386916117d29185916158d0565b8161239281612dd4565b61239a612e20565b6001600160a01b03821660009081526004602090815260408083208684529091529020805460ff16156123cd5750505050565b6000848152600560205260409020600201805460631901612401576040516305a48e0f60e01b815260040160405180910390fd5b8154600160ff1990911681178355815490810182556000828152602090200180546001600160a01b0319166001600160a01b03861617905560405185907f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e19061246b908790614e09565b60405180910390a25050505050565b60008160405160200161248d919061590d565b604051602081830303815290604052805190602001209050919050565b816124b481612dd4565b6124bc612e20565b6124c5836110f7565b156124e357604051631685ecdd60e31b815260040160405180910390fd5b6001600160a01b038216600090815260046020908152604080832086845290915290205460ff1661252b5782826040516379bfd40160e01b8152600401610a7292919061573c565b60008381526005602090815260408083206002018054825181850281018501909352808352919290919083018282801561258e57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612570575b505050505090506000600182516125a591906155af565b905060005b82518110156126ae57846001600160a01b03168382815181106125cf576125cf61552f565b60200260200101516001600160a01b03160361269e5760008383815181106125f9576125f961552f565b602002602001015190508060056000898152602001908152602001600020600201838154811061262b5761262b61552f565b600091825260208083209190910180546001600160a01b0319166001600160a01b0394909416939093179092558881526005909152604090206002018054806126765761267661591b565b600082815260209020810160001990810180546001600160a01b0319169055019055506126ae565b6126a781615545565b90506125aa565b506001600160a01b038416600090815260046020908152604080832088845290915290819020805460ff191690555185907f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a79061246b908790614e09565b600061271a82840184615948565b9050806000015160ff1660011461274a57805160405163237d181f60e21b8152610a7291906001906004016158bc565b8060a001516001600160601b0316341461278e5760a08101516040516306acf13560e41b81523460048201526001600160601b039091166024820152604401610a72565b6020808201516000908152600590915260409020546001600160a01b0316156127ca576040516326afa43560e11b815260040160405180910390fd5b60005b8160600151518110156128c357604051806060016040528060011515815260200160006001600160401b0316815260200160006001600160401b031681525060046000846060015184815181106128265761282661552f565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081208684015182528352819020835181549385015194909201516001600160481b0319909316911515610100600160481b031916919091176101006001600160401b039485160217600160481b600160881b031916600160481b9390921692909202179055806128bb81615545565b9150506127cd565b50604080516060808201835260808401516001600160601b03908116835260a0850151811660208085019182526000858701818152828901805183526006845288832097518854955192516001600160401b0316600160c01b026001600160c01b03938816600160601b026001600160c01b0319909716919097161794909417169390931790945584518084018652868601516001600160a01b03908116825281860184815294880151828801908152925184526005865295909220825181549087166001600160a01b03199182161782559351600182018054919097169416939093179094559251805191926129c292600285019290910190614c68565b5050506080810151600a80546000906129e59084906001600160601b0316615709565b92506101000a8154816001600160601b0302191690836001600160601b031602179055508060a00151600a600c8282829054906101000a90046001600160601b0316612a319190615709565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550612a6d8160200151600861395090919063ffffffff16565b50505050565b600e8181548110612a8357600080fd5b600091825260209091200154905081565b81612a9e81612dd4565b612aa6612e20565b600083815260056020526040902060018101546001600160a01b03848116911614612a6d576001810180546001600160a01b0319166001600160a01b03851617905560405184907f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a190612b1c90339087906158d0565b60405180910390a250505050565b600081815260056020526040812054819081906001600160a01b03166060612b518261372d565b600086815260066020908152604080832054600583529281902060020180548251818502810185019093528083526001600160601b0380861695600160601b810490911694600160c01b9091046001600160401b0316938893929091839190830182828015612be957602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612bcb575b505050505090509450945094509450945091939590929450565b612c0b612f2a565b6002546001600160a01b0316612c345760405163c1f0c0a160e01b815260040160405180910390fd5b6002546040516370a0823160e01b81526000916001600160a01b0316906370a0823190612c65903090600401614e09565b602060405180830381865afa158015612c82573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ca69190615a73565b600a549091506001600160601b031681811115612cda5780826040516354ced18160e11b8152600401610a72929190615521565b81811015610994576000612cee82846155af565b60025460405163a9059cbb60e01b81529192506001600160a01b03169063a9059cbb90612d2190879085906004016156d9565b6020604051808303816000875af1158015612d40573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d6491906156bc565b612d8157604051631f01ff1360e21b815260040160405180910390fd5b7f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b4366008482604051612db29291906156d9565b60405180910390a150505050565b612dc8612f2a565b612dd181613972565b50565b6000818152600560205260409020546001600160a01b0316612df58161372d565b336001600160a01b03821614610a365780604051636c51fda960e11b8152600401610a729190614e09565b600c54600160301b900460ff1615612e4b5760405163769dd35360e11b815260040160405180910390fd5b565b600080612e5984613585565b60025491935091506001600160a01b031615801590612e8057506001600160601b03821615155b15612ec25760025460405163a9059cbb60e01b8152612ec2916001600160a01b03169063a9059cbb906112809087906001600160601b038816906004016156d9565b612ed583826001600160601b0316612f9b565b604080516001600160a01b03851681526001600160601b03808516602083015283169181019190915284907f8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c490606001612b1c565b6000546001600160a01b03163314612e4b5760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610a72565b80612dd157604051631e9acf1760e31b815260040160405180910390fd5b6000826001600160a01b03168260405160006040518083038185875af1925050503d8060008114612fe8576040519150601f19603f3d011682016040523d82523d6000602084013e612fed565b606091505b50509050806109945760405163950b247960e01b815260040160405180910390fd5b6040805160a08101825260006060820181815260808301829052825260208201819052918101919091526000613048846000015161247a565b6000818152600d602090815260409182902082518084019093525460ff811615158084526101009091046001600160401b031691830191909152919250906130a657604051631dfd6e1360e21b815260048101839052602401610a72565b60008286608001516040516020016130bf929190615521565b60408051601f1981840301815291815281516020928301206000818152600f909352908220549092509081900361310957604051631b44092560e11b815260040160405180910390fd5b85516020808801516040808a015160608b015160808c015160a08d01519351613138978a979096959101615a8c565b60405160208183030381529060405280519060200120811461316d5760405163354a450b60e21b815260040160405180910390fd5b600061317c8760000151613a15565b905080613234578651604051631d2827a760e31b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169163e9413d38916131d09190600401615215565b602060405180830381865afa1580156131ed573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906132119190615a73565b90508061323457865160405163175dadad60e01b8152610a729190600401615215565b6000886080015182604051602001613256929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c9050600061327d8a83613ae3565b604080516060810182529788526020880196909652948601949094525092979650505050505050565b6000816001600160401b03163a11156132ec5782156132cf57506001600160401b0381166122a2565b3a8260405163435e532d60e11b8152600401610a729291906156f2565b503a92915050565b6000806000631fe543e360e01b8685604051602401613314929190615ae0565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252600c805460ff60301b1916600160301b1790559086015160808701519192506133789163ffffffff9091169083613b4e565b600c805460ff60301b191690559695505050505050565b60008083156133ae576133a3868685613b9a565b6000915091506133be565b6133b9868685613cab565b915091505b94509492505050565b600081815260066020526040902082156134865780546001600160601b03600160601b909104811690851681101561341257604051631e9acf1760e31b815260040160405180910390fd5b61341c8582615501565b8254600160601b600160c01b031916600160601b6001600160601b039283168102919091178455600b805488939192600c9261345c928692900416615709565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555050612a6d565b80546001600160601b039081169085168110156134b657604051631e9acf1760e31b815260040160405180910390fd5b6134c08582615501565b82546001600160601b0319166001600160601b03918216178355600b805487926000916134ef91859116615709565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505050505050565b601154600090815b8181101561357b57836001600160a01b0316601182815481106135475761354761552f565b6000918252602090912001546001600160a01b03160361356b575060019392505050565b61357481615545565b9050613522565b5060009392505050565b60008181526005602090815260408083206006909252822054600290910180546001600160601b0380841694600160601b90940416925b8181101561362757600460008483815481106135da576135da61552f565b60009182526020808320909101546001600160a01b031683528281019390935260409182018120898252909252902080546001600160881b031916905561362081615545565b90506135bc565b50600085815260056020526040812080546001600160a01b0319908116825560018201805490911690559061365f6002830182614ccd565b505060008581526006602052604081205561367b600886613e9c565b506001600160601b038416156136ce57600a80548591906000906136a99084906001600160601b0316615501565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b6001600160601b038316156137265782600a600c8282829054906101000a90046001600160601b03166137019190615501565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b5050915091565b6001600160a01b038116612dd157604051630fb532db60e11b815260040160405180910390fd5b60408051602081018690526001600160a01b03851691810191909152606081018390526001600160401b03821660808201526000908190819060a00160408051601f1981840301815290829052805160209182012092506137b9918991849101615521565b60408051808303601f19018152919052805160209091012097909650945050505050565b604080516020810190915260008152600082900361380a57506040805160208101909152600081526122a2565b63125fa26760e31b61381c8385615b01565b6001600160e01b0319161461384457604051632923fee760e11b815260040160405180910390fd5b6138518260048186615b31565b810190610d4d9190615b5b565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa8260405160240161389791511515815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b03199093169290921790915292915050565b6000466138db81613ea8565b156139495760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561391f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139439190615a73565b91505090565b4391505090565b6000610d4d8383613ecb565b60006122a2825490565b6000610d4d8383613f1a565b336001600160a01b038216036139c45760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610a72565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600046613a2181613ea8565b15613ad457610100836001600160401b0316613a3b6138cf565b613a4591906155af565b1180613a615750613a546138cf565b836001600160401b031610155b15613a6f5750600092915050565b6040516315a03d4160e11b8152606490632b407a8290613a93908690600401615215565b602060405180830381865afa158015613ab0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d4d9190615a73565b50506001600160401b03164090565b6000613b178360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151613f44565b60038360200151604051602001613b2f929190615ba6565b60408051601f1981840301815291905280516020909101209392505050565b60005a611388811015613b6057600080fd5b611388810390508460408204820311613b7857600080fd5b50823b613b8457600080fd5b60008083516020850160008789f1949350505050565b600080613bdd6000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061415f92505050565b905060005a600c54613bfd908890600160581b900463ffffffff16615729565b613c0791906155af565b613c119086615bba565b600c54909150600090613c3690600160781b900463ffffffff1664e8d4a51000615bba565b90508415613c8257600c548190606490600160b81b900460ff16613c5a8587615729565b613c649190615bba565b613c6e9190615be7565b613c789190615729565b9350505050610d4d565b600c548190606490613c9e90600160b81b900460ff1682615bfb565b60ff16613c5a8587615729565b600080600080613cb9614232565b9150915060008213613ce1576040516321ea67b360e11b815260048101839052602401610a72565b6000613d236000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061415f92505050565b9050600083825a600c54613d45908d90600160581b900463ffffffff16615729565b613d4f91906155af565b613d59908b615bba565b613d639190615729565b613d7590670de0b6b3a7640000615bba565b613d7f9190615be7565b600c54909150600090613da89063ffffffff600160981b8204811691600160781b900416615c14565b613dbd9063ffffffff1664e8d4a51000615bba565b9050600085613dd483670de0b6b3a7640000615bba565b613dde9190615be7565b905060008915613e1f57600c548290606490613e0490600160c01b900460ff1687615bba565b613e0e9190615be7565b613e189190615729565b9050613e5f565b600c548290606490613e3b90600160c01b900460ff1682615bfb565b613e489060ff1687615bba565b613e529190615be7565b613e5c9190615729565b90505b676765c793fa10079d601b1b811115613e8b5760405163e80fa38160e01b815260040160405180910390fd5b9b949a509398505050505050505050565b6000610d4d83836142f9565b600061a4b1821480613ebc575062066eed82145b806122a257505062066eee1490565b6000818152600183016020526040812054613f12575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556122a2565b5060006122a2565b6000826000018281548110613f3157613f3161552f565b9060005260206000200154905092915050565b613f4d896143f3565b613f965760405162461bcd60e51b815260206004820152601a6024820152797075626c6963206b6579206973206e6f74206f6e20637572766560301b6044820152606401610a72565b613f9f886143f3565b613fe35760405162461bcd60e51b815260206004820152601560248201527467616d6d61206973206e6f74206f6e20637572766560581b6044820152606401610a72565b613fec836143f3565b6140385760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610a72565b614041826143f3565b61408c5760405162461bcd60e51b815260206004820152601c60248201527b73486173685769746e657373206973206e6f74206f6e20637572766560201b6044820152606401610a72565b614098878a88876144b6565b6140e05760405162461bcd60e51b81526020600482015260196024820152786164647228632a706b2b732a6729213d5f755769746e65737360381b6044820152606401610a72565b60006140ec8a876145ca565b905060006140ff898b878b86898961462e565b90506000614110838d8d8a8661474d565b9050808a146141515760405162461bcd60e51b815260206004820152600d60248201526c34b73b30b634b210383937b7b360991b6044820152606401610a72565b505050505050505050505050565b60004661416b81613ea8565b156141af57606c6001600160a01b031663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613ab0573d6000803e3d6000fd5b6141b88161478d565b1561422957600f602160991b016001600160a01b03166349948e0e84604051806080016040528060478152602001615d87604791396040516020016141fe929190615c31565b6040516020818303038152906040526040518263ffffffff1660e01b8152600401613a939190615687565b50600092915050565b600c5460035460408051633fabe5a360e21b815290516000938493600160381b90910463ffffffff169284926001600160a01b039092169163feaf968c9160048082019260a0929091908290030181865afa158015614295573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906142b99190615c77565b50919650909250505063ffffffff8216158015906142e557506142dc81426155af565b8263ffffffff16105b925082156142f35760105493505b50509091565b600081815260018301602052604081205480156143e257600061431d6001836155af565b8554909150600090614331906001906155af565b90508181146143965760008660000182815481106143515761435161552f565b90600052602060002001549050808760000184815481106143745761437461552f565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806143a7576143a761591b565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506122a2565b60009150506122a2565b5092915050565b80516000906401000003d019116144415760405162461bcd60e51b8152602060048201526012602482015271696e76616c696420782d6f7264696e61746560701b6044820152606401610a72565b60208201516401000003d0191161448f5760405162461bcd60e51b8152602060048201526012602482015271696e76616c696420792d6f7264696e61746560701b6044820152606401610a72565b60208201516401000003d0199080096144af8360005b60200201516147d4565b1492915050565b60006001600160a01b0382166144fc5760405162461bcd60e51b815260206004820152600b60248201526a626164207769746e65737360a81b6044820152606401610a72565b60208401516000906001161561451357601c614516565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe199182039250600091908909875160408051600080825260209091019182905292935060019161458091869188918790615cc7565b6020604051602081039080840390855afa1580156145a2573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b6145d2614ceb565b6145ff600184846040516020016145eb93929190615ce5565b6040516020818303038152906040526147f8565b90505b61460b816143f3565b6122a257805160408051602081019290925261462791016145eb565b9050614602565b614636614ceb565b825186516401000003d01991829006919006036146955760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610a72565b6146a0878988614845565b6146e55760405162461bcd60e51b8152602060048201526016602482015275119a5c9cdd081b5d5b0818da1958dac819985a5b195960521b6044820152606401610a72565b6146f0848685614845565b6147365760405162461bcd60e51b815260206004820152601760248201527614d958dbdb99081b5d5b0818da1958dac819985a5b1959604a1b6044820152606401610a72565b614741868484614963565b98975050505050505050565b60006002868686858760405160200161476b96959493929190615d06565b60408051601f1981840301815291905280516020909101209695505050505050565b6000600a82148061479f57506101a482145b806147ac575062aa37dc82145b806147b8575061210582145b806147c5575062014a3382145b806122a257505062014a341490565b6000806401000003d01980848509840990506401000003d019600782089392505050565b614800614ceb565b61480982614a26565b815261481e6148198260006144a5565b614a61565b6020820181905260029006600103611b93576020810180516401000003d019039052919050565b6000826000036148855760405162461bcd60e51b815260206004820152600b60248201526a3d32b9379039b1b0b630b960a91b6044820152606401610a72565b8351602085015160009061489b90600290615d60565b156148a757601c6148aa565b601b5b9050600070014551231950b75fc4402da1732fc9bebe198387096040805160008082526020909101918290529192506001906148ed908390869088908790615cc7565b6020604051602081039080840390855afa15801561490f573d6000803e3d6000fd5b50505060206040510351905060008660405160200161492e9190615d74565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b61496b614ceb565b83516020808601518551918601516000938493849361498c93909190614a81565b919450925090506401000003d0198582096001146149e85760405162461bcd60e51b815260206004820152601960248201527834b73b2d1036bab9ba1031329034b73b32b939b29037b3103d60391b6044820152606401610a72565b60405180604001604052806401000003d01980614a0757614a07615bd1565b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d0198110611b9357604080516020808201939093528151808203840181529082019091528051910120614a2e565b60006122a2826002614a7a6401000003d0196001615729565b901c614b61565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000614ac183838585614bfb565b9098509050614ad288828e88614c1f565b9098509050614ae388828c87614c1f565b90985090506000614af68d878b85614c1f565b9098509050614b0788828686614bfb565b9098509050614b1888828e89614c1f565b9098509050818114614b4d576401000003d019818a0998506401000003d01982890997506401000003d0198183099650614b51565b8196505b5050505050509450945094915050565b600080614b6c614d09565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152614b9e614d27565b60208160c0846005600019fa925082600003614bf15760405162461bcd60e51b81526020600482015260126024820152716269674d6f64457870206661696c7572652160701b6044820152606401610a72565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215614cbd579160200282015b82811115614cbd57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190614c88565b50614cc9929150614d45565b5090565b5080546000825590600052602060002090810190612dd19190614d45565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b80821115614cc95760008155600101614d46565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b81811015614dab57845183529383019391830191600101614d8f565b509098975050505050505050565b6001600160a01b0381168114612dd157600080fd5b8035611b9381614db9565b60008060408385031215614dec57600080fd5b823591506020830135614dfe81614db9565b809150509250929050565b6001600160a01b0391909116815260200190565b600060208284031215614e2f57600080fd5b8135610d4d81614db9565b634e487b7160e01b600052604160045260246000fd5b60405160c081016001600160401b0381118282101715614e7257614e72614e3a565b60405290565b60405161012081016001600160401b0381118282101715614e7257614e72614e3a565b604051601f8201601f191681016001600160401b0381118282101715614ec357614ec3614e3a565b604052919050565b600082601f830112614edc57600080fd5b604080519081016001600160401b0381118282101715614efe57614efe614e3a565b8060405250806040840185811115614f1557600080fd5b845b81811015614f2f578035835260209283019201614f17565b509195945050505050565b80356001600160401b0381168114611b9357600080fd5b803563ffffffff81168114611b9357600080fd5b600060c08284031215614f7757600080fd5b614f7f614e50565b9050614f8a82614f3a565b815260208083013581830152614fa260408401614f51565b6040830152614fb360608401614f51565b60608301526080830135614fc681614db9565b608083015260a08301356001600160401b0380821115614fe557600080fd5b818501915085601f830112614ff957600080fd5b81358181111561500b5761500b614e3a565b61501d601f8201601f19168501614e9b565b9150808252868482850101111561503357600080fd5b80848401858401376000848284010152508060a085015250505092915050565b8015158114612dd157600080fd5b8035611b9381615053565b60008060008385036101e081121561508357600080fd5b6101a08082121561509357600080fd5b61509b614e78565b91506150a78787614ecb565b82526150b68760408801614ecb565b60208301526080860135604083015260a0860135606083015260c086013560808301526150e560e08701614dce565b60a08301526101006150f988828901614ecb565b60c084015261510c886101408901614ecb565b60e0840152610180870135908301529093508401356001600160401b0381111561513557600080fd5b61514186828701614f65565b9250506151516101c08501615061565b90509250925092565b60006020828403121561516c57600080fd5b5035919050565b6000806040838503121561518657600080fd5b823561519181614db9565b91506020830135614dfe81614db9565b600080606083850312156151b457600080fd5b60408301848111156151c557600080fd5b8392506151d181614f3a565b9150509250929050565b6000602082840312156151ed57600080fd5b81356001600160401b0381111561520357600080fd5b820160c08185031215610d4d57600080fd5b6001600160401b0391909116815260200190565b60008083601f84011261523b57600080fd5b5081356001600160401b0381111561525257600080fd5b60208301915083602082850101111561526a57600080fd5b9250929050565b6000806000806060858703121561528757600080fd5b843561529281614db9565b93506020850135925060408501356001600160401b038111156152b457600080fd5b6152c087828801615229565b95989497509550505050565b803561ffff81168114611b9357600080fd5b803560ff81168114611b9357600080fd5b60008060008060008060008060006101208a8c03121561530e57600080fd5b6153178a6152cc565b985061532560208b01614f51565b975061533360408b01614f51565b965061534160608b01614f51565b955060808a0135945061535660a08b01614f51565b935061536460c08b01614f51565b925061537260e08b016152de565b91506153816101008b016152de565b90509295985092959850929598565b600080604083850312156153a357600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b838110156153e2578151875295820195908201906001016153c6565b509495945050505050565b602081526000610d4d60208301846153b2565b60006040828403121561541257600080fd5b610d4d8383614ecb565b6000806020838503121561542f57600080fd5b82356001600160401b0381111561544557600080fd5b61545185828601615229565b90969095509350505050565b600081518084526020808501945080840160005b838110156153e25781516001600160a01b031687529582019590820190600101615471565b6001600160601b038681168252851660208201526001600160401b03841660408201526001600160a01b038316606082015260a0608082018190526000906154e09083018461545d565b979650505050505050565b634e487b7160e01b600052601160045260246000fd5b6001600160601b038281168282160390808211156143ec576143ec6154eb565b918252602082015260400190565b634e487b7160e01b600052603260045260246000fd5b600060018201615557576155576154eb565b5060010190565b60006001600160401b038281166002600160401b03198101615582576155826154eb565b6001019392505050565b60006001600160401b038216806155a5576155a56154eb565b6000190192915050565b818103818111156122a2576122a26154eb565b6020815260ff82511660208201526020820151604082015260018060a01b0360408301511660608201526000606083015160c0608084015261560760e084018261545d565b60808501516001600160601b0390811660a0868101919091529095015190941660c0909301929092525090919050565b60005b8381101561565257818101518382015260200161563a565b50506000910152565b60008151808452615673816020860160208601615637565b601f01601f19169290920160200192915050565b602081526000610d4d602083018461565b565b6001600160a01b039290921682526001600160601b0316602082015260400190565b6000602082840312156156ce57600080fd5b8151610d4d81615053565b6001600160a01b03929092168252602082015260400190565b9182526001600160401b0316602082015260400190565b6001600160601b038181168382160190808211156143ec576143ec6154eb565b808201808211156122a2576122a26154eb565b9182526001600160a01b0316602082015260400190565b60006020828403121561576557600080fd5b610d4d826152cc565b61ffff93841681529183166020830152909116604082015260600190565b60006020828403121561579e57600080fd5b610d4d82614f51565b63ffffffff92831681529116602082015260400190565b6000808335601e198436030181126157d557600080fd5b8301803591506001600160401b038211156157ef57600080fd5b60200191503681900382131561526a57600080fd5b878152602081018790526040810186905263ffffffff8581166060830152841660808201526001600160a01b03831660a082015260e060c082018190526000906158509083018461565b565b9998505050505050505050565b86815285602082015261ffff85166040820152600063ffffffff808616606084015280851660808401525060c060a083015261474160c083018461565b565b6001600160401b038181168382160190808211156143ec576143ec6154eb565b60ff92831681529116602082015260400190565b6001600160a01b0392831681529116602082015260400190565b8060005b6002811015612a6d5781518452602093840193909101906001016158ee565b604081016122a282846158ea565b634e487b7160e01b600052603160045260246000fd5b80356001600160601b0381168114611b9357600080fd5b6000602080838503121561595b57600080fd5b82356001600160401b038082111561597257600080fd5b9084019060c0828703121561598657600080fd5b61598e614e50565b615997836152de565b8152838301358482015260408301356159af81614db9565b60408201526060830135828111156159c657600080fd5b8301601f810188136159d757600080fd5b8035838111156159e9576159e9614e3a565b8060051b93506159fa868501614e9b565b818152938201860193868101908a861115615a1457600080fd5b928701925b85841015615a3e5783359250615a2e83614db9565b8282529287019290870190615a19565b606085015250615a5391505060808401615931565b6080820152615a6460a08401615931565b60a08201529695505050505050565b600060208284031215615a8557600080fd5b5051919050565b8781526001600160401b03871660208201526040810186905263ffffffff8581166060830152841660808201526001600160a01b03831660a082015260e060c082018190526000906158509083018461565b565b828152604060208201526000615af960408301846153b2565b949350505050565b6001600160e01b03198135818116916004851015615b295780818660040360031b1b83161692505b505092915050565b60008085851115615b4157600080fd5b83861115615b4e57600080fd5b5050820193919092039150565b600060208284031215615b6d57600080fd5b604051602081016001600160401b0381118282101715615b8f57615b8f614e3a565b6040528235615b9d81615053565b81529392505050565b82815260608101610d4d60208301846158ea565b80820281158282048414176122a2576122a26154eb565b634e487b7160e01b600052601260045260246000fd5b600082615bf657615bf6615bd1565b500490565b60ff81811683821601908111156122a2576122a26154eb565b63ffffffff8281168282160390808211156143ec576143ec6154eb565b60008351615c43818460208801615637565b835190830190615c57818360208801615637565b01949350505050565b80516001600160501b0381168114611b9357600080fd5b600080600080600060a08688031215615c8f57600080fd5b615c9886615c60565b9450602086015193506040860151925060608601519150615cbb60808701615c60565b90509295509295909350565b93845260ff9290921660208401526040830152606082015260800190565b838152615cf560208201846158ea565b606081019190915260800192915050565b868152615d1660208201876158ea565b615d2360608201866158ea565b615d3060a08201856158ea565b615d3d60e08201846158ea565b60609190911b6001600160601b0319166101208201526101340195945050505050565b600082615d6f57615d6f615bd1565b500690565b615d7e81836158ea565b60400191905056feffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa164736f6c6343000813000a",
}

var VRFCoordinatorV2PlusUpgradedVersionABI = VRFCoordinatorV2PlusUpgradedVersionMetaData.ABI

var VRFCoordinatorV2PlusUpgradedVersionBin = VRFCoordinatorV2PlusUpgradedVersionMetaData.Bin

func DeployVRFCoordinatorV2PlusUpgradedVersion(auth *bind.TransactOpts, backend bind.ContractBackend, blockhashStore common.Address) (common.Address, *types.Transaction, *VRFCoordinatorV2PlusUpgradedVersion, error) {
	parsed, err := VRFCoordinatorV2PlusUpgradedVersionMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorV2PlusUpgradedVersionBin), backend, blockhashStore)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorV2PlusUpgradedVersion{address: address, abi: *parsed, VRFCoordinatorV2PlusUpgradedVersionCaller: VRFCoordinatorV2PlusUpgradedVersionCaller{contract: contract}, VRFCoordinatorV2PlusUpgradedVersionTransactor: VRFCoordinatorV2PlusUpgradedVersionTransactor{contract: contract}, VRFCoordinatorV2PlusUpgradedVersionFilterer: VRFCoordinatorV2PlusUpgradedVersionFilterer{contract: contract}}, nil
}

type VRFCoordinatorV2PlusUpgradedVersion struct {
	address common.Address
	abi     abi.ABI
	VRFCoordinatorV2PlusUpgradedVersionCaller
	VRFCoordinatorV2PlusUpgradedVersionTransactor
	VRFCoordinatorV2PlusUpgradedVersionFilterer
}

type VRFCoordinatorV2PlusUpgradedVersionCaller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2PlusUpgradedVersionTransactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2PlusUpgradedVersionFilterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2PlusUpgradedVersionSession struct {
	Contract     *VRFCoordinatorV2PlusUpgradedVersion
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV2PlusUpgradedVersionCallerSession struct {
	Contract *VRFCoordinatorV2PlusUpgradedVersionCaller
	CallOpts bind.CallOpts
}

type VRFCoordinatorV2PlusUpgradedVersionTransactorSession struct {
	Contract     *VRFCoordinatorV2PlusUpgradedVersionTransactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV2PlusUpgradedVersionRaw struct {
	Contract *VRFCoordinatorV2PlusUpgradedVersion
}

type VRFCoordinatorV2PlusUpgradedVersionCallerRaw struct {
	Contract *VRFCoordinatorV2PlusUpgradedVersionCaller
}

type VRFCoordinatorV2PlusUpgradedVersionTransactorRaw struct {
	Contract *VRFCoordinatorV2PlusUpgradedVersionTransactor
}

func NewVRFCoordinatorV2PlusUpgradedVersion(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorV2PlusUpgradedVersion, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorV2PlusUpgradedVersionABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFCoordinatorV2PlusUpgradedVersion(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersion{address: address, abi: abi, VRFCoordinatorV2PlusUpgradedVersionCaller: VRFCoordinatorV2PlusUpgradedVersionCaller{contract: contract}, VRFCoordinatorV2PlusUpgradedVersionTransactor: VRFCoordinatorV2PlusUpgradedVersionTransactor{contract: contract}, VRFCoordinatorV2PlusUpgradedVersionFilterer: VRFCoordinatorV2PlusUpgradedVersionFilterer{contract: contract}}, nil
}

func NewVRFCoordinatorV2PlusUpgradedVersionCaller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorV2PlusUpgradedVersionCaller, error) {
	contract, err := bindVRFCoordinatorV2PlusUpgradedVersion(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionCaller{contract: contract}, nil
}

func NewVRFCoordinatorV2PlusUpgradedVersionTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorV2PlusUpgradedVersionTransactor, error) {
	contract, err := bindVRFCoordinatorV2PlusUpgradedVersion(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionTransactor{contract: contract}, nil
}

func NewVRFCoordinatorV2PlusUpgradedVersionFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorV2PlusUpgradedVersionFilterer, error) {
	contract, err := bindVRFCoordinatorV2PlusUpgradedVersion(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionFilterer{contract: contract}, nil
}

func bindVRFCoordinatorV2PlusUpgradedVersion(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFCoordinatorV2PlusUpgradedVersionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.VRFCoordinatorV2PlusUpgradedVersionCaller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.VRFCoordinatorV2PlusUpgradedVersionTransactor.contract.Transfer(opts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.VRFCoordinatorV2PlusUpgradedVersionTransactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.contract.Transfer(opts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "BLOCKHASH_STORE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) LINK() (common.Address, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.LINK(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) LINK() (common.Address, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.LINK(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) LINKNATIVEFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "LINK_NATIVE_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) LINKNATIVEFEED() (common.Address, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.LINKNATIVEFEED(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) LINKNATIVEFEED() (common.Address, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.LINKNATIVEFEED(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.MAXCONSUMERS(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.MAXCONSUMERS(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) MAXNUMWORDS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "MAX_NUM_WORDS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.MAXNUMWORDS(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.MAXNUMWORDS(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "MAX_REQUEST_CONFIRMATIONS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) GetActiveSubscriptionIds(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "getActiveSubscriptionIds", startIndex, maxCount)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) GetActiveSubscriptionIds(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.GetActiveSubscriptionIds(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, startIndex, maxCount)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) GetActiveSubscriptionIds(startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.GetActiveSubscriptionIds(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, startIndex, maxCount)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) GetRequestConfig(opts *bind.CallOpts) (uint16, uint32, [][32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "getRequestConfig")

	if err != nil {
		return *new(uint16), *new(uint32), *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)
	out2 := *abi.ConvertType(out[2], new([][32]byte)).(*[][32]byte)

	return out0, out1, out2, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.GetRequestConfig(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.GetRequestConfig(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "getSubscription", subId)

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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.GetSubscription(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, subId)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) GetSubscription(subId *big.Int) (GetSubscription,

	error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.GetSubscription(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, subId)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "hashOfKey", publicKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.HashOfKey(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, publicKey)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.HashOfKey(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, publicKey)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) MigrationVersion(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "migrationVersion")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) MigrationVersion() (uint8, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.MigrationVersion(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) MigrationVersion() (uint8, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.MigrationVersion(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) Owner() (common.Address, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.Owner(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) Owner() (common.Address, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.Owner(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) PendingRequestExists(opts *bind.CallOpts, subId *big.Int) (bool, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "pendingRequestExists", subId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) PendingRequestExists(subId *big.Int) (bool, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.PendingRequestExists(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, subId)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) PendingRequestExists(subId *big.Int) (bool, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.PendingRequestExists(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, subId)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) SConfig(opts *bind.CallOpts) (SConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "s_config")

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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) SConfig() (SConfig,

	error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SConfig(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) SConfig() (SConfig,

	error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SConfig(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) SCurrentSubNonce(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "s_currentSubNonce")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) SCurrentSubNonce() (uint64, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SCurrentSubNonce(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) SCurrentSubNonce() (uint64, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SCurrentSubNonce(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) SFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "s_fallbackWeiPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SFallbackWeiPerUnitLink(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SFallbackWeiPerUnitLink(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) SProvingKeyHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "s_provingKeyHashes", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) SProvingKeyHashes(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SProvingKeyHashes(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, arg0)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) SProvingKeyHashes(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SProvingKeyHashes(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, arg0)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) SProvingKeys(opts *bind.CallOpts, arg0 [32]byte) (SProvingKeys,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "s_provingKeys", arg0)

	outstruct := new(SProvingKeys)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.MaxGas = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) SProvingKeys(arg0 [32]byte) (SProvingKeys,

	error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SProvingKeys(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, arg0)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) SProvingKeys(arg0 [32]byte) (SProvingKeys,

	error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SProvingKeys(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, arg0)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) SRequestCommitments(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "s_requestCommitments", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) SRequestCommitments(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SRequestCommitments(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, arg0)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) SRequestCommitments(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SRequestCommitments(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts, arg0)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) STotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "s_totalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) STotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.STotalBalance(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) STotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.STotalBalance(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCaller) STotalNativeBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2PlusUpgradedVersion.contract.Call(opts, &out, "s_totalNativeBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) STotalNativeBalance() (*big.Int, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.STotalNativeBalance(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionCallerSession) STotalNativeBalance() (*big.Int, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.STotalNativeBalance(&_VRFCoordinatorV2PlusUpgradedVersion.CallOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "acceptOwnership")
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.AcceptOwnership(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.AcceptOwnership(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) AcceptSubscriptionOwnerTransfer(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "addConsumer", subId, consumer)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.AddConsumer(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) AddConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.AddConsumer(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "cancelSubscription", subId, to)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.CancelSubscription(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId, to)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) CancelSubscription(subId *big.Int, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.CancelSubscription(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId, to)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "createSubscription")
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.CreateSubscription(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.CreateSubscription(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "fulfillRandomWords", proof, rc, onlyPremium)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) FulfillRandomWords(proof VRFProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.FulfillRandomWords(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, proof, rc, onlyPremium)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) FulfillRandomWords(proof VRFProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.FulfillRandomWords(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, proof, rc, onlyPremium)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) FundSubscriptionWithNative(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "fundSubscriptionWithNative", subId)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) FundSubscriptionWithNative(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.FundSubscriptionWithNative(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) FundSubscriptionWithNative(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.FundSubscriptionWithNative(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) Migrate(opts *bind.TransactOpts, subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "migrate", subId, newCoordinator)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) Migrate(subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.Migrate(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId, newCoordinator)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) Migrate(subId *big.Int, newCoordinator common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.Migrate(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId, newCoordinator)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) OnMigration(opts *bind.TransactOpts, encodedData []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "onMigration", encodedData)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) OnMigration(encodedData []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.OnMigration(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, encodedData)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) OnMigration(encodedData []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.OnMigration(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, encodedData)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.OnTokenTransfer(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.OnTokenTransfer(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) OwnerCancelSubscription(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "ownerCancelSubscription", subId)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) OwnerCancelSubscription(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.OwnerCancelSubscription(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) OwnerCancelSubscription(subId *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.OwnerCancelSubscription(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "recoverFunds", to)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RecoverFunds(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, to)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RecoverFunds(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, to)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) RecoverNativeFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "recoverNativeFunds", to)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) RecoverNativeFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RecoverNativeFunds(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, to)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) RecoverNativeFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RecoverNativeFunds(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, to)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) RegisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "registerMigratableCoordinator", target)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) RegisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RegisterMigratableCoordinator(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, target)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) RegisterMigratableCoordinator(target common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RegisterMigratableCoordinator(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, target)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) RegisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int, maxGas uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "registerProvingKey", publicProvingKey, maxGas)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) RegisterProvingKey(publicProvingKey [2]*big.Int, maxGas uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RegisterProvingKey(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, publicProvingKey, maxGas)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) RegisterProvingKey(publicProvingKey [2]*big.Int, maxGas uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RegisterProvingKey(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, publicProvingKey, maxGas)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) RemoveConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "removeConsumer", subId, consumer)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RemoveConsumer(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) RemoveConsumer(subId *big.Int, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RemoveConsumer(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) RequestRandomWords(opts *bind.TransactOpts, req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "requestRandomWords", req)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RequestRandomWords(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, req)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RequestRandomWords(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, req)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) RequestSubscriptionOwnerTransfer(subId *big.Int, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeNativePPM uint32, fulfillmentFlatFeeLinkDiscountPPM uint32, nativePremiumPercentage uint8, linkPremiumPercentage uint8) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fulfillmentFlatFeeNativePPM, fulfillmentFlatFeeLinkDiscountPPM, nativePremiumPercentage, linkPremiumPercentage)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeNativePPM uint32, fulfillmentFlatFeeLinkDiscountPPM uint32, nativePremiumPercentage uint8, linkPremiumPercentage uint8) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SetConfig(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fulfillmentFlatFeeNativePPM, fulfillmentFlatFeeLinkDiscountPPM, nativePremiumPercentage, linkPremiumPercentage)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, fulfillmentFlatFeeNativePPM uint32, fulfillmentFlatFeeLinkDiscountPPM uint32, nativePremiumPercentage uint8, linkPremiumPercentage uint8) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SetConfig(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, fulfillmentFlatFeeNativePPM, fulfillmentFlatFeeLinkDiscountPPM, nativePremiumPercentage, linkPremiumPercentage)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) SetLINKAndLINKNativeFeed(opts *bind.TransactOpts, link common.Address, linkNativeFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "setLINKAndLINKNativeFeed", link, linkNativeFeed)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) SetLINKAndLINKNativeFeed(link common.Address, linkNativeFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SetLINKAndLINKNativeFeed(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, link, linkNativeFeed)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) SetLINKAndLINKNativeFeed(link common.Address, linkNativeFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SetLINKAndLINKNativeFeed(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, link, linkNativeFeed)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.TransferOwnership(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, to)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.TransferOwnership(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, to)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) Withdraw(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "withdraw", recipient)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) Withdraw(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.Withdraw(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, recipient)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) Withdraw(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.Withdraw(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, recipient)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) WithdrawNative(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "withdrawNative", recipient)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) WithdrawNative(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.WithdrawNative(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, recipient)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) WithdrawNative(recipient common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.WithdrawNative(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, recipient)
}

type VRFCoordinatorV2PlusUpgradedVersionConfigSetIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionConfigSet)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionConfigSet)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionConfigSet struct {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionConfigSetIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionConfigSet)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseConfigSet(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionConfigSet, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionConfigSet)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegisteredIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegistered)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegistered)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegistered struct {
	CoordinatorAddress common.Address
	Raw                types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegisteredIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "CoordinatorRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "CoordinatorRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegistered)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseCoordinatorRegistered(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegistered, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegistered)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "CoordinatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsedIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsed)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsed)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsed struct {
	RequestId              *big.Int
	FallbackWeiPerUnitLink *big.Int
	Raw                    types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterFallbackWeiPerUnitLinkUsed(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsedIterator, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "FallbackWeiPerUnitLinkUsed")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsedIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "FallbackWeiPerUnitLinkUsed", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchFallbackWeiPerUnitLinkUsed(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsed) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "FallbackWeiPerUnitLinkUsed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsed)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "FallbackWeiPerUnitLinkUsed", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseFallbackWeiPerUnitLinkUsed(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsed, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsed)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "FallbackWeiPerUnitLinkUsed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionFundsRecoveredIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionFundsRecovered)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionFundsRecovered)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionFundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionFundsRecoveredIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionFundsRecovered)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseFundsRecovered(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionFundsRecovered, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionFundsRecovered)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionMigrationCompletedIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionMigrationCompletedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionMigrationCompletedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionMigrationCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted struct {
	NewCoordinator common.Address
	SubId          *big.Int
	Raw            types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterMigrationCompleted(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionMigrationCompletedIterator, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "MigrationCompleted")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionMigrationCompletedIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "MigrationCompleted", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchMigrationCompleted(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "MigrationCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseMigrationCompleted(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "MigrationCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecoveredIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecovered)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecovered)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterNativeFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "NativeFundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecoveredIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "NativeFundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchNativeFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "NativeFundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecovered)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "NativeFundsRecovered", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseNativeFundsRecovered(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecovered, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecovered)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "NativeFundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequestedIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequested)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequested)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequestedIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequested)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequested, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequested)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferredIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferred)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferred)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferredIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferred)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferred, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferred)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegisteredIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegistered)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegistered)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegistered struct {
	KeyHash [32]byte
	MaxGas  uint64
	Raw     types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterProvingKeyRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegisteredIterator, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "ProvingKeyRegistered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegisteredIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "ProvingKeyRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegistered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "ProvingKeyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegistered)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegistered, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegistered)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilledIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled struct {
	RequestId     *big.Int
	OutputSeed    *big.Int
	SubId         *big.Int
	Payment       *big.Int
	NativePayment bool
	Success       bool
	OnlyPremium   bool
	Raw           types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule, subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilledIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled, requestId []*big.Int, subId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule, subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequestedIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested struct {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (*VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequestedIterator, error) {

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

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequestedIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceledIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceled)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceled)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceled struct {
	SubId        *big.Int
	To           common.Address
	AmountLink   *big.Int
	AmountNative *big.Int
	Raw          types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceledIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceled, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceled)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceled, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceled)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAddedIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAdded)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAdded)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAdded struct {
	SubId    *big.Int
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAddedIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAdded, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAdded)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAdded, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAdded)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemovedIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemoved)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemoved)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemoved struct {
	SubId    *big.Int
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemovedIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemoved, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemoved)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemoved, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemoved)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreatedIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreated)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreated)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreated struct {
	SubId *big.Int
	Owner common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreatedIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreated, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreated)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreated, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreated)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionSubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionFunded)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionFunded)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionFunded struct {
	SubId      *big.Int
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionFunded, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionFunded)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionFunded, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionFunded)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNativeIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNative

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNativeIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNative)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNative)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNativeIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNativeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNative struct {
	SubId            *big.Int
	OldNativeBalance *big.Int
	NewNativeBalance *big.Int
	Raw              types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterSubscriptionFundedWithNative(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNativeIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "SubscriptionFundedWithNative", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNativeIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "SubscriptionFundedWithNative", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchSubscriptionFundedWithNative(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNative, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "SubscriptionFundedWithNative", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNative)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionFundedWithNative", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseSubscriptionFundedWithNative(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNative, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNative)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionFundedWithNative", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequestedIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequested)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequested)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequested struct {
	SubId *big.Int
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequestedIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequested, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequested)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequested, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequested)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferredIterator struct {
	Event *VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferred)
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
		it.Event = new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferred)
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

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferred struct {
	SubId *big.Int
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferredIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferred, subId []*big.Int) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferred)
				if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferred, error) {
	event := new(VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferred)
	if err := _VRFCoordinatorV2PlusUpgradedVersion.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersion) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["ConfigSet"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseConfigSet(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["CoordinatorRegistered"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseCoordinatorRegistered(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["FallbackWeiPerUnitLinkUsed"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseFallbackWeiPerUnitLinkUsed(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["FundsRecovered"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseFundsRecovered(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["MigrationCompleted"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseMigrationCompleted(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["NativeFundsRecovered"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseNativeFundsRecovered(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseOwnershipTransferRequested(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["OwnershipTransferred"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseOwnershipTransferred(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["ProvingKeyRegistered"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseProvingKeyRegistered(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseRandomWordsFulfilled(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["RandomWordsRequested"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseRandomWordsRequested(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["SubscriptionCanceled"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseSubscriptionCanceled(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["SubscriptionConsumerAdded"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseSubscriptionConsumerAdded(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseSubscriptionConsumerRemoved(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["SubscriptionCreated"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseSubscriptionCreated(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["SubscriptionFunded"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseSubscriptionFunded(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["SubscriptionFundedWithNative"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseSubscriptionFundedWithNative(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseSubscriptionOwnerTransferRequested(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseSubscriptionOwnerTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFCoordinatorV2PlusUpgradedVersionConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2c6b6b12413678366b05b145c5f00745bdd00e739131ab5de82484a50c9d78b6")
}

func (VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegistered) Topic() common.Hash {
	return common.HexToHash("0xb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af01625")
}

func (VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsed) Topic() common.Hash {
	return common.HexToHash("0x6ca648a381f22ead7e37773d934e64885dcf861fbfbb26c40354cbf0c4662d1a")
}

func (VRFCoordinatorV2PlusUpgradedVersionFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted) Topic() common.Hash {
	return common.HexToHash("0xd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be4187")
}

func (VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x4aed7c8eed0496c8c19ea2681fcca25741c1602342e38b045d9f1e8e905d2e9c")
}

func (VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegistered) Topic() common.Hash {
	return common.HexToHash("0x9b911b2c240bfbef3b6a8f7ed6ee321d1258bb2a3fe6becab52ac1cd3210afd3")
}

func (VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0xaeb4b4786571e184246d39587f659abf0e26f41f6a3358692250382c0cdb47b7")
}

func (VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested) Topic() common.Hash {
	return common.HexToHash("0xeb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e")
}

func (VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0x8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c4")
}

func (VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e1")
}

func (VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a7")
}

func (VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d")
}

func (VRFCoordinatorV2PlusUpgradedVersionSubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0x1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a")
}

func (VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNative) Topic() common.Hash {
	return common.HexToHash("0x7603b205d03651ee812f803fccde89f1012e545a9c99f0abfea9cedd0fd8e902")
}

func (VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a1")
}

func (VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0xd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c9386")
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersion) Address() common.Address {
	return _VRFCoordinatorV2PlusUpgradedVersion.address
}

type VRFCoordinatorV2PlusUpgradedVersionInterface interface {
	BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	LINKNATIVEFEED(opts *bind.CallOpts) (common.Address, error)

	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	MAXNUMWORDS(opts *bind.CallOpts) (uint32, error)

	MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error)

	GetActiveSubscriptionIds(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)

	GetRequestConfig(opts *bind.CallOpts) (uint16, uint32, [][32]byte, error)

	GetSubscription(opts *bind.CallOpts, subId *big.Int) (GetSubscription,

		error)

	HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error)

	MigrationVersion(opts *bind.CallOpts) (uint8, error)

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

	FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFTypesRequestCommitmentV2Plus, onlyPremium bool) (*types.Transaction, error)

	FundSubscriptionWithNative(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	Migrate(opts *bind.TransactOpts, subId *big.Int, newCoordinator common.Address) (*types.Transaction, error)

	OnMigration(opts *bind.TransactOpts, encodedData []byte) (*types.Transaction, error)

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

	FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionConfigSet, error)

	FilterCoordinatorRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegisteredIterator, error)

	WatchCoordinatorRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegistered) (event.Subscription, error)

	ParseCoordinatorRegistered(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegistered, error)

	FilterFallbackWeiPerUnitLinkUsed(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsedIterator, error)

	WatchFallbackWeiPerUnitLinkUsed(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsed) (event.Subscription, error)

	ParseFallbackWeiPerUnitLinkUsed(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionFallbackWeiPerUnitLinkUsed, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionFundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionFundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionFundsRecovered, error)

	FilterMigrationCompleted(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionMigrationCompletedIterator, error)

	WatchMigrationCompleted(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted) (event.Subscription, error)

	ParseMigrationCompleted(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionMigrationCompleted, error)

	FilterNativeFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecoveredIterator, error)

	WatchNativeFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecovered) (event.Subscription, error)

	ParseNativeFundsRecovered(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionNativeFundsRecovered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionOwnershipTransferred, error)

	FilterProvingKeyRegistered(opts *bind.FilterOpts) (*VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegisteredIterator, error)

	WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegistered) (event.Subscription, error)

	ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionProvingKeyRegistered, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled, requestId []*big.Int, subId []*big.Int) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled, error)

	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (*VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequestedIterator, error)

	WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested, keyHash [][32]byte, subId []*big.Int, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionRandomWordsRequested, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceled, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAdded, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemoved, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreated, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionFunded, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionFunded, error)

	FilterSubscriptionFundedWithNative(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNativeIterator, error)

	WatchSubscriptionFundedWithNative(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNative, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionFundedWithNative(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionFundedWithNative, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequested, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferred, subId []*big.Int) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV2PlusUpgradedVersionSubscriptionOwnerTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
