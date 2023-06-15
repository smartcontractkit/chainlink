// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_coordinator_v2_5

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

type VRFCoordinatorV25FeeConfig struct {
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

type VRFCoordinatorV25RequestCommitment struct {
	BlockNum         uint64
	SubId            uint64
	CallbackGasLimit uint32
	NumWords         uint32
	Sender           common.Address
	NativePayment    bool
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

var VRFCoordinatorV25MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"blockhashStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"internalBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"externalBalance\",\"type\":\"uint256\"}],\"name\":\"BalanceInvariantViolated\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"BlockhashNotInStore\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToSendEther\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"GasLimitTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectCommitment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"have\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"want\",\"type\":\"uint256\"}],\"name\":\"InsufficientGasForConsumer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCalldata\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"InvalidConsumer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"linkWei\",\"type\":\"int256\"}],\"name\":\"InvalidLinkWeiPrice\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"have\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"min\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"max\",\"type\":\"uint16\"}],\"name\":\"InvalidRequestConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSubscription\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LinkAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposedOwner\",\"type\":\"address\"}],\"name\":\"MustBeRequestedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"MustBeSubOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCorrespondingRequest\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"NoSuchProvingKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"have\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"want\",\"type\":\"uint32\"}],\"name\":\"NumWordsTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableFromLink\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingRequestExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"}],\"name\":\"ProvingKeyAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyConsumers\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structVRFCoordinatorV2_5.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EthFundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountLink\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountEth\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldEthBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newEthBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFundedWithEth\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKHASH_STORE\",\"outputs\":[{\"internalType\":\"contractBlockhashStoreInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK\",\"outputs\":[{\"internalType\":\"contractLinkTokenInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LINK_ETH_FEED\",\"outputs\":[{\"internalType\":\"contractAggregatorV3Interface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONSUMERS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUM_WORDS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_REQUEST_CONFIRMATIONS\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"acceptSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"addConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"cancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createSubscription\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"deregisterProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"gamma\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"uWitness\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"cGammaWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"sHashWitness\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"zInv\",\"type\":\"uint256\"}],\"internalType\":\"structVRF.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"blockNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"nativePayment\",\"type\":\"bool\"}],\"internalType\":\"structVRFCoordinatorV2_5.RequestCommitment\",\"name\":\"rc\",\"type\":\"tuple\"}],\"name\":\"fulfillRandomWords\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"fundSubscriptionWithEth\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"}],\"name\":\"getFeeTier\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRequestConfig\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"getSubscription\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"balance\",\"type\":\"uint96\"},{\"internalType\":\"uint96\",\"name\":\"ethBalance\",\"type\":\"uint96\"},{\"internalType\":\"uint64\",\"name\":\"reqCount\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"consumers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"publicKey\",\"type\":\"uint256[2]\"}],\"name\":\"hashOfKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onTokenTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"oracleWithdrawEth\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"ownerCancelSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"}],\"name\":\"pendingRequestExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverEthFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"publicProvingKey\",\"type\":\"uint256[2]\"}],\"name\":\"registerProvingKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"removeConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"payInEth\",\"type\":\"bool\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"requestSubscriptionOwnerTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_config\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_currentSubId\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_fallbackWeiPerUnitLink\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_feeConfig\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_provingKeyHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"s_provingKeys\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requestCommitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_totalEthBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"internalType\":\"structVRFCoordinatorV2_5.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"link\",\"type\":\"address\"}],\"name\":\"setLINK\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"linkEthFeed\",\"type\":\"address\"}],\"name\":\"setLinkEthFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162005bc538038062005bc5833981016040819052620000349162000188565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be81620000dc565b505060016002555060601b6001600160601b031916608052620001ba565b6001600160a01b038116331415620001375760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000602082840312156200019b57600080fd5b81516001600160a01b0381168114620001b357600080fd5b9392505050565b60805160601c6159e5620001e0600039600081816105880152613b6901526159e56000f3fe6080604052600436106102e65760003560e01c806386fe91c711610184578063caf70c4a116100d6578063e72f6e301161008a578063ee9d2d3811610064578063ee9d2d38146109e1578063efcf1d9414610a0e578063f2fde38b14610a2e57600080fd5b8063e72f6e301461096a578063e82ad7d41461098a578063e95704bd146109ba57600080fd5b8063d7ae1d30116100bb578063d7ae1d30146108f4578063d98e620e14610914578063da2f26101461093457600080fd5b8063caf70c4a146108b4578063d2f9f9a7146108d457600080fd5b8063a02e061611610138578063a4c0ed3611610112578063a4c0ed3614610854578063a8cb447b14610874578063ad1783611461089457600080fd5b8063a02e0616146107ee578063a21a23e41461080e578063a47c76961461082357600080fd5b80638da5cb5b116101695780638da5cb5b14610790578063981837d5146107ae5780639f87fad7146107ce57600080fd5b806386fe91c71461072c5780638741e3ee1461075857600080fd5b806346d8d4861161023d5780636b6feccc116101f157806379ba5097116101cb57806379ba5097146106bf57806382359740146106d4578063827f1c5d146106f457600080fd5b80636b6feccc146105aa5780636f64f03f1461067f5780637341c10c1461069f57600080fd5b806364d51a2a1161022257806364d51a2a1461054157806366316d8d14610556578063689c45171461057657600080fd5b806346d8d486146105015780634cb48a541461052157600080fd5b806308821d581161029f5780631b6b6d23116102795780631b6b6d231461048b5780633697af8b146104c357806340d6bb82146104d657600080fd5b806308821d58146103f757806315c48b8414610417578063181f5a771461043f57600080fd5b8063043bd6ae116102d0578063043bd6ae1461033a57806304c357cb1461035e578063088070f51461037e57600080fd5b8062012291146102eb57806302bcc5b614610318575b600080fd5b3480156102f757600080fd5b50610300610a4e565b60405161030f939291906155ac565b60405180910390f35b34801561032457600080fd5b506103386103333660046153ec565b610aca565b005b34801561034657600080fd5b50610350600e5481565b60405190815260200161030f565b34801561036a57600080fd5b50610338610379366004615407565b610b41565b34801561038a57600080fd5b50600f546103c69061ffff81169063ffffffff620100008204811691660100000000000081048216916a01000000000000000000009091041684565b6040805161ffff909516855263ffffffff93841660208601529183169184019190915216606082015260800161030f565b34801561040357600080fd5b5061033861041236600461510d565b610caa565b34801561042357600080fd5b5061042c60c881565b60405161ffff909116815260200161030f565b34801561044b57600080fd5b50604080518082018252601681527f565246436f6f7264696e61746f72563220322e302e30000000000000000000006020820152905161030f9190615557565b34801561049757600080fd5b506003546104ab906001600160a01b031681565b6040516001600160a01b03909116815260200161030f565b6103386104d13660046153ec565b610e3e565b3480156104e257600080fd5b506104ec6101f481565b60405163ffffffff909116815260200161030f565b34801561050d57600080fd5b5061033861051c366004615023565b610fc0565b34801561052d57600080fd5b5061033861053c3660046152cd565b61115b565b34801561054d57600080fd5b5061042c606481565b34801561056257600080fd5b50610338610571366004615023565b611453565b34801561058257600080fd5b506104ab7f000000000000000000000000000000000000000000000000000000000000000081565b3480156105b657600080fd5b506010546106219063ffffffff808216916401000000008104821691680100000000000000008204811691600160601b8104821691600160801b8204169062ffffff600160a01b8204811691600160b81b8104821691600160d01b8204811691600160e81b90041689565b6040805163ffffffff9a8b168152988a1660208a015296891696880196909652938716606087015291909516608085015262ffffff94851660a0850152841660c0840152831660e0830152919091166101008201526101200161030f565b34801561068b57600080fd5b5061033861069a366004615058565b611610565b3480156106ab57600080fd5b506103386106ba366004615407565b611728565b3480156106cb57600080fd5b5061033861190b565b3480156106e057600080fd5b506103386106ef3660046153ec565b6119bc565b34801561070057600080fd5b5061071461070f366004615205565b611b43565b6040516001600160601b03909116815260200161030f565b34801561073857600080fd5b50600754610714906801000000000000000090046001600160601b031681565b34801561076457600080fd5b50600754610778906001600160401b031681565b6040516001600160401b03909116815260200161030f565b34801561079c57600080fd5b506000546001600160a01b03166104ab565b3480156107ba57600080fd5b506103386107c9366004615006565b612074565b3480156107da57600080fd5b506103386107e9366004615407565b61209e565b3480156107fa57600080fd5b50610338610809366004615006565b61240e565b34801561081a57600080fd5b5061077861247b565b34801561082f57600080fd5b5061084361083e3660046153ec565b61268c565b60405161030f95949392919061574a565b34801561086057600080fd5b5061033861086f366004615085565b612795565b34801561088057600080fd5b5061033861088f366004615006565b612991565b3480156108a057600080fd5b50600a546104ab906001600160a01b031681565b3480156108c057600080fd5b506103506108cf366004615129565b612aad565b3480156108e057600080fd5b506104ec6108ef3660046153ec565b612add565b34801561090057600080fd5b5061033861090f366004615407565b612c5f565b34801561092057600080fd5b5061035061092f366004615162565b612d4a565b34801561094057600080fd5b506104ab61094f366004615162565b600b602052600090815260409020546001600160a01b031681565b34801561097657600080fd5b50610338610985366004615006565b612d6b565b34801561099657600080fd5b506109aa6109a53660046153ec565b612f33565b604051901515815260200161030f565b3480156109c657600080fd5b5060075461071490600160a01b90046001600160601b031681565b3480156109ed57600080fd5b506103506109fc366004615162565b600d6020526000908152604090205481565b348015610a1a57600080fd5b50610350610a29366004615194565b613153565b348015610a3a57600080fd5b50610338610a49366004615006565b6134d2565b600f54600c805460408051602080840282018101909252828152600094859460609461ffff8316946201000090930463ffffffff16939192839190830182828015610ab857602002820191906000526020600020905b815481526020019060010190808311610aa4575b50505050509050925092509250909192565b610ad26134e3565b6001600160401b0381166000908152600560205260409020546001600160a01b0316610b1157604051630fb532db60e11b815260040160405180910390fd5b6001600160401b038116600090815260056020526040902054610b3e9082906001600160a01b031661353f565b50565b6001600160401b03821660009081526005602052604090205482906001600160a01b031680610b8357604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b03821614610bbc57604051636c51fda960e11b81526001600160a01b03821660048201526024015b60405180910390fd5b600280541415610bfc5760405162461bcd60e51b815260206004820152601f60248201526000805160206159b98339815191526044820152606401610bb3565b600280556001600160401b0384166000908152600560205260409020600101546001600160a01b03848116911614610c9f576001600160401b03841660008181526005602090815260409182902060010180546001600160a01b0319166001600160a01b0388169081179091558251338152918201527f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be91015b60405180910390a25b505060016002555050565b610cb26134e3565b604080518082018252600091610ce1919084906002908390839080828437600092019190915250612aad915050565b6000818152600b60205260409020549091506001600160a01b031680610d1d57604051631dfd6e1360e21b815260048101839052602401610bb3565b6000828152600b6020526040812080546001600160a01b03191690555b600c54811015610df55782600c8281548110610d5857610d58615969565b90600052602060002001541415610de357600c805460009190610d7d90600190615892565b81548110610d8d57610d8d615969565b9060005260206000200154905080600c8381548110610dae57610dae615969565b600091825260209091200155600c805480610dcb57610dcb615953565b60019003818190600052602060002001600090559055505b80610ded816158d1565b915050610d3a565b50806001600160a01b03167f72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d83604051610e3191815260200190565b60405180910390a2505050565b600280541415610e7e5760405162461bcd60e51b815260206004820152601f60248201526000805160206159b98339815191526044820152606401610bb3565b600280556001600160401b0381166000908152600560205260409020546001600160a01b0316610ec157604051630fb532db60e11b815260040160405180910390fd5b6001600160401b03811660009081526006602052604090208054600160601b90046001600160601b0316903490600c610efa838561583d565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555034600760148282829054906101000a90046001600160601b0316610f42919061583d565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550816001600160401b03167f4e4421a74ab9b76003de9d9527153d002f1338c36c1bbc180069671631cf0958823484610f9e91906157fa565b604080519283526020830191909152015b60405180910390a250506001600255565b6002805414156110005760405162461bcd60e51b815260206004820152601f60248201526000805160206159b98339815191526044820152606401610bb3565b60028055336000908152600960205260409020546001600160601b038083169116101561104057604051631e9acf1760e31b815260040160405180910390fd5b33600090815260096020526040812080548392906110689084906001600160601b03166158a9565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600760148282829054906101000a90046001600160601b03166110b091906158a9565b92506101000a8154816001600160601b0302191690836001600160601b031602179055506000826001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d806000811461112a576040519150601f19603f3d011682016040523d82523d6000602084013e61112f565b606091505b505090508061115157604051630dcf35db60e41b815260040160405180910390fd5b5050600160025550565b6111636134e3565b60c861ffff8716111561119d5760405163539c34bb60e11b815261ffff871660048201819052602482015260c86044820152606401610bb3565b600082136111c1576040516321ea67b360e11b815260048101839052602401610bb3565b604080516080808201835261ffff891680835263ffffffff89811660208086018290528a83168688018190528a84166060978801819052600f805465ffffffffffff191690961762010000909402939093176dffffffffffffffff000000000000191666010000000000009091026dffffffff000000000000000000001916176a0100000000000000000000909202919091179092558551601080549388015188880151968901519589015160a08a015160c08b015160e08c01516101008d015196881667ffffffffffffffff199099169890981764010000000094881694909402939093177fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff1668010000000000000000998716999099027fffffffffffffffffffffffffffffffff00000000ffffffffffffffffffffffff1698909817600160601b97861697909702969096177fffffffffffffffffff00000000000000ffffffffffffffffffffffffffffffff16600160801b96909416959095027fffffffffffffffffff000000ffffffffffffffffffffffffffffffffffffffff1692909217600160a01b62ffffff96871602177fffffff000000000000ffffffffffffffffffffffffffffffffffffffffffffff16600160b81b948616949094027fffffff000000ffffffffffffffffffffffffffffffffffffffffffffffffffff1693909317600160d01b92851692909202919091177cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff16600160e81b9390921692909202178155600e84905590517fc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb29161144391899189918991899189919061560b565b60405180910390a1505050505050565b6002805414156114935760405162461bcd60e51b815260206004820152601f60248201526000805160206159b98339815191526044820152606401610bb3565b60028055336000908152600860205260409020546001600160601b03808316911610156114d357604051631e9acf1760e31b815260040160405180910390fd5b33600090815260086020526040812080548392906114fb9084906001600160601b03166158a9565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600760088282829054906101000a90046001600160601b031661154391906158a9565b82546101009290920a6001600160601b0381810219909316918316021790915560035460405163a9059cbb60e01b81526001600160a01b03868116600483015292851660248201529116915063a9059cbb90604401602060405180830381600087803b1580156115b257600080fd5b505af11580156115c6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115ea9190615145565b61160757604051631e9acf1760e31b815260040160405180910390fd5b50506001600255565b6116186134e3565b604080518082018252600091611647919084906002908390839080828437600092019190915250612aad915050565b6000818152600b60205260409020549091506001600160a01b03161561169c576040517f4a0b8fa700000000000000000000000000000000000000000000000000000000815260048101829052602401610bb3565b6000818152600b6020908152604080832080546001600160a01b0319166001600160a01b038816908117909155600c805460018101825594527fdf6966c971051c3d54ec59162606531493a51404a002842f56009d7e5cf4a8c7909301849055518381527fe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b89101610e31565b6001600160401b03821660009081526005602052604090205482906001600160a01b03168061176a57604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b0382161461179e57604051636c51fda960e11b81526001600160a01b0382166004820152602401610bb3565b6002805414156117de5760405162461bcd60e51b815260206004820152601f60248201526000805160206159b98339815191526044820152606401610bb3565b60028080556001600160401b0385166000908152600560205260409020015460641415611837576040517f05a48e0f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b03831660009081526004602090815260408083206001600160401b038089168552925290912054161561187057610c9f565b6001600160a01b03831660008181526004602090815260408083206001600160401b038916808552908352818420805467ffffffffffffffff19166001908117909155600584528285206002018054918201815585529383902090930180546001600160a01b031916851790555192835290917f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e09101610c96565b6001546001600160a01b031633146119655760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610bb3565b60008054336001600160a01b0319808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6002805414156119fc5760405162461bcd60e51b815260206004820152601f60248201526000805160206159b98339815191526044820152606401610bb3565b600280556001600160401b0381166000908152600560205260409020546001600160a01b0316611a3f57604051630fb532db60e11b815260040160405180910390fd5b6001600160401b0381166000908152600560205260409020600101546001600160a01b03163314611ac5576001600160401b038116600090815260056020526040908190206001015490517fd084e9750000000000000000000000000000000000000000000000000000000081526001600160a01b039091166004820152602401610bb3565b6001600160401b0381166000818152600560209081526040918290208054336001600160a01b0319808316821784556001909301805490931690925583516001600160a01b03909116808252928101919091529092917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f09101610faf565b6000600280541415611b855760405162461bcd60e51b815260206004820152601f60248201526000805160206159b98339815191526044820152606401610bb3565b6002805560005a90506000611b9a8585613977565b90506000846060015163ffffffff166001600160401b03811115611bc057611bc061597f565b604051908082528060200260200182016040528015611be9578160200160208202803683370190505b50905060005b856060015163ffffffff16811015611c6957826040015181604051602001611c21929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c828281518110611c4c57611c4c615969565b602090810291909101015280611c61816158d1565b915050611bef565b50602080830180516000908152600d90925260408083208390559051905182917f1fe543e30000000000000000000000000000000000000000000000000000000091611cba919086906024016156fc565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090506000611d34886040015163ffffffff16896080015184613ca8565b6020808a0180516001600160401b03908116600090815260069093526040808420549251821684529092208054939450600160c01b918290048316936001939192601892611d86928692900416615812565b92506101000a8154816001600160401b0302191690836001600160401b031602179055506000611de088600f600001600a9054906101000a900463ffffffff1663ffffffff16611dd585612add565b3a8e60a00151613cf6565b90508960a0015115611eff576020808b01516001600160401b03166000908152600690915260409020546001600160601b03808316600160601b909204161015611e3d57604051631e9acf1760e31b815260040160405180910390fd5b60208a8101516001600160401b031660009081526006909152604090208054829190600c90611e7d908490600160601b90046001600160601b03166158a9565b82546101009290920a6001600160601b0381810219909316918316021790915588516000908152600b60209081526040808320546001600160a01b031683526009909152812080548594509092611ed69185911661583d565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550611ffd565b6020808b01516001600160401b03166000908152600690915260409020546001600160601b0380831691161015611f4957604051631e9acf1760e31b815260040160405180910390fd5b6020808b01516001600160401b031660009081526006909152604081208054839290611f7f9084906001600160601b03166158a9565b82546101009290920a6001600160601b0381810219909316918316021790915588516000908152600b60209081526040808320546001600160a01b031683526008909152812080548594509092611fd89185911661583d565b92506101000a8154816001600160601b0302191690836001600160601b031602179055505b6020808801516040808a015160a08e015182519182526001600160601b03861694820194909452921515908301528415156060830152907f04624b1adcd18b732732ed3009cbdcc7dfc5f0340105149d0c193e8361d955e19060800160405180910390a260016002559a9950505050505050505050565b61207c6134e3565b600a80546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160401b03821660009081526005602052604090205482906001600160a01b0316806120e057604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b0382161461211457604051636c51fda960e11b81526001600160a01b0382166004820152602401610bb3565b6002805414156121545760405162461bcd60e51b815260206004820152601f60248201526000805160206159b98339815191526044820152606401610bb3565b6002805561216184612f33565b1561217f57604051631685ecdd60e31b815260040160405180910390fd5b6001600160a01b03831660009081526004602090815260408083206001600160401b038089168552925290912054166121e557604051637800cff360e11b81526001600160401b03851660048201526001600160a01b0384166024820152604401610bb3565b6001600160401b03841660009081526005602090815260408083206002018054825181850281018501909352808352919290919083018282801561225257602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612234575b505050505090506000600182516122699190615892565b905060005b825181101561239057856001600160a01b031683828151811061229357612293615969565b60200260200101516001600160a01b0316141561237e5760008383815181106122be576122be615969565b6020026020010151905080600560008a6001600160401b03166001600160401b03168152602001908152602001600020600201838154811061230257612302615969565b600091825260208083209190910180546001600160a01b0319166001600160a01b0394909416939093179092556001600160401b038a16815260059091526040902060020180548061235657612356615953565b600082815260209020810160001990810180546001600160a01b031916905501905550612390565b80612388816158d1565b91505061226e565b506001600160a01b03851660008181526004602090815260408083206001600160401b038b1680855290835292819020805467ffffffffffffffff191690555192835290917f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b91015b60405180910390a25050600160025550505050565b6124166134e3565b6003546001600160a01b031615612459576040517f2d118a6e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600380546001600160a01b0319166001600160a01b0392909216919091179055565b60006002805414156124bd5760405162461bcd60e51b815260206004820152601f60248201526000805160206159b98339815191526044820152606401610bb3565b60028055600780546001600160401b03169060006124da836158ec565b82546101009290920a6001600160401b0381810219909316918316021790915560075416905060008060405190808252806020026020018201604052801561252c578160200160208202803683370190505b506040805160608082018352600080835260208084018281528486018381526001600160401b038a8116808652600685528886209751885494519351909216600160c01b0277ffffffffffffffffffffffffffffffffffffffffffffffff6001600160601b03948516600160601b027fffffffffffffffff0000000000000000000000000000000000000000000000009096169390941692909217939093179190911617909455845192830185523383528281018281528386018781529483526005825294909120825181546001600160a01b039182166001600160a01b03199182161783559551600183018054919092169616959095179094559151805194955090936126409260028501920190614d5a565b50506040513381526001600160401b03841691507f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf9060200160405180910390a2509050600160025590565b6001600160401b0381166000908152600560205260408120548190819081906060906001600160a01b03166126d457604051630fb532db60e11b815260040160405180910390fd5b6001600160401b0380871660009081526006602090815260408083205460058352928190208054600290910180548351818602810186019094528084526001600160601b0380871697600160601b880490911696600160c01b900416946001600160a01b0390931693919283919083018282801561277b57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161275d575b505050505090509450945094509450945091939590929450565b6002805414156127d55760405162461bcd60e51b815260206004820152601f60248201526000805160206159b98339815191526044820152606401610bb3565b600280556003546001600160a01b0316331461281d576040517f44b0e3c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208114612857576040517f8129bbcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000612865828401846153ec565b6001600160401b0381166000908152600560205260409020549091506001600160a01b03166128a757604051630fb532db60e11b815260040160405180910390fd5b6001600160401b038116600090815260066020526040812080546001600160601b0316918691906128d8838561583d565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555084600760088282829054906101000a90046001600160601b0316612920919061583d565b92506101000a8154816001600160601b0302191690836001600160601b03160217905550816001600160401b03167fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f882878461297c91906157fa565b604080519283526020830191909152016123f9565b6129996134e3565b6007544790600160a01b90046001600160601b0316818111156129d9576040516354ced18160e11b81526004810182905260248101839052604401610bb3565b81811015612aa85760006129ed8284615892565b90506000846001600160a01b03168260405160006040518083038185875af1925050503d8060008114612a3c576040519150601f19603f3d011682016040523d82523d6000602084013e612a41565b606091505b5050905080612a6357604051630dcf35db60e41b815260040160405180910390fd5b604080516001600160a01b0387168152602081018490527f879c9ea2b9d5345b84ccd12610b032602808517cebdb795007f3dcb4df377317910160405180910390a150505b505050565b600081604051602001612ac09190615549565b604051602081830303815290604052805190602001209050919050565b604080516101208101825260105463ffffffff8082168352640100000000820481166020840152680100000000000000008204811693830193909352600160601b810483166060830152600160801b8104909216608082015262ffffff600160a01b8304811660a08301819052600160b81b8404821660c0840152600160d01b8404821660e0840152600160e81b909304166101008201526000916001600160401b03841611612b8e575192915050565b826001600160401b03168160a0015162ffffff16108015612bc157508060c0015162ffffff16836001600160401b031611155b15612bd0576020015192915050565b826001600160401b03168160c0015162ffffff16108015612c0357508060e0015162ffffff16836001600160401b031611155b15612c12576040015192915050565b826001600160401b03168160e0015162ffffff16108015612c46575080610100015162ffffff16836001600160401b031611155b15612c55576060015192915050565b6080015192915050565b6001600160401b03821660009081526005602052604090205482906001600160a01b031680612ca157604051630fb532db60e11b815260040160405180910390fd5b336001600160a01b03821614612cd557604051636c51fda960e11b81526001600160a01b0382166004820152602401610bb3565b600280541415612d155760405162461bcd60e51b815260206004820152601f60248201526000805160206159b98339815191526044820152606401610bb3565b60028055612d2284612f33565b15612d4057604051631685ecdd60e31b815260040160405180910390fd5b610c9f848461353f565b600c8181548110612d5a57600080fd5b600091825260209091200154905081565b612d736134e3565b6003546040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000916001600160a01b0316906370a082319060240160206040518083038186803b158015612dd057600080fd5b505afa158015612de4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e08919061517b565b6007549091506801000000000000000090046001600160601b031681811115612e4e576040516354ced18160e11b81526004810182905260248101839052604401610bb3565b81811015612aa8576000612e628284615892565b60035460405163a9059cbb60e01b81526001600160a01b0387811660048301526024820184905292935091169063a9059cbb90604401602060405180830381600087803b158015612eb257600080fd5b505af1158015612ec6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612eea9190615145565b50604080516001600160a01b0386168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a150505050565b6001600160401b0381166000908152600560209081526040808320815160608101835281546001600160a01b0390811682526001830154168185015260028201805484518187028101870186528181528796939586019390929190830182828015612fc757602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311612fa9575b505050505081525050905060005b8160400151518110156131495760005b600c548110156131365760006130ff600c838154811061300757613007615969565b90600052602060002001548560400151858151811061302857613028615969565b602002602001015188600460008960400151898151811061304b5761304b615969565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081206001600160401b03808f16835293522054166040805160208082018790526001600160a01b0395909516818301526001600160401b039384166060820152919092166080808301919091528251808303909101815260a08201835280519084012060c082019490945260e080820185905282518083039091018152610100909101909152805191012091565b506000818152600d6020526040902054909150156131235750600195945050505050565b508061312e816158d1565b915050612fe5565b5080613141816158d1565b915050612fd5565b5060009392505050565b60006002805414156131955760405162461bcd60e51b815260206004820152601f60248201526000805160206159b98339815191526044820152606401610bb3565b600280556001600160401b0386166000908152600560205260409020546001600160a01b03166131d857604051630fb532db60e11b815260040160405180910390fd5b3360009081526004602090815260408083206001600160401b03808b168552925290912054168061322d57604051637800cff360e11b81526001600160401b0388166004820152336024820152604401610bb3565b600f5461ffff9081169087161080613249575060c861ffff8716115b1561328057600f5460405163539c34bb60e11b815261ffff8089166004830152909116602482015260c86044820152606401610bb3565b600f5463ffffffff62010000909104811690861611156132e757600f546040517ff5d7e01e00000000000000000000000000000000000000000000000000000000815263ffffffff8088166004830152620100009092049091166024820152604401610bb3565b6101f463ffffffff85161115613339576040517f47386bec00000000000000000000000000000000000000000000000000000000815263ffffffff851660048201526101f46024820152604401610bb3565b6000613346826001615812565b6040805160208082018d905233828401526001600160401b03808d16606084015284166080808401919091528351808403909101815260a08301845280519082012060c083018e905260e0808401829052845180850390910181526101009093019093528151910120919250816133bb613d29565b6040805160208101939093528201526001600160401b038b16606082015263ffffffff808a166080830152881660a08201523360c082015260e00160408051601f1981840301815282825280516020918201206000868152600d835283902055848352820183905261ffff8b169082015263ffffffff808a1660608301528816608082015286151560a082015233906001600160401b038c16908d907fad8d0b36e0dd1bc233432f9e3c28d376b9713fc9d26ac2dc972077be458d86fe9060c00160405180910390a4503360009081526004602090815260408083206001600160401b03808e168552925290912080549190931667ffffffffffffffff199091161790915591505060016002559695505050505050565b6134da6134e3565b610b3e81613dc2565b6000546001600160a01b0316331461353d5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610bb3565b565b60028054141561357f5760405162461bcd60e51b815260206004820152601f60248201526000805160206159b98339815191526044820152606401610bb3565b60028080556001600160401b0383166000908152600560209081526040808320815160608101835281546001600160a01b039081168252600183015416818501529481018054835181860281018601855281815295969592949386019383018282801561361557602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116135f7575b505050919092525050506001600160401b038085166000908152600660209081526040808320815160608101835290546001600160601b03808216808452600160601b8304909116948301859052600160c01b90910490951691810191909152939450905b8460400151518110156136fb5760046000866040015183815181106136a1576136a1615969565b6020908102919091018101516001600160a01b0316825281810192909252604090810160009081206001600160401b038b1682529092529020805467ffffffffffffffff19169055806136f3816158d1565b91505061367a565b506001600160401b038616600090815260056020526040812080546001600160a01b0319908116825560018201805490911690559061373d6002830182614dbf565b50506001600160401b038616600090815260066020526040812055600780548391906008906137829084906801000000000000000090046001600160601b03166158a9565b92506101000a8154816001600160601b0302191690836001600160601b0316021790555080600760148282829054906101000a90046001600160601b03166137ca91906158a9565b82546101009290920a6001600160601b0381810219909316918316021790915560035460405163a9059cbb60e01b81526001600160a01b03898116600483015292861660248201529116915063a9059cbb90604401602060405180830381600087803b15801561383957600080fd5b505af115801561384d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138719190615145565b61388e57604051631e9acf1760e31b815260040160405180910390fd5b6000856001600160a01b0316826001600160601b031660405160006040518083038185875af1925050503d80600081146138e4576040519150601f19603f3d011682016040523d82523d6000602084013e6138e9565b606091505b505090508061390b57604051630dcf35db60e41b815260040160405180910390fd5b604080516001600160a01b03881681526001600160601b03858116602083015284168183015290516001600160401b038916917fb3a76e2c52bb7a484e325118e4cadc4346fdaa16a7df0621c5525e1f7ab4d495919081900360600190a2505060016002555050505050565b604080516060810182526000808252602082018190529181019190915260006139a38460000151612aad565b6000818152600b60205260409020549091506001600160a01b0316806139df57604051631dfd6e1360e21b815260048101839052602401610bb3565b6000828660800151604051602001613a01929190918252602082015260400190565b60408051601f1981840301815291815281516020928301206000818152600d90935291205490915080613a60576040517f3688124a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b85516020808801516040808a015160608b015160808c01519251613acb96899690959491019586526001600160401b03948516602087015292909316604085015263ffffffff90811660608501529190911660808301526001600160a01b031660a082015260c00190565b604051602081830303815290604052805190602001208114613b19576040517fd529142c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000613b288760000151613e6c565b905080613c325786516040517fe9413d380000000000000000000000000000000000000000000000000000000081526001600160401b0390911660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063e9413d389060240160206040518083038186803b158015613bb357600080fd5b505afa158015613bc7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613beb919061517b565b905080613c325786516040517f175dadad0000000000000000000000000000000000000000000000000000000081526001600160401b039091166004820152602401610bb3565b6000886080015182604051602001613c54929190918252602082015260400190565b6040516020818303038152906040528051906020012060001c90506000613c7b8a83613f6f565b90506040518060600160405280888152602001868152602001828152509750505050505050505b92915050565b60005a611388811015613cba57600080fd5b611388810390508460408204820311613cd257600080fd5b50823b613cde57600080fd5b60008083516020850160008789f190505b9392505050565b60008115613d1157613d0a86868686613fda565b9050613d20565b613d1d86868686614044565b90505b95945050505050565b60004661a4b1811480613d3e575062066eed81145b15613dbb5760646001600160a01b031663a3b1b31d6040518163ffffffff1660e01b815260040160206040518083038186803b158015613d7d57600080fd5b505afa158015613d91573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613db5919061517b565b91505090565b4391505090565b6001600160a01b038116331415613e1b5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610bb3565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60004661a4b1811480613e81575062066eed81145b15613f6057610100836001600160401b0316613e9b613d29565b613ea59190615892565b1180613ec15750613eb4613d29565b836001600160401b031610155b15613ecf5750600092915050565b6040517f2b407a820000000000000000000000000000000000000000000000000000000081526001600160401b0384166004820152606490632b407a829060240160206040518083038186803b158015613f2857600080fd5b505afa158015613f3c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613cef919061517b565b50506001600160401b03164090565b6000613fa38360000151846020015185604001518660600151868860a001518960c001518a60e001518b610100015161414b565b60038360200151604051602001613fbb9291906156e8565b60408051601f1981840301815291905280516020909101209392505050565b600080613fe5614386565b905060005a613ff488886157fa565b613ffe9190615892565b6140089085615873565b9050600061402163ffffffff871664e8d4a51000615873565b90508261402e82846157fa565b61403891906157fa565b98975050505050505050565b60008061404f6143e2565b905060008113614075576040516321ea67b360e11b815260048101829052602401610bb3565b600061407f614386565b9050600082825a6140908b8b6157fa565b61409a9190615892565b6140a49088615873565b6140ae91906157fa565b6140c090670de0b6b3a7640000615873565b6140ca919061585f565b905060006140e363ffffffff881664e8d4a51000615873565b90506140fb816b033b2e3c9fd0803ce8000000615892565b821115614134576040517fe80fa38100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61413e81836157fa565b9998505050505050505050565b614154896144ca565b6141a05760405162461bcd60e51b815260206004820152601a60248201527f7075626c6963206b6579206973206e6f74206f6e2063757276650000000000006044820152606401610bb3565b6141a9886144ca565b6141f55760405162461bcd60e51b815260206004820152601560248201527f67616d6d61206973206e6f74206f6e20637572766500000000000000000000006044820152606401610bb3565b6141fe836144ca565b61424a5760405162461bcd60e51b815260206004820152601d60248201527f6347616d6d615769746e657373206973206e6f74206f6e2063757276650000006044820152606401610bb3565b614253826144ca565b61429f5760405162461bcd60e51b815260206004820152601c60248201527f73486173685769746e657373206973206e6f74206f6e206375727665000000006044820152606401610bb3565b6142ab878a88876145a3565b6142f75760405162461bcd60e51b815260206004820152601960248201527f6164647228632a706b2b732a6729213d5f755769746e657373000000000000006044820152606401610bb3565b60006143038a876146d8565b90506000614316898b878b86898961473c565b90506000614327838d8d8a8661485c565b9050808a146143785760405162461bcd60e51b815260206004820152600d60248201527f696e76616c69642070726f6f66000000000000000000000000000000000000006044820152606401610bb3565b505050505050505050505050565b60004661a4b181148061439b575062066eed81145b156143da57606c6001600160a01b031663c6f7de0e6040518163ffffffff1660e01b815260040160206040518083038186803b158015613d7d57600080fd5b600091505090565b600f54600a54604080517ffeaf968c00000000000000000000000000000000000000000000000000000000815290516000936601000000000000900463ffffffff169283151592859283926001600160a01b03169163feaf968c9160048083019260a0929190829003018186803b15801561445c57600080fd5b505afa158015614470573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614494919061543e565b5094509092508491505080156144b857506144af8242615892565b8463ffffffff16105b156144c25750600e545b949350505050565b80516000906401000003d019116145235760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420782d6f7264696e61746500000000000000000000000000006044820152606401610bb3565b60208201516401000003d0191161457c5760405162461bcd60e51b815260206004820152601260248201527f696e76616c696420792d6f7264696e61746500000000000000000000000000006044820152606401610bb3565b60208201516401000003d01990800961459c8360005b602002015161489c565b1492915050565b60006001600160a01b0382166145fb5760405162461bcd60e51b815260206004820152600b60248201527f626164207769746e6573730000000000000000000000000000000000000000006044820152606401610bb3565b60208401516000906001161561461257601c614615565b601b5b9050600070014551231950b75fc4402da1732fc9bebe1985876000602002015109865170014551231950b75fc4402da1732fc9bebe19918203925060009190890987516040805160008082526020820180845287905260ff88169282019290925260608101929092526080820183905291925060019060a0016020604051602081039080840390855afa1580156146b0573d6000803e3d6000fd5b5050604051601f1901516001600160a01b039081169088161495505050505050949350505050565b6146e0614ddd565b61470d600184846040516020016146f993929190615528565b6040516020818303038152906040526148c0565b90505b614719816144ca565b613ca257805160408051602081019290925261473591016146f9565b9050614710565b614744614ddd565b825186516401000003d01990819006910614156147a35760405162461bcd60e51b815260206004820152601e60248201527f706f696e747320696e2073756d206d7573742062652064697374696e637400006044820152606401610bb3565b6147ae87898861490f565b6147fa5760405162461bcd60e51b815260206004820152601660248201527f4669727374206d756c20636865636b206661696c6564000000000000000000006044820152606401610bb3565b61480584868561490f565b6148515760405162461bcd60e51b815260206004820152601760248201527f5365636f6e64206d756c20636865636b206661696c65640000000000000000006044820152606401610bb3565b614038868484614a49565b60006002868686858760405160200161487a969594939291906154c9565b60408051601f1981840301815291905280516020909101209695505050505050565b6000806401000003d01980848509840990506401000003d019600782089392505050565b6148c8614ddd565b6148d182614b10565b81526148e66148e1826000614592565b614b4b565b60208201819052600290066001141561490a576020810180516401000003d0190390525b919050565b60008261495e5760405162461bcd60e51b815260206004820152600b60248201527f7a65726f207363616c61720000000000000000000000000000000000000000006044820152606401610bb3565b8351602085015160009061497490600290615913565b1561498057601c614983565b601b5b9050600070014551231950b75fc4402da1732fc9bebe198387096040805160008082526020820180845281905260ff86169282019290925260608101869052608081018390529192509060019060a0016020604051602081039080840390855afa1580156149f5573d6000803e3d6000fd5b505050602060405103519050600086604051602001614a1491906154b7565b60408051601f1981840301815291905280516020909101206001600160a01b0392831692169190911498975050505050505050565b614a51614ddd565b835160208086015185519186015160009384938493614a7293909190614b6b565b919450925090506401000003d019858209600114614ad25760405162461bcd60e51b815260206004820152601960248201527f696e765a206d75737420626520696e7665727365206f66207a000000000000006044820152606401610bb3565b60405180604001604052806401000003d01980614af157614af161593d565b87860981526020016401000003d0198785099052979650505050505050565b805160208201205b6401000003d019811061490a57604080516020808201939093528151808203840181529082019091528051910120614b18565b6000613ca2826002614b646401000003d01960016157fa565b901c614c4b565b60008080600180826401000003d019896401000003d019038808905060006401000003d0198b6401000003d019038a0890506000614bab83838585614ced565b9098509050614bbc88828e88614d11565b9098509050614bcd88828c87614d11565b90985090506000614be08d878b85614d11565b9098509050614bf188828686614ced565b9098509050614c0288828e89614d11565b9098509050818114614c37576401000003d019818a0998506401000003d01982890997506401000003d0198183099650614c3b565b8196505b5050505050509450945094915050565b600080614c56614dfb565b6020808252818101819052604082015260608101859052608081018490526401000003d01960a0820152614c88614e19565b60208160c0846005600019fa925082614ce35760405162461bcd60e51b815260206004820152601260248201527f6269674d6f64457870206661696c7572652100000000000000000000000000006044820152606401610bb3565b5195945050505050565b6000806401000003d0198487096401000003d0198487099097909650945050505050565b600080806401000003d019878509905060006401000003d01987876401000003d019030990506401000003d0198183086401000003d01986890990999098509650505050505050565b828054828255906000526020600020908101928215614daf579160200282015b82811115614daf57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190614d7a565b50614dbb929150614e37565b5090565b5080546000825590600052602060002090810190610b3e9190614e37565b60405180604001604052806002906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b5b80821115614dbb5760008155600101614e38565b803561490a81615995565b8060408101831015613ca257600080fd5b600082601f830112614e7957600080fd5b604051604081018181106001600160401b0382111715614e9b57614e9b61597f565b8060405250808385604086011115614eb257600080fd5b60005b6002811015614ed4578135835260209283019290910190600101614eb5565b509195945050505050565b600060c08284031215614ef157600080fd5b60405160c081018181106001600160401b0382111715614f1357614f1361597f565b604052905080614f2283614fbe565b8152614f3060208401614fbe565b6020820152614f4160408401614faa565b6040820152614f5260608401614faa565b60608201526080830135614f6581615995565b608082015260a0830135614f78816159aa565b60a0919091015292915050565b803561ffff8116811461490a57600080fd5b803562ffffff8116811461490a57600080fd5b803563ffffffff8116811461490a57600080fd5b80356001600160401b038116811461490a57600080fd5b805169ffffffffffffffffffff8116811461490a57600080fd5b80356001600160601b038116811461490a57600080fd5b60006020828403121561501857600080fd5b8135613cef81615995565b6000806040838503121561503657600080fd5b823561504181615995565b915061504f60208401614fef565b90509250929050565b6000806060838503121561506b57600080fd5b823561507681615995565b915061504f8460208501614e57565b6000806000806060858703121561509b57600080fd5b84356150a681615995565b93506020850135925060408501356001600160401b03808211156150c957600080fd5b818701915087601f8301126150dd57600080fd5b8135818111156150ec57600080fd5b8860208285010111156150fe57600080fd5b95989497505060200194505050565b60006040828403121561511f57600080fd5b613cef8383614e57565b60006040828403121561513b57600080fd5b613cef8383614e68565b60006020828403121561515757600080fd5b8151613cef816159aa565b60006020828403121561517457600080fd5b5035919050565b60006020828403121561518d57600080fd5b5051919050565b60008060008060008060c087890312156151ad57600080fd5b863595506151bd60208801614fbe565b94506151cb60408801614f85565b93506151d960608801614faa565b92506151e760808801614faa565b915060a08701356151f7816159aa565b809150509295509295509295565b60008082840361026081121561521a57600080fd5b6101a08082121561522a57600080fd5b6152326157d1565b915061523e8686614e68565b825261524d8660408701614e68565b60208301526080850135604083015260a0850135606083015260c0850135608083015261527c60e08601614e4c565b60a083015261010061529087828801614e68565b60c08401526152a3876101408801614e68565b60e084015261018086013581840152508193506152c286828701614edf565b925050509250929050565b6000806000806000808688036101c08112156152e857600080fd5b6152f188614f85565b96506152ff60208901614faa565b955061530d60408901614faa565b945061531b60608901614faa565b93506080880135925061012080609f198301121561533857600080fd5b6153406157d1565b915061534e60a08a01614faa565b825261535c60c08a01614faa565b602083015261536d60e08a01614faa565b6040830152610100615380818b01614faa565b6060840152615390828b01614faa565b60808401526153a26101408b01614f97565b60a08401526153b46101608b01614f97565b60c08401526153c66101808b01614f97565b60e08401526153d86101a08b01614f97565b818401525050809150509295509295509295565b6000602082840312156153fe57600080fd5b613cef82614fbe565b6000806040838503121561541a57600080fd5b61542383614fbe565b9150602083013561543381615995565b809150509250929050565b600080600080600060a0868803121561545657600080fd5b61545f86614fd5565b945060208601519350604086015192506060860151915061548260808701614fd5565b90509295509295909350565b8060005b60028110156154b1578151845260209384019390910190600101615492565b50505050565b6154c1818361548e565b604001919050565b8681526154d9602082018761548e565b6154e6606082018661548e565b6154f360a082018561548e565b61550060e082018461548e565b60609190911b6bffffffffffffffffffffffff19166101208201526101340195945050505050565b838152615538602082018461548e565b606081019190915260800192915050565b60408101613ca2828461548e565b600060208083528351808285015260005b8181101561558457858101830151858201604001528201615568565b81811115615596576000604083870101525b50601f01601f1916929092016040019392505050565b60006060820161ffff86168352602063ffffffff86168185015260606040850152818551808452608086019150828701935060005b818110156155fd578451835293830193918301916001016155e1565b509098975050505050505050565b60006101c08201905061ffff8816825263ffffffff808816602084015280871660408401528086166060840152846080840152835481811660a085015261565f60c08501838360201c1663ffffffff169052565b61567660e08501838360401c1663ffffffff169052565b61568e6101008501838360601c1663ffffffff169052565b6156a66101208501838360801c1663ffffffff169052565b62ffffff60a082901c811661014086015260b882901c811661016086015260d082901c1661018085015260e81c6101a090930192909252979650505050505050565b82815260608101613cef602083018461548e565b6000604082018483526020604081850152818551808452606086019150828701935060005b8181101561573d57845183529383019391830191600101615721565b5090979650505050505050565b600060a082016001600160601b0380891684526020818916818601526001600160401b03881660408601526001600160a01b039150818716606086015260a0608086015282865180855260c087019150828801945060005b818110156157c05785518516835294830194918301916001016157a2565b50909b9a5050505050505050505050565b60405161012081016001600160401b03811182821017156157f4576157f461597f565b60405290565b6000821982111561580d5761580d615927565b500190565b60006001600160401b0380831681851680830382111561583457615834615927565b01949350505050565b60006001600160601b0380831681851680830382111561583457615834615927565b60008261586e5761586e61593d565b500490565b600081600019048311821515161561588d5761588d615927565b500290565b6000828210156158a4576158a4615927565b500390565b60006001600160601b03838116908316818110156158c9576158c9615927565b039392505050565b60006000198214156158e5576158e5615927565b5060010190565b60006001600160401b038083168181141561590957615909615927565b6001019392505050565b6000826159225761592261593d565b500690565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052603160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b0381168114610b3e57600080fd5b8015158114610b3e57600080fdfe5265656e7472616e637947756172643a207265656e7472616e742063616c6c00a164736f6c6343000806000a",
}

