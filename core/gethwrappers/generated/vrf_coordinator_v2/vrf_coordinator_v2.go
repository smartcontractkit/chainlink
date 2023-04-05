// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_coordinator_v2

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
)

type VRFCoordinatorV2FeeConfig struct {
	FulfillmentFlatFeeLinkPPMTier1 uint32
	FulfillmentFlatFeeLinkPPMTier2 uint32
	FulfillmentFlatFeeLinkPPMTier3 uint32
	FulfillmentFlatFeeLinkPPMTier4 uint32
	FulfillmentFlatFeeLinkPPMTier5 uint32
	ReqsForTier2                   *big.Int
	ReqsForTier3                   *big.Int
	ReqsForTier4                   *big.Int
	ReqsForTier5                   *big.Int
}

type VRFCoordinatorV2RequestCommitment struct {
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

var VRFCoordinatorV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"want\",\"type\":\"uint256\"}],\"name\":\"InsufficientGasForConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structVRFCoordinatorV2.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structVRFCoordinatorV2.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"getCommitment\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentSubId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"}],\"name\":\"getFeeTier\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"internalType\":\"structVRFCoordinatorV2.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162006237380380620062378339810160408190526200003491620001b1565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000e8565b5050506001600160601b0319606093841b811660805290831b811660a052911b1660c052620001fb565b6001600160a01b038116331415620001435760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001ac57600080fd5b919050565b600080600060608486031215620001c757600080fd5b620001d28462000194565b9250620001e26020850162000194565b9150620001f26040850162000194565b90509250925092565b60805160601c60a05160601c60c05160601c615fd262000265600039600081816105260152613c7e01526000818161061d01526141cf01526000818161036d0152818161159401528181612591015281816130270152818161317d01526138310152615fd26000f3fe608060405234801561001057600080fd5b506004361061025b5760003560e01c80636f64f03f11610145578063ad178361116100bd578063d2f9f9a71161008c578063e72f6e3011610071578063e72f6e30146106fa578063e82ad7d41461070d578063f2fde38b1461073057600080fd5b8063d2f9f9a7146106d4578063d7ae1d30146106e757600080fd5b8063ad17836114610618578063af198b971461063f578063c3f909d41461066f578063caf70c4a146106c157600080fd5b80638da5cb5b11610114578063a21a23e4116100f9578063a21a23e4146105da578063a47c7696146105e2578063a4c0ed361461060557600080fd5b80638da5cb5b146105a95780639f87fad7146105c757600080fd5b80636f64f03f146105685780637341c10c1461057b57806379ba50971461058e578063823597401461059657600080fd5b8063356dac71116101d85780635fbbc0d2116101a757806366316d8d1161018c57806366316d8d1461050e578063689c45171461052157806369bcdb7d1461054857600080fd5b80635fbbc0d21461040057806364d51a2a1461050657600080fd5b8063356dac71146103b457806340d6bb82146103bc5780634cb48a54146103da5780635d3b1d30146103ed57600080fd5b806308821d581161022f57806315c48b841161021457806315c48b841461030e578063181f5a77146103295780631b6b6d231461036857600080fd5b806308821d58146102cf57806312b58349146102e257600080fd5b80620122911461026057806302bcc5b61461028057806304c357cb1461029557806306bfa637146102a8575b600080fd5b610268610743565b60405161027793929190615b0f565b60405180910390f35b61029361028e36600461593d565b6107bf565b005b6102936102a3366004615958565b61086b565b60055467ffffffffffffffff165b60405167ffffffffffffffff9091168152602001610277565b6102936102dd36600461564e565b610a60565b6005546801000000000000000090046bffffffffffffffffffffffff165b604051908152602001610277565b61031660c881565b60405161ffff9091168152602001610277565b604080518082018252601681527f565246436f6f7264696e61746f72563220312e302e3000000000000000000000602082015290516102779190615a9c565b61038f7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610277565b600a54610300565b6103c56101f481565b60405163ffffffff9091168152602001610277565b6102936103e83660046157e7565b610c3f565b6103006103fb3660046156c1565b611036565b600c546040805163ffffffff80841682526401000000008404811660208301526801000000000000000084048116928201929092526c010000000000000000000000008304821660608201527001000000000000000000000000000000008304909116608082015262ffffff740100000000000000000000000000000000000000008304811660a0830152770100000000000000000000000000000000000000000000008304811660c08301527a0100000000000000000000000000000000000000000000000000008304811660e08301527d01000000000000000000000000000000000000000000000000000000000090920490911661010082015261012001610277565b610316606481565b61029361051c366004615606565b61143f565b61038f7f000000000000000000000000000000000000000000000000000000000000000081565b610300610556366004615924565b60009081526009602052604090205490565b61029361057636600461554b565b6116a8565b610293610589366004615958565b6117f2565b610293611a80565b6102936105a436600461593d565b611b7d565b60005473ffffffffffffffffffffffffffffffffffffffff1661038f565b6102936105d5366004615958565b611d77565b6102b66121f8565b6105f56105f036600461593d565b6123e8565b6040516102779493929190615cad565b61029361061336600461557f565b612532565b61038f7f000000000000000000000000000000000000000000000000000000000000000081565b61065261064d36600461571f565b6127a3565b6040516bffffffffffffffffffffffff9091168152602001610277565b600b546040805161ffff8316815263ffffffff6201000084048116602083015267010000000000000084048116928201929092526b010000000000000000000000909204166060820152608001610277565b6103006106cf36600461566a565b612c68565b6103c56106e236600461593d565b612c98565b6102936106f5366004615958565b612e8d565b610293610708366004615530565b612fee565b61072061071b36600461593d565b613252565b6040519015158152602001610277565b61029361073e366004615530565b6134a9565b600b546007805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff169391928391908301828280156107ad57602002820191906000526020600020905b815481526020019060010190808311610799575b50505050509050925092509250909192565b6107c76134ba565b67ffffffffffffffff811660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1661082d576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526003602052604090205461086890829073ffffffffffffffffffffffffffffffffffffffff1661353d565b50565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff16806108d4576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614610940576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024015b60405180910390fd5b600b546601000000000000900460ff1615610987576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff841660009081526003602052604090206001015473ffffffffffffffffffffffffffffffffffffffff848116911614610a5a5767ffffffffffffffff841660008181526003602090815260409182902060010180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88169081179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25b50505050565b610a686134ba565b604080518082018252600091610a97919084906002908390839080828437600092019190915250612c68915050565b60008181526006602052604090205490915073ffffffffffffffffffffffffffffffffffffffff1680610af9576040517f77f5b84c00000000000000000000000000000000000000000000000000000000815260048101839052602401610937565b600082815260066020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b600754811015610be9578260078281548110610b4c57610b4c615f67565b90600052602060002001541415610bd7576007805460009190610b7190600190615e21565b81548110610b8157610b81615f67565b906000526020600020015490508060078381548110610ba257610ba2615f67565b6000918252602090912001556007805480610bbf57610bbf615f38565b60019003818190600052602060002001600090559055505b80610be181615e65565b915050610b2e565b508073ffffffffffffffffffffffffffffffffffffffff167f72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d83604051610c3291815260200190565b60405180910390a2505050565b610c476134ba565b60c861ffff87161115610c9a576040517fa738697600000000000000000000000000000000000000000000000000000000815261ffff871660048201819052602482015260c86044820152606401610937565b60008213610cd7576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101839052602401610937565b6040805160a0808201835261ffff891680835263ffffffff89811660208086018290526000868801528a831660608088018290528b85166080988901819052600b80547fffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000001690971762010000909502949094177fffffffffffffffffffffffffffffffffff000000000000000000ffffffffffff166701000000000000009092027fffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffff16919091176b010000000000000000000000909302929092179093558651600c80549489015189890151938a0151978a0151968a015160c08b015160e08c01516101008d01519588167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009099169890981764010000000093881693909302929092177fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff1668010000000000000000958716959095027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff16949094176c0100000000000000000000000098861698909802979097177fffffffffffffffffff00000000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000096909416959095027fffffffffffffffffff000000ffffffffffffffffffffffffffffffffffffffff16929092177401000000000000000000000000000000000000000062ffffff92831602177fffffff000000000000ffffffffffffffffffffffffffffffffffffffffffffff1677010000000000000000000000000000000000000000000000958216959095027fffffff000000ffffffffffffffffffffffffffffffffffffffffffffffffffff16949094177a01000000000000000000000000000000000000000000000000000092851692909202919091177cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167d0100000000000000000000000000000000000000000000000000000000009390911692909202919091178155600a84905590517fc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb291611026918991899189918991899190615b6e565b60405180910390a1505050505050565b600b546000906601000000000000900460ff1615611080576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff851660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff166110e6576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b33600090815260026020908152604080832067ffffffffffffffff808a1685529252909120541680611156576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff87166004820152336024820152604401610937565b600b5461ffff9081169086161080611172575060c861ffff8616115b156111c257600b546040517fa738697600000000000000000000000000000000000000000000000000000000815261ffff8088166004830152909116602482015260c86044820152606401610937565b600b5463ffffffff620100009091048116908516111561122957600b546040517ff5d7e01e00000000000000000000000000000000000000000000000000000000815263ffffffff8087166004830152620100009092049091166024820152604401610937565b6101f463ffffffff8416111561127b576040517f47386bec00000000000000000000000000000000000000000000000000000000815263ffffffff841660048201526101f46024820152604401610937565b6000611288826001615d7d565b6040805160208082018c9052338284015267ffffffffffffffff808c16606084015284166080808401919091528351808403909101815260a08301845280519082012060c083018d905260e0808401829052845180850390910181526101009093019093528151910120919250816112fe6139b0565b60408051602081019390935282015267ffffffffffffffff8a16606082015263ffffffff8089166080830152871660a08201523360c082015260e001604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152828252805160209182012060008681526009835283902055848352820183905261ffff8a169082015263ffffffff808916606083015287166080820152339067ffffffffffffffff8b16908c907f63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a97729060a00160405180910390a45033600090815260026020908152604080832067ffffffffffffffff808d16855292529091208054919093167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009091161790915591505095945050505050565b600b546601000000000000900460ff1615611486576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600860205260409020546bffffffffffffffffffffffff808316911610156114e0576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600860205260408120805483929061150d9084906bffffffffffffffffffffffff16615e38565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080600560088282829054906101000a90046bffffffffffffffffffffffff166115649190615e38565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb83836040518363ffffffff1660e01b815260040161161c92919073ffffffffffffffffffffffffffffffffffffffff9290921682526bffffffffffffffffffffffff16602082015260400190565b602060405180830381600087803b15801561163657600080fd5b505af115801561164a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061166e9190615686565b6116a4576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b6116b06134ba565b6040805180820182526000916116df919084906002908390839080828437600092019190915250612c68915050565b60008181526006602052604090205490915073ffffffffffffffffffffffffffffffffffffffff1615611741576040517f4a0b8fa700000000000000000000000000000000000000000000000000000000815260048101829052602401610937565b600081815260066020908152604080832080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88169081179091556007805460018101825594527fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c688909301849055518381527fe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b89101610c32565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff168061185b576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff8216146118c2576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610937565b600b546601000000000000900460ff1615611909576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff841660009081526003602052604090206002015460641415611960576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8316600090815260026020908152604080832067ffffffffffffffff808916855292529091205416156119a757610a5a565b73ffffffffffffffffffffffffffffffffffffffff8316600081815260026020818152604080842067ffffffffffffffff8a1680865290835281852080547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000166001908117909155600384528286209094018054948501815585529382902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001685179055905192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e09101610a51565b60015473ffffffffffffffffffffffffffffffffffffffff163314611b01576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610937565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600b546601000000000000900460ff1615611bc4576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff16611c2a576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526003602052604090206001015473ffffffffffffffffffffffffffffffffffffffff163314611ccc5767ffffffffffffffff8116600090815260036020526040908190206001015490517fd084e97500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610937565b67ffffffffffffffff81166000818152600360209081526040918290208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560019093018054909316909255835173ffffffffffffffffffffffffffffffffffffffff909116808252928101919091529092917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0910160405180910390a25050565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680611de0576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614611e47576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610937565b600b546601000000000000900460ff1615611e8e576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8316600090815260026020908152604080832067ffffffffffffffff808916855292529091205416611f29576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8516600482015273ffffffffffffffffffffffffffffffffffffffff84166024820152604401610937565b67ffffffffffffffff8416600090815260036020908152604080832060020180548251818502810185019093528083529192909190830182828015611fa457602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611f79575b50505050509050600060018251611fbb9190615e21565b905060005b825181101561215a578573ffffffffffffffffffffffffffffffffffffffff16838281518110611ff257611ff2615f67565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16141561214857600083838151811061202a5761202a615f67565b6020026020010151905080600360008a67ffffffffffffffff1667ffffffffffffffff168152602001908152602001600020600201838154811061207057612070615f67565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff949094169390931790925567ffffffffffffffff8a1681526003909152604090206002018054806120ea576120ea615f38565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190555061215a565b8061215281615e65565b915050611fc0565b5073ffffffffffffffffffffffffffffffffffffffff8516600081815260026020908152604080832067ffffffffffffffff8b168085529083529281902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b91015b60405180910390a2505050505050565b600b546000906601000000000000900460ff1615612242576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005805467ffffffffffffffff1690600061225c83615e9e565b82546101009290920a67ffffffffffffffff8181021990931691831602179091556005541690506000806040519080825280602002602001820160405280156122af578160200160208202803683370190505b506040805180820182526000808252602080830182815267ffffffffffffffff888116808552600484528685209551865493516bffffffffffffffffffffffff9091167fffffffffffffffffffffffff0000000000000000000000000000000000000000948516176c010000000000000000000000009190931602919091179094558451606081018652338152808301848152818701888152958552600384529590932083518154831673ffffffffffffffffffffffffffffffffffffffff918216178255955160018201805490931696169590951790559151805194955090936123a09260028501920190615270565b505060405133815267ffffffffffffffff841691507f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a250905090565b67ffffffffffffffff81166000908152600360205260408120548190819060609073ffffffffffffffffffffffffffffffffffffffff16612455576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff80861660009081526004602090815260408083205460038352928190208054600290910180548351818602810186019094528084526bffffffffffffffffffffffff8616966c010000000000000000000000009096049095169473ffffffffffffffffffffffffffffffffffffffff90921693909291839183018282801561251c57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116124f1575b5050505050905093509350935093509193509193565b600b546601000000000000900460ff1615612579576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146125e8576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114612622576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006126308284018461593d565b67ffffffffffffffff811660009081526003602052604090205490915073ffffffffffffffffffffffffffffffffffffffff16612699576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8116600090815260046020526040812080546bffffffffffffffffffffffff16918691906126d08385615da9565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555084600560088282829054906101000a90046bffffffffffffffffffffffff166127279190615da9565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f882878461278e9190615d65565b604080519283526020830191909152016121e8565b600b546000906601000000000000900460ff16156127ed576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005a905060008060006128018787613a56565b9250925092506000866060015163ffffffff1667ffffffffffffffff81111561282c5761282c615f96565b604051908082528060200260200182016040528015612855578160200160208202803683370190505b50905060005b876060015163ffffffff168110156128c95760408051602081018590529081018290526060016040516020818303038152906040528051906020012060001c8282815181106128ac576128ac615f67565b6020908102919091010152806128c181615e65565b91505061285b565b506000838152600960205260408082208290555181907f1fe543e300000000000000000000000000000000000000000000000000000000906129119087908690602401615c5f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090941693909317909252600b80547fffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffff166601000000000000179055908a015160808b01519192506000916129df9163ffffffff169084613da9565b600b80547fffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffff1690556020808c01805167ffffffffffffffff9081166000908152600490935260408084205492518216845290922080549394506c01000000000000000000000000918290048316936001939192600c92612a63928692900416615d7d565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506000612aba8a600b600001600b9054906101000a900463ffffffff1663ffffffff16612ab485612c98565b3a613df7565b6020808e015167ffffffffffffffff166000908152600490915260409020549091506bffffffffffffffffffffffff80831691161015612b26576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6020808d015167ffffffffffffffff1660009081526004909152604081208054839290612b629084906bffffffffffffffffffffffff16615e38565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008b81526006602090815260408083205473ffffffffffffffffffffffffffffffffffffffff1683526008909152812080548594509092612bcb91859116615da9565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550877f7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4888386604051612c4e939291909283526bffffffffffffffffffffffff9190911660208301521515604082015260600190565b60405180910390a299505050505050505050505b92915050565b600081604051602001612c7b9190615a8e565b604051602081830303815290604052805190602001209050919050565b6040805161012081018252600c5463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c010000000000000000000000008104831660608301527001000000000000000000000000000000008104909216608082015262ffffff740100000000000000000000000000000000000000008304811660a08301819052770100000000000000000000000000000000000000000000008404821660c08401527a0100000000000000000000000000000000000000000000000000008404821660e08401527d0100000000000000000000000000000000000000000000000000000000009093041661010082015260009167ffffffffffffffff841611612db6575192915050565b8267ffffffffffffffff168160a0015162ffffff16108015612deb57508060c0015162ffffff168367ffffffffffffffff1611155b15612dfa576020015192915050565b8267ffffffffffffffff168160c0015162ffffff16108015612e2f57508060e0015162ffffff168367ffffffffffffffff1611155b15612e3e576040015192915050565b8267ffffffffffffffff168160e0015162ffffff16108015612e74575080610100015162ffffff168367ffffffffffffffff1611155b15612e83576060015192915050565b6080015192915050565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680612ef6576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614612f5d576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610937565b600b546601000000000000900460ff1615612fa4576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612fad84613252565b15612fe4576040517fb42f66e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a5a848461353d565b612ff66134ba565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b15801561307e57600080fd5b505afa158015613092573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130b691906156a8565b6005549091506801000000000000000090046bffffffffffffffffffffffff168181111561311a576040517fa99da3020000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610937565b8181101561324d57600061312e8284615e21565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152602482018390529192507f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb90604401602060405180830381600087803b1580156131c357600080fd5b505af11580156131d7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906131fb9190615686565b506040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a1505b505050565b67ffffffffffffffff811660009081526003602090815260408083208151606081018352815473ffffffffffffffffffffffffffffffffffffffff9081168252600183015416818501526002820180548451818702810187018652818152879693958601939092919083018282801561330157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116132d6575b505050505081525050905060005b81604001515181101561349f5760005b60075481101561348c5760006134556007838154811061334157613341615f67565b90600052602060002001548560400151858151811061336257613362615f67565b602002602001015188600260008960400151898151811061338557613385615f67565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff808f168352935220541660408051602080820187905273ffffffffffffffffffffffffffffffffffffffff959095168183015267ffffffffffffffff9384166060820152919092166080808301919091528251808303909101815260a08201835280519084012060c082019490945260e080820185905282518083039091018152610100909101909152805191012091565b50600081815260096020526040902054909150156134795750600195945050505050565b508061348481615e65565b91505061331f565b508061349781615e65565b91505061330f565b5060009392505050565b6134b16134ba565b61086881613eff565b60005473ffffffffffffffffffffffffffffffffffffffff16331461353b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610937565b565b600b546601000000000000900460ff1615613584576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff821660009081526003602090815260408083208151606081018352815473ffffffffffffffffffffffffffffffffffffffff90811682526001830154168185015260028201805484518187028101870186528181529295939486019383018282801561362f57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311613604575b5050509190925250505067ffffffffffffffff80851660009081526004602090815260408083208151808301909252546bffffffffffffffffffffffff81168083526c01000000000000000000000000909104909416918101919091529293505b8360400151518110156137365760026000856040015183815181106136b7576136b7615f67565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff8a168252909252902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690558061372e81615e65565b915050613690565b5067ffffffffffffffff8516600090815260036020526040812080547fffffffffffffffffffffffff0000000000000000000000000000000000000000908116825560018201805490911690559061379160028301826152fa565b505067ffffffffffffffff8516600090815260046020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055600580548291906008906138019084906801000000000000000090046bffffffffffffffffffffffff16615e38565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb85836bffffffffffffffffffffffff166040518363ffffffff1660e01b81526004016138b992919073ffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b602060405180830381600087803b1580156138d357600080fd5b505af11580156138e7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061390b9190615686565b613941576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805173ffffffffffffffffffffffffffffffffffffffff861681526bffffffffffffffffffffffff8316602082015267ffffffffffffffff8716917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815910160405180910390a25050505050565b60004661a4b18114806139c5575062066eed81145b15613a4f57606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015613a1157600080fd5b505afa158015613a25573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a4991906156a8565b91505090565b4391505090565b6000806000613a688560000151612c68565b60008181526006602052604090205490935073ffffffffffffffffffffffffffffffffffffffff1680613aca576040517f77f5b84c00000000000000000000000000000000000000000000000000000000815260048101859052602401610937565b6080860151604051613ae9918691602001918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291815281516020928301206000818152600990935291205490935080613b66576040517f3688124a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b85516020808801516040808a015160608b015160808c01519251613bdf968b96909594910195865267ffffffffffffffff948516602087015292909316604085015263ffffffff908116606085015291909116608083015273ffffffffffffffffffffffffffffffffffffffff1660a082015260c00190565b604051602081830303815290604052805190602001208114613c2d576040517fd529142c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000613c3c8760000151613ff5565b905080613d555786516040517fe9413d3800000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063e9413d389060240160206040518083038186803b158015613cd557600080fd5b505afa158015613ce9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613d0d91906156a8565b905080613d555786516040517f175dadad00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610937565b6000886080015182604051602001613d77929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c9050613d9c89826140fb565b9450505050509250925092565b60005a611388811015613dbb57600080fd5b611388810390508460408204820311613dd357600080fd5b50823b613ddf57600080fd5b60008083516020850160008789f190505b9392505050565b600080613e02614184565b905060008113613e41576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101829052602401610937565b6000815a613e4f8989615d65565b613e599190615e21565b613e6b86670de0b6b3a7640000615de4565b613e759190615de4565b613e7f9190615dd0565b90506000613e9863ffffffff871664e8d4a51000615de4565b9050613eb0816b033b2e3c9fd0803ce8000000615e21565b821115613ee9576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b613ef38183615d65565b98975050505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116331415613f7f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610937565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60004661a4b181148061400a575062066eed81145b156140eb576101008367ffffffffffffffff166140256139b0565b61402f9190615e21565b118061404b575061403e6139b0565b8367ffffffffffffffff16115b156140595750600092915050565b6040517f2b407a8200000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152606490632b407a829060240160206040518083038186803b1580156140b357600080fd5b505afa1580156140c7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613df091906156a8565b505067ffffffffffffffff164090565b600061412f8360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151614298565b60038360200151604051602001614147929190615c4b565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209392505050565b600b54604080517ffeaf968c0000000000000000000000000000000000000000000000000000000081529051600092670100000000000000900463ffffffff169182151591849182917f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169163feaf968c9160048083019260a0929190829003018186803b15801561422a57600080fd5b505afa15801561423e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906142629190615982565b509450909250849150508015614286575061427d8242615e21565b8463ffffffff16105b156142905750600a545b949350505050565b6142a18961456f565b614307576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e2063757276650000000000006044820152606401610937565b6143108861456f565b614376576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f67616d6d61206973206e6f74206f6e20637572766500000000000000000000006044820152606401610937565b61437f8361456f565b6143e5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610937565b6143ee8261456f565b614454576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e206375727665000000006044820152606401610937565b614460878a88876146ca565b6144c6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e657373000000000000006044820152606401610937565b60006144d28a8761486d565b905060006144e5898b878b8689896148d1565b905060006144f6838d8d8a86614a59565b9050808a14614561576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f696e76616c69642070726f6f66000000000000000000000000000000000000006044820152606401610937565b505050505050505050505050565b80516000907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f116145fc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420782d6f7264696e61746500000000000000000000000000006044820152606401610937565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f11614689576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420792d6f7264696e61746500000000000000000000000000006044820152606401610937565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f9080096146c38360005b6020020151614ab7565b1492915050565b600073ffffffffffffffffffffffffffffffffffffffff8216614749576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f626164207769746e6573730000000000000000000000000000000000000000006044820152606401610937565b60208401516000906001161561476057601c614763565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418587600060200201510986517ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa15801561481a573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff9081169088161495505050505050949350505050565b614875615318565b6148a26001848460405160200161488e93929190615a6d565b604051602081830303815290604052614b0f565b90505b6148ae8161456f565b612c625780516040805160208101929092526148ca910161488e565b90506148a5565b6148d9615318565b825186517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f908190069106141561496c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610937565b614977878988614b78565b6149dd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4669727374206d756c20636865636b206661696c6564000000000000000000006044820152606401610937565b6149e8848685614b78565b614a4e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c65640000000000000000006044820152606401610937565b613ef3868484614d05565b600060028686868587604051602001614a77969594939291906159fb565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209695505050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f600782089392505050565b614b17615318565b614b2082614e34565b8152614b35614b308260006146b9565b614e89565b602082018190526002900660011415614b73576020810180517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f0390525b919050565b600082614be1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f7a65726f207363616c61720000000000000000000000000000000000000000006044820152606401610937565b83516020850151600090614bf790600290615ec6565b15614c0357601c614c06565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa158015614c86573d6000803e3d6000fd5b505050602060405103519050600086604051602001614ca591906159e9565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052805160209091012073ffffffffffffffffffffffffffffffffffffffff92831692169190911498975050505050505050565b614d0d615318565b835160208086015185519186015160009384938493614d2e93909190614ec3565b919450925090507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f858209600114614dc2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a000000000000006044820152606401610937565b60405180604001604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80614dfb57614dfb615f09565b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8785099052979650505050505050565b805160208201205b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8110614b7357604080516020808201939093528151808203840181529082019091528051910120614e3c565b6000612c62826002614ebc7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f6001615d65565b901c615059565b60008080600180827ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a0890506000614f6b8383858561514d565b9098509050614f7c88828e886151a5565b9098509050614f8d88828c876151a5565b90985090506000614fa08d878b856151a5565b9098509050614fb18882868661514d565b9098509050614fc288828e896151a5565b9098509050818114615045577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183099650615049565b8196505b5050505050509450945094915050565b600080615064615336565b6020808252818101819052604082015260608101859052608081018490527ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f60a08201526150b0615354565b60208160c08460057ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa925082615143576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f6269674d6f64457870206661696c7572652100000000000000000000000000006044820152606401610937565b5195945050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487099097909650945050505050565b600080807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f86890990999098509650505050505050565b8280548282559060005260206000209081019282156152ea579160200282015b828111156152ea57825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190615290565b506152f6929150615372565b5090565b50805460008255906000526020600020908101906108689190615372565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b808211156152f65760008155600101615373565b803573ffffffffffffffffffffffffffffffffffffffff81168114614b7357600080fd5b8060408101831015612c6257600080fd5b600082601f8301126153cd57600080fd5b6040516040810181811067ffffffffffffffff821117156153f0576153f0615f96565b806040525080838560408601111561540757600080fd5b60005b600281101561542957813583526020928301929091019060010161540a565b509195945050505050565b600060a0828403121561544657600080fd5b60405160a0810181811067ffffffffffffffff8211171561546957615469615f96565b604052905080615478836154fe565b8152615486602084016154fe565b6020820152615497604084016154ea565b60408201526154a8606084016154ea565b60608201526154b960808401615387565b60808201525092915050565b803561ffff81168114614b7357600080fd5b803562ffffff81168114614b7357600080fd5b803563ffffffff81168114614b7357600080fd5b803567ffffffffffffffff81168114614b7357600080fd5b805169ffffffffffffffffffff81168114614b7357600080fd5b60006020828403121561554257600080fd5b613df082615387565b6000806060838503121561555e57600080fd5b61556783615387565b915061557684602085016153ab565b90509250929050565b6000806000806060858703121561559557600080fd5b61559e85615387565b935060208501359250604085013567ffffffffffffffff808211156155c257600080fd5b818701915087601f8301126155d657600080fd5b8135818111156155e557600080fd5b8860208285010111156155f757600080fd5b95989497505060200194505050565b6000806040838503121561561957600080fd5b61562283615387565b915060208301356bffffffffffffffffffffffff8116811461564357600080fd5b809150509250929050565b60006040828403121561566057600080fd5b613df083836153ab565b60006040828403121561567c57600080fd5b613df083836153bc565b60006020828403121561569857600080fd5b81518015158114613df057600080fd5b6000602082840312156156ba57600080fd5b5051919050565b600080600080600060a086880312156156d957600080fd5b853594506156e9602087016154fe565b93506156f7604087016154c5565b9250615705606087016154ea565b9150615713608087016154ea565b90509295509295909350565b60008082840361024081121561573457600080fd5b6101a08082121561574457600080fd5b61574c615d3b565b915061575886866153bc565b825261576786604087016153bc565b60208301526080850135604083015260a0850135606083015260c0850135608083015261579660e08601615387565b60a08301526101006157aa878288016153bc565b60c08401526157bd8761014088016153bc565b60e084015261018086013581840152508193506157dc86828701615434565b925050509250929050565b6000806000806000808688036101c081121561580257600080fd5b61580b886154c5565b9650615819602089016154ea565b9550615827604089016154ea565b9450615835606089016154ea565b935060808801359250610120807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff608301121561587057600080fd5b615878615d3b565b915061588660a08a016154ea565b825261589460c08a016154ea565b60208301526158a560e08a016154ea565b60408301526101006158b8818b016154ea565b60608401526158c8828b016154ea565b60808401526158da6101408b016154d7565b60a08401526158ec6101608b016154d7565b60c08401526158fe6101808b016154d7565b60e08401526159106101a08b016154d7565b818401525050809150509295509295509295565b60006020828403121561593657600080fd5b5035919050565b60006020828403121561594f57600080fd5b613df0826154fe565b6000806040838503121561596b57600080fd5b615974836154fe565b915061557660208401615387565b600080600080600060a0868803121561599a57600080fd5b6159a386615516565b945060208601519350604086015192506060860151915061571360808701615516565b8060005b6002811015610a5a5781518452602093840193909101906001016159ca565b6159f381836159c6565b604001919050565b868152615a0b60208201876159c6565b615a1860608201866159c6565b615a2560a08201856159c6565b615a3260e08201846159c6565b60609190911b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166101208201526101340195945050505050565b838152615a7d60208201846159c6565b606081019190915260800192915050565b60408101612c6282846159c6565b600060208083528351808285015260005b81811015615ac957858101830151858201604001528201615aad565b81811115615adb576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b81811015615b6057845183529383019391830191600101615b44565b509098975050505050505050565b60006101c08201905061ffff8816825263ffffffff808816602084015280871660408401528086166060840152846080840152835481811660a0850152615bc260c08501838360201c1663ffffffff169052565b615bd960e08501838360401c1663ffffffff169052565b615bf16101008501838360601c1663ffffffff169052565b615c096101208501838360801c1663ffffffff169052565b62ffffff60a082901c811661014086015260b882901c811661016086015260d082901c1661018085015260e81c6101a090930192909252979650505050505050565b82815260608101613df060208301846159c6565b6000604082018483526020604081850152818551808452606086019150828701935060005b81811015615ca057845183529383019391830191600101615c84565b5090979650505050505050565b6000608082016bffffffffffffffffffffffff87168352602067ffffffffffffffff87168185015273ffffffffffffffffffffffffffffffffffffffff80871660408601526080606086015282865180855260a087019150838801945060005b81811015615d2b578551841683529484019491840191600101615d0d565b50909a9950505050505050505050565b604051610120810167ffffffffffffffff81118282101715615d5f57615d5f615f96565b60405290565b60008219821115615d7857615d78615eda565b500190565b600067ffffffffffffffff808316818516808303821115615da057615da0615eda565b01949350505050565b60006bffffffffffffffffffffffff808316818516808303821115615da057615da0615eda565b600082615ddf57615ddf615f09565b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615615e1c57615e1c615eda565b500290565b600082821015615e3357615e33615eda565b500390565b60006bffffffffffffffffffffffff83811690831681811015615e5d57615e5d615eda565b039392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415615e9757615e97615eda565b5060010190565b600067ffffffffffffffff80831681811415615ebc57615ebc615eda565b6001019392505050565b600082615ed557615ed5615f09565b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFCoordinatorV2ABI = VRFCoordinatorV2MetaData.ABI

var VRFCoordinatorV2Bin = VRFCoordinatorV2MetaData.Bin

func DeployVRFCoordinatorV2(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, blockhashStore common.Address, linkEthFeed common.Address) (common.Address, *types.Transaction, *VRFCoordinatorV2, error) {
	parsed, err := VRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorV2Bin), backend, link, blockhashStore, linkEthFeed)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorV2{VRFCoordinatorV2Caller: VRFCoordinatorV2Caller{contract: contract}, VRFCoordinatorV2Transactor: VRFCoordinatorV2Transactor{contract: contract}, VRFCoordinatorV2Filterer: VRFCoordinatorV2Filterer{contract: contract}}, nil
}

