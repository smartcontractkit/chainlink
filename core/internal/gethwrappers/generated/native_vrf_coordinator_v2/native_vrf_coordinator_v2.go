// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package native_vrf_coordinator_v2

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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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

type NativeVRFCoordinatorV2FeeConfig struct {
	FulfillmentFlatFeeGWeiTier1 uint32
	FulfillmentFlatFeeGWeiTier2 uint32
	FulfillmentFlatFeeGWeiTier3 uint32
	FulfillmentFlatFeeGWeiTier4 uint32
	FulfillmentFlatFeeGWeiTier5 uint32
	ReqsForTier2                *big.Int
	ReqsForTier3                *big.Int
	ReqsForTier4                *big.Int
	ReqsForTier5                *big.Int
}

type NativeVRFCoordinatorV2RequestCommitment struct {
	BlockNum         uint64
	SubId            uint64
	CallbackGasLimit uint32
	NumWords         uint32
	Sender           common.Address
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

var NativeVRFCoordinatorV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"want\",\"type\":\"uint256\"}],\"name\":\"InsufficientGasForConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structNativeVRFCoordinatorV2.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structNativeVRFCoordinatorV2.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"fundSubscription\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"getCommitment\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentSubId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"}],\"name\":\"getFeeTier\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeGWeiTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"internalType\":\"structNativeVRFCoordinatorV2.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620058c7380380620058c7833981016040819052620000349162000183565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000d7565b50505060601b6001600160601b031916608052620001b5565b6001600160a01b038116331415620001325760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019657600080fd5b81516001600160a01b0381168114620001ae57600080fd5b9392505050565b60805160601c6156ec620001db6000396000818161051c015261371e01526156ec6000f3fe6080604052600436106101e25760003560e01c80636f64f03f11610102578063af198b9711610095578063d7ae1d3011610064578063d7ae1d301461076f578063dfaba7941461078f578063e82ad7d4146107a2578063f2fde38b146107d257600080fd5b8063af198b9714610695578063c3f909d4146106d2578063caf70c4a1461072f578063d2f9f9a71461074f57600080fd5b80638da5cb5b116100d15780638da5cb5b146106055780639f87fad714610630578063a21a23e414610650578063a47c76961461066557600080fd5b80636f64f03f146105905780637341c10c146105b057806379ba5097146105d057806382359740146105e557600080fd5b80631d16919b1161017a57806364d51a2a1161014957806364d51a2a146104d557806366316d8d146104ea578063689c45171461050a57806369bcdb7d1461056357600080fd5b80631d16919b1461035757806340d6bb82146103775780635d3b1d30146103a25780635fbbc0d2146103c257600080fd5b806308821d58116101b657806308821d581461028a57806312b58349146102aa57806315c48b84146102e3578063181f5a771461030b57600080fd5b8062012291146101e757806302bcc5b61461021457806304c357cb1461023657806306bfa63714610256575b600080fd5b3480156101f357600080fd5b506101fc6107f2565b60405161020b93929190615206565b60405180910390f35b34801561022057600080fd5b5061023461022f366004615076565b61086e565b005b34801561024257600080fd5b50610234610251366004615091565b61091a565b34801561026257600080fd5b5060055467ffffffffffffffff165b60405167ffffffffffffffff909116815260200161020b565b34801561029657600080fd5b506102346102a5366004614dc7565b610b0f565b3480156102b657600080fd5b506005546801000000000000000090046bffffffffffffffffffffffff165b60405190815260200161020b565b3480156102ef57600080fd5b506102f860c881565b60405161ffff909116815260200161020b565b34801561031757600080fd5b50604080518082018252601781527f565246436f6f7264696e61746f7256323120312e302e300000000000000000006020820152905161020b9190615193565b34801561036357600080fd5b50610234610372366004614f3e565b610cee565b34801561038357600080fd5b5061038d6101f481565b60405163ffffffff909116815260200161020b565b3480156103ae57600080fd5b506102d56103bd366004614e18565b61105a565b3480156103ce57600080fd5b50600b546040805163ffffffff80841682526401000000008404811660208301526801000000000000000084048116928201929092526c010000000000000000000000008304821660608201527001000000000000000000000000000000008304909116608082015262ffffff740100000000000000000000000000000000000000008304811660a0830152770100000000000000000000000000000000000000000000008304811660c08301527a0100000000000000000000000000000000000000000000000000008304811660e08301527d0100000000000000000000000000000000000000000000000000000000009092049091166101008201526101200161020b565b3480156104e157600080fd5b506102f8606481565b3480156104f657600080fd5b50610234610505366004614d47565b611468565b34801561051657600080fd5b5061053e7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161020b565b34801561056f57600080fd5b506102d561057e36600461505d565b60009081526009602052604090205490565b34801561059c57600080fd5b506102346105ab366004614d91565b611614565b3480156105bc57600080fd5b506102346105cb366004615091565b61175e565b3480156105dc57600080fd5b506102346119ec565b3480156105f157600080fd5b50610234610600366004615076565b611ae9565b34801561061157600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff1661053e565b34801561063c57600080fd5b5061023461064b366004615091565b611ce4565b34801561065c57600080fd5b50610271612164565b34801561067157600080fd5b50610685610680366004615076565b612354565b60405161020b9493929190615392565b3480156106a157600080fd5b506106b56106b0366004614e76565b61249e565b6040516bffffffffffffffffffffffff909116815260200161020b565b3480156106de57600080fd5b50610708600a5461ffff81169163ffffffff62010000830481169267010000000000000090041690565b6040805161ffff909416845263ffffffff928316602085015291169082015260600161020b565b34801561073b57600080fd5b506102d561074a366004614de3565b612963565b34801561075b57600080fd5b5061038d61076a366004615076565b612993565b34801561077b57600080fd5b5061023461078a366004615091565b612b88565b61023461079d366004615076565b612ce9565b3480156107ae57600080fd5b506107c26107bd366004615076565b612e59565b604051901515815260200161020b565b3480156107de57600080fd5b506102346107ed366004614d2a565b6130b0565b600a546007805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff1693919283919083018282801561085c57602002820191906000526020600020905b815481526020019060010190808311610848575b50505050509050925092509250909192565b6108766130c1565b67ffffffffffffffff811660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff166108dc576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526003602052604090205461091790829073ffffffffffffffffffffffffffffffffffffffff16613144565b50565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680610983576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff8216146109ef576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024015b60405180910390fd5b600a546601000000000000900460ff1615610a36576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff841660009081526003602052604090206001015473ffffffffffffffffffffffffffffffffffffffff848116911614610b095767ffffffffffffffff841660008181526003602090815260409182902060010180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88169081179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25b50505050565b610b176130c1565b604080518082018252600091610b46919084906002908390839080828437600092019190915250612963915050565b60008181526006602052604090205490915073ffffffffffffffffffffffffffffffffffffffff1680610ba8576040517f77f5b84c000000000000000000000000000000000000000000000000000000008152600481018390526024016109e6565b600082815260066020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b600754811015610c98578260078281548110610bfb57610bfb61565f565b90600052602060002001541415610c86576007805460009190610c20906001906154f2565b81548110610c3057610c3061565f565b906000526020600020015490508060078381548110610c5157610c5161565f565b6000918252602090912001556007805480610c6e57610c6e615630565b60019003818190600052602060002001600090559055505b80610c9081615536565b915050610bdd565b508073ffffffffffffffffffffffffffffffffffffffff167f72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d83604051610ce191815260200190565b60405180910390a2505050565b610cf66130c1565b60c861ffff85161115610d49576040517fa738697600000000000000000000000000000000000000000000000000000000815261ffff851660048201819052602482015260c860448201526064016109e6565b604080516080808201835261ffff871680835263ffffffff87811660208086018290526000868801528883166060968701819052600a80547fffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000001690951762010000909302929092177fffffffffffffffffffffffffffffffffffffffffff0000000000ffffffffffff16670100000000000000909202919091179092558551600b80549388015188880151968901519589015160a08a015160c08b015160e08c01516101008d01519688167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009099169890981764010000000094881694909402939093177fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff1668010000000000000000998716999099027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff16989098176c0100000000000000000000000097861697909702969096177fffffffffffffffffff00000000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000096909416959095027fffffffffffffffffff000000ffffffffffffffffffffffffffffffffffffffff16929092177401000000000000000000000000000000000000000062ffffff96871602177fffffff000000000000ffffffffffffffffffffffffffffffffffffffffffffff1677010000000000000000000000000000000000000000000000948616949094027fffffff000000ffffffffffffffffffffffffffffffffffffffffffffffffffff16939093177a01000000000000000000000000000000000000000000000000000092851692909202919091177cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167d010000000000000000000000000000000000000000000000000000000000939092169290920217815590517f3248fab4375f32e0d851d39a71c0750d4652d98bcc7d32cec9d178c9824d796b9161104c9187918791879190615265565b60405180910390a150505050565b600a546000906601000000000000900460ff16156110a4576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff851660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1661110a576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b33600090815260026020908152604080832067ffffffffffffffff808a168552925290912054168061117a576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff871660048201523360248201526044016109e6565b600a5461ffff9081169086161080611196575060c861ffff8616115b156111e657600a546040517fa738697600000000000000000000000000000000000000000000000000000000815261ffff8088166004830152909116602482015260c860448201526064016109e6565b600a5463ffffffff620100009091048116908516111561124d57600a546040517ff5d7e01e00000000000000000000000000000000000000000000000000000000815263ffffffff80871660048301526201000090920490911660248201526044016109e6565b6101f463ffffffff8416111561129f576040517f47386bec00000000000000000000000000000000000000000000000000000000815263ffffffff841660048201526101f460248201526044016109e6565b60006112ac826001615462565b6040805160208082018c9052338284015267ffffffffffffffff808c16606084015284166080808401919091528351808403909101815260a08301845280519082012060c083018d905260e080840182905284518085039091018152610100909301909352815191012091925060009182916040805160208101849052439181019190915267ffffffffffffffff8c16606082015263ffffffff808b166080830152891660a08201523360c0820152919350915060e001604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152828252805160209182012060008681526009835283902055848352820183905261ffff8a169082015263ffffffff808916606083015287166080820152339067ffffffffffffffff8b16908c907f63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a97729060a00160405180910390a45033600090815260026020908152604080832067ffffffffffffffff808d16855292529091208054919093167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009091161790915591505095945050505050565b600a546601000000000000900460ff16156114af576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600860205260409020546bffffffffffffffffffffffff80831691161015611509576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b33600090815260086020526040812080548392906115369084906bffffffffffffffffffffffff16615509565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080600560088282829054906101000a90046bffffffffffffffffffffffff1661158d9190615509565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff166108fc826bffffffffffffffffffffffff169081150290604051600060405180830381858888f1935050505015801561160f573d6000803e3d6000fd5b505050565b61161c6130c1565b60408051808201825260009161164b919084906002908390839080828437600092019190915250612963915050565b60008181526006602052604090205490915073ffffffffffffffffffffffffffffffffffffffff16156116ad576040517f4a0b8fa7000000000000000000000000000000000000000000000000000000008152600481018290526024016109e6565b600081815260066020908152604080832080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88169081179091556007805460018101825594527fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c688909301849055518381527fe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b89101610ce1565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff16806117c7576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff82161461182e576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016109e6565b600a546601000000000000900460ff1615611875576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8416600090815260036020526040902060020154606414156118cc576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8316600090815260026020908152604080832067ffffffffffffffff8089168552925290912054161561191357610b09565b73ffffffffffffffffffffffffffffffffffffffff8316600081815260026020818152604080842067ffffffffffffffff8a1680865290835281852080547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000166001908117909155600384528286209094018054948501815585529382902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001685179055905192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e09101610b00565b60015473ffffffffffffffffffffffffffffffffffffffff163314611a6d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016109e6565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600a546601000000000000900460ff1615611b30576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff16611b96576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526003602052604090206001015473ffffffffffffffffffffffffffffffffffffffff163314611c385767ffffffffffffffff8116600090815260036020526040908190206001015490517fd084e97500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff90911660048201526024016109e6565b67ffffffffffffffff81166000818152600360209081526040918290208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560019093018054909316909255835173ffffffffffffffffffffffffffffffffffffffff909116808252928101919091529092917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f091015b60405180910390a25050565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680611d4d576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614611db4576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016109e6565b600a546601000000000000900460ff1615611dfb576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8316600090815260026020908152604080832067ffffffffffffffff808916855292529091205416611e96576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8516600482015273ffffffffffffffffffffffffffffffffffffffff841660248201526044016109e6565b67ffffffffffffffff8416600090815260036020908152604080832060020180548251818502810185019093528083529192909190830182828015611f1157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611ee6575b50505050509050600060018251611f2891906154f2565b905060005b82518110156120c7578573ffffffffffffffffffffffffffffffffffffffff16838281518110611f5f57611f5f61565f565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1614156120b5576000838381518110611f9757611f9761565f565b6020026020010151905080600360008a67ffffffffffffffff1667ffffffffffffffff1681526020019081526020016000206002018381548110611fdd57611fdd61565f565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff949094169390931790925567ffffffffffffffff8a16815260039091526040902060020180548061205757612057615630565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055019055506120c7565b806120bf81615536565b915050611f2d565b5073ffffffffffffffffffffffffffffffffffffffff8516600081815260026020908152604080832067ffffffffffffffff8b168085529083529281902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b910160405180910390a2505050505050565b600a546000906601000000000000900460ff16156121ae576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005805467ffffffffffffffff169060006121c88361556f565b82546101009290920a67ffffffffffffffff81810219909316918316021790915560055416905060008060405190808252806020026020018201604052801561221b578160200160208202803683370190505b506040805180820182526000808252602080830182815267ffffffffffffffff888116808552600484528685209551865493516bffffffffffffffffffffffff9091167fffffffffffffffffffffffff0000000000000000000000000000000000000000948516176c010000000000000000000000009190931602919091179094558451606081018652338152808301848152818701888152958552600384529590932083518154831673ffffffffffffffffffffffffffffffffffffffff9182161782559551600182018054909316961695909517905591518051949550909361230c9260028501920190614a9a565b505060405133815267ffffffffffffffff841691507f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a250905090565b67ffffffffffffffff81166000908152600360205260408120548190819060609073ffffffffffffffffffffffffffffffffffffffff166123c1576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff80861660009081526004602090815260408083205460038352928190208054600290910180548351818602810186019094528084526bffffffffffffffffffffffff8616966c010000000000000000000000009096049095169473ffffffffffffffffffffffffffffffffffffffff90921693909291839183018282801561248857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161245d575b5050505050905093509350935093509193509193565b600a546000906601000000000000900460ff16156124e8576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005a905060008060006124fc87876134fa565b9250925092506000866060015163ffffffff1667ffffffffffffffff8111156125275761252761568e565b604051908082528060200260200182016040528015612550578160200160208202803683370190505b50905060005b876060015163ffffffff168110156125c45760408051602081018590529081018290526060016040516020818303038152906040528051906020012060001c8282815181106125a7576125a761565f565b6020908102919091010152806125bc81615536565b915050612556565b506000838152600960205260408082208290555181907f1fe543e3000000000000000000000000000000000000000000000000000000009061260c9087908690602401615344565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090941693909317909252600a80547fffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffff166601000000000000179055908a015160808b01519192506000916126da9163ffffffff169084613849565b600a80547fffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffff1690556020808c01805167ffffffffffffffff9081166000908152600490935260408084205492518216845290922080549394506c01000000000000000000000000918290048316936001939192600c9261275e928692900416615462565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060006127b58a600a60000160079054906101000a900463ffffffff1663ffffffff166127af85612993565b3a613897565b6020808e015167ffffffffffffffff166000908152600490915260409020549091506bffffffffffffffffffffffff80831691161015612821576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6020808d015167ffffffffffffffff166000908152600490915260408120805483929061285d9084906bffffffffffffffffffffffff16615509565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008b81526006602090815260408083205473ffffffffffffffffffffffffffffffffffffffff16835260089091528120805485945090926128c69185911661548e565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550877f7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4888386604051612949939291909283526bffffffffffffffffffffffff9190911660208301521515604082015260600190565b60405180910390a299505050505050505050505b92915050565b6000816040516020016129769190615185565b604051602081830303815290604052805190602001209050919050565b6040805161012081018252600b5463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c010000000000000000000000008104831660608301527001000000000000000000000000000000008104909216608082015262ffffff740100000000000000000000000000000000000000008304811660a08301819052770100000000000000000000000000000000000000000000008404821660c08401527a0100000000000000000000000000000000000000000000000000008404821660e08401527d0100000000000000000000000000000000000000000000000000000000009093041661010082015260009167ffffffffffffffff841611612ab1575192915050565b8267ffffffffffffffff168160a0015162ffffff16108015612ae657508060c0015162ffffff168367ffffffffffffffff1611155b15612af5576020015192915050565b8267ffffffffffffffff168160c0015162ffffff16108015612b2a57508060e0015162ffffff168367ffffffffffffffff1611155b15612b39576040015192915050565b8267ffffffffffffffff168160e0015162ffffff16108015612b6f575080610100015162ffffff168367ffffffffffffffff1611155b15612b7e576060015192915050565b6080015192915050565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680612bf1576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614612c58576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024016109e6565b600a546601000000000000900460ff1615612c9f576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612ca884612e59565b15612cdf576040517fb42f66e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610b098484613144565b67ffffffffffffffff811660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff16612d4f576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8116600090815260046020526040812080546bffffffffffffffffffffffff1691349190612d86838561548e565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555034600560088282829054906101000a90046bffffffffffffffffffffffff16612ddd919061548e565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8823484612e44919061544a565b60408051928352602083019190915201611cd8565b67ffffffffffffffff811660009081526003602090815260408083208151606081018352815473ffffffffffffffffffffffffffffffffffffffff90811682526001830154168185015260028201805484518187028101870186528181528796939586019390929190830182828015612f0857602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311612edd575b505050505081525050905060005b8160400151518110156130a65760005b60075481101561309357600061305c60078381548110612f4857612f4861565f565b906000526020600020015485604001518581518110612f6957612f6961565f565b6020026020010151886002600089604001518981518110612f8c57612f8c61565f565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff808f168352935220541660408051602080820187905273ffffffffffffffffffffffffffffffffffffffff959095168183015267ffffffffffffffff9384166060820152919092166080808301919091528251808303909101815260a08201835280519084012060c082019490945260e080820185905282518083039091018152610100909101909152805191012091565b50600081815260096020526040902054909150156130805750600195945050505050565b508061308b81615536565b915050612f26565b508061309e81615536565b915050612f16565b5060009392505050565b6130b86130c1565b61091781613937565b60005473ffffffffffffffffffffffffffffffffffffffff163314613142576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016109e6565b565b600a546601000000000000900460ff161561318b576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff821660009081526003602090815260408083208151606081018352815473ffffffffffffffffffffffffffffffffffffffff90811682526001830154168185015260028201805484518187028101870186528181529295939486019383018282801561323657602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161320b575b5050509190925250505067ffffffffffffffff80851660009081526004602090815260408083208151808301909252546bffffffffffffffffffffffff81168083526c01000000000000000000000000909104909416918101919091529293505b83604001515181101561333d5760026000856040015183815181106132be576132be61565f565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff8a168252909252902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690558061333581615536565b915050613297565b5067ffffffffffffffff8516600090815260036020526040812080547fffffffffffffffffffffffff000000000000000000000000000000000000000090811682556001820180549091169055906133986002830182614b24565b505067ffffffffffffffff8516600090815260046020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055600580548291906008906134089084906801000000000000000090046bffffffffffffffffffffffff16615509565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508373ffffffffffffffffffffffffffffffffffffffff166108fc826bffffffffffffffffffffffff169081150290604051600060405180830381858888f1935050505015801561348a573d6000803e3d6000fd5b506040805173ffffffffffffffffffffffffffffffffffffffff861681526bffffffffffffffffffffffff8316602082015267ffffffffffffffff8716917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815910160405180910390a25050505050565b600080600061350c8560000151612963565b60008181526006602052604090205490935073ffffffffffffffffffffffffffffffffffffffff168061356e576040517f77f5b84c000000000000000000000000000000000000000000000000000000008152600481018590526024016109e6565b608086015160405161358d918691602001918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181528151602092830120600081815260099093529120549093508061360a576040517f3688124a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b85516020808801516040808a015160608b015160808c01519251613683968b96909594910195865267ffffffffffffffff948516602087015292909316604085015263ffffffff908116606085015291909116608083015273ffffffffffffffffffffffffffffffffffffffff1660a082015260c00190565b6040516020818303038152906040528051906020012081146136d1576040517fd529142c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b855167ffffffffffffffff1640806137f55786516040517fe9413d3800000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063e9413d389060240160206040518083038186803b15801561377557600080fd5b505afa158015613789573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137ad9190614dff565b9050806137f55786516040517f175dadad00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201526024016109e6565b6000886080015182604051602001613817929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c905061383c8982613a2d565b9450505050509250925092565b60005a61138881101561385b57600080fd5b61138881039050846040820482031161387357600080fd5b50823b61387f57600080fd5b60008083516020850160008789f190505b9392505050565b6000805a6138a590876154f2565b6138af908661544a565b6138b990846154b5565b905060006138d163ffffffff8616633b9aca006154b5565b90506138e9816b03eb91e20b699321ead800006154f2565b821115613922576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61392c818361544a565b979650505050505050565b73ffffffffffffffffffffffffffffffffffffffff81163314156139b7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016109e6565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000613a618360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151613ab6565b60038360200151604051602001613a79929190615330565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209392505050565b613abf89613d8d565b613b25576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e20637572766500000000000060448201526064016109e6565b613b2e88613d8d565b613b94576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f67616d6d61206973206e6f74206f6e206375727665000000000000000000000060448201526064016109e6565b613b9d83613d8d565b613c03576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e20637572766500000060448201526064016109e6565b613c0c82613d8d565b613c72576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e2063757276650000000060448201526064016109e6565b613c7e878a8887613ee8565b613ce4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e6573730000000000000060448201526064016109e6565b6000613cf08a8761408b565b90506000613d03898b878b8689896140ef565b90506000613d14838d8d8a86614283565b9050808a14613d7f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f696e76616c69642070726f6f660000000000000000000000000000000000000060448201526064016109e6565b505050505050505050505050565b80516000907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f11613e1a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420782d6f7264696e617465000000000000000000000000000060448201526064016109e6565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f11613ea7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420792d6f7264696e617465000000000000000000000000000060448201526064016109e6565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f908009613ee18360005b60200201516142e1565b1492915050565b600073ffffffffffffffffffffffffffffffffffffffff8216613f67576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f626164207769746e65737300000000000000000000000000000000000000000060448201526064016109e6565b602084015160009060011615613f7e57601c613f81565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418587600060200201510986517ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa158015614038573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff9081169088161495505050505050949350505050565b614093614b42565b6140c0600184846040516020016140ac93929190615164565b604051602081830303815290604052614339565b90505b6140cc81613d8d565b61295d5780516040805160208101929092526140e891016140ac565b90506140c3565b6140f7614b42565b825186517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f908190069106141561418a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e6374000060448201526064016109e6565b6141958789886143a2565b6141fb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4669727374206d756c20636865636b206661696c65640000000000000000000060448201526064016109e6565b6142068486856143a2565b61426c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c656400000000000000000060448201526064016109e6565b61427786848461452f565b98975050505050505050565b6000600286868685876040516020016142a1969594939291906150f2565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209695505050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f600782089392505050565b614341614b42565b61434a8261465e565b815261435f61435a826000613ed7565b6146b3565b60208201819052600290066001141561439d576020810180517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f0390525b919050565b60008261440b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f7a65726f207363616c617200000000000000000000000000000000000000000060448201526064016109e6565b8351602085015160009061442190600290615597565b1561442d57601c614430565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa1580156144b0573d6000803e3d6000fd5b5050506020604051035190506000866040516020016144cf91906150e0565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052805160209091012073ffffffffffffffffffffffffffffffffffffffff92831692169190911498975050505050505050565b614537614b42565b835160208086015185519186015160009384938493614558939091906146ed565b919450925090507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8582096001146145ec576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a0000000000000060448201526064016109e6565b60405180604001604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8061462557614625615601565b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8785099052979650505050505050565b805160208201205b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f811061439d57604080516020808201939093528151808203840181529082019091528051910120614666565b600061295d8260026146e67ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f600161544a565b901c614883565b60008080600180827ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a089050600061479583838585614977565b90985090506147a688828e886149cf565b90985090506147b788828c876149cf565b909850905060006147ca8d878b856149cf565b90985090506147db88828686614977565b90985090506147ec88828e896149cf565b909850905081811461486f577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183099650614873565b8196505b5050505050509450945094915050565b60008061488e614b60565b6020808252818101819052604082015260608101859052608081018490527ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f60a08201526148da614b7e565b60208160c08460057ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa92508261496d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f6269674d6f64457870206661696c75726521000000000000000000000000000060448201526064016109e6565b5195945050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487099097909650945050505050565b600080807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f86890990999098509650505050505050565b828054828255906000526020600020908101928215614b14579160200282015b82811115614b1457825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614aba565b50614b20929150614b9c565b5090565b50805460008255906000526020600020908101906109179190614b9c565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b80821115614b205760008155600101614b9d565b803561439d816156bd565b806040810183101561295d57600080fd5b600082601f830112614bde57600080fd5b6040516040810181811067ffffffffffffffff82111715614c0157614c0161568e565b8060405250808385604086011115614c1857600080fd5b60005b6002811015614c3a578135835260209283019290910190600101614c1b565b509195945050505050565b600060a08284031215614c5757600080fd5b60405160a0810181811067ffffffffffffffff82111715614c7a57614c7a61568e565b604052905080614c8983614d12565b8152614c9760208401614d12565b6020820152614ca860408401614cfe565b6040820152614cb960608401614cfe565b60608201526080830135614ccc816156bd565b6080919091015292915050565b803561ffff8116811461439d57600080fd5b803562ffffff8116811461439d57600080fd5b803563ffffffff8116811461439d57600080fd5b803567ffffffffffffffff8116811461439d57600080fd5b600060208284031215614d3c57600080fd5b8135613890816156bd565b60008060408385031215614d5a57600080fd5b8235614d65816156bd565b915060208301356bffffffffffffffffffffffff81168114614d8657600080fd5b809150509250929050565b60008060608385031215614da457600080fd5b8235614daf816156bd565b9150614dbe8460208501614bbc565b90509250929050565b600060408284031215614dd957600080fd5b6138908383614bbc565b600060408284031215614df557600080fd5b6138908383614bcd565b600060208284031215614e1157600080fd5b5051919050565b600080600080600060a08688031215614e3057600080fd5b85359450614e4060208701614d12565b9350614e4e60408701614cd9565b9250614e5c60608701614cfe565b9150614e6a60808701614cfe565b90509295509295909350565b600080828403610240811215614e8b57600080fd5b6101a080821215614e9b57600080fd5b614ea3615420565b9150614eaf8686614bcd565b8252614ebe8660408701614bcd565b60208301526080850135604083015260a0850135606083015260c08501356080830152614eed60e08601614bb1565b60a0830152610100614f0187828801614bcd565b60c0840152614f14876101408801614bcd565b60e08401526101808601358184015250819350614f3386828701614c45565b925050509250929050565b600080600080848603610180811215614f5657600080fd5b614f5f86614cd9565b9450614f6d60208701614cfe565b9350614f7b60408701614cfe565b9250610120807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa083011215614faf57600080fd5b614fb7615420565b9150614fc560608801614cfe565b8252614fd360808801614cfe565b6020830152614fe460a08801614cfe565b6040830152614ff560c08801614cfe565b606083015261500660e08801614cfe565b6080830152610100615019818901614ceb565b60a0840152615029828901614ceb565b60c084015261503b6101408901614ceb565b60e084015261504d6101608901614ceb565b9083015250939692955090935050565b60006020828403121561506f57600080fd5b5035919050565b60006020828403121561508857600080fd5b61389082614d12565b600080604083850312156150a457600080fd5b6150ad83614d12565b91506020830135614d86816156bd565b8060005b6002811015610b095781518452602093840193909101906001016150c1565b6150ea81836150bd565b604001919050565b86815261510260208201876150bd565b61510f60608201866150bd565b61511c60a08201856150bd565b61512960e08201846150bd565b60609190911b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166101208201526101340195945050505050565b83815261517460208201846150bd565b606081019190915260800192915050565b6040810161295d82846150bd565b600060208083528351808285015260005b818110156151c0578581018301518582016040015282016151a4565b818111156151d2576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b818110156152575784518352938301939183019160010161523b565b509098975050505050505050565b60006101808201905061ffff8616825263ffffffff80861660208401528085166040840152835481811660608501526152ab60808501838360201c1663ffffffff169052565b6152c260a08501838360401c1663ffffffff169052565b6152d960c08501838360601c1663ffffffff169052565b6152f060e08501838360801c1663ffffffff169052565b62ffffff60a082901c811661010086015260b882901c811661012086015260d082901c1661014085015260e81c6101609093019290925295945050505050565b8281526060810161389060208301846150bd565b6000604082018483526020604081850152818551808452606086019150828701935060005b8181101561538557845183529383019391830191600101615369565b5090979650505050505050565b6000608082016bffffffffffffffffffffffff87168352602067ffffffffffffffff87168185015273ffffffffffffffffffffffffffffffffffffffff80871660408601526080606086015282865180855260a087019150838801945060005b818110156154105785518416835294840194918401916001016153f2565b50909a9950505050505050505050565b604051610120810167ffffffffffffffff811182821017156154445761544461568e565b60405290565b6000821982111561545d5761545d6155d2565b500190565b600067ffffffffffffffff808316818516808303821115615485576154856155d2565b01949350505050565b60006bffffffffffffffffffffffff808316818516808303821115615485576154856155d2565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156154ed576154ed6155d2565b500290565b600082821015615504576155046155d2565b500390565b60006bffffffffffffffffffffffff8381169083168181101561552e5761552e6155d2565b039392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415615568576155686155d2565b5060010190565b600067ffffffffffffffff8083168181141561558d5761558d6155d2565b6001019392505050565b6000826155cd577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff8116811461091757600080fdfea164736f6c6343000806000a",
}