var VRFCoordinatorV25ABI = VRFCoordinatorV25MetaData.ABI

var VRFCoordinatorV25Bin = VRFCoordinatorV25MetaData.Bin

func DeployVRFCoordinatorV25(auth *bind.TransactOpts, backend bind.ContractBackend, blockhashStore common.Address) (common.Address, *types.Transaction, *VRFCoordinatorV25, error) {
	parsed, err := VRFCoordinatorV25MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorV25Bin), backend, blockhashStore)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorV25{VRFCoordinatorV25Caller: VRFCoordinatorV25Caller{contract: contract}, VRFCoordinatorV25Transactor: VRFCoordinatorV25Transactor{contract: contract}, VRFCoordinatorV25Filterer: VRFCoordinatorV25Filterer{contract: contract}}, nil
}

type VRFCoordinatorV25 struct {
	address common.Address
	abi     abi.ABI
	VRFCoordinatorV25Caller
	VRFCoordinatorV25Transactor
	VRFCoordinatorV25Filterer
}

type VRFCoordinatorV25Caller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV25Transactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV25Filterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV25Session struct {
	Contract     *VRFCoordinatorV25
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV25CallerSession struct {
	Contract *VRFCoordinatorV25Caller
	CallOpts bind.CallOpts
}

type VRFCoordinatorV25TransactorSession struct {
	Contract     *VRFCoordinatorV25Transactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV25Raw struct {
	Contract *VRFCoordinatorV25
}

type VRFCoordinatorV25CallerRaw struct {
	Contract *VRFCoordinatorV25Caller
}

type VRFCoordinatorV25TransactorRaw struct {
	Contract *VRFCoordinatorV25Transactor
}

func NewVRFCoordinatorV25(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorV25, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorV25ABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFCoordinatorV25(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25{address: address, abi: abi, VRFCoordinatorV25Caller: VRFCoordinatorV25Caller{contract: contract}, VRFCoordinatorV25Transactor: VRFCoordinatorV25Transactor{contract: contract}, VRFCoordinatorV25Filterer: VRFCoordinatorV25Filterer{contract: contract}}, nil
}

func NewVRFCoordinatorV25Caller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorV25Caller, error) {
	contract, err := bindVRFCoordinatorV25(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25Caller{contract: contract}, nil
}

func NewVRFCoordinatorV25Transactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorV25Transactor, error) {
	contract, err := bindVRFCoordinatorV25(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25Transactor{contract: contract}, nil
}

func NewVRFCoordinatorV25Filterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorV25Filterer, error) {
	contract, err := bindVRFCoordinatorV25(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25Filterer{contract: contract}, nil
}

func bindVRFCoordinatorV25(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFCoordinatorV25MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV25.Contract.VRFCoordinatorV25Caller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.VRFCoordinatorV25Transactor.contract.Transfer(opts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.VRFCoordinatorV25Transactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV25.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.contract.Transfer(opts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "BLOCKHASH_STORE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV25.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) BLOCKHASHSTORE() (common.Address, error) {
	return _VRFCoordinatorV25.Contract.BLOCKHASHSTORE(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) LINK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "LINK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) LINK() (common.Address, error) {
	return _VRFCoordinatorV25.Contract.LINK(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) LINK() (common.Address, error) {
	return _VRFCoordinatorV25.Contract.LINK(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) LINKETHFEED(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "LINK_ETH_FEED")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) LINKETHFEED() (common.Address, error) {
	return _VRFCoordinatorV25.Contract.LINKETHFEED(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) LINKETHFEED() (common.Address, error) {
	return _VRFCoordinatorV25.Contract.LINKETHFEED(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) MAXCONSUMERS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "MAX_CONSUMERS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV25.Contract.MAXCONSUMERS(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) MAXCONSUMERS() (uint16, error) {
	return _VRFCoordinatorV25.Contract.MAXCONSUMERS(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) MAXNUMWORDS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "MAX_NUM_WORDS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV25.Contract.MAXNUMWORDS(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) MAXNUMWORDS() (uint32, error) {
	return _VRFCoordinatorV25.Contract.MAXNUMWORDS(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "MAX_REQUEST_CONFIRMATIONS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV25.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) MAXREQUESTCONFIRMATIONS() (uint16, error) {
	return _VRFCoordinatorV25.Contract.MAXREQUESTCONFIRMATIONS(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) GetFeeTier(opts *bind.CallOpts, reqCount uint64) (uint32, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "getFeeTier", reqCount)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) GetFeeTier(reqCount uint64) (uint32, error) {
	return _VRFCoordinatorV25.Contract.GetFeeTier(&_VRFCoordinatorV25.CallOpts, reqCount)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) GetFeeTier(reqCount uint64) (uint32, error) {
	return _VRFCoordinatorV25.Contract.GetFeeTier(&_VRFCoordinatorV25.CallOpts, reqCount)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) GetRequestConfig(opts *bind.CallOpts) (uint16, uint32, [][32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "getRequestConfig")

	if err != nil {
		return *new(uint16), *new(uint32), *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	out1 := *abi.ConvertType(out[1], new(uint32)).(*uint32)
	out2 := *abi.ConvertType(out[2], new([][32]byte)).(*[][32]byte)

	return out0, out1, out2, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _VRFCoordinatorV25.Contract.GetRequestConfig(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) GetRequestConfig() (uint16, uint32, [][32]byte, error) {
	return _VRFCoordinatorV25.Contract.GetRequestConfig(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) GetSubscription(opts *bind.CallOpts, subId uint64) (GetSubscription,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "getSubscription", subId)

	outstruct := new(GetSubscription)
	if err != nil {
		return *outstruct, err
	}

	outstruct.Balance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.EthBalance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ReqCount = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Owner = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Consumers = *abi.ConvertType(out[4], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) GetSubscription(subId uint64) (GetSubscription,

	error) {
	return _VRFCoordinatorV25.Contract.GetSubscription(&_VRFCoordinatorV25.CallOpts, subId)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) GetSubscription(subId uint64) (GetSubscription,

	error) {
	return _VRFCoordinatorV25.Contract.GetSubscription(&_VRFCoordinatorV25.CallOpts, subId)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "hashOfKey", publicKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25.Contract.HashOfKey(&_VRFCoordinatorV25.CallOpts, publicKey)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) HashOfKey(publicKey [2]*big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25.Contract.HashOfKey(&_VRFCoordinatorV25.CallOpts, publicKey)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) Owner() (common.Address, error) {
	return _VRFCoordinatorV25.Contract.Owner(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) Owner() (common.Address, error) {
	return _VRFCoordinatorV25.Contract.Owner(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) PendingRequestExists(opts *bind.CallOpts, subId uint64) (bool, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "pendingRequestExists", subId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) PendingRequestExists(subId uint64) (bool, error) {
	return _VRFCoordinatorV25.Contract.PendingRequestExists(&_VRFCoordinatorV25.CallOpts, subId)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) PendingRequestExists(subId uint64) (bool, error) {
	return _VRFCoordinatorV25.Contract.PendingRequestExists(&_VRFCoordinatorV25.CallOpts, subId)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) SConfig(opts *bind.CallOpts) (SConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "s_config")

	outstruct := new(SConfig)
	if err != nil {
		return *outstruct, err
	}

	outstruct.MinimumRequestConfirmations = *abi.ConvertType(out[0], new(uint16)).(*uint16)
	outstruct.MaxGasLimit = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.StalenessSeconds = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.GasAfterPaymentCalculation = *abi.ConvertType(out[3], new(uint32)).(*uint32)

	return *outstruct, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) SConfig() (SConfig,

	error) {
	return _VRFCoordinatorV25.Contract.SConfig(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) SConfig() (SConfig,

	error) {
	return _VRFCoordinatorV25.Contract.SConfig(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) SCurrentSubId(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "s_currentSubId")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) SCurrentSubId() (uint64, error) {
	return _VRFCoordinatorV25.Contract.SCurrentSubId(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) SCurrentSubId() (uint64, error) {
	return _VRFCoordinatorV25.Contract.SCurrentSubId(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) SFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "s_fallbackWeiPerUnitLink")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV25.Contract.SFallbackWeiPerUnitLink(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) SFallbackWeiPerUnitLink() (*big.Int, error) {
	return _VRFCoordinatorV25.Contract.SFallbackWeiPerUnitLink(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) SFeeConfig(opts *bind.CallOpts) (SFeeConfig,

	error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "s_feeConfig")

	outstruct := new(SFeeConfig)
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) SFeeConfig() (SFeeConfig,

	error) {
	return _VRFCoordinatorV25.Contract.SFeeConfig(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) SFeeConfig() (SFeeConfig,

	error) {
	return _VRFCoordinatorV25.Contract.SFeeConfig(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) SProvingKeyHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "s_provingKeyHashes", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) SProvingKeyHashes(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25.Contract.SProvingKeyHashes(&_VRFCoordinatorV25.CallOpts, arg0)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) SProvingKeyHashes(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25.Contract.SProvingKeyHashes(&_VRFCoordinatorV25.CallOpts, arg0)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) SProvingKeys(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "s_provingKeys", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) SProvingKeys(arg0 [32]byte) (common.Address, error) {
	return _VRFCoordinatorV25.Contract.SProvingKeys(&_VRFCoordinatorV25.CallOpts, arg0)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) SProvingKeys(arg0 [32]byte) (common.Address, error) {
	return _VRFCoordinatorV25.Contract.SProvingKeys(&_VRFCoordinatorV25.CallOpts, arg0)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) SRequestCommitments(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "s_requestCommitments", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) SRequestCommitments(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25.Contract.SRequestCommitments(&_VRFCoordinatorV25.CallOpts, arg0)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) SRequestCommitments(arg0 *big.Int) ([32]byte, error) {
	return _VRFCoordinatorV25.Contract.SRequestCommitments(&_VRFCoordinatorV25.CallOpts, arg0)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) STotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "s_totalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) STotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV25.Contract.STotalBalance(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) STotalBalance() (*big.Int, error) {
	return _VRFCoordinatorV25.Contract.STotalBalance(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) STotalEthBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "s_totalEthBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) STotalEthBalance() (*big.Int, error) {
	return _VRFCoordinatorV25.Contract.STotalEthBalance(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) STotalEthBalance() (*big.Int, error) {
	return _VRFCoordinatorV25.Contract.STotalEthBalance(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Caller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VRFCoordinatorV25.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) TypeAndVersion() (string, error) {
	return _VRFCoordinatorV25.Contract.TypeAndVersion(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25CallerSession) TypeAndVersion() (string, error) {
	return _VRFCoordinatorV25.Contract.TypeAndVersion(&_VRFCoordinatorV25.CallOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "acceptOwnership")
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.AcceptOwnership(&_VRFCoordinatorV25.TransactOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.AcceptOwnership(&_VRFCoordinatorV25.TransactOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "acceptSubscriptionOwnerTransfer", subId)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV25.TransactOpts, subId)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) AcceptSubscriptionOwnerTransfer(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.AcceptSubscriptionOwnerTransfer(&_VRFCoordinatorV25.TransactOpts, subId)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) AddConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "addConsumer", subId, consumer)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.AddConsumer(&_VRFCoordinatorV25.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) AddConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.AddConsumer(&_VRFCoordinatorV25.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) CancelSubscription(opts *bind.TransactOpts, subId uint64, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "cancelSubscription", subId, to)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) CancelSubscription(subId uint64, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.CancelSubscription(&_VRFCoordinatorV25.TransactOpts, subId, to)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) CancelSubscription(subId uint64, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.CancelSubscription(&_VRFCoordinatorV25.TransactOpts, subId, to)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "createSubscription")
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.CreateSubscription(&_VRFCoordinatorV25.TransactOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) CreateSubscription() (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.CreateSubscription(&_VRFCoordinatorV25.TransactOpts)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "deregisterProvingKey", publicProvingKey)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.DeregisterProvingKey(&_VRFCoordinatorV25.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) DeregisterProvingKey(publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.DeregisterProvingKey(&_VRFCoordinatorV25.TransactOpts, publicProvingKey)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFCoordinatorV25RequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "fulfillRandomWords", proof, rc)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorV25RequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.FulfillRandomWords(&_VRFCoordinatorV25.TransactOpts, proof, rc)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) FulfillRandomWords(proof VRFProof, rc VRFCoordinatorV25RequestCommitment) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.FulfillRandomWords(&_VRFCoordinatorV25.TransactOpts, proof, rc)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) FundSubscriptionWithEth(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "fundSubscriptionWithEth", subId)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) FundSubscriptionWithEth(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.FundSubscriptionWithEth(&_VRFCoordinatorV25.TransactOpts, subId)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) FundSubscriptionWithEth(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.FundSubscriptionWithEth(&_VRFCoordinatorV25.TransactOpts, subId)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "onTokenTransfer", arg0, amount, data)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.OnTokenTransfer(&_VRFCoordinatorV25.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) OnTokenTransfer(arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.OnTokenTransfer(&_VRFCoordinatorV25.TransactOpts, arg0, amount, data)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "oracleWithdraw", recipient, amount)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.OracleWithdraw(&_VRFCoordinatorV25.TransactOpts, recipient, amount)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) OracleWithdraw(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.OracleWithdraw(&_VRFCoordinatorV25.TransactOpts, recipient, amount)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) OracleWithdrawEth(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "oracleWithdrawEth", recipient, amount)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) OracleWithdrawEth(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.OracleWithdrawEth(&_VRFCoordinatorV25.TransactOpts, recipient, amount)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) OracleWithdrawEth(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.OracleWithdrawEth(&_VRFCoordinatorV25.TransactOpts, recipient, amount)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "ownerCancelSubscription", subId)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.OwnerCancelSubscription(&_VRFCoordinatorV25.TransactOpts, subId)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) OwnerCancelSubscription(subId uint64) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.OwnerCancelSubscription(&_VRFCoordinatorV25.TransactOpts, subId)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) RecoverEthFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "recoverEthFunds", to)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) RecoverEthFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.RecoverEthFunds(&_VRFCoordinatorV25.TransactOpts, to)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) RecoverEthFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.RecoverEthFunds(&_VRFCoordinatorV25.TransactOpts, to)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "recoverFunds", to)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.RecoverFunds(&_VRFCoordinatorV25.TransactOpts, to)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) RecoverFunds(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.RecoverFunds(&_VRFCoordinatorV25.TransactOpts, to)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "registerProvingKey", oracle, publicProvingKey)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.RegisterProvingKey(&_VRFCoordinatorV25.TransactOpts, oracle, publicProvingKey)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) RegisterProvingKey(oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.RegisterProvingKey(&_VRFCoordinatorV25.TransactOpts, oracle, publicProvingKey)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) RemoveConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "removeConsumer", subId, consumer)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) RemoveConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.RemoveConsumer(&_VRFCoordinatorV25.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) RemoveConsumer(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.RemoveConsumer(&_VRFCoordinatorV25.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, payInEth bool) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "requestRandomWords", keyHash, subId, requestConfirmations, callbackGasLimit, numWords, payInEth)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) RequestRandomWords(keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, payInEth bool) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.RequestRandomWords(&_VRFCoordinatorV25.TransactOpts, keyHash, subId, requestConfirmations, callbackGasLimit, numWords, payInEth)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) RequestRandomWords(keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, payInEth bool) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.RequestRandomWords(&_VRFCoordinatorV25.TransactOpts, keyHash, subId, requestConfirmations, callbackGasLimit, numWords, payInEth)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "requestSubscriptionOwnerTransfer", subId, newOwner)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV25.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) RequestSubscriptionOwnerTransfer(subId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.RequestSubscriptionOwnerTransfer(&_VRFCoordinatorV25.TransactOpts, subId, newOwner)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV25FeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "setConfig", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV25FeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.SetConfig(&_VRFCoordinatorV25.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) SetConfig(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV25FeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.SetConfig(&_VRFCoordinatorV25.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) SetLINK(opts *bind.TransactOpts, link common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "setLINK", link)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) SetLINK(link common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.SetLINK(&_VRFCoordinatorV25.TransactOpts, link)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) SetLINK(link common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.SetLINK(&_VRFCoordinatorV25.TransactOpts, link)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) SetLinkEthFeed(opts *bind.TransactOpts, linkEthFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "setLinkEthFeed", linkEthFeed)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) SetLinkEthFeed(linkEthFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.SetLinkEthFeed(&_VRFCoordinatorV25.TransactOpts, linkEthFeed)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) SetLinkEthFeed(linkEthFeed common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.SetLinkEthFeed(&_VRFCoordinatorV25.TransactOpts, linkEthFeed)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Transactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.contract.Transact(opts, "transferOwnership", to)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Session) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.TransferOwnership(&_VRFCoordinatorV25.TransactOpts, to)
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25TransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV25.Contract.TransferOwnership(&_VRFCoordinatorV25.TransactOpts, to)
}

