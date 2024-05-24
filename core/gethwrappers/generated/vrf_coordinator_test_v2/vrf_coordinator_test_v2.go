// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_coordinator_test_v2

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

type VRFCoordinatorTestV2FeeConfig struct {
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

type VRFCoordinatorTestV2RequestCommitment struct {
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

var VRFCoordinatorTestV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"want\",\"type\":\"uint256\"}],\"name\":\"InsufficientGasForConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structVRFCoordinatorTestV2.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structVRFCoordinatorTestV2.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"getCommitment\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentSubId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"}],\"name\":\"getFeeTier\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"internalType\":\"structVRFCoordinatorTestV2.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162005755380380620057558339810160408190526200003491620001b1565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000e8565b5050506001600160601b0319606093841b811660805290831b811660a052911b1660c052620001fb565b6001600160a01b038116331415620001435760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001ac57600080fd5b919050565b600080600060608486031215620001c757600080fd5b620001d28462000194565b9250620001e26020850162000194565b9150620001f26040850162000194565b90509250925092565b60805160601c60a05160601c60c05160601c6154f0620002656000396000818161051901526138180152600081816106030152613c0401526000818161036d015281816114df0152818161233c01528181612d7301528181612eaf01526134d401526154f06000f3fe608060405234801561001057600080fd5b506004361061025b5760003560e01c80636f64f03f11610145578063ad178361116100bd578063d2f9f9a71161008c578063e72f6e3011610071578063e72f6e30146106e0578063e82ad7d4146106f3578063f2fde38b1461071657600080fd5b8063d2f9f9a7146106ba578063d7ae1d30146106cd57600080fd5b8063ad178361146105fe578063af198b9714610625578063c3f909d414610655578063caf70c4a146106a757600080fd5b80638da5cb5b11610114578063a21a23e4116100f9578063a21a23e4146105c0578063a47c7696146105c8578063a4c0ed36146105eb57600080fd5b80638da5cb5b1461059c5780639f87fad7146105ad57600080fd5b80636f64f03f1461055b5780637341c10c1461056e57806379ba509714610581578063823597401461058957600080fd5b8063356dac71116101d85780635fbbc0d2116101a757806366316d8d1161018c57806366316d8d14610501578063689c45171461051457806369bcdb7d1461053b57600080fd5b80635fbbc0d2146103f357806364d51a2a146104f957600080fd5b8063356dac71146103a757806340d6bb82146103af5780634cb48a54146103cd5780635d3b1d30146103e057600080fd5b806308821d581161022f57806315c48b841161021457806315c48b841461030e578063181f5a77146103295780631b6b6d231461036857600080fd5b806308821d58146102cf57806312b58349146102e257600080fd5b80620122911461026057806302bcc5b61461028057806304c357cb1461029557806306bfa637146102a8575b600080fd5b610268610729565b6040516102779392919061503a565b60405180910390f35b61029361028e366004614e86565b6107a5565b005b6102936102a3366004614ea1565b610837565b60055467ffffffffffffffff165b60405167ffffffffffffffff9091168152602001610277565b6102936102dd366004614b97565b6109eb565b6005546801000000000000000090046bffffffffffffffffffffffff165b604051908152602001610277565b61031660c881565b60405161ffff9091168152602001610277565b604080518082018252601681527f565246436f6f7264696e61746f72563220312e302e3000000000000000000000602082015290516102779190614fe5565b61038f7f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610277565b600a54610300565b6103b86101f481565b60405163ffffffff9091168152602001610277565b6102936103db366004614d30565b610bb0565b6103006103ee366004614c0a565b610fa7565b600c546040805163ffffffff80841682526401000000008404811660208301526801000000000000000084048116928201929092526c010000000000000000000000008304821660608201527001000000000000000000000000000000008304909116608082015262ffffff740100000000000000000000000000000000000000008304811660a0830152770100000000000000000000000000000000000000000000008304811660c08301527a0100000000000000000000000000000000000000000000000000008304811660e08301527d01000000000000000000000000000000000000000000000000000000000090920490911661010082015261012001610277565b610316606481565b61029361050f366004614b4f565b61138a565b61038f7f000000000000000000000000000000000000000000000000000000000000000081565b610300610549366004614e6d565b60009081526009602052604090205490565b610293610569366004614a94565b6115d9565b61029361057c366004614ea1565b611709565b610293611956565b610293610597366004614e86565b611a1f565b6000546001600160a01b031661038f565b6102936105bb366004614ea1565b611be5565b6102b6611fe4565b6105db6105d6366004614e86565b6121c7565b60405161027794939291906151d8565b6102936105f9366004614ac8565b6122ea565b61038f7f000000000000000000000000000000000000000000000000000000000000000081565b610638610633366004614c68565b612541565b6040516bffffffffffffffffffffffff9091168152602001610277565b600b546040805161ffff8316815263ffffffff6201000084048116602083015267010000000000000084048116928201929092526b010000000000000000000000909204166060820152608001610277565b6103006106b5366004614bb3565b6129db565b6103b86106c8366004614e86565b612a0b565b6102936106db366004614ea1565b612c00565b6102936106ee366004614a79565b612d3a565b610706610701366004614e86565b612f77565b6040519015158152602001610277565b610293610724366004614a79565b61319a565b600b546007805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff1693919283919083018282801561079357602002820191906000526020600020905b81548152602001906001019080831161077f575b50505050509050925092509250909192565b6107ad6131ab565b67ffffffffffffffff81166000908152600360205260409020546001600160a01b0316610806576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152600360205260409020546108349082906001600160a01b0316613207565b50565b67ffffffffffffffff821660009081526003602052604090205482906001600160a01b031680610893576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336001600160a01b038216146108e5576040517fd8a3fb520000000000000000000000000000000000000000000000000000000081526001600160a01b03821660048201526024015b60405180910390fd5b600b546601000000000000900460ff161561092c576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff84166000908152600360205260409020600101546001600160a01b038481169116146109e55767ffffffffffffffff841660008181526003602090815260409182902060010180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0388169081179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25b50505050565b6109f36131ab565b604080518082018252600091610a229190849060029083908390808284376000920191909152506129db915050565b6000818152600660205260409020549091506001600160a01b031680610a77576040517f77f5b84c000000000000000000000000000000000000000000000000000000008152600481018390526024016108dc565b600082815260066020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b600754811015610b67578260078281548110610aca57610aca615485565b90600052602060002001541415610b55576007805460009190610aef9060019061533f565b81548110610aff57610aff615485565b906000526020600020015490508060078381548110610b2057610b20615485565b6000918252602090912001556007805480610b3d57610b3d615456565b60019003818190600052602060002001600090559055505b80610b5f81615383565b915050610aac565b50806001600160a01b03167f72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d83604051610ba391815260200190565b60405180910390a2505050565b610bb86131ab565b60c861ffff87161115610c0b576040517fa738697600000000000000000000000000000000000000000000000000000000815261ffff871660048201819052602482015260c860448201526064016108dc565b60008213610c48576040517f43d4cf66000000000000000000000000000000000000000000000000000000008152600481018390526024016108dc565b6040805160a0808201835261ffff891680835263ffffffff89811660208086018290526000868801528a831660608088018290528b85166080988901819052600b80547fffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000001690971762010000909502949094177fffffffffffffffffffffffffffffffffff000000000000000000ffffffffffff166701000000000000009092027fffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffff16919091176b010000000000000000000000909302929092179093558651600c80549489015189890151938a0151978a0151968a015160c08b015160e08c01516101008d01519588167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009099169890981764010000000093881693909302929092177fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff1668010000000000000000958716959095027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff16949094176c0100000000000000000000000098861698909802979097177fffffffffffffffffff00000000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000096909416959095027fffffffffffffffffff000000ffffffffffffffffffffffffffffffffffffffff16929092177401000000000000000000000000000000000000000062ffffff92831602177fffffff000000000000ffffffffffffffffffffffffffffffffffffffffffffff1677010000000000000000000000000000000000000000000000958216959095027fffffff000000ffffffffffffffffffffffffffffffffffffffffffffffffffff16949094177a01000000000000000000000000000000000000000000000000000092851692909202919091177cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167d0100000000000000000000000000000000000000000000000000000000009390911692909202919091178155600a84905590517fc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb291610f97918991899189918991899190615099565b60405180910390a1505050505050565b600b546000906601000000000000900460ff1615610ff1576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff85166000908152600360205260409020546001600160a01b031661104a576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b33600090815260026020908152604080832067ffffffffffffffff808a16855292529091205416806110ba576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff871660048201523360248201526044016108dc565b600b5461ffff90811690861610806110d6575060c861ffff8616115b1561112657600b546040517fa738697600000000000000000000000000000000000000000000000000000000815261ffff8088166004830152909116602482015260c860448201526064016108dc565b600b5463ffffffff620100009091048116908516111561118d57600b546040517ff5d7e01e00000000000000000000000000000000000000000000000000000000815263ffffffff80871660048301526201000090920490911660248201526044016108dc565b6101f463ffffffff841611156111df576040517f47386bec00000000000000000000000000000000000000000000000000000000815263ffffffff841660048201526101f460248201526044016108dc565b60006111ec82600161529b565b6040805160208082018c9052338284015267ffffffffffffffff808c16606084015284166080808401919091528351808403909101815260a08301845280519082012060c083018d905260e080840182905284518085039091018152610100909301909352815191012091925060009182916040805160208101849052439181019190915267ffffffffffffffff8c16606082015263ffffffff808b166080830152891660a08201523360c0820152919350915060e00160408051808303601f19018152828252805160209182012060008681526009835283902055848352820183905261ffff8a169082015263ffffffff808916606083015287166080820152339067ffffffffffffffff8b16908c907f63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a97729060a00160405180910390a45033600090815260026020908152604080832067ffffffffffffffff808d16855292529091208054919093167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009091161790915591505095945050505050565b600b546601000000000000900460ff16156113d1576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600860205260409020546bffffffffffffffffffffffff8083169116101561142b576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b33600090815260086020526040812080548392906114589084906bffffffffffffffffffffffff16615356565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080600560088282829054906101000a90046bffffffffffffffffffffffff166114af9190615356565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb83836040518363ffffffff1660e01b815260040161154d9291906001600160a01b039290921682526bffffffffffffffffffffffff16602082015260400190565b602060405180830381600087803b15801561156757600080fd5b505af115801561157b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061159f9190614bcf565b6115d5576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b6115e16131ab565b6040805180820182526000916116109190849060029083908390808284376000920191909152506129db915050565b6000818152600660205260409020549091506001600160a01b031615611665576040517f4a0b8fa7000000000000000000000000000000000000000000000000000000008152600481018290526024016108dc565b600081815260066020908152604080832080547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0388169081179091556007805460018101825594527fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c688909301849055518381527fe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b89101610ba3565b67ffffffffffffffff821660009081526003602052604090205482906001600160a01b031680611765576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336001600160a01b038216146117b2576040517fd8a3fb520000000000000000000000000000000000000000000000000000000081526001600160a01b03821660048201526024016108dc565b600b546601000000000000900460ff16156117f9576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff841660009081526003602052604090206002015460641415611850576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b038316600090815260026020908152604080832067ffffffffffffffff8089168552925290912054161561188a576109e5565b6001600160a01b038316600081815260026020818152604080842067ffffffffffffffff8a1680865290835281852080547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000166001908117909155600384528286209094018054948501815585529382902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001685179055905192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e091016109dc565b6001546001600160a01b031633146119b05760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016108dc565b60008054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600b546601000000000000900460ff1615611a66576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152600360205260409020546001600160a01b0316611abf576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff81166000908152600360205260409020600101546001600160a01b03163314611b475767ffffffffffffffff8116600090815260036020526040908190206001015490517fd084e9750000000000000000000000000000000000000000000000000000000081526001600160a01b0390911660048201526024016108dc565b67ffffffffffffffff81166000818152600360209081526040918290208054337fffffffffffffffffffffffff0000000000000000000000000000000000000000808316821784556001909301805490931690925583516001600160a01b03909116808252928101919091529092917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0910160405180910390a25050565b67ffffffffffffffff821660009081526003602052604090205482906001600160a01b031680611c41576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336001600160a01b03821614611c8e576040517fd8a3fb520000000000000000000000000000000000000000000000000000000081526001600160a01b03821660048201526024016108dc565b600b546601000000000000900460ff1615611cd5576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b038316600090815260026020908152604080832067ffffffffffffffff808916855292529091205416611d56576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff851660048201526001600160a01b03841660248201526044016108dc565b67ffffffffffffffff8416600090815260036020908152604080832060020180548251818502810185019093528083529192909190830182828015611dc457602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611da6575b50505050509050600060018251611ddb919061533f565b905060005b8251811015611f5357856001600160a01b0316838281518110611e0557611e05615485565b60200260200101516001600160a01b03161415611f41576000838381518110611e3057611e30615485565b6020026020010151905080600360008a67ffffffffffffffff1667ffffffffffffffff1681526020019081526020016000206002018381548110611e7657611e76615485565b600091825260208083209190910180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03949094169390931790925567ffffffffffffffff8a168152600390915260409020600201805480611ee357611ee3615456565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905501905550611f53565b80611f4b81615383565b915050611de0565b506001600160a01b038516600081815260026020908152604080832067ffffffffffffffff8b168085529083529281902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b91015b60405180910390a2505050505050565b600b546000906601000000000000900460ff161561202e576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005805467ffffffffffffffff16906000612048836153bc565b82546101009290920a67ffffffffffffffff81810219909316918316021790915560055416905060008060405190808252806020026020018201604052801561209b578160200160208202803683370190505b506040805180820182526000808252602080830182815267ffffffffffffffff888116808552600484528685209551865493516bffffffffffffffffffffffff9091167fffffffffffffffffffffffff0000000000000000000000000000000000000000948516176c01000000000000000000000000919093160291909117909455845160608101865233815280830184815281870188815295855260038452959093208351815483166001600160a01b039182161782559551600182018054909316961695909517905591518051949550909361217f92600285019201906147d3565b505060405133815267ffffffffffffffff841691507f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a250905090565b67ffffffffffffffff8116600090815260036020526040812054819081906060906001600160a01b0316612227576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff80861660009081526004602090815260408083205460038352928190208054600290910180548351818602810186019094528084526bffffffffffffffffffffffff8616966c01000000000000000000000000909604909516946001600160a01b039092169390929183918301828280156122d457602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116122b6575b5050505050905093509350935093509193509193565b600b546601000000000000900460ff1615612331576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614612393576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602081146123cd576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006123db82840184614e86565b67ffffffffffffffff81166000908152600360205260409020549091506001600160a01b0316612437576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8116600090815260046020526040812080546bffffffffffffffffffffffff169186919061246e83856152c7565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555084600560088282829054906101000a90046bffffffffffffffffffffffff166124c591906152c7565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f882878461252c9190615283565b60408051928352602083019190915201611fd4565b600b546000906601000000000000900460ff161561258b576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005a9050600080600061259f878761362c565b9250925092506000866060015163ffffffff1667ffffffffffffffff8111156125ca576125ca6154b4565b6040519080825280602002602001820160405280156125f3578160200160208202803683370190505b50905060005b876060015163ffffffff168110156126675760408051602081018590529081018290526060016040516020818303038152906040528051906020012060001c82828151811061264a5761264a615485565b60209081029190910101528061265f81615383565b9150506125f9565b506000838152600960205260408082208290555181907f1fe543e300000000000000000000000000000000000000000000000000000000906126af908790869060240161518a565b60408051601f198184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090941693909317909252600b80547fffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffff166601000000000000179055908a015160808b015191925060009161275f9163ffffffff169084613936565b600b80547fffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffff1690556020808c01805167ffffffffffffffff9081166000908152600490935260408084205492518216845290922080549394506c01000000000000000000000000918290048316936001939192600c926127e392869290041661529b565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550600061283a8a600b600001600b9054906101000a900463ffffffff1663ffffffff1661283485612a0b565b3a613984565b6020808e015167ffffffffffffffff166000908152600490915260409020549091506bffffffffffffffffffffffff808316911610156128a6576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6020808d015167ffffffffffffffff16600090815260049091526040812080548392906128e29084906bffffffffffffffffffffffff16615356565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008b8152600660209081526040808320546001600160a01b03168352600890915281208054859450909261293e918591166152c7565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550877f7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e48883866040516129c1939291909283526bffffffffffffffffffffffff9190911660208301521515604082015260600190565b60405180910390a299505050505050505050505b92915050565b6000816040516020016129ee9190614fd7565b604051602081830303815290604052805190602001209050919050565b6040805161012081018252600c5463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c010000000000000000000000008104831660608301527001000000000000000000000000000000008104909216608082015262ffffff740100000000000000000000000000000000000000008304811660a08301819052770100000000000000000000000000000000000000000000008404821660c08401527a0100000000000000000000000000000000000000000000000000008404821660e08401527d0100000000000000000000000000000000000000000000000000000000009093041661010082015260009167ffffffffffffffff841611612b29575192915050565b8267ffffffffffffffff168160a0015162ffffff16108015612b5e57508060c0015162ffffff168367ffffffffffffffff1611155b15612b6d576020015192915050565b8267ffffffffffffffff168160c0015162ffffff16108015612ba257508060e0015162ffffff168367ffffffffffffffff1611155b15612bb1576040015192915050565b8267ffffffffffffffff168160e0015162ffffff16108015612be7575080610100015162ffffff168367ffffffffffffffff1611155b15612bf6576060015192915050565b6080015192915050565b67ffffffffffffffff821660009081526003602052604090205482906001600160a01b031680612c5c576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336001600160a01b03821614612ca9576040517fd8a3fb520000000000000000000000000000000000000000000000000000000081526001600160a01b03821660048201526024016108dc565b600b546601000000000000900460ff1615612cf0576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612cf984612f77565b15612d30576040517fb42f66e800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109e58484613207565b612d426131ab565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906370a082319060240160206040518083038186803b158015612dbd57600080fd5b505afa158015612dd1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612df59190614bf1565b6005549091506801000000000000000090046bffffffffffffffffffffffff1681811115612e59576040517fa99da30200000000000000000000000000000000000000000000000000000000815260048101829052602481018390526044016108dc565b81811015612f72576000612e6d828461533f565b6040517fa9059cbb0000000000000000000000000000000000000000000000000000000081526001600160a01b038681166004830152602482018390529192507f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb90604401602060405180830381600087803b158015612ef557600080fd5b505af1158015612f09573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f2d9190614bcf565b50604080516001600160a01b0386168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a1505b505050565b67ffffffffffffffff81166000908152600360209081526040808320815160608101835281546001600160a01b039081168252600183015416818501526002820180548451818702810187018652818152879693958601939092919083018282801561300c57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612fee575b505050505081525050905060005b8160400151518110156131905760005b60075481101561317d5760006131466007838154811061304c5761304c615485565b90600052602060002001548560400151858151811061306d5761306d615485565b602002602001015188600260008960400151898151811061309057613090615485565b6020908102919091018101516001600160a01b03168252818101929092526040908101600090812067ffffffffffffffff808f16835293522054166040805160208082018790526001600160a01b03959095168183015267ffffffffffffffff9384166060820152919092166080808301919091528251808303909101815260a08201835280519084012060c082019490945260e080820185905282518083039091018152610100909101909152805191012091565b506000818152600960205260409020549091501561316a5750600195945050505050565b508061317581615383565b91505061302a565b508061318881615383565b91505061301a565b5060009392505050565b6131a26131ab565b61083481613a8c565b6000546001600160a01b031633146132055760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016108dc565b565b600b546601000000000000900460ff161561324e576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff82166000908152600360209081526040808320815160608101835281546001600160a01b039081168252600183015416818501526002820180548451818702810187018652818152929593948601938301828280156132df57602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116132c1575b5050509190925250505067ffffffffffffffff80851660009081526004602090815260408083208151808301909252546bffffffffffffffffffffffff81168083526c01000000000000000000000000909104909416918101919091529293505b8360400151518110156133d957600260008560400151838151811061336757613367615485565b6020908102919091018101516001600160a01b03168252818101929092526040908101600090812067ffffffffffffffff8a168252909252902080547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000169055806133d181615383565b915050613340565b5067ffffffffffffffff8516600090815260036020526040812080547fffffffffffffffffffffffff000000000000000000000000000000000000000090811682556001820180549091169055906134346002830182614850565b505067ffffffffffffffff8516600090815260046020526040902080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055600580548291906008906134a49084906801000000000000000090046bffffffffffffffffffffffff16615356565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a9059cbb85836bffffffffffffffffffffffff166040518363ffffffff1660e01b81526004016135429291906001600160a01b03929092168252602082015260400190565b602060405180830381600087803b15801561355c57600080fd5b505af1158015613570573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135949190614bcf565b6135ca576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516001600160a01b03861681526bffffffffffffffffffffffff8316602082015267ffffffffffffffff8716917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815910160405180910390a25050505050565b600080600061363e85600001516129db565b6000818152600660205260409020549093506001600160a01b031680613693576040517f77f5b84c000000000000000000000000000000000000000000000000000000008152600481018590526024016108dc565b60808601516040516136b2918691602001918252602082015260400190565b60408051601f1981840301815291815281516020928301206000818152600990935291205490935080613711576040517f3688124a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b85516020808801516040808a015160608b015160808c0151925161377d968b96909594910195865267ffffffffffffffff948516602087015292909316604085015263ffffffff90811660608501529190911660808301526001600160a01b031660a082015260c00190565b6040516020818303038152906040528051906020012081146137cb576040517fd529142c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b855167ffffffffffffffff1640806138e25786516040517fe9413d3800000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063e9413d389060240160206040518083038186803b15801561386257600080fd5b505afa158015613876573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061389a9190614bf1565b9050806138e25786516040517f175dadad00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201526024016108dc565b6000886080015182604051602001613904929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c90506139298982613b4e565b9450505050509250925092565b60005a61138881101561394857600080fd5b61138881039050846040820482031161396057600080fd5b50823b61396c57600080fd5b60008083516020850160008789f190505b9392505050565b60008061398f613bb9565b9050600081136139ce576040517f43d4cf66000000000000000000000000000000000000000000000000000000008152600481018290526024016108dc565b6000815a6139dc8989615283565b6139e6919061533f565b6139f886670de0b6b3a7640000615302565b613a029190615302565b613a0c91906152ee565b90506000613a2563ffffffff871664e8d4a51000615302565b9050613a3d816b033b2e3c9fd0803ce800000061533f565b821115613a76576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b613a808183615283565b98975050505050505050565b6001600160a01b038116331415613ae55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016108dc565b600180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000613b828360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151613cc0565b60038360200151604051602001613b9a929190615176565b60408051601f1981840301815291905280516020909101209392505050565b600b54604080517ffeaf968c0000000000000000000000000000000000000000000000000000000081529051600092670100000000000000900463ffffffff169182151591849182917f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169163feaf968c9160048083019260a0929190829003018186803b158015613c5257600080fd5b505afa158015613c66573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613c8a9190614ecb565b509450909250849150508015613cae5750613ca5824261533f565b8463ffffffff16105b15613cb85750600a545b949350505050565b613cc989613efb565b613d155760405162461bcd60e51b815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e20637572766500000000000060448201526064016108dc565b613d1e88613efb565b613d6a5760405162461bcd60e51b815260206004820152601560248201527f67616d6d61206973206e6f74206f6e206375727665000000000000000000000060448201526064016108dc565b613d7383613efb565b613dbf5760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e20637572766500000060448201526064016108dc565b613dc882613efb565b613e145760405162461bcd60e51b815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e2063757276650000000060448201526064016108dc565b613e20878a8887613fd4565b613e6c5760405162461bcd60e51b815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e6573730000000000000060448201526064016108dc565b6000613e788a87614125565b90506000613e8b898b878b868989614189565b90506000613e9c838d8d8a866142a9565b9050808a14613eed5760405162461bcd60e51b815260206004820152600d60248201527f696e76616c69642070726f6f660000000000000000000000000000000000000060448201526064016108dc565b505050505050505050505050565b80516000906401000003d01911613f545760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420782d6f7264696e617465000000000000000000000000000060448201526064016108dc565b60208201516401000003d01911613fad5760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420792d6f7264696e617465000000000000000000000000000060448201526064016108dc565b60208201516401000003d019908009613fcd8360005b60200201516142e9565b1492915050565b60006001600160a01b03821661402c5760405162461bcd60e51b815260206004820152600b60248201527f626164207769746e65737300000000000000000000000000000000000000000060448201526064016108dc565b60208401516000906001161561404357601c614046565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418587600060200201510986517ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa1580156140fd573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b61412d61486e565b61415a6001848460405160200161414693929190614fb6565b60405160208183030381529060405261430d565b90505b61416681613efb565b6129d55780516040805160208101929092526141829101614146565b905061415d565b61419161486e565b825186516401000003d01990819006910614156141f05760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e6374000060448201526064016108dc565b6141fb87898861435c565b6142475760405162461bcd60e51b815260206004820152601660248201527f4669727374206d756c20636865636b206661696c65640000000000000000000060448201526064016108dc565b61425284868561435c565b61429e5760405162461bcd60e51b815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c656400000000000000000060448201526064016108dc565b613a808684846144a4565b6000600286868685876040516020016142c796959493929190614f44565b60408051601f1981840301815291905280516020909101209695505050505050565b6000806401000003d01980848509840990506401000003d019600782089392505050565b61431561486e565b61431e8261456b565b815261433361432e826000613fc3565b6145a6565b602082018190526002900660011415614357576020810180516401000003d0190390525b919050565b6000826143ab5760405162461bcd60e51b815260206004820152600b60248201527f7a65726f207363616c617200000000000000000000000000000000000000000060448201526064016108dc565b835160208501516000906143c1906002906153e4565b156143cd57601c6143d0565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa158015614450573d6000803e3d6000fd5b50505060206040510351905060008660405160200161446f9190614f32565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b6144ac61486e565b8351602080860151855191860151600093849384936144cd939091906145c6565b919450925090506401000003d01985820960011461452d5760405162461bcd60e51b815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a0000000000000060448201526064016108dc565b60405180604001604052806401000003d0198061454c5761454c615427565b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d019811061435757604080516020808201939093528151808203840181529082019091528051910120614573565b60006129d58260026145bf6401000003d0196001615283565b901c6146a6565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a089050600061460683838585614766565b909850905061461788828e8861478a565b909850905061462888828c8761478a565b9098509050600061463b8d878b8561478a565b909850905061464c88828686614766565b909850905061465d88828e8961478a565b9098509050818114614692576401000003d019818a0998506401000003d01982890997506401000003d0198183099650614696565b8196505b5050505050509450945094915050565b6000806146b161488c565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a08201526146e36148aa565b60208160c08460057ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa92508261475c5760405162461bcd60e51b815260206004820152601260248201527f6269674d6f64457870206661696c75726521000000000000000000000000000060448201526064016108dc565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215614840579160200282015b8281111561484057825182547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b039091161782556020909201916001909101906147f3565b5061484c9291506148c8565b5090565b508054600082559060005260206000209081019061083491906148c8565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b8082111561484c57600081556001016148c9565b80356001600160a01b038116811461435757600080fd5b80604081018310156129d557600080fd5b600082601f83011261491657600080fd5b6040516040810181811067ffffffffffffffff82111715614939576149396154b4565b806040525080838560408601111561495057600080fd5b60005b6002811015614972578135835260209283019290910190600101614953565b509195945050505050565b600060a0828403121561498f57600080fd5b60405160a0810181811067ffffffffffffffff821117156149b2576149b26154b4565b6040529050806149c183614a47565b81526149cf60208401614a47565b60208201526149e060408401614a33565b60408201526149f160608401614a33565b6060820152614a02608084016148dd565b60808201525092915050565b803561ffff8116811461435757600080fd5b803562ffffff8116811461435757600080fd5b803563ffffffff8116811461435757600080fd5b803567ffffffffffffffff8116811461435757600080fd5b805169ffffffffffffffffffff8116811461435757600080fd5b600060208284031215614a8b57600080fd5b61397d826148dd565b60008060608385031215614aa757600080fd5b614ab0836148dd565b9150614abf84602085016148f4565b90509250929050565b60008060008060608587031215614ade57600080fd5b614ae7856148dd565b935060208501359250604085013567ffffffffffffffff80821115614b0b57600080fd5b818701915087601f830112614b1f57600080fd5b813581811115614b2e57600080fd5b886020828501011115614b4057600080fd5b95989497505060200194505050565b60008060408385031215614b6257600080fd5b614b6b836148dd565b915060208301356bffffffffffffffffffffffff81168114614b8c57600080fd5b809150509250929050565b600060408284031215614ba957600080fd5b61397d83836148f4565b600060408284031215614bc557600080fd5b61397d8383614905565b600060208284031215614be157600080fd5b8151801515811461397d57600080fd5b600060208284031215614c0357600080fd5b5051919050565b600080600080600060a08688031215614c2257600080fd5b85359450614c3260208701614a47565b9350614c4060408701614a0e565b9250614c4e60608701614a33565b9150614c5c60808701614a33565b90509295509295909350565b600080828403610240811215614c7d57600080fd5b6101a080821215614c8d57600080fd5b614c95615259565b9150614ca18686614905565b8252614cb08660408701614905565b60208301526080850135604083015260a0850135606083015260c08501356080830152614cdf60e086016148dd565b60a0830152610100614cf387828801614905565b60c0840152614d06876101408801614905565b60e08401526101808601358184015250819350614d258682870161497d565b925050509250929050565b6000806000806000808688036101c0811215614d4b57600080fd5b614d5488614a0e565b9650614d6260208901614a33565b9550614d7060408901614a33565b9450614d7e60608901614a33565b935060808801359250610120807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6083011215614db957600080fd5b614dc1615259565b9150614dcf60a08a01614a33565b8252614ddd60c08a01614a33565b6020830152614dee60e08a01614a33565b6040830152610100614e01818b01614a33565b6060840152614e11828b01614a33565b6080840152614e236101408b01614a20565b60a0840152614e356101608b01614a20565b60c0840152614e476101808b01614a20565b60e0840152614e596101a08b01614a20565b818401525050809150509295509295509295565b600060208284031215614e7f57600080fd5b5035919050565b600060208284031215614e9857600080fd5b61397d82614a47565b60008060408385031215614eb457600080fd5b614ebd83614a47565b9150614abf602084016148dd565b600080600080600060a08688031215614ee357600080fd5b614eec86614a5f565b9450602086015193506040860151925060608601519150614c5c60808701614a5f565b8060005b60028110156109e5578151845260209384019390910190600101614f13565b614f3c8183614f0f565b604001919050565b868152614f546020820187614f0f565b614f616060820186614f0f565b614f6e60a0820185614f0f565b614f7b60e0820184614f0f565b60609190911b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166101208201526101340195945050505050565b838152614fc66020820184614f0f565b606081019190915260800192915050565b604081016129d58284614f0f565b600060208083528351808285015260005b8181101561501257858101830151858201604001528201614ff6565b81811115615024576000604083870101525b50601f01601f1916929092016040019392505050565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b8181101561508b5784518352938301939183019160010161506f565b509098975050505050505050565b60006101c08201905061ffff8816825263ffffffff808816602084015280871660408401528086166060840152846080840152835481811660a08501526150ed60c08501838360201c1663ffffffff169052565b61510460e08501838360401c1663ffffffff169052565b61511c6101008501838360601c1663ffffffff169052565b6151346101208501838360801c1663ffffffff169052565b62ffffff60a082901c811661014086015260b882901c811661016086015260d082901c1661018085015260e81c6101a090930192909252979650505050505050565b8281526060810161397d6020830184614f0f565b6000604082018483526020604081850152818551808452606086019150828701935060005b818110156151cb578451835293830193918301916001016151af565b5090979650505050505050565b6000608082016bffffffffffffffffffffffff87168352602067ffffffffffffffff8716818501526001600160a01b0380871660408601526080606086015282865180855260a087019150838801945060005b8181101561524957855184168352948401949184019160010161522b565b50909a9950505050505050505050565b604051610120810167ffffffffffffffff8111828210171561527d5761527d6154b4565b60405290565b60008219821115615296576152966153f8565b500190565b600067ffffffffffffffff8083168185168083038211156152be576152be6153f8565b01949350505050565b60006bffffffffffffffffffffffff8083168185168083038211156152be576152be6153f8565b6000826152fd576152fd615427565b500490565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561533a5761533a6153f8565b500290565b600082821015615351576153516153f8565b500390565b60006bffffffffffffffffffffffff8381169083168181101561537b5761537b6153f8565b039392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156153b5576153b56153f8565b5060010190565b600067ffffffffffffffff808316818114156153da576153da6153f8565b6001019392505050565b6000826153f3576153f3615427565b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var VRFCoordinatorTestV2ABI = VRFCoordinatorTestV2MetaData.ABI

var VRFCoordinatorTestV2Bin = VRFCoordinatorTestV2MetaData.Bin

func DeployVRFCoordinatorTestV2(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, blockhashStore common.Address, linkEthFeed common.Address) (common.Address, *types.Transaction, *VRFCoordinatorTestV2, error) {
	parsed, err := VRFCoordinatorTestV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorTestV2Bin), backend, link, blockhashStore, linkEthFeed)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorTestV2{address: address, abi: *parsed, VRFCoordinatorTestV2Caller: VRFCoordinatorTestV2Caller{contract: contract}, VRFCoordinatorTestV2Transactor: VRFCoordinatorTestV2Transactor{contract: contract}, VRFCoordinatorTestV2Filterer: VRFCoordinatorTestV2Filterer{contract: contract}}, nil
}

