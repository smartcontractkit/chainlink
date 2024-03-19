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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendNative\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToTransferLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"transferredValue\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"expectedValue\",\"type\":\"uint96\"}],\"name\":\"InvalidNativeBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"requestVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"expectedVersion\",\"type\":\"uint8\"}],\"name\":\"InvalidVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SubscriptionIDCollisionFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativeFundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountNative\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldNativeBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newNativeBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFundedWithNative\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_NATIVE_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFCoordinatorV2PlusUpgradedVersion.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"fundSubscriptionWithNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveSubscriptionIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrationVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedData\",\"type\":\"bytes\"}],\"name\":\"onMigration\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverNativeFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"registerMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"reentrancyLock\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_currentSubNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_provingKeyHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestCommitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalNativeBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkDiscountPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"nativePremiumPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"linkPremiumPercentage\",\"type\":\"uint8\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkNativeFeed\",\"type\":\"address\"}],\"name\":\"setLINKAndLINKNativeFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162005ed638038062005ed6833981016040819052620000349162000183565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000d7565b50505060601b6001600160601b031916608052620001b5565b6001600160a01b038116331415620001325760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019657600080fd5b81516001600160a01b0381168114620001ae57600080fd5b9392505050565b60805160601c615cfb620001db600039600081816104f401526134bb0152615cfb6000f3fe6080604052600436106101e05760003560e01c8062012291146101e5578063088070f5146102125780630ae09540146102e057806315c48b841461030257806318e3dd271461032a5780631b6b6d2314610369578063294daa49146103965780632f622e6b146103b2578063301f42e9146103d2578063405b84fa146103f257806340d6bb821461041257806341af6c871461043d57806351cff8d91461046d5780635d06b4ab1461048d57806364d51a2a146104ad57806365982744146104c2578063689c4517146104e257806372e9d5651461051657806379ba5097146105365780637bce14d11461054b5780638402595e1461056b57806386fe91c71461058b5780638da5cb5b146105ab57806395b55cfc146105c95780639b1c385e146105dc5780639d40a6fd1461060a578063a21a23e414610637578063a4c0ed361461064c578063a63e0bfb1461066c578063aa433aff1461068c578063aefb212f146106ac578063b2a7cac5146106d9578063bec4c08c146106f9578063caf70c4a14610719578063cb63179714610739578063ce3f471914610759578063d98e620e1461076c578063dac83d291461078c578063dc311dd3146107ac578063e72f6e30146107dd578063ee9d2d38146107fd578063f2fde38b1461082a575b600080fd5b3480156101f157600080fd5b506101fa61084a565b604051610209939291906157a8565b60405180910390f35b34801561021e57600080fd5b50600c546102839061ffff81169063ffffffff62010000820481169160ff600160301b8204811692600160381b8304811692600160581b8104821692600160781b8204831692600160981b83041691600160b81b8104821691600160c01b9091041689565b6040805161ffff909a168a5263ffffffff98891660208b01529615159689019690965293861660608801529185166080870152841660a08601529290921660c084015260ff91821660e08401521661010082015261012001610209565b3480156102ec57600080fd5b506103006102fb36600461541b565b6108c6565b005b34801561030e57600080fd5b5061031760c881565b60405161ffff9091168152602001610209565b34801561033657600080fd5b50600a5461035190600160601b90046001600160601b031681565b6040516001600160601b039091168152602001610209565b34801561037557600080fd5b50600254610389906001600160a01b031681565b604051610209919061564c565b3480156103a257600080fd5b5060405160028152602001610209565b3480156103be57600080fd5b506103006103cd366004614f35565b61090e565b3480156103de57600080fd5b506103516103ed3660046150ec565b610a5a565b3480156103fe57600080fd5b5061030061040d36600461541b565b610ef0565b34801561041e57600080fd5b506104286101f481565b60405163ffffffff9091168152602001610209565b34801561044957600080fd5b5061045d610458366004615402565b6112c1565b6040519015158152602001610209565b34801561047957600080fd5b50610300610488366004614f35565b611467565b34801561049957600080fd5b506103006104a8366004614f35565b6115f5565b3480156104b957600080fd5b50610317606481565b3480156104ce57600080fd5b506103006104dd366004614f52565b6116ac565b3480156104ee57600080fd5b506103897f000000000000000000000000000000000000000000000000000000000000000081565b34801561052257600080fd5b50600354610389906001600160a01b031681565b34801561054257600080fd5b5061030061170c565b34801561055757600080fd5b50610300610566366004614fe6565b6117b6565b34801561057757600080fd5b50610300610586366004614f35565b6118af565b34801561059757600080fd5b50600a54610351906001600160601b031681565b3480156105b757600080fd5b506000546001600160a01b0316610389565b6103006105d7366004615402565b6119bb565b3480156105e857600080fd5b506105fc6105f73660046151da565b611adc565b604051908152602001610209565b34801561061657600080fd5b5060075461062a906001600160401b031681565b6040516102099190615941565b34801561064357600080fd5b506105fc611ead565b34801561065857600080fd5b50610300610667366004614f8b565b612080565b34801561067857600080fd5b50610300610687366004615361565b6121fa565b34801561069857600080fd5b506103006106a7366004615402565b612403565b3480156106b857600080fd5b506106cc6106c7366004615440565b61244b565b60405161020991906156c3565b3480156106e557600080fd5b506103006106f4366004615402565b61254d565b34801561070557600080fd5b5061030061071436600461541b565b612642565b34801561072557600080fd5b506105fc61073436600461500e565b612734565b34801561074557600080fd5b5061030061075436600461541b565b612764565b610300610767366004615060565b6129cf565b34801561077857600080fd5b506105fc610787366004615402565b612d3f565b34801561079857600080fd5b506103006107a736600461541b565b612d60565b3480156107b857600080fd5b506107cc6107c7366004615402565b612df6565b604051610209959493929190615955565b3480156107e957600080fd5b506103006107f8366004614f35565b612ee4565b34801561080957600080fd5b506105fc610818366004615402565b600f6020526000908152604090205481565b34801561083657600080fd5b50610300610845366004614f35565b6130bf565b600c54600e805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff169391928391908301828280156108b457602002820191906000526020600020905b8154815260200190600101908083116108a0575b50505050509050925092509250909192565b816108d0816130d3565b6108d8613134565b6108e1836112c1565b156108ff57604051631685ecdd60e31b815260040160405180910390fd5b6109098383613161565b505050565b610916613134565b61091e613315565b600b54600160601b90046001600160601b031661094e57604051631e9acf1760e31b815260040160405180910390fd5b600b8054600160601b90046001600160601b0316908190600c6109718380615b3b565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600a600c8282829054906101000a90046001600160601b03166109b99190615b3b565b92506101000a8154816001600160601b0302191690836001600160601b031602179055506000826001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d8060008114610a33576040519150601f19603f3d011682016040523d82523d6000602084013e610a38565b606091505b50509050806109095760405163950b247960e01b815260040160405180910390fd5b6000610a64613134565b60005a90506000610a758686613368565b90506000856060015163ffffffff166001600160401b03811115610a9b57610a9b615c6d565b604051908082528060200260200182016040528015610ac4578160200160208202803683370190505b50905060005b866060015163ffffffff16811015610b3b57826040015181604051602001610af39291906156d6565b6040516020818303038152906040528051906020012060001c828281518110610b1e57610b1e615c57565b602090810291909101015280610b3381615bbf565b915050610aca565b50602080830180516000908152600f9092526040808320839055905190518291631fe543e360e01b91610b7391908690602401615832565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252600c805460ff60301b1916600160301b1790559089015160808a0151919250600091610bd89163ffffffff1690846135d3565b600c805460ff60301b1916905560208a810151600090815260069091526040902054909150600160c01b90046001600160401b0316610c18816001615aad565b6020808c0151600090815260069091526040812080546001600160401b0393909316600160c01b026001600160c01b039093169290921790915560a08b01518051610c6590600190615b24565b81518110610c7557610c75615c57565b602091010151600c5460f89190911c6001149150600090610ca6908a90600160581b900463ffffffff163a8561361f565b90508115610d9e576020808d01516000908152600690915260409020546001600160601b03808316600160601b909204161015610cf657604051631e9acf1760e31b815260040160405180910390fd5b60208c81015160009081526006909152604090208054829190600c90610d2d908490600160601b90046001600160601b0316615b3b565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600b600c8282829054906101000a90046001600160601b0316610d759190615acf565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550610e79565b6020808d01516000908152600690915260409020546001600160601b0380831691161015610ddf57604051631e9acf1760e31b815260040160405180910390fd5b6020808d015160009081526006909152604081208054839290610e0c9084906001600160601b0316615b3b565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600b60008282829054906101000a90046001600160601b0316610e549190615acf565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b8b6020015188602001517f49580fdfd9497e1ed5c1b1cec0495087ae8e3f1267470ec2fb015db32e3d6aa78a604001518488604051610ed6939291909283526001600160601b039190911660208301521515604082015260600190565b60405180910390a3985050505050505050505b9392505050565b610ef8613134565b610f018161366e565b610f295780604051635428d44960e01b8152600401610f20919061564c565b60405180910390fd5b600080600080610f3886612df6565b945094505093509350336001600160a01b0316826001600160a01b031614610f9b5760405162461bcd60e51b81526020600482015260166024820152752737ba1039bab139b1b934b83a34b7b71037bbb732b960511b6044820152606401610f20565b610fa4866112c1565b15610fea5760405162461bcd60e51b815260206004820152601660248201527550656e64696e6720726571756573742065786973747360501b6044820152606401610f20565b60006040518060c00160405280610fff600290565b60ff168152602001888152602001846001600160a01b03168152602001838152602001866001600160601b03168152602001856001600160601b031681525090506000816040516020016110539190615715565b604051602081830303815290604052905061106d886136d8565b505060405163ce3f471960e01b81526001600160a01b0388169063ce3f4719906001600160601b038816906110a6908590600401615702565b6000604051808303818588803b1580156110bf57600080fd5b505af11580156110d3573d6000803e3d6000fd5b50506002546001600160a01b0316158015935091506110fc905057506001600160601b03861615155b156111c65760025460405163a9059cbb60e01b81526001600160a01b039091169063a9059cbb90611133908a908a90600401615693565b602060405180830381600087803b15801561114d57600080fd5b505af1158015611161573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611185919061502a565b6111c65760405162461bcd60e51b8152602060048201526012602482015271696e73756666696369656e742066756e647360701b6044820152606401610f20565b600c805460ff60301b1916600160301b17905560005b835181101561126f578381815181106111f7576111f7615c57565b60200260200101516001600160a01b0316638ea98117896040518263ffffffff1660e01b815260040161122a919061564c565b600060405180830381600087803b15801561124457600080fd5b505af1158015611258573d6000803e3d6000fd5b50505050808061126790615bbf565b9150506111dc565b50600c805460ff60301b191690556040517fd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be4187906112af9089908b90615660565b60405180910390a15050505050505050565b6000818152600560209081526040808320815160608101835281546001600160a01b039081168252600183015416818501526002820180548451818702810187018652818152879693958601939092919083018282801561134b57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161132d575b505050505081525050905060005b81604001515181101561145d5760005b600e5481101561144a576000611413600e838154811061138b5761138b615c57565b9060005260206000200154856040015185815181106113ac576113ac615c57565b60200260200101518860046000896040015189815181106113cf576113cf615c57565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081208d825290925290205461010090046001600160401b0316613880565b506000818152600f6020526040902054909150156114375750600195945050505050565b508061144281615bbf565b915050611369565b508061145581615bbf565b915050611359565b5060009392505050565b61146f613134565b611477613315565b6002546001600160a01b03166114a05760405163c1f0c0a160e01b815260040160405180910390fd5b600b546001600160601b03166114c957604051631e9acf1760e31b815260040160405180910390fd5b600b80546001600160601b031690819060006114e58380615b3b565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600a60008282829054906101000a90046001600160601b031661152d9190615b3b565b82546001600160601b039182166101009390930a92830291909202199091161790555060025460405163a9059cbb60e01b81526001600160a01b039091169063a9059cbb906115829085908590600401615693565b602060405180830381600087803b15801561159c57600080fd5b505af11580156115b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115d4919061502a565b6115f157604051631e9acf1760e31b815260040160405180910390fd5b5050565b6115fd613315565b6116068161366e565b15611626578060405163ac8a27ef60e01b8152600401610f20919061564c565b601280546001810182556000919091527fbb8a6a4669ba250d26cd7a459eca9d215f8307e33aebe50379bc5a3617ec34440180546001600160a01b0319166001600160a01b0383161790556040517fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af01625906116a190839061564c565b60405180910390a150565b6116b4613315565b6002546001600160a01b0316156116de57604051631688c53760e11b815260040160405180910390fd5b600280546001600160a01b039384166001600160a01b03199182161790915560038054929093169116179055565b6001546001600160a01b0316331461175f5760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b6044820152606401610f20565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6117be613315565b6040805180820182526000916117ed919084906002908390839080828437600092019190915250612734915050565b6000818152600d602052604090205490915060ff161561182357604051634a0b8fa760e01b815260048101829052602401610f20565b6000818152600d6020526040808220805460ff19166001908117909155600e805491820181559092527fbb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fd909101829055517fc9583fd3afa3d7f16eb0b88d0268e7d05c09bafa4b21e092cbd1320e1bc8089d906118a39083815260200190565b60405180910390a15050565b6118b7613315565b600a544790600160601b90046001600160601b0316818111156118f15780826040516354ced18160e11b8152600401610f209291906156d6565b818110156109095760006119058284615b24565b90506000846001600160a01b03168260405160006040518083038185875af1925050503d8060008114611954576040519150601f19603f3d011682016040523d82523d6000602084013e611959565b606091505b505090508061197b5760405163950b247960e01b815260040160405180910390fd5b7f4aed7c8eed0496c8c19ea2681fcca25741c1602342e38b045d9f1e8e905d2e9c85836040516119ac929190615660565b60405180910390a15050505050565b6119c3613134565b6000818152600560205260409020546001600160a01b03166119f857604051630fb532db60e11b815260040160405180910390fd5b60008181526006602052604090208054600160601b90046001600160601b0316903490600c611a278385615acf565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555034600a600c8282829054906101000a90046001600160601b0316611a6f9190615acf565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f7603b205d03651ee812f803fccde89f1012e545a9c99f0abfea9cedd0fd8e902823484611ac29190615a95565b604051611ad09291906156d6565b60405180910390a25050565b6000611ae6613134565b6020808301356000908152600590915260409020546001600160a01b0316611b2157604051630fb532db60e11b815260040160405180910390fd5b336000908152600460209081526040808320858301358452808352928190208151606081018352905460ff811615158083526001600160401b036101008304811695840195909552600160481b9091049093169181019190915290611ba1578360200135336040516379bfd40160e01b8152600401610f20929190615807565b600c5461ffff16611bb86060860160408701615346565b61ffff161080611bdb575060c8611bd56060860160408701615346565b61ffff16115b15611c1557611bf06060850160408601615346565b600c5460405163539c34bb60e11b8152610f20929161ffff169060c89060040161578a565b600c5462010000900463ffffffff16611c346080860160608701615462565b63ffffffff161115611c7a57611c506080850160608601615462565b600c54604051637aebf00f60e11b8152610f20929162010000900463ffffffff169060040161592a565b6101f4611c8d60a0860160808701615462565b63ffffffff161115611cc757611ca960a0850160808601615462565b6101f46040516311ce1afb60e21b8152600401610f2092919061592a565b806020018051611cd690615bda565b6001600160401b031690526020818101516000918291611cfe9188359133918a013590613880565b90925090506000611d1a611d1560a08901896159aa565b613909565b90506000611d2782613986565b905083611d326139f7565b60208a0135611d4760808c0160608d01615462565b611d5760a08d0160808e01615462565b3386604051602001611d6f979695949392919061588a565b60405160208183030381529060405280519060200120600f600086815260200190815260200160002081905550336001600160a01b0316886020013589600001357feb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e87878d6040016020810190611de69190615346565b8e6060016020810190611df99190615462565b8f6080016020810190611e0c9190615462565b89604051611e1f9695949392919061584b565b60405180910390a45050336000908152600460209081526040808320898301358452825291829020855181549287015193909601516001600160401b03908116600160481b02600160481b600160881b03199190941661010002610100600160481b0319971515979097166001600160481b031990931692909217959095171617909255925050505b919050565b6000611eb7613134565b6007546001600160401b031633611ecf600143615b24565b6040516001600160601b0319606093841b81166020830152914060348201523090921b1660548201526001600160c01b031960c083901b16606882015260700160408051601f1981840301815291905280516020909101209150611f34816001615aad565b600780546001600160401b0319166001600160401b03928316179055604080516000808252608082018352602080830182815283850183815260608086018581528a86526006855287862093518454935191516001600160601b039182166001600160c01b031990951694909417600160601b9190921602176001600160c01b0316600160c01b9290981691909102969096179055835194850184523385528481018281528585018481528884526005835294909220855181546001600160a01b03199081166001600160a01b0392831617835593516001830180549095169116179092559251805192949391926120329260028501920190614bf9565b5061204291506008905084613a87565b50827f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d33604051612073919061564c565b60405180910390a2505090565b612088613134565b6002546001600160a01b031633146120b3576040516344b0e3c360e01b815260040160405180910390fd5b602081146120d457604051638129bbcd60e01b815260040160405180910390fd5b60006120e282840184615402565b6000818152600560205260409020549091506001600160a01b031661211a57604051630fb532db60e11b815260040160405180910390fd5b600081815260066020526040812080546001600160601b0316918691906121418385615acf565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600a60008282829054906101000a90046001600160601b03166121899190615acf565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a8287846121dc9190615a95565b6040516121ea9291906156d6565b60405180910390a2505050505050565b612202613315565b60c861ffff8a16111561222f57888960c860405163539c34bb60e11b8152600401610f209392919061578a565b60008513612253576040516321ea67b360e11b815260048101869052602401610f20565b604080516101208101825261ffff8b1680825263ffffffff808c16602084018190526000848601528b8216606085018190528b8316608086018190528a841660a08701819052938a1660c0870181905260ff808b1660e08901819052908a16610100909801889052600c8054600160c01b90990260ff60c01b19600160b81b9093029290921661ffff60b81b19600160981b90940263ffffffff60981b19600160781b90990298909816600160781b600160b81b0319600160581b90960263ffffffff60581b19600160381b90980297909716600160301b600160781b03196201000090990265ffffffffffff19909c16909a179a909a1796909616979097179390931791909116959095179290921793909316929092179190911790556010869055517f95cb2ddab6d2297c29a4861691de69b3969c464aa4a9c44258b101ff02ff375a906123f0908b908b908b908b908b908990899061ffff97909716875263ffffffff95861660208801529385166040870152919093166060850152608084019290925260ff91821660a08401521660c082015260e00190565b60405180910390a1505050505050505050565b61240b613315565b6000818152600560205260409020546001600160a01b03168061244157604051630fb532db60e11b815260040160405180910390fd5b6115f18282613161565b606060006124596008613a93565b905080841061247b57604051631390f2a160e01b815260040160405180910390fd5b60006124878486615a95565b905081811180612495575083155b61249f57806124a1565b815b905060006124af8683615b24565b9050806001600160401b038111156124c9576124c9615c6d565b6040519080825280602002602001820160405280156124f2578160200160208202803683370190505b50935060005b818110156125425761251561250d8883615a95565b600890613a9d565b85828151811061252757612527615c57565b602090810291909101015261253b81615bbf565b90506124f8565b505050505b92915050565b612555613134565b6000818152600560205260409020546001600160a01b03168061258b57604051630fb532db60e11b815260040160405180910390fd5b6000828152600560205260409020600101546001600160a01b031633146125e2576000828152600560205260409081902060010154905163d084e97560e01b8152610f20916001600160a01b03169060040161564c565b600082815260056020526040908190208054336001600160a01b031991821681178355600190920180549091169055905183917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c938691611ad0918591615679565b8161264c816130d3565b612654613134565b6001600160a01b03821660009081526004602090815260408083208684529091529020805460ff16156126875750505050565b60008481526005602052604090206002018054606414156126bb576040516305a48e0f60e01b815260040160405180910390fd5b8154600160ff1990911681178355815490810182556000828152602090200180546001600160a01b0319166001600160a01b03861617905560405185907f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e19061272590879061564c565b60405180910390a25050505050565b60008160405160200161274791906156b5565b604051602081830303815290604052805190602001209050919050565b8161276e816130d3565b612776613134565b61277f836112c1565b1561279d57604051631685ecdd60e31b815260040160405180910390fd5b6001600160a01b038216600090815260046020908152604080832086845290915290205460ff166127e55782826040516379bfd40160e01b8152600401610f20929190615807565b60008381526005602090815260408083206002018054825181850281018501909352808352919290919083018282801561284857602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161282a575b5050505050905060006001825161285f9190615b24565b905060005b825181101561296b57846001600160a01b031683828151811061288957612889615c57565b60200260200101516001600160a01b031614156129595760008383815181106128b4576128b4615c57565b60200260200101519050806005600089815260200190815260200160002060020183815481106128e6576128e6615c57565b600091825260208083209190910180546001600160a01b0319166001600160a01b03949094169390931790925588815260059091526040902060020180548061293157612931615c41565b600082815260209020810160001990810180546001600160a01b03191690550190555061296b565b8061296381615bbf565b915050612864565b506001600160a01b03841660009081526004602090815260408083208884529091529081902080546001600160881b03191690555185907f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a79061272590879061564c565b60006129dd82840184615214565b9050806000015160ff16600114612a1657805160405163237d181f60e21b815260ff909116600482015260016024820152604401610f20565b8060a001516001600160601b03163414612a5a5760a08101516040516306acf13560e41b81523460048201526001600160601b039091166024820152604401610f20565b6020808201516000908152600590915260409020546001600160a01b031615612a96576040516326afa43560e11b815260040160405180910390fd5b60005b816060015151811015612b8f57604051806060016040528060011515815260200160006001600160401b0316815260200160006001600160401b03168152506004600084606001518481518110612af257612af2615c57565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081208684015182528352819020835181549385015194909201516001600160481b0319909316911515610100600160481b031916919091176101006001600160401b039485160217600160481b600160881b031916600160481b939092169290920217905580612b8781615bbf565b915050612a99565b50604080516060808201835260808401516001600160601b03908116835260a0850151811660208085019182526000858701818152828901805183526006845288832097518854955192516001600160401b0316600160c01b026001600160c01b03938816600160601b026001600160c01b0319909716919097161794909417169390931790945584518084018652868601516001600160a01b03908116825281860184815294880151828801908152925184526005865295909220825181549087166001600160a01b0319918216178255935160018201805491909716941693909317909455925180519192612c8e92600285019290910190614bf9565b5050506080810151600a8054600090612cb19084906001600160601b0316615acf565b92506101000a8154816001600160601b0302191690836001600160601b031602179055508060a00151600a600c8282829054906101000a90046001600160601b0316612cfd9190615acf565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550612d3981602001516008613a8790919063ffffffff16565b50505050565b600e8181548110612d4f57600080fd5b600091825260209091200154905081565b81612d6a816130d3565b612d72613134565b600083815260056020526040902060018101546001600160a01b03848116911614612d39576001810180546001600160a01b0319166001600160a01b03851617905560405184907f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a190612de89033908790615679565b60405180910390a250505050565b600081815260056020526040812054819081906001600160a01b0316606081612e3257604051630fb532db60e11b815260040160405180910390fd5b600086815260066020908152604080832054600583529281902060020180548251818502810185019093528083526001600160601b0380861695600160601b810490911694600160c01b9091046001600160401b0316938893929091839190830182828015612eca57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612eac575b505050505090509450945094509450945091939590929450565b612eec613315565b6002546001600160a01b0316612f155760405163c1f0c0a160e01b815260040160405180910390fd5b6002546040516370a0823160e01b81526000916001600160a01b0316906370a0823190612f4690309060040161564c565b60206040518083038186803b158015612f5e57600080fd5b505afa158015612f72573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f969190615047565b600a549091506001600160601b031681811115612fca5780826040516354ced18160e11b8152600401610f209291906156d6565b81811015610909576000612fde8284615b24565b60025460405163a9059cbb60e01b81529192506001600160a01b03169063a9059cbb906130119087908590600401615660565b602060405180830381600087803b15801561302b57600080fd5b505af115801561303f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613063919061502a565b61308057604051631f01ff1360e21b815260040160405180910390fd5b7f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b43660084826040516130b1929190615660565b60405180910390a150505050565b6130c7613315565b6130d081613aa9565b50565b6000818152600560205260409020546001600160a01b03168061310957604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b038216146115f15780604051636c51fda960e11b8152600401610f20919061564c565b600c54600160301b900460ff161561315f5760405163769dd35360e11b815260040160405180910390fd5b565b60008061316d846136d8565b60025491935091506001600160a01b03161580159061319457506001600160601b03821615155b156132435760025460405163a9059cbb60e01b81526001600160a01b039091169063a9059cbb906131d49086906001600160601b03871690600401615660565b602060405180830381600087803b1580156131ee57600080fd5b505af1158015613202573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613226919061502a565b61324357604051631e9acf1760e31b815260040160405180910390fd5b6000836001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d8060008114613299576040519150601f19603f3d011682016040523d82523d6000602084013e61329e565b606091505b50509050806132c05760405163950b247960e01b815260040160405180910390fd5b604080516001600160a01b03861681526001600160601b03808616602083015284169181019190915285907f8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c490606001612725565b6000546001600160a01b0316331461315f5760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b6044820152606401610f20565b604080516060810182526000808252602082018190529181019190915260006133948460000151612734565b6000818152600d602052604090205490915060ff166133c957604051631dfd6e1360e21b815260048101829052602401610f20565b60008185608001516040516020016133e29291906156d6565b60408051601f1981840301815291815281516020928301206000818152600f9093529120549091508061342857604051631b44092560e11b815260040160405180910390fd5b845160208087015160408089015160608a015160808b015160a08c01519351613457978a9790969591016158d6565b60405160208183030381529060405280519060200120811461348c5760405163354a450b60e21b815260040160405180910390fd5b600061349b8660000151613b4d565b905080613562578551604051631d2827a760e31b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169163e9413d38916134ef9190600401615941565b60206040518083038186803b15801561350757600080fd5b505afa15801561351b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061353f9190615047565b90508061356257855160405163175dadad60e01b8152610f209190600401615941565b6000876080015182604051602001613584929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c905060006135ab8983613c2a565b6040805160608101825297885260208801969096529486019490945250929695505050505050565b60005a6113888110156135e557600080fd5b6113888103905084604082048203116135fd57600080fd5b50823b61360957600080fd5b60008083516020850160008789f1949350505050565b6000811561364c576011546136459086908690600160201b900463ffffffff1686613c95565b9050613666565b601154613663908690869063ffffffff1686613d37565b90505b949350505050565b6000805b6012548110156136cf57826001600160a01b03166012828154811061369957613699615c57565b6000918252602090912001546001600160a01b031614156136bd5750600192915050565b806136c781615bbf565b915050613672565b50600092915050565b60008181526005602090815260408083206006909252822054600290910180546001600160601b0380841694600160601b90940416925b8181101561377a576004600084838154811061372d5761372d615c57565b60009182526020808320909101546001600160a01b031683528281019390935260409182018120898252909252902080546001600160881b031916905561377381615bbf565b905061370f565b50600085815260056020526040812080546001600160a01b031990811682556001820180549091169055906137b26002830182614c5e565b50506000858152600660205260408120556137ce600886613e5c565b506001600160601b0384161561382157600a80548591906000906137fc9084906001600160601b0316615b3b565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b6001600160601b038316156138795782600a600c8282829054906101000a90046001600160601b03166138549190615b3b565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b5050915091565b60408051602081018690526001600160a01b03851691810191909152606081018390526001600160401b03821660808201526000908190819060a00160408051601f1981840301815290829052805160209182012092506138e59189918491016156d6565b60408051808303601f19018152919052805160209091012097909650945050505050565b604080516020810190915260008152816139325750604080516020810190915260008152612547565b63125fa26760e31b6139448385615b63565b6001600160e01b0319161461396c57604051632923fee760e11b815260040160405180910390fd5b6139798260048186615a6b565b810190610ee991906150a1565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa826040516024016139bf91511515815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b03199093169290921790915292915050565b600046613a0381613e68565b15613a805760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015613a4257600080fd5b505afa158015613a56573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a7a9190615047565b91505090565b4391505090565b6000610ee98383613e8b565b6000612547825490565b6000610ee98383613eda565b6001600160a01b038116331415613afc5760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b6044820152606401610f20565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600046613b5981613e68565b15613c1b57610100836001600160401b0316613b736139f7565b613b7d9190615b24565b1180613b995750613b8c6139f7565b836001600160401b031610155b15613ba75750600092915050565b6040516315a03d4160e11b8152606490632b407a8290613bcb908690600401615941565b60206040518083038186803b158015613be357600080fd5b505afa158015613bf7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ee99190615047565b50506001600160401b03164090565b6000613c5e8360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151613f04565b60038360200151604051602001613c7692919061581e565b60408051601f1981840301815291905280516020909101209392505050565b600080613cd86000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061411f92505050565b905060005a613ce78888615a95565b613cf19190615b24565b613cfb9085615b05565b90506000613d1463ffffffff871664e8d4a51000615b05565b905082613d218284615a95565b613d2b9190615a95565b98975050505050505050565b600080613d426141e4565b905060008113613d68576040516321ea67b360e11b815260048101829052602401610f20565b6000613daa6000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061411f92505050565b9050600082825a613dbb8b8b615a95565b613dc59190615b24565b613dcf9088615b05565b613dd99190615a95565b613deb90670de0b6b3a7640000615b05565b613df59190615af1565b90506000613e0e63ffffffff881664e8d4a51000615b05565b9050613e2581676765c793fa10079d601b1b615b24565b821115613e455760405163e80fa38160e01b815260040160405180910390fd5b613e4f8183615a95565b9998505050505050505050565b6000610ee983836142af565b600061a4b1821480613e7c575062066eed82145b8061254757505062066eee1490565b6000818152600183016020526040812054613ed257508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155612547565b506000612547565b6000826000018281548110613ef157613ef1615c57565b9060005260206000200154905092915050565b613f0d896143a2565b613f565760405162461bcd60e51b815260206004820152601a6024820152797075626c6963206b6579206973206e6f74206f6e20637572766560301b6044820152606401610f20565b613f5f886143a2565b613fa35760405162461bcd60e51b815260206004820152601560248201527467616d6d61206973206e6f74206f6e20637572766560581b6044820152606401610f20565b613fac836143a2565b613ff85760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610f20565b614001826143a2565b61404c5760405162461bcd60e51b815260206004820152601c60248201527b73486173685769746e657373206973206e6f74206f6e20637572766560201b6044820152606401610f20565b614058878a8887614465565b6140a05760405162461bcd60e51b81526020600482015260196024820152786164647228632a706b2b732a6729213d5f755769746e65737360381b6044820152606401610f20565b60006140ac8a87614579565b905060006140bf898b878b8689896145dd565b905060006140d0838d8d8a866146f0565b9050808a146141115760405162461bcd60e51b815260206004820152600d60248201526c34b73b30b634b210383937b7b360991b6044820152606401610f20565b505050505050505050505050565b60004661412b81613e68565b1561416a57606c6001600160a01b031663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b158015613be357600080fd5b61417381614730565b156136cf57600f602160991b016001600160a01b03166349948e0e84604051806080016040528060488152602001615ca7604891396040516020016141b99291906155a2565b6040516020818303038152906040526040518263ffffffff1660e01b8152600401613bcb9190615702565b600c5460035460408051633fabe5a360e21b81529051600093600160381b900463ffffffff169283151592859283926001600160a01b03169163feaf968c9160048083019260a0929190829003018186803b15801561424257600080fd5b505afa158015614256573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061427a919061547d565b50945090925084915050801561429e57506142958242615b24565b8463ffffffff16105b156136665750601054949350505050565b600081815260018301602052604081205480156143985760006142d3600183615b24565b85549091506000906142e790600190615b24565b905081811461434c57600086600001828154811061430757614307615c57565b906000526020600020015490508087600001848154811061432a5761432a615c57565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061435d5761435d615c41565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050612547565b6000915050612547565b80516000906401000003d019116143f05760405162461bcd60e51b8152602060048201526012602482015271696e76616c696420782d6f7264696e61746560701b6044820152606401610f20565b60208201516401000003d0191161443e5760405162461bcd60e51b8152602060048201526012602482015271696e76616c696420792d6f7264696e61746560701b6044820152606401610f20565b60208201516401000003d01990800961445e8360005b602002015161476a565b1492915050565b60006001600160a01b0382166144ab5760405162461bcd60e51b815260206004820152600b60248201526a626164207769746e65737360a81b6044820152606401610f20565b6020840151600090600116156144c257601c6144c5565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe199182039250600091908909875160408051600080825260209091019182905292935060019161452f918691889187906156e4565b6020604051602081039080840390855afa158015614551573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b614581614c7c565b6145ae6001848460405160200161459a9392919061562b565b60405160208183030381529060405261478e565b90505b6145ba816143a2565b6125475780516040805160208101929092526145d6910161459a565b90506145b1565b6145e5614c7c565b825186516401000003d01990819006910614156146445760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610f20565b61464f8789886147dc565b6146945760405162461bcd60e51b8152602060048201526016602482015275119a5c9cdd081b5d5b0818da1958dac819985a5b195960521b6044820152606401610f20565b61469f8486856147dc565b6146e55760405162461bcd60e51b815260206004820152601760248201527614d958dbdb99081b5d5b0818da1958dac819985a5b1959604a1b6044820152606401610f20565b613d2b8684846148f7565b60006002868686858760405160200161470e969594939291906155d1565b60408051601f1981840301815291905280516020909101209695505050505050565b6000600a82148061474257506101a482145b8061474f575062aa37dc82145b8061475b575061210582145b8061254757505062014a331490565b6000806401000003d01980848509840990506401000003d019600782089392505050565b614796614c7c565b61479f826149ba565b81526147b46147af826000614454565b6149f5565b602082018190526002900660011415611ea8576020810180516401000003d019039052919050565b6000826148195760405162461bcd60e51b815260206004820152600b60248201526a3d32b9379039b1b0b630b960a91b6044820152606401610f20565b8351602085015160009061482f90600290615c01565b1561483b57601c61483e565b601b5b9050600070014551231950b75fc4402da1732fc9bebe198387096040805160008082526020909101918290529192506001906148819083908690889087906156e4565b6020604051602081039080840390855afa1580156148a3573d6000803e3d6000fd5b5050506020604051035190506000866040516020016148c29190615590565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b6148ff614c7c565b83516020808601518551918601516000938493849361492093909190614a15565b919450925090506401000003d01985820960011461497c5760405162461bcd60e51b815260206004820152601960248201527834b73b2d1036bab9ba1031329034b73b32b939b29037b3103d60391b6044820152606401610f20565b60405180604001604052806401000003d0198061499b5761499b615c2b565b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d0198110611ea8576040805160208082019390935281518082038401815290820190915280519101206149c2565b6000612547826002614a0e6401000003d0196001615a95565b901c614af5565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000614a5583838585614b8c565b9098509050614a6688828e88614bb0565b9098509050614a7788828c87614bb0565b90985090506000614a8a8d878b85614bb0565b9098509050614a9b88828686614b8c565b9098509050614aac88828e89614bb0565b9098509050818114614ae1576401000003d019818a0998506401000003d01982890997506401000003d0198183099650614ae5565b8196505b5050505050509450945094915050565b600080614b00614c9a565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152614b32614cb8565b60208160c0846005600019fa925082614b825760405162461bcd60e51b81526020600482015260126024820152716269674d6f64457870206661696c7572652160701b6044820152606401610f20565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215614c4e579160200282015b82811115614c4e57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190614c19565b50614c5a929150614cd6565b5090565b50805460008255906000526020600020908101906130d09190614cd6565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b80821115614c5a5760008155600101614cd7565b8035611ea881615c83565b600082601f830112614d0757600080fd5b604080519081016001600160401b0381118282101715614d2957614d29615c6d565b8060405250808385604086011115614d4057600080fd5b60005b6002811015614d62578135835260209283019290910190600101614d43565b509195945050505050565b8035611ea881615c98565b60008083601f840112614d8a57600080fd5b5081356001600160401b03811115614da157600080fd5b602083019150836020828501011115614db957600080fd5b9250929050565b600082601f830112614dd157600080fd5b81356001600160401b03811115614dea57614dea615c6d565b614dfd601f8201601f1916602001615a3b565b818152846020838601011115614e1257600080fd5b816020850160208301376000918101602001919091529392505050565b600060c08284031215614e4157600080fd5b614e496159f0565b905081356001600160401b038082168214614e6357600080fd5b81835260208401356020840152614e7c60408501614ee2565b6040840152614e8d60608501614ee2565b6060840152614e9e60808501614ceb565b608084015260a0840135915080821115614eb757600080fd5b50614ec484828501614dc0565b60a08301525092915050565b803561ffff81168114611ea857600080fd5b803563ffffffff81168114611ea857600080fd5b803560ff81168114611ea857600080fd5b80516001600160501b0381168114611ea857600080fd5b80356001600160601b0381168114611ea857600080fd5b600060208284031215614f4757600080fd5b8135610ee981615c83565b60008060408385031215614f6557600080fd5b8235614f7081615c83565b91506020830135614f8081615c83565b809150509250929050565b60008060008060608587031215614fa157600080fd5b8435614fac81615c83565b93506020850135925060408501356001600160401b03811115614fce57600080fd5b614fda87828801614d78565b95989497509550505050565b600060408284031215614ff857600080fd5b8260408301111561500857600080fd5b50919050565b60006040828403121561502057600080fd5b610ee98383614cf6565b60006020828403121561503c57600080fd5b8151610ee981615c98565b60006020828403121561505957600080fd5b5051919050565b6000806020838503121561507357600080fd5b82356001600160401b0381111561508957600080fd5b61509585828601614d78565b90969095509350505050565b6000602082840312156150b357600080fd5b604051602081016001600160401b03811182821017156150d5576150d5615c6d565b60405282356150e381615c98565b81529392505050565b60008060008385036101e081121561510357600080fd5b6101a08082121561511357600080fd5b61511b615a18565b91506151278787614cf6565b82526151368760408801614cf6565b60208301526080860135604083015260a0860135606083015260c0860135608083015261516560e08701614ceb565b60a083015261010061517988828901614cf6565b60c084015261518c886101408901614cf6565b60e0840152610180870135908301529093508401356001600160401b038111156151b557600080fd5b6151c186828701614e2f565b9250506151d16101c08501614d6d565b90509250925092565b6000602082840312156151ec57600080fd5b81356001600160401b0381111561520257600080fd5b820160c08185031215610ee957600080fd5b6000602080838503121561522757600080fd5b82356001600160401b038082111561523e57600080fd5b9084019060c0828703121561525257600080fd5b61525a6159f0565b61526383614ef6565b81528383013584820152604083013561527b81615c83565b604082015260608301358281111561529257600080fd5b8301601f810188136152a357600080fd5b8035838111156152b5576152b5615c6d565b8060051b93506152c6868501615a3b565b8181528681019083880186850189018c10156152e157600080fd5b600096505b8387101561531057803594506152fb85615c83565b848352600196909601959188019188016152e6565b5060608501525061532691505060808401614f1e565b608082015261533760a08401614f1e565b60a08201529695505050505050565b60006020828403121561535857600080fd5b610ee982614ed0565b60008060008060008060008060006101208a8c03121561538057600080fd5b6153898a614ed0565b985061539760208b01614ee2565b97506153a560408b01614ee2565b96506153b360608b01614ee2565b955060808a013594506153c860a08b01614ee2565b93506153d660c08b01614ee2565b92506153e460e08b01614ef6565b91506153f36101008b01614ef6565b90509295985092959850929598565b60006020828403121561541457600080fd5b5035919050565b6000806040838503121561542e57600080fd5b823591506020830135614f8081615c83565b6000806040838503121561545357600080fd5b50508035926020909101359150565b60006020828403121561547457600080fd5b610ee982614ee2565b600080600080600060a0868803121561549557600080fd5b61549e86614f07565b94506020860151935060408601519250606086015191506154c160808701614f07565b90509295509295909350565b600081518084526020808501945080840160005b838110156155065781516001600160a01b0316875295820195908201906001016154e1565b509495945050505050565b8060005b6002811015612d39578151845260209384019390910190600101615515565b600081518084526020808501945080840160005b8381101561550657815187529582019590820190600101615548565b6000815180845261557c816020860160208601615b93565b601f01601f19169290920160200192915050565b61559a8183615511565b604001919050565b600083516155b4818460208801615b93565b8351908301906155c8818360208801615b93565b01949350505050565b8681526155e16020820187615511565b6155ee6060820186615511565b6155fb60a0820185615511565b61560860e0820184615511565b60609190911b6001600160601b0319166101208201526101340195945050505050565b83815261563b6020820184615511565b606081019190915260800192915050565b6001600160a01b0391909116815260200190565b6001600160a01b03929092168252602082015260400190565b6001600160a01b0392831681529116602082015260400190565b6001600160a01b039290921682526001600160601b0316602082015260400190565b604081016125478284615511565b602081526000610ee96020830184615534565b918252602082015260400190565b93845260ff9290921660208401526040830152606082015260800190565b602081526000610ee96020830184615564565b6020815260ff82511660208201526020820151604082015260018060a01b0360408301511660608201526000606083015160c0608084015261575a60e08401826154cd565b60808501516001600160601b0390811660a0868101919091529095015190941660c0909301929092525090919050565b61ffff93841681529183166020830152909116604082015260600190565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b818110156157f9578451835293830193918301916001016157dd565b509098975050505050505050565b9182526001600160a01b0316602082015260400190565b82815260608101610ee96020830184615511565b8281526040602082015260006136666040830184615534565b86815285602082015261ffff85166040820152600063ffffffff808616606084015280851660808401525060c060a0830152613d2b60c0830184615564565b878152602081018790526040810186905263ffffffff8581166060830152841660808201526001600160a01b03831660a082015260e060c08201819052600090613e4f90830184615564565b8781526001600160401b03871660208201526040810186905263ffffffff8581166060830152841660808201526001600160a01b03831660a082015260e060c08201819052600090613e4f90830184615564565b63ffffffff92831681529116602082015260400190565b6001600160401b0391909116815260200190565b6001600160601b038681168252851660208201526001600160401b03841660408201526001600160a01b038316606082015260a06080820181905260009061599f908301846154cd565b979650505050505050565b6000808335601e198436030181126159c157600080fd5b8301803591506001600160401b038211156159db57600080fd5b602001915036819003821315614db957600080fd5b60405160c081016001600160401b0381118282101715615a1257615a12615c6d565b60405290565b60405161012081016001600160401b0381118282101715615a1257615a12615c6d565b604051601f8201601f191681016001600160401b0381118282101715615a6357615a63615c6d565b604052919050565b60008085851115615a7b57600080fd5b83861115615a8857600080fd5b5050820193919092039150565b60008219821115615aa857615aa8615c15565b500190565b60006001600160401b038281168482168083038211156155c8576155c8615c15565b60006001600160601b038281168482168083038211156155c8576155c8615c15565b600082615b0057615b00615c2b565b500490565b6000816000190483118215151615615b1f57615b1f615c15565b500290565b600082821015615b3657615b36615c15565b500390565b60006001600160601b0383811690831681811015615b5b57615b5b615c15565b039392505050565b6001600160e01b03198135818116916004851015615b8b5780818660040360031b1b83161692505b505092915050565b60005b83811015615bae578181015183820152602001615b96565b83811115612d395750506000910152565b6000600019821415615bd357615bd3615c15565b5060010190565b60006001600160401b0382811680821415615bf757615bf7615c15565b6001019392505050565b600082615c1057615c10615c2b565b500690565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052603160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b03811681146130d057600080fd5b80151581146130d057600080fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a",
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
	outstruct.Owner = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
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
	Owner         common.Address
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
