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

type VRFCoordinatorV2PlusUpgradedVersionRequestCommitment struct {
	BlockNum         uint64
	SubId            *big.Int
	CallbackGasLimit uint32
	NumWords         uint32
	Sender           common.Address
	ExtraArgs        []byte
}

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

type VRFV2PlusClientRandomWordsRequest struct {
	KeyHash              [32]byte
	SubId                *big.Int
	RequestConfirmations uint16
	CallbackGasLimit     uint32
	NumWords             uint32
	ExtraArgs            []byte
}

var VRFCoordinatorV2PlusUpgradedVersionMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendNative\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToTransferLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"transferredValue\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"expectedValue\",\"type\":\"uint96\"}],\"name\":\"InvalidNativeBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"requestVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"expectedVersion\",\"type\":\"uint8\"}],\"name\":\"InvalidVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SubscriptionIDCollisionFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativeFundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountNative\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldNativeBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newNativeBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFundedWithNative\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_NATIVE_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFCoordinatorV2PlusUpgradedVersion.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"fundSubscriptionWithNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveSubscriptionIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"subOwner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrationVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedData\",\"type\":\"bytes\"}],\"name\":\"onMigration\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverNativeFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"registerMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"reentrancyLock\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_currentSubNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_provingKeyHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestCommitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalNativeBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkNativeFeed\",\"type\":\"address\"}],\"name\":\"setLINKAndLINKNativeFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162005e3438038062005e3483398101604081905262000034916200017e565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000d3565b5050506001600160a01b0316608052620001b0565b336001600160a01b038216036200012d5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019157600080fd5b81516001600160a01b0381168114620001a957600080fd5b9392505050565b608051615c61620001d3600039600081816104f401526134790152615c616000f3fe6080604052600436106101e05760003560e01c8062012291146101e5578063088070f5146102125780630ae09540146102e057806315c48b841461030257806318e3dd271461032a5780631b6b6d2314610369578063294daa49146103965780632f622e6b146103b2578063301f42e9146103d2578063405b84fa146103f257806340d6bb821461041257806341af6c871461043d57806351cff8d91461046d5780635d06b4ab1461048d57806364d51a2a146104ad57806365982744146104c2578063689c4517146104e257806372e9d5651461051657806379ba5097146105365780637bce14d11461054b5780638402595e1461056b57806386fe91c71461058b5780638da5cb5b146105ab57806395b55cfc146105c95780639b1c385e146105dc5780639d40a6fd1461060a578063a21a23e414610637578063a4c0ed361461064c578063a63e0bfb1461066c578063aa433aff1461068c578063aefb212f146106ac578063b2a7cac5146106d9578063bec4c08c146106f9578063caf70c4a14610719578063cb63179714610739578063ce3f471914610759578063d98e620e1461076c578063dac83d291461078c578063dc311dd3146107ac578063e72f6e30146107dd578063ee9d2d38146107fd578063f2fde38b1461082a575b600080fd5b3480156101f157600080fd5b506101fa61084a565b60405161020993929190614c80565b60405180910390f35b34801561021e57600080fd5b50600c546102839061ffff81169063ffffffff62010000820481169160ff600160301b8204811692600160381b8304811692600160581b8104821692600160781b8204831692600160981b83041691600160b81b8104821691600160c01b9091041689565b6040805161ffff909a168a5263ffffffff98891660208b01529615159689019690965293861660608801529185166080870152841660a08601529290921660c084015260ff91821660e08401521661010082015261012001610209565b3480156102ec57600080fd5b506103006102fb366004614cff565b6108c6565b005b34801561030e57600080fd5b5061031760c881565b60405161ffff9091168152602001610209565b34801561033657600080fd5b50600a5461035190600160601b90046001600160601b031681565b6040516001600160601b039091168152602001610209565b34801561037557600080fd5b50600254610389906001600160a01b031681565b6040516102099190614d2f565b3480156103a257600080fd5b5060405160028152602001610209565b3480156103be57600080fd5b506103006103cd366004614d43565b61090e565b3480156103de57600080fd5b506103516103ed366004614f9d565b610a5d565b3480156103fe57600080fd5b5061030061040d366004614cff565b610ef3565b34801561041e57600080fd5b506104286101f481565b60405163ffffffff9091168152602001610209565b34801561044957600080fd5b5061045d61045836600461508b565b6112b5565b6040519015158152602001610209565b34801561047957600080fd5b50610300610488366004614d43565b61145b565b34801561049957600080fd5b506103006104a8366004614d43565b6115dd565b3480156104b957600080fd5b50610317606481565b3480156104ce57600080fd5b506103006104dd3660046150a4565b611694565b3480156104ee57600080fd5b506103897f000000000000000000000000000000000000000000000000000000000000000081565b34801561052257600080fd5b50600354610389906001600160a01b031681565b34801561054257600080fd5b506103006116f4565b34801561055757600080fd5b506103006105663660046150d2565b61179e565b34801561057757600080fd5b50610300610586366004614d43565b611897565b34801561059757600080fd5b50600a54610351906001600160601b031681565b3480156105b757600080fd5b506000546001600160a01b0316610389565b6103006105d736600461508b565b6119a3565b3480156105e857600080fd5b506105fc6105f73660046150fa565b611ac4565b604051908152602001610209565b34801561061657600080fd5b5060075461062a906001600160401b031681565b6040516102099190615134565b34801561064357600080fd5b506105fc611e95565b34801561065857600080fd5b50610300610667366004615190565b612068565b34801561067857600080fd5b5061030061068736600461520e565b6121e2565b34801561069857600080fd5b506103006106a736600461508b565b6123eb565b3480156106b857600080fd5b506106cc6106c73660046152af565b612433565b604051610209919061530c565b3480156106e557600080fd5b506103006106f436600461508b565b612535565b34801561070557600080fd5b50610300610714366004614cff565b61262a565b34801561072557600080fd5b506105fc61073436600461531f565b61271c565b34801561074557600080fd5b50610300610754366004614cff565b61274c565b61030061076736600461533b565b6129b6565b34801561077857600080fd5b506105fc61078736600461508b565b612d26565b34801561079857600080fd5b506103006107a7366004614cff565b612d47565b3480156107b857600080fd5b506107cc6107c736600461508b565b612ddd565b6040516102099594939291906153b5565b3480156107e957600080fd5b506103006107f8366004614d43565b612ecb565b34801561080957600080fd5b506105fc61081836600461508b565b600f6020526000908152604090205481565b34801561083657600080fd5b50610300610845366004614d43565b613088565b600c54600e805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff169391928391908301828280156108b457602002820191906000526020600020905b8154815260200190600101908083116108a0575b50505050509050925092509250909192565b816108d08161309c565b6108d86130fd565b6108e1836112b5565b156108ff57604051631685ecdd60e31b815260040160405180910390fd5b610909838361312a565b505050565b6109166130fd565b61091e6132cf565b600b54600160601b90046001600160601b031660000361095157604051631e9acf1760e31b815260040160405180910390fd5b600b8054600160601b90046001600160601b0316908190600c6109748380615420565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600a600c8282829054906101000a90046001600160601b03166109bc9190615420565b92506101000a8154816001600160601b0302191690836001600160601b031602179055506000826001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d8060008114610a36576040519150601f19603f3d011682016040523d82523d6000602084013e610a3b565b606091505b50509050806109095760405163950b247960e01b815260040160405180910390fd5b6000610a676130fd565b60005a90506000610a788686613322565b90506000856060015163ffffffff166001600160401b03811115610a9e57610a9e614d60565b604051908082528060200260200182016040528015610ac7578160200160208202803683370190505b50905060005b866060015163ffffffff16811015610b3e57826040015181604051602001610af6929190615440565b6040516020818303038152906040528051906020012060001c828281518110610b2157610b2161544e565b602090810291909101015280610b3681615464565b915050610acd565b50602080830180516000908152600f9092526040808320839055905190518291631fe543e360e01b91610b769190869060240161547d565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252600c805460ff60301b1916600160301b1790559089015160808a0151919250600091610bdb9163ffffffff169084613582565b600c805460ff60301b1916905560208a810151600090815260069091526040902054909150600160c01b90046001600160401b0316610c1b816001615496565b6020808c0151600090815260069091526040812080546001600160401b0393909316600160c01b026001600160c01b039093169290921790915560a08b01518051610c68906001906154b6565b81518110610c7857610c7861544e565b602091010151600c5460f89190911c6001149150600090610ca9908a90600160581b900463ffffffff163a856135ce565b90508115610da1576020808d01516000908152600690915260409020546001600160601b03808316600160601b909204161015610cf957604051631e9acf1760e31b815260040160405180910390fd5b60208c81015160009081526006909152604090208054829190600c90610d30908490600160601b90046001600160601b0316615420565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600b600c8282829054906101000a90046001600160601b0316610d7891906154c9565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550610e7c565b6020808d01516000908152600690915260409020546001600160601b0380831691161015610de257604051631e9acf1760e31b815260040160405180910390fd5b6020808d015160009081526006909152604081208054839290610e0f9084906001600160601b0316615420565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600b60008282829054906101000a90046001600160601b0316610e5791906154c9565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b8b6020015188602001517f49580fdfd9497e1ed5c1b1cec0495087ae8e3f1267470ec2fb015db32e3d6aa78a604001518488604051610ed9939291909283526001600160601b039190911660208301521515604082015260600190565b60405180910390a3985050505050505050505b9392505050565b610efb6130fd565b610f048161361d565b610f2c5780604051635428d44960e01b8152600401610f239190614d2f565b60405180910390fd5b600080600080610f3b86612ddd565b945094505093509350336001600160a01b0316826001600160a01b031614610f9e5760405162461bcd60e51b81526020600482015260166024820152752737ba1039bab139b1b934b83a34b7b71037bbb732b960511b6044820152606401610f23565b610fa7866112b5565b15610fed5760405162461bcd60e51b815260206004820152601660248201527550656e64696e6720726571756573742065786973747360501b6044820152606401610f23565b60006040518060c00160405280611002600290565b60ff168152602001888152602001846001600160a01b03168152602001838152602001866001600160601b03168152602001856001600160601b0316815250905060008160405160200161105691906154e9565b604051602081830303815290604052905061107088613686565b505060405163ce3f471960e01b81526001600160a01b0388169063ce3f4719906001600160601b038816906110a99085906004016155ae565b6000604051808303818588803b1580156110c257600080fd5b505af11580156110d6573d6000803e3d6000fd5b50506002546001600160a01b0316158015935091506110ff905057506001600160601b03861615155b156111ba5760025460405163a9059cbb60e01b81526001600160a01b039091169063a9059cbb90611136908a908a906004016155c1565b6020604051808303816000875af1158015611155573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061117991906155e3565b6111ba5760405162461bcd60e51b8152602060048201526012602482015271696e73756666696369656e742066756e647360701b6044820152606401610f23565b600c805460ff60301b1916600160301b17905560005b8351811015611263578381815181106111eb576111eb61544e565b60200260200101516001600160a01b0316638ea98117896040518263ffffffff1660e01b815260040161121e9190614d2f565b600060405180830381600087803b15801561123857600080fd5b505af115801561124c573d6000803e3d6000fd5b50505050808061125b90615464565b9150506111d0565b50600c805460ff60301b191690556040517fd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be4187906112a39089908b90615600565b60405180910390a15050505050505050565b6000818152600560209081526040808320815160608101835281546001600160a01b039081168252600183015416818501526002820180548451818702810187018652818152879693958601939092919083018282801561133f57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611321575b505050505081525050905060005b8160400151518110156114515760005b600e5481101561143e576000611407600e838154811061137f5761137f61544e565b9060005260206000200154856040015185815181106113a0576113a061544e565b60200260200101518860046000896040015189815181106113c3576113c361544e565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081208d825290925290205461010090046001600160401b031661382e565b506000818152600f60205260409020549091501561142b5750600195945050505050565b508061143681615464565b91505061135d565b508061144981615464565b91505061134d565b5060009392505050565b6114636130fd565b61146b6132cf565b6002546001600160a01b03166114945760405163c1f0c0a160e01b815260040160405180910390fd5b600b546001600160601b03166000036114c057604051631e9acf1760e31b815260040160405180910390fd5b600b80546001600160601b031690819060006114dc8380615420565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600a60008282829054906101000a90046001600160601b03166115249190615420565b82546001600160601b039182166101009390930a92830291909202199091161790555060025460405163a9059cbb60e01b81526001600160a01b039091169063a9059cbb9061157990859085906004016155c1565b6020604051808303816000875af1158015611598573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115bc91906155e3565b6115d957604051631e9acf1760e31b815260040160405180910390fd5b5050565b6115e56132cf565b6115ee8161361d565b1561160e578060405163ac8a27ef60e01b8152600401610f239190614d2f565b601280546001810182556000919091527fbb8a6a4669ba250d26cd7a459eca9d215f8307e33aebe50379bc5a3617ec34440180546001600160a01b0319166001600160a01b0383161790556040517fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af0162590611689908390614d2f565b60405180910390a150565b61169c6132cf565b6002546001600160a01b0316156116c657604051631688c53760e11b815260040160405180910390fd5b600280546001600160a01b039384166001600160a01b03199182161790915560038054929093169116179055565b6001546001600160a01b031633146117475760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b6044820152606401610f23565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6117a66132cf565b6040805180820182526000916117d591908490600290839083908082843760009201919091525061271c915050565b6000818152600d602052604090205490915060ff161561180b57604051634a0b8fa760e01b815260048101829052602401610f23565b6000818152600d6020526040808220805460ff19166001908117909155600e805491820181559092527fbb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fd909101829055517fc9583fd3afa3d7f16eb0b88d0268e7d05c09bafa4b21e092cbd1320e1bc8089d9061188b9083815260200190565b60405180910390a15050565b61189f6132cf565b600a544790600160601b90046001600160601b0316818111156118d95780826040516354ced18160e11b8152600401610f23929190615440565b818110156109095760006118ed82846154b6565b90506000846001600160a01b03168260405160006040518083038185875af1925050503d806000811461193c576040519150601f19603f3d011682016040523d82523d6000602084013e611941565b606091505b50509050806119635760405163950b247960e01b815260040160405180910390fd5b7f4aed7c8eed0496c8c19ea2681fcca25741c1602342e38b045d9f1e8e905d2e9c8583604051611994929190615600565b60405180910390a15050505050565b6119ab6130fd565b6000818152600560205260409020546001600160a01b03166119e057604051630fb532db60e11b815260040160405180910390fd5b60008181526006602052604090208054600160601b90046001600160601b0316903490600c611a0f83856154c9565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555034600a600c8282829054906101000a90046001600160601b0316611a5791906154c9565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f7603b205d03651ee812f803fccde89f1012e545a9c99f0abfea9cedd0fd8e902823484611aaa9190615619565b604051611ab8929190615440565b60405180910390a25050565b6000611ace6130fd565b6020808301356000908152600590915260409020546001600160a01b0316611b0957604051630fb532db60e11b815260040160405180910390fd5b336000908152600460209081526040808320858301358452808352928190208151606081018352905460ff811615158083526001600160401b036101008304811695840195909552600160481b9091049093169181019190915290611b89578360200135336040516379bfd40160e01b8152600401610f2392919061562c565b600c5461ffff16611ba06060860160408701615643565b61ffff161080611bc3575060c8611bbd6060860160408701615643565b61ffff16115b15611bfd57611bd86060850160408601615643565b600c5460405163539c34bb60e11b8152610f23929161ffff169060c89060040161565e565b600c5462010000900463ffffffff16611c1c608086016060870161567c565b63ffffffff161115611c6257611c38608085016060860161567c565b600c54604051637aebf00f60e11b8152610f23929162010000900463ffffffff1690600401615697565b6101f4611c7560a086016080870161567c565b63ffffffff161115611caf57611c9160a085016080860161567c565b6101f46040516311ce1afb60e21b8152600401610f23929190615697565b806020018051611cbe906156ae565b6001600160401b031690526020818101516000918291611ce69188359133918a01359061382e565b90925090506000611d02611cfd60a08901896156dc565b6138b7565b90506000611d0f82613938565b905083611d1a6139a9565b60208a0135611d2f60808c0160608d0161567c565b611d3f60a08d0160808e0161567c565b3386604051602001611d579796959493929190615722565b60405160208183030381529060405280519060200120600f600086815260200190815260200160002081905550336001600160a01b0316886020013589600001357feb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e87878d6040016020810190611dce9190615643565b8e6060016020810190611de1919061567c565b8f6080016020810190611df4919061567c565b89604051611e079695949392919061576e565b60405180910390a45050336000908152600460209081526040808320898301358452825291829020855181549287015193909601516001600160401b03908116600160481b02600160481b600160881b03199190941661010002610100600160481b0319971515979097166001600160481b031990931692909217959095171617909255925050505b919050565b6000611e9f6130fd565b6007546001600160401b031633611eb76001436154b6565b6040516001600160601b0319606093841b81166020830152914060348201523090921b1660548201526001600160c01b031960c083901b16606882015260700160408051601f1981840301815291905280516020909101209150611f1c816001615496565b600780546001600160401b0319166001600160401b03928316179055604080516000808252608082018352602080830182815283850183815260608086018581528a86526006855287862093518454935191516001600160601b039182166001600160c01b031990951694909417600160601b9190921602176001600160c01b0316600160c01b9290981691909102969096179055835194850184523385528481018281528585018481528884526005835294909220855181546001600160a01b03199081166001600160a01b03928316178355935160018301805490951691161790925592518051929493919261201a9260028501920190614b8e565b5061202a91506008905084613a2a565b50827f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d3360405161205b9190614d2f565b60405180910390a2505090565b6120706130fd565b6002546001600160a01b0316331461209b576040516344b0e3c360e01b815260040160405180910390fd5b602081146120bc57604051638129bbcd60e01b815260040160405180910390fd5b60006120ca8284018461508b565b6000818152600560205260409020549091506001600160a01b031661210257604051630fb532db60e11b815260040160405180910390fd5b600081815260066020526040812080546001600160601b03169186919061212983856154c9565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600a60008282829054906101000a90046001600160601b031661217191906154c9565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a8287846121c49190615619565b6040516121d2929190615440565b60405180910390a2505050505050565b6121ea6132cf565b60c861ffff8a16111561221757888960c860405163539c34bb60e11b8152600401610f239392919061565e565b6000851361223b576040516321ea67b360e11b815260048101869052602401610f23565b604080516101208101825261ffff8b1680825263ffffffff808c16602084018190526000848601528b8216606085018190528b8316608086018190528a841660a08701819052938a1660c0870181905260ff808b1660e08901819052908a16610100909801889052600c8054600160c01b90990260ff60c01b19600160b81b9093029290921661ffff60b81b19600160981b90940263ffffffff60981b19600160781b90990298909816600160781b600160b81b0319600160581b90960263ffffffff60581b19600160381b90980297909716600160301b600160781b03196201000090990265ffffffffffff19909c16909a179a909a1796909616979097179390931791909116959095179290921793909316929092179190911790556010869055517f95cb2ddab6d2297c29a4861691de69b3969c464aa4a9c44258b101ff02ff375a906123d8908b908b908b908b908b908990899061ffff97909716875263ffffffff95861660208801529385166040870152919093166060850152608084019290925260ff91821660a08401521660c082015260e00190565b60405180910390a1505050505050505050565b6123f36132cf565b6000818152600560205260409020546001600160a01b03168061242957604051630fb532db60e11b815260040160405180910390fd5b6115d9828261312a565b606060006124416008613a36565b905080841061246357604051631390f2a160e01b815260040160405180910390fd5b600061246f8486615619565b90508181118061247d575083155b6124875780612489565b815b9050600061249786836154b6565b9050806001600160401b038111156124b1576124b1614d60565b6040519080825280602002602001820160405280156124da578160200160208202803683370190505b50935060005b8181101561252a576124fd6124f58883615619565b600890613a40565b85828151811061250f5761250f61544e565b602090810291909101015261252381615464565b90506124e0565b505050505b92915050565b61253d6130fd565b6000818152600560205260409020546001600160a01b03168061257357604051630fb532db60e11b815260040160405180910390fd5b6000828152600560205260409020600101546001600160a01b031633146125ca576000828152600560205260409081902060010154905163d084e97560e01b8152610f23916001600160a01b031690600401614d2f565b600082815260056020526040908190208054336001600160a01b031991821681178355600190920180549091169055905183917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c938691611ab89185916157ad565b816126348161309c565b61263c6130fd565b6001600160a01b03821660009081526004602090815260408083208684529091529020805460ff161561266f5750505050565b60008481526005602052604090206002018054606319016126a3576040516305a48e0f60e01b815260040160405180910390fd5b8154600160ff1990911681178355815490810182556000828152602090200180546001600160a01b0319166001600160a01b03861617905560405185907f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e19061270d908790614d2f565b60405180910390a25050505050565b60008160405160200161272f91906157ea565b604051602081830303815290604052805190602001209050919050565b816127568161309c565b61275e6130fd565b612767836112b5565b1561278557604051631685ecdd60e31b815260040160405180910390fd5b6001600160a01b038216600090815260046020908152604080832086845290915290205460ff166127cd5782826040516379bfd40160e01b8152600401610f2392919061562c565b60008381526005602090815260408083206002018054825181850281018501909352808352919290919083018282801561283057602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612812575b5050505050905060006001825161284791906154b6565b905060005b825181101561295257846001600160a01b03168382815181106128715761287161544e565b60200260200101516001600160a01b03160361294057600083838151811061289b5761289b61544e565b60200260200101519050806005600089815260200190815260200160002060020183815481106128cd576128cd61544e565b600091825260208083209190910180546001600160a01b0319166001600160a01b039490941693909317909255888152600590915260409020600201805480612918576129186157f8565b600082815260209020810160001990810180546001600160a01b031916905501905550612952565b8061294a81615464565b91505061284c565b506001600160a01b03841660009081526004602090815260408083208884529091529081902080546001600160881b03191690555185907f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a79061270d908790614d2f565b60006129c482840184615825565b9050806000015160ff166001146129fd57805160405163237d181f60e21b815260ff909116600482015260016024820152604401610f23565b8060a001516001600160601b03163414612a415760a08101516040516306acf13560e41b81523460048201526001600160601b039091166024820152604401610f23565b6020808201516000908152600590915260409020546001600160a01b031615612a7d576040516326afa43560e11b815260040160405180910390fd5b60005b816060015151811015612b7657604051806060016040528060011515815260200160006001600160401b0316815260200160006001600160401b03168152506004600084606001518481518110612ad957612ad961544e565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081208684015182528352819020835181549385015194909201516001600160481b0319909316911515610100600160481b031916919091176101006001600160401b039485160217600160481b600160881b031916600160481b939092169290920217905580612b6e81615464565b915050612a80565b50604080516060808201835260808401516001600160601b03908116835260a0850151811660208085019182526000858701818152828901805183526006845288832097518854955192516001600160401b0316600160c01b026001600160c01b03938816600160601b026001600160c01b0319909716919097161794909417169390931790945584518084018652868601516001600160a01b03908116825281860184815294880151828801908152925184526005865295909220825181549087166001600160a01b0319918216178255935160018201805491909716941693909317909455925180519192612c7592600285019290910190614b8e565b5050506080810151600a8054600090612c989084906001600160601b03166154c9565b92506101000a8154816001600160601b0302191690836001600160601b031602179055508060a00151600a600c8282829054906101000a90046001600160601b0316612ce491906154c9565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550612d2081602001516008613a2a90919063ffffffff16565b50505050565b600e8181548110612d3657600080fd5b600091825260209091200154905081565b81612d518161309c565b612d596130fd565b600083815260056020526040902060018101546001600160a01b03848116911614612d20576001810180546001600160a01b0319166001600160a01b03851617905560405184907f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a190612dcf90339087906157ad565b60405180910390a250505050565b600081815260056020526040812054819081906001600160a01b0316606081612e1957604051630fb532db60e11b815260040160405180910390fd5b600086815260066020908152604080832054600583529281902060020180548251818502810185019093528083526001600160601b0380861695600160601b810490911694600160c01b9091046001600160401b0316938893929091839190830182828015612eb157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612e93575b505050505090509450945094509450945091939590929450565b612ed36132cf565b6002546001600160a01b0316612efc5760405163c1f0c0a160e01b815260040160405180910390fd5b6002546040516370a0823160e01b81526000916001600160a01b0316906370a0823190612f2d903090600401614d2f565b602060405180830381865afa158015612f4a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f6e9190615950565b600a549091506001600160601b031681811115612fa25780826040516354ced18160e11b8152600401610f23929190615440565b81811015610909576000612fb682846154b6565b60025460405163a9059cbb60e01b81529192506001600160a01b03169063a9059cbb90612fe99087908590600401615600565b6020604051808303816000875af1158015613008573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061302c91906155e3565b61304957604051631f01ff1360e21b815260040160405180910390fd5b7f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600848260405161307a929190615600565b60405180910390a150505050565b6130906132cf565b61309981613a4c565b50565b6000818152600560205260409020546001600160a01b0316806130d257604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b038216146115d95780604051636c51fda960e11b8152600401610f239190614d2f565b600c54600160301b900460ff16156131285760405163769dd35360e11b815260040160405180910390fd5b565b60008061313684613686565b60025491935091506001600160a01b03161580159061315d57506001600160601b03821615155b156131fd5760025460405163a9059cbb60e01b81526001600160a01b039091169063a9059cbb9061319d9086906001600160601b03871690600401615600565b6020604051808303816000875af11580156131bc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906131e091906155e3565b6131fd57604051631e9acf1760e31b815260040160405180910390fd5b6000836001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d8060008114613253576040519150601f19603f3d011682016040523d82523d6000602084013e613258565b606091505b505090508061327a5760405163950b247960e01b815260040160405180910390fd5b604080516001600160a01b03861681526001600160601b03808616602083015284169181019190915285907f8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c49060600161270d565b6000546001600160a01b031633146131285760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610f23565b6040805160608101825260008082526020820181905291810191909152600061334e846000015161271c565b6000818152600d602052604090205490915060ff1661338357604051631dfd6e1360e21b815260048101829052602401610f23565b600081856080015160405160200161339c929190615440565b60408051601f1981840301815291815281516020928301206000818152600f90935290822054909250908190036133e657604051631b44092560e11b815260040160405180910390fd5b845160208087015160408089015160608a015160808b015160a08c01519351613415978a979096959101615969565b60405160208183030381529060405280519060200120811461344a5760405163354a450b60e21b815260040160405180910390fd5b60006134598660000151613aef565b905080613511578551604051631d2827a760e31b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169163e9413d38916134ad9190600401615134565b602060405180830381865afa1580156134ca573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134ee9190615950565b90508061351157855160405163175dadad60e01b8152610f239190600401615134565b6000876080015182604051602001613533929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c9050600061355a8983613bbd565b6040805160608101825297885260208801969096529486019490945250929695505050505050565b60005a61138881101561359457600080fd5b6113888103905084604082048203116135ac57600080fd5b50823b6135b857600080fd5b60008083516020850160008789f1949350505050565b600081156135fb576011546135f49086908690600160201b900463ffffffff1686613c28565b9050613615565b601154613612908690869063ffffffff1686613cca565b90505b949350505050565b6000805b60125481101561367d57826001600160a01b0316601282815481106136485761364861544e565b6000918252602090912001546001600160a01b03160361366b5750600192915050565b8061367581615464565b915050613621565b50600092915050565b60008181526005602090815260408083206006909252822054600290910180546001600160601b0380841694600160601b90940416925b8181101561372857600460008483815481106136db576136db61544e565b60009182526020808320909101546001600160a01b031683528281019390935260409182018120898252909252902080546001600160881b031916905561372181615464565b90506136bd565b50600085815260056020526040812080546001600160a01b031990811682556001820180549091169055906137606002830182614bf3565b505060008581526006602052604081205561377c600886613def565b506001600160601b038416156137cf57600a80548591906000906137aa9084906001600160601b0316615420565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b6001600160601b038316156138275782600a600c8282829054906101000a90046001600160601b03166138029190615420565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b5050915091565b60408051602081018690526001600160a01b03851691810191909152606081018390526001600160401b03821660808201526000908190819060a00160408051601f198184030181529082905280516020918201209250613893918991849101615440565b60408051808303601f19018152919052805160209091012097909650945050505050565b60408051602081019091526000815260008290036138e4575060408051602081019091526000815261252f565b63125fa26760e31b6138f683856159bd565b6001600160e01b0319161461391e57604051632923fee760e11b815260040160405180910390fd5b61392b82600481866159ed565b810190610eec9190615a17565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa8260405160240161397191511515815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b03199093169290921790915292915050565b6000466139b581613dfb565b15613a235760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156139f9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a1d9190615950565b91505090565b4391505090565b6000610eec8383613e1e565b600061252f825490565b6000610eec8383613e6d565b336001600160a01b03821603613a9e5760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610f23565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600046613afb81613dfb565b15613bae57610100836001600160401b0316613b156139a9565b613b1f91906154b6565b1180613b3b5750613b2e6139a9565b836001600160401b031610155b15613b495750600092915050565b6040516315a03d4160e11b8152606490632b407a8290613b6d908690600401615134565b602060405180830381865afa158015613b8a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610eec9190615950565b50506001600160401b03164090565b6000613bf18360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151613e97565b60038360200151604051602001613c09929190615a62565b60408051601f1981840301815291905280516020909101209392505050565b600080613c6b6000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506140b292505050565b905060005a613c7a8888615619565b613c8491906154b6565b613c8e9085615a76565b90506000613ca763ffffffff871664e8d4a51000615a76565b905082613cb48284615619565b613cbe9190615619565b98975050505050505050565b600080613cd561417c565b905060008113613cfb576040516321ea67b360e11b815260048101829052602401610f23565b6000613d3d6000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506140b292505050565b9050600082825a613d4e8b8b615619565b613d5891906154b6565b613d629088615a76565b613d6c9190615619565b613d7e90670de0b6b3a7640000615a76565b613d889190615aa3565b90506000613da163ffffffff881664e8d4a51000615a76565b9050613db881676765c793fa10079d601b1b6154b6565b821115613dd85760405163e80fa38160e01b815260040160405180910390fd5b613de28183615619565b9998505050505050505050565b6000610eec8383614238565b600061a4b1821480613e0f575062066eed82145b8061252f57505062066eee1490565b6000818152600183016020526040812054613e655750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561252f565b50600061252f565b6000826000018281548110613e8457613e8461544e565b9060005260206000200154905092915050565b613ea089614332565b613ee95760405162461bcd60e51b815260206004820152601a6024820152797075626c6963206b6579206973206e6f74206f6e20637572766560301b6044820152606401610f23565b613ef288614332565b613f365760405162461bcd60e51b815260206004820152601560248201527467616d6d61206973206e6f74206f6e20637572766560581b6044820152606401610f23565b613f3f83614332565b613f8b5760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610f23565b613f9482614332565b613fdf5760405162461bcd60e51b815260206004820152601c60248201527b73486173685769746e657373206973206e6f74206f6e20637572766560201b6044820152606401610f23565b613feb878a88876143f5565b6140335760405162461bcd60e51b81526020600482015260196024820152786164647228632a706b2b732a6729213d5f755769746e65737360381b6044820152606401610f23565b600061403f8a87614509565b90506000614052898b878b86898961456d565b90506000614063838d8d8a86614680565b9050808a146140a45760405162461bcd60e51b815260206004820152600d60248201526c34b73b30b634b210383937b7b360991b6044820152606401610f23565b505050505050505050505050565b6000466140be81613dfb565b1561410257606c6001600160a01b031663c6f7de0e6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613b8a573d6000803e3d6000fd5b61410b816146c0565b1561367d57600f602160991b016001600160a01b03166349948e0e84604051806080016040528060488152602001615c0d60489139604051602001614151929190615ab7565b6040516020818303038152906040526040518263ffffffff1660e01b8152600401613b6d91906155ae565b600c5460035460408051633fabe5a360e21b81529051600093600160381b900463ffffffff169283151592859283926001600160a01b03169163feaf968c9160048083019260a09291908290030181865afa1580156141df573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906142039190615afd565b509450909250849150508015614227575061421e82426154b6565b8463ffffffff16105b156136155750601054949350505050565b6000818152600183016020526040812054801561432157600061425c6001836154b6565b8554909150600090614270906001906154b6565b90508181146142d55760008660000182815481106142905761429061544e565b90600052602060002001549050808760000184815481106142b3576142b361544e565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806142e6576142e66157f8565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061252f565b600091505061252f565b5092915050565b80516000906401000003d019116143805760405162461bcd60e51b8152602060048201526012602482015271696e76616c696420782d6f7264696e61746560701b6044820152606401610f23565b60208201516401000003d019116143ce5760405162461bcd60e51b8152602060048201526012602482015271696e76616c696420792d6f7264696e61746560701b6044820152606401610f23565b60208201516401000003d0199080096143ee8360005b60200201516146fa565b1492915050565b60006001600160a01b03821661443b5760405162461bcd60e51b815260206004820152600b60248201526a626164207769746e65737360a81b6044820152606401610f23565b60208401516000906001161561445257601c614455565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe19918203925060009190890987516040805160008082526020909101918290529293506001916144bf91869188918790615b4d565b6020604051602081039080840390855afa1580156144e1573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b614511614c11565b61453e6001848460405160200161452a93929190615b6b565b60405160208183030381529060405261471e565b90505b61454a81614332565b61252f578051604080516020810192909252614566910161452a565b9050614541565b614575614c11565b825186516401000003d01991829006919006036145d45760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610f23565b6145df87898861476b565b6146245760405162461bcd60e51b8152602060048201526016602482015275119a5c9cdd081b5d5b0818da1958dac819985a5b195960521b6044820152606401610f23565b61462f84868561476b565b6146755760405162461bcd60e51b815260206004820152601760248201527614d958dbdb99081b5d5b0818da1958dac819985a5b1959604a1b6044820152606401610f23565b613cbe868484614889565b60006002868686858760405160200161469e96959493929190615b8c565b60408051601f1981840301815291905280516020909101209695505050505050565b6000600a8214806146d257506101a482145b806146df575062aa37dc82145b806146eb575061210582145b8061252f57505062014a331490565b6000806401000003d01980848509840990506401000003d019600782089392505050565b614726614c11565b61472f8261494c565b815261474461473f8260006143e4565b614987565b6020820181905260029006600103611e90576020810180516401000003d019039052919050565b6000826000036147ab5760405162461bcd60e51b815260206004820152600b60248201526a3d32b9379039b1b0b630b960a91b6044820152606401610f23565b835160208501516000906147c190600290615be6565b156147cd57601c6147d0565b601b5b9050600070014551231950b75fc4402da1732fc9bebe19838709604080516000808252602090910191829052919250600190614813908390869088908790615b4d565b6020604051602081039080840390855afa158015614835573d6000803e3d6000fd5b5050506020604051035190506000866040516020016148549190615bfa565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b614891614c11565b8351602080860151855191860151600093849384936148b2939091906149a7565b919450925090506401000003d01985820960011461490e5760405162461bcd60e51b815260206004820152601960248201527834b73b2d1036bab9ba1031329034b73b32b939b29037b3103d60391b6044820152606401610f23565b60405180604001604052806401000003d0198061492d5761492d615a8d565b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d0198110611e9057604080516020808201939093528151808203840181529082019091528051910120614954565b600061252f8260026149a06401000003d0196001615619565b901c614a87565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a08905060006149e783838585614b21565b90985090506149f888828e88614b45565b9098509050614a0988828c87614b45565b90985090506000614a1c8d878b85614b45565b9098509050614a2d88828686614b21565b9098509050614a3e88828e89614b45565b9098509050818114614a73576401000003d019818a0998506401000003d01982890997506401000003d0198183099650614a77565b8196505b5050505050509450945094915050565b600080614a92614c2f565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152614ac4614c4d565b60208160c0846005600019fa925082600003614b175760405162461bcd60e51b81526020600482015260126024820152716269674d6f64457870206661696c7572652160701b6044820152606401610f23565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215614be3579160200282015b82811115614be357825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190614bae565b50614bef929150614c6b565b5090565b50805460008255906000526020600020908101906130999190614c6b565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b80821115614bef5760008155600101614c6c565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b81811015614cd157845183529383019391830191600101614cb5565b509098975050505050505050565b6001600160a01b038116811461309957600080fd5b8035611e9081614cdf565b60008060408385031215614d1257600080fd5b823591506020830135614d2481614cdf565b809150509250929050565b6001600160a01b0391909116815260200190565b600060208284031215614d5557600080fd5b8135610eec81614cdf565b634e487b7160e01b600052604160045260246000fd5b60405160c081016001600160401b0381118282101715614d9857614d98614d60565b60405290565b60405161012081016001600160401b0381118282101715614d9857614d98614d60565b604051601f8201601f191681016001600160401b0381118282101715614de957614de9614d60565b604052919050565b600082601f830112614e0257600080fd5b604080519081016001600160401b0381118282101715614e2457614e24614d60565b8060405250806040840185811115614e3b57600080fd5b845b81811015614e55578035835260209283019201614e3d565b509195945050505050565b803563ffffffff81168114611e9057600080fd5b600082601f830112614e8557600080fd5b81356001600160401b03811115614e9e57614e9e614d60565b614eb1601f8201601f1916602001614dc1565b818152846020838601011115614ec657600080fd5b816020850160208301376000918101602001919091529392505050565b600060c08284031215614ef557600080fd5b614efd614d76565b905081356001600160401b038082168214614f1757600080fd5b81835260208401356020840152614f3060408501614e60565b6040840152614f4160608501614e60565b6060840152614f5260808501614cf4565b608084015260a0840135915080821115614f6b57600080fd5b50614f7884828501614e74565b60a08301525092915050565b801515811461309957600080fd5b8035611e9081614f84565b60008060008385036101e0811215614fb457600080fd5b6101a080821215614fc457600080fd5b614fcc614d9e565b9150614fd88787614df1565b8252614fe78760408801614df1565b60208301526080860135604083015260a0860135606083015260c0860135608083015261501660e08701614cf4565b60a083015261010061502a88828901614df1565b60c084015261503d886101408901614df1565b60e0840152610180870135908301529093508401356001600160401b0381111561506657600080fd5b61507286828701614ee3565b9250506150826101c08501614f92565b90509250925092565b60006020828403121561509d57600080fd5b5035919050565b600080604083850312156150b757600080fd5b82356150c281614cdf565b91506020830135614d2481614cdf565b6000604082840312156150e457600080fd5b826040830111156150f457600080fd5b50919050565b60006020828403121561510c57600080fd5b81356001600160401b0381111561512257600080fd5b820160c08185031215610eec57600080fd5b6001600160401b0391909116815260200190565b60008083601f84011261515a57600080fd5b5081356001600160401b0381111561517157600080fd5b60208301915083602082850101111561518957600080fd5b9250929050565b600080600080606085870312156151a657600080fd5b84356151b181614cdf565b93506020850135925060408501356001600160401b038111156151d357600080fd5b6151df87828801615148565b95989497509550505050565b803561ffff81168114611e9057600080fd5b803560ff81168114611e9057600080fd5b60008060008060008060008060006101208a8c03121561522d57600080fd5b6152368a6151eb565b985061524460208b01614e60565b975061525260408b01614e60565b965061526060608b01614e60565b955060808a0135945061527560a08b01614e60565b935061528360c08b01614e60565b925061529160e08b016151fd565b91506152a06101008b016151fd565b90509295985092959850929598565b600080604083850312156152c257600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b83811015615301578151875295820195908201906001016152e5565b509495945050505050565b602081526000610eec60208301846152d1565b60006040828403121561533157600080fd5b610eec8383614df1565b6000806020838503121561534e57600080fd5b82356001600160401b0381111561536457600080fd5b61537085828601615148565b90969095509350505050565b600081518084526020808501945080840160005b838110156153015781516001600160a01b031687529582019590820190600101615390565b6001600160601b038681168252851660208201526001600160401b03841660408201526001600160a01b038316606082015260a0608082018190526000906153ff9083018461537c565b979650505050505050565b634e487b7160e01b600052601160045260246000fd5b6001600160601b0382811682821603908082111561432b5761432b61540a565b918252602082015260400190565b634e487b7160e01b600052603260045260246000fd5b6000600182016154765761547661540a565b5060010190565b82815260406020820152600061361560408301846152d1565b6001600160401b0381811683821601908082111561432b5761432b61540a565b8181038181111561252f5761252f61540a565b6001600160601b0381811683821601908082111561432b5761432b61540a565b6020815260ff82511660208201526020820151604082015260018060a01b0360408301511660608201526000606083015160c0608084015261552e60e084018261537c565b60808501516001600160601b0390811660a0868101919091529095015190941660c0909301929092525090919050565b60005b83811015615579578181015183820152602001615561565b50506000910152565b6000815180845261559a81602086016020860161555e565b601f01601f19169290920160200192915050565b602081526000610eec6020830184615582565b6001600160a01b039290921682526001600160601b0316602082015260400190565b6000602082840312156155f557600080fd5b8151610eec81614f84565b6001600160a01b03929092168252602082015260400190565b8082018082111561252f5761252f61540a565b9182526001600160a01b0316602082015260400190565b60006020828403121561565557600080fd5b610eec826151eb565b61ffff93841681529183166020830152909116604082015260600190565b60006020828403121561568e57600080fd5b610eec82614e60565b63ffffffff92831681529116602082015260400190565b60006001600160401b038281166002600160401b031981016156d2576156d261540a565b6001019392505050565b6000808335601e198436030181126156f357600080fd5b8301803591506001600160401b0382111561570d57600080fd5b60200191503681900382131561518957600080fd5b878152602081018790526040810186905263ffffffff8581166060830152841660808201526001600160a01b03831660a082015260e060c08201819052600090613de290830184615582565b86815285602082015261ffff85166040820152600063ffffffff808616606084015280851660808401525060c060a0830152613cbe60c0830184615582565b6001600160a01b0392831681529116602082015260400190565b8060005b6002811015612d205781518452602093840193909101906001016157cb565b6040810161252f82846157c7565b634e487b7160e01b600052603160045260246000fd5b80356001600160601b0381168114611e9057600080fd5b6000602080838503121561583857600080fd5b82356001600160401b038082111561584f57600080fd5b9084019060c0828703121561586357600080fd5b61586b614d76565b615874836151fd565b81528383013584820152604083013561588c81614cdf565b60408201526060830135828111156158a357600080fd5b8301601f810188136158b457600080fd5b8035838111156158c6576158c6614d60565b8060051b93506158d7868501614dc1565b818152938201860193868101908a8611156158f157600080fd5b928701925b8584101561591b578335925061590b83614cdf565b82825292870192908701906158f6565b6060850152506159309150506080840161580e565b608082015261594160a0840161580e565b60a08201529695505050505050565b60006020828403121561596257600080fd5b5051919050565b8781526001600160401b03871660208201526040810186905263ffffffff8581166060830152841660808201526001600160a01b03831660a082015260e060c08201819052600090613de290830184615582565b6001600160e01b031981358181169160048510156159e55780818660040360031b1b83161692505b505092915050565b600080858511156159fd57600080fd5b83861115615a0a57600080fd5b5050820193919092039150565b600060208284031215615a2957600080fd5b604051602081016001600160401b0381118282101715615a4b57615a4b614d60565b6040528235615a5981614f84565b81529392505050565b82815260608101610eec60208301846157c7565b808202811582820484141761252f5761252f61540a565b634e487b7160e01b600052601260045260246000fd5b600082615ab257615ab2615a8d565b500490565b60008351615ac981846020880161555e565b835190830190615add81836020880161555e565b01949350505050565b80516001600160501b0381168114611e9057600080fd5b600080600080600060a08688031215615b1557600080fd5b615b1e86615ae6565b9450602086015193506040860151925060608601519150615b4160808701615ae6565b90509295509295909350565b93845260ff9290921660208401526040830152606082015260800190565b838152615b7b60208201846157c7565b606081019190915260800192915050565b868152615b9c60208201876157c7565b615ba960608201866157c7565b615bb660a08201856157c7565b615bc360e08201846157c7565b60609190911b6001600160601b0319166101208201526101340195945050505050565b600082615bf557615bf5615a8d565b500690565b615c0481836157c7565b60400191905056fe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000813000a",
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFCoordinatorV2PlusUpgradedVersionRequestCommitment, arg2 bool) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "fulfillRandomWords", proof, rc, arg2)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorV2PlusUpgradedVersionRequestCommitment, arg2 bool) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.FulfillRandomWords(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, proof, rc, arg2)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorV2PlusUpgradedVersionRequestCommitment, arg2 bool) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.FulfillRandomWords(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, proof, rc, arg2)
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) RegisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "registerProvingKey", publicProvingKey)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) RegisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RegisterProvingKey(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) RegisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.RegisterProvingKey(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, publicProvingKey)
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
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
	FallbackWeiPerUnitLink      *big.Int
	NativePremiumPercentage     uint8
	LinkPremiumPercentage       uint8
	Raw                         types.Log
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
	RequestId  *big.Int
	OutputSeed *big.Int
	SubID      *big.Int
	Payment    *big.Int
	Success    bool
	Raw        types.Log
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subID []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIDRule []interface{}
	for _, subIDItem := range subID {
		subIDRule = append(subIDRule, subIDItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule, subIDRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilledIterator{contract: _VRFCoordinatorV2PlusUpgradedVersion.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled, requestId []*big.Int, subID []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	var subIDRule []interface{}
	for _, subIDItem := range subID {
		subIDRule = append(subIDRule, subIDItem)
	}

	logs, sub, err := _VRFCoordinatorV2PlusUpgradedVersion.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule, subIDRule)
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersion) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["ConfigSet"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseConfigSet(log)
	case _VRFCoordinatorV2PlusUpgradedVersion.abi.Events["CoordinatorRegistered"].ID:
		return _VRFCoordinatorV2PlusUpgradedVersion.ParseCoordinatorRegistered(log)
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
	return common.HexToHash("0x95cb2ddab6d2297c29a4861691de69b3969c464aa4a9c44258b101ff02ff375a")
}

func (VRFCoordinatorV2PlusUpgradedVersionCoordinatorRegistered) Topic() common.Hash {
	return common.HexToHash("0xb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af01625")
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
	return common.HexToHash("0xc9583fd3afa3d7f16eb0b88d0268e7d05c09bafa4b21e092cbd1320e1bc8089d")
}

func (VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x49580fdfd9497e1ed5c1b1cec0495087ae8e3f1267470ec2fb015db32e3d6aa7")
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

	SProvingKeyHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	SRequestCommitments(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	STotalBalance(opts *bind.CallOpts) (*big.Int, error)

	STotalNativeBalance(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, subId *big.Int, consumer common.Address) (*types.Transaction, error)

	CancelSubscription(opts *bind.TransactOpts, subId *big.Int, to common.Address) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFCoordinatorV2PlusUpgradedVersionRequestCommitment, arg2 bool) (*types.Transaction, error)

	FundSubscriptionWithNative(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	Migrate(opts *bind.TransactOpts, subId *big.Int, newCoordinator common.Address) (*types.Transaction, error)

	OnMigration(opts *bind.TransactOpts, encodedData []byte) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subId *big.Int) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RecoverNativeFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RegisterMigratableCoordinator(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error)

	RegisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error)

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

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int, subID []*big.Int) (*VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2PlusUpgradedVersionRandomWordsFulfilled, requestId []*big.Int, subID []*big.Int) (event.Subscription, error)

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
