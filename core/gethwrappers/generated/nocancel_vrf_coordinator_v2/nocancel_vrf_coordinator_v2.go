// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nocancel_vrf_coordinator_v2

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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

type NoCancelVRFCoordinatorV2FeeConfig struct {
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

type NoCancelVRFCoordinatorV2RequestCommitment struct {
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

var NoCancelVRFCoordinatorV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"want\",\"type\":\"uint256\"}],\"name\":\"InsufficientGasForConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Reentrant\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structNoCancelVRFCoordinatorV2.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structNoCancelVRFCoordinatorV2.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"getCommitment\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentSubId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"}],\"name\":\"getFeeTier\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"internalType\":\"structNoCancelVRFCoordinatorV2.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162005ea138038062005ea18339810160408190526200003491620001b1565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000e8565b5050506001600160601b0319606093841b811660805290831b811660a052911b1660c052620001fb565b6001600160a01b038116331415620001435760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001ac57600080fd5b919050565b600080600060608486031215620001c757600080fd5b620001d28462000194565b9250620001e26020850162000194565b9150620001f26040850162000194565b90509250925092565b60805160601c60a05160601c60c05160601c615c436200025e600039600081816105260152613bde0152600061061d01526000818161036d01528181611586015281816125830152818161303101528181613187015261383b0152615c436000f3fe608060405234801561001057600080fd5b506004361061025b5760003560e01c80636f64f03f11610145578063ad178361116100bd578063d2f9f9a71161008c578063e72f6e3011610071578063e72f6e30146106fa578063e82ad7d41461070d578063f2fde38b1461073057600080fd5b8063d2f9f9a7146106d4578063d7ae1d30146106e757600080fd5b8063ad17836114610618578063af198b971461063f578063c3f909d41461066f578063caf70c4a146106c157600080fd5b80638da5cb5b11610114578063a21a23e4116100f9578063a21a23e4146105da578063a47c7696146105e2578063a4c0ed361461060557600080fd5b80638da5cb5b146105a95780639f87fad7146105c757600080fd5b80636f64f03f146105685780637341c10c1461057b57806379ba50971461058e578063823597401461059657600080fd5b8063356dac71116101d85780635fbbc0d2116101a757806366316d8d1161018c57806366316d8d1461050e578063689c45171461052157806369bcdb7d1461054857600080fd5b80635fbbc0d21461040057806364d51a2a1461050657600080fd5b8063356dac71146103b457806340d6bb82146103bc5780634cb48a54146103da5780635d3b1d30146103ed57600080fd5b806308821d581161022f57806315c48b841161021457806315c48b841461030e578063181f5a77146103295780631b6b6d231461036857600080fd5b806308821d58146102cf57806312b58349146102e257600080fd5b80620122911461026057806302bcc5b61461028057806304c357cb1461029557806306bfa637146102a8575b600080fd5b610268610743565b6040516102779392919061576d565b60405180910390f35b61029361028e3660046155df565b6107bf565b005b6102936102a33660046155fa565b610858565b60055467ffffffffffffffff165b60405167ffffffffffffffff9091168152602001610277565b6102936102dd3660046152f0565b610a4d565b6005546801000000000000000090046bffffffffffffffffffffffff165b604051908152602001610277565b61031660c881565b60405161ffff9091168152602001610277565b604080518082018252601e81527f4e6f43616e63656c565246436f6f7264696e61746f72563220312e302e3000006020820152905161027791906156fa565b61038f7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610277565b600a54610300565b6103c56101f481565b60405163ffffffff9091168152602001610277565b6102936103e8366004615489565b610c2c565b6103006103fb366004615363565b611023565b600c546040805163ffffffff80841682526401000000008404811660208301526801000000000000000084048116928201929092526c010000000000000000000000008304821660608201527001000000000000000000000000000000008304909116608082015262ffffff740100000000000000000000000000000000000000008304811660a0830152770100000000000000000000000000000000000000000000008304811660c08301527a0100000000000000000000000000000000000000000000000000008304811660e08301527d01000000000000000000000000000000000000000000000000000000000090920490911661010082015261012001610277565b610316606481565b61029361051c3660046152a8565b611431565b61038f7f000000000000000000000000000000000000000000000000000000000000000081565b6103006105563660046155c6565b60009081526009602052604090205490565b6102936105763660046151ed565b61169a565b6102936105893660046155fa565b6117e4565b610293611a72565b6102936105a43660046155df565b611b6f565b60005473ffffffffffffffffffffffffffffffffffffffff1661038f565b6102936105d53660046155fa565b611d69565b6102b66121ea565b6105f56105f03660046155df565b6123da565b604051610277949392919061590b565b610293610613366004615221565b612524565b61038f7f000000000000000000000000000000000000000000000000000000000000000081565b61065261064d3660046153c1565b612795565b6040516bffffffffffffffffffffffff9091168152602001610277565b600b546040805161ffff8316815263ffffffff6201000084048116602083015267010000000000000084048116928201929092526b010000000000000000000000909204166060820152608001610277565b6103006106cf36600461530c565b612c5a565b6103c56106e23660046155df565b612c8a565b6102936106f53660046155fa565b612e7f565b6102936107083660046151d2565b612ff8565b61072061071b3660046155df565b61325c565b6040519015158152602001610277565b61029361073e3660046151d2565b6134b3565b600b546007805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff169391928391908301828280156107ad57602002820191906000526020600020905b815481526020019060010190808311610799575b50505050509050925092509250909192565b6107c76134c4565b67ffffffffffffffff811660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1661082d576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108558161085060005473ffffffffffffffffffffffffffffffffffffffff1690565b613547565b50565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff16806108c1576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff82161461092d576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024015b60405180910390fd5b600b546601000000000000900460ff1615610974576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff841660009081526003602052604090206001015473ffffffffffffffffffffffffffffffffffffffff848116911614610a475767ffffffffffffffff841660008181526003602090815260409182902060010180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88169081179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25b50505050565b610a556134c4565b604080518082018252600091610a84919084906002908390839080828437600092019190915250612c5a915050565b60008181526006602052604090205490915073ffffffffffffffffffffffffffffffffffffffff1680610ae6576040517f77f5b84c00000000000000000000000000000000000000000000000000000000815260048101839052602401610924565b600082815260066020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690555b600754811015610bd6578260078281548110610b3957610b39615bd8565b90600052602060002001541415610bc4576007805460009190610b5e90600190615a6b565b81548110610b6e57610b6e615bd8565b906000526020600020015490508060078381548110610b8f57610b8f615bd8565b6000918252602090912001556007805480610bac57610bac615ba9565b60019003818190600052602060002001600090559055505b80610bce81615aaf565b915050610b1b565b508073ffffffffffffffffffffffffffffffffffffffff167f72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d83604051610c1f91815260200190565b60405180910390a2505050565b610c346134c4565b60c861ffff87161115610c87576040517fa738697600000000000000000000000000000000000000000000000000000000815261ffff871660048201819052602482015260c86044820152606401610924565b60008213610cc4576040517f43d4cf6600000000000000000000000000000000000000000000000000000000815260048101839052602401610924565b6040805160a0808201835261ffff891680835263ffffffff89811660208086018290526000868801528a831660608088018290528b85166080988901819052600b80547fffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000001690971762010000909502949094177fffffffffffffffffffffffffffffffffff000000000000000000ffffffffffff166701000000000000009092027fffffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffff16919091176b010000000000000000000000909302929092179093558651600c80549489015189890151938a0151978a0151968a015160c08b015160e08c01516101008d01519588167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009099169890981764010000000093881693909302929092177fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff1668010000000000000000958716959095027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff16949094176c0100000000000000000000000098861698909802979097177fffffffffffffffffff00000000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000096909416959095027fffffffffffffffffff000000ffffffffffffffffffffffffffffffffffffffff16929092177401000000000000000000000000000000000000000062ffffff92831602177fffffff000000000000ffffffffffffffffffffffffffffffffffffffffffffff1677010000000000000000000000000000000000000000000000958216959095027fffffff000000ffffffffffffffffffffffffffffffffffffffffffffffffffff16949094177a01000000000000000000000000000000000000000000000000000092851692909202919091177cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167d0100000000000000000000000000000000000000000000000000000000009390911692909202919091178155600a84905590517fc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb2916110139189918991899189918991906157cc565b60405180910390a1505050505050565b600b546000906601000000000000900460ff161561106d576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff851660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff166110d3576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b33600090815260026020908152604080832067ffffffffffffffff808a1685529252909120541680611143576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff87166004820152336024820152604401610924565b600b5461ffff908116908616108061115f575060c861ffff8616115b156111af57600b546040517fa738697600000000000000000000000000000000000000000000000000000000815261ffff8088166004830152909116602482015260c86044820152606401610924565b600b5463ffffffff620100009091048116908516111561121657600b546040517ff5d7e01e00000000000000000000000000000000000000000000000000000000815263ffffffff8087166004830152620100009092049091166024820152604401610924565b6101f463ffffffff84161115611268576040517f47386bec00000000000000000000000000000000000000000000000000000000815263ffffffff841660048201526101f46024820152604401610924565b60006112758260016159db565b6040805160208082018c9052338284015267ffffffffffffffff808c16606084015284166080808401919091528351808403909101815260a08301845280519082012060c083018d905260e080840182905284518085039091018152610100909301909352815191012091925060009182916040805160208101849052439181019190915267ffffffffffffffff8c16606082015263ffffffff808b166080830152891660a08201523360c0820152919350915060e001604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152828252805160209182012060008681526009835283902055848352820183905261ffff8a169082015263ffffffff808916606083015287166080820152339067ffffffffffffffff8b16908c907f63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a97729060a00160405180910390a45033600090815260026020908152604080832067ffffffffffffffff808d16855292529091208054919093167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009091161790915591505095945050505050565b600b546601000000000000900460ff1615611478576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b336000908152600860205260409020546bffffffffffffffffffffffff808316911610156114d2576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b33600090815260086020526040812080548392906114ff9084906bffffffffffffffffffffffff16615a82565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555080600560088282829054906101000a90046bffffffffffffffffffffffff166115569190615a82565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb83836040518363ffffffff1660e01b815260040161160e92919073ffffffffffffffffffffffffffffffffffffffff9290921682526bffffffffffffffffffffffff16602082015260400190565b602060405180830381600087803b15801561162857600080fd5b505af115801561163c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116609190615328565b611696576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b6116a26134c4565b6040805180820182526000916116d1919084906002908390839080828437600092019190915250612c5a915050565b60008181526006602052604090205490915073ffffffffffffffffffffffffffffffffffffffff1615611733576040517f4a0b8fa700000000000000000000000000000000000000000000000000000000815260048101829052602401610924565b600081815260066020908152604080832080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff88169081179091556007805460018101825594527fa66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a8736c688909301849055518381527fe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b89101610c1f565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff168061184d576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff8216146118b4576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610924565b600b546601000000000000900460ff16156118fb576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff841660009081526003602052604090206002015460641415611952576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8316600090815260026020908152604080832067ffffffffffffffff8089168552925290912054161561199957610a47565b73ffffffffffffffffffffffffffffffffffffffff8316600081815260026020818152604080842067ffffffffffffffff8a1680865290835281852080547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000166001908117909155600384528286209094018054948501815585529382902090920180547fffffffffffffffffffffffff00000000000000000000000000000000000000001685179055905192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e09101610a3e565b60015473ffffffffffffffffffffffffffffffffffffffff163314611af3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610924565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600b546601000000000000900460ff1615611bb6576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff16611c1c576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff811660009081526003602052604090206001015473ffffffffffffffffffffffffffffffffffffffff163314611cbe5767ffffffffffffffff8116600090815260036020526040908190206001015490517fd084e97500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610924565b67ffffffffffffffff81166000818152600360209081526040918290208054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560019093018054909316909255835173ffffffffffffffffffffffffffffffffffffffff909116808252928101919091529092917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0910160405180910390a25050565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680611dd2576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614611e39576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610924565b600b546601000000000000900460ff1615611e80576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8316600090815260026020908152604080832067ffffffffffffffff808916855292529091205416611f1b576040517ff0019fe600000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8516600482015273ffffffffffffffffffffffffffffffffffffffff84166024820152604401610924565b67ffffffffffffffff8416600090815260036020908152604080832060020180548251818502810185019093528083529192909190830182828015611f9657602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff168152600190910190602001808311611f6b575b50505050509050600060018251611fad9190615a6b565b905060005b825181101561214c578573ffffffffffffffffffffffffffffffffffffffff16838281518110611fe457611fe4615bd8565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16141561213a57600083838151811061201c5761201c615bd8565b6020026020010151905080600360008a67ffffffffffffffff1667ffffffffffffffff168152602001908152602001600020600201838154811061206257612062615bd8565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff949094169390931790925567ffffffffffffffff8a1681526003909152604090206002018054806120dc576120dc615ba9565b60008281526020902081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff00000000000000000000000000000000000000001690550190555061214c565b8061214481615aaf565b915050611fb2565b5073ffffffffffffffffffffffffffffffffffffffff8516600081815260026020908152604080832067ffffffffffffffff8b168085529083529281902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b91015b60405180910390a2505050505050565b600b546000906601000000000000900460ff1615612234576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005805467ffffffffffffffff1690600061224e83615ae8565b82546101009290920a67ffffffffffffffff8181021990931691831602179091556005541690506000806040519080825280602002602001820160405280156122a1578160200160208202803683370190505b506040805180820182526000808252602080830182815267ffffffffffffffff888116808552600484528685209551865493516bffffffffffffffffffffffff9091167fffffffffffffffffffffffff0000000000000000000000000000000000000000948516176c010000000000000000000000009190931602919091179094558451606081018652338152808301848152818701888152958552600384529590932083518154831673ffffffffffffffffffffffffffffffffffffffff918216178255955160018201805490931696169590951790559151805194955090936123929260028501920190614f2c565b505060405133815267ffffffffffffffff841691507f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a250905090565b67ffffffffffffffff81166000908152600360205260408120548190819060609073ffffffffffffffffffffffffffffffffffffffff16612447576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff80861660009081526004602090815260408083205460038352928190208054600290910180548351818602810186019094528084526bffffffffffffffffffffffff8616966c010000000000000000000000009096049095169473ffffffffffffffffffffffffffffffffffffffff90921693909291839183018282801561250e57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116124e3575b5050505050905093509350935093509193509193565b600b546601000000000000900460ff161561256b576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016146125da576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114612614576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612622828401846155df565b67ffffffffffffffff811660009081526003602052604090205490915073ffffffffffffffffffffffffffffffffffffffff1661268b576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8116600090815260046020526040812080546bffffffffffffffffffffffff16918691906126c28385615a07565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff16021790555084600560088282829054906101000a90046bffffffffffffffffffffffff166127199190615a07565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055508167ffffffffffffffff167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f882878461278091906159c3565b604080519283526020830191909152016121da565b600b546000906601000000000000900460ff16156127df576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005a905060008060006127f387876139ba565b9250925092506000866060015163ffffffff1667ffffffffffffffff81111561281e5761281e615c07565b604051908082528060200260200182016040528015612847578160200160208202803683370190505b50905060005b876060015163ffffffff168110156128bb5760408051602081018590529081018290526060016040516020818303038152906040528051906020012060001c82828151811061289e5761289e615bd8565b6020908102919091010152806128b381615aaf565b91505061284d565b506000838152600960205260408082208290555181907f1fe543e3000000000000000000000000000000000000000000000000000000009061290390879086906024016158bd565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090941693909317909252600b80547fffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffff166601000000000000179055908a015160808b01519192506000916129d19163ffffffff169084613d09565b600b80547fffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffff1690556020808c01805167ffffffffffffffff9081166000908152600490935260408084205492518216845290922080549394506c01000000000000000000000000918290048316936001939192600c92612a559286929004166159db565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506000612aac8a600b600001600b9054906101000a900463ffffffff1663ffffffff16612aa685612c8a565b3a613d57565b6020808e015167ffffffffffffffff166000908152600490915260409020549091506bffffffffffffffffffffffff80831691161015612b18576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6020808d015167ffffffffffffffff1660009081526004909152604081208054839290612b549084906bffffffffffffffffffffffff16615a82565b82546101009290920a6bffffffffffffffffffffffff81810219909316918316021790915560008b81526006602090815260408083205473ffffffffffffffffffffffffffffffffffffffff1683526008909152812080548594509092612bbd91859116615a07565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff160217905550877f7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4888386604051612c40939291909283526bffffffffffffffffffffffff9190911660208301521515604082015260600190565b60405180910390a299505050505050505050505b92915050565b600081604051602001612c6d91906156ec565b604051602081830303815290604052805190602001209050919050565b6040805161012081018252600c5463ffffffff80821683526401000000008204811660208401526801000000000000000082048116938301939093526c010000000000000000000000008104831660608301527001000000000000000000000000000000008104909216608082015262ffffff740100000000000000000000000000000000000000008304811660a08301819052770100000000000000000000000000000000000000000000008404821660c08401527a0100000000000000000000000000000000000000000000000000008404821660e08401527d0100000000000000000000000000000000000000000000000000000000009093041661010082015260009167ffffffffffffffff841611612da8575192915050565b8267ffffffffffffffff168160a0015162ffffff16108015612ddd57508060c0015162ffffff168367ffffffffffffffff1611155b15612dec576020015192915050565b8267ffffffffffffffff168160c0015162ffffff16108015612e2157508060e0015162ffffff168367ffffffffffffffff1611155b15612e30576040015192915050565b8267ffffffffffffffff168160e0015162ffffff16108015612e66575080610100015162ffffff168367ffffffffffffffff1611155b15612e75576060015192915050565b6080015192915050565b67ffffffffffffffff8216600090815260036020526040902054829073ffffffffffffffffffffffffffffffffffffffff1680612ee8576040517f1f6a65b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff821614612f4f576040517fd8a3fb5200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610924565b600b546601000000000000900460ff1615612f96576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7375622063616e63656c6c6174696f6e206e6f7420616c6c6f776564000000006044820152606401610924565b6130006134c4565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a082319060240160206040518083038186803b15801561308857600080fd5b505afa15801561309c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130c0919061534a565b6005549091506801000000000000000090046bffffffffffffffffffffffff1681811115613124576040517fa99da3020000000000000000000000000000000000000000000000000000000081526004810182905260248101839052604401610924565b818110156132575760006131388284615a6b565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152602482018390529192507f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb90604401602060405180830381600087803b1580156131cd57600080fd5b505af11580156131e1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906132059190615328565b506040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a1505b505050565b67ffffffffffffffff811660009081526003602090815260408083208151606081018352815473ffffffffffffffffffffffffffffffffffffffff9081168252600183015416818501526002820180548451818702810187018652818152879693958601939092919083018282801561330b57602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116132e0575b505050505081525050905060005b8160400151518110156134a95760005b60075481101561349657600061345f6007838154811061334b5761334b615bd8565b90600052602060002001548560400151858151811061336c5761336c615bd8565b602002602001015188600260008960400151898151811061338f5761338f615bd8565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff808f168352935220541660408051602080820187905273ffffffffffffffffffffffffffffffffffffffff959095168183015267ffffffffffffffff9384166060820152919092166080808301919091528251808303909101815260a08201835280519084012060c082019490945260e080820185905282518083039091018152610100909101909152805191012091565b50600081815260096020526040902054909150156134835750600195945050505050565b508061348e81615aaf565b915050613329565b50806134a181615aaf565b915050613319565b5060009392505050565b6134bb6134c4565b61085581613dc9565b60005473ffffffffffffffffffffffffffffffffffffffff163314613545576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610924565b565b600b546601000000000000900460ff161561358e576040517fed3ba6a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff821660009081526003602090815260408083208151606081018352815473ffffffffffffffffffffffffffffffffffffffff90811682526001830154168185015260028201805484518187028101870186528181529295939486019383018282801561363957602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161360e575b5050509190925250505067ffffffffffffffff80851660009081526004602090815260408083208151808301909252546bffffffffffffffffffffffff81168083526c01000000000000000000000000909104909416918101919091529293505b8360400151518110156137405760026000856040015183815181106136c1576136c1615bd8565b60209081029190910181015173ffffffffffffffffffffffffffffffffffffffff168252818101929092526040908101600090812067ffffffffffffffff8a168252909252902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001690558061373881615aaf565b91505061369a565b5067ffffffffffffffff8516600090815260036020526040812080547fffffffffffffffffffffffff0000000000000000000000000000000000000000908116825560018201805490911690559061379b6002830182614fb6565b505067ffffffffffffffff8516600090815260046020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690556005805482919060089061380b9084906801000000000000000090046bffffffffffffffffffffffff16615a82565b92506101000a8154816bffffffffffffffffffffffff02191690836bffffffffffffffffffffffff1602179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb85836bffffffffffffffffffffffff166040518363ffffffff1660e01b81526004016138c392919073ffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b602060405180830381600087803b1580156138dd57600080fd5b505af11580156138f1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139159190615328565b61394b576040517ff4d678b800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805173ffffffffffffffffffffffffffffffffffffffff861681526bffffffffffffffffffffffff8316602082015267ffffffffffffffff8716917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815910160405180910390a25050505050565b60008060006139cc8560000151612c5a565b60008181526006602052604090205490935073ffffffffffffffffffffffffffffffffffffffff1680613a2e576040517f77f5b84c00000000000000000000000000000000000000000000000000000000815260048101859052602401610924565b6080860151604051613a4d918691602001918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291815281516020928301206000818152600990935291205490935080613aca576040517f3688124a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b85516020808801516040808a015160608b015160808c01519251613b43968b96909594910195865267ffffffffffffffff948516602087015292909316604085015263ffffffff908116606085015291909116608083015273ffffffffffffffffffffffffffffffffffffffff1660a082015260c00190565b604051602081830303815290604052805190602001208114613b91576040517fd529142c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b855167ffffffffffffffff164080613cb55786516040517fe9413d3800000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063e9413d389060240160206040518083038186803b158015613c3557600080fd5b505afa158015613c49573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613c6d919061534a565b905080613cb55786516040517f175dadad00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610924565b6000886080015182604051602001613cd7929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c9050613cfc8982613ebf565b9450505050509250925092565b60005a611388811015613d1b57600080fd5b611388810390508460408204820311613d3357600080fd5b50823b613d3f57600080fd5b60008083516020850160008789f190505b9392505050565b600080613d6f63ffffffff851664e8d4a51000615a2e565b9050613d87816b033b2e3c9fd0803ce8000000615a6b565b811115613dc0576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b95945050505050565b73ffffffffffffffffffffffffffffffffffffffff8116331415613e49576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610924565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000613ef38360000151846020015185604001518660600151868860a001518960c001518a60e001518b6101000151613f48565b60038360200151604051602001613f0b9291906158a9565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209392505050565b613f518961421f565b613fb7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e2063757276650000000000006044820152606401610924565b613fc08861421f565b614026576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f67616d6d61206973206e6f74206f6e20637572766500000000000000000000006044820152606401610924565b61402f8361421f565b614095576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610924565b61409e8261421f565b614104576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e206375727665000000006044820152606401610924565b614110878a888761437a565b614176576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e657373000000000000006044820152606401610924565b60006141828a8761451d565b90506000614195898b878b868989614581565b905060006141a6838d8d8a86614715565b9050808a14614211576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f696e76616c69642070726f6f66000000000000000000000000000000000000006044820152606401610924565b505050505050505050505050565b80516000907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f116142ac576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420782d6f7264696e61746500000000000000000000000000006044820152606401610924565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f11614339576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f696e76616c696420792d6f7264696e61746500000000000000000000000000006044820152606401610924565b60208201517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f9080096143738360005b6020020151614773565b1492915050565b600073ffffffffffffffffffffffffffffffffffffffff82166143f9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f626164207769746e6573730000000000000000000000000000000000000000006044820152606401610924565b60208401516000906001161561441057601c614413565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418587600060200201510986517ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa1580156144ca573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff9081169088161495505050505050949350505050565b614525614fd4565b6145526001848460405160200161453e939291906156cb565b6040516020818303038152906040526147cb565b90505b61455e8161421f565b612c5457805160408051602081019290925261457a910161453e565b9050614555565b614589614fd4565b825186517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f908190069106141561461c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610924565b614627878988614834565b61468d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4669727374206d756c20636865636b206661696c6564000000000000000000006044820152606401610924565b614698848685614834565b6146fe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c65640000000000000000006044820152606401610924565b6147098684846149c1565b98975050505050505050565b60006002868686858760405160200161473396959493929190615659565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209695505050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80848509840990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f600782089392505050565b6147d3614fd4565b6147dc82614af0565b81526147f16147ec826000614369565b614b45565b60208201819052600290066001141561482f576020810180517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f0390525b919050565b60008261489d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f7a65726f207363616c61720000000000000000000000000000000000000000006044820152606401610924565b835160208501516000906148b390600290615b10565b156148bf57601c6148c2565b601b5b905060007ffffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd03641418387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa158015614942573d6000803e3d6000fd5b5050506020604051035190506000866040516020016149619190615647565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052805160209091012073ffffffffffffffffffffffffffffffffffffffff92831692169190911498975050505050505050565b6149c9614fd4565b8351602080860151855191860151600093849384936149ea93909190614b7f565b919450925090507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f858209600114614a7e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a000000000000006044820152606401610924565b60405180604001604052807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f80614ab757614ab7615b7a565b87860981526020017ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8785099052979650505050505050565b805160208201205b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f811061482f57604080516020808201939093528151808203840181529082019091528051910120614af8565b6000612c54826002614b787ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f60016159c3565b901c614d15565b60008080600180827ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f897ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038808905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f038a0890506000614c2783838585614e09565b9098509050614c3888828e88614e61565b9098509050614c4988828c87614e61565b90985090506000614c5c8d878b85614e61565b9098509050614c6d88828686614e09565b9098509050614c7e88828e89614e61565b9098509050818114614d01577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f818a0998507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f82890997507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183099650614d05565b8196505b5050505050509450945094915050565b600080614d20614ff2565b6020808252818101819052604082015260608101859052608081018490527ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f60a0820152614d6c615010565b60208160c08460057ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa925082614dff576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f6269674d6f64457870206661696c7572652100000000000000000000000000006044820152606401610924565b5195945050505050565b6000807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487097ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8487099097909650945050505050565b600080807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f878509905060007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f87877ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f030990507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8183087ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f86890990999098509650505050505050565b828054828255906000526020600020908101928215614fa6579160200282015b82811115614fa657825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190614f4c565b50614fb292915061502e565b5090565b5080546000825590600052602060002090810190610855919061502e565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b80821115614fb2576000815560010161502f565b803573ffffffffffffffffffffffffffffffffffffffff8116811461482f57600080fd5b8060408101831015612c5457600080fd5b600082601f83011261508957600080fd5b6040516040810181811067ffffffffffffffff821117156150ac576150ac615c07565b80604052508083856040860111156150c357600080fd5b60005b60028110156150e55781358352602092830192909101906001016150c6565b509195945050505050565b600060a0828403121561510257600080fd5b60405160a0810181811067ffffffffffffffff8211171561512557615125615c07565b604052905080615134836151ba565b8152615142602084016151ba565b6020820152615153604084016151a6565b6040820152615164606084016151a6565b606082015261517560808401615043565b60808201525092915050565b803561ffff8116811461482f57600080fd5b803562ffffff8116811461482f57600080fd5b803563ffffffff8116811461482f57600080fd5b803567ffffffffffffffff8116811461482f57600080fd5b6000602082840312156151e457600080fd5b613d5082615043565b6000806060838503121561520057600080fd5b61520983615043565b91506152188460208501615067565b90509250929050565b6000806000806060858703121561523757600080fd5b61524085615043565b935060208501359250604085013567ffffffffffffffff8082111561526457600080fd5b818701915087601f83011261527857600080fd5b81358181111561528757600080fd5b88602082850101111561529957600080fd5b95989497505060200194505050565b600080604083850312156152bb57600080fd5b6152c483615043565b915060208301356bffffffffffffffffffffffff811681146152e557600080fd5b809150509250929050565b60006040828403121561530257600080fd5b613d508383615067565b60006040828403121561531e57600080fd5b613d508383615078565b60006020828403121561533a57600080fd5b81518015158114613d5057600080fd5b60006020828403121561535c57600080fd5b5051919050565b600080600080600060a0868803121561537b57600080fd5b8535945061538b602087016151ba565b935061539960408701615181565b92506153a7606087016151a6565b91506153b5608087016151a6565b90509295509295909350565b6000808284036102408112156153d657600080fd5b6101a0808212156153e657600080fd5b6153ee615999565b91506153fa8686615078565b82526154098660408701615078565b60208301526080850135604083015260a0850135606083015260c0850135608083015261543860e08601615043565b60a083015261010061544c87828801615078565b60c084015261545f876101408801615078565b60e0840152610180860135818401525081935061547e868287016150f0565b925050509250929050565b6000806000806000808688036101c08112156154a457600080fd5b6154ad88615181565b96506154bb602089016151a6565b95506154c9604089016151a6565b94506154d7606089016151a6565b935060808801359250610120807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff608301121561551257600080fd5b61551a615999565b915061552860a08a016151a6565b825261553660c08a016151a6565b602083015261554760e08a016151a6565b604083015261010061555a818b016151a6565b606084015261556a828b016151a6565b608084015261557c6101408b01615193565b60a084015261558e6101608b01615193565b60c08401526155a06101808b01615193565b60e08401526155b26101a08b01615193565b818401525050809150509295509295509295565b6000602082840312156155d857600080fd5b5035919050565b6000602082840312156155f157600080fd5b613d50826151ba565b6000806040838503121561560d57600080fd5b615616836151ba565b915061521860208401615043565b8060005b6002811015610a47578151845260209384019390910190600101615628565b6156518183615624565b604001919050565b8681526156696020820187615624565b6156766060820186615624565b61568360a0820185615624565b61569060e0820184615624565b60609190911b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166101208201526101340195945050505050565b8381526156db6020820184615624565b606081019190915260800192915050565b60408101612c548284615624565b600060208083528351808285015260005b818110156157275785810183015185820160400152820161570b565b81811115615739576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b818110156157be578451835293830193918301916001016157a2565b509098975050505050505050565b60006101c08201905061ffff8816825263ffffffff808816602084015280871660408401528086166060840152846080840152835481811660a085015261582060c08501838360201c1663ffffffff169052565b61583760e08501838360401c1663ffffffff169052565b61584f6101008501838360601c1663ffffffff169052565b6158676101208501838360801c1663ffffffff169052565b62ffffff60a082901c811661014086015260b882901c811661016086015260d082901c1661018085015260e81c6101a090930192909252979650505050505050565b82815260608101613d506020830184615624565b6000604082018483526020604081850152818551808452606086019150828701935060005b818110156158fe578451835293830193918301916001016158e2565b5090979650505050505050565b6000608082016bffffffffffffffffffffffff87168352602067ffffffffffffffff87168185015273ffffffffffffffffffffffffffffffffffffffff80871660408601526080606086015282865180855260a087019150838801945060005b8181101561598957855184168352948401949184019160010161596b565b50909a9950505050505050505050565b604051610120810167ffffffffffffffff811182821017156159bd576159bd615c07565b60405290565b600082198211156159d6576159d6615b4b565b500190565b600067ffffffffffffffff8083168185168083038211156159fe576159fe615b4b565b01949350505050565b60006bffffffffffffffffffffffff8083168185168083038211156159fe576159fe615b4b565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615615a6657615a66615b4b565b500290565b600082821015615a7d57615a7d615b4b565b500390565b60006bffffffffffffffffffffffff83811690831681811015615aa757615aa7615b4b565b039392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415615ae157615ae1615b4b565b5060010190565b600067ffffffffffffffff80831681811415615b0657615b06615b4b565b6001019392505050565b600082615b46577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var NoCancelVRFCoordinatorV2ABI = NoCancelVRFCoordinatorV2MetaData.ABI

var NoCancelVRFCoordinatorV2Bin = NoCancelVRFCoordinatorV2MetaData.Bin

func DeployNoCancelVRFCoordinatorV2(auth *bind.TransactOpts, backend bind.ContractBackend, link common.Address, blockhashStore common.Address, linkEthFeed common.Address) (common.Address, *types.Transaction, *NoCancelVRFCoordinatorV2, error) {
	parsed, err := NoCancelVRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NoCancelVRFCoordinatorV2Bin), backend, link, blockhashStore, linkEthFeed)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NoCancelVRFCoordinatorV2{NoCancelVRFCoordinatorV2Caller: NoCancelVRFCoordinatorV2Caller{contract: contract}, NoCancelVRFCoordinatorV2Transactor: NoCancelVRFCoordinatorV2Transactor{contract: contract}, NoCancelVRFCoordinatorV2Filterer: NoCancelVRFCoordinatorV2Filterer{contract: contract}}, nil
}

