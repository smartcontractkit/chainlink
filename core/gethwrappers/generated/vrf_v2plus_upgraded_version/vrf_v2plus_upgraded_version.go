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

type VRFCoordinatorV2PlusUpgradedVersionFeeConfig struct {
	FulfillmentFlatFeeLinkPPM   uint32
	FulfillmentFlatFeeNativePPM uint32
}

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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorNotRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendNative\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToTransferLink\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IndexOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidExtraArgsTag\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"transferredValue\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"expectedValue\",\"type\":\"uint96\"}],\"name\":\"InvalidNativeBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"requestVersion\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"expectedVersion\",\"type\":\"uint8\"}],\"name\":\"InvalidVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkNotSet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SubscriptionIDCollisionFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"}],\"indexed\":false,\"internalType\":\"structVRFCoordinatorV2PlusUpgradedVersion.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinatorAddress\",\"type\":\"address\"}],\"name\":\"CoordinatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"MigrationCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativeFundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountNative\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldNativeBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newNativeBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFundedWithNative\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_NATIVE_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFCoordinatorV2PlusUpgradedVersion.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"fundSubscriptionWithNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getActiveSubscriptionIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"nativeBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newCoordinator\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrationVersion\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encodedData\",\"type\":\"bytes\"}],\"name\":\"onMigration\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverNativeFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"registerMigratableCoordinator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"reentrancyLock\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_currentSubNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_provingKeyHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestCommitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalNativeBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPM\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeNativePPM\",\"type\":\"uint32\"}],\"internalType\":\"structVRFCoordinatorV2PlusUpgradedVersion.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkNativeFeed\",\"type\":\"address\"}],\"name\":\"setLINKAndLINKNativeFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200610538038062006105833981016040819052620000349162000183565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000d7565b50505060601b6001600160601b031916608052620001b5565b6001600160a01b038116331415620001325760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019657600080fd5b81516001600160a01b0381168114620001ae57600080fd5b9392505050565b60805160601c615f2a620001db600039600081816104a601526136510152615f2a6000f3fe6080604052600436106101e05760003560e01c8062012291146101e5578063088070f5146102125780630ae095401461029257806315c48b84146102b457806318e3dd27146102dc5780631b6b6d231461031b578063294daa49146103485780632f622e6b14610364578063330987b314610384578063405b84fa146103a457806340d6bb82146103c457806341af6c87146103ef57806351cff8d91461041f5780635d06b4ab1461043f57806364d51a2a1461045f5780636598274414610474578063689c45171461049457806372e9d565146104c857806379ba5097146104e85780637bce14d1146104fd5780638402595e1461051d57806386fe91c71461053d5780638da5cb5b1461055d57806395b55cfc1461057b5780639b1c385e1461058e5780639d40a6fd146105bc578063a21a23e4146105e9578063a4c0ed36146105fe578063aa433aff1461061e578063aefb212f1461063e578063b08c87951461066b578063b2a7cac51461068b578063bec4c08c146106ab578063caf70c4a146106cb578063cb631797146106eb578063ce3f47191461070b578063d98e620e1461071e578063dac83d291461073e578063dc311dd31461075e578063e72f6e301461078f578063ee9d2d38146107af578063f2fde38b146107dc575b600080fd5b3480156101f157600080fd5b506101fa6107fc565b604051610209939291906159b5565b60405180910390f35b34801561021e57600080fd5b50600c5461025a9061ffff81169063ffffffff62010000820481169160ff600160301b82041691600160381b8204811691600160581b90041685565b6040805161ffff909616865263ffffffff9485166020870152921515928501929092528216606084015216608082015260a001610209565b34801561029e57600080fd5b506102b26102ad366004615628565b610878565b005b3480156102c057600080fd5b506102c960c881565b60405161ffff9091168152602001610209565b3480156102e857600080fd5b50600a5461030390600160601b90046001600160601b031681565b6040516001600160601b039091168152602001610209565b34801561032757600080fd5b5060025461033b906001600160a01b031681565b6040516102099190615859565b34801561035457600080fd5b5060405160028152602001610209565b34801561037057600080fd5b506102b261037f3660046151c6565b610946565b34801561039057600080fd5b5061030361039f36600461537d565b610aba565b3480156103b057600080fd5b506102b26103bf366004615628565b610f73565b3480156103d057600080fd5b506103da6101f481565b60405163ffffffff9091168152602001610209565b3480156103fb57600080fd5b5061040f61040a36600461560f565b61135e565b6040519015158152602001610209565b34801561042b57600080fd5b506102b261043a3660046151c6565b6114ff565b34801561044b57600080fd5b506102b261045a3660046151c6565b6116b0565b34801561046b57600080fd5b506102c9606481565b34801561048057600080fd5b506102b261048f3660046151e3565b611767565b3480156104a057600080fd5b5061033b7f000000000000000000000000000000000000000000000000000000000000000081565b3480156104d457600080fd5b5060035461033b906001600160a01b031681565b3480156104f457600080fd5b506102b26117c7565b34801561050957600080fd5b506102b2610518366004615277565b611871565b34801561052957600080fd5b506102b26105383660046151c6565b61196a565b34801561054957600080fd5b50600a54610303906001600160601b031681565b34801561056957600080fd5b506000546001600160a01b031661033b565b6102b261058936600461560f565b611a76565b34801561059a57600080fd5b506105ae6105a936600461545a565b611bba565b604051908152602001610209565b3480156105c857600080fd5b506007546105dc906001600160401b031681565b6040516102099190615b4e565b3480156105f557600080fd5b506105ae611f2a565b34801561060a57600080fd5b506102b261061936600461521c565b612178565b34801561062a57600080fd5b506102b261063936600461560f565b612315565b34801561064a57600080fd5b5061065e61065936600461564d565b612378565b60405161020991906158d0565b34801561067757600080fd5b506102b2610686366004615571565b612479565b34801561069757600080fd5b506102b26106a636600461560f565b6125ed565b3480156106b757600080fd5b506102b26106c6366004615628565b612711565b3480156106d757600080fd5b506105ae6106e636600461529f565b6128a8565b3480156106f757600080fd5b506102b2610706366004615628565b6128d8565b6102b26107193660046152f1565b612bc5565b34801561072a57600080fd5b506105ae61073936600461560f565b612ed6565b34801561074a57600080fd5b506102b2610759366004615628565b612ef7565b34801561076a57600080fd5b5061077e61077936600461560f565b613007565b604051610209959493929190615b62565b34801561079b57600080fd5b506102b26107aa3660046151c6565b613102565b3480156107bb57600080fd5b506105ae6107ca36600461560f565b600f6020526000908152604090205481565b3480156107e857600080fd5b506102b26107f73660046151c6565b6132dd565b600c54600e805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff1693919283919083018282801561086657602002820191906000526020600020905b815481526020019060010190808311610852575b50505050509050925092509250909192565b60008281526005602052604090205482906001600160a01b0316806108b057604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b038216146108e45780604051636c51fda960e11b81526004016108db9190615859565b60405180910390fd5b600c54600160301b900460ff161561090f5760405163769dd35360e11b815260040160405180910390fd5b6109188461135e565b1561093657604051631685ecdd60e31b815260040160405180910390fd5b61094084846132ee565b50505050565b600c54600160301b900460ff16156109715760405163769dd35360e11b815260040160405180910390fd5b6109796134a9565b600b54600160601b90046001600160601b03166109a957604051631e9acf1760e31b815260040160405180910390fd5b600b8054600160601b90046001600160601b0316908190600c6109cc8380615d6a565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600a600c8282829054906101000a90046001600160601b0316610a149190615d6a565b92506101000a8154816001600160601b0302191690836001600160601b031602179055506000826001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d8060008114610a8e576040519150601f19603f3d011682016040523d82523d6000602084013e610a93565b606091505b5050905080610ab55760405163950b247960e01b815260040160405180910390fd5b505050565b600c54600090600160301b900460ff1615610ae85760405163769dd35360e11b815260040160405180910390fd5b60005a90506000610af985856134fe565b90506000846060015163ffffffff166001600160401b03811115610b1f57610b1f615e9c565b604051908082528060200260200182016040528015610b48578160200160208202803683370190505b50905060005b856060015163ffffffff16811015610bbf57826040015181604051602001610b779291906158e3565b6040516020818303038152906040528051906020012060001c828281518110610ba257610ba2615e86565b602090810291909101015280610bb781615dee565b915050610b4e565b50602080830180516000908152600f9092526040808320839055905190518291631fe543e360e01b91610bf791908690602401615a3f565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252600c805460ff60301b1916600160301b179055908801516080890151919250600091610c5c9163ffffffff169084613769565b600c805460ff60301b19169055602089810151600090815260069091526040902054909150600160c01b90046001600160401b0316610c9c816001615cdc565b6020808b0151600090815260069091526040812080546001600160401b0393909316600160c01b026001600160c01b039093169290921790915560a08a01518051610ce990600190615d53565b81518110610cf957610cf9615e86565b602091010151600c5460f89190911c6001149150600090610d2a908a90600160581b900463ffffffff163a856137b7565b90508115610e22576020808c01516000908152600690915260409020546001600160601b03808316600160601b909204161015610d7a57604051631e9acf1760e31b815260040160405180910390fd5b60208b81015160009081526006909152604090208054829190600c90610db1908490600160601b90046001600160601b0316615d6a565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600b600c8282829054906101000a90046001600160601b0316610df99190615cfe565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550610efd565b6020808c01516000908152600690915260409020546001600160601b0380831691161015610e6357604051631e9acf1760e31b815260040160405180910390fd5b6020808c015160009081526006909152604081208054839290610e909084906001600160601b0316615d6a565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600b60008282829054906101000a90046001600160601b0316610ed89190615cfe565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b8a6020015188602001517f49580fdfd9497e1ed5c1b1cec0495087ae8e3f1267470ec2fb015db32e3d6aa78a604001518488604051610f5a939291909283526001600160601b039190911660208301521515604082015260600190565b60405180910390a3985050505050505050505b92915050565b600c54600160301b900460ff1615610f9e5760405163769dd35360e11b815260040160405180910390fd5b610fa781613806565b610fc65780604051635428d44960e01b81526004016108db9190615859565b600080600080610fd586613007565b945094505093509350336001600160a01b0316826001600160a01b0316146110385760405162461bcd60e51b81526020600482015260166024820152752737ba1039bab139b1b934b83a34b7b71037bbb732b960511b60448201526064016108db565b6110418661135e565b156110875760405162461bcd60e51b815260206004820152601660248201527550656e64696e6720726571756573742065786973747360501b60448201526064016108db565b60006040518060c0016040528061109c600290565b60ff168152602001888152602001846001600160a01b03168152602001838152602001866001600160601b03168152602001856001600160601b031681525090506000816040516020016110f09190615922565b604051602081830303815290604052905061110a88613870565b505060405163ce3f471960e01b81526001600160a01b0388169063ce3f4719906001600160601b0388169061114390859060040161590f565b6000604051808303818588803b15801561115c57600080fd5b505af1158015611170573d6000803e3d6000fd5b50506002546001600160a01b031615801593509150611199905057506001600160601b03861615155b156112635760025460405163a9059cbb60e01b81526001600160a01b039091169063a9059cbb906111d0908a908a906004016158a0565b602060405180830381600087803b1580156111ea57600080fd5b505af11580156111fe573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061122291906152bb565b6112635760405162461bcd60e51b8152602060048201526012602482015271696e73756666696369656e742066756e647360701b60448201526064016108db565b600c805460ff60301b1916600160301b17905560005b835181101561130c5783818151811061129457611294615e86565b60200260200101516001600160a01b0316638ea98117896040518263ffffffff1660e01b81526004016112c79190615859565b600060405180830381600087803b1580156112e157600080fd5b505af11580156112f5573d6000803e3d6000fd5b50505050808061130490615dee565b915050611279565b50600c805460ff60301b191690556040517fd63ca8cb945956747ee69bfdc3ea754c24a4caf7418db70e46052f7850be41879061134c9089908b9061586d565b60405180910390a15050505050505050565b6000818152600560209081526040808320815160608101835281546001600160a01b03908116825260018301541681850152600282018054845181870281018701865281815287969395860193909291908301828280156113e857602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116113ca575b505050505081525050905060005b8160400151518110156114f55760005b600e548110156114e25760006114ab600e838154811061142857611428615e86565b90600052602060002001548560400151858151811061144957611449615e86565b602002602001015188600460008960400151898151811061146c5761146c615e86565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081208d82529092529020546001600160401b0316613abe565b506000818152600f6020526040902054909150156114cf5750600195945050505050565b50806114da81615dee565b915050611406565b50806114ed81615dee565b9150506113f6565b5060009392505050565b600c54600160301b900460ff161561152a5760405163769dd35360e11b815260040160405180910390fd5b6115326134a9565b6002546001600160a01b031661155b5760405163c1f0c0a160e01b815260040160405180910390fd5b600b546001600160601b031661158457604051631e9acf1760e31b815260040160405180910390fd5b600b80546001600160601b031690819060006115a08380615d6a565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600a60008282829054906101000a90046001600160601b03166115e89190615d6a565b82546001600160601b039182166101009390930a92830291909202199091161790555060025460405163a9059cbb60e01b81526001600160a01b039091169063a9059cbb9061163d90859085906004016158a0565b602060405180830381600087803b15801561165757600080fd5b505af115801561166b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061168f91906152bb565b6116ac57604051631e9acf1760e31b815260040160405180910390fd5b5050565b6116b86134a9565b6116c181613806565b156116e1578060405163ac8a27ef60e01b81526004016108db9190615859565b601280546001810182556000919091527fbb8a6a4669ba250d26cd7a459eca9d215f8307e33aebe50379bc5a3617ec34440180546001600160a01b0319166001600160a01b0383161790556040517fb7cabbfc11e66731fc77de0444614282023bcbd41d16781c753a431d0af016259061175c908390615859565b60405180910390a150565b61176f6134a9565b6002546001600160a01b03161561179957604051631688c53760e11b815260040160405180910390fd5b600280546001600160a01b039384166001600160a01b03199182161790915560038054929093169116179055565b6001546001600160a01b0316331461181a5760405162461bcd60e51b815260206004820152601660248201527526bab9ba10313290383937b837b9b2b21037bbb732b960511b60448201526064016108db565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6118796134a9565b6040805180820182526000916118a89190849060029083908390808284376000920191909152506128a8915050565b6000818152600d602052604090205490915060ff16156118de57604051634a0b8fa760e01b8152600481018290526024016108db565b6000818152600d6020526040808220805460ff19166001908117909155600e805491820181559092527fbb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fd909101829055517fc9583fd3afa3d7f16eb0b88d0268e7d05c09bafa4b21e092cbd1320e1bc8089d9061195e9083815260200190565b60405180910390a15050565b6119726134a9565b600a544790600160601b90046001600160601b0316818111156119ac5780826040516354ced18160e11b81526004016108db9291906158e3565b81811015610ab55760006119c08284615d53565b90506000846001600160a01b03168260405160006040518083038185875af1925050503d8060008114611a0f576040519150601f19603f3d011682016040523d82523d6000602084013e611a14565b606091505b5050905080611a365760405163950b247960e01b815260040160405180910390fd5b7f4aed7c8eed0496c8c19ea2681fcca25741c1602342e38b045d9f1e8e905d2e9c8583604051611a6792919061586d565b60405180910390a15050505050565b600c54600160301b900460ff1615611aa15760405163769dd35360e11b815260040160405180910390fd5b6000818152600560205260409020546001600160a01b0316611ad657604051630fb532db60e11b815260040160405180910390fd5b60008181526006602052604090208054600160601b90046001600160601b0316903490600c611b058385615cfe565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555034600a600c8282829054906101000a90046001600160601b0316611b4d9190615cfe565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f7603b205d03651ee812f803fccde89f1012e545a9c99f0abfea9cedd0fd8e902823484611ba09190615cc4565b604051611bae9291906158e3565b60405180910390a25050565b600c54600090600160301b900460ff1615611be85760405163769dd35360e11b815260040160405180910390fd5b6020808301356000908152600590915260409020546001600160a01b0316611c2357604051630fb532db60e11b815260040160405180910390fd5b3360009081526004602090815260408083208583013584529091529020546001600160401b031680611c70578260200135336040516379bfd40160e01b81526004016108db929190615a14565b600c5461ffff16611c876060850160408601615556565b61ffff161080611caa575060c8611ca46060850160408601615556565b61ffff16115b15611ce457611cbf6060840160408501615556565b600c5460405163539c34bb60e11b81526108db929161ffff169060c890600401615997565b600c5462010000900463ffffffff16611d03608085016060860161566f565b63ffffffff161115611d4957611d1f608084016060850161566f565b600c54604051637aebf00f60e11b81526108db929162010000900463ffffffff1690600401615b37565b6101f4611d5c60a085016080860161566f565b63ffffffff161115611d9657611d7860a084016080850161566f565b6101f46040516311ce1afb60e21b81526004016108db929190615b37565b6000611da3826001615cdc565b9050600080611db9863533602089013586613abe565b90925090506000611dd5611dd060a0890189615bb7565b613b47565b90506000611de282613bc4565b905083611ded613c35565b60208a0135611e0260808c0160608d0161566f565b611e1260a08d0160808e0161566f565b3386604051602001611e2a9796959493929190615a97565b60405160208183030381529060405280519060200120600f600086815260200190815260200160002081905550336001600160a01b0316886020013589600001357feb0e3652e0f44f417695e6e90f2f42c99b65cd7169074c5a654b16b9748c3a4e87878d6040016020810190611ea19190615556565b8e6060016020810190611eb4919061566f565b8f6080016020810190611ec7919061566f565b89604051611eda96959493929190615a58565b60405180910390a45050336000908152600460209081526040808320898301358452909152902080546001600160401b0319166001600160401b039490941693909317909255925050505b919050565b600c54600090600160301b900460ff1615611f585760405163769dd35360e11b815260040160405180910390fd5b600033611f66600143615d53565b600754604051606093841b6001600160601b03199081166020830152924060348201523090931b909116605483015260c01b6001600160c01b031916606882015260700160408051601f198184030181529190528051602090910120600780549192506001600160401b03909116906000611fe083615e09565b91906101000a8154816001600160401b0302191690836001600160401b03160217905550506000806001600160401b0381111561201f5761201f615e9c565b604051908082528060200260200182016040528015612048578160200160208202803683370190505b506040805160608082018352600080835260208084018281528486018381528984526006835286842095518654925191516001600160601b039182166001600160c01b031990941693909317600160601b9190921602176001600160c01b0316600160c01b6001600160401b039092169190910217909355835191820184523382528183018181528285018681528883526005855294909120825181546001600160a01b03199081166001600160a01b0392831617835592516001830180549094169116179091559251805194955090936121299260028501920190614e37565b5061213991506008905083613cc5565b50817f1d3015d7ba850fa198dc7b1a3f5d42779313a681035f77c8c03764c61005518d3360405161216a9190615859565b60405180910390a250905090565b600c54600160301b900460ff16156121a35760405163769dd35360e11b815260040160405180910390fd5b6002546001600160a01b031633146121ce576040516344b0e3c360e01b815260040160405180910390fd5b602081146121ef57604051638129bbcd60e01b815260040160405180910390fd5b60006121fd8284018461560f565b6000818152600560205260409020549091506001600160a01b031661223557604051630fb532db60e11b815260040160405180910390fd5b600081815260066020526040812080546001600160601b03169186919061225c8385615cfe565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600a60008282829054906101000a90046001600160601b03166122a49190615cfe565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550817f1ced9348ff549fceab2ac57cd3a9de38edaaab274b725ee82c23e8fc8c4eec7a8287846122f79190615cc4565b6040516123059291906158e3565b60405180910390a2505050505050565b61231d6134a9565b6000818152600560205260409020546001600160a01b031661235257604051630fb532db60e11b815260040160405180910390fd5b6000818152600560205260409020546123759082906001600160a01b03166132ee565b50565b606060006123866008613cd1565b90508084106123a857604051631390f2a160e01b815260040160405180910390fd5b60006123b48486615cc4565b9050818111806123c2575083155b6123cc57806123ce565b815b905060006123dc8683615d53565b6001600160401b038111156123f3576123f3615e9c565b60405190808252806020026020018201604052801561241c578160200160208202803683370190505b50905060005b815181101561246f576124406124388883615cc4565b600890613cdb565b82828151811061245257612452615e86565b60209081029190910101528061246781615dee565b915050612422565b5095945050505050565b6124816134a9565b60c861ffff871611156124ae57858660c860405163539c34bb60e11b81526004016108db93929190615997565b600082136124d2576040516321ea67b360e11b8152600481018390526024016108db565b6040805160a0808201835261ffff891680835263ffffffff89811660208086018290526000868801528a831660608088018290528b85166080988901819052600c805465ffffffffffff1916881762010000870217600160301b600160781b031916600160381b850263ffffffff60581b191617600160581b83021790558a51601180548d8701519289166001600160401b031990911617600160201b92891692909202919091179081905560108d90558a519788528785019590955298860191909152840196909652938201879052838116928201929092529190921c90911660c08201527f777357bb93f63d088f18112d3dba38457aec633eb8f1341e1d418380ad328e789060e00160405180910390a1505050505050565b600c54600160301b900460ff16156126185760405163769dd35360e11b815260040160405180910390fd5b6000818152600560205260409020546001600160a01b031661264d57604051630fb532db60e11b815260040160405180910390fd5b6000818152600560205260409020600101546001600160a01b031633146126a4576000818152600560205260409081902060010154905163d084e97560e01b81526108db916001600160a01b031690600401615859565b6000818152600560205260409081902080546001600160a01b031980821633908117845560019093018054909116905591516001600160a01b039092169183917fd4114ab6e9af9f597c52041f32d62dc57c5c4e4c0d4427006069635e216c938691611bae918591615886565b60008281526005602052604090205482906001600160a01b03168061274957604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b038216146127745780604051636c51fda960e11b81526004016108db9190615859565b600c54600160301b900460ff161561279f5760405163769dd35360e11b815260040160405180910390fd5b600084815260056020526040902060020154606414156127d2576040516305a48e0f60e01b815260040160405180910390fd5b6001600160a01b03831660009081526004602090815260408083208784529091529020546001600160401b03161561280957610940565b6001600160a01b0383166000818152600460209081526040808320888452825280832080546001600160401b031916600190811790915560058352818420600201805491820181558452919092200180546001600160a01b0319169092179091555184907f1e980d04aa7648e205713e5e8ea3808672ac163d10936d36f91b2c88ac1575e19061289a908690615859565b60405180910390a250505050565b6000816040516020016128bb91906158c2565b604051602081830303815290604052805190602001209050919050565b60008281526005602052604090205482906001600160a01b03168061291057604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b0382161461293b5780604051636c51fda960e11b81526004016108db9190615859565b600c54600160301b900460ff16156129665760405163769dd35360e11b815260040160405180910390fd5b61296f8461135e565b1561298d57604051631685ecdd60e31b815260040160405180910390fd5b6001600160a01b03831660009081526004602090815260408083208784529091529020546001600160401b03166129db5783836040516379bfd40160e01b81526004016108db929190615a14565b600084815260056020908152604080832060020180548251818502810185019093528083529192909190830182828015612a3e57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612a20575b50505050509050600060018251612a559190615d53565b905060005b8251811015612b6157856001600160a01b0316838281518110612a7f57612a7f615e86565b60200260200101516001600160a01b03161415612b4f576000838381518110612aaa57612aaa615e86565b6020026020010151905080600560008a81526020019081526020016000206002018381548110612adc57612adc615e86565b600091825260208083209190910180546001600160a01b0319166001600160a01b039490941693909317909255898152600590915260409020600201805480612b2757612b27615e70565b600082815260209020810160001990810180546001600160a01b031916905501905550612b61565b80612b5981615dee565b915050612a5a565b506001600160a01b03851660009081526004602090815260408083208984529091529081902080546001600160401b03191690555186907f32158c6058347c1601b2d12bc696ac6901d8a9a9aa3ba10c27ab0a983e8425a790612305908890615859565b6000612bd382840184615494565b9050806000015160ff16600114612c0c57805160405163237d181f60e21b815260ff9091166004820152600160248201526044016108db565b8060a001516001600160601b03163414612c505760a08101516040516306acf13560e41b81523460048201526001600160601b0390911660248201526044016108db565b6020808201516000908152600590915260409020546001600160a01b031615612c8c576040516326afa43560e11b815260040160405180910390fd5b60005b816060015151811015612d2c5760016004600084606001518481518110612cb857612cb8615e86565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060008460200151815260200190815260200160002060006101000a8154816001600160401b0302191690836001600160401b031602179055508080612d2490615dee565b915050612c8f565b50604080516060808201835260808401516001600160601b03908116835260a0850151811660208085019182526000858701818152828901805183526006845288832097518854955192516001600160401b0316600160c01b026001600160c01b03938816600160601b026001600160c01b0319909716919097161794909417169390931790945584518084018652868601516001600160a01b03908116825281860184815294880151828801908152925184526005865295909220825181549087166001600160a01b0319918216178255935160018201805491909716941693909317909455925180519192612e2b92600285019290910190614e37565b5050506080810151600a8054600090612e4e9084906001600160601b0316615cfe565b92506101000a8154816001600160601b0302191690836001600160601b031602179055508060a00151600a600c8282829054906101000a90046001600160601b0316612e9a9190615cfe565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555061094081602001516008613cc590919063ffffffff16565b600e8181548110612ee657600080fd5b600091825260209091200154905081565b60008281526005602052604090205482906001600160a01b031680612f2f57604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b03821614612f5a5780604051636c51fda960e11b81526004016108db9190615859565b600c54600160301b900460ff1615612f855760405163769dd35360e11b815260040160405180910390fd5b6000848152600560205260409020600101546001600160a01b03848116911614610940576000848152600560205260409081902060010180546001600160a01b0319166001600160a01b0386161790555184907f21a4dad170a6bf476c31bbcf4a16628295b0e450672eec25d7c93308e05344a19061289a9033908790615886565b6000818152600560205260408120548190819081906060906001600160a01b031661304557604051630fb532db60e11b815260040160405180910390fd5b60008681526006602090815260408083205460058352928190208054600290910180548351818602810186019094528084526001600160601b0380871696600160601b810490911695600160c01b9091046001600160401b0316946001600160a01b03909416939183918301828280156130e857602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116130ca575b505050505090509450945094509450945091939590929450565b61310a6134a9565b6002546001600160a01b03166131335760405163c1f0c0a160e01b815260040160405180910390fd5b6002546040516370a0823160e01b81526000916001600160a01b0316906370a0823190613164903090600401615859565b60206040518083038186803b15801561317c57600080fd5b505afa158015613190573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906131b491906152d8565b600a549091506001600160601b0316818111156131e85780826040516354ced18160e11b81526004016108db9291906158e3565b81811015610ab55760006131fc8284615d53565b60025460405163a9059cbb60e01b81529192506001600160a01b03169063a9059cbb9061322f908790859060040161586d565b602060405180830381600087803b15801561324957600080fd5b505af115801561325d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061328191906152bb565b61329e57604051631f01ff1360e21b815260040160405180910390fd5b7f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b43660084826040516132cf92919061586d565b60405180910390a150505050565b6132e56134a9565b61237581613ce7565b6000806132fa84613870565b60025491935091506001600160a01b03161580159061332157506001600160601b03821615155b156133d05760025460405163a9059cbb60e01b81526001600160a01b039091169063a9059cbb906133619086906001600160601b0387169060040161586d565b602060405180830381600087803b15801561337b57600080fd5b505af115801561338f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133b391906152bb565b6133d057604051631e9acf1760e31b815260040160405180910390fd5b6000836001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d8060008114613426576040519150601f19603f3d011682016040523d82523d6000602084013e61342b565b606091505b505090508061344d5760405163950b247960e01b815260040160405180910390fd5b604080516001600160a01b03861681526001600160601b038581166020830152841681830152905186917f8c74ce8b8cf87f5eb001275c8be27eb34ea2b62bfab6814fcc62192bb63e81c4919081900360600190a25050505050565b6000546001600160a01b031633146134fc5760405162461bcd60e51b815260206004820152601660248201527527b7363c9031b0b63630b1363290313c9037bbb732b960511b60448201526064016108db565b565b6040805160608101825260008082526020820181905291810191909152600061352a84600001516128a8565b6000818152600d602052604090205490915060ff1661355f57604051631dfd6e1360e21b8152600481018290526024016108db565b60008185608001516040516020016135789291906158e3565b60408051601f1981840301815291815281516020928301206000818152600f909352912054909150806135be57604051631b44092560e11b815260040160405180910390fd5b845160208087015160408089015160608a015160808b015160a08c015193516135ed978a979096959101615ae3565b6040516020818303038152906040528051906020012081146136225760405163354a450b60e21b815260040160405180910390fd5b60006136318660000151613d8b565b9050806136f8578551604051631d2827a760e31b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169163e9413d38916136859190600401615b4e565b60206040518083038186803b15801561369d57600080fd5b505afa1580156136b1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136d591906152d8565b9050806136f857855160405163175dadad60e01b81526108db9190600401615b4e565b600087608001518260405160200161371a929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c905060006137418983613e68565b6040805160608101825297885260208801969096529486019490945250929695505050505050565b60005a61138881101561377b57600080fd5b61138881039050846040820482031161379357600080fd5b50823b61379f57600080fd5b60008083516020850160008789f190505b9392505050565b600081156137e4576011546137dd9086908690600160201b900463ffffffff1686613ed3565b90506137fe565b6011546137fb908690869063ffffffff1686613f75565b90505b949350505050565b6000805b60125481101561386757826001600160a01b03166012828154811061383157613831615e86565b6000918252602090912001546001600160a01b031614156138555750600192915050565b8061385f81615dee565b91505061380a565b50600092915050565b6000818152600560209081526040808320815160608101835281546001600160a01b039081168252600183015416818501526002820180548451818702810187018652818152879687969495948601939192908301828280156138fc57602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116138de575b505050919092525050506000858152600660209081526040808320815160608101835290546001600160601b03808216808452600160601b8304909116948301859052600160c01b9091046001600160401b0316928201929092529096509094509192505b8260400151518110156139d857600460008460400151838151811061398857613988615e86565b6020908102919091018101516001600160a01b031682528181019290925260409081016000908120898252909252902080546001600160401b0319169055806139d081615dee565b915050613961565b50600085815260056020526040812080546001600160a01b03199081168255600182018054909116905590613a106002830182614e9c565b5050600085815260066020526040812055613a2c60088661409a565b50600a8054859190600090613a4b9084906001600160601b0316615d6a565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555082600a600c8282829054906101000a90046001600160601b0316613a939190615d6a565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505050915091565b60408051602081018690526001600160a01b03851691810191909152606081018390526001600160401b03821660808201526000908190819060a00160408051601f198184030181529082905280516020918201209250613b239189918491016158e3565b60408051808303601f19018152919052805160209091012097909650945050505050565b60408051602081019091526000815281613b705750604080516020810190915260008152610f6d565b63125fa26760e31b613b828385615d92565b6001600160e01b03191614613baa57604051632923fee760e11b815260040160405180910390fd5b613bb78260048186615c9a565b8101906137b09190615332565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa82604051602401613bfd91511515815260200190565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b03199093169290921790915292915050565b600046613c41816140a6565b15613cbe5760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015613c8057600080fd5b505afa158015613c94573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613cb891906152d8565b91505090565b4391505090565b60006137b083836140c9565b6000610f6d825490565b60006137b08383614118565b6001600160a01b038116331415613d3a5760405162461bcd60e51b815260206004820152601760248201527621b0b73737ba103a3930b739b332b9103a379039b2b63360491b60448201526064016108db565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600046613d97816140a6565b15613e5957610100836001600160401b0316613db1613c35565b613dbb9190615d53565b1180613dd75750613dca613c35565b836001600160401b031610155b15613de55750600092915050565b6040516315a03d4160e11b8152606490632b407a8290613e09908690600401615b4e565b60206040518083038186803b158015613e2157600080fd5b505afa158015613e35573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137b091906152d8565b50506001600160401b03164090565b6000613e9c8360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151614142565b60038360200151604051602001613eb4929190615a2b565b60408051601f1981840301815291905280516020909101209392505050565b600080613f166000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061435d92505050565b905060005a613f258888615cc4565b613f2f9190615d53565b613f399085615d34565b90506000613f5263ffffffff871664e8d4a51000615d34565b905082613f5f8284615cc4565b613f699190615cc4565b98975050505050505050565b600080613f80614422565b905060008113613fa6576040516321ea67b360e11b8152600481018290526024016108db565b6000613fe86000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061435d92505050565b9050600082825a613ff98b8b615cc4565b6140039190615d53565b61400d9088615d34565b6140179190615cc4565b61402990670de0b6b3a7640000615d34565b6140339190615d20565b9050600061404c63ffffffff881664e8d4a51000615d34565b905061406381676765c793fa10079d601b1b615d53565b8211156140835760405163e80fa38160e01b815260040160405180910390fd5b61408d8183615cc4565b9998505050505050505050565b60006137b083836144ed565b600061a4b18214806140ba575062066eed82145b80610f6d57505062066eee1490565b600081815260018301602052604081205461411057508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610f6d565b506000610f6d565b600082600001828154811061412f5761412f615e86565b9060005260206000200154905092915050565b61414b896145e0565b6141945760405162461bcd60e51b815260206004820152601a6024820152797075626c6963206b6579206973206e6f74206f6e20637572766560301b60448201526064016108db565b61419d886145e0565b6141e15760405162461bcd60e51b815260206004820152601560248201527467616d6d61206973206e6f74206f6e20637572766560581b60448201526064016108db565b6141ea836145e0565b6142365760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e20637572766500000060448201526064016108db565b61423f826145e0565b61428a5760405162461bcd60e51b815260206004820152601c60248201527b73486173685769746e657373206973206e6f74206f6e20637572766560201b60448201526064016108db565b614296878a88876146a3565b6142de5760405162461bcd60e51b81526020600482015260196024820152786164647228632a706b2b732a6729213d5f755769746e65737360381b60448201526064016108db565b60006142ea8a876147b7565b905060006142fd898b878b86898961481b565b9050600061430e838d8d8a8661492e565b9050808a1461434f5760405162461bcd60e51b815260206004820152600d60248201526c34b73b30b634b210383937b7b360991b60448201526064016108db565b505050505050505050505050565b600046614369816140a6565b156143a857606c6001600160a01b031663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b158015613e2157600080fd5b6143b18161496e565b1561386757600f602160991b016001600160a01b03166349948e0e84604051806080016040528060488152602001615ed6604891396040516020016143f79291906157af565b6040516020818303038152906040526040518263ffffffff1660e01b8152600401613e09919061590f565b600c5460035460408051633fabe5a360e21b81529051600093600160381b900463ffffffff169283151592859283926001600160a01b03169163feaf968c9160048083019260a0929190829003018186803b15801561448057600080fd5b505afa158015614494573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906144b8919061568a565b5094509092508491505080156144dc57506144d38242615d53565b8463ffffffff16105b156137fe5750601054949350505050565b600081815260018301602052604081205480156145d6576000614511600183615d53565b855490915060009061452590600190615d53565b905081811461458a57600086600001828154811061454557614545615e86565b906000526020600020015490508087600001848154811061456857614568615e86565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061459b5761459b615e70565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610f6d565b6000915050610f6d565b80516000906401000003d0191161462e5760405162461bcd60e51b8152602060048201526012602482015271696e76616c696420782d6f7264696e61746560701b60448201526064016108db565b60208201516401000003d0191161467c5760405162461bcd60e51b8152602060048201526012602482015271696e76616c696420792d6f7264696e61746560701b60448201526064016108db565b60208201516401000003d01990800961469c8360005b60200201516149a8565b1492915050565b60006001600160a01b0382166146e95760405162461bcd60e51b815260206004820152600b60248201526a626164207769746e65737360a81b60448201526064016108db565b60208401516000906001161561470057601c614703565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe199182039250600091908909875160408051600080825260209091019182905292935060019161476d918691889187906158f1565b6020604051602081039080840390855afa15801561478f573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b6147bf614eba565b6147ec600184846040516020016147d893929190615838565b6040516020818303038152906040526149cc565b90505b6147f8816145e0565b610f6d57805160408051602081019290925261481491016147d8565b90506147ef565b614823614eba565b825186516401000003d01990819006910614156148825760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e6374000060448201526064016108db565b61488d878988614a1a565b6148d25760405162461bcd60e51b8152602060048201526016602482015275119a5c9cdd081b5d5b0818da1958dac819985a5b195960521b60448201526064016108db565b6148dd848685614a1a565b6149235760405162461bcd60e51b815260206004820152601760248201527614d958dbdb99081b5d5b0818da1958dac819985a5b1959604a1b60448201526064016108db565b613f69868484614b35565b60006002868686858760405160200161494c969594939291906157de565b60408051601f1981840301815291905280516020909101209695505050505050565b6000600a82148061498057506101a482145b8061498d575062aa37dc82145b80614999575061210582145b80610f6d57505062014a331490565b6000806401000003d01980848509840990506401000003d019600782089392505050565b6149d4614eba565b6149dd82614bf8565b81526149f26149ed826000614692565b614c33565b602082018190526002900660011415611f25576020810180516401000003d019039052919050565b600082614a575760405162461bcd60e51b815260206004820152600b60248201526a3d32b9379039b1b0b630b960a91b60448201526064016108db565b83516020850151600090614a6d90600290615e30565b15614a7957601c614a7c565b601b5b9050600070014551231950b75fc4402da1732fc9bebe19838709604080516000808252602090910191829052919250600190614abf9083908690889087906158f1565b6020604051602081039080840390855afa158015614ae1573d6000803e3d6000fd5b505050602060405103519050600086604051602001614b00919061579d565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b614b3d614eba565b835160208086015185519186015160009384938493614b5e93909190614c53565b919450925090506401000003d019858209600114614bba5760405162461bcd60e51b815260206004820152601960248201527834b73b2d1036bab9ba1031329034b73b32b939b29037b3103d60391b60448201526064016108db565b60405180604001604052806401000003d01980614bd957614bd9615e5a565b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d0198110611f2557604080516020808201939093528151808203840181529082019091528051910120614c00565b6000610f6d826002614c4c6401000003d0196001615cc4565b901c614d33565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000614c9383838585614dca565b9098509050614ca488828e88614dee565b9098509050614cb588828c87614dee565b90985090506000614cc88d878b85614dee565b9098509050614cd988828686614dca565b9098509050614cea88828e89614dee565b9098509050818114614d1f576401000003d019818a0998506401000003d01982890997506401000003d0198183099650614d23565b8196505b5050505050509450945094915050565b600080614d3e614ed8565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152614d70614ef6565b60208160c0846005600019fa925082614dc05760405162461bcd60e51b81526020600482015260126024820152716269674d6f64457870206661696c7572652160701b60448201526064016108db565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215614e8c579160200282015b82811115614e8c57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190614e57565b50614e98929150614f14565b5090565b50805460008255906000526020600020908101906123759190614f14565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b80821115614e985760008155600101614f15565b8035611f2581615eb2565b600082601f830112614f4557600080fd5b813560206001600160401b03821115614f6057614f60615e9c565b8160051b614f6f828201615c6a565b838152828101908684018388018501891015614f8a57600080fd5b600093505b85841015614fb6578035614fa281615eb2565b835260019390930192918401918401614f8f565b50979650505050505050565b600082601f830112614fd357600080fd5b614fdb615bfd565b808385604086011115614fed57600080fd5b60005b600281101561500f578135845260209384019390910190600101614ff0565b509095945050505050565b60008083601f84011261502c57600080fd5b5081356001600160401b0381111561504357600080fd5b60208301915083602082850101111561505b57600080fd5b9250929050565b600082601f83011261507357600080fd5b81356001600160401b0381111561508c5761508c615e9c565b61509f601f8201601f1916602001615c6a565b8181528460208386010111156150b457600080fd5b816020850160208301376000918101602001919091529392505050565b600060c082840312156150e357600080fd5b6150eb615c25565b905081356001600160401b03808216821461510557600080fd5b8183526020840135602084015261511e60408501615184565b604084015261512f60608501615184565b606084015261514060808501614f29565b608084015260a084013591508082111561515957600080fd5b5061516684828501615062565b60a08301525092915050565b803561ffff81168114611f2557600080fd5b803563ffffffff81168114611f2557600080fd5b80516001600160501b0381168114611f2557600080fd5b80356001600160601b0381168114611f2557600080fd5b6000602082840312156151d857600080fd5b81356137b081615eb2565b600080604083850312156151f657600080fd5b823561520181615eb2565b9150602083013561521181615eb2565b809150509250929050565b6000806000806060858703121561523257600080fd5b843561523d81615eb2565b93506020850135925060408501356001600160401b0381111561525f57600080fd5b61526b8782880161501a565b95989497509550505050565b60006040828403121561528957600080fd5b8260408301111561529957600080fd5b50919050565b6000604082840312156152b157600080fd5b6137b08383614fc2565b6000602082840312156152cd57600080fd5b81516137b081615ec7565b6000602082840312156152ea57600080fd5b5051919050565b6000806020838503121561530457600080fd5b82356001600160401b0381111561531a57600080fd5b6153268582860161501a565b90969095509350505050565b60006020828403121561534457600080fd5b604051602081016001600160401b038111828210171561536657615366615e9c565b604052823561537481615ec7565b81529392505050565b6000808284036101c081121561539257600080fd5b6101a0808212156153a257600080fd5b6153aa615c47565b91506153b68686614fc2565b82526153c58660408701614fc2565b60208301526080850135604083015260a0850135606083015260c085013560808301526153f460e08601614f29565b60a083015261010061540887828801614fc2565b60c084015261541b876101408801614fc2565b60e0840152610180860135908301529092508301356001600160401b0381111561544457600080fd5b615450858286016150d1565b9150509250929050565b60006020828403121561546c57600080fd5b81356001600160401b0381111561548257600080fd5b820160c081850312156137b057600080fd5b6000602082840312156154a657600080fd5b81356001600160401b03808211156154bd57600080fd5b9083019060c082860312156154d157600080fd5b6154d9615c25565b823560ff811681146154ea57600080fd5b81526020838101359082015261550260408401614f29565b604082015260608301358281111561551957600080fd5b61552587828601614f34565b606083015250615537608084016151af565b608082015261554860a084016151af565b60a082015295945050505050565b60006020828403121561556857600080fd5b6137b082615172565b60008060008060008086880360e081121561558b57600080fd5b61559488615172565b96506155a260208901615184565b95506155b060408901615184565b94506155be60608901615184565b9350608088013592506040609f19820112156155d957600080fd5b506155e2615bfd565b6155ee60a08901615184565b81526155fc60c08901615184565b6020820152809150509295509295509295565b60006020828403121561562157600080fd5b5035919050565b6000806040838503121561563b57600080fd5b82359150602083013561521181615eb2565b6000806040838503121561566057600080fd5b50508035926020909101359150565b60006020828403121561568157600080fd5b6137b082615184565b600080600080600060a086880312156156a257600080fd5b6156ab86615198565b94506020860151935060408601519250606086015191506156ce60808701615198565b90509295509295909350565b600081518084526020808501945080840160005b838110156157135781516001600160a01b0316875295820195908201906001016156ee565b509495945050505050565b8060005b6002811015610940578151845260209384019390910190600101615722565b600081518084526020808501945080840160005b8381101561571357815187529582019590820190600101615755565b60008151808452615789816020860160208601615dc2565b601f01601f19169290920160200192915050565b6157a7818361571e565b604001919050565b600083516157c1818460208801615dc2565b8351908301906157d5818360208801615dc2565b01949350505050565b8681526157ee602082018761571e565b6157fb606082018661571e565b61580860a082018561571e565b61581560e082018461571e565b60609190911b6001600160601b0319166101208201526101340195945050505050565b838152615848602082018461571e565b606081019190915260800192915050565b6001600160a01b0391909116815260200190565b6001600160a01b03929092168252602082015260400190565b6001600160a01b0392831681529116602082015260400190565b6001600160a01b039290921682526001600160601b0316602082015260400190565b60408101610f6d828461571e565b6020815260006137b06020830184615741565b918252602082015260400190565b93845260ff9290921660208401526040830152606082015260800190565b6020815260006137b06020830184615771565b6020815260ff82511660208201526020820151604082015260018060a01b0360408301511660608201526000606083015160c0608084015261596760e08401826156da565b60808501516001600160601b0390811660a0868101919091529095015190941660c0909301929092525090919050565b61ffff93841681529183166020830152909116604082015260600190565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b81811015615a06578451835293830193918301916001016159ea565b509098975050505050505050565b9182526001600160a01b0316602082015260400190565b828152606081016137b0602083018461571e565b8281526040602082015260006137fe6040830184615741565b86815285602082015261ffff85166040820152600063ffffffff808616606084015280851660808401525060c060a0830152613f6960c0830184615771565b878152602081018790526040810186905263ffffffff8581166060830152841660808201526001600160a01b03831660a082015260e060c0820181905260009061408d90830184615771565b8781526001600160401b03871660208201526040810186905263ffffffff8581166060830152841660808201526001600160a01b03831660a082015260e060c0820181905260009061408d90830184615771565b63ffffffff92831681529116602082015260400190565b6001600160401b0391909116815260200190565b6001600160601b038681168252851660208201526001600160401b03841660408201526001600160a01b038316606082015260a060808201819052600090615bac908301846156da565b979650505050505050565b6000808335601e19843603018112615bce57600080fd5b8301803591506001600160401b03821115615be857600080fd5b60200191503681900382131561505b57600080fd5b604080519081016001600160401b0381118282101715615c1f57615c1f615e9c565b60405290565b60405160c081016001600160401b0381118282101715615c1f57615c1f615e9c565b60405161012081016001600160401b0381118282101715615c1f57615c1f615e9c565b604051601f8201601f191681016001600160401b0381118282101715615c9257615c92615e9c565b604052919050565b60008085851115615caa57600080fd5b83861115615cb757600080fd5b5050820193919092039150565b60008219821115615cd757615cd7615e44565b500190565b60006001600160401b038281168482168083038211156157d5576157d5615e44565b60006001600160601b038281168482168083038211156157d5576157d5615e44565b600082615d2f57615d2f615e5a565b500490565b6000816000190483118215151615615d4e57615d4e615e44565b500290565b600082821015615d6557615d65615e44565b500390565b60006001600160601b0383811690831681811015615d8a57615d8a615e44565b039392505050565b6001600160e01b03198135818116916004851015615dba5780818660040360031b1b83161692505b505092915050565b60005b83811015615ddd578181015183820152602001615dc5565b838111156109405750506000910152565b6000600019821415615e0257615e02615e44565b5060010190565b60006001600160401b0382811680821415615e2657615e26615e44565b6001019392505050565b600082615e3f57615e3f615e5a565b500690565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052603160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b038116811461237557600080fd5b801515811461237557600080fdfe307866666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666a164736f6c6343000806000a",
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFCoordinatorV2PlusUpgradedVersionRequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "fulfillRandomWords", proof, rc)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorV2PlusUpgradedVersionRequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.FulfillRandomWords(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, proof, rc)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorV2PlusUpgradedVersionRequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.FulfillRandomWords(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, proof, rc)
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

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2PlusUpgradedVersionFeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2PlusUpgradedVersionFeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SetConfig(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2PlusUpgradedVersion *VRFCoordinatorV2PlusUpgradedVersionTransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2PlusUpgradedVersionFeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2PlusUpgradedVersion.Contract.SetConfig(&_VRFCoordinatorV2PlusUpgradedVersion.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
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
	FeeConfig                   VRFCoordinatorV2PlusUpgradedVersionFeeConfig
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
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	ReentrancyLock              bool
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
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
	return common.HexToHash("0x777357bb93f63d088f18112d3dba38457aec633eb8f1341e1d418380ad328e78")
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

	FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFCoordinatorV2PlusUpgradedVersionRequestCommitment) (*types.Transaction, error)

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

	SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2PlusUpgradedVersionFeeConfig) (*types.Transaction, error)

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