var NativeVRFCoordinatorV2ABI = NativeVRFCoordinatorV2MetaData.ABI

var NativeVRFCoordinatorV2Bin = NativeVRFCoordinatorV2MetaData.Bin

func DeployNativeVRFCoordinatorV2(auth *bind.TransactOpts, backend bind.ContractBackend, blockhashStore common.Address) (common.Address, *types.Transaction, *NativeVRFCoordinatorV2, error) {
	parsed, err := NativeVRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NativeVRFCoordinatorV2Bin), backend, blockhashStore)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NativeVRFCoordinatorV2{NativeVRFCoordinatorV2Caller: NativeVRFCoordinatorV2Caller{contract: contract}, NativeVRFCoordinatorV2Transactor: NativeVRFCoordinatorV2Transactor{contract: contract}, NativeVRFCoordinatorV2Filterer: NativeVRFCoordinatorV2Filterer{contract: contract}}, nil
}

type NativeVRFCoordinatorV2 struct {
	address common.Address
	abi     abi.ABI
	NativeVRFCoordinatorV2Caller
	NativeVRFCoordinatorV2Transactor
	NativeVRFCoordinatorV2Filterer
}

type NativeVRFCoordinatorV2Caller struct {
	contract *bind.BoundContract
}

type NativeVRFCoordinatorV2Transactor struct {
	contract *bind.BoundContract
}