type VRFCoordinatorTestV2 struct {
	address common.Address
	abi     abi.ABI
	VRFCoordinatorTestV2Caller
	VRFCoordinatorTestV2Transactor
	VRFCoordinatorTestV2Filterer
}

type VRFCoordinatorTestV2Caller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorTestV2Transactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorTestV2Filterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorTestV2Session struct {
	Contract     *VRFCoordinatorTestV2
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorTestV2CallerSession struct {
	Contract *VRFCoordinatorTestV2Caller
	CallOpts bind.CallOpts
}

type VRFCoordinatorTestV2TransactorSession struct {
	Contract     *VRFCoordinatorTestV2Transactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorTestV2Raw struct {
	Contract *VRFCoordinatorTestV2
}

type VRFCoordinatorTestV2CallerRaw struct {
	Contract *VRFCoordinatorTestV2Caller
}

type VRFCoordinatorTestV2TransactorRaw struct {
	Contract *VRFCoordinatorTestV2Transactor
}

func NewVRFCoordinatorTestV2(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorTestV2, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorTestV2ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFCoordinatorTestV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2{address: address, abi: abi, VRFCoordinatorTestV2Caller: VRFCoordinatorTestV2Caller{contract: contract}, VRFCoordinatorTestV2Transactor: VRFCoordinatorTestV2Transactor{contract: contract}, VRFCoordinatorTestV2Filterer: VRFCoordinatorTestV2Filterer{contract: contract}}, nil
}

func NewVRFCoordinatorTestV2Caller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorTestV2Caller, error) {
	contract, err := bindVRFCoordinatorTestV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2Caller{contract: contract}, nil
}

func NewVRFCoordinatorTestV2Transactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorTestV2Transactor, error) {
	contract, err := bindVRFCoordinatorTestV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2Transactor{contract: contract}, nil
}