type NoCancelVRFCoordinatorV2 struct {
	address common.Address
	abi     abi.ABI
	NoCancelVRFCoordinatorV2Caller
	NoCancelVRFCoordinatorV2Transactor
	NoCancelVRFCoordinatorV2Filterer
}

type NoCancelVRFCoordinatorV2Caller struct {
	contract *bind.BoundContract
}

type NoCancelVRFCoordinatorV2Transactor struct {
	contract *bind.BoundContract
}

type NoCancelVRFCoordinatorV2Filterer struct {
	contract *bind.BoundContract
}

type NoCancelVRFCoordinatorV2Session struct {
	Contract     *NoCancelVRFCoordinatorV2
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type NoCancelVRFCoordinatorV2CallerSession struct {
	Contract *NoCancelVRFCoordinatorV2Caller
	CallOpts bind.CallOpts
}

type NoCancelVRFCoordinatorV2TransactorSession struct {
	Contract     *NoCancelVRFCoordinatorV2Transactor
	TransactOpts bind.TransactOpts
}

type NoCancelVRFCoordinatorV2Raw struct {
	Contract *NoCancelVRFCoordinatorV2
}

type NoCancelVRFCoordinatorV2CallerRaw struct {
	Contract *NoCancelVRFCoordinatorV2Caller
}

type NoCancelVRFCoordinatorV2TransactorRaw struct {
	Contract *NoCancelVRFCoordinatorV2Transactor
}

func NewNoCancelVRFCoordinatorV2(address common.Address, backend bind.ContractBackend) (*NoCancelVRFCoordinatorV2, error) {
	abi, err := abi.JSON(strings.NewReader(NoCancelVRFCoordinatorV2ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindNoCancelVRFCoordinatorV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2{address: address, abi: abi, NoCancelVRFCoordinatorV2Caller: NoCancelVRFCoordinatorV2Caller{contract: contract}, NoCancelVRFCoordinatorV2Transactor: NoCancelVRFCoordinatorV2Transactor{contract: contract}, NoCancelVRFCoordinatorV2Filterer: NoCancelVRFCoordinatorV2Filterer{contract: contract}}, nil
}

func NewNoCancelVRFCoordinatorV2Caller(address common.Address, caller bind.ContractCaller) (*NoCancelVRFCoordinatorV2Caller, error) {
	contract, err := bindNoCancelVRFCoordinatorV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2Caller{contract: contract}, nil
}

func NewNoCancelVRFCoordinatorV2Transactor(address common.Address, transactor bind.ContractTransactor) (*NoCancelVRFCoordinatorV2Transactor, error) {
	contract, err := bindNoCancelVRFCoordinatorV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2Transactor{contract: contract}, nil
}

func NewNoCancelVRFCoordinatorV2Filterer(address common.Address, filterer bind.ContractFilterer) (*NoCancelVRFCoordinatorV2Filterer, error) {
	contract, err := bindNoCancelVRFCoordinatorV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2Filterer{contract: contract}, nil
}

func bindNoCancelVRFCoordinatorV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NoCancelVRFCoordinatorV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NoCancelVRFCoordinatorV2.Contract.NoCancelVRFCoordinatorV2Caller.contract.Call(opts, result, method, params...)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.NoCancelVRFCoordinatorV2Transactor.contract.Transfer(opts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.NoCancelVRFCoordinatorV2Transactor.contract.Transact(opts, method, params...)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NoCancelVRFCoordinatorV2.Contract.contract.Call(opts, result, method, params...)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.contract.Transfer(opts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.contract.Transact(opts, method, params...)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "BLOCKHASH_STORE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) BLOCKHASHSTORE() (common.Address, error) {
	return _NoCancelVRFCoordinatorV2.Contract.BLOCKHASHSTORE(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) BLOCKHASHSTORE() (common.Address, error) {
	return _NoCancelVRFCoordinatorV2.Contract.BLOCKHASHSTORE(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) LINK() (common.Address, error) {
	return _NoCancelVRFCoordinatorV2.Contract.LINK(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) LINK() (common.Address, error) {
	return _NoCancelVRFCoordinatorV2.Contract.LINK(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) LINKETHFEED() (common.Address, error) {
	return _NoCancelVRFCoordinatorV2.Contract.LINKETHFEED(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) LINKETHFEED() (common.Address, error) {
	return _NoCancelVRFCoordinatorV2.Contract.LINKETHFEED(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) MAXCONSUMERS() (uint16, error) {
	return _NoCancelVRFCoordinatorV2.Contract.MAXCONSUMERS(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) MAXCONSUMERS() (uint16, error) {
	return _NoCancelVRFCoordinatorV2.Contract.MAXCONSUMERS(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) MAXNUMWORDS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "MAX_NUM_WORDS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) MAXNUMWORDS() (uint32, error) {
	return _NoCancelVRFCoordinatorV2.Contract.MAXNUMWORDS(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) MAXNUMWORDS() (uint32, error) {
	return _NoCancelVRFCoordinatorV2.Contract.MAXNUMWORDS(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "MAX_REQUEST_CONFIRMATIONS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _NoCancelVRFCoordinatorV2.Contract.MAXREQUESTCONFIRMATIONS(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _NoCancelVRFCoordinatorV2.Contract.MAXREQUESTCONFIRMATIONS(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) GetCommitment(opts *bind.CallOpts, requestId *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "getCommitment", requestId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) GetCommitment(requestId *big.Int) ([32]byte, error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetCommitment(&_NoCancelVRFCoordinatorV2.CallOpts, requestId)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) GetCommitment(requestId *big.Int) ([32]byte, error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetCommitment(&_NoCancelVRFCoordinatorV2.CallOpts, requestId)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) GetConfig(opts *bind.CallOpts) (GetConfig,

	error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "getConfig")

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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) GetConfig() (GetConfig,

	error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetConfig(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) GetConfig() (GetConfig,

	error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetConfig(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) GetCurrentSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "getCurrentSubId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) GetCurrentSubId() (uint64, error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetCurrentSubId(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) GetCurrentSubId() (uint64, error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetCurrentSubId(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) GetFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "getFallbackWeiPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) GetFallbackWeiPerUnitLink() (*big.Int, error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetFallbackWeiPerUnitLink(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) GetFallbackWeiPerUnitLink() (*big.Int, error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetFallbackWeiPerUnitLink(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) GetFeeConfig(opts *bind.CallOpts) (GetFeeConfig,

	error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "getFeeConfig")

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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) GetFeeConfig() (GetFeeConfig,

	error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetFeeConfig(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) GetFeeConfig() (GetFeeConfig,

	error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetFeeConfig(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) GetFeeTier(opts *bind.CallOpts, reqCount uint64) (uint32, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "getFeeTier", reqCount)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) GetFeeTier(reqCount uint64) (uint32, error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetFeeTier(&_NoCancelVRFCoordinatorV2.CallOpts, reqCount)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) GetFeeTier(reqCount uint64) (uint32, error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetFeeTier(&_NoCancelVRFCoordinatorV2.CallOpts, reqCount)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) GetRequestConfig(opts *bind.CallOpts) (uint16, uint32, [][32]byte, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "getRequestConfig")

	if err != nil {
		return *new(uint16), *new(uint32), *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)
	out2 := *abi.ConvertType(out[2], new([][32]byte)).(*[][32]byte)

	return out0, out1, out2, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetRequestConfig(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetRequestConfig(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) GetSubscription(opts *bind.CallOpts, subId uint64) (GetSubscription,

	error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "getSubscription", subId)

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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) GetSubscription(subId uint64) (GetSubscription,

	error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetSubscription(&_NoCancelVRFCoordinatorV2.CallOpts, subId)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) GetSubscription(subId uint64) (GetSubscription,

	error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetSubscription(&_NoCancelVRFCoordinatorV2.CallOpts, subId)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) GetTotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "getTotalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) GetTotalBalance() (*big.Int, error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetTotalBalance(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) GetTotalBalance() (*big.Int, error) {
	return _NoCancelVRFCoordinatorV2.Contract.GetTotalBalance(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "hashOfKey", publicKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _NoCancelVRFCoordinatorV2.Contract.HashOfKey(&_NoCancelVRFCoordinatorV2.CallOpts, publicKey)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _NoCancelVRFCoordinatorV2.Contract.HashOfKey(&_NoCancelVRFCoordinatorV2.CallOpts, publicKey)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) Owner() (common.Address, error) {
	return _NoCancelVRFCoordinatorV2.Contract.Owner(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) Owner() (common.Address, error) {
	return _NoCancelVRFCoordinatorV2.Contract.Owner(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) PendingRequestExists(opts *bind.CallOpts, subId uint64) (bool, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "pendingRequestExists", subId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) PendingRequestExists(subId uint64) (bool, error) {
	return _NoCancelVRFCoordinatorV2.Contract.PendingRequestExists(&_NoCancelVRFCoordinatorV2.CallOpts, subId)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) PendingRequestExists(subId uint64) (bool, error) {
	return _NoCancelVRFCoordinatorV2.Contract.PendingRequestExists(&_NoCancelVRFCoordinatorV2.CallOpts, subId)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NoCancelVRFCoordinatorV2.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) TypeAndVersion() (string, error) {
	return _NoCancelVRFCoordinatorV2.Contract.TypeAndVersion(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2CallerSession) TypeAndVersion() (string, error) {
	return _NoCancelVRFCoordinatorV2.Contract.TypeAndVersion(&_NoCancelVRFCoordinatorV2.CallOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "acceptOwnership")
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) AcceptOwnership() (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.AcceptOwnership(&_NoCancelVRFCoordinatorV2.TransactOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.AcceptOwnership(&_NoCancelVRFCoordinatorV2.TransactOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.AcceptSubscriptionOwnerTransfer(&_NoCancelVRFCoordinatorV2.TransactOpts, subId)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.AcceptSubscriptionOwnerTransfer(&_NoCancelVRFCoordinatorV2.TransactOpts, subId)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) AddConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "addConsumer", subId, consumer)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.AddConsumer(&_NoCancelVRFCoordinatorV2.TransactOpts, subId, consumer)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.AddConsumer(&_NoCancelVRFCoordinatorV2.TransactOpts, subId, consumer)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) CancelSubscription(opts *bind.TransactOpts, subId uint64, to common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "cancelSubscription", subId, to)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) CancelSubscription(subId uint64, to common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.CancelSubscription(&_NoCancelVRFCoordinatorV2.TransactOpts, subId, to)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) CancelSubscription(subId uint64, to common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.CancelSubscription(&_NoCancelVRFCoordinatorV2.TransactOpts, subId, to)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "createSubscription")
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) CreateSubscription() (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.CreateSubscription(&_NoCancelVRFCoordinatorV2.TransactOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.CreateSubscription(&_NoCancelVRFCoordinatorV2.TransactOpts)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "deregisterProvingKey", publicProvingKey)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.DeregisterProvingKey(&_NoCancelVRFCoordinatorV2.TransactOpts, publicProvingKey)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.DeregisterProvingKey(&_NoCancelVRFCoordinatorV2.TransactOpts, publicProvingKey)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc NoCancelVRFCoordinatorV2RequestCommitment) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "fulfillRandomWords", proof, rc)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) FulfillRandomWords(proof VRFProof, rc NoCancelVRFCoordinatorV2RequestCommitment) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.FulfillRandomWords(&_NoCancelVRFCoordinatorV2.TransactOpts, proof, rc)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) FulfillRandomWords(proof VRFProof, rc NoCancelVRFCoordinatorV2RequestCommitment) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.FulfillRandomWords(&_NoCancelVRFCoordinatorV2.TransactOpts, proof, rc)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.OnTokenTransfer(&_NoCancelVRFCoordinatorV2.TransactOpts, arg0, amount, data)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.OnTokenTransfer(&_NoCancelVRFCoordinatorV2.TransactOpts, arg0, amount, data)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "oracleWithdraw", recipient, amount)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.OracleWithdraw(&_NoCancelVRFCoordinatorV2.TransactOpts, recipient, amount)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.OracleWithdraw(&_NoCancelVRFCoordinatorV2.TransactOpts, recipient, amount)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "ownerCancelSubscription", subId)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.OwnerCancelSubscription(&_NoCancelVRFCoordinatorV2.TransactOpts, subId)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.OwnerCancelSubscription(&_NoCancelVRFCoordinatorV2.TransactOpts, subId)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "recoverFunds", to)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.RecoverFunds(&_NoCancelVRFCoordinatorV2.TransactOpts, to)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.RecoverFunds(&_NoCancelVRFCoordinatorV2.TransactOpts, to)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "registerProvingKey", oracle, publicProvingKey)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.RegisterProvingKey(&_NoCancelVRFCoordinatorV2.TransactOpts, oracle, publicProvingKey)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.RegisterProvingKey(&_NoCancelVRFCoordinatorV2.TransactOpts, oracle, publicProvingKey)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) RemoveConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "removeConsumer", subId, consumer)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) RemoveConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.RemoveConsumer(&_NoCancelVRFCoordinatorV2.TransactOpts, subId, consumer)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) RemoveConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.RemoveConsumer(&_NoCancelVRFCoordinatorV2.TransactOpts, subId, consumer)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "requestRandomWords", keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) RequestRandomWords(keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.RequestRandomWords(&_NoCancelVRFCoordinatorV2.TransactOpts, keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) RequestRandomWords(keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.RequestRandomWords(&_NoCancelVRFCoordinatorV2.TransactOpts, keyHash, subId, requestConfirmations, callbackGasLimit, numWords)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.RequestSubscriptionOwnerTransfer(&_NoCancelVRFCoordinatorV2.TransactOpts, subId, newOwner)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.RequestSubscriptionOwnerTransfer(&_NoCancelVRFCoordinatorV2.TransactOpts, subId, newOwner)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig NoCancelVRFCoordinatorV2FeeConfig) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig NoCancelVRFCoordinatorV2FeeConfig) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.SetConfig(&_NoCancelVRFCoordinatorV2.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig NoCancelVRFCoordinatorV2FeeConfig) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.SetConfig(&_NoCancelVRFCoordinatorV2.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.contract.Transact(opts, "transferOwnership", to)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.TransferOwnership(&_NoCancelVRFCoordinatorV2.TransactOpts, to)
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _NoCancelVRFCoordinatorV2.Contract.TransferOwnership(&_NoCancelVRFCoordinatorV2.TransactOpts, to)
}

type NoCancelVRFCoordinatorV2ConfigSetIterator struct {
	Event *NoCancelVRFCoordinatorV2ConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2ConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2ConfigSet)
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
		it.Event = new(NoCancelVRFCoordinatorV2ConfigSet)
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

func (it *NoCancelVRFCoordinatorV2ConfigSetIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2ConfigSet struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
	FallbackWeiPerUnitLink      *big.Int
	FeeConfig                   NoCancelVRFCoordinatorV2FeeConfig
	Raw                         types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterConfigSet(opts *bind.FilterOpts) (*NoCancelVRFCoordinatorV2ConfigSetIterator, error) {

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2ConfigSetIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2ConfigSet) (event.Subscription, error) {

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2ConfigSet)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseConfigSet(log types.Log) (*NoCancelVRFCoordinatorV2ConfigSet, error) {
	event := new(NoCancelVRFCoordinatorV2ConfigSet)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2FundsRecoveredIterator struct {
	Event *NoCancelVRFCoordinatorV2FundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2FundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2FundsRecovered)
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
		it.Event = new(NoCancelVRFCoordinatorV2FundsRecovered)
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

func (it *NoCancelVRFCoordinatorV2FundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2FundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2FundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterFundsRecovered(opts *bind.FilterOpts) (*NoCancelVRFCoordinatorV2FundsRecoveredIterator, error) {

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2FundsRecoveredIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2FundsRecovered) (event.Subscription, error) {

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2FundsRecovered)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseFundsRecovered(log types.Log) (*NoCancelVRFCoordinatorV2FundsRecovered, error) {
	event := new(NoCancelVRFCoordinatorV2FundsRecovered)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2OwnershipTransferRequestedIterator struct {
	Event *NoCancelVRFCoordinatorV2OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2OwnershipTransferRequested)
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
		it.Event = new(NoCancelVRFCoordinatorV2OwnershipTransferRequested)
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

func (it *NoCancelVRFCoordinatorV2OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NoCancelVRFCoordinatorV2OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2OwnershipTransferRequestedIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2OwnershipTransferRequested)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseOwnershipTransferRequested(log types.Log) (*NoCancelVRFCoordinatorV2OwnershipTransferRequested, error) {
	event := new(NoCancelVRFCoordinatorV2OwnershipTransferRequested)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2OwnershipTransferredIterator struct {
	Event *NoCancelVRFCoordinatorV2OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2OwnershipTransferred)
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
		it.Event = new(NoCancelVRFCoordinatorV2OwnershipTransferred)
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

func (it *NoCancelVRFCoordinatorV2OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NoCancelVRFCoordinatorV2OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2OwnershipTransferredIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2OwnershipTransferred)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseOwnershipTransferred(log types.Log) (*NoCancelVRFCoordinatorV2OwnershipTransferred, error) {
	event := new(NoCancelVRFCoordinatorV2OwnershipTransferred)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2ProvingKeyDeregisteredIterator struct {
	Event *NoCancelVRFCoordinatorV2ProvingKeyDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2ProvingKeyDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2ProvingKeyDeregistered)
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
		it.Event = new(NoCancelVRFCoordinatorV2ProvingKeyDeregistered)
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

func (it *NoCancelVRFCoordinatorV2ProvingKeyDeregisteredIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2ProvingKeyDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2ProvingKeyDeregistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*NoCancelVRFCoordinatorV2ProvingKeyDeregisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2ProvingKeyDeregisteredIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "ProvingKeyDeregistered", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2ProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2ProvingKeyDeregistered)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseProvingKeyDeregistered(log types.Log) (*NoCancelVRFCoordinatorV2ProvingKeyDeregistered, error) {
	event := new(NoCancelVRFCoordinatorV2ProvingKeyDeregistered)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2ProvingKeyRegisteredIterator struct {
	Event *NoCancelVRFCoordinatorV2ProvingKeyRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2ProvingKeyRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2ProvingKeyRegistered)
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
		it.Event = new(NoCancelVRFCoordinatorV2ProvingKeyRegistered)
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

func (it *NoCancelVRFCoordinatorV2ProvingKeyRegisteredIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2ProvingKeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2ProvingKeyRegistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*NoCancelVRFCoordinatorV2ProvingKeyRegisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2ProvingKeyRegisteredIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "ProvingKeyRegistered", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2ProvingKeyRegistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2ProvingKeyRegistered)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseProvingKeyRegistered(log types.Log) (*NoCancelVRFCoordinatorV2ProvingKeyRegistered, error) {
	event := new(NoCancelVRFCoordinatorV2ProvingKeyRegistered)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2RandomWordsFulfilledIterator struct {
	Event *NoCancelVRFCoordinatorV2RandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2RandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2RandomWordsFulfilled)
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
		it.Event = new(NoCancelVRFCoordinatorV2RandomWordsFulfilled)
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

func (it *NoCancelVRFCoordinatorV2RandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2RandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2RandomWordsFulfilled struct {
	RequestId  *big.Int
	OutputSeed *big.Int
	Payment    *big.Int
	Success    bool
	Raw        types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*NoCancelVRFCoordinatorV2RandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2RandomWordsFulfilledIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2RandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2RandomWordsFulfilled)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseRandomWordsFulfilled(log types.Log) (*NoCancelVRFCoordinatorV2RandomWordsFulfilled, error) {
	event := new(NoCancelVRFCoordinatorV2RandomWordsFulfilled)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2RandomWordsRequestedIterator struct {
	Event *NoCancelVRFCoordinatorV2RandomWordsRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2RandomWordsRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2RandomWordsRequested)
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
		it.Event = new(NoCancelVRFCoordinatorV2RandomWordsRequested)
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

func (it *NoCancelVRFCoordinatorV2RandomWordsRequestedIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2RandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2RandomWordsRequested struct {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*NoCancelVRFCoordinatorV2RandomWordsRequestedIterator, error) {

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

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2RandomWordsRequestedIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2RandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2RandomWordsRequested)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseRandomWordsRequested(log types.Log) (*NoCancelVRFCoordinatorV2RandomWordsRequested, error) {
	event := new(NoCancelVRFCoordinatorV2RandomWordsRequested)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2SubscriptionCanceledIterator struct {
	Event *NoCancelVRFCoordinatorV2SubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2SubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2SubscriptionCanceled)
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
		it.Event = new(NoCancelVRFCoordinatorV2SubscriptionCanceled)
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

func (it *NoCancelVRFCoordinatorV2SubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2SubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2SubscriptionCanceled struct {
	SubId  uint64
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2SubscriptionCanceledIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionCanceled, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2SubscriptionCanceled)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseSubscriptionCanceled(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionCanceled, error) {
	event := new(NoCancelVRFCoordinatorV2SubscriptionCanceled)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2SubscriptionConsumerAddedIterator struct {
	Event *NoCancelVRFCoordinatorV2SubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2SubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2SubscriptionConsumerAdded)
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
		it.Event = new(NoCancelVRFCoordinatorV2SubscriptionConsumerAdded)
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

func (it *NoCancelVRFCoordinatorV2SubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2SubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2SubscriptionConsumerAdded struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2SubscriptionConsumerAddedIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionConsumerAdded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2SubscriptionConsumerAdded)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseSubscriptionConsumerAdded(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionConsumerAdded, error) {
	event := new(NoCancelVRFCoordinatorV2SubscriptionConsumerAdded)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2SubscriptionConsumerRemovedIterator struct {
	Event *NoCancelVRFCoordinatorV2SubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2SubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2SubscriptionConsumerRemoved)
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
		it.Event = new(NoCancelVRFCoordinatorV2SubscriptionConsumerRemoved)
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

func (it *NoCancelVRFCoordinatorV2SubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2SubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2SubscriptionConsumerRemoved struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2SubscriptionConsumerRemovedIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2SubscriptionConsumerRemoved)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseSubscriptionConsumerRemoved(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionConsumerRemoved, error) {
	event := new(NoCancelVRFCoordinatorV2SubscriptionConsumerRemoved)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2SubscriptionCreatedIterator struct {
	Event *NoCancelVRFCoordinatorV2SubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2SubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2SubscriptionCreated)
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
		it.Event = new(NoCancelVRFCoordinatorV2SubscriptionCreated)
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

func (it *NoCancelVRFCoordinatorV2SubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2SubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2SubscriptionCreated struct {
	SubId uint64
	Owner common.Address
	Raw   types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2SubscriptionCreatedIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionCreated, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2SubscriptionCreated)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseSubscriptionCreated(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionCreated, error) {
	event := new(NoCancelVRFCoordinatorV2SubscriptionCreated)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2SubscriptionFundedIterator struct {
	Event *NoCancelVRFCoordinatorV2SubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2SubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2SubscriptionFunded)
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
		it.Event = new(NoCancelVRFCoordinatorV2SubscriptionFunded)
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

func (it *NoCancelVRFCoordinatorV2SubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2SubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2SubscriptionFunded struct {
	SubId      uint64
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2SubscriptionFundedIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionFunded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2SubscriptionFunded)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseSubscriptionFunded(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionFunded, error) {
	event := new(NoCancelVRFCoordinatorV2SubscriptionFunded)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator struct {
	Event *NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequested)
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
		it.Event = new(NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequested)
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

func (it *NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequested struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequested)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequested, error) {
	event := new(NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequested)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type NoCancelVRFCoordinatorV2SubscriptionOwnerTransferredIterator struct {
	Event *NoCancelVRFCoordinatorV2SubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *NoCancelVRFCoordinatorV2SubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NoCancelVRFCoordinatorV2SubscriptionOwnerTransferred)
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
		it.Event = new(NoCancelVRFCoordinatorV2SubscriptionOwnerTransferred)
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

func (it *NoCancelVRFCoordinatorV2SubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *NoCancelVRFCoordinatorV2SubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type NoCancelVRFCoordinatorV2SubscriptionOwnerTransferred struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &NoCancelVRFCoordinatorV2SubscriptionOwnerTransferredIterator{contract: _NoCancelVRFCoordinatorV2.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _NoCancelVRFCoordinatorV2.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(NoCancelVRFCoordinatorV2SubscriptionOwnerTransferred)
				if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2Filterer) ParseSubscriptionOwnerTransferred(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionOwnerTransferred, error) {
	event := new(NoCancelVRFCoordinatorV2SubscriptionOwnerTransferred)
	if err := _NoCancelVRFCoordinatorV2.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _NoCancelVRFCoordinatorV2.abi.Events["ConfigSet"].ID:
		return _NoCancelVRFCoordinatorV2.ParseConfigSet(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["FundsRecovered"].ID:
		return _NoCancelVRFCoordinatorV2.ParseFundsRecovered(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["OwnershipTransferRequested"].ID:
		return _NoCancelVRFCoordinatorV2.ParseOwnershipTransferRequested(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["OwnershipTransferred"].ID:
		return _NoCancelVRFCoordinatorV2.ParseOwnershipTransferred(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["ProvingKeyDeregistered"].ID:
		return _NoCancelVRFCoordinatorV2.ParseProvingKeyDeregistered(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["ProvingKeyRegistered"].ID:
		return _NoCancelVRFCoordinatorV2.ParseProvingKeyRegistered(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["RandomWordsFulfilled"].ID:
		return _NoCancelVRFCoordinatorV2.ParseRandomWordsFulfilled(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["RandomWordsRequested"].ID:
		return _NoCancelVRFCoordinatorV2.ParseRandomWordsRequested(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["SubscriptionCanceled"].ID:
		return _NoCancelVRFCoordinatorV2.ParseSubscriptionCanceled(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["SubscriptionConsumerAdded"].ID:
		return _NoCancelVRFCoordinatorV2.ParseSubscriptionConsumerAdded(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _NoCancelVRFCoordinatorV2.ParseSubscriptionConsumerRemoved(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["SubscriptionCreated"].ID:
		return _NoCancelVRFCoordinatorV2.ParseSubscriptionCreated(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["SubscriptionFunded"].ID:
		return _NoCancelVRFCoordinatorV2.ParseSubscriptionFunded(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _NoCancelVRFCoordinatorV2.ParseSubscriptionOwnerTransferRequested(log)
	case _NoCancelVRFCoordinatorV2.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _NoCancelVRFCoordinatorV2.ParseSubscriptionOwnerTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (NoCancelVRFCoordinatorV2ConfigSet) Topic() common.Hash {
	return common.HexToHash("0xc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb2")
}

func (NoCancelVRFCoordinatorV2FundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (NoCancelVRFCoordinatorV2OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (NoCancelVRFCoordinatorV2OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (NoCancelVRFCoordinatorV2ProvingKeyDeregistered) Topic() common.Hash {
	return common.HexToHash("0x72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d")
}

func (NoCancelVRFCoordinatorV2ProvingKeyRegistered) Topic() common.Hash {
	return common.HexToHash("0xe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b8")
}

func (NoCancelVRFCoordinatorV2RandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4")
}

func (NoCancelVRFCoordinatorV2RandomWordsRequested) Topic() common.Hash {
	return common.HexToHash("0x63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772")
}

func (NoCancelVRFCoordinatorV2SubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0xe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815")
}

func (NoCancelVRFCoordinatorV2SubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0")
}

func (NoCancelVRFCoordinatorV2SubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b")
}

func (NoCancelVRFCoordinatorV2SubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf")
}

func (NoCancelVRFCoordinatorV2SubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8")
}

func (NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be")
}

func (NoCancelVRFCoordinatorV2SubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0")
}

func (_NoCancelVRFCoordinatorV2 *NoCancelVRFCoordinatorV2) Address() common.Address {
	return _NoCancelVRFCoordinatorV2.address
}

type NoCancelVRFCoordinatorV2Interface interface {
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

	FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc NoCancelVRFCoordinatorV2RequestCommitment) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64, newOwner common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig NoCancelVRFCoordinatorV2FeeConfig) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*NoCancelVRFCoordinatorV2ConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2ConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*NoCancelVRFCoordinatorV2ConfigSet, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*NoCancelVRFCoordinatorV2FundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2FundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*NoCancelVRFCoordinatorV2FundsRecovered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NoCancelVRFCoordinatorV2OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*NoCancelVRFCoordinatorV2OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*NoCancelVRFCoordinatorV2OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*NoCancelVRFCoordinatorV2OwnershipTransferred, error)

	FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*NoCancelVRFCoordinatorV2ProvingKeyDeregisteredIterator, error)

	WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2ProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyDeregistered(log types.Log) (*NoCancelVRFCoordinatorV2ProvingKeyDeregistered, error)

	FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*NoCancelVRFCoordinatorV2ProvingKeyRegisteredIterator, error)

	WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2ProvingKeyRegistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyRegistered(log types.Log) (*NoCancelVRFCoordinatorV2ProvingKeyRegistered, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*NoCancelVRFCoordinatorV2RandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2RandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*NoCancelVRFCoordinatorV2RandomWordsFulfilled, error)

	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*NoCancelVRFCoordinatorV2RandomWordsRequestedIterator, error)

	WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2RandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsRequested(log types.Log) (*NoCancelVRFCoordinatorV2RandomWordsRequested, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionCanceled, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionConsumerAdded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionCreated, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionFunded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionFunded, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*NoCancelVRFCoordinatorV2SubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *NoCancelVRFCoordinatorV2SubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*NoCancelVRFCoordinatorV2SubscriptionOwnerTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