type NativeVRFCoordinatorV2Filterer struct {
	contract *bind.BoundContract
}

type NativeVRFCoordinatorV2Session struct {
	Contract     *NativeVRFCoordinatorV2
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type NativeVRFCoordinatorV2CallerSession struct {
	Contract *NativeVRFCoordinatorV2Caller
	CallOpts bind.CallOpts
}

type NativeVRFCoordinatorV2TransactorSession struct {
	Contract     *NativeVRFCoordinatorV2Transactor
	TransactOpts bind.TransactOpts
}

type NativeVRFCoordinatorV2Raw struct {
	Contract *NativeVRFCoordinatorV2
}

type NativeVRFCoordinatorV2CallerRaw struct {
	Contract *NativeVRFCoordinatorV2Caller
}

type NativeVRFCoordinatorV2TransactorRaw struct {
	Contract *NativeVRFCoordinatorV2Transactor
}

func NewNativeVRFCoordinatorV2(address common.Address, backend bind.ContractBackend) (*NativeVRFCoordinatorV2, error) {
	abi, err := abi.JSON(strings.NewReader(NativeVRFCoordinatorV2ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindNativeVRFCoordinatorV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2{address: address, abi: abi, NativeVRFCoordinatorV2Caller: NativeVRFCoordinatorV2Caller{contract: contract}, NativeVRFCoordinatorV2Transactor: NativeVRFCoordinatorV2Transactor{contract: contract}, NativeVRFCoordinatorV2Filterer: NativeVRFCoordinatorV2Filterer{contract: contract}}, nil
}

func NewNativeVRFCoordinatorV2Caller(address common.Address, caller bind.ContractCaller) (*NativeVRFCoordinatorV2Caller, error) {
	contract, err := bindNativeVRFCoordinatorV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2Caller{contract: contract}, nil
}

func NewNativeVRFCoordinatorV2Transactor(address common.Address, transactor bind.ContractTransactor) (*NativeVRFCoordinatorV2Transactor, error) {
	contract, err := bindNativeVRFCoordinatorV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2Transactor{contract: contract}, nil
}

func NewNativeVRFCoordinatorV2Filterer(address common.Address, filterer bind.ContractFilterer) (*NativeVRFCoordinatorV2Filterer, error) {
	contract, err := bindNativeVRFCoordinatorV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2Filterer{contract: contract}, nil
}

func bindNativeVRFCoordinatorV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NativeVRFCoordinatorV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NativeVRFCoordinatorV2.Contract.NativeVRFCoordinatorV2Caller.contract.Call(opts, result, method, params...)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.NativeVRFCoordinatorV2Transactor.contract.Transfer(opts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.NativeVRFCoordinatorV2Transactor.contract.Transact(opts, method, params...)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NativeVRFCoordinatorV2.Contract.contract.Call(opts, result, method, params...)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.contract.Transfer(opts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.contract.Transact(opts, method, params...)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "BLOCKHASH_STORE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) BLOCKHASHSTORE() (common.Address, error) {
	return _NativeVRFCoordinatorV2.Contract.BLOCKHASHSTORE(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) BLOCKHASHSTORE() (common.Address, error) {
	return _NativeVRFCoordinatorV2.Contract.BLOCKHASHSTORE(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) MAXCONSUMERS() (uint16, error) {
	return _NativeVRFCoordinatorV2.Contract.MAXCONSUMERS(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) MAXCONSUMERS() (uint16, error) {
	return _NativeVRFCoordinatorV2.Contract.MAXCONSUMERS(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) MAXNUMWORDS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "MAX_NUM_WORDS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) MAXNUMWORDS() (uint32, error) {
	return _NativeVRFCoordinatorV2.Contract.MAXNUMWORDS(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) MAXNUMWORDS() (uint32, error) {
	return _NativeVRFCoordinatorV2.Contract.MAXNUMWORDS(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "MAX_REQUEST_CONFIRMATIONS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _NativeVRFCoordinatorV2.Contract.MAXREQUESTCONFIRMATIONS(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _NativeVRFCoordinatorV2.Contract.MAXREQUESTCONFIRMATIONS(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) GetCommitment(opts *bind.CallOpts, requestId *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "getCommitment", requestId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) GetCommitment(requestId *big.Int) ([32]byte, error) {
	return _NativeVRFCoordinatorV2.Contract.GetCommitment(&_NativeVRFCoordinatorV2.CallOpts, requestId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) GetCommitment(requestId *big.Int) ([32]byte, error) {
	return _NativeVRFCoordinatorV2.Contract.GetCommitment(&_NativeVRFCoordinatorV2.CallOpts, requestId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) GetConfig(opts *bind.CallOpts) (GetConfig,

	error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "getConfig")

	outstruct := new(GetConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MinimumRequestConfirmations = *abi.ConvertType(out[0], new(uint16)).(*uint16)
	outstruct.MaxGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.GasAfterPaymentCalculation = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) GetConfig() (GetConfig,

	error) {
	return _NativeVRFCoordinatorV2.Contract.GetConfig(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) GetConfig() (GetConfig,

	error) {
	return _NativeVRFCoordinatorV2.Contract.GetConfig(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) GetCurrentSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "getCurrentSubId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) GetCurrentSubId() (uint64, error) {
	return _NativeVRFCoordinatorV2.Contract.GetCurrentSubId(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) GetCurrentSubId() (uint64, error) {
	return _NativeVRFCoordinatorV2.Contract.GetCurrentSubId(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) GetFeeConfig(opts *bind.CallOpts) (GetFeeConfig,

	error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "getFeeConfig")

	outstruct := new(GetFeeConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.FulfillmentFlatFeeGWeiTier1 = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeGWeiTier2 = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeGWeiTier3 = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeGWeiTier4 = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeGWeiTier5 = *abi.ConvertType(out[4], new(uint32)).(*uint32)
	outstruct.ReqsForTier2 = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.ReqsForTier3 = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.ReqsForTier4 = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.ReqsForTier5 = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) GetFeeConfig() (GetFeeConfig,

	error) {
	return _NativeVRFCoordinatorV2.Contract.GetFeeConfig(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) GetFeeConfig() (GetFeeConfig,

	error) {
	return _NativeVRFCoordinatorV2.Contract.GetFeeConfig(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) GetFeeTier(opts *bind.CallOpts, reqCount uint64) (uint32, error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "getFeeTier", reqCount)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) GetFeeTier(reqCount uint64) (uint32, error) {
	return _NativeVRFCoordinatorV2.Contract.GetFeeTier(&_NativeVRFCoordinatorV2.CallOpts, reqCount)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) GetFeeTier(reqCount uint64) (uint32, error) {
	return _NativeVRFCoordinatorV2.Contract.GetFeeTier(&_NativeVRFCoordinatorV2.CallOpts, reqCount)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) GetRequestConfig(opts *bind.CallOpts) (uint16, uint32, [][32]byte, error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "getRequestConfig")

	if err != nil {
		return *new(uint16), *new(uint32), *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)
	out2 := *abi.ConvertType(out[2], new([][32]byte)).(*[][32]byte)

	return out0, out1, out2, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _NativeVRFCoordinatorV2.Contract.GetRequestConfig(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _NativeVRFCoordinatorV2.Contract.GetRequestConfig(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) GetSubscription(opts *bind.CallOpts, subId uint64) (GetSubscription,

	error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "getSubscription", subId)

	outstruct := new(GetSubscription)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Balance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ReqCount = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.Owner = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) GetSubscription(subId uint64) (GetSubscription,

	error) {
	return _NativeVRFCoordinatorV2.Contract.GetSubscription(&_NativeVRFCoordinatorV2.CallOpts, subId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) GetSubscription(subId uint64) (GetSubscription,

	error) {
	return _NativeVRFCoordinatorV2.Contract.GetSubscription(&_NativeVRFCoordinatorV2.CallOpts, subId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) GetTotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "getTotalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) GetTotalBalance() (*big.Int, error) {
	return _NativeVRFCoordinatorV2.Contract.GetTotalBalance(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) GetTotalBalance() (*big.Int, error) {
	return _NativeVRFCoordinatorV2.Contract.GetTotalBalance(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "hashOfKey", publicKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _NativeVRFCoordinatorV2.Contract.HashOfKey(&_NativeVRFCoordinatorV2.CallOpts, publicKey)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _NativeVRFCoordinatorV2.Contract.HashOfKey(&_NativeVRFCoordinatorV2.CallOpts, publicKey)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) Owner() (common.Address, error) {
	return _NativeVRFCoordinatorV2.Contract.Owner(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) Owner() (common.Address, error) {
	return _NativeVRFCoordinatorV2.Contract.Owner(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) PendingRequestExists(opts *bind.CallOpts, subId uint64) (bool, error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "pendingRequestExists", subId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) PendingRequestExists(subId uint64) (bool, error) {
	return _NativeVRFCoordinatorV2.Contract.PendingRequestExists(&_NativeVRFCoordinatorV2.CallOpts, subId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) PendingRequestExists(subId uint64) (bool, error) {
	return _NativeVRFCoordinatorV2.Contract.PendingRequestExists(&_NativeVRFCoordinatorV2.CallOpts, subId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NativeVRFCoordinatorV2.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) TypeAndVersion() (string, error) {
	return _NativeVRFCoordinatorV2.Contract.TypeAndVersion(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2CallerSession) TypeAndVersion() (string, error) {
	return _NativeVRFCoordinatorV2.Contract.TypeAndVersion(&_NativeVRFCoordinatorV2.CallOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "acceptOwnership")
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) AcceptOwnership() (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.AcceptOwnership(&_NativeVRFCoordinatorV2.TransactOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.AcceptOwnership(&_NativeVRFCoordinatorV2.TransactOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.AcceptSubscriptionOwnerTransfer(&_NativeVRFCoordinatorV2.TransactOpts, subId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.AcceptSubscriptionOwnerTransfer(&_NativeVRFCoordinatorV2.TransactOpts, subId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) AddConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "addConsumer", subId, consumer)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.AddConsumer(&_NativeVRFCoordinatorV2.TransactOpts, subId, consumer)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.AddConsumer(&_NativeVRFCoordinatorV2.TransactOpts, subId, consumer)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) CancelSubscription(opts *bind.TransactOpts, subId uint64, to common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "cancelSubscription", subId, to)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) CancelSubscription(subId uint64, to common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.CancelSubscription(&_NativeVRFCoordinatorV2.TransactOpts, subId, to)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) CancelSubscription(subId uint64, to common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.CancelSubscription(&_NativeVRFCoordinatorV2.TransactOpts, subId, to)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "createSubscription")
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) CreateSubscription() (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.CreateSubscription(&_NativeVRFCoordinatorV2.TransactOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.CreateSubscription(&_NativeVRFCoordinatorV2.TransactOpts)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "deregisterProvingKey", publicProvingKey)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.DeregisterProvingKey(&_NativeVRFCoordinatorV2.TransactOpts, publicProvingKey)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.DeregisterProvingKey(&_NativeVRFCoordinatorV2.TransactOpts, publicProvingKey)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc NativeVRFCoordinatorV2RequestCommitment) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "fulfillRandomWords", proof, rc)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) FulfillRandomWords(proof VRFProof, rc NativeVRFCoordinatorV2RequestCommitment) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.FulfillRandomWords(&_NativeVRFCoordinatorV2.TransactOpts, proof, rc)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) FulfillRandomWords(proof VRFProof, rc NativeVRFCoordinatorV2RequestCommitment) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.FulfillRandomWords(&_NativeVRFCoordinatorV2.TransactOpts, proof, rc)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) FundSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "fundSubscription", subId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) FundSubscription(subId uint64) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.FundSubscription(&_NativeVRFCoordinatorV2.TransactOpts, subId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) FundSubscription(subId uint64) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.FundSubscription(&_NativeVRFCoordinatorV2.TransactOpts, subId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "oracleWithdraw", recipient, amount)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.OracleWithdraw(&_NativeVRFCoordinatorV2.TransactOpts, recipient, amount)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.OracleWithdraw(&_NativeVRFCoordinatorV2.TransactOpts, recipient, amount)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "ownerCancelSubscription", subId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.OwnerCancelSubscription(&_NativeVRFCoordinatorV2.TransactOpts, subId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.OwnerCancelSubscription(&_NativeVRFCoordinatorV2.TransactOpts, subId)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "registerProvingKey", oracle, publicProvingKey)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.RegisterProvingKey(&_NativeVRFCoordinatorV2.TransactOpts, oracle, publicProvingKey)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.RegisterProvingKey(&_NativeVRFCoordinatorV2.TransactOpts, oracle, publicProvingKey)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) RemoveConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "removeConsumer", subId, consumer)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) RemoveConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.RemoveConsumer(&_NativeVRFCoordinatorV2.TransactOpts, subId, consumer)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) RemoveConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.RemoveConsumer(&_NativeVRFCoordinatorV2.TransactOpts, subId, consumer)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "requestRandomWords", keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) RequestRandomWords(keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.RequestRandomWords(&_NativeVRFCoordinatorV2.TransactOpts, keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) RequestRandomWords(keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.RequestRandomWords(&_NativeVRFCoordinatorV2.TransactOpts, keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.RequestSubscriptionOwnerTransfer(&_NativeVRFCoordinatorV2.TransactOpts, subId, newOwner)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.RequestSubscriptionOwnerTransfer(&_NativeVRFCoordinatorV2.TransactOpts, subId, newOwner)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, gasAfterPaymentCalculation uint32, feeConfig NativeVRFCoordinatorV2FeeConfig) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, gasAfterPaymentCalculation, feeConfig)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, gasAfterPaymentCalculation uint32, feeConfig NativeVRFCoordinatorV2FeeConfig) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.SetConfig(&_NativeVRFCoordinatorV2.TransactOpts, minimumRequestConfirmations, maxGasLimit, gasAfterPaymentCalculation, feeConfig)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, gasAfterPaymentCalculation uint32, feeConfig NativeVRFCoordinatorV2FeeConfig) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.SetConfig(&_NativeVRFCoordinatorV2.TransactOpts, minimumRequestConfirmations, maxGasLimit, gasAfterPaymentCalculation, feeConfig)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.contract.Transact(opts, "transferOwnership", to)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.TransferOwnership(&_NativeVRFCoordinatorV2.TransactOpts, to)
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _NativeVRFCoordinatorV2.Contract.TransferOwnership(&_NativeVRFCoordinatorV2.TransactOpts, to)
}

type NativeVRFCoordinatorV2ConfigSetIterator struct {
	Event *NativeVRFCoordinatorV2ConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2ConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2ConfigSet)
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
		it.Event = new(NativeVRFCoordinatorV2ConfigSet)
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

func (it *NativeVRFCoordinatorV2ConfigSetIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2ConfigSet struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	GasAfterPaymentCalculation  uint32
	FeeConfig                   NativeVRFCoordinatorV2FeeConfig
	Raw                         types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterConfigSet(opts *bind.FilterOpts) (*NativeVRFCoordinatorV2ConfigSetIterator, error) {

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2ConfigSetIterator{contract: _NativeVRFCoordinatorV2.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2ConfigSet) (event.Subscription, error) {

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2ConfigSet)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseConfigSet(log types.Log) (*NativeVRFCoordinatorV2ConfigSet, error) {
	event := new(NativeVRFCoordinatorV2ConfigSet)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2FundsRecoveredIterator struct {
	Event *NativeVRFCoordinatorV2FundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2FundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2FundsRecovered)
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
		it.Event = new(NativeVRFCoordinatorV2FundsRecovered)
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

func (it *NativeVRFCoordinatorV2FundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2FundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2FundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterFundsRecovered(opts *bind.FilterOpts) (*NativeVRFCoordinatorV2FundsRecoveredIterator, error) {

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2FundsRecoveredIterator{contract: _NativeVRFCoordinatorV2.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2FundsRecovered) (event.Subscription, error) {

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2FundsRecovered)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseFundsRecovered(log types.Log) (*NativeVRFCoordinatorV2FundsRecovered, error) {
	event := new(NativeVRFCoordinatorV2FundsRecovered)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2OwnershipTransferRequestedIterator struct {
	Event *NativeVRFCoordinatorV2OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2OwnershipTransferRequested)
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
		it.Event = new(NativeVRFCoordinatorV2OwnershipTransferRequested)
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

func (it *NativeVRFCoordinatorV2OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NativeVRFCoordinatorV2OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2OwnershipTransferRequestedIterator{contract: _NativeVRFCoordinatorV2.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2OwnershipTransferRequested)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseOwnershipTransferRequested(log types.Log) (*NativeVRFCoordinatorV2OwnershipTransferRequested, error) {
	event := new(NativeVRFCoordinatorV2OwnershipTransferRequested)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2OwnershipTransferredIterator struct {
	Event *NativeVRFCoordinatorV2OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2OwnershipTransferred)
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
		it.Event = new(NativeVRFCoordinatorV2OwnershipTransferred)
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

func (it *NativeVRFCoordinatorV2OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NativeVRFCoordinatorV2OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2OwnershipTransferredIterator{contract: _NativeVRFCoordinatorV2.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2OwnershipTransferred)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseOwnershipTransferred(log types.Log) (*NativeVRFCoordinatorV2OwnershipTransferred, error) {
	event := new(NativeVRFCoordinatorV2OwnershipTransferred)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2ProvingKeyDeregisteredIterator struct {
	Event *NativeVRFCoordinatorV2ProvingKeyDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2ProvingKeyDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2ProvingKeyDeregistered)
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
		it.Event = new(NativeVRFCoordinatorV2ProvingKeyDeregistered)
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

func (it *NativeVRFCoordinatorV2ProvingKeyDeregisteredIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2ProvingKeyDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2ProvingKeyDeregistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*NativeVRFCoordinatorV2ProvingKeyDeregisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2ProvingKeyDeregisteredIterator{contract: _NativeVRFCoordinatorV2.contract, event: "ProvingKeyDeregistered", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2ProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2ProvingKeyDeregistered)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseProvingKeyDeregistered(log types.Log) (*NativeVRFCoordinatorV2ProvingKeyDeregistered, error) {
	event := new(NativeVRFCoordinatorV2ProvingKeyDeregistered)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2ProvingKeyRegisteredIterator struct {
	Event *NativeVRFCoordinatorV2ProvingKeyRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2ProvingKeyRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2ProvingKeyRegistered)
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
		it.Event = new(NativeVRFCoordinatorV2ProvingKeyRegistered)
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

func (it *NativeVRFCoordinatorV2ProvingKeyRegisteredIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2ProvingKeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2ProvingKeyRegistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*NativeVRFCoordinatorV2ProvingKeyRegisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2ProvingKeyRegisteredIterator{contract: _NativeVRFCoordinatorV2.contract, event: "ProvingKeyRegistered", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2ProvingKeyRegistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2ProvingKeyRegistered)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseProvingKeyRegistered(log types.Log) (*NativeVRFCoordinatorV2ProvingKeyRegistered, error) {
	event := new(NativeVRFCoordinatorV2ProvingKeyRegistered)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2RandomWordsFulfilledIterator struct {
	Event *NativeVRFCoordinatorV2RandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2RandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2RandomWordsFulfilled)
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
		it.Event = new(NativeVRFCoordinatorV2RandomWordsFulfilled)
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

func (it *NativeVRFCoordinatorV2RandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2RandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2RandomWordsFulfilled struct {
	RequestId  *big.Int
	OutputSeed *big.Int
	Payment    *big.Int
	Success    bool
	Raw        types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*NativeVRFCoordinatorV2RandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2RandomWordsFulfilledIterator{contract: _NativeVRFCoordinatorV2.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2RandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2RandomWordsFulfilled)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseRandomWordsFulfilled(log types.Log) (*NativeVRFCoordinatorV2RandomWordsFulfilled, error) {
	event := new(NativeVRFCoordinatorV2RandomWordsFulfilled)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2RandomWordsRequestedIterator struct {
	Event *NativeVRFCoordinatorV2RandomWordsRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2RandomWordsRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2RandomWordsRequested)
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
		it.Event = new(NativeVRFCoordinatorV2RandomWordsRequested)
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

func (it *NativeVRFCoordinatorV2RandomWordsRequestedIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2RandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2RandomWordsRequested struct {
	KeyHash                     [32]byte
	RequestId                   *big.Int
	PreSeed                     *big.Int
	SubId                       uint64
	MinimumRequestConfirmations uint16
	CallbackGasLimit            uint32
	NumWords                    uint32
	Sender                      common.Address
	Raw                         types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*NativeVRFCoordinatorV2RandomWordsRequestedIterator, error) {

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

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2RandomWordsRequestedIterator{contract: _NativeVRFCoordinatorV2.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2RandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2RandomWordsRequested)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseRandomWordsRequested(log types.Log) (*NativeVRFCoordinatorV2RandomWordsRequested, error) {
	event := new(NativeVRFCoordinatorV2RandomWordsRequested)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2SubscriptionCanceledIterator struct {
	Event *NativeVRFCoordinatorV2SubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2SubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2SubscriptionCanceled)
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
		it.Event = new(NativeVRFCoordinatorV2SubscriptionCanceled)
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

func (it *NativeVRFCoordinatorV2SubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2SubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2SubscriptionCanceled struct {
	SubId  uint64
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2SubscriptionCanceledIterator{contract: _NativeVRFCoordinatorV2.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionCanceled, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2SubscriptionCanceled)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseSubscriptionCanceled(log types.Log) (*NativeVRFCoordinatorV2SubscriptionCanceled, error) {
	event := new(NativeVRFCoordinatorV2SubscriptionCanceled)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2SubscriptionConsumerAddedIterator struct {
	Event *NativeVRFCoordinatorV2SubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2SubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2SubscriptionConsumerAdded)
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
		it.Event = new(NativeVRFCoordinatorV2SubscriptionConsumerAdded)
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

func (it *NativeVRFCoordinatorV2SubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2SubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2SubscriptionConsumerAdded struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2SubscriptionConsumerAddedIterator{contract: _NativeVRFCoordinatorV2.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionConsumerAdded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2SubscriptionConsumerAdded)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseSubscriptionConsumerAdded(log types.Log) (*NativeVRFCoordinatorV2SubscriptionConsumerAdded, error) {
	event := new(NativeVRFCoordinatorV2SubscriptionConsumerAdded)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2SubscriptionConsumerRemovedIterator struct {
	Event *NativeVRFCoordinatorV2SubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2SubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2SubscriptionConsumerRemoved)
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
		it.Event = new(NativeVRFCoordinatorV2SubscriptionConsumerRemoved)
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

func (it *NativeVRFCoordinatorV2SubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2SubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2SubscriptionConsumerRemoved struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2SubscriptionConsumerRemovedIterator{contract: _NativeVRFCoordinatorV2.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2SubscriptionConsumerRemoved)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseSubscriptionConsumerRemoved(log types.Log) (*NativeVRFCoordinatorV2SubscriptionConsumerRemoved, error) {
	event := new(NativeVRFCoordinatorV2SubscriptionConsumerRemoved)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2SubscriptionCreatedIterator struct {
	Event *NativeVRFCoordinatorV2SubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2SubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2SubscriptionCreated)
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
		it.Event = new(NativeVRFCoordinatorV2SubscriptionCreated)
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

func (it *NativeVRFCoordinatorV2SubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2SubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2SubscriptionCreated struct {
	SubId uint64
	Owner common.Address
	Raw   types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2SubscriptionCreatedIterator{contract: _NativeVRFCoordinatorV2.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionCreated, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2SubscriptionCreated)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseSubscriptionCreated(log types.Log) (*NativeVRFCoordinatorV2SubscriptionCreated, error) {
	event := new(NativeVRFCoordinatorV2SubscriptionCreated)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2SubscriptionFundedIterator struct {
	Event *NativeVRFCoordinatorV2SubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2SubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2SubscriptionFunded)
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
		it.Event = new(NativeVRFCoordinatorV2SubscriptionFunded)
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

func (it *NativeVRFCoordinatorV2SubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2SubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2SubscriptionFunded struct {
	SubId      uint64
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2SubscriptionFundedIterator{contract: _NativeVRFCoordinatorV2.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionFunded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2SubscriptionFunded)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseSubscriptionFunded(log types.Log) (*NativeVRFCoordinatorV2SubscriptionFunded, error) {
	event := new(NativeVRFCoordinatorV2SubscriptionFunded)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator struct {
	Event *NativeVRFCoordinatorV2SubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2SubscriptionOwnerTransferRequested)
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
		it.Event = new(NativeVRFCoordinatorV2SubscriptionOwnerTransferRequested)
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

func (it *NativeVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2SubscriptionOwnerTransferRequested struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator{contract: _NativeVRFCoordinatorV2.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2SubscriptionOwnerTransferRequested)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*NativeVRFCoordinatorV2SubscriptionOwnerTransferRequested, error) {
	event := new(NativeVRFCoordinatorV2SubscriptionOwnerTransferRequested)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NativeVRFCoordinatorV2SubscriptionOwnerTransferredIterator struct {
	Event *NativeVRFCoordinatorV2SubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NativeVRFCoordinatorV2SubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NativeVRFCoordinatorV2SubscriptionOwnerTransferred)
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
		it.Event = new(NativeVRFCoordinatorV2SubscriptionOwnerTransferred)
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

func (it *NativeVRFCoordinatorV2SubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *NativeVRFCoordinatorV2SubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NativeVRFCoordinatorV2SubscriptionOwnerTransferred struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NativeVRFCoordinatorV2SubscriptionOwnerTransferredIterator{contract: _NativeVRFCoordinatorV2.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NativeVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NativeVRFCoordinatorV2SubscriptionOwnerTransferred)
				if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2Filterer) ParseSubscriptionOwnerTransferred(log types.Log) (*NativeVRFCoordinatorV2SubscriptionOwnerTransferred, error) {
	event := new(NativeVRFCoordinatorV2SubscriptionOwnerTransferred)
	if err := _NativeVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetConfig struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	GasAfterPaymentCalculation  uint32
}
type GetFeeConfig struct {
	FulfillmentFlatFeeGWeiTier1 uint32
	FulfillmentFlatFeeGWeiTier2 uint32
	FulfillmentFlatFeeGWeiTier3 uint32
	FulfillmentFlatFeeGWeiTier4 uint32
	FulfillmentFlatFeeGWeiTier5 uint32
	ReqsForTier2                *big.Int
	ReqsForTier3                *big.Int
	ReqsForTier4                *big.Int
	ReqsForTier5                *big.Int
}
type GetSubscription struct {
	Balance   *big.Int
	ReqCount  uint64
	Owner     common.Address
	Consumers []common.Address
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _NativeVRFCoordinatorV2.abi.Events["ConfigSet"].ID:
		return _NativeVRFCoordinatorV2.ParseConfigSet(log)
	case _NativeVRFCoordinatorV2.abi.Events["FundsRecovered"].ID:
		return _NativeVRFCoordinatorV2.ParseFundsRecovered(log)
	case _NativeVRFCoordinatorV2.abi.Events["OwnershipTransferRequested"].ID:
		return _NativeVRFCoordinatorV2.ParseOwnershipTransferRequested(log)
	case _NativeVRFCoordinatorV2.abi.Events["OwnershipTransferred"].ID:
		return _NativeVRFCoordinatorV2.ParseOwnershipTransferred(log)
	case _NativeVRFCoordinatorV2.abi.Events["ProvingKeyDeregistered"].ID:
		return _NativeVRFCoordinatorV2.ParseProvingKeyDeregistered(log)
	case _NativeVRFCoordinatorV2.abi.Events["ProvingKeyRegistered"].ID:
		return _NativeVRFCoordinatorV2.ParseProvingKeyRegistered(log)
	case _NativeVRFCoordinatorV2.abi.Events["RandomWordsFulfilled"].ID:
		return _NativeVRFCoordinatorV2.ParseRandomWordsFulfilled(log)
	case _NativeVRFCoordinatorV2.abi.Events["RandomWordsRequested"].ID:
		return _NativeVRFCoordinatorV2.ParseRandomWordsRequested(log)
	case _NativeVRFCoordinatorV2.abi.Events["SubscriptionCanceled"].ID:
		return _NativeVRFCoordinatorV2.ParseSubscriptionCanceled(log)
	case _NativeVRFCoordinatorV2.abi.Events["SubscriptionConsumerAdded"].ID:
		return _NativeVRFCoordinatorV2.ParseSubscriptionConsumerAdded(log)
	case _NativeVRFCoordinatorV2.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _NativeVRFCoordinatorV2.ParseSubscriptionConsumerRemoved(log)
	case _NativeVRFCoordinatorV2.abi.Events["SubscriptionCreated"].ID:
		return _NativeVRFCoordinatorV2.ParseSubscriptionCreated(log)
	case _NativeVRFCoordinatorV2.abi.Events["SubscriptionFunded"].ID:
		return _NativeVRFCoordinatorV2.ParseSubscriptionFunded(log)
	case _NativeVRFCoordinatorV2.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _NativeVRFCoordinatorV2.ParseSubscriptionOwnerTransferRequested(log)
	case _NativeVRFCoordinatorV2.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _NativeVRFCoordinatorV2.ParseSubscriptionOwnerTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (NativeVRFCoordinatorV2ConfigSet) Topic() common.Hash {
	return common.HexToHash("0x3248fab4375f32e0d851d39a71c0750d4652d98bcc7d32cec9d178c9824d796b")
}

func (NativeVRFCoordinatorV2FundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (NativeVRFCoordinatorV2OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (NativeVRFCoordinatorV2OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (NativeVRFCoordinatorV2ProvingKeyDeregistered) Topic() common.Hash {
	return common.HexToHash("0x72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d")
}

func (NativeVRFCoordinatorV2ProvingKeyRegistered) Topic() common.Hash {
	return common.HexToHash("0xe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b8")
}

func (NativeVRFCoordinatorV2RandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4")
}

func (NativeVRFCoordinatorV2RandomWordsRequested) Topic() common.Hash {
	return common.HexToHash("0x63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772")
}

func (NativeVRFCoordinatorV2SubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0xe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815")
}

func (NativeVRFCoordinatorV2SubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0")
}

func (NativeVRFCoordinatorV2SubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b")
}

func (NativeVRFCoordinatorV2SubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf")
}

func (NativeVRFCoordinatorV2SubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8")
}

func (NativeVRFCoordinatorV2SubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be")
}

func (NativeVRFCoordinatorV2SubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0")
}

func (_NativeVRFCoordinatorV2 *NativeVRFCoordinatorV2) Address() common.Address {
	return _NativeVRFCoordinatorV2.address
}

type NativeVRFCoordinatorV2Interface interface {
	BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error)

	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	MAXNUMWORDS(opts *bind.CallOpts) (uint32, error)

	MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error)

	GetCommitment(opts *bind.CallOpts, requestId *big.Int) ([32]byte, error)

	GetConfig(opts *bind.CallOpts) (GetConfig,

		error)

	GetCurrentSubId(opts *bind.CallOpts) (uint64, error)

	GetFeeConfig(opts *bind.CallOpts) (GetFeeConfig,

		error)

	GetFeeTier(opts *bind.CallOpts, reqCount uint64) (uint32, error)

	GetRequestConfig(opts *bind.CallOpts) (uint16, uint32, [][32]byte, error)

	GetSubscription(opts *bind.CallOpts, subId uint64) (GetSubscription,

		error)

	GetTotalBalance(opts *bind.CallOpts) (*big.Int, error)

	HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	PendingRequestExists(opts *bind.CallOpts, subId uint64) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error)

	CancelSubscription(opts *bind.TransactOpts, subId uint64, to common.Address) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error)

	FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc NativeVRFCoordinatorV2RequestCommitment) (*types.Transaction, error)

	FundSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64, newOwner common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, gasAfterPaymentCalculation uint32, feeConfig NativeVRFCoordinatorV2FeeConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*NativeVRFCoordinatorV2ConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2ConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*NativeVRFCoordinatorV2ConfigSet, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*NativeVRFCoordinatorV2FundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2FundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*NativeVRFCoordinatorV2FundsRecovered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NativeVRFCoordinatorV2OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*NativeVRFCoordinatorV2OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NativeVRFCoordinatorV2OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*NativeVRFCoordinatorV2OwnershipTransferred, error)

	FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*NativeVRFCoordinatorV2ProvingKeyDeregisteredIterator, error)

	WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2ProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyDeregistered(log types.Log) (*NativeVRFCoordinatorV2ProvingKeyDeregistered, error)

	FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*NativeVRFCoordinatorV2ProvingKeyRegisteredIterator, error)

	WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2ProvingKeyRegistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyRegistered(log types.Log) (*NativeVRFCoordinatorV2ProvingKeyRegistered, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*NativeVRFCoordinatorV2RandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2RandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*NativeVRFCoordinatorV2RandomWordsFulfilled, error)

	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*NativeVRFCoordinatorV2RandomWordsRequestedIterator, error)

	WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2RandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsRequested(log types.Log) (*NativeVRFCoordinatorV2RandomWordsRequested, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionCanceled, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*NativeVRFCoordinatorV2SubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionConsumerAdded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*NativeVRFCoordinatorV2SubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*NativeVRFCoordinatorV2SubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionCreated, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*NativeVRFCoordinatorV2SubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionFunded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*NativeVRFCoordinatorV2SubscriptionFunded, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*NativeVRFCoordinatorV2SubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*NativeVRFCoordinatorV2SubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *NativeVRFCoordinatorV2SubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*NativeVRFCoordinatorV2SubscriptionOwnerTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