func NewVRFCoordinatorTestV2Filterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorTestV2Filterer, error) {
	contract, err := bindVRFCoordinatorTestV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2Filterer{contract: contract}, nil
}

func bindVRFCoordinatorTestV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFCoordinatorTestV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorTestV2.Contract.VRFCoordinatorTestV2Caller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.VRFCoordinatorTestV2Transactor.contract.Transfer(opts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.VRFCoordinatorTestV2Transactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorTestV2.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.contract.Transfer(opts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "BLOCKHASH_STORE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorTestV2.Contract.BLOCKHASHSTORE(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorTestV2.Contract.BLOCKHASHSTORE(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) LINK() (common.Address, error) {
	return _VRFCoordinatorTestV2.Contract.LINK(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) LINK() (common.Address, error) {
	return _VRFCoordinatorTestV2.Contract.LINK(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) LINKETHFEED() (common.Address, error) {
	return _VRFCoordinatorTestV2.Contract.LINKETHFEED(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) LINKETHFEED() (common.Address, error) {
	return _VRFCoordinatorTestV2.Contract.LINKETHFEED(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorTestV2.Contract.MAXCONSUMERS(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorTestV2.Contract.MAXCONSUMERS(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) MAXNUMWORDS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "MAX_NUM_WORDS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorTestV2.Contract.MAXNUMWORDS(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorTestV2.Contract.MAXNUMWORDS(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "MAX_REQUEST_CONFIRMATIONS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorTestV2.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorTestV2.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) GetCommitment(opts *bind.CallOpts, requestId *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "getCommitment", requestId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) GetCommitment(requestId *big.Int) ([32]byte, error) {
	return _VRFCoordinatorTestV2.Contract.GetCommitment(&_VRFCoordinatorTestV2.CallOpts, requestId)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) GetCommitment(requestId *big.Int) ([32]byte, error) {
	return _VRFCoordinatorTestV2.Contract.GetCommitment(&_VRFCoordinatorTestV2.CallOpts, requestId)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) GetConfig(opts *bind.CallOpts) (GetConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "getConfig")

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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) GetConfig() (GetConfig,

	error) {
	return _VRFCoordinatorTestV2.Contract.GetConfig(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) GetConfig() (GetConfig,

	error) {
	return _VRFCoordinatorTestV2.Contract.GetConfig(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) GetCurrentSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "getCurrentSubId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) GetCurrentSubId() (uint64, error) {
	return _VRFCoordinatorTestV2.Contract.GetCurrentSubId(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) GetCurrentSubId() (uint64, error) {
	return _VRFCoordinatorTestV2.Contract.GetCurrentSubId(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) GetFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "getFallbackWeiPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) GetFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorTestV2.Contract.GetFallbackWeiPerUnitLink(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) GetFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorTestV2.Contract.GetFallbackWeiPerUnitLink(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) GetFeeConfig(opts *bind.CallOpts) (GetFeeConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "getFeeConfig")

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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) GetFeeConfig() (GetFeeConfig,

	error) {
	return _VRFCoordinatorTestV2.Contract.GetFeeConfig(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) GetFeeConfig() (GetFeeConfig,

	error) {
	return _VRFCoordinatorTestV2.Contract.GetFeeConfig(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) GetFeeTier(opts *bind.CallOpts, reqCount uint64) (uint32, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "getFeeTier", reqCount)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) GetFeeTier(reqCount uint64) (uint32, error) {
	return _VRFCoordinatorTestV2.Contract.GetFeeTier(&_VRFCoordinatorTestV2.CallOpts, reqCount)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) GetFeeTier(reqCount uint64) (uint32, error) {
	return _VRFCoordinatorTestV2.Contract.GetFeeTier(&_VRFCoordinatorTestV2.CallOpts, reqCount)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) GetRequestConfig(opts *bind.CallOpts) (uint16, uint32, [][32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "getRequestConfig")

	if err != nil {
		return *new(uint16), *new(uint32), *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)
	out2 := *abi.ConvertType(out[2], new([][32]byte)).(*[][32]byte)

	return out0, out1, out2, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _VRFCoordinatorTestV2.Contract.GetRequestConfig(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _VRFCoordinatorTestV2.Contract.GetRequestConfig(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) GetSubscription(opts *bind.CallOpts, subId uint64) (GetSubscription,

	error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "getSubscription", subId)

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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) GetSubscription(subId uint64) (GetSubscription,

	error) {
	return _VRFCoordinatorTestV2.Contract.GetSubscription(&_VRFCoordinatorTestV2.CallOpts, subId)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) GetSubscription(subId uint64) (GetSubscription,

	error) {
	return _VRFCoordinatorTestV2.Contract.GetSubscription(&_VRFCoordinatorTestV2.CallOpts, subId)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) GetTotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "getTotalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) GetTotalBalance() (*big.Int, error) {
	return _VRFCoordinatorTestV2.Contract.GetTotalBalance(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) GetTotalBalance() (*big.Int, error) {
	return _VRFCoordinatorTestV2.Contract.GetTotalBalance(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "hashOfKey", publicKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorTestV2.Contract.HashOfKey(&_VRFCoordinatorTestV2.CallOpts, publicKey)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorTestV2.Contract.HashOfKey(&_VRFCoordinatorTestV2.CallOpts, publicKey)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) Owner() (common.Address, error) {
	return _VRFCoordinatorTestV2.Contract.Owner(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) Owner() (common.Address, error) {
	return _VRFCoordinatorTestV2.Contract.Owner(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) PendingRequestExists(opts *bind.CallOpts, subId uint64) (bool, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "pendingRequestExists", subId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) PendingRequestExists(subId uint64) (bool, error) {
	return _VRFCoordinatorTestV2.Contract.PendingRequestExists(&_VRFCoordinatorTestV2.CallOpts, subId)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) PendingRequestExists(subId uint64) (bool, error) {
	return _VRFCoordinatorTestV2.Contract.PendingRequestExists(&_VRFCoordinatorTestV2.CallOpts, subId)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VRFCoordinatorTestV2.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) TypeAndVersion() (string, error) {
	return _VRFCoordinatorTestV2.Contract.TypeAndVersion(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2CallerSession) TypeAndVersion() (string, error) {
	return _VRFCoordinatorTestV2.Contract.TypeAndVersion(&_VRFCoordinatorTestV2.CallOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "acceptOwnership")
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.AcceptOwnership(&_VRFCoordinatorTestV2.TransactOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.AcceptOwnership(&_VRFCoordinatorTestV2.TransactOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorTestV2.TransactOpts, subId)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorTestV2.TransactOpts, subId)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) AddConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "addConsumer", subId, consumer)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.AddConsumer(&_VRFCoordinatorTestV2.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.AddConsumer(&_VRFCoordinatorTestV2.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) CancelSubscription(opts *bind.TransactOpts, subId uint64, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "cancelSubscription", subId, to)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) CancelSubscription(subId uint64, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.CancelSubscription(&_VRFCoordinatorTestV2.TransactOpts, subId, to)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) CancelSubscription(subId uint64, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.CancelSubscription(&_VRFCoordinatorTestV2.TransactOpts, subId, to)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "createSubscription")
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.CreateSubscription(&_VRFCoordinatorTestV2.TransactOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.CreateSubscription(&_VRFCoordinatorTestV2.TransactOpts)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "deregisterProvingKey", publicProvingKey)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.DeregisterProvingKey(&_VRFCoordinatorTestV2.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.DeregisterProvingKey(&_VRFCoordinatorTestV2.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFCoordinatorTestV2RequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "fulfillRandomWords", proof, rc)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorTestV2RequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.FulfillRandomWords(&_VRFCoordinatorTestV2.TransactOpts, proof, rc)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorTestV2RequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.FulfillRandomWords(&_VRFCoordinatorTestV2.TransactOpts, proof, rc)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.OnTokenTransfer(&_VRFCoordinatorTestV2.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.OnTokenTransfer(&_VRFCoordinatorTestV2.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "oracleWithdraw", recipient, amount)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.OracleWithdraw(&_VRFCoordinatorTestV2.TransactOpts, recipient, amount)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.OracleWithdraw(&_VRFCoordinatorTestV2.TransactOpts, recipient, amount)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "ownerCancelSubscription", subId)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.OwnerCancelSubscription(&_VRFCoordinatorTestV2.TransactOpts, subId)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.OwnerCancelSubscription(&_VRFCoordinatorTestV2.TransactOpts, subId)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "recoverFunds", to)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.RecoverFunds(&_VRFCoordinatorTestV2.TransactOpts, to)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.RecoverFunds(&_VRFCoordinatorTestV2.TransactOpts, to)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "registerProvingKey", oracle, publicProvingKey)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.RegisterProvingKey(&_VRFCoordinatorTestV2.TransactOpts, oracle, publicProvingKey)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.RegisterProvingKey(&_VRFCoordinatorTestV2.TransactOpts, oracle, publicProvingKey)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) RemoveConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "removeConsumer", subId, consumer)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) RemoveConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.RemoveConsumer(&_VRFCoordinatorTestV2.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) RemoveConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.RemoveConsumer(&_VRFCoordinatorTestV2.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "requestRandomWords", keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) RequestRandomWords(keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.RequestRandomWords(&_VRFCoordinatorTestV2.TransactOpts, keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) RequestRandomWords(keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.RequestRandomWords(&_VRFCoordinatorTestV2.TransactOpts, keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorTestV2.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorTestV2.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorTestV2FeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorTestV2FeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.SetConfig(&_VRFCoordinatorTestV2.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorTestV2FeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.SetConfig(&_VRFCoordinatorTestV2.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.TransferOwnership(&_VRFCoordinatorTestV2.TransactOpts, to)
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorTestV2.Contract.TransferOwnership(&_VRFCoordinatorTestV2.TransactOpts, to)
}

type VRFCoordinatorTestV2ConfigSetIterator struct {
	Event *VRFCoordinatorTestV2ConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2ConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2ConfigSet)
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
		it.Event = new(VRFCoordinatorTestV2ConfigSet)
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

func (it *VRFCoordinatorTestV2ConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2ConfigSet struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
	FallbackWeiPerUnitLink      *big.Int
	FeeConfig                   VRFCoordinatorTestV2FeeConfig
	Raw                         types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorTestV2ConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2ConfigSetIterator{contract: _VRFCoordinatorTestV2.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2ConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2ConfigSet)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseConfigSet(log types.Log) (*VRFCoordinatorTestV2ConfigSet, error) {
	event := new(VRFCoordinatorTestV2ConfigSet)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2FundsRecoveredIterator struct {
	Event *VRFCoordinatorTestV2FundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2FundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2FundsRecovered)
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
		it.Event = new(VRFCoordinatorTestV2FundsRecovered)
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

func (it *VRFCoordinatorTestV2FundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2FundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2FundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorTestV2FundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2FundsRecoveredIterator{contract: _VRFCoordinatorTestV2.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2FundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2FundsRecovered)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseFundsRecovered(log types.Log) (*VRFCoordinatorTestV2FundsRecovered, error) {
	event := new(VRFCoordinatorTestV2FundsRecovered)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2OwnershipTransferRequestedIterator struct {
	Event *VRFCoordinatorTestV2OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2OwnershipTransferRequested)
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
		it.Event = new(VRFCoordinatorTestV2OwnershipTransferRequested)
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

func (it *VRFCoordinatorTestV2OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorTestV2OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2OwnershipTransferRequestedIterator{contract: _VRFCoordinatorTestV2.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2OwnershipTransferRequested)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorTestV2OwnershipTransferRequested, error) {
	event := new(VRFCoordinatorTestV2OwnershipTransferRequested)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2OwnershipTransferredIterator struct {
	Event *VRFCoordinatorTestV2OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2OwnershipTransferred)
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
		it.Event = new(VRFCoordinatorTestV2OwnershipTransferred)
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

func (it *VRFCoordinatorTestV2OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorTestV2OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2OwnershipTransferredIterator{contract: _VRFCoordinatorTestV2.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2OwnershipTransferred)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorTestV2OwnershipTransferred, error) {
	event := new(VRFCoordinatorTestV2OwnershipTransferred)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2ProvingKeyDeregisteredIterator struct {
	Event *VRFCoordinatorTestV2ProvingKeyDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2ProvingKeyDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2ProvingKeyDeregistered)
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
		it.Event = new(VRFCoordinatorTestV2ProvingKeyDeregistered)
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

func (it *VRFCoordinatorTestV2ProvingKeyDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2ProvingKeyDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2ProvingKeyDeregistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorTestV2ProvingKeyDeregisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2ProvingKeyDeregisteredIterator{contract: _VRFCoordinatorTestV2.contract, event: "ProvingKeyDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2ProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2ProvingKeyDeregistered)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorTestV2ProvingKeyDeregistered, error) {
	event := new(VRFCoordinatorTestV2ProvingKeyDeregistered)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2ProvingKeyRegisteredIterator struct {
	Event *VRFCoordinatorTestV2ProvingKeyRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2ProvingKeyRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2ProvingKeyRegistered)
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
		it.Event = new(VRFCoordinatorTestV2ProvingKeyRegistered)
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

func (it *VRFCoordinatorTestV2ProvingKeyRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2ProvingKeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2ProvingKeyRegistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorTestV2ProvingKeyRegisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2ProvingKeyRegisteredIterator{contract: _VRFCoordinatorTestV2.contract, event: "ProvingKeyRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2ProvingKeyRegistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2ProvingKeyRegistered)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorTestV2ProvingKeyRegistered, error) {
	event := new(VRFCoordinatorTestV2ProvingKeyRegistered)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2RandomWordsFulfilledIterator struct {
	Event *VRFCoordinatorTestV2RandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2RandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2RandomWordsFulfilled)
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
		it.Event = new(VRFCoordinatorTestV2RandomWordsFulfilled)
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

func (it *VRFCoordinatorTestV2RandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2RandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2RandomWordsFulfilled struct {
	RequestId  *big.Int
	OutputSeed *big.Int
	Payment    *big.Int
	Success    bool
	Raw        types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*VRFCoordinatorTestV2RandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2RandomWordsFulfilledIterator{contract: _VRFCoordinatorTestV2.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2RandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2RandomWordsFulfilled)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorTestV2RandomWordsFulfilled, error) {
	event := new(VRFCoordinatorTestV2RandomWordsFulfilled)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2RandomWordsRequestedIterator struct {
	Event *VRFCoordinatorTestV2RandomWordsRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2RandomWordsRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2RandomWordsRequested)
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
		it.Event = new(VRFCoordinatorTestV2RandomWordsRequested)
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

func (it *VRFCoordinatorTestV2RandomWordsRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2RandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2RandomWordsRequested struct {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*VRFCoordinatorTestV2RandomWordsRequestedIterator, error) {

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

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2RandomWordsRequestedIterator{contract: _VRFCoordinatorTestV2.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2RandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2RandomWordsRequested)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorTestV2RandomWordsRequested, error) {
	event := new(VRFCoordinatorTestV2RandomWordsRequested)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2SubscriptionCanceledIterator struct {
	Event *VRFCoordinatorTestV2SubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2SubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2SubscriptionCanceled)
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
		it.Event = new(VRFCoordinatorTestV2SubscriptionCanceled)
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

func (it *VRFCoordinatorTestV2SubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2SubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2SubscriptionCanceled struct {
	SubId  uint64
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2SubscriptionCanceledIterator{contract: _VRFCoordinatorTestV2.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionCanceled, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2SubscriptionCanceled)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorTestV2SubscriptionCanceled, error) {
	event := new(VRFCoordinatorTestV2SubscriptionCanceled)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2SubscriptionConsumerAddedIterator struct {
	Event *VRFCoordinatorTestV2SubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2SubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2SubscriptionConsumerAdded)
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
		it.Event = new(VRFCoordinatorTestV2SubscriptionConsumerAdded)
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

func (it *VRFCoordinatorTestV2SubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2SubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2SubscriptionConsumerAdded struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2SubscriptionConsumerAddedIterator{contract: _VRFCoordinatorTestV2.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionConsumerAdded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2SubscriptionConsumerAdded)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorTestV2SubscriptionConsumerAdded, error) {
	event := new(VRFCoordinatorTestV2SubscriptionConsumerAdded)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2SubscriptionConsumerRemovedIterator struct {
	Event *VRFCoordinatorTestV2SubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2SubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2SubscriptionConsumerRemoved)
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
		it.Event = new(VRFCoordinatorTestV2SubscriptionConsumerRemoved)
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

func (it *VRFCoordinatorTestV2SubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2SubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2SubscriptionConsumerRemoved struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2SubscriptionConsumerRemovedIterator{contract: _VRFCoordinatorTestV2.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2SubscriptionConsumerRemoved)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorTestV2SubscriptionConsumerRemoved, error) {
	event := new(VRFCoordinatorTestV2SubscriptionConsumerRemoved)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2SubscriptionCreatedIterator struct {
	Event *VRFCoordinatorTestV2SubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2SubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2SubscriptionCreated)
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
		it.Event = new(VRFCoordinatorTestV2SubscriptionCreated)
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

func (it *VRFCoordinatorTestV2SubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2SubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2SubscriptionCreated struct {
	SubId uint64
	Owner common.Address
	Raw   types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2SubscriptionCreatedIterator{contract: _VRFCoordinatorTestV2.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionCreated, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2SubscriptionCreated)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorTestV2SubscriptionCreated, error) {
	event := new(VRFCoordinatorTestV2SubscriptionCreated)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2SubscriptionFundedIterator struct {
	Event *VRFCoordinatorTestV2SubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2SubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2SubscriptionFunded)
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
		it.Event = new(VRFCoordinatorTestV2SubscriptionFunded)
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

func (it *VRFCoordinatorTestV2SubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2SubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2SubscriptionFunded struct {
	SubId      uint64
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2SubscriptionFundedIterator{contract: _VRFCoordinatorTestV2.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionFunded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2SubscriptionFunded)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorTestV2SubscriptionFunded, error) {
	event := new(VRFCoordinatorTestV2SubscriptionFunded)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2SubscriptionOwnerTransferRequestedIterator struct {
	Event *VRFCoordinatorTestV2SubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2SubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2SubscriptionOwnerTransferRequested)
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
		it.Event = new(VRFCoordinatorTestV2SubscriptionOwnerTransferRequested)
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

func (it *VRFCoordinatorTestV2SubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2SubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2SubscriptionOwnerTransferRequested struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2SubscriptionOwnerTransferRequestedIterator{contract: _VRFCoordinatorTestV2.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2SubscriptionOwnerTransferRequested)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorTestV2SubscriptionOwnerTransferRequested, error) {
	event := new(VRFCoordinatorTestV2SubscriptionOwnerTransferRequested)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorTestV2SubscriptionOwnerTransferredIterator struct {
	Event *VRFCoordinatorTestV2SubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorTestV2SubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorTestV2SubscriptionOwnerTransferred)
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
		it.Event = new(VRFCoordinatorTestV2SubscriptionOwnerTransferred)
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

func (it *VRFCoordinatorTestV2SubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorTestV2SubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorTestV2SubscriptionOwnerTransferred struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorTestV2SubscriptionOwnerTransferredIterator{contract: _VRFCoordinatorTestV2.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorTestV2.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorTestV2SubscriptionOwnerTransferred)
				if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2Filterer) ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorTestV2SubscriptionOwnerTransferred, error) {
	event := new(VRFCoordinatorTestV2SubscriptionOwnerTransferred)
	if err := _VRFCoordinatorTestV2.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinatorTestV2.abi.Events["ConfigSet"].ID:
		return _VRFCoordinatorTestV2.ParseConfigSet(log)
	case _VRFCoordinatorTestV2.abi.Events["FundsRecovered"].ID:
		return _VRFCoordinatorTestV2.ParseFundsRecovered(log)
	case _VRFCoordinatorTestV2.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFCoordinatorTestV2.ParseOwnershipTransferRequested(log)
	case _VRFCoordinatorTestV2.abi.Events["OwnershipTransferred"].ID:
		return _VRFCoordinatorTestV2.ParseOwnershipTransferred(log)
	case _VRFCoordinatorTestV2.abi.Events["ProvingKeyDeregistered"].ID:
		return _VRFCoordinatorTestV2.ParseProvingKeyDeregistered(log)
	case _VRFCoordinatorTestV2.abi.Events["ProvingKeyRegistered"].ID:
		return _VRFCoordinatorTestV2.ParseProvingKeyRegistered(log)
	case _VRFCoordinatorTestV2.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFCoordinatorTestV2.ParseRandomWordsFulfilled(log)
	case _VRFCoordinatorTestV2.abi.Events["RandomWordsRequested"].ID:
		return _VRFCoordinatorTestV2.ParseRandomWordsRequested(log)
	case _VRFCoordinatorTestV2.abi.Events["SubscriptionCanceled"].ID:
		return _VRFCoordinatorTestV2.ParseSubscriptionCanceled(log)
	case _VRFCoordinatorTestV2.abi.Events["SubscriptionConsumerAdded"].ID:
		return _VRFCoordinatorTestV2.ParseSubscriptionConsumerAdded(log)
	case _VRFCoordinatorTestV2.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _VRFCoordinatorTestV2.ParseSubscriptionConsumerRemoved(log)
	case _VRFCoordinatorTestV2.abi.Events["SubscriptionCreated"].ID:
		return _VRFCoordinatorTestV2.ParseSubscriptionCreated(log)
	case _VRFCoordinatorTestV2.abi.Events["SubscriptionFunded"].ID:
		return _VRFCoordinatorTestV2.ParseSubscriptionFunded(log)
	case _VRFCoordinatorTestV2.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _VRFCoordinatorTestV2.ParseSubscriptionOwnerTransferRequested(log)
	case _VRFCoordinatorTestV2.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _VRFCoordinatorTestV2.ParseSubscriptionOwnerTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFCoordinatorTestV2ConfigSet) Topic() common.Hash {
	return common.HexToHash("0xc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb2")
}

func (VRFCoordinatorTestV2FundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (VRFCoordinatorTestV2OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFCoordinatorTestV2OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFCoordinatorTestV2ProvingKeyDeregistered) Topic() common.Hash {
	return common.HexToHash("0x72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d")
}

func (VRFCoordinatorTestV2ProvingKeyRegistered) Topic() common.Hash {
	return common.HexToHash("0xe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b8")
}

func (VRFCoordinatorTestV2RandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4")
}

func (VRFCoordinatorTestV2RandomWordsRequested) Topic() common.Hash {
	return common.HexToHash("0x63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772")
}

func (VRFCoordinatorTestV2SubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0xe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815")
}

func (VRFCoordinatorTestV2SubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0")
}

func (VRFCoordinatorTestV2SubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b")
}

func (VRFCoordinatorTestV2SubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf")
}

func (VRFCoordinatorTestV2SubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8")
}

func (VRFCoordinatorTestV2SubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be")
}

func (VRFCoordinatorTestV2SubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0")
}

func (_VRFCoordinatorTestV2 *VRFCoordinatorTestV2) Address() common.Address {
	return _VRFCoordinatorTestV2.address
}

type VRFCoordinatorTestV2Interface interface {
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

	FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFCoordinatorTestV2RequestCommitment) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64, newOwner common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorTestV2FeeConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorTestV2ConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2ConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFCoordinatorTestV2ConfigSet, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorTestV2FundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2FundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*VRFCoordinatorTestV2FundsRecovered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorTestV2OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorTestV2OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorTestV2OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorTestV2OwnershipTransferred, error)

	FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorTestV2ProvingKeyDeregisteredIterator, error)

	WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2ProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorTestV2ProvingKeyDeregistered, error)

	FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorTestV2ProvingKeyRegisteredIterator, error)

	WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2ProvingKeyRegistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorTestV2ProvingKeyRegistered, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*VRFCoordinatorTestV2RandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2RandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorTestV2RandomWordsFulfilled, error)

	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*VRFCoordinatorTestV2RandomWordsRequestedIterator, error)

	WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2RandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorTestV2RandomWordsRequested, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionCanceled, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorTestV2SubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionConsumerAdded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorTestV2SubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorTestV2SubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionCreated, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorTestV2SubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionFunded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorTestV2SubscriptionFunded, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorTestV2SubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorTestV2SubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorTestV2SubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorTestV2SubscriptionOwnerTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