type VRFCoordinatorV2 struct {
	address common.Address
	abi     abi.ABI
	VRFCoordinatorV2Caller
	VRFCoordinatorV2Transactor
	VRFCoordinatorV2Filterer
}

type VRFCoordinatorV2Caller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2Transactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2Filterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2Session struct {
	Contract     *VRFCoordinatorV2
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV2CallerSession struct {
	Contract *VRFCoordinatorV2Caller
	CallOpts bind.CallOpts
}

type VRFCoordinatorV2TransactorSession struct {
	Contract     *VRFCoordinatorV2Transactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV2Raw struct {
	Contract *VRFCoordinatorV2
}

type VRFCoordinatorV2CallerRaw struct {
	Contract *VRFCoordinatorV2Caller
}

type VRFCoordinatorV2TransactorRaw struct {
	Contract *VRFCoordinatorV2Transactor
}

func NewVRFCoordinatorV2(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorV2, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorV2ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFCoordinatorV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2{address: address, abi: abi, VRFCoordinatorV2Caller: VRFCoordinatorV2Caller{contract: contract}, VRFCoordinatorV2Transactor: VRFCoordinatorV2Transactor{contract: contract}, VRFCoordinatorV2Filterer: VRFCoordinatorV2Filterer{contract: contract}}, nil
}

func NewVRFCoordinatorV2Caller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorV2Caller, error) {
	contract, err := bindVRFCoordinatorV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2Caller{contract: contract}, nil
}

func NewVRFCoordinatorV2Transactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorV2Transactor, error) {
	contract, err := bindVRFCoordinatorV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2Transactor{contract: contract}, nil
}