type VRFCoordinatorV25ConfigSetIterator struct {
	Event *VRFCoordinatorV25ConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ConfigSet)
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
		it.Event = new(VRFCoordinatorV25ConfigSet)
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

func (it *VRFCoordinatorV25ConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ConfigSet struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
	FallbackWeiPerUnitLink      *big.Int
	FeeConfig                   VRFCoordinatorV25FeeConfig
	Raw                         types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV25ConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ConfigSetIterator{contract: _VRFCoordinatorV25.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ConfigSet)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseConfigSet(log types.Log) (*VRFCoordinatorV25ConfigSet, error) {
	event := new(VRFCoordinatorV25ConfigSet)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25EthFundsRecoveredIterator struct {
	Event *VRFCoordinatorV25EthFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25EthFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25EthFundsRecovered)
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
		it.Event = new(VRFCoordinatorV25EthFundsRecovered)
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

func (it *VRFCoordinatorV25EthFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25EthFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25EthFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterEthFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV25EthFundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "EthFundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25EthFundsRecoveredIterator{contract: _VRFCoordinatorV25.contract, event: "EthFundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchEthFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25EthFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "EthFundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25EthFundsRecovered)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "EthFundsRecovered", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseEthFundsRecovered(log types.Log) (*VRFCoordinatorV25EthFundsRecovered, error) {
	event := new(VRFCoordinatorV25EthFundsRecovered)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "EthFundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25FundsRecoveredIterator struct {
	Event *VRFCoordinatorV25FundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25FundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25FundsRecovered)
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
		it.Event = new(VRFCoordinatorV25FundsRecovered)
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

func (it *VRFCoordinatorV25FundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25FundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25FundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV25FundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25FundsRecoveredIterator{contract: _VRFCoordinatorV25.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25FundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25FundsRecovered)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseFundsRecovered(log types.Log) (*VRFCoordinatorV25FundsRecovered, error) {
	event := new(VRFCoordinatorV25FundsRecovered)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OwnershipTransferRequestedIterator struct {
	Event *VRFCoordinatorV25OwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OwnershipTransferRequested)
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
		it.Event = new(VRFCoordinatorV25OwnershipTransferRequested)
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

func (it *VRFCoordinatorV25OwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV25OwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OwnershipTransferRequestedIterator{contract: _VRFCoordinatorV25.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OwnershipTransferRequested)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV25OwnershipTransferRequested, error) {
	event := new(VRFCoordinatorV25OwnershipTransferRequested)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25OwnershipTransferredIterator struct {
	Event *VRFCoordinatorV25OwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25OwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25OwnershipTransferred)
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
		it.Event = new(VRFCoordinatorV25OwnershipTransferred)
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

func (it *VRFCoordinatorV25OwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25OwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV25OwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25OwnershipTransferredIterator{contract: _VRFCoordinatorV25.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25OwnershipTransferred)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV25OwnershipTransferred, error) {
	event := new(VRFCoordinatorV25OwnershipTransferred)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ProvingKeyDeregisteredIterator struct {
	Event *VRFCoordinatorV25ProvingKeyDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ProvingKeyDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ProvingKeyDeregistered)
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
		it.Event = new(VRFCoordinatorV25ProvingKeyDeregistered)
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

func (it *VRFCoordinatorV25ProvingKeyDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ProvingKeyDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ProvingKeyDeregistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV25ProvingKeyDeregisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ProvingKeyDeregisteredIterator{contract: _VRFCoordinatorV25.contract, event: "ProvingKeyDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ProvingKeyDeregistered)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorV25ProvingKeyDeregistered, error) {
	event := new(VRFCoordinatorV25ProvingKeyDeregistered)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25ProvingKeyRegisteredIterator struct {
	Event *VRFCoordinatorV25ProvingKeyRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25ProvingKeyRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25ProvingKeyRegistered)
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
		it.Event = new(VRFCoordinatorV25ProvingKeyRegistered)
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

func (it *VRFCoordinatorV25ProvingKeyRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25ProvingKeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25ProvingKeyRegistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV25ProvingKeyRegisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25ProvingKeyRegisteredIterator{contract: _VRFCoordinatorV25.contract, event: "ProvingKeyRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ProvingKeyRegistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25ProvingKeyRegistered)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV25ProvingKeyRegistered, error) {
	event := new(VRFCoordinatorV25ProvingKeyRegistered)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25RandomWordsFulfilledIterator struct {
	Event *VRFCoordinatorV25RandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25RandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25RandomWordsFulfilled)
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
		it.Event = new(VRFCoordinatorV25RandomWordsFulfilled)
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

func (it *VRFCoordinatorV25RandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25RandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25RandomWordsFulfilled struct {
	RequestId     *big.Int
	OutputSeed    *big.Int
	Payment       *big.Int
	NativePayment bool
	Success       bool
	Raw           types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*VRFCoordinatorV25RandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25RandomWordsFulfilledIterator{contract: _VRFCoordinatorV25.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25RandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25RandomWordsFulfilled)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV25RandomWordsFulfilled, error) {
	event := new(VRFCoordinatorV25RandomWordsFulfilled)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25RandomWordsRequestedIterator struct {
	Event *VRFCoordinatorV25RandomWordsRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25RandomWordsRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25RandomWordsRequested)
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
		it.Event = new(VRFCoordinatorV25RandomWordsRequested)
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

func (it *VRFCoordinatorV25RandomWordsRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25RandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25RandomWordsRequested struct {
	KeyHash                     [32]byte
	RequestId                   *big.Int
	PreSeed                     *big.Int
	SubId                       uint64
	MinimumRequestConfirmations uint16
	CallbackGasLimit            uint32
	NumWords                    uint32
	NativePayment               bool
	Sender                      common.Address
	Raw                         types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*VRFCoordinatorV25RandomWordsRequestedIterator, error) {

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

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25RandomWordsRequestedIterator{contract: _VRFCoordinatorV25.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25RandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25RandomWordsRequested)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV25RandomWordsRequested, error) {
	event := new(VRFCoordinatorV25RandomWordsRequested)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25SubscriptionCanceledIterator struct {
	Event *VRFCoordinatorV25SubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25SubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25SubscriptionCanceled)
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
		it.Event = new(VRFCoordinatorV25SubscriptionCanceled)
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

func (it *VRFCoordinatorV25SubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25SubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25SubscriptionCanceled struct {
	SubId      uint64
	To         common.Address
	AmountLink *big.Int
	AmountEth  *big.Int
	Raw        types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25SubscriptionCanceledIterator{contract: _VRFCoordinatorV25.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionCanceled, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25SubscriptionCanceled)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV25SubscriptionCanceled, error) {
	event := new(VRFCoordinatorV25SubscriptionCanceled)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25SubscriptionConsumerAddedIterator struct {
	Event *VRFCoordinatorV25SubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25SubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25SubscriptionConsumerAdded)
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
		it.Event = new(VRFCoordinatorV25SubscriptionConsumerAdded)
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

func (it *VRFCoordinatorV25SubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25SubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25SubscriptionConsumerAdded struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25SubscriptionConsumerAddedIterator{contract: _VRFCoordinatorV25.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionConsumerAdded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25SubscriptionConsumerAdded)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV25SubscriptionConsumerAdded, error) {
	event := new(VRFCoordinatorV25SubscriptionConsumerAdded)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25SubscriptionConsumerRemovedIterator struct {
	Event *VRFCoordinatorV25SubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25SubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25SubscriptionConsumerRemoved)
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
		it.Event = new(VRFCoordinatorV25SubscriptionConsumerRemoved)
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

func (it *VRFCoordinatorV25SubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25SubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25SubscriptionConsumerRemoved struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25SubscriptionConsumerRemovedIterator{contract: _VRFCoordinatorV25.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25SubscriptionConsumerRemoved)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV25SubscriptionConsumerRemoved, error) {
	event := new(VRFCoordinatorV25SubscriptionConsumerRemoved)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25SubscriptionCreatedIterator struct {
	Event *VRFCoordinatorV25SubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25SubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25SubscriptionCreated)
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
		it.Event = new(VRFCoordinatorV25SubscriptionCreated)
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

func (it *VRFCoordinatorV25SubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25SubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25SubscriptionCreated struct {
	SubId uint64
	Owner common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25SubscriptionCreatedIterator{contract: _VRFCoordinatorV25.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionCreated, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25SubscriptionCreated)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV25SubscriptionCreated, error) {
	event := new(VRFCoordinatorV25SubscriptionCreated)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25SubscriptionFundedIterator struct {
	Event *VRFCoordinatorV25SubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25SubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25SubscriptionFunded)
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
		it.Event = new(VRFCoordinatorV25SubscriptionFunded)
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

func (it *VRFCoordinatorV25SubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25SubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25SubscriptionFunded struct {
	SubId      uint64
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25SubscriptionFundedIterator{contract: _VRFCoordinatorV25.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionFunded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25SubscriptionFunded)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV25SubscriptionFunded, error) {
	event := new(VRFCoordinatorV25SubscriptionFunded)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25SubscriptionFundedWithEthIterator struct {
	Event *VRFCoordinatorV25SubscriptionFundedWithEth

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25SubscriptionFundedWithEthIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25SubscriptionFundedWithEth)
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
		it.Event = new(VRFCoordinatorV25SubscriptionFundedWithEth)
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

func (it *VRFCoordinatorV25SubscriptionFundedWithEthIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25SubscriptionFundedWithEthIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25SubscriptionFundedWithEth struct {
	SubId         uint64
	OldEthBalance *big.Int
	NewEthBalance *big.Int
	Raw           types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterSubscriptionFundedWithEth(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionFundedWithEthIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "SubscriptionFundedWithEth", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25SubscriptionFundedWithEthIterator{contract: _VRFCoordinatorV25.contract, event: "SubscriptionFundedWithEth", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchSubscriptionFundedWithEth(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionFundedWithEth, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "SubscriptionFundedWithEth", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25SubscriptionFundedWithEth)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionFundedWithEth", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseSubscriptionFundedWithEth(log types.Log) (*VRFCoordinatorV25SubscriptionFundedWithEth, error) {
	event := new(VRFCoordinatorV25SubscriptionFundedWithEth)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionFundedWithEth", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25SubscriptionOwnerTransferRequestedIterator struct {
	Event *VRFCoordinatorV25SubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25SubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25SubscriptionOwnerTransferRequested)
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
		it.Event = new(VRFCoordinatorV25SubscriptionOwnerTransferRequested)
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

func (it *VRFCoordinatorV25SubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25SubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25SubscriptionOwnerTransferRequested struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25SubscriptionOwnerTransferRequestedIterator{contract: _VRFCoordinatorV25.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25SubscriptionOwnerTransferRequested)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV25SubscriptionOwnerTransferRequested, error) {
	event := new(VRFCoordinatorV25SubscriptionOwnerTransferRequested)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV25SubscriptionOwnerTransferredIterator struct {
	Event *VRFCoordinatorV25SubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV25SubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV25SubscriptionOwnerTransferred)
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
		it.Event = new(VRFCoordinatorV25SubscriptionOwnerTransferred)
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

func (it *VRFCoordinatorV25SubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV25SubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV25SubscriptionOwnerTransferred struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV25SubscriptionOwnerTransferredIterator{contract: _VRFCoordinatorV25.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV25.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV25SubscriptionOwnerTransferred)
				if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25Filterer) ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV25SubscriptionOwnerTransferred, error) {
	event := new(VRFCoordinatorV25SubscriptionOwnerTransferred)
	if err := _VRFCoordinatorV25.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type GetSubscription struct {
	Balance    *big.Int
	EthBalance *big.Int
	ReqCount   uint64
	Owner      common.Address
	Consumers  []common.Address
}
type SConfig struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
}
type SFeeConfig struct {
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

func (_VRFCoordinatorV25 *VRFCoordinatorV25) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinatorV25.abi.Events["ConfigSet"].ID:
		return _VRFCoordinatorV25.ParseConfigSet(log)
	case _VRFCoordinatorV25.abi.Events["EthFundsRecovered"].ID:
		return _VRFCoordinatorV25.ParseEthFundsRecovered(log)
	case _VRFCoordinatorV25.abi.Events["FundsRecovered"].ID:
		return _VRFCoordinatorV25.ParseFundsRecovered(log)
	case _VRFCoordinatorV25.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFCoordinatorV25.ParseOwnershipTransferRequested(log)
	case _VRFCoordinatorV25.abi.Events["OwnershipTransferred"].ID:
		return _VRFCoordinatorV25.ParseOwnershipTransferred(log)
	case _VRFCoordinatorV25.abi.Events["ProvingKeyDeregistered"].ID:
		return _VRFCoordinatorV25.ParseProvingKeyDeregistered(log)
	case _VRFCoordinatorV25.abi.Events["ProvingKeyRegistered"].ID:
		return _VRFCoordinatorV25.ParseProvingKeyRegistered(log)
	case _VRFCoordinatorV25.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFCoordinatorV25.ParseRandomWordsFulfilled(log)
	case _VRFCoordinatorV25.abi.Events["RandomWordsRequested"].ID:
		return _VRFCoordinatorV25.ParseRandomWordsRequested(log)
	case _VRFCoordinatorV25.abi.Events["SubscriptionCanceled"].ID:
		return _VRFCoordinatorV25.ParseSubscriptionCanceled(log)
	case _VRFCoordinatorV25.abi.Events["SubscriptionConsumerAdded"].ID:
		return _VRFCoordinatorV25.ParseSubscriptionConsumerAdded(log)
	case _VRFCoordinatorV25.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _VRFCoordinatorV25.ParseSubscriptionConsumerRemoved(log)
	case _VRFCoordinatorV25.abi.Events["SubscriptionCreated"].ID:
		return _VRFCoordinatorV25.ParseSubscriptionCreated(log)
	case _VRFCoordinatorV25.abi.Events["SubscriptionFunded"].ID:
		return _VRFCoordinatorV25.ParseSubscriptionFunded(log)
	case _VRFCoordinatorV25.abi.Events["SubscriptionFundedWithEth"].ID:
		return _VRFCoordinatorV25.ParseSubscriptionFundedWithEth(log)
	case _VRFCoordinatorV25.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _VRFCoordinatorV25.ParseSubscriptionOwnerTransferRequested(log)
	case _VRFCoordinatorV25.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _VRFCoordinatorV25.ParseSubscriptionOwnerTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFCoordinatorV25ConfigSet) Topic() common.Hash {
	return common.HexToHash("0xc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb2")
}

func (VRFCoordinatorV25EthFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x879c9ea2b9d5345b84ccd12610b032602808517cebdb795007f3dcb4df377317")
}

func (VRFCoordinatorV25FundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (VRFCoordinatorV25OwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFCoordinatorV25OwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFCoordinatorV25ProvingKeyDeregistered) Topic() common.Hash {
	return common.HexToHash("0x72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d")
}

func (VRFCoordinatorV25ProvingKeyRegistered) Topic() common.Hash {
	return common.HexToHash("0xe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b8")
}

func (VRFCoordinatorV25RandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x04624b1adcd18b732732ed3009cbdcc7dfc5f0340105149d0c193e8361d955e1")
}

func (VRFCoordinatorV25RandomWordsRequested) Topic() common.Hash {
	return common.HexToHash("0xad8d0b36e0dd1bc233432f9e3c28d376b9713fc9d26ac2dc972077be458d86fe")
}

func (VRFCoordinatorV25SubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0xb3a76e2c52bb7a484e325118e4cadc4346fdaa16a7df0621c5525e1f7ab4d495")
}

func (VRFCoordinatorV25SubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0")
}

func (VRFCoordinatorV25SubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b")
}

func (VRFCoordinatorV25SubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf")
}

func (VRFCoordinatorV25SubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8")
}

func (VRFCoordinatorV25SubscriptionFundedWithEth) Topic() common.Hash {
	return common.HexToHash("0x4e4421a74ab9b76003de9d9527153d002f1338c36c1bbc180069671631cf0958")
}

func (VRFCoordinatorV25SubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be")
}

func (VRFCoordinatorV25SubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0")
}

func (_VRFCoordinatorV25 *VRFCoordinatorV25) Address() common.Address {
	return _VRFCoordinatorV25.address
}

type VRFCoordinatorV25Interface interface {
	BLOCKHASHSTORE(opts *bind.CallOpts) (common.Address, error)

	LINK(opts *bind.CallOpts) (common.Address, error)

	LINKETHFEED(opts *bind.CallOpts) (common.Address, error)

	MAXCONSUMERS(opts *bind.CallOpts) (uint16, error)

	MAXNUMWORDS(opts *bind.CallOpts) (uint32, error)

	MAXREQUESTCONFIRMATIONS(opts *bind.CallOpts) (uint16, error)

	GetFeeTier(opts *bind.CallOpts, reqCount uint64) (uint32, error)

	GetRequestConfig(opts *bind.CallOpts) (uint16, uint32, [][32]byte, error)

	GetSubscription(opts *bind.CallOpts, subId uint64) (GetSubscription,

		error)

	HashOfKey(opts *bind.CallOpts, publicKey [2]*big.Int) ([32]byte, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	PendingRequestExists(opts *bind.CallOpts, subId uint64) (bool, error)

	SConfig(opts *bind.CallOpts) (SConfig,

		error)

	SCurrentSubId(opts *bind.CallOpts) (uint64, error)

	SFallbackWeiPerUnitLink(opts *bind.CallOpts) (*big.Int, error)

	SFeeConfig(opts *bind.CallOpts) (SFeeConfig,

		error)

	SProvingKeyHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	SProvingKeys(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error)

	SRequestCommitments(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	STotalBalance(opts *bind.CallOpts) (*big.Int, error)

	STotalEthBalance(opts *bind.CallOpts) (*big.Int, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AcceptSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	AddConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error)

	CancelSubscription(opts *bind.TransactOpts, subId uint64, to common.Address) (*types.Transaction, error)

	CreateSubscription(opts *bind.TransactOpts) (*types.Transaction, error)

	DeregisterProvingKey(opts *bind.TransactOpts, publicProvingKey [2]*big.Int) (*types.Transaction, error)

	FulfillRandomWords(opts *bind.TransactOpts, proof VRFProof, rc VRFCoordinatorV25RequestCommitment) (*types.Transaction, error)

	FundSubscriptionWithEth(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	OnTokenTransfer(opts *bind.TransactOpts, arg0 common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	OracleWithdraw(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OracleWithdrawEth(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	OwnerCancelSubscription(opts *bind.TransactOpts, subId uint64) (*types.Transaction, error)

	RecoverEthFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RecoverFunds(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	RegisterProvingKey(opts *bind.TransactOpts, oracle common.Address, publicProvingKey [2]*big.Int) (*types.Transaction, error)

	RemoveConsumer(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error)

	RequestRandomWords(opts *bind.TransactOpts, keyHash [32]byte, subId uint64, requestConfirmations uint16, callbackGasLimit uint32, numWords uint32, payInEth bool) (*types.Transaction, error)

	RequestSubscriptionOwnerTransfer(opts *bind.TransactOpts, subId uint64, newOwner common.Address) (*types.Transaction, error)

	SetConfig(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV25FeeConfig) (*types.Transaction, error)

	SetLINK(opts *bind.TransactOpts, link common.Address) (*types.Transaction, error)

	SetLinkEthFeed(opts *bind.TransactOpts, linkEthFeed common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV25ConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFCoordinatorV25ConfigSet, error)

	FilterEthFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV25EthFundsRecoveredIterator, error)

	WatchEthFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25EthFundsRecovered) (event.Subscription, error)

	ParseEthFundsRecovered(log types.Log) (*VRFCoordinatorV25EthFundsRecovered, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV25FundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25FundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*VRFCoordinatorV25FundsRecovered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV25OwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV25OwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV25OwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25OwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV25OwnershipTransferred, error)

	FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV25ProvingKeyDeregisteredIterator, error)

	WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorV25ProvingKeyDeregistered, error)

	FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV25ProvingKeyRegisteredIterator, error)

	WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25ProvingKeyRegistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV25ProvingKeyRegistered, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*VRFCoordinatorV25RandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25RandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV25RandomWordsFulfilled, error)

	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*VRFCoordinatorV25RandomWordsRequestedIterator, error)

	WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25RandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV25RandomWordsRequested, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionCanceled, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV25SubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionConsumerAdded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV25SubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV25SubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionCreated, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV25SubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionFunded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV25SubscriptionFunded, error)

	FilterSubscriptionFundedWithEth(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionFundedWithEthIterator, error)

	WatchSubscriptionFundedWithEth(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionFundedWithEth, subId []uint64) (event.Subscription, error)

	ParseSubscriptionFundedWithEth(log types.Log) (*VRFCoordinatorV25SubscriptionFundedWithEth, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV25SubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV25SubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV25SubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV25SubscriptionOwnerTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