func NewVRFCoordinatorV2Filterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorV2Filterer, error) {
	contract, err := bindVRFCoordinatorV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2Filterer{contract: contract}, nil
}

func bindVRFCoordinatorV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFCoordinatorV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2.Contract.VRFCoordinatorV2Caller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.VRFCoordinatorV2Transactor.contract.Transfer(opts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.VRFCoordinatorV2Transactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.contract.Transfer(opts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "BLOCKHASH_STORE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) LINK() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.LINK(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) LINK() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.LINK(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) LINKETHFEED() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.LINKETHFEED(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) LINKETHFEED() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.LINKETHFEED(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV2.Contract.MAXCONSUMERS(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV2.Contract.MAXCONSUMERS(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) MAXNUMWORDS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "MAX_NUM_WORDS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV2.Contract.MAXNUMWORDS(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV2.Contract.MAXNUMWORDS(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "MAX_REQUEST_CONFIRMATIONS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV2.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV2.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetCommitment(opts *bind.CallOpts, requestId *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getCommitment", requestId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetCommitment(requestId *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2.Contract.GetCommitment(&_VRFCoordinatorV2.CallOpts, requestId)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetCommitment(requestId *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2.Contract.GetCommitment(&_VRFCoordinatorV2.CallOpts, requestId)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetConfig(opts *bind.CallOpts) (GetConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getConfig")

	outstruct := new(GetConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MinimumRequestConfirmations = *abi.ConvertType(out[0], new(uint16)).(*uint16)
	outstruct.MaxGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.StalenessSeconds = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.GasAfterPaymentCalculation = *abi.ConvertType(out[3], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetConfig() (GetConfig,

	error) {
	return _VRFCoordinatorV2.Contract.GetConfig(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetConfig() (GetConfig,

	error) {
	return _VRFCoordinatorV2.Contract.GetConfig(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetCurrentSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getCurrentSubId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetCurrentSubId() (uint64, error) {
	return _VRFCoordinatorV2.Contract.GetCurrentSubId(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetCurrentSubId() (uint64, error) {
	return _VRFCoordinatorV2.Contract.GetCurrentSubId(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getFallbackWeiPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV2.Contract.GetFallbackWeiPerUnitLink(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV2.Contract.GetFallbackWeiPerUnitLink(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetFeeConfig(opts *bind.CallOpts) (GetFeeConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getFeeConfig")

	outstruct := new(GetFeeConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.FulfillmentFlatFeeLinkPPMTier1 = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeLinkPPMTier2 = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeLinkPPMTier3 = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeLinkPPMTier4 = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	outstruct.FulfillmentFlatFeeLinkPPMTier5 = *abi.ConvertType(out[4], new(uint32)).(*uint32)
	outstruct.ReqsForTier2 = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.ReqsForTier3 = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.ReqsForTier4 = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.ReqsForTier5 = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetFeeConfig() (GetFeeConfig,

	error) {
	return _VRFCoordinatorV2.Contract.GetFeeConfig(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetFeeConfig() (GetFeeConfig,

	error) {
	return _VRFCoordinatorV2.Contract.GetFeeConfig(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetFeeTier(opts *bind.CallOpts, reqCount uint64) (uint32, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getFeeTier", reqCount)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetFeeTier(reqCount uint64) (uint32, error) {
	return _VRFCoordinatorV2.Contract.GetFeeTier(&_VRFCoordinatorV2.CallOpts, reqCount)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetFeeTier(reqCount uint64) (uint32, error) {
	return _VRFCoordinatorV2.Contract.GetFeeTier(&_VRFCoordinatorV2.CallOpts, reqCount)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetRequestConfig(opts *bind.CallOpts) (uint16, uint32, [][32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getRequestConfig")

	if err != nil {
		return *new(uint16), *new(uint32), *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)
	out2 := *abi.ConvertType(out[2], new([][32]byte)).(*[][32]byte)

	return out0, out1, out2, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _VRFCoordinatorV2.Contract.GetRequestConfig(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _VRFCoordinatorV2.Contract.GetRequestConfig(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetSubscription(opts *bind.CallOpts, subId uint64) (GetSubscription,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getSubscription", subId)

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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetSubscription(subId uint64) (GetSubscription,

	error) {
	return _VRFCoordinatorV2.Contract.GetSubscription(&_VRFCoordinatorV2.CallOpts, subId)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetSubscription(subId uint64) (GetSubscription,

	error) {
	return _VRFCoordinatorV2.Contract.GetSubscription(&_VRFCoordinatorV2.CallOpts, subId)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) GetTotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "getTotalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) GetTotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV2.Contract.GetTotalBalance(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) GetTotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV2.Contract.GetTotalBalance(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "hashOfKey", publicKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2.Contract.HashOfKey(&_VRFCoordinatorV2.CallOpts, publicKey)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV2.Contract.HashOfKey(&_VRFCoordinatorV2.CallOpts, publicKey)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) Owner() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.Owner(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) Owner() (common.Address, error) {
	return _VRFCoordinatorV2.Contract.Owner(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) PendingRequestExists(opts *bind.CallOpts, subId uint64) (bool, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "pendingRequestExists", subId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) PendingRequestExists(subId uint64) (bool, error) {
	return _VRFCoordinatorV2.Contract.PendingRequestExists(&_VRFCoordinatorV2.CallOpts, subId)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) PendingRequestExists(subId uint64) (bool, error) {
	return _VRFCoordinatorV2.Contract.PendingRequestExists(&_VRFCoordinatorV2.CallOpts, subId)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) TypeAndVersion() (string, error) {
	return _VRFCoordinatorV2.Contract.TypeAndVersion(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) TypeAndVersion() (string, error) {
	return _VRFCoordinatorV2.Contract.TypeAndVersion(&_VRFCoordinatorV2.CallOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "acceptOwnership")
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.AcceptOwnership(&_VRFCoordinatorV2.TransactOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.AcceptOwnership(&_VRFCoordinatorV2.TransactOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV2.TransactOpts, subId)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV2.TransactOpts, subId)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) AddConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "addConsumer", subId, consumer)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.AddConsumer(&_VRFCoordinatorV2.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.AddConsumer(&_VRFCoordinatorV2.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) CancelSubscription(opts *bind.TransactOpts, subId uint64, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "cancelSubscription", subId, to)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) CancelSubscription(subId uint64, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.CancelSubscription(&_VRFCoordinatorV2.TransactOpts, subId, to)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) CancelSubscription(subId uint64, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.CancelSubscription(&_VRFCoordinatorV2.TransactOpts, subId, to)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "createSubscription")
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.CreateSubscription(&_VRFCoordinatorV2.TransactOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.CreateSubscription(&_VRFCoordinatorV2.TransactOpts)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "deregisterProvingKey", publicProvingKey)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.DeregisterProvingKey(&_VRFCoordinatorV2.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.DeregisterProvingKey(&_VRFCoordinatorV2.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFCoordinatorV2RequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "fulfillRandomWords", proof, rc)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorV2RequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.FulfillRandomWords(&_VRFCoordinatorV2.TransactOpts, proof, rc)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorV2RequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.FulfillRandomWords(&_VRFCoordinatorV2.TransactOpts, proof, rc)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.OnTokenTransfer(&_VRFCoordinatorV2.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.OnTokenTransfer(&_VRFCoordinatorV2.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "oracleWithdraw", recipient, amount)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.OracleWithdraw(&_VRFCoordinatorV2.TransactOpts, recipient, amount)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.OracleWithdraw(&_VRFCoordinatorV2.TransactOpts, recipient, amount)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "ownerCancelSubscription", subId)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.OwnerCancelSubscription(&_VRFCoordinatorV2.TransactOpts, subId)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.OwnerCancelSubscription(&_VRFCoordinatorV2.TransactOpts, subId)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "recoverFunds", to)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RecoverFunds(&_VRFCoordinatorV2.TransactOpts, to)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RecoverFunds(&_VRFCoordinatorV2.TransactOpts, to)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "registerProvingKey", oracle, publicProvingKey)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RegisterProvingKey(&_VRFCoordinatorV2.TransactOpts, oracle, publicProvingKey)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RegisterProvingKey(&_VRFCoordinatorV2.TransactOpts, oracle, publicProvingKey)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) RemoveConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "removeConsumer", subId, consumer)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) RemoveConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RemoveConsumer(&_VRFCoordinatorV2.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) RemoveConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RemoveConsumer(&_VRFCoordinatorV2.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "requestRandomWords", keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) RequestRandomWords(keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RequestRandomWords(&_VRFCoordinatorV2.TransactOpts, keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) RequestRandomWords(keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RequestRandomWords(&_VRFCoordinatorV2.TransactOpts, keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV2.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV2.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2FeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2FeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.SetConfig(&_VRFCoordinatorV2.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2FeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.SetConfig(&_VRFCoordinatorV2.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.TransferOwnership(&_VRFCoordinatorV2.TransactOpts, to)
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.TransferOwnership(&_VRFCoordinatorV2.TransactOpts, to)
}

type VRFCoordinatorV2ConfigSetIterator struct {
	Event *VRFCoordinatorV2ConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2ConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2ConfigSet)
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
		it.Event = new(VRFCoordinatorV2ConfigSet)
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

func (it *VRFCoordinatorV2ConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2ConfigSet struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
	FallbackWeiPerUnitLink      *big.Int
	FeeConfig                   VRFCoordinatorV2FeeConfig
	Raw                         types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV2ConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2ConfigSetIterator{contract: _VRFCoordinatorV2.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2ConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2ConfigSet)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseConfigSet(log types.Log) (*VRFCoordinatorV2ConfigSet, error) {
	event := new(VRFCoordinatorV2ConfigSet)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2FundsRecoveredIterator struct {
	Event *VRFCoordinatorV2FundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2FundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2FundsRecovered)
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
		it.Event = new(VRFCoordinatorV2FundsRecovered)
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

func (it *VRFCoordinatorV2FundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2FundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2FundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV2FundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2FundsRecoveredIterator{contract: _VRFCoordinatorV2.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2FundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2FundsRecovered)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseFundsRecovered(log types.Log) (*VRFCoordinatorV2FundsRecovered, error) {
	event := new(VRFCoordinatorV2FundsRecovered)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2OwnershipTransferRequestedIterator struct {
	Event *VRFCoordinatorV2OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2OwnershipTransferRequested)
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
		it.Event = new(VRFCoordinatorV2OwnershipTransferRequested)
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

func (it *VRFCoordinatorV2OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2OwnershipTransferRequestedIterator{contract: _VRFCoordinatorV2.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2OwnershipTransferRequested)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV2OwnershipTransferRequested, error) {
	event := new(VRFCoordinatorV2OwnershipTransferRequested)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2OwnershipTransferredIterator struct {
	Event *VRFCoordinatorV2OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2OwnershipTransferred)
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
		it.Event = new(VRFCoordinatorV2OwnershipTransferred)
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

func (it *VRFCoordinatorV2OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2OwnershipTransferredIterator{contract: _VRFCoordinatorV2.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2OwnershipTransferred)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV2OwnershipTransferred, error) {
	event := new(VRFCoordinatorV2OwnershipTransferred)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2ProvingKeyDeregisteredIterator struct {
	Event *VRFCoordinatorV2ProvingKeyDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2ProvingKeyDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2ProvingKeyDeregistered)
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
		it.Event = new(VRFCoordinatorV2ProvingKeyDeregistered)
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

func (it *VRFCoordinatorV2ProvingKeyDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2ProvingKeyDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2ProvingKeyDeregistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2ProvingKeyDeregisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2ProvingKeyDeregisteredIterator{contract: _VRFCoordinatorV2.contract, event: "ProvingKeyDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2ProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2ProvingKeyDeregistered)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorV2ProvingKeyDeregistered, error) {
	event := new(VRFCoordinatorV2ProvingKeyDeregistered)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2ProvingKeyRegisteredIterator struct {
	Event *VRFCoordinatorV2ProvingKeyRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2ProvingKeyRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2ProvingKeyRegistered)
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
		it.Event = new(VRFCoordinatorV2ProvingKeyRegistered)
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

func (it *VRFCoordinatorV2ProvingKeyRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2ProvingKeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2ProvingKeyRegistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2ProvingKeyRegisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2ProvingKeyRegisteredIterator{contract: _VRFCoordinatorV2.contract, event: "ProvingKeyRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2ProvingKeyRegistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2ProvingKeyRegistered)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV2ProvingKeyRegistered, error) {
	event := new(VRFCoordinatorV2ProvingKeyRegistered)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2RandomWordsFulfilledIterator struct {
	Event *VRFCoordinatorV2RandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2RandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2RandomWordsFulfilled)
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
		it.Event = new(VRFCoordinatorV2RandomWordsFulfilled)
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

func (it *VRFCoordinatorV2RandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2RandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2RandomWordsFulfilled struct {
	RequestId  *big.Int
	OutputSeed *big.Int
	Payment    *big.Int
	Success    bool
	Raw        types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*VRFCoordinatorV2RandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2RandomWordsFulfilledIterator{contract: _VRFCoordinatorV2.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2RandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2RandomWordsFulfilled)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV2RandomWordsFulfilled, error) {
	event := new(VRFCoordinatorV2RandomWordsFulfilled)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2RandomWordsRequestedIterator struct {
	Event *VRFCoordinatorV2RandomWordsRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2RandomWordsRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2RandomWordsRequested)
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
		it.Event = new(VRFCoordinatorV2RandomWordsRequested)
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

func (it *VRFCoordinatorV2RandomWordsRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2RandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2RandomWordsRequested struct {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*VRFCoordinatorV2RandomWordsRequestedIterator, error) {

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

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2RandomWordsRequestedIterator{contract: _VRFCoordinatorV2.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2RandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2RandomWordsRequested)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV2RandomWordsRequested, error) {
	event := new(VRFCoordinatorV2RandomWordsRequested)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2SubscriptionCanceledIterator struct {
	Event *VRFCoordinatorV2SubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2SubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionCanceled)
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
		it.Event = new(VRFCoordinatorV2SubscriptionCanceled)
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

func (it *VRFCoordinatorV2SubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2SubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2SubscriptionCanceled struct {
	SubId  uint64
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionCanceledIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionCanceled, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2SubscriptionCanceled)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV2SubscriptionCanceled, error) {
	event := new(VRFCoordinatorV2SubscriptionCanceled)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2SubscriptionConsumerAddedIterator struct {
	Event *VRFCoordinatorV2SubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2SubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionConsumerAdded)
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
		it.Event = new(VRFCoordinatorV2SubscriptionConsumerAdded)
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

func (it *VRFCoordinatorV2SubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2SubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2SubscriptionConsumerAdded struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionConsumerAddedIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionConsumerAdded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2SubscriptionConsumerAdded)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV2SubscriptionConsumerAdded, error) {
	event := new(VRFCoordinatorV2SubscriptionConsumerAdded)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2SubscriptionConsumerRemovedIterator struct {
	Event *VRFCoordinatorV2SubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2SubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionConsumerRemoved)
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
		it.Event = new(VRFCoordinatorV2SubscriptionConsumerRemoved)
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

func (it *VRFCoordinatorV2SubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2SubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2SubscriptionConsumerRemoved struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionConsumerRemovedIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2SubscriptionConsumerRemoved)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV2SubscriptionConsumerRemoved, error) {
	event := new(VRFCoordinatorV2SubscriptionConsumerRemoved)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2SubscriptionCreatedIterator struct {
	Event *VRFCoordinatorV2SubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2SubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionCreated)
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
		it.Event = new(VRFCoordinatorV2SubscriptionCreated)
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

func (it *VRFCoordinatorV2SubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2SubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2SubscriptionCreated struct {
	SubId uint64
	Owner common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionCreatedIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionCreated, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2SubscriptionCreated)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV2SubscriptionCreated, error) {
	event := new(VRFCoordinatorV2SubscriptionCreated)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2SubscriptionFundedIterator struct {
	Event *VRFCoordinatorV2SubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2SubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionFunded)
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
		it.Event = new(VRFCoordinatorV2SubscriptionFunded)
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

func (it *VRFCoordinatorV2SubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2SubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2SubscriptionFunded struct {
	SubId      uint64
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionFundedIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionFunded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2SubscriptionFunded)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV2SubscriptionFunded, error) {
	event := new(VRFCoordinatorV2SubscriptionFunded)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator struct {
	Event *VRFCoordinatorV2SubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionOwnerTransferRequested)
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
		it.Event = new(VRFCoordinatorV2SubscriptionOwnerTransferRequested)
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

func (it *VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2SubscriptionOwnerTransferRequested struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2SubscriptionOwnerTransferRequested)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV2SubscriptionOwnerTransferRequested, error) {
	event := new(VRFCoordinatorV2SubscriptionOwnerTransferRequested)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2SubscriptionOwnerTransferredIterator struct {
	Event *VRFCoordinatorV2SubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2SubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2SubscriptionOwnerTransferred)
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
		it.Event = new(VRFCoordinatorV2SubscriptionOwnerTransferred)
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

func (it *VRFCoordinatorV2SubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2SubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2SubscriptionOwnerTransferred struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2SubscriptionOwnerTransferredIterator{contract: _VRFCoordinatorV2.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2SubscriptionOwnerTransferred)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV2SubscriptionOwnerTransferred, error) {
	event := new(VRFCoordinatorV2SubscriptionOwnerTransferred)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetConfig struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
}
type GetFeeConfig struct {
	FulfillmentFlatFeeLinkPPMTier1 uint32
	FulfillmentFlatFeeLinkPPMTier2 uint32
	FulfillmentFlatFeeLinkPPMTier3 uint32
	FulfillmentFlatFeeLinkPPMTier4 uint32
	FulfillmentFlatFeeLinkPPMTier5 uint32
	ReqsForTier2                   *big.Int
	ReqsForTier3                   *big.Int
	ReqsForTier4                   *big.Int
	ReqsForTier5                   *big.Int
}
type GetSubscription struct {
	Balance   *big.Int
	ReqCount  uint64
	Owner     common.Address
	Consumers []common.Address
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinatorV2.abi.Events["ConfigSet"].ID:
		return _VRFCoordinatorV2.ParseConfigSet(log)
	case _VRFCoordinatorV2.abi.Events["FundsRecovered"].ID:
		return _VRFCoordinatorV2.ParseFundsRecovered(log)
	case _VRFCoordinatorV2.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFCoordinatorV2.ParseOwnershipTransferRequested(log)
	case _VRFCoordinatorV2.abi.Events["OwnershipTransferred"].ID:
		return _VRFCoordinatorV2.ParseOwnershipTransferred(log)
	case _VRFCoordinatorV2.abi.Events["ProvingKeyDeregistered"].ID:
		return _VRFCoordinatorV2.ParseProvingKeyDeregistered(log)
	case _VRFCoordinatorV2.abi.Events["ProvingKeyRegistered"].ID:
		return _VRFCoordinatorV2.ParseProvingKeyRegistered(log)
	case _VRFCoordinatorV2.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFCoordinatorV2.ParseRandomWordsFulfilled(log)
	case _VRFCoordinatorV2.abi.Events["RandomWordsRequested"].ID:
		return _VRFCoordinatorV2.ParseRandomWordsRequested(log)
	case _VRFCoordinatorV2.abi.Events["SubscriptionCanceled"].ID:
		return _VRFCoordinatorV2.ParseSubscriptionCanceled(log)
	case _VRFCoordinatorV2.abi.Events["SubscriptionConsumerAdded"].ID:
		return _VRFCoordinatorV2.ParseSubscriptionConsumerAdded(log)
	case _VRFCoordinatorV2.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _VRFCoordinatorV2.ParseSubscriptionConsumerRemoved(log)
	case _VRFCoordinatorV2.abi.Events["SubscriptionCreated"].ID:
		return _VRFCoordinatorV2.ParseSubscriptionCreated(log)
	case _VRFCoordinatorV2.abi.Events["SubscriptionFunded"].ID:
		return _VRFCoordinatorV2.ParseSubscriptionFunded(log)
	case _VRFCoordinatorV2.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _VRFCoordinatorV2.ParseSubscriptionOwnerTransferRequested(log)
	case _VRFCoordinatorV2.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _VRFCoordinatorV2.ParseSubscriptionOwnerTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFCoordinatorV2ConfigSet) Topic() common.Hash {
	return common.HexToHash("0xc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb2")
}

func (VRFCoordinatorV2FundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (VRFCoordinatorV2OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFCoordinatorV2OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFCoordinatorV2ProvingKeyDeregistered) Topic() common.Hash {
	return common.HexToHash("0x72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d")
}

func (VRFCoordinatorV2ProvingKeyRegistered) Topic() common.Hash {
	return common.HexToHash("0xe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b8")
}

func (VRFCoordinatorV2RandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4")
}

func (VRFCoordinatorV2RandomWordsRequested) Topic() common.Hash {
	return common.HexToHash("0x63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772")
}

func (VRFCoordinatorV2SubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0xe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815")
}

func (VRFCoordinatorV2SubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0")
}

func (VRFCoordinatorV2SubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b")
}

func (VRFCoordinatorV2SubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf")
}

func (VRFCoordinatorV2SubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8")
}

func (VRFCoordinatorV2SubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be")
}

func (VRFCoordinatorV2SubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0")
}

func (_VRFCoordinatorV2 *VRFCoordinatorV2) Address() common.Address {
	return _VRFCoordinatorV2.address
}

type VRFCoordinatorV2Interface interface {
	BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	LINKETHFEED(opts *bind.CallOpts) (common.Address, error)

	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	MAXNUMWORDS(opts *bind.CallOpts) (uint32, error)

	MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error)

	GetCommitment(opts *bind.CallOpts, requestId *big.Int) ([32]byte, error)

	GetConfig(opts *bind.CallOpts) (GetConfig,

		error)

	GetCurrentSubId(opts *bind.CallOpts) (uint64, error)

	GetFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error)

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

	FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFCoordinatorV2RequestCommitment) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64, newOwner common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2FeeConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV2ConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2ConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFCoordinatorV2ConfigSet, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV2FundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2FundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*VRFCoordinatorV2FundsRecovered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV2OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV2OwnershipTransferred, error)

	FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2ProvingKeyDeregisteredIterator, error)

	WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2ProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorV2ProvingKeyDeregistered, error)

	FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2ProvingKeyRegisteredIterator, error)

	WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2ProvingKeyRegistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV2ProvingKeyRegistered, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*VRFCoordinatorV2RandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2RandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV2RandomWordsFulfilled, error)

	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*VRFCoordinatorV2RandomWordsRequestedIterator, error)

	WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2RandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV2RandomWordsRequested, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionCanceled, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV2SubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionConsumerAdded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV2SubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV2SubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionCreated, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV2SubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionFunded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV2SubscriptionFunded, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV2SubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2SubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2SubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV2SubscriptionOwnerTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
